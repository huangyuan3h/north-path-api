package verify

import (
	"net/mail"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsInLengthRange(text string, minLength int, maxLength int) bool {
	textLength := len(text)
	return textLength >= minLength && textLength <= maxLength
}
