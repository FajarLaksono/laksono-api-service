package handler

import (
	"errors"
	"net/http"
	"net/mail"

	"fajarlaksono.github.io/laksono-api-service/app/constant/apievent"
	constError "fajarlaksono.github.io/laksono-api-service/app/constant/error"
	"fajarlaksono.github.io/laksono-api-service/app/model/api/response"
	model "fajarlaksono.github.io/laksono-api-service/app/model/postgres"
	"github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
)

// CreateUser handles to create a user.
func (service *APIService) CreateUser(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	userUUID := uuid.New()

	// Read Request Body
	var requestData model.CreateUserRequest
	err := req.ReadEntity(&requestData)
	if err != nil {
		response.WriteError(req, resp, http.StatusBadRequest, apievent.TypeDS, err, &response.Error{
			ErrorCode:    apievent.CreateUserInvalidBody,
			ErrorMessage: "unable process request: malformed request: " + err.Error(),
			ErrorLogMsg:  "Malformed request: invalid body",
		})

		return
	}

	// Validate eamil address
	_, err = mail.ParseAddress(requestData.Email)
	if err != nil {
		response.WriteError(req, resp, http.StatusBadRequest, apievent.TypeDS, err, &response.Error{
			ErrorCode:    apievent.CreateUserInvalidBodyEmailFormat,
			ErrorMessage: "Malformed request",
			ErrorLogMsg:  "Malformed request: bad emails format",
		})

		return
	}

	userData := requestData.ConvertToUser(userUUID)

	rowAffected, err := service.DAOPostgres.CreateUser(ctx, userData)
	if err != nil {
		if errors.Is(err, constError.ErrUniqueConstraintViolated) || errors.Is(err, constError.ErrForeignKeyConstraintViolated) {
			response.WriteError(req, resp, http.StatusConflict, apievent.TypeDS, err, &response.Error{
				ErrorCode:    apievent.CreateUserConflict,
				ErrorMessage: "malformed request: " + err.Error(),
				ErrorLogMsg:  "unable process request: conflict",
			})

			return
		}
		response.WriteError(req, resp, http.StatusInternalServerError, apievent.TypeDS, err, &response.Error{
			ErrorCode:    apievent.CreateUserInternalServerError,
			ErrorMessage: "internal server error " + err.Error(),
			ErrorLogMsg:  "unable process request: internal server error",
		})

		return
	}

	if rowAffected < 1 {
		response.WriteError(req, resp, http.StatusInternalServerError, apievent.TypeDS, err, &response.Error{
			ErrorCode:    apievent.CreateUserInsertingDBError,
			ErrorMessage: "internal server error, failed creating user",
			ErrorLogMsg:  "unable process request: internal server error, fail creating user",
		})

		return
	}

	responseData := userData.ConvertToCreateUserResponse()

	response.Write(req, resp, http.StatusCreated, apievent.TypeDS,
		apievent.CreateUserSuccess, "success creating a user", responseData)
}

// GetUsers handles to create a user.
func (service *APIService) GetUsers(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	users, err := service.DAOPostgres.GetUsers(ctx)
	if err != nil {
		if errors.Is(err, constError.ErrNotFound) {
			response.WriteError(req, resp, http.StatusNotFound, apievent.TypeDS, err, &response.Error{
				ErrorCode:    apievent.GetUsersDataNotFound,
				ErrorMessage: "malformed request: " + err.Error(),
				ErrorLogMsg:  "unable process request: conflict",
			})

			return
		}
		response.WriteError(req, resp, http.StatusInternalServerError, apievent.TypeDS, err, &response.Error{
			ErrorCode:    apievent.GetUsersInternalServerError,
			ErrorMessage: "internal server error " + err.Error(),
			ErrorLogMsg:  "unable process request: internal server error",
		})

		return
	}

	responseData := users.ConvertToGetUsersResponse()

	response.Write(req, resp, http.StatusOK, apievent.TypeDS,
		apievent.GetUsersSuccess, "success getting users", responseData)
}
