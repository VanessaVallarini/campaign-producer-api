package service

import (
	"context"
	"time"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
)

type MerchantDao interface {
	Fetch(context.Context, uuid.UUID) (model.Merchant, error)
}

type MerchantService struct {
	merchantDao  MerchantDao
	lc           LocalCache
	producer     KafkaProducer
	timeLocation *time.Location
}

func NewMerchantService(merchantDao MerchantDao, lc LocalCache, producer KafkaProducer, timeLocation *time.Location) MerchantService {
	return MerchantService{
		merchantDao:  merchantDao,
		lc:           lc,
		producer:     producer,
		timeLocation: timeLocation,
	}
}

func (m MerchantService) Fetch(ctx context.Context, id uuid.UUID) (model.Merchant, error) {
	value, exists := m.lc.Get(id.String())
	if exists {
		merchant, ok := value.(model.Merchant)
		if ok {

			return merchant, nil
		}
	}

	merchant, err := m.merchantDao.Fetch(ctx, id)
	if err != nil {

		return model.Merchant{}, err
	}

	return merchant, nil
}

func (r MerchantService) Create(ctx context.Context, req model.MerchantCreateRequest) error {
	merchant := model.Merchant{
		Id:        uuid.New(),
		OwnerId:   req.OwnerId,
		RegionId:  req.RegionId,
		Slugs:     req.Slugs,
		Name:      req.Name,
		Status:    string(model.Active),
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.CreatedBy,
		CreatedAt: time.Now().In(r.timeLocation),
		UpdatedAt: time.Now().In(r.timeLocation),
	}

	err := r.producer.Send(merchant.Id.String(), merchant)
	if err != nil {

		return err
	}

	r.lc.Set(merchant.Id.String(), merchant)

	return nil
}
