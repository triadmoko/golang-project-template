package repository

import (
	"app/internal/features/product/domain/entity"
	"app/internal/features/product/domain/repository"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// productRepository implements repository.ProductRepository interface
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

// Create creates a new product
func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	result := r.db.WithContext(ctx).Create(product)
	if result.Error != nil {
		return fmt.Errorf("failed to create product: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a product by ID
func (r *productRepository) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	var product entity.Product
	result := r.db.WithContext(ctx).First(&product, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, result.Error
	}
	return &product, nil
}

// Update updates a product
func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	result := r.db.WithContext(ctx).Save(product)
	if result.Error != nil {
		return fmt.Errorf("failed to update product: %w", result.Error)
	}
	return nil
}

// Delete deletes a product (soft delete)
func (r *productRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&entity.Product{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}
	return nil
}

// List retrieves a list of products with pagination
func (r *productRepository) List(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	var products []*entity.Product
	result := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list products: %w", result.Error)
	}

	return products, nil
}

// GetByCategory retrieves products by category with pagination
func (r *productRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*entity.Product, error) {
	var products []*entity.Product
	result := r.db.WithContext(ctx).
		Where("category = ?", category).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get products by category: %w", result.Error)
	}

	return products, nil
}

// Search searches products by query with pagination
func (r *productRepository) Search(ctx context.Context, query string, limit, offset int) ([]*entity.Product, error) {
	var products []*entity.Product
	searchTerm := "%" + strings.ToLower(query) + "%"
	result := r.db.WithContext(ctx).
		Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(category) LIKE ?",
			searchTerm, searchTerm, searchTerm).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to search products: %w", result.Error)
	}

	return products, nil
}
