// +build integration

package postory_server

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ShippoCarrier    = "shippo"
	ShippoPreTransit = "SHIPPO_PRE_TRANSIT"
	ShippoTransit    = "SHIPPO_TRANSIT"
	ShippoDelivered  = "SHIPPO_DELIVERED"
	ShippoReturned   = "SHIPPO_RETURNED"
	ShippoFailure    = "SHIPPO_FAILURE"
	ShippoUnknown    = "SHIPPO_UNKNOWN"
)

var shippoToken = os.Getenv("SHIPPO_TOKEN")
var adapter = NewShippoAdapter(shippoToken)

func TestShippoAdapter_GetTrackingInfo_Transit(t *testing.T) {
	resp := adapter.GetTrackingInfo(ShippoCarrier, ShippoTransit)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, resp.Error)

	data := resp.Data
	assert.Nil(t, data.AddressFrom)
	assert.Nil(t, data.AddressTo)
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

func TestShippoAdapter_GetTrackingInfo_Delivered(t *testing.T) {
	resp := adapter.GetTrackingInfo(ShippoCarrier, ShippoDelivered)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, resp.Error)

	data := resp.Data
	assert.Nil(t, data.AddressFrom)
	assert.Nil(t, data.AddressTo)
	assert.Equal(t, "Priority Mail", data.ServiceLevel.Name)
	assert.Equal(t, "shippo_priority", data.ServiceLevel.Token)
	assert.Equal(t, "", data.ServiceLevel.Terms)
	assert.Equal(t, "DELIVERED", data.TrackingStatus.Status)
	assert.Equal(t, "Your shipment has been delivered.", data.TrackingStatus.StatusDetails)
	assert.Equal(t, "Chicago", data.TrackingStatus.Location.City)
	assert.Equal(t, "IL", data.TrackingStatus.Location.State)
	assert.Equal(t, "60611", data.TrackingStatus.Location.Zip)
	assert.Equal(t, "US", data.TrackingStatus.Location.Country)
}

func TestShippoAdapter_GetTrackingInfo_Failure(t *testing.T) {
	resp := adapter.GetTrackingInfo(ShippoCarrier, ShippoFailure)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, resp.Error)

	data := resp.Data
	assert.Nil(t, data.AddressFrom)
	assert.Nil(t, data.AddressTo)
	assert.Equal(t, "Priority Mail", data.ServiceLevel.Name)
	assert.Equal(t, "shippo_priority", data.ServiceLevel.Token)
	assert.Equal(t, "", data.ServiceLevel.Terms)
	assert.Equal(t, "FAILURE", data.TrackingStatus.Status)
	assert.Equal(t, "The Postal Service has identified a problem with the processing of this item and you should contact support to get further information.", data.TrackingStatus.StatusDetails)
	assert.Equal(t, "Memphis", data.TrackingStatus.Location.City)
	assert.Equal(t, "TN", data.TrackingStatus.Location.State)
	assert.Equal(t, "37501", data.TrackingStatus.Location.Zip)
	assert.Equal(t, "US", data.TrackingStatus.Location.Country)
}

func TestShippoAdapter_GetTrackingInfo_UnsupportedCarrier(t *testing.T) {
	resp := adapter.GetTrackingInfo("", "")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NotNil(t, resp.Error)
}

func TestShippoAdapter_GetTrackingInfo_ShippoError(t *testing.T) {
	resp := adapter.GetTrackingInfo(ShippoCarrier, "")

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NotNil(t, resp.Error)
}

func TestShippoAdapter_GetTrackingInfoHistory_PreTransit(t *testing.T) {
	resp := adapter.GetTrackingInfoHistory(ShippoCarrier, ShippoPreTransit)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, resp.Error)

	assert.Empty(t, resp.Data)
}

func TestShippoAdapter_GetTrackingInfoHistory_Transit(t *testing.T) {
	resp := adapter.GetTrackingInfoHistory(ShippoCarrier, ShippoTransit)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, resp.Error)

	data := resp.Data
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

func TestShippoAdapter_GetTrackingInfoHistory_UnsupportedCarrier(t *testing.T) {
	resp := adapter.GetTrackingInfoHistory("", "")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NotNil(t, resp.Error)
}

func TestShippoAdapter_GetTrackingInfoHistory_ShippoError(t *testing.T) {
	resp := adapter.GetTrackingInfoHistory(ShippoCarrier, "")

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NotNil(t, resp.Error)
}
