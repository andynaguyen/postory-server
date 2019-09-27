// +build integration

package postory_server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_TrackingInfo(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/tracking_info/shippo/SHIPPO_TRANSIT")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var trackingInfo TrackingInfoResponse
	json.NewDecoder(resp.Body).Decode(&trackingInfo)
	assert.Nil(t, trackingInfo.Error)
	assert.Equal(t, http.StatusOK, trackingInfo.StatusCode)

	data := trackingInfo.Data
	assert.Equal(t, "Priority Mail", data.ServiceLevel.Name)
	assert.Equal(t, "shippo_priority", data.ServiceLevel.Token)
	assert.Equal(t, "", data.ServiceLevel.Terms)
	assert.Equal(t, "TRANSIT", data.TrackingStatus.Status)
	assert.Equal(t, "Your shipment has departed from the origin.", data.TrackingStatus.StatusDetails)
	assert.Equal(t, "San Francisco", data.TrackingStatus.Location.City)
	assert.Equal(t, "CA", data.TrackingStatus.Location.State)
	assert.Equal(t, "94103", data.TrackingStatus.Location.Zip)
	assert.Equal(t, "US", data.TrackingStatus.Location.Country)
}

func TestServer_TrackingInfoHistory(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/tracking_info_history/shippo/SHIPPO_TRANSIT")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var trackingInfoHistory TrackingInfoHistoryResponse
	json.NewDecoder(resp.Body).Decode(&trackingInfoHistory)
	assert.Nil(t, trackingInfoHistory.Error)
	assert.Equal(t, http.StatusOK, trackingInfoHistory.StatusCode)

	data := trackingInfoHistory.Data
	assert.Len(t, data, 2)

	assert.Equal(t, "UNKNOWN", data[0].Status)
	assert.Equal(t, "The carrier has received the electronic shipment information.", data[0].StatusDetails)
	assert.Equal(t, "San Francisco", data[0].Location.City)
	assert.Equal(t, "CA", data[0].Location.State)
	assert.Equal(t, "94103", data[0].Location.Zip)
	assert.Equal(t, "US", data[0].Location.Country)
	assert.Equal(t, "TRANSIT", data[1].Status)
	assert.Equal(t, "Your shipment has departed from the origin.", data[1].StatusDetails)
	assert.Equal(t, "San Francisco", data[1].Location.City)
	assert.Equal(t, "CA", data[1].Location.State)
	assert.Equal(t, "94103", data[1].Location.Zip)
	assert.Equal(t, "US", data[1].Location.Country)
}
