package v1

import (
	"auth/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.signUp)
		users.POST("/sign-in", h.signIn)
		users.POST("/gpthub-backend/refresh", h.userRefresh)

		authenticated := users.Group("/", h.userIdentity)
		{
			authenticated.GET("/:id", h.getUser)
		}
	}
}

type userSignUpInput struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

// @Summary User SignUp
// @Tags users-gpthub-backend
// @Description create user account
// @ModuleID signUp
// @Accept  json
// @Produce  json
// @Param input body userSignUpInput true "sign up info"
// @Success 201 {string} string "ok"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/users/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input userSignUpInput // maybe change
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Users.SignUp(c, service.UserSignUpInput{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// @Summary User SignIn
// @Tags users-gpthub-backend
// @Description user sign in
// @ModuleID signUp
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign in info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/users/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.services.SignIn(c, service.UserSignInInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

// @Summary User Refresh Tokens
// @Tags users-gpthub-backend
// @Description user refresh tokens
// @Accept  json
// @Produce  json
// @Param input body refreshInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/users/gpthub-backend/refresh [post]
func (h *Handler) userRefresh(c *gin.Context) {
	var input refreshInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.services.Users.RefreshTokens(c, input.Token)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *Handler) getUser(c *gin.Context) {
	logrus.Info("Not implemented yet")
}
