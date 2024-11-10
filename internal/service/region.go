package service

import (
	"context"
	"time"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
)

type RegionDao interface {
	Fetch(context.Context, uuid.UUID) (model.Region, error)
}

type RegionHistoryDao interface {
	List(context.Context, uuid.UUID, model.ListingFilters) ([]model.RegionHistory, model.Paging, error)
}

type RegionService struct {
	regionDao        RegionDao
	regionHistoryDao RegionHistoryDao
	lc               LocalCache
	producer         KafkaProducer
	timeLocation     *time.Location
}

func NewRegionService(regionDao RegionDao, regionHistoryDao RegionHistoryDao, lc LocalCache, producer KafkaProducer, timeLocation *time.Location) RegionService {
	return RegionService{
		regionDao:        regionDao,
		regionHistoryDao: regionHistoryDao,
		lc:               lc,
		producer:         producer,
		timeLocation:     timeLocation,
	}
}

func (r RegionService) Fetch(ctx context.Context, id uuid.UUID) (model.Region, error) {
	value, exists := r.lc.Get(id.String())
	if exists {
		region, ok := value.(model.Region)
		if ok {

			return region, nil
		}
	}

	region, err := r.regionDao.Fetch(ctx, id)
	if err != nil {

		return model.Region{}, err
	}

	return region, nil
}

func (r RegionService) Create(ctx context.Context, req model.RegionCreateRequest) error {
	region := model.Region{
		Id:        uuid.New(),
		Name:      req.Name,
		Status:    string(model.Active),
		Lat:       req.Lat,
		Long:      req.Long,
		Cost:      req.Cost,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.CreatedBy,
		CreatedAt: time.Now().In(r.timeLocation),
		UpdatedAt: time.Now().In(r.timeLocation),
	}

	err := r.producer.Send(region.Id.String(), region)
	if err != nil {

		return err
	}

	r.lc.Set(region.Id.String(), region)

	return nil
}

func (r RegionService) ListHistory(ctx context.Context, id uuid.UUID, filters model.ListingFilters) ([]model.RegionHistory, model.Paging, error) {

	return r.regionHistoryDao.List(ctx, id, filters)
}
