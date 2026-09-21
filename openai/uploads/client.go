/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package uploads implements OpenAI large-upload lifecycle operations.
package uploads

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls upload endpoints.
type Client struct{ transport *transport.Client }

// New creates an upload client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create starts an upload.
func (client *Client) Create(ctx context.Context, input CreateRequest) (Upload, error) {
	var output Upload
	err := request.JSON(ctx, client.transport, request.Post, "/uploads", "uploads.create", input, &output)
	return output, err
}

// Get returns upload status.
func (client *Client) Get(ctx context.Context, id string) (Upload, error) {
	var output Upload
	err := request.JSON(ctx, client.transport, request.Get, "/uploads/"+url.PathEscape(id), "uploads.get", nil, &output)
	return output, err
}

// AddPart uploads one part without buffering the source reader in memory.
func (client *Client) AddPart(ctx context.Context, id string, input PartInput) (Part, error) {
	var output Part
	body, contentType, err := multipartReader(input)
	if err != nil {
		return output, err
	}
	data, requestErr := request.Reader(ctx, client.transport, http.MethodPost, "/uploads/"+url.PathEscape(id)+"/parts", "uploads.add_part", body, contentType)
	bodyErr := body.Close()
	if requestErr != nil {
		return output, requestErr
	}
	if bodyErr != nil {
		return output, bodyErr
	}
	return request.Decode[Part](data)
}

// Complete finishes an upload.
func (client *Client) Complete(ctx context.Context, id string, input CompleteRequest) (Upload, error) {
	var output Upload
	err := request.JSON(ctx, client.transport, request.Post, "/uploads/"+url.PathEscape(id)+"/complete", "uploads.complete", input, &output)
	return output, err
}

// Cancel cancels an upload.
func (client *Client) Cancel(ctx context.Context, id string) (DeleteResponse, error) {
	var output DeleteResponse
	err := request.JSON(ctx, client.transport, request.Post, "/uploads/"+url.PathEscape(id)+"/cancel", "uploads.cancel", nil, &output)
	return output, err
}

func multipartReader(input PartInput) (*transport.MultipartBody, string, error) {
	if input.Filename == "" || input.File == nil {
		return nil, "", errors.New("filename and file are required")
	}
	body, contentType := transport.NewMultipartBody(func(multipartWriter *multipart.Writer) error {
		part, err := multipartWriter.CreateFormFile("data", input.Filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(part, input.File)
		return err
	})
	return body, contentType, nil
}
