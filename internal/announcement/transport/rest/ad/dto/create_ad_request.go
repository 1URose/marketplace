package dto

type CreateAdRequest struct {
	Title       string `json:"title" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=1000"`
	ImageURL    string `json:"image_url" binding:"required,url"`
	Price       int64  `json:"price" binding:"required,gt=0"`
}
