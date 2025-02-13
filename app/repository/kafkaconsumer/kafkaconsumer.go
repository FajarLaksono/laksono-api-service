package KafkaConsumer

import (
	worker "fajarlaksono.github.io/laksono-api-service/app/service/workerservice"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// Consumer represents an interface that define the behavior of stream consumer.
type Consumer interface {
	Consume() error
	ProcessCommitRecord()
}

type KafkaConsumerService struct {
	Consumer          Consumer
	RetryConsumer     Consumer
	NumberOfProcessor int
	RecordFinishChan  chan *kafka.Message
	Realm             string
	Processors        []worker.WorkerService
}

// Start starts processing data from the stream
func (workerService *KafkaConsumerService) Start() error {
	for _, jobProcessor := range workerService.Processors {
		go jobProcessor.Start()
	}

	go workerService.Consumer.ProcessCommitRecord()
	go workerService.RetryConsumer.ProcessCommitRecord()

	go func() {
		err := workerService.RetryConsumer.Consume()
		if err != nil {
			logrus.Error(err)
		}
	}()
	err := workerService.Consumer.Consume()

	for _, jobProcessor := range workerService.Processors {
		go jobProcessor.Stop()
	}

	return err
}
