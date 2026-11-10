// Package service is the business layer. It validates input, enforces rules,
// and orchestrates repositories. It knows nothing about HTTP or SQL.
package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/starman011/urbanism/backend/internal/domain"
)

const (
	maxTitle   = 120
	maxText    = 4000
	maxName    = 80
	maxTags    = 8
	maxOptions = 6
	maxNotes   = 12
)

// Services bundles every use-case group for the controller layer.
type Services struct {
	Sessions  *SessionService
	Entries   *EntryService
	Replies   *ReplyService
	Materials *MaterialService
	Votes     *VoteService
	Reactions *ReactionService
	Journal   *JournalService
	State     *StateService
}

func New(r domain.Repositories, now func() time.Time) *Services {
	if now == nil {
		now = time.Now
	}
	s := &Services{
		Sessions:  &SessionService{repo: r.Sessions, now: now},
		Entries:   &EntryService{repo: r.Entries, sessions: r.Sessions, now: now},
		Replies:   &ReplyService{repo: r.Replies, entries: r.Entries, now: now},
		Materials: &MaterialService{repo: r.Materials, now: now},
		Votes:     &VoteService{repo: r.Votes, sessions: r.Sessions},
		Reactions: &ReactionService{repo: r.Reactions, entries: r.Entries},
		Journal:   &JournalService{repo: r.Journal},
	}
	s.State = &StateService{r: r}
	return s
}

// ---- helpers ----

func clean(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) > max {
		s = s[:max]
	}
	return s
}

func cleanList(in []string, max, maxLen int) []string {
	out := make([]string, 0, len(in))
	seen := map[string]bool{}
	for _, v := range in {
		v = clean(v, maxLen)
		if v == "" || seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
		if len(out) == max {
			break
		}
	}
	return out
}

func newID(prefix string) string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("%s_%x", prefix, b)
}

func isDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

