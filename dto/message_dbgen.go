// !! Autogenerated file, do not edit!
package dto

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"time"

	"github.com/jackc/pgx"
)

// MessageField TBD
type MessageField int

const (
	// MessageFieldID TBD
	MessageFieldID MessageField = iota

	// MessageFieldUserID TBD
	MessageFieldUserID

	// MessageFieldUserUniqID TBD
	MessageFieldUserUniqID

	// MessageFieldPeerID TBD
	MessageFieldPeerID

	// MessageFieldPeerType TBD
	MessageFieldPeerType

	// MessageFieldTypo TBD
	MessageFieldTypo

	// MessageFieldCreationTime TBD
	MessageFieldCreationTime

	// MessageFieldMessageExtra TBD
	MessageFieldMessageExtra
)

// MessageFieldsList TBD
type MessageFieldsList []MessageField

// MessageAllFields TBD
var MessageAllFields = MessageFieldsList{
	MessageFieldID,
	MessageFieldUserID,
	MessageFieldUserUniqID,
	MessageFieldPeerID,
	MessageFieldPeerType,
	MessageFieldTypo,
	MessageFieldCreationTime,
	MessageFieldMessageExtra,
}

// MessageFieldsNames TBD
var MessageFieldsNames = [...]string{
	"id",
	"user_id",
	"user_uniq_id",
	"peer_id",
	"peer_type",
	"type",
	"creation_time",
	"extra",
}

// Name TBD
func (f MessageField) Name() string {
	if int(f) > len(MessageFieldsNames)-1 {
		return "unknown"
	}
	return MessageFieldsNames[f]
}

// JoinedNames TBD
func (l MessageFieldsList) JoinedNames() string {
	return l.JoinedNamesWithAlias("")
}

// JoinedNamesWithAlias TBD
func (l MessageFieldsList) JoinedNamesWithAlias(alias string) string {
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
func (l MessageFieldsList) Placeholders() string {
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
func (l MessageFieldsList) Del(fields ...MessageField) MessageFieldsList {
	var res = make(MessageFieldsList, 0, len(l))
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
func (l MessageFieldsList) PushBack(fields ...MessageField) MessageFieldsList {
	var res = make(MessageFieldsList, 0, len(l)+len(fields))
	for _, f := range l {
		res = append(res, f)
	}
	for _, f := range fields {
		res = append(res, f)
	}
	return res
}

// PushFront TBD
func (l MessageFieldsList) PushFront(fields ...MessageField) MessageFieldsList {
	var res = make(MessageFieldsList, 0, len(l)+len(fields))
	for _, f := range fields {
		res = append(res, f)
	}
	for _, f := range l {
		res = append(res, f)
	}
	return res
}

// FieldsValues TBD
func (x *Message) FieldsValues(fields MessageFieldsList) []interface{} {
	values := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		values = append(values, x.FieldValue(f))
	}
	return values
}

// FieldValue TBD
func (x *Message) FieldValue(f MessageField) interface{} {
	switch f {
	case MessageFieldID:
		return x.ID
	case MessageFieldUserID:
		return x.UserID
	case MessageFieldUserUniqID:
		return x.UserUniqID
	case MessageFieldPeerID:
		return x.PeerID
	case MessageFieldPeerType:
		return int16(x.PeerType)
	case MessageFieldTypo:
		return int16(x.Typo)
	case MessageFieldCreationTime:
		return x.CreationTime
	case MessageFieldMessageExtra:
		data, _ := json.Marshal(x.MessageExtra)
		return string(data)
	}
	return nil
}

// SetFieldValue TBD
func (x *Message) SetFieldValue(f MessageField, v interface{}) error {
	switch f {
	case MessageFieldID:
		x.ID = v.(int64)
	case MessageFieldUserID:
		x.UserID = v.(int64)
	case MessageFieldUserUniqID:
		x.UserUniqID = v.(int64)
	case MessageFieldPeerID:
		x.PeerID = v.(int64)
	case MessageFieldPeerType:
		x.PeerType = PeerType(v.(int16))
	case MessageFieldTypo:
		x.Typo = MessageType(v.(int16))
	case MessageFieldCreationTime:
		x.CreationTime = v.(time.Time)
	case MessageFieldMessageExtra:
		if v != nil {
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				TagName: "json",
				Result:  &x.MessageExtra,
			})
			if err != nil {
				return err
			}
			err = decoder.Decode(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// MessageQuerier TBD
type MessageQuerier interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}

// MessageRows TBD
type MessageRows struct {
	rows   *pgx.Rows
	fields MessageFieldsList
}

// Close TBD
func (r *MessageRows) Close() {
	if r.rows == nil {
		return
	}
	r.rows.Close()
}

// Next TBD
func (r *MessageRows) Next() bool {
	if r.rows == nil {
		return false
	}

	return r.rows.Next()
}

// Err TBD
func (r *MessageRows) Err() error {
	if r.rows == nil {
		return errors.New("empty rows")
	}
	return r.rows.Err()
}

// ScanTo TBD
func (r *MessageRows) ScanTo(x *Message) error {
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
func (r *MessageRows) Scan() (x *Message, err error) {
	x = new(Message)
	err = r.ScanTo(x)
	return
}

// ScanAll TBD
func (r *MessageRows) ScanAll(sizeHint int) ([]*Message, error) {
	defer r.Close()

	if sizeHint == 0 {
		sizeHint = 10
	}

	var list = make([]*Message, 0, sizeHint)

	for r.Next() {
		item, err := r.Scan()
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}

	return list, r.Err()
}

// MessageRow TBD
type MessageRow MessageRows

// ScanTo TBD
func (r *MessageRow) ScanTo(x *Message) error {
	rows := (*MessageRows)(r)
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
func (r *MessageRow) Scan() (x *Message, err error) {
	x = new(Message)
	err = r.ScanTo(x)
	return
}

// QueryMessage TBD
func QueryMessage(q MessageQuerier, fields MessageFieldsList, sql string, args ...interface{}) (*MessageRows, error) {
	pgxRows, err := q.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	return &MessageRows{rows: pgxRows, fields: fields}, nil
}

// QueryMessageRow TBD
func QueryMessageRow(q MessageQuerier, fields MessageFieldsList, sql string, args ...interface{}) *MessageRow {
	pgxRows, _ := q.Query(sql, args...)
	return &MessageRow{rows: pgxRows, fields: fields}
}

// ScanMessageList TBD
func ScanMessageList(q MessageQuerier, fields MessageFieldsList, sql string, args ...interface{}) ([]*Message, error) {
	rows, err := QueryMessage(q, fields, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows.ScanAll(0)
}

// ScanMessage TBD
func ScanMessage(q MessageQuerier, fields MessageFieldsList, sql string, args ...interface{}) (*Message, error) {
	x, err := QueryMessageRow(q, fields, sql, args...).Scan()
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return x, err
}
