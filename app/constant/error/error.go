// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package error

import (
	"errors"
	"time"
)

var (
	ErrNotFound                          = errors.New("data not found")
	ErrUnableToParse                     = errors.New("unable to parse data")
	ErrUniqueConstraintViolated          = errors.New("unique constraint violated")
	ErrForeignKeyConstraintViolated      = errors.New("foreign Key Constraint Violated Error")
	ErrUniqueConstraintViolation         = errors.New("unique constraint violated")
	UniqueConstraintViolatedErrorMessage = "unique constraint violated"
	DBSlowThreshold                      = time.Second
	DBTimeout                            = time.Second * 5
)
