package dto

import (
	"github.com/1URose/marketplace/internal/announcement/domain/ad/entity"
	"time"
)

type AdBaseResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Price       int    `json:"price"`
	AuthorID    int    `json:"author_id,omitempty"`
	AuthorEmail string `json:"author_email"`
	CreatedAt   string `json:"created_at"`
	IsMine      bool   `json:"is_mine,omitempty"`
}

func NewAdBaseResponse(ad *entity.Ad) AdBaseResponse {
	return AdBaseResponse{
		ID:          ad.ID,
		Title:       ad.Title,
		Description: ad.Description,
		ImageURL:    ad.ImageURL,
		Price:       ad.Price,
		AuthorEmail: ad.AuthorEmail,
		CreatedAt:   ad.CreatedAt.Format(time.RFC3339),
	}
}
