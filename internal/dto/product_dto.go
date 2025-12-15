package dto

import "time"

// CreateProductRequest is the partial product data for creation
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=3"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}

// UpdateProductRequest is the partial product data for updates
type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" binding:"omitempty,min=3"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" binding:"omitempty,gt=0"`
}

// ProductResponse is the full product data returned to clients
type ProductResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
