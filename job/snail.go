package job

import (
	"context"
	"fmt"
	"github.com/seantcanavan/notification-step-machine/enum"
	"github.com/seantcanavan/notification-step-machine/util"
	"net/http"
	"time"
)

type Address struct {
	City            string  `json:"city,omitempty" dynamodbav:"city,omitempty"`             // US city
	Formatted       string  `json:"formatted,omitempty" dynamodbav:"formatted,omitempty"`   // Google Maps JavaScript API formatted string
	Latitude        float32 `json:"latitude,omitempty" dynamodbav:"latitude,omitempty"`     // Google Maps JavaScript API Latitude value
	Longitude       float32 `json:"longitude,omitempty" dynamodbav:"longitude,omitempty"`   // Google Maps JavaScript API Longitude value
	NumberAndStreet string  `json:"numberAndStreet,omitempty" dynamodbav:"numberAndStreet"` // The physical number of the address and the street
	Plus            int     `json:"plus,omitempty" dynamodbav:"plus,omitempty"`             // The 4 digit 'Plus' code after the zip. E.G. {zip}-{plus}
	State           string  `json:"state,omitempty" dynamodbav:"state,omitempty"`           // The two digit state code
	Unit            string  `json:"unit,omitempty" dynamodbav:"unit,omitempty"`             // optional apartment, unit, suite, etc
	Zip             int     `json:"zip,omitempty" dynamodbav:"zip,omitempty"`               // The 5 digit Zip code
}

type Snail struct {
	Address Address `json:"address,omitempty" dynamodbav:"address,omitempty"`
}

func HandleSnail(ctx context.Context, jobInstance *Instance) (int, error) {
	num := util.GenerateRandomNumber()
	uReq := &UpdateReq{
		Email:          jobInstance.Email,
		ID:             jobInstance.ID,
		PreviousStatus: jobInstance.Status,
		SMS:            jobInstance.SMS,
		Snail:          jobInstance.Snail,
		Variables:      jobInstance.Variables,
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

func NudgeSnail(_ context.Context) (int, error) {
	fmt.Println("running NudgeSnail for snail.go")
	return 0, nil
}

func GenerateRandomAddress() Address {
	return Address{
		City:            util.GenerateRandomString(2),
		Formatted:       util.GenerateRandomString(25),
		Latitude:        util.GenerateRandomFloat(),
		Longitude:       util.GenerateRandomFloat(),
		NumberAndStreet: util.GenerateRandomString(10),
		Plus:            util.GenerateNumberWithLength(4),
		State:           util.GenerateRandomString(2),
		Unit:            util.GenerateRandomString(8),
		Zip:             util.GenerateNumberWithLength(5),
	}
}

func GenerateRandomSnail() *Instance {
	now := time.Now()
	return &Instance{
		Created:   now,
		From:      util.GenerateRandomString(10),
		ID:        util.NewUUID(),
		Snail:     &Snail{Address: GenerateRandomAddress()},
		Status:    enum.Created,
		Template:  util.GenerateRandomString(10),
		To:        util.GenerateRandomString(10),
		Type:      enum.Snail,
		Updated:   now,
		Variables: GenerateRandomSnailVariables(),
	}
}

func GenerateRandomSnailVariables() map[string]interface{} {
	return map[string]interface{}{
		"address":   GenerateRandomAddress(),
		"firstName": util.GenerateRandomString(10),
		"lastName":  util.GenerateRandomString(10),
		"offerCode": util.GenerateRandomString(5),
	}
}
