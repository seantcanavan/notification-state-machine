package queue

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/seantcanavan/lambda_jwt_router/lambda_util"
	"log"
	"os"
	"strings"
	"time"
)

type SQSQueue struct {
	client  *sqs.SQS
	session *session.Session
	URL     string
}

func (sqsQueue *SQSQueue) CreateSession() {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String("us-east-2"),
			Credentials: credentials.NewEnvCredentials(),
		},
	})

	if err != nil {
		log.Fatalf("error opening SQS session %+v", err)
	}

	sqsQueue.session = sess
	sqsQueue.URL = os.Getenv("ASYNC_SQS_URL")

	if sqsQueue.URL == "" {
		log.Fatalf("cannot proceed without ASYNC_SQS_URL environment variable")
	}

	sqsQueue.client = sqs.New(sqsQueue.session)
}

func (sqsQueue *SQSQueue) SendMessage(message *Message) error {
	deduplicationID := generateDeduplicationID(message)

	jsonBytes, _ := json.Marshal(message)
	messageBody := string(jsonBytes)
	_, err := sqsQueue.client.SendMessage(&sqs.SendMessageInput{
		QueueUrl:               &sqsQueue.URL,
		MessageBody:            aws.String(messageBody),
		MessageGroupId:         aws.String(lambda_util.GenerateRandomString(24)),
		MessageDeduplicationId: aws.String(deduplicationID),
	})

	return err
}

func (sqsQueue *SQSQueue) DeleteMessage(receiptHandle string) error {
	_, err := sqsQueue.client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(sqsQueue.URL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}

func generateDeduplicationID(message *Message) string {
	var sb strings.Builder
	sb.WriteString("_type_")
	sb.WriteString(message.Type.String())
	sb.WriteString("_resource_")
	sb.WriteString(message.Resource)
	sb.WriteString("_id_")
	sb.WriteString(message.ID)
	sb.WriteString("_at_")
	sb.WriteString(time.Now().Format(time.RFC3339Nano))
	return sb.String()
}

func (sqsQueue *SQSQueue) GetMessages() map[string][]*Message {
	return nil
}
