package database_job

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"os"
)

var Client *dynamodb.DynamoDB
var TableName *string

func Connect() {
	region := os.Getenv("REGION_AWS")
	stage := os.Getenv("STAGE")

	// use staging for now in dev for simplicity's sake
	if stage == "development" {
		TableName = aws.String(os.Getenv("APP_NAME") + "-" + "jobs" + "-" + "staging")
	} else {
		TableName = aws.String(os.Getenv("APP_NAME") + "-" + "jobs" + "-" + os.Getenv("STAGE"))
	}

	mySession := session.Must(session.NewSession())
	if region == "" {
		log.Fatal("REGION_AWS is not set")
	}
	Client = dynamodb.New(mySession, aws.NewConfig().WithRegion(region))
}

func Disconnect() {
}
