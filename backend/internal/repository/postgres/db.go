// Package postgres implements the domain repositories with pgx. It is the only
// package that knows SQL; nothing above it imports pgx.
package postgres

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/starman011/urbanism/backend/internal/domain"
)

//go:embed migrations/*.sql migrations/seed.json
var assets embed.FS

type DB struct{ pool *pgxpool.Pool }

func Open(ctx context.Context, url string) (*DB, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("open pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return &DB{pool: pool}, nil
}

func (d *DB) Close() { d.pool.Close() }

// Migrate applies every embedded migrations/*.sql in name order, once each.
func (d *DB) Migrate(ctx context.Context) error {
	if _, err := d.pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil {
		return err
	}
	entries, err := assets.ReadDir("migrations")
	if err != nil {
		return err
	}
	names := []string{}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for _, n := range names {
		var exists bool
		if err := d.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`, n).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}
		sqlBytes, _ := assets.ReadFile("migrations/" + n)
		tx, err := d.pool.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("migration %s: %w", n, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES($1)`, n); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Seed inserts the prototype's reference content if the tables are empty.
func (d *DB) Seed(ctx context.Context) error {
	var n int
	if err := d.pool.QueryRow(ctx, `SELECT count(*) FROM sessions`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	raw, _ := assets.ReadFile("migrations/seed.json")
	var seed struct {
		Sessions []domain.Session     `json:"sessions"`
		Entries  []domain.Entry       `json:"entries"`
		Journal  []domain.JournalPost `json:"journal"`
	}
	if err := json.Unmarshal(raw, &seed); err != nil {
		return fmt.Errorf("seed.json: %w", err)
	}
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	for _, s := range seed.Sessions {
		if _, err := tx.Exec(ctx, insertSession, s.ID, s.Code, s.Title, s.Date, s.EndDate, s.Place, s.Neighborhood, s.Desc, s.Art, js(s.Tags), js(s.Notes), s.Poll, js(s.Options), s.Sample); err != nil {
			return err
		}
	}
	for _, e := range seed.Entries {
		if _, err := tx.Exec(ctx, insertEntry, e.ID, e.SessionID, string(e.Type), e.Name, e.Hood, e.Text, js(e.Tags), e.Location, e.Art, e.Created); err != nil {
			return err
		}
	}
	for _, p := range seed.Journal {
		if _, err := tx.Exec(ctx, `INSERT INTO journal_posts(id,number,category,title,excerpt,date,read_time,art,author,session_code,body) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT (id) DO NOTHING`,
			p.ID, p.Number, p.Category, p.Title, p.Excerpt, p.Date, p.ReadTime, p.Art, p.Author, p.SessionCode, js(p.Body)); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (d *DB) Repositories() domain.Repositories {
	return domain.Repositories{
		Sessions:  &sessionRepo{d.pool},
		Entries:   &entryRepo{d.pool},
		Replies:   &replyRepo{d.pool},
		Materials: &materialRepo{d.pool},
		Votes:     &voteRepo{d.pool},
		Reactions: &reactionRepo{d.pool},
		Journal:   &journalRepo{d.pool},
	}
}

func js(v []string) []byte {
	if v == nil {
		v = []string{}
	}
	b, _ := json.Marshal(v)
	return b
}
