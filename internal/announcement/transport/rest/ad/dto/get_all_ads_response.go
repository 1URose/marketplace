package dto

import "github.com/1URose/marketplace/internal/announcement/domain/ad/entity"

type GetAllAdsResponse struct {
	Ads        []AdBaseResponse `json:"ads"`
	CountPages int              `json:"count_pages"`
}

func NewGetAllAdsResponse(ads []*entity.Ad, userID, countPages int) *GetAllAdsResponse {
	resp := make([]AdBaseResponse, len(ads))
	for i, a := range ads {
		base := NewAdBaseResponse(a)

		if userID != 0 && a.AuthorID == userID {
			base.IsMine = true
		}
		resp[i] = base
	}
	return &GetAllAdsResponse{Ads: resp, CountPages: countPages}
}
