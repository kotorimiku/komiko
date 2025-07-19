package repo

import (
	"komiko/model"

	"gorm.io/gorm"
)

type LibraryRepo struct {
	baseRepo[model.Library]
}

func NewLibraryRepo(db *gorm.DB) *LibraryRepo {
	return &LibraryRepo{baseRepo: baseRepo[model.Library]{db: db}}
}
