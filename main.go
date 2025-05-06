package main

import (
	"ToDoProject/internal/handler"
	middleware1 "ToDoProject/internal/middleware"
	"ToDoProject/internal/repository"
	"ToDoProject/internal/service"
	jwtPkg "ToDoProject/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var users = map[string]struct {
	Password string
	Role     string
}{
	"admin": {Password: "admin123", Role: "admin"},
	"user1": {Password: "user123", Role: "user"},
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.POST("/token", func(c *gin.Context) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		user, exists := users[input.Username]
		if !exists || user.Password != input.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := jwtPkg.GenerateToken(input.Username, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	r.GET("/admin", middleware1.AuthMiddleware("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome Admin!"})
	})

	r.GET("/user", middleware1.AuthMiddleware("user"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome User!"})
	})

	todoRepo := repository.NewInMemoryToDoRepository()
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)
	todoHandler.RegisterRoutes(r)

	todoStepRepo := repository.NewInMemoryTodoStepRepository()
	todoStepService := service.NewTodoStepService(todoStepRepo)
	todoStepHandler := handler.NewTodoStepHandler(todoStepService)
	todoStepHandler.RegisterRoutes(r)

	err := r.Run(":8080")
	if err != nil {
		return
	}

}
