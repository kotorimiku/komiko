package model

type Person struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"unique"`
	Cover       string `json:"cover"`
	Description string `json:"description"`
}
