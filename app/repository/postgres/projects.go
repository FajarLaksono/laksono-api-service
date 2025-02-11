// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package postgres

import (
	"context"
	"strings"

	consterror "fajarlaksono.github.io/laksono-api-service/app/constant/error"
	model "fajarlaksono.github.io/laksono-api-service/app/model/postgres"
	"github.com/pkg/errors"
)

func (db *PostgreClient) CreateProjects(ctx context.Context, data *model.Projects) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, consterror.DBTimeout)
	defer cancel()

	// Start a new transaction
	tx := db.Client.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// Perform the bulk insert within the transaction
	tx = tx.WithContext(ctx).Create(data)

	if tx.Error != nil {
		// Rollback the transaction in case of error
		tx.Rollback()

		if strings.Contains(tx.Error.Error(), postgresUniqueKeyViolation) {
			return 0, errors.New(consterror.UniqueConstraintViolatedErrorMessage)
		}

		return 0, tx.Error
	}

	// Commit the transaction if no errors
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return tx.RowsAffected, nil
}
