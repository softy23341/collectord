package dto

// EntityType TBD
type EntityType int16

// Entity types
const (
	ObjectEntityType     EntityType = 0
	CollectionEntityType            = 10
	GroupEntityType                 = 20

	ActorEntityType          = 30
	BadgeEntityType          = 40
	MaterialEntityType       = 50
	OriginLocationEntityType = 60
	NamedIntervalEntityType  = 70
)

// EntityRef TBD
type EntityRef struct {
	Typo EntityType `json:"typo"`
	ID   int64      `json:"id"`
}
