package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"sync"
	driver "ttsService/adapters/repository"
	awstts "ttsService/awsPolly/repository"
	"ttsService/configs"
	"ttsService/logger"
	ttshttpHandler "ttsService/ttsRequestHandler/repository"
)

var onceRest sync.Once

func main() {
	onceRest.Do(func() {
		e := echo.New()
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		//Setting up the config
		config := configs.GetConfig()

		//getting the awsPollyAdapter
		awspollyAdapter, err := driver.NewAwsPollySessionAdapter(config)
		if err != nil {
			logger.Log().WithFields(logrus.Fields{
				"message": "AWS polly session connection failed",
			}).Error(err)
			panic(err)
			//handle error
		}
		/*
		minioAdapter, err := driver.NewMinioSessionAdapter(config)
		if err != nil {
			logger.Log().WithFields(logrus.Fields{
				"message": "Minio session connection failed",
			}).Error(err)
			panic(err)
		}*/
		awsttsRepo := awstts.NewPollyRequestRepository(awspollyAdapter)
		if err != nil {
			logger.Log().WithFields(logrus.Fields{
				"message": "ttsRepo Adapter failed",
			}).Error(err)
			panic(err)
		}
		ttshttpHandler.NewRequestController(e, awsttsRepo)
		if err := e.Start(config.EchoConfig.HostPort); err != nil {
			logger.Log().WithFields(logrus.Fields{
				"message": "Echo Server Failed to Start",
			}).Error(err)
			panic(err)
			//handle error
		}
	})
}
