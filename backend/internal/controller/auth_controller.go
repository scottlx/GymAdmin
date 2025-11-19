package controller

import (
	"gym-admin/internal/config"
	"gym-admin/internal/service"
	"gym-admin/pkg/jwt"
	"gym-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService *service.UserService
}

func NewAuthController() *AuthController {
	return &AuthController{
		userService: service.NewUserService(),
	}
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// Login handles user login
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// TODO: Implement password verification
	// For now, just check if user exists
	user, err := ctrl.userService.GetUserByPhone(req.Phone)
	if err != nil {
		response.Unauthorized(c, "Invalid phone or password")
		return
	}

	// Load config for JWT expiration
	cfg, _ := config.LoadConfig()

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, "user", cfg.JWT.ExpireTime)
	if err != nil {
		response.InternalServerError(c, "Failed to generate token")
		return
	}

	response.Success(c, LoginResponse{
		Token: token,
		User:  user,
	})
}

// Register handles user registration
func (ctrl *AuthController) Register(c *gin.Context) {
	var user struct {
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// TODO: Hash password before storing
	// For now, just create user without password field
	newUser := &struct {
		Name  string
		Phone string
	}{
		Name:  user.Name,
		Phone: user.Phone,
	}

	response.SuccessWithMessage(c, "Registration successful", newUser)
}
