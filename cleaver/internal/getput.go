package internal

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/BurntSushi/toml"
)

var (
	errNotFound = errors.New("not found")
)

// ImageGetter TODO
type ImageGetter interface {
	Get(url.URL) ([]byte, error)
}

// ImagePutter TODO
type ImagePutter interface {
	Put(url.URL, []byte) error
}

// NewGetterFunc TODO
type NewGetterFunc func(config *toml.Primitive) (ImageGetter, error)

// NewPutterFunc TODO
type NewPutterFunc func(config *toml.Primitive) (ImagePutter, error)

var (
	getterRegistry = make(map[string]NewGetterFunc)
	putterRegistry = make(map[string]NewPutterFunc)
)

// RegisterGetter TODO
func RegisterGetter(typ string, factory NewGetterFunc) {
	getterRegistry[typ] = factory
}

// RegisterPutter TODO
func RegisterPutter(typ string, factory NewPutterFunc) {
	putterRegistry[typ] = factory
}

// CreateGetter TODO
func CreateGetter(typ string, config *toml.Primitive) (ImageGetter, error) {
	factory := getterRegistry[typ]
	if factory == nil {
		return nil, fmt.Errorf("unregistered getter type '%s'", typ)
	}
	return factory(config)
}

// CreatePutter TODO
func CreatePutter(typ string, config *toml.Primitive) (ImagePutter, error) {
	factory := putterRegistry[typ]
	if factory == nil {
		return nil, fmt.Errorf("unregistered putter type '%s'", typ)
	}
	return factory(config)
}

type schemeGetterPutter struct {
	getters map[string]ImageGetter
	putters map[string]ImagePutter
}

func newSchemeGetterPutter() *schemeGetterPutter {
	return &schemeGetterPutter{
		getters: make(map[string]ImageGetter),
		putters: make(map[string]ImagePutter),
	}
}

func (sgp *schemeGetterPutter) addGetter(scheme string, getter ImageGetter) {
	sgp.getters[scheme] = getter
}

func (sgp *schemeGetterPutter) addPutter(scheme string, putter ImagePutter) {
	sgp.putters[scheme] = putter
}

func (sgp *schemeGetterPutter) Get(src url.URL) ([]byte, error) {
	scheme := src.Scheme
	getter, ok := sgp.getters[scheme]
	if !ok {
		return nil, fmt.Errorf("unknown getter scheme:'%s'", scheme)
	}
	return getter.Get(src)
}

func (sgp *schemeGetterPutter) Put(dst url.URL, blob []byte) error {
	scheme := dst.Scheme
	putter, ok := sgp.putters[scheme]
	if !ok {
		return fmt.Errorf("unknown putter scheme:'%s'", scheme)
	}
	return putter.Put(dst, blob)
}
