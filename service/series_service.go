package service

import (
	"fmt"
	"komiko/model"
	"komiko/repo"
)

type SeriesService struct {
	baseService[model.Series, *repo.SeriesRepo]
	repos *repo.Repo
}

func NewSeriesService(repos *repo.Repo) *SeriesService {
	return &SeriesService{baseService: baseService[model.Series, *repo.SeriesRepo]{repo: repos.SeriesRepo}, repos: repos}
}

func (s *SeriesService) Create(series *model.Series) error {
	return s.repo.Create(series)
}

func (s *SeriesService) Save(series *model.Series) error {
	existingSeries, err := s.repo.GetBy("dir", series.Dir)
	if err != nil {
		return fmt.Errorf("get series by dir error: %w", err)
	}
	if existingSeries != nil {
		series.ID = existingSeries.ID
		return s.repo.Update(series)
	}
	return s.repo.Create(series)
}
