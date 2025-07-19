package repo

import (
	"komiko/model"

	"gorm.io/gorm"
)

type SeriesRepo struct {
	baseRepo[model.Series]
}

func NewSeriesRepo(db *gorm.DB) *SeriesRepo {
	return &SeriesRepo{baseRepo: baseRepo[model.Series]{db: db}}
}

func (r *SeriesRepo) FirstOrCreateByDir(series *model.Series) error {
	return r.FirstOrCreate(series, &model.Series{Dir: series.Dir})
}

// not update cover
func (r *SeriesRepo) UpdateOrCreateByDir(series *model.Series) error {
	return r.UpdateOrCreateBy(series, "dir", []string{"title", "description", "library_id"})
}

func (r *SeriesRepo) GetByLibraryID(libraryID uint) ([]*model.Series, error) {
	return r.GetList(ListQueryParams{
		Query: "library_id = ?",
		Args:  []interface{}{libraryID},
	})
}
