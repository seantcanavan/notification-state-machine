package thaw

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/seantcanavan/notification-step-machine/database_job"
	"github.com/seantcanavan/notification-step-machine/database_ttl"
	"time"
)

func init() {
	time.Local = time.UTC // force GoLang to use UTC as the default time zone for all calculations

	database_ttl.Connect()
	database_job.Connect()
}

func main() {
	lambda.Start(handler)
}

func handler(_ context.Context, event events.DynamoDBEvent) error {
	fmt.Println(fmt.Sprintf("event is [%+v]", event))
	return nil
}
