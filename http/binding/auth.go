package binding

type ReqAuthSignin struct {
	Username    string `json:"username"    binding:"required"`
	Password    string `json:"password"    binding:"required"`
	Description string `json:"description" binding:"required"`
}

var ResErrorInvalidLogin ResError = ResError{Error: "invalid username or password"}

type ResAuthSignin struct {
	Token string `json:"token"`
}

type ReqAuthSignup struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ReqAuthChangePassword struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type UserDetails struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	SessionID uint   `json:"sessionId"`
	IsAdmin   bool   `json:"isAdmin"`
}

type ResAuthDetails struct {
	UserDetails UserDetails `json:"userDetails"`
}
