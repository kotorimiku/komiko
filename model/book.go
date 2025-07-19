package model

type Book struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Cover       string    `json:"cover"`
	Description string    `json:"description"`
	Path        string    `json:"path" gorm:"uniqueIndex;not null"`
	PageCount   uint      `json:"pageCount"`
	Pages       []string  `json:"pages" gorm:"type:json;serializer:json"`
	Images      []*Image  `json:"images" gorm:"type:json;serializer:json"`
	Number      float32   `json:"number"`
	Type        string    `json:"type" gorm:"check:type IN ('comic', 'novel')"`
	SeriesID    uint      `json:"seriesId"`
	Series      *Series   `json:"series" gorm:"foreignKey:SeriesID;constraint:OnDelete:CASCADE"`
	Author      []*Person `json:"author" gorm:"many2many:book_author;"`
	Artist      []*Person `json:"artist" gorm:"many2many:book_artist;"`
	Genres      []*Genre  `json:"genre" gorm:"many2many:book_genre;"`
	Publisher   []*Person `json:"publisher" gorm:"many2many:book_publisher;"`
}
