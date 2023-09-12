package internal

import (
	"fmt"

	"git.softndit.com/collector/backend/cleaver"
)

// ImageData TBD
type ImageData interface {
	Geometry() cleaver.Geometry
	Clone() ImageData
	Destroy()
}

// ImageEngine TBD
type ImageEngine interface {
	Load(blob []byte) (ImageData, error)
	Save(image ImageData, quality float32, strip bool) ([]byte, error)
	Resize(image ImageData, geom cleaver.Geometry) error
	Crop(image ImageData, geom cleaver.Geometry, x uint, y uint) error
}

// NewImageEngineFunc TBD
type NewImageEngineFunc func() ImageEngine

var imageEngineRegistry = make(map[string]NewImageEngineFunc)

// RegisterImageEngine TBD
func RegisterImageEngine(name string, factory NewImageEngineFunc) {
	imageEngineRegistry[name] = factory
}

// CreateImageEngine TBD
func CreateImageEngine(name string) (ImageEngine, error) {
	if factory, ok := imageEngineRegistry[name]; ok {
		return factory(), nil
	}
	return nil, fmt.Errorf("unregistered image engine '%s'", name)
}
