package internal

import (
	"bytes"
	"net/url"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
)

func init() {
	RegisterGetter("swift", NewSwiftGetter)
	RegisterPutter("swift", NewSwiftPutter)
}

type swiftConfig struct {
	Container      string
	Region         string
	FileBasePrefix string

	// gophercloud.AuthOptions
	Version             string
	IdentityEndpointURL string

	Username string
	UserID   string

	Password string
	APIKey   string

	DomainID   string
	DomainName string

	TenantID   string
	TenantName string

	AllowReauth bool
	TokenID     string
}

type swiftBase struct {
	container      string
	fileBasePrefix *url.URL
	client         *gophercloud.ServiceClient
}

func (s *swiftBase) addBasePrefix(dst url.URL) string {
	return path.Join(s.fileBasePrefix.Path, dst.Path)
}

type swiftGetter struct {
	*swiftBase
}

func newSwiftBase(config swiftConfig) (*swiftBase, error) {
	identityEndpointURL, err := url.Parse(config.IdentityEndpointURL)
	if err != nil {
		return nil, err
	}

	cfgFileBasePrefix := "/"
	if config.FileBasePrefix != "" {
		cfgFileBasePrefix = config.FileBasePrefix
	}
	fileBasePrefix, err := url.Parse(cfgFileBasePrefix)
	if err != nil {
		return nil, err
	}

	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: identityEndpointURL.String(),

		Username: config.Username,
		UserID:   config.UserID,

		Password: config.Password,
		APIKey:   config.APIKey,

		DomainID:   config.DomainID,
		DomainName: config.DomainName,

		TenantID:   config.TenantID,
		TenantName: config.TenantName,

		AllowReauth: config.AllowReauth,
	}
	//TokenID:     config.TokenID, removed from authOpts on 08/30

	provider, err := openstack.NewClient(identityEndpointURL.String())
	if err != nil {
		return nil, err
	}

	if config.Version == "v2" {
		if err := openstack.AuthenticateV2(provider, authOpts); err != nil {
			return nil, err
		}
	} else if config.Version == "v3" {
		if err := openstack.AuthenticateV3(provider, authOpts); err != nil {
			return nil, err
		}
	} else {
		if err := openstack.Authenticate(provider, authOpts); err != nil {
			return nil, err
		}
	}

	client, err := openstack.NewObjectStorageV1(provider, gophercloud.EndpointOpts{
		Region: config.Region,
	})
	if err != nil {
		return nil, err
	}

	// check for container existence
	_, err = containers.Get(client, config.Container).ExtractMetadata() //i changed on 08/30
	if err != nil {
		return nil, err
	}

	return &swiftBase{
		client:         client,
		container:      config.Container,
		fileBasePrefix: fileBasePrefix,
	}, err
}

// NewSwiftGetter TBD
func NewSwiftGetter(rawConfig *toml.Primitive) (ImageGetter, error) {
	var config swiftConfig

	if err := toml.PrimitiveDecode(*rawConfig, &config); err != nil {
		return nil, err
	}

	swiftBase, err := newSwiftBase(config)
	if err != nil {
		return nil, err
	}

	return &swiftGetter{swiftBase: swiftBase}, nil
}

func (s *swiftGetter) Get(src url.URL) ([]byte, error) {
	opts := objects.DownloadOpts{}

	res := objects.Download(s.client, s.container, s.addBasePrefix(src), opts)

	return res.ExtractContent()
}

type swiftPutter struct {
	*swiftBase
}

// NewSwiftPutter TBD
func NewSwiftPutter(rawConfig *toml.Primitive) (ImagePutter, error) {
	var config swiftConfig

	if err := toml.PrimitiveDecode(*rawConfig, &config); err != nil {
		return nil, err
	}

	swiftBase, err := newSwiftBase(config)
	if err != nil {
		return nil, err
	}

	return &swiftPutter{swiftBase: swiftBase}, nil
}

func (s *swiftPutter) Put(dst url.URL, blob []byte) error {
	opts := objects.CreateOpts{
		ContentType: "image/jpeg",
	}

	res := objects.Create(s.client, s.container, s.addBasePrefix(dst), bytes.NewReader(blob), opts)
	_, err := res.ExtractHeader()

	return err
}
