package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/seantcanavan/error_group"
	"github.com/seantcanavan/notification-step-machine/database_job"
	"github.com/seantcanavan/notification-step-machine/database_ttl"
	"github.com/seantcanavan/notification-step-machine/enum"
	"net/http"
	"time"
)

func init() {
	time.Local = time.UTC // force GoLang to use UTC as the default time zone for all calculations

	database_ttl.Connect()
	database_job.Connect()
}

func main() {
	lambda.Start(Handler)
}

type ServerlessInput struct {
	Type enum.Type `json:"type"`
}

func Handler(ctx context.Context, input ServerlessInput) (int, error) {
	if input.Type == "" {
		return http.StatusBadRequest, fmt.Errorf("empty Type [%+v]", input.Type)
	} else if input.Type == enum.Email {

	} else if input.Type == enum.SMS {

	} else if input.Type == enum.Snail {

	} else {
		return http.StatusBadRequest, fmt.Errorf("invalid Type [%+v]", input.Type)
	}

	esg := error_group.NewErrorStatusGroup()

	return esg.ToStatusAndError()
}
