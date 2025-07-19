package utils

import (
	"komiko/model"
	"strings"
)

func ToPersons(str string) []*model.Person {
	names := strings.Split(str, ",")
	var persons []*model.Person
	for _, name := range names {
		persons = append(persons, &model.Person{Name: name})
	}
	return persons
}

func ToGenres(str string) []*model.Genre {
	names := strings.Split(str, ",")
	var genres []*model.Genre
	for _, name := range names {
		genres = append(genres, &model.Genre{Name: name})
	}
	return genres
}

func ToGenresFromList(list []string) []*model.Genre {
	var genres []*model.Genre
	for _, name := range list {
		genres = append(genres, &model.Genre{Name: name})
	}
	return genres
}
