package binding

import "cloudgobrrr/backend/pkg/structs"

type ReqFilesList struct {
	Path string `form:"path" binding:"required"`
}

type ResFilesList struct {
	Files []structs.File `json:"files"`
}
