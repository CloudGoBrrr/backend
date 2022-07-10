package binding

type ReqFileDownloadCreateSecret struct {
	Path string `json:"path" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type ResFileDownloadCreateSecret struct {
	Secret string `json:"secret"`
}

type ReqFileDownloadWithSecret struct {
	Secret string `form:"secret" binding:"required"`
}
