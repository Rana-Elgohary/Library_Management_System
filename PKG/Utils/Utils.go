package utils

import "regexp"

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ^: Asserts the position at the start of a line. This ensures that the pattern must start matching from the beginning of the string.
// +: This means the local part (before the @) of the email address can contain any combination of these characters. (one or more)
