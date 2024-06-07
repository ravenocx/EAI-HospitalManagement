package utils

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func NewValidator() *validator.Validate {
	validate := validator.New()

	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
			return true
		}
		return false
	})

	_ = validate.RegisterValidation("img_url", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()

		re := `^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+(?:.jpg|.jpeg|.png)+$`
		regex := regexp.MustCompile(re)

		return regex.MatchString(field)
	})

	_ = validate.RegisterValidation("phone_number", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()

		re := `^\+62`
		regex := regexp.MustCompile(re)

		return regex.MatchString(field)
	})

	_ = validate.RegisterValidation("birth_date", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()

		re := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})?$`
		regex := regexp.MustCompile(re)

		return regex.MatchString(field)
	})

	_ = validate.RegisterValidation("identity_number", func(fl validator.FieldLevel) bool {
		field := fl.Field().Int()

		strNum := strconv.FormatInt(field, 10)

		return len(strNum) == 16
	})

	return validate
}

func ValidatorErrors(err error) map[string]string {
	fields := map[string]string{}

	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}
