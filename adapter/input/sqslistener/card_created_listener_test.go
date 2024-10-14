package sqslistener

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
	"github.com/tromanini125/go-testcontainer-localstack-example/application/service"
	"github.com/tromanini125/go-testcontainer-localstack-example/configuration"
)

const (
	accesskey            = "a"
	secretkey            = "b"
	token                = "c"
	region               = "us-east-1"
	cardCreatedQueueName = "card-created-queue"
)

func init() {
	configuration.LoadConfig()
}

func createSQSClient(ctx context.Context, l *localstack.LocalStackContainer) (*sqs.Client, error) {
	mappedPort, err := l.MappedPort(ctx, nat.Port("4566/tcp"))
	if err != nil {
		return nil, err
	}

	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		return nil, err
	}
	defer provider.Close()

	host, err := provider.DaemonHost(ctx)
	if err != nil {
		return nil, err
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accesskey, secretkey, token)),
	)
	if err != nil {
		return nil, err
	}

	client := sqs.NewFromConfig(awsCfg, func(o *sqs.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider(accesskey, secretkey, token)
		o.Region = region
		o.BaseEndpoint = aws.String("http://" + host + ":" + mappedPort.Port())
	})

	return client, nil
}

func TestS3(t *testing.T) {
	ctx := context.Background()

	ctr, err := localstack.Run(ctx, "localstack/localstack:latest")
	if err != nil {
		t.Fatal(err)
	}

	sqsClient, err := createSQSClient(ctx, ctr)
	require.NoError(t, err)

	t.Run("Create Queue", func(t *testing.T) {
		outputQueue, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(cardCreatedQueueName),
		})
		require.NoError(t, err)
		assert.NotNil(t, outputQueue)
		configuration.Config.CardCreatedQueue.URL = *outputQueue.QueueUrl
	})

	t.Run("Publish Message", func(t *testing.T) {
		_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:     aws.String(cardCreatedQueueName),
			MessageBody:  aws.String("Hello, SQS!"),
			DelaySeconds: 0,
		})
		require.NoError(t, err)
	})

	t.Run("Receive Message", func(t *testing.T) {
		service := service.NewCardService()
		sqslistener := NewCardCreatedListener(sqsClient, service)
		sqslistener.FetchMessages(ctx)
	})
}
