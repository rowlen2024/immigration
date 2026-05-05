package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type migration struct {
	Name string
	SQL  string
}

func RunMigrations(db *gorm.DB, migrationsDir string) error {
	if err := db.Exec("CREATE TABLE IF NOT EXISTS `migrations` (`name` VARCHAR(255) PRIMARY KEY, `applied_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP)").Error; err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []migration
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".up.sql") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(migrationsDir, f.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", f.Name(), err)
		}
		migrations = append(migrations, migration{Name: f.Name(), SQL: string(content)})
	}

	sort.Slice(migrations, func(i, j int) bool { return migrations[i].Name < migrations[j].Name })

	for _, m := range migrations {
		var count int64
		if err := db.Raw("SELECT COUNT(*) FROM `migrations` WHERE `name` = ?", m.Name).Scan(&count).Error; err != nil {
			return fmt.Errorf("failed to check migration %s: %w", m.Name, err)
		}
		if count > 0 {
			log.Printf("Migration %s already applied, skipping", m.Name)
			continue
		}

		log.Printf("Applying migration: %s", m.Name)
		if err := db.Transaction(func(tx *gorm.DB) error {
			for _, stmt := range splitStatements(m.SQL) {
				stmt = strings.TrimSpace(stmt)
				if stmt == "" {
					continue
				}
				if err := tx.Exec(stmt).Error; err != nil {
					return fmt.Errorf("statement failed: %w\nSQL: %s", err, stmt)
				}
			}
			return tx.Exec("INSERT INTO `migrations` (`name`) VALUES (?)", m.Name).Error
		}); err != nil {
			return fmt.Errorf("migration %s failed: %w", m.Name, err)
		}
		log.Printf("Migration %s applied successfully", m.Name)
	}

	return nil
}

func splitStatements(sql string) []string {
	var stmts []string
	var current []byte
	inString := false
	var stringChar byte

	for i := 0; i < len(sql); i++ {
		c := sql[i]
		if inString {
			current = append(current, c)
			if c == stringChar && (i == 0 || sql[i-1] != '\\') {
				inString = false
			}
		} else if c == '\'' || c == '"' {
			inString = true
			stringChar = c
			current = append(current, c)
		} else if c == ';' {
			stmts = append(stmts, string(current))
			current = current[:0]
		} else {
			current = append(current, c)
		}
	}
	remainder := strings.TrimSpace(string(current))
	if remainder != "" {
		stmts = append(stmts, remainder)
	}
	return stmts
}
