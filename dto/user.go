package dto

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Cover    string `json:"cover"`
	Role     string `json:"role"`
}
