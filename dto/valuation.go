package dto

import (
	"time"

	"github.com/jackc/pgx/pgtype"
)

//go:generate dbgen -type Valuation

// Valuation TBD
type (
	Valuation struct {
		ID         int64             `db:"id"`
		ObjectID   int64             `db:"object_id"`
		Name       string            `db:"name"`
		Comment    string            `db:"comment"`
		Date       *time.Time        `db:"date"`
		Price      *pgtype.Int8range `db:"price"`
		CurrencyID *int64            `db:"price_currency_id"`
	}

	// ValuationList TBD
	ValuationList []*Valuation
)

// ObjectToOneValuationMap TBD
func (vl ValuationList) ObjectToOneValuationMap() map[int64][]*Valuation {
	objectIDValuationRef := make(map[int64][]*Valuation)
	for _, ref := range vl {
		objectIDValuationRef[ref.ObjectID] = append(objectIDValuationRef[ref.ObjectID], ref)
	}
	return objectIDValuationRef
}
