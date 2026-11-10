package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/starman011/urbanism/backend/internal/domain"
)

func seeded() (*memRepos, *Services) {
	m := newMem()
	m.sessions = []domain.Session{{ID: "s1", Code: "ATL24", Title: "Atlantic", Date: "2026-10-18", EndDate: "2026-10-18", Options: []string{"a", "b", "c"}, Tags: []string{}, Notes: []string{}}}
	m.entries = []domain.Entry{{ID: "e1", SessionID: "s1", Type: "text", Text: "hi", Tags: []string{}}}
	return m, newTestServices(m)
}

func wantValidation(t *testing.T, err error, frag string) {
	t.Helper()
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("want validation error, got %v", err)
	}
	if !strings.Contains(err.Error(), frag) {
		t.Fatalf("want message containing %q, got %q", frag, err.Error())
	}
}

func TestSessionCreate(t *testing.T) {
	ctx := context.Background()
	_, svc := seeded()

	t.Run("normalizes code, defaults, trims", func(t *testing.T) {
		s, err := svc.Sessions.Create(ctx, domain.Session{Code: " gow-09 ", Title: "  Gowanus  ", Date: "2026-12-01", Tags: []string{"A", "A", " ", "B"}})
		if err != nil {
			t.Fatal(err)
		}
		if s.Code != "GOW09" || s.Title != "Gowanus" || s.EndDate != "2026-12-01" || s.Art != "blocks" {
			t.Fatalf("bad normalization: %+v", s)
		}
		if len(s.Tags) != 2 || s.Sample || s.ID == "" {
			t.Fatalf("bad defaults: %+v", s)
		}
		if s.Options == nil || s.Notes == nil {
			t.Fatal("slices must be non-nil for the client validator")
		}
	})
	t.Run("rejects bad input", func(t *testing.T) {
		_, err := svc.Sessions.Create(ctx, domain.Session{Code: "AB", Title: "x", Date: "2026-01-01"})
		wantValidation(t, err, "code")
		_, err = svc.Sessions.Create(ctx, domain.Session{Code: "ABC", Title: "", Date: "2026-01-01"})
		wantValidation(t, err, "title")
		_, err = svc.Sessions.Create(ctx, domain.Session{Code: "ABC", Title: "x", Date: "01/01/2026"})
		wantValidation(t, err, "date")
		_, err = svc.Sessions.Create(ctx, domain.Session{Code: "ABC", Title: "x", Date: "2026-01-02", EndDate: "2026-01-01"})
		wantValidation(t, err, "endDate")
		_, err = svc.Sessions.Create(ctx, domain.Session{Code: "ABC", Title: "x", Date: "2026-01-01", Poll: "q?", Options: []string{"only"}})
		wantValidation(t, err, "two options")
	})
	t.Run("duplicate code conflicts", func(t *testing.T) {
		_, err := svc.Sessions.Create(ctx, domain.Session{Code: "atl24", Title: "dup", Date: "2026-01-01"})
		if !errors.Is(err, domain.ErrConflict) {
			t.Fatalf("want conflict, got %v", err)
		}
	})
	t.Run("get by code normalizes", func(t *testing.T) {
		s, err := svc.Sessions.GetByCode(ctx, " atl24 ")
		if err != nil || s.ID != "s1" {
			t.Fatalf("got %v %v", s, err)
		}
		_, err = svc.Sessions.GetByCode(ctx, "   ")
		wantValidation(t, err, "code")
	})
}

func TestEntryCreate(t *testing.T) {
	ctx := context.Background()
	_, svc := seeded()
	e, err := svc.Entries.Create(ctx, domain.Entry{SessionID: "s1", Text: "  note  ", Name: " ", Created: "garbage"})
	if err != nil {
		t.Fatal(err)
	}
	if e.Type != domain.EntryText || e.Name != "Anonymous neighbor" || e.Created != "2026-11-11T12:00:00Z" || e.Tags == nil {
		t.Fatalf("bad defaults: %+v", e)
	}
	_, err = svc.Entries.Create(ctx, domain.Entry{SessionID: "s1", Text: "x", Type: "video"})
	wantValidation(t, err, "type")
	_, err = svc.Entries.Create(ctx, domain.Entry{SessionID: "s1", Text: "   "})
	wantValidation(t, err, "text")
	_, err = svc.Entries.Create(ctx, domain.Entry{Text: "x"})
	wantValidation(t, err, "sessionId")
	long, _ := svc.Entries.Create(ctx, domain.Entry{SessionID: "s1", Text: strings.Repeat("a", maxText+50)})
	if len(long.Text) != maxText {
		t.Fatalf("text not capped: %d", len(long.Text))
	}
}

