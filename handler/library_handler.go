package handler

import (
	"komiko/model"
	"komiko/response"
	"komiko/service"
	"komiko/utils"

	"github.com/gin-gonic/gin"
)

type LibraryHandler struct {
	baseHandler[model.Library, *service.LibraryService]
}

func NewLibraryHandler(services *service.Service) *LibraryHandler {
	return &LibraryHandler{
		baseHandler: baseHandler[model.Library, *service.LibraryService]{
			service: services.LibraryService,
		},
	}
}

func (h *LibraryHandler) ScanUpdate(c *gin.Context) {
	id, err := utils.ParamToUint(c.Param("id"))
	if err != nil {
		response.Error(c, "Invalid ID")
		return
	}

	err = h.service.ScanUpdate(id)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, "Scan Update Success")
}

func (h *LibraryHandler) ScanCreate(c *gin.Context) {
	id, err := utils.ParamToUint(c.Param("id"))
	if err != nil {
		response.Error(c, "Invalid ID")
		return
	}

	err = h.service.ScanCreate(id)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, "Scan Create Success")
}

func (h *LibraryHandler) UpdateCover(c *gin.Context) {
	id, err := utils.ParamToUint(c.Param("id"))
	if err != nil {
		response.Error(c, "Invalid ID")
		return
	}

	err = h.service.UpdateCover(id)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, "Update Cover Success")
}
