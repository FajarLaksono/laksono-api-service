package kafka

import (
	"context"
	"time"

	model "fajarlaksono.github.io/laksono-api-service/app/model/kafka"
	modelkafka "fajarlaksono.github.io/laksono-api-service/app/model/kafka"
	"github.com/segmentio/kafka-go"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Consumer represents a consumer that consume AWS Kinesis stream.
type Consumer struct {
	reader           *kafka.Reader
	RecordChan       chan *modelkafka.Record
	RecordFinishChan chan *kafka.Message
	ctx              context.Context
	retryWaitingTime int
}

// NewConsumer creates a new Kinesis consumer
func NewConsumer(ctx context.Context, kafkaReader *kafka.Reader, recordChan chan *model.Record,
	recordFinishChan chan *kafka.Message, retryWaitingTime int,
) *Consumer {
	return &Consumer{
		reader:           kafkaReader,
		RecordChan:       recordChan,
		RecordFinishChan: recordFinishChan,
		ctx:              ctx,
		retryWaitingTime: retryWaitingTime,
	}
}

func (consumer *Consumer) Consume() error {
	for {
		message, err := consumer.reader.FetchMessage(consumer.ctx)
		if err != nil {
			logrus.Error("Unable to fetch message")

			return errors.Wrap(err, "Unable to fetch message")
		}

		if consumer.retryWaitingTime != 0 {
			delay := message.Time.Add(time.Duration(consumer.retryWaitingTime) * time.Minute).UTC().Sub(time.Now().UTC())
			time.Sleep(delay)
		}

		consumer.RecordChan <- &model.Record{
			Message: &message,
			Retry:   consumer.retryWaitingTime != 0,
		}
	}
}

// ProcessCommitRecord the process for commit record
func (consumer *Consumer) ProcessCommitRecord() {
	for {
		finishedRecord := <-consumer.RecordFinishChan
		err := consumer.reader.CommitMessages(consumer.ctx, *finishedRecord)
		if err != nil {
			logrus.Warnf("Unable to Commit Messages err: %s", err)
		}
	}
}
