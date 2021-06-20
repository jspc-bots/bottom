package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jspc/bottom"
)

var (
	Username  = os.Getenv("SASL_USER")
	Password  = os.Getenv("SASL_PASSWORD")
	Server    = os.Getenv("SERVER")
	VerifyTLS = os.Getenv("VERIFY_TLS") == "true"
)

func main() {
	b, err := bottom.New(Username, Password, Server, VerifyTLS)
	if err != nil {
		log.Fatal(err)
	}

	router := bottom.NewRouter()
	router.AddRoute("Hello example", func(sender, channel string, groups []string) error {
		b.Client.Cmd.Messagef(channel, "And hello to you too %s", sender)
		return nil
	})

	router.AddRoute("show an error", func(sender, channel string, groups []string) error {
		return fmt.Errorf("an error :(")
	})

	b.Middlewares.Push(router)

	log.Fatal(b.Client.Connect())
}
