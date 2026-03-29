package helpers

import (
	"errors"
	"regexp"
)

func ValidatePassword(password string) error {
	var (
		minLength  = 8
		hasUpper   = regexp.MustCompile(`[A-Z]`)
		hasLower   = regexp.MustCompile(`[a-z]`)
		hasNumber  = regexp.MustCompile(`[0-9]`)
		hasSpecial = regexp.MustCompile(`[!@#~$%^&*()_\-+={}\[\]|:;"'<>,.?/]`)
	)

	if len(password) < minLength {
		return errors.New("password must be at least 8 characters long")
	}
	if !hasUpper.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber.MatchString(password) {
		return errors.New("password must contain at least one digit")
	}
	if !hasSpecial.MatchString(password) {
		return errors.New("password must contain at least one special character")
	}
	return nil
}
