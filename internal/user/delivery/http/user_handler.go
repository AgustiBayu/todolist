package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func (h *UserHandlerImpl) GithubLogin(c *gin.Context) {
	conf := helper.GetGitHubOauthConfig()
	state := os.Getenv("GITHUB_STATE")

	url := conf.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *UserHandlerImpl) GithubCallback(c *gin.Context) {
	stateQuery := c.Query("state")
	if stateQuery != os.Getenv("GITHUB_STATE") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State mismatch"})
		return
	}
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}
	conf := helper.GetGitHubOauthConfig()
	token, err := conf.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Exchange failed"})
		return
	}
	client := conf.Client(c.Request.Context(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed getting user info"})
		return
	}
	defer resp.Body.Close()

	var githubUser struct {
		Name string `json:"name"`
	}
	json.NewDecoder(resp.Body).Decode(&githubUser)
	respEmail, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed getting emails"})
		return
	}
	defer respEmail.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	json.NewDecoder(respEmail.Body).Decode(&emails)
	var primaryEmail string
	for _, e := range emails {
		if e.Primary {
			primaryEmail = e.Email
			break
		}
	}
	if primaryEmail == "" && len(emails) > 0 {
		primaryEmail = emails[0].Email
	}
	userResponse, err := h.UserUsecase.LoginOrRegisterOAuth(c.Request.Context(), primaryEmail, githubUser.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jwtToken, _ := helper.GenereteToken(userResponse.ID)
	helper.NewHandleLoginSuccess(c, userResponse, jwtToken)
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
	err := h.UserUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		helper.NewHandleError(c, err)
		return
	}
	helper.NewHandleSuccess(c, "Register success!")
}

func (h *UserHandlerImpl) Login(c *gin.Context) {
	var req domain.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.NewHandleError(c, err)
		return
	}
	userResponse, err := h.UserUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		helper.NewHandleError(c, err)
		return
	}
	token, err := helper.GenereteToken(userResponse.ID)
	fmt.Println("TOKEN LOGIN:", token)
	if err != nil {
		helper.NewHandleError(c, err)
		return
	}
	helper.NewHandleLoginSuccess(c, userResponse, token)
}
