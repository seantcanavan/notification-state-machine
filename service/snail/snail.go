package snail

import (
	"context"
	"fmt"
	"github.com/seantcanavan/notification-step-machine/enum"
	"github.com/seantcanavan/notification-step-machine/job"
	"github.com/seantcanavan/notification-step-machine/util"
	"net/http"
	"time"
)

type Address struct {
	City            string  // US city
	Formatted       string  // Google Maps JavaScript API formatted string
	Latitude        float32 // Google Maps JavaScript API Latitude value
	Longitude       float32 // Google Maps JavaScript API Longitude value
	NumberAndStreet string  // The physical number of the address and the street
	Plus            int     // The 4 digit 'Plus' code after the zip. E.G. {zip}-{plus}
	State           string  // The two digit state code
	Unit            string  // optional apartment, unit, suite, etc
	Zip             int     // The 5 digit Zip code
}

type Instance struct {
	Address Address
}

func Handle(ctx context.Context, jobInstance job.Instance) (int, error) {
	num := util.GenerateRandomNumber()
	uReq := &job.UpdateReq{
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
		_, updateStatus, updateErr := job.Update(ctx, uReq)
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

	_, updateStatus, updateErr := job.Update(ctx, uReq)
	return updateStatus, updateErr
}

func Nudge(_ context.Context) (int, error) {
	fmt.Println("running Nudge for snail.go")
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

func GenerateRandom() *job.Instance {
	now := time.Now()
	return &job.Instance{
		Created:   now,
		From:      util.GenerateRandomString(10),
		ID:        util.NewUUID(),
		Snail:     &Instance{Address: GenerateRandomAddress()},
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
