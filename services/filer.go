package services

import (
	"fmt"
	"io"
	"path"

	"github.com/BurntSushi/toml"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

// FileStorage TBD
type FileStorage interface {
	Save(name string, src io.Reader) (*FileLocation, error)
	GetFile(l *FileLocation) (io.Reader, error)
	GetFileLocation(URL string) (*FileLocation, error)
	Remove(l *FileLocation) error
	IsExist(l *FileLocation) bool
}

// TomlFileStorageContext TBD
type TomlFileStorageContext struct {
	Log    log15.Logger
	Config toml.Primitive
}

// FileLocation TBD
type FileLocation struct {
	BaseURL      string
	BasePath     string
	RelativePath string
	Name         string
}

// FullRelativePath TBD
func (f *FileLocation) FullRelativePath() string {
	return f.FullRelativePathWithPrefix("")
}

// FullRelativePathWithPrefix TBD
func (f *FileLocation) FullRelativePathWithPrefix(pr string) string {
	return path.Join(f.RelativePath, pr+f.Name)
}

// FullURL TBD
func (f *FileLocation) FullURL() string {
	return path.Join(f.BaseURL, f.RelativePath, f.Name)
}

// FullFsPath TBD
func (f *FileLocation) FullFsPath() string {
	return path.Join(f.BasePath, f.RelativePath, f.Name)
}

// NewTomlFilerFunc TBD
type NewTomlFilerFunc func(context *TomlFileStorageContext) (FileStorage, error)

var tomlFilerRegistry = make(map[string]NewTomlFilerFunc)

// RegesterTomlFiler TBD
func RegesterTomlFiler(typo string, factoryMethod NewTomlFilerFunc) {
	tomlFilerRegistry[typo] = factoryMethod
}

// GetTomlFilerBuilder TBD
func GetTomlFilerBuilder(typo string) (NewTomlFilerFunc, error) {
	factoryMethod := tomlFilerRegistry[typo]
	if factoryMethod == nil {
		return nil, fmt.Errorf("Unregestered filer name")
	}
	return factoryMethod, nil
}
