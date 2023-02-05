package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	sessionV1 "github.com/aws/aws-sdk-go/aws/session"
	stsV1 "github.com/aws/aws-sdk-go/service/sts"
)

type GetCallerIdentityAPI interface {
	GetCallerIdentity(ctx context.Context,
		params *sts.GetCallerIdentityInput,
		optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

type GetCallerIdentityV1API interface {
	GetCallerIdentity(params *stsV1.GetCallerIdentityInput) (*stsV1.GetCallerIdentityOutput, error)
}

func GetCallerIdentity(ctx context.Context, client GetCallerIdentityAPI, params *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return client.GetCallerIdentity(ctx, params)
}

func GetCallerIdentityV1(client GetCallerIdentityV1API, params *stsV1.GetCallerIdentityInput) (*stsV1.GetCallerIdentityOutput, error) {
	return client.GetCallerIdentity(params)
}

func V1() error {

	sess := sessionV1.Must(sessionV1.NewSessionWithOptions(sessionV1.Options{
		SharedConfigState: sessionV1.SharedConfigEnable,
	}))
	sess.Config.Region = aws.String("us-east-2")

	svc := stsV1.New(sess)

	result, err := GetCallerIdentityV1(svc, &stsV1.GetCallerIdentityInput{})
	if err != nil {
		return err
	}

	fmt.Printf("\nVersion 1\n")
	fmt.Printf("Account ID: %s\n", *result.Account)
	fmt.Printf("User ID: %s\n", *result.UserId)
	fmt.Printf("ARN: %s\n", *result.Arn)
	fmt.Printf("\n")
	return nil

}

func V1NoShared() error {

	sess := sessionV1.Must(sessionV1.NewSessionWithOptions(sessionV1.Options{}))
	sess.Config.Region = aws.String("us-east-2")

	svc := stsV1.New(sess)

	result, err := GetCallerIdentityV1(svc, &stsV1.GetCallerIdentityInput{})
	if err != nil {
		return err
	}

	fmt.Printf("\nNO AWS SSO/Identity: Version 1 NoShared\n")
	fmt.Printf("Account ID: %s\n", *result.Account)
	fmt.Printf("User ID: %s\n", *result.UserId)
	fmt.Printf("ARN: %s\n", *result.Arn)
	fmt.Printf("\n")
	return nil

}

func V2() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	cfg.Region = "us-east-2"
	client := sts.NewFromConfig(cfg)

	result, err := GetCallerIdentity(context.TODO(), client, &sts.GetCallerIdentityInput{})
	if err != nil {

		return err
	}
	fmt.Printf("\nVersion 2\n")
	fmt.Printf("Account ID: %s\n", *result.Account)
	fmt.Printf("User ID: %s\n", *result.UserId)
	fmt.Printf("ARN: %s\n", *result.Arn)
	fmt.Printf("\n")

	return nil
}

func main() {

	err := V1()
	if err != nil {
		fmt.Println(err)
	}
	err = V2()
	if err != nil {
		fmt.Println(err)
	}

	err = V1NoShared()
	if err != nil {
		fmt.Printf("You appear to be picking up AWS SSO\n\n")
		fmt.Println(err)
	}
}
