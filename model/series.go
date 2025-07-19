package model

import "time"

type Series struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title"`
	Cover       string     `json:"cover"`
	Description string     `json:"description"`
	Dir         string     `json:"dir" gorm:"uniqueIndex;not null"`
	Author      []*Person  `json:"author" gorm:"many2many:series_author;"`
	Artist      []*Person  `json:"artist" gorm:"many2many:series_artist;"`
	Genres      []*Genre   `json:"genre" gorm:"many2many:series_genre;"`
	Publisher   []*Person  `json:"publisher" gorm:"many2many:series_publisher;"`
	LibraryID   uint       `json:"libraryId"`
	Library     *Library   `json:"library" gorm:"foreignKey:LibraryID;constraint:OnDelete:CASCADE"`
	CreatedAt   *time.Time `json:"createdAt" gorm:"autoCreateTime:true"`
	UpdatedAt   *time.Time `json:"updatedAt" gorm:"autoUpdateTime:true"`
}
