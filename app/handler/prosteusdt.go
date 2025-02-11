package handler

import (
	"net/http"

	"fajarlaksono.github.io/laksono-api-service/app/constant/apievent"
	constError "fajarlaksono.github.io/laksono-api-service/app/constant/error"
	modelapirequest "fajarlaksono.github.io/laksono-api-service/app/model/api/request"
	modelapiresponse "fajarlaksono.github.io/laksono-api-service/app/model/api/response"
	"github.com/emicklei/go-restful/v3"
	"github.com/pkg/errors"
)

func (service *APIService) AddNewProjects(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	// Read Request Body
	var requestData modelapirequest.CreateProjectsRequest
	err := req.ReadEntity(&requestData)
	if err != nil {
		WriteError(req, resp, http.StatusBadRequest, apievent.TypeDS, err, &Error{
			ErrorCode:    apievent.CreateProjectsInvalidBody,
			ErrorMessage: "unable process request: malformed request: " + err.Error(),
			ErrorLogMsg:  "Malformed request: invalid body",
		})

		return
	}

	insertData := requestData.ConvertToProject()

	rowAffected, err := service.DAOPostgres.CreateProjects(ctx, insertData)
	if err != nil {
		if errors.Is(err, constError.ErrUniqueConstraintViolated) || errors.Is(err, constError.ErrForeignKeyConstraintViolated) {
			WriteError(req, resp, http.StatusConflict, apievent.TypeDS, err, &Error{
				ErrorCode:    apievent.CreateProjectsConflict,
				ErrorMessage: "unable process request: " + err.Error(),
				ErrorLogMsg:  "unable process request: conflict",
			})

			return
		}
		WriteError(req, resp, http.StatusInternalServerError, apievent.TypeDS, err, &Error{
			ErrorCode:    apievent.CreateProjectsInternalServerError,
			ErrorMessage: "internal server error " + err.Error(),
			ErrorLogMsg:  "unable process request: internal server error",
		})

		return
	}

	if rowAffected < 1 {
		WriteError(req, resp, http.StatusInternalServerError, apievent.TypeDS, err, &Error{
			ErrorCode:    apievent.CreateProjectsInsertingDBError,
			ErrorMessage: "internal server error, failed creating projects",
			ErrorLogMsg:  "unable process request: internal server error, fail creating projects",
		})

		return
	}

	responseData := modelapiresponse.ConvertToCreateProjectsResponseFromProjects(insertData)

	Write(req, resp, http.StatusCreated, apievent.TypeDS,
		apievent.CreateProjectsSuccess, "success creating projects", responseData)
}

func (service *APIService) ListProjects(req *restful.Request, resp *restful.Response) {

}

func (service *APIService) GetProjectByID(req *restful.Request, resp *restful.Response) {

}

func (service *APIService) PatchProjectByIDBulk(req *restful.Request, resp *restful.Response) {

}

func (service *APIService) DeleteProjectByIDs(req *restful.Request, resp *restful.Response) {

}
