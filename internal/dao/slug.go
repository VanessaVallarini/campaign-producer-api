package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type SlugDao struct {
	pool *pgxpool.Pool
}

func NewSlugDao(pool *pgxpool.Pool) SlugDao {

	return SlugDao{
		pool: pool,
	}
}

const allSlugFields = `
	id, 
	name, 
	status, 
	cost, 
	created_by, 
	updated_by, 
	created_at, 
	updated_at
`

func (sd SlugDao) Fetch(ctx context.Context, id uuid.UUID) (model.Slug, error) {
	var slug model.Slug

	query := `SELECT ` + allSlugFields + ` from slug WHERE id = $1`

	row := sd.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&slug.Id, &slug.Name, &slug.Status, &slug.Cost, &slug.CreatedBy,
		&slug.UpdatedBy, &slug.CreatedAt, &slug.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {

			return model.Slug{}, model.ErrNotFound
		}

		return model.Slug{}, errors.Wrap(err, "Failed to fetch slug in database")
	}

	return slug, nil
}
