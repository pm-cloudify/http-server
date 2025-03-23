package models

// login request body
// can be from a form or a json body
type LoginRequest struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}
