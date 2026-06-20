package service

import (
	"strings"

	"mygo-immigration/backend/internal/repository"
)

func tableVersion(repo *repository.PublicVersionRepo, table, where string, args ...interface{}) (repository.PublicVersion, error) {
	query := "SELECT MAX(updated_at) AS updated_at, COUNT(*) AS count FROM " + table
	if strings.TrimSpace(where) != "" {
		query += " WHERE " + where
	}
	return repo.VersionFromQuery(query, args...)
}

func publicSlug(key, prefix string) string {
	return strings.TrimPrefix(key, prefix)
}
