package validation

import (
	"github.com/docker/distribution/reference"
	"github.com/golang/glog"
)

// ImageReferenceValiditySpec is a specification of an image reference validation request
type ImageReferenceValiditySpec struct {
	// Reference of the image
	Reference string `json:"reference"`
}

type ImageReferenceValidity struct {
	// True when the image reference is valid
	Valid bool `json:"valid"`

	// Error reason when image reference is valid
	Reason string `json:"reason"`
}

func ValidateImageReference(spec *ImageReferenceValiditySpec) (*ImageReferenceValidity, error) {
	glog.Infof("Validating %s as an image reference", spec.Reference)

	s := spec.Reference
	_, err := reference.ParseNamed(s)
	if err != nil {
		return &ImageReferenceValidity{Valid: false, Reason: err.Error()}, nil
	}
	return &ImageReferenceValidity{Valid: true}, nil
}
