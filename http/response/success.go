package response

type Success struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
