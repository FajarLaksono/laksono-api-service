package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"fajarlaksono.github.io/laksono-api-service/app/constant"
	"fajarlaksono.github.io/laksono-api-service/app/constant/apievent"
	constError "fajarlaksono.github.io/laksono-api-service/app/constant/error"
	modelapirequest "fajarlaksono.github.io/laksono-api-service/app/model/api/request"
	modelapiresponse "fajarlaksono.github.io/laksono-api-service/app/model/api/response"
	modelkafka "fajarlaksono.github.io/laksono-api-service/app/model/kafka"
	"github.com/cenkalti/backoff/v3"
	"github.com/emicklei/go-restful/v3"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func (service *APIService) AddNewProjects(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	// Read Request Body
	var requestData modelapirequest.CreateProjectsRequest
	err := req.ReadEntity(&requestData)
	if err != nil {
		WriteError(req, resp, http.StatusBadRequest, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.CreateProjectsInvalidBody,
			ErrorMessage: "unable process request: malformed request: " + err.Error(),
			ErrorLogMsg:  "Malformed request: invalid body",
		})

		return
	}

	insertData, err := requestData.ConvertToProject()
	if err != nil {
		WriteError(req, resp, http.StatusBadRequest, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.CreateProjectsInvalidBody,
			ErrorMessage: "unable process request: malformed request: " + err.Error(),
			ErrorLogMsg:  "Malformed request: invalid body: invalid date",
		})

		return
	}

	rowAffected, err := service.DAOPostgres.CreateProjects(ctx, insertData)
	if err != nil {
		if errors.Is(err, constError.ErrUniqueConstraintViolated) || errors.Is(err, constError.ErrForeignKeyConstraintViolated) {
			WriteError(req, resp, http.StatusConflict, apievent.ServiceNumber, err, &Error{
				ErrorCode:    apievent.CreateProjectsConflict,
				ErrorMessage: "unable process request: " + err.Error(),
				ErrorLogMsg:  "unable process request: conflict",
			})

			return
		}
		WriteError(req, resp, http.StatusInternalServerError, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.CreateProjectsInternalServerError,
			ErrorMessage: "internal server error " + err.Error(),
			ErrorLogMsg:  "unable process request: internal server error",
		})

		return
	}

	if rowAffected < 1 {
		WriteError(req, resp, http.StatusInternalServerError, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.CreateProjectsInsertingDBError,
			ErrorMessage: "internal server error, failed creating projects",
			ErrorLogMsg:  "unable process request: internal server error, fail creating projects",
		})

		return
	}

	responseData := modelapiresponse.ConvertToCreateProjectsResponseFromProjects(insertData)

	StartDate := time.Now()
	EndDate := time.Now()
	for _, data := range *insertData {
		if StartDate.After(data.StartDate) {
			StartDate = data.StartDate
		}

		if EndDate.After(data.EndDate) {
			EndDate = data.EndDate
		}
	}
	kafkaMessage := modelkafka.ProjectMessage{
		StartDate: StartDate.Format("2006-01-02"),
		EndDate:   EndDate.Format("2006-01-02"),
	}
	go service.pushKafkaMessage(kafkaMessage)

	Write(req, resp, http.StatusCreated, apievent.ServiceNumber,
		apievent.CreateProjectsSuccess, "success creating projects", responseData)
}

