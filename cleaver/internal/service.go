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

// ServiceContext TBD
type ServiceContext struct {
	name     string
	config   *toml.Primitive
	executor *Executor
	log      log15.Logger
}

// NewServiceFunc TBD
type NewServiceFunc func(ctx *ServiceContext) (Service, error)

var serviceRegistry = make(map[string]NewServiceFunc)

// RegisterService TBD
func RegisterService(typ string, factory NewServiceFunc) {
	serviceRegistry[typ] = factory
}

// CreateService TBD
func CreateService(typ string, ctx *ServiceContext) (Service, error) {
	factory := serviceRegistry[typ]
	if factory == nil {
		return nil, fmt.Errorf("unregestered service type '%s'", typ)
	}

	return factory(ctx)
}
