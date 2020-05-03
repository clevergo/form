// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package form

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type login struct {
	validated bool
	Username  string `schema:"username" json:"username" xml:"username"`
	Password  string `schema:"password" json:"password" xml:"password"`
}

func (l *login) Validate() error {
	l.validated = true
	return nil
}

func TestRegister(t *testing.T) {
	invoked := false
	decoder := func(r *http.Request, v interface{}) error {
		invoked = true
		return nil
	}
	contentType := "content/type"
	Register(contentType, decoder)
	actual, ok := (*defaultDecoders)[contentType]
	assert.True(t, ok)
	err := actual(httptest.NewRequest(http.MethodGet, "/", nil), nil)
	assert.Nil(t, err)
	assert.True(t, invoked)
}

func TestJSON(t *testing.T) {
	tests := []struct {
		shouldErr bool
		body      []byte
	}{
		{true, []byte(`{"invalid json"}`)},
		{false, []byte(`{}`)},
		{false, []byte(`{"username": "foo"}`)},
		{false, []byte(`{"password": "bar"}`)},
		{false, []byte(`{"username": "foo","password": "bar"}`)},
	}

	for _, test := range tests {
		actual := login{}
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(test.body))
		r.Header.Set("Content-Type", ContentTypeJSON)
		err := Decode(r, &actual)
		if test.shouldErr {
			assert.NotNil(t, err)
			continue
		}
		assert.Nil(t, err)
		expected := login{validated: true}
		assert.Nil(t, json.Unmarshal(test.body, &expected))
		assert.Equal(t, expected, actual)
	}
}

func TestXML(t *testing.T) {
	tests := []struct {
		shouldErr bool
		body      []byte
	}{
		{true, []byte(`invalid xml`)},
		{false, []byte(`<xml></xml>`)},
		{false, []byte(`<xml><username>foo</username></xml>`)},
		{false, []byte(`<xml><password>bar></password></xml>`)},
		{false, []byte(`<xml><username>foo</username><password>bar></password></xml>`)},
	}

	for _, test := range tests {
		actual := login{}
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(test.body))
		r.Header.Set("Content-Type", ContentTypeXML)
		err := Decode(r, &actual)
		if test.shouldErr {
			assert.NotNil(t, err)
			continue
		}
		assert.Nil(t, err)
		expected := login{validated: true}
		assert.Nil(t, xml.Unmarshal(test.body, &expected))
		assert.Equal(t, expected, actual)
	}
}

var formData = map[string][]string{
	"username": {"foo"},
	"password": {"bar"},
}

func TestForm(t *testing.T) {
	expected := login{validated: true}
	if err := defaultDecoder.Decode(&expected, formData); err != nil {
		t.Fatal(err.Error())
	}

	actual := login{validated: true}
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("username=foo&password=bar"))
	r.Header.Set("Content-Type", ContentTypeForm)
	assert.Nil(t, Decode(r, &actual))
	assert.Equal(t, expected, actual)
}

func TestMultipartForm(t *testing.T) {
	expected := login{}
	err := defaultDecoder.Decode(&expected, formData)
	assert.Nil(t, err)

	actual := login{}
	postData :=
		`--xxx
Content-Disposition: form-data; name="username"

foo
--xxx
Content-Disposition: form-data; name="password"

bar
--xxx--
`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(postData))
	r.Header.Set("Content-Type", ContentTypeMultipartForm+"; boundary=xxx")
	assert.Nil(t, Decode(r, &actual))
	assert.True(t, actual.validated)
	assert.Equal(t, expected.Username, actual.Username)
	assert.Equal(t, expected.Password, actual.Password)
}

func TestParseContentType(t *testing.T) {
	tests := []struct {
		contentType         string
		shouldError         bool
		expectedContentType string
	}{
		{"application/json", false, ContentTypeJSON},
		{"application/json; charset=utf-8", false, ContentTypeJSON},
		{"application/xml", false, ContentTypeXML},
		{"application/xml; charset=utf-8", false, ContentTypeXML},
		{"application/x-www-form-urlencoded", false, ContentTypeForm},
		{"multipart/form-data", false, ContentTypeMultipartForm},
		{"", true, ""},
	}

	for _, test := range tests {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set(ContentType, test.contentType)
		contentType, err := parseContentType(r)
		if test.shouldError {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, test.expectedContentType, contentType)
		}
	}
}

func TestDecode(t *testing.T) {
	delete((*defaultDecoders), "application/json")
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.Header.Set("Content-Type", "application/json")
	v := login{}
	err := Decode(req, &v)
	assert.NotNil(t, err)
}

func TestNewMultipartForm(t *testing.T) {
	tests := []struct {
		maxMemory   int64
		contentType string
		shouldError bool
	}{
		{1024, ContentTypeMultipartForm + "; boundary=xxx", false},
		{1024, "", true},
	}
	for _, test := range tests {
		decoder := NewMultipartForm(test.maxMemory)
		v := login{}
		body := strings.NewReader(`--xxx--`)
		req := httptest.NewRequest(http.MethodPost, "/login", body)
		req.Header.Set("Content-Type", test.contentType)
		err := decoder(req, &v)
		if test.shouldError {
			assert.NotNil(t, err)
		}
	}
}
