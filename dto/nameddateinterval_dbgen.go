// !! Autogenerated file, do not edit!
package dto

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
)

// NamedDateIntervalField TBD
type NamedDateIntervalField int

const (
	// NamedDateIntervalFieldID TBD
	NamedDateIntervalFieldID NamedDateIntervalField = iota

	// NamedDateIntervalFieldRootID TBD
	NamedDateIntervalFieldRootID

	// NamedDateIntervalFieldProductionDateIntervalFrom TBD
	NamedDateIntervalFieldProductionDateIntervalFrom

	// NamedDateIntervalFieldProductionDateIntervalTo TBD
	NamedDateIntervalFieldProductionDateIntervalTo

	// NamedDateIntervalFieldName TBD
	NamedDateIntervalFieldName

	// NamedDateIntervalFieldNormalName TBD
	NamedDateIntervalFieldNormalName
)

// NamedDateIntervalFieldsList TBD
type NamedDateIntervalFieldsList []NamedDateIntervalField

// NamedDateIntervalAllFields TBD
var NamedDateIntervalAllFields = NamedDateIntervalFieldsList{
	NamedDateIntervalFieldID,
	NamedDateIntervalFieldRootID,
	NamedDateIntervalFieldProductionDateIntervalFrom,
	NamedDateIntervalFieldProductionDateIntervalTo,
	NamedDateIntervalFieldName,
	NamedDateIntervalFieldNormalName,
}

// NamedDateIntervalFieldsNames TBD
var NamedDateIntervalFieldsNames = [...]string{
	"id",
	"root_id",
	"production_date_interval_from",
	"production_date_interval_to",
	"name",
	"normal_name",
}

// Name TBD
func (f NamedDateIntervalField) Name() string {
	if int(f) > len(NamedDateIntervalFieldsNames)-1 {
		return "unknown"
	}
	return NamedDateIntervalFieldsNames[f]
}

// JoinedNames TBD
func (l NamedDateIntervalFieldsList) JoinedNames() string {
	return l.JoinedNamesWithAlias("")
}

