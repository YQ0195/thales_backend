package main

// Product represents an item in our catalog.
type Product struct {
	ID          int64   `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Type        string  `db:"type" json:"type"`
	Price       float64 `db:"price" json:"price"`
	Description string  `db:"description" json:"description"`
	PictureURL  string  `db:"picture_url" json:"picture_url"`
}
