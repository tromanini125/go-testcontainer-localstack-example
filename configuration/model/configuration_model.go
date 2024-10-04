package model

type ConfingurationModel struct {
	CardCreatedQueue AWSQueue
}

type AWSQueue struct {
	URL                 string
	MaxNumberOfMessages int32
	WaitTimeSeconds     int32
}
