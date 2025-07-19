package repo

import (
	"gorm.io/gorm"
)

type Repo struct {
	BookRepo     *BookRepo
	LibraryRepo  *LibraryRepo
	SeriesRepo   *SeriesRepo
	GenreRepo    *GenreRepo
	PersonRepo   *PersonRepo
	UserRepo     *UserRepo
	ProgressRepo *ProgressRepo
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		BookRepo:     NewBookRepo(db),
		LibraryRepo:  NewLibraryRepo(db),
		SeriesRepo:   NewSeriesRepo(db),
		GenreRepo:    NewGenreRepo(db),
		PersonRepo:   NewPersonRepo(db),
		UserRepo:     NewUserRepo(db),
		ProgressRepo: NewProgressRepo(db),
	}
}
