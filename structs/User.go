package structs

import "github.com/oklog/ulid/v2"

type User struct {
	ID       ulid.ULID
	Username string
	IsAdmin  bool
}
