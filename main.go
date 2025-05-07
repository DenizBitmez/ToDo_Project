package main

import (
	"ToDoProject/internal/handler"
	"ToDoProject/internal/model"
	"ToDoProject/internal/repository"
	"ToDoProject/internal/service"
	jwtPkg "ToDoProject/pkg/jwt"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	// CORS ayarlarını güncelle
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:63342", "http://localhost:8080", "http://127.0.0.1:8080"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login",
		})
	})

	r.POST("/token", func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			fmt.Printf("JSON bağlama hatası: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek formatı"})
			return
		}

		fmt.Printf("Giriş denemesi - Kullanıcı: %s\n", input.Username)

		user, exists := users[input.Username]
		if !exists {
			fmt.Printf("Kullanıcı bulunamadı: %s\n", input.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz kullanıcı adı"})
			return
		}

		if user.Password != input.Password {
			fmt.Printf("Hatalı şifre - Kullanıcı: %s\n", input.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Hatalı şifre"})
			return
		}

		token, err := jwtPkg.GenerateToken(input.Username, user.Role)
		if err != nil {
			fmt.Printf("Token oluşturma hatası - Kullanıcı: %s, Hata: %v\n", input.Username, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token oluşturulamadı"})
			return
		}

		fmt.Printf("Giriş başarılı - Kullanıcı: %s, Rol: %s\n", input.Username, user.Role)

		// Token'ı cookie'ye ekle
		c.SetCookie("token", token, 3600, "/", "", false, true)

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{
			"token":    token,
			"role":     user.Role,
			"username": input.Username,
		})
	})

	// Admin dashboard
	r.GET("/admin", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			var err error
			token, err = c.Cookie("token")
			if err != nil {
				fmt.Println("Cookie'den token alınamadı:", err)
				c.Redirect(http.StatusFound, "/login")
				return
			}
		}

		fmt.Println("Token alındı:", token)

		claims, err := jwtPkg.ValidateToken(token)
		if err != nil {
			fmt.Println("Token doğrulanamadı:", err)
			c.Redirect(http.StatusFound, "/login")
			return
		}

		fmt.Println("Token doğrulandı, kullanıcı:", claims.Username, "rol:", claims.Role)

		if claims.Role != "admin" {
			fmt.Println("Yetkisiz rol:", claims.Role)
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.Header("Authorization", "Bearer "+token)

		todoRepo := repository.NewInMemoryToDoRepository()
		todoService := service.NewTodoService(todoRepo)

		todoStepRepo := repository.NewInMemoryTodoStepRepository()
		todoStepService := service.NewTodoStepService(todoStepRepo)

		todos, _ := todoService.GetAll()
		steps, _ := todoStepService.GetAllStepsForTodo()

		c.HTML(http.StatusOK, "admin-dashboard.html", gin.H{
			"Todos": todos,
			"Steps": steps,
		})
	})

	r.GET("/dashboard", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			var err error
			token, err = c.Cookie("token")
			if err != nil {
				fmt.Println("Cookie'den token alınamadı:", err)
				c.Redirect(http.StatusFound, "/login")
				return
			}
		}

		fmt.Println("Token alındı:", token)

		claims, err := jwtPkg.ValidateToken(token)
		if err != nil {
			fmt.Println("Token doğrulanamadı:", err)
			c.Redirect(http.StatusFound, "/login")
			return
		}

		fmt.Println("Token doğrulandı, kullanıcı:", claims.Username, "rol:", claims.Role)

		c.Header("Authorization", "Bearer "+token)

		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"username": claims.Username,
		})
	})

	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "", false, true)
		c.Redirect(http.StatusFound, "/login")
	})

	todoRepo := repository.NewInMemoryToDoRepository()
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)
	todoHandler.RegisterRoutes(r)

	todoStepRepo := repository.NewInMemoryTodoStepRepository()
	todoStepService := service.NewTodoStepService(todoStepRepo)
	todoStepHandler := handler.NewTodoStepHandler(todoStepService)
	todoStepHandler.RegisterRoutes(r)

	r.GET("/todos/:id/steps", func(c *gin.Context) {

		token := c.GetHeader("Authorization")
		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			var err error
			token, err = c.Cookie("token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token bulunamadı"})
				return
			}
		}

		claims, err := jwtPkg.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token"})
			return
		}

		todoId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz todo ID"})
			return
		}

		todo, err := todoService.GetById(todoId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo bulunamadı"})
			return
		}

		if claims.Role != "admin" && todo.Username != claims.Username {
			c.JSON(http.StatusForbidden, gin.H{"error": "Bu todo'ya erişim yetkiniz yok"})
			return
		}

		steps, err := todoStepService.GetAllStepsForTodo()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Adımlar yüklenirken bir hata oluştu"})
			return
		}

		var todoSteps []model.TodoStep
		for _, step := range steps {
			if step.TODOID == todoId {
				todoSteps = append(todoSteps, step)
			}
		}

		c.JSON(http.StatusOK, todoSteps)
	})

	r.POST("/todos/:id/steps", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			var err error
			token, err = c.Cookie("token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token bulunamadı"})
				return
			}
		}

		claims, err := jwtPkg.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token"})
			return
		}

		todoId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz todo ID"})
			return
		}

		todo, err := todoService.GetById(todoId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo bulunamadı"})
			return
		}

		if claims.Role != "admin" && todo.Username != claims.Username {
			c.JSON(http.StatusForbidden, gin.H{"error": "Bu todo'ya erişim yetkiniz yok"})
			return
		}

		var input struct {
			Title string `json:"title" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek formatı"})
			return
		}

		step := model.TodoStep{
			TODOID:    todoId,
			Content:   input.Title,
			Username:  claims.Username,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		createdStep, err := todoStepService.CreateStep(step)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Adım oluşturulurken bir hata oluştu"})
			return
		}

		c.JSON(http.StatusCreated, createdStep)
	})

	r.PUT("/todos/:id/steps/:stepId", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			var err error
			token, err = c.Cookie("token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token bulunamadı"})
				return
			}
		}

		claims, err := jwtPkg.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token"})
			return
		}

		todoId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz todo ID"})
			return
		}

		stepId, err := strconv.Atoi(c.Param("stepId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz step ID"})
			return
		}

		todo, err := todoService.GetById(todoId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo bulunamadı"})
			return
		}

		if claims.Role != "admin" && todo.Username != claims.Username {
			c.JSON(http.StatusForbidden, gin.H{"error": "Bu todo'ya erişim yetkiniz yok"})
			return
		}

		var input struct {
			Title string `json:"title" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek formatı"})
			return
		}

		step := model.TodoStep{
			ID:        stepId,
			TODOID:    todoId,
			Content:   input.Title,
			UpdatedAt: time.Now(),
		}

		updatedStep, err := todoStepService.UpdateStep(step)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Adım güncellenirken bir hata oluştu"})
			return
		}

		c.JSON(http.StatusOK, updatedStep)
	})

	r.DELETE("/todos/:id/steps/:stepId", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			var err error
			token, err = c.Cookie("token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token bulunamadı"})
				return
			}
		}

		claims, err := jwtPkg.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token"})
			return
		}

		todoId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz todo ID"})
			return
		}

		stepId, err := strconv.Atoi(c.Param("stepId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz step ID"})
			return
		}

		todo, err := todoService.GetById(todoId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo bulunamadı"})
			return
		}

		if claims.Role != "admin" && todo.Username != claims.Username {
			c.JSON(http.StatusForbidden, gin.H{"error": "Bu todo'ya erişim yetkiniz yok"})
			return
		}

		err = todoStepService.DeleteStep(stepId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Adım silinirken bir hata oluştu"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Adım başarıyla silindi"})
	})

	r.PUT("/todos/:id/steps/:stepId/toggle", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			var err error
			token, err = c.Cookie("token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token bulunamadı"})
				return
			}
		}

		claims, err := jwtPkg.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token"})
			return
		}

		todoId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz todo ID"})
			return
		}

		stepId, err := strconv.Atoi(c.Param("stepId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz step ID"})
			return
		}

		todo, err := todoService.GetById(todoId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo bulunamadı"})
			return
		}

		if claims.Role != "admin" && todo.Username != claims.Username {
			c.JSON(http.StatusForbidden, gin.H{"error": "Bu todo'ya erişim yetkiniz yok"})
			return
		}

		steps, err := todoStepService.GetAllStepsForTodo()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Adımlar yüklenirken bir hata oluştu"})
			return
		}

		var targetStep *model.TodoStep
		for _, step := range steps {
			if step.ID == stepId && step.TODOID == todoId {
				targetStep = &step
				break
			}
		}

		if targetStep == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Adım bulunamadı"})
			return
		}

		if targetStep.Status == 0 {
			targetStep.Status = 1
		} else {
			targetStep.Status = 0
		}
		targetStep.UpdatedAt = time.Now()

		updatedStep, err := todoStepService.UpdateStep(*targetStep)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Adım durumu değiştirilirken bir hata oluştu"})
			return
		}

		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, updatedStep)
	})

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
