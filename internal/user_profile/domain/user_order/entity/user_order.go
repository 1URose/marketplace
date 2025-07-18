package entity

import (
	"github.com/1URose/marketplace/internal/user_profile/infrastructure/config/kafka/contracts"
	"log"
	"strconv"
)

type UserOrder struct {
	UserId  int `json:"user_id"`
	OrderId int `json:"order_id"`
}

func NewUserOrder(userId, orderId int) *UserOrder {
	return &UserOrder{
		UserId:  userId,
		OrderId: orderId,
	}
}

func NewUserOrderFromKafka(dto *contracts.UserOrderKafka) (*UserOrder, error) {

	log.Printf("[entity:user_order] NewUserOrderFromKafka called: dto=%+v", dto)

	userId, err := strconv.Atoi(dto.ClientId)

	if err != nil {

		log.Printf("[entity:user_order][ERROR] parsing ClientId %q: %v", dto.ClientId, err)

		return nil, err
	}

	orderId, err := strconv.Atoi(dto.OrderId)

	if err != nil {

		log.Printf("[entity:user_order][ERROR] parsing OrderId %q: %v", dto.OrderId, err)

		return nil, err
	}

	log.Printf("[entity:user_order] parsed IDs: userId=%d, orderId=%d", userId, orderId)

	return &UserOrder{
		UserId:  userId,
		OrderId: orderId,
	}, nil
}
