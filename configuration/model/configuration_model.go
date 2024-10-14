package model

type ConfingurationModel struct {
	CardCreatedQueue AWSQueue
	DBConfig         DBConfig
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

type AWSQueue struct {
	URL                 string
	MaxNumberOfMessages int32
	WaitTimeSeconds     int32
}
