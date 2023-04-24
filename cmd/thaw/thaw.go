package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/seantcanavan/notification-step-machine/database_job"
	"github.com/seantcanavan/notification-step-machine/database_ttl"
	"github.com/seantcanavan/notification-step-machine/job"
	"github.com/seantcanavan/notification-step-machine/util"
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
func lambdaHandler(ctx context.Context, event util.DynamoDBEvent) error {
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
