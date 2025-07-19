package dto

import "github.com/1URose/marketplace/internal/announcement/domain/ad/entity"

type CreateAdResponse struct {
	AdBaseResponse
}

func NewCreateAdResponse(ad *entity.Ad) *CreateAdResponse {
	return &CreateAdResponse{
		AdBaseResponse: NewAdBaseResponse(ad),
	}
}
