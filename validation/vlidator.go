package validation

import (
	"fmt"
	"net/mail"
	"regexp"
)

const (
	USERNAME_MIN_LENGTH = 3
	USERNAME_MAX_LENGTH = 100
	FULLNAME_MIN_LENGTH = 3
	FULLNAME_MAX_LENGTH = 100
	PASWORD_MIN_LENGTH  = 6
	PASSWORD_MAX_LENGTH = 100
	EMAIL_MIN_LENGTH    = 3
	EMAIL_MAX_LENGTH    = 200
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
    isValidFullName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain %d to %d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, USERNAME_MIN_LENGTH, USERNAME_MAX_LENGTH); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lowercase letters, digits or underscore")
	}

	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, FULLNAME_MIN_LENGTH, FULLNAME_MAX_LENGTH); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, PASWORD_MIN_LENGTH, PASSWORD_MAX_LENGTH)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, EMAIL_MIN_LENGTH, EMAIL_MAX_LENGTH); err != nil {
		return err
	}

    if _, err := mail.ParseAddress(value); err != nil {
        return fmt.Errorf("%s is not a valid email address", value)
    }

    return nil
}
