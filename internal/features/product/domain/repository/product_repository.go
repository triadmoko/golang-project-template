package repository

import (
	"app/internal/features/product/domain/entity"
	"context"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Product, error)
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*entity.Product, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*entity.Product, error)
}
