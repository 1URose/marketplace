package user

import (
	"github.com/1URose/marketplace/internal/user_profile/use_cases"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserService *use_cases.UserService
}

func NewUserHandler(userService *use_cases.UserService) *UserHandler {
	log.Println("[handler:user] NewUserHandler initialized")

	return &UserHandler{
		UserService: userService,
	}
}

// GetUserByID godoc
// @Summary Получение пользователя по ID
// @Description Возвращает информацию о пользователе по его ID
// @Tags user
// @Param id path int true "ID пользователя"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user/{id} [get]
func (u *UserHandler) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	log.Printf("[handler:user] GetUserByID called with id=%q", idStr)

	userID, err := strconv.Atoi(idStr)

	if err != nil {

		log.Printf("[handler:user][ERROR] invalid id: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "id": idStr})

		return
	}

	user, err := u.UserService.GetUserByID(ctx, userID)

	if err != nil {

		log.Printf("[handler:user][ERROR] service.GetUserByID: %v", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user", "details": err.Error()})

		return
	}

	if user == nil {

		log.Printf("[handler:user] user not found id=%d", userID)

		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})

		return
	}

	log.Printf("[handler:user] GetUserByID succeeded id=%d", userID)
	ctx.JSON(http.StatusOK, gin.H{"message": "User found", "user": user})
}

// GetAllUsers godoc
// @Summary Получение всех пользователей
// @Description Возвращает список всех пользователей
// @Tags user
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user [get]
func (u *UserHandler) GetAllUsers(ctx *gin.Context) {
	log.Println("[handler:user] GetAllUsers called")

	users, err := u.UserService.GetAllUsers(ctx)

	if err != nil {

		log.Printf("[handler:user][ERROR] service.GetAllUsers: %v", err)

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users", "details": err.Error()})

		return
	}

	log.Printf("[handler:user] GetAllUsers succeeded: count=%d", len(users))

	ctx.JSON(http.StatusOK, gin.H{"message": "All users", "users": users})
}
