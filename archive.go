package postory_server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type ArchiveItem struct {
	Id           string
	TrackingInfo TrackingInfo
}

type Archive struct {
	svc *dynamodb.DynamoDB
}

var logger = NewLogger()

func NewArchive() *Archive {
	sess, err := session.NewSession()
	if err != nil {
		logger.Error().Err(err).Msg("could not start aws session")
		return nil
	}
	return &Archive{dynamodb.New(sess)}
}

const TableName = "archive"

func (a Archive) GetInfo(carrier string, trackingNumber string) *TrackingInfo {
	log := logger.With().
		Str("carrier", carrier).
		Str("trackingNumber", trackingNumber).
		Logger()

	output, err := a.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {S: aws.String(getId(carrier, trackingNumber))},
		},
	})

	if err != nil {
		log.Error().Err(err).Msg("error getting from dynamodb table")
		return nil
	}
	if len(output.Item) == 0 {
		log.Info().Msg("nothing in archive")
		return nil
	}

	log.Info().Msg("found archived tracking info")
	item := &ArchiveItem{}
	err = dynamodbattribute.UnmarshalMap(output.Item, item)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal item from dynamo attribute")
		return nil
	}
	return &item.TrackingInfo
}

func (a Archive) PutInfo(info TrackingInfo) {
	log := logger.With().
		Str("carrier", info.Carrier).
		Str("trackingNumber", info.TrackingNumber).
		Logger()

	item := ArchiveItem{
		Id:           getId(info.Carrier, info.TrackingNumber),
		TrackingInfo: info,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal item to dynamo attribute")
		return
	}

	_, err = a.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      av,
	})
	if err != nil {
		log.Error().Err(err).Msg("error putting to dynamo table")
	}
}

func getId(carrier string, trackingNumber string) string {
	return carrier + "#" + trackingNumber
}
