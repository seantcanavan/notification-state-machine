package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/seantcanavan/lambda_jwt_router/lambda_jwt"
	"github.com/seantcanavan/lambda_jwt_router/lambda_router"
	"github.com/seantcanavan/notification-step-machine/database_job"
	"github.com/seantcanavan/notification-step-machine/database_ttl"
	"github.com/seantcanavan/notification-step-machine/job"
	"github.com/seantcanavan/notification-step-machine/job/audit"
	"github.com/seantcanavan/notification-step-machine/queue"
	"log"
	"net/http"
	"os"
	"time"
)

var router *lambda_router.Router

func init() {
	time.Local = time.UTC // force GoLang to use UTC as the default time zone for all calculations

	database_ttl.Connect()
	database_job.Connect()

	environment := os.Getenv("STAGE")
	if environment == "staging" || environment == "production" {
		queue.NewClient(true)
	} else {
		queue.NewClient(false)
	}

	router = lambda_router.NewRouter("/api", lambda_jwt.LogRequestMW)

	// Notifications endpoints
	router.Route("GET", "/notifications/:id", job.GetLambda, lambda_jwt.DecodeStandard)
	router.Route("POST", "/notifications", job.CreateLambda, lambda_jwt.DecodeStandard)

	// Audit endpoints
	router.Route("GET", "/audit/:jobId", audit.GetLambda, lambda_jwt.DecodeStandard)
	router.Route("POST", "/audit", audit.CreateLambda, lambda_jwt.DecodeStandard)
}

func main() {
	environment := os.Getenv("STAGE")
	if environment == "staging" || environment == "production" {
		lambda.Start(router.Handler)
	} else {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		fmt.Println(fmt.Sprintf("Ready to listen and serve on port %s", port))
		err := http.ListenAndServe(":"+port, http.HandlerFunc(router.ServeHTTP))
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
}
