package response

type AuthToken struct {
	Token   string `json:"token"`
	Session string `json:"session"`
}
