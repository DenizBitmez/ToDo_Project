package handler

import (
	"ToDoProject/internal/middleware"
	"ToDoProject/internal/model"
	"ToDoProject/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type todoHandler struct {
	service service.TodoListService
}

func NewTodoHandler(service service.TodoListService) *todoHandler {
	return &todoHandler{service: service}
}

func (h *todoHandler) RegisterRoutes(router *gin.Engine) {
	userGroup := router.Group("/todos", middleware.AuthMiddleware("user"))
	{
		userGroup.GET("/", h.GetAll)
		userGroup.POST("/", h.Create)
		userGroup.GET("/:id", h.GetById)
		userGroup.PUT("/:id", h.Update)
		userGroup.DELETE("/:id", h.Delete)
	}

	adminGroup := router.Group("/admin", middleware.AuthMiddleware("admin"))
	{
		adminGroup.GET("/todos", h.GetAll)
	}
}

func (h *todoHandler) GetAll(c *gin.Context) {
	username := c.GetString("username")
	role := c.GetString("role")
	var todos []model.TodoList
	var err error

	if role == "admin" {
		todos, err = h.service.GetAll()
	} else {
		todos, err = h.service.GetAllByUsername(username)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *todoHandler) GetById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	todo, err := h.service.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *todoHandler) Create(c *gin.Context) {
	var todo model.TodoList
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")

	todo.Username = username

	created := h.service.Create(todo)
	c.JSON(http.StatusCreated, created)
}

func (h *todoHandler) Update(c *gin.Context) {
	var todo model.TodoList
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.Update(todo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update successfully"})
}

func (h *todoHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	err := h.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Delete successfully"})
}
