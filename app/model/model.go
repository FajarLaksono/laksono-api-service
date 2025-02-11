// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package model

import (
	"time"
)

// BaseModel holds base model of table.
type DBBaseModel struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

