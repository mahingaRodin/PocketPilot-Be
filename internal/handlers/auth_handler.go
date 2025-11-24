package handlers

import (
	"net/http"
	"pocketpilot-api/internal/models"
	"pocketpilot-api/internal/services"
	"pocketpilot-api/internal/utils"

	"github.com/gin-gonic/gin"
)


type AuthHandler struct {
	authService *services.AuthService
}


func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register (c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request data"))
		return
	}

	authResponse, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	    c.JSON(http.StatusCreated, utils.SuccessResponse("User registered successfully", authResponse))
}


func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request data"))
		return
	}

	authResponse, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error()))
		return
	}

	    c.JSON(http.StatusCreated, utils.SuccessResponse("User Logged in successfully", authResponse))
}


func (h *AuthHandler) GetProfile(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated"))
        return
    }

    user, err := h.authService.GetUserProfile(userID.(string))
    if err != nil {
        c.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
        return
    }

    c.JSON(http.StatusOK, utils.SuccessResponse("Profile retrieved successfully", user))
}