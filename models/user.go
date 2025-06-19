package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"-"` 
}

type UserInput struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
}
