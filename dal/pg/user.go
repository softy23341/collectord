package dalpg

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx"

	"git.softndit.com/collector/backend/dto"
)

// GetUsersByEmail TBD
func (m *Manager) GetUsersByEmail(emails []string) (dto.UserList, error) {
	flds := dto.UserAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM "user"
                WHERE email = any($1::character varying[])
        `, flds.JoinedNames())

	return dto.ScanUserList(m.p, flds, sql, emails)
}

// GetUsersByIDs TBD
func (m *Manager) GetUsersByIDs(userIDs []int64) (dto.UserList, error) {
	flds := dto.UserAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM "user"
                WHERE id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanUserList(m.p, flds, sql, userIDs)
}

// CreateUser TBD
func (m *Manager) CreateUser(user *dto.User) error {
	flds := dto.UserAllFields.Del(dto.UserFieldID)

	sql := fmt.Sprintf(`INSERT INTO "user"(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		flds.JoinedNames(),
		flds.Placeholders(),
		dto.UserFieldID.Name(),
	)

	return dto.QueryUserRow(
		m.p,
		dto.UserFieldsList{dto.UserFieldID},
		sql,
		(*dto.User)(user).FieldsValues(flds)...,
	).ScanTo(user)
}

// UpdateUser TBD
func (m *Manager) UpdateUser(user *dto.User) error {
	sql := `
	  UPDATE "user"
          SET
            first_name = $1
          , last_name = $2
          , email = $3
          , avatar_media_id = $4
          , description = $5
          , encrypted_password = $6
          , tags = $7
          , email_verified = $8
          , is_anonymous = $9
          , speciality = $10
          WHERE id = $11
        `

	_, err := m.p.Exec(sql,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AvatarMediaID,
		user.Description,
		user.EncryptedPassword,
		user.Tags,
		user.EmailVerified,
		user.IsAnonymous,
		user.Speciality,
		user.ID,
	)
	return err
}

// GetUsersByRootID TBD
func (m *Manager) GetUsersByRootID(rootID int64) (dto.UserList, error) {
	flds := dto.UserAllFields
	sql := fmt.Sprintf(`
           SELECT %s FROM "user" AS u
           INNER JOIN user_root_ref AS ur
             ON ur.user_id = u.id
           WHERE ur.root_id = $1
        `, flds.JoinedNamesWithAlias("u"))

	return dto.ScanUserList(m.p, flds, sql, rootID)
}

// SearchUsersByName TBD
func (m *Manager) SearchUsersByName(udi int64, name string, page, perPage int16) (int64, dto.UserList, error) {
	flds := dto.UserAllFields

	lowerName := strings.ToLower(name)
	nameParts := strings.Split(lowerName, " ")
	var args []interface{}

	queryParts := make([]string, 0, len(nameParts))
	totalPlaceholders := 0
	for i, namePart := range nameParts {
		totalPlaceholders++
		queryParts = append(queryParts, fmt.Sprintf(" user_search_name(u) LIKE $%d ", i+1))
		args = append(args, fmt.Sprintf("%%  %s%%", namePart))
	}

	// except requester user;
	totalPlaceholders++
	queryParts = append(queryParts, fmt.Sprintf(" u.id != $%d ", totalPlaceholders))
	args = append(args, udi)
	// ---

	queryPart := strings.Join(queryParts, "AND")

	// do not show system users;
	queryPart += " AND u.system_user IS false"

	// do not show anonymous users;
	queryPart += " AND u.is_anonymous IS false"

	var cnt int64
	err := m.p.QueryRow(fmt.Sprintf(`SELECT COUNT(u) FROM "user" AS u WHERE %s`, queryPart), args...).Scan(&cnt)
	if err != nil {
		return 0, nil, err
	}

	sql := fmt.Sprintf(`
          SELECT %s
          FROM "user" AS u
          WHERE %s
          ORDER BY u.first_name, u.email
          OFFSET $%d LIMIT $%d;
        `, flds.JoinedNamesWithAlias("u"), queryPart, totalPlaceholders+1, totalPlaceholders+2)

	args = append(args, page*perPage, perPage)

	list, err := dto.ScanUserList(m.p, flds, sql, args...)
	if err != nil {
		return 0, nil, err
	}

	return cnt, list, nil
}

// UpdateUserLastEventSeqNo TBD
func (m *Manager) UpdateUserLastEventSeqNo(userID int64, delta int32) (seqNo *int64, err error) {
	sql := `UPDATE "user"
                SET last_event_seq_no = last_event_seq_no + $1
                WHERE id = $2
                RETURNING last_event_seq_no`
	var value int64
	if err := m.p.QueryRow(sql, delta, userID).Scan(&value); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &value, nil
}

// UpdateUserNUnreadMessages TBD
func (m *Manager) UpdateUserNUnreadMessages(userID int64, delta int32) (*int32, error) {
	sql := `UPDATE "user"
                SET n_unread_messages = n_unread_messages + $1
                WHERE id = $2
                RETURNING n_unread_messages`
	var value int32
	if err := m.p.QueryRow(sql, delta, userID).Scan(&value); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &value, nil
}

// UpdateUserNUnreadNotifications TBD
func (m *Manager) UpdateUserNUnreadNotifications(userID int64, delta int32) (*int32, error) {
	sql := `UPDATE "user"
                SET n_unread_notifications = n_unread_notifications + $1
                WHERE id = $2
                RETURNING n_unread_notifications`
	var value int32
	if err := m.p.QueryRow(sql, delta, userID).Scan(&value); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &value, nil
}

// GetSystemUser TBD
func (m *Manager) GetSystemUser() (*dto.User, error) {
	flds := dto.UserAllFields
	sql := fmt.Sprintf(`
SELECT %s
FROM "user"
WHERE system_user
        `, flds.JoinedNames())

	return dto.ScanUser(m.p, flds, sql)
}

// GetRootOwner TBD
func (m *Manager) GetRootOwner(rootID int64) (*dto.User, error) {
	flds := dto.UserAllFields
	sql := fmt.Sprintf(`
SELECT %s
FROM "user" AS u
INNER JOIN user_root_ref AS urr
  ON urr.user_id = u.id
WHERE urr.root_id = $1 AND urr.typo = $2

        `, flds.JoinedNamesWithAlias("u"))

	return dto.ScanUser(m.p, flds, sql, rootID, int16(dto.UserRootTypeOwner))
}

// GetPopularUserTags TBD
func (m *Manager) GetPopularUserTags(max int) ([]string, error) {
	sql := `
SELECT tags.name AS name, COUNT(*) AS cnt
FROM (
  SELECT unnest(tags) AS name
  FROM "user"
) AS tags
GROUP BY tags.name
ORDER BY cnt DESC
LIMIT $1;
`

	rows, err := m.p.Query(sql, max)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		name string
		cnt  int64
	)

	tags := make([]string, 0, max)
	for rows.Next() {
		if err := rows.Scan(&name, &cnt); err != nil {
			return nil, err
		}
		tags = append(tags, name)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return tags, nil
}

// GetUserByEmailWithCustomFields TBD
func (m *Manager) GetUserByEmailWithCustomFields(email string, flds dto.UserFieldsList) (*dto.User, error) {
	sql := fmt.Sprintf(`SELECT %s FROM "user" WHERE email = $1 Limit 1`, flds.JoinedNames())
	return dto.ScanUser(m.p, flds, sql, email)
}
