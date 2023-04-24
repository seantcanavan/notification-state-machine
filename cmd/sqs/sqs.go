package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func lambdaHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	fmt.Println(fmt.Sprintf("sqs.lambdaHandler invoked"))
	for _, record := range sqsEvent.Records {
		var cReq job.CreateReq
		err := json.Unmarshal([]byte(record.Body), &cReq)
		if err != nil {
			fmt.Println(fmt.Sprintf("error unmarshalling record [%+v]", record))
			continue
		}

		_, httpStatus, err := job.Create(ctx, &cReq)
		if err != nil {
			return fmt.Errorf("code [%d] error [%+v] creating job [%+v]", httpStatus, err, cReq)
		}
	}

	return nil
}
