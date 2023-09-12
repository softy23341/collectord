package client

import (
	"git.softndit.com/collector/backend/cleaver"
)

// Client TBD
type Client interface {
	Resize(*cleaver.ResizeTask) ([]*cleaver.TransformResult, error)
	Copy(*cleaver.CopyTask) (*cleaver.CopyResult, error)
}

// ConnectClient TBD
type ConnectClient interface {
	Client
	Connect() error
}
