package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	log.Println(dbURL)
	if dbURL == "" {
		dbURL = "postgres://postgres:573636@localhost:5432/go_postgres"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := runMigrations(ctx, pool, "db"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("All migrations applied successfully.")
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	// 1. Ensure migrations table exists
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// 2. Scan migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	// 3. Get already applied migrations
	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return fmt.Errorf("failed to fetch applied migrations: %w", err)
	}

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return err
		}
		applied[version] = true
	}
	rows.Close()

	// 4. Apply unapplied migrations sequentially
	for _, filename := range sqlFiles {
		if applied[filename] {
			continue // skip already applied
		}

		log.Printf("Applying migration: %s", filename)

		content, err := os.ReadFile(filepath.Join(migrationsDir, filename))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Split statements:
		//   ALTER TYPE ... ADD VALUE cannot run inside a BEGIN/COMMIT block in
		//   PostgreSQL (the new value isn't visible within that transaction, and
		//   some Postgres versions error outright).  Execute those statements
		//   directly on the pool (auto-commit), then wrap everything else in a
		//   transaction for atomicity.
		autoCommitStmts, txStmts := splitStatements(string(content))

		// 4a. Auto-commit statements (ALTER TYPE ADD VALUE)
		for _, stmt := range autoCommitStmts {
			if _, err := pool.Exec(ctx, stmt); err != nil {
				return fmt.Errorf("error executing auto-commit stmt in %s: %w\nStatement: %s", filename, err, stmt)
			}
		}

		// 4b. Transactional statements
		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", filename, err)
		}

		if len(txStmts) > 0 {
			if _, err := tx.Exec(ctx, strings.Join(txStmts, ";\n")); err != nil {
				_ = tx.Rollback(ctx)
				return fmt.Errorf("error executing %s: %w", filename, err)
			}
		}

		// Record completion inside the same transaction so it's atomic
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, filename); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("error recording migration %s: %w", filename, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", filename, err)
		}

		log.Printf("Successfully applied: %s", filename)
	}

	return nil
}

// splitStatements separates ALTER TYPE ... ADD VALUE statements (which must run
// outside a transaction in Postgres) from all other SQL statements.
//
// It strips single-line SQL comments (--) before splitting to avoid false
// positives from semicolons that appear inside comment text.
func splitStatements(sql string) (autoCommit []string, transactional []string) {
	// Remove single-line comments first so semicolons inside comments don't
	// confuse the splitter.
	var cleaned strings.Builder
	for _, line := range strings.Split(sql, "\n") {
		// Strip any single-line comment (from -- to end of line)
		idx := strings.Index(line, "--")
		if idx >= 0 {
			line = line[:idx]
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		cleaned.WriteString(line)
		cleaned.WriteString("\n")
	}

	for _, raw := range strings.Split(cleaned.String(), ";") {
		stmt := strings.TrimSpace(raw)
		if stmt == "" {
			continue
		}
		upper := strings.ToUpper(stmt)
		if strings.Contains(upper, "ALTER TYPE") && strings.Contains(upper, "ADD VALUE") {
			autoCommit = append(autoCommit, stmt)
		} else {
			transactional = append(transactional, stmt)
		}
	}
	return
}
