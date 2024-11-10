package service

import (
	"context"
	"fmt"
	"time"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
)

type SpentDao interface {
	FetchByMerchantIdAndBucket(context.Context, uuid.UUID, string) (model.Spent, error)
}

type SpentService struct {
	spentDao     SpentDao
	producer     KafkaProducer
	timeLocation *time.Location
}

func NewSpentService(spentDao SpentDao, producer KafkaProducer, timeLocation *time.Location) SpentService {
	return SpentService{
		spentDao:     spentDao,
		producer:     producer,
		timeLocation: timeLocation,
	}
}

func (s SpentService) FetchByMerchantIdAndBucket(ctx context.Context, merchantId uuid.UUID, bucket string) (model.Spent, error) {
	spent, err := s.spentDao.FetchByMerchantIdAndBucket(ctx, merchantId, bucket)
	if err != nil {

		return model.Spent{}, err
	}

	return spent, nil
}

func (s SpentService) Create(ctx context.Context, req model.SpentEvent) error {
	req.EventTime = time.Now().In(s.timeLocation)
	fmt.Println(req.CampaignId)
	fmt.Println(req.MerchantId)
	fmt.Println(req.SessionId)
	fmt.Println(req.SlugName)
	fmt.Println(req.UserId)
	fmt.Println(req.EventType)
	fmt.Println(req.EventTime)
	fmt.Println(req.IP)

	err := s.producer.Send(req.SessionId.String(), req)
	if err != nil {

		return err
	}

	return nil
}
