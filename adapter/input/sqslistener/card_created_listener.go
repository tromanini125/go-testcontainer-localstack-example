package sqslistener

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/tromanini125/go-testcontainer-localstack-example/adapter/input/sqslistener/model"
	"github.com/tromanini125/go-testcontainer-localstack-example/application/domain"
	"github.com/tromanini125/go-testcontainer-localstack-example/application/input"
	"github.com/tromanini125/go-testcontainer-localstack-example/configuration"
)

var (
	ticker *time.Ticker
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

func (cl *cardCreatedListener) Listen(ctx context.Context) {
	ticker = time.NewTicker(time.Second * 1)
	func() {
		for range ticker.C {
			cl.FetchMessages(ctx)
		}
	}()
}

func (actor *cardCreatedListener) FetchMessages(ctx context.Context) {
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
	for _, message := range messages {
		log.Printf("Got new message from queue")
		log.Printf("message: %s, body: %s", *message.MessageId, *message.Body)

		cardCreatedEvent, err := parseMessage(message.Body)
		if err != nil {
			log.Printf("Couldn't parse message body into CardCreatedEvent. Skipping message.")
		} else {
			cardDomain := mapEventToDomain(cardCreatedEvent)
			err = actor.service.Execute(ctx, cardDomain)
			if err != nil {
				log.Printf("Couldn't Process message %s", *message.MessageId)
			} else {
				log.Printf("Message %s processed successfully", *message.MessageId)
			}
		}

		actor.sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(configuration.Config.CardCreatedQueue.URL),
			ReceiptHandle: message.ReceiptHandle,
		})
	}
}

func parseMessage(messageBody *string) (*model.CardCreatedEvent, error) {
	var cardCreatedEvent model.CardCreatedEvent
	err := json.Unmarshal([]byte(*messageBody), &cardCreatedEvent)
	if err != nil {
		log.Printf("Couldn't parse message body into CardCreatedEvent. Here's why: %v\n", err)
		return nil, err
	}
	log.Printf("Message parsed, body: %+v", &cardCreatedEvent)

	return &cardCreatedEvent, nil
}

func mapEventToDomain(event *model.CardCreatedEvent) *domain.Card {
	return &domain.Card{
		CardHolderName: event.CardHolderName,
		CardNumber:     event.CardNumber,
		CVV:            event.CVV,
		ExpiryDate:     event.ExpiryDate,
	}
}
