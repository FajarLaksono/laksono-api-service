// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package route

import (
	"fmt"
	"net/http"

	"fajarlaksono.github.io/laksono-api-service/app/constant"
	"fajarlaksono.github.io/laksono-api-service/app/handler"
	apirequestmodel "fajarlaksono.github.io/laksono-api-service/app/model/api/request"
	apiresponsemodel "fajarlaksono.github.io/laksono-api-service/app/model/api/response"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

func AddProteusDTRouteGroup(service *handler.APIService) {
	miscTags := []string{"ProteusDT' Project Management"}

	// Create
	service.WebService.Route(
		service.WebService.POST(
			fmt.Sprintf("/project")).
			To(service.AddNewProjects).
			Doc("Add new projects").
			Notes(
				"<p>This endpoint is adding new projects.</p>").
			Reads(apirequestmodel.CreateProjectsRequest{}).
			Returns(http.StatusCreated, http.StatusText(http.StatusCreated), apiresponsemodel.CreateProjectsResponse{}).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), ErrorResponse{}).
			Returns(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), ErrorResponse{}).
			Metadata(restfulspec.KeyOpenAPITags, miscTags))

	// Read
	service.WebService.Route(
		service.WebService.GET(
			fmt.Sprintf("/project")).
			To(service.GetProjects).
			Doc("Get and List projects").
			Notes(
				"<p>This endpoint lists projects.</p>").
			Param(service.WebService.QueryParameter(constant.ProjectNameParameter, "Search "+constant.ProjectNameParameterString).DataType("string")).
			Param(service.WebService.QueryParameter(constant.ProjectIsOverlappingParameter, constant.ProjectIsOverlappingParameterString).DataType("bool")).
			Param(service.WebService.QueryParameter(constant.ProjectStartDateParameter, constant.ProjectStartDateParameterString).DataType("string")).
			Param(service.WebService.QueryParameter(constant.ProjectEndDateParameter, constant.ProjectEndDateParameterString).DataType("string")).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), apiresponsemodel.GetProjectsResponse{}).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), ErrorResponse{}).
			Returns(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), ErrorResponse{}).
			Metadata(restfulspec.KeyOpenAPITags, miscTags))

	service.WebService.Route(
		service.WebService.GET(
			fmt.Sprintf("/project/{%s}", constant.ProjectIDParameter)).
			To(service.GetProjectByID).
			Doc("Get a project detail by ID").
			Notes(
				"<p>This endpoint gets a project by ID.</p>").
			Param(service.WebService.PathParameter(constant.ProjectIDParameter, constant.ProjectIDParameterString).DataType("string").Required(true)).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), apiresponsemodel.GetProjectResponse{}).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), ErrorResponse{}).
			Returns(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), ErrorResponse{}).
			Metadata(restfulspec.KeyOpenAPITags, miscTags))

	// Update
	service.WebService.Route(
		service.WebService.PATCH(
			fmt.Sprintf("/project")).
			To(service.PatchProjectByIDBulk).
			Doc("updates a project by IDs").
			Notes(
				"<p>This endpoint updates a project by ID.</p>").
			Reads(apirequestmodel.UpdateProjectsRequest{}).
			Returns(http.StatusAccepted, http.StatusText(http.StatusAccepted), apiresponsemodel.UpdatedProjectResponse{}).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), ErrorResponse{}).
			Returns(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), ErrorResponse{}).
			Metadata(restfulspec.KeyOpenAPITags, miscTags))

	// Delete
	service.WebService.Route(
		service.WebService.DELETE(
			fmt.Sprintf("/project")).
			To(service.DeleteProjectByIDs).
			Doc("Delete project by IDs").
			Notes(
				"<p>This endpoint deletes projects by IDs.</p>").
			Reads(apirequestmodel.DeleteProjectsByIDs{}).
			Returns(http.StatusNoContent, http.StatusText(http.StatusOK), apiresponsemodel.DeletedProjectResponse{}).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), ErrorResponse{}).
			Returns(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), ErrorResponse{}).
			Metadata(restfulspec.KeyOpenAPITags, miscTags))

}
