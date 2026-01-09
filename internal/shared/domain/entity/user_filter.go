package entity

import "time"

// FilterUser represents the filtering options for user queries
type FilterUser struct {
	// Basic filters
	ID        string `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"` // Pointer to distinguish between false and not set

	// Extended filters (add these fields to User entity if needed)
	Phone     *string    `json:"phone,omitempty"`
	Status    string     `json:"status,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Gender    string     `json:"gender,omitempty"`
	Role      string     `json:"role,omitempty"`
	Provider  string     `json:"provider,omitempty"`

	// Array filters for IN queries
	Genders []string `json:"genders,omitempty"`
	Roles   []string `json:"roles,omitempty"`

	// Pagination
	Offset  int `json:"offset"`
	PerPage int `json:"per_page"`
}
