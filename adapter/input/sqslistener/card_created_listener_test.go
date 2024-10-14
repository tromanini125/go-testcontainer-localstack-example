package sqslistener

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"
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
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/tromanini125/go-testcontainer-localstack-example/adapter/input/sqslistener/model"
	"github.com/tromanini125/go-testcontainer-localstack-example/adapter/output/persistence"
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

var (
	sqsClient *sqs.Client
)

func TestMain(m *testing.M) {
	configuration.LoadConfig()
	ctx := context.Background()

	ctr, err := localstack.Run(ctx, "localstack/localstack:latest")
	if err != nil {
		panic(err)
	}

	sqsClient, err = createSQSClient(ctx, ctr)
	if err != nil {
		panic(err)
	}

	mysqlcontainer, err := createMysqlContainer(ctx)
	if err != nil {
		panic(err)
	}
	defer mysqlcontainer.Terminate(ctx)

}

func createMysqlContainer(ctx context.Context) (*mysql.MySQLContainer, error) {
	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithDatabase("cards"),
		mysql.WithUsername("root"),
		mysql.WithPassword("root"),
		mysql.WithScripts(filepath.Join("testdata", "cards.sql")),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err
	}
	host, _ := mysqlContainer.Host(ctx)
	port, _ := mysqlContainer.MappedPort(ctx, nat.Port("3306/tcp"))

	configuration.Config.DBConfig.Database = "cards"
	configuration.Config.DBConfig.Host = host
	configuration.Config.DBConfig.Port = port.Port()
	configuration.Config.DBConfig.User = "root"
	configuration.Config.DBConfig.Password = "root"
	return mysqlContainer, nil
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

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	persistence := persistence.NewCardRepository()
	cardNumber := "1234567890123456"

	t.Run("Create Queue", func(t *testing.T) {
		outputQueue, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
			QueueName: aws.String(cardCreatedQueueName),
		})
		require.NoError(t, err)
		assert.NotNil(t, outputQueue)
		configuration.Config.CardCreatedQueue.URL = *outputQueue.QueueUrl
	})

	t.Run("Publish Message", func(t *testing.T) {

		event := createEvent(cardNumber)
		jsonEvent, _ := json.Marshal(event)

		_, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:     aws.String(cardCreatedQueueName),
			MessageBody:  aws.String(string(jsonEvent)),
			DelaySeconds: 0,
		})
		require.NoError(t, err)
	})

	t.Run("Receive Message", func(t *testing.T) {
		service := service.NewCardService(persistence)
		sqslistener := NewCardCreatedListener(sqsClient, service)
		sqslistener.FetchMessages(ctx)
	})

	t.Run("Validate db", func(t *testing.T) {
		card, err := persistence.FindCardByNumber(ctx, cardNumber)
		require.NoError(t, err)

		assert.Equal(t, cardNumber, card.CardNumber)
	})
}

func createEvent(cardNumber string) model.CardCreatedEvent {
	return model.CardCreatedEvent{
		CardHolderName: "John Doe",
		CardNumber:     cardNumber,
		CVV:            "123",
		ExpiryDate:     "01/23",
	}
}
