package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/seantcanavan/error_group"
	"github.com/seantcanavan/notification-step-machine/database_job"
	"github.com/seantcanavan/notification-step-machine/database_ttl"
	"github.com/seantcanavan/notification-step-machine/enum"
	"github.com/seantcanavan/notification-step-machine/job"
	"github.com/seantcanavan/notification-step-machine/util"
	"net/http"
	"sync"
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
	randomID := util.GenerateRandomString(10)
	fmt.Println(fmt.Sprintf("RandomID: [%s] delegator.lambdaHandler invoked", randomID))
	esg := error_group.NewErrorStatusGroup()

	var wg sync.WaitGroup
	wg.Add(len(event.Records))

	for _, currentRecord := range event.Records {
		go func(record util.DynamoDBEventRecord) {
			defer wg.Done()
			change := record.Change
			newImage := change.NewImage

			var jobInstance job.Instance
			unmarshalErr := dynamodbattribute.UnmarshalMap(newImage, &jobInstance)
			if unmarshalErr != nil {
				esg.AddStatusAndError(http.StatusInternalServerError, unmarshalErr)
			}

			fmt.Println(fmt.Sprintf("RandomID: [%s] ID: [%s] jobInstance [%+v]", randomID, jobInstance.ID, jobInstance))

			if jobInstance.Type == "" {
				return
			}

			fmt.Println(fmt.Sprintf("RandomID: [%s] ID: [%s] StreamViewType type is [%s]", randomID, jobInstance.ID, record.Change.StreamViewType))
			fmt.Println(fmt.Sprintf("RandomID: [%s] ID: [%s] delegator.go jobInstance is [%+v]", randomID, jobInstance.ID, jobInstance))
			esg.AddStatusAndError(delegate(ctx, randomID, &jobInstance))

		}(currentRecord)
	}

	wg.Wait()

	fmt.Println(fmt.Sprintf("RandomID: [%s] delegator.go error [%+v]", randomID, esg.ToError()))

	return esg.ToError()
}

func delegate(ctx context.Context, randomID string, jobInstance *job.Instance) (int, error) {
	if jobInstance.Type == enum.SMS {
		return job.HandleSMS(ctx, jobInstance)
	} else if jobInstance.Type == enum.Email {
		return job.HandleEmail(ctx, jobInstance)
	} else if jobInstance.Type == enum.Snail {
		return job.HandleSnail(ctx, jobInstance)
	} else {
		return http.StatusBadRequest, fmt.Errorf("RandomID: [%s] ID: [%s] unknown jobInstance.Type [%+v]", randomID, jobInstance.ID, jobInstance.Type)
	}
}
