package entity

import "time"

type Ad struct {
	ID          int
	Title       string
	Description string
	ImageURL    string
	Price       int
	AuthorID    int
	AuthorEmail string
	CreatedAt   time.Time
}

func NewAd(title, description, imageURL string, price int, authorID int) *Ad {
	return &Ad{
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
		Price:       price,
		AuthorID:    authorID,
	}
}
