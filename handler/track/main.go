package main

import (
	"encoding/json"
	"net/http"
	"os"

	postory "github.com/andynaguyen/postory-server"
	"github.com/andynaguyen/postory-server/handler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var shippoToken = os.Getenv("SHIPPO_TOKEN")
var shippoAdapter = postory.NewShippoAdapter(shippoToken)

func handle(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}

	carrier := event.QueryStringParameters["carrier"]
	trackingNumber := event.QueryStringParameters["trackingNumber"]
	if err := handler.ValidateInput(carrier); err != nil {
		println(err.Error())
		response.StatusCode = http.StatusBadRequest
		return response, err
	}

	trackingInfo, err := shippoAdapter.GetTrackingInfo(carrier, trackingNumber)
	bodyBytes, err := json.Marshal(trackingInfo)
	if err != nil {
		println(err.Error())
		response.StatusCode = http.StatusInternalServerError
		return response, err
	}

	response.StatusCode = http.StatusOK
	response.Body = string(bodyBytes)
	return response, nil
}

func main() {
	lambda.Start(handle)
}
