package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/starman011/urbanism/backend/internal/domain"
)

const (
	insertSession = `INSERT INTO sessions(id,code,title,date,end_date,place,neighborhood,description,art,tags,notes,poll,options,sample)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`
	selectSession = `SELECT id,code,title,to_char(date,'YYYY-MM-DD'),to_char(end_date,'YYYY-MM-DD'),place,neighborhood,description,art,tags,notes,poll,options,sample,created_at FROM sessions`
	insertEntry   = `INSERT INTO entries(id,session_id,type,name,hood,text,tags,location,art,created) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	selectEntry   = `SELECT id,session_id,type,name,hood,text,tags,location,art,created FROM entries`
)

func mapErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	var pg *pgconn.PgError
	if errors.As(err, &pg) {
		switch pg.Code {
		case "23505":
			return domain.ErrConflict
		case "23503":
			return domain.ErrNotFound
		}
	}
	return err
}

// ---- sessions ----

type sessionRepo struct{ pool *pgxpool.Pool }

func scanSession(row pgx.Row) (*domain.Session, error) {
	var s domain.Session
	if err := row.Scan(&s.ID, &s.Code, &s.Title, &s.Date, &s.EndDate, &s.Place, &s.Neighborhood, &s.Desc, &s.Art, &s.Tags, &s.Notes, &s.Poll, &s.Options, &s.Sample, &s.CreatedAt); err != nil {
		return nil, mapErr(err)
	}
	return &s, nil
}

func (r *sessionRepo) List(ctx context.Context) ([]domain.Session, error) {
	rows, err := r.pool.Query(ctx, selectSession+` ORDER BY date DESC, created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Session{}
	for rows.Next() {
		s, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *s)
	}
	return out, rows.Err()
}

func (r *sessionRepo) GetByCode(ctx context.Context, code string) (*domain.Session, error) {
	return scanSession(r.pool.QueryRow(ctx, selectSession+` WHERE code=$1`, code))
}

func (r *sessionRepo) Create(ctx context.Context, s *domain.Session) error {
	_, err := r.pool.Exec(ctx, insertSession, s.ID, s.Code, s.Title, s.Date, s.EndDate, s.Place, s.Neighborhood, s.Desc, s.Art, js(s.Tags), js(s.Notes), s.Poll, js(s.Options), s.Sample)
	return mapErr(err)
}

// ---- entries ----

type entryRepo struct{ pool *pgxpool.Pool }

func (r *entryRepo) List(ctx context.Context, sessionID string) ([]domain.Entry, error) {
	q := selectEntry + ` ORDER BY created ASC`
	args := []any{}
	if sessionID != "" {
		q = selectEntry + ` WHERE session_id=$1 ORDER BY created ASC`
		args = append(args, sessionID)
	}
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Entry{}
	for rows.Next() {
		var e domain.Entry
		var created time.Time
		if err := rows.Scan(&e.ID, &e.SessionID, &e.Type, &e.Name, &e.Hood, &e.Text, &e.Tags, &e.Location, &e.Art, &created); err != nil {
			return nil, err
		}
		e.Created = created.UTC().Format(time.RFC3339)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *entryRepo) Create(ctx context.Context, e *domain.Entry) error {
	_, err := r.pool.Exec(ctx, insertEntry, e.ID, e.SessionID, string(e.Type), e.Name, e.Hood, e.Text, js(e.Tags), e.Location, e.Art, e.Created)
	return mapErr(err)
}

func (r *entryRepo) Exists(ctx context.Context, id string) (bool, error) {
	var ok bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM entries WHERE id=$1)`, id).Scan(&ok)
	return ok, err
}

// ---- replies ----

type replyRepo struct{ pool *pgxpool.Pool }

func (r *replyRepo) List(ctx context.Context) ([]domain.Reply, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,entry_id,name,text,created FROM replies ORDER BY created ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Reply{}
	for rows.Next() {
		var x domain.Reply
		var created time.Time
		if err := rows.Scan(&x.ID, &x.EntryID, &x.Name, &x.Text, &created); err != nil {
			return nil, err
		}
		x.Created = created.UTC().Format(time.RFC3339)
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *replyRepo) Create(ctx context.Context, x *domain.Reply) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO replies(id,entry_id,name,text,created) VALUES($1,$2,$3,$4,$5)`, x.ID, x.EntryID, x.Name, x.Text, x.Created)
	return mapErr(err)
}

