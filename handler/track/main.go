package main

import (
	"encoding/json"
	"net/http"

	postory "github.com/andynaguyen/postory-server"
	"github.com/andynaguyen/postory-server/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/coldbrewcloud/go-shippo/models"
)

var shippoAdapter = postory.NewShippoAdapter()
var archive = postory.NewArchive()
var logger = postory.NewLogger()

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
		logger.Error().Err(err).Msg("error validating input")
		response.StatusCode = http.StatusBadRequest
		return response, err
	}

	// check archive first, otherwise call shippo
	trackingInfo := archive.GetInfo(carrier, trackingNumber)
	if trackingInfo == nil {
		logger.Info().Msg("not found in archive, callling shippo")
		trackingInfo, err = shippoAdapter.GetTrackingInfo(carrier, trackingNumber)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get tracking info")
			response.StatusCode = http.StatusInternalServerError
			return response, err
		}
	}

	logger.Printf("tracking info: %+v\n", trackingInfo)
	// If terminal status, archive the tracking info
	if trackingInfo != nil && trackingInfo.TrackingStatus != nil {
		isTerminal := trackingInfo.TrackingStatus.Status == models.TrackingStatusStatusDelivered ||
			trackingInfo.TrackingStatus.Status == models.TrackingStatusStatusFailure ||
			trackingInfo.TrackingStatus.Status == models.TrackingStatusStatusReturned
		if isTerminal {
			archive.PutInfo(*trackingInfo)
		}
	}

	bodyBytes, err := json.Marshal(trackingInfo)
	if err != nil {
		logger.Error().Err(err).Msg("error marshalling tracking info")
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
