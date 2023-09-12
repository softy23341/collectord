package dto

//go:generate dbgen -type NamedDateInterval

// NamedDateInterval TBD
type NamedDateInterval struct {
	ID                         int64  `db:"id"`
	RootID                     *int64 `db:"root_id"`
	ProductionDateIntervalFrom int64  `db:"production_date_interval_from"`
	ProductionDateIntervalTo   int64  `db:"production_date_interval_to"`
	Name                       string `db:"name"`
	NormalName                 string `db:"normal_name"`
}

// NamedDateIntervalList TBD
type NamedDateIntervalList []*NamedDateInterval

// IDToNamedDateInterval TBD
func (ni NamedDateIntervalList) IDToNamedDateInterval() map[int64]*NamedDateInterval {
	id2date := make(map[int64]*NamedDateInterval, 0)
	for _, interval := range ni {
		id2date[interval.ID] = interval
	}
	return id2date
}
