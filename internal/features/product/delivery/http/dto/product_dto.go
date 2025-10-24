package dto

// CreateProductRequest represents the request for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"required,min=0"`
	Category    string  `json:"category" binding:"required"`
}

// UpdateProductRequest represents the request for updating a product
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	Category    string  `json:"category"`
	IsActive    *bool   `json:"is_active"`
}
