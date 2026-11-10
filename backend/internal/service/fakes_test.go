package service

import (
	"context"
	"sync"
	"time"

	"github.com/starman011/urbanism/backend/internal/domain"
)

// In-memory repositories: the service layer is tested without a database.

type memRepos struct {
	mu        sync.Mutex
	sessions  []domain.Session
	entries   []domain.Entry
	replies   []domain.Reply
	materials []domain.Material
	votes     map[string]map[string]int // client -> session -> option
	likes     map[string]map[string]bool
	journal   []domain.JournalPost
}

func newMem() *memRepos {
	return &memRepos{votes: map[string]map[string]int{}, likes: map[string]map[string]bool{}}
}

func (m *memRepos) List(ctx context.Context) ([]domain.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]domain.Session(nil), m.sessions...), nil
}
func (m *memRepos) GetByCode(ctx context.Context, code string) (*domain.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range m.sessions {
		if s.Code == code {
			s := s
			return &s, nil
		}
	}
	return nil, domain.ErrNotFound
}
func (m *memRepos) Create(ctx context.Context, s *domain.Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, x := range m.sessions {
		if x.Code == s.Code {
			return domain.ErrConflict
		}
	}
	m.sessions = append(m.sessions, *s)
	return nil
}

func (m *memRepos) listEntries(sessionID string) []domain.Entry {
	out := []domain.Entry{}
	for _, e := range m.entries {
		if sessionID == "" || e.SessionID == sessionID {
			out = append(out, e)
		}
	}
	return out
}

// EntryRepository
func (m *memRepos) Exists(ctx context.Context, id string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.entries {
		if e.ID == id {
			return true, nil
		}
	}
	return false, nil
}

// Disambiguate overlapping method names via wrapper types.
type entryAdapter struct{ *memRepos }

func (a entryAdapter) List(ctx context.Context, sessionID string) ([]domain.Entry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.listEntries(sessionID), nil
}
func (a entryAdapter) Create(ctx context.Context, e *domain.Entry) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = append(a.entries, *e)
	return nil
}
func (a entryAdapter) Exists(ctx context.Context, id string) (bool, error) {
	return a.memRepos.Exists(ctx, id)
}

type replyAdapter struct{ *memRepos }

func (a replyAdapter) List(ctx context.Context) ([]domain.Reply, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]domain.Reply{}, a.replies...), nil
}
func (a replyAdapter) Create(ctx context.Context, r *domain.Reply) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.replies = append(a.replies, *r)
	return nil
}

type materialAdapter struct{ *memRepos }

func (a materialAdapter) List(ctx context.Context) ([]domain.Material, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]domain.Material{}, a.materials...), nil
}
func (a materialAdapter) Create(ctx context.Context, x *domain.Material) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.materials = append(a.materials, *x)
	return nil
}

// VoteRepository
func (m *memRepos) ForClient(ctx context.Context, clientID string) (map[string]int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := map[string]int{}
	for k, v := range m.votes[clientID] {
		out[k] = v
	}
	return out, nil
}
func (m *memRepos) Upsert(ctx context.Context, v domain.Vote) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.votes[v.ClientID] == nil {
		m.votes[v.ClientID] = map[string]int{}
	}
	m.votes[v.ClientID][v.SessionID] = v.Option
	return nil
}
func (m *memRepos) Tally(ctx context.Context, sessionID string) (map[int]int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := map[int]int{}
	for _, per := range m.votes {
		if o, ok := per[sessionID]; ok {
			out[o]++
		}
	}
	return out, nil
}

type reactionAdapter struct{ *memRepos }

func (a reactionAdapter) ForClient(ctx context.Context, clientID string) ([]string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := []string{}
	for id := range a.likes[clientID] {
		out = append(out, id)
	}
	return out, nil
}
func (a reactionAdapter) Add(ctx context.Context, r domain.Reaction) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.likes[r.ClientID] == nil {
		a.likes[r.ClientID] = map[string]bool{}
	}
	a.likes[r.ClientID][r.EntryID] = true
	return nil
}
func (a reactionAdapter) Remove(ctx context.Context, r domain.Reaction) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.likes[r.ClientID], r.EntryID)
	return nil
}

type journalAdapter struct{ *memRepos }

func (a journalAdapter) List(ctx context.Context) ([]domain.JournalPost, error) {
	return a.journal, nil
}

func (m *memRepos) all() domain.Repositories {
	return domain.Repositories{
		Sessions:  m,
		Entries:   entryAdapter{m},
		Replies:   replyAdapter{m},
		Materials: materialAdapter{m},
		Votes:     m,
		Reactions: reactionAdapter{m},
		Journal:   journalAdapter{m},
	}
}

func fixedNow() time.Time { return time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC) }

func newTestServices(m *memRepos) *Services { return New(m.all(), fixedNow) }
