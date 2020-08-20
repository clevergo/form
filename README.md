# Form Decoder 
[![Build Status](https://img.shields.io/travis/clevergo/form?style=for-the-badge)](https://travis-ci.org/clevergo/form)
[![Coverage Status](https://img.shields.io/coveralls/github/clevergo/form?style=for-the-badge)](https://coveralls.io/github/clevergo/form?branch=master)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/clevergo.tech/form?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/clevergo/form?style=for-the-badge)](https://goreportcard.com/report/github.com/clevergo/form)
[![Release](https://img.shields.io/github/release/clevergo/form.svg?style=for-the-badge)](https://github.com/clevergo/form/releases)
[![Downloads](https://img.shields.io/endpoint?url=https://pkg.clevergo.tech/api/badges/downloads/month/clevergo.tech/form&style=for-the-badge)](https://pkg.clevergo.tech/)

A form decoder that decode request body of any types(xml, json, form, multipart form...) into a sctruct by same codebase.

By default, form decoder can handles the following content types:

- Form(application/x-www-form-urlencoded)
- Multipart Form(multipart/form-data)
- JSON(application/json)
- XML(application/xml)

> Form and multipart form are built on top of gorilla [schema](https://github.com/gorilla/schema), tag name is `schema`.

[Register](https://pkg.go.dev/clevergo.tech/form?tab=doc#Decoders.Register) allow to register particular decoder or replace default decoder 
for the specified content type.

## Installation

```go
$ go get clevergo.tech/form
```

## Usage

```go
import (
	"net/http"

	"clevergo.tech/form"
)

var decoders = form.New()

type user struct {
	Username string `schema:"username" json:"username" xml:"username"`
	Password string `schema:"password" json:"password" xml:"password"`
}

func init() {
	// replaces multipart form decoder.
	decoders.Register(form.ContentTypeMultipartForm, form.NewMultipartForm(10*1024*1024))
	// registers other decoder
	// decoders.Register(contentType, decoder)
}

func(w http.ResponseWriter, r *http.Request) {
	u := user{}
	if err := decoders.Decode(r, &u); err != nil {
		http.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	// ...
}
```

### Example

Checkout [example](https://github.com/clevergo/examples/tree/master/form) for details.
