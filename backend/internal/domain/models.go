// Package domain holds the business entities and the repository contracts.
// It has no dependency on HTTP, SQL, or any framework.
package domain

import "time"

// Session is an assembly, walk, lab, or forum that neighbors can respond to.
type Session struct {
	ID           string   `json:"id"`
	Code         string   `json:"code"`
	Title        string   `json:"title"`
	Date         string   `json:"date"`    // YYYY-MM-DD (America/New_York)
	EndDate      string   `json:"endDate"` // YYYY-MM-DD
	Place        string   `json:"place"`
	Neighborhood string   `json:"neighborhood"`
	Desc         string   `json:"desc"`
	Art          string   `json:"art"`
	Tags         []string `json:"tags"`
	Notes        []string `json:"notes"`
	Poll         string   `json:"poll"`
	Options      []string `json:"options"`
	Sample       bool     `json:"sample"`
	CreatedAt    time.Time `json:"-"`
}

// EntryType is the medium of a field note.
type EntryType string

const (
	EntryText  EntryType = "text"
	EntryVoice EntryType = "voice"
	EntryPhoto EntryType = "photo"
)

// Entry is a single field note left by a neighbor on a session.
type Entry struct {
	ID        string    `json:"id"`
	SessionID string    `json:"sessionId"`
	Type      EntryType `json:"type"`
	Name      string    `json:"name"`
	Hood      string    `json:"hood"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags"`
	Location  string    `json:"location,omitempty"`
	Art       string    `json:"art,omitempty"`
	Created   string    `json:"created"` // RFC3339, kept as string to match the client contract
}

// Reply is a response to an entry.
type Reply struct {
	ID      string `json:"id"`
	EntryID string `json:"entryId"`
	Name    string `json:"name"`
	Text    string `json:"text"`
	Created string `json:"created"`
}

// MaterialCategory classifies uploaded session materials.
type MaterialCategory string

var MaterialCategories = []MaterialCategory{"panel", "slides", "handout", "report", "other"}

// Material is metadata about a file shared for a session. Bytes live in object
// storage (phase 2); phase 1 keeps bytes in the browser and syncs metadata only.
type Material struct {
	ID          string           `json:"id"`
	SessionID   string           `json:"sessionId"`
	Title       string           `json:"title"`
	Category    MaterialCategory `json:"category"`
	Description string           `json:"description"`
	CreditedTo  string           `json:"creditedTo"`
	FileName    string           `json:"fileName"`
	MimeType    string           `json:"mimeType"`
	MediaID     string           `json:"mediaId"`
	Size        int64            `json:"size"`
	Created     string           `json:"created"`
}

// Vote records one client's choice on a session poll.
type Vote struct {
	ClientID  string
	SessionID string
	Option    int
}

// Reaction records one client's appreciation of an entry.
type Reaction struct {
	ClientID string
	EntryID  string
}

// JournalPost is an editorial piece published by the committee.
type JournalPost struct {
	ID          string   `json:"id"`
	Number      string   `json:"number"`
	Category    string   `json:"category"`
	Title       string   `json:"title"`
	Excerpt     string   `json:"excerpt"`
	Date        string   `json:"date"`
	ReadTime    string   `json:"readTime"`
	Art         string   `json:"art"`
	Author      string   `json:"author"`
	SessionCode string   `json:"sessionCode,omitempty"`
	Body        []string `json:"body"`
}

// State is the aggregate document the client boots from. Its shape mirrors the
// prototype's localStorage document (version 1) so the UI needs no changes.
type State struct {
	Version   int                 `json:"version"`
	Sessions  []Session           `json:"sessions"`
	Entries   []Entry             `json:"entries"`
	Replies   []Reply             `json:"replies"`
	Materials []Material          `json:"materials"`
	Votes     map[string]int      `json:"votes"` // sessionId -> option index, for the requesting client
	Liked     []string            `json:"liked"` // entry ids, for the requesting client
}
