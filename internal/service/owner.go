package service

import (
	"context"
	"time"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
)

type OwnerDao interface {
	Fetch(context.Context, uuid.UUID) (model.Owner, error)
}

type OwnerService struct {
	ownerDao     OwnerDao
	lc           LocalCache
	producer     KafkaProducer
	timeLocation *time.Location
}

func NewOwnerService(ownerDao OwnerDao, lc LocalCache, producer KafkaProducer, timeLocation *time.Location) OwnerService {
	return OwnerService{
		ownerDao:     ownerDao,
		lc:           lc,
		producer:     producer,
		timeLocation: timeLocation,
	}
}

func (o OwnerService) Fetch(ctx context.Context, id uuid.UUID) (model.Owner, error) {
	value, exists := o.lc.Get(id.String())
	if exists {
		owner, ok := value.(model.Owner)
		if ok {

			return owner, nil
		}
	}

	owner, err := o.ownerDao.Fetch(ctx, id)
	if err != nil {

		return model.Owner{}, err
	}

	return owner, nil
}

func (o OwnerService) Create(ctx context.Context, req model.OwnerCreateRequest) error {
	owner := model.Owner{
		Id:        uuid.New(),
		Email:     req.Email,
		Status:    string(model.Active),
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.CreatedBy,
		CreatedAt: time.Now().In(o.timeLocation),
		UpdatedAt: time.Now().In(o.timeLocation),
	}

	err := o.producer.Send(owner.Id.String(), owner)
	if err != nil {

		return err
	}

	o.lc.Set(owner.Id.String(), owner)

	return nil
}
