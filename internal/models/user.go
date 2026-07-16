package models

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Role  Role   `json:"role"`
	Email string `json:"email"`
}
