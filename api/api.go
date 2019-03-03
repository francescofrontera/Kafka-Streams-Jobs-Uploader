package api

import (
	"fmt"
	"github.com/francescofrontera/ks-job-uploader/services"
	"net/http"
	"os"
	"strings"

	"github.com/francescofrontera/ks-job-uploader/models"
	"github.com/gin-gonic/gin"
)

type JarHandler struct {
	dockerClient *services.DockerClientResult
	redisClient *services.RedisConf
}

func NewJarHandlerFunction(dClient *services.DockerClientResult, rClient *services.RedisConf) *JarHandler {
	return &JarHandler{
		dockerClient: dClient,
		redisClient: rClient,
	}
}

func (jarHandler *JarHandler) UploadHandler(gctx *gin.Context) {
	file, err := gctx.FormFile("uploadFile"); if err != nil {
		panic(err) //dont do this
	}

	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dst := strings.Join([]string{workdir, "jars", file.Filename}, "/")
	if err := gctx.SaveUploadedFile(file, dst); err != nil {
		gctx.JSON(http.StatusBadRequest, gin.H{
			"error":  fmt.Sprintf("upload file error %s", err.Error()),
			"status": 400,
		})
		return
	}

	gctx.JSON(http.StatusAccepted, gin.H{"fileName": file.Filename, "status": http.StatusAccepted})
}

func (jarHandler *JarHandler) RunKSJob(gctx *gin.Context) {
	dockerClient := jarHandler.dockerClient

	runJar := &models.RunJar{}
	if error := gctx.BindJSON(runJar); error == nil {
		containerID, dockerClientError := dockerClient.RunContainer(runJar.JarName, runJar.MainClass)
		if dockerClientError != nil {
			responseStatus := http.StatusInternalServerError
			gctx.JSON(responseStatus, gin.H{"error": dockerClientError.Error(), "status": responseStatus})
			return
		}

		responseStatus := http.StatusAccepted
		gctx.JSON(responseStatus, gin.H{"containerId": containerID, "status": responseStatus})
	}
}
