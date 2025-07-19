package model

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Cover    string `json:"cover"`
	Role     string `json:"role" gorm:"check:role IN ('admin', 'user');default:'user'"`
}
