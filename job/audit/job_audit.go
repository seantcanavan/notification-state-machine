package audit

import (
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
