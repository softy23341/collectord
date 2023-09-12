package dal

import (
	"time"

	"git.softndit.com/collector/backend/dto"
)

// Tx TBD
type Tx interface {
	Commit() error
	Rollback() error
}

// TxManager TBD
type TxManager interface {
	Tx
	Manager
	// Material
	CreateObjectMaterialsRefs(objectID int64, materialIDs []int64) error
	// Media
	CreateObjectMediasRefs(objectID int64, mediaIDs []int64) error
	// actors
	CreateObjectActorsRefs(objectID int64, actorIDs []int64) error
	// origin location
	CreateObjectOriginLocationsRefs(objectID int64, originLocationIDs []int64) error
	// badge
	CreateObjectBadgesRefs(objectID int64, badgeIDs []int64) error

	// collections
	GetCollectionsByIDsForUpdate(collectionsIDs []int64) (dto.CollectionList, error)

	// groups
	CreateGroupCollectionsRefs(groupID int64, collectionsIDs []int64) error
	CreateCollectionGroupsRefs(collectionID int64, groupsIDs []int64) error

	// tasks
	CreateTaskObjectsRefs(taskID int64, objectsIDs []int64) error
	CreateTaskCollectionsRefs(taskID int64, collectionsIDs []int64) error
	CreateTaskGroupsRefs(taskID int64, groupsIDs []int64) error
	CreateTaskMediasRefs(taskID int64, mediasIDs []int64) error

	// user rights
	PutUserRight(right *dto.UserEntityRight) error
}

// Transactionable TBD
type Transactionable interface {
	BeginTx() (TxManager, error)
}

// TrManager TBD
type TrManager interface {
	Transactionable
	Manager
}

