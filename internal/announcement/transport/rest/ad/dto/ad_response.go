package dto

import (
	"github.com/1URose/marketplace/internal/announcement/domain/ad/entity"
)

type AdResponse struct {
	AdBaseResponse
	IsMine bool `json:"is_mine"`
}

func NewAdResponse(ad *entity.Ad, userID int) *AdResponse {
	return &AdResponse{
		AdBaseResponse: NewAdBaseResponse(ad),
		IsMine:         userID == ad.AuthorID,
	}
}
