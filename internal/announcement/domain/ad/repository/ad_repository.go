package repository

import (
	"context"
	"github.com/1URose/marketplace/internal/announcement/domain/ad/entity"
	entityAF "github.com/1URose/marketplace/internal/announcement/domain/ad_filter/entity"
)

type AdRepository interface {
	CreateAd(ctx context.Context, ad *entity.Ad) (*entity.Ad, error)
	GetAllAds(ctx context.Context, filter *entityAF.AdFilter) ([]*entity.Ad, error)
	CountAds(ctx context.Context) (int, error)
}
