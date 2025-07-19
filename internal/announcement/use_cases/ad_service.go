package use_cases

import (
	"context"
	"github.com/1URose/marketplace/internal/announcement/domain/ad/entity"
	"github.com/1URose/marketplace/internal/announcement/domain/ad/repository"
	entityAF "github.com/1URose/marketplace/internal/announcement/domain/ad_filter/entity"
	"github.com/1URose/marketplace/internal/announcement/transport/rest/ad/dto"
	"log"
)

type AdService struct {
	adRepo repository.AdRepository
}

func NewAdService(adRepo repository.AdRepository) *AdService {
	return &AdService{
		adRepo: adRepo,
	}
}

func (as *AdService) CreateAd(ctx context.Context, userId int, ad *dto.CreateAdRequest) (*entity.Ad, error) {
	log.Println("[usecase:ad] CreateAd called")

	newAd := entity.NewAd(ad.Title, ad.Description, ad.ImageURL, ad.Price, userId)

	createdAd, err := as.adRepo.CreateAd(ctx, newAd)

	if err != nil {
		return nil, err
	}

	return createdAd, nil
}

func (as *AdService) GetAllAds(ctx context.Context, req *dto.GetAllAdsRequest) ([]*entity.Ad, error) {
	const pageSize = 10

	adFilter := entityAF.NewAdFilter(req.Page, pageSize, req.SortBy, req.SortOrder, req.MinPrice, req.MaxPrice)

	return as.adRepo.GetAllAds(ctx, &adFilter)
}

func (as *AdService) GetAdByID(ctx context.Context, id int) (*entity.Ad, error) {

	ad, err := as.adRepo.GetAdByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ad == nil {
		return nil, nil
	}
	return ad, nil
}
