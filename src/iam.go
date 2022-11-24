package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/kataras/golog"
)

func getAwsGroups(groupName string) *iam.GetGroupOutput {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_DEFAULT_REGION")),
	})
	if err != nil {
		golog.Info("Session Olusturulamadi!")
		os.Exit(1)
	}

	iamClient := iam.New(sess)

	group, err := iamClient.GetGroup(&iam.GetGroupInput{
		GroupName: aws.String(groupName),
		MaxItems:  aws.Int64(500),
	})
	if err != nil {
		golog.Errorf("There is a problem, %s", err)
		os.Exit(1)
	}
	return group
}
