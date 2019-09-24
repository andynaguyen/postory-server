package postory_server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/coldbrewcloud/go-shippo/models"
	"github.com/go-chi/chi"
)

type TrackingInfo struct {
	AddressFrom    *models.TrackingStatusLocation `json:"address_from,omitempty"`
	AddressTo      *models.TrackingStatusLocation `json:"address_from,omitempty"`
	ETA            time.Time                      `json:"eta"`
	ServiceLevel   *models.ServiceLevel           `json:"servicelevel,omitempty"`
	TrackingStatus *models.TrackingStatusDict     `json:"tracking_status,omitempty"`
}

type TrackingInfoResponse struct {
	Error      *string       `json:"error"`
	StatusCode int           `json:"status_code"`
	Data       *TrackingInfo `json:"data"`
}

type TrackingInfoHistoryResponse struct {
	Error      *string                      `json:"error"`
	StatusCode int                          `json:"status_code"`
	Data       []*models.TrackingStatusDict `json:"tracking_history"`
}

type PackageTrackingClient interface {
	GetTrackingUpdate(string, string) (*models.TrackingStatus, error)
}

type Postory struct {
	Adapter TrackingAdapter
}

func (p Postory) TrackingInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trackingNumber := chi.URLParam(r, "trackingNumber")
		carrier := chi.URLParam(r, "carrier")
		trackingInfo := p.Adapter.GetTrackingInfo(carrier, trackingNumber)
		w.WriteHeader(trackingInfo.StatusCode)

		bytes, err := json.Marshal(trackingInfo)
		if err != nil {
			println(err.Error())
		}
		w.Write(bytes)
	}
}

func (p Postory) TrackingInfoHistoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trackingNumber := chi.URLParam(r, "trackingNumber")
		carrier := chi.URLParam(r, "carrier")
		trackingInfoHistory := p.Adapter.GetTrackingInfoHistory(carrier, trackingNumber)
		w.WriteHeader(trackingInfoHistory.StatusCode)

		bytes, err := json.Marshal(trackingInfoHistory)
		if err != nil {
			println(err.Error())
		}
		w.Write(bytes)
	}
}
