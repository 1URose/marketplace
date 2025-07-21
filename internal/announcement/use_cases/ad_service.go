package use_cases

import (
	"context"
	"fmt"
	"github.com/1URose/marketplace/internal/announcement/domain/ad/entity"
	"github.com/1URose/marketplace/internal/announcement/domain/ad/repository"
	entityAF "github.com/1URose/marketplace/internal/announcement/domain/ad_filter/entity"
	"github.com/1URose/marketplace/internal/announcement/transport/rest/ad/dto"
	"log"
	"math"
)

type AdService struct {
	adRepo   repository.AdRepository
	pageSize int
}

func NewAdService(adRepo repository.AdRepository, pageSize int) *AdService {
	log.Printf("[usecase:ad] NewAdService initialized: pageSize=%d", pageSize)
	return &AdService{
		adRepo:   adRepo,
		pageSize: pageSize,
	}
}

func (as *AdService) CreateAd(ctx context.Context, userId int, req *dto.CreateAdRequest) (*entity.Ad, error) {
	log.Printf("[usecase:ad] CreateAd called: userId=%d title=%q price=%d", userId, req.Title, req.Price)

	newAd := entity.NewAd(req.Title, req.Description, req.ImageURL, req.Price, userId)
	createdAd, err := as.adRepo.CreateAd(ctx, newAd)
	if err != nil {
		log.Printf("[usecase:ad][ERROR] CreateAd failed: %v", err)
		return nil, err
	}

	log.Printf("[usecase:ad] CreateAd succeeded: adID=%d createdAt=%s", createdAd.ID, createdAd.CreatedAt)
	return createdAd, nil
}

func (as *AdService) GetAllAds(ctx context.Context, req *dto.GetAllAdsRequest) ([]*entity.Ad, int, error) {
	log.Printf("[usecase:ad] GetAllAds called: page=%d sortBy=%s sortOrder=%s minPrice=%v maxPrice=%v",
		req.Page, req.SortBy, req.SortOrder, req.MinPrice, req.MaxPrice,
	)

	filter := entityAF.NewAdFilter(
		req.Page,
		as.pageSize,
		req.SortBy,
		req.SortOrder,
		req.MinPrice,
		req.MaxPrice,
	)

	total, err := as.adRepo.CountAds(ctx)
	if err != nil {
		log.Printf("[usecase:ad][ERROR] CountAds failed: %v", err)
		return nil, 0, fmt.Errorf("count ads: %w", err)
	}

	if total == 0 {
		log.Printf("[usecase:ad] GetAllAds: no ads found")
		return nil, 0, nil
	}

	log.Printf("[usecase:ad] CountAds succeeded: total=%d", total)

	countPages := int(math.Ceil(float64(total) / float64(as.pageSize)))

	if req.Page > countPages {
		log.Printf("[usecase:ad][ERROR] %v", err)
		return nil, 0, fmt.Errorf("invalid page number: %d > %d", req.Page, countPages)
	}

	ads, err := as.adRepo.GetAllAds(ctx, filter)
	if err != nil {
		log.Printf("[usecase:ad][ERROR] GetAllAds failed: %v", err)
		return nil, 0, fmt.Errorf("query ads: %w", err)
	}

	log.Printf("[usecase:ad] GetAllAds succeeded: returned=%d countPages=%d", len(ads), countPages)
	return ads, countPages, nil
}
