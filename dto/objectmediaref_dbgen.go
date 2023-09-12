// !! Autogenerated file, do not edit!
package dto

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
)

// ObjectMediaRefField TBD
type ObjectMediaRefField int

const (
	// ObjectMediaRefFieldObjectID TBD
	ObjectMediaRefFieldObjectID ObjectMediaRefField = iota

	// ObjectMediaRefFieldMediaID TBD
	ObjectMediaRefFieldMediaID

	// ObjectMediaRefFieldMediaPosition TBD
	ObjectMediaRefFieldMediaPosition
)

// ObjectMediaRefFieldsList TBD
type ObjectMediaRefFieldsList []ObjectMediaRefField

// ObjectMediaRefAllFields TBD
var ObjectMediaRefAllFields = ObjectMediaRefFieldsList{
	ObjectMediaRefFieldObjectID,
	ObjectMediaRefFieldMediaID,
	ObjectMediaRefFieldMediaPosition,
}

// ObjectMediaRefFieldsNames TBD
var ObjectMediaRefFieldsNames = [...]string{
	"object_id",
	"media_id",
	"media_position",
}

// Name TBD
func (f ObjectMediaRefField) Name() string {
	if int(f) > len(ObjectMediaRefFieldsNames)-1 {
		return "unknown"
	}
	return ObjectMediaRefFieldsNames[f]
}

// JoinedNames TBD
func (l ObjectMediaRefFieldsList) JoinedNames() string {
	return l.JoinedNamesWithAlias("")
}

// JoinedNamesWithAlias TBD
func (l ObjectMediaRefFieldsList) JoinedNamesWithAlias(alias string) string {
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
func (l ObjectMediaRefFieldsList) Placeholders() string {
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
func (l ObjectMediaRefFieldsList) Del(fields ...ObjectMediaRefField) ObjectMediaRefFieldsList {
	var res = make(ObjectMediaRefFieldsList, 0, len(l))
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
func (l ObjectMediaRefFieldsList) PushBack(fields ...ObjectMediaRefField) ObjectMediaRefFieldsList {
	var res = make(ObjectMediaRefFieldsList, 0, len(l)+len(fields))
	for _, f := range l {
		res = append(res, f)
	}
	for _, f := range fields {
		res = append(res, f)
	}
	return res
}

// PushFront TBD
func (l ObjectMediaRefFieldsList) PushFront(fields ...ObjectMediaRefField) ObjectMediaRefFieldsList {
	var res = make(ObjectMediaRefFieldsList, 0, len(l)+len(fields))
	for _, f := range fields {
		res = append(res, f)
	}
	for _, f := range l {
		res = append(res, f)
	}
	return res
}

// FieldsValues TBD
func (x *ObjectMediaRef) FieldsValues(fields ObjectMediaRefFieldsList) []interface{} {
	values := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		values = append(values, x.FieldValue(f))
	}
	return values
}

// FieldValue TBD
func (x *ObjectMediaRef) FieldValue(f ObjectMediaRefField) interface{} {
	switch f {
	case ObjectMediaRefFieldObjectID:
		return x.ObjectID
	case ObjectMediaRefFieldMediaID:
		return x.MediaID
	case ObjectMediaRefFieldMediaPosition:
		return x.MediaPosition
	}
	return nil
}

// SetFieldValue TBD
func (x *ObjectMediaRef) SetFieldValue(f ObjectMediaRefField, v interface{}) error {
	switch f {
	case ObjectMediaRefFieldObjectID:
		x.ObjectID = v.(int64)
	case ObjectMediaRefFieldMediaID:
		x.MediaID = v.(int64)
	case ObjectMediaRefFieldMediaPosition:
		x.MediaPosition = v.(int32)
	}
	return nil
}

// ObjectMediaRefQuerier TBD
type ObjectMediaRefQuerier interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}

// ObjectMediaRefRows TBD
type ObjectMediaRefRows struct {
	rows   *pgx.Rows
	fields ObjectMediaRefFieldsList
}

// Close TBD
func (r *ObjectMediaRefRows) Close() {
	if r.rows == nil {
		return
	}
	r.rows.Close()
}

// Next TBD
func (r *ObjectMediaRefRows) Next() bool {
	if r.rows == nil {
		return false
	}

	return r.rows.Next()
}

// Err TBD
func (r *ObjectMediaRefRows) Err() error {
	if r.rows == nil {
		return errors.New("empty rows")
	}
	return r.rows.Err()
}

// ScanTo TBD
func (r *ObjectMediaRefRows) ScanTo(x *ObjectMediaRef) error {
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
func (r *ObjectMediaRefRows) Scan() (x *ObjectMediaRef, err error) {
	x = new(ObjectMediaRef)
	err = r.ScanTo(x)
	return
}

// ScanAll TBD
func (r *ObjectMediaRefRows) ScanAll(sizeHint int) ([]*ObjectMediaRef, error) {
	defer r.Close()

	if sizeHint == 0 {
		sizeHint = 10
	}

	var list = make([]*ObjectMediaRef, 0, sizeHint)

	for r.Next() {
		item, err := r.Scan()
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}

	return list, r.Err()
}

// ObjectMediaRefRow TBD
type ObjectMediaRefRow ObjectMediaRefRows

// ScanTo TBD
func (r *ObjectMediaRefRow) ScanTo(x *ObjectMediaRef) error {
	rows := (*ObjectMediaRefRows)(r)
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
func (r *ObjectMediaRefRow) Scan() (x *ObjectMediaRef, err error) {
	x = new(ObjectMediaRef)
	err = r.ScanTo(x)
	return
}

// QueryObjectMediaRef TBD
func QueryObjectMediaRef(q ObjectMediaRefQuerier, fields ObjectMediaRefFieldsList, sql string, args ...interface{}) (*ObjectMediaRefRows, error) {
	pgxRows, err := q.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	return &ObjectMediaRefRows{rows: pgxRows, fields: fields}, nil
}

// QueryObjectMediaRefRow TBD
func QueryObjectMediaRefRow(q ObjectMediaRefQuerier, fields ObjectMediaRefFieldsList, sql string, args ...interface{}) *ObjectMediaRefRow {
	pgxRows, _ := q.Query(sql, args...)
	return &ObjectMediaRefRow{rows: pgxRows, fields: fields}
}

// ScanObjectMediaRefList TBD
func ScanObjectMediaRefList(q ObjectMediaRefQuerier, fields ObjectMediaRefFieldsList, sql string, args ...interface{}) ([]*ObjectMediaRef, error) {
	rows, err := QueryObjectMediaRef(q, fields, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows.ScanAll(0)
}

// ScanObjectMediaRef TBD
func ScanObjectMediaRef(q ObjectMediaRefQuerier, fields ObjectMediaRefFieldsList, sql string, args ...interface{}) (*ObjectMediaRef, error) {
	x, err := QueryObjectMediaRefRow(q, fields, sql, args...).Scan()
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return x, err
}
