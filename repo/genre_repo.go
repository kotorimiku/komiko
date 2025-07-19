package repo

import (
	"komiko/model"

	"gorm.io/gorm"
)

type GenreRepo struct {
	baseRepo[model.Genre]
}

func NewGenreRepo(db *gorm.DB) *GenreRepo {
	return &GenreRepo{baseRepo: baseRepo[model.Genre]{db: db}}
}
