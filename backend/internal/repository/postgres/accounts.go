package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/starman011/urbanism/backend/internal/domain"
)

func (d *DB) Accounts() domain.Accounts { return d }
func (d *DB) Get(ctx context.Context, email string) (domain.Account, error) {
	var a domain.Account
	err := d.pool.QueryRow(ctx, `SELECT email,role,active,password_hash FROM accounts WHERE email=$1`, email).Scan(&a.Email, &a.Role, &a.Active, &a.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		err = domain.ErrNotFound
	}
	return a, err
}
func (d *DB) List(ctx context.Context) ([]domain.Account, error) {
	rows, err := d.pool.Query(ctx, `SELECT email,role,active FROM accounts ORDER BY email`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Account{}
	for rows.Next() {
		var a domain.Account
		if err := rows.Scan(&a.Email, &a.Role, &a.Active); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
func (d *DB) Invite(ctx context.Context, email, hash, actor string) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `INSERT INTO accounts(email,role,invite_hash,invite_expires) VALUES($1,'board',$2,now()+interval '24 hours') ON CONFLICT(email) DO UPDATE SET invite_hash=$2,invite_expires=now()+interval '24 hours' WHERE accounts.role='board' AND NOT accounts.active`, email, hash)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return domain.ErrConflict
	}
	if _, err = tx.Exec(ctx, `INSERT INTO account_audit(actor,target,action) VALUES($1,$2,'invite')`, actor, email); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (d *DB) Activate(ctx context.Context, email, tokenHash, password string, bootstrap bool) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var count int64
	if bootstrap {
		r, e := tx.Exec(ctx, `INSERT INTO accounts(email,role,active,password_hash) VALUES($1,'super_admin',true,$2) ON CONFLICT DO NOTHING`, email, password)
		err = e
		count = r.RowsAffected()
	} else {
		r, e := tx.Exec(ctx, `UPDATE accounts SET active=true,password_hash=$3,invite_hash='',invite_expires=NULL WHERE email=$1 AND invite_hash=$2 AND invite_expires>now() AND NOT active AND role='board'`, email, tokenHash, password)
		err = e
		count = r.RowsAffected()
	}
	if err != nil {
		return err
	}
	if count != 1 {
		return domain.ErrConflict
	}
	if _, err = tx.Exec(ctx, `INSERT INTO account_audit(actor,target,action) VALUES($1,$1,'activate')`, email); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (d *DB) Revoke(ctx context.Context, email, actor string) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	r, err := tx.Exec(ctx, `UPDATE accounts SET active=false,password_hash='',invite_hash='',invite_expires=NULL WHERE email=$1 AND role='board'`, email)
	if err != nil {
		return err
	}
	if r.RowsAffected() != 1 {
		return domain.ErrNotFound
	}
	if _, err = tx.Exec(ctx, `INSERT INTO account_audit(actor,target,action) VALUES($1,$2,'revoke')`, actor, email); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
