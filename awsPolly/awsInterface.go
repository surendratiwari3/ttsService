package awsPolly

import (
	"github.com/google/uuid"
	"ttsService/awsPolly/models"
)
// Repository - repo layer for request
type AwsPollyRepository interface {
	AWSSynthesizeSpeech(request *models.AwsTTSRequest,uuidAudio uuid.UUID) (string,error)
}
