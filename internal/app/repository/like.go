package repository

import (
	"patients_categories/internal/app/ds"
)

func (r *Repository) GetLikesCounts() (map[int]int, error) {
	var rows []struct {
		CalculationID uint
		Cnt           int
	}
	err := r.db.Model(&ds.Like{}).
		Select("calculation_id, COUNT(*) AS cnt").
		Group("calculation_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	res := make(map[int]int, len(rows))
	for _, row := range rows {
		res[int(row.CalculationID)] = row.Cnt
	}
	return res, nil
}

func (r *Repository) GetLikesCount(id uint) (int, error) {
	var count int64
	err := r.db.Model(&ds.Like{}).Where("calculation_id = ?", id).Count(&count).Error
	return int(count), err
}
