package dto

type AdResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Price       int64  `json:"price"`
	AuthorID    int    `json:"author_id"`
	CreatedAt   string `json:"created_at"`
	IsMine      bool   `json:"is_mine,omitempty"`
}
