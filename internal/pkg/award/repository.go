package award

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type Repository interface {
	GetUserAwards(ctx context.Context, id int) ([]models.Award, error)
}
