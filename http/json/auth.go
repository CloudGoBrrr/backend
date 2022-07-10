package json

type ReqSessionChangeDescription struct {
	Id             uint   `json:"sessionId" binding:"required"`
	NewDescription string `json:"newDescription" binding:"required"`
}

type ResSessionChangeDescription struct {
	OldDescription string `json:"oldDescription"`
}
