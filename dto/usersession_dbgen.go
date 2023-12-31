// !! Autogenerated file, do not edit!
package dto

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

// UserSessionField TBD
type UserSessionField int

const (
	// UserSessionFieldID TBD
	UserSessionFieldID UserSessionField = iota

	// UserSessionFieldUserID TBD
	UserSessionFieldUserID

	// UserSessionFieldAuthToken TBD
	UserSessionFieldAuthToken

	// UserSessionFieldCreationTime TBD
	UserSessionFieldCreationTime
)

// UserSessionFieldsList TBD
type UserSessionFieldsList []UserSessionField

// UserSessionAllFields TBD
var UserSessionAllFields = UserSessionFieldsList{
	UserSessionFieldID,
	UserSessionFieldUserID,
	UserSessionFieldAuthToken,
	UserSessionFieldCreationTime,
}

// UserSessionFieldsNames TBD
var UserSessionFieldsNames = [...]string{
	"id",
	"user_id",
	"auth_token",
	"creation_time",
}

// Name TBD
func (f UserSessionField) Name() string {
	if int(f) > len(UserSessionFieldsNames)-1 {
		return "unknown"
	}
	return UserSessionFieldsNames[f]
}

// JoinedNames TBD
func (l UserSessionFieldsList) JoinedNames() string {
	return l.JoinedNamesWithAlias("")
}

// JoinedNamesWithAlias TBD
func (l UserSessionFieldsList) JoinedNamesWithAlias(alias string) string {
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
func (l UserSessionFieldsList) Placeholders() string {
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
func (l UserSessionFieldsList) Del(fields ...UserSessionField) UserSessionFieldsList {
	var res = make(UserSessionFieldsList, 0, len(l))
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
func (l UserSessionFieldsList) PushBack(fields ...UserSessionField) UserSessionFieldsList {
	var res = make(UserSessionFieldsList, 0, len(l)+len(fields))
	for _, f := range l {
		res = append(res, f)
	}
	for _, f := range fields {
		res = append(res, f)
	}
	return res
}

// PushFront TBD
func (l UserSessionFieldsList) PushFront(fields ...UserSessionField) UserSessionFieldsList {
	var res = make(UserSessionFieldsList, 0, len(l)+len(fields))
	for _, f := range fields {
		res = append(res, f)
	}
	for _, f := range l {
		res = append(res, f)
	}
	return res
}

// FieldsValues TBD
func (x *UserSession) FieldsValues(fields UserSessionFieldsList) []interface{} {
	values := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		values = append(values, x.FieldValue(f))
	}
	return values
}

// FieldValue TBD
func (x *UserSession) FieldValue(f UserSessionField) interface{} {
	switch f {
	case UserSessionFieldID:
		return x.ID
	case UserSessionFieldUserID:
		return x.UserID
	case UserSessionFieldAuthToken:
		return x.AuthToken
	case UserSessionFieldCreationTime:
		return &x.CreationTime
	}
	return nil
}

// SetFieldValue TBD
func (x *UserSession) SetFieldValue(f UserSessionField, v interface{}) error {
	switch f {
	case UserSessionFieldID:
		x.ID = v.(int64)
	case UserSessionFieldUserID:
		x.UserID = v.(int64)
	case UserSessionFieldAuthToken:
		x.AuthToken = v.(string)
	case UserSessionFieldCreationTime:
		if v == nil {
			x.CreationTime = nil
		} else {
			value := v.(time.Time)
			x.CreationTime = &value
		}
	}
	return nil
}

// UserSessionQuerier TBD
type UserSessionQuerier interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}

// UserSessionRows TBD
type UserSessionRows struct {
	rows   *pgx.Rows
	fields UserSessionFieldsList
}

// Close TBD
func (r *UserSessionRows) Close() {
	if r.rows == nil {
		return
	}
	r.rows.Close()
}

// Next TBD
func (r *UserSessionRows) Next() bool {
	if r.rows == nil {
		return false
	}

	return r.rows.Next()
}

// Err TBD
func (r *UserSessionRows) Err() error {
	if r.rows == nil {
		return errors.New("empty rows")
	}
	return r.rows.Err()
}

// ScanTo TBD
func (r *UserSessionRows) ScanTo(x *UserSession) error {
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
func (r *UserSessionRows) Scan() (x *UserSession, err error) {
	x = new(UserSession)
	err = r.ScanTo(x)
	return
}

// ScanAll TBD
func (r *UserSessionRows) ScanAll(sizeHint int) ([]*UserSession, error) {
	defer r.Close()

	if sizeHint == 0 {
		sizeHint = 10
	}

	var list = make([]*UserSession, 0, sizeHint)

	for r.Next() {
		item, err := r.Scan()
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}

	return list, r.Err()
}

// UserSessionRow TBD
type UserSessionRow UserSessionRows

// ScanTo TBD
func (r *UserSessionRow) ScanTo(x *UserSession) error {
	rows := (*UserSessionRows)(r)
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
func (r *UserSessionRow) Scan() (x *UserSession, err error) {
	x = new(UserSession)
	err = r.ScanTo(x)
	return
}

// QueryUserSession TBD
func QueryUserSession(q UserSessionQuerier, fields UserSessionFieldsList, sql string, args ...interface{}) (*UserSessionRows, error) {
	pgxRows, err := q.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	return &UserSessionRows{rows: pgxRows, fields: fields}, nil
}

// QueryUserSessionRow TBD
func QueryUserSessionRow(q UserSessionQuerier, fields UserSessionFieldsList, sql string, args ...interface{}) *UserSessionRow {
	pgxRows, _ := q.Query(sql, args...)
	return &UserSessionRow{rows: pgxRows, fields: fields}
}

// ScanUserSessionList TBD
func ScanUserSessionList(q UserSessionQuerier, fields UserSessionFieldsList, sql string, args ...interface{}) ([]*UserSession, error) {
	rows, err := QueryUserSession(q, fields, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows.ScanAll(0)
}

// ScanUserSession TBD
func ScanUserSession(q UserSessionQuerier, fields UserSessionFieldsList, sql string, args ...interface{}) (*UserSession, error) {
	x, err := QueryUserSessionRow(q, fields, sql, args...).Scan()
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return x, err
}
