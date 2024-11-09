package model

type ListingFilters struct {
	Page      int    `query:"page"`
	Size      int    `query:"size"`
	Status    string `query:"status"`
	StartDate string `query:"start"`
	EndDate   string `query:"end"`
}