func normalizeCode(code string) string {
	var sb strings.Builder
	for _, r := range strings.ToUpper(strings.TrimSpace(code)) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// ---- Sessions ----

type SessionService struct {
	repo domain.SessionRepository
	now  func() time.Time
}

func (s *SessionService) List(ctx context.Context) ([]domain.Session, error) {
	return s.repo.List(ctx)
}

func (s *SessionService) GetByCode(ctx context.Context, code string) (*domain.Session, error) {
	code = normalizeCode(code)
	if code == "" {
		return nil, domain.Invalid("session code is required")
	}
	return s.repo.GetByCode(ctx, code)
}

func (s *SessionService) Create(ctx context.Context, in domain.Session) (*domain.Session, error) {
	in.Title = clean(in.Title, maxTitle)
	in.Code = normalizeCode(in.Code)
	if in.Title == "" {
		return nil, domain.Invalid("title is required")
	}
	if len(in.Code) < 3 || len(in.Code) > 12 {
		return nil, domain.Invalid("code must be 3 to 12 letters or digits")
	}
	if !isDate(in.Date) {
		return nil, domain.Invalid("date must be YYYY-MM-DD")
	}
	if in.EndDate == "" {
		in.EndDate = in.Date
	}
	if !isDate(in.EndDate) || in.EndDate < in.Date {
		return nil, domain.Invalid("endDate must be YYYY-MM-DD and not before date")
	}
	in.Place = clean(in.Place, maxTitle)
	in.Neighborhood = clean(in.Neighborhood, maxTitle)
	in.Desc = clean(in.Desc, maxText)
	in.Art = clean(in.Art, 32)
	if in.Art == "" {
		in.Art = "blocks"
	}
	in.Tags = cleanList(in.Tags, maxTags, 40)
	in.Notes = cleanList(in.Notes, maxNotes, maxText)
	in.Poll = clean(in.Poll, maxTitle)
	in.Options = cleanList(in.Options, maxOptions, 80)
	if in.Poll != "" && len(in.Options) < 2 {
		return nil, domain.Invalid("a poll needs at least two options")
	}
	if in.ID == "" {
		in.ID = newID("s")
	}
	in.Sample = false
	in.CreatedAt = s.now()
	if err := s.repo.Create(ctx, &in); err != nil {
		return nil, err
	}
	return &in, nil
}

// ---- Entries ----

type EntryService struct {
	repo     domain.EntryRepository
	sessions domain.SessionRepository
	now      func() time.Time
}

func (s *EntryService) List(ctx context.Context, sessionID string) ([]domain.Entry, error) {
	return s.repo.List(ctx, sessionID)
}

func (s *EntryService) Create(ctx context.Context, in domain.Entry) (*domain.Entry, error) {
	in.Text = clean(in.Text, maxText)
	if in.Text == "" {
		return nil, domain.Invalid("text is required")
	}
	if in.SessionID == "" {
		return nil, domain.Invalid("sessionId is required")
	}
	switch in.Type {
	case domain.EntryText, domain.EntryVoice, domain.EntryPhoto:
	case "":
		in.Type = domain.EntryText
	default:
		return nil, domain.Invalid("type must be text, voice, or photo")
	}
	in.Name = clean(in.Name, maxName)
	if in.Name == "" {
		in.Name = "Anonymous neighbor"
	}
	in.Hood = clean(in.Hood, maxName)
	in.Location = clean(in.Location, maxTitle)
	in.Art = clean(in.Art, 32)
	in.Tags = cleanList(in.Tags, maxTags, 40)
	if in.ID == "" {
		in.ID = newID("e")
	}
	if _, err := time.Parse(time.RFC3339, in.Created); err != nil {
		in.Created = s.now().UTC().Format(time.RFC3339)
	}
	if err := s.repo.Create(ctx, &in); err != nil {
		return nil, err
	}
	return &in, nil
}

// ---- Replies ----

type ReplyService struct {
	repo    domain.ReplyRepository
	entries domain.EntryRepository
	now     func() time.Time
}

func (s *ReplyService) List(ctx context.Context) ([]domain.Reply, error) { return s.repo.List(ctx) }

func (s *ReplyService) Create(ctx context.Context, in domain.Reply) (*domain.Reply, error) {
	in.Text = clean(in.Text, maxText)
	if in.Text == "" {
		return nil, domain.Invalid("text is required")
	}
	ok, err := s.entries.Exists(ctx, in.EntryID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrNotFound
	}
	in.Name = clean(in.Name, maxName)
	if in.Name == "" {
		in.Name = "Anonymous neighbor"
	}
	if in.ID == "" {
		in.ID = newID("r")
	}
	if _, err := time.Parse(time.RFC3339, in.Created); err != nil {
		in.Created = s.now().UTC().Format(time.RFC3339)
	}
	if err := s.repo.Create(ctx, &in); err != nil {
		return nil, err
	}
	return &in, nil
}

// ---- Materials (metadata only in phase 1) ----

type MaterialService struct {
	repo domain.MaterialRepository
	now  func() time.Time
}

func (s *MaterialService) List(ctx context.Context) ([]domain.Material, error) { return s.repo.List(ctx) }

func (s *MaterialService) Create(ctx context.Context, in domain.Material) (*domain.Material, error) {
	in.Title = clean(in.Title, maxTitle)
	if in.Title == "" || in.SessionID == "" || in.FileName == "" || in.MimeType == "" || in.MediaID == "" {
		return nil, domain.Invalid("title, sessionId, fileName, mimeType and mediaId are required")
	}
	valid := false
	for _, c := range domain.MaterialCategories {
		if c == in.Category {
			valid = true
		}
	}
	if !valid {
		return nil, domain.Invalid("category must be panel, slides, handout, report, or other")
	}
	if in.Size <= 0 {
		return nil, domain.Invalid("size must be positive")
	}
	in.Description = clean(in.Description, maxText)
	in.CreditedTo = clean(in.CreditedTo, maxName)
	if in.ID == "" {
		in.ID = newID("m")
	}
	if _, err := time.Parse(time.RFC3339, in.Created); err != nil {
		in.Created = s.now().UTC().Format(time.RFC3339)
	}
	if err := s.repo.Create(ctx, &in); err != nil {
		return nil, err
	}
	return &in, nil
}

// ---- Votes ----

type VoteService struct {
	repo     domain.VoteRepository
	sessions domain.SessionRepository
}

func (s *VoteService) Cast(ctx context.Context, clientID, sessionID string, option int) error {
	if clientID == "" {
		return domain.Invalid("client id is required")
	}
	sessions, err := s.sessions.List(ctx)
	if err != nil {
		return err
	}
	for _, sess := range sessions {
		if sess.ID == sessionID {
			if option < 0 || option >= len(sess.Options) {
				return domain.Invalid("option out of range")
			}
			return s.repo.Upsert(ctx, domain.Vote{ClientID: clientID, SessionID: sessionID, Option: option})
		}
	}
	return domain.ErrNotFound
}

func (s *VoteService) Tally(ctx context.Context, sessionID string) (map[int]int, error) {
	return s.repo.Tally(ctx, sessionID)
}

// ---- Reactions ----

type ReactionService struct {
	repo    domain.ReactionRepository
	entries domain.EntryRepository
}

func (s *ReactionService) Set(ctx context.Context, clientID, entryID string, liked bool) error {
	if clientID == "" {
		return domain.Invalid("client id is required")
	}
	ok, err := s.entries.Exists(ctx, entryID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrNotFound
	}
	r := domain.Reaction{ClientID: clientID, EntryID: entryID}
	if liked {
		return s.repo.Add(ctx, r)
	}
	return s.repo.Remove(ctx, r)
}

// ---- Journal ----

type JournalService struct{ repo domain.JournalRepository }

func (s *JournalService) List(ctx context.Context) ([]domain.JournalPost, error) { return s.repo.List(ctx) }

// ---- Aggregate state ----

type StateService struct{ r domain.Repositories }

// Get builds the version-1 client document for one anonymous client.
func (s *StateService) Get(ctx context.Context, clientID string) (*domain.State, error) {
	st := &domain.State{Version: 1, Votes: map[string]int{}, Liked: []string{}}
	var err error
	if st.Sessions, err = s.r.Sessions.List(ctx); err != nil {
		return nil, err
	}
	if st.Entries, err = s.r.Entries.List(ctx, ""); err != nil {
		return nil, err
	}
	if st.Replies, err = s.r.Replies.List(ctx); err != nil {
		return nil, err
	}
	if st.Materials, err = s.r.Materials.List(ctx); err != nil {
		return nil, err
	}
	if clientID != "" {
		if st.Votes, err = s.r.Votes.ForClient(ctx, clientID); err != nil {
			return nil, err
		}
		if st.Liked, err = s.r.Reactions.ForClient(ctx, clientID); err != nil {
			return nil, err
		}
	}
	// The client validator requires arrays, never null.
	for i := range st.Sessions {
		nz(&st.Sessions[i].Tags); nz(&st.Sessions[i].Notes); nz(&st.Sessions[i].Options)
	}
	for i := range st.Entries {
		nz(&st.Entries[i].Tags)
	}
	nz(&st.Liked)
	if st.Replies == nil { st.Replies = []domain.Reply{} }
	if st.Materials == nil { st.Materials = []domain.Material{} }
	if st.Entries == nil { st.Entries = []domain.Entry{} }
	if st.Sessions == nil { st.Sessions = []domain.Session{} }
	if st.Votes == nil { st.Votes = map[string]int{} }
	return st, nil
}

func nz(p *[]string) {
	if *p == nil {
		*p = []string{}
	}
}
