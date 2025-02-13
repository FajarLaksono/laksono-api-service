package worker

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"fajarlaksono.github.io/laksono-api-service/app/config"
	model "fajarlaksono.github.io/laksono-api-service/app/model/kafka"
	modelkafka "fajarlaksono.github.io/laksono-api-service/app/model/kafka"
	"fajarlaksono.github.io/laksono-api-service/app/repository"
	"fajarlaksono.github.io/laksono-api-service/app/repository/websocketclient"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// Worker interface
type WorkerService interface {
	Start()
	Stop()
}

// Processor handles decoding and parsing Kafka record and send them to Storage and CrashReporter service
type Worker struct {
	WorkerService
	Config                *config.Config
	RecordChan            chan *modelkafka.Record
	RecordFinishChan      chan *kafka.Message
	RecordRetryFinishChan chan *kafka.Message
	QuitChan              chan bool
	RetryWriter           *kafka.Writer
	DAOPostgres           repository.PostgresDAO
	DAOWebsocket          *websocketclient.Client
	LogBuffer             *bytes.Buffer
}

// New creates new processor
func New(conf *config.Config, recordChan chan *modelkafka.Record,
	recordFinishChan chan *kafka.Message, quitChan chan bool, retryWriter *kafka.Writer,
	recordRetryFinishChan chan *kafka.Message, daoPostgres repository.PostgresDAO, wsclient *websocketclient.Client,
) *Worker {
	w := Worker{
		Config:                conf,
		RecordChan:            recordChan,
		RecordFinishChan:      recordFinishChan,
		RecordRetryFinishChan: recordRetryFinishChan,
		QuitChan:              quitChan,
		RetryWriter:           retryWriter,
		LogBuffer:             new(bytes.Buffer),
		DAOPostgres:           daoPostgres,
		DAOWebsocket:          wsclient,
	}

	return &w
}

// Start stops the worker
func (service *Worker) Start() {
	logrus.Info("worker start")

	if service.LogBuffer == nil {
		service.LogBuffer = new(bytes.Buffer)
		mw := io.MultiWriter(os.Stdout, service.LogBuffer)
		logrus.SetOutput(mw)
	}

	for {
		select {
		case <-service.QuitChan:
			logrus.Info("received quit signal, quitting processor...")
			close(service.RecordChan)
			close(service.QuitChan)

			return

		case record := <-service.RecordChan:
			service.LogBuffer.Reset()
			projectMessage := &modelkafka.ProjectMessage{}

			err := json.Unmarshal(record.Message.Value, projectMessage)
			if err != nil {
				logrus.Error("unable to unmarshal ProjectMessage: ", err)
				service.finishJob(record)

				continue
			}

			log := logrus.WithFields(logrus.Fields{
				"start_date": projectMessage.StartDate,
				"end_date":   projectMessage.EndDate,
			})

			err = service.evaluateProjects(projectMessage, log)
			if err != nil {
				log.Error("unable proces projectMessage: ", err.Error())
				service.finishJob(record)

				continue
			}
		}
	}
}

// Stop stops the processor
func (service *Worker) Stop() {
	service.QuitChan <- true
}

func (service *Worker) finishJob(record *model.Record) {
	if record.Retry {
		service.RecordRetryFinishChan <- record.Message
	} else {
		service.RecordFinishChan <- record.Message
	}
}