func (service *APIService) GetProjects(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	//Get Querie Parameters
	var err error
	var projectNameQP *string
	if req.QueryParameter(constant.ProjectNameParameter) != "" {
		projectName := req.QueryParameter(constant.ProjectNameParameter)
		projectNameQP = &projectName
	}

	var isOverlappingQP *bool
	if req.QueryParameter(constant.ProjectIsOverlappingParameter) != "" {
		isOverlapping, err := strconv.ParseBool(req.QueryParameter(constant.ProjectIsOverlappingParameter))
		isOverlappingQP = &isOverlapping
		if err != nil {
			WriteError(req, resp, http.StatusBadRequest, apievent.ServiceNumber, err, &Error{
				ErrorCode:    apievent.ListProjectsInvalidQueryParameter,
				ErrorMessage: "unable process request: malformed request: invalid boolean value " + err.Error(),
				ErrorLogMsg:  "Malformed request: invalid query parameter: is overlapping",
			})

			return
		}
	}

	var startDateQP *time.Time
	if req.QueryParameter(constant.ProjectStartDateParameter) != "" {
		startDate, err := time.Parse("2006-01-02", req.QueryParameter(constant.ProjectStartDateParameter))
		startDateQP = &startDate
		if err != nil {
			WriteError(req, resp, http.StatusBadRequest, apievent.ServiceNumber, err, &Error{
				ErrorCode:    apievent.ListProjectsInvalidQueryParameter,
				ErrorMessage: "unable process request: malformed request: use YYYY-MM-DD" + err.Error(),
				ErrorLogMsg:  "Malformed request: invalid query parameter: start date (use YYYY-MM-DD)",
			})

			return
		}
	}

	var endDateQP *time.Time
	if req.QueryParameter(constant.ProjectEndDateParameter) != "" {
		endDate, err := time.Parse("2006-01-02", req.QueryParameter(constant.ProjectEndDateParameter))
		endDateQP = &endDate
		if err != nil {
			WriteError(req, resp, http.StatusBadRequest, apievent.ServiceNumber, err, &Error{
				ErrorCode:    apievent.ListProjectsInvalidQueryParameter,
				ErrorMessage: "unable process request: malformed request: use YYYY-MM-DD" + err.Error(),
				ErrorLogMsg:  "Malformed request: invalid query parameter: end date (use YYYY-MM-DD)",
			})

			return
		}
	}

	total, result, err := service.DAOPostgres.GetProjects(ctx, projectNameQP, isOverlappingQP, startDateQP, endDateQP)
	if err != nil {
		WriteError(req, resp, http.StatusInternalServerError, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.GetProjectsDBError,
			ErrorMessage: "internal server error " + err.Error(),
			ErrorLogMsg:  "unable process request: internal server error",
		})
		return
	}

	responseData := modelapiresponse.ConvertToGetProjectsResponseFromProjects(total, result)

	Write(req, resp, http.StatusOK, apievent.ServiceNumber,
		apievent.GetProjectsSuccess, "success getting projects", responseData)
}

func (service *APIService) GetProjectByID(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	projectID := req.PathParameter(constant.ProjectIDParameter)

	result, err := service.DAOPostgres.GetProjectByID(ctx, projectID)
	if err != nil {
		WriteError(req, resp, http.StatusInternalServerError, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.GetProjectDBError,
			ErrorMessage: "internal server error " + err.Error(),
			ErrorLogMsg:  "unable process request: internal server error",
		})
		return
	}

	responseData := modelapiresponse.ConvertToGetProjectResponseFromProject(result)

	Write(req, resp, http.StatusOK, apievent.ServiceNumber,
		apievent.GetProjectSuccess, "success getting project detail", responseData)
}

func (service *APIService) PatchProjectByIDBulk(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	// Read Request Body
	var requestData modelapirequest.UpdateProjectsRequest
	err := req.ReadEntity(&requestData)
	if err != nil {
		WriteError(req, resp, http.StatusBadRequest, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.PatchProjectsInvalidBody,
			ErrorMessage: "unable process request: malformed request: " + err.Error(),
			ErrorLogMsg:  "Malformed request: invalid body",
		})

		return
	}

	rowAffected, err := service.DAOPostgres.PatchProjects(ctx, requestData)
	if err != nil {
		if errors.Is(err, constError.ErrUniqueConstraintViolated) || errors.Is(err, constError.ErrForeignKeyConstraintViolated) {
			WriteError(req, resp, http.StatusConflict, apievent.ServiceNumber, err, &Error{
				ErrorCode:    apievent.PatchProjectsConflict,
				ErrorMessage: "unable process request: " + err.Error(),
				ErrorLogMsg:  "unable process request: conflict",
			})

			return
		}
		WriteError(req, resp, http.StatusInternalServerError, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.PatchProjectsInternalServerError,
			ErrorMessage: "internal server error " + err.Error(),
			ErrorLogMsg:  "unable process request: internal server error",
		})

		return
	}

	if rowAffected < 1 {
		WriteError(req, resp, http.StatusInternalServerError, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.PatchProjectsInsertingDBError,
			ErrorMessage: "internal server error, failed patching projects",
			ErrorLogMsg:  "unable process request: internal server error, fail patching projects",
		})

		return
	}

	responseData := modelapiresponse.UpdatedProjectResponse{
		TotalRowUpdated: rowAffected,
	}

	StartDate := time.Now()
	EndDate := time.Now()
	for _, data := range requestData {
		if data.StartDate != nil {
			input, _ := time.Parse("2006-01-02", *data.StartDate)
			if StartDate.After(input) {
				StartDate = input
			}
		}

		if data.EndDate != nil {
			input, _ := time.Parse("2006-01-02", *data.EndDate)
			if EndDate.After(input) {
				EndDate = input
			}
		}
	}
	kafkaMessage := modelkafka.ProjectMessage{
		StartDate: StartDate.Format("2006-01-02"),
		EndDate:   EndDate.Format("2006-01-02"),
	}
	go service.pushKafkaMessage(kafkaMessage)

	Write(req, resp, http.StatusAccepted, apievent.ServiceNumber,
		apievent.PatchingProjectsSuccess, "success patching projects", responseData)
}

