// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package model

import (
	"errors"
	"fmt"
	"time"

	"fajarlaksono.github.io/laksono-api-service/app/model"
	modelpostgres "fajarlaksono.github.io/laksono-api-service/app/model/postgres"
	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (data *CreateProjectRequest) ConvertToProject() (*modelpostgres.Project, error) {
	var result *modelpostgres.Project

	startDate, err := time.Parse("2006-01-02", data.StartDate)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid start date format: %v", err))
	}

	endDate, err := time.Parse("2006-01-02", data.EndDate)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid start date format: %v", err))
	}

	generatedUUID := uuid.New()
	result = &modelpostgres.Project{
		ID:        generatedUUID,
		Name:      data.Name,
		StartDate: startDate,
		EndDate:   endDate,
		DBBaseModel: model.DBBaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return result, nil
}

type CreateProjectsRequest []CreateProjectRequest

func (data *CreateProjectsRequest) ConvertToProject() (*modelpostgres.Projects, error) {
	var result modelpostgres.Projects

	for _, projectRequest := range *data {
		generatedUUID := uuid.New()

		startDate, err := time.Parse("2006-01-02", projectRequest.StartDate)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("invalid start date format: %v", err))
		}

		endDate, err := time.Parse("2006-01-02", projectRequest.EndDate)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("invalid start date format: %v", err))
		}

		project := modelpostgres.Project{
			ID:        generatedUUID,
			Name:      projectRequest.Name,
			StartDate: startDate,
			EndDate:   endDate,
			DBBaseModel: model.DBBaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		result = append(result, project)
	}

	return &result, nil
}

type UpdateProjectRequest struct {
	ID        string  `json:"id"`
	Name      *string `json:"name"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

type UpdateProjectsRequest []UpdateProjectRequest

type DeleteProjectsByIDs struct {
	IDs []uuid.UUID
}
