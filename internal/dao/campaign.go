package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type CampaignDao struct {
	pool *pgxpool.Pool
}

func NewCampaignDao(pool *pgxpool.Pool) CampaignDao {
	return CampaignDao{
		pool: pool,
	}
}

const allCampaignFields = `
	id, 
	merchant_id, 
	status,
	budget,
	created_by,
	updated_by, 
	created_at, 
	updated_at
`

func (c CampaignDao) Fetch(ctx context.Context, id uuid.UUID) (model.Campaign, error) {
	var campaign model.Campaign

	query := `SELECT ` + allCampaignFields + ` from campaign WHERE id = $1`

	row := c.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&campaign.Id, &campaign.MerchantId, &campaign.Status,
		&campaign.Budget, &campaign.CreatedBy, &campaign.UpdatedBy,
		&campaign.CreatedAt, &campaign.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {

			return model.Campaign{}, model.ErrNotFound
		}

		return model.Campaign{}, errors.Wrap(err, "Failed to fetch campaign in database")
	}

	return campaign, nil
}
