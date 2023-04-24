package job

import (
	"context"
	"fmt"
	"github.com/seantcanavan/notification-step-machine/enum"
	"github.com/seantcanavan/notification-step-machine/util"
	"net/http"
	"time"
)

type Email struct {
	SesID string
}

func HandleEmail(ctx context.Context, jobInstance Instance) (int, error) {
	num := util.GenerateRandomNumber()
	uReq := &UpdateReq{
		Email:     jobInstance.Email,
		ID:        jobInstance.ID,
		SMS:       jobInstance.SMS,
		Snail:     jobInstance.Snail,
		Status:    jobInstance.Status,
		Variables: jobInstance.Variables,
	}

	if num < 3 { // stay in the same state and get 'nudged' later
		return http.StatusOK, nil
	} else if num < 4 { // move to the error state
		uReq.Status = enum.Error
		_, updateStatus, updateErr := Update(ctx, uReq)
		return updateStatus, updateErr
	}

	if jobInstance.Status == enum.Created {
		fmt.Println(fmt.Printf("jobInstance [%+v] is in Created state\n", jobInstance))
		uReq.Status = enum.Queued
	} else if jobInstance.Status == enum.Error {
		fmt.Println(fmt.Sprintf("jobInstance [%+v] is in Error state", jobInstance))
		return http.StatusOK, nil
	} else if jobInstance.Status == enum.Processing {
		fmt.Println(fmt.Sprintf("jobInstance [%+v] is in Processing state", jobInstance))
		uReq.Status = enum.Sent
	} else if jobInstance.Status == enum.Queued {
		fmt.Println(fmt.Sprintf("jobInstance [%+v] is in Queued state", jobInstance))
		uReq.Status = enum.Processing
	} else if jobInstance.Status == enum.Sent {
		fmt.Println(fmt.Sprintf("jobInstance [%+v] is in Sent state", jobInstance))
		return http.StatusOK, nil
	}

	_, updateStatus, updateErr := Update(ctx, uReq)
	return updateStatus, updateErr
}

func NudgeEmail(_ context.Context) (int, error) {
	fmt.Println("running NudgeEmail for email.go")
	return 0, nil
}

func GenerateRandomEmail() *Instance {
	now := time.Now()
	return &Instance{
		Created:   now,
		Email:     &Email{SesID: util.GenerateRandomString(10)},
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
