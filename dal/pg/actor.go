package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetObjectsActorRefs TBD
func (m *Manager) GetObjectsActorRefs(objectsIDs []int64) (dto.ObjectActorRefList, error) {
	flds := dto.ObjectActorRefAllFields

	sql := fmt.Sprintf(`
	        SELECT %s
                FROM object_actor_ref
                WHERE object_id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanObjectActorRefList(m.p, flds, sql, objectsIDs)
}

// GetActorsByIDs TBD
func (m *Manager) GetActorsByIDs(actorsIDs []int64) (dto.ActorList, error) {
	flds := dto.ActorAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM actor
                WHERE id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanActorList(m.p, flds, sql, actorsIDs)
}

// GetOrCreateActorByNormalName TBD
func (m *Manager) GetOrCreateActorByNormalName(inActor *dto.Actor) (*dto.Actor, error) {
	outActors, err := m.GetActorsByNormalNames(
		*inActor.RootID,
		[]string{inActor.NormalName})

	if err != nil {
		return nil, err
	}
	if len(outActors) > 0 {
		return outActors[0], nil
	}
	if err := m.CreateActor(inActor); err != nil {
		return nil, err
	}
	return inActor, nil
}

// GetActorsByNormalNames TBD
func (m *Manager) GetActorsByNormalNames(rootID int64, normalNames []string) (dto.ActorList, error) {
	flds := dto.ActorAllFields
	sql := fmt.Sprintf(`
	  SELECT %s
	  FROM actor
	  WHERE %s = $1 AND %s = any($2::varchar[])
	`, flds.JoinedNames(), dto.ActorFieldRootID.Name(), dto.ActorFieldNormalName.Name())

	return dto.ScanActorList(m.p, flds, sql, rootID, normalNames)
}

// CreateActor TBD
func (m *Manager) CreateActor(actor *dto.Actor) error {
	flds := dto.ActorAllFields.Del(dto.ActorFieldID)

	sql := fmt.Sprintf(`INSERT INTO actor(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		flds.JoinedNames(),
		flds.Placeholders(),
		dto.ActorFieldID.Name(),
	)

	return dto.QueryActorRow(
		m.p,
		dto.ActorFieldsList{dto.ActorFieldID},
		sql,
		(*dto.Actor)(actor).FieldsValues(flds)...,
	).ScanTo(actor)
}

// CreateObjectActorsRefs TBD
func (m *TxManager) CreateObjectActorsRefs(objectID int64, actorIDs []int64) error {
	// TODO make batch insert
	for _, actorID := range actorIDs {
		ref := &dto.ObjectActorRef{
			ObjectID: objectID,
			ActorID:  actorID,
		}
		if err := m.createObjectActorRef(ref); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) createObjectActorRef(ref *dto.ObjectActorRef) error {
	flds := dto.ObjectActorRefAllFields.Del(dto.ObjectActorRefFieldID)

	sql := fmt.Sprintf(`INSERT INTO object_actor_ref(%[1]s) VALUES(%[2]s)`,
		flds.JoinedNames(),
		flds.Placeholders(),
	)

	_, err := m.p.Exec(sql, (*dto.ObjectActorRef)(ref).FieldsValues(flds)...)
	return err
}

// DeleteObjectActorsRefs TBD
func (m *Manager) DeleteObjectActorsRefs(objectID int64) error {
	sql := fmt.Sprintf(`DELETE FROM object_actor_ref WHERE %s = $1`,
		dto.ObjectActorRefFieldObjectID.Name(),
	)

	_, err := m.p.Exec(sql, objectID)
	return err
}

// UpdateActor TBD
func (m *Manager) UpdateActor(actor *dto.Actor) error {
	sql := `
          UPDATE actor
          SET
            name = $1,
            normal_name = $2
          WHERE id = $3
        `

	_, err := m.p.Exec(sql,
		actor.Name,
		actor.NormalName,
		actor.ID,
	)
	return err
}

// GetActorsByRootID TBD
func (m *Manager) GetActorsByRootID(rootID int64) (dto.ActorList, error) {
	flds := dto.ActorAllFields

	sql := fmt.Sprintf(`
	        SELECT %s
                FROM actor
                WHERE root_id = $1
                ORDER BY id DESC
        `, flds.JoinedNames())

	return dto.ScanActorList(m.p, flds, sql, rootID)
}

// DeleteActors TBD
func (m *Manager) DeleteActors(actorsIDs []int64) error {
	sql := `DELETE FROM actor WHERE id = any($1::bigint[])`

	_, err := m.p.Exec(sql, actorsIDs)
	return err
}
