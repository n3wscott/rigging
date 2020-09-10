package manifest

import (
	"context"
	"github.com/n3wscott/rigging/pkg/manifest"
	"knative.dev/pkg/injection/clients/dynamicclient"
)

const (
	MountPath = "/var/bindings/manifests" // filepath.Join isn't const.
)

// New
func New(ctx context.Context) (manifest.Manifest, error) {
	return manifest.NewYamlManifest(MountPath, true, dynamicclient.Get(ctx))
}
