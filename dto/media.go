package dto

import "encoding/json"

//go:generate dbgen -type Media

// MediaType TBD
type MediaType int16

// IsPhoto TBD
func (m MediaType) IsPhoto() bool {
	return m > 0 && m < 100
}

// IsDocument TBD
func (m MediaType) IsDocument() bool {
	return m > 100 && m < 200
}

// IsDocument TBD
func (m MediaType) IsVideo() bool {
	return m > 200 && m < 300
}

// MediaTypeSet TBD
type MediaTypeSet []MediaType

// Media types
const (
	MediaTypePhoto        MediaType = 1
	MediaTypePhotoObject            = 2
	MediaTypeAvatar                 = 10
	MediaTypeCollection             = 15
	MediaTypeMsg                    = 20
	MediaTypeGenericPhoto           = 30
	MediaTypeDocument               = 101
	MediaTypeVideo                  = 201
	MediaTypeMapPoint               = 301
)

// MediaTypeList TBD
var MediaTypeList = MediaTypeSet{
	MediaTypePhoto,
	MediaTypePhotoObject,
	MediaTypeAvatar,
	MediaTypeCollection,
	MediaTypeMsg,
	MediaTypeGenericPhoto,
	MediaTypeDocument,
	MediaTypeVideo,
	MediaTypeMapPoint,
}

// Contain TBD
func (m MediaTypeSet) Contain(typo MediaType) bool {
	for _, mediaType := range m {
		if mediaType == typo {
			return true
		}
	}
	return false
}

// Media TBD
type Media struct {
	ID         int64     `db:"id"`
	UserID     *int64    `db:"user_id"`
	UserUniqID *int64    `db:"user_uniq_id"`
	Type       MediaType `db:"type"`
	RootID     *int64    `db:"root_id"`
	MediaUnion `db:"extra,json"`
}

// ExtraJSON TBD
func (m *Media) ExtraJSON() []byte {
	data, _ := json.Marshal(m.MediaUnion)
	return data
}

// MediaUnion TBD
type MediaUnion struct {
	Photo    *Photo    `json:"photo,omitempty"`
	Document *Document `json:"document,omitempty"`
}

// MediaList TBD
type MediaList []*Media

// GetIDs ids slice
func (m MediaList) GetIDs() []int64 {
	o := make([]int64, len(m))
	for i := range m {
		o[i] = m[i].ID
	}
	return o
}

// GetOwnersIDs TBD
func (m MediaList) GetOwnersIDs() []int64 {
	usersIDs := make([]int64, 0, len(m))
	for _, media := range m {
		if media.UserID != nil {
			usersIDs = append(usersIDs, *media.UserID)
		}
	}
	return usersIDs
}

// IDToUser TBD
func (m MediaList) IDToMedia() map[int64]*Media {
	id2media := make(map[int64]*Media, 0)
	for _, media := range m {
		id2media[media.ID] = media
	}
	return id2media
}

// Document TBD
type Document struct {
	URI      string `json:"uri"`
	ByteSize int64  `json:"byte_size"`
	MimeType string `json:"mime_type"`
	Name     string `json:"name"`
}

// MediaPhotoType TBD
type MediaPhotoType int16

// photo types
const (
	PhotoTypeObject MediaPhotoType = iota
)

// Photo TBD
type Photo struct {
	Name     string          `json:"name"`
	Variants []*PhotoVariant `json:"variants"`
}

// PhotoVariantType TBD
type PhotoVariantType int16

// Photo variant types
const (
	PhotoVariantOther              PhotoVariantType = 0
	PhotoVariantOriginal                            = 1
	PhotoVariantSmallThumb                          = 3
	PhotoVariantGenericSmallThumb                   = 4
	PhotoVariantGenericMedium                       = 5
	PhotoVariantGenericLarge                        = 6
	PhotoVariantObjectPreview                       = 100
	PhotoVariantObjectGallery                       = 101
	PhotoVariantObjectHD                            = 102
	PhotoVariantObjectPreviewSmall                  = 150
	PhotoVariantObjectGallerySmall                  = 151
	PhotoVariantAvatar                              = 200
	PhotoVariantMsgPreview                          = 300
)

// Original TBD
func (t PhotoVariant) Original() bool {
	return t.Type == PhotoVariantOriginal
}

// PhotoVariantTypeList TBD
type PhotoVariantTypeList []PhotoVariantType

// ToSet TBD
func (l PhotoVariantTypeList) ToSet() (set PhotoVariantTypeSet) {
	set = make(PhotoVariantTypeSet, len(l))
	for _, t := range l {
		set[t] = struct{}{}
	}
	return
}

// PhotoVariantTypeSet TBD
type PhotoVariantTypeSet map[PhotoVariantType]struct{}

// ToList TBD
func (s PhotoVariantTypeSet) ToList() (l PhotoVariantTypeList) {
	l = make(PhotoVariantTypeList, 0, len(s))
	for t := range s {
		l = append(l, t)
	}
	return
}

// PhotoVariant TBD
type PhotoVariant struct {
	URI      string           `json:"uri"`
	Type     PhotoVariantType `json:"type,omitempty"`
	Width    int32            `json:"width"`
	Height   int32            `json:"height"`
	ByteSize int64            `json:"byte_size"`
}
