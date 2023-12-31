// !! Autogenerated file, do not edit!
package dto

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

// TaskField TBD
type TaskField int

const (
	// TaskFieldID TBD
	TaskFieldID TaskField = iota

	// TaskFieldTitle TBD
	TaskFieldTitle

	// TaskFieldDescription TBD
	TaskFieldDescription

	// TaskFieldCreatorUserID TBD
	TaskFieldCreatorUserID

	// TaskFieldAssignedUserID TBD
	TaskFieldAssignedUserID

	// TaskFieldDeadline TBD
	TaskFieldDeadline

	// TaskFieldStatus TBD
	TaskFieldStatus

	// TaskFieldArchive TBD
	TaskFieldArchive

	// TaskFieldCreationTime TBD
	TaskFieldCreationTime

	// TaskFieldUserUniqID TBD
	TaskFieldUserUniqID
)

// TaskFieldsList TBD
type TaskFieldsList []TaskField

// TaskAllFields TBD
var TaskAllFields = TaskFieldsList{
	TaskFieldID,
	TaskFieldTitle,
	TaskFieldDescription,
	TaskFieldCreatorUserID,
	TaskFieldAssignedUserID,
	TaskFieldDeadline,
	TaskFieldStatus,
	TaskFieldArchive,
	TaskFieldCreationTime,
	TaskFieldUserUniqID,
}

// TaskFieldsNames TBD
var TaskFieldsNames = [...]string{
	"id",
	"title",
	"description",
	"creator_user_id",
	"assigned_user_id",
	"deadline",
	"status",
	"archive",
	"creation_time",
	"user_uniq_id",
}

// Name TBD
func (f TaskField) Name() string {
	if int(f) > len(TaskFieldsNames)-1 {
		return "unknown"
	}
	return TaskFieldsNames[f]
}

// JoinedNames TBD
func (l TaskFieldsList) JoinedNames() string {
	return l.JoinedNamesWithAlias("")
}

// JoinedNamesWithAlias TBD
func (l TaskFieldsList) JoinedNamesWithAlias(alias string) string {
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
func (l TaskFieldsList) Placeholders() string {
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
func (l TaskFieldsList) Del(fields ...TaskField) TaskFieldsList {
	var res = make(TaskFieldsList, 0, len(l))
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
func (l TaskFieldsList) PushBack(fields ...TaskField) TaskFieldsList {
	var res = make(TaskFieldsList, 0, len(l)+len(fields))
	for _, f := range l {
		res = append(res, f)
	}
	for _, f := range fields {
		res = append(res, f)
	}
	return res
}

// PushFront TBD
func (l TaskFieldsList) PushFront(fields ...TaskField) TaskFieldsList {
	var res = make(TaskFieldsList, 0, len(l)+len(fields))
	for _, f := range fields {
		res = append(res, f)
	}
	for _, f := range l {
		res = append(res, f)
	}
	return res
}

// FieldsValues TBD
func (x *Task) FieldsValues(fields TaskFieldsList) []interface{} {
	values := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		values = append(values, x.FieldValue(f))
	}
	return values
}

// FieldValue TBD
func (x *Task) FieldValue(f TaskField) interface{} {
	switch f {
	case TaskFieldID:
		return x.ID
	case TaskFieldTitle:
		return x.Title
	case TaskFieldDescription:
		return x.Description
	case TaskFieldCreatorUserID:
		return x.CreatorUserID
	case TaskFieldAssignedUserID:
		return &x.AssignedUserID
	case TaskFieldDeadline:
		return &x.Deadline
	case TaskFieldStatus:
		return string(x.Status)
	case TaskFieldArchive:
		return x.Archive
	case TaskFieldCreationTime:
		return x.CreationTime
	case TaskFieldUserUniqID:
		return x.UserUniqID
	}
	return nil
}

// SetFieldValue TBD
func (x *Task) SetFieldValue(f TaskField, v interface{}) error {
	switch f {
	case TaskFieldID:
		x.ID = v.(int64)
	case TaskFieldTitle:
		x.Title = v.(string)
	case TaskFieldDescription:
		x.Description = v.(string)
	case TaskFieldCreatorUserID:
		x.CreatorUserID = v.(int64)
	case TaskFieldAssignedUserID:
		if v == nil {
			x.AssignedUserID = nil
		} else {
			value := v.(int64)
			x.AssignedUserID = &value
		}
	case TaskFieldDeadline:
		if v == nil {
			x.Deadline = nil
		} else {
			value := v.(time.Time)
			x.Deadline = &value
		}
	case TaskFieldStatus:
		x.Status = TaskStatus(v.(string))
	case TaskFieldArchive:
		x.Archive = v.(bool)
	case TaskFieldCreationTime:
		x.CreationTime = v.(time.Time)
	case TaskFieldUserUniqID:
		x.UserUniqID = v.(int64)
	}
	return nil
}

// TaskQuerier TBD
type TaskQuerier interface {
	Query(sql string, args ...interface{}) (*pgx.Rows, error)
	QueryRow(sql string, args ...interface{}) *pgx.Row
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
}

// TaskRows TBD
type TaskRows struct {
	rows   *pgx.Rows
	fields TaskFieldsList
}

// Close TBD
func (r *TaskRows) Close() {
	if r.rows == nil {
		return
	}
	r.rows.Close()
}

// Next TBD
func (r *TaskRows) Next() bool {
	if r.rows == nil {
		return false
	}

	return r.rows.Next()
}

// Err TBD
func (r *TaskRows) Err() error {
	if r.rows == nil {
		return errors.New("empty rows")
	}
	return r.rows.Err()
}

// ScanTo TBD
func (r *TaskRows) ScanTo(x *Task) error {
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
func (r *TaskRows) Scan() (x *Task, err error) {
	x = new(Task)
	err = r.ScanTo(x)
	return
}

// ScanAll TBD
func (r *TaskRows) ScanAll(sizeHint int) ([]*Task, error) {
	defer r.Close()

	if sizeHint == 0 {
		sizeHint = 10
	}

	var list = make([]*Task, 0, sizeHint)

	for r.Next() {
		item, err := r.Scan()
		if err != nil {
			return list, err
		}
		list = append(list, item)
	}

	return list, r.Err()
}

// TaskRow TBD
type TaskRow TaskRows

// ScanTo TBD
func (r *TaskRow) ScanTo(x *Task) error {
	rows := (*TaskRows)(r)
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
func (r *TaskRow) Scan() (x *Task, err error) {
	x = new(Task)
	err = r.ScanTo(x)
	return
}

// QueryTask TBD
func QueryTask(q TaskQuerier, fields TaskFieldsList, sql string, args ...interface{}) (*TaskRows, error) {
	pgxRows, err := q.Query(sql, args...)
	if err != nil {
		return nil, err
	}

	return &TaskRows{rows: pgxRows, fields: fields}, nil
}

// QueryTaskRow TBD
func QueryTaskRow(q TaskQuerier, fields TaskFieldsList, sql string, args ...interface{}) *TaskRow {
	pgxRows, _ := q.Query(sql, args...)
	return &TaskRow{rows: pgxRows, fields: fields}
}

// ScanTaskList TBD
func ScanTaskList(q TaskQuerier, fields TaskFieldsList, sql string, args ...interface{}) ([]*Task, error) {
	rows, err := QueryTask(q, fields, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows.ScanAll(0)
}

// ScanTask TBD
func ScanTask(q TaskQuerier, fields TaskFieldsList, sql string, args ...interface{}) (*Task, error) {
	x, err := QueryTaskRow(q, fields, sql, args...).Scan()
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return x, err
}
