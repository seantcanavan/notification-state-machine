package job

import (
	"context"
	"fmt"
	"github.com/seantcanavan/notification-step-machine/enum"
	"github.com/seantcanavan/notification-step-machine/util"
	"net/http"
	"time"
)

var OldThreshold = time.Now().AddDate(0, 0, -1)

type SMS struct {
	SnsID string `json:"snsId,omitempty" dynamodbav:"snsId,omitempty"`
}

func HandleSMS(ctx context.Context, jobInstance *Instance) (int, error) {
	num := util.GenerateRandomNumber()
	uReq := &UpdateReq{
		ID:             jobInstance.ID,
		PreviousStatus: jobInstance.Status,
		SMS:            jobInstance.SMS,
		Variables:      jobInstance.Variables,
	}

	if uReq.SMS == nil {
		uReq.SMS = &SMS{SnsID: util.GenerateRandomString(10)}
	}

	if num < 3 { // stay in the same state and get 'nudged' later
		return http.StatusOK, nil
	} else if num < 4 { // move to the error state
		uReq.NextStatus = enum.Error
		_, updateStatus, updateErr := Update(ctx, uReq)
		return updateStatus, updateErr
	}

	if jobInstance.Status == enum.Created {
		fmt.Println(fmt.Printf("jobInstance [%+v] is in Created state\n", jobInstance))
		uReq.NextStatus = enum.Queued
	} else if jobInstance.Status == enum.Error {
		fmt.Println(fmt.Sprintf("jobInstance [%+v] is in Error state", jobInstance))
		return http.StatusOK, nil
	} else if jobInstance.Status == enum.Processing {
		fmt.Println(fmt.Sprintf("jobInstance [%+v] is in Processing state", jobInstance))
		uReq.NextStatus = enum.Sent
	} else if jobInstance.Status == enum.Queued {
		fmt.Println(fmt.Sprintf("jobInstance [%+v] is in Queued state", jobInstance))
		uReq.NextStatus = enum.Processing
	} else if jobInstance.Status == enum.Sent {
		fmt.Println(fmt.Sprintf("jobInstance [%+v] is in Sent state", jobInstance))
		return http.StatusOK, nil
	}

	_, updateStatus, updateErr := Update(ctx, uReq)
	return updateStatus, updateErr
}

func NudgeSMS(_ context.Context) (int, error) {
	fmt.Println("running NudgeSMS for sms.go")
	return 0, nil
}

func GenerateRandomSMS() *Instance {
	now := time.Now()
	return &Instance{
		Created:   now,
		From:      util.GenerateRandomString(10),
		ID:        util.NewUUID(),
		SMS:       &SMS{SnsID: util.GenerateRandomString(10)},
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
