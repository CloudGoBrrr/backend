package binding

type ReqFileDelete struct {
	Path string `form:"path" binding:"required"`
	Name string `form:"name" binding:"required"`
}
