package repository

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/aws/session"
	"ttsService/configs"
	"ttsService/adapters"
)

type awsPollySessionAdapter struct {
	awsPolly *polly.Polly
}

func newAwsPollySession(config *configs.Config) (*polly.Polly,error){

	token := ""

	creds := credentials.NewStaticCredentials(config.AwsConfig.AccessKeyID, config.AwsConfig.SecretAcessKey, token)
	_, err := creds.Get()
	if err != nil {
		return nil,err
		// handle error
	}
	cfg := aws.NewConfig().WithRegion(config.AwsConfig.Region).WithCredentials(creds)

	svc := polly.New(session.New(), cfg)
	return svc,nil
}

func NewAwsPollySessionAdapter(config *configs.Config) (adapters.AwsAdapter, error) {
	pollySession,err := newAwsPollySession(config)
	return &awsPollySessionAdapter{
		awsPolly: pollySession,
	}, err
}

// Synthesizes plain text or SSML into a file of human-like speech.
func (awsPollySession *awsPollySessionAdapter) Polly_SynthesizeSpeech(input *polly.SynthesizeSpeechInput) (*polly.SynthesizeSpeechOutput,error){
	result, err := awsPollySession.awsPolly.SynthesizeSpeech(input)
	if err != nil {
		return nil,err
	}
	return result,err
}