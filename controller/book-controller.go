package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang-rest-api-gin-gorm-mysql-jwt/dto"
	"golang-rest-api-gin-gorm-mysql-jwt/entity"
	"golang-rest-api-gin-gorm-mysql-jwt/helper"
	"golang-rest-api-gin-gorm-mysql-jwt/service"
	"net/http"
	"strconv"
)

type BookController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type bookController struct{
	bookService service.BookService
	jwtService service.JWTService
}

func NewBookController(bookService service.BookService, jwtService service.JWTService) BookController {
	return &bookController{
		bookService: bookService,
		jwtService: jwtService,
	}
}

func (c *bookController) All(context *gin.Context) {
	var books []entity.Book = c.bookService.All()
	res := helper.BuildResponse(true, "OK!", books)
	context.JSON(http.StatusOK, res)
}

func (c *bookController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("no param id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var book entity.Book = c.bookService.FindByID(id)
	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("Data not found", "no data with given id", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "OK!", book)
		context.JSON(http.StatusOK, res)
	}
}

func (c *bookController) Insert(context *gin.Context) {
	var bookCreatedDTO dto.BookCreateDTO
	errDTO := context.ShouldBind(&bookCreatedDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		autHeader := context.GetHeader("Authorization")
		userID := c.getUserIDByToken(autHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			bookCreatedDTO.UserID = convertedUserID
		}
		response := c.bookService.Insert(bookCreatedDTO)
		context.JSON(http.StatusCreated, response)
	}

}

func (c *bookController) Update(context *gin.Context) {
	var bookUpdateDTO dto.BookUpdateDTO
	errDTO := context.ShouldBind(bookUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userID, bookUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			bookUpdateDTO.UserID = id
		}

		result := c.bookService.Update(bookUpdateDTO)
		response := helper.BuildResponse(true, "OK!", result)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *bookController) Delete(context *gin.Context) {
	var book entity.Book
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to get id", "No param id were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	book.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.bookService.IsAllowedToEdit(userID, book.ID) {
		c.bookService.Delete(book)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *bookController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}