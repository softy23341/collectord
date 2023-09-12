package resource

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	sauth "git.softndit.com/collector/backend/auth"
	cleaverClient "git.softndit.com/collector/backend/cleaver/client"
	"git.softndit.com/collector/backend/dal"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/auth"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/inconshreveable/log15.v2"
)

// Registration TBD
func (a *Auth) Registration(params auth.PostAuthRegParams) middleware.Responder {
	logger := sauth.RequestLogger(params.HTTPRequest)
	logger.Debug("Registration")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return auth.NewPostAuthRegDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	alreadyPresentResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return auth.NewPostAuthRegConflict().WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	// check email
	regParams := params.RRegistration
	users, err := a.Context.DBM.GetUsersByEmail([]string{*regParams.Email})
	if err != nil {
		logger.Error("cant check users present GetUsersByEmail", "err", err)
		return errorResponse(500, err.Error())
	}
	if len(users) > 0 {
		logger.Error("user already present", "user email", *regParams.Email)
		return alreadyPresentResponse(409, "user already present")
	}

	// gen password
	//password := strings.ToLower(*regParams.Password)
	password := *regParams.Password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("cant crypt password", "err", err)
		return errorResponse(500, err.Error())
	}

	user := &dto.User{
		Email:             *regParams.Email,
		EncryptedPassword: string(encryptedPassword),
		Tags:              []string{},
	}

	// save user locale
	headerLocale := strings.Split(strings.Split(params.HTTPRequest.Header.Get("Accept-Language"), ",")[0], "-")[0]
	if len(headerLocale) != 0 {
		user.Locale = headerLocale
	}

	// creation
	tx, err := a.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// user
	if err := tx.CreateUser(user); err != nil {
		logger.Error("cant create user", "user", fmt.Sprintf("%+v", user))
		return errorResponse(500, err.Error())
	}

	// root
	root := &dto.Root{}
	if err := tx.CreateRoot(root); err != nil {
		logger.Error("cant create root", "err", err)
		return errorResponse(500, err.Error())
	}

	rootRef := dto.NewOwnerUserRootRef(user.ID, root.ID)
	if err := tx.CreateUserRootRef(rootRef); err != nil {
		logger.Error("cant AttachUserToRoot", "err", err)
		return errorResponse(500, err.Error())
	}

	// trash collection
	//trashCollection := dto.NewTrashCollection(root.ID)
	//if err := tx.CreateCollection(trashCollection); err != nil {
	//	logger.Error("cant create trash collection", "err", err)
	//	return errorResponse(500, err.Error())
	//}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "err", err.Error())
		return errorResponse(500, err.Error())
	}

	//send confirm_email
	if regParams.Invite == nil || len(*regParams.Invite) == 0 {
		if err := a.sendConfirmEmail(user, params.HTTPRequest.Host); err != nil {
			logger.Error("send email err", err.Error())
		}
	}

	// check invite code
	if regParams.Invite != nil && len(*regParams.Invite) != 0 {
		if err := a.checkInvitation(a.Context.DBM, user, *regParams.Invite); err != nil {
			logger.Error("add invite ref err", err.Error())
		}
	}

	return auth.NewPostAuthRegNoContent()
}

func (a *Auth) sendConfirmEmail(user *dto.User, host string) error {
	confirmEmailToken, err := util.GenerateInviteToken()
	if err != nil {
		return err
	}
	if err := a.EmailTokenStorage.Set(confirmEmailToken, user.Email); err != nil {
		return err
	}

	scheme := "https"
	confirmEmailURL := fmt.Sprintf("%v://%v/confirm-email/%v",
		scheme,
		host,
		confirmEmailToken)

	a.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		subject, html := NewConfirmEmailMail(a.Templates, user, confirmEmailURL)
		a.Context.MailClient.Send(&services.Mail{
			To:      []string{user.Email},
			From:    services.SystemMailFrom,
			Subject: subject,
			Body:    html,
		})
	}))

	return nil
}

