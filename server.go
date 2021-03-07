package main

import (
	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang-rest-api-gin-gorm-mysql-jwt/config"
	"golang-rest-api-gin-gorm-mysql-jwt/controller"
	_ "golang-rest-api-gin-gorm-mysql-jwt/docs"
	"golang-rest-api-gin-gorm-mysql-jwt/middleware"
	"golang-rest-api-gin-gorm-mysql-jwt/repository"
	"golang-rest-api-gin-gorm-mysql-jwt/service"
	"os"
	"strings"
)

var (
	db = config.SetupDatabaseConnection()
	userRepository = repository.NewUserRepository(db)
	bookRepository = repository.NewBookRepository(db)
	jwtService = service.NewJWTService()
	userService = service.NewUserService(userRepository)
	bookService = service.NewBookService(bookRepository)
	authService = service.NewAuthService(userRepository)
	authController = controller.NewAuthController(authService, jwtService)
	userController = controller.NewUserController(userService, jwtService)
	bookController = controller.NewBookController(bookService, jwtService)
)

// @title Swagger Example API Change
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/
// @schemes http https

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host 127.0.0.1:8080
// @BasePath /api/v1
// @query.collection.format multi
func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	bookRoutes := r.Group("api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("/", bookController.All)
		bookRoutes.POST("/", bookController.Insert)
		bookRoutes.GET("/:id", bookController.FindByID)
		bookRoutes.PUT("/:id", bookController.Update)
		bookRoutes.DELETE("/:id", bookController.Delete)
	}

	if !strings.EqualFold(os.Getenv("ENVIRONMENT"), "production") {
		r.GET("/swagger/*any", func(c *gin.Context) {
			httpSwagger.WrapHandler(c.Writer, c.Request)
		})
	}

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}