package json

type ResError struct {
	Error string `json:"error"`
}

var ResErrorInvalidRequest ResError = ResError{Error: "invalid request"}
var ResErrorInternalServerError ResError = ResError{Error: "internal server error"}
