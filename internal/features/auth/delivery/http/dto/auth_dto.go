package dto

import (
	"app/internal/shared/constants"
	"fmt"
)

// RegisterRequest represents the request for user registration
type RegisterRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Validate validates RegisterRequest fields
func (r *RegisterRequest) Validate(lang constants.Lang) map[string][]string {
	errors := make(map[string][]string)

	// Email validation
	if r.Email == "" {
		errors["email"] = append(errors["email"], fmt.Sprintf(constants.GetValidationMessage(constants.Required, lang), "email"))
	} else if !constants.IsValidEmail(r.Email) {
		errors["email"] = append(errors["email"], constants.GetValidationMessage(constants.InvalidEmail, lang))
	}

	// Username validation
	if r.Username == "" {
		errors["username"] = append(errors["username"], fmt.Sprintf(constants.GetValidationMessage(constants.Required, lang), "username"))
	} else {
		if !constants.MinLength(r.Username, 3) {
			errors["username"] = append(errors["username"], fmt.Sprintf(constants.GetValidationMessage(constants.UsernameTooShort, lang), 3))
		}
		if !constants.MaxLength(r.Username, 20) {
			errors["username"] = append(errors["username"], fmt.Sprintf(constants.GetValidationMessage(constants.UsernameTooLong, lang), 20))
		}
	}

	// Password validation
	if r.Password == "" {
		errors["password"] = append(errors["password"], fmt.Sprintf(constants.GetValidationMessage(constants.Required, lang), "password"))
	} else if !constants.MinLength(r.Password, 6) {
		errors["password"] = append(errors["password"], fmt.Sprintf(constants.GetValidationMessage(constants.PasswordTooShort, lang), 6))
	}

	// FirstName validation
	if r.FirstName == "" {
		errors["first_name"] = append(errors["first_name"], fmt.Sprintf(constants.GetValidationMessage(constants.Required, lang), "first_name"))
	}

	// LastName validation
	if r.LastName == "" {
		errors["last_name"] = append(errors["last_name"], fmt.Sprintf(constants.GetValidationMessage(constants.Required, lang), "last_name"))
	}

	return errors
}

// LoginRequest represents the request for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates LoginRequest fields
func (r *LoginRequest) Validate(lang constants.Lang) map[string][]string {
	errors := make(map[string][]string)

	// Email validation
	if r.Email == "" {
		errors["email"] = append(errors["email"], fmt.Sprintf(constants.GetValidationMessage(constants.Required, lang), "email"))
	} else if !constants.IsValidEmail(r.Email) {
		errors["email"] = append(errors["email"], constants.GetValidationMessage(constants.InvalidEmail, lang))
	}

	// Password validation
	if r.Password == "" {
		errors["password"] = append(errors["password"], fmt.Sprintf(constants.GetValidationMessage(constants.Required, lang), "password"))
	}

	return errors
}

// LoginResponse represents the response for user login
type LoginResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}
