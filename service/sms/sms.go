package sms

import (
	"context"
	"fmt"
	"github.com/seantcanavan/notification-step-machine/enum"
	"github.com/seantcanavan/notification-step-machine/job"
	"github.com/seantcanavan/notification-step-machine/util"
	"net/http"
	"time"
)

type Instance struct {
	SnsID string
}

func Handle(ctx context.Context, job job.Instance) (int, error) {
	if job.Status == enum.Created {
		fmt.Println(fmt.Printf("job [%+v] has reached Created state\n", job))
	} else if job.Status == enum.Error {
		fmt.Println(fmt.Sprintf("job [%+v] has reached Error state", job))
	} else if job.Status == enum.Processing {
		fmt.Println(fmt.Sprintf("job [%+v] has reached Processing state", job))
	} else if job.Status == enum.Queued {
		fmt.Println(fmt.Sprintf("job [%+v] has reached Queued state", job))
	} else if job.Status == enum.Sent {
		fmt.Println(fmt.Sprintf("job [%+v] has reached Sent state", job))
	}

	return http.StatusInternalServerError, fmt.Errorf("unknown job.Stats [%+v]", job.Status)
}

func GenerateRandom() *job.Instance {
	now := time.Now()
	return &job.Instance{
		Created:   now,
		From:      util.GenerateRandomString(10),
		ID:        util.NewUUID(),
		SMS:       &Instance{SnsID: util.GenerateRandomString(10)},
		Status:    enum.Created,
		Template:  util.GenerateRandomString(10),
		To:        util.GenerateRandomString(10),
		Type:      enum.SMS,
		Updated:   now,
		Variables: GenerateRandomSMSVariables(),
	}
}

func GenerateRandomSMSVariables() map[string]interface{} {
	return map[string]interface{}{
		"amount":    util.GenerateRandomFloat(),
		"code":      util.GenerateRandomNumber(),
		"firstName": util.GenerateRandomString(10),
		"lastName":  util.GenerateRandomString(10),
	}
}
