package postory_server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/coldbrewcloud/go-shippo"
	"github.com/coldbrewcloud/go-shippo/client"
	"github.com/coldbrewcloud/go-shippo/models"
)

type TrackingInfo struct {
	AddressFrom    *models.TrackingStatusLocation `json:"address_from,omitempty"`
	AddressTo      *models.TrackingStatusLocation `json:"address_from,omitempty"`
	ETA            time.Time                      `json:"eta"`
	ServiceLevel   *models.ServiceLevel           `json:"servicelevel,omitempty"`
	TrackingStatus *models.TrackingStatusDict     `json:"tracking_status,omitempty"`
}

type TrackingInfoResponse struct {
	Error      error         `json:"error"`
	StatusCode int           `json:"status_code"`
	Data       *TrackingInfo `json:"data"`
}

type TrackingInfoHistory []*models.TrackingStatusDict

type TrackingInfoHistoryResponse struct {
	Error      error               `json:"error"`
	StatusCode int                 `json:"status_code"`
	Data       TrackingInfoHistory `json:"tracking_history"`
}

type TrackingAdapter interface {
	GetTrackingInfo(string, string) TrackingInfoResponse
	GetTrackingInfoHistory(string, string) TrackingInfoHistoryResponse
}

type ShippoAdapter struct {
	cl *client.Client
}

func NewShippoAdapter(token string) ShippoAdapter {
	return ShippoAdapter{shippo.NewClient(token)}
}

func (adapter ShippoAdapter) GetTrackingInfo(carrier string, trackingNumber string) TrackingInfoResponse {
	if !isCarrierSupported(carrier) {
		msg := fmt.Sprintf("carrier is unsupported: %s", carrier)
		return TrackingInfoResponse{errors.New(msg), http.StatusBadRequest, nil}
	}

	trackingStatus, err := adapter.cl.GetTrackingUpdate(carrier, trackingNumber)
	if err != nil {
		msg := fmt.Sprintf("could not get tracking update for carrier \"%s\" and tracking number \"%s\"", carrierTokensByName[carrier], trackingNumber)
		return TrackingInfoResponse{errors.New(msg), http.StatusInternalServerError, nil}
	}

	return TrackingInfoResponse{nil, http.StatusOK, &TrackingInfo{
		AddressFrom:    trackingStatus.AddressFrom,
		AddressTo:      trackingStatus.AddressTo,
		ETA:            trackingStatus.ETA,
		ServiceLevel:   trackingStatus.ServiceLevel,
		TrackingStatus: trackingStatus.TrackingStatus,
	}}
}

func (adapter ShippoAdapter) GetTrackingInfoHistory(carrier string, trackingNumber string) TrackingInfoHistoryResponse {
	if !isCarrierSupported(carrier) {
		msg := fmt.Sprintf("carrier is unsupported: %s", carrier)
		return TrackingInfoHistoryResponse{errors.New(msg), http.StatusBadRequest, nil}
	}

	trackingStatus, err := adapter.cl.GetTrackingUpdate(carrier, trackingNumber)
	if err != nil {
		msg := fmt.Sprintf("could not get tracking update for carrier \"%s\" and tracking number \"%s\"", carrierTokensByName[carrier], trackingNumber)
		return TrackingInfoHistoryResponse{errors.New(msg), http.StatusInternalServerError, nil}
	}

	return TrackingInfoHistoryResponse{nil, http.StatusOK, trackingStatus.TrackingHistory}
}
