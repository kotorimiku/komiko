package handler

import (
	"komiko/model"
	"komiko/service"
)

type SeriesHandler struct {
	baseHandler[model.Series, *service.SeriesService]
}

func NewSeriesHandler(services *service.Service) *SeriesHandler {
	return &SeriesHandler{
		baseHandler: baseHandler[model.Series, *service.SeriesService]{
			service: services.SeriesService,
		},
	}
}
