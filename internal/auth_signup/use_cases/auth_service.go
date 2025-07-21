package use_cases

import (
	"context"
	"fmt"
	redisRepo "github.com/1URose/marketplace/internal/auth_signup/domain/redis/repository"
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest/auth/dto"
	"github.com/1URose/marketplace/internal/common/password"
	"github.com/1URose/marketplace/internal/user_profile/domain/user/entity"
	userRepo "github.com/1URose/marketplace/internal/user_profile/domain/user/repository"
	"log"
)

type AuthService struct {
	RedisRepo redisRepo.RedisRepository
	UserRepo  userRepo.UserRepository
}

func NewAccountService(
	redisRepo redisRepo.RedisRepository,
	userRepo userRepo.UserRepository,
) *AuthService {

	log.Println("[auth] initializing AuthService")

	svc := &AuthService{
		RedisRepo: redisRepo,
		UserRepo:  userRepo,
	}

	log.Println("[auth] AuthService initialized")

	return svc
}

func (as *AuthService) SingUp(ctx context.Context, req dto.SignUpRequest) (*entity.User, error) {
	log.Printf("[auth] SingUp called: req=%+v", req)

	exists, err := as.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("[auth][ERROR] get user by email: %v", err)
		return nil, err
	}
	if exists != nil {
		log.Printf("[auth][ERROR] user with email=%s already exists", req.Email)
		return nil, fmt.Errorf("user with email=%s already exists", req.Email)
	}

	passwordHash, err := password.HashPassword(req.Password)
	if err != nil {
		log.Printf("[auth][ERROR] hashing password: %v", err)
		return nil, err
	}

	user := entity.NewUser(req.Email, passwordHash)

	createdUser, err := as.UserRepo.CreateUser(ctx, user)

	log.Printf("[auth] SingUp succesful: user=%+v", createdUser)
	return createdUser, nil
}

func (as *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*entity.User, error) {

	log.Printf("[auth] Email called: req=%+v", req)

	existsUser, err := as.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("[auth][ERROR] get user by email: %v", err)
		return nil, err
	}

	if existsUser == nil {
		log.Printf("[auth][ERROR] user with email=%s not found", req.Email)
		return nil, fmt.Errorf("user with email=%s not found", req.Email)
	}

	ok := password.CheckPasswordHash(req.Password, existsUser.PasswordHash)
	if !ok {
		log.Printf("[auth][ERROR] password verification failed")
		return nil, fmt.Errorf("password verification failed")
	}

	log.Printf("[auth] Email succesful")

	return existsUser, nil
}
