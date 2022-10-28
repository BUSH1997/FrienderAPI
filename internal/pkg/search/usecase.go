package search

import "github.com/BUSH1997/FrienderAPI/internal/pkg/models"

type UseCase interface {
	Search(words []string) ([]models.Event, error)
}
