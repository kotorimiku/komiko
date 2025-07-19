package handler

import (
	"komiko/model"
	"komiko/response"
	"komiko/service"
	"komiko/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	baseHandler[model.User, *service.UserService]
}

func NewUserHandler(services *service.Service) *UserHandler {
	return &UserHandler{
		baseHandler: baseHandler[model.User, *service.UserService]{
			service: services.UserService,
		},
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, "Invalid input")
		return
	}

	err := h.service.Register(&user)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, "success")
}

func (h *UserHandler) Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, "Invalid input")
		return
	}

	token, err := h.service.Login(&user)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, token)
}

func (h *UserHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, "Invalid input")
		return
	}
	if err := h.service.CreateUser(userID, &user); err != nil {
		response.Error(c, err.Error())
		return
	}
	user.Password = ""
	response.Success(c, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, "Invalid input")
		return
	}
	if err := h.service.UpdateUser(userID, &user); err != nil {
		response.Error(c, err.Error())
		return
	}
	user.Password = ""
	response.Success(c, user)
}

func (h *UserHandler) DeleteByID(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := utils.ParamToUint(c.Param("id"))
	if err != nil {
		response.Error(c, "Invalid ID")
		return
	}
	if err := h.service.DeleteUserByID(userID, id); err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.service.GetUserDtoByID(userID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.service.GetAllUserDto()
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, users)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := utils.ParamToUint(c.Param("id"))
	if err != nil {
		response.Error(c, "Invalid ID")
		return
	}
	user, err := h.service.GetUserDtoByID(id)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) AllowRegister(c *gin.Context) {
	response.Success(c, h.service.AllowedRegister())
}
