package main

import (
	"encoding/json"
	"net/http"

	postory "github.com/andynaguyen/postory-server"
	"github.com/andynaguyen/postory-server/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var shippoAdapter = postory.NewShippoAdapter()
var logger = postory.NewLogger()

func handle(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}

	carrier := event.QueryStringParameters["carrier"]
	trackingNumber := event.QueryStringParameters["trackingNumber"]
	if err := handler.ValidateInput(carrier); err != nil {
		logger.Error().Err(err).Msg("error validating input")
		response.StatusCode = http.StatusBadRequest
		return response, err
	}

	trackingInfoHistory, err := shippoAdapter.GetTrackingInfoHistory(carrier, trackingNumber)
	logger.Printf("returning tracking info history: %+v\n", trackingInfoHistory)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get tracking info history")
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}

	bodyBytes, err := json.Marshal(trackingInfoHistory)
	if err != nil {
		logger.Error().Err(err).Msg("error marshalling tracking info history")
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
