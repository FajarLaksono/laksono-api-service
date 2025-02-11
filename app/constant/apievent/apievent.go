// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package apievent

const (
	TypeDS = 001

	ServiceHealthy   = 0
	ServiceUnhealthy = 99

	CreateUserSuccess                = 001001200
	CreateUserInvalidBody            = 001001401
	CreateUserInvalidBodyEmailFormat = 001001402
	CreateUserConflict               = 001001403
	CreateUserInternalServerError    = 001001504
	CreateUserInsertingDBError       = 001001505
	CreateUserInvalidWorkspace       = 001001406

	CreateProjectsSuccess             = 001002200
	CreateProjectsInvalidBody         = 001002401
	CreateProjectsConflict            = 001002402
	CreateProjectsInternalServerError = 001003501
	CreateProjectsInsertingDBError    = 001004502

	GetUsersSuccess             = 001002200
	GetUsersDataNotFound        = 001002501
	GetUsersInternalServerError = 001002502
	GetUsersInvalidWorkspace    = 001002403
)
