package service

import (
	"komiko/model"
	"komiko/repo"
	"komiko/utils"
	"strconv"
)

type ProgressService struct {
	baseService[model.Progress, *repo.ProgressRepo]
	repos *repo.Repo
}

func NewProgressService(repos *repo.Repo) *ProgressService {
	return &ProgressService{baseService: baseService[model.Progress, *repo.ProgressRepo]{repo: repos.ProgressRepo}, repos: repos}
}

func (s *ProgressService) GetListBySeriesID(userID uint, seriesID uint) ([]*model.Progress, error) {
	return s.repo.GetListBySeriesID(seriesID, userID)
}

func (s *ProgressService) GetBookProgresses(userID uint, seriesIdStr string, limit string, offset string) ([]*model.Progress, error) {
	var (
		query  = "user_id = ?"
		params []interface{}
	)
	params = append(params, userID)

	if seriesIdStr != "" {
		seriesID, err := utils.ParamToUint(seriesIdStr)
		if err != nil {
			return nil, err
		}
		query += " AND series_id = ?"
		params = append(params, seriesID)
	}

	var limitInt, offsetInt int
	var err error
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
	}
	if offset != "" {
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			return nil, err
		}
	}

	return s.repo.GetList(repo.ListQueryParams{
		Sort:   "updated_at",
		Desc:   true,
		Limit:  limitInt,
		Offset: offsetInt,
		Query:  query,
		Args:   params,
	})
}

func (s *ProgressService) GetSeriesProgresses(userID uint, libraryIdStr string, limit string, offset string) ([]*model.Progress, error) {
	var (
		libraryID uint
		err       error
	)
	if libraryIdStr != "" {
		libraryID, err = utils.ParamToUint(libraryIdStr)
		if err != nil {
			return nil, err
		}
	}
	var limitInt, offsetInt uint
	if limit != "" {
		limitInt, err = utils.ParamToUint(limit)
		if err != nil {
			return nil, err
		}
	}
	if offset != "" {
		offsetInt, err = utils.ParamToUint(offset)
		if err != nil {
			return nil, err
		}
	}
	return s.repo.GetSeriesProgresses(userID, libraryID, limitInt, offsetInt)
}

func (s *ProgressService) UpdateOrCreateByBookID(progress *model.Progress) error {
	return s.repo.UpdateOrCreateBy(progress, "book_id", []string{})
}
