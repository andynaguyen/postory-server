package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/andynaguyen/postory-server"
	"github.com/coldbrewcloud/go-shippo"
	"github.com/go-chi/chi"
)

func main() {
	shippoToken := os.Getenv("SHIPPO_TOKEN")
	if shippoToken == "" {
		panic(errors.New("please set $SHIPPO_TOKEN with your Shippo API private token"))
	}

	shippoClient := shippo.NewClient(shippoToken)
	api := postory_server.PostoryApi{shippoClient}

	router := chi.NewRouter()
	router.Get("/track/{carrier}/{trackingNumber}", api.TrackHandler())
	http.ListenAndServe(":3000", router)
}
