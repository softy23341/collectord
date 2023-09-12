package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/dal"
	"git.softndit.com/collector/backend/dto"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

// AccessRights TBD
type AccessRights struct {
	DBM    dal.Manager
	logger log15.Logger
}

// NewAccessRightsChecker TBD
func NewAccessRightsChecker(DBM dal.Manager, log log15.Logger) *AccessRights {
	return &AccessRights{
		DBM:    DBM,
		logger: log,
	}
}

// HasUserRightsForObjects TBD
func (ar *AccessRights) HasUserRightsForObjects(userID int64, level dto.RightEntityLevel, objectsID []int64) (bool, error) {
	return ar.DBM.HasUserRightsForObjects(userID, level, objectsID)
}

// HasUserRightsForCollections TBD
func (ar *AccessRights) HasUserRightsForCollections(userID int64, level dto.RightEntityLevel, collectionsID []int64) (bool, error) {
	return ar.DBM.HasUserRightsForCollections(userID, level, collectionsID)
}

// HasUserRightsForGroups TBD
func (ar *AccessRights) HasUserRightsForGroups(userID int64, level dto.RightEntityLevel, groupsID []int64) (bool, error) {
	return ar.DBM.HasUserRightsForGroups(userID, level, groupsID)
}

// IsUserRelatedToRoots is a generic light check
func (ar *AccessRights) IsUserRelatedToRoots(userID int64, rootsIDs []int64) (bool, error) {
	rootsList, err := ar.DBM.GetRootsByUserID(userID)
	if err != nil {
		return false, fmt.Errorf("rights check GetRootsByUserID: %v", err)
	}

	targetHits := len(rootsIDs)
	hits := 0
	for _, targetRootID := range rootsIDs {
		for _, root := range rootsList {
			if root.ID == targetRootID {
				hits++
				break
			}
		}
	}

	return targetHits == hits, nil
}

// IsUserInRoots is a strict check
func (ar *AccessRights) IsUserInRoots(userID int64, rootsIDs []int64) (bool, error) {
	for _, rootID := range rootsIDs {
		rootUserList, err := ar.DBM.GetUsersByRootID(rootID)
		if err != nil {
			return false, fmt.Errorf("cant check rights: %v", err)
		}
		for _, user := range rootUserList {
			if user.ID == userID {
				return true, nil
			}
		}
	}

	return false, nil
}

// IsUserRootOwner is a strict check
func (ar *AccessRights) IsUserRootOwner(userID int64, rootID int64) (bool, error) {
	rootsRefList, err := ar.DBM.GetUserRelatedRootRefs([]int64{userID})
	if err != nil {
		ar.logger.Error("GetUserRelatedRootRefs", "err", err)
		return false, nil
	}

	for _, rootRef := range rootsRefList {
		if rootRef.RootID == rootID && rootRef.UserID == userID && rootRef.IsOwner() {
			return true, nil
		}
	}

	return false, nil
}

// IsUserRelatedToTask is a generic light check
func (ar *AccessRights) IsUserRelatedToTask(userID int64, task *dto.Task) (bool, error) {
	for _, taskUserID := range task.GetUsersIDs() {
		if taskUserID == userID {
			return true, nil
		}
	}

	return false, nil
}

// IsUserTaskOwner is a generic light check
func (ar *AccessRights) IsUserTaskOwner(userID int64, task *dto.Task) (bool, error) {
	return task.CreatorUserID == userID, nil
}

// IsUserOwnerOfMedias TBD
func (ar *AccessRights) IsUserOwnerOfMedias(userID int64, medias dto.MediaList) bool {
	//mediaOwners := medias.GetOwnersIDs()
	//for _, mediaOwnerID := range mediaOwners {
	//	if mediaOwnerID != userID {
	//		return false
	//	}
	//}
	return true
}

// WasUserInvited TBD
func (ar *AccessRights) WasUserInvited(userID int64, invite *dto.Invite) bool {
	return *invite.ToUserID == userID
}

// IsUserInviteOwner TBD
func (ar *AccessRights) IsUserInviteOwner(userID int64, invite *dto.Invite) bool {
	return invite.CreatorUserID == userID
}

// IsUserInChat TBD
func (ar *AccessRights) IsUserInChat(userID int64, chatID int64) (bool, error) {
	// get users in chat
	conversationUsers, err := ar.DBM.GetChatConversationUsers(chatID)
	if err != nil {
		ar.logger.Error("GetChatConversationUsers", "err", err)
		return false, err
	}

	for _, conversationUser := range conversationUsers {
		if conversationUser.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}
