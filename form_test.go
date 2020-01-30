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
	"reflect"
	"strings"
	"testing"
)

type login struct {
	Username string `schema:"username" json:"username" xml:"username"`
	Password string `schema:"password" json:"password" xml:"password"`
}

func TestRegister(t *testing.T) {
	invoked := false
	decoder := func(r *http.Request, v interface{}) error {
		invoked = true
		return nil
	}
	contentType := "content/type"
	Register(contentType, decoder)
	actual, ok := decoders[contentType]
	if !ok {
		t.Error("Failed to register decoder")
	}

	err := actual(httptest.NewRequest(http.MethodGet, "/", nil), nil)
	if err != nil || !invoked {
		t.Error("Failed to register decoder")
	}
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
			if err == nil {
				t.Fatal("Expected decode error")
			}
			return
		}

		if err != nil {
			t.Fatal(err.Error())
		}

		expected := login{}
		if err := json.Unmarshal(test.body, &expected); err != nil {
			t.Fatal(err.Error())
		}

		if !test.shouldErr && !reflect.DeepEqual(expected, actual) {
			t.Error("Failed to decode JSON")
		}
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
			if err == nil {
				t.Fatal("Expected decode error")
			}
			return
		}

		if err != nil {
			t.Fatal(err.Error())
		}

		expected := login{}
		if err := xml.Unmarshal(test.body, &expected); err != nil {
			t.Fatal(err.Error())
		}

		if !test.shouldErr && !reflect.DeepEqual(expected, actual) {
			t.Error("Failed to decode XML")
		}
	}
}

var formData = map[string][]string{
	"username": []string{"foo"},
	"password": []string{"bar"},
}

func TestForm(t *testing.T) {
	expected := login{}
	if err := defaultDecoder.Decode(&expected, formData); err != nil {
		t.Fatal(err.Error())
	}

	actual := login{}
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("username=foo&password=bar"))
	r.Header.Set("Content-Type", ContentTypeForm)
	if err := Decode(r, &actual); err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Failed to decode post form")
	}
}

func TestMultipartForm(t *testing.T) {
	expected := login{}
	if err := defaultDecoder.Decode(&expected, formData); err != nil {
		t.Fatal(err.Error())
	}

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
	if err := Decode(r, &actual); err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Failed to decode multipart form")
	}
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
			if err == nil {
				t.Error("expected an error, got nil")
			}
			continue
		}

		if contentType != test.expectedContentType {
			t.Errorf("expected content type %q, ot %q", test.expectedContentType, contentType)
		}
	}
}

func TestDecode(t *testing.T) {
	delete(decoders, "application/json")
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.Header.Set("Content-Type", "application/json")
	v := login{}
	err := Decode(req, &v)
	if err == nil {
		t.Error("expected an error, got nil")
	}
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
			if err == nil {
				t.Error("expected an error, got nil")
			}
		}
	}
}
