package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Simplou/openai/v1"
)

const fileName = "hello.mp3"

var (
	audioFilePath = fmt.Sprintf("./temp/%s", fileName)
)

func tts() {
	audio, err := openai.TextToSpeech(client, httpClient, &openai.SpeechRequestBody{
		Model: "tts-1",
		Input: "Hello",
		Voice: openai.SpeechVoices.Onyx,
	})
	if err != nil {
		panic(err)
	}
	defer audio.Close()
	b, err := io.ReadAll(audio)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(audioFilePath, b, os.ModePerm); err != nil {
		panic(err)
	}
}

func whisper() {
	transcription, err := openai.Transcription(client, httpClient, &openai.TranscriptionsRequestBody{
		Model:         openai.DefaultTranscriptionModel,
		Filename:      fileName,
		AudioFilePath: audioFilePath,
	})
	if err != nil {
		panic(err)
	}
	println(transcription.Text)
}