// Manager TBD
type Manager interface {
	// Object
	GetObjectsByIDs(objectIDs []int64) (dto.ObjectList, error)
	GetObjectsPreviewByIDs(objectIDs []int64) (dto.ObjectPreviewList, error)
	GetCollectionsObjectPreviews(collectionsIDs []int64, orders *dto.ObjectOrders, paginator *dto.PagePaginator) (dto.ObjectPreviewList, error)
	GetObjectIDByUserUniqID(userID, uniqID int64) (*int64, error)
	GetObjects(paginator *dto.PagePaginator) (dto.ObjectList, error)

	GetObjectsByUserIDsWithCustomFields(userIDs []int64, flds dto.ObjectFieldsList) (dto.ObjectList, error)

	CreateObject(object *dto.Object) error

	UpdateObject(object *dto.Object) error
	ChangeObjectsCollection(fromID, toID int64) error
	ChangeObjectsCollectionByIDs(objectsIDs []int64, collectionID int64) error

	DeleteObjectsByIDs(objectIDs []int64) error
	DeleteObjectsByCollectionsIDs(collectionsIDs []int64) error

	// media
	GetMediaByUserUniqID(userID, uniqID int64) (*dto.Media, error)
	GetMediasByIDs(mediasIDs []int64) (dto.MediaList, error)
	GetObjectsMediaRefs(objectsIDs []int64) (dto.ObjectMediaRefList, error)
	GetObjectsMediaRefsPhoto(objectsIDs []int64) (dto.ObjectMediaRefList, error)
	GetMediaByPage(types []dto.MediaType, p *dto.PagePaginator) (dto.MediaList, error)
	IsMediaUsed(mediaID int64) (bool, error)
	GetMediaByIDAndVariantURI(mediaID int64, variantURI string) (*dto.Media, error)
	CanUserGetMediaByAuthToken(authToken string, mediaID int64) (bool, error)

	CreateMedia(*dto.Media) error
	UpdateMedia(media *dto.Media) error

	DeleteObjectMediasRefs(objectID int64) error
	DeleteMedias(mediasIDs []int64) error

	// materials
	GetMaterialsByIDs(objectIDs []int64) (dto.MaterialList, error)
	GetMaterialsByNormalNames(rootID int64, normalNames []string) (dto.MaterialList, error)
	GetMaterialsByObjectID(objectID int64) (dto.MaterialList, error)
	GetMaterialRefsByObjectsIDs(objectIDs []int64) (dto.ObjectMaterialRefList, error)
	GetMaterialsByRootID(rootID int64) (dto.MaterialList, error)

	GetOrCreateMaterialByNormalName(mat *dto.Material) (*dto.Material, error)

	CreateMaterial(material *dto.Material) error
	CopyMaterialsToRoot(rootID int64) error
	UpdateMaterial(material *dto.Material) error

	DeleteObjectMaterialsRefs(objectID int64) error
	DeleteMaterials(materialsIDs []int64) error

	// actors
	GetActorsByIDs(actorsIDs []int64) (dto.ActorList, error)
	GetObjectsActorRefs(objectsIDs []int64) (dto.ObjectActorRefList, error)
	GetActorsByNormalNames(rootID int64, normalNames []string) (dto.ActorList, error)
	GetActorsByRootID(rootID int64) (dto.ActorList, error)

	GetOrCreateActorByNormalName(mat *dto.Actor) (*dto.Actor, error)

	CreateActor(actor *dto.Actor) error
	UpdateActor(actor *dto.Actor) error

	DeleteObjectActorsRefs(objectID int64) error
	DeleteActors(actorsIDs []int64) error

	// origin location
	GetOriginLocationsByIDs(originLocationsIDs []int64) (dto.OriginLocationList, error)
	GetObjectsOriginLocationRefs(objectsIDs []int64) (dto.ObjectOriginLocationRefList, error)
	GetOriginLocationsByNormalNames(rootID int64, normalNames []string) (dto.OriginLocationList, error)
	GetOriginLocationsByRootID(rootID int64) (dto.OriginLocationList, error)

	GetOrCreateOriginLocationByNormalName(mat *dto.OriginLocation) (*dto.OriginLocation, error)

	CreateOriginLocation(originLocation *dto.OriginLocation) error
	CopyOriginLocationToRoot(rootID int64) error
	UpdateOriginLocation(originLocation *dto.OriginLocation) error

	DeleteObjectOriginLocationsRefs(objectID int64) error
	DeleteOriginLocations(originLocationsIDs []int64) error

	// named date intervals
	GetNamedDateIntervalsByIDs(ids []int64) (dto.NamedDateIntervalList, error)
	GetNamedDateIntervalsForRoots(rootIDs []int64) (dto.NamedDateIntervalList, error)
	GetNamedDayeIntervalsByNormalNames(rootID int64, names []string) (dto.NamedDateIntervalList, error)

	CreateNamedDateInterval(interval *dto.NamedDateInterval) error
	CopyNamedDateIntervalsToRoot(rootID int64) error
	UpdateNamedDateInterval(interval *dto.NamedDateInterval) error

	DeleteNamedDateIntervalsByIDs(ids []int64) error

	// badges
	GetBadgesByIDs(badgesIDs []int64) (dto.BadgeList, error)
	GetObjectsBadgeRefs(objectsIDs []int64) (dto.ObjectBadgeRefList, error)
	GetBadgesByRootID(rootID int64) (dto.BadgeList, error)
	GetBadgesByNormalNamesOrColors(rootID int64, normalNames, colors []string) (dto.BadgeList, error)

	GetOrCreateBadgeByNormalNameAndColor(mat *dto.Badge) (*dto.Badge, error)

	CreateBadge(badge *dto.Badge) error
	UpdateBadge(badge *dto.Badge) error

	DeleteObjectBadgesRefs(objectID int64) error
	DeleteBadges(badgesIDs []int64) error

	// collections
	GetTypedCollection(rootID int64, typo dto.CollectionTypo) (*dto.Collection, error)
	GetCollectionsByIDs(collectionsIDs []int64) (dto.CollectionList, error)
	GetCollectionIDByUserUniqID(userID, uniqID int64) (*int64, error)
	GetCollectionsByRootID(rootID int64) (dto.CollectionList, error)
	GetCollectionsByGroupIDs(groupIDs []int64) (dto.CollectionList, error)
	GetCollectionsByObjectIDs(objectIDs []int64) (dto.CollectionList, error)
	GetObjectsCnt(collectionsIDs []int64) (int64, error)
	GetObjectsCntByCollections(collectionsIDs []int64) (map[int64]int64, error)
	GetPublicCollections(query *string, rootID *int64, paginator *dto.PagePaginator) (int64, dto.CollectionList, error)

	CreateCollection(*dto.Collection) error

	UpdateCollection(collection *dto.Collection) error

	DeleteCollections(collectionsIDs []int64) error
	DeleteCollectionGroupsRefs(collectionID int64) error
	DeleteCollectionsGroupRefsByGroupAndCollections(collectionsIDs []int64, groupID int64) error

	GetCollectionsByUserIDsWithCustomFields(userIDs []int64, flds dto.CollectionFieldsList) (dto.CollectionList, error)
	GetCollectionsByIDsWithCustomFields(collectionsIDs []int64, flds dto.CollectionFieldsList) (dto.CollectionList, error)

	// currencies
	GetCurrenciesByIDs(currencyIDs []int64) (dto.CurrencyList, error)
	GetCurrencies() (dto.CurrencyList, error)

	// groups
	GetGroupsByIDs(groupsIDs []int64) (dto.GroupList, error)
	GetGroupByRootIDAndName(rootID int64, name string) (*dto.Group, error)
	GetGroupsByRootID(rootID int64) (dto.GroupList, error)
	GetCollectionsGroupRefs(collectionsIDs []int64) (dto.CollectionGroupRefList, error)
	GetCollectionsGroupRefsByGroups(groupsIDs []int64) (dto.CollectionGroupRefList, error)

	CreateGroup(group *dto.Group) error
	UpdateGroup(group *dto.Group) error

	DeleteCollectionsGroupRefs(groupID int64) error
	DeleteGroups(groupsIDs []int64) error

	// object statuses
	GetObjectStatuses() (dto.ObjectStatusList, error)
	GetCurrentObjectsStatusesRefs(objectsIDs []int64) (dto.ObjectStatusRefList, error)
	GetObjectStatusByIDs(IDs []int64) (dto.ObjectStatusList, error)

	CreateObjectStatus(objectStatus *dto.ObjectStatus) error
	CreateObjectStatusRef(ref *dto.ObjectStatusRef) error

	// valuations
	GetObjectValuations(objectsIDs []int64) (dto.ValuationList, error)
	CreateValuation(valuation *dto.Valuation) (*dto.Valuation, error)
	DeleteValuationsByObjectID(objectID int64) error
	GetValuationsByCollectionIDs(collectionsIDs []int64) (map[int64]int64, error)

	// user
	GetUsersByEmail(emails []string) (dto.UserList, error)
	GetUsersByIDs(userIDs []int64) (dto.UserList, error)
	GetUsersByRootID(rootID int64) (dto.UserList, error)
	GetRootOwner(rootID int64) (*dto.User, error)
	GetSystemUser() (*dto.User, error)
	GetPopularUserTags(max int) ([]string, error)

	SearchUsersByName(udi int64, name string, page, perPage int16) (int64, dto.UserList, error)

	GetUserByEmailWithCustomFields(email string, flds dto.UserFieldsList) (*dto.User, error)

	CreateUser(user *dto.User) error
	UpdateUser(user *dto.User) error

	// user counters
	UpdateUserLastEventSeqNo(userID int64, delta int32) (seqNo *int64, err error)
	UpdateUserNUnreadMessages(userID int64, delta int32) (n *int32, err error)
	UpdateUserNUnreadNotifications(userID int64, delta int32) (n *int32, err error)

	// Root
	GetRootsByUserID(userID int64) (dto.RootList, error)
	GetRootsByIDs(rootIDs []int64) (dto.RootList, error)
	GetMainUserRoot(userID int64) (dto.RootList, error)
	CreateRoot(root *dto.Root) error

	// Invite
	GetInvitesByIDs(inviteIDs []int64) (dto.InviteList, error)
	GetInviteByToken(token string, status dto.InviteStatus) (*dto.Invite, error)
	GetInvitesByRoot(rootID int64, status dto.InviteStatus) (dto.InviteList, error)
	GetInviteByUserRoot(fromUserID, toUserID, rootID int64, status dto.InviteStatus) (*dto.Invite, error)

	CreateInvite(invite *dto.Invite) error
	ChangeInviteStatus(inviteID int64, status dto.InviteStatus) error
	ChangeInviteToUserID(inviteID int64, toUserID int64) error

	// user root ref
	GetUserRelatedRootRefs(userIDs []int64) (dto.UserRootRefList, error)
	GetUserRootRefs(rootsIDs []int64) (dto.UserRootRefList, error)

	CreateUserRootRef(ref *dto.UserRootRef) error

	DeleteUserRootRef(userID, rootID int64) error

	// events
	GetEvents(userID, beginSeqNo int64, endSeqNo *int64, statuses []dto.EventStatus) (dto.EventList, error)
	GetOldEventsCnt(ago time.Duration) (int32, error)

	SaveEvent(event *dto.Event) error
	ConfirmEvents(userID int64, endSeqNo int64) error

	DeleteOldEvents(ago time.Duration) error

	// messages
	GetMessageByUserUniqID(userID, uniqID int64) (*dto.Message, error)
	GetMessagesByIDs(ids []int64) (dto.MessageList, error)
	GetDialogMessagesByRange(firstUserID, secondUserID int64, p *dto.RangePaginator) (dto.MessageList, error)
	GetChatMessagesByRange(chatID int64, p *dto.RangePaginator) (dto.MessageList, error)

	CreateMessage(m *dto.Message) error

	// unread messages
	SaveUnreadMessage(unreadMessage *dto.UnreadMessage) error

	DelUnreadMessages(userID int64, peer *dto.Peer, maxMessageID int64) (int32, int64, error)
	DelAllUnreadMessagesByPeer(userID int64, peer *dto.Peer) (int32, error)

	// dialogs
	GetDialog(firstUserID, secondUserID int64) (*dto.Dialog, error)

	CreateDialog(dialog *dto.Dialog) error

	// conversation user
	GetChatConversationUsers(chatID int64) (dto.ConversationUserList, error)
	GetDialogConversationUsers(firstUserID, secondUserID int64) (dto.ConversationUserList, error)
	GetUserPeers(userID int64) (dto.PeerList, error)

	CreateConversationUser(user *dto.ConversationUser) error

	UpdateConversationUserNUnreadMessages(userID int64, peer *dto.Peer, delta int32) (*int32, error)
	UpdateConversationUserLastReadMessageID(userID int64, peer *dto.Peer, lastReadMessageID int64) (*int64, error)

	DelConversationUser(userID int64, peer dto.Peer) error

	// chat
	GetChatByUserUniqID(creatorUserID int64, uniqID int64) (*dto.Chat, error)
	GetChatsByIDs(IDs []int64) (dto.ChatList, error)
	GetChatMembersCntByChats(chatsIDs []int64) (map[int64]int64, error)
	GetChatLastReadMessageID(chatID int64) (*int64, error)
	GetChatMessagesAuthors(chatID, startMessageID, endMessageID int64) ([]int64, error)

	CreateChat(chat *dto.Chat) error
	UpdateChat(chat *dto.Chat) error
	SetChatLastReadMessageIDIfGreater(chatID, lastReadMessageID int64) (realLastReadMessageID *int64, err error)

	// auth
	CreateUserSession(userSession *dto.UserSession) error
	GetUsersBySessionTokens(tokens []string) (dto.UserList, error)
	GetSessionsBySessionTokens(tokens []string) (dto.UserSessionList, error)
	DestroyUserSessionByTokens(tokens []string) error

	// session device
	DeleteSessionDevices(sessionToken []string) error
	SetSessionDeviceByAuthToken(*dto.UserSessionDevice) error
	GetUserActiveDevices(userID int64) (dto.UserSessionDeviceList, error)

	// Transaction
	GetCurrentTransaction() (int64, error)

	// task
	GetTasksByIDs(IDs []int64) (dto.TaskList, error)
	GetUserRelatedTasks(userID int64) (dto.TaskList, error)
	GetAssignedToUserTasksFrom(creatorID int64, userID int64) (dto.TaskList, error)
	GetAssignedToUserTasks(userID int64) (dto.TaskList, error)
	GetUserRelatedTasksArchive(userID int64, paginator *dto.PagePaginator) (dto.TaskList, error)

	GetTaskObjects(taskID int64) (dto.ObjectList, error)
	GetTaskCollections(taskID int64) (dto.CollectionList, error)
	GetTaskGroups(taskID int64) (dto.GroupList, error)
	GetTaskMedias(taskID int64) (dto.MediaList, error)

	CreateTask(task *dto.Task) error
	UpdateTask(taskID int64, editedTask *dto.Task, changedFields dto.TaskFieldsList) (*dto.Task, error)

	DeleteTask(taskID int64) error

	// task relations
	GetMediasCntByTaskID(tasksIDs []int64) (map[int64]int64, error)
	GetTasksCntByUserID(creatorID int64, usersIDs []int64) (map[int64]int64, error)

	DeleteTaskObjectsRefs(taskID int64) error
	DeleteTaskCollectionsRefs(taskID int64) error
	DeleteTaskGroupsRefs(taskID int64) error
	DeleteTaskMediasRefs(taskID int64) error

	// user rights
	HasUserRightsForObjects(userID int64, level dto.RightEntityLevel, objectsID []int64) (bool, error)
	HasUserRightsForCollections(userID int64, level dto.RightEntityLevel, collectionsID []int64) (bool, error)
	HasUserRightsForGroups(userID int64, level dto.RightEntityLevel, groupsID []int64) (bool, error)

	GetUserRightsForObjects(userID int64, objectsID []int64) (dto.ShortUserEntityRightList, error)

	GetUserRightsForCollections(userID int64, collectionsID []int64) (dto.ShortUserEntityRightList, error)
	GetUsersRightsForCollection(collection int64) (dto.ShortUserEntityRightList, error)
	GetUserRightsForCollectionsInRoot(userID int64, rootID int64) (dto.ShortUserEntityRightList, error)
	GetUserRightsForGroups(userID int64, groupsID []int64) (dto.ShortUserEntityRightList, error)

	DeleteEntityRights(entityType dto.RightEntityType, entityID int64) error

	CheckUserTmpAccessToObject(userID, messageID, objectID int64) (bool, error)
	CheckUserTmpAccessToMedia(authToken string, messageID, mediaID int64) (bool, error)

	// ban users
	CreateUserBan(ban *dto.UserBan) error
	DeleteUserBan(creatorUserID, userID int64) error
	GetUserBanList(userID int64) (dto.UserBanList, error)
	IsUserBanned(creatorUserID, userID int64) (bool, error)
}
