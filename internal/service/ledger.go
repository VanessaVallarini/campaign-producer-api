package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
)

type LedgerDao interface {
	List(context.Context, uuid.UUID, model.ListingFilters) ([]model.Ledger, model.Paging, error)
}

type LedgerService struct {
	ledgerDao LedgerDao
}

func NewLedgerService(ledgerDao LedgerDao) LedgerService {
	return LedgerService{
		ledgerDao: ledgerDao,
	}
}

func (l LedgerService) List(ctx context.Context, campaignId uuid.UUID, filters model.ListingFilters) ([]model.Ledger, model.Paging, error) {

	return l.ledgerDao.List(ctx, campaignId, filters)
}
