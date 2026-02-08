package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strings"
)

//go:embed *.sql
var migrationsFS embed.FS

func Apply(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `create table if not exists schema_migrations (version text primary key)`); err != nil {
		return err
	}
	entries, err := migrationsFS.ReadDir(".")
	if err != nil {
		return err
	}
	var files []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	for _, file := range files {
		var exists bool
		if err := db.QueryRowContext(ctx, `select exists(select 1 from schema_migrations where version = $1)`, file).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}
		b, err := migrationsFS.ReadFile(file)
		if err != nil {
			return err
		}
		if _, err := db.ExecContext(ctx, string(b)); err != nil {
			return fmt.Errorf("migration %s failed: %w", file, err)
		}
		if _, err := db.ExecContext(ctx, `insert into schema_migrations (version) values ($1)`, file); err != nil {
			return err
		}
	}
	return nil
}
