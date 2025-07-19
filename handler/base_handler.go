package handler

import (
	"encoding/json"
	"komiko/repo"
	"komiko/response"
	"komiko/service"
	"komiko/utils"

	"github.com/gin-gonic/gin"
)

type baseHandler[T any, R service.BaseService[T, repo.BaseRepo[T]]] struct {
	service R
}

func NewBaseHandler[T any, R service.BaseService[T, repo.BaseRepo[T]]](service R) *baseHandler[T, R] {
	return &baseHandler[T, R]{
		service: service,
	}
}

func (h *baseHandler[T, R]) GetByID(c *gin.Context) {
	id, err := utils.ParamToUint(c.Param("id"))
	if err != nil {
		response.Error(c, "Invalid ID")
		return
	}
	item, err := h.service.GetByID(id)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *baseHandler[T, R]) GetAll(c *gin.Context) {
	items, err := h.service.GetAll()
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *baseHandler[T, R]) Create(c *gin.Context) {
	var item T
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Error(c, "Invalid input")
		return
	}
	if err := h.service.Create(&item); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *baseHandler[T, R]) Delete(c *gin.Context) {
	var item T
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Error(c, "Invalid input")
		return
	}
	if err := h.service.Delete(&item); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *baseHandler[T, R]) DeleteByID(c *gin.Context) {
	id, err := utils.ParamToUint(c.Param("id"))
	if err != nil {
		response.Error(c, "Invalid ID")
		return
	}
	if err := h.service.DeleteByID(id); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *baseHandler[T, R]) Update(c *gin.Context) {
	var item T
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Error(c, "Invalid input")
		return
	}
	if err := h.service.Update(&item); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *baseHandler[T, R]) UpdateOrCreate(c *gin.Context) {
	var item T
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Error(c, "Invalid input")
		return
	}
	if err := h.service.UpdateOrCreate(&item); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *baseHandler[T, R]) GetList(c *gin.Context) {
	sort := c.Query("sort")
	desc := c.Query("desc")
	limit := c.Query("limit")
	offset := c.Query("offset")
	query := c.Query("query")
	var queryObj *T = new(T)
	if query != "" {
		if err := json.Unmarshal([]byte(query), queryObj); err != nil {
			response.Error(c, "Invalid query")
			return
		}
	} else {
		queryObj = nil
	}
	items, err := h.service.GetList(sort, desc, limit, offset, queryObj)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, items)
}
