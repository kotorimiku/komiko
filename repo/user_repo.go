package repo

import (
	"komiko/dto"
	"komiko/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	baseRepo[model.User]
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{baseRepo: baseRepo[model.User]{db: db}}
}

func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	return r.baseRepo.GetBy("username = ?", username)
}

func (r *UserRepo) GetAllDto() ([]*dto.User, error) {
	var items []*dto.User
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *UserRepo) GetDtoByID(id uint) (*dto.User, error) {
	var item *dto.User
	if err := r.db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}
