package handler

import (
	"ToDoProject/internal/model"
	"ToDoProject/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type todoStepHandler struct {
	service service.TodoStepService
}

func NewTodoStepHandler(service service.TodoStepService) *todoStepHandler {
	return &todoStepHandler{service: service}
}

func (h *todoStepHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/todo_steps")
	{
		group.GET("/:todoId", h.GetAllStepsForTodo)
		group.POST("/", h.CreateStep)
		group.PUT("/:id", h.UpdateStep)
		group.DELETE("/:id", h.DeleteStep)
	}
}

func (h *todoStepHandler) GetAllStepsForTodo(c *gin.Context) {
	todoIdStr := c.Param("todoId")
	todoId, err := strconv.Atoi(todoIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todoId"})
		return
	}

	steps, err := h.service.GetAllStepsForTodo(todoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve steps"})
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

	step.CreatedAt = time.Now()
	step.UpdatedAt = time.Now()
	createdStep, err := h.service.CreateStep(step)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create step"})
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
