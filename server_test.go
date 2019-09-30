// +build integration

package postory_server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/coldbrewcloud/go-shippo/models"
	"github.com/stretchr/testify/assert"
)

func TestServer_TrackingInfo(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/track?carrier=shippo&trackingNumber=SHIPPO_TRANSIT")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var trackingInfo TrackingInfo
	json.NewDecoder(resp.Body).Decode(&trackingInfo)
	assert.Equal(t, "San Francisco", trackingInfo.AddressFrom.City)
	assert.Equal(t, "CA", trackingInfo.AddressFrom.State)
	assert.Equal(t, "94103", trackingInfo.AddressFrom.Zip)
	assert.Equal(t, "US", trackingInfo.AddressFrom.Country)
	assert.Equal(t, "Chicago", trackingInfo.AddressTo.City)
	assert.Equal(t, "IL", trackingInfo.AddressTo.State)
	assert.Equal(t, "60611", trackingInfo.AddressTo.Zip)
	assert.Equal(t, "US", trackingInfo.AddressTo.Country)
	assert.Equal(t, "Priority Mail", trackingInfo.ServiceLevel.Name)
	assert.Equal(t, "shippo_priority", trackingInfo.ServiceLevel.Token)
	assert.Equal(t, "", trackingInfo.ServiceLevel.Terms)
	assert.Equal(t, "TRANSIT", trackingInfo.TrackingStatus.Status)
	assert.Equal(t, "Your shipment has departed from the origin.", trackingInfo.TrackingStatus.StatusDetails)
	assert.Equal(t, "San Francisco", trackingInfo.TrackingStatus.Location.City)
	assert.Equal(t, "CA", trackingInfo.TrackingStatus.Location.State)
	assert.Equal(t, "94103", trackingInfo.TrackingStatus.Location.Zip)
	assert.Equal(t, "US", trackingInfo.TrackingStatus.Location.Country)
}

func TestServer_TrackingInfoHistory(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/history?carrier=shippo&trackingNumber=SHIPPO_TRANSIT")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var trackingInfoHistory []*models.TrackingStatusDict
	json.NewDecoder(resp.Body).Decode(&trackingInfoHistory)
	assert.Len(t, trackingInfoHistory, 2)
	assert.Equal(t, "UNKNOWN", trackingInfoHistory[0].Status)
	assert.Equal(t, "The carrier has received the electronic shipment information.", trackingInfoHistory[0].StatusDetails)
	assert.Equal(t, "San Francisco", trackingInfoHistory[0].Location.City)
	assert.Equal(t, "CA", trackingInfoHistory[0].Location.State)
	assert.Equal(t, "94103", trackingInfoHistory[0].Location.Zip)
	assert.Equal(t, "US", trackingInfoHistory[0].Location.Country)
	assert.Equal(t, "TRANSIT", trackingInfoHistory[1].Status)
	assert.Equal(t, "Your shipment has departed from the origin.", trackingInfoHistory[1].StatusDetails)
	assert.Equal(t, "San Francisco", trackingInfoHistory[1].Location.City)
	assert.Equal(t, "CA", trackingInfoHistory[1].Location.State)
	assert.Equal(t, "94103", trackingInfoHistory[1].Location.Zip)
	assert.Equal(t, "US", trackingInfoHistory[1].Location.Country)
}
