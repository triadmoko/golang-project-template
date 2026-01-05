package dto

import (
	"monorepo/libs/errors"
)

// UpdateProfileRequest represents the request for updating user profile
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Validate validates UpdateProfileRequest fields
func (r *UpdateProfileRequest) Validate(lang errors.Lang) map[string][]string {
	errs := make(map[string][]string)

	// At least one field should be provided
	if r.FirstName == "" && r.LastName == "" {
		errs["first_name"] = append(errs["first_name"], errors.GetValidationMessage(errors.Required, lang))
		errs["last_name"] = append(errs["last_name"], errors.GetValidationMessage(errors.Required, lang))
	}

	return errs
}
