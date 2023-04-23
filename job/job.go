package job

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/seantcanavan/notification-step-machine/database_job"
	"github.com/seantcanavan/notification-step-machine/database_ttl"
	"github.com/seantcanavan/notification-step-machine/enum"
	"github.com/seantcanavan/notification-step-machine/metadata"
	"github.com/seantcanavan/notification-step-machine/util"
	"net/http"
	"time"
)

type CreateReq struct {
	ExpiresAt *time.Time             `json:"expiresAt,omitempty"`
	From      string                 `json:"from,omitempty"`
	Template  string                 `json:"template,omitempty"`
	To        string                 `json:"to,omitempty"`
	Type      enum.Type              `json:"type,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type Job struct {
	Created       time.Time              `json:"created,omitempty" dynamodbav:"created,omitempty"`
	EmailMetadata *metadata.Email        `json:"emailMetadata,omitempty" dynamodbav:"emailMetadata,omitempty"`
	From          string                 `json:"from,omitempty" dynamodbav:"from,omitempty"`
	ID            string                 `json:"id,omitempty" dynamodbav:"id,omitempty"`
	SMSMetadata   *metadata.SMS          `json:"smsMetadata,omitempty" dynamodbav:"smsMetadata,omitempty"`
	SnailMetadata *metadata.Snail        `json:"snailMetadata,omitempty" dynamodbav:"snailMetadata,omitempty"`
	Status        enum.Status            `json:"status,omitempty" dynamodbav:"status,omitempty"`
	Template      string                 `json:"template,omitempty" dynamodbav:"template,omitempty"`
	To            string                 `json:"to,omitempty" dynamodbav:"to,omitempty"`
	Type          enum.Type              `json:"type,omitempty" dynamodbav:"type,omitempty"`
	Updated       time.Time              `json:"updated,omitempty" dynamodbav:"updated,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty" dynamodbav:"variables,omitempty"`
}

func Create(ctx context.Context, cReq *CreateReq) (*Job, int, error) {
	cReq, httpStatus, err := validateCreateReq(cReq)
	if err != nil {
		return nil, httpStatus, err
	}

	now := time.Now()

	job := &Job{
		Created:   now,
		From:      cReq.From,
		ID:        util.NewUUID(),
		Status:    enum.Created,
		Template:  cReq.Template,
		To:        cReq.To,
		Type:      cReq.Type,
		Updated:   now,
		Variables: cReq.Variables,
	}

	marshalled, marshalErr := dynamodbattribute.MarshalMap(job)
	if marshalErr != nil {
		return nil, http.StatusInternalServerError, marshalErr
	}

	fmt.Println(fmt.Sprintf("cReq.ExpiresAt is [%+v]", *cReq.ExpiresAt))

	if cReq.ExpiresAt != nil {
		_, err = database_job.Client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
			ConditionExpression: aws.String("attribute_not_exists(id)"),
			//ExpressionAttributeNames:    nil,
			//ExpressionAttributeValues:   nil,
			Item:      marshalled,
			TableName: database_ttl.TableName,
		})
	} else {
		_, err = database_job.Client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
			ConditionExpression: aws.String("attribute_not_exists(id)"),
			//ExpressionAttributeNames:    nil,
			//ExpressionAttributeValues:   nil,
			Item:      marshalled,
			TableName: database_job.TableName,
		})
	}

	if err != nil {
		return nil, util.DecodeAWSErr(err), err
	}

	return job, http.StatusOK, nil
}

func Get(ctx context.Context, id string) (*Job, int, error) {
	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("parameter id [%s] is required", id)
	}

	gio, err := database_job.Client.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		AttributesToGet:          nil,
		ConsistentRead:           nil,
		ExpressionAttributeNames: nil,
		Key:                      nil,
		ProjectionExpression:     nil,
		ReturnConsumedCapacity:   nil,
		TableName:                nil,
	})
	if err != nil {
		return nil, util.DecodeAWSErr(err), err
	}

	var job Job
	httpStatus, err := util.ParseGIO(gio, id, &job)
	return &job, httpStatus, err
}

func GenerateRandom() *CreateReq {
	return &CreateReq{
		From:      util.GenerateRandomString(15),
		Template:  util.GenerateRandomString(15),
		To:        util.GenerateRandomEmail(),
		Type:      enum.Email,
		Variables: metadata.GenerateRandomEmailVariables(),
	}
}
