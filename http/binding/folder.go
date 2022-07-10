package binding

type ReqFolderCreate struct {
	Path string `json:"path" binding:"required"`
	Name string `json:"name" binding:"required"`
}

var ResErrorFolderAlreadyExists ResError = ResError{Error: "folder already exists"}
