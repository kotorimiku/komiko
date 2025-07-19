package model

type Image struct {
	FileName   string `json:"fileName"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	PageNumber uint   `json:"pageNumber"`
}
