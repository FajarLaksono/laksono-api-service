// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package model

import (
	"time"

	modelpostgres "fajarlaksono.github.io/laksono-api-service/app/model/postgres"
	"github.com/google/uuid"
)

type CreateProjectResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsOverlapping bool      `json:"is_overlapping"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ConvertToCreateProjectResponseFromProject(
	data *modelpostgres.Project) *CreateProjectResponse {
	var result *CreateProjectResponse

	generatedUUID := uuid.New()
	result = &CreateProjectResponse{
		ID:        generatedUUID,
		Name:      data.Name,
		StartDate: data.StartDate,
		EndDate:   data.EndDate,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return result
}

type CreateProjectsResponse []CreateProjectResponse

func ConvertToCreateProjectsResponseFromProjects(data *modelpostgres.Projects) *CreateProjectsResponse {
	var result CreateProjectsResponse

	for _, projectRequest := range *data {
		generatedUUID := uuid.New()
		project := CreateProjectResponse{
			ID:        generatedUUID,
			Name:      projectRequest.Name,
			StartDate: projectRequest.StartDate,
			EndDate:   projectRequest.EndDate,
			CreatedAt: projectRequest.CreatedAt,
			UpdatedAt: projectRequest.UpdatedAt,
		}
		result = append(result, project)
	}

	return &result
}

type GetProjectResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsOverlapping bool      `json:"is_overlapping"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ConvertToGetProjectResponseFromProject(
	data *modelpostgres.Project) *GetProjectResponse {
	var result *GetProjectResponse

	generatedUUID := uuid.New()
	result = &GetProjectResponse{
		ID:        generatedUUID,
		Name:      data.Name,
		StartDate: data.StartDate,
		EndDate:   data.EndDate,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return result
}

type GetProjectsResponse []GetProjectResponse

func ConvertToGetProjectsResponseFromProjects(data *modelpostgres.Projects) *GetProjectsResponse {
	var result GetProjectsResponse

	for _, projectRequest := range *data {
		generatedUUID := uuid.New()
		project := GetProjectResponse{
			ID:        generatedUUID,
			Name:      projectRequest.Name,
			StartDate: projectRequest.StartDate,
			EndDate:   projectRequest.EndDate,
			CreatedAt: projectRequest.CreatedAt,
			UpdatedAt: projectRequest.UpdatedAt,
		}
		result = append(result, project)
	}

	return &result
}

type UpdatedProjectResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsOverlapping bool      `json:"is_overlapping"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ConvertToUpdatedProjectResponseFromProject(
	data *modelpostgres.Project) *UpdatedProjectResponse {
	var result *UpdatedProjectResponse

	generatedUUID := uuid.New()
	result = &UpdatedProjectResponse{
		ID:        generatedUUID,
		Name:      data.Name,
		StartDate: data.StartDate,
		EndDate:   data.EndDate,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return result
}

type UpdatedProjectsResponse []UpdatedProjectResponse

func ConvertToUpdatedProjectsResponseFromProjects(data *modelpostgres.Projects) *UpdatedProjectsResponse {
	var result UpdatedProjectsResponse

	for _, projectRequest := range *data {
		generatedUUID := uuid.New()
		project := UpdatedProjectResponse{
			ID:        generatedUUID,
			Name:      projectRequest.Name,
			StartDate: projectRequest.StartDate,
			EndDate:   projectRequest.EndDate,
			CreatedAt: projectRequest.CreatedAt,
			UpdatedAt: projectRequest.UpdatedAt,
		}
		result = append(result, project)
	}

	return &result
}
