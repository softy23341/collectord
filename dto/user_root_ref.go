package dto

//go:generate dbgen -type UserRootRef

// UserRootRefType TBD
type UserRootRefType int16

// User Root ref types
const (
	UserRootTypeRegular UserRootRefType = 0
	UserRootTypeOwner   UserRootRefType = 30
)

// NewRegularUserRootRef TBD
func NewRegularUserRootRef(userID, rootID int64) *UserRootRef {
	return &UserRootRef{
		UserID: userID,
		RootID: rootID,
		Typo:   UserRootTypeRegular,
	}
}

// NewOwnerUserRootRef TBD
func NewOwnerUserRootRef(userID, rootID int64) *UserRootRef {
	return &UserRootRef{
		UserID: userID,
		RootID: rootID,
		Typo:   UserRootTypeOwner,
	}
}

// UserRootRef TBD
type UserRootRef struct {
	UserID int64           `db:"user_id"`
	RootID int64           `db:"root_id"`
	Typo   UserRootRefType `db:"typo"`
}

// RootUsers TBD
type RootUsers struct {
	RootID   int64
	OwnerID  int64
	UsersIDs []int64
}

// IsOwner TBD
func (u *UserRootRef) IsOwner() bool {
	return u.Typo == UserRootTypeOwner
}

// UserRootRefList TBD
type UserRootRefList []*UserRootRef

// GetRootOwnerID TBD
func (u UserRootRefList) GetRootOwnerID(rootID int64) int64 {
	for _, ref := range u {
		if ref.RootID == rootID && ref.IsOwner() {
			return ref.UserID
		}
	}
	return 0
}

// GetRootOwnerIDs TBD
func (u UserRootRefList) GetRootOwnerIDs() []int64 {
	usersIDs := make([]int64, 0, len(u))
	for _, ref := range u {
		if ref.IsOwner() {
			usersIDs = append(usersIDs, ref.UserID)
		}
	}
	return usersIDs
}

// UsersList TBD
func (u UserRootRefList) UsersList() []int64 {
	m := make(map[int64]struct{})
	for _, ref := range u {
		m[ref.UserID] = struct{}{}
	}
	ids := make([]int64, 0, len(m))
	for k := range m {
		ids = append(ids, k)
	}
	return ids
}

// GetRootIDs TBD
func (u UserRootRefList) GetRootIDs() []int64 {
	m := make(map[int64]struct{})
	for _, ref := range u {
		m[ref.RootID] = struct{}{}
	}
	ids := make([]int64, 0, len(m))
	for k := range m {
		ids = append(ids, k)
	}
	return ids
}

// RootToUsers TBD
func (u UserRootRefList) RootToUsers() map[int64]*RootUsers {
	m := make(map[int64]*RootUsers)
	for _, ref := range u {
		rootUsers, found := m[ref.RootID]
		if !found {
			rootUsers = &RootUsers{
				RootID: ref.RootID,
			}
		}
		if ref.IsOwner() {
			rootUsers.OwnerID = ref.UserID
		} else {
			rootUsers.UsersIDs = append(rootUsers.UsersIDs, ref.UserID)
		}
		m[ref.RootID] = rootUsers
	}
	return m
}
