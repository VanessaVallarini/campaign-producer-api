package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type RegionDao struct {
	pool *pgxpool.Pool
}

func NewRegionDao(pool *pgxpool.Pool) RegionDao {
	return RegionDao{
		pool: pool,
	}
}

const allRegionFields = `
	id,
	name,
	status,
	lat,
	long,
	cost,
	created_by,
	updated_by,
	created_at,
	updated_at
`

func (r RegionDao) Fetch(ctx context.Context, id uuid.UUID) (model.Region, error) {
	var region model.Region

	query := `SELECT ` + allRegionFields + ` from region WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&region.Id, &region.Name, &region.Status, &region.Lat,
		&region.Long, &region.Cost, &region.CreatedBy,
		&region.UpdatedBy, &region.CreatedAt, &region.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {

			return model.Region{}, model.ErrNotFound
		}

		return model.Region{}, errors.Wrap(err, "Failed to fetch region in database")
	}

	return region, nil
}
