package repository

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"ttsService/awsPolly"
	awsttsModel "ttsService/awsPolly/models"
	"ttsService/logger"
)

// Controller - struct to logically bind all the controller functions
type Controller struct {
	AWSPollyRepo awsPolly.AwsPollyRepository
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}


func NewRequestController(e *echo.Echo, awsPollyRepo awsPolly.AwsPollyRepository) {
	e.Validator = &CustomValidator{Validator: validator.New()}
	requestHandler := &Controller{
		AWSPollyRepo: awsPollyRepo,
	}
	e.POST("v1/tts/aws/", requestHandler.ttsRequest)
}

func (a *Controller) ttsRequest(c echo.Context) error {
	response := make(map[string]interface{})
	uuidAudio, err := uuid.NewUUID()
	logger.Log().WithFields(logrus.Fields{
		"uuid": uuidAudio,
	}).Info("Request Received")
	if err != nil {
		handleError(err,uuidAudio)
		response["error"] = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}
	response["request_id"]=uuidAudio

	ttsDetails := new(awsttsModel.AwsTTSRequest)

	if err := c.Bind(ttsDetails); err != nil {
		handleError(err,uuidAudio)
		response["error"] = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := c.Validate(ttsDetails); err != nil {
		handleError(err,uuidAudio)
		response["error"] = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	resp, err := a.AWSPollyRepo.AWSSynthesizeSpeech(ttsDetails,uuidAudio)
	if err != nil {
		handleError(err,uuidAudio)
		response["error"] = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}
	response["tts_file"] = resp
	logger.Log().WithFields(logrus.Fields{
		"uuid": uuidAudio,
		"response":response,
	}).Info("success")
	return c.JSON(http.StatusOK, response)
}

func handleError(err error, uuidAudio uuid.UUID)  {
	logger.Log().WithFields(logrus.Fields{
		"uuid": uuidAudio,
	}).Error(err)
}