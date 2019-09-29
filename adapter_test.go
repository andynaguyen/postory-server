// +build integration

package postory_server

import (
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

var adapter = NewShippoAdapter(shippoToken)

func TestShippoAdapter_GetTrackingInfo_Transit(t *testing.T) {
	data, err := adapter.GetTrackingInfo(ShippoCarrier, ShippoTransit)

	assert.Nil(t, err)
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
	data, err := adapter.GetTrackingInfo(ShippoCarrier, ShippoDelivered)

	assert.Nil(t, err)
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
	data, err := adapter.GetTrackingInfo(ShippoCarrier, ShippoFailure)

	assert.Nil(t, err)
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
	_, err := adapter.GetTrackingInfo("", "")

	assert.Error(t, err)
}

func TestShippoAdapter_GetTrackingInfo_ShippoError(t *testing.T) {
	_, err := adapter.GetTrackingInfo(ShippoCarrier, "")

	assert.Error(t, err)
}

func TestShippoAdapter_GetTrackingInfoHistory_PreTransit(t *testing.T) {
	data, err := adapter.GetTrackingInfoHistory(ShippoCarrier, ShippoPreTransit)

	assert.Nil(t, err)
	assert.Empty(t, data)
}

func TestShippoAdapter_GetTrackingInfoHistory_Transit(t *testing.T) {
	data, err := adapter.GetTrackingInfoHistory(ShippoCarrier, ShippoTransit)

	assert.Nil(t, err)
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
	_, err := adapter.GetTrackingInfoHistory("", "")

	assert.Error(t, err)
}

func TestShippoAdapter_GetTrackingInfoHistory_ShippoError(t *testing.T) {
	_, err := adapter.GetTrackingInfoHistory(ShippoCarrier, "")

	assert.Error(t, err)
}
