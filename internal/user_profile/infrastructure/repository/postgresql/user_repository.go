package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/1URose/marketplace/internal/common/db/postgresql"
	"github.com/1URose/marketplace/internal/user_profile/domain/user/entity"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

type UserRepository struct {
	Connection *postgresql.Client
}

func NewUserRepository(connection *postgresql.Client) *UserRepository {

	log.Println("[postgresql:user_repo] initializing UserRepository")

	return &UserRepository{
		Connection: connection,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	log.Printf("[postgresql:user_repo] CreateUser called: email=%q", user.Email)

	query := `
        INSERT INTO users (email, password_hash)
        VALUES ($1, $2)
        RETURNING id, created_at
    `
	var userID int

	var createdAt time.Time

	if err := ur.Connection.GetPool().QueryRow(ctx, query, user.Email, user.PasswordHash).
		Scan(&userID, &createdAt); err != nil {

		log.Printf("[postgresql:user_repo][ERROR] insert query failed: %v", err)

		return nil, fmt.Errorf("failed to execute insert query: %w", err)
	}

	log.Printf("[postgresql:user_repo] CreateUser succeeded: id=%d", userID)

	user.ID = userID
	user.CreatedAt = createdAt

	return user, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	log.Printf("[postgresql:user_repo] GetUserByEmail called: email=%q", email)
	return ur.fetchOne(ctx, "email = $1", email)
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	log.Printf("[postgresql:user_repo] GetUserByID called: id=%d", id)
	return ur.fetchOne(ctx, "id = $1", id)
}

func (ur *UserRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	log.Println("[postgresql:user_repo] GetAllUsers called")

	query := `
        SELECT id, email, password_hash, created_at
        FROM users
    `
	rows, err := ur.Connection.GetPool().Query(ctx, query)

	if err != nil {

		log.Printf("[postgresql:user_repo][ERROR] query failed: %v", err)

		return nil, fmt.Errorf("get all users query failed: %w", err)
	}

	defer rows.Close()

	var users []entity.User

	for rows.Next() {

		var u entity.User

		if err = rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {

			log.Printf("[postgresql:user_repo][ERROR] scan failed: %v", err)

			return nil, fmt.Errorf("scan user failed: %w", err)
		}

		users = append(users, u)
	}

	log.Printf("[postgresql:user_repo] GetAllUsers succeeded: count=%d", len(users))

	return users, nil
}

func (ur *UserRepository) UpdateEmail(ctx context.Context, id int, email string) error {
	const query = `
        UPDATE users
        SET email = $1
        WHERE id = $2
    `
	tag, err := ur.Connection.GetPool().Exec(ctx, query, email, id)
	if err != nil {
		return fmt.Errorf("update email failed: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no user to update with id=%d", id)
	}
	return nil
}

func (ur *UserRepository) fetchOne(ctx context.Context, where string, args ...interface{}) (*entity.User, error) {
	const baseQuery = `
        SELECT id, email, password_hash, created_at
        FROM users
        WHERE %s
    `
	query := fmt.Sprintf(baseQuery, where)

	var u entity.User
	err := ur.Connection.GetPool().
		QueryRow(ctx, query, args...).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("[postgresql:user_repo] fetchOne: no rows for %q", where)
		return nil, nil
	}
	if err != nil {
		log.Printf("[postgresql:user_repo][ERROR] fetchOne failed (%s): %v", where, err)
		return nil, fmt.Errorf("fetch user failed: %w", err)
	}
	return &u, nil
}
