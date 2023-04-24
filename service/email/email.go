package email

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
	SesID string
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
		Email:     &Instance{SesID: util.GenerateRandomString(10)},
		From:      util.GenerateRandomString(10),
		ID:        util.NewUUID(),
		Status:    enum.Created,
		Template:  util.GenerateRandomString(10),
		To:        util.GenerateRandomString(10),
		Type:      enum.Email,
		Updated:   now,
		Variables: GenerateRandomEmailVariables(),
	}
}

func GenerateRandomEmailVariables() map[string]interface{} {
	return map[string]interface{}{
		"firstName": util.GenerateRandomString(10),
		"footer":    util.GenerateRandomString(10),
		"header":    util.GenerateRandomString(10),
		"lastName":  util.GenerateRandomString(10),
	}
}
