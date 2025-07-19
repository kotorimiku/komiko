package service

import (
	"errors"
	"komiko/repo"
	"strconv"
)

type BaseService[T any, R repo.BaseRepo[T]] interface {
	Create(item *T) error
	GetByID(id uint) (*T, error)
	GetAll() ([]*T, error)
	Delete(item *T) error
	DeleteByID(id uint) error
	Update(item *T) error
	GetList(sort string, descStr string, limitStr string, offsetStr string, query interface{}, args ...interface{}) ([]*T, error)
	UpdateOrCreate(item *T) error
}

type baseService[T any, R repo.BaseRepo[T]] struct {
	repo R
}

func (s *baseService[T, R]) Create(item *T) error {
	if item == nil {
		return errors.New("item is nil")
	}
	return s.repo.Create(item)
}

func (s *baseService[T, R]) GetByID(id uint) (*T, error) {
	return s.repo.GetByID(id)
}

func (s *baseService[T, R]) GetAll() ([]*T, error) {
	return s.repo.GetAll()
}

func (s *baseService[T, R]) UpdateOrCreate(item *T) error {
	return s.repo.UpdateOrCreate(item)
}

func (s *baseService[T, R]) UpdateOrCreateBy(item *T, column string, doUpdates []string) error {
	return s.repo.UpdateOrCreateBy(item, column, doUpdates)
}

func (s *baseService[T, R]) Delete(item *T) error {
	return s.repo.Delete(item)
}

func (s *baseService[T, R]) DeleteByID(id uint) error {
	return s.repo.DeleteByID(id)
}

func (s *baseService[T, R]) Exists(query interface{}, args ...interface{}) (bool, error) {
	return s.repo.Exists(query, args...)
}

func (s *baseService[T, R]) Update(item *T) error {
	if item == nil {
		return nil
	}
	return s.repo.Update(item)
}

func (s *baseService[T, R]) GetList(sort string, descStr string, limitStr string, offsetStr string, query interface{}, args ...interface{}) ([]*T, error) {
	var desc bool
	if descStr == "true" {
		desc = true
	} else {
		desc = false
	}
	limit := 0
	offset := 0
	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return nil, err
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return nil, err
		}
	}
	return s.repo.GetList(repo.ListQueryParams{
		Sort:   sort,
		Desc:   desc,
		Limit:  limit,
		Offset: offset,
		Query:  query,
		Args:   args,
	})
}
