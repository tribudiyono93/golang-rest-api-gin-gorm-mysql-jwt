package main

import (
	"github.com/gin-gonic/gin"
	"golang-rest-api-gin-gorm-mysql-jwt/config"
	"golang-rest-api-gin-gorm-mysql-jwt/controller"
	"golang-rest-api-gin-gorm-mysql-jwt/middleware"
	"golang-rest-api-gin-gorm-mysql-jwt/repository"
	"golang-rest-api-gin-gorm-mysql-jwt/service"
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

	r.Run()
}