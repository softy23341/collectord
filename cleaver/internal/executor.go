package internal

import (
	"math"
	"net/url"
	"sync"

	"git.softndit.com/collector/backend/cleaver"
)

// Executor TBD
type Executor struct {
	engine ImageEngine
	getter ImageGetter
	putter ImagePutter
}

// NewExecutor TBD
func NewExecutor(engine ImageEngine, getter ImageGetter, putter ImagePutter) *Executor {
	return &Executor{
		engine: engine,
		getter: getter,
		putter: putter,
	}
}

func (ex *Executor) resizeTransform(imageBlob []byte, image ImageData, tr *cleaver.Transform, dst *url.URL) (*cleaver.TransformResult, error) {
	var (
		originalGeometry = image.Geometry()
		newGeometry      = ex.calculateGeometry(originalGeometry, tr.Geometry, tr.Fit)
		upscale          = !originalGeometry.Contains(newGeometry)
	)

	// no upscale, just skip
	if upscale && tr.UpscaleMode == cleaver.UpscaleDisabled {
		return nil, nil
	}

	if tr.SkipIfLessThan != nil {
		w, h := tr.SkipIfLessThan.Width, tr.SkipIfLessThan.Height
		if (w != 0 && newGeometry.Width < w) || (h != 0 && newGeometry.Height < h) {
			return nil, nil
		}
	}

	doCopy := (tr.Geometry.Empty()) ||
		(tr.CopyEqual && newGeometry == originalGeometry) ||
		(upscale && tr.UpscaleMode == cleaver.UpscaleCopy)

	var newBlob []byte
	if doCopy {
		newBlob = make([]byte, len(imageBlob))
		copy(newBlob, imageBlob)
	} else {
		if err := ex.resizeImage(image, tr.Geometry, tr.Fit); err != nil {
			return nil, cleaver.ErrCategoryTransform.WrapRaw(err)
		}
		blob, err := ex.engine.Save(image, tr.Quality, !tr.NoStrip)
		if err != nil {
			return nil, cleaver.ErrCategoryEncode.WrapRaw(err)
		}
		newBlob = blob
	}

	if err := ex.putter.Put(*dst, newBlob); err != nil {
		return nil, cleaver.ErrCategoryPut.WrapRaw(err)
	}

	return &cleaver.TransformResult{
		Target:           tr.Target,
		OriginalGeometry: originalGeometry,
		OriginalByteSize: uint32(len(imageBlob)),
		NewGeometry:      image.Geometry(),
		NewByteSize:      uint32(len(newBlob)),
	}, nil
}

// Resize TBD
func (ex *Executor) Resize(task *cleaver.ResizeTask) ([]*cleaver.TransformResult, error) {
	if len(task.Transforms) == 0 {
		return nil, nil
	}

	srcURL, err := url.Parse(task.Source)
	if err != nil {
		return nil, cleaver.ErrCategoryInvalidTask.WrapRaw(err)
	}

	var targetURLs = make([]*url.URL, len(task.Transforms))
	for idx, tr := range task.Transforms {
		u, err := url.Parse(tr.Target)
		if err != nil {
			return nil, cleaver.ErrCategoryInvalidTask.WrapRaw(err)
		}
		targetURLs[idx] = u
	}

	var fallbackURL *url.URL
	if task.Fallback != nil {
		u, err := url.Parse(task.Fallback.Target)
		if err != nil {
			return nil, cleaver.ErrCategoryInvalidTask.WrapRaw(err)
		}
		fallbackURL = u
	}

	srcBlob, err := ex.getter.Get(*srcURL)
	if err != nil {
		return nil, cleaver.ErrCategoryGet.WrapRaw(err)
	}

	srcImage, err := ex.engine.Load(srcBlob)
	if err != nil {
		return nil, cleaver.ErrCategoryDecode.WrapRaw(err)
	}
	defer srcImage.Destroy()

	var (
		wg       = new(sync.WaitGroup)
		mu       = new(sync.Mutex) // protect firstErr and list
		firstErr error
		list     []*cleaver.TransformResult
	)
	for idx, transform := range task.Transforms {
		wg.Add(1)
		go func(imageBlob []byte, image ImageData, tr *cleaver.Transform, dst *url.URL) {
			defer image.Destroy()
			defer wg.Done()
			res, err := ex.resizeTransform(imageBlob, image, tr, dst)
			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				if res != nil {
					list = append(list, res)
				}
			} else {
				if firstErr == nil {
					firstErr = err
				}
			}
		}(srcBlob, srcImage.Clone(), transform, targetURLs[idx])
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	if len(list) == 0 && task.Fallback != nil {
		res, err := ex.resizeTransform(srcBlob, srcImage, task.Fallback, fallbackURL)
		if err != nil {
			return nil, err
		}
		if res != nil {
			list = append(list, res)
		}
	}

	return list, nil
}

