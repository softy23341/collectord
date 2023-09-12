package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx"
)

// EventServiceConfig TBD
type EventServiceConfig struct {
	Port int `toml:"ws_port"`
}

// CleaverClientConfig TBD
type CleaverClientConfig struct {
	RabbitRef string `toml:"rabbit_ref"`
}

// SearchClientConfig TBD
type SearchClientConfig struct {
	ElasticSearchRef string `toml:"elastic_search_ref"`
	ObjectIndex      string `toml:"object_index"`
}

// PusherClientConfig TBD
type PusherClientConfig struct {
	RabbitRef string `toml:"rabbit_ref"`
}

// DBMConfig TBD
type DBMConfig struct {
	PgRef string `toml:"pg_ref"`
}

// FileServiceConfig TBD
type FileServiceConfig struct {
	BaseNamePrefix string `toml:"base_name_prefix"`
	BasePath       string `toml:"base_path"`
	BaseURL        string `toml:"base_URL"`
}

// SwiftFileServiceConfig TBD
type SwiftFileServiceConfig struct {
	IdentityEndpointURL string `toml:"identity_endpoint_URL"`
	Container           string `toml:"container"`
	Username            string `toml:"username"`
	Password            string `toml:"password"`
	BaseURL             string `toml:"base_URL"`
	Region              string `toml:"region"`
	BaseNamePrefix      string `toml:"base_name_prefix"`
}

// FilerGenericConfig TBD
type FilerGenericConfig struct {
	Type   string
	Config toml.Primitive
}

// FilerClientConfig TBD
type FilerClientConfig struct {
	FilerRef string `toml:"filer_ref"`
}

// MailClientConfig TBD
type MailClientConfig struct {
	Password   string `toml:"password"`
	Username   string `toml:"username"`
	ServerName string `toml:"server_name"`
	Port       int    `toml:"port"`
}

// RedisClientConfig TBD
type RedisConfig struct {
	Url string `toml:"url"`
}

// ServerConfig TBD
type ServerConfig struct {
	Port         int    `toml:"port"`
	I18nPath     string `toml:"i18n_path"`
	TemplatePath string `toml:"template_path"`

	DBM           DBMConfig           `toml:"DBM"`
	Redis         RedisConfig         `toml:"redis"`
	EventService  EventServiceConfig  `toml:"event_service"`
	CleaverClient CleaverClientConfig `toml:"cleaver_client"`
	SearchClient  SearchClientConfig  `toml:"search_client"`
	PusherClient  PusherClientConfig  `toml:"pusher_client"`
	FilerClient   FilerClientConfig   `toml:"filer_client"`
	MailClient    MailClientConfig    `toml:"mail_client"`
}

// Config TBD
type Config struct {
	Rabbit        map[string]toml.Primitive `toml:"rabbit"`
	Pg            map[string]toml.Primitive `toml:"pg"`
	ElasticSearch map[string]toml.Primitive `toml:"elastic_search"`
	Filer         map[string]toml.Primitive `toml:"filer"`

	ServerConfig ServerConfig `toml:"server"`
}

// ReadConfig TBD
func (c *Config) ReadConfig(configPath string) error {
	_, err := toml.DecodeFile(configPath, c)
	return err
}

// RabbitConfig TBD
type RabbitConfig struct {
	URL string `toml:"url"`
}

// ElasticSearchConfig TBD
type ElasticSearchConfig struct {
	URL string `toml:"url"`
}

// PgConfig TBD
type PgConfig struct {
	Host           string
	User           string
	Password       string
	DB             string `toml:"db"`
	MaxConnections int    `toml:"max_connections"`
}

// NewInstanceHolder TBD
func NewInstanceHolder(c *Config) (*InstanceHolder, error) {
	i := &InstanceHolder{}
	err := i.Configure(c)
	return i, err
}

// InstanceHolder TBD
type InstanceHolder struct {
	RabbitConfig        map[string]*RabbitConfig
	ElasticSearchConfig map[string]*ElasticSearchConfig
	PgConfig            map[string]*pgx.ConnPoolConfig
	FilerConfig         map[string]*FilerGenericConfig
}

