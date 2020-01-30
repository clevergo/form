# Form Decoder [![Build Status](https://travis-ci.org/clevergo/form.svg?branch=master)](https://travis-ci.org/clevergo/form) [![Coverage Status](https://coveralls.io/repos/github/clevergo/form/badge.svg?branch=master)](https://coveralls.io/github/clevergo/form?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/clevergo/form)](https://goreportcard.com/report/github.com/clevergo/form) [![GoDoc](https://godoc.org/github.com/clevergo/form?status.svg)](http://godoc.org/github.com/clevergo/form)

A form deocder that decode request body into a scturct, you don't need to care about what type the request data is.

## Usage

By default, form decoder can handles the following content types:

- Form(application/x-www-form-urlencoded)
- Multipart Form(multipart/form-data)
- JSON(application/json)
- XML(application/xml)

[Register](https://godoc.org/github.com/clevergo/form#Register) allow to register particular decoder or replace default decoder 
for the specified content type.


```go
package main

import (
	"fmt"
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/form"
)

type User struct {
	Username string `json:"username" xml:"username"`
	Password string `json:"password" xml:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
	user := User{}
	err := form.Decode(r, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	fmt.Fprintf(w, "username: %s, password: %s", user.Username, user.Password)
}

func main() {
	app := clevergo.New("localhost:1234")
	app.Post("/login", login)
	app.ListenAndServe()
}
```

```shell
$ curl -XPOST -d "username=foo&password=bar"  http://localhost:1234/login
username: foo, password: bar

$ curl -XPOST -H "Content-Type: application/json" -d '{"username":"foo", "password": "bar"}' http://localhost:1234/login
username: foo, password: bar

$ curl -XPOST -H "Content-Type: application/xml" -d '<xml><username>foo</username><password>bar</password></xml>' http://localhost:1234/login
username: foo, password: bar

$ curl -XPOST -F "username=foo" -F "password=bar" http://localhost:1234/login
username: foo, password: bar
```