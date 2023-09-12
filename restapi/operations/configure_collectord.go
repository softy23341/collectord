package restapi

import (
	"bytes"
	"context"
	"crypto/tls"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.softndit.com/collector/backend/auth/storage"
	resource "git.softndit.com/collector/backend/resources"

	goruntime "runtime"

	cleaver "git.softndit.com/collector/backend/cleaver/client"
	"git.softndit.com/collector/backend/config"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/models"
	npusherclient "git.softndit.com/collector/backend/npusher/client"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/jackc/pgx/log/log15adapter"
	"github.com/tylerb/graceful"
	"gopkg.in/inconshreveable/log15.v2"
	logext "gopkg.in/inconshreveable/log15.v2/ext"
	"gopkg.in/redis.v5"

	dalpg "git.softndit.com/collector/backend/dal/pg"
	"git.softndit.com/collector/backend/restapi/operations"
	"git.softndit.com/collector/backend/restapi/operations/actors"
	"git.softndit.com/collector/backend/restapi/operations/auth"
	"git.softndit.com/collector/backend/restapi/operations/badges"
	"git.softndit.com/collector/backend/restapi/operations/chat"
	"git.softndit.com/collector/backend/restapi/operations/collections"
	"git.softndit.com/collector/backend/restapi/operations/dashboard"
	"git.softndit.com/collector/backend/restapi/operations/events"
	"git.softndit.com/collector/backend/restapi/operations/groups"
	"git.softndit.com/collector/backend/restapi/operations/invites"
	"git.softndit.com/collector/backend/restapi/operations/materials"
	"git.softndit.com/collector/backend/restapi/operations/medias"
	"git.softndit.com/collector/backend/restapi/operations/messages"
	"git.softndit.com/collector/backend/restapi/operations/nameddateintervals"
	"git.softndit.com/collector/backend/restapi/operations/objects"
	"git.softndit.com/collector/backend/restapi/operations/objectstatus"
	"git.softndit.com/collector/backend/restapi/operations/originlocations"
	"git.softndit.com/collector/backend/restapi/operations/public_collections"
	"git.softndit.com/collector/backend/restapi/operations/public_objects"
	"git.softndit.com/collector/backend/restapi/operations/references"
	"git.softndit.com/collector/backend/restapi/operations/rights"
	"git.softndit.com/collector/backend/restapi/operations/roots"
	"git.softndit.com/collector/backend/restapi/operations/session"
	"git.softndit.com/collector/backend/restapi/operations/tasks"
	"git.softndit.com/collector/backend/restapi/operations/teams"
	"git.softndit.com/collector/backend/restapi/operations/users"
	"git.softndit.com/collector/backend/restapi/operations/users_ban_list"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"git.softndit.com/collector/backend/wsserver"
	"github.com/rs/cors"
)

var (
	logger   = log15.New("app", "collector")
	dumpBody = true
)

var cliOptions struct {
	ConfigPath string `long:"config" description:"base config" required:"1" env:"COLLECTORD_CFG"`
}

func die(msg string, err error) {
	if err != nil {
		logger.Crit(msg, "err", err)
	} else {
		logger.Crit(msg)
	}
	os.Exit(1)
}

func configureFlags(api *operations.CollectordAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{
			Options: &cliOptions,
		},
	}
}

func configureAPI(api *operations.CollectordAPI) http.Handler {
	configPath := cliOptions.ConfigPath

	mainConfig := &config.Config{}
	if err := mainConfig.ReadConfig(configPath); err != nil {
		die("can't load config", err)
	}

	instanceHolder, err := config.NewInstanceHolder(mainConfig)
	if err != nil {
		die("can't init config", err)
	}

	serverConfig := mainConfig.ServerConfig

	// DBM
	pgConfig, err := instanceHolder.GetPgConfig(serverConfig.DBM.PgRef)

	if os.Getenv("COLLECTOR_SQL_LOGS") == "1" {
		pgConfig.Logger = log15adapter.NewLogger(logger.New("service", "pg-raw"))
	}

	if err != nil {
		die("dbm err", err)
	}
	DBM, err := dalpg.NewManager(&dalpg.ManagerContext{
		PoolConfig: pgConfig,
		Log:        logger.New("service", "pg"),
	})
	if err != nil {
		die("dbm err", err)
	}

	// file storage
	filerCfg, err := instanceHolder.GetGenericFilerConfig(serverConfig.FilerClient.FilerRef)
	if err != nil {
		die("can't load file storage", err)
	}
	builder, err := services.GetTomlFilerBuilder(filerCfg.Type)
	if err != nil {
		die("can't get toml filer builder", err)
	}

	fileStorage, err := builder(&services.TomlFileStorageContext{
		Log:    logger.New("service", "filer"),
		Config: filerCfg.Config,
	})
	if err != nil {
		die("can't init storage", err)
	}

	// cleaver
	rabbitConfig, err := instanceHolder.GetRabbitConfig(serverConfig.CleaverClient.RabbitRef)
	if err != nil {
		die("can't init rabbitmq server", err)
	}
	cleaverClient := cleaver.NewRMQClient(rabbitConfig.URL)
	if err := cleaverClient.Connect(); err != nil {
		die("can't init cleaver", err)
	}

	// search // elastic
	elasticConfig, err := instanceHolder.GetElasticSearchConfig(serverConfig.SearchClient.ElasticSearchRef)
	if err != nil {
		die("can't init elastic", err)
	}

	searchClient, err := services.NewElasticSearchClient(
		&services.ElasticSearchClientContext{
			URL:         elasticConfig.URL,
			ObjectIndex: serverConfig.SearchClient.ObjectIndex,
		},
	)
	if err != nil {
		die("can't init elastic", err)
	}

	// Autheticator
	dbAuth := services.DBAutheticator{
		DBM: DBM,
		Log: logger,
	}

	// probably make infinity
	eventCh := make(chan *models.UserPayloadedEvent, wsserver.EventChanBuffer)

	// job pool
	jobPool := delayedjob.NewPool(10, 1*time.Second)
	if err := jobPool.AsyncRun(); err != nil {
		die("can't init jobpool", err)
	}

	// push client
	rabbitConfig, err = instanceHolder.GetRabbitConfig(serverConfig.PusherClient.RabbitRef)
	if err != nil {
		die("can't init push client", err)
	}
	pushClient := npusherclient.NewRMQClient(rabbitConfig.URL)
	if err := pushClient.Connect(); err != nil {
		die("can't init push client", err)
	}

	// mail client
	mailConfig := &services.SMTPMailConfig{}
	if err := mailConfig.ParseServerConfig(&serverConfig.MailClient); err != nil {
		die("can't init mail client", err)
	}
	emailClient, err := services.NewSMTPEmail(mailConfig)
	if err != nil {
		die("can't init mail client", err)
	}

	// init i18n
	if err := util.InitTranslations(mainConfig.ServerConfig.I18nPath); err != nil {
		panic(err)
	}

	// init templates
	templates := loadTemplates(mainConfig.ServerConfig.TemplatePath)

	// redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: serverConfig.Redis.Url,
	})
	if err := redisClient.Ping().Err(); err != nil {
		panic(err)
	}
	tokenStorage := storage.NewRedisStorage(redisClient, time.Hour*24)

	// resource context
	resourceContext := resource.Context{
		DBM:           DBM,
		Log:           logger,
		FileStorage:   fileStorage,
		CleaverClient: cleaverClient,
		SearchClient:  searchClient,
		PushClient:    pushClient,
		EventSender: services.NewEventSender(eventCh, &services.PayloadedEventFiller{
			Log: logger.New("service", "filler"),
			DBM: DBM,
		}),
		MessengerClient: &services.MessengerClient{
			Log:        logger.New("service", "messanger"),
			DBM:        DBM,
			PushClient: pushClient,
		},
		MailClient: emailClient,
		JobPool:    jobPool,
	}

	// hub handles new connections from wsserver and routes events to them
	hub := &wsserver.Hub{
		NewSessionCh:   make(chan *wsserver.Session),
		CloseSessionCh: make(chan *wsserver.Session),
		EventCh:        eventCh,
		Log:            logger.New("service", "hub"),
	}

	go hub.Run()
	go wsserver.SetUpWebsocketServer(&wsserver.WebSocketContext{
		Log:  logger,
		Auth: &dbAuth,
		Hub:  hub,
		Port: serverConfig.EventService.Port,
	})

	// objects
	objectResource := &resource.Object{Context: resourceContext}
	api.ObjectsGetObjectsIDHandler = objects.GetObjectsIDHandlerFunc(objectResource.GetObjectByID)
	api.ObjectsPostObjectsHandler = objects.PostObjectsHandlerFunc(objectResource.CreateObject)
	api.ObjectsPutObjectsIDHandler = objects.PutObjectsIDHandlerFunc(objectResource.UpdateObject)
	api.ObjectsPostObjectsSearchHandler = objects.PostObjectsSearchHandlerFunc(objectResource.SearchObjects)
	api.ObjectsDeleteObjectsIDDeleteHandler = objects.DeleteObjectsIDDeleteHandlerFunc(objectResource.DeleteObject)
	api.ObjectsPostObjectsMoveHandler = objects.PostObjectsMoveHandlerFunc(objectResource.MoveObjects)
	api.ObjectsPostObjectsDeleteHandler = objects.PostObjectsDeleteHandlerFunc(objectResource.ButchDeleteObject)

	// badges
	badgeResource := &resource.Badge{Context: resourceContext}
	api.BadgesGetBadgeGetHandler = badges.GetBadgeGetHandlerFunc(badgeResource.GetBadgesForRoot)
	api.BadgesPostBadgeNewHandler = badges.PostBadgeNewHandlerFunc(badgeResource.CreateBadge)
	api.BadgesPostBadgeIDUpdateHandler = badges.PostBadgeIDUpdateHandlerFunc(badgeResource.UpdateBadge)
	api.BadgesDeleteBadgeIDDeleteHandler = badges.DeleteBadgeIDDeleteHandlerFunc(badgeResource.DeleteBadge)

	// medias
	mediasResource := &resource.Media{Context: resourceContext}
	api.MediasPostMediasHandler = medias.PostMediasHandlerFunc(mediasResource.CreateMedia)
	api.MediasGetMediasByIdsHandler = medias.GetMediasByIdsHandlerFunc(mediasResource.GetMediasByIDs)
	api.MediasGetMediasHandler = medias.GetMediasHandlerFunc(mediasResource.GetMediaLocation)

	// references
	referenceResource := &resource.Reference{Context: resourceContext}
	api.ReferencesGetReferencesHandler = references.GetReferencesHandlerFunc(referenceResource.GetReferences)

	// collections
	collectionResource := &resource.Collection{Context: resourceContext}
	api.CollectionsPostCollectionsHandler = collections.PostCollectionsHandlerFunc(collectionResource.CreateCollection)
	api.CollectionsPutCollectionsIDHandler = collections.PutCollectionsIDHandlerFunc(collectionResource.UpdateCollection)
	api.CollectionsPostCollectionsObjectsHandler = collections.PostCollectionsObjectsHandlerFunc(collectionResource.GetCollectionsObjects)
	api.CollectionsGetCollectionsDraftHandler = collections.GetCollectionsDraftHandlerFunc(collectionResource.GetDraftCollection)
	api.CollectionsDeleteCollectionsIDDeleteHandler = collections.DeleteCollectionsIDDeleteHandlerFunc(collectionResource.DeleteCollection)
	api.CollectionsPostCollectionsAddToGroupHandler = collections.PostCollectionsAddToGroupHandlerFunc(collectionResource.AddCollectionsToGroup)
	api.CollectionsPostCollectionsRemoveFromGroupHandler = collections.PostCollectionsRemoveFromGroupHandlerFunc(collectionResource.RemoveCollectionsFromGroup)
	api.CollectionsGetCollectionsIDHandler = collections.GetCollectionsIDHandlerFunc(collectionResource.GetCollection)

	// groups
	groupResource := &resource.Group{Context: resourceContext}
	api.GroupsGetGroupsIDHandler = groups.GetGroupsIDHandlerFunc(groupResource.GetGroupByID)
	api.GroupsPostGroupsHandler = groups.PostGroupsHandlerFunc(groupResource.CreateGroup)
	api.GroupsPutGroupsIDUpdateHandler = groups.PutGroupsIDUpdateHandlerFunc(groupResource.UpdateGroup)
	api.GroupsDeleteGroupsIDDeleteHandler = groups.DeleteGroupsIDDeleteHandlerFunc(groupResource.DeleteGroup)

	// dashboard
	dashboardResource := &resource.Dashboard{Context: resourceContext}
	api.DashboardGetDashboardHandler = dashboard.GetDashboardHandlerFunc(dashboardResource.GetDashboard)

	// auth
	authResources := &resource.Auth{Context: resourceContext, EmailTokenStorage: tokenStorage, Templates: templates}
	api.AuthPostAuthRegHandler = auth.PostAuthRegHandlerFunc(authResources.Registration)
	api.AuthPostAuthLoginHandler = auth.PostAuthLoginHandlerFunc(authResources.Login)
	api.AuthPostAuthLogoutHandler = auth.PostAuthLogoutHandlerFunc(authResources.Logout)
	api.AuthPostAuthPasswordRecoveryHandler = auth.PostAuthPasswordRecoveryHandlerFunc(authResources.RecoveryPassword)
	api.AuthPostAuthPasswordResetTokenHandler = auth.PostAuthPasswordResetTokenHandlerFunc(authResources.ResetPassword)

	api.AuthGetAuthRegConfirmEmailHandler = auth.GetAuthRegConfirmEmailHandlerFunc(authResources.ConfirmEmail)
	api.AuthPostAuthRegConfirmEmailHandler = auth.PostAuthRegConfirmEmailHandlerFunc(authResources.SendConfirmEmailToken)

	// user
	userResources := &resource.User{Context: resourceContext}
	api.UsersPutUserHandler = users.PutUserHandlerFunc(userResources.UpdateUser)
	api.UsersGetUserAboutHandler = users.GetUserAboutHandlerFunc(userResources.UserInfo)
	api.UsersPostUserSearchByNameHandler = users.PostUserSearchByNameHandlerFunc(userResources.SearchByName)
	api.UsersPutUserUpdateLocaleHandler = users.PutUserUpdateLocaleHandlerFunc(userResources.UpdateUserLocale)

	// root
	rootResource := &resource.Root{Context: resourceContext}
	api.RootsGetRootIDHandler = roots.GetRootIDHandlerFunc(rootResource.RootInfo)
	api.RootsGetRootByUserHandler = roots.GetRootByUserHandlerFunc(rootResource.UserRoots)
	api.RootsPostRootAddUserHandler = roots.PostRootAddUserHandlerFunc(rootResource.AddUserToRoot)
	api.RootsPostRootRemoveUserHandler = roots.PostRootRemoveUserHandlerFunc(rootResource.RemoveUserFromRoot)

	// event
	eventResource := &resource.Event{Context: resourceContext}
	api.EventsPostEventHandler = events.PostEventHandlerFunc(eventResource.GetEvents)
	api.EventsPostEventConfirmHandler = events.PostEventConfirmHandlerFunc(eventResource.ConfirmEvents)

	// message
	messageResource := &resource.Message{Context: resourceContext, Templates: templates}
	api.MessagesPostMessageHandler = messages.PostMessageHandlerFunc(messageResource.SendMessage)
	api.MessagesPostMessageByidsHandler = messages.PostMessageByidsHandlerFunc(messageResource.GetMessages)
	api.MessagesPostMessageRangeHandler = messages.PostMessageRangeHandlerFunc(messageResource.GetMessagesRange)
	api.MessagesPostMessageAllConversationHandler = messages.PostMessageAllConversationHandlerFunc(messageResource.GetAllConversations)
	api.MessagesPostMessageReadHistoryHandler = messages.PostMessageReadHistoryHandlerFunc(messageResource.ReadHistory)

	// chat
	chatResource := &resource.Chat{Context: resourceContext}
	api.ChatPostChatHandler = chat.PostChatHandlerFunc(chatResource.CreateChat)
	api.ChatPostChatAddUserHandler = chat.PostChatAddUserHandlerFunc(chatResource.AddUser)
	api.ChatPostChatRemoveUserHandler = chat.PostChatRemoveUserHandlerFunc(chatResource.RemoveUser)
	api.ChatPostChatIDChangeAvatarHandler = chat.PostChatIDChangeAvatarHandlerFunc(chatResource.ChangeAvatar)
	api.ChatPostChatIDChangeNameHandler = chat.PostChatIDChangeNameHandlerFunc(chatResource.ChangeName)
	api.ChatGetChatIDHandler = chat.GetChatIDHandlerFunc(chatResource.GetChat)

	// session
	sessionResource := &resource.Session{Context: resourceContext}
	api.SessionPostSessionRegisterDeviceTokenHandler = session.PostSessionRegisterDeviceTokenHandlerFunc(sessionResource.SetDeviceToken)

	// actors
	actorResource := &resource.Actor{Context: resourceContext}
	api.ActorsPostActorNewHandler = actors.PostActorNewHandlerFunc(actorResource.CreateActor)
	api.ActorsPostActorIDUpdateHandler = actors.PostActorIDUpdateHandlerFunc(actorResource.UpdateActor)
	api.ActorsGetActorGetHandler = actors.GetActorGetHandlerFunc(actorResource.GetActorsForRoot)
	api.ActorsDeleteActorIDDeleteHandler = actors.DeleteActorIDDeleteHandlerFunc(actorResource.DeleteActor)

	// materials
	materialResource := &resource.Material{Context: resourceContext}
	api.MaterialsPostMaterialNewHandler = materials.PostMaterialNewHandlerFunc(materialResource.CreateMaterial)
	api.MaterialsPostMaterialIDUpdateHandler = materials.PostMaterialIDUpdateHandlerFunc(materialResource.UpdateMaterial)
	api.MaterialsGetMaterialGetHandler = materials.GetMaterialGetHandlerFunc(materialResource.GetMaterialsForRoot)
	api.MaterialsDeleteMaterialIDDeleteHandler = materials.DeleteMaterialIDDeleteHandlerFunc(materialResource.DeleteMaterial)

	// originLocations
	originLocationResource := &resource.OriginLocation{Context: resourceContext}
	api.OriginlocationsPostOriginLocationNewHandler = originlocations.PostOriginLocationNewHandlerFunc(originLocationResource.CreateOriginLocation)
	api.OriginlocationsPostOriginLocationIDUpdateHandler = originlocations.PostOriginLocationIDUpdateHandlerFunc(originLocationResource.UpdateOriginLocation)
	api.OriginlocationsGetOriginLocationGetHandler = originlocations.GetOriginLocationGetHandlerFunc(originLocationResource.GetOriginLocationsForRoot)
	api.OriginlocationsDeleteOriginLocationIDDeleteHandler = originlocations.DeleteOriginLocationIDDeleteHandlerFunc(originLocationResource.DeleteOriginLocation)

	// named date interval
	namedDateIntervalResource := &resource.NamedDateInterval{Context: resourceContext}
	api.NameddateintervalsPostNamedDateIntervalIDUpdateHandler = nameddateintervals.PostNamedDateIntervalIDUpdateHandlerFunc(namedDateIntervalResource.UpdateNamedDateInterval)
	api.NameddateintervalsGetNamedDateIntervalGetHandler = nameddateintervals.GetNamedDateIntervalGetHandlerFunc(namedDateIntervalResource.GetNamedDatesIntervalsForRoot)
	api.NameddateintervalsDeleteNamedDateIntervalIDDeleteHandler = nameddateintervals.DeleteNamedDateIntervalIDDeleteHandlerFunc(namedDateIntervalResource.DeleteNamedDateInterval)
	api.NameddateintervalsPostNamedDateIntervalNewHandler = nameddateintervals.PostNamedDateIntervalNewHandlerFunc(namedDateIntervalResource.CreateNamedDateInterval)

	// object status
	objectStatusResource := &resource.ObjectStatus{Context: resourceContext}
	api.ObjectstatusPostObjectStatusNewHandler = objectstatus.PostObjectStatusNewHandlerFunc(objectStatusResource.CreateObjectStatus)

	// team
	teamResource := &resource.Team{Context: resourceContext}
	api.TeamsGetTeamByRootIDHandler = teams.GetTeamByRootIDHandlerFunc(teamResource.GetTeamByRootID)

	// invite
	inviteResource := &resource.Invite{Context: resourceContext, Templates: templates}
	api.InvitesPostInviteNewHandler = invites.PostInviteNewHandlerFunc(inviteResource.CreateInvite)
	api.InvitesPostInviteIDAcceptHandler = invites.PostInviteIDAcceptHandlerFunc(inviteResource.AcceptInvite)
	api.InvitesPostInviteIDRejectHandler = invites.PostInviteIDRejectHandlerFunc(inviteResource.RejectInvite)
	api.InvitesPostInviteIDCancelHandler = invites.PostInviteIDCancelHandlerFunc(inviteResource.CancelInvite)

	// task
	taskResource := &resource.Task{Context: resourceContext, Templates: templates}
	api.TasksGetTaskIDHandler = tasks.GetTaskIDHandlerFunc(taskResource.GetTaskByID)
	api.TasksPostTaskHandler = tasks.PostTaskHandlerFunc(taskResource.CreateTask)
	api.TasksPutTaskIDHandler = tasks.PutTaskIDHandlerFunc(taskResource.EditTask)
	api.TasksDeleteTaskIDHandler = tasks.DeleteTaskIDHandlerFunc(taskResource.DeleteTask)
	api.TasksGetTaskMyListHandler = tasks.GetTaskMyListHandlerFunc(taskResource.GetMyList)
	api.TasksPostTaskIDChangeStatusHandler = tasks.PostTaskIDChangeStatusHandlerFunc(taskResource.ChangeStatus)
	api.TasksPostTaskIDArchiveHandler = tasks.PostTaskIDArchiveHandlerFunc(taskResource.Archive)
	api.TasksPostTaskMyArchiveListHandler = tasks.PostTaskMyArchiveListHandlerFunc(taskResource.GetMyArchiveList)
	api.TasksPostTaskIDAssignToHandler = tasks.PostTaskIDAssignToHandlerFunc(taskResource.AssignTo)

	// user entity rights
	userEntityRightResource := &resource.UserEntityRight{Context: resourceContext}
	api.RightsPutRightHandler = rights.PutRightHandlerFunc(userEntityRightResource.SetRight)
	api.RightsGetRightHandler = rights.GetRightHandlerFunc(userEntityRightResource.GetRights)

	// public collections
	publicCollectionsResource := &resource.PublicCollections{Context: resourceContext}
	api.PublicCollectionsPostPublicCollectionsHandler = public_collections.PostPublicCollectionsHandlerFunc(publicCollectionsResource.GetList)
	api.PublicCollectionsPostPublicCollectionsObjectsHandler = public_collections.PostPublicCollectionsObjectsHandlerFunc(publicCollectionsResource.GetObjects)
	api.PublicCollectionsGetPublicCollectionsIDHandler = public_collections.GetPublicCollectionsIDHandlerFunc(publicCollectionsResource.GetCollection)

	// public objects
	publicObjectsResource := &resource.PublicObjects{Context: resourceContext}
	api.PublicObjectsGetPublicObjectsIDHandler = public_objects.GetPublicObjectsIDHandlerFunc(publicObjectsResource.GetPublicObjectByID)

	// users ban list
	usersBanResource := &resource.UsersBan{Context: resourceContext}
	api.UsersBanListGetUsersBanListHandler = users_ban_list.GetUsersBanListHandlerFunc(usersBanResource.GetUsersBanList)
	api.UsersBanListPostUsersBanListAddHandler = users_ban_list.PostUsersBanListAddHandlerFunc(usersBanResource.UsersBanListAdd)
	api.UsersBanListPostUsersBanListRemoveHandler = users_ban_list.PostUsersBanListRemoveHandlerFunc(usersBanResource.UsersBanListRemove)

	// invite by email
	inviteByEmailResource := &resource.InviteByEmail{Context: resourceContext, Templates: templates}
	api.InvitesPostInviteByEmailHandler = invites.PostInviteByEmailHandlerFunc(inviteByEmailResource.Send)

	// shutdown
	api.ServerShutdown = func() {}

	// auth
	api.APIKeyAuth = func(token string) (interface{}, error) {
		return dbAuth.Auth(token)
	}

	// configure the api here
	api.ServeError = errors.ServeError
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return addLogging(handler)
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "DELETE", "PUT", "PATCH"},
		Debug:            true,
	}).Handler

	return handleCORS(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err == nil {
					return
				}
				stack := make([]byte, 1024*32)
				stack = stack[:goruntime.Stack(stack, false)]
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error("panic", "err", err, "stack", string(stack))
			}()
			handler.ServeHTTP(w, r)
		}))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

