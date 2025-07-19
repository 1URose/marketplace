package auth

import (
	"fmt"
	"github.com/1URose/marketplace/internal/auth_signup/domain/redis/entity"
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest/auth/dto"
	"github.com/1URose/marketplace/internal/auth_signup/use_cases"
	"github.com/1URose/marketplace/internal/common/jwt"
	"github.com/gin-gonic/gin"
	"strconv"
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

	existsUser, err := ah.AuthService.Login(ctx, loginReq)

	if err != nil {

		log.Printf("[handler:auth][ERROR] Email failed: %v", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})

		return
	}

	if existsUser == nil {

		log.Printf("[handler:auth] invalid credentials for %q", loginReq.Email)

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})

		return
	}

	log.Printf("[handler:auth] authentication succeeded for %q", loginReq.Email)

	accessToken, err := jwt.GenerateAccessToken(existsUser.Email, existsUser.ID)

	if err != nil {

		log.Printf("[handler:auth][ERROR] GenerateAccessToken failed: %v", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})

		return
	}

	refreshToken, err := jwt.GenerateRefreshToken(existsUser.Email, existsUser.ID)

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
func (ah *Handler) Refresh(ctx *gin.Context) {
	if authH := ctx.GetHeader("Authorization"); authH != "" {
		parts := strings.SplitN(authH, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			if _, err := jwt.ValidateAccessToken(parts[1]); err == nil {
				ctx.JSON(http.StatusOK, gin.H{"status": "access token still valid"})
				return
			}
		}
	}

	raw := ctx.GetHeader("X-Refresh-Token")
	if raw == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token"})
		return
	}
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token format"})
		return
	}
	refreshToken := parts[1]

	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	session, err := ah.AuthService.RedisRepo.Get(ctx, claims.Email)
	if err != nil || session == nil || session.RefreshToken != refreshToken {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Session not found or token mismatch"})
		return
	}
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Printf("[handler:auth][ERROR] strconv.Atoi failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
		return
	}
	newAccess, err := jwt.GenerateAccessToken(claims.Email, userId)
	if err != nil {
		log.Printf("[handler:auth][ERROR] GenerateAccessToken failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}
	newRefresh, err := jwt.GenerateRefreshToken(claims.Email, userId)
	if err != nil {
		log.Printf("[handler:auth][ERROR] GenerateRefreshToken failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh token"})
		return
	}

	session.RefreshToken = newRefresh
	if err := ah.AuthService.RedisRepo.Set(ctx, session); err != nil {
		log.Printf("[handler:auth][ERROR] Redis Set failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
	})
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = parseAndSetClaims(ctx)
		ctx.Next()
	}
}

func RequireAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := parseAndSetClaims(ctx); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}
		ctx.Next()
	}
}

func parseAndSetClaims(ctx *gin.Context) error {
	auth := ctx.GetHeader("Authorization")
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fmt.Errorf("no bearer token")
	}
	claims, err := jwt.ValidateAccessToken(parts[1])
	if err != nil {
		return err
	}

	//ctx.Set("userEmail", claims.Email)
	if uid, err := strconv.Atoi(claims.Subject); err == nil {
		ctx.Set("userId", uid)
	}
	return nil
}
