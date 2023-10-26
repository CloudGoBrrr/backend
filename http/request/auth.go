package request

type (
	AuthSignin struct {
		User     string `json:"user" validate:"required"`
		Password string `json:"password" validate:"required"`
		Remember bool   `json:"remember"`
	}

	AuthSignup struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	AuthSignout struct {
		Session string `json:"session" validate:"required"`
	}

	AuthRefresh struct {
		Session string `json:"session" validate:"required"`
	}
)
