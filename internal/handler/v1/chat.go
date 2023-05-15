package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) initChatRoutes(api *gin.RouterGroup) {
	chat := api.Group("/chat", h.userIdentity)
	{
		chat.POST("", h.communicate)
	}
}

type Prompt struct {
	Payload  string `json:"payload" binding:"required"`
	Username string `json:"username" binding:"required,max=32"`
	Email    string `json:"email" binding:"required,email,max=64"`
}

// @Summary Communicate with OpenAI
// @Security UsersAuth
// @Tags chat
// @Description chat with OpenAI Model
// @ModuleID communicate
// @Accept  json
// @Produce  json
// @Param input body Prompt true "prompt to AI"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/chat [post]
func (h *Handler) communicate(c *gin.Context) {
	var prompt Prompt
	if err := c.BindJSON(&prompt); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.services.Communicate(prompt.Payload)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"response": resp,
	})
}
