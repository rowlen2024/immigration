package repository

import (
	"time"

	"gorm.io/gorm"
)

// CountByModel returns the total count of rows for the given model type.
func CountByModel[T any](db *gorm.DB) (int64, error) {
	var c int64
	err := db.Model(new(T)).Count(&c).Error
	return c, err
}

// CountByModelRange returns the count of rows created between start (inclusive) and end (exclusive).
func CountByModelRange[T any](db *gorm.DB, start, end time.Time) (int64, error) {
	var c int64
	err := db.Model(new(T)).Where("created_at >= ? AND created_at < ?", start, end).Count(&c).Error
	return c, err
}

// PluckUploadsByColumn returns non-empty column values containing /uploads/ references (unscoped).
func PluckUploadsByColumn[T any](db *gorm.DB, column string) ([]string, error) {
	var urls []string
	err := db.Unscoped().Model(new(T)).
		Where(column+" LIKE ?", "%/uploads/%").
		Pluck(column, &urls).Error
	return urls, err
}