// Copy TBD
func (ex *Executor) Copy(task *cleaver.CopyTask) (*cleaver.CopyResult, error) {
	srcURL, err := url.Parse(task.Source)
	if err != nil {
		return nil, err
	}

	dstURL, err := url.Parse(task.Target)
	if err != nil {
		return nil, err
	}

	data, err := ex.getter.Get(*srcURL)
	if err != nil {
		return nil, err
	}

	if err := ex.putter.Put(*dstURL, data); err != nil {
		return nil, err
	}

	return &cleaver.CopyResult{ByteSize: int32(len(data)), Target: task.Target}, nil
}

func (ex *Executor) calculateGeometry(src, dst cleaver.Geometry, fit bool) cleaver.Geometry {
	var aspectRatio = src.AspectRatio()
	if dst.Width == 0 {
		dst.Width = uint(math.Max(1.0, float64(dst.Height)*aspectRatio))
	} else if dst.Height == 0 {
		dst.Height = uint(math.Max(1.0, float64(dst.Width)/aspectRatio))
	} else if fit {
		if aspectRatio > dst.AspectRatio() {
			dst.Height = uint(math.Max(1.0, float64(dst.Width)/aspectRatio))
		} else {
			dst.Width = uint(math.Max(1.0, float64(dst.Height)*aspectRatio))
		}
	}
	return dst
}

func (ex *Executor) resizeImage(image ImageData, dst cleaver.Geometry, fit bool) error {
	if dst.Height == 0 || dst.Width == 0 || fit {
		return ex.resizeImageToFit(image, dst)
	}
	return ex.resizeImageToFill(image, dst)
}

func (ex *Executor) resizeImageToFit(image ImageData, dst cleaver.Geometry) error {
	src := image.Geometry()
	dst = ex.calculateGeometry(src, dst, true)

	if dst == src {
		return nil
	}

	if err := ex.engine.Resize(image, dst); err != nil {
		return err
	}

	return nil
}

func (ex *Executor) cropImageCenter(image ImageData, dst cleaver.Geometry) error {
	src := image.Geometry()

	if dst.Width > src.Width {
		dst.Width = src.Width
	}

	if dst.Height > src.Height {
		dst.Height = src.Height
	}

	if dst == src {
		return nil
	}

	x := (src.Width - dst.Width) / 2
	y := (src.Height - dst.Height) / 2

	return ex.engine.Crop(image, dst, x, y)
}

func (ex *Executor) resizeImageToFill(image ImageData, dst cleaver.Geometry) error {
	src := image.Geometry()
	if dst == src {
		return nil
	}

	resizeGeom := dst
	if src.AspectRatio() < dst.AspectRatio() {
		resizeGeom.Height = 0
	} else {
		resizeGeom.Width = 0
	}

	if err := ex.resizeImageToFit(image, resizeGeom); err != nil {
		return err
	}

	return ex.cropImageCenter(image, dst)
}
