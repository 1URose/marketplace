package validator

import (
	"fmt"
	"github.com/1URose/marketplace/internal/announcement/transport/rest/ad/dto"
)

var allowedSortFields = map[string]struct{}{
	"created_at": {},
	"price":      {},
}

var allowedSortOrders = map[string]struct{}{
	"asc":  {},
	"desc": {},
}

func ValidateGetAllAdsRequest(req *dto.GetAllAdsRequest) error {
	if req.Page < 1 {
		return fmt.Errorf("page must be ≥ 1, got %d", req.Page)
	}
	if _, ok := allowedSortFields[req.SortBy]; !ok {
		return fmt.Errorf("unsupported sort_by: %q", req.SortBy)
	}
	if _, ok := allowedSortOrders[req.SortOrder]; !ok {
		return fmt.Errorf("unsupported sort_order: %q", req.SortOrder)
	}
	if req.MinPrice != nil {
		if *req.MinPrice < 0 {
			return fmt.Errorf("min_price must be ≥ 0, got %d", *req.MinPrice)
		}
	}
	if req.MaxPrice != nil {
		if *req.MaxPrice < 0 {
			return fmt.Errorf("max_price must be ≥ 0, got %d", *req.MaxPrice)
		}
	}
	if req.MinPrice != nil && req.MaxPrice != nil {
		if *req.MinPrice > *req.MaxPrice {
			return fmt.Errorf("min_price (%d) cannot be greater than max_price (%d)", *req.MinPrice, *req.MaxPrice)
		}
	}
	return nil
}
