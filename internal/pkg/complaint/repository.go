package complaint

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type Repository interface {
	Create(ctx context.Context, event models.Complaint) error
}