// JoinedNamesWithAlias TBD
func (l NamedDateIntervalFieldsList) JoinedNamesWithAlias(alias string) string {
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
func (l NamedDateIntervalFieldsList) Placeholders() string {
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
func (l NamedDateIntervalFieldsList) Del(fields ...NamedDateIntervalField) NamedDateIntervalFieldsList {
	var res = make(NamedDateIntervalFieldsList, 0, len(l))
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
func (l NamedDateIntervalFieldsList) PushBack(fields ...NamedDateIntervalField) NamedDateIntervalFieldsList {
	var res = make(NamedDateIntervalFieldsList, 0, len(l)+len(fields))
	for _, f := range l {
		res = append(res, f)
	}
	for _, f := range fields {
		res = append(res, f)
	}
	return res
}

// PushFront TBD
func (l NamedDateIntervalFieldsList) PushFront(fields ...NamedDateIntervalField) NamedDateIntervalFieldsList {
	var res = make(NamedDateIntervalFieldsList, 0, len(l)+len(fields))
	for _, f := range fields {
		res = append(res, f)
	}
	for _, f := range l {
		res = append(res, f)
	}
	return res
}

// FieldsValues TBD
func (x *NamedDateInterval) FieldsValues(fields NamedDateIntervalFieldsList) []interface{} {
	values := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		values = append(values, x.FieldValue(f))
	}
	return values
}

// FieldValue TBD
func (x *NamedDateInterval) FieldValue(f NamedDateIntervalField) interface{} {
	switch f {
	case NamedDateIntervalFieldID:
		return x.ID
	case NamedDateIntervalFieldRootID:
		return &x.RootID
	case NamedDateIntervalFieldProductionDateIntervalFrom:
		return x.ProductionDateIntervalFrom
	case NamedDateIntervalFieldProductionDateIntervalTo:
		return x.ProductionDateIntervalTo
	case NamedDateIntervalFieldName:
		return x.Name
	case NamedDateIntervalFieldNormalName:
		return x.NormalName
	}
	return nil
}

// SetFieldValue TBD
func (x *NamedDateInterval) SetFieldValue(f NamedDateIntervalField, v interface{}) error {
	switch f {
	case NamedDateIntervalFieldID:
		x.ID = v.(int64)
	case NamedDateIntervalFieldRootID:
		if v == nil {
			x.RootID = nil
		} else {
			value := v.(int64)
			x.RootID = &value
		}
	case NamedDateIntervalFieldProductionDateIntervalFrom:
		x.ProductionDateIntervalFrom = v.(int64)
	case NamedDateIntervalFieldProductionDateIntervalTo:
		x.ProductionDateIntervalTo = v.(int64)
	case NamedDateIntervalFieldName:
		x.Name = v.(string)
	case NamedDateIntervalFieldNormalName:
		x.NormalName = v.(string)
	}
	return nil
}

// NamedDateIntervalQuerier TBD
type NamedDateIntervalQuerier interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}

// NamedDateIntervalRows TBD
type NamedDateIntervalRows struct {
	rows   *pgx.Rows
	fields NamedDateIntervalFieldsList
}

// Close TBD
func (r *NamedDateIntervalRows) Close() {
	if r.rows == nil {
		return
	}
	r.rows.Close()
}

// Next TBD
func (r *NamedDateIntervalRows) Next() bool {
	if r.rows == nil {
		return false
	}

	return r.rows.Next()
}

// Err TBD
func (r *NamedDateIntervalRows) Err() error {
	if r.rows == nil {
		return errors.New("empty rows")
	}
	return r.rows.Err()
}

// ScanTo TBD
func (r *NamedDateIntervalRows) ScanTo(x *NamedDateInterval) error {
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
func (r *NamedDateIntervalRows) Scan() (x *NamedDateInterval, err error) {
	x = new(NamedDateInterval)
	err = r.ScanTo(x)
	return
}

// ScanAll TBD
func (r *NamedDateIntervalRows) ScanAll(sizeHint int) ([]*NamedDateInterval, error) {
	defer r.Close()

	if sizeHint == 0 {
		sizeHint = 10
	}

	var list = make([]*NamedDateInterval, 0, sizeHint)

	for r.Next() {
		item, err := r.Scan()
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}

	return list, r.Err()
}

// NamedDateIntervalRow TBD
type NamedDateIntervalRow NamedDateIntervalRows

// ScanTo TBD
func (r *NamedDateIntervalRow) ScanTo(x *NamedDateInterval) error {
	rows := (*NamedDateIntervalRows)(r)
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
func (r *NamedDateIntervalRow) Scan() (x *NamedDateInterval, err error) {
	x = new(NamedDateInterval)
	err = r.ScanTo(x)
	return
}

// QueryNamedDateInterval TBD
func QueryNamedDateInterval(q NamedDateIntervalQuerier, fields NamedDateIntervalFieldsList, sql string, args ...interface{}) (*NamedDateIntervalRows, error) {
	pgxRows, err := q.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	return &NamedDateIntervalRows{rows: pgxRows, fields: fields}, nil
}

// QueryNamedDateIntervalRow TBD
func QueryNamedDateIntervalRow(q NamedDateIntervalQuerier, fields NamedDateIntervalFieldsList, sql string, args ...interface{}) *NamedDateIntervalRow {
	pgxRows, _ := q.Query(sql, args...)
	return &NamedDateIntervalRow{rows: pgxRows, fields: fields}
}

// ScanNamedDateIntervalList TBD
func ScanNamedDateIntervalList(q NamedDateIntervalQuerier, fields NamedDateIntervalFieldsList, sql string, args ...interface{}) ([]*NamedDateInterval, error) {
	rows, err := QueryNamedDateInterval(q, fields, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows.ScanAll(0)
}

// ScanNamedDateInterval TBD
func ScanNamedDateInterval(q NamedDateIntervalQuerier, fields NamedDateIntervalFieldsList, sql string, args ...interface{}) (*NamedDateInterval, error) {
	x, err := QueryNamedDateIntervalRow(q, fields, sql, args...).Scan()
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return x, err
}
