// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package model

import (
	"time"

	"fajarlaksono.github.io/laksono-api-service/app/model"
	modelpostgres "fajarlaksono.github.io/laksono-api-service/app/model/postgres"
	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func (data *CreateProjectRequest) ConvertToProject() *modelpostgres.Project {
	var result *modelpostgres.Project

	generatedUUID := uuid.New()
	result = &modelpostgres.Project{
		ID:        generatedUUID,
		Name:      data.Name,
		StartDate: data.StartDate,
		EndDate:   data.EndDate,
		DBBaseModel: model.DBBaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return result
}

type CreateProjectsRequest []CreateProjectRequest

func (data *CreateProjectsRequest) ConvertToProject() *modelpostgres.Projects {
	var result modelpostgres.Projects

	for _, projectRequest := range *data {
		generatedUUID := uuid.New()
		project := modelpostgres.Project{
			ID:        generatedUUID,
			Name:      projectRequest.Name,
			StartDate: projectRequest.StartDate,
			EndDate:   projectRequest.EndDate,
			DBBaseModel: model.DBBaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		result = append(result, project)
	}

	return &result
}

type UpdateProjectRequest struct {
	ID        string     `json:"id"`
	Name      *string    `json:"name"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

type UpdateProjectsRequest []UpdateProjectRequest
