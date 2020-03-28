# Form Decoder [![Build Status](https://travis-ci.org/clevergo/form.svg?branch=master)](https://travis-ci.org/clevergo/form) [![Coverage Status](https://coveralls.io/repos/github/clevergo/form/badge.svg?branch=master)](https://coveralls.io/github/clevergo/form?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/clevergo/form)](https://goreportcard.com/report/github.com/clevergo/form) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/clevergo/form) [![Release](https://img.shields.io/github/release/clevergo/form.svg?style=flat-square)](https://github.com/clevergo/form/releases)

A form deocder that decode request body of any types(xml, json, form, multipart form...) into a scturct by same codebase.

By default, form decoder can handles the following content types:

- Form(application/x-www-form-urlencoded)
- Multipart Form(multipart/form-data)
- JSON(application/json)
- XML(application/xml)

> Form and multipart form are built on top of gorilla [schema](https://github.com/gorilla/schema), tag name is `json`.

[Register](https://pkg.go.dev/github.com/clevergo/form#Register) allow to register particular decoder or replace default decoder 
for the specified content type.

## Installation

```go
$ go get github.com/clevergo/form
```

## Usage

```go
import (
	"net/http"

	"github.com/clevergo/form"

var decoders = form.New()

func init() {
	// replaces multipart form decoder.
	decoders.Register(form.ContentTypeMultipartForm, form.NewMultipartForm(10*1024*1024))
	// registers other decoder
	// decoders.Register(contentType, decoder)
}

func(w http.ResponseWriter, r *http.Request) {
	err := form
	if err := decoders.Decode(ctx.Request, &u); err != nil {
		http.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	// ...
}
```

### Example

See [Example](example).