package metadata

import "github.com/seantcanavan/notification-step-machine/util"

type Email struct {
	SESJobID string
}

type SMS struct {
}

type Snail struct {
	Address Address
}

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

func GenerateRandomEmailVariables() map[string]interface{} {
	return map[string]interface{}{
		"firstName": util.GenerateRandomString(10),
		"footer":    util.GenerateRandomString(10),
		"header":    util.GenerateRandomString(10),
		"lastName":  util.GenerateRandomString(10),
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

func GenerateRandomSnailVariables() map[string]interface{} {
	return map[string]interface{}{
		"address":   GenerateRandomAddress(),
		"firstName": util.GenerateRandomString(10),
		"lastName":  util.GenerateRandomString(10),
		"offerCode": util.GenerateRandomString(5),
	}
}

func GenerateRandomAddress() *Address {
	return &Address{
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