// GetPgConfig TBD
func (i *InstanceHolder) GetPgConfig(ref string) (*pgx.ConnPoolConfig, error) {
	c, find := i.PgConfig[ref]
	if !find {
		return c, fmt.Errorf("cant find pg config by key: %s", ref)
	}
	return c, nil
}

// GetRabbitConfig TBD
func (i *InstanceHolder) GetRabbitConfig(ref string) (*RabbitConfig, error) {
	c, find := i.RabbitConfig[ref]
	if !find {
		return c, fmt.Errorf("cant find rabbit config by key: %s", ref)
	}
	return c, nil
}

// GetElasticSearchConfig TBD
func (i *InstanceHolder) GetElasticSearchConfig(ref string) (*ElasticSearchConfig, error) {
	c, find := i.ElasticSearchConfig[ref]
	if !find {
		return c, fmt.Errorf("cant find elastic config by key: %s", ref)
	}
	return c, nil
}

// GetGenericFilerConfig TBD
func (i *InstanceHolder) GetGenericFilerConfig(ref string) (*FilerGenericConfig, error) {
	c, found := i.FilerConfig[ref]
	if !found {
		return nil, fmt.Errorf("cant find filer config by key: %s", ref)
	}
	return c, nil
}

// Configure TBD
func (i *InstanceHolder) Configure(c *Config) error {
	if err := i.ConfigureRabbits(c); err != nil {
		return err
	}

	if err := i.ConfigureElasticSearch(c); err != nil {
		return err
	}

	if err := i.ConfigurePg(c); err != nil {
		return err
	}

	if err := i.ConfigureFiler(c); err != nil {
		return err
	}

	return nil
}

// ConfigureFiler TBD
func (i *InstanceHolder) ConfigureFiler(c *Config) error {
	i.FilerConfig = make(map[string]*FilerGenericConfig, len(c.Filer))

	for refName, tomlConfig := range c.Filer {
		typeCfg := &struct {
			Type string `toml:"type"`
		}{}
		if err := toml.PrimitiveDecode(tomlConfig, typeCfg); err != nil {
			return err
		}
		i.FilerConfig[refName] = &FilerGenericConfig{
			Type:   typeCfg.Type,
			Config: tomlConfig,
		}
	}

	return nil
}

// ConfigureRabbits TBD
func (i *InstanceHolder) ConfigureRabbits(c *Config) error {
	i.RabbitConfig = make(map[string]*RabbitConfig, len(c.Rabbit))

	for refName, tomlConfig := range c.Rabbit {
		rabbitCfg := &RabbitConfig{}
		toml.PrimitiveDecode(tomlConfig, rabbitCfg)
		i.RabbitConfig[refName] = rabbitCfg
	}
	return nil
}

// ConfigureElasticSearch TBD
func (i *InstanceHolder) ConfigureElasticSearch(c *Config) error {
	i.ElasticSearchConfig = make(map[string]*ElasticSearchConfig, len(c.ElasticSearch))

	for refName, tomlConfig := range c.ElasticSearch {
		elasticSearchConfig := &ElasticSearchConfig{}
		toml.PrimitiveDecode(tomlConfig, elasticSearchConfig)
		i.ElasticSearchConfig[refName] = elasticSearchConfig
	}
	return nil
}

// ConfigurePg TBD
func (i *InstanceHolder) ConfigurePg(c *Config) error {
	i.PgConfig = make(map[string]*pgx.ConnPoolConfig, len(c.Pg))

	for refName, tomlConfig := range c.Pg {
		pgConfig := &PgConfig{}
		toml.PrimitiveDecode(tomlConfig, pgConfig)
		i.PgConfig[refName] = &pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     pgConfig.Host,
				Database: pgConfig.DB,
				User:     pgConfig.User,
				Password: pgConfig.Password,
			},
			MaxConnections: pgConfig.MaxConnections,
		}

	}

	return nil
}
