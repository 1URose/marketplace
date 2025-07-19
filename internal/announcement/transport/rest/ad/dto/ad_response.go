package dto

import (
	"github.com/1URose/marketplace/internal/announcement/domain/ad/entity"
)

type AdResponse struct {
	AdBaseResponse
	IsMine bool `json:"is_mine"`
}

func NewAdResponse(ad *entity.Ad, userID int) *AdResponse {
	var isMine bool
	if userID == ad.AuthorID {
		isMine = true
	}
	return &AdResponse{
		AdBaseResponse: NewAdBaseResponse(ad),
		IsMine:         isMine,
	}
}
