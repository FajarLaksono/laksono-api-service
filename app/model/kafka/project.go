package model

import "github.com/segmentio/kafka-go"

type Record struct {
	Message *kafka.Message
	Retry   bool
}

type ProjectMessage struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
