package dto

// PeerType TBD
type PeerType int16

// CubeUserID TBD
const CubeUserID int64 = 7

const (
	// PeerTypeUser TBD
	PeerTypeUser PeerType = iota + 1
	// PeerTypeChat TBD
	PeerTypeChat
)

// IsUser TBD
func (p *PeerType) IsUser() bool {
	return *p == PeerTypeUser
}

// IsChat TBD
func (p *PeerType) IsChat() bool {
	return *p == PeerTypeChat
}

// Peer TBD
type Peer struct {
	ID   int64    `json:"id"`
	Type PeerType `json:"type"`
}

// IsUser TBD
func (p *Peer) IsUser() bool {
	return p.Type.IsUser()
}

// IsChat TBD
func (p *Peer) IsChat() bool {
	return p.Type.IsChat()
}

// PeerList TBD
type PeerList []*Peer

// ChatIDs TBD
func (p PeerList) ChatIDs() []int64 {
	chatIDs := make([]int64, 0, len(p))
	for _, peer := range p {
		if peer.Type.IsChat() {
			chatIDs = append(chatIDs, peer.ID)
		}
	}
	return chatIDs
}
