package api

import (
	"fmt"
	"github.com/francescofrontera/ks-job-uploader/models"
	"github.com/francescofrontera/ks-job-uploader/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

//FIXME: Move this const in conf file
const WorkDir = "/go/src/github.com/francescofrontera/ks-job-uploader"

type JarHandler struct {
	dockerClient *services.DockerClientResult
	redisClient *services.RedisConf
	log *log.Logger
}

func NewJarHandlerFunction(dClient *services.DockerClientResult, rClient *services.RedisConf, logger *log.Logger) *JarHandler {
	return &JarHandler{
		dockerClient: dClient,
		redisClient: rClient,
		log: logger,
	}
}

func (jarHandler *JarHandler) UploadHandler(gctx *gin.Context) {
	jarsPath := strings.Join([]string{WorkDir, "jars"}, "/")

	if _, error := os.Stat(jarsPath); os.IsNotExist(error) {
		os.Mkdir(jarsPath, os.ModePerm)
	}

	file, err := gctx.FormFile("uploadFile"); if err != nil {
		jarHandler.log.Fatalf("An error occured when upload JAR file %v", err) //dont do this
	}

	dst := strings.Join([]string{WorkDir, "jars", file.Filename}, "/")
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
