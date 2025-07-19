package repo

import (
	"komiko/model"

	"gorm.io/gorm"
)

type PersonRepo struct {
	baseRepo[model.Person]
}

func NewPersonRepo(db *gorm.DB) *PersonRepo {
	return &PersonRepo{baseRepo: baseRepo[model.Person]{db: db}}
}
