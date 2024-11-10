package service

import (
	"context"
	"time"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
)

type SlugDao interface {
	Fetch(context.Context, uuid.UUID) (model.Slug, error)
}

type SlugHistoryDao interface {
	List(context.Context, uuid.UUID, model.ListingFilters) ([]model.SlugHistory, model.Paging, error)
}

type SlugService struct {
	slugDao        SlugDao
	slugHistoryDao SlugHistoryDao
	lc             LocalCache
	producer       KafkaProducer
	timeLocation   *time.Location
}

func NewSlugService(slugDao SlugDao, slugHistoryDao SlugHistoryDao, lc LocalCache, producer KafkaProducer, timeLocation *time.Location) SlugService {
	return SlugService{
		slugDao:        slugDao,
		slugHistoryDao: slugHistoryDao,
		lc:             lc,
		producer:       producer,
		timeLocation:   timeLocation,
	}
}

func (s SlugService) Fetch(ctx context.Context, id uuid.UUID) (model.Slug, error) {
	value, exists := s.lc.Get(id.String())
	if exists {
		slug, ok := value.(model.Slug)
		if ok {

			return slug, nil
		}
	}

	slug, err := s.slugDao.Fetch(ctx, id)
	if err != nil {

		return model.Slug{}, err
	}

	return slug, nil
}

func (s SlugService) Create(ctx context.Context, req model.SlugCreateRequest) error {
	slug := model.Slug{
		Id:        uuid.New(),
		Name:      req.Name,
		Status:    string(model.Active),
		Cost:      req.Cost,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.CreatedBy,
		CreatedAt: time.Now().In(s.timeLocation),
		UpdatedAt: time.Now().In(s.timeLocation),
	}

	err := s.producer.Send(slug.Id.String(), slug)
	if err != nil {

		return err
	}

	s.lc.Set(slug.Id.String(), slug)

	return nil
}

func (s SlugService) ListHistory(ctx context.Context, id uuid.UUID, filters model.ListingFilters) ([]model.SlugHistory, model.Paging, error) {

	return s.slugHistoryDao.List(ctx, id, filters)
}
