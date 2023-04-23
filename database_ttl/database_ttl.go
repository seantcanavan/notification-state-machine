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

func Connect() {
	fmt.Println("attempting to connect to ttl database")

	mySession := session.Must(session.NewSession())
	region := os.Getenv("REGION_AWS")
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
