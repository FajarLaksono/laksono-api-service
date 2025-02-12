// Copyright (c) 2023 Fajar Laksono. All Rights Reserved.

package api

import (
	"fmt"
	"net"
	"net/http"

	"fajarlaksono.github.io/laksono-api-service/app/config"
	"fajarlaksono.github.io/laksono-api-service/app/handler"
	"fajarlaksono.github.io/laksono-api-service/app/repository"
	"fajarlaksono.github.io/laksono-api-service/app/service/api/route"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

//nolint:funlen
func NewService(Config *config.Config, DAOPostgres repository.PostgresDAO, kafkaWriter *kafka.Writer) (*Service, error) {
	apiservice := handler.New(DAOPostgres, kafkaWriter)
	route.AddRoute(Config.BasePath, apiservice)

	goRestfulContainer := restful.NewContainer()
	goRestfulContainer.Add(apiservice.WebService)
	goRestfulContainer.Add(addSwagger(Config.BasePath, goRestfulContainer))

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{},
		AllowedHeaders: Config.AllowedHeaders,
		AllowedDomains: Config.AllowedOrigins,
		AllowedMethods: Config.AllowedMethods,
		CookiesAllowed: true,
		Container:      restful.DefaultContainer,
	}
	goRestfulContainer.Filter(cors.Filter)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", Config.Port))
	if err != nil {
		return nil, err
	}

	return &Service{
		config:           Config,
		restfulContainer: goRestfulContainer,
		listener:         listener,
	}, nil
}

type Service struct {
	config           *config.Config
	restfulContainer *restful.Container
	listener         net.Listener
}

func (s *Service) Start() error {
	return http.Serve(s.listener, s.restfulContainer)
}

func (s *Service) Stop() {
	if err := s.listener.Close(); err != nil {
		logrus.WithError(err).Error("unable to clos HTTP listener")
	}
}

func (s *Service) GetHostPort() (string, string, error) {
	return net.SplitHostPort(s.listener.Addr().String())
}

func addSwagger(basePath string, serviceContainer *restful.Container,
) *restful.WebService {
	swaggerConfig := restfulspec.Config{
		WebServices: serviceContainer.RegisteredWebServices(),
		APIPath:     basePath + "/apidocs/api.json",
		PostBuildSwaggerObjectHandler: func(s *spec.Swagger) {
			s.Info = &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       "Laksono API Service",
					Description: "API for Laksono' practice and playground",
					Version:     "0.0.1",
					Contact: &spec.ContactInfo{
						ContactInfoProps: spec.ContactInfoProps{
							Name:  "Fajar Laksono",
							Email: "fajrlaksono@gmail.com",
							URL:   "https://fajar.laksono.github.io",
						},
					},
				},
			}
			s.SecurityDefinitions = map[string]*spec.SecurityScheme{
				"authorization": spec.APIKeyAuth("Authorization", "header"),
			}
			s.Security = []map[string][]string{
				{"authorization": {}},
			}
		},
	}

	serviceContainer.ServeMux.Handle(basePath+"/apidocs/",
		http.StripPrefix(basePath+"/apidocs", http.FileServer(http.Dir("/srv/contents/swagger-ui"))))

	return restfulspec.NewOpenAPIService(swaggerConfig)
}
