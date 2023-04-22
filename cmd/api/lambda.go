package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/seantcanavan/lambda_jwt_router/lambda_router"
	"github.com/seantcanavan/notification-step-machine/job"
	"github.com/seantcanavan/notification-step-machine/service"
	"net/http"
)

func CreateLambda(ctx context.Context, lambdaReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var cReq service.CreateReq
	err := lambda_router.UnmarshalReq(lambdaReq, true, &cReq)
	if err != nil {
		return lambda_router.StatusAndErrorRes(http.StatusInternalServerError, err)
	}

	job, httpStatus, err := job.Create(ctx, &cReq)
	if err != nil {
		return lambda_router.StatusAndErrorRes(httpStatus, err)
	}

	return lambda_router.SuccessRes(job)
}

func GetLambda(ctx context.Context, lambdaReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	job, httpStatus, err := job.Get(ctx, lambdaReq.PathParameters["id"])
	if err != nil {
		return lambda_router.StatusAndErrorRes(httpStatus, err)
	}

	return lambda_router.SuccessRes(job)
}
