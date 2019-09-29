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

type ArchiveDao struct {
	svc *dynamodb.DynamoDB
}

func NewArchiveDao() *ArchiveDao {
	sess, err := session.NewSession()
	if err != nil {
		println(err.Error())
		return nil
	}
	return &ArchiveDao{dynamodb.New(sess)}
}

const TableName = "archive"

func (dao ArchiveDao) GetArchivedInfo(carrier string, trackingNumber string) *TrackingInfo {
	output, err := dao.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {S: aws.String(getId(carrier, trackingNumber))},
		},
	})

	if err != nil {
		println(err.Error())
		return nil
	}
	if len(output.Item) == 0 {
		return nil
	}

	println("found archived tracking info")
	item := &ArchiveItem{}
	err = dynamodbattribute.UnmarshalMap(output.Item, item)
	if err != nil {
		println(err.Error())
		return nil
	}
	return &item.TrackingInfo
}

func (dao ArchiveDao) PutArchivedInfo(carrier string, trackingNumber string, info TrackingInfo) error {
	item := ArchiveItem{
		Id:           getId(carrier, trackingNumber),
		TrackingInfo: info,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = dao.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      av,
	})
	return err
}

func getId(carrier string, trackingNumber string) string {
	return carrier + "#" + trackingNumber
}
