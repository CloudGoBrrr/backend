package binding

type ReqFileChunkedUploadFinish struct {
	Path string `json:"path" binding:"required"`
	Name string `json:"name" binding:"required"`
}
