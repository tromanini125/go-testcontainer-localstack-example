package sqslistener

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/tromanini125/go-testcontainer-localstack-example/application/input"
	"github.com/tromanini125/go-testcontainer-localstack-example/configuration"
)

func NewCardCreatedListener(sqsClient *sqs.Client, svc input.SaveNewCardUseCase) *cardCreatedListener {
	return &cardCreatedListener{
		sqsClient: sqsClient,
		service:   svc,
	}
}

type cardCreatedListener struct {
	sqsClient *sqs.Client
	service   input.SaveNewCardUseCase
}

func (l *cardCreatedListener) Listen() error {
	l.GetMessages(context.Background())
	return nil
}

// GetMessages uses the ReceiveMessage action to get messages from an Amazon SQS queue.
func (actor *cardCreatedListener) GetMessages(ctx context.Context) ([]types.Message, error) {
	var messages []types.Message
	result, err := actor.sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(configuration.Config.CardCreatedQueue.URL),
		MaxNumberOfMessages: configuration.Config.CardCreatedQueue.MaxNumberOfMessages,
		WaitTimeSeconds:     configuration.Config.CardCreatedQueue.WaitTimeSeconds,
	})
	if err != nil {
		log.Printf("Couldn't get messages from queue %v. Here's why: %v\n", configuration.Config.CardCreatedQueue.URL, err)
	} else {
		messages = result.Messages
	}
	log.Printf("Got messages from queue")
	for _, message := range messages {
		log.Printf("message: %s, body: %v", *message.MessageId, message.Body)
	}
	return messages, err
}
