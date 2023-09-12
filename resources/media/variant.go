package mediavariant

import (
	"git.softndit.com/collector/backend/cleaver"
	"git.softndit.com/collector/backend/dto"
)

// MediaTransform TBD
type MediaTransform struct {
	Type       dto.MediaType
	Transforms []*TypedTransform
}

// TypedTransform TBD
type TypedTransform struct {
	Type      dto.PhotoVariantType
	Transform *cleaver.Transform
}

// TransformsByMediaType TBD
func TransformsByMediaType(typo dto.MediaType) ([]*TypedTransform, error) {
	for _, mTransform := range MediaTransformsList() {
		if typo == mTransform.Type {
			return mTransform.Transforms, nil
		}
	}

	return nil, nil
}

// MediaTransformsList TBD
func MediaTransformsList() []*MediaTransform {
	return []*MediaTransform{
		&MediaTransform{
			Type:       dto.MediaTypePhoto,
			Transforms: []*TypedTransform{},
		},
		PhotoObjectTransform(),
		AvatarTransform(),
		CollectionTransform(),
		MessengerMediaTransform(),
	}
}

// MessengerMediaTransform TBD
func MessengerMediaTransform() *MediaTransform {
	return &MediaTransform{
		Type: dto.MediaTypeMsg,
		Transforms: []*TypedTransform{
			&TypedTransform{
				Type: dto.PhotoVariantMsgPreview,
				Transform: &cleaver.Transform{
					Target:    "var_400x400_",
					Geometry:  cleaver.Geometry{Width: 400, Height: 400},
					Quality:   0.70,
					Fit:       true,
					CopyEqual: true,
				},
			},
		},
	}
}

// CollectionTransform TBD
func CollectionTransform() *MediaTransform {
	return &MediaTransform{
		Type: dto.MediaTypeCollection,
		Transforms: []*TypedTransform{
			&TypedTransform{
				Type: dto.PhotoVariantSmallThumb,
				Transform: &cleaver.Transform{
					Target:    "var_210x210_",
					Geometry:  cleaver.Geometry{Width: 210, Height: 210},
					Quality:   0.70,
					CopyEqual: true,
				},
			},
		},
	}
}

// PhotoObjectTransform TBD
func PhotoObjectTransform() *MediaTransform {
	return &MediaTransform{
		Type: dto.MediaTypePhotoObject,
		Transforms: []*TypedTransform{
			&TypedTransform{
				Type: dto.PhotoVariantObjectPreview,
				Transform: &cleaver.Transform{
					Target:    "var_460x780_",
					Geometry:  cleaver.Geometry{Width: 460, Height: 780},
					Quality:   0.70,
					Fit:       true,
					CopyEqual: true,
				},
			},
			&TypedTransform{
				Type: dto.PhotoVariantObjectGallery,
				Transform: &cleaver.Transform{
					Target:    "var_1360x_",
					Geometry:  cleaver.Geometry{Width: 1360},
					Quality:   0.70,
					Fit:       true,
					CopyEqual: true,
				},
			},
			&TypedTransform{
				Type: dto.PhotoVariantObjectPreviewSmall,
				Transform: &cleaver.Transform{
					Target:    "var_210x260_",
					Geometry:  cleaver.Geometry{Width: 210, Height: 260},
					Quality:   0.70,
					Fit:       true,
					CopyEqual: true,
				},
			},
			&TypedTransform{
				Type: dto.PhotoVariantObjectGallerySmall,
				Transform: &cleaver.Transform{
					Target:    "var_250x_",
					Geometry:  cleaver.Geometry{Width: 250},
					Quality:   0.70,
					Fit:       true,
					CopyEqual: true,
				},
			},
			&TypedTransform{
				Type: dto.PhotoVariantObjectHD,
				Transform: &cleaver.Transform{
					Target:    "var_1280x_",
					Geometry:  cleaver.Geometry{Width: 1280},
					Quality:   0.70,
					Fit:       true,
					CopyEqual: true,
				},
			},
		},
	}
}

// AvatarTransform TBD
func AvatarTransform() *MediaTransform {
	return &MediaTransform{
		Type: dto.MediaTypeAvatar,
		Transforms: []*TypedTransform{
			&TypedTransform{
				Type: dto.PhotoVariantAvatar,
				Transform: &cleaver.Transform{
					Target:    "var_240x240_",
					Geometry:  cleaver.Geometry{Width: 240, Height: 240},
					Quality:   0.70,
					CopyEqual: true,
				},
			},
		},
	}
}

// GenericTransform TBD
func GenericTransform() *MediaTransform {
	return &MediaTransform{
		Type: dto.MediaTypeGenericPhoto,
		Transforms: []*TypedTransform{
			&TypedTransform{
				Type: dto.PhotoVariantGenericSmallThumb,
				Transform: &cleaver.Transform{
					Target:    "var_300x300_",
					Geometry:  cleaver.Geometry{Width: 300, Height: 300},
					Quality:   0.70,
					Fit:       true,
					CopyEqual: true,
				},
			},
			&TypedTransform{
				Type: dto.PhotoVariantGenericMedium,
				Transform: &cleaver.Transform{
					Target:    "var_600x_",
					Geometry:  cleaver.Geometry{Width: 600},
					Quality:   0.70,
					CopyEqual: true,
				},
			},
			&TypedTransform{
				Type: dto.PhotoVariantGenericLarge,
				Transform: &cleaver.Transform{
					Target:    "var_1280x_",
					Geometry:  cleaver.Geometry{Width: 1280},
					Quality:   0.70,
					CopyEqual: true,
				},
			},
		},
	}
}
