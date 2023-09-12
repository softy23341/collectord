package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetTasksByIDs TBD
func (m *Manager) GetTasksByIDs(IDs []int64) (dto.TaskList, error) {
	flds := dto.TaskAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM "task"
                WHERE id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanTaskList(m.p, flds, sql, IDs)
}

// objects
// CreateTaskObjectsRefs TBD
func (m *Manager) CreateTaskObjectsRefs(taskID int64, objectsIDs []int64) error {
	// TODO make batch insert
	for i, objectID := range objectsIDs {
		ref := &dto.TaskObjectRef{
			TaskID:   taskID,
			ObjectID: objectID,
			Position: int32(i),
		}
		if err := m.CreateTaskObjectRef(ref); err != nil {
			return err
		}
	}
	return nil
}

// CreateTaskObjectRef TBD
func (m *Manager) CreateTaskObjectRef(ref *dto.TaskObjectRef) error {
	flds := dto.TaskObjectRefAllFields

	sql := fmt.Sprintf(`INSERT INTO task_object_ref(%[1]s) VALUES(%[2]s)`,
		flds.JoinedNames(),
		flds.Placeholders(),
	)

	_, err := m.p.Exec(sql, (*dto.TaskObjectRef)(ref).FieldsValues(flds)...)
	return err
}

// collections
// CreateTaskCollectionsRefs TBD
func (m *Manager) CreateTaskCollectionsRefs(taskID int64, collectionsIDs []int64) error {
	// TODO make batch insert
	for i, collectionID := range collectionsIDs {
		ref := &dto.TaskCollectionRef{
			TaskID:       taskID,
			CollectionID: collectionID,
			Position:     int32(i),
		}
		if err := m.CreateTaskCollectionRef(ref); err != nil {
			return err
		}
	}
	return nil
}

// CreateTaskCollectionRef TBD
func (m *Manager) CreateTaskCollectionRef(ref *dto.TaskCollectionRef) error {
	flds := dto.TaskCollectionRefAllFields

	sql := fmt.Sprintf(`INSERT INTO task_collection_ref(%[1]s) VALUES(%[2]s)`,
		flds.JoinedNames(),
		flds.Placeholders(),
	)

	_, err := m.p.Exec(sql, (*dto.TaskCollectionRef)(ref).FieldsValues(flds)...)
	return err
}

// groups
// CreateTaskGroupsRefs TBD
func (m *Manager) CreateTaskGroupsRefs(taskID int64, groupsIDs []int64) error {
	// TODO make batch insert
	for i, groupID := range groupsIDs {
		ref := &dto.TaskGroupRef{
			TaskID:   taskID,
			GroupID:  groupID,
			Position: int32(i),
		}
		if err := m.CreateTaskGroupRef(ref); err != nil {
			return err
		}
	}
	return nil
}

// CreateTaskGroupRef TBD
func (m *Manager) CreateTaskGroupRef(ref *dto.TaskGroupRef) error {
	flds := dto.TaskGroupRefAllFields

	sql := fmt.Sprintf(`INSERT INTO task_group_ref(%[1]s) VALUES(%[2]s)`,
		flds.JoinedNames(),
		flds.Placeholders(),
	)

	_, err := m.p.Exec(sql, (*dto.TaskGroupRef)(ref).FieldsValues(flds)...)
	return err
}

// medias
// CreateTaskMediasRefs TBD
func (m *Manager) CreateTaskMediasRefs(taskID int64, mediasIDs []int64) error {
	// TODO make batch insert
	for i, mediaID := range mediasIDs {
		ref := &dto.TaskMediaRef{
			TaskID:   taskID,
			MediaID:  mediaID,
			Position: int32(i),
		}
		if err := m.CreateTaskMediaRef(ref); err != nil {
			return err
		}
	}
	return nil
}

// CreateTaskMediaRef TBD
func (m *Manager) CreateTaskMediaRef(ref *dto.TaskMediaRef) error {
	flds := dto.TaskMediaRefAllFields

	sql := fmt.Sprintf(`INSERT INTO task_media_ref(%[1]s) VALUES(%[2]s)`,
		flds.JoinedNames(),
		flds.Placeholders(),
	)

	_, err := m.p.Exec(sql, (*dto.TaskMediaRef)(ref).FieldsValues(flds)...)
	return err
}

