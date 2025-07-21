package dto

import "github.com/1URose/marketplace/internal/announcement/domain/ad/entity"

type GetAllAdsResponse struct {
	Ads        []*AdResponse `json:"ads"`
	CountPages int           `json:"count_pages"`
}

func NewGetAllAdsResponse(ads []*entity.Ad, userId, countPages int) GetAllAdsResponse {
	adsResponse := make([]*AdResponse, 0, len(ads))

	for _, ad := range ads {
		adsResponse = append(adsResponse, NewAdResponse(ad, userId))
	}
	return GetAllAdsResponse{
		Ads:        adsResponse,
		CountPages: countPages,
	}
}
