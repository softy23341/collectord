package internal

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/inconshreveable/log15"
)

// TokenStorager TBD
type tokenStorager interface {
	get(deviceToken string) (arn string, err error)
	set(deviceToken string, arn string) (err error)
}

type tokenStorageCtx struct {
	log    log15.Logger
	config *toml.Primitive
}

type newTokenStorageFunc func(ctx *tokenStorageCtx) (tokenStorager, error)

var tokenStorageRegistry = make(map[string]newTokenStorageFunc)

func registerTokenStorage(typo string, factoryMethod newTokenStorageFunc) {
	tokenStorageRegistry[typo] = factoryMethod
}

func getTokenStorageFactory(typo string) (newTokenStorageFunc, error) {
	factoryMethod := tokenStorageRegistry[typo]
	if factoryMethod == nil {
		return nil, fmt.Errorf("unregistered token storage type: %s", typo)
	}
	return factoryMethod, nil
}

func createTokenStorage(log log15.Logger, config *toml.Primitive) (tokenStorager, error) {
	storageInfo := &struct {
		Typo string `toml:"type"`
	}{}
	if err := toml.PrimitiveDecode(*config, storageInfo); err != nil {
		return nil, err
	}
	factory, err := getTokenStorageFactory(storageInfo.Typo)
	if err != nil {
		return nil, err
	}
	return factory(&tokenStorageCtx{
		config: config,
		log:    log.New("storage", storageInfo.Typo),
	})
}
