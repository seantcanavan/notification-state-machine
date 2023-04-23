package job

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/seantcanavan/notification-step-machine/database_job"
	"github.com/seantcanavan/notification-step-machine/enum"
	"github.com/seantcanavan/notification-step-machine/metadata"
	"github.com/seantcanavan/notification-step-machine/util"
	"net/http"
	"time"
)

type CreateReq struct {
	ExpiresAt time.Time              `json:"expiresAt,omitempty"`
	From      string                 `json:"from,omitempty"`
	Template  string                 `json:"template,omitempty"`
	To        string                 `json:"to,omitempty"`
	Type      enum.Type              `json:"type,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type Job struct {
	Created       time.Time
	EmailMetadata *metadata.Email
	From          string
	ID            string
	SMSMetadata   *metadata.SMS
	SnailMetadata *metadata.Snail
	Status        enum.Status
	Template      string
	To            string
	Type          enum.Type
	Updated       time.Time
	Variables     map[string]interface{}
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

	_, err = database_job.Client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		ConditionExpression:         nil,
		ConditionalOperator:         nil,
		Expected:                    nil,
		ExpressionAttributeNames:    nil,
		ExpressionAttributeValues:   nil,
		Item:                        nil,
		ReturnConsumedCapacity:      nil,
		ReturnItemCollectionMetrics: nil,
		ReturnValues:                nil,
		TableName:                   nil,
	})

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
