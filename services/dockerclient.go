package services

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/francescofrontera/ks-job-uploader/utils"
	"io"
	"log"
	"os"
)

type DockerClientResult struct {
	dockerClient *client.Client
	ctx context.Context
	log *log.Logger
}

const BasicImageName = "ks_base_image"

func InitClient(clientVersion string, logger *log.Logger)  *DockerClientResult {
	cli, error := client.NewClientWithOpts(client.WithVersion(clientVersion)); if error != nil {
		logger.Fatalf("An error occured during init of docker client: %v", error)
	}

	ctx := context.Background()

	return &DockerClientResult{
		dockerClient: cli,
		ctx: ctx,
		log: logger,
	}
}

func getDockerFileCtx() (*os.File, error) {
	ctx, error := os.Open("/go/src/github.com/francescofrontera/ks-job-uploader/docker/docker_as_t.tar.gz")
	return ctx, error
}

func (dcb *DockerClientResult) BuildImage() error {
	dockerBuildContext, errF := getDockerFileCtx(); if errF != nil {
		return errF
	}

	defer dockerBuildContext.Close()

	cli := dcb.dockerClient
	ctx := dcb.ctx

	buildOptions := types.ImageBuildOptions{
		Tags: []string{BasicImageName},
		Dockerfile: "docker/Dockerfile",
		Context: dockerBuildContext,
	}

	response, err := cli.ImageBuild(ctx, dockerBuildContext, buildOptions); if err != nil {
		return err
	}

	io.Copy(os.Stdout, response.Body)

	defer response.Body.Close()
	return nil
}



func (dcb *DockerClientResult) RunContainer(jarToMount, mainClass string) (string, error) {
	cli := dcb.dockerClient
	ctx := dcb.ctx

	containerConfig := &container.Config{
		Image: BasicImageName,
		Tty:   true,
		Env: []string{
			fmt.Sprintf("JAR_TO_EXECUTE=%s", jarToMount),
			fmt.Sprintf("MAIN_CLASS=%s", mainClass),
		},
	}

	sourcePath, targetPath := utils.GetPathToJar(jarToMount)
	log.Printf("SourcePath: %s, TargetPath: %s", sourcePath, targetPath)

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type: mount.TypeBind,
				Source: sourcePath,
				Target: targetPath,
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx,
		containerConfig,
		hostConfig,
		nil,
		""); if err != nil {
		return "", err
	}

	containerId := resp.ID

	if err := cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return containerId, nil
}

