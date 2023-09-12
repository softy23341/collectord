// !! Autogenerated file, do not edit!
package dto

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
)

// ObjectActorRefField TBD
type ObjectActorRefField int

const (
	// ObjectActorRefFieldID TBD
	ObjectActorRefFieldID ObjectActorRefField = iota

	// ObjectActorRefFieldObjectID TBD
	ObjectActorRefFieldObjectID

	// ObjectActorRefFieldActorID TBD
	ObjectActorRefFieldActorID
)

// ObjectActorRefFieldsList TBD
type ObjectActorRefFieldsList []ObjectActorRefField

// ObjectActorRefAllFields TBD
var ObjectActorRefAllFields = ObjectActorRefFieldsList{
	ObjectActorRefFieldID,
	ObjectActorRefFieldObjectID,
	ObjectActorRefFieldActorID,
}

// ObjectActorRefFieldsNames TBD
var ObjectActorRefFieldsNames = [...]string{
	"id",
	"object_id",
	"actor_id",
}

// Name TBD
func (f ObjectActorRefField) Name() string {
	if int(f) > len(ObjectActorRefFieldsNames)-1 {
		return "unknown"
	}
	return ObjectActorRefFieldsNames[f]
}

// JoinedNames TBD
func (l ObjectActorRefFieldsList) JoinedNames() string {
	return l.JoinedNamesWithAlias("")
}

// JoinedNamesWithAlias TBD
func (l ObjectActorRefFieldsList) JoinedNamesWithAlias(alias string) string {
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
func (l ObjectActorRefFieldsList) Placeholders() string {
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
func (l ObjectActorRefFieldsList) Del(fields ...ObjectActorRefField) ObjectActorRefFieldsList {
	var res = make(ObjectActorRefFieldsList, 0, len(l))
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
func (l ObjectActorRefFieldsList) PushBack(fields ...ObjectActorRefField) ObjectActorRefFieldsList {
	var res = make(ObjectActorRefFieldsList, 0, len(l)+len(fields))
	for _, f := range l {
		res = append(res, f)
	}
	for _, f := range fields {
		res = append(res, f)
	}
	return res
}

// PushFront TBD
func (l ObjectActorRefFieldsList) PushFront(fields ...ObjectActorRefField) ObjectActorRefFieldsList {
	var res = make(ObjectActorRefFieldsList, 0, len(l)+len(fields))
	for _, f := range fields {
		res = append(res, f)
	}
	for _, f := range l {
		res = append(res, f)
	}
	return res
}

// FieldsValues TBD
func (x *ObjectActorRef) FieldsValues(fields ObjectActorRefFieldsList) []interface{} {
	values := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		values = append(values, x.FieldValue(f))
	}
	return values
}

// FieldValue TBD
func (x *ObjectActorRef) FieldValue(f ObjectActorRefField) interface{} {
	switch f {
	case ObjectActorRefFieldID:
		return x.ID
	case ObjectActorRefFieldObjectID:
		return x.ObjectID
	case ObjectActorRefFieldActorID:
		return x.ActorID
	}
	return nil
}

// SetFieldValue TBD
func (x *ObjectActorRef) SetFieldValue(f ObjectActorRefField, v interface{}) error {
	switch f {
	case ObjectActorRefFieldID:
		x.ID = v.(int64)
	case ObjectActorRefFieldObjectID:
		x.ObjectID = v.(int64)
	case ObjectActorRefFieldActorID:
		x.ActorID = v.(int64)
	}
	return nil
}

// ObjectActorRefQuerier TBD
type ObjectActorRefQuerier interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}

// ObjectActorRefRows TBD
type ObjectActorRefRows struct {
	rows   *pgx.Rows
	fields ObjectActorRefFieldsList
}

// Close TBD
func (r *ObjectActorRefRows) Close() {
	if r.rows == nil {
		return
	}
	r.rows.Close()
}

// Next TBD
func (r *ObjectActorRefRows) Next() bool {
	if r.rows == nil {
		return false
	}

	return r.rows.Next()
}

// Err TBD
func (r *ObjectActorRefRows) Err() error {
	if r.rows == nil {
		return errors.New("empty rows")
	}
	return r.rows.Err()
}

// ScanTo TBD
func (r *ObjectActorRefRows) ScanTo(x *ObjectActorRef) error {
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
func (r *ObjectActorRefRows) Scan() (x *ObjectActorRef, err error) {
	x = new(ObjectActorRef)
	err = r.ScanTo(x)
	return
}

// ScanAll TBD
func (r *ObjectActorRefRows) ScanAll(sizeHint int) ([]*ObjectActorRef, error) {
	defer r.Close()

	if sizeHint == 0 {
		sizeHint = 10
	}

	var list = make([]*ObjectActorRef, 0, sizeHint)

	for r.Next() {
		item, err := r.Scan()
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}

	return list, r.Err()
}

// ObjectActorRefRow TBD
type ObjectActorRefRow ObjectActorRefRows

// ScanTo TBD
func (r *ObjectActorRefRow) ScanTo(x *ObjectActorRef) error {
	rows := (*ObjectActorRefRows)(r)
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
func (r *ObjectActorRefRow) Scan() (x *ObjectActorRef, err error) {
	x = new(ObjectActorRef)
	err = r.ScanTo(x)
	return
}

// QueryObjectActorRef TBD
func QueryObjectActorRef(q ObjectActorRefQuerier, fields ObjectActorRefFieldsList, sql string, args ...interface{}) (*ObjectActorRefRows, error) {
	pgxRows, err := q.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	return &ObjectActorRefRows{rows: pgxRows, fields: fields}, nil
}

// QueryObjectActorRefRow TBD
func QueryObjectActorRefRow(q ObjectActorRefQuerier, fields ObjectActorRefFieldsList, sql string, args ...interface{}) *ObjectActorRefRow {
	pgxRows, _ := q.Query(sql, args...)
	return &ObjectActorRefRow{rows: pgxRows, fields: fields}
}

// ScanObjectActorRefList TBD
func ScanObjectActorRefList(q ObjectActorRefQuerier, fields ObjectActorRefFieldsList, sql string, args ...interface{}) ([]*ObjectActorRef, error) {
	rows, err := QueryObjectActorRef(q, fields, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows.ScanAll(0)
}

// ScanObjectActorRef TBD
func ScanObjectActorRef(q ObjectActorRefQuerier, fields ObjectActorRefFieldsList, sql string, args ...interface{}) (*ObjectActorRef, error) {
	x, err := QueryObjectActorRefRow(q, fields, sql, args...).Scan()
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return x, err
}
