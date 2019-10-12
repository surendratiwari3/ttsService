package repository

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	_ "io/ioutil"
	"os"
	"ttsService/adapters"
	"ttsService/awsPolly"
	"ttsService/awsPolly/models"
	"ttsService/logger"
)

type awsRespository struct{
	awsPollySessionAdapter adapters.AwsAdapter
}



// NewRequestRepository - Requests repository object
func NewPollyRequestRepository(awsPollyadapter adapters.AwsAdapter) awsPolly.AwsPollyRepository{
	return &awsRespository{
		awsPollySessionAdapter: awsPollyadapter,
	}
}


func createTTS_Request(outputFormat string, sampleRate string, text string, voiceId string) *polly.SynthesizeSpeechInput {
	input := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String(outputFormat),
		SampleRate:   aws.String(sampleRate),
		Text:         aws.String(text),
		TextType:     aws.String("text"),
		VoiceId:      aws.String(voiceId),
	}
	return input
}

// Synthesizes plain text or SSML into a file of human-like speech.
func (r *awsRespository) AWSSynthesizeSpeech(request *models.AwsTTSRequest,uuidAudio uuid.UUID) (string,error){
	logger.Log().WithFields(logrus.Fields{
		"uuid": uuidAudio,
	}).Info("requesting tts with outputformat: " + request.OutputFormat)

	//audioFileName := uuidAudio.String() + "." + request.OutputFormat
	input := createTTS_Request(request.OutputFormat,request.SampleRate,request.InputText, request.VoiceId)
	response,err := r.awsPollySessionAdapter.Polly_SynthesizeSpeech(input)
	//Error Handler
	if err != nil {
		ttserrorHandle(err)
		logger.Log().WithFields(logrus.Fields{
			"uuid": uuidAudio,
		}).Error(err)
		return "",err
	}

	/*
	//upload to minio
	presignedURL, err := r.minioSessionAdapter.PutStreamObject("test2", audioFileName, response.AudioStream)
	if err != nil {
		logger.Log().WithFields(logrus.Fields{
			"uuid": uuidAudio,
		}).Error(err)
		return nil,err
	}

	logger.Log().WithFields(logrus.Fields{
		"uuid": uuidAudio,
		"url": presignedURL,
	}).Info("tts generated successfully")

*/
	//writing locally files
	audioFile, err := readStreamToFile(response.AudioStream,request.OutputFormat)
	if err!=nil{
		return "",err
	}
	return audioFile,err
}


func readStreamToFile(closer io.ReadCloser,outputFormat string) (string,error) {
	uuidAudio, err := uuid.NewUUID()
	if err != nil {
		return "",err
	}
	audioFileName := "/opt/tts_files/"+ uuidAudio.String() + "." + outputFormat
	// open output file
	fo, err := os.Create(audioFileName)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := closer.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
	return audioFileName,nil
}


func ttserrorHandle(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case polly.ErrCodeTextLengthExceededException:
			fmt.Println(polly.ErrCodeTextLengthExceededException, aerr.Error())
		case polly.ErrCodeInvalidSampleRateException:
			fmt.Println(polly.ErrCodeInvalidSampleRateException, aerr.Error())
		case polly.ErrCodeInvalidSsmlException:
			fmt.Println(polly.ErrCodeInvalidSsmlException, aerr.Error())
		case polly.ErrCodeLexiconNotFoundException:
			fmt.Println(polly.ErrCodeLexiconNotFoundException, aerr.Error())
		case polly.ErrCodeServiceFailureException:
			fmt.Println(polly.ErrCodeServiceFailureException, aerr.Error())
		case polly.ErrCodeMarksNotSupportedForFormatException:
			fmt.Println(polly.ErrCodeMarksNotSupportedForFormatException, aerr.Error())
		case polly.ErrCodeSsmlMarksNotSupportedForTextTypeException:
			fmt.Println(polly.ErrCodeSsmlMarksNotSupportedForTextTypeException, aerr.Error())
		case polly.ErrCodeLanguageNotSupportedException:
			fmt.Println(polly.ErrCodeLanguageNotSupportedException, aerr.Error())
		default:
			fmt.Println(aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
	}
	return
}
