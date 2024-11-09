package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
)

type CampaignDao interface {
	Fetch(context.Context, uuid.UUID) (model.Campaign, error)
}

type CampaignHistoryDao interface {
	List(context.Context, uuid.UUID, model.ListingFilters) ([]model.CampaignHistory, model.Paging, error)
}

type CampaignService struct {
	campaignDao        CampaignDao
	campaignHistoryDao CampaignHistoryDao
	lc                 LocalCache
}

func NewCampaignService(campaignDao CampaignDao, campaignHistoryDao CampaignHistoryDao, lc LocalCache) CampaignService {
	return CampaignService{
		campaignDao:        campaignDao,
		campaignHistoryDao: campaignHistoryDao,
		lc:                 lc,
	}
}

func (c CampaignService) Fetch(ctx context.Context, id uuid.UUID) (model.Campaign, error) {
	value, exists := c.lc.Get(id.String())
	if exists {
		campaign, ok := value.(model.Campaign)
		if ok {

			return campaign, nil
		}
	}

	campaign, err := c.campaignDao.Fetch(ctx, id)
	if err != nil {

		return model.Campaign{}, err
	}

	c.lc.Set(id.String(), campaign)

	return campaign, nil
}

func (c CampaignService) List(ctx context.Context, campaignId uuid.UUID, filters model.ListingFilters) ([]model.CampaignHistory, model.Paging, error) {

	return c.campaignHistoryDao.List(ctx, campaignId, filters)
}
