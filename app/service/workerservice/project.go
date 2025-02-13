// Copyright (c) 2023 Fajar Laksono. All Rights Reserved.

package worker

import (
	"context"
	"fmt"

	modelkafka "fajarlaksono.github.io/laksono-api-service/app/model/kafka"
	"github.com/sirupsen/logrus"
)

const (
	jiraCommentPublicProperty = "sd.public.comment"
)

func (service *Worker) evaluateProjects(data *modelkafka.ProjectMessage, log *logrus.Entry) error {
	ctx := context.Background()
	log.Info("Project Evaluation started")

	rowNumberNonOverlapUpdated, err := service.DAOPostgres.EvaluateNonOverlapProjects(ctx)
	if err != nil {
		log.Errorf("EvaluateNonOverlapProjects error: %+v", err)
	}
	nonOverlapedMessage := fmt.Sprintf("%d new non overlapping projects.", rowNumberNonOverlapUpdated)
	log.Info(nonOverlapedMessage)
	service.DAOWebsocket.SendMessage(nonOverlapedMessage)

	rowNumberOverlapUpdated, err := service.DAOPostgres.EvaluateOverlapProjects(ctx)
	if err != nil {
		log.Errorf("EvaluateOverlapProjects error: %+v", err)
	}
	overlapedMessage := fmt.Sprintf("%d new overlapping projects.", rowNumberOverlapUpdated)
	log.Info(overlapedMessage)
	service.DAOWebsocket.SendMessage(overlapedMessage)

	log.Info("Project Evaluation finished")
	return nil
}
