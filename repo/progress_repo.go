package repo

import (
	"komiko/model"
	"reflect"

	"gorm.io/gorm"
)

type ProgressRepo struct {
	baseRepo[model.Progress]
}

func NewProgressRepo(db *gorm.DB) *ProgressRepo {
	return &ProgressRepo{baseRepo: baseRepo[model.Progress]{db: db}}
}

func (r *ProgressRepo) GetListBySeriesID(userID uint, seriesID uint) ([]*model.Progress, error) {
	return r.GetList(ListQueryParams{
		Sort:  "updated_at",
		Desc:  true,
		Query: "series_id = ? AND user_id = ?",
		Args:  []interface{}{seriesID, userID},
	})
}

func (r *ProgressRepo) GetSeriesProgresses(userID uint, libraryID uint, limit uint, offset uint) ([]*model.Progress, error) {
	var progresses []*model.Progress = make([]*model.Progress, 0)

	progress := model.Progress{}
	structName := reflect.TypeOf(progress).Name()
	tableName := r.db.NamingStrategy.TableName(structName)

	if libraryID == 0 {
		r.db.Raw(`
			SELECT p.*
			FROM `+tableName+` p
			JOIN (
				SELECT series_id, MAX(updated_at) as max_updated
				FROM `+tableName+`
				WHERE user_id = ?
				GROUP BY series_id
			) t ON p.series_id = t.series_id AND p.updated_at = t.max_updated
			WHERE p.user_id = ?
			ORDER BY p.updated_at DESC
			LIMIT ? OFFSET ?
		`, userID, userID, limit, offset).Scan(&progresses)
	} else {
		r.db.Raw(`
			SELECT p.*
			FROM `+tableName+` p
			JOIN (
				SELECT series_id, MAX(updated_at) as max_updated
				FROM `+tableName+`
				WHERE user_id = ? AND library_id = ?
				GROUP BY series_id
			) t ON p.series_id = t.series_id AND p.updated_at = t.max_updated
			WHERE p.user_id = ? AND p.library_id = ?
			ORDER BY p.updated_at DESC
			LIMIT ? OFFSET ?
		`, userID, libraryID, userID, libraryID, limit, offset).Scan(&progresses)
	}

	seriesIDs := make([]uint, 0, len(progresses))
	for _, p := range progresses {
		seriesIDs = append(seriesIDs, p.SeriesID)
	}
	if len(seriesIDs) > 0 {
		var seriesList []*model.Series
		r.db.Where("id IN ?", seriesIDs).Find(&seriesList)
		seriesMap := make(map[uint]*model.Series)
		for _, s := range seriesList {
			seriesMap[s.ID] = s
		}
		for _, p := range progresses {
			p.Series = seriesMap[p.SeriesID]
		}
	}

	return progresses, nil
}
