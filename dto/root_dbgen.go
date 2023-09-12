// !! Autogenerated file, do not edit!
package dto

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
)

// RootField TBD
type RootField int

const (
	// RootFieldID TBD
	RootFieldID RootField = iota
)

// RootFieldsList TBD
type RootFieldsList []RootField

// RootAllFields TBD
var RootAllFields = RootFieldsList{
	RootFieldID,
}

// RootFieldsNames TBD
var RootFieldsNames = [...]string{
	"id",
}

// Name TBD
func (f RootField) Name() string {
	if int(f) > len(RootFieldsNames)-1 {
		return "unknown"
	}
	return RootFieldsNames[f]
}

// JoinedNames TBD
func (l RootFieldsList) JoinedNames() string {
	return l.JoinedNamesWithAlias("")
}

// JoinedNamesWithAlias TBD
func (l RootFieldsList) JoinedNamesWithAlias(alias string) string {
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
func (l RootFieldsList) Placeholders() string {
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
func (l RootFieldsList) Del(fields ...RootField) RootFieldsList {
	var res = make(RootFieldsList, 0, len(l))
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
func (l RootFieldsList) PushBack(fields ...RootField) RootFieldsList {
	var res = make(RootFieldsList, 0, len(l)+len(fields))
	for _, f := range l {
		res = append(res, f)
	}
	for _, f := range fields {
		res = append(res, f)
	}
	return res
}

// PushFront TBD
func (l RootFieldsList) PushFront(fields ...RootField) RootFieldsList {
	var res = make(RootFieldsList, 0, len(l)+len(fields))
	for _, f := range fields {
		res = append(res, f)
	}
	for _, f := range l {
		res = append(res, f)
	}
	return res
}

// FieldsValues TBD
func (x *Root) FieldsValues(fields RootFieldsList) []interface{} {
	values := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		values = append(values, x.FieldValue(f))
	}
	return values
}

// FieldValue TBD
func (x *Root) FieldValue(f RootField) interface{} {
	switch f {
	case RootFieldID:
		return x.ID
	}
	return nil
}

// SetFieldValue TBD
func (x *Root) SetFieldValue(f RootField, v interface{}) error {
	switch f {
	case RootFieldID:
		x.ID = v.(int64)
	}
	return nil
}

// RootQuerier TBD
type RootQuerier interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}

// RootRows TBD
type RootRows struct {
	rows   *pgx.Rows
	fields RootFieldsList
}

// Close TBD
func (r *RootRows) Close() {
	if r.rows == nil {
		return
	}
	r.rows.Close()
}

// Next TBD
func (r *RootRows) Next() bool {
	if r.rows == nil {
		return false
	}

	return r.rows.Next()
}

// Err TBD
func (r *RootRows) Err() error {
	if r.rows == nil {
		return errors.New("empty rows")
	}
	return r.rows.Err()
}

// ScanTo TBD
func (r *RootRows) ScanTo(x *Root) error {
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
func (r *RootRows) Scan() (x *Root, err error) {
	x = new(Root)
	err = r.ScanTo(x)
	return
}

// ScanAll TBD
func (r *RootRows) ScanAll(sizeHint int) ([]*Root, error) {
	defer r.Close()

	if sizeHint == 0 {
		sizeHint = 10
	}

	var list = make([]*Root, 0, sizeHint)

	for r.Next() {
		item, err := r.Scan()
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}

	return list, r.Err()
}

// RootRow TBD
type RootRow RootRows

// ScanTo TBD
func (r *RootRow) ScanTo(x *Root) error {
	rows := (*RootRows)(r)
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
func (r *RootRow) Scan() (x *Root, err error) {
	x = new(Root)
	err = r.ScanTo(x)
	return
}

// QueryRoot TBD
func QueryRoot(q RootQuerier, fields RootFieldsList, sql string, args ...interface{}) (*RootRows, error) {
	pgxRows, err := q.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	return &RootRows{rows: pgxRows, fields: fields}, nil
}

// QueryRootRow TBD
func QueryRootRow(q RootQuerier, fields RootFieldsList, sql string, args ...interface{}) *RootRow {
	pgxRows, _ := q.Query(sql, args...)
	return &RootRow{rows: pgxRows, fields: fields}
}

// ScanRootList TBD
func ScanRootList(q RootQuerier, fields RootFieldsList, sql string, args ...interface{}) ([]*Root, error) {
	rows, err := QueryRoot(q, fields, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows.ScanAll(0)
}

// ScanRoot TBD
func ScanRoot(q RootQuerier, fields RootFieldsList, sql string, args ...interface{}) (*Root, error) {
	x, err := QueryRootRow(q, fields, sql, args...).Scan()
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return x, err
}