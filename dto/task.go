package dto

import "time"

//go:generate dbgen -type Task
//go:generate dbgen -type TaskMediaRef
//go:generate dbgen -type TaskGroupRef
//go:generate dbgen -type TaskCollectionRef
//go:generate dbgen -type TaskObjectRef

// TaskStatus TBD
type TaskStatus string

// TaskStatus TBD
const (
	TaskStatusTODO       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

// IsTaskStatusValid TBD
func IsTaskStatusValid(s string) bool {
	switch TaskStatus(s) {
	case TaskStatusTODO, TaskStatusInProgress, TaskStatusDone:
		return true
	}
	return false
}

func (s TaskStatus) Text() string {
	switch s {
	case TaskStatusTODO:
		return "TODO"
	case TaskStatusInProgress:
		return "In progress"
	case TaskStatusDone:
		return "Done"
	default:
		return ""
	}
}

type (
	// Task TBD
	Task struct {
		ID             int64      `db:"id"`
		Title          string     `db:"title"`
		Description    string     `db:"description"`
		CreatorUserID  int64      `db:"creator_user_id"`
		AssignedUserID *int64     `db:"assigned_user_id"`
		Deadline       *time.Time `db:"deadline"`

		Status  TaskStatus `db:"status"`
		Archive bool       `db:"archive"`

		CreationTime time.Time `db:"creation_time"`
		UserUniqID   int64     `db:"user_uniq_id"`
	}

	// TaskList TBD
	TaskList []*Task

	// TaskMediaRef TBD
	TaskMediaRef struct {
		TaskID   int64 `db:"task_id"`
		MediaID  int64 `db:"media_id"`
		Position int32 `db:"position"`
	}

	// TaskGroupRef TBD
	TaskGroupRef struct {
		TaskID   int64 `db:"task_id"`
		GroupID  int64 `db:"group_id"`
		Position int32 `db:"position"`
	}

	// TaskCollectionRef TBD
	TaskCollectionRef struct {
		TaskID       int64 `db:"task_id"`
		CollectionID int64 `db:"collection_id"`
		Position     int32 `db:"position"`
	}

	// TaskObjectRef TBD
	TaskObjectRef struct {
		TaskID   int64 `db:"task_id"`
		ObjectID int64 `db:"object_id"`
		Position int32 `db:"position"`
	}
)

// IsValid TBD
func (t *Task) IsValid() bool {
	nameOK := len(t.Title) > 0 && len(t.Title) < 100
	statusOK := IsTaskStatusValid(string(t.Status))

	allOK := nameOK && statusOK
	return allOK
}

// GetUsersIDs TBD
func (t *Task) GetUsersIDs() []int64 {
	usersIDs := make([]int64, 0, 2)
	usersIDs = append(usersIDs, t.CreatorUserID)
	if t.AssignedUserID != nil {
		usersIDs = append(usersIDs, *t.AssignedUserID)
	}

	return usersIDs
}

// GetIDs TBD
func (tl TaskList) GetIDs() []int64 {
	ids := make([]int64, 0, len(tl))
	for _, task := range tl {
		ids = append(ids, task.ID)
	}
	return ids
}

// GetUsersIDs TBD
func (tl TaskList) GetUsersIDs() []int64 {
	usersIDsMap := make(map[int64]struct{})
	for _, task := range tl {
		for _, userID := range task.GetUsersIDs() {
			usersIDsMap[userID] = struct{}{}
		}
	}

	usersIDs := make([]int64, 0, len(usersIDsMap))
	for userID := range usersIDsMap {
		usersIDs = append(usersIDs, userID)
	}

	return usersIDs
}
