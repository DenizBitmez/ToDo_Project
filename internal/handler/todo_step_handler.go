package handler

import (
	"ToDoProject/internal/middleware"
	"ToDoProject/internal/model"
	"ToDoProject/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type todoStepHandler struct {
	service service.TodoStepService
}

func NewTodoStepHandler(service service.TodoStepService) *todoStepHandler {
	return &todoStepHandler{service: service}
}

func (h *todoStepHandler) RegisterRoutes(r *gin.Engine) {
	userGroup := r.Group("/steps", middleware.AuthMiddleware("user"))
	{
		userGroup.GET("/", h.GetAllSteps)
		userGroup.POST("/", h.CreateStep)
		userGroup.PUT("/", h.UpdateStep)
		userGroup.DELETE("/", h.DeleteStep)
	}

	adminGroup := r.Group("/admin/steps", middleware.AuthMiddleware("admin"))
	{
		adminGroup.GET("/", h.GetAllSteps)
	}
}

func (h *todoStepHandler) GetAllSteps(c *gin.Context) {
	username := c.GetString("username")
	role := c.GetString("role")
	var steps []model.TodoStep
	var err error

	if role == "admin" {
		steps, err = h.service.GetAllStepsForTodo()
	} else {
		steps, err = h.service.GetAllSteps(username)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, steps)
}

func (h *todoStepHandler) CreateStep(c *gin.Context) {
	var step model.TodoStep
	if err := c.ShouldBindJSON(&step); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")
	step.Username = username

	createdStep, err := h.service.CreateStep(step)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdStep)
}

func (h *todoStepHandler) UpdateStep(c *gin.Context) {
	id := c.Param("id")
	var step model.TodoStep
	if err := c.ShouldBindJSON(&step); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	step.ID, _ = strconv.Atoi(id)
	updatedStep, err := h.service.UpdateStep(step)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Step not found"})
		return
	}
	c.JSON(http.StatusOK, updatedStep)
}

func (h *todoStepHandler) DeleteStep(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service.DeleteStep(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Step not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Step deleted successfully"})
}
