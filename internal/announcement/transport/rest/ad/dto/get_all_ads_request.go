package dto

type GetAllAdsRequest struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	SortBy    string `form:"sort_by,default=created_at"`
	SortOrder string `form:"sort_order,default=desc" binding:"oneof=asc desc"`
	MinPrice  *int   `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice  *int   `form:"max_price" binding:"omitempty,min=0"`
}
