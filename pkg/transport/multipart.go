/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package transport

import (
	"errors"
	"io"
	"mime/multipart"
	"sync"
)

// MultipartBody streams a multipart form through an owned pipe. Closing the
// body closes the reader side and waits for the writer goroutine to finish.
type MultipartBody struct {
	reader    *io.PipeReader
	done      chan struct{}
	closeOnce sync.Once
	writeErr  error
}

// NewMultipartBody creates a bounded-lifetime multipart request body.
// The write function owns the multipart fields and file readers; it must not
// call multipart.Writer.Close because the helper closes the writer itself.
func NewMultipartBody(write func(*multipart.Writer) error) (*MultipartBody, string) {
	reader, writer := io.Pipe()
	multipartWriter := multipart.NewWriter(writer)
	body := &MultipartBody{reader: reader, done: make(chan struct{})}
	go func() {
		var err error
		if write == nil {
			err = errors.New("multipart writer function is nil")
		} else {
			err = write(multipartWriter)
		}
		if err == nil {
			err = multipartWriter.Close()
		}
		body.writeErr = err
		if err != nil {
			_ = writer.CloseWithError(err)
		} else {
			_ = writer.Close()
		}
		close(body.done)
	}()
	return body, multipartWriter.FormDataContentType()
}

// Read reads the encoded multipart body.
func (body *MultipartBody) Read(buffer []byte) (int, error) {
	return body.reader.Read(buffer)
}

// Close closes the reader side and waits for multipart encoding to finish.
func (body *MultipartBody) Close() error {
	body.closeOnce.Do(func() { _ = body.reader.Close() })
	<-body.done
	return body.writeErr
}
