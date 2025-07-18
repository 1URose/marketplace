package repository

import (
	"context"
	"github.com/1URose/marketplace/internal/announcement/domain/ad/entity"
)

type AdRepository interface {
	CreateAd(ctx context.Context, ad *entity.Ad) (*entity.Ad, error)
	GetAllAds(ctx context.Context) ([]*entity.Ad, error)
	GetAdByID(ctx context.Context, id int) (*entity.Ad, error)
}
