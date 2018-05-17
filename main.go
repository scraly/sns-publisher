package main

import (
	"os"

	"github.com/scraly/sns-publisher/core"
)

func main() {

	//Retrieve environment variables
	topic := os.Getenv("SNS_TOPIC")
	awsRegion := os.Getenv("AWS_REGION")

	message := "salut!!!"

	//publish a message to a sns topic
	core.Publish(topic, awsRegion, message)
}
