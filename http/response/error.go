package response

type Error struct {
	Success bool           `json:"success"`
	Errors  []ErrorElement `json:"errors"`
}

type ErrorElement struct {
	Message string `json:"msg"`
}
