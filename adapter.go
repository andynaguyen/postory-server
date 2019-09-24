package postory_server

import (
	"fmt"
	"net/http"

	"github.com/coldbrewcloud/go-shippo/client"
)

type TrackingAdapter interface {
	GetTrackingInfo(string, string) TrackingInfoResponse
	GetTrackingInfoHistory(string, string) TrackingInfoHistoryResponse
}

type ShippoAdapter struct {
	*client.Client
}

func (adapter ShippoAdapter) GetTrackingInfo(carrier string, trackingNumber string) TrackingInfoResponse {
	if !isCarrierSupported(carrier) {
		msg := fmt.Sprintf("carrier is unsupported: %s", carrier)
		return TrackingInfoResponse{&msg, http.StatusBadRequest, nil}
	}

	trackingStatus, err := adapter.GetTrackingUpdate(carrier, trackingNumber)
	if err != nil {
		msg := fmt.Sprintf("could not get tracking update for carrier [%s] and tracking number [%s]\n", carrier, trackingNumber)
		return TrackingInfoResponse{&msg, http.StatusInternalServerError, nil}
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
		return TrackingInfoHistoryResponse{&msg, http.StatusBadRequest, nil}
	}

	trackingStatus, err := adapter.GetTrackingUpdate(carrier, trackingNumber)
	if err != nil {
		msg := fmt.Sprintf("could not get tracking update for carrier [%s] and tracking number [%s]\n", carrier, trackingNumber)
		return TrackingInfoHistoryResponse{&msg, http.StatusInternalServerError, nil}
	}

	return TrackingInfoHistoryResponse{nil, http.StatusOK, trackingStatus.TrackingHistory}
}
