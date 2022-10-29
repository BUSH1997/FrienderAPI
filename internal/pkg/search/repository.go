package search

import "context"

type Repository interface {
	GetEventUIDs(ctx context.Context, terms []string) ([]string, error)
}
