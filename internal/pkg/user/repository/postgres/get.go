package postgres

import (
	"context"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
)

func (r UserRepository) GetRefresh(ctx context.Context) (models.RefreshToken, error) {
	dbUser := db_models.AuthUser{}
	res := r.db.Take(&dbUser, "uid = ?", contextlib.GetUser(ctx))
	if err := res.Error; err != nil {
		return models.RefreshToken{}, errors.Wrap(err, "failed to get user")
	}

	return models.RefreshToken{
		Value: dbUser.RefreshToken,
	}, nil
}
