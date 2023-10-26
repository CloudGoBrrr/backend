package utils

import "cloudgobrrr/http/response"

func ConvertErrorsToErrorResponse(err ...error) response.Error {
	var errors []response.ErrorElement
	for _, e := range err {
		errors = append(errors, response.ErrorElement{Message: e.Error()})
	}
	return response.Error{Success: false, Errors: errors}
}

func ConvertStringsToErrorResponse(err ...string) response.Error {
	var errors []response.ErrorElement
	for _, e := range err {
		errors = append(errors, response.ErrorElement{Message: e})
	}
	return response.Error{Success: false, Errors: errors}
}
