// !! Autogenerated file, do not edit!
package dto

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
)

// ObjectBadgeRefField TBD
type ObjectBadgeRefField int

const (
	// ObjectBadgeRefFieldObjectID TBD
	ObjectBadgeRefFieldObjectID ObjectBadgeRefField = iota

	// ObjectBadgeRefFieldBadgeID TBD
	ObjectBadgeRefFieldBadgeID
)

// ObjectBadgeRefFieldsList TBD
type ObjectBadgeRefFieldsList []ObjectBadgeRefField

// ObjectBadgeRefAllFields TBD
var ObjectBadgeRefAllFields = ObjectBadgeRefFieldsList{
	ObjectBadgeRefFieldObjectID,
	ObjectBadgeRefFieldBadgeID,
}

// ObjectBadgeRefFieldsNames TBD
var ObjectBadgeRefFieldsNames = [...]string{
	"object_id",
	"badge_id",
}

// Name TBD
func (f ObjectBadgeRefField) Name() string {
	if int(f) > len(ObjectBadgeRefFieldsNames)-1 {
		return "unknown"
	}
	return ObjectBadgeRefFieldsNames[f]
}

// JoinedNames TBD
func (l ObjectBadgeRefFieldsList) JoinedNames() string {
	return l.JoinedNamesWithAlias("")
}

// JoinedNamesWithAlias TBD
func (l ObjectBadgeRefFieldsList) JoinedNamesWithAlias(alias string) string {
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
func (l ObjectBadgeRefFieldsList) Placeholders() string {
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
func (l ObjectBadgeRefFieldsList) Del(fields ...ObjectBadgeRefField) ObjectBadgeRefFieldsList {
	var res = make(ObjectBadgeRefFieldsList, 0, len(l))
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
func (l ObjectBadgeRefFieldsList) PushBack(fields ...ObjectBadgeRefField) ObjectBadgeRefFieldsList {
	var res = make(ObjectBadgeRefFieldsList, 0, len(l)+len(fields))
	for _, f := range l {
		res = append(res, f)
	}
	for _, f := range fields {
		res = append(res, f)
	}
	return res
}

// PushFront TBD
func (l ObjectBadgeRefFieldsList) PushFront(fields ...ObjectBadgeRefField) ObjectBadgeRefFieldsList {
	var res = make(ObjectBadgeRefFieldsList, 0, len(l)+len(fields))
	for _, f := range fields {
		res = append(res, f)
	}
	for _, f := range l {
		res = append(res, f)
	}
	return res
}

// FieldsValues TBD
func (x *ObjectBadgeRef) FieldsValues(fields ObjectBadgeRefFieldsList) []interface{} {
	values := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		values = append(values, x.FieldValue(f))
	}
	return values
}

// FieldValue TBD
func (x *ObjectBadgeRef) FieldValue(f ObjectBadgeRefField) interface{} {
	switch f {
	case ObjectBadgeRefFieldObjectID:
		return x.ObjectID
	case ObjectBadgeRefFieldBadgeID:
		return x.BadgeID
	}
	return nil
}

// SetFieldValue TBD
func (x *ObjectBadgeRef) SetFieldValue(f ObjectBadgeRefField, v interface{}) error {
	switch f {
	case ObjectBadgeRefFieldObjectID:
		x.ObjectID = v.(int64)
	case ObjectBadgeRefFieldBadgeID:
		x.BadgeID = v.(int64)
	}
	return nil
}

// ObjectBadgeRefQuerier TBD
type ObjectBadgeRefQuerier interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}

// ObjectBadgeRefRows TBD
type ObjectBadgeRefRows struct {
	rows   *pgx.Rows
	fields ObjectBadgeRefFieldsList
}

// Close TBD
func (r *ObjectBadgeRefRows) Close() {
	if r.rows == nil {
		return
	}
	r.rows.Close()
}

// Next TBD
func (r *ObjectBadgeRefRows) Next() bool {
	if r.rows == nil {
		return false
	}

	return r.rows.Next()
}

// Err TBD
func (r *ObjectBadgeRefRows) Err() error {
	if r.rows == nil {
		return errors.New("empty rows")
	}
	return r.rows.Err()
}

// ScanTo TBD
func (r *ObjectBadgeRefRows) ScanTo(x *ObjectBadgeRef) error {
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
func (r *ObjectBadgeRefRows) Scan() (x *ObjectBadgeRef, err error) {
	x = new(ObjectBadgeRef)
	err = r.ScanTo(x)
	return
}

// ScanAll TBD
func (r *ObjectBadgeRefRows) ScanAll(sizeHint int) ([]*ObjectBadgeRef, error) {
	defer r.Close()

	if sizeHint == 0 {
		sizeHint = 10
	}

	var list = make([]*ObjectBadgeRef, 0, sizeHint)

	for r.Next() {
		item, err := r.Scan()
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}

	return list, r.Err()
}

// ObjectBadgeRefRow TBD
type ObjectBadgeRefRow ObjectBadgeRefRows

// ScanTo TBD
func (r *ObjectBadgeRefRow) ScanTo(x *ObjectBadgeRef) error {
	rows := (*ObjectBadgeRefRows)(r)
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
func (r *ObjectBadgeRefRow) Scan() (x *ObjectBadgeRef, err error) {
	x = new(ObjectBadgeRef)
	err = r.ScanTo(x)
	return
}

// QueryObjectBadgeRef TBD
func QueryObjectBadgeRef(q ObjectBadgeRefQuerier, fields ObjectBadgeRefFieldsList, sql string, args ...interface{}) (*ObjectBadgeRefRows, error) {
	pgxRows, err := q.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	return &ObjectBadgeRefRows{rows: pgxRows, fields: fields}, nil
}

// QueryObjectBadgeRefRow TBD
func QueryObjectBadgeRefRow(q ObjectBadgeRefQuerier, fields ObjectBadgeRefFieldsList, sql string, args ...interface{}) *ObjectBadgeRefRow {
	pgxRows, _ := q.Query(sql, args...)
	return &ObjectBadgeRefRow{rows: pgxRows, fields: fields}
}

// ScanObjectBadgeRefList TBD
func ScanObjectBadgeRefList(q ObjectBadgeRefQuerier, fields ObjectBadgeRefFieldsList, sql string, args ...interface{}) ([]*ObjectBadgeRef, error) {
	rows, err := QueryObjectBadgeRef(q, fields, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows.ScanAll(0)
}

// ScanObjectBadgeRef TBD
func ScanObjectBadgeRef(q ObjectBadgeRefQuerier, fields ObjectBadgeRefFieldsList, sql string, args ...interface{}) (*ObjectBadgeRef, error) {
	x, err := QueryObjectBadgeRefRow(q, fields, sql, args...).Scan()
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return x, err
}
