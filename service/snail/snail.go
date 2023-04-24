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
