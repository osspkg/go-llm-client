// Package audio implements OpenAI audio operations.
package audio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

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

// CreateVoiceConsent uploads a recording that authorizes creation of a custom voice.
func (client *Client) CreateVoiceConsent(ctx context.Context, input VoiceConsentInput) (VoiceConsent, error) {
	var output VoiceConsent
	body, contentType, err := voiceConsentBody(input)
	if err != nil {
		return output, err
	}
	data, err := request.Reader(ctx, client.transport, request.Post, "/audio/voice_consents", "audio.voice_consents.create", body, contentType)
	if err != nil {
		return output, err
	}
	return request.Decode[VoiceConsent](data)
}

// ListVoiceConsents returns the uploaded voice-consent recordings.
func (client *Client) ListVoiceConsents(ctx context.Context) (VoiceConsentListResponse, error) {
	var output VoiceConsentListResponse
	err := request.JSON(ctx, client.transport, request.Get, "/audio/voice_consents", "audio.voice_consents.list", nil, &output)
	return output, err
}

// GetVoiceConsent returns a voice-consent recording.
func (client *Client) GetVoiceConsent(ctx context.Context, id string) (VoiceConsent, error) {
	var output VoiceConsent
	err := request.JSON(ctx, client.transport, request.Get, "/audio/voice_consents/"+url.PathEscape(id), "audio.voice_consents.get", nil, &output)
	return output, err
}

// UpdateVoiceConsent changes the label of a voice-consent recording.
func (client *Client) UpdateVoiceConsent(ctx context.Context, id string, input UpdateVoiceConsentRequest) (VoiceConsent, error) {
	var output VoiceConsent
	err := request.JSON(ctx, client.transport, request.Post, "/audio/voice_consents/"+url.PathEscape(id), "audio.voice_consents.update", input, &output)
	return output, err
}

// DeleteVoiceConsent removes a voice-consent recording.
func (client *Client) DeleteVoiceConsent(ctx context.Context, id string) (DeleteVoiceConsentResponse, error) {
	var output DeleteVoiceConsentResponse
	err := request.JSON(ctx, client.transport, request.Delete, "/audio/voice_consents/"+url.PathEscape(id), "audio.voice_consents.delete", nil, &output)
	return output, err
}

// CreateVoice creates a custom voice from a consented audio sample.
func (client *Client) CreateVoice(ctx context.Context, input VoiceInput) (Voice, error) {
	var output Voice
	body, contentType, err := voiceBody(input)
	if err != nil {
		return output, err
	}
	data, err := request.Reader(ctx, client.transport, request.Post, "/audio/voices", "audio.voices.create", body, contentType)
	if err != nil {
		return output, err
	}
	return request.Decode[Voice](data)
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

func voiceConsentBody(input VoiceConsentInput) (io.ReadCloser, string, error) {
	if input.Name == "" || input.Language == "" || input.Filename == "" || input.Recording == nil {
		return nil, "", errors.New("name, language, filename, and recording are required")
	}
	return multipartFileBody("recording", input.Filename, input.Recording, func(writer *multipart.Writer) error {
		if err := writer.WriteField("name", input.Name); err != nil {
			return err
		}
		return writer.WriteField("language", input.Language)
	})
}

func voiceBody(input VoiceInput) (io.ReadCloser, string, error) {
	if input.Name == "" || input.ConsentID == "" || input.Filename == "" || input.AudioSample == nil {
		return nil, "", errors.New("name, consent id, filename, and audio sample are required")
	}
	return multipartFileBody("audio_sample", input.Filename, input.AudioSample, func(writer *multipart.Writer) error {
		if err := writer.WriteField("name", input.Name); err != nil {
			return err
		}
		return writer.WriteField("consent", input.ConsentID)
	})
}

func multipartFileBody(field, filename string, file io.Reader, fields func(*multipart.Writer) error) (io.ReadCloser, string, error) {
	reader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)
	go func() {
		err := fields(multipartWriter)
		if err == nil {
			var part io.Writer
			part, err = multipartWriter.CreateFormFile(field, filename)
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
