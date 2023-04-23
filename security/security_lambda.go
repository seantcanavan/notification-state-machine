package security

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/seantcanavan/lambda_jwt_router/lambda_router"
	"net/http"
)

func LoginLambda(_ context.Context, lambdaReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("i'm here")
	var loginReq LoginReq
	err := lambda_router.UnmarshalReq(lambdaReq, true, &loginReq)
	if err != nil {
		return lambda_router.StatusAndErrorRes(http.StatusInternalServerError, err)
	}

	fmt.Println(fmt.Sprintf("loginReq [%+v]", loginReq))

	loginRes, httpStatus, err := Login(&loginReq)
	if err != nil {
		return lambda_router.StatusAndErrorRes(httpStatus, err)
	}

	return lambda_router.SuccessRes(loginRes)
}
