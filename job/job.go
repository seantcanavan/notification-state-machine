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
	"github.com/seantcanavan/notification-step-machine/job/audit"
	"github.com/seantcanavan/notification-step-machine/util"
	"log"
	"net/http"
	"time"
)

type Instance struct {
	Created   time.Time              `json:"created,omitempty" dynamodbav:"created,omitempty"`
	Email     *Email                 `json:"email,omitempty" dynamodbav:"email,omitempty"`
	From      string                 `json:"from,omitempty" dynamodbav:"from,omitempty"`
	ID        string                 `json:"id,omitempty" dynamodbav:"id,omitempty"`
	SMS       *SMS                   `json:"sms,omitempty" dynamodbav:"sms,omitempty"`
	Snail     *Snail                 `json:"snail,omitempty" dynamodbav:"snail,omitempty"`
	Status    enum.Status            `json:"status,omitempty" dynamodbav:"status,omitempty"`
	Template  string                 `json:"template,omitempty" dynamodbav:"template,omitempty"`
	To        string                 `json:"to,omitempty" dynamodbav:"to,omitempty"`
	Type      enum.Type              `json:"type,omitempty" dynamodbav:"type,omitempty"`
	Updated   time.Time              `json:"updated,omitempty" dynamodbav:"updated,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty" dynamodbav:"variables,omitempty"`
}

type CreateReq struct {
	Created   *time.Time             `json:"created,omitempty" dynamodbav:"created,omitempty"`
	ExpiresAt *time.Time             `json:"expiresAt,omitempty" dynamodbav:"expiresAt,omitempty"`
	From      string                 `json:"from,omitempty" dynamodbav:"from,omitempty"`
	ID        *string                `json:"id,omitempty" dynamodbav:"id,omitempty"`
	Template  string                 `json:"template,omitempty" dynamodbav:"template,omitempty"`
	To        string                 `json:"to,omitempty" dynamodbav:"to,omitempty"`
	Type      enum.Type              `json:"type,omitempty" dynamodbav:"type,omitempty"`
	Updated   *time.Time             `json:"updated,omitempty" dynamodbav:"updated,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty" dynamodbav:"variables,omitempty"`
	TTL       int64                  `json:"-" dynamodbav:"ttl,omitempty"`
}

func Create(ctx context.Context, cReq *CreateReq) (*Instance, int, error) {
	cReq, httpStatus, err := validateCreateReq(cReq)
	if err != nil {
		return nil, httpStatus, err
	}

	now := time.Now()

	job := &Instance{
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

	_, err = database_job.Client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(id)"),
		Item:                marshalled,
		TableName:           database_job.TableName,
	})

	if err != nil {
		return nil, util.DecodeAWSErr(err), err
	}

	audit.CreateSilent(ctx, &audit.CreateReq{
		JobID:          job.ID,
		NextStatus:     enum.Created,
		Operation:      enum.Create,
		PreviousStatus: enum.NA,
	})

	return job, http.StatusOK, nil
}

func Freeze(ctx context.Context, cReq *CreateReq) (*CreateReq, int, error) {
	cReq, httpStatus, err := validateCreateReq(cReq)
	if err != nil {
		return nil, httpStatus, err
	}

	now := time.Now()

	if cReq.ExpiresAt.Before(now) {
		return nil, http.StatusBadRequest, fmt.Errorf("cannot expire [%+v] before now [%+v]", cReq.ExpiresAt, now)
	}

	cReq.Created = &now
	cReq.ID = aws.String(util.NewUUID())
	cReq.Updated = &now
	cReq.TTL = cReq.ExpiresAt.Unix()

	marshalled, marshalErr := dynamodbattribute.MarshalMap(cReq)
	if marshalErr != nil {
		return nil, http.StatusInternalServerError, marshalErr
	}

	_, err = database_job.Client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(id)"),
		Item:                marshalled,
		TableName:           database_ttl.TableName,
	})

	if err != nil {
		return nil, util.DecodeAWSErr(err), err
	}

	return cReq, http.StatusOK, nil
}

func Get(ctx context.Context, id string) (*Instance, int, error) {
	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("parameter id [%s] is required", id)
	}

	gio, err := database_job.Client.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		TableName: database_job.TableName,
	})
	if err != nil {
		return nil, util.DecodeAWSErr(err), err
	}

	var job Instance
	httpStatus, err := util.ParseGIO(gio, id, &job)
	return &job, httpStatus, err
}

type UpdateReq struct {
	Email          *Email                 `json:"email,omitempty"`
	ID             string                 `json:"id,omitempty"`
	SMS            *SMS                   `json:"sms,omitempty"`
	Snail          *Snail                 `json:"snail,omitempty"`
	Variables      map[string]interface{} `json:"variables,omitempty"`
	NextStatus     enum.Status            `json:"nextStatus,omitempty" `
	PreviousStatus enum.Status            `json:"previousStatus,omitempty"`
}

func Update(ctx context.Context, uReq *UpdateReq) (*Instance, int, error) {
	uReq, httpStatus, err := validateUpdateReq(uReq)
	if err != nil {
		return nil, httpStatus, err
	}

	variables, variablesErr := dynamodbattribute.MarshalMap(uReq.Variables)
	if variablesErr != nil {
		return nil, http.StatusInternalServerError, variablesErr
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(uReq.NextStatus.String()),
			},
			":variables": {
				M: variables,
			},
		},
		TableName: database_job.TableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(uReq.ID),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET status = :status, variables = :variables"),
	}

	if uReq.Email != nil {
		emailVal, emailErr := dynamodbattribute.MarshalMap(uReq.Email)
		if emailErr != nil {
			return nil, http.StatusInternalServerError, emailErr
		}
		input.ExpressionAttributeValues[":email"] = &dynamodb.AttributeValue{M: emailVal}
		current := *input.UpdateExpression
		input.UpdateExpression = aws.String(current + ", email = :email")
	}

	if uReq.SMS != nil {
		smsVal, smsErr := dynamodbattribute.MarshalMap(uReq.SMS)
		if smsErr != nil {
			return nil, http.StatusInternalServerError, smsErr
		}
		input.ExpressionAttributeValues[":sms"] = &dynamodb.AttributeValue{M: smsVal}
		current := *input.UpdateExpression
		input.UpdateExpression = aws.String(current + ", sms = :sms")
	}

	if uReq.Snail != nil {
		snailVal, snailErr := dynamodbattribute.MarshalMap(uReq.Snail)
		if snailErr != nil {
			return nil, http.StatusInternalServerError, snailErr
		}
		input.ExpressionAttributeValues[":snail"] = &dynamodb.AttributeValue{M: snailVal}
		current := *input.UpdateExpression
		input.UpdateExpression = aws.String(current + ", snail = :snail")
	}

	uio, err := database_job.Client.UpdateItemWithContext(ctx, input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	}

	audit.CreateSilent(ctx, &audit.CreateReq{
		JobID:          uReq.ID,
		NextStatus:     uReq.NextStatus,
		Operation:      enum.Update,
		PreviousStatus: uReq.PreviousStatus,
	})

	var jobInstance Instance
	httpStatus, err = util.ParseUIO(uio, uReq.ID, &jobInstance)
	if err != nil {
		return nil, httpStatus, err
	}

	return &jobInstance, http.StatusOK, nil
}

func GenerateRandom() *CreateReq {
	return &CreateReq{
		From:      util.GenerateRandomString(15),
		Template:  util.GenerateRandomString(15),
		To:        util.GenerateRandomEmail(),
		Type:      enum.Email,
		Variables: GenerateRandomEmailVariables(),
	}
}
