package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/seantcanavan/notification-step-machine/database_job"
	"github.com/seantcanavan/notification-step-machine/database_ttl"
	"github.com/seantcanavan/notification-step-machine/job"
	"time"
)

func init() {
	time.Local = time.UTC // force GoLang to use UTC as the default time zone for all calculations

	database_ttl.Connect()
	database_job.Connect()
}

func main() {
	lambda.Start(lambdaHandler)
}

// changed type of event from: events.DynamoDBEvent to DynamoDBEvent (see below)
func lambdaHandler(ctx context.Context, event DynamoDBEvent) error {

	for _, record := range event.Records {
		var beforeCReq job.CreateReq

		change := record.Change
		oldImage := change.OldImage // now of type: map[string]*dynamodb.AttributeValue

		err1 := dynamodbattribute.UnmarshalMap(oldImage, &beforeCReq)
		if err1 != nil {
			return err1
		}

		thawed, httpStatus, err := job.Create(ctx, &beforeCReq)
		if err != nil {
			fmt.Println(fmt.Sprintf("got status [%d] and err [%+v] thawing beforeCReq [%+v]", httpStatus, err, thawed))
		}
	}

	return nil
}

type DynamoDBEvent struct {
	Records []DynamoDBEventRecord `json:"Records"`
}

type DynamoDBEventRecord struct {
	AWSRegion      string                       `json:"awsRegion"`
	Change         DynamoDBStreamRecord         `json:"dynamodb"`
	EventID        string                       `json:"eventID"`
	EventName      string                       `json:"eventName"`
	EventSource    string                       `json:"eventSource"`
	EventVersion   string                       `json:"eventVersion"`
	EventSourceArn string                       `json:"eventSourceARN"`
	UserIdentity   *events.DynamoDBUserIdentity `json:"userIdentity,omitempty"`
}

type DynamoDBStreamRecord struct {
	ApproximateCreationDateTime events.SecondsEpochTime `json:"ApproximateCreationDateTime,omitempty"`
	// changed to map[string]*dynamodb.AttributeValue
	Keys map[string]*dynamodb.AttributeValue `json:"Keys,omitempty"`
	// changed to map[string]*dynamodb.AttributeValue
	NewImage map[string]*dynamodb.AttributeValue `json:"NewImage,omitempty"`
	// changed to map[string]*dynamodb.AttributeValue
	OldImage       map[string]*dynamodb.AttributeValue `json:"OldImage,omitempty"`
	SequenceNumber string                              `json:"SequenceNumber"`
	SizeBytes      int64                               `json:"SizeBytes"`
	StreamViewType string                              `json:"StreamViewType"`
}

//func main() {
//	lambda.Start(handler)
//}
//
//func handler(_ context.Context, event events.DynamoDBEvent) error {
//	for _, record := range event.Records {
//		fmt.Println(fmt.Sprintf("record is [%+v]", record))
//		var beforeCReq job.CreateReq
//		var afterCReq job.CreateReq
//		oldErr := dynamodbattribute.UnmarshalMap(record.Change.OldImage, &beforeCReq)
//		if oldErr != nil {
//			fmt.Println(fmt.Sprintf("oldErr [%+v]", oldErr))
//		}
//
//		newErr := dynamodbattribute.UnmarshalMap(record.Change.NewImage, &afterCReq)
//		if newErr != nil {
//			fmt.Println(fmt.Sprintf("newErr [%+v]", newErr))
//		}
//
//		fmt.Println(fmt.Sprintf("beforeCReq [%+v]", beforeCReq))
//		fmt.Println(fmt.Sprintf("afterCReq [%+v]", afterCReq))
//
//	}
//	return nil
//}
