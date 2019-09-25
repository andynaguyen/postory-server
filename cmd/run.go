package main

import (
	"errors"
	"os"

	postory "github.com/andynaguyen/postory-server"
)

func main() {
	shippoToken := os.Getenv("SHIPPO_TOKEN")
	if shippoToken == "" {
		panic(errors.New("please set $SHIPPO_TOKEN with your Shippo API private token"))
	}

	config := &postory.ServerConfig{
		ShippoToken: shippoToken,
	}

	postory.StartServer(config)
}
