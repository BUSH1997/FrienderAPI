package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
)

func (uc complaintUsecase) Create(ctx context.Context, complaint models.Complaint) error {
	err := uc.complaintRepository.Create(ctx, complaint)
	if err != nil {
		return errors.Wrap(err, "failed to create complaint in usecase")
	}

	return nil
}
