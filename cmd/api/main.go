package main

import (
	"fmt"
	"oAuthTest/internal/oauth"
	"oAuthTest/internal/server"
)

func main() {

	oauth.NewAuth()
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
