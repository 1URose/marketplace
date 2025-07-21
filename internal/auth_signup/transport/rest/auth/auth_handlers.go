package auth

import (
	"github.com/1URose/marketplace/internal/auth_signup/domain/redis/entity"
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest/auth/dto"
	"github.com/1URose/marketplace/internal/auth_signup/use_cases"
	dtoErr "github.com/1URose/marketplace/internal/common/transport/rest/dto"

	"github.com/1URose/marketplace/internal/common/jwt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"

	"log"
	"net/http"
)

type Handler struct {
	AuthService *use_cases.AuthService
	JWTManager  *jwt.Manager
}

func NewAuthHandler(authService *use_cases.AuthService, jwtManager *jwt.Manager) *Handler {
	log.Println("[handler:auth] NewAuthHandler initialized")

	return &Handler{
		AuthService: authService,
		JWTManager:  jwtManager,
	}
}

// SignUp godoc
// @Summary      Регистрация нового пользователя
// @Description  Создать обычного пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user_input  body      dto.SignUpRequest     true  "Регистрационные данные"
// @Success      201         {object}  dto.SignUpResponse   "Успешная регистрация"
// @Failure      400         {object}  dto.ErrorResponse    "Ошибка валидации"
// @Failure      500         {object}  dto.ErrorResponse    "Внутренняя ошибка сервера"
// @Router       /auth/signup [post]
func (ah *Handler) SignUp(ctx *gin.Context) {

	log.Printf("[handler:auth] SignUpRequest called")

	var signUpReq dto.SignUpRequest

	if err := ctx.ShouldBindJSON(&signUpReq); err != nil {
		log.Printf("[handler:auth][ERROR] bind body: %v", err)

		ctx.JSON(http.StatusBadRequest, dtoErr.ErrorResponse{
			Error:  "Invalid request body",
			Detail: err.Error(),
		})
		return
	}

	createUserReq, err := ah.AuthService.SingUp(ctx, signUpReq)
	if err != nil {

		log.Printf("[handler:auth][ERROR] SingUp failed: %v", err)

		ctx.JSON(http.StatusInternalServerError, dtoErr.ErrorResponse{
			Error: "Failed to create user",
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

	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		log.Printf("[handler:auth][ERROR] bind body: %v", err)

		ctx.JSON(http.StatusBadRequest, dtoErr.ErrorResponse{
			Error:  "Invalid request body",
			Detail: err.Error(),
		})
		return
	}

	existsUser, err := ah.AuthService.Login(ctx, loginReq)

	if err != nil {

		log.Printf("[handler:auth][ERROR] Email failed: %v", err)

		ctx.JSON(http.StatusInternalServerError, dtoErr.ErrorResponse{
			Error:  "Failed to login",
			Detail: err.Error(),
		})

		return
	}

	if existsUser == nil {

		log.Printf("[handler:auth] invalid credentials for %q", loginReq.Email)

		ctx.JSON(http.StatusUnauthorized, dtoErr.ErrorResponse{
			Error: "Invalid credentials",
		})

		return
	}

	log.Printf("[handler:auth] authentication succeeded for %q", loginReq.Email)

	accessToken, err := ah.JWTManager.GenerateAccessToken(existsUser.Email, existsUser.ID)

	if err != nil {

		log.Printf("[handler:auth][ERROR] GenerateAccessToken failed: %v", err)

		ctx.JSON(http.StatusInternalServerError, dtoErr.ErrorResponse{
			Error: "Failed to generate access token",
		})

		return
	}

	refreshToken, err := ah.JWTManager.GenerateRefreshToken(existsUser.Email, existsUser.ID)

	if err != nil {

		log.Printf("[handler:auth][ERROR] GenerateRefreshToken failed: %v", err)

		ctx.JSON(http.StatusInternalServerError, dtoErr.ErrorResponse{
			Error: "Failed to generate refresh token",
		})

		return
	}

	session := entity.NewSession(loginReq.Email, refreshToken)

	if err = ah.AuthService.RedisRepo.Set(ctx, session); err != nil {

		log.Printf("[handler:auth][ERROR] Redis Set failed: %v", err)

		ctx.JSON(http.StatusInternalServerError, dtoErr.ErrorResponse{
			Error: "Failed to set session",
		})

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
// @Success      200             {object}  dto.LoginResponse      "Новые access и refresh токены"
// @Success      200             {object}  dto.StillValidResponse "Access токен ещё действует"
// @Failure      401             {object}  dto.ErrorResponse      "Invalid or missing token"
// @Failure      500             {object}  dto.ErrorResponse      "Server error"
// @Router       /auth/refresh   [post]
func (ah *Handler) Refresh(ctx *gin.Context) {
	log.Printf("[handler:auth] Refresh called")

	if authH := ctx.GetHeader("Authorization"); authH != "" {
		parts := strings.SplitN(authH, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			if _, err := ah.JWTManager.ValidateAccessToken(parts[1]); err == nil {
				log.Printf("[handler:auth] Refresh: access token still valid")
				ctx.JSON(http.StatusOK, dto.StillValidResponse{
					StillValid: true,
					Detail:     "Access token is still valid",
				})
				return
			}
		}
	}

	raw := ctx.GetHeader("X-Refresh-Token")
	if raw == "" {
		log.Printf("[handler:auth][ERROR] Missing refresh token")
		ctx.JSON(http.StatusUnauthorized, dtoErr.ErrorResponse{
			Error: "Missing refresh token",
		})
		return
	}
	parts := strings.SplitN(raw, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		log.Printf("[handler:auth][ERROR] Invalid refresh token format: %q", raw)
		ctx.JSON(http.StatusUnauthorized, dtoErr.ErrorResponse{
			Error: "Invalid refresh token",
		})
		return
	}
	refreshToken := parts[1]

	claims, err := ah.JWTManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		log.Printf("[handler:auth][ERROR] ValidateRefreshToken failed: %v", err)
		ctx.JSON(http.StatusUnauthorized, dtoErr.ErrorResponse{
			Error: "Invalid refresh token",
		})
		return
	}
	log.Printf("[handler:auth] Refresh: valid refresh token for email=%s", claims.Email)

	session, err := ah.AuthService.RedisRepo.Get(ctx, claims.Email)
	if err != nil {
		log.Printf("[handler:auth][ERROR] Redis Get failed for email=%s: %v", claims.Email, err)
		ctx.JSON(http.StatusUnauthorized, dtoErr.ErrorResponse{
			Error: "Session not found or token mismatch",
		})
		return
	}
	if session == nil || session.RefreshToken != refreshToken {
		log.Printf("[handler:auth][ERROR] Session not found or token mismatch: session=%v", session)
		ctx.JSON(http.StatusUnauthorized, dtoErr.ErrorResponse{
			Error: "Session not found or token mismatch",
		})
		return
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Printf("[handler:auth][ERROR] strconv.Atoi failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, dtoErr.ErrorResponse{
			Error: "Failed to parse user id",
		})
		return
	}
	newAccess, err := ah.JWTManager.GenerateAccessToken(claims.Email, userId)
	if err != nil {
		log.Printf("[handler:auth][ERROR] GenerateAccessToken failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, dtoErr.ErrorResponse{
			Error: "Failed to generate new access token",
		})
		return
	}
	newRefresh, err := ah.JWTManager.GenerateRefreshToken(claims.Email, userId)
	if err != nil {
		log.Printf("[handler:auth][ERROR] GenerateRefreshToken failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, dtoErr.ErrorResponse{
			Error: "Failed to generate new refresh token",
		})
		return
	}

	log.Printf("[handler:auth] Session updated in Redis for email=%s", claims.Email)

	session.RefreshToken = newRefresh
	if err := ah.AuthService.RedisRepo.Set(ctx, session); err != nil {
		log.Printf("[handler:auth][ERROR] Redis Set failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, dtoErr.ErrorResponse{
			Error: "Failed to update session in Redis",
		})
		return
	}

	response := dto.NewLoginResponse(newAccess, newRefresh)

	ctx.JSON(http.StatusOK, response)

	log.Printf("[handler:auth] Refresh succeeded: userID=%d", userId)
}
