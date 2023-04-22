package main

import (
	lf_queue "github.com/seantcanavan/notification-step-machine/queue"
	"os"
)

func init() {
	environment := os.Getenv("STAGE")
	if environment == "staging" || environment == "production" {
		lf_queue.NewClient(true)
	} else {
		lf_queue.NewClient(false)
	}
}

func main() {

}
