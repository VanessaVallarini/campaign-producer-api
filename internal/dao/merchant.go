package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	easyzap "github.com/lockp111/go-easyzap"
	"github.com/pkg/errors"
)

type MerchantDao struct {
	pool *pgxpool.Pool
}

func NewMerchantDao(pool *pgxpool.Pool) MerchantDao {

	return MerchantDao{
		pool: pool,
	}
}

const allMerchantFields = `
	id, 
	owner_id, 
	region_id, 
	slugs, 
	name,
	status,
	created_by,
	updated_by, 
	created_at, 
	updated_at
`

func (m MerchantDao) Fetch(ctx context.Context, id uuid.UUID) (model.Merchant, error) {
	var merchant model.Merchant

	query := `SELECT ` + allMerchantFields + ` from merchant WHERE id = $1`

	row := m.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&merchant.Id, &merchant.OwnerId, &merchant.RegionId, &merchant.Slugs,
		&merchant.Name, &merchant.Status, &merchant.CreatedBy,
		&merchant.UpdatedBy, &merchant.CreatedAt, &merchant.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {

			return model.Merchant{}, model.ErrNotFound
		}
		easyzap.Error(ctx, err, "failed to fetch merchant in database")

		return model.Merchant{}, errors.Wrap(err, "Failed to fetch merchant in database")
	}

	return merchant, nil
}
