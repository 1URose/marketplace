package dto

type CreateAdRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"image_url" binding:"required,url"`
	Price       int    `json:"price" binding:"required"`
}
