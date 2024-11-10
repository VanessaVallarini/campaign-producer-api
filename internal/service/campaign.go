package service

import (
	"context"
	"time"

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
	producer           KafkaProducer
	timeLocation       *time.Location
}

func NewCampaignService(campaignDao CampaignDao, campaignHistoryDao CampaignHistoryDao, lc LocalCache, producer KafkaProducer, timeLocation *time.Location) CampaignService {
	return CampaignService{
		campaignDao:        campaignDao,
		campaignHistoryDao: campaignHistoryDao,
		lc:                 lc,
		producer:           producer,
		timeLocation:       timeLocation,
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

	return campaign, nil
}

func (c CampaignService) Create(ctx context.Context, req model.CampaignCreateRequest) error {
	campaign := model.Campaign{
		Id:         uuid.New(),
		MerchantId: req.MerchantId,
		Status:     string(model.Active),
		Budget:     req.Budget,
		CreatedBy:  req.CreatedBy,
		UpdatedBy:  req.CreatedBy,
		CreatedAt:  time.Now().In(c.timeLocation),
		UpdatedAt:  time.Now().In(c.timeLocation),
	}

	err := c.producer.Send(campaign.Id.String(), campaign)
	if err != nil {

		return err
	}

	c.lc.Set(campaign.Id.String(), campaign)

	return nil
}

func (c CampaignService) ListHistory(ctx context.Context, campaignId uuid.UUID, filters model.ListingFilters) ([]model.CampaignHistory, model.Paging, error) {

	return c.campaignHistoryDao.List(ctx, campaignId, filters)
}
