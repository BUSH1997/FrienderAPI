package status

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type Repository interface {
	GetUserCurrentStatus(ctx context.Context, id int) (models.Status, error)
	GetAllUserStatuses(ctx context.Context, id int) ([]models.Status, error)
}
