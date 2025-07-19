package handler

import (
	"komiko/response"
	"komiko/service"

	"github.com/gin-gonic/gin"
)

func GetTaskHandler(c *gin.Context) {
	id := c.Param("id")
	task := service.GetTaskManager().GetTask(id)
	if task == nil {
		response.Error(c, "task not found")
		return
	}
	response.Success(c, task)
}

func StopTaskHandler(c *gin.Context) {
	id := c.Param("id")
	service.GetTaskManager().StopTask(id)
	response.Success(c, nil)
}

func ListTasksHandler(c *gin.Context) {
	tasks := service.GetTaskManager().ListTasks()
	response.Success(c, tasks)
}
