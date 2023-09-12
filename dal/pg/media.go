package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetMediasByObjectID TBD
func (m *Manager) GetMediasByObjectID(objectID int64) (dto.MediaList, error) {
	flds := dto.MediaAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM media AS m
                JOIN object_media_ref AS omr ON omr.media_id = m.id
                WHERE omr.object_id = $1;
        `, flds.JoinedNamesWithAlias("m"))

	return dto.ScanMediaList(m.p, flds, sql, objectID)
}

// GetMediaByUserUniqID TBD
func (m *Manager) GetMediaByUserUniqID(userID, uniqID int64) (*dto.Media, error) {
	flds := dto.MediaAllFields

	sqlQuery := fmt.Sprintf(`
                SELECT %s
                FROM media
                WHERE user_id = $1
                AND user_uniq_id = $2
        `, flds.JoinedNames())

	return dto.ScanMedia(m.p, flds, sqlQuery, userID, uniqID)
}

// CreateMedia TBD
func (m *Manager) CreateMedia(media *dto.Media) error {
	flds := dto.MediaAllFields.Del(dto.MediaFieldID)

	sql := fmt.Sprintf(`INSERT INTO media(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		flds.JoinedNames(),
		flds.Placeholders(),
		dto.MediaFieldID.Name(),
	)

	return dto.QueryMediaRow(
		m.p,
		dto.MediaFieldsList{dto.MediaFieldID},
		sql,
		(*dto.Media)(media).FieldsValues(flds)...,
	).ScanTo(media)
}

// GetMediasByIDs TBD
func (m *Manager) GetMediasByIDs(mediasIDs []int64) (dto.MediaList, error) {
	if len(mediasIDs) == 0 {
		return nil, nil
	}

	flds := dto.MediaAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM media
                WHERE id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanMediaList(m.p, flds, sql, mediasIDs)

}

