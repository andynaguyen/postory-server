package main

import (
	"encoding/json"
	"os"

	postory "github.com/andynaguyen/postory-server"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var shippoToken = os.Getenv("SHIPPO_TOKEN")
var shippoAdapter = postory.NewShippoAdapter(shippoToken)

func handle(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	carrier := event.QueryStringParameters["carrier"]
	trackingNumber := event.QueryStringParameters["trackingNumber"]
	trackingInfoHistory := shippoAdapter.GetTrackingInfoHistory(carrier, trackingNumber)

	respBytes, err := json.Marshal(trackingInfoHistory.Data)
	if err != nil {
		println(err.Error())
	}

	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      trackingInfoHistory.StatusCode,
		Headers:         headers,
		Body:            string(respBytes),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handle)
}
