package binding

import "github.com/gin-gonic/gin"

var ResEmpty gin.H = gin.H{}

type ResError struct {
	Error string `json:"error"`
}

var ResErrorInvalidRequest ResError = ResError{Error: "invalid request"}
var ResErrorInternalServerError ResError = ResError{Error: "internal server error"}
var ResErrorUnauthorized ResError = ResError{Error: "unauthorized"}
var ResErrorInvalidPath ResError = ResError{Error: "invalid path"}
var ResErrorNotFound ResError = ResError{Error: "not found"}
