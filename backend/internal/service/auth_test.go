package service

import (
	"context"
	"github.com/starman011/urbanism/backend/internal/domain"
	"strings"
	"testing"
	"time"
)

type authRepo struct {
	users     map[string]domain.Account
	activated bool
}

func (r *authRepo) Get(_ context.Context, e string) (domain.Account, error) {
	u, ok := r.users[e]
	if !ok {
		return u, domain.ErrNotFound
	}
	return u, nil
}
func (r *authRepo) List(context.Context) ([]domain.Account, error)       { return []domain.Account{}, nil }
func (r *authRepo) Invite(context.Context, string, string, string) error { return nil }
func (r *authRepo) Activate(_ context.Context, e, h, p string, b bool) error {
	if r.activated {
		return domain.ErrConflict
	}
	r.activated = true
	return nil
}
func (r *authRepo) Revoke(_ context.Context, e, actor string) error {
	u := r.users[e]
	u.Active = false
	r.users[e] = u
	return nil
}
func TestJWTBoundaries(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	hash, err := passwordHash("test-password-long-enough")
	if err != nil {
		t.Fatal(err)
	}
	repo := &authRepo{users: map[string]domain.Account{"board@example.org": {Email: "board@example.org", Role: "board", Active: true, PasswordHash: hash}}}
	a := &Auth{Repo: repo, Secret: strings.Repeat("s", 32), Now: func() time.Time { return now }}
	token, err := a.Login(ctx, "BOARD@example.org", "test-password-long-enough")
	if err != nil {
		t.Fatal(err)
	}
	if u, e := a.Verify(ctx, token); e != nil || !CanCreateSession(u) {
		t.Fatalf("valid token rejected: %v", e)
	}
	p := strings.Split(token, ".")
	p[1] = "e30"
	if _, e := a.Verify(ctx, strings.Join(p, ".")); e == nil {
		t.Fatal("tampered token accepted")
	}
	now = now.Add(2 * time.Hour)
	if _, e := a.Verify(ctx, token); e == nil {
		t.Fatal("expired token accepted")
	}
	now = now.Add(-2 * time.Hour)
	u := repo.users["board@example.org"]
	u.Active = false
	repo.users[u.Email] = u
	if _, e := a.Verify(ctx, token); e == nil {
		t.Fatal("revoked token accepted")
	}
	if _, e := a.Login(ctx, u.Email, "wrong-password"); e == nil {
		t.Fatal("wrong password accepted")
	}
}
func TestRoleAndBootstrapBoundaries(t *testing.T) {
	ctx := context.Background()
	r := &authRepo{users: map[string]domain.Account{}}
	a := &Auth{Repo: r, Secret: strings.Repeat("s", 32), Bootstrap: strings.Repeat("b", 32), Now: time.Now}
	board := domain.Account{Email: "board@example.org", Role: "board", Active: true}
	if _, e := a.Invite(ctx, board, "other@example.org"); e != ErrForbidden {
		t.Fatal("board can invite")
	}
	if e := a.Revoke(ctx, board, SuperAdminEmail); e != ErrForbidden {
		t.Fatal("board can revoke admin")
	}
	if e := a.Activate(ctx, SuperAdminEmail, strings.Repeat("x", 32), "long-test-password"); e == nil {
		t.Fatal("email alone bootstraps admin")
	}
	if e := a.Activate(ctx, SuperAdminEmail, a.Bootstrap, "long-test-password"); e != nil {
		t.Fatal(e)
	}
	if e := a.Activate(ctx, SuperAdminEmail, a.Bootstrap, "long-test-password"); e == nil {
		t.Fatal("bootstrap reused")
	}
	if CanCreateSession(domain.Account{Role: "super_admin"}) {
		t.Fatal("inactive admin allowed")
	}
	if _, e := passwordHash("short"); e == nil {
		t.Fatal("weak password accepted")
	}
}
