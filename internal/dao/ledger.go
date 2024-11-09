package dao

import (
	"context"
	"fmt"

	"github.com/VanessaVallarini/campaign-producer-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type LedgerDao struct {
	pool *pgxpool.Pool
}

func NewLedgerDao(pool *pgxpool.Pool) LedgerDao {
	return LedgerDao{
		pool: pool,
	}
}

const allLedgerFields = `
	id,
	spent_id,
	campaign_id,
	merchant_id,
	slug_name,
	region_name,
	user_id,
	session_id,
	event_type,
	cost,
	ip,
	lat,
	long, 
	created_at, 
	event_time
`

func (l LedgerDao) List(ctx context.Context, campaignId uuid.UUID, filters model.ListingFilters) ([]model.Ledger, model.Paging, error) {
	dataQuery := `SELECT ` + allLedgerFields + ` from ledger WHERE 1=1 and campaign_id = $1`
	countQuery := `SELECT count(id) as "total" from ledger WHERE 1=1 and campaign_id = $1`

	paramCount := 2
	extraParams := make([]interface{}, 0, 6)
	extraParams = append(extraParams, campaignId)

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
	countRow := l.pool.QueryRow(ctx, countQuery, extraParams...)
	err := countRow.Scan(&countingResult.Total)
	if err != nil {

		return []model.Ledger{}, model.Paging{}, errors.Wrap(err, "Failed to get total items count for ledger table")
	}

	// pagination params
	dataQuery = dataQuery + fmt.Sprintf(" order by created_at desc limit $%d offset $%d", paramCount, paramCount+1)
	paramCount += 2
	extraParams = append(extraParams, filters.Size, filters.Size*filters.Page)

	// Select paged data
	result := []model.Ledger{}
	rows, err := l.pool.Query(ctx, dataQuery, extraParams...)
	if err != nil {

		return []model.Ledger{}, model.Paging{}, errors.Wrap(err, "Failed to get ledger list data")
	}
	defer rows.Close()

	for rows.Next() {
		var ledger model.Ledger
		err := rows.Scan(&ledger.SpentId, &ledger.CampaignId, &ledger.MerchantId, &ledger.SlugName,
			&ledger.RegionName, &ledger.UserId, &ledger.SessionId, &ledger.EventType,
			&ledger.Cost, &ledger.Ip, &ledger.Lat, &ledger.Long, &ledger.CreatedAt, &ledger.EventTime)
		if err != nil {

			return nil, model.Paging{}, errors.Wrap(err, "Failed to parse ledger when listing by campaign")
		}

		result = append(result, ledger)
	}
	if err := rows.Err(); err != nil {

		return nil, model.Paging{}, errors.Wrap(err, "Failed to read row when listing ledgers data")
	}

	return result, model.Paging{Page: filters.Page, Size: len(result), TotalItems: countingResult.Total}, nil
}
