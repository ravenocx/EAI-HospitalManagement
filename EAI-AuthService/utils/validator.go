package utils

import (
	"reflect"
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
			return false
		}
		return true
	})

	_ = validate.RegisterValidation("nip_admin", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		var nip string

		switch field.Kind() {
		case reflect.String:
			nip = field.String()
		case reflect.Int64:
			nip = strconv.FormatInt(field.Int(), 10)
		default:
			return false
		}

		re := `^615[12](200[0-9]|201[0-9]|202[0-4])(0[1-9]|1[0-2])([0-9]{3,5})$`
		regex := regexp.MustCompile(re)

		return regex.MatchString(nip)
	})

	_ = validate.RegisterValidation("nip_nurse", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		var nip string

		switch field.Kind() {
		case reflect.String:
			nip = field.String()
		case reflect.Int64:
			nip = strconv.FormatInt(field.Int(), 10)
		default:
			return false
		}

		re := `^303[12](200[0-9]|201[0-9]|202[0-4])(0[1-9]|1[0-2])([0-9]{3,5})$`
		regex := regexp.MustCompile(re)

		return regex.MatchString(nip)
	})

	_ = validate.RegisterValidation("img_url", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()

		re := `^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+(?:.jpg|.jpeg|.png)+$`
		regex := regexp.MustCompile(re)

		return regex.MatchString(field)
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
