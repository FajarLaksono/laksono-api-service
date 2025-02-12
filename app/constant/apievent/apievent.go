// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package apievent

const (
	ServiceNumber = 001

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
	CreateProjectsInternalServerError = 001002501
	CreateProjectsInsertingDBError    = 001002502

	GetProjectsSuccess                = 001003200
	ListProjectsInvalidQueryParameter = 001003401
	GetProjectsDBError                = 001003502

	GetProjectSuccess = 001004200
	GetProjectDBError = 001004502

	PatchingProjectsSuccess          = 001005200
	PatchProjectsInsertingDBError    = 001005501
	PatchProjectsInternalServerError = 001005502
	PatchProjectsConflict            = 001005403
	PatchProjectsInvalidBody         = 001005404

	DeleteProjectsByIDsSuccess             = 001006200
	DeleteProjectsByIDsInsertingDBError    = 001006501
	DeleteProjectsByIDsInternalServerError = 001006502
	DeleteProjectsByIDsConflict            = 001006403
	DeleteProjectsByIDsInvalidBody         = 001006404

	GetUsersSuccess             = 001002200
	GetUsersDataNotFound        = 001002501
	GetUsersInternalServerError = 001002502
	GetUsersInvalidWorkspace    = 001002403
)
