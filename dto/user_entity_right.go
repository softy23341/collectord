package dto

import "time"

//go:generate dbgen -type UserEntityRight

// entity types and levels
const (
	RightEntityTypeCollection RightEntityType = "collection"
	RightEntityTypeGroup      RightEntityType = "group"
	RightEntityTypeObject     RightEntityType = "object"

	RightEntityLevelAdmin RightEntityLevel = "admin"
	RightEntityLevelWrite RightEntityLevel = "write"
	RightEntityLevelRead  RightEntityLevel = "read"
	RightEntityLevelNone  RightEntityLevel = "none"
)

var levelComparation = map[RightEntityLevel]int{
	RightEntityLevelAdmin: 100,
	RightEntityLevelWrite: 10,
	RightEntityLevelRead:  1,
	RightEntityLevelNone:  0,
}

// AccessLevelToNum TBD
func AccessLevelToNum(level RightEntityLevel) int {
	num, found := levelComparation[level]
	if !found {
		panic("undefined level")
	}
	return num
}

// NumToAccessLevel TBD
func NumToAccessLevel(num int) RightEntityLevel {
	for key, value := range levelComparation {
		if value == num {
			return key
		}
	}
	panic("undefined access level num")
}

type (
	// RightEntityType TBD
	RightEntityType string
	// RightEntityLevel TBD
	RightEntityLevel string

	// UserEntityRight TBD
	UserEntityRight struct {
		ID         int64           `db:"id"`
		UserID     int64           `db:"user_id"`
		EntityType RightEntityType `db:"entity_type"`
		EntityID   int64           `db:"entity_id"`

		Level RightEntityLevel `db:"level"`

		CreationTime time.Time `db:"creation_time"`
		RootID       int64     `db:"root_id"`
	}

	// UserEntityRightList TBD
	UserEntityRightList []*UserEntityRight

	// ShortUserEntityRight TBD
	ShortUserEntityRight struct {
		UserID     int64           `db:"user_id"`
		EntityType RightEntityType `db:"entity_type"`
		EntityID   int64           `db:"entity_id"`

		Level RightEntityLevel `db:"level"`
	}

	// ShortUserEntityRightList TBD
	ShortUserEntityRightList []*ShortUserEntityRight
)

// ToShort TBD
func (uer *UserEntityRight) ToShort() *ShortUserEntityRight {
	return &ShortUserEntityRight{
		UserID:     uer.UserID,
		EntityType: uer.EntityType,
		EntityID:   uer.EntityID,
		Level:      uer.Level,
	}
}

// UnderLevel TBD
func (ur ShortUserEntityRightList) UnderLevel(lvl RightEntityLevel) ShortUserEntityRightList {
	filtered := make(ShortUserEntityRightList, 0, len(ur))
	for _, right := range ur {
		if right.Level.IsStrongerOrEqual(lvl) {
			continue
		}

		filtered = append(filtered, right)
	}

	return filtered
}

// GetEntityIDToLevel TBD
func (ur ShortUserEntityRightList) GetEntityIDToLevel(typo RightEntityType) map[int64]RightEntityLevel {
	solver := make(map[int64]RightEntityLevel, len(ur))
	for _, right := range ur {
		if right.EntityType != typo {
			continue
		}

		solver[right.EntityID] = right.Level
	}

	return solver
}

// GetUserIDToLevel TBD
func (ur ShortUserEntityRightList) GetUserIDToLevel() map[int64]RightEntityLevel {
	solver := make(map[int64]RightEntityLevel, len(ur))
	for _, right := range ur {
		solver[right.UserID] = right.Level
	}

	return solver
}

// GetEntityIDs TBD
func (ur ShortUserEntityRightList) GetEntityIDs() []int64 {
	ids := make([]int64, 0, len(ur))
	for _, right := range ur {
		ids = append(ids, right.EntityID)
	}

	return ids
}

// GetUsersIDs TBD
func (ur ShortUserEntityRightList) GetUsersIDs() []int64 {
	ids := make([]int64, 0, len(ur))
	for _, right := range ur {
		ids = append(ids, right.UserID)
	}

	return ids
}

// IsStronger TBD
func (lvl RightEntityLevel) IsStronger(other RightEntityLevel) bool {
	return levelComparation[lvl] > levelComparation[other]
}

// IsStrongerOrEqual TBD
func (lvl RightEntityLevel) IsStrongerOrEqual(other RightEntityLevel) bool {
	return levelComparation[lvl] >= levelComparation[other]
}

// IsRightEntityTypeValid TBD
func IsRightEntityTypeValid(typo string) bool {
	switch RightEntityType(typo) {
	case RightEntityTypeCollection:
		return true
	}
	return false
}

// IsRightEntityLevelValid TBD
func IsRightEntityLevelValid(level string) bool {
	switch RightEntityLevel(level) {
	case RightEntityLevelAdmin, RightEntityLevelWrite, RightEntityLevelRead, RightEntityLevelNone:
		return true
	}
	return false
}
