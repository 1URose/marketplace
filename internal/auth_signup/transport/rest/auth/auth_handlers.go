package auth

import (
	"github.com/1URose/marketplace/internal/auth_signup/domain/redis/entity"
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest/auth/dto"
	"github.com/1URose/marketplace/internal/auth_signup/use_cases"
	"github.com/1URose/marketplace/internal/common/jwt"
	"github.com/gin-gonic/gin"
	"strings"

	"log"
	"net/http"
)

type Handler struct {
	AuthService *use_cases.AuthService
}

func NewAuthHandler(authService *use_cases.AuthService) *Handler {
	log.Println("[handler:auth] NewAuthHandler initialized")

	return &Handler{
		AuthService: authService,
	}
}

// SignUp godoc
// @Summary      Регистрация нового пользователя
// @Description  Создать обычного пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user_input  body      dto.SignUpRequest     true  "User registration payload"
// @Success      201         {object}  dto.SignUpResponse   "Успешная регистрация"
// @Failure      400         {object}  dto.ErrorResponse    "Ошибка валидации"
// @Failure      500         {object}  dto.ErrorResponse    "Внутренняя ошибка сервера"
// @Router       /auth/register [post]
func (ah *Handler) SignUp(ctx *gin.Context) {

	log.Printf("[handler:auth] SignUpRequest called")

	var signUpReq dto.SignUpRequest

	if err := ctx.BindJSON(&signUpReq); err != nil {
		log.Printf("[handler:auth][ERROR] bind body: %v", err)

		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	createUserReq, err := ah.AuthService.SingUp(ctx, signUpReq)
	if err != nil {

		log.Printf("[handler:auth][ERROR] SingUp failed: %v", err)

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: err.Error(),
		})

		return
	}

	log.Printf("[handler:auth] registration successful for %q", createUserReq.Email)

	response := dto.NewSignUpResponse(createUserReq)

	ctx.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary      Вход в систему
// @Description  Аутентификация пользователя и возврат токенов JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user_input  body      dto.LoginRequest      true  "User login payload"
// @Success      200         {object}  dto.LoginResponse     "Access и Refresh токены"
// @Failure      400         {object}  dto.ErrorResponse     "Неверный запрос"
// @Failure      401         {object}  dto.ErrorResponse     "Неверные учетные данные"
// @Failure      500         {object}  dto.ErrorResponse     "Внутренняя ошибка сервера"
// @Router       /auth/login [post]
func (ah *Handler) Login(ctx *gin.Context) {

	var loginReq dto.LoginRequest

	if err := ctx.BindJSON(&loginReq); err != nil {
		log.Printf("[handler:auth][ERROR] bind body: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ok, err := ah.AuthService.Login(ctx, loginReq)

	if err != nil {

		log.Printf("[handler:auth][ERROR] Email failed: %v", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})

		return
	}

	if !ok {

		log.Printf("[handler:auth] invalid credentials for %q", loginReq.Email)

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})

		return
	}

	log.Printf("[handler:auth] authentication succeeded for %q", loginReq.Email)

	accessToken, err := jwt.GenerateAccessToken(loginReq.Email)

	if err != nil {

		log.Printf("[handler:auth][ERROR] GenerateAccessToken failed: %v", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})

		return
	}

	refreshToken, err := jwt.GenerateRefreshToken(loginReq.Email)

	if err != nil {

		log.Printf("[handler:auth][ERROR] GenerateRefreshToken failed: %v", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})

		return
	}

	session := entity.NewSession(loginReq.Email, refreshToken)

	if err = ah.AuthService.RedisRepo.Set(ctx, session); err != nil {

		log.Printf("[handler:auth][ERROR] Redis Set failed: %v", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token"})

		return
	}

	response := dto.NewLoginResponse(accessToken, refreshToken)

	ctx.JSON(http.StatusOK, response)
}

// Refresh godoc
// @Summary      Обновление токенов
// @Description  При истечении срока действия access-токена позволяет получить новую пару токенов.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization    header    string                false  "Текущий access-токен: Bearer <token>"
// @Param        X-Refresh-Token  header    string                true   "Refresh-токен: Bearer <token>"
// @Success      200             {object}  dto.LoginResponse     "Новые access и refresh токены"
// @Failure      401             {object}  dto.ErrorResponse     "Invalid or missing token"
// @Failure      500             {object}  dto.ErrorResponse     "Server error"
// @Router       /auth/refresh   [post]
func (ah *Handler) Refresh(c *gin.Context) {
	if authH := c.GetHeader("Authorization"); authH != "" {
		parts := strings.SplitN(authH, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			if _, err := jwt.ValidateAccessToken(parts[1]); err == nil {
				c.JSON(http.StatusOK, gin.H{"status": "access token still valid"})
				return
			}
		}
	}

	raw := c.GetHeader("X-Refresh-Token")
	if raw == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token"})
		return
	}
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token format"})
		return
	}
	refreshToken := parts[1]

	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	session, err := ah.AuthService.RedisRepo.Get(c, claims.Email)
	if err != nil || session == nil || session.RefreshToken != refreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session not found or token mismatch"})
		return
	}

	newAccess, err := jwt.GenerateAccessToken(claims.Email)
	if err != nil {
		log.Printf("[handler:auth][ERROR] GenerateAccessToken failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}
	newRefresh, err := jwt.GenerateRefreshToken(claims.Email)
	if err != nil {
		log.Printf("[handler:auth][ERROR] GenerateRefreshToken failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh token"})
		return
	}

	session.RefreshToken = newRefresh
	if err := ah.AuthService.RedisRepo.Set(c, session); err != nil {
		log.Printf("[handler:auth][ERROR] Redis Set failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
	})
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		accessToken := parts[1]

		claims, err := jwt.ValidateAccessToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired access token"})
			return
		}

		c.Set("userEmail", claims.Email)

		c.Next()
	}
}

func OptionalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if parts := strings.SplitN(auth, " ", 2); len(parts) == 2 && parts[0] == "Bearer" {
			if claims, err := jwt.ValidateAccessToken(parts[1]); err == nil {
				c.Set("userEmail", claims.Email)
			}
		}
		c.Next()
	}
}
