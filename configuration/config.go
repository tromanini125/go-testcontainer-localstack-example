package configuration

import (
	"os"

	"github.com/tromanini125/go-testcontainer-localstack-example/configuration/model"
)

var (
	Config *model.ConfingurationModel
)

func LoadConfig() {
	// Load configuration from a file or environment variables
	Config = &model.ConfingurationModel{
		CardCreatedQueue: model.AWSQueue{
			URL:                 os.Getenv("CARD_CREATED_QUEUE_URL"),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     3,
		},
	}
}
