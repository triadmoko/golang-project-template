package repository

import (
	"app/internal/shared/domain/entity"
	"app/internal/shared/domain/repository"
	"app/pkg"
	"context"

	"gorm.io/gorm"
)

// userRepository implements repository.UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetByID retrieves a user by ID (UUID string)
func (r *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, filter entity.FilterUser, user *entity.User) error {
	result := r.db.WithContext(ctx).Where("deleted_at IS NULL AND id = ?", filter.ID).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete deletes a user (soft delete)
func (r *userRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.User{}).Error; err != nil {
		return err
	}
	return nil
}

// List retrieves a list of users with pagination and filtering
func (r *userRepository) List(ctx context.Context, filter entity.FilterUser) ([]*entity.User, int, error) {
	// Build scopes for dynamic query construction
	scopes := []func(db *gorm.DB) *gorm.DB{
		// Soft delete filter - only get non-deleted records
		func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL")
		},
	}

	// Basic field filters
	if filter.ID != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("id = ?", filter.ID)
		})
	}
	if filter.Email != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("email = ?", filter.Email)
		})
	}
	if filter.Username != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("username = ?", filter.Username)
		})
	}
	if filter.FirstName != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("first_name LIKE ?", "%"+filter.FirstName+"%")
		})
	}
	if filter.LastName != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("last_name LIKE ?", "%"+filter.LastName+"%")
		})
	}
	if filter.IsActive != nil {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", *filter.IsActive)
		})
	}

	// Extended filters (add these columns to your User entity if needed)
	if filter.Phone != nil {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("phone = ?", *filter.Phone)
		})
	}
	if filter.Status != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", filter.Status)
		})
	}
	if filter.BirthDate != nil {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("birth_date = ?", filter.BirthDate)
		})
	}
	if filter.Gender != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("gender = ?", filter.Gender)
		})
	}
	if filter.Role != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("role = ?", filter.Role)
		})
	}
	if filter.Provider != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("provider = ?", filter.Provider)
		})
	}

	// Array filters for IN queries
	if len(filter.Genders) > 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("gender IN ?", filter.Genders)
		})
	}
	if len(filter.Roles) > 0 {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("role IN ?", filter.Roles)
		})
	}

	// Query with pagination and filters
	var users []*entity.User
	err := r.db.WithContext(ctx).
		Scopes(pkg.Paginate(filter.Offset, filter.PerPage, r.db)).
		Scopes(scopes...).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	var totalRows int64
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Scopes(scopes...).Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	return users, int(totalRows), nil
}