func (a *Auth) checkInvitation(dbm dal.TrManager, user *dto.User, inviteToken string) error {
	var (
		events dto.EventList
		jobs   []delayedjob.Job
	)

	invite, err := dbm.GetInviteByToken(inviteToken, dto.InviteCreated)
	if err != nil {
		return err
	}
	if invite == nil {
		return nil
	}

	//// email must be the same
	//if invite.ToUserEmail == nil || *invite.ToUserEmail != user.Email {
	//	return nil
	//}

	tx, err := dbm.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	invite.ToUserID = &user.ID
	if err := tx.ChangeInviteToUserID(invite.ID, user.ID); err != nil {
		return err
	}
	if err := tx.ChangeInviteStatus(invite.ID, dto.InviteAccepted); err != nil {
		return err
	}

	systemUser, err := tx.GetSystemUser()
	if err != nil {
		return err
	}

	// get creator user
	var creatorUser *dto.User
	if usersList, err := dbm.GetUsersByIDs([]int64{invite.CreatorUserID}); err != nil {
		return err
	} else if len(usersList) != 1 {
		return errors.New("invalid user id")
	} else {
		creatorUser = usersList[0]
	}

	if len(creatorUser.Email) != 0 {
		sendEmailJob := delayedjob.NewJob(delayedjob.Immideate, func() {
			subject, body := InviteWasAcceptedMail(a.Templates, user, creatorUser)
			a.Context.MailClient.Send(&services.Mail{
				To:      []string{creatorUser.Email},
				From:    services.SystemMailFrom,
				Subject: subject,
				Body:    body,
			})
		})
		jobs = append(jobs, sendEmailJob)
	}

	inviteAcceptedMessage := &dto.Message{
		UserID:     systemUser.ID,
		UserUniqID: util.NextUniqID(),
		PeerID:     invite.CreatorUserID,
		PeerType:   dto.PeerTypeUser,
		Typo:       dto.MessageTypeService,
		MessageExtra: dto.MessageExtra{
			Service: &dto.ServiceMessage{
				Type: dto.ServiceMessageTypeInviteStatusChanged,
				InviteStatusChanged: &dto.ServiceMessageInviteStatusChanged{
					InviteID:     invite.ID,
					InviteStatus: dto.InviteAccepted,
				},
			},
		},
	}

	rootRef := dto.NewRegularUserRootRef(*invite.ToUserID, invite.RootID)
	if err := tx.CreateUserRootRef(rootRef); err != nil {
		return err
	}

	user.EmailVerified = true
	if err := tx.UpdateUser(user); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	a.EmailTokenStorage.Del(inviteToken)

	result, err := a.Context.MessengerClient.
		NewMessageSender(tx, inviteAcceptedMessage, &services.MessageInfo{}).
		Send()

	if err != nil {
		return err
	}

	events = append(events, result.Events...)
	jobs = append(jobs, result.Jobs...)

	a.Context.EventSender.Send(events...)
	a.Context.JobPool.Enqueue(jobs...)

	return nil
}

type userFiller struct {
	DBM           dal.TxManager
	Storage       services.FileStorage
	CleaverClient cleaverClient.ConnectClient
	SearchClient  services.SearchClient
	JobPool       *delayedjob.Pool
	log           log15.Logger

	userID int64
	rootID int64
}

func (u *userFiller) fillUser() error {
	// first collection
	firstRegularCollection := dto.NewRegularCollection(u.rootID)
	firstRegularCollection.Name = "Tables"
	if err := u.DBM.CreateCollection(firstRegularCollection); err != nil {
		u.log.Error("cant create collection", "err", err)
		return err
	}

	// second collection
	secondRegularCollection := dto.NewRegularCollection(u.rootID)
	secondRegularCollection.Name = "Chairs"
	if err := u.DBM.CreateCollection(secondRegularCollection); err != nil {
		u.log.Error("cant create collection 2", "err", err)
		return err
	}

	// group
	nextID := util.NextUniqID()
	newGroup := &dto.Group{
		Name:       "Furniture",
		RootID:     u.rootID,
		UserID:     &u.userID,
		UserUniqID: &nextID,
	}
	if err := u.DBM.CreateGroup(newGroup); err != nil {
		u.log.Error("cant create group", "err", err)
		return err
	}

	// add collections to group
	collectionsIDs := []int64{firstRegularCollection.ID, secondRegularCollection.ID}
	if err := u.DBM.CreateGroupCollectionsRefs(newGroup.ID, collectionsIDs); err != nil {
		u.log.Error("cant create group refs", "err", err)
		return err

	}

	// add objects

	// materials
	materialsNames := []string{"wood", "iron"}
	var (
		materialList dto.MaterialList
		err          error
	)
	materialCreator := &materialManager{
		DBM:    u.DBM,
		log:    u.log,
		rootID: u.rootID,
	}

	materialList, err = materialCreator.createMaterialsByNormalNames(materialsNames)
	if err != nil {
		u.log.Error("cant create materials", "err", err)
		return err
	}

	// actors
	actorsNames := []string{"Tom Smith", "John Joiner"}
	var actorList dto.ActorList
	actorCreator := &actorManager{
		DBM:    u.DBM,
		log:    u.log,
		rootID: u.rootID,
	}

	actorList, err = actorCreator.createActorsByNormalNames(actorsNames)
	if err != nil {
		u.log.Error("cant create actors", "err", err)
		return err
	}

	tablesMedia := dto.MediaList{}
	chairMedia := dto.MediaList{}

	// media
	{
		file, err := os.Open("/collector-uploads/demo/table1.jpg")
		if err != nil {
			u.log.Error("cant open file", "err", err)
			return err
		}

		mediaSaver := &MediaSaver{
			DBM:           u.DBM,
			Logger:        u.log,
			Storage:       u.Storage,
			CleaverClient: u.CleaverClient,
		}
		media, err := mediaSaver.SaveMedia(&InputFileData{
			OriginalFileName: "table1.jpg",
			Typo:             dto.MediaTypePhotoObject,
			UserID:           &u.userID,
			UniqID:           util.NextUniqID(),
			Content:          file,
		})
		if err != nil {
			u.log.Error("cant save media", "err", err)
			return err
		}
		tablesMedia = append(tablesMedia, media)
	}

	// media
	{
		file, err := os.Open("/collector-uploads/demo/table2.jpg")
		if err != nil {
			u.log.Error("cant open file", "err", err)
			return err
		}

		mediaSaver := &MediaSaver{
			DBM:           u.DBM,
			Logger:        u.log,
			Storage:       u.Storage,
			CleaverClient: u.CleaverClient,
		}
		media, err := mediaSaver.SaveMedia(&InputFileData{
			OriginalFileName: "table2.jpg",
			Typo:             dto.MediaTypePhotoObject,
			UserID:           &u.userID,
			UniqID:           util.NextUniqID(),
			Content:          file,
		})
		if err != nil {
			u.log.Error("cant save media", "err", err)
			return err
		}
		tablesMedia = append(tablesMedia, media)
	}

	// media
	{
		file, err := os.Open("/collector-uploads/demo/chair1.jpg")
		if err != nil {
			u.log.Error("cant open file", "err", err)
			return err
		}

		mediaSaver := &MediaSaver{
			DBM:           u.DBM,
			Logger:        u.log,
			Storage:       u.Storage,
			CleaverClient: u.CleaverClient,
		}
		media, err := mediaSaver.SaveMedia(&InputFileData{
			OriginalFileName: "chair1.jpg",
			Typo:             dto.MediaTypePhotoObject,
			UserID:           &u.userID,
			UniqID:           util.NextUniqID(),
			Content:          file,
		})
		if err != nil {
			u.log.Error("cant save media", "err", err)
			return err
		}
		chairMedia = append(chairMedia, media)
	}

	// Object creation
	{
		// table 1

		object := dto.NewObject()
		object.Name = "Solid Table"
		object.UserID = u.userID
		object.UserUniqID = util.NextUniqID()
		object.CollectionID = firstRegularCollection.ID
		object.Description = "Solid table"

		err = u.DBM.CreateObject(object)
		if err != nil {
			u.log.Error("create object", "error", err)
			return err
		}

		version, err := u.DBM.GetCurrentTransaction()
		if err != nil {
			u.log.Error("GetCurrentTransaction", "error", err)
			return err
		}

		forIndex := object.SearchObject().
			WithRootID(u.rootID).
			SetExternalVersion(version)

		if err := u.DBM.CreateObjectMediasRefs(object.ID, []int64{tablesMedia[0].ID}); err != nil {
			u.log.Error("CreateObjectMediasRefs", "error", err)
			return err
		}

		forIndex.WithMaterialsIDs(materialList.IDs())
		if err := u.DBM.CreateObjectMaterialsRefs(object.ID, materialList.IDs()); err != nil {
			u.log.Error("CreateObjectMaterialsRefs", "error", err)
			return err
		}

		forIndex.WithActorsIDs([]int64{actorList[0].ID})
		if err := u.DBM.CreateObjectActorsRefs(object.ID, []int64{actorList[0].ID}); err != nil {
			u.log.Error("CreateObjectActorsRefs", "error", err, "objectID", object.ID, "actors", []int64{actorList[0].ID})
			return err
		}

		u.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
			err := u.SearchClient.IndexObject(forIndex)
			if err != nil {
				u.log.Error("insert into search", "err", err)
			}
		}))
	}

	// Object creation
	{
		// table 2

		object := dto.NewObject()
		object.Name = "Old Table"
		object.UserID = u.userID
		object.UserUniqID = util.NextUniqID()
		object.CollectionID = firstRegularCollection.ID
		object.Description = "Old table"

		err = u.DBM.CreateObject(object)
		if err != nil {
			u.log.Error("create object", "error", err)
			return err
		}

		version, err := u.DBM.GetCurrentTransaction()
		if err != nil {
			u.log.Error("GetCurrentTransaction", "error", err)
			return err
		}

		forIndex := object.SearchObject().
			WithRootID(u.rootID).
			SetExternalVersion(version)

		if err := u.DBM.CreateObjectMediasRefs(object.ID, []int64{tablesMedia[1].ID}); err != nil {
			u.log.Error("CreateObjectMediasRefs", "error", err)
			return err
		}

		forIndex.WithMaterialsIDs([]int64{materialList[1].ID})
		if err := u.DBM.CreateObjectMaterialsRefs(object.ID, []int64{materialList[1].ID}); err != nil {
			u.log.Error("CreateObjectMaterialsRefs", "error", err)
			return err
		}

		forIndex.WithActorsIDs(actorList.IDs())
		if err := u.DBM.CreateObjectActorsRefs(object.ID, actorList.IDs()); err != nil {
			u.log.Error("CreateObjectActorsRefs", "error", err, "objectID", object.ID, "actors", []int64{actorList[0].ID})
			return err
		}

		u.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
			err := u.SearchClient.IndexObject(forIndex)
			if err != nil {
				u.log.Error("insert into search", "err", err)
			}
		}))
	}

	// Object creation
	{
		// chair1

		object := dto.NewObject()
		object.Name = "Nice Chair"
		object.UserID = u.userID
		object.UserUniqID = util.NextUniqID()
		object.CollectionID = secondRegularCollection.ID
		object.Description = "Nice Chair"

		err = u.DBM.CreateObject(object)
		if err != nil {
			u.log.Error("create object", "error", err)
			return err
		}

		version, err := u.DBM.GetCurrentTransaction()
		if err != nil {
			u.log.Error("GetCurrentTransaction", "error", err)
			return err
		}

		forIndex := object.SearchObject().
			WithRootID(u.rootID).
			SetExternalVersion(version)

		if err := u.DBM.CreateObjectMediasRefs(object.ID, []int64{chairMedia[0].ID}); err != nil {
			u.log.Error("CreateObjectMediasRefs", "error", err)
			return err
		}

		forIndex.WithMaterialsIDs(materialList.IDs())
		if err := u.DBM.CreateObjectMaterialsRefs(object.ID, materialList.IDs()); err != nil {
			u.log.Error("CreateObjectMaterialsRefs", "error", err)
			return err
		}

		forIndex.WithActorsIDs(actorList.IDs())
		if err := u.DBM.CreateObjectActorsRefs(object.ID, actorList.IDs()); err != nil {
			u.log.Error("CreateObjectActorsRefs", "error", err, "objectID", object.ID, "actors", []int64{actorList[0].ID})
			return err
		}

		u.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
			err := u.SearchClient.IndexObject(forIndex)
			if err != nil {
				u.log.Error("insert into search", "err", err)
			}
		}))
	}

	return nil
}

