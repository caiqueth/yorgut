package controller

import (
	"net/http"
	"strconv"
	"yorgut/internal/data/request"
	"yorgut/internal/data/response"
	"yorgut/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{
		userService: service,
	}
}

// CreateUsers		godoc
// @Symmary			Create user
// @Decription		Save user in Db.
// @Param			user body request.CreateUserRequest true "Create User"
// @Produce			application/json
// @Tags			user
// @Success			200 {object} response.Response{}
// @Router			/user [post]
func (controller *UserController) Create(ctx *gin.Context) {
	log.Info().Msg("Create user")
	createUserRequest := request.CreateUserRequest{}
	err := ctx.ShouldBindJSON(&createUserRequest)
	if err != nil {
		panic(err)
	}

	controller.userService.Create(createUserRequest)
	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   nil,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, webResponse)
}

// FindByIdUser 		godoc
// @Summary				Get Single user by id.
// @Param				userId path string true "get user by id"
// @Description			Returns the user matching the requested id
// @Produce				application/json
// @Tags				tags
// @Success				200 {object} response.FindUserByIdResponse{}
// @Router				/user/{userId} [get]
func (controller *UserController) FindUserById(ctx *gin.Context) {
	log.Info().Msg("Find user by id")
	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		panic(err)
	}

	user := controller.userService.FindById(id)

	webResponse := response.Response{
		Code:   http.StatusOK,
		Status: "Ok",
		Data:   user,
	}
	ctx.JSON(http.StatusOK, webResponse)

}
