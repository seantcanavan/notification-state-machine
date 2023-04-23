package audit

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/seantcanavan/notification-step-machine/database_job"
	"github.com/seantcanavan/notification-step-machine/enum"
	"github.com/seantcanavan/notification-step-machine/util"
	"net/http"
	"time"
)

type JobAudit struct {
	Created        time.Time      `json:"created,omitempty" dynamodbav:"created,omitempty"`
	GSI1RK1Created time.Time      `json:"-" dynamodbav:"gsi1-rk1,omitempty"`
	ID             string         `json:"id,omitempty" dynamodbav:"id,omitempty"`
	JobID          string         `json:"jobId,omitempty" dynamodbav:"gsi1,omitempty"`
	NextStatus     enum.Status    `json:"nextStatus,omitempty" dynamodbav:"nextStatus"`
	Operation      enum.Operation `json:"operation,omitempty" dynamodbav:"operation"`
	PreviousStatus enum.Status    `json:"previousStatus,omitempty" dynamodbav:"previousStatus"`
	Updated        time.Time      `json:"updated,omitempty" dynamodbav:"updated"`
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
		GSI1RK1Created: now,
		ID:             util.NewUUID(),
		JobID:          "AUD|" + cReq.JobID, // single table define. prefix with A| for audit entries
		NextStatus:     cReq.NextStatus,
		Operation:      cReq.Operation,
		PreviousStatus: cReq.PreviousStatus,
		Updated:        now,
	}

	marshalled, marshalErr := dynamodbattribute.MarshalMap(jobAudit)
	if marshalErr != nil {
		return nil, http.StatusInternalServerError, marshalErr
	}

	_, err = database_job.Client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(id)"),
		Item:                marshalled,
		TableName:           database_job.TableName,
	})

	if err != nil {
		return nil, util.DecodeAWSErr(err), err
	}

	return jobAudit, http.StatusOK, nil
}

func CreateSilent(ctx context.Context, cReq *CreateReq) {
	_, auditStatus, auditErr := Create(ctx, cReq)
	if auditErr != nil {
		fmt.Println(fmt.Sprintf("encountered status [%d] while creating audit [%+v] with error [%+v]", auditStatus, *cReq, auditErr))
	}
}

func Get(ctx context.Context, jobID string) ([]*JobAudit, int, error) {
	if jobID == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("parameter jobID [%s] is required", jobID)
	}

	keyCond := expression.Key("gsi1").Equal(expression.Value("AUD|" + jobID))
	builder := expression.NewBuilder().WithKeyCondition(keyCond)
	expr, exprErr := builder.Build()
	if exprErr != nil {
		return nil, http.StatusInternalServerError, exprErr
	}

	queryInput := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		IndexName:                 aws.String("gsi1"),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 database_job.TableName,
	}

	qo, queryErr := database_job.Client.QueryWithContext(ctx, queryInput)
	if queryErr != nil {
		return nil, util.DecodeAWSErr(queryErr), queryErr
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
