// !! Autogenerated file, do not edit!
package dto

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
)

// CurrencyField TBD
type CurrencyField int

const (
	// CurrencyFieldID TBD
	CurrencyFieldID CurrencyField = iota

	// CurrencyFieldSymbol TBD
	CurrencyFieldSymbol

	// CurrencyFieldCode TBD
	CurrencyFieldCode

	// CurrencyFieldNum TBD
	CurrencyFieldNum

	// CurrencyFieldE TBD
	CurrencyFieldE

	// CurrencyFieldCurrency TBD
	CurrencyFieldCurrency
)

// CurrencyFieldsList TBD
type CurrencyFieldsList []CurrencyField

// CurrencyAllFields TBD
var CurrencyAllFields = CurrencyFieldsList{
	CurrencyFieldID,
	CurrencyFieldSymbol,
	CurrencyFieldCode,
	CurrencyFieldNum,
	CurrencyFieldE,
	CurrencyFieldCurrency,
}

// CurrencyFieldsNames TBD
var CurrencyFieldsNames = [...]string{
	"id",
	"symbol",
	"code",
	"num",
	"e",
	"currency",
}

// Name TBD
func (f CurrencyField) Name() string {
	if int(f) > len(CurrencyFieldsNames)-1 {
		return "unknown"
	}
	return CurrencyFieldsNames[f]
}

// JoinedNames TBD
func (l CurrencyFieldsList) JoinedNames() string {
	return l.JoinedNamesWithAlias("")
}

// JoinedNamesWithAlias TBD
func (l CurrencyFieldsList) JoinedNamesWithAlias(alias string) string {
	var aliasPrefix string
	if alias != "" {
		aliasPrefix = alias + "."
	}

	var buf bytes.Buffer
	for idx, f := range l {
		if idx != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(aliasPrefix)
		buf.WriteString("\"")
		buf.WriteString(f.Name())
		buf.WriteString("\"")
	}
	return buf.String()
}

// Placeholders TBD
func (l CurrencyFieldsList) Placeholders() string {
	var buf bytes.Buffer
	for idx := range l {
		if idx != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprint("$", idx+1))
	}
	return buf.String()
}

// Del TBD
func (l CurrencyFieldsList) Del(fields ...CurrencyField) CurrencyFieldsList {
	var res = make(CurrencyFieldsList, 0, len(l))
	for _, srcFld := range l {
		remove := false
		for _, delFld := range fields {
			if srcFld == delFld {
				remove = true
				break
			}
		}
		if !remove {
			res = append(res, srcFld)
		}
	}
	return res
}

// PushBack TBD
func (l CurrencyFieldsList) PushBack(fields ...CurrencyField) CurrencyFieldsList {
	var res = make(CurrencyFieldsList, 0, len(l)+len(fields))
	for _, f := range l {
		res = append(res, f)
	}
	for _, f := range fields {
		res = append(res, f)
	}
	return res
}

// PushFront TBD
func (l CurrencyFieldsList) PushFront(fields ...CurrencyField) CurrencyFieldsList {
	var res = make(CurrencyFieldsList, 0, len(l)+len(fields))
	for _, f := range fields {
		res = append(res, f)
	}
	for _, f := range l {
		res = append(res, f)
	}
	return res
}

// FieldsValues TBD
func (x *Currency) FieldsValues(fields CurrencyFieldsList) []interface{} {
	values := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		values = append(values, x.FieldValue(f))
	}
	return values
}

// FieldValue TBD
func (x *Currency) FieldValue(f CurrencyField) interface{} {
	switch f {
	case CurrencyFieldID:
		return x.ID
	case CurrencyFieldSymbol:
		return x.Symbol
	case CurrencyFieldCode:
		return x.Code
	case CurrencyFieldNum:
		return x.Num
	case CurrencyFieldE:
		return x.E
	case CurrencyFieldCurrency:
		return x.Currency
	}
	return nil
}

// SetFieldValue TBD
func (x *Currency) SetFieldValue(f CurrencyField, v interface{}) error {
	switch f {
	case CurrencyFieldID:
		x.ID = v.(int64)
	case CurrencyFieldSymbol:
		x.Symbol = v.(string)
	case CurrencyFieldCode:
		x.Code = v.(string)
	case CurrencyFieldNum:
		x.Num = v.(string)
	case CurrencyFieldE:
		x.E = v.(int16)
	case CurrencyFieldCurrency:
		x.Currency = v.(string)
	}
	return nil
}

// CurrencyQuerier TBD
type CurrencyQuerier interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}

// CurrencyRows TBD
type CurrencyRows struct {
	rows   *pgx.Rows
	fields CurrencyFieldsList
}

// Close TBD
func (r *CurrencyRows) Close() {
	if r.rows == nil {
		return
	}
	r.rows.Close()
}

// Next TBD
func (r *CurrencyRows) Next() bool {
	if r.rows == nil {
		return false
	}

	return r.rows.Next()
}

// Err TBD
func (r *CurrencyRows) Err() error {
	if r.rows == nil {
		return errors.New("empty rows")
	}
	return r.rows.Err()
}

// ScanTo TBD
func (r *CurrencyRows) ScanTo(x *Currency) error {
	values, err := r.rows.Values()
	if err != nil {
		return err
	}
	for idx, v := range values {
		if err := x.SetFieldValue(r.fields[idx], v); err != nil {
			return err
		}
	}

	return nil
}

// Scan TBD
func (r *CurrencyRows) Scan() (x *Currency, err error) {
	x = new(Currency)
	err = r.ScanTo(x)
	return
}

// ScanAll TBD
func (r *CurrencyRows) ScanAll(sizeHint int) ([]*Currency, error) {
	defer r.Close()

	if sizeHint == 0 {
		sizeHint = 10
	}

	var list = make([]*Currency, 0, sizeHint)

	for r.Next() {
		item, err := r.Scan()
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}

	return list, r.Err()
}

// CurrencyRow TBD
type CurrencyRow CurrencyRows

// ScanTo TBD
func (r *CurrencyRow) ScanTo(x *Currency) error {
	rows := (*CurrencyRows)(r)
	defer rows.Close()

	if rows.Err() != nil {
		return rows.Err()
	}

	if !rows.Next() {
		if rows.Err() == nil {
			return pgx.ErrNoRows
		}
		return rows.Err()
	}

	return rows.ScanTo(x)
}

// Scan TBD
func (r *CurrencyRow) Scan() (x *Currency, err error) {
	x = new(Currency)
	err = r.ScanTo(x)
	return
}

// QueryCurrency TBD
func QueryCurrency(q CurrencyQuerier, fields CurrencyFieldsList, sql string, args ...interface{}) (*CurrencyRows, error) {
	pgxRows, err := q.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	return &CurrencyRows{rows: pgxRows, fields: fields}, nil
}

// QueryCurrencyRow TBD
func QueryCurrencyRow(q CurrencyQuerier, fields CurrencyFieldsList, sql string, args ...interface{}) *CurrencyRow {
	pgxRows, _ := q.Query(sql, args...)
	return &CurrencyRow{rows: pgxRows, fields: fields}
}

// ScanCurrencyList TBD
func ScanCurrencyList(q CurrencyQuerier, fields CurrencyFieldsList, sql string, args ...interface{}) ([]*Currency, error) {
	rows, err := QueryCurrency(q, fields, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows.ScanAll(0)
}

// ScanCurrency TBD
func ScanCurrency(q CurrencyQuerier, fields CurrencyFieldsList, sql string, args ...interface{}) (*Currency, error) {
	x, err := QueryCurrencyRow(q, fields, sql, args...).Scan()
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return x, err
}
