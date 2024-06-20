package valid

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	validateUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	validateFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLen, maxLen int) error {
	n := len(value)
	if n < minLen || n > maxLen {
		return fmt.Errorf("should be %d-%d long string", minLen, maxLen)
	}
	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 25); err != nil {
		return err
	}
	if !validateUsername(username) {
		return fmt.Errorf("can only contain letters, numbers, and underscores")
	}
	return nil
}

func ValidateFullname(fullname string) error {
	if err := ValidateString(fullname, 3, 40); err != nil {
		return err
	}
	if !validateFullname(fullname) {
		return fmt.Errorf("can only contain letters and spaces")
	}
	return nil
}

func ValidatePassword(password string) error {
	return ValidateString(password, 6, 100)
}

func ValidateEmail(email string) error {
	if err := ValidateString(email, 3, 320); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}
	return nil
}
