package db

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Migrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TEXT NOT NULL)`).Error; err != nil {
		return err
	}

	paths, err := fs.Glob(migrationsFS, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(paths)

	for _, path := range paths {
		version := strings.TrimPrefix(path, "migrations/")
		version = strings.TrimSuffix(version, ".sql")

		var count int
		if err := db.Raw(`SELECT COUNT(*) FROM schema_migrations WHERE version=?`, version).Scan(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		content, err := migrationsFS.ReadFile(path)
		if err != nil {
			return err
		}
		statements := splitSQL(string(content))
		appliedAt := time.Now().UTC().Format(time.RFC3339)

		if err := db.Transaction(func(tx *gorm.DB) error {
			for _, stmt := range statements {
				stmt = strings.TrimSpace(stmt)
				if stmt == "" {
					continue
				}
				if err := tx.Exec(stmt).Error; err != nil {
					return fmt.Errorf("migration %s failed: %w", version, err)
				}
			}

			if err := tx.Exec(`INSERT INTO schema_migrations (version, applied_at) VALUES (?,?)`, version, appliedAt).Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func splitSQL(sql string) []string {
	parts := strings.Split(sql, ";")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		stmt := strings.TrimSpace(part)
		if stmt == "" {
			continue
		}
		out = append(out, stmt)
	}
	return out
}
