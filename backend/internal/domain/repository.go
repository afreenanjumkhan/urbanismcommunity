package domain

import "context"

// Repositories are the persistence contracts the service layer depends on.
// Implementations live under internal/repository and must not leak SQL types.

type SessionRepository interface {
	List(ctx context.Context) ([]Session, error)
	GetByCode(ctx context.Context, code string) (*Session, error)
	Create(ctx context.Context, s *Session) error
}

type EntryRepository interface {
	List(ctx context.Context, sessionID string) ([]Entry, error)
	Create(ctx context.Context, e *Entry) error
	Exists(ctx context.Context, id string) (bool, error)
}

type ReplyRepository interface {
	List(ctx context.Context) ([]Reply, error)
	Create(ctx context.Context, r *Reply) error
}

type MaterialRepository interface {
	List(ctx context.Context) ([]Material, error)
	Create(ctx context.Context, m *Material) error
}

type VoteRepository interface {
	ForClient(ctx context.Context, clientID string) (map[string]int, error)
	Upsert(ctx context.Context, v Vote) error
	Tally(ctx context.Context, sessionID string) (map[int]int, error)
}

type ReactionRepository interface {
	ForClient(ctx context.Context, clientID string) ([]string, error)
	Add(ctx context.Context, r Reaction) error
	Remove(ctx context.Context, r Reaction) error
}

type JournalRepository interface {
	List(ctx context.Context) ([]JournalPost, error)
}

// Repositories bundles every contract for dependency injection.
type Repositories struct {
	Sessions  SessionRepository
	Entries   EntryRepository
	Replies   ReplyRepository
	Materials MaterialRepository
	Votes     VoteRepository
	Reactions ReactionRepository
	Journal   JournalRepository
}
