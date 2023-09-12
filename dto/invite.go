package dto

import (
	"time"

	"git.softndit.com/collector/backend/util"
)

//go:generate dbgen -type Invite

// InviteStatus TBD
type InviteStatus int16

// Invite statuses
const (
	InviteCreated  InviteStatus = 0
	InviteAccepted InviteStatus = 10
	InviteRejected InviteStatus = 20
	InviteCanceled InviteStatus = 101
)

type (
	// Invite TBD
	Invite struct {
		ID            int64        `db:"id"`
		CreatorUserID int64        `db:"creator_user_id"`
		RootID        int64        `db:"root_id"`
		ToUserID      *int64       `db:"to_user_id"`
		ToUserEmail   *string      `db:"to_user_email"`
		Token         string       `db:"token"`
		CreationTime  time.Time    `db:"creation_time"`
		Status        InviteStatus `db:"status"`
	}

	// InviteList TBD
	InviteList []*Invite
)

// NewInvite TBD
func NewInvite(rootID int64) (*Invite, error) {
	token, err := util.GenerateInviteToken()
	if err != nil {
		return nil, err
	}

	return &Invite{
		RootID:       rootID,
		Token:        token,
		Status:       InviteCreated,
		CreationTime: time.Now(),
	}, nil
}

// IsCreated TBD
func (i *Invite) IsCreated() bool {
	return i.Status == InviteCreated
}

// IsAccepted TBD
func (i *Invite) IsAccepted() bool {
	return i.Status == InviteAccepted
}

// IsRejected TBD
func (i *Invite) IsRejected() bool {
	return i.Status == InviteRejected
}

// IsCanceled TBD
func (i *Invite) IsCanceled() bool {
	return i.Status == InviteCanceled
}

// FromUser TBD
func (i *Invite) FromUser(userID int64) *Invite {
	i.CreatorUserID = userID
	return i
}

// ToUser TBD
func (i *Invite) ToUser(userID int64) *Invite {
	i.ToUserID = &userID
	return i
}

// UsersIDs TBD
func (i *Invite) UsersIDs() []int64 {
	usersIDs := make([]int64, 0, 2)
	usersIDs = append(usersIDs, i.CreatorUserID)
	if i.ToUserID != nil {
		usersIDs = append(usersIDs, *i.ToUserID)
	}
	return usersIDs
}

// UsersIDs TBD
func (ii InviteList) UsersIDs() []int64 {
	m := make(map[int64]struct{})
	for _, i := range ii {
		for _, userID := range i.UsersIDs() {
			m[userID] = struct{}{}
		}
	}

	ids := make([]int64, 0, len(m))
	for k := range m {
		ids = append(ids, k)
	}
	return ids
}

// InvitedUsersIDs TBD
func (ii InviteList) InvitedUsersIDs() []int64 {
	m := make(map[int64]struct{})
	for _, i := range ii {
		if i.ToUserID != nil {
			m[*i.ToUserID] = struct{}{}
		}
	}

	ids := make([]int64, 0, len(m))
	for k := range m {
		ids = append(ids, k)
	}
	return ids
}

// RootsIDs TBD
func (ii InviteList) RootsIDs() []int64 {
	m := make(map[int64]struct{})
	for _, i := range ii {
		m[i.RootID] = struct{}{}
	}

	ids := make([]int64, 0, len(m))
	for k := range m {
		ids = append(ids, k)
	}
	return ids
}

// IDToInvite TBD
func (ii InviteList) IDToInvite() map[int64]*Invite {
	id2invite := make(map[int64]*Invite, 0)
	for _, invite := range ii {
		id2invite[invite.ID] = invite
	}
	return id2invite
}

// InvitedUserIDToInvite TBD
func (ii InviteList) InvitedUserIDToInvite() map[int64]*Invite {
	id2invite := make(map[int64]*Invite, 0)
	for _, invite := range ii {
		if invite.ToUserID != nil {
			id2invite[*invite.ToUserID] = invite
		}
	}
	return id2invite
}
