package dto

//go:generate dbgen -type UnreadMessage

// UnreadMessage TBD
type UnreadMessage struct {
	ID        int64    `db:"id"`
	MessageID int64    `db:"message_id"`
	UserID    int64    `db:"user_id"`
	PeerID    int64    `db:"peer_id"`
	PeerType  PeerType `db:"peer_type"`
}
