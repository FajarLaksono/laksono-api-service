// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	consterror "fajarlaksono.github.io/laksono-api-service/app/constant/error"
	modelapirequest "fajarlaksono.github.io/laksono-api-service/app/model/api/request"
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

func (db *PostgreClient) GetProjects(ctx context.Context, projectNameFilter *string, isOverlappingFilter *bool, startDateFilter, endDateFilter *time.Time) (
	int64, *model.Projects, error) {
	ctx, cancel := context.WithTimeout(ctx, consterror.DBTimeout)
	defer cancel()

	var result *model.Projects
	var totalData int64

	tx := db.Client.WithContext(ctx).
		Table("projects AS p").
		Select("p.id, " +
			"p.name, " +
			"p.start_date, " +
			"p.end_date, " +
			"p.created_at, " +
			"p.updated_at, " +
			"p.deleted_at").
		Where("p.deleted_at IS NULL").
		Order("p.created_at DESC")

	if projectNameFilter != nil {
		tx = tx.Where("p.name ILIKE ?", "%"+*projectNameFilter+"%")
	}

	if isOverlappingFilter != nil {
		tx = tx.Where("p.is_overlapping = ?", strconv.FormatBool(*isOverlappingFilter))
	}

	if startDateFilter != nil {
		tx = tx.Where("p.start_Date >= ?", startDateFilter.Format("2006-01-02"))
	}

	if endDateFilter != nil {
		tx = tx.Where("p.end_date <= ?", endDateFilter.Format("2006-01-02"))
	}

	err := tx.Find(&result).Error
	if err != nil {
		return 0, nil, err
	}
	err = tx.Count(&totalData).Error

	return totalData, result, err
}

func (db *PostgreClient) GetProjectByID(ctx context.Context, projectID string) (*model.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, consterror.DBTimeout)
	defer cancel()

	var result *model.Project

	tx := db.Client.WithContext(ctx).
		Table("projects AS p").
		Select("p.id, " +
			"p.name, " +
			"p.start_date, " +
			"p.end_date, " +
			"p.created_at, " +
			"p.updated_at, " +
			"p.deleted_at").
		Order("p.created_at DESC")

	err := tx.First(&result).Error

	return result, err
}

func (db *PostgreClient) PatchProjects(ctx context.Context, input modelapirequest.UpdateProjectsRequest) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, consterror.DBTimeout)
	defer cancel()

	var totalRowsAffected int64

	fmt.Printf("\n Fajar 2.1 \n")
	tx := db.Client.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// Perform batch updates
	for _, data := range input {

		update := make(map[string]interface{})
		if data.Name != nil {
			update["name"] = data.Name
		}

		if data.StartDate != nil {
			update["start_date"] = data.StartDate
		}

		if data.EndDate != nil {
			update["end_date"] = data.EndDate
		}

		result := tx.WithContext(ctx).Model(&model.Project{}).
			Where("id = ?", data.ID).
			Updates(update)

		if result.Error != nil {
			tx.Rollback()
			return 0, result.Error
		}

		totalRowsAffected += result.RowsAffected
	}

	// Commit the transaction if no errors
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return totalRowsAffected, nil
}

func (db *PostgreClient) DeleteProjects(ctx context.Context, input modelapirequest.DeleteProjectsByIDs) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, consterror.DBTimeout)
	defer cancel()

	tx := db.Client.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// Perform bulk deletion
	result := tx.WithContext(ctx).
		Table("projects").
		Where("id IN ?", input.IDs).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
		})

	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return result.RowsAffected, nil
}
