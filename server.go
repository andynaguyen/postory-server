package postory_server

import (
	"net/http"

	"github.com/go-chi/chi"
)

func enableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

type ServerConfig struct {
	ShippoToken string
}

func StartServer(config *ServerConfig) {
	shippoAdapter := NewShippoAdapter(config.ShippoToken)
	api := Postory{shippoAdapter}

	router := chi.NewRouter()
	router.Use(enableCorsMiddleware)
	router.Get("/tracking_info/{carrier}/{trackingNumber}", api.TrackingInfoHandler())
	router.Get("/tracking_info_history/{carrier}/{trackingNumber}", api.TrackingInfoHistoryHandler())
	http.ListenAndServe(":3000", router)
}
