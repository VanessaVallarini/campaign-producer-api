package dao

import (
	"context"
	"fmt"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type RegionHistoryDao struct {
	pool *pgxpool.Pool
}

func NewRegionHistoryDao(pool *pgxpool.Pool) RegionHistoryDao {
	return RegionHistoryDao{
		pool: pool,
	}
}

const allRegionHistoryFields = `
	id, 
	region_id, 
	status,
	description,
	created_by,
	created_at
`

func (rh RegionHistoryDao) List(ctx context.Context, regionId uuid.UUID, filters model.ListingFilters) ([]model.RegionHistory, model.Paging, error) {
	dataQuery := `SELECT ` + allRegionHistoryFields + ` from region_history WHERE 1=1 and region_id = $1`
	countQuery := `SELECT count(id) as "total" from region_history WHERE 1=1 and region_id = $1`

	paramCount := 2
	extraParams := make([]interface{}, 0, 6)
	extraParams = append(extraParams, regionId)

	if string(filters.Status) != "" {
		extra := fmt.Sprintf(" and status = $%d ", paramCount)
		paramCount++
		dataQuery, countQuery = dataQuery+extra, countQuery+extra
		extraParams = append(extraParams, filters.Status)
	}

	if filters.StartDate != "" && filters.EndDate != "" {
		extra := fmt.Sprintf(" and DATE(created_at) BETWEEN $%d AND $%d ", paramCount, paramCount+1)
		paramCount += 2
		dataQuery, countQuery = dataQuery+extra, countQuery+extra
		extraParams = append(extraParams, filters.StartDate, filters.EndDate)
	}

	if filters.StartDate != "" && filters.EndDate == "" {
		extra := fmt.Sprintf(" and DATE(created_at) >= $%d ", paramCount)
		paramCount++
		dataQuery, countQuery = dataQuery+extra, countQuery+extra
		extraParams = append(extraParams, filters.StartDate)
	}

	if filters.EndDate != "" && filters.StartDate == "" {
		extra := fmt.Sprintf(" and DATE(created_at) <= $%d ", paramCount)
		paramCount++
		dataQuery, countQuery = dataQuery+extra, countQuery+extra
		extraParams = append(extraParams, filters.EndDate)
	}

	// default limit size
	if filters.Size == 0 {
		filters.Size = 50
	}

	// Select total count
	type counting struct {
		Total int `db:"total"`
	}
	countingResult := counting{}
	countRow := rh.pool.QueryRow(ctx, countQuery, extraParams...)
	err := countRow.Scan(&countingResult.Total)
	if err != nil {

		return []model.RegionHistory{}, model.Paging{}, errors.Wrap(err, "Failed to get total items count for region_history table")
	}

	// pagination params
	dataQuery = dataQuery + fmt.Sprintf(" order by created_at desc limit $%d offset $%d", paramCount, paramCount+1)
	paramCount += 2
	extraParams = append(extraParams, filters.Size, filters.Size*filters.Page)

	// Select paged data
	result := []model.RegionHistory{}
	rows, err := rh.pool.Query(ctx, dataQuery, extraParams...)
	if err != nil {

		return []model.RegionHistory{}, model.Paging{}, errors.Wrap(err, "Failed to get region_history list data")
	}
	defer rows.Close()

	for rows.Next() {
		var history model.RegionHistory
		err := rows.Scan(&history.Id, &history.RegionId, &history.Status, &history.Description, &history.CreatedBy, &history.CreatedAt)
		if err != nil {

			return nil, model.Paging{}, errors.Wrap(err, "Failed to parse history when listing by region")
		}

		result = append(result, history)
	}
	if err := rows.Err(); err != nil {

		return nil, model.Paging{}, errors.Wrap(err, "Failed to read row when listing regions data")
	}

	return result, model.Paging{Page: filters.Page, Size: len(result), TotalItems: countingResult.Total}, nil
}
