package audit

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/seantcanavan/notification-step-machine/database_job"
	"github.com/seantcanavan/notification-step-machine/enum"
	"github.com/seantcanavan/notification-step-machine/util"
	"net/http"
	"time"
)

type JobAudit struct {
	Created        time.Time      `json:"created,omitempty" dynamodbav:"created,omitempty"`
	ID             string         `json:"id,omitempty" dynamodbav:"id,omitempty"`
	JobID          string         `json:"jobId,omitempty" dynamodbav:"jobId,omitempty"`
	NextStatus     enum.Status    `json:"nextStatus,omitempty" dynamodbav:"nextStatus"`
	Operation      enum.Operation `json:"operation,omitempty" dynamodbav:"operation"`
	PreviousStatus enum.Status    `json:"previousStatus,omitempty" dynamodbav:"previousStatus"`
	Updated        time.Time      `json:"updted,omitempty" dynamodbav:"updated"`
}

type CreateReq struct {
	JobID          string
	NextStatus     enum.Status
	Operation      enum.Operation
	PreviousStatus enum.Status
}

func Create(ctx context.Context, cReq *CreateReq) (*JobAudit, int, error) {
	cReq, httpStatus, err := validateCreateReq(cReq)
	if err != nil {
		return nil, httpStatus, err
	}

	now := time.Now()

	jobAudit := &JobAudit{
		Created:        now,
		ID:             util.NewUUID(),
		JobID:          cReq.JobID,
		NextStatus:     cReq.NextStatus,
		Operation:      cReq.Operation,
		PreviousStatus: cReq.PreviousStatus,
		Updated:        now,
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

	return jobAudit, http.StatusOK, nil
}

func Get(ctx context.Context, jobID string) ([]*JobAudit, int, error) {
	if jobID == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("parameter jobID [%s] is required", jobID)
	}

	qo, err := database_job.Client.QueryWithContext(ctx, &dynamodb.QueryInput{
		AttributesToGet:           nil,
		ConditionalOperator:       nil,
		ConsistentRead:            nil,
		ExclusiveStartKey:         nil,
		ExpressionAttributeNames:  nil,
		ExpressionAttributeValues: nil,
		FilterExpression:          nil,
		IndexName:                 nil,
		KeyConditionExpression:    nil,
		KeyConditions:             nil,
		Limit:                     nil,
		ProjectionExpression:      nil,
		QueryFilter:               nil,
		ReturnConsumedCapacity:    nil,
		ScanIndexForward:          nil,
		Select:                    nil,
		TableName:                 nil,
	})
	if err != nil {
		return nil, util.DecodeAWSErr(err), err
	}

	var jobAudit []*JobAudit
	httpStatus, err := util.ParseQO(qo, jobID, &jobAudit)
	return jobAudit, httpStatus, err
}

func GenerateRandom() *CreateReq {
	return &CreateReq{
		JobID:          util.NewUUID(),
		NextStatus:     enum.Sent,
		Operation:      enum.Update,
		PreviousStatus: enum.Processing,
	}
}
