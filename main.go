package main

import (
	"github.com/francescofrontera/ks-job-uploader/api"
	"github.com/francescofrontera/ks-job-uploader/services"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
)

var (
	serverAddress = os.Getenv("SERVER_ADDRESS")
	redisAddress =  os.Getenv("REDIS_ADDRESS")
	redisPassword = os.Getenv("REDIS_PWD")
	redisDB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
)

func main() {
	appLogger := log.New(os.Stdout, "KS-JOB-UPLOADER ", log.LstdFlags | log.Lshortfile )

	/* Services */
	db := services.RedisClient(redisAddress, redisPassword, redisDB, appLogger)
	dockerClient := services.InitClient("1.38", appLogger)


	/* HANDLERS */
	jarHandler := api.NewJarHandlerFunction(dockerClient, db, appLogger)

	route := gin.Default()

	v1 := route.Group("/v1/api")
	{
		v1.POST("/upload", jarHandler.UploadHandler)
		v1.POST("/run", jarHandler.RunKSJob)
	}

	buildImgError := dockerClient.BuildImage();  if buildImgError != nil {
		appLogger.Fatalf("Error during build base image: %v", buildImgError)
	}

	defer route.Run(serverAddress)
}
