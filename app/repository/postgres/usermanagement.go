// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package postgres

import (
	"context"
	"strings"

	consterror "fajarlaksono.github.io/laksono-api-service/app/constant/error"
	model "fajarlaksono.github.io/laksono-api-service/app/model/postgres"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (db *PostgreClient) CreateUser(ctx context.Context, data *model.User) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, consterror.DBTimeout)
	defer cancel()

	tx := db.Client.WithContext(ctx).Create(data)

	if tx.Error != nil {
		if strings.Contains(tx.Error.Error(), postgresUniqueKeyViolation) {
			return 0, errors.New(consterror.UniqueConstraintViolatedErrorMessage)
		}

		return 0, tx.Error
	}

	return tx.RowsAffected, nil
}

func (db *PostgreClient) GetUsers(ctx context.Context) (*model.Users, error) {
	ctx, cancel := context.WithTimeout(ctx, consterror.DBTimeout)
	defer cancel()

	var result model.Users

	err := db.Client.WithContext(ctx).
		Select("t.id, " +
			"t.username, " +
			"t.firstname, " +
			"t.lastname, " +
			"t.email ").
		Table("users AS t").
		Find(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, consterror.ErrNotFound
		}

		return nil, errors.Wrap(err, "unable to query")
	}

	return &result, nil
}
