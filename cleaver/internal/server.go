package internal

import (
	"errors"
	"fmt"

	"github.com/inconshreveable/log15"

	"github.com/BurntSushi/toml"
)

// NewServer TBD
func NewServer(configFile string, log log15.Logger) (*Server, error) {
	server := &Server{
		log:      log,
		services: make(map[string]Service),
	}

	if _, err := toml.DecodeFile(configFile, &server.config); err != nil {
		return nil, err
	}

	if err := server.configureExecutor(); err != nil {
		return nil, err
	}

	if err := server.configureServices(); err != nil {
		return nil, err
	}

	return server, nil
}

// Server TBD
type Server struct {
	log log15.Logger

	config struct {
		Engine   string
		Getters  map[string]toml.Primitive
		Putters  map[string]toml.Primitive
		Services map[string]toml.Primitive
	}

	executor *Executor
	services map[string]Service
}

// Run TBD
func (s *Server) Run() error {
	s.log.Info("Running server")
	for name, svc := range s.services {
		s.log.Debug("Running service", "svcname", name)
		if err := svc.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) configureExecutor() error {
	engine, err := CreateImageEngine(s.config.Engine)
	if err != nil {
		return err
	}

	if len(s.config.Getters) == 0 {
		return errors.New("no getters specified")
	}

	if len(s.config.Putters) == 0 {
		return errors.New("no putters specified")
	}

	sgp := newSchemeGetterPutter()

	type commonCfg struct {
		Type    string
		Default bool
	}

	for scheme, config := range s.config.Getters {
		var base commonCfg
		if err := toml.PrimitiveDecode(config, &base); err != nil {
			return fmt.Errorf("can't configure getter for scheme '%s': %v", scheme, err)
		}

		getter, err := CreateGetter(base.Type, &config)
		if err != nil {
			return fmt.Errorf("can't configure getter for scheme '%s': %v", scheme, err)
		}

		sgp.addGetter(scheme, getter)
		if base.Default {
			sgp.addGetter("", getter)
		}
	}

	for scheme, config := range s.config.Putters {
		var base commonCfg
		if err := toml.PrimitiveDecode(config, &base); err != nil {
			return fmt.Errorf("can't configure putter for scheme '%s': %v", scheme, err)
		}

		putter, err := CreatePutter(base.Type, &config)
		if err != nil {
			return fmt.Errorf("can't configure putter for scheme '%s': %v", scheme, err)
		}

		sgp.addPutter(scheme, putter)
		if base.Default {
			sgp.addPutter("", putter)
		}
	}

	s.executor = NewExecutor(engine, sgp, sgp)
	return nil
}

func (s *Server) configureServices() error {
	if len(s.config.Services) == 0 {
		return errors.New("no services specified")
	}

	for name, config := range s.config.Services {
		var base struct{ Type string }
		if err := toml.PrimitiveDecode(config, &base); err != nil {
			return fmt.Errorf("invalid config for service '%s': %v", name, err)
		}

		ctx := ServiceContext{
			name:     name,
			config:   &config,
			executor: s.executor,
			log:      s.log.New("svctype", base.Type, "svcname", name),
		}

		sv, err := CreateService(base.Type, &ctx)
		if err != nil {
			return fmt.Errorf("cant create service '%s': %v", name, err)
		}

		s.services[name] = sv
	}

	return nil
}
