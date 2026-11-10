// Package http is the controller layer: it decodes requests, calls services,
// and encodes responses. No business rules live here.
package http

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/starman011/urbanism/backend/internal/domain"
	"github.com/starman011/urbanism/backend/internal/middleware"
	"github.com/starman011/urbanism/backend/internal/service"
)

type Controller struct {
	svc *service.Services
	log *slog.Logger
}

func NewRouter(svc *service.Services, log *slog.Logger, origins []string) http.Handler {
	c := &Controller{svc: svc, log: log}
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Recover(log), middleware.Logger(log), middleware.CORS(origins), middleware.ClientID, middleware.BodyLimit(1<<20))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, map[string]string{"status": "ok"}) })

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/state", c.getState)
		r.Get("/journal", c.listJournal)
		r.Get("/export", c.exportJSON)
		r.Get("/export.csv", c.exportCSV)
		r.Get("/sessions", c.listSessions)
		r.Post("/sessions", c.createSession)
		r.Get("/sessions/{code}", c.getSession)
		r.Get("/sessions/{id}/tally", c.tally)
		r.Put("/sessions/{id}/vote", c.vote)
		r.Get("/entries", c.listEntries)
		r.Post("/entries", c.createEntry)
		r.Put("/entries/{id}/like", c.like(true))
		r.Delete("/entries/{id}/like", c.like(false))
		r.Get("/replies", c.listReplies)
		r.Post("/replies", c.createReply)
		r.Get("/materials", c.listMaterials)
		r.Post("/materials", c.createMaterial)
	})
	return r
}

// ---- helpers ----

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (c *Controller) fail(w http.ResponseWriter, r *http.Request, err error) {
	var ve *domain.ValidationError
	switch {
	case errors.As(err, &ve):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": ve.Msg})
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	case errors.Is(err, domain.ErrConflict):
		writeJSON(w, http.StatusConflict, map[string]string{"error": "already exists"})
	default:
		c.log.Error("handler", "id", middleware.RequestIDFrom(r.Context()), "err", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
}

func decode(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return domain.Invalid("invalid JSON body: " + err.Error())
	}
	return nil
}

// ---- handlers ----

func (c *Controller) getState(w http.ResponseWriter, r *http.Request) {
	st, err := c.svc.State.Get(r.Context(), middleware.ClientIDFrom(r.Context()))
	if err != nil {
		c.fail(w, r, err)
		return
	}
	writeJSON(w, 200, st)
}

func (c *Controller) listJournal(w http.ResponseWriter, r *http.Request) {
	v, err := c.svc.Journal.List(r.Context())
	if err != nil {
		c.fail(w, r, err)
		return
	}
	writeJSON(w, 200, v)
}

func (c *Controller) listSessions(w http.ResponseWriter, r *http.Request) {
	v, err := c.svc.Sessions.List(r.Context())
	if err != nil {
		c.fail(w, r, err)
		return
	}
	writeJSON(w, 200, v)
}

func (c *Controller) getSession(w http.ResponseWriter, r *http.Request) {
	v, err := c.svc.Sessions.GetByCode(r.Context(), chi.URLParam(r, "code"))
	if err != nil {
		c.fail(w, r, err)
		return
	}
	writeJSON(w, 200, v)
}

func (c *Controller) createSession(w http.ResponseWriter, r *http.Request) {
	var in domain.Session
	if err := decode(r, &in); err != nil {
		c.fail(w, r, err)
		return
	}
	v, err := c.svc.Sessions.Create(r.Context(), in)
	if err != nil {
		c.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, v)
}

func (c *Controller) tally(w http.ResponseWriter, r *http.Request) {
	v, err := c.svc.Votes.Tally(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		c.fail(w, r, err)
		return
	}
	out := map[string]int{}
	for k, n := range v {
		out[strconv.Itoa(k)] = n
	}
	writeJSON(w, 200, out)
}

func (c *Controller) vote(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Option int `json:"option"`
	}
	if err := decode(r, &in); err != nil {
		c.fail(w, r, err)
		return
	}
	if err := c.svc.Votes.Cast(r.Context(), middleware.ClientIDFrom(r.Context()), chi.URLParam(r, "id"), in.Option); err != nil {
		c.fail(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) listEntries(w http.ResponseWriter, r *http.Request) {
	v, err := c.svc.Entries.List(r.Context(), r.URL.Query().Get("session"))
	if err != nil {
		c.fail(w, r, err)
		return
	}
	writeJSON(w, 200, v)
}

func (c *Controller) createEntry(w http.ResponseWriter, r *http.Request) {
	var in domain.Entry
	if err := decode(r, &in); err != nil {
		c.fail(w, r, err)
		return
	}
	v, err := c.svc.Entries.Create(r.Context(), in)
	if err != nil {
		c.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, v)
}

func (c *Controller) like(on bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := c.svc.Reactions.Set(r.Context(), middleware.ClientIDFrom(r.Context()), chi.URLParam(r, "id"), on); err != nil {
			c.fail(w, r, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (c *Controller) listReplies(w http.ResponseWriter, r *http.Request) {
	v, err := c.svc.Replies.List(r.Context())
	if err != nil {
		c.fail(w, r, err)
		return
	}
	writeJSON(w, 200, v)
}

func (c *Controller) createReply(w http.ResponseWriter, r *http.Request) {
	var in domain.Reply
	if err := decode(r, &in); err != nil {
		c.fail(w, r, err)
		return
	}
	v, err := c.svc.Replies.Create(r.Context(), in)
	if err != nil {
		c.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, v)
}

func (c *Controller) listMaterials(w http.ResponseWriter, r *http.Request) {
	v, err := c.svc.Materials.List(r.Context())
	if err != nil {
		c.fail(w, r, err)
		return
	}
	writeJSON(w, 200, v)
}

func (c *Controller) createMaterial(w http.ResponseWriter, r *http.Request) {
	var in domain.Material
	if err := decode(r, &in); err != nil {
		c.fail(w, r, err)
		return
	}
	v, err := c.svc.Materials.Create(r.Context(), in)
	if err != nil {
		c.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, v)
}

// ---- research export ----

func (c *Controller) exportJSON(w http.ResponseWriter, r *http.Request) {
	v, err := c.svc.Export.Snapshot(r.Context())
	if err != nil {
		c.fail(w, r, err)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename=\"field-notes-export.json\"")
	writeJSON(w, 200, v)
}

// exportCSV flattens entries with their session for spreadsheet work.
func (c *Controller) exportCSV(w http.ResponseWriter, r *http.Request) {
	v, err := c.svc.Export.Snapshot(r.Context())
	if err != nil {
		c.fail(w, r, err)
		return
	}
	sessions := map[string]service.SessionExport{}
	for _, s := range v.Sessions {
		sessions[s.ID] = s
	}
	replies := map[string]int{}
	for _, rp := range v.Replies {
		replies[rp.EntryID]++
	}
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"field-notes-entries.csv\"")
	cw := csv.NewWriter(w)
	_ = cw.Write([]string{"entry_id", "created", "session_code", "session_title", "session_date", "neighborhood", "type", "name", "hood", "location", "tags", "text", "reply_count", "sample_session"})
	for _, e := range v.Entries {
		s := sessions[e.SessionID]
		_ = cw.Write([]string{e.ID, e.Created, s.Code, s.Title, s.Date, s.Neighborhood, string(e.Type), e.Name, e.Hood, e.Location, strings.Join(e.Tags, "|"), e.Text, fmt.Sprint(replies[e.ID]), fmt.Sprint(s.Sample)})
	}
	cw.Flush()
}
