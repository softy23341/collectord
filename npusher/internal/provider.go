package internal

import (
	"fmt"

	"git.softndit.com/collector/backend/npusher"

	"github.com/BurntSushi/toml"
	"github.com/inconshreveable/log15"
)

// Provider TBD
type Provider interface {
	SupportedTypes() []string
	Send(task *npusher.NotificationTask) error
}

// ProviderCtx TBD
type ProviderCtx struct {
	Log    log15.Logger
	Config *toml.Primitive
}

// NewProviderFunc TBD
type NewProviderFunc func(ctx *ProviderCtx) (Provider, error)

var providerRegistry = make(map[string]NewProviderFunc)

// RegisterProvider TBD
func RegisterProvider(typo string, factoryMethod NewProviderFunc) {
	providerRegistry[typo] = factoryMethod
}

// GetProviderFactory TBD
func GetProviderFactory(typo string) (NewProviderFunc, error) {
	factoryMethod := providerRegistry[typo]
	if factoryMethod == nil {
		return nil, fmt.Errorf("unregestered provider name")
	}
	return factoryMethod, nil
}
