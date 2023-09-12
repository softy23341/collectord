package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetCurrenciesByIDs TBD
func (m *Manager) GetCurrenciesByIDs(currencyIDs []int64) (dto.CurrencyList, error) {
	flds := dto.CurrencyAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM currency
                WHERE id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanCurrencyList(m.p, flds, sql, currencyIDs)
}

// GetCurrencies TBD
func (m *Manager) GetCurrencies() (dto.CurrencyList, error) {
	flds := dto.CurrencyAllFields

	sql := fmt.Sprintf(`SELECT %s FROM currency ORDER BY id ASC`, flds.JoinedNames())

	return dto.ScanCurrencyList(m.p, flds, sql)
}
