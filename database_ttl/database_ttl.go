package database_ttl

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"os"
)

var Client *dynamodb.DynamoDB
var TableName *string

func Connect() {
	fmt.Println("attempting to connect to ttl database")
	region := os.Getenv("REGION_AWS")
	stage := os.Getenv("STAGE")

	// use staging for now in dev for simplicity's sake
	if stage == "development" {
		TableName = aws.String(os.Getenv("APP_NAME") + "-" + "jobsTTL3" + "-" + "staging")
	} else {
		TableName = aws.String(os.Getenv("APP_NAME") + "-" + "jobsTTL3" + "-" + os.Getenv("STAGE"))
	}

	mySession := session.Must(session.NewSession())
	if region == "" {
		log.Fatal("REGION_AWS is not set")
	}
	Client = dynamodb.New(mySession, aws.NewConfig().WithRegion(region))

	fmt.Println("successfully connected to ttl database")
}

func Disconnect() {
	fmt.Println("attempting to disconnect from ttl database")
	// any cleanup code required here? probably not because dynamo is awesome
	fmt.Println("successfully disconnected from ttl database")
}
