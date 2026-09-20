// Package audio implements OpenAI audio operations.
package audio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/codec"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls audio endpoints.
type Client struct{ transport *transport.Client }

// New creates an audio client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// TranscriptionInput is a multipart transcription request.
type TranscriptionInput struct {
	Model    string
	FileName string
	File     io.Reader
	Language string
}

// Transcribe transcribes an audio file.
func (client *Client) Transcribe(ctx context.Context, input TranscriptionInput) (Transcription, error) {
	var output Transcription
	body, contentType, err := multipartBody(input.Model, input.FileName, input.File, input.Language)
	if err != nil {
		return output, err
	}
	data, err := request.Reader(ctx, client.transport, http.MethodPost, "/audio/transcriptions", "audio.transcriptions", body, contentType)
	if err != nil {
		return output, err
	}
	return request.Decode[Transcription](data)
}

// Translate translates an audio file to English.
func (client *Client) Translate(ctx context.Context, input TranscriptionInput) (Translation, error) {
	var output Translation
	body, contentType, err := multipartBody(input.Model, input.FileName, input.File, "")
	if err != nil {
		return output, err
	}
	data, err := request.Reader(ctx, client.transport, http.MethodPost, "/audio/translations", "audio.translations", body, contentType)
	if err != nil {
		return output, err
	}
	return request.Decode[Translation](data)
}

// Speech generates audio bytes. The caller owns the returned byte slice.
func (client *Client) Speech(ctx context.Context, input SpeechRequest) ([]byte, error) {
	body, err := codec.Marshal(input)
	if err != nil {
		return nil, err
	}
	return request.Bytes(ctx, client.transport, http.MethodPost, "/audio/speech", "audio.speech", body, "application/json")
}

func multipartBody(model, fileName string, file io.Reader, language string) (io.ReadCloser, string, error) {
	if model == "" || fileName == "" || file == nil {
		return nil, "", errors.New("model, file name, and file are required")
	}
	reader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)
	go func() {
		err := multipartWriter.WriteField("model", model)
		if err == nil && language != "" {
			err = multipartWriter.WriteField("language", language)
		}
		if err == nil {
			var part io.Writer
			part, err = multipartWriter.CreateFormFile("file", fileName)
			if err == nil {
				_, err = io.Copy(part, file)
			}
		}
		if err == nil {
			err = multipartWriter.Close()
		}
		if err != nil {
			_ = pipeWriter.CloseWithError(fmt.Errorf("write audio multipart: %w", err))
			return
		}
		_ = pipeWriter.Close()
	}()
	return reader, multipartWriter.FormDataContentType(), nil
}
