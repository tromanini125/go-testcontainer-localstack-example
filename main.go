package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/tromanini125/go-testcontainer-localstack-example/adapter/input/sqslistener"
	"github.com/tromanini125/go-testcontainer-localstack-example/application/service"
	"github.com/tromanini125/go-testcontainer-localstack-example/configuration"
)

func main() {
	configuration.LoadConfig()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	sqsClient := sqs.NewFromConfig(cfg)
	service := service.NewCardService()
	sqslistener := sqslistener.NewCardCreatedListener(sqsClient, service)
	sqslistener.Listen()
}