// UpdateTask TBD
func (m *Manager) UpdateTask(taskID int64, editedTask *dto.Task, changedFields dto.TaskFieldsList) (*dto.Task, error) {
	var (
		retFlds = dto.TaskAllFields
		sql     string
		args    []interface{}
	)
	if len(changedFields) > 0 {
		pholders := MakePlaceholders(len(changedFields), 2)
		sql = fmt.Sprintf(`UPDATE task SET (%s) = (%s) WHERE id = $1 RETURNING %s`,
			changedFields.JoinedNames(), pholders, retFlds.JoinedNames())
		args = append([]interface{}{taskID}, editedTask.FieldsValues(changedFields)...)
	} else {
		sql = fmt.Sprintf(`SELECT %s FROM task WHERE id = $1`, retFlds.JoinedNames())
		args = []interface{}{taskID}
	}

	return dto.ScanTask(m.p, retFlds, sql, args...)
}

// CreateTask TBD
func (m *Manager) CreateTask(task *dto.Task) error {
	retFlds := dto.TaskFieldsList{
		dto.TaskFieldID,
		dto.TaskFieldCreationTime,
	}
	insFlds := dto.TaskAllFields.Del(retFlds...)

	sql := fmt.Sprintf(`INSERT INTO "task"(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		insFlds.JoinedNames(), insFlds.Placeholders(), retFlds.JoinedNames())

	return dto.QueryTaskRow(m.p, retFlds, sql, task.FieldsValues(insFlds)...).ScanTo(task)
}

// DeleteTaskObjectsRefs TBD
func (m *Manager) DeleteTaskObjectsRefs(taskID int64) error {
	sql := fmt.Sprintf(`DELETE FROM task_object_ref WHERE %s = $1`,
		dto.TaskObjectRefFieldTaskID.Name(),
	)

	_, err := m.p.Exec(sql, taskID)
	return err
}

// DeleteTaskCollectionsRefs TBD
func (m *Manager) DeleteTaskCollectionsRefs(taskID int64) error {
	sql := fmt.Sprintf(`DELETE FROM task_collection_ref WHERE %s = $1`,
		dto.TaskCollectionRefFieldTaskID.Name(),
	)

	_, err := m.p.Exec(sql, taskID)
	return err
}

// DeleteTaskGroupsRefs TBD
func (m *Manager) DeleteTaskGroupsRefs(taskID int64) error {
	sql := fmt.Sprintf(`DELETE FROM task_group_ref WHERE %s = $1`,
		dto.TaskGroupRefFieldTaskID.Name(),
	)

	_, err := m.p.Exec(sql, taskID)
	return err
}

// DeleteTaskMediasRefs TBD
func (m *Manager) DeleteTaskMediasRefs(taskID int64) error {
	sql := fmt.Sprintf(`DELETE FROM task_media_ref WHERE %s = $1`,
		dto.TaskMediaRefFieldTaskID.Name(),
	)

	_, err := m.p.Exec(sql, taskID)
	return err
}

// GetTaskObjects TBD
func (m *Manager) GetTaskObjects(taskID int64) (dto.ObjectList, error) {
	flds := dto.ObjectAllFields
	sql := fmt.Sprintf(`
SELECT %s
FROM object AS o
INNER JOIN task_object_ref AS tor
  ON tor.object_id = o.id
WHERE tor.task_id = $1
ORDER BY o.id ASC
`, flds.JoinedNamesWithAlias("o"))

	return dto.ScanObjectList(m.p, flds, sql, taskID)
}

// GetTaskCollections TBD
func (m *Manager) GetTaskCollections(taskID int64) (dto.CollectionList, error) {
	flds := dto.CollectionAllFields
	sql := fmt.Sprintf(`
SELECT %s
FROM collection AS c
INNER JOIN task_collection_ref AS tcr
  ON tcr.collection_id = c.id
WHERE tcr.task_id = $1
ORDER BY c.id ASC
`, flds.JoinedNamesWithAlias("c"))

	return dto.ScanCollectionList(m.p, flds, sql, taskID)
}

// GetTaskGroups TBD
func (m *Manager) GetTaskGroups(taskID int64) (dto.GroupList, error) {
	flds := dto.GroupAllFields
	sql := fmt.Sprintf(`
SELECT %s
FROM "group" AS g
INNER JOIN task_group_ref AS tgr
  ON tgr.group_id = g.id
WHERE tgr.task_id = $1
ORDER BY g.id ASC
`, flds.JoinedNamesWithAlias("g"))

	return dto.ScanGroupList(m.p, flds, sql, taskID)
}

// GetTaskMedias TBD
func (m *Manager) GetTaskMedias(taskID int64) (dto.MediaList, error) {
	flds := dto.MediaAllFields
	sql := fmt.Sprintf(`
SELECT %s
FROM media AS m
INNER JOIN task_media_ref AS tmr
  ON tmr.media_id = m.id
WHERE tmr.task_id = $1
ORDER BY m.id ASC
`, flds.JoinedNamesWithAlias("m"))

	return dto.ScanMediaList(m.p, flds, sql, taskID)
}

// DeleteTask TBD
func (m *Manager) DeleteTask(taskID int64) error {
	sql := "DELETE FROM task WHERE id = $1"

	_, err := m.p.Exec(sql, taskID)
	return err
}

// GetUserRelatedTasks TBD
func (m *Manager) GetUserRelatedTasks(userID int64) (dto.TaskList, error) {
	flds := dto.TaskAllFields
	sql := fmt.Sprintf(`
SELECT %s
FROM task
WHERE (creator_user_id = $1 OR assigned_user_id = $1)
AND NOT archive
ORDER BY status, creation_time DESC
`, flds.JoinedNames())

	return dto.ScanTaskList(m.p, flds, sql, userID)
}

// GetAssignedToUserTasks TBD
func (m *Manager) GetAssignedToUserTasksFrom(creatorID int64, userID int64) (dto.TaskList, error) {
	flds := dto.TaskAllFields
	sql := fmt.Sprintf(`
SELECT %s
FROM task
WHERE creator_user_id = $1 AND assigned_user_id = $2
AND NOT archive
ORDER BY status, creation_time DESC
`, flds.JoinedNames())

	return dto.ScanTaskList(m.p, flds, sql, creatorID, userID)
}

// GetAssignedToUserTasks TBD
func (m *Manager) GetAssignedToUserTasks(userID int64) (dto.TaskList, error) {
	flds := dto.TaskAllFields
	sql := fmt.Sprintf(`
SELECT %s
FROM task
WHERE assigned_user_id = $1
AND NOT archive
ORDER BY status, creation_time DESC
`, flds.JoinedNames())

	return dto.ScanTaskList(m.p, flds, sql, userID)
}

// GetUserRelatedTasksArchive TBD
func (m *Manager) GetUserRelatedTasksArchive(userID int64, paginator *dto.PagePaginator) (dto.TaskList, error) {
	flds := dto.TaskAllFields
	sql := fmt.Sprintf(`
SELECT %s
FROM task
WHERE (creator_user_id = $1 OR assigned_user_id = $1)
AND archive
ORDER BY status, creation_time DESC
LIMIT $2 OFFSET $3
`, flds.JoinedNames())

	return dto.ScanTaskList(m.p, flds, sql, userID, paginator.Limit(), paginator.Offset())
}

// GetMediasCntByTaskID TBD
func (m *Manager) GetMediasCntByTaskID(tasksIDs []int64) (map[int64]int64, error) {
	sql := `
SELECT t.id, COUNT(tmr) AS cnt
FROM task AS t
LEFT JOIN task_media_ref AS tmr
  ON tmr.task_id = t.id
WHERE t.id = any($1::bigint[])
GROUP BY t.id
        `

	rows, err := m.p.Query(sql, tasksIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		taskID int64
		cnt    int64
	)
	task2cnt := make(map[int64]int64)
	for rows.Next() {
		if err := rows.Scan(&taskID, &cnt); err != nil {
			return nil, err
		}
		task2cnt[taskID] = cnt
	}
	return task2cnt, nil
}

// GetTasksCntByUserID TBD
func (m *Manager) GetTasksCntByUserID(creatorID int64, usersIDs []int64) (map[int64]int64, error) {
	sql := `
		SELECT t.assigned_user_id as user_id, COUNT(t) AS cnt
		FROM task AS t
		WHERE t.creator_user_id = $1 AND t.assigned_user_id = any($2::bigint[]) AND t.archive IS false
		GROUP BY t.assigned_user_id;`

	rows, err := m.p.Query(sql, creatorID, usersIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		userID int64
		cnt    int64
	)
	user2taskCnt := make(map[int64]int64)
	for rows.Next() {
		if err := rows.Scan(&userID, &cnt); err != nil {
			return nil, err
		}
		user2taskCnt[userID] = cnt
	}
	return user2taskCnt, nil
}
