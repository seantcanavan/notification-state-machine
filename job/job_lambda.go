package job

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/seantcanavan/lambda_jwt_router/lambda_router"
	"net/http"
)

func CreateLambda(ctx context.Context, lambdaReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var cReq CreateReq
	err := lambda_router.UnmarshalReq(lambdaReq, true, &cReq)
	if err != nil {
		return lambda_router.StatusAndErrorRes(http.StatusInternalServerError, err)
	}

	job, httpStatus, err := Create(ctx, &cReq)
	if err != nil {
		return lambda_router.StatusAndErrorRes(httpStatus, err)
	}

	return lambda_router.SuccessRes(job)
}

func GetLambda(ctx context.Context, lambdaReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	job, httpStatus, err := Get(ctx, lambdaReq.PathParameters["id"])
	if err != nil {
		return lambda_router.StatusAndErrorRes(httpStatus, err)
	}

	return lambda_router.SuccessRes(job)
}