func TestReplyCreate(t *testing.T) {
	ctx := context.Background()
	_, svc := seeded()
	_, err := svc.Replies.Create(ctx, domain.Reply{EntryID: "missing", Text: "x"})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("want not found, got %v", err)
	}
	r, err := svc.Replies.Create(ctx, domain.Reply{EntryID: "e1", Text: "x"})
	if err != nil || r.Name != "Anonymous neighbor" || r.ID == "" {
		t.Fatalf("got %+v %v", r, err)
	}
}

func TestMaterialCreate(t *testing.T) {
	ctx := context.Background()
	_, svc := seeded()
	base := domain.Material{SessionID: "s1", Title: "Panel", Category: "panel", FileName: "p.pdf", MimeType: "application/pdf", MediaID: "m1", Size: 10}
	if _, err := svc.Materials.Create(ctx, base); err != nil {
		t.Fatal(err)
	}
	bad := base
	bad.Category = "video"
	_, err := svc.Materials.Create(ctx, bad)
	wantValidation(t, err, "category")
	bad = base
	bad.Size = 0
	_, err = svc.Materials.Create(ctx, bad)
	wantValidation(t, err, "size")
	bad = base
	bad.MediaID = ""
	_, err = svc.Materials.Create(ctx, bad)
	wantValidation(t, err, "required")
}

func TestVotesAndReactions(t *testing.T) {
	ctx := context.Background()
	_, svc := seeded()
	if err := svc.Votes.Cast(ctx, "", "s1", 0); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("want validation for empty client, got %v", err)
	}
	if err := svc.Votes.Cast(ctx, "c1", "s1", 3); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("want out of range, got %v", err)
	}
	if err := svc.Votes.Cast(ctx, "c1", "nope", 0); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("want not found, got %v", err)
	}
	_ = svc.Votes.Cast(ctx, "c1", "s1", 1)
	_ = svc.Votes.Cast(ctx, "c1", "s1", 2) // change vote
	_ = svc.Votes.Cast(ctx, "c2", "s1", 2)
	tally, _ := svc.Votes.Tally(ctx, "s1")
	if tally[2] != 2 || tally[1] != 0 {
		t.Fatalf("tally %v", tally)
	}
	if err := svc.Reactions.Set(ctx, "c1", "missing", true); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("want not found, got %v", err)
	}
	_ = svc.Reactions.Set(ctx, "c1", "e1", true)
	st, _ := svc.State.Get(ctx, "c1")
	if len(st.Liked) != 1 || st.Votes["s1"] != 2 {
		t.Fatalf("state %+v", st)
	}
	_ = svc.Reactions.Set(ctx, "c1", "e1", false)
	st, _ = svc.State.Get(ctx, "c1")
	if len(st.Liked) != 0 {
		t.Fatalf("like not removed: %+v", st.Liked)
	}
}

func TestStateShape(t *testing.T) {
	ctx := context.Background()
	svc := newTestServices(newMem())
	st, err := svc.State.Get(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if st.Version != 1 || st.Sessions == nil || st.Entries == nil || st.Replies == nil || st.Materials == nil || st.Votes == nil || st.Liked == nil {
		t.Fatalf("nil slices leak to the client: %+v", st)
	}
}

func TestExportSnapshot(t *testing.T) {
	ctx := context.Background()
	_, svc := seeded()
	_ = svc.Votes.Cast(ctx, "c1", "s1", 2)
	_ = svc.Votes.Cast(ctx, "c2", "s1", 2)
	_ = svc.Votes.Cast(ctx, "c3", "s1", 0)
	ex, err := svc.Export.Snapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if ex.ExportedAt != "2026-11-11T12:00:00Z" || len(ex.Sessions) != 1 {
		t.Fatalf("bad export %+v", ex)
	}
	s := ex.Sessions[0]
	if len(s.Tally) != 3 || s.Tally[2] != 2 || s.Tally[0] != 1 || s.TotalVotes != 3 || s.EntryCount != 1 {
		t.Fatalf("bad tally %+v", s)
	}
}
