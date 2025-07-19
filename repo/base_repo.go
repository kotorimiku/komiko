package repo

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseRepo[T any] interface {
	GetByID(uint) (*T, error)
	GetBy(interface{}, ...interface{}) (*T, error)
	GetList(ListQueryParams) ([]*T, error)
	GetAll() ([]*T, error)
	Exists(interface{}, ...interface{}) (bool, error)
	FirstOrCreate(*T, *T) error
	Create(*T) error
	UpdateOrCreate(*T) error
	UpdateOrCreateBy(*T, string, []string) error
	Update(*T) error
	Delete(*T) error
	DeleteByID(uint) error
	Count() (int64, error)
}

type ListQueryParams struct {
	Sort   string
	Desc   bool
	Limit  int
	Offset int
	Query  interface{}
	Args   []interface{}
}

type baseRepo[T any] struct {
	db *gorm.DB
}

func (r *baseRepo[T]) GetBy(query interface{}, args ...interface{}) (*T, error) {
	var obj T
	result := r.db.Where(query, args...).Take(&obj)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &obj, nil
}

func (r *baseRepo[T]) GetList(params ListQueryParams) ([]*T, error) {
	var objs []*T
	db := r.db
	if params.Sort != "" {
		db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: params.Sort}, Desc: params.Desc})
	}
	if params.Query != nil {
		db = db.Where(params.Query, params.Args...)
	}
	if params.Limit > 0 {
		db = db.Limit(params.Limit)
	}
	if params.Offset > 0 {
		db = db.Offset(params.Offset)
	}
	result := db.Find(&objs)
	if result.Error != nil {
		return nil, result.Error
	}
	return objs, nil
}

func (r *baseRepo[T]) GetByID(id uint) (*T, error) {
	return r.GetBy("id", id)
}

func (r *baseRepo[T]) GetAll() ([]*T, error) {
	var items []*T
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *baseRepo[T]) Exists(query interface{}, args ...interface{}) (bool, error) {
	var count int64
	var t *T
	err := r.db.Model(t).Where(query, args...).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *baseRepo[T]) FirstOrCreate(item *T, query *T) error {
	return r.db.Where(query).FirstOrCreate(item).Error
}

func (r *baseRepo[T]) Create(item *T) error {
	return r.db.Create(item).Error
}

func (r *baseRepo[T]) UpdateOrCreate(item *T) error {
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(item).Error
}

func (r *baseRepo[T]) UpdateOrCreateBy(item *T, column string, doUpdates []string) error {
	upDateAll := true
	if len(doUpdates) > 0 {
		upDateAll = false
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: column}},
		DoUpdates: clause.AssignmentColumns(doUpdates),
		UpdateAll: upDateAll,
	}).Create(item).Error
}

func (r *baseRepo[T]) Update(item *T) error {
	return r.db.Model(item).Updates(item).Error
}

func (r *baseRepo[T]) Delete(item *T) error {
	return r.db.Delete(item).Error
}

func (r *baseRepo[T]) DeleteByID(id uint) error {
	var item T
	return r.db.Delete(&item, id).Error
}

func (r *baseRepo[T]) Count() (int64, error) {
	var count int64
	var t *T
	err := r.db.Model(t).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
