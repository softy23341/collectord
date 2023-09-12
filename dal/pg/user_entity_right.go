package dalpg

import (
	"fmt"

	"time"

	"git.softndit.com/collector/backend/dto"
)

// PutUserRight TBD
func (m *Manager) PutUserRight(right *dto.UserEntityRight) error {
	if err := m.deleteUserRight(right.UserID, right.EntityType, right.EntityID); err != nil {
		return err
	}

	if right.Level == dto.RightEntityLevelNone {
		return nil
	}

	return m.createUserRight(right)
}

// DeleteEntityRights TBD
func (m *Manager) DeleteEntityRights(entityType dto.RightEntityType, entityID int64) error {
	sql := `
DELETE FROM user_entity_right
WHERE entity_type = $1 AND entity_id = $2`

	_, err := m.p.Exec(sql, string(entityType), entityID)
	return err
}

// TBD DeleteUserRight
func (m *Manager) deleteUserRight(userID int64, entityType dto.RightEntityType, entityID int64) error {
	sql := `
DELETE FROM user_entity_right
WHERE user_id = $1 AND entity_type = $2 AND entity_id = $3`

	_, err := m.p.Exec(sql, userID, string(entityType), entityID)
	return err
}

func (m *Manager) createUserRight(right *dto.UserEntityRight) error {
	retFlds := dto.UserEntityRightFieldsList{
		dto.UserEntityRightFieldID,
		dto.UserEntityRightFieldCreationTime,
	}
	insFlds := dto.UserEntityRightAllFields.Del(retFlds...)

	sql := fmt.Sprintf(`INSERT INTO "user_entity_right"(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		insFlds.JoinedNames(), insFlds.Placeholders(), retFlds.JoinedNames())

	return dto.QueryUserEntityRightRow(m.p, retFlds, sql,
		right.FieldsValues(insFlds)...).ScanTo(right)
}

// HasUserRightsForObjects TBD
func (m *Manager) HasUserRightsForObjects(targetUserID int64, targetLevel dto.RightEntityLevel, objectsID []int64) (bool, error) {
	sql := fmt.Sprintf(`
SELECT o.id, COALESCE(uer.level, $1)
FROM object AS o
INNER JOIN collection AS c
  ON o.collection_id = c.id
LEFT JOIN user_entity_right AS uer
  ON uer.entity_id = c.id AND uer.entity_type = $2 AND uer.user_id = $4
WHERE o.id = any($3::bigint[])
`)

	rows, err := m.p.Query(sql, dto.RightEntityLevelNone, dto.RightEntityTypeCollection, objectsID, targetUserID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var (
		entityID int64
		level    string
	)
	for rows.Next() {
		if err := rows.Scan(&entityID, &level); err != nil {
			return false, err
		}
		if targetLevel.IsStronger(dto.RightEntityLevel(level)) {
			return false, nil
		}
	}

	if rows.Err() != nil {
		return false, err
	}

	return true, nil
}

// HasUserRightsForCollections TBD
func (m *Manager) HasUserRightsForCollections(targetUserID int64, targetLevel dto.RightEntityLevel, collectionsID []int64) (bool, error) {
	sql := fmt.Sprintf(`
SELECT c.id, COALESCE(uer.level, $1)
FROM collection AS c
LEFT JOIN user_entity_right AS uer
  ON uer.entity_id = c.id AND uer.entity_type = $2 AND uer.user_id = $4
WHERE c.id = any($3::bigint[])
`)

	rows, err := m.p.Query(
		sql,
		dto.RightEntityLevelNone,
		dto.RightEntityTypeCollection,
		collectionsID,
		targetUserID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var (
		entityID int64
		level    string
	)
	for rows.Next() {
		if err := rows.Scan(&entityID, &level); err != nil {
			return false, err
		}
		if targetLevel.IsStronger(dto.RightEntityLevel(level)) {
			return false, nil
		}
	}

	if rows.Err() != nil {
		return false, err
	}

	return true, nil
}

// HasUserRightsForGroups TBD
func (m *Manager) HasUserRightsForGroups(targetUserID int64, targetLevel dto.RightEntityLevel, groupsID []int64) (bool, error) {
	sql := fmt.Sprintf(`
SELECT cgr.group_id, COALESCE(uer.level, $1)
FROM collection_group_ref AS cgr
LEFT JOIN user_entity_right AS uer
  ON uer.entity_id = cgr.collection_id AND uer.entity_type = $2 AND uer.user_id = $4
WHERE cgr.group_id = any($3::bigint[])
`)

	rows, err := m.p.Query(sql, dto.RightEntityLevelNone, dto.RightEntityTypeCollection, groupsID, targetUserID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var (
		entityID int64
		level    string
	)
	for rows.Next() {
		if err := rows.Scan(&entityID, &level); err != nil {
			return false, err
		}
		// XXX
		level = string(dto.RightEntityLevelAdmin)
		if targetLevel.IsStronger(dto.RightEntityLevel(level)) {
			return false, nil
		}
	}

	if rows.Err() != nil {
		return false, err
	}

	return true, nil
}

// GetUserRightsForObjects TBD
func (m *Manager) GetUserRightsForObjects(targetUserID int64, objectsID []int64) (dto.ShortUserEntityRightList, error) {
	sql := fmt.Sprintf(`
SELECT o.id, COALESCE(uer.level, '%s')
FROM object AS o
INNER JOIN collection AS c
  ON o.collection_id = c.id
LEFT JOIN user_entity_right AS uer
  ON uer.entity_id = c.id AND uer.entity_type = $1
WHERE o.id = any($2::bigint[]) AND uer.user_id = $3
ORDER BY uer.id
`, dto.RightEntityLevelNone)

	rows, err := m.p.Query(sql, dto.RightEntityTypeCollection, objectsID, targetUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		entityID int64
		level    string
	)

	userEntityRightMap := make(map[int64]*dto.ShortUserEntityRight, len(objectsID))
	for rows.Next() {
		if err := rows.Scan(&entityID, &level); err != nil {
			return nil, err
		}

		userEntityRightMap[entityID] = &dto.ShortUserEntityRight{
			UserID:     targetUserID,
			EntityType: dto.RightEntityTypeObject,
			EntityID:   entityID,
			Level:      dto.RightEntityLevel(level),
		}
	}

	if rows.Err() != nil {
		return nil, err
	}

	userEntityRightList := make(dto.ShortUserEntityRightList, 0, len(objectsID))
	for _, entityID := range objectsID {
		right, found := userEntityRightMap[entityID]
		if !found {
			right = &dto.ShortUserEntityRight{
				UserID:     targetUserID,
				EntityType: dto.RightEntityTypeObject,
				EntityID:   entityID,
				Level:      dto.RightEntityLevelNone,
			}
		}
		userEntityRightList = append(userEntityRightList, right)
	}

	return userEntityRightList, nil

}

// GetUserRightsForCollections TBD
func (m *Manager) GetUserRightsForCollections(targetUserID int64, collectionsID []int64) (dto.ShortUserEntityRightList, error) {
	sql := fmt.Sprintf(`
SELECT c.id, COALESCE(uer.level, '%s')
FROM collection AS c
LEFT JOIN user_entity_right AS uer
  ON uer.entity_id = c.id AND uer.entity_type = $1
WHERE c.id = any($2::bigint[]) AND uer.user_id = $3
ORDER BY uer.id
`, dto.RightEntityLevelNone)

	rows, err := m.p.Query(sql, dto.RightEntityTypeCollection, collectionsID, targetUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		entityID int64
		level    string
	)

	userEntityRightMap := make(map[int64]*dto.ShortUserEntityRight, len(collectionsID))
	for rows.Next() {
		if err := rows.Scan(&entityID, &level); err != nil {
			return nil, err
		}

		userEntityRightMap[entityID] = &dto.ShortUserEntityRight{
			UserID:     targetUserID,
			EntityType: dto.RightEntityTypeCollection,
			EntityID:   entityID,
			Level:      dto.RightEntityLevel(level),
		}
	}

	if rows.Err() != nil {
		return nil, err
	}

	userEntityRightList := make(dto.ShortUserEntityRightList, 0, len(collectionsID))
	for _, entityID := range collectionsID {
		right, found := userEntityRightMap[entityID]
		if !found {
			right = &dto.ShortUserEntityRight{
				UserID:     targetUserID,
				EntityType: dto.RightEntityTypeCollection,
				EntityID:   entityID,
				Level:      dto.RightEntityLevelNone,
			}
		}
		userEntityRightList = append(userEntityRightList, right)
	}

	return userEntityRightList, nil
}

// GetUserRightsForGroups TBD
func (m *Manager) GetUserRightsForGroups(targetUserID int64, groupsID []int64) (dto.ShortUserEntityRightList, error) {

	// TODO remove a terrible code piece
	sql := fmt.Sprintf(`
SELECT
  cgr.group_id AS entity_id,
  MIN(CASE
    WHEN uer.level = '%s' THEN %d
    WHEN uer.level = '%s' THEN %d
    WHEN uer.level = '%s' THEN %d
    WHEN uer.level = '%s' THEN %d
    ELSE 0
  END)
FROM collection_group_ref AS cgr
LEFT JOIN user_entity_right AS uer
  ON uer.entity_id = cgr.collection_id AND uer.entity_type = $1
WHERE cgr.group_id = any($2::bigint[]) AND uer.user_id = $3
GROUP BY cgr.group_id
`,
		string(dto.RightEntityLevelAdmin), dto.AccessLevelToNum(dto.RightEntityLevelAdmin),
		string(dto.RightEntityLevelWrite), dto.AccessLevelToNum(dto.RightEntityLevelWrite),
		string(dto.RightEntityLevelRead), dto.AccessLevelToNum(dto.RightEntityLevelRead),
		string(dto.RightEntityLevelNone), dto.AccessLevelToNum(dto.RightEntityLevelNone),
	)

	rows, err := m.p.Query(sql, dto.RightEntityTypeCollection, groupsID, targetUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		entityID int64
		level    int
	)

	userEntityRightMap := make(map[int64]*dto.ShortUserEntityRight, len(groupsID))
	for rows.Next() {
		if err := rows.Scan(&entityID, &level); err != nil {
			return nil, err
		}

		userEntityRightMap[entityID] = &dto.ShortUserEntityRight{
			UserID:     targetUserID,
			EntityType: dto.RightEntityTypeGroup,
			EntityID:   entityID,
			// XXX
			Level: dto.RightEntityLevelAdmin, //dto.NumToAccessLevel(level),
		}
	}

	if rows.Err() != nil {
		return nil, err
	}

	userEntityRightList := make(dto.ShortUserEntityRightList, 0, len(groupsID))
	for _, entityID := range groupsID {
		right, found := userEntityRightMap[entityID]
		if !found {
			right = &dto.ShortUserEntityRight{
				UserID:     targetUserID,
				EntityType: dto.RightEntityTypeGroup,
				EntityID:   entityID,
				// XXX
				Level: dto.RightEntityLevelAdmin, // dto.RightEntityLevelNone,
			}
		}
		userEntityRightList = append(userEntityRightList, right)
	}

	return userEntityRightList, nil
}

// GetUserRightsForCollectionsInRoot TBD
func (m *Manager) GetUserRightsForCollectionsInRoot(targetUserID int64, rootID int64) (dto.ShortUserEntityRightList, error) {
	sql := fmt.Sprintf(`
SELECT c.id, COALESCE(uer.level, '%s')
FROM collection AS c
LEFT JOIN user_entity_right AS uer
  ON uer.entity_id = c.id AND uer.entity_type = $1 AND uer.user_id = $3
WHERE c.root_id = $2
ORDER BY uer.id
`, dto.RightEntityLevelNone)

	rows, err := m.p.Query(sql, dto.RightEntityTypeCollection, rootID, targetUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		entityID int64
		level    string
	)

	var userEntityRightList dto.ShortUserEntityRightList
	for rows.Next() {
		if err := rows.Scan(&entityID, &level); err != nil {
			return nil, err
		}

		userEntityRightList = append(userEntityRightList, &dto.ShortUserEntityRight{
			UserID:     targetUserID,
			EntityType: dto.RightEntityTypeCollection,
			EntityID:   entityID,
			Level:      dto.RightEntityLevel(level),
		})
	}

	if rows.Err() != nil {
		return nil, err
	}

	return userEntityRightList, nil
}

// GetUsersRightsForCollection TBD
func (m *Manager) GetUsersRightsForCollection(collectionID int64) (dto.ShortUserEntityRightList, error) {
	sql := `
SELECT u.id, COALESCE(uer.level, $1)
FROM "user" AS u
INNER JOIN collection AS c ON
  c.id = $2
LEFT JOIN user_entity_right AS uer ON
  uer.entity_type = 'collection' AND uer.entity_id = c.id AND uer.user_id = u.id
INNER JOIN user_root_ref AS urr ON
  urr.root_id = c.root_id AND u.id = urr.user_id
`

	rows, err := m.p.Query(sql, dto.RightEntityLevelNone, collectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		userID int64
		level  string
	)

	var userEntityRightList dto.ShortUserEntityRightList
	for rows.Next() {
		if err := rows.Scan(&userID, &level); err != nil {
			return nil, err
		}

		userEntityRightList = append(userEntityRightList, &dto.ShortUserEntityRight{
			UserID:     userID,
			EntityType: dto.RightEntityTypeCollection,
			EntityID:   collectionID,
			Level:      dto.RightEntityLevel(level),
		})
	}

	if rows.Err() != nil {
		return nil, err
	}

	return userEntityRightList, nil
}

// CheckUserTmpAccessToObject TBD
func (m *Manager) CheckUserTmpAccessToObject(userID, messageID, objectID int64) (bool, error) {
	flds := dto.MessageAllFields

	sql := fmt.Sprintf(`
		SELECT %s
		FROM "message" as m
		WHERE m.id = $1 AND m.peer_id = $2 AND 
		m.peer_type = $3 AND m.extra -> 'object' ->> 'ObjectID' = ($4)::text  
	`, flds.JoinedNamesWithAlias("m"))

	message, err := dto.ScanMessage(m.p, flds, sql, messageID, userID, dto.PeerTypeUser, fmt.Sprintf("%v", objectID))
	if err != nil {
		println(err.Error())
		return false, err
	}
	if message == nil {
		return false, nil
	}
	if time.Now().UTC().Sub(message.CreationTime.UTC()).Hours() > 24 {
		return false, nil
	}
	return true, nil
}

// CheckUserTmpAccessToMedia TBD
func (m *Manager) CheckUserTmpAccessToMedia(authToken string, messageID, mediaID int64) (bool, error) {
	flds := dto.MessageAllFields

	sql := fmt.Sprintf(`
		WITH media_object AS (
  SELECT o.id as id
  FROM "object" AS o
  JOIN object_media_ref
    ON object_media_ref.object_id = o.id
  WHERE object_media_ref.media_id = $1
),requester_user AS (
  SELECT user_session.user_id AS id 
  FROM user_session WHERE user_session.auth_token = $2
)
SELECT %s
FROM "message" as m, media_object, requester_user
WHERE m.id = $3 AND m.peer_id = requester_user.id AND m.peer_type = $4 AND m.extra -> 'object' ->> 'ObjectID' = media_object.id::text
	`, flds.JoinedNamesWithAlias("m"))

	message, err := dto.ScanMessage(m.p, flds, sql, mediaID, authToken, messageID, dto.PeerTypeUser)
	if err != nil {
		return false, err
	}
	if message == nil {
		return false, nil
	}
	if time.Now().UTC().Sub(message.CreationTime.UTC()).Hours() > 24 {
		return false, nil
	}
	return true, nil
}
