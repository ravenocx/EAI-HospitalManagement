package utils

import (
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"

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

	_ = validate.RegisterValidation("img_file", func(fl validator.FieldLevel) bool {
		field := fl.Field().Interface()

		fileHeader, ok := field.(multipart.FileHeader)
		if !ok {
			return false
		}

		file, err := fileHeader.Open()
		if err != nil {
			return false
		}
		defer file.Close()

		// Check file extension
		ext := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, ".")+1:])
		if ext != "jpg" && ext != "jpeg" && ext != "png" {
			return false
		}

		return true
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
