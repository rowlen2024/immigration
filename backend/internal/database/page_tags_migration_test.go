package database

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPageTagsMigrationFiles(t *testing.T) {
	up, err := os.ReadFile(filepath.Join("..", "..", "database", "migrations", "000029_add_page_tags.up.sql"))
	if err != nil {
		t.Fatalf("read page tags up migration: %v", err)
	}
	down, err := os.ReadFile(filepath.Join("..", "..", "database", "migrations", "000029_add_page_tags.down.sql"))
	if err != nil {
		t.Fatalf("read page tags down migration: %v", err)
	}

	if !strings.Contains(strings.ToLower(string(up)), "add column `tags` json") {
		t.Fatalf("up migration must add pages.tags JSON: %s", up)
	}
	if !strings.Contains(strings.ToLower(string(down)), "drop column `tags`") {
		t.Fatalf("down migration must drop pages.tags: %s", down)
	}
}
