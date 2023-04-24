package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/seantcanavan/notification-step-machine/database_job"
	"github.com/seantcanavan/notification-step-machine/database_ttl"
	"github.com/seantcanavan/notification-step-machine/enum"
	"github.com/seantcanavan/notification-step-machine/job"
	"net/http"
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

type ServerlessInput struct {
	Type enum.Type `json:"type"`
}

func lambdaHandler(ctx context.Context, input ServerlessInput) (int, error) {
	fmt.Println(fmt.Sprintf("cleanup.lambdaHandler invoked"))
	if input.Type == "" {
		return http.StatusBadRequest, fmt.Errorf("empty Type [%+v]", input.Type)
	} else if input.Type == enum.Email {
		return job.NudgeEmail(ctx)
	} else if input.Type == enum.SMS {
		return job.NudgeSMS(ctx)
	} else if input.Type == enum.Snail {
		return job.NudgeSnail(ctx)
	} else {
		return http.StatusBadRequest, fmt.Errorf("invalid Type [%+v]", input.Type)
	}
}
