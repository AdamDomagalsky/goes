package gapi

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) (err error) {
	if err = ValidateString(value, 3, 64); err != nil {
		return
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only letters, numbers, dashes, underscores or dots")
	}
	return nil
}

func ValidateFullname(value string) (err error) {
	if err = ValidateString(value, 3, 64); err != nil {
		return
	}
	if !isValidFullname(value) {
		return fmt.Errorf("must contain only letters, spaces")
	}
	return nil
}

func ValidatePassword(value string) (err error) {
	if err = ValidateString(value, 6, 128); err != nil {
		return
	}
	return nil
}

func ValidateEmail(value string) (err error) {
	if err = ValidateString(value, 3, 256); err != nil {
		return
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}

	return nil
}
