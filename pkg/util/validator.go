package util

import (
	"errors"
	"regexp"

	"github.com/Hilaladiii/aureus/pkg/exception"
	"github.com/google/uuid"
)

func ValidateStrongPassword(pass string) error {
	if len(pass) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(pass) {
		return errors.New("Password must contain at least one lowercase letter")
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(pass) {
		return errors.New("Password must contain at least one uppercase letter")
	}

	if !regexp.MustCompile(`\d`).MatchString(pass) {
		return errors.New("Password must contain at least one number")
	}

	if !regexp.MustCompile(`[@$!%*?&#]`).MatchString(pass) {
		return errors.New("Password must contain at least one special character")
	}

	return nil
}

func ValidateUUID(id string) error {
	if err := uuid.Validate(id); err != nil {
		return exception.NewBadRequestError("invalid uuid")
	}
	return nil
}
