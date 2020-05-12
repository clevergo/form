// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

// Package form is a form decoder that decode request body into a struct.
package form

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"net/http"

	"github.com/gorilla/schema"
)

// Validatable indicates whether a value can be validated.
type Validatable interface {
	Validate() error
}

// Content type constants.
const (
	ContentType              = "Content-Type"
	ContentTypeForm          = "application/x-www-form-urlencoded"
	ContentTypeMultipartForm = "multipart/form-data"
	ContentTypeJSON          = "application/json"
	ContentTypeXML           = "application/xml"
)

var (
	defaultDecoders  *Decoders
	defaultDecoder         = schema.NewDecoder()
	defaultMaxMemory int64 = 10 * 1024 * 1024
)

func init() {
	defaultDecoder.IgnoreUnknownKeys(true)

	defaultDecoders = New()
}

// New returns a decoders.
func New() *Decoders {
	d := &Decoders{}
	d.Register(ContentTypeForm, NewForm(defaultDecoder))
	d.Register(ContentTypeMultipartForm, NewMultipartForm(defaultMaxMemory))
	d.Register(ContentTypeJSON, JSON)
	d.Register(ContentTypeXML, XML)
	return d
}

// Register a decoder for the given content type.
func Register(contentType string, decoder Decoder) {
	defaultDecoders.Register(contentType, decoder)
}

// Decode data from a request into v, v should be a pointer.
func Decode(r *http.Request, v interface{}) error {
	return defaultDecoders.Decode(r, v)
}

// Decoders is a map that mapping from content type to decoder.
type Decoders map[string]Decoder

// Register a decoder for the given content type.
func (d *Decoders) Register(contentType string, decoder Decoder) {
	(*d)[contentType] = decoder
}

// Decode data from a request into v, v should be a pointer.
func (d *Decoders) Decode(r *http.Request, v interface{}) error {
	contentType, err := parseContentType(r)
	if err != nil {
		return err
	}
	decoder, ok := (*d)[contentType]
	if !ok {
		return errors.New("Unsupported content type: " + contentType)
	}
	if err = decoder(r, v); err != nil {
		return err
	}
	if vv, ok := v.(Validatable); ok {
		err = vv.Validate()
	}

	return err
}

// Decoder is a function that decode data from request into v.
type Decoder func(req *http.Request, v interface{}) error

func parseContentType(r *http.Request) (string, error) {
	header := r.Header.Get(ContentType)
	contentType, _, err := mime.ParseMediaType(header)
	if err != nil {
		return "", err
	}

	return contentType, nil
}

// JSON is a JSON decoder.
func JSON(r *http.Request, v interface{}) error {
	var buf bytes.Buffer
	reader := io.TeeReader(r.Body, &buf)
	r.Body = ioutil.NopCloser(&buf)
	return json.NewDecoder(reader).Decode(v)
}

// XML is an XML decoder.
func XML(r *http.Request, v interface{}) error {
	var buf bytes.Buffer
	reader := io.TeeReader(r.Body, &buf)
	r.Body = ioutil.NopCloser(&buf)
	return xml.NewDecoder(reader).Decode(v)
}

// NewForm returns a post form decoder with the given schema decoder.
func NewForm(decoder *schema.Decoder) Decoder {
	return func(r *http.Request, v interface{}) error {
		err := r.ParseForm()
		if err != nil {
			return err
		}

		return decoder.Decode(v, r.PostForm)
	}
}

// NewMultipartForm returns a multipart form decoder with the given schema decoder.
func NewMultipartForm(maxMemory int64) Decoder {
	return func(r *http.Request, v interface{}) error {
		err := r.ParseMultipartForm(maxMemory)
		if err != nil {
			return err
		}

		return defaultDecoder.Decode(v, r.MultipartForm.Value)
	}
}
