package core

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sts"
)

// CallerIdentity used to map to AWS caller identity
type CallerIdentity struct {
	Account string
	Arn     string
	UserId  string
}

// Publish function push a message in  AWS SNS topic
func Publish(topic string, awsRegion string, message string) {

	log.Println("Will publish a SNS notification on topic", topic, "with message", message)

	//Create a session object to talk to SNS (also make sure you have your key and secret setup in your .aws/credentials file)
	svc := sns.New(session.New())
	// params will be sent to the publish call included here is the bare minimum params to send a message.
	params := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String("arn:aws:sns:" + awsRegion + ":" + getAWSAccountID() + ":" + topic), //Get this from the Topic in the AWS console.
	}

	resp, err := svc.Publish(params) //Call to publish the message

	if err != nil { //Check for errors
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Println("[Publisher] Error:", err.Error())
		return
	}

	log.Println("[Publisher] Message publish successfully:", resp.String())
}

func getAWSAccountID() string {
	svc := sts.New(session.New())
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and Message from an error.
			log.Println(err.Error())
		}
	}

	// {
	//   Account: "350659893639",
	//   Arn: "arn:aws:iam::350659893639:user/terraform",
	//   UserId: "AIDAI7FKJ7X5OYRWY2QEM"
	// }
	callerIdentity := result.GoString()
	callerIdentity = strings.Replace(callerIdentity, "Account:", "\"Account\":", -1)
	callerIdentity = strings.Replace(callerIdentity, "Arn:", "\"Arn\":", -1)
	callerIdentity = strings.Replace(callerIdentity, "UserId:", "\"UserId\":", -1)

	//Return only AWS Account ID
	var m CallerIdentity
	errUnmarshal := json.Unmarshal([]byte(callerIdentity), &m)
	if errUnmarshal != nil {
		fmt.Println("err: ", errUnmarshal)
	}

	return m.Account
}
