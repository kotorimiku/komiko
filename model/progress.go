package model

import "time"

type Progress struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"userId"`
	BookID     uint      `json:"bookId" gorm:"uniqueIndex"`
	SeriesID   uint      `json:"seriesId"`
	Page       uint      `json:"page"`
	BookScroll string    `json:"bookScroll"`
	Book       *Book     `json:"book" gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
	Series     *Series   `json:"series" gorm:"foreignKey:SeriesID;constraint:OnDelete:CASCADE"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime:true"`
}
