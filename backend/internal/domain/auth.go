package domain

import "context"

type Account struct {
	Email        string `json:"email"`
	Role         string `json:"role"`
	Active       bool   `json:"active"`
	PasswordHash string `json:"-"`
}
type Accounts interface {
	Get(context.Context, string) (Account, error)
	List(context.Context) ([]Account, error)
	Invite(context.Context, string, string, string) error
	Activate(context.Context, string, string, string, bool) error
	Revoke(context.Context, string, string) error
}
