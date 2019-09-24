package main

import (
	"errors"
	"net/http"
	"os"

	postory "github.com/andynaguyen/postory-server"
	"github.com/coldbrewcloud/go-shippo"
	"github.com/go-chi/chi"
)

func enableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func main() {
	shippoToken := os.Getenv("SHIPPO_TOKEN")
	if shippoToken == "" {
		panic(errors.New("please set $SHIPPO_TOKEN with your Shippo API private token"))
	}

	shippoClient := shippo.NewClient(shippoToken)
	shippoProxy := postory.ShippoAdapter{shippoClient}
	api := postory.Postory{shippoProxy}

	router := chi.NewRouter()
	router.Use(enableCorsMiddleware)
	router.Get("/tracking_info/{carrier}/{trackingNumber}", api.TrackingInfoHandler())
	router.Get("/tracking_info_history/{carrier}/{trackingNumber}", api.TrackingInfoHistoryHandler())
	http.ListenAndServe(":3000", router)
}