func loadTemplates(templatePath string) *template.Template {
	var templates []string

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			templates = append(templates, path)
		}
		return nil
	}

	err := filepath.Walk(templatePath, fn)
	if err != nil {
		return nil
	}

	return template.Must(template.ParseFiles(templates...))
}

func addLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.New("req_id", logext.RandId(8))

		begin := time.Now()

		var reqDump []byte
		if dumpBody && r.Body != nil && r.ContentLength > 0 && r.Header.Get("Content-Type") == "application/json" {
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(r.Body); err != nil {
				log.Debug("Can't read request body", "err", err)
				return
			}
			if err := r.Body.Close(); err != nil {
				log.Debug("Can't close request body", "err", err)
				return
			}
			r.Body = io.NopCloser(&buf)
			reqDump = buf.Bytes()
		}

		log.Debug("new request",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"content-len", r.ContentLength,
			"dump", string(reqDump))

		rw := newDumpResponseWriter(w, dumpBody)
		next.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), "logger", log)))

		log.Debug("request processed",
			"status", rw.status,
			"elapsed", time.Since(begin),
			"dump", string(rw.body.Bytes()))

	})
}

func addTranslateFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "T", util.GetTranslationFunc(r.Header.Get("Accept-Language")))))
	})
}

// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme string, socketpath string) {
}

func newDumpResponseWriter(w http.ResponseWriter, recordBody bool) *dumpResponseWriter {
	var buf *bytes.Buffer
	if recordBody {
		buf = &bytes.Buffer{}
	}
	drw := &dumpResponseWriter{ResponseWriter: w, body: buf, recordBody: recordBody}
	return drw
}

type dumpResponseWriter struct {
	http.ResponseWriter
	recordBody bool
	body       *bytes.Buffer
	status     int
}

func (rw *dumpResponseWriter) WriteHeader(s int) {
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}
func (rw *dumpResponseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	if rw.recordBody {
		rw.body.Write(b[:n])
	}
	return n, err
}

func (rw *dumpResponseWriter) Flush() {
	rw.ResponseWriter.(http.Flusher).Flush()
}

func (rw *dumpResponseWriter) CloseNotify() <-chan bool {
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
