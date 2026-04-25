package controller

import (
	"strconv"

	"github.com/friedrichad/golang_web_api_demo/internal/common"
	"github.com/friedrichad/golang_web_api_demo/internal/model"
	"github.com/friedrichad/golang_web_api_demo/internal/service"
	"github.com/friedrichad/golang_web_api_demo/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController[model.User]
	userService service.IUserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

func (controller *UserController) GetAllUsers(c *gin.Context) {
	var query model.UserRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		controller.Error(c, common.RequestInvalid, nil)
		return
	}

	users, total, err := controller.userService.GetAllUsers(c, query)
	if err != nil {
		controller.Error(c, err, nil)
		return
	}

	response := model.Page[model.User]{
		Content: users,
		Total:   total,
	}

	controller.Success(c, response)
}

func (controller *UserController) GetUserById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		controller.Error(c, common.RequestInvalid, nil)
		return
	}

	user, errService := controller.userService.GetUserByUuid(c, id)
	if errService != nil {
		controller.Error(c, errService, nil)
		return
	}

	controller.Success(c, user)
}

func (controller *UserController) CreateUser(c *gin.Context) {
	var req model.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		controller.Error(c, common.RequestInvalid, nil)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		controller.Error(c, common.SystemError, nil)
		return
	}

	// Convert DTO to model
	user := &model.User{
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		StatusInt:    req.StatusInt,
	}

	newUser, errService := controller.userService.CreateUser(c, user)
	if errService != nil {
		controller.Error(c, errService, nil)
		return
	}

	controller.Success(c, newUser)
}

func (controller *UserController) UpdateUser(c *gin.Context) {
	var user model.UserUpdate
	if err := c.ShouldBindJSON(&user); err != nil {
		controller.Error(c, common.RequestInvalid, nil)
		return
	}

	errService := controller.userService.UpdateUser(c, &user)
	if errService != nil {
		controller.Error(c, errService, nil)
		return
	}

	controller.Success(c, nil)
}

func (controller *UserController) DeleteUser(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		controller.Error(c, common.RequestInvalid, nil)
		return
	}

	errService := controller.userService.DeleteUser(c, ids)
	if errService != nil {
		controller.Error(c, errService, nil)
		return
	}

	controller.Success(c, nil)
}

func (controller *UserController) GetUserAuthorities(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		controller.Error(c, common.RequestInvalid, nil)
		return
	}

	authorities, errService := controller.userService.GetUserAuthorities(c, id)
	if errService != nil {
		controller.Error(c, errService, nil)
		return
	}

	controller.Success(c, authorities)
}
