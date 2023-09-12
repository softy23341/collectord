package dto

//go:generate dbgen -type Currency

// Currency TBD
type Currency struct {
	ID       int64  `db:"id"`
	Symbol   string `db:"symbol"`
	Code     string `db:"code"`
	Num      string `db:"num"`
	E        int16  `db:"e"`
	Currency string `db:"currency"`
}

// CurrencyList TBD
type CurrencyList []*Currency
