package use_cases

import (
	"context"
	"fmt"
	"github.com/1URose/marketplace/internal/user_profile/domain/user/entity"
	"github.com/1URose/marketplace/internal/user_profile/domain/user/repository"
	"log"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {

	log.Println("[usecase:user] initializing UserService")

	return &UserService{Repo: repo}
}

func (us *UserService) GetUserByID(ctx context.Context, id int) (*entity.User, error) {

	log.Printf("[usecase:user] GetUserByID called: id=%d", id)

	user, err := us.Repo.GetUserByID(ctx, id)

	if err != nil {

		log.Printf("[usecase:user][ERROR] GetUserByID: %v", err)

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	if user == nil {

		log.Printf("[usecase:user] GetUserByID: no user found for id=%d", id)

	} else {

		log.Printf("[usecase:user] GetUserByID succeeded: user=%+v", user)
	}

	return user, nil
}

func (us *UserService) GetAllUsers(ctx context.Context) ([]entity.User, error) {

	log.Println("[usecase:user] GetAllUsers called")

	users, err := us.Repo.GetAllUsers(ctx)

	if err != nil {

		log.Printf("[usecase:user][ERROR] GetAllUsers: %v", err)

		return nil, fmt.Errorf("get all users: %w", err)

	}

	log.Printf("[usecase:user] GetAllUsers succeeded: count=%d", len(users))

	return users, nil
}
