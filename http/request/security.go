package request

// SecurityTokenCreate is the request body for the token creation route
type SecurityTokenCreate struct {
	Description string `json:"description" validate:"required"`
}
