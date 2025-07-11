package image

import (
	"context"
	"fmt"

	"github.com/loft-sh/image/docker/reference"
	"github.com/loft-sh/image/internal/imagesource"
	"github.com/loft-sh/image/internal/private"
	"github.com/loft-sh/image/manifest"
	"github.com/loft-sh/image/types"
	"github.com/opencontainers/go-digest"
)

// UnparsedImage implements types.UnparsedImage .
// An UnparsedImage is a pair of (ImageSource, instance digest); it can represent either a manifest list or a single image instance.
//
// This is publicly visible as c/image/image.UnparsedImage.
type UnparsedImage struct {
	src            private.ImageSource
	instanceDigest *digest.Digest
	cachedManifest []byte // A private cache for Manifest(); nil if not yet known.
	// A private cache for Manifest(), may be the empty string if guessing failed.
	// Valid iff cachedManifest is not nil.
	cachedManifestMIMEType string
}

// UnparsedInstance returns a types.UnparsedImage implementation for (source, instanceDigest).
// If instanceDigest is not nil, it contains a digest of the specific manifest instance to retrieve (when the primary manifest is a manifest list).
//
// This implementation of [types.UnparsedImage] ensures that [types.UnparsedImage.Manifest] validates the image
// against instanceDigest if set, or, if not, a digest implied by src.Reference, if any.
//
// The UnparsedImage must not be used after the underlying ImageSource is Close()d.
//
// This is publicly visible as c/image/image.UnparsedInstance.
func UnparsedInstance(src types.ImageSource, instanceDigest *digest.Digest) *UnparsedImage {
	return &UnparsedImage{
		src:            imagesource.FromPublic(src),
		instanceDigest: instanceDigest,
	}
}

// Reference returns the reference used to set up this source, _as specified by the user_
// (not as the image itself, or its underlying storage, claims).  This can be used e.g. to determine which public keys are trusted for this image.
func (i *UnparsedImage) Reference() types.ImageReference {
	// Note that this does not depend on instanceDigest; e.g. all instances within a manifest list need to be signed with the manifest list identity.
	return i.src.Reference()
}

// Manifest is like ImageSource.GetManifest, but the result is cached; it is OK to call this however often you need.
//
// Users of UnparsedImage are promised that this validates the image
// against either i.instanceDigest if set, or against a digest included in i.src.Reference.
func (i *UnparsedImage) Manifest(ctx context.Context) ([]byte, string, error) {
	if i.cachedManifest == nil {
		m, mt, err := i.src.GetManifest(ctx, i.instanceDigest)
		if err != nil {
			return nil, "", err
		}

		// ImageSource.GetManifest does not do digest verification, but we do;
		// this immediately protects also any user of types.Image.
		if digest, haveDigest := i.expectedManifestDigest(); haveDigest {
			matches, err := manifest.MatchesDigest(m, digest)
			if err != nil {
				return nil, "", fmt.Errorf("computing manifest digest: %w", err)
			}
			if !matches {
				return nil, "", fmt.Errorf("Manifest does not match provided manifest digest %s", digest)
			}
		}

		i.cachedManifest = m
		i.cachedManifestMIMEType = mt
	}
	return i.cachedManifest, i.cachedManifestMIMEType, nil
}

// expectedManifestDigest returns a the expected value of the manifest digest, and an indicator whether it is known.
// The bool return value seems redundant with digest != ""; it is used explicitly
// to refuse (unexpected) situations when the digest exists but is "".
func (i *UnparsedImage) expectedManifestDigest() (digest.Digest, bool) {
	if i.instanceDigest != nil {
		return *i.instanceDigest, true
	}
	ref := i.Reference().DockerReference()
	if ref != nil {
		if canonical, ok := ref.(reference.Canonical); ok {
			return canonical.Digest(), true
		}
	}
	return "", false
}
