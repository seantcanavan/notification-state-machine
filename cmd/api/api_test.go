package main

import (
	"github.com/joho/godotenv"
	"github.com/seantcanavan/notification-step-machine/database_job"
	"github.com/seantcanavan/notification-step-machine/database_ttl"
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tearDown()
}

func setup() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Unable to load .env file: %s", err)
	}

	database_job.Connect()
	database_ttl.Connect()
	time.Local = time.UTC
}

func tearDown() {
	database_job.Disconnect()
	database_ttl.Disconnect()
}
