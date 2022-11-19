package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/complaint"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
)

type complaintUsecase struct {
	complaintRepository complaint.Repository
	logger              hardlogger.Logger
}

func New(
	complaintRepository complaint.Repository,
	logger hardlogger.Logger,
) complaint.Usecase {
	return &complaintUsecase{
		complaintRepository: complaintRepository,
		logger:              logger,
	}
}
