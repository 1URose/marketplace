package repository

import (
	"context"
	"github.com/1URose/marketplace/internal/user_profile/domain/user_order/entity"
)

type UserOrderRepository interface {
	CreateUserOrder(ctx context.Context, userOrder *entity.UserOrder) error
	GetUserOrder(ctx context.Context, userID, orderID int) (*entity.UserOrder, error)
	GetUserOrders(ctx context.Context, userID int) ([]int, error)
}
