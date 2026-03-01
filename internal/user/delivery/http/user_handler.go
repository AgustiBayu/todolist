package http

import (
	"todolist/internal/domain"
	"todolist/internal/helper"

	"github.com/gin-gonic/gin"
)

type UserHandlerImpl struct {
	UserUsecase domain.UserUsecase
}

func NewUserHandler(UserUsecase domain.UserUsecase) *UserHandlerImpl {
	return &UserHandlerImpl{
		UserUsecase: UserUsecase,
	}
}

func (h *UserHandlerImpl) GetProfile(c *gin.Context) {
	response, err := h.UserUsecase.GetProfile(c.Request.Context())
	if err != nil {
		helper.NewHandleError(c, err)
		return
	}
	helper.NewHandleSuccess(c, response)
}

func (h *UserHandlerImpl) Register(c *gin.Context) {
	var req domain.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.NewHandleError(c, err)
		return
	}
	err := h.UserUsecase.Register(c.Request.Context(), req)
	if err != nil {
		helper.NewHandleError(c, err)
		return
	}
	helper.NewHandleSuccess(c, "Register success!")
}
