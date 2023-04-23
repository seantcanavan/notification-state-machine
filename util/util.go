package util

import (
	"crypto/rand"
	"log"
	"math/big"
	"strconv"
	"strings"
)

func GenerateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}

func GenerateRandomNumber() int {
	const letters = "123456789"
	ret := make([]byte, 1)
	for i := 0; i < 1; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return 0
		}
		ret[i] = letters[num.Int64()]
	}

	byteToInt, _ := strconv.Atoi(string(ret))
	return byteToInt
}

func GenerateNumberWithLength(length int) int {
	digits := ""
	for x := 0; x < length; x++ {
		digits += strconv.FormatInt(int64(GenerateRandomNumber()), 10)
	}

	parseInt, parseErr := strconv.Atoi(digits)
	if parseErr != nil {
		log.Fatal("should not error here")
	}

	return parseInt
}

func GenerateRandomFloat() float32 {
	return float32(GenerateRandomNumber()) / float32(GenerateRandomNumber())
}

func GenerateRandomAlphaNumeric(length int) string {
	const alphaNumeric = "01234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphaNumeric))))
		if err != nil {
			return ""
		}
		ret[i] = alphaNumeric[num.Int64()]
	}

	return string(ret)
}

// GenerateRandomEmail will generate a subset of an email address to Sean
// so that when random emails are used in staging - those emails are not
// potentially sent to external users. This should be a very minor edge case
// but I wanted to account for it just in case.
func GenerateRandomEmail() string {
	sb := strings.Builder{}
	sb.WriteString("SEanTcaNavaN+") // caps here help us test LowerCase() algorithms
	sb.WriteString(GenerateRandomString(12))
	sb.WriteString("@gmail.com")
	return sb.String()
}
