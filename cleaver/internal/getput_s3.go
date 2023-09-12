package internal

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/amz.v3/aws"
	"gopkg.in/amz.v3/s3"
)

func init() {
	RegisterGetter("s3", NewS3Getter)
	RegisterPutter("s3", NewS3Putter)
}

const (
	defaultRegion = "eu-central-1"
	defaultAcl    = "public-read"

	contentType = "image/jpeg"

	envLabel = "@env"

	envClGetAk  = "CL_GET_AWS_ACCESS_KEY_ID"
	envClGetSk  = "CL_GET_AWS_SECRET_KEY"
	envClGetReg = "CL_GET_AWS_DEFAULT_REGION"

	envClPutAk  = "CL_PUT_AWS_ACCESS_KEY_ID"
	envClPutSk  = "CL_PUT_AWS_SECRET_KEY"
	envClPutReg = "CL_PUT_AWS_DEFAULT_REGION"

	envClAk  = "CL_AWS_ACCESS_KEY_ID"
	envClSk  = "CL_AWS_SECRET_KEY"
	envClReg = "CL_AWS_DEFAULT_REGION"

	envAk  = "AWS_ACCESS_KEY_ID"
	envSk  = "AWS_SECRET_KEY"
	envReg = "AWS_DEFAULT_REGION"
)

type getConfig struct {
	Bucket    string
	BaseURL   string
	AccessKey string
	SecretKey string
	Region    string
}

type putConfig struct {
	*getConfig
	ACL string
}

func getEnv(keys []string) string {
	for _, k := range keys {
		v := os.Getenv(k)
		if len(v) > 0 {
			return v
		}
	}

	return ""
}

func fillGetConfig(c *getConfig) error {
	if c.Bucket == "" {
		return errors.New("no bucket specified")
	}

	if c.AccessKey == envLabel {
		env := []string{envClPutAk, envClAk, envAk}
		c.AccessKey = getEnv(env)
	}
	if c.AccessKey == "" {
		return errors.New("no aws access key specified")
	}

	if c.SecretKey == envLabel {
		env := []string{envClPutSk, envClSk, envSk}
		c.SecretKey = getEnv(env)
	}
	if c.SecretKey == "" {
		return errors.New("no aws secret key specified")
	}

	if c.Region == envLabel {
		env := []string{envClPutReg, envClReg, envReg}
		c.Region = getEnv(env)
	}

	if c.Region == "" {
		c.Region = defaultRegion
	}
	if _, ok := aws.Regions[c.Region]; !ok {
		return fmt.Errorf("unknown aws region: '%s'", c.Region)
	}

	return nil
}

func fillPutConfig(c *putConfig) error {
	if err := fillGetConfig(c.getConfig); err != nil {
		return err
	}

	if c.ACL == "" {
		c.ACL = defaultAcl
	}

	return nil
}

func normalizeS3Path(path string) string {
	for {
		path = strings.TrimPrefix(path, "/")
		if !strings.HasPrefix(path, "/") {
			break
		}
	}
	return path
}

type s3Getter struct {
	bucket  *s3.Bucket
	baseURL url.URL
	acl     s3.ACL
}

// NewS3Getter TBD
func NewS3Getter(rawConfig *toml.Primitive) (ImageGetter, error) {
	var config getConfig

	if err := toml.PrimitiveDecode(*rawConfig, &config); err != nil {
		return nil, err
	}

	if err := fillGetConfig(&config); err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, err
	}

	auth := aws.Auth{
		AccessKey: config.AccessKey,
		SecretKey: config.SecretKey,
	}

	bucket, err := s3.New(auth, aws.Regions[config.Region]).Bucket(config.Bucket)
	if err != nil {
		return nil, err
	}

	getter := &s3Getter{
		bucket:  bucket,
		baseURL: *baseURL,
	}

	return getter, nil
}

func (s *s3Getter) Get(src url.URL) ([]byte, error) {
	p := normalizeS3Path(path.Join(s.baseURL.Path, src.Path))
	return s.bucket.Get(p)
}

type s3Putter struct {
	bucket  *s3.Bucket
	baseURL url.URL
	acl     s3.ACL
}

// NewS3Putter TBD
func NewS3Putter(rawConfig *toml.Primitive) (ImagePutter, error) {
	var config putConfig

	if err := toml.PrimitiveDecode(*rawConfig, &config); err != nil {
		return nil, err
	}

	if err := fillPutConfig(&config); err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, err
	}

	auth := aws.Auth{
		AccessKey: config.AccessKey,
		SecretKey: config.SecretKey,
	}

	bucket, err := s3.New(auth, aws.Regions[config.Region]).Bucket(config.Bucket)
	if err != nil {
		return nil, err
	}

	putter := &s3Putter{
		bucket:  bucket,
		baseURL: *baseURL,
		acl:     s3.ACL(config.ACL),
	}

	return putter, nil
}

func (s *s3Putter) Put(dst url.URL, blob []byte) error {
	p := normalizeS3Path(path.Join(s.baseURL.Path, dst.Path))
	return s.bucket.Put(p, blob, contentType, s3.PublicRead)
}
