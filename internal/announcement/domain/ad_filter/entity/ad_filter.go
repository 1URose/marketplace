package entity

type AdFilter struct {
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
	MinPrice  *int
	MaxPrice  *int
}

func NewAdFilter(page, pageSize int, sortBy, sortOrder string, minPrice, maxPrice *int) *AdFilter {
	return &AdFilter{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		MinPrice:  minPrice,
		MaxPrice:  maxPrice,
	}
}
