package dto

import (
	"app/internal/shared/constants"
)

// UpdateProfileRequest represents the request for updating user profile
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Validate validates UpdateProfileRequest fields
func (r *UpdateProfileRequest) Validate(lang constants.Lang) map[string][]string {
	errors := make(map[string][]string)

	// At least one field should be provided
	if r.FirstName == "" && r.LastName == "" {
		errors["first_name"] = append(errors["first_name"], constants.GetValidationMessage(constants.Required, lang))
		errors["last_name"] = append(errors["last_name"], constants.GetValidationMessage(constants.Required, lang))
	}

	return errors
}
