package postory_server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

func enableCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func trackingInfoHandler(adapter TrackingAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trackingNumber := chi.URLParam(r, "trackingNumber")
		carrier := chi.URLParam(r, "carrier")
		trackingInfo := adapter.GetTrackingInfo(carrier, trackingNumber)
		w.WriteHeader(trackingInfo.StatusCode)

		bytes, err := json.Marshal(trackingInfo)
		if err != nil {
			println(err.Error())
		}
		w.Write(bytes)
	}
}

func trackingInfoHistoryHandler(adapter TrackingAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trackingNumber := chi.URLParam(r, "trackingNumber")
		carrier := chi.URLParam(r, "carrier")
		trackingInfoHistory := adapter.GetTrackingInfoHistory(carrier, trackingNumber)
		w.WriteHeader(trackingInfoHistory.StatusCode)

		bytes, err := json.Marshal(trackingInfoHistory)
		if err != nil {
			println(err.Error())
		}
		w.Write(bytes)
	}
}

type ServerConfig struct {
	ShippoToken string
}

func StartServer(config *ServerConfig) {
	shippoAdapter := NewShippoAdapter(config.ShippoToken)

	router := chi.NewRouter()
	router.Use(enableCorsMiddleware)
	router.Get("/tracking_info/{carrier}/{trackingNumber}", trackingInfoHandler(shippoAdapter))
	router.Get("/tracking_info_history/{carrier}/{trackingNumber}", trackingInfoHistoryHandler(shippoAdapter))
	println("listening on port 3000")
	http.ListenAndServe(":3000", router)
}
