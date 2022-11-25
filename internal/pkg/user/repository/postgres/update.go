package postgres

import (
	"context"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
)

func (r UserRepository) UpdateRefresh(ctx context.Context, token models.RefreshToken) error {
	res := r.db.Model(&db_models.AuthUser{}).Where("uid = ?", contextlib.GetUser(ctx)).
		Updates(map[string]interface{}{
			"refresh_token": token.Value,
			"fingerprint":   token.FingerPrint,
		})
	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}
