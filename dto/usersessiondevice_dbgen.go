// !! Autogenerated file, do not edit!
package dto

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
)

// UserSessionDeviceField TBD
type UserSessionDeviceField int

const (
	// UserSessionDeviceFieldID TBD
	UserSessionDeviceFieldID UserSessionDeviceField = iota

	// UserSessionDeviceFieldSessionID TBD
	UserSessionDeviceFieldSessionID

	// UserSessionDeviceFieldTypo TBD
	UserSessionDeviceFieldTypo

	// UserSessionDeviceFieldToken TBD
	UserSessionDeviceFieldToken

	// UserSessionDeviceFieldSandbox TBD
	UserSessionDeviceFieldSandbox
)

// UserSessionDeviceFieldsList TBD
type UserSessionDeviceFieldsList []UserSessionDeviceField

// UserSessionDeviceAllFields TBD
var UserSessionDeviceAllFields = UserSessionDeviceFieldsList{
	UserSessionDeviceFieldID,
	UserSessionDeviceFieldSessionID,
	UserSessionDeviceFieldTypo,
	UserSessionDeviceFieldToken,
	UserSessionDeviceFieldSandbox,
}

// UserSessionDeviceFieldsNames TBD
var UserSessionDeviceFieldsNames = [...]string{
	"id",
	"session_id",
	"typo",
	"token",
	"push_notification_sandbox",
}

// Name TBD
func (f UserSessionDeviceField) Name() string {
	if int(f) > len(UserSessionDeviceFieldsNames)-1 {
		return "unknown"
	}
	return UserSessionDeviceFieldsNames[f]
}

// JoinedNames TBD
func (l UserSessionDeviceFieldsList) JoinedNames() string {
	return l.JoinedNamesWithAlias("")
}

// JoinedNamesWithAlias TBD
func (l UserSessionDeviceFieldsList) JoinedNamesWithAlias(alias string) string {
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
func (l UserSessionDeviceFieldsList) Placeholders() string {
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
func (l UserSessionDeviceFieldsList) Del(fields ...UserSessionDeviceField) UserSessionDeviceFieldsList {
	var res = make(UserSessionDeviceFieldsList, 0, len(l))
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
func (l UserSessionDeviceFieldsList) PushBack(fields ...UserSessionDeviceField) UserSessionDeviceFieldsList {
	var res = make(UserSessionDeviceFieldsList, 0, len(l)+len(fields))
	for _, f := range l {
		res = append(res, f)
	}
	for _, f := range fields {
		res = append(res, f)
	}
	return res
}

// PushFront TBD
func (l UserSessionDeviceFieldsList) PushFront(fields ...UserSessionDeviceField) UserSessionDeviceFieldsList {
	var res = make(UserSessionDeviceFieldsList, 0, len(l)+len(fields))
	for _, f := range fields {
		res = append(res, f)
	}
	for _, f := range l {
		res = append(res, f)
	}
	return res
}

// FieldsValues TBD
func (x *UserSessionDevice) FieldsValues(fields UserSessionDeviceFieldsList) []interface{} {
	values := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		values = append(values, x.FieldValue(f))
	}
	return values
}

// FieldValue TBD
func (x *UserSessionDevice) FieldValue(f UserSessionDeviceField) interface{} {
	switch f {
	case UserSessionDeviceFieldID:
		return x.ID
	case UserSessionDeviceFieldSessionID:
		return x.SessionID
	case UserSessionDeviceFieldTypo:
		return int16(x.Typo)
	case UserSessionDeviceFieldToken:
		return x.Token
	case UserSessionDeviceFieldSandbox:
		return x.Sandbox
	}
	return nil
}

// SetFieldValue TBD
func (x *UserSessionDevice) SetFieldValue(f UserSessionDeviceField, v interface{}) error {
	switch f {
	case UserSessionDeviceFieldID:
		x.ID = v.(int64)
	case UserSessionDeviceFieldSessionID:
		x.SessionID = v.(int64)
	case UserSessionDeviceFieldTypo:
		x.Typo = UserSessionDeviceType(v.(int16))
	case UserSessionDeviceFieldToken:
		x.Token = v.(string)
	case UserSessionDeviceFieldSandbox:
		x.Sandbox = v.(bool)
	}
	return nil
}

// UserSessionDeviceQuerier TBD
type UserSessionDeviceQuerier interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}

// UserSessionDeviceRows TBD
type UserSessionDeviceRows struct {
	rows   *pgx.Rows
	fields UserSessionDeviceFieldsList
}

// Close TBD
func (r *UserSessionDeviceRows) Close() {
	if r.rows == nil {
		return
	}
	r.rows.Close()
}

// Next TBD
func (r *UserSessionDeviceRows) Next() bool {
	if r.rows == nil {
		return false
	}

	return r.rows.Next()
}

// Err TBD
func (r *UserSessionDeviceRows) Err() error {
	if r.rows == nil {
		return errors.New("empty rows")
	}
	return r.rows.Err()
}

// ScanTo TBD
func (r *UserSessionDeviceRows) ScanTo(x *UserSessionDevice) error {
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
func (r *UserSessionDeviceRows) Scan() (x *UserSessionDevice, err error) {
	x = new(UserSessionDevice)
	err = r.ScanTo(x)
	return
}

// ScanAll TBD
func (r *UserSessionDeviceRows) ScanAll(sizeHint int) ([]*UserSessionDevice, error) {
	defer r.Close()

	if sizeHint == 0 {
		sizeHint = 10
	}

	var list = make([]*UserSessionDevice, 0, sizeHint)

	for r.Next() {
		item, err := r.Scan()
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}

	return list, r.Err()
}

// UserSessionDeviceRow TBD
type UserSessionDeviceRow UserSessionDeviceRows

// ScanTo TBD
func (r *UserSessionDeviceRow) ScanTo(x *UserSessionDevice) error {
	rows := (*UserSessionDeviceRows)(r)
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
func (r *UserSessionDeviceRow) Scan() (x *UserSessionDevice, err error) {
	x = new(UserSessionDevice)
	err = r.ScanTo(x)
	return
}

// QueryUserSessionDevice TBD
func QueryUserSessionDevice(q UserSessionDeviceQuerier, fields UserSessionDeviceFieldsList, sql string, args ...interface{}) (*UserSessionDeviceRows, error) {
	pgxRows, err := q.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	return &UserSessionDeviceRows{rows: pgxRows, fields: fields}, nil
}

// QueryUserSessionDeviceRow TBD
func QueryUserSessionDeviceRow(q UserSessionDeviceQuerier, fields UserSessionDeviceFieldsList, sql string, args ...interface{}) *UserSessionDeviceRow {
	pgxRows, _ := q.Query(sql, args...)
	return &UserSessionDeviceRow{rows: pgxRows, fields: fields}
}

// ScanUserSessionDeviceList TBD
func ScanUserSessionDeviceList(q UserSessionDeviceQuerier, fields UserSessionDeviceFieldsList, sql string, args ...interface{}) ([]*UserSessionDevice, error) {
	rows, err := QueryUserSessionDevice(q, fields, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows.ScanAll(0)
}

// ScanUserSessionDevice TBD
func ScanUserSessionDevice(q UserSessionDeviceQuerier, fields UserSessionDeviceFieldsList, sql string, args ...interface{}) (*UserSessionDevice, error) {
	x, err := QueryUserSessionDeviceRow(q, fields, sql, args...).Scan()
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return x, err
}
