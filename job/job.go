package job

import (
	"context"
	"fmt"
	"github.com/seantcanavan/notification-step-machine/enum"
	"github.com/seantcanavan/notification-step-machine/metadata"
	"github.com/seantcanavan/notification-step-machine/util"
	"net/http"
	"time"
)

type CreateReq struct {
	ExpiresAt time.Time
	From      string
	Template  string
	To        string
	Type      enum.Type
	Variables map[string]interface{}
}

type Job struct {
	Created       time.Time
	EmailMetadata *metadata.Email
	From          string
	ID            string
	SMSMetadata   *metadata.SMS
	SnailMetadata *metadata.Snail
	Status        enum.Status
	Template      string
	To            string
	Type          enum.Type
	Updated       time.Time
	Variables     map[string]interface{}
}

func Create(ctx context.Context, cReq *CreateReq) (*Job, int, error) {
	cReq, httpStatus, err := validateCreateReq(cReq)
	if err != nil {
		return nil, httpStatus, err
	}
}

func Get(ctx context.Context, id string) (*Job, int, error) {
	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("parameter id [%s] is required", id)
	}
}

func GenerateRandom() *CreateReq {
	return &CreateReq{
		From:      util.GenerateRandomString(15),
		Template:  util.GenerateRandomString(15),
		To:        util.GenerateRandomEmail(),
		Type:      enum.Email,
		Variables: metadata.GenerateRandomEmailVariables(),
	}
}
