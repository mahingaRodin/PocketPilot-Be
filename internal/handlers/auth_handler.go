package handlers

import (
	"net/http"
	"pocketpilot/internal/models"
	"pocketpilot/internal/services"
	"pocketpilot/internal/utils"

	"github.com/gin-gonic/gin"
)


type AuthHandler struct {
	authService *services.AuthService
}


func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param register body models.RegisterRequest true "Register payload"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/auth/register [post]
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

// @Summary Login user
// @Description Authenticate user and return JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login payload"
// @Success 200 {object} models.AuthResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/auth/login [post]
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

// @Summary Get user profile
// @Description Retrieve authenticated user's profile
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} models.ErrorResponse
// @Router /api/auth/profile [get]
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