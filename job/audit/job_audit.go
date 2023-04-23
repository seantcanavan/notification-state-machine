package audit

import (
	"context"
	"github.com/seantcanavan/notification-step-machine/enum"
	"time"
)

type JobAudit struct {
	Created        time.Time
	ID             string
	JobID          string
	NextStatus     enum.Status
	Operation      enum.Operation
	PreviousStatus enum.Status
	Updated        time.Time
}

type CreateReq struct {
	JobID          string
	NextStatus     enum.Status
	Operation      enum.Operation
	PreviousStatus enum.Status
}

func Create(ctx context.Context, cReq *CreateReq) (*JobAudit, int, error) {

}

func Get(ctx context.Context, jobID string) ([]*JobAudit, int, error) {

}

func GenerateRandom() *CreateReq {
}
