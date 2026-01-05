package dto

import (
	"fmt"

	"monorepo/libs/errors"
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
func (r *RegisterRequest) Validate(lang errors.Lang) map[string][]string {
	errs := make(map[string][]string)

	// Email validation
	if r.Email == "" {
		errs["email"] = append(errs["email"], fmt.Sprintf(errors.GetValidationMessage(errors.Required, lang), "email"))
	} else if !errors.IsValidEmail(r.Email) {
		errs["email"] = append(errs["email"], errors.GetValidationMessage(errors.InvalidEmail, lang))
	}

	// Username validation
	if r.Username == "" {
		errs["username"] = append(errs["username"], fmt.Sprintf(errors.GetValidationMessage(errors.Required, lang), "username"))
	} else {
		if !errors.MinLength(r.Username, 3) {
			errs["username"] = append(errs["username"], fmt.Sprintf(errors.GetValidationMessage(errors.UsernameTooShort, lang), 3))
		}
		if !errors.MaxLength(r.Username, 20) {
			errs["username"] = append(errs["username"], fmt.Sprintf(errors.GetValidationMessage(errors.UsernameTooLong, lang), 20))
		}
	}

	// Password validation
	if r.Password == "" {
		errs["password"] = append(errs["password"], fmt.Sprintf(errors.GetValidationMessage(errors.Required, lang), "password"))
	} else if !errors.MinLength(r.Password, 6) {
		errs["password"] = append(errs["password"], fmt.Sprintf(errors.GetValidationMessage(errors.PasswordTooShort, lang), 6))
	}

	// FirstName validation
	if r.FirstName == "" {
		errs["first_name"] = append(errs["first_name"], fmt.Sprintf(errors.GetValidationMessage(errors.Required, lang), "first_name"))
	}

	// LastName validation
	if r.LastName == "" {
		errs["last_name"] = append(errs["last_name"], fmt.Sprintf(errors.GetValidationMessage(errors.Required, lang), "last_name"))
	}

	return errs
}

// LoginRequest represents the request for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates LoginRequest fields
func (r *LoginRequest) Validate(lang errors.Lang) map[string][]string {
	errs := make(map[string][]string)

	// Email validation
	if r.Email == "" {
		errs["email"] = append(errs["email"], fmt.Sprintf(errors.GetValidationMessage(errors.Required, lang), "email"))
	} else if !errors.IsValidEmail(r.Email) {
		errs["email"] = append(errs["email"], errors.GetValidationMessage(errors.InvalidEmail, lang))
	}

	// Password validation
	if r.Password == "" {
		errs["password"] = append(errs["password"], fmt.Sprintf(errors.GetValidationMessage(errors.Required, lang), "password"))
	}

	return errs
}

// LoginResponse represents the response for user login
type LoginResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}
