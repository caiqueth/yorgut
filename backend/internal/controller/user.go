package controller

import (
	"net/http"
	"strconv"
	"yorgut/internal/dto/request"
	"yorgut/internal/dto/response"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{
		userService: service,
	}
}

// Create			godoc
// @Summary			Create user
// @Decription		Save user in Db.
// @Param			user body request.CreateUserRequest true "Create User"
// @Produce			application/json
// @Tags			User
// @Success			200 {object} response.Response{}
// @Router			/user [post]
func (uc *UserController) Create(ctx *gin.Context) {
	createUserRequest := request.CreateUserRequest{}

	err := ctx.ShouldBindJSON(&createUserRequest)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	uc.userService.Create(createUserRequest)
	webResponse := response.Response{
		Data: nil,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// FindUserById 		godoc
// @Summary				Get Single user by id.
// @Param				userId path string true "get user by id"
// @Description			Returns the user matching the requested id
// @Produce				application/json
// @Tags				User
// @Success				200 {object} response.FindUserByIdResponse{}
// @Router				/user/{userId} [get]
func (uc *UserController) FindUserById(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		panic(err)
	}

	user := uc.userService.FindById(id)

	webResponse := response.Response{
		Data: user,
	}
	ctx.JSON(http.StatusOK, webResponse)

}

// ListUsers			godoc
// @Summary				Gets a list of users.
// @Description			Returns a list of users
// @Produce				application/json
// @Tags				User
// @Success				200 {object} response.Response{}
// @Router				/user [get]
func (uc *UserController) FindMany(ctx *gin.Context) {
	users := uc.userService.FindMany()

	webResponse := response.Response{
		Data: users,
	}
	ctx.JSON(http.StatusOK, webResponse)

}

// DeleteUser 		godoc
// @Summary			Delete User by id
// @Param			userId path string true "delete user by id"
// @Description		Deletes the user matching the requested id
// @Produce			application/json
// @Tags			User
// @Success			200 {object} response.Response{}
// @Router			/user [delete]
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		panic(err)
	}

	uc.userService.Delete(id)

	webResponse := response.Response{
		Data: nil,
	}
	ctx.JSON(http.StatusOK, webResponse)

}

// UpdateUser 		godoc
// @Summary			Update user
// @Param			userId path string true "update user"
// @Description		Returns null
// @Produce			application/json
// @Tags			User
// @Success			200 {object} response.Response{}
// @Router			/user/{userId} [put]
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	panic("not implemented")
}
