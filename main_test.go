package main

import (
	"fmt"
	"github.com/jgroeneveld/trial/assert"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tearDown()
}

func setup() {
	fmt.Println("setting up main_test.go")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Unable to load .env file: %s", err)
	}

	time.Local = time.UTC
}

func tearDown() {
	fmt.Println("tearing down main_test.go")
}

func TestENV(t *testing.T) {
	stage := os.Getenv("STAGE")
	assert.Equal(t, "development", stage)
}

func TestXYZ(t *testing.T) {
	res := XYZ()
	assert.Equal(t, "abc", res)
}
