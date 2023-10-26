package validator

import (
	"cloudgobrrr/http/response"

	"github.com/go-playground/validator/v10"
)

type (
	Validator struct {
		innerValidator *validator.Validate
	}
	ErrorData struct {
		FailedField string
		Tag         string
		Value       interface{}
	}
)

// Validate validates the data object
func (v Validator) Validate(data interface{}) []ErrorData {
	validationErrors := []ErrorData{}

	errs := v.innerValidator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorData

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

// ConvertToResponse converts the validation errors to a response
func (v Validator) ConvertToResponse(errs []ErrorData) *response.Error {
	res := &response.Error{
		Success: false,
	}

	for _, err := range errs {
		if err.Value.(string) == "" {
			res.Errors = append(res.Errors, response.ErrorElement{
				Message: "field " + err.FailedField + " " + err.Tag + " validation failed",
			})
			continue
		}
		res.Errors = append(res.Errors, response.ErrorElement{
			Message: "field " + err.FailedField + " " + err.Tag + " validation failed on value " + err.Value.(string),
		})
	}

	return res
}

var validate *Validator

func init() {
	validate = &Validator{innerValidator: validator.New()}
}

// Get returns the validator
func Get() *Validator {
	return validate
}
