package service

import (
	"context"
	"crypto/hmac"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/starman011/urbanism/backend/internal/domain"
	"net/mail"
	"strings"
	"time"
)

var ErrUnauthorized = errors.New("sign in required or credentials invalid")
var ErrForbidden = errors.New("permission denied")

const SuperAdminEmail = "akhannyc72@gmail.com"

type Auth struct {
	Repo      domain.Accounts
	Secret    string
	Bootstrap string
	Now       func() time.Time
}
type claims struct {
	Sub             string `json:"sub"`
	Iss             string `json:"iss"`
	Exp             int64  `json:"exp"`
	PasswordVersion string `json:"pv"`
}

func (a *Auth) Enabled() bool       { return a != nil && len(a.Secret) >= 32 }
func normalize(email string) string { return strings.ToLower(strings.TrimSpace(email)) }
func hashToken(s string) string     { v := sha256.Sum256([]byte(s)); return hex.EncodeToString(v[:]) }
func passwordHash(password string) (string, error) {
	if len(password) < 12 || len(password) > 128 {
		return "", domain.Invalid("Password must be 12–128 bytes.")
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key, err := pbkdf2.Key(sha256.New, password, salt, 600000, 32)
	return hex.EncodeToString(salt) + ":" + hex.EncodeToString(key), err
}
func passwordOK(password, encoded string) bool {
	p := strings.Split(encoded, ":")
	if len(p) != 2 {
		return false
	}
	salt, e := hex.DecodeString(p[0])
	if e != nil {
		return false
	}
	want, e := hex.DecodeString(p[1])
	if e != nil {
		return false
	}
	got, e := pbkdf2.Key(sha256.New, password, salt, 600000, 32)
	return e == nil && subtle.ConstantTimeCompare(got, want) == 1
}
func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	if !a.Enabled() || len(password) > 128 {
		return "", ErrUnauthorized
	}
	u, err := a.Repo.Get(ctx, normalize(email))
	if err != nil || !u.Active {
		// Comparable work for absent and existing users.
		_, _ = pbkdf2.Key(sha256.New, password, make([]byte, 16), 600000, 32)
		return "", ErrUnauthorized
	}
	if !passwordOK(password, u.PasswordHash) {
		return "", ErrUnauthorized
	}
	body, _ := json.Marshal(claims{u.Email, "field-notes", a.Now().Add(time.Hour).Unix(), hashToken(u.PasswordHash)})
	head := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	unsigned := head + "." + base64.RawURLEncoding.EncodeToString(body)
	h := hmac.New(sha256.New, []byte(a.Secret))
	h.Write([]byte(unsigned))
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(h.Sum(nil)), nil
}
func (a *Auth) Verify(ctx context.Context, token string) (domain.Account, error) {
	var empty domain.Account
	if !a.Enabled() || len(token) > 4096 {
		return empty, ErrUnauthorized
	}
	p := strings.Split(token, ".")
	if len(p) != 3 {
		return empty, ErrUnauthorized
	}
	header, e := base64.RawURLEncoding.DecodeString(p[0])
	if e != nil || string(header) != `{"alg":"HS256","typ":"JWT"}` {
		return empty, ErrUnauthorized
	}
	sig, e := base64.RawURLEncoding.DecodeString(p[2])
	if e != nil {
		return empty, ErrUnauthorized
	}
	h := hmac.New(sha256.New, []byte(a.Secret))
	h.Write([]byte(p[0] + "." + p[1]))
	if !hmac.Equal(sig, h.Sum(nil)) {
		return empty, ErrUnauthorized
	}
	body, e := base64.RawURLEncoding.DecodeString(p[1])
	var c claims
	if e != nil || json.Unmarshal(body, &c) != nil || c.Iss != "field-notes" || c.Exp <= a.Now().Unix() {
		return empty, ErrUnauthorized
	}
	u, e := a.Repo.Get(ctx, c.Sub)
	if e != nil || !u.Active || c.PasswordVersion != hashToken(u.PasswordHash) {
		return empty, ErrUnauthorized
	}
	return u, nil
}
func (a *Auth) Activate(ctx context.Context, email, token, password string) error {
	if !a.Enabled() || len(token) < 32 || len(token) > 256 {
		return ErrUnauthorized
	}
	email = normalize(email)
	bootstrap := email == SuperAdminEmail
	if bootstrap && (len(a.Bootstrap) < 32 || subtle.ConstantTimeCompare([]byte(token), []byte(a.Bootstrap)) != 1) {
		return ErrUnauthorized
	}
	hash, err := passwordHash(password)
	if err != nil {
		return err
	}
	if err = a.Repo.Activate(ctx, email, hashToken(token), hash, bootstrap); err != nil {
		return ErrUnauthorized
	}
	return nil
}
func (a *Auth) Invite(ctx context.Context, actor domain.Account, email string) (string, error) {
	if !actor.Active || actor.Role != "super_admin" {
		return "", ErrForbidden
	}
	email = normalize(email)
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email || email == SuperAdminEmail || len(email) > 254 {
		return "", domain.Invalid("Enter a valid board member email.")
	}
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(b)
	err = a.Repo.Invite(ctx, email, hashToken(token), actor.Email)
	return token, err
}
func (a *Auth) Revoke(ctx context.Context, actor domain.Account, email string) error {
	if !actor.Active || actor.Role != "super_admin" || normalize(email) == SuperAdminEmail {
		return ErrForbidden
	}
	return a.Repo.Revoke(ctx, normalize(email), actor.Email)
}
func CanCreateSession(u domain.Account) bool {
	return u.Active && (u.Role == "board" || u.Role == "super_admin")
}

func (a *Auth) Members(ctx context.Context, actor domain.Account) ([]domain.Account, error) {
	if !actor.Active || actor.Role != "super_admin" {
		return nil, ErrForbidden
	}
	return a.Repo.List(ctx)
}
