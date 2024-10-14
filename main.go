package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/tromanini125/go-testcontainer-localstack-example/adapter/input/sqslistener"
	"github.com/tromanini125/go-testcontainer-localstack-example/application/service"
	"github.com/tromanini125/go-testcontainer-localstack-example/configuration"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configuration.LoadConfig()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	sqsClient := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider("test", "test", "")
		o.Region = "us-east-1"
		o.BaseEndpoint = aws.String("http://localhost:4566")
	})

	service := service.NewCardService()
	cardCreatedListener := sqslistener.NewCardCreatedListener(sqsClient, service)
	cardCreatedListener.Listen(ctx)
}
