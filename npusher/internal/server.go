package internal

import (
	"fmt"

	"github.com/inconshreveable/log15"

	"github.com/BurntSushi/toml"
)

// Server TBD
type Server struct {
	provider Provider
	services []Service
}

// ServerConfig TBD
type ServerConfig struct {
	Log      log15.Logger
	Provider map[string]toml.Primitive
	Services map[string]toml.Primitive
}

// ReadFromToml TBD
func (config *ServerConfig) ReadFromToml(path string) error {
	_, err := toml.DecodeFile(path, &config)
	return err
}

// NewServer TBD
func NewServer(config *ServerConfig) (*Server, error) {
	server := &Server{}

	if err := server.configureProvider(config); err != nil {
		return nil, err
	}
	if err := server.configureServices(config); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) configureProvider(config *ServerConfig) error {
	if len(config.Provider) != 1 {
		return fmt.Errorf("no providers specified")
	}
	for name, options := range config.Provider {
		providerFactory, err := GetProviderFactory(name)
		if err != nil {
			return err
		}
		provider, err := providerFactory(&ProviderCtx{
			Config: &options,
			Log:    config.Log.New("provider", name),
		})
		if err != nil {
			return err
		}
		s.provider = provider
	}
	return nil
}

func (s *Server) configureServices(config *ServerConfig) error {
	if len(config.Services) < 1 {
		return fmt.Errorf("no services specified")
	}
	if s.provider == nil {
		return fmt.Errorf("no provider specified")
	}

	for serviceName, serviceConfig := range config.Services {
		builder, err := GetServiceBuilder(serviceName)
		if err != nil {
			return err
		}
		context := &ServiceCtx{
			Config:   &serviceConfig,
			Log:      config.Log.New("service", serviceName),
			Provider: s.provider,
		}
		readyService, err := builder(context)
		if err != nil {
			return err
		}
		s.services = append(s.services, readyService)
	}

	if len(s.services) < 1 {
		return fmt.Errorf("no services was init")
	}
	return nil
}

// Run TBD
func (s *Server) Run() error {
	for _, readyService := range s.services {
		if err := readyService.Run(); err == nil {
			return err
		}
	}
	return nil
}
