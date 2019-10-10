package adapters

import "github.com/aws/aws-sdk-go/service/polly"

// Repository - repo layer for request
type AwsAdapter interface {
	Polly_SynthesizeSpeech(input *polly.SynthesizeSpeechInput) (*polly.SynthesizeSpeechOutput,error)
}
