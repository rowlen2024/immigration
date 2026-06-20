package repository

import (
	"database/sql"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type PublicVersion struct {
	UpdatedAt time.Time
	Count     int64
}

func (v PublicVersion) String() string {
	return v.UpdatedAt.UTC().Format(time.RFC3339Nano) + ":" + strconv.FormatInt(v.Count, 10)
}

func MergePublicVersions(versions ...PublicVersion) PublicVersion {
	var merged PublicVersion
	for _, v := range versions {
		if v.UpdatedAt.After(merged.UpdatedAt) {
			merged.UpdatedAt = v.UpdatedAt
		}
		merged.Count += v.Count
	}
	return merged
}

type PublicVersionRepo struct {
	db *gorm.DB
}

func (r *PublicVersionRepo) VersionFromQuery(query string, args ...interface{}) (PublicVersion, error) {
	var row struct {
		UpdatedAt sql.NullTime `gorm:"column:updated_at"`
		Count     int64        `gorm:"column:count"`
	}
	if err := r.db.Raw(query, args...).Scan(&row).Error; err != nil {
		return PublicVersion{}, err
	}
	if !row.UpdatedAt.Valid {
		return PublicVersion{Count: row.Count}, nil
	}
	return PublicVersion{UpdatedAt: row.UpdatedAt.Time, Count: row.Count}, nil
}
