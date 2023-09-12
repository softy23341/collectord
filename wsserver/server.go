package wsserver

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	uauth "git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/services"
	openapierrors "github.com/go-openapi/errors"
	"github.com/gorilla/websocket"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

// EventChanBuffer TBD
const (
	EventChanBuffer = 100000
	writeWait       = 10 * time.Second
	pongWait        = 60 * time.Second
	pingPeriod      = (pongWait * 9) / 10
)

// WebSocketContext TBD
type WebSocketContext struct {
	Log  log15.Logger
	Auth services.Autheticator
	Hub  *Hub
	Port int
}

// Session TBD
type Session struct {
	userContext *uauth.UserContext
	eventCh     chan *models.UserPayloadedEvent

	conn *websocket.Conn
}

func (s *Session) sendEvent(e *models.UserPayloadedEvent) {
	s.eventCh <- e
}

type sessionInfo struct {
	begin time.Time
}

// Hub TBD
type Hub struct {
	NewSessionCh   chan *Session
	CloseSessionCh chan *Session
	EventCh        chan *models.UserPayloadedEvent

	Log log15.Logger
}

// NewSession TBD
func (h *Hub) NewSession(s *Session) {
	h.NewSessionCh <- s
}

// CloseSession TBD
func (h *Hub) CloseSession(s *Session) {
	h.CloseSessionCh <- s
}

type sessionReg struct {
	sessions map[int64]map[*Session]*sessionInfo
}

func newSessionReg() *sessionReg {
	return &sessionReg{
		sessions: make(map[int64]map[*Session]*sessionInfo),
	}
}

func (cr *sessionReg) addSession(s *Session) {
	if cr.sessions[s.userContext.User.ID] == nil {
		cr.sessions[s.userContext.User.ID] = make(map[*Session]*sessionInfo)
	}
	cr.sessions[s.userContext.User.ID][s] = &sessionInfo{begin: time.Now()}
}

func (cr *sessionReg) removeSession(s *Session) {
	if cr.sessions[s.userContext.User.ID] != nil {
		delete(cr.sessions[s.userContext.User.ID], s)
	}

	if i, find := cr.sessions[s.userContext.User.ID]; find && i == nil {
		delete(cr.sessions, s.userContext.User.ID)
	}
}

func (cr *sessionReg) sessionsByUserID(userID int64) []*Session {
	sessions := cr.sessions[userID]
	if len(sessions) == 0 {
		return nil
	}
	oSessions := make([]*Session, 0, len(sessions))
	for session := range sessions {
		oSessions = append(oSessions, session)
	}
	return oSessions
}

// Run TBD
func (h *Hub) Run() {
	newSessionCh := h.NewSessionCh
	closeSessionCh := h.CloseSessionCh
	eventCh := h.EventCh

	sessionReg := newSessionReg()
	for {
		select {
		case session, open := <-newSessionCh:
			if !open {
				newSessionCh = nil
			}
			sessionReg.addSession(session)
		case session, open := <-closeSessionCh:
			if !open {
				closeSessionCh = nil
			}
			sessionReg.removeSession(session)
		case event, open := <-eventCh:
			if !open {
				eventCh = nil
			}
			for _, s := range sessionReg.sessionsByUserID(event.UserID) {
				s.sendEvent(event)
			}
		}
	}
}

// SetUpWebsocketServer TBD
func SetUpWebsocketServer(context *WebSocketContext) error {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		authToken := r.URL.Query().Get("auth-token")
		versionString := r.URL.Query().Get("version")
		version, _ := strconv.Atoi(versionString)

		userContext, err := context.Auth.Auth(authToken)
		if err != nil {
			if openapiError, ok := err.(openapierrors.Error); ok {
				context.Log.Error("token is wrong", "err", openapiError.Error())
				http.Error(w, openapiError.Error(), int(openapiError.Code()))
				return
			}
			context.Log.Error("cant auth", "err", err)
			http.Error(w, "", 500)
			return
		}

		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			context.Log.Error("cant upgrade conn", "err", err)
			return
		}
		defer conn.Close()
		defer context.Log.Warn("close socket connection", "authToken", authToken)

		session := &Session{
			userContext: userContext,
			eventCh:     make(chan *models.UserPayloadedEvent, EventChanBuffer),
			conn:        conn,
		}

		context.Hub.NewSession(session)
		defer context.Hub.CloseSession(session)

		for {
			select {
			case userPayloadedEvent, open := <-session.eventCh:
				if !open {
					return
				}

				if version == 0 && *userPayloadedEvent.PayloadedEvent.Typo > 10 {
					break
				}
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteJSON(userPayloadedEvent.PayloadedEvent); err != nil {
					context.Log.Warn("close socket connection by writeJSON",
						"authToken", authToken,
						"err", err,
					)
					return
				}
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					context.Log.Warn("close socket connection by timeout",
						"authToken", authToken,
						"err", err,
					)
					return
				}
			}
		}
	}

	http.HandleFunc("/ws", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", context.Port), nil)

	return nil
}
