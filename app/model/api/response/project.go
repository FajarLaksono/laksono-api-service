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
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	IsOverlapping bool      `json:"is_overlapping"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ConvertToCreateProjectResponseFromProject(
	data *modelpostgres.Project) *CreateProjectResponse {
	var result *CreateProjectResponse

	result = &CreateProjectResponse{
		ID:        data.ID,
		Name:      data.Name,
		StartDate: data.StartDate.Format("2006-01-02"),
		EndDate:   data.EndDate.Format("2006-01-02"),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return result
}

type CreateProjectsResponse []CreateProjectResponse

func ConvertToCreateProjectsResponseFromProjects(data *modelpostgres.Projects) *CreateProjectsResponse {
	var result CreateProjectsResponse

	for _, row := range *data {
		project := CreateProjectResponse{
			ID:        row.ID,
			Name:      row.Name,
			StartDate: row.StartDate.Format("2006-01-02"),
			EndDate:   row.EndDate.Format("2006-01-02"),
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
		result = append(result, project)
	}

	return &result
}

type GetProjectResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	IsOverlapping bool      `json:"is_overlapping"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ConvertToGetProjectResponseFromProject(
	data *modelpostgres.Project) *GetProjectResponse {
	var result *GetProjectResponse

	result = &GetProjectResponse{
		ID:        data.ID,
		Name:      data.Name,
		StartDate: data.StartDate.Format("2006-01-02"),
		EndDate:   data.EndDate.Format("2006-01-02"),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt, // "created_at": "2025-02-12T07:21:04.232215Z", "updated_at": "2025-02-12T07:21:04.232215Z"
	}

	return result
}

type GetProjectsResponse struct {
	Total int64                `json:"total"`
	Data  []GetProjectResponse `json:"data"`
}

func ConvertToGetProjectsResponseFromProjects(total int64, data *modelpostgres.Projects) *GetProjectsResponse {
	var result GetProjectsResponse

	for _, row := range *data {
		project := GetProjectResponse{
			ID:        row.ID,
			Name:      row.Name,
			StartDate: row.StartDate.Format("2006-01-02"),
			EndDate:   row.EndDate.Format("2006-01-02"),
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
		result.Data = append(result.Data, project)
	}

	result.Total = total

	return &result
}

type UpdatedProjectResponse struct {
	TotalRowUpdated int64 `json:"total_row_updated"`
}

type DeletedProjectResponse struct {
	TotalRowDeleted int64 `json:"total_row_deleted"`
}
