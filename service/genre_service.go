package service

import (
	"komiko/model"
	"komiko/repo"
	"strings"
)

type GenreService struct {
	baseService[model.Genre, *repo.GenreRepo]
	repos *repo.Repo
}

func NewGenreService(repos *repo.Repo) *GenreService {
	return &GenreService{baseService: baseService[model.Genre, *repo.GenreRepo]{repo: repos.GenreRepo}, repos: repos}
}

func (s *GenreService) ToGenres(str string) ([]*model.Genre, error) {
	list := strings.Split(str, ",")
	var genres []*model.Genre
	for _, name := range list {
		name = strings.TrimSpace(name)
		genre, err := s.repos.GenreRepo.GetBy(model.Genre{Name: name})
		if err != nil {
			return nil, err
		}
		if genre == nil {
			genre = &model.Genre{Name: name}
		}
		genres = append(genres, genre)
	}

	return genres, nil
}

func (s *GenreService) ToGenresFromList(list []string) ([]*model.Genre, error) {
	var genres []*model.Genre
	for _, name := range list {
		name = strings.TrimSpace(name)
		genre, err := s.repos.GenreRepo.GetBy(model.Genre{Name: name})
		if err != nil {
			return nil, err
		}
		if genre == nil {
			genre = &model.Genre{Name: name}
		}
		genres = append(genres, genre)
	}

	return genres, nil
}
