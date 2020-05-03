package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/form"
)

var decoders = form.New()

type user struct {
	Username string `schema:"username" json:"username" xml:"username"`
	Password string `schema:"password" json:"password" xml:"password"`
}

func login(ctx *clevergo.Context) error {
	u := user{}
	if err := decoders.Decode(ctx.Request, &u); err != nil {
		return err
	}
	ctx.WriteString(fmt.Sprintf("username: %s, password: %s", u.Username, u.Password))
	return nil
}

func main() {
	router := clevergo.NewRouter()
	router.Post("/login", login)
	log.Println(http.ListenAndServe(":12345", router))
}
