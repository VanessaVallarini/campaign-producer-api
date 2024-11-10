package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type SpentDao struct {
	pool *pgxpool.Pool
}

func NewSpentDao(pool *pgxpool.Pool) SpentDao {
	return SpentDao{
		pool: pool,
	}
}

const allSpentFields = `
	id,
	campaign_id,
	merchant_id,
	bucket,
	total_spent,
	total_clicks,
	total_impressions
`

func (sd SpentDao) FetchByMerchantIdAndBucket(ctx context.Context, merchantId uuid.UUID, bucket string) (model.Spent, error) {
	var spent model.Spent

	query := `SELECT` + allSpentFields + ` from spent WHERE merchant_id = $1 AND bucket = $2`

	row := sd.pool.QueryRow(ctx, query, merchantId, bucket)
	err := row.Scan(
		&spent.Id, &spent.CampaignId, &spent.MerchantId, &spent.Bucket,
		&spent.TotalSpent, &spent.TotalClicks, &spent.TotalImpressions,
	)

	if err != nil {
		if err == pgx.ErrNoRows {

			return model.Spent{}, model.ErrNotFound
		}

		return model.Spent{}, errors.Wrap(err, "Failed to fetch spent by merchantId id in database")
	}

	return spent, nil
}
