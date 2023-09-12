package cleaver

import "git.softndit.com/collector/backend/erry"

var (
	ErrCleaver = erry.NewCategory()

	ErrCategoryInvalidTask = ErrCleaver.NewSubCategory()
	ErrCategoryGetNotFound = ErrCleaver.NewSubCategory()
	ErrCategoryGet         = ErrCleaver.NewSubCategory()
	ErrCategoryDecode      = ErrCleaver.NewSubCategory()
	ErrCategoryTransform   = ErrCleaver.NewSubCategory()
	ErrCategoryEncode      = ErrCleaver.NewSubCategory()
	ErrCategoryPut         = ErrCleaver.NewSubCategory()
	ErrCategoryOther       = ErrCleaver.NewSubCategory()
)

type UpscaleMode int

const (
	UpscaleDisabled UpscaleMode = iota
	UpscaleEnabled
	UpscaleCopy
)

// Transform TBD
type Transform struct {
	Target         string      `json:"target"`
	Geometry       Geometry    `json:"geometry"`
	Quality        float32     `json:"quality,omitempty"`
	Fit            bool        `json:"fit,omitempty"`
	NoStrip        bool        `json:"no_strip,omitempty"`
	CopyEqual      bool        `json:"copy_equal,omitempty"`
	UpscaleMode    UpscaleMode `json:"upscale_mode,omitempty"`
	SkipIfLessThan *Geometry   `json:"skip_if_less_than,omitempty"`
}

// ResizeTask TBD
type ResizeTask struct {
	Source     string       `json:"source"`
	Transforms []*Transform `json:"transforms,omitempty"`
	Fallback   *Transform   `json:"fallback,omitempty"`
}

// TransformResult TBD
type TransformResult struct {
	Target           string   `json:"target"`
	OriginalGeometry Geometry `json:"original_geometry"`
	OriginalByteSize uint32   `json:"original_byte_size,omitempty"`
	NewGeometry      Geometry `json:"new_geometry"`
	NewByteSize      uint32   `json:"new_byte_size,omitempty"`
}

// CopyTask TBD
type CopyTask struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// CopyResult TBD
type CopyResult struct {
	Target   string `json:"target"`
	ByteSize int32  `json:"byte_size"`
}

// Geometry specifies image size
type Geometry struct {
	Width  uint `json:"width,omitempty"`
	Height uint `json:"height,omitempty"`
}

// AspectRatio TBD
func (g Geometry) AspectRatio() float64 {
	return float64(g.Width) / float64(g.Height)
}

// Empty TBD
func (g Geometry) Empty() bool {
	return g.Width == 0 && g.Height == 0
}

// Contains TBD
func (g Geometry) Contains(other Geometry) bool {
	return g.Width >= other.Width && g.Height >= other.Height
}

type StatusCode int16

const (
	StatusOK          StatusCode = 0
	StatusInvalidTask StatusCode = 1
	StatusGetNotFound StatusCode = 2
	StatusGetError    StatusCode = 3
	StatusDecodeError StatusCode = 4
	StatusEncodeError StatusCode = 5
	StatusPutError    StatusCode = 6
	StatusOther       StatusCode = 7
)

var errStatusMap = map[*erry.Category]StatusCode{
	ErrCategoryInvalidTask: StatusInvalidTask,
	ErrCategoryGetNotFound: StatusGetNotFound,
	ErrCategoryGet:         StatusGetError,
	ErrCategoryDecode:      StatusDecodeError,
	ErrCategoryEncode:      StatusEncodeError,
	ErrCategoryPut:         StatusPutError,
	ErrCategoryOther:       StatusOther,
}

var statusErrMap = make(map[StatusCode]*erry.Category)

func init() {
	for err, status := range errStatusMap {
		statusErrMap[status] = err
	}
}

// NewStatusFromErr TBD
func NewStatusFromErr(e error) Status {
	if e == nil {
		return Status{Code: StatusOK}
	}
	var (
		code StatusCode
		ok   bool
	)
	if code, ok = errStatusMap[erry.GetCategory(e)]; !ok {
		code = StatusOther
	}
	return Status{Code: code, Text: e.Error()}
}

// Status TBD
type Status struct {
	Code StatusCode `json:"code"`
	Text string     `json:"text,omitempty"`
}

// ToErr TDB
func (s Status) ToErr() error {
	if s.Code == StatusOK {
		return nil
	}
	var (
		cat *erry.Category
		ok  bool
	)
	if cat, ok = statusErrMap[s.Code]; !ok {
		cat = ErrCategoryOther
	}
	return cat.New(s.Text)
}

// StatusedResizeResult TBD
type StatusedResizeResult struct {
	Status     Status             `json:"status"`
	Transforms []*TransformResult `json:"results,omitempty"`
}

// StatusedCopyResult TBD
type StatusedCopyResult struct {
	Status Status      `json:"status"`
	Result *CopyResult `json:"result,omitempty"`
}

const (
	// RMQQueueResizeJSON TBD
	RMQQueueResizeJSON = "cl-resize.json"
	// RMQQueueCopyJSON TBD
	RMQQueueCopyJSON = "cl-copy.json"
)
