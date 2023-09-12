package internal

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

const defaultFileMode = 0644

func init() {
	RegisterGetter("fs", NewFileSystemGetter)
	RegisterPutter("fs", NewFileSystemPutter)
}

type fsGetter struct {
	baseURL url.URL
}

// NewFileSystemGetter TBD
func NewFileSystemGetter(config *toml.Primitive) (ImageGetter, error) {
	var data struct {
		BaseURL string
	}
	if err := toml.PrimitiveDecode(*config, &data); err != nil {
		return nil, err
	}

	url, err := url.Parse(data.BaseURL)
	if err != nil {
		return nil, err
	}

	return &fsGetter{*url}, nil
}

func (fg *fsGetter) Get(url url.URL) ([]byte, error) {
	p := path.Join(fg.baseURL.Path, url.Path)
	return os.ReadFile(p)
}

type fileMode uint32

func (fm *fileMode) UnmarshalText(text []byte) error {
	str := string(text)
	_, err := fmt.Sscanf(str, "%o", fm)
	return err
}

func (fm *fileMode) String() string {
	return fmt.Sprintf("0%o", *fm)
}

type fsPutter struct {
	baseURL url.URL
	mode    fileMode
}

// NewFileSystemPutter TBD
func NewFileSystemPutter(config *toml.Primitive) (ImagePutter, error) {
	var data struct {
		BaseURL  string
		FileMode fileMode
	}
	if err := toml.PrimitiveDecode(*config, &data); err != nil {
		return nil, err
	}

	mode := data.FileMode
	if mode == 0 {
		mode = defaultFileMode
	}

	url, err := url.Parse(data.BaseURL)
	if err != nil {
		return nil, err
	}

	return &fsPutter{*url, mode}, nil
}

func (fp *fsPutter) Put(url url.URL, blob []byte) error {
	p := path.Join(fp.baseURL.Path, url.Path)
	return os.WriteFile(p, blob, os.FileMode(fp.mode))
}
