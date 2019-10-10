package models


// AwsTTSRequest
type (
	AwsTTSRequest struct {
		OutputFormat string `json:"outputFormat" form:"outputFormat" query:"outputFormat" validate:"oneof=mp3 ogg pcm"`
		SampleRate   string `json:"sampleRate" form:"sampleRate" query:"sampleRate" validate:"oneof=8000 16000 22050"`
		InputText    string `json:"inputText" form:"inputText" query:"inputText" validate:"required"`
		VoiceId      string `json:"voiceId" form:"voiceId" query:"voiceId"`
	}
)
