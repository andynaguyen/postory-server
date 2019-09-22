package postory_server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andynaguyen/postory-server/carriers"
	"github.com/coldbrewcloud/go-shippo/client"
	"github.com/coldbrewcloud/go-shippo/models"
	"github.com/go-chi/chi"
)

type TrackedPackage struct {
	Error      *string                `json:"error"`
	StatusCode int                    `json:"status_code"`
	Data       *models.TrackingStatus `json:"data"`
}

type PostoryApi struct {
	ShippoClient *client.Client
}

func (p PostoryApi) TrackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		trackingNumber := chi.URLParam(r, "trackingNumber")
		carrier := chi.URLParam(r, "carrier")
		trackedPackage := p.getTrackedPackage(carrier, trackingNumber)
		w.WriteHeader(trackedPackage.StatusCode)

		bytes, err := json.Marshal(trackedPackage)
		if err != nil {
			println(err.Error())
		}
		w.Write(bytes)
	}
}

func (p PostoryApi) getTrackedPackage(carrier string, trackingNumber string) TrackedPackage {
	if !carriers.IsSupported(carrier) {
		msg := fmt.Sprintf("carrier is unsupported: %s", carrier)
		return TrackedPackage{&msg, http.StatusBadRequest, nil}
	}

	trackingStatus, err := p.ShippoClient.GetTrackingUpdate(carrier, trackingNumber)
	if err != nil {
		msg := fmt.Sprintf("could not get tracking update for carrier [%s] and tracking number [%s]\n", carrier, trackingNumber)
		return TrackedPackage{&msg, http.StatusInternalServerError, trackingStatus}
	}

	return TrackedPackage{nil, http.StatusOK, trackingStatus}
}
