// Copyright (c) 2023 Fajar Laksono. All Rights Reserved.

package route

import (
	"fajarlaksono.github.io/laksono-api-service/app/handler"
	"github.com/emicklei/go-restful/v3"
)

func AddRoute(basePath string, service *handler.APIService) {
	webService := new(restful.WebService)
	service.WebService = webService

	webService.
		Path(basePath).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Laksono API Service").
		ApiVersion("0.0.0.0")

	AddUserManagementRoute(service)
}
