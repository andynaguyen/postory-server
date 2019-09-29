package main

import (
	"encoding/json"
	"net/http"
	"os"

	postory "github.com/andynaguyen/postory-server"
	"github.com/andynaguyen/postory-server/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/coldbrewcloud/go-shippo/models"
)

var shippoToken = os.Getenv("SHIPPO_TOKEN")
var shippoAdapter = postory.NewShippoAdapter(shippoToken)
var archiveDao = postory.NewArchiveDao()

func handle(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}

	carrier := event.QueryStringParameters["carrier"]
	trackingNumber := event.QueryStringParameters["trackingNumber"]
	err := handler.ValidateInput(carrier)
	if err != nil {
		println(err.Error())
		response.StatusCode = http.StatusBadRequest
		return response, err
	}

	// check archive first, otherwise call shippo
	trackingInfo := archiveDao.GetArchivedInfo(carrier, trackingNumber)
	if trackingInfo == nil {
		println("not found in archive, callling shippo")
		trackingInfo, err = shippoAdapter.GetTrackingInfo(carrier, trackingNumber)
		if err != nil {
			println(err.Error())
			response.StatusCode = http.StatusInternalServerError
			return response, err
		}
	}

	// If terminal status, archive the tracking info
	if trackingInfo.TrackingStatus.Status == models.TrackingStatusStatusDelivered ||
		trackingInfo.TrackingStatus.Status == models.TrackingStatusStatusFailure ||
		trackingInfo.TrackingStatus.Status == models.TrackingStatusStatusReturned {
		if err = archiveDao.PutArchivedInfo(carrier, trackingNumber, *trackingInfo); err != nil {
			println(err.Error())
		}
	}

	bodyBytes, err := json.Marshal(trackingInfo)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}
	response.Body = string(bodyBytes)
	response.StatusCode = http.StatusOK
	return response, nil
}

func main() {
	lambda.Start(handle)
}
