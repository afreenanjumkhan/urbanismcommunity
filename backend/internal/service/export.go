package service

import (
	"context"
	"time"

	"github.com/starman011/urbanism/backend/internal/domain"
)

// Export is the research snapshot: the whole shared record plus poll tallies.
// Personal slices and client ids are never included.
type Export struct {
	ExportedAt string            `json:"exportedAt"`
	Sessions   []SessionExport   `json:"sessions"`
	Entries    []domain.Entry    `json:"entries"`
	Replies    []domain.Reply    `json:"replies"`
	Materials  []domain.Material `json:"materials"`
}

type SessionExport struct {
	domain.Session
	Tally      []int `json:"tally"` // votes per option index
	TotalVotes int   `json:"totalVotes"`
	EntryCount int   `json:"entryCount"`
}

type ExportService struct {
	r   domain.Repositories
	now func() time.Time
}

func (s *ExportService) Snapshot(ctx context.Context) (*Export, error) {
	st, err := (&StateService{r: s.r}).Get(ctx, "")
	if err != nil {
		return nil, err
	}
	counts := map[string]int{}
	for _, e := range st.Entries {
		counts[e.SessionID]++
	}
	out := &Export{ExportedAt: s.now().UTC().Format(time.RFC3339), Entries: st.Entries, Replies: st.Replies, Materials: st.Materials, Sessions: []SessionExport{}}
	for _, sess := range st.Sessions {
		tally, err := s.r.Votes.Tally(ctx, sess.ID)
		if err != nil {
			return nil, err
		}
		se := SessionExport{Session: sess, Tally: make([]int, len(sess.Options)), EntryCount: counts[sess.ID]}
		for i := range sess.Options {
			se.Tally[i] = tally[i]
			se.TotalVotes += tally[i]
		}
		out.Sessions = append(out.Sessions, se)
	}
	return out, nil
}
