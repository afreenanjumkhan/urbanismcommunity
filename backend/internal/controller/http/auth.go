package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/starman011/urbanism/backend/internal/domain"
	"github.com/starman011/urbanism/backend/internal/service"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

func (c *Controller) identity(r *http.Request) (domain.Account, error) {
	if c.auth == nil {
		return domain.Account{}, service.ErrUnauthorized
	}
	return c.auth.Verify(r.Context(), strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
}
func (c *Controller) authRoutes(r chi.Router) {
	r.Get("/auth/status", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]bool{"enabled": c.auth.Enabled()})
	})
	r.Get("/auth/me", func(w http.ResponseWriter, r *http.Request) {
		u, e := c.identity(r)
		if e != nil {
			c.fail(w, r, e)
			return
		}
		writeJSON(w, 200, u)
	})
	r.With(authRateLimit).Post("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if e := decode(r, &in); e != nil {
			c.fail(w, r, e)
			return
		}
		if c.auth == nil {
			c.fail(w, r, service.ErrUnauthorized)
			return
		}
		token, e := c.auth.Login(r.Context(), in.Email, in.Password)
		if e != nil {
			c.fail(w, r, e)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, 200, map[string]string{"token": token})
	})
	r.With(authRateLimit).Post("/auth/activate", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Token    string `json:"token"`
		}
		if e := decode(r, &in); e != nil {
			c.fail(w, r, e)
			return
		}
		if c.auth == nil {
			c.fail(w, r, service.ErrUnauthorized)
			return
		}
		if e := c.auth.Activate(r.Context(), in.Email, in.Token, in.Password); e != nil {
			c.fail(w, r, e)
			return
		}
		w.WriteHeader(204)
	})
	r.Get("/admin/members", func(w http.ResponseWriter, r *http.Request) {
		u, e := c.identity(r)
		if e != nil {
			c.fail(w, r, e)
			return
		}
		if u.Role != "super_admin" {
			c.fail(w, r, service.ErrForbidden)
			return
		}
		users, e := c.auth.Members(r.Context(), u)
		if e != nil {
			c.fail(w, r, e)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, 200, users)
	})
	r.Post("/admin/members", func(w http.ResponseWriter, r *http.Request) {
		u, e := c.identity(r)
		if e != nil {
			c.fail(w, r, e)
			return
		}
		var in struct {
			Email string `json:"email"`
		}
		if e = decode(r, &in); e != nil {
			c.fail(w, r, e)
			return
		}
		token, e := c.auth.Invite(r.Context(), u, in.Email)
		if e != nil {
			c.fail(w, r, e)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, 201, map[string]string{"invitation": token})
	})
	r.Delete("/admin/members/{email}", func(w http.ResponseWriter, r *http.Request) {
		u, e := c.identity(r)
		if e != nil {
			c.fail(w, r, e)
			return
		}
		if e = c.auth.Revoke(r.Context(), u, chi.URLParam(r, "email")); e != nil {
			c.fail(w, r, e)
			return
		}
		w.WriteHeader(204)
	})
}

// Per-process limiter deliberately does not trust client-supplied forwarded headers.
// Deploy an edge limiter too when scaling across instances.
var attempts = struct {
	sync.Mutex
	entries map[string]struct {
		count int
		until time.Time
	}
}{entries: make(map[string]struct {
	count int
	until time.Time
})}

func authRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		now := time.Now()
		attempts.Lock()
		for k, v := range attempts.entries {
			if now.After(v.until) {
				delete(attempts.entries, k)
			}
		}
		v := attempts.entries[host]
		if v.until.IsZero() {
			v.until = now.Add(time.Minute)
		}
		v.count++
		allowed := v.count <= 10 && len(attempts.entries) < 10000
		if allowed {
			attempts.entries[host] = v
		}
		attempts.Unlock()
		if !allowed {
			w.Header().Set("Retry-After", "60")
			writeJSON(w, 429, map[string]string{"error": "Too many attempts. Try again in a minute."})
			return
		}
		next.ServeHTTP(w, r)
	})
}
