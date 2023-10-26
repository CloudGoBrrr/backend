package response

import "github.com/oklog/ulid/v2"

// SecurityTokenCreate is the response body for the token creation route
type SecurityTokenCreate struct {
	ID       ulid.ULID `json:"id"`
	Password string    `json:"password"`
}

// SecurityTokenGetAll is the response body for the token get all route
type SecurityTokenGetAll struct {
	Tokens []SecurityTokenElement `json:"tokens"`
}

type SecurityTokenElement struct {
	ID          ulid.ULID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   int64     `json:"createdAt"`
	UpdatedAt   int64     `json:"updatedAt"`
}

type SecuritySessionGetAll struct {
	Sessions []SecuritySessionElement `json:"sessions"`
}

type SecuritySessionElement struct {
	ID          ulid.ULID `json:"id"`
	Description string    `json:"description"`
	Remember    bool      `json:"remember"`
	CreatedAt   int64     `json:"createdAt"`
	UpdatedAt   int64     `json:"updatedAt"`
}
