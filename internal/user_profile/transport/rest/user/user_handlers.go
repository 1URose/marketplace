package user

import (
	dtoErr "github.com/1URose/marketplace/internal/common/transport/rest/dto"
	"github.com/1URose/marketplace/internal/user_profile/transport/rest/user/dto"
	"github.com/1URose/marketplace/internal/user_profile/use_cases"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
	UserService *use_cases.UserService
}

func NewUserHandler(userService *use_cases.UserService) *Handler {
	log.Println("[handler:user] NewUserHandler initialized")

	return &Handler{
		UserService: userService,
	}
}

// GetAllUsers godoc
// @Summary Получение всех пользователей
// @Description Возвращает список всех пользователей вместе с их данными(Хэш пароль в частности - для тестирования)
// @Tags user
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user [get]
func (u *Handler) GetAllUsers(ctx *gin.Context) {
	log.Println("[handler:user] GetAllUsers called")

	users, err := u.UserService.GetAllUsers(ctx)

	if err != nil {

		log.Printf("[handler:user][ERROR] service.GetAllUsers: %v", err)

		ctx.JSON(http.StatusInternalServerError, dtoErr.ErrorResponse{
			Error:  "Failed to get all users",
			Detail: err.Error(),
		})

		return
	}
	log.Printf("[handler:user] GetAllUsers succeeded: count=%d", len(users))

	usersResponse := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		currentUser := dto.NewUserResponse(&user)
		currentUser.PasswordHash = user.PasswordHash

		usersResponse = append(usersResponse, currentUser)
	}

	ctx.JSON(http.StatusOK, usersResponse)
}
