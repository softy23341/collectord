package internal

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/inconshreveable/log15"
)

// Service TBD
type Service interface {
	Run() error
}

// ServiceCtx TBD
type ServiceCtx struct {
	Config   *toml.Primitive
	Log      log15.Logger
	Provider Provider
}

// NewServiceFunc TBD
type NewServiceFunc func(config *ServiceCtx) (Service, error)

var servicesRegistry = make(map[string]NewServiceFunc)

// RegesterService TBD
func RegesterService(typo string, factoryMethod NewServiceFunc) {
	servicesRegistry[typo] = factoryMethod
}

// GetServiceBuilder TBD
func GetServiceBuilder(typo string) (NewServiceFunc, error) {
	factoryMethod := servicesRegistry[typo]
	if factoryMethod == nil {
		return nil, fmt.Errorf("unregestered service name")
	}
	return factoryMethod, nil
}
