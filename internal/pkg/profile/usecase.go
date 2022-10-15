package profile

import "github.com/BUSH1997/FrienderAPI/internal/pkg/models"

type UseCase interface {
	GetOneProfile(id string) (models.Profile, error)
	GetAllStatusesUser(id string) ([]models.Status, error)
	ChangeProfile(profile models.ChangeProfile) error
	ChangePriorityEvent(eventPriority models.UidEventPriority) error
}