// CreateObjectMediasRefs TBD
func (m *TxManager) CreateObjectMediasRefs(objectID int64, mediaIDs []int64) error {
	// TODO make batch insert
	for position, mediaID := range mediaIDs {
		ref := &dto.ObjectMediaRef{
			ObjectID:      objectID,
			MediaID:       mediaID,
			MediaPosition: int32(position),
		}

		if err := m.createObjectMediaRef(ref); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) createObjectMediaRef(ref *dto.ObjectMediaRef) error {
	flds := dto.ObjectMediaRefAllFields

	sql := fmt.Sprintf(`INSERT INTO object_media_ref(%[1]s) VALUES(%[2]s)`,
		flds.JoinedNames(),
		flds.Placeholders(),
	)

	_, err := m.p.Exec(sql, (*dto.ObjectMediaRef)(ref).FieldsValues(flds)...)
	return err
}

// DeleteObjectMediasRefs TBD
func (m *Manager) DeleteObjectMediasRefs(objectID int64) error {
	sql := fmt.Sprintf(`DELETE FROM object_media_ref WHERE %s = $1`,
		dto.ObjectMediaRefFieldObjectID.Name(),
	)

	_, err := m.p.Exec(sql, objectID)
	return err
}

// GetObjectsMediaRefs TBD
func (m *Manager) GetObjectsMediaRefs(objectsIDs []int64) (dto.ObjectMediaRefList, error) {
	flds := dto.ObjectMediaRefAllFields

	sql := fmt.Sprintf(`
	        SELECT %s
                FROM object_media_ref
                WHERE object_id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanObjectMediaRefList(m.p, flds, sql, objectsIDs)
}

// GetObjectsMediaRefsPhoto TBD
func (m *Manager) GetObjectsMediaRefsPhoto(objectsIDs []int64) (dto.ObjectMediaRefList, error) {
	flds := dto.ObjectMediaRefAllFields

	sql := fmt.Sprintf(`
	        SELECT %s
                FROM object_media_ref AS omr
                INNER JOIN media AS m
                  ON m.id = omr.media_id
                WHERE object_id = any($1::bigint[])
                AND m.type BETWEEN 1 AND 100
        `, flds.JoinedNamesWithAlias("omr"))

	return dto.ScanObjectMediaRefList(m.p, flds, sql, objectsIDs)
}

// GetMediaByPage TBD
func (m *Manager) GetMediaByPage(types []dto.MediaType, p *dto.PagePaginator) (dto.MediaList, error) {
	flds := dto.MediaAllFields

	int16types := make([]int16, len(types))
	for i := range types {
		int16types[i] = int16(types[i])
	}

	sql := fmt.Sprintf(`
                SELECT %s
                FROM media
                WHERE type = any($1::smallint[])
                ORDER BY id ASC
                OFFSET $2 LIMIT $3
        `, flds.JoinedNames())

	return dto.ScanMediaList(m.p, flds, sql, int16types, p.Offset(), p.Limit())
}

// UpdateMedia TBD
func (m *Manager) UpdateMedia(media *dto.Media) error {
	sql := `
	  UPDATE "media"
          SET
            user_id = $1
          , user_uniq_id = $2
          , type = $3
          , root_id = $4
          , extra = $5
          WHERE id = $6
        `

	_, err := m.p.Exec(sql,
		media.UserID,
		media.UserUniqID,
		int16(media.Type),
		media.RootID,
		media.ExtraJSON(),
		media.ID,
	)
	return err
}

// IsMediaUsed TBD
func (m *Manager) IsMediaUsed(mediaID int64) (bool, error) {
	sql := `
         SELECT COUNT(medias.*) AS cnt FROM (
           SELECT avatar_media_id AS media_id FROM "user" WHERE avatar_media_id IS NOT NULL
           UNION
           SELECT avatar_media_id AS media_id FROM chat WHERE avatar_media_id IS NOT NULL
           UNION
           SELECT image_media_id AS media_id FROM collection WHERE image_media_id IS NOT NULL
           UNION
           SELECT image_media_id AS media_id FROM object_status WHERE image_media_id IS NOT NULL
           UNION
           SELECT media_id AS media_id FROM object_media_ref WHERE media_id IS NOT NULL
         ) AS medias WHERE media_id = $1
        `
	rows, err := m.p.Query(sql, mediaID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var (
		cnt int64
	)
	for rows.Next() {
		if err := rows.Scan(&cnt); err != nil {
			return false, err
		}
	}

	if rows.Err() != nil {
		return false, err
	}

	return cnt > 0, nil
}

// DeleteMedias TBD
func (m *Manager) DeleteMedias(mediasIDs []int64) error {
	sql := fmt.Sprintf(`DELETE FROM media WHERE id = any($1::bigint[])`)

	_, err := m.p.Exec(sql, mediasIDs)
	return err
}

// GetMediaByIDAndVariantURI TBD
func (m *Manager) GetMediaByIDAndVariantURI(mediaID int64, variantURI string) (*dto.Media, error) {
	flds := dto.MediaAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM media
                WHERE media.ID = $1 AND (
    				(media.extra -> 'document' ->> 'uri' = $2) OR
    				(SELECT $2 IN (SELECT jsonb_array_elements(media.extra -> 'photo' -> 'variants') ->> 'uri')
				)
  			)
	`, flds.JoinedNames())
	return dto.ScanMedia(m.p, flds, sql, mediaID, variantURI)
}

// CanUserGetMediaByAuthToken TBD
func (m *Manager) CanUserGetMediaByAuthToken(authToken string, mediaID int64) (bool, error) {
	sql := `
WITH user_media AS (
  SELECT media.*, user_session.user_id AS requester_user_id
  FROM media AS media
  INNER JOIN user_session AS user_session
    ON user_session.auth_token = $1
  WHERE media.ID = $2 
)
SELECT
  SUM(object_is_present) AS object_is_present,
  SUM(object_is_ok) AS object_is_ok,
  SUM(object_is_public) AS object_is_public,
  SUM(collection_is_present) AS collection_is_present,
  SUM(collection_is_ok) AS collection_is_ok,
  SUM(collection_is_public) AS collection_is_public
FROM (
  -- init values
  SELECT 0 AS object_is_present, 0 AS object_is_ok, 0 AS object_is_public, 0 AS collection_is_present, 0 AS collection_is_ok, 0 AS collection_is_public
  UNION
  -- check objects
  SELECT 
    COALESCE(object.id, 0) AS object_is_present, 
    COALESCE(collection_uer.id, 0) AS object_is_ok, 
    0 AS object_is_public, 
    0 AS collection_is_present, 
    0 AS collection_is_ok,
    0 AS collection_is_public
  FROM user_media AS user_media
  LEFT JOIN object_media_ref AS omr
    ON omr.media_id = user_media.id
  LEFT JOIN object AS object
    ON object.id = omr.object_id
  LEFT JOIN collection AS collection
    ON collection.id = object.collection_id
  LEFT JOIN user_entity_right AS collection_uer
    ON (
      collection_uer.entity_type = 'collection' AND
      collection_uer.entity_id = collection.id  AND
      collection_uer.user_id = user_media.requester_user_id AND
      collection_uer.level IN ('admin', 'write', 'read')
   )
  UNION
  -- check collections
  SELECT 
    0 AS object_is_present, 
    0 AS object_is_ok,
    0 AS object_is_public,
    COALESCE(collection.id, 0) AS collection_is_present,
    COALESCE(collection_uer.id, 0) AS collection_is_ok,
	0 AS collection_is_public
  FROM user_media AS user_media
  LEFT JOIN collection AS collection
    ON collection.image_media_id = user_media.id
  LEFT JOIN user_entity_right AS collection_uer
    ON (
      collection_uer.entity_type = 'collection' AND
      collection_uer.entity_id = collection.id  AND
      collection_uer.user_id = user_media.requester_user_id AND
      collection_uer.level IN ('admin', 'write', 'read')
   )
  UNION
  -- check collections
  SELECT
    0 AS object_is_present,
    0 AS object_is_ok,
    0 AS object_is_public,
    0 AS collection_is_present,
    0 AS collection_is_ok,
    COALESCE(collection.public::int, 0) AS collection_is_public
  FROM collection AS collection 
  WHERE collection.image_media_id = $2
  UNION 
  -- check objects
  SELECT
    0 AS object_is_present,
    0 AS object_is_ok,
    COALESCE(collection.public::int, 0) AS object_is_public,
    0 AS collection_is_present,
    0 AS collection_is_ok,
    0 AS collection_is_public 
  FROM object_media_ref AS omr
  LEFT JOIN object AS object
    ON object.id = omr.object_id
  LEFT JOIN collection AS collection
    ON collection.id = object.collection_id
  WHERE omr.media_id = $2
) AS total
`
	rows, err := m.p.Query(sql, authToken, mediaID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var (
		objectIsPresent int64
		objectIsOk      int64
		objectIsPublic  int64

		collectionIsPresent int64
		collectionIsOk      int64
		collectionIsPublic  int64
	)
	for rows.Next() {
		if err := rows.Scan(&objectIsPresent, &objectIsOk, &objectIsPublic, &collectionIsPresent, &collectionIsOk, &collectionIsPublic); err != nil {
			return false, err
		}
	}
	if rows.Err() != nil {
		return false, err
	}

	objectIsAllowed := objectIsPresent > 0 && objectIsOk > 0
	collectionIsAllowed := collectionIsPresent > 0 && collectionIsOk > 0
	isPublic := collectionIsPublic > 0 || objectIsPublic > 0

	allowed := isPublic || objectIsAllowed || collectionIsAllowed
	return allowed, nil
}
