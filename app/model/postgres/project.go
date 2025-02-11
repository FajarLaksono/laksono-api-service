// Copyright (c) 2024 FajarLaksono All Rights Reserved.

package model

import (
	"time"

	"fajarlaksono.github.io/laksono-api-service/app/model"
	"github.com/google/uuid"
)

// Project
type Project struct {
	ID            uuid.UUID `json:"id" grom:"primaryKey"`
	Name          string    `json:"name"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsOverlapping bool      `json:"is_overlapping"`
	model.DBBaseModel
}

// Projects
type Projects []Project