func cpreferences() {
	// cp references
	// materials
	//if err := tx.CopyMaterialsToRoot(root.ID); err != nil {
	//	logger.Error("cant copy materials to root", "err", err)
	//	return errorResponse(500, err.Error())
	//}

	// origin location
	//if err := tx.CopyOriginLocationToRoot(root.ID); err != nil {
	//	logger.Error("cant copy origin location to root", "err", err)
	//	return errorResponse(500, err.Error())
	//}

	// named date intervals
	//if err := tx.CopyNamedDateIntervalsToRoot(root.ID); err != nil {
	//	logger.Error("cant copy named date intervals", "err", err)
	//	return errorResponse(500, err.Error())
	//}

	// YEAH commented code; There is not place for drafts in this app
	// draft
	// draftCollection := dto.NewDraftCollection(root.ID)
	// if err := tx.CreateCollection(draftCollection); err != nil {
	// 	logger.Error("cant create collection", "err", err)
	// 	return errorResponse(500, err.Error())
	// }

	// fillUser with test data
	/*
		filler := &userFiller{
			DBM:           tx,
			log:           logger,
			Storage:       a.Context.FileStorage, // XXX TODO
			CleaverClient: a.Context.CleaverClient,
			SearchClient:  a.Context.SearchClient,
			JobPool:       a.Context.JobPool,

			rootID: root.ID,
			userID: user.ID,
		}
		if err := filler.fillUser(); err != nil {
			logger.Error("cant fill user", "err", err)
			return errorResponse(500, err.Error())
		}
	*/

	// session auth token
	/*
		token, err := util.GenerateToken()
		if err != nil {
			logger.Error("cant create user", "user", fmt.Sprintf("%+v", user))
			return errorResponse(500, err.Error())
		}
		session := &dto.UserSession{
			UserID:    user.ID,
			AuthToken: token,
		}
		if err := tx.CreateUserSession(session); err != nil {
			logger.Error("error", "err", err.Error())
			return errorResponse(500, err.Error())
		}
	*/

	//return newCookieAuthRespond(auth.NewPostAuthRegOK().WithPayload(&models.ARegistration{
	//	AuthToken: &models.AuthToken{
	//		Token: &token,
	//	},
	//}), token)

}
