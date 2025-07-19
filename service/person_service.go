package service

import (
	"komiko/model"
	"komiko/repo"
	"strings"
)

type PersonService struct {
	baseService[model.Person, *repo.PersonRepo]
	repos *repo.Repo
}

func NewPersonService(repos *repo.Repo) *PersonService {
	return &PersonService{baseService: baseService[model.Person, *repo.PersonRepo]{repo: repos.PersonRepo}, repos: repos}
}

func (s *PersonService) ToPersons(str string) ([]*model.Person, error) {
	list := strings.Split(str, ",")
	var persons []*model.Person
	for _, name := range list {
		name = strings.TrimSpace(name)
		person, err := s.repos.PersonRepo.GetBy(model.Person{Name: name})
		if err != nil {
			return nil, err
		}
		if person == nil {
			person = &model.Person{Name: name}
		}
		persons = append(persons, person)
	}

	return persons, nil
}