// ---- materials ----

type materialRepo struct{ pool *pgxpool.Pool }

func (r *materialRepo) List(ctx context.Context) ([]domain.Material, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,session_id,title,category,description,credited_to,file_name,mime_type,media_id,size,created FROM materials ORDER BY created ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Material{}
	for rows.Next() {
		var m domain.Material
		var created time.Time
		if err := rows.Scan(&m.ID, &m.SessionID, &m.Title, &m.Category, &m.Description, &m.CreditedTo, &m.FileName, &m.MimeType, &m.MediaID, &m.Size, &created); err != nil {
			return nil, err
		}
		m.Created = created.UTC().Format(time.RFC3339)
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *materialRepo) Create(ctx context.Context, m *domain.Material) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO materials(id,session_id,title,category,description,credited_to,file_name,mime_type,media_id,size,created) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		m.ID, m.SessionID, m.Title, string(m.Category), m.Description, m.CreditedTo, m.FileName, m.MimeType, m.MediaID, m.Size, m.Created)
	return mapErr(err)
}

// ---- votes ----

type voteRepo struct{ pool *pgxpool.Pool }

func (r *voteRepo) ForClient(ctx context.Context, clientID string) (map[string]int, error) {
	rows, err := r.pool.Query(ctx, `SELECT session_id,option_idx FROM votes WHERE client_id=$1`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var sid string
		var opt int
		if err := rows.Scan(&sid, &opt); err != nil {
			return nil, err
		}
		out[sid] = opt
	}
	return out, rows.Err()
}

func (r *voteRepo) Upsert(ctx context.Context, v domain.Vote) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO votes(client_id,session_id,option_idx) VALUES($1,$2,$3)
		ON CONFLICT (client_id,session_id) DO UPDATE SET option_idx=EXCLUDED.option_idx, updated_at=now()`, v.ClientID, v.SessionID, v.Option)
	return mapErr(err)
}

func (r *voteRepo) Tally(ctx context.Context, sessionID string) (map[int]int, error) {
	rows, err := r.pool.Query(ctx, `SELECT option_idx,count(*) FROM votes WHERE session_id=$1 GROUP BY option_idx`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int]int{}
	for rows.Next() {
		var opt, n int
		if err := rows.Scan(&opt, &n); err != nil {
			return nil, err
		}
		out[opt] = n
	}
	return out, rows.Err()
}

// ---- reactions ----

type reactionRepo struct{ pool *pgxpool.Pool }

func (r *reactionRepo) ForClient(ctx context.Context, clientID string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT entry_id FROM reactions WHERE client_id=$1`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (r *reactionRepo) Add(ctx context.Context, x domain.Reaction) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO reactions(client_id,entry_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, x.ClientID, x.EntryID)
	return mapErr(err)
}

func (r *reactionRepo) Remove(ctx context.Context, x domain.Reaction) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM reactions WHERE client_id=$1 AND entry_id=$2`, x.ClientID, x.EntryID)
	return mapErr(err)
}

// ---- journal ----

type journalRepo struct{ pool *pgxpool.Pool }

func (r *journalRepo) List(ctx context.Context) ([]domain.JournalPost, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,number,category,title,excerpt,to_char(date,'YYYY-MM-DD'),read_time,art,author,session_code,body FROM journal_posts ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.JournalPost{}
	for rows.Next() {
		var p domain.JournalPost
		if err := rows.Scan(&p.ID, &p.Number, &p.Category, &p.Title, &p.Excerpt, &p.Date, &p.ReadTime, &p.Art, &p.Author, &p.SessionCode, &p.Body); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
