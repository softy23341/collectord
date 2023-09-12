package dto

import "time"

//go:generate dbgen -type Event

// EventType TBD
type EventType int16

// Event types
const (
	EventTypeNewMessage  EventType = 0
	EventTypeHistoryRead           = 10

	EventTypeNewObject     = 20
	EventTypeMovedObject   = 21
	EventTypeDeletedObject = 22

	EventTypeAddedToGroupCollection     = 30
	EventTypeRemovedFromGroupCollection = 31

	EventTypeNewEntity     = 40
	EventTypeEditedEntity  = 41
	EventTypeDeletedEntity = 42
)

// EventStatus TBD
type EventStatus int16

// Event statuses
const (
	EventStatusNew       EventStatus = 0
	EventStatusConfirmed EventStatus = 10
)

type (
	// Event TBD
	Event struct {
		ID           int64       `db:"id"`
		UserID       int64       `db:"user_id"`
		SeqNo        int64       `db:"seq_no"`
		Type         EventType   `db:"type"`
		Status       EventStatus `db:"status"`
		CreationTime time.Time   `db:"creation_time"`
		EventUnion   `db:"extra,json"`
	}

	// EventUnion TBD
	EventUnion struct {
		NewMessage  *EventNewMessage  `json:"new_message,omitempty"`
		HistoryRead *EventHistoryRead `json:"history_read,omitempty"`

		NewObject     *EventNewObject     `json:"new_object,omitempty"`
		MovedObject   *EventMovedObject   `json:"moved_object,omitempty"`
		DeletedObject *EventDeletedObject `json:"deleted_object,omitempty"`

		AddedToGroupCollection     *EventAddedToGroupCollection     `json:"added_to_group_collection,omitempty"`
		RemovedFromGroupCollection *EventRemovedFromGroupCollection `json:"removed_from_group_collection,omitempty"`

		NewEntity     *EntityRef `json:"new_entity,omitempty"`
		EditedEntity  *EntityRef `json:"edited_entity,omitempty"`
		DeletedEntity *EntityRef `json:"deleted_entity,omitempty"`
	}

	// EventNewMessage TBD
	EventNewMessage struct {
		MessageID int64 `json:"message_id"`
	}

	// EventHistoryRead TBD
	EventHistoryRead struct {
		Peer              Peer  `json:"peer"`
		LastReadMessageID int64 `json:"last_read_message_id"`
	}

	// EventList TBD
	EventList []*Event

	// EventNewObject TBD
	EventNewObject struct {
		ObjectID     int64 `json:"object_id"`
		CollectionID int64 `json:"collection_id"`
	}

	// EventMovedObject TBD
	EventMovedObject struct {
		ObjectID             int64 `json:"object_id"`
		OriginalCollectionID int64 `json:"original_collection_id"`
		NewCollectionID      int64 `json:"new_collection_id"`
	}

	// EventDeletedObject TBD
	EventDeletedObject struct {
		ObjectID             int64 `json:"object_id"`
		OriginalCollectionID int64 `json:"original_collection_id"`
	}

	// EventAddedToGroupCollection TBD
	EventAddedToGroupCollection struct {
		CollectionID int64 `json:"collection_id"`
		GroupID      int64 `json:"group_id"`
	}

	// EventRemovedFromGroupCollection TBD
	EventRemovedFromGroupCollection struct {
		CollectionID int64 `json:"collection_id"`
		GroupID      int64 `json:"group_id"`
	}
)

// TargetUserID TBD
func (e *Event) TargetUserID() int64 {
	return e.UserID
}
