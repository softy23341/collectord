package internal

import (
	"git.softndit.com/collector/backend/cleaver"
	"gopkg.in/gographics/imagick.v3/imagick"
)

const (
	defaultResizeFilter = imagick.FILTER_LANCZOS
	defaultResizeBlur   = 1
)

func init() {
	RegisterImageEngine("imagick", newImagickEngine)
}

func newImagickEngine() ImageEngine {
	return &imagickEngine{}
}

type imagickData struct {
	wand *imagick.MagickWand
}

func (id imagickData) Clone() ImageData {
	if id.wand == nil {
		return imagickData{}
	}
	return imagickData{wand: id.wand.Clone()}
}

func (id imagickData) Destroy() {
	if id.wand == nil {
		return
	}
	id.wand.Destroy()
	id.wand = nil
}

func (id imagickData) Geometry() cleaver.Geometry {
	if id.wand == nil {
		return cleaver.Geometry{}
	}
	return cleaver.Geometry{Width: id.wand.GetImageWidth(), Height: id.wand.GetImageHeight()}
}

func getWand(image ImageData) *imagick.MagickWand {
	if id, ok := image.(imagickData); ok {
		return id.wand
	}
	panic("imagick_engine: invalid ImageData type")
}

type imagickEngine struct{}

func (im *imagickEngine) Load(blob []byte) (ImageData, error) {
	wand := imagick.NewMagickWand()

	if err := wand.ReadImageBlob(blob); err != nil {
		return imagickData{}, err
	}

	return imagickData{wand: wand}, nil
}

func (im *imagickEngine) Save(image ImageData, quality float32, strip bool) ([]byte, error) {
	wand := getWand(image)

	wand.ResetIterator()

	if err := wand.AutoOrientImage(); err != nil {
		return nil, err
	}

	if strip {
		if err := wand.StripImage(); err != nil {
			return nil, err
		}
	}

	if err := wand.SetInterlaceScheme(imagick.INTERLACE_PLANE); err != nil {
		return nil, err
	}

	if quality > 0.0 {
		if err := wand.SetImageCompressionQuality(uint(quality * 100)); err != nil {
			return nil, err
		}
	}

	if err := wand.SetImageFormat("JPEG"); err != nil {
		return nil, err
	}

	blob := wand.GetImageBlob()

	return blob, nil
}

func (im *imagickEngine) Resize(image ImageData, geom cleaver.Geometry) error {
	wand := getWand(image)
	return wand.ResizeImage(geom.Width, geom.Height, defaultResizeFilter)
}

func (im *imagickEngine) Crop(image ImageData, geom cleaver.Geometry, x uint, y uint) error {
	wand := getWand(image)
	return wand.CropImage(geom.Width, geom.Height, int(x), int(y))
}
