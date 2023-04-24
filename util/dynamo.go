package util

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/private/protocol"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"net/http"
)

func NewUUID() string {
	return protocol.GetIdempotencyToken()
}

func ParseGIO(gio *dynamodb.GetItemOutput, id string, out interface{}) (int, error) {
	if gio.Item == nil {
		return http.StatusNotFound, fmt.Errorf("could not find item with ID [%s]", id)
	}

	err := dynamodbattribute.UnmarshalMap(gio.Item, out)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("could not unmarshal item [%+v] into out [%+v]", gio.Item, out)
	}

	return http.StatusOK, nil
}

func ParseQO(qo *dynamodb.QueryOutput, id string, out interface{}) (int, error) {
	if qo.Items == nil {
		return http.StatusNotFound, fmt.Errorf("could not query item(s) with ID [%s]", id)
	}

	err := dynamodbattribute.UnmarshalListOfMaps(qo.Items, out)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("could not unmarshal items [%+v] into out [%+v]", qo.Items, out)
	}

	return http.StatusOK, nil
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
