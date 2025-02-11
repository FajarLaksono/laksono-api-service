// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package route

import (
	"fmt"
	"net/http"

	"fajarlaksono.github.io/laksono-api-service/app/handler"
	"fajarlaksono.github.io/laksono-api-service/app/model"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// AddUserManagementRoute endpoints group
func AddUserManagementRoute(service *handler.APIService) {
	miscTags := []string{"User Management"}

	service.WebService.Route(
		service.WebService.POST(
			fmt.Sprintf("/user")).
			To(service.CreateUser).
			Doc("Create a new User").
			Notes(
				"<p>This endpoint is creating a new user.</p>").
			Reads(model.CreateUserRequest{}).
			Returns(http.StatusCreated, http.StatusText(http.StatusCreated), model.CreateUserResponse{}).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), ErrorResponse{}).
			Returns(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), ErrorResponse{}).
			Metadata(restfulspec.KeyOpenAPITags, miscTags))

	service.WebService.Route(
		service.WebService.GET(
			fmt.Sprintf("/users")).
			To(service.GetUsers).
			Doc("Get list of user").
			Notes(
				"<p>This endpoint is getting list of users.</p>").
			Returns(http.StatusCreated, http.StatusText(http.StatusCreated), model.GetUsersResponse{}).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), ErrorResponse{}).
			Returns(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), ErrorResponse{}).
			Metadata(restfulspec.KeyOpenAPITags, miscTags))
}
