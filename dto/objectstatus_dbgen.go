// !! Autogenerated file, do not edit!
package dto

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
)

// ObjectStatusField TBD
type ObjectStatusField int

const (
	// ObjectStatusFieldID TBD
	ObjectStatusFieldID ObjectStatusField = iota

	// ObjectStatusFieldName TBD
	ObjectStatusFieldName

	// ObjectStatusFieldDescription TBD
	ObjectStatusFieldDescription

	// ObjectStatusFieldImageMediaID TBD
	ObjectStatusFieldImageMediaID
)

// ObjectStatusFieldsList TBD
type ObjectStatusFieldsList []ObjectStatusField

// ObjectStatusAllFields TBD
var ObjectStatusAllFields = ObjectStatusFieldsList{
	ObjectStatusFieldID,
	ObjectStatusFieldName,
	ObjectStatusFieldDescription,
	ObjectStatusFieldImageMediaID,
}

// ObjectStatusFieldsNames TBD
var ObjectStatusFieldsNames = [...]string{
	"id",
	"name",
	"description",
	"image_media_id",
}

// Name TBD
func (f ObjectStatusField) Name() string {
	if int(f) > len(ObjectStatusFieldsNames)-1 {
		return "unknown"
	}
	return ObjectStatusFieldsNames[f]
}

// JoinedNames TBD
func (l ObjectStatusFieldsList) JoinedNames() string {
	return l.JoinedNamesWithAlias("")
}

// JoinedNamesWithAlias TBD
func (l ObjectStatusFieldsList) JoinedNamesWithAlias(alias string) string {
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
func (l ObjectStatusFieldsList) Placeholders() string {
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
func (l ObjectStatusFieldsList) Del(fields ...ObjectStatusField) ObjectStatusFieldsList {
	var res = make(ObjectStatusFieldsList, 0, len(l))
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
func (l ObjectStatusFieldsList) PushBack(fields ...ObjectStatusField) ObjectStatusFieldsList {
	var res = make(ObjectStatusFieldsList, 0, len(l)+len(fields))
	for _, f := range l {
		res = append(res, f)
	}
	for _, f := range fields {
		res = append(res, f)
	}
	return res
}

// PushFront TBD
func (l ObjectStatusFieldsList) PushFront(fields ...ObjectStatusField) ObjectStatusFieldsList {
	var res = make(ObjectStatusFieldsList, 0, len(l)+len(fields))
	for _, f := range fields {
		res = append(res, f)
	}
	for _, f := range l {
		res = append(res, f)
	}
	return res
}

// FieldsValues TBD
func (x *ObjectStatus) FieldsValues(fields ObjectStatusFieldsList) []interface{} {
	values := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		values = append(values, x.FieldValue(f))
	}
	return values
}

// FieldValue TBD
func (x *ObjectStatus) FieldValue(f ObjectStatusField) interface{} {
	switch f {
	case ObjectStatusFieldID:
		return x.ID
	case ObjectStatusFieldName:
		return x.Name
	case ObjectStatusFieldDescription:
		return x.Description
	case ObjectStatusFieldImageMediaID:
		return &x.ImageMediaID
	}
	return nil
}

// SetFieldValue TBD
func (x *ObjectStatus) SetFieldValue(f ObjectStatusField, v interface{}) error {
	switch f {
	case ObjectStatusFieldID:
		x.ID = v.(int64)
	case ObjectStatusFieldName:
		x.Name = v.(string)
	case ObjectStatusFieldDescription:
		x.Description = v.(string)
	case ObjectStatusFieldImageMediaID:
		if v == nil {
			x.ImageMediaID = nil
		} else {
			value := v.(int64)
			x.ImageMediaID = &value
		}
	}
	return nil
}

// ObjectStatusQuerier TBD
type ObjectStatusQuerier interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}

// ObjectStatusRows TBD
type ObjectStatusRows struct {
	rows   *pgx.Rows
	fields ObjectStatusFieldsList
}

// Close TBD
func (r *ObjectStatusRows) Close() {
	if r.rows == nil {
		return
	}
	r.rows.Close()
}

// Next TBD
func (r *ObjectStatusRows) Next() bool {
	if r.rows == nil {
		return false
	}

	return r.rows.Next()
}

// Err TBD
func (r *ObjectStatusRows) Err() error {
	if r.rows == nil {
		return errors.New("empty rows")
	}
	return r.rows.Err()
}

// ScanTo TBD
func (r *ObjectStatusRows) ScanTo(x *ObjectStatus) error {
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
func (r *ObjectStatusRows) Scan() (x *ObjectStatus, err error) {
	x = new(ObjectStatus)
	err = r.ScanTo(x)
	return
}

// ScanAll TBD
func (r *ObjectStatusRows) ScanAll(sizeHint int) ([]*ObjectStatus, error) {
	defer r.Close()

	if sizeHint == 0 {
		sizeHint = 10
	}

	var list = make([]*ObjectStatus, 0, sizeHint)

	for r.Next() {
		item, err := r.Scan()
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}

	return list, r.Err()
}

// ObjectStatusRow TBD
type ObjectStatusRow ObjectStatusRows

// ScanTo TBD
func (r *ObjectStatusRow) ScanTo(x *ObjectStatus) error {
	rows := (*ObjectStatusRows)(r)
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
func (r *ObjectStatusRow) Scan() (x *ObjectStatus, err error) {
	x = new(ObjectStatus)
	err = r.ScanTo(x)
	return
}

// QueryObjectStatus TBD
func QueryObjectStatus(q ObjectStatusQuerier, fields ObjectStatusFieldsList, sql string, args ...interface{}) (*ObjectStatusRows, error) {
	pgxRows, err := q.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	return &ObjectStatusRows{rows: pgxRows, fields: fields}, nil
}

// QueryObjectStatusRow TBD
func QueryObjectStatusRow(q ObjectStatusQuerier, fields ObjectStatusFieldsList, sql string, args ...interface{}) *ObjectStatusRow {
	pgxRows, _ := q.Query(sql, args...)
	return &ObjectStatusRow{rows: pgxRows, fields: fields}
}

// ScanObjectStatusList TBD
func ScanObjectStatusList(q ObjectStatusQuerier, fields ObjectStatusFieldsList, sql string, args ...interface{}) ([]*ObjectStatus, error) {
	rows, err := QueryObjectStatus(q, fields, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows.ScanAll(0)
}

// ScanObjectStatus TBD
func ScanObjectStatus(q ObjectStatusQuerier, fields ObjectStatusFieldsList, sql string, args ...interface{}) (*ObjectStatus, error) {
	x, err := QueryObjectStatusRow(q, fields, sql, args...).Scan()
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return x, err
}