package usecase

import (
	"app/internal/features/product/domain/entity"
	"app/internal/features/product/domain/repository"
	domainError "app/internal/shared/domain/error"
	"context"
)

// ProductUsecase defines the interface for product use cases
type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*entity.Product, error)
	GetProduct(ctx context.Context, productID string) (*entity.Product, error)
	UpdateProduct(ctx context.Context, productID string, req *UpdateProductRequest) (*entity.Product, error)
	DeleteProduct(ctx context.Context, productID string) error
	GetProducts(ctx context.Context, limit, offset int) ([]*entity.Product, error)
	GetProductsByCategory(ctx context.Context, category string, limit, offset int) ([]*entity.Product, error)
	SearchProducts(ctx context.Context, query string, limit, offset int) ([]*entity.Product, error)
}

// productUsecase implements ProductUsecase interface
type productUsecase struct {
	productRepo repository.ProductRepository
}

// NewProductUsecase creates a new product usecase
func NewProductUsecase(productRepo repository.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

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

// CreateProduct creates a new product
func (p *productUsecase) CreateProduct(ctx context.Context, req *CreateProductRequest) (*entity.Product, error) {
	// Create product entity
	product := entity.NewProduct(req.Name, req.Description, req.Category, req.Price, req.Stock)

	// Save product
	if err := p.productRepo.Create(ctx, product); err != nil {
		return nil, domainError.NewCustomError(500, "failed to create product", err)
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (p *productUsecase) GetProduct(ctx context.Context, productID string) (*entity.Product, error) {
	product, err := p.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, domainError.NewCustomError(404, "product not found", domainError.ErrProductNotFound)
	}

	return product, nil
}

// UpdateProduct updates a product
func (p *productUsecase) UpdateProduct(ctx context.Context, productID string, req *UpdateProductRequest) (*entity.Product, error) {
	product, err := p.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, domainError.NewCustomError(404, "product not found", domainError.ErrProductNotFound)
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price >= 0 {
		product.Price = req.Price
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	// Save updated product
	if err := p.productRepo.Update(ctx, product); err != nil {
		return nil, domainError.NewCustomError(500, "failed to update product", err)
	}

	return product, nil
}

// DeleteProduct deletes a product
func (p *productUsecase) DeleteProduct(ctx context.Context, productID string) error {
	// Check if product exists
	_, err := p.productRepo.GetByID(ctx, productID)
	if err != nil {
		return domainError.NewCustomError(404, "product not found", domainError.ErrProductNotFound)
	}

	// Delete product
	if err := p.productRepo.Delete(ctx, productID); err != nil {
		return domainError.NewCustomError(500, "failed to delete product", err)
	}

	return nil
}

// GetProducts retrieves list of products
func (p *productUsecase) GetProducts(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	products, err := p.productRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, domainError.NewCustomError(500, "failed to get products", err)
	}

	return products, nil
}

// GetProductsByCategory retrieves products by category
func (p *productUsecase) GetProductsByCategory(ctx context.Context, category string, limit, offset int) ([]*entity.Product, error) {
	products, err := p.productRepo.GetByCategory(ctx, category, limit, offset)
	if err != nil {
		return nil, domainError.NewCustomError(500, "failed to get products by category", err)
	}

	return products, nil
}

// SearchProducts searches products by query
func (p *productUsecase) SearchProducts(ctx context.Context, query string, limit, offset int) ([]*entity.Product, error) {
	products, err := p.productRepo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, domainError.NewCustomError(500, "failed to search products", err)
	}

	return products, nil
}
