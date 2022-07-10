package binding

type ReqSessionChangeDescription struct {
	Id             uint   `json:"sessionId" binding:"required"`
	NewDescription string `json:"newDescription" binding:"required"`
}

type ResSessionChangeDescription struct {
	OldDescription string `json:"oldDescription"`
}

type ReqSessionCreateBasicAuth struct {
	Description string `json:"description" binding:"required"`
}

type ResSessionCreateBasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResSessionList struct {
	Sessions []map[string]interface{} `json:"sessions"`
}

type ReqSessionDeleteWithID struct {
	ID uint `form:"id" binding:"required"`
}
