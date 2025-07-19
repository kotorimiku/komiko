package repo

import (
	"komiko/model"

	"gorm.io/gorm"
)

type BookRepo struct {
	baseRepo[model.Book]
}

func NewBookRepo(db *gorm.DB) *BookRepo {
	return &BookRepo{baseRepo: baseRepo[model.Book]{db: db}}
}

func (r *BookRepo) GetByPath(path string) (*model.Book, error) {
	return r.GetBy("path", path)
}

func (r *BookRepo) GetBySeriesID(seriesID uint) ([]*model.Book, error) {
	return r.GetList(ListQueryParams{
		Sort:  "number",
		Desc:  false,
		Query: "series_id = ?",
		Args:  []interface{}{seriesID},
	})
}

func (r *BookRepo) GetByLibraryID(libraryID uint) ([]*model.Book, error) {
	books := []*model.Book{}
	err := r.db.Joins("JOIN series ON books.series_id = series.id").
		Where("series.library_id = ?", libraryID).
		Preload("Series").
		Preload("Series.Library").
		Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

// not update cover
func (r *BookRepo) UpdateOrCreateByPath(book *model.Book) error {
	return r.UpdateOrCreateBy(book, "path", []string{"title", "description", "page_count", "pages", "images", "number", "type", "series_id"})
}

func (r *BookRepo) FirstOrCreateByPath(book *model.Book) error {
	return r.FirstOrCreate(book, &model.Book{Path: book.Path})
}

func (r *BookRepo) GetCoversByLibraryID(libraryID uint) ([]string, error) {
	covers := []string{}
	err := r.db.Model(&model.Book{}).
		Joins("JOIN series ON books.series_id = series.id").
		Select("books.cover").
		Where("series.library_id = ?", libraryID).
		Find(&covers).Error
	if err != nil {
		return nil, err
	}
	return covers, nil
}
