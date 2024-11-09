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

type OwnerDao struct {
	pool *pgxpool.Pool
}

func NewOwnerDao(pool *pgxpool.Pool) OwnerDao {
	return OwnerDao{
		pool: pool,
	}
}

const allOwnerFields = `
	id, 
	email, 
	status, 
	created_by, 
	updated_by, 
	created_at, 
	updated_at
`

func (o OwnerDao) Fetch(ctx context.Context, id uuid.UUID) (model.Owner, error) {
	var owner model.Owner

	query := `SELECT ` + allOwnerFields + ` from owner WHERE id = $1`

	row := o.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&owner.Id, &owner.Email, &owner.Status,
		&owner.CreatedBy, &owner.UpdatedBy,
		&owner.CreatedAt, &owner.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {

			return model.Owner{}, model.ErrNotFound
		}
		easyzap.Error(ctx, err, "failed to fetch owner in database")

		return model.Owner{}, errors.Wrap(err, "Failed to fetch owner in database")
	}

	return owner, nil
}