func (service *APIService) DeleteProjectByIDs(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	// Read Request Body
	var requestData modelapirequest.DeleteProjectsByIDs
	err := req.ReadEntity(&requestData)
	if err != nil {
		WriteError(req, resp, http.StatusBadRequest, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.DeleteProjectsByIDsInvalidBody,
			ErrorMessage: "unable process request: malformed request: " + err.Error(),
			ErrorLogMsg:  "Malformed request: invalid body",
		})

		return
	}

	if len(requestData.IDs) == 0 {
		WriteError(req, resp, http.StatusBadRequest, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.DeleteProjectsByIDsInvalidBody,
			ErrorMessage: "unable process request: malformed request: no ids provided" + err.Error(),
			ErrorLogMsg:  "Malformed request: invalid body: no ids provided",
		})

		return
	}

	rowAffected, err := service.DAOPostgres.DeleteProjects(ctx, requestData)
	if err != nil {
		if errors.Is(err, constError.ErrUniqueConstraintViolated) || errors.Is(err, constError.ErrForeignKeyConstraintViolated) {
			WriteError(req, resp, http.StatusConflict, apievent.ServiceNumber, err, &Error{
				ErrorCode:    apievent.DeleteProjectsByIDsConflict,
				ErrorMessage: "unable process request: " + err.Error(),
				ErrorLogMsg:  "unable process request: conflict",
			})

			return
		}
		WriteError(req, resp, http.StatusInternalServerError, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.DeleteProjectsByIDsInternalServerError,
			ErrorMessage: "internal server error " + err.Error(),
			ErrorLogMsg:  "unable process request: internal server error",
		})

		return
	}

	if rowAffected < 1 {
		WriteError(req, resp, http.StatusInternalServerError, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.DeleteProjectsByIDsInsertingDBError,
			ErrorMessage: "internal server error, failed deleting projects",
			ErrorLogMsg:  "unable process request: internal server error, fail deleting projects",
		})

		return
	}

	responseData := &modelapiresponse.DeletedProjectResponse{
		TotalRowDeleted: rowAffected,
	}

	kafkaMessage := modelkafka.ProjectMessage{
		StartDate: "",
		EndDate:   "",
	}
	go service.pushKafkaMessage(kafkaMessage)

	Write(req, resp, http.StatusOK, apievent.ServiceNumber,
		apievent.DeleteProjectsByIDsSuccess, "success deleting projects", responseData)
}

func (service *APIService) PatchEvaluateOverlapProjects(req *restful.Request, resp *restful.Response) {
	ctx := req.Request.Context()

	_, err := service.DAOPostgres.EvaluateNonOverlapProjects(ctx)
	if err != nil {
		WriteError(req, resp, http.StatusInternalServerError, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.PatchProjectsInternalServerError,
			ErrorMessage: "internal server error " + err.Error(),
			ErrorLogMsg:  "unable process request: internal server error",
		})

		return
	}

	_, err = service.DAOPostgres.EvaluateOverlapProjects(ctx)
	if err != nil {
		WriteError(req, resp, http.StatusInternalServerError, apievent.ServiceNumber, err, &Error{
			ErrorCode:    apievent.PatchProjectsInternalServerError,
			ErrorMessage: "internal server error " + err.Error(),
			ErrorLogMsg:  "unable process request: internal server error",
		})

		return
	}

	Write(req, resp, http.StatusNoContent, apievent.ServiceNumber,
		apievent.PatchingProjectsSuccess, "success evaluating overlapping projects", nil)
}

func (service *APIService) pushKafkaMessage(data modelkafka.ProjectMessage) {
	log := logrus.WithFields(logrus.Fields{
		"start_Daate": data.StartDate,
		"end_date":    data.EndDate,
	})

	message, err := json.Marshal(data)
	if err != nil {
		log.Errorf("unable to marshal data, err: %s", err.Error())

		return
	}

	err = backoff.RetryNotify(func() error {
		return service.KafkaWriter.WriteMessages(context.Background(), kafka.Message{Value: message})
	}, backoff.NewExponentialBackOff(), func(e error, duration time.Duration) {
		log.Warnf("Retry duration: %v sec, err: %s", duration.Seconds(), e.Error())
	})
	if err != nil {
		log.Errorf("unable to send project data to kafka, err: %s", err.Error())

		return
	}
	log.Info("successfully push project data to kafka")
}
