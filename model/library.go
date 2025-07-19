package model

type Library struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Path string `json:"path" gorm:"unique;not null"`
	Type string `json:"type" gorm:"check:type IN ('comic', 'novel')"`
}

const (
	Comic = "comic"
	Novel = "novel"
)
