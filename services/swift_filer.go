package services

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"git.softndit.com/collector/backend/config"

	"github.com/BurntSushi/toml"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"

	log15 "gopkg.in/inconshreveable/log15.v2"
)

func init() {
	RegesterTomlFiler("swift", NewSwiftFilerFromToml)
}

// SwiftStorageContext TBD
type SwiftStorageContext struct {
	Log    log15.Logger
	Config *SwiftConfig
}

// NewSwiftFiler TBD
func NewSwiftFiler(context *SwiftStorageContext) (*SwiftFiler, error) {
	base, err := newSwiftBase(context.Config)
	if err != nil {
		return nil, err
	}
	return &SwiftFiler{
		swiftBase: base,
		Log:       context.Log,
	}, nil
}

// NewSwiftFilerFromToml TBD
func NewSwiftFilerFromToml(context *TomlFileStorageContext) (FileStorage, error) {
	swiftConfig := &config.SwiftFileServiceConfig{}

	//if err := toml.PrimitiveDecode(context.Config, swiftConfig); err != nil {
	if err := toml.PrimitiveDecode(context.Config, swiftConfig); err != nil {
		return nil, err
	}
	serviceSwiftConfig := &SwiftConfig{
		IdentityEndpointURL: swiftConfig.IdentityEndpointURL,
		Container:           swiftConfig.Container,

		Username: swiftConfig.Username,
		Password: swiftConfig.Password,
		Region:   swiftConfig.Region,

		BaseURL:        swiftConfig.BaseURL,
		BaseNamePrefix: swiftConfig.BaseNamePrefix,
	}

	return NewSwiftFiler(&SwiftStorageContext{Log: context.Log, Config: serviceSwiftConfig})
}

// SwiftConfig TBD
type SwiftConfig struct {
	IdentityEndpointURL string
	Container           string

	Username string
	Password string
	Region   string

	BaseURL        string
	BaseNamePrefix string
}

// SwiftFiler TBD
type SwiftFiler struct {
	Log log15.Logger
	*swiftBase
}

func (s *swiftBase) addBaseURL(dst string) string {
	return path.Join(s.baseURL.Path, dst)
}

type swiftBase struct {
	container      string
	baseURL        *url.URL
	BaseNamePrefix string
	client         *gophercloud.ServiceClient
}

func newSwiftBase(config *SwiftConfig) (*swiftBase, error) {
	identityEndpointURL, err := url.Parse(config.IdentityEndpointURL)
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, err
	}

	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: identityEndpointURL.String(),
		Username:         config.Username,
		Password:         config.Password,
		AllowReauth:      true,
	}

	provider, err := openstack.NewClient(identityEndpointURL.String())
	if err != nil {
		return nil, err
	}

	if err := openstack.AuthenticateV2(provider, authOpts); err != nil {
		return nil, err
	}

	client, err := openstack.NewObjectStorageV1(provider, gophercloud.EndpointOpts{
		Region: config.Region,
	})
	if err != nil {
		return nil, err
	}

	// check for container existence
	_, err = containers.Get(client, config.Container).ExtractHeader() //e.m. changed 08/30
	if err != nil {
		return nil, err
	}

	return &swiftBase{
		client:         client,
		container:      config.Container,
		BaseNamePrefix: config.BaseNamePrefix,
		baseURL:        baseURL,
	}, err
}

// PathFor TBD
func (s *SwiftFiler) PathFor(name string) string {
	date := time.Now()
	return path.Join(
		s.BaseNamePrefix,
		fmt.Sprintf("%d", date.Year()),
		fmt.Sprintf("%d", date.Month()),
		fmt.Sprintf("%d", date.Day()),
		fmt.Sprintf("%d", rand.Int()%distrRatio),
	)
}

// Save TBD
func (s *SwiftFiler) Save(name string, src io.Reader) (*FileLocation, error) {
	relativePath := s.PathFor(name)
	fullPath := path.Join(s.baseURL.Path, relativePath, name)

	if err := s.SaveByPath(fullPath, src); err != nil {
		return nil, err
	}
	return &FileLocation{
		BasePath: s.baseURL.Path,
		BaseURL:  s.baseURL.Path,

		RelativePath: relativePath,
		Name:         name,
	}, nil
}

// SaveByPath TBD
func (s *SwiftFiler) SaveByPath(fullPath string, src io.Reader) error {
	opts := objects.CreateOpts{}

	blob, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	res := objects.Create(s.client, s.container, fullPath, bytes.NewReader(blob), opts)
	_, err = res.ExtractHeader()

	if err != nil {
		return err

	}

	return nil
}

// GetFileLocation TBD
func (s *SwiftFiler) GetFileLocation(URL string) (*FileLocation, error) {
	sfFilePath, err := s.remapURLToRelative(URL)
	if err != nil {
		return nil, err
	}

	dir, fileName := filepath.Split(sfFilePath)

	return &FileLocation{
		BasePath: s.baseURL.String(),
		BaseURL:  s.baseURL.String(),

		RelativePath: dir,
		Name:         fileName,
	}, nil
}

func (s *SwiftFiler) remapURLToRelative(URL string) (string, error) {
	// XXX
	if !strings.HasPrefix(URL, s.baseURL.String()) {
		return "", fmt.Errorf("cant remap %s", URL)
	}
	relativeFilePath := strings.Replace(URL, s.baseURL.String(), "", 1)

	return path.Join(relativeFilePath), nil
}

// Remove TBD
func (s *SwiftFiler) Remove(l *FileLocation) error {
	res := objects.Delete(s.client, s.container, l.FullURL(), objects.DeleteOpts{})
	_, err := res.ExtractHeader()
	return err
}

// IsExist TBD
func (s *SwiftFiler) IsExist(l *FileLocation) bool {
	//opts := objects.DownloadOpts{IfMatch: "etag"}

	//result := objects.Download(s.client, s.container, l.FullURL(), opts)
	//_, err := result.ExtractHeader() //e.m. changed on 08/30

	//return err == nil
	return false
}

// GetFile TBD
func (s *SwiftFiler) GetFile(l *FileLocation) (io.Reader, error) {
	return nil, nil
}

var _ FileStorage = (*SwiftFiler)(nil)
