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
		adminGroup.POST("/todos/:id/restore", h.Restore)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"})
		return
	}

	todo, err := h.service.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo bulunamadı"})
		return
	}

	username := c.GetString("username")
	if todo.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu todo'yu görüntüleme yetkiniz yok"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *todoHandler) Create(c *gin.Context) {
	var input struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek formatı"})
		return
	}

	username := c.GetString("username")
	todo := model.TodoList{
		Title:    input.Title,
		Username: username,
	}

	created := h.service.Create(todo)
	c.JSON(http.StatusCreated, created)
}

func (h *todoHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"})
		return
	}

	var input struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek formatı"})
		return
	}

	todo, err := h.service.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo bulunamadı"})
		return
	}

	username := c.GetString("username")
	if todo.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu todo'yu düzenleme yetkiniz yok"})
		return
	}

	todo.Title = input.Title
	if err := h.service.Update(*todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Todo güncellenirken bir hata oluştu"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *todoHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"})
		return
	}

	todo, err := h.service.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo bulunamadı"})
		return
	}

	username := c.GetString("username")
	role := c.GetString("role")

	if role != "admin" && todo.Username != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu todo'yu silme yetkiniz yok"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Todo silinirken bir hata oluştu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo başarıyla silindi"})
}

func (h *todoHandler) Restore(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"})
		return
	}

	if err := h.service.Restore(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo başarıyla geri getirildi"})
}
