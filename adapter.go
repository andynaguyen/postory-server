package postory_server

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/coldbrewcloud/go-shippo"
	"github.com/coldbrewcloud/go-shippo/client"
	"github.com/coldbrewcloud/go-shippo/models"
)

type TrackingInfo struct {
	Carrier        string                         `json:"carrier"`
	TrackingNumber string                         `json:"tracking_number"`
	AddressFrom    *models.TrackingStatusLocation `json:"address_from"`
	AddressTo      *models.TrackingStatusLocation `json:"address_to"`
	ETA            time.Time                      `json:"eta"`
	TrackingStatus *models.TrackingStatusDict     `json:"tracking_status"`
}

type ShippoAdapter struct {
	cl *client.Client
}

func NewShippoAdapter() *ShippoAdapter {
	shippoToken := os.Getenv("SHIPPO_TOKEN")
	return &ShippoAdapter{shippo.NewClient(shippoToken)}
}

func (adapter ShippoAdapter) GetTrackingInfo(carrier string, trackingNumber string) (*TrackingInfo, error) {
	trackingStatus, err := adapter.cl.GetTrackingUpdate(carrier, trackingNumber)
	if err != nil {
		msg := fmt.Sprintf("could not get tracking update for carrier \"%s\" and tracking number \"%s\"", carrierTokensByName[carrier], trackingNumber)
		return nil, errors.New(msg)
	}

	return &TrackingInfo{
		Carrier:        carrier,
		TrackingNumber: trackingNumber,
		AddressFrom:    trackingStatus.AddressFrom,
		AddressTo:      trackingStatus.AddressTo,
		ETA:            trackingStatus.ETA,
		TrackingStatus: trackingStatus.TrackingStatus,
	}, nil
}

func (adapter ShippoAdapter) GetTrackingInfoHistory(carrier string, trackingNumber string) ([]*models.TrackingStatusDict, error) {
	trackingStatus, err := adapter.cl.GetTrackingUpdate(carrier, trackingNumber)
	if err != nil {
		msg := fmt.Sprintf("could not get tracking update for carrier \"%s\" and tracking number \"%s\"", carrierTokensByName[carrier], trackingNumber)
		return nil, errors.New(msg)
	}

	return trackingStatus.TrackingHistory, nil
}
