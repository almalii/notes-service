package validators

import (
	"github.com/go-playground/validator/v10"
	"net/mail"
	"regexp"
	"unsafe"
)

var (
	minPasswordLength = 6
	passwordRegex     = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()-_=+]+$`)
)

func validatePasswordSecurity(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < minPasswordLength {
		return false
	}

	hasDigit := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasDigit = true
			break
		}
	}

	hasUpper := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
			break
		}
	}

	hasLower := false
	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLower = true
			break
		}
	}

	hasSpecial := false
	for _, char := range password {
		if char == '!' || char == '@' || char == '#' || char == '$' || char == '%' ||
			char == '^' || char == '&' || char == '*' || char == '(' || char == ')' ||
			char == '-' || char == '_' || char == '=' || char == '+' || char == '[' ||
			char == ']' || char == '{' || char == '}' || char == '|' || char == ':' ||
			char == ';' || char == '"' || char == '\'' || char == '<' || char == '>' ||
			char == ',' || char == '.' || char == '?' || char == '/' || char == '`' ||
			char == '~' {
			hasSpecial = true
			break
		}
	}

	if !hasDigit || !hasUpper || !hasLower || !hasSpecial {
		return false
	}

	if !passwordRegex.MatchString(password) {
		return false
	}

	return true
}

func validateBytesize(fl validator.FieldLevel) bool {
	fieldValue := fl.Field()

	size := int64(unsafe.Sizeof(fieldValue.String()))
	maxSize := int64(30000000) // 30 mb

	return size <= maxSize
}

func validateEmailRFC(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	_, err := mail.ParseAddress(email)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(email) && err == nil
}

func RegisterCustomValidation(v *validator.Validate) {
	_ = v.RegisterValidation("bytesize", validateBytesize)
	_ = v.RegisterValidation("emailRFC", validateEmailRFC)
	_ = v.RegisterValidation("security", validatePasswordSecurity)
}
