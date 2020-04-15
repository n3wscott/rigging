package manifest

import (
	"context"
	"github.com/mattmoor/bindings/pkg/bindings"
	"github.com/n3wscott/rigging/pkg/manifest"
	"knative.dev/pkg/injection/clients/dynamicclient"
)

const (
	MountPath = bindings.MountPath + "/manifests" // filepath.Join isn't const.
)

// New
func New(ctx context.Context) (manifest.Manifest, error) {
	return manifest.NewYamlManifest(MountPath, true, dynamicclient.Get(ctx))
}
