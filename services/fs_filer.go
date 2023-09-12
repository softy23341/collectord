package services

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"git.softndit.com/collector/backend/config"
	"github.com/BurntSushi/toml"

	"gopkg.in/inconshreveable/log15.v2"
)

func init() {
	RegesterTomlFiler("fs", NewFSFilerFromToml)
}

const distrRatio = 10

func init() {
	rand.Seed(time.Now().Unix())
}

// FSStorageContext TBD
type FSStorageContext struct {
	Log            log15.Logger
	BaseNamePrefix string
	BasePath       string
	BaseURL        string
}

// NewFSFilerFromToml TBD
func NewFSFilerFromToml(context *TomlFileStorageContext) (FileStorage, error) {
	filerConfig := &config.FileServiceConfig{}
	if err := toml.PrimitiveDecode(context.Config, filerConfig); err != nil {
		return nil, err
	}

	return &FSFiler{
		log:            context.Log,
		baseNamePrefix: filerConfig.BaseNamePrefix,
		basePath:       filerConfig.BasePath,
		baseURL:        filerConfig.BaseURL,
	}, nil
}

// NewFSFiler TBD
func NewFSFiler(context *FSStorageContext) (*FSFiler, error) {
	return &FSFiler{
		log:            context.Log,
		baseNamePrefix: context.BaseNamePrefix,
		basePath:       context.BasePath,
		baseURL:        context.BaseURL,
	}, nil
}

// FSFiler TBD
type FSFiler struct {
	baseNamePrefix string
	basePath       string
	baseURL        string
	log            log15.Logger
}

// PathFor TBD
func (f *FSFiler) PathFor(name string) string {
	date := time.Now()
	return path.Join(
		f.baseNamePrefix,
		fmt.Sprintf("%d", date.Year()),
		fmt.Sprintf("%d", date.Month()),
		fmt.Sprintf("%d", date.Day()),
		fmt.Sprintf("%d", rand.Int()%distrRatio),
	)
}

// Save TBD
func (f *FSFiler) Save(name string, src io.Reader) (*FileLocation, error) {
	dirPath := f.PathFor(name)
	osPath := path.Join(f.basePath, dirPath)
	fileRelativPath := path.Join(dirPath, name)

	if err := os.MkdirAll(osPath, 0755); err != nil {
		return nil, fmt.Errorf("can't save file: %s", err.Error())
	}
	filePath := path.Join(f.basePath, fileRelativPath)

	out, err := os.Create(filePath)
	if err != nil {
		f.log.Error("open file", "err", err)
		return nil, err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	if err != nil {
		f.log.Error("copy error", "err", err)
		return nil, err
	}

	return &FileLocation{
		BasePath: f.basePath,
		BaseURL:  f.baseURL,

		RelativePath: dirPath,
		Name:         name,
	}, nil
}

// GetFileLocation TBD
func (f *FSFiler) GetFileLocation(URL string) (*FileLocation, error) {
	sfFilePath, err := f.remapURLToRelative(URL)
	if err != nil {
		return nil, err
	}

	dir, fileName := filepath.Split(sfFilePath)

	return &FileLocation{
		BasePath: f.basePath,
		BaseURL:  f.baseURL,

		RelativePath: dir,
		Name:         fileName,
	}, nil
}

// Remove TBD
func (f *FSFiler) Remove(l *FileLocation) error {
	return os.Remove(l.FullFsPath())
}

// IsExist TBD
func (f *FSFiler) IsExist(l *FileLocation) bool {
	_, err := os.Stat(l.FullFsPath())
	return !os.IsNotExist(err)
}

func (f *FSFiler) remapURLToRelative(URL string) (string, error) {
	// XXX
	if !strings.HasPrefix(URL, f.baseURL) {
		return "", fmt.Errorf("cant remap %s", URL)
	}
	relativeFilePath := strings.Replace(URL, f.baseURL, "", 1)

	return path.Join(relativeFilePath), nil
}

// GetFile TBD
func (f *FSFiler) GetFile(l *FileLocation) (io.Reader, error) {
	fileData, err := os.ReadFile(l.FullFsPath())
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(fileData)

	return r, nil
}

var _ FileStorage = (*FSFiler)(nil)
