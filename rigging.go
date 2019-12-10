package rigging

import (
	"errors"
	"fmt"
	yaml "github.com/jcrossley3/manifestival/pkg/manifestival"
	"github.com/n3wscott/rigging/pkg/installer"
	"github.com/n3wscott/rigging/pkg/lifecycle"
	v1 "k8s.io/api/core/v1"
	kntest "knative.dev/pkg/test"
	"knative.dev/pkg/test/helpers"
	"path/filepath"
	"runtime"
	"strings"
)

type Rigging interface {
	Install(config map[string]string, files ...string) error
	Uninstall() error
	Objects() []v1.ObjectReference
	Namespace() string
}

type Option func(Rigging) error

func New(opts ...Option) (Rigging, error) {
	r := &riggingImpl{
		configDirTemplate: DefaultConfigDirTemplate,
		runInParallel:     true,
	}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}

	// Set other defaults if not set.
	if r.rootDir == "" {
		_, filename, _, _ := runtime.Caller(1)
		r.rootDir = filepath.Dir(filename)
	}
	if r.name == "" {
		r.name = "rigging"
	}

	return r, nil
}

const (
	DefaultConfigDirTemplate = "%s/config/%s"
)

// WithConfigDirTemplate lets you change the configuration directory template.
// This value will be used to produce the fill file path of manifest files. And
// used like sprintf(dir, rootDir, name).
// By default:
//  - dir this value will be "%s/config/%s".
//  - rootDir is the directory of the caller to New.
//  - name is provided in a list for install
func WithConfigDirTemplate(dir string) Option {
	return func(r Rigging) error {
		if ri, ok := r.(*riggingImpl); ok {
			if !strings.HasPrefix(dir, "%s") || !strings.HasSuffix(dir, "%s") {
				return errors.New("invalid format; WithConfigDirTemplate(dir), expecting dir like '%s/<custom>/%s'")
			}

			ri.configDirTemplate = dir
			return nil
		}
		return errors.New("unknown rigging implementation")
	}
}

func WithRootDir(dir string) Option {
	return func(r Rigging) error {
		if ri, ok := r.(*riggingImpl); ok {
			ri.rootDir = dir
			return nil
		}
		return errors.New("unknown rigging implementation")
	}
}

type riggingImpl struct {
	configDirTemplate string
	rootDir           string
	runInParallel     bool
	namespace         string
	name              string

	client   *lifecycle.Client
	manifest yaml.Manifest
}

// Install implements Rigging.Install
// Install will do the following:
//  1. Create testing client and testing namespace.
//  2. Produce all images registered.
//  3. Pass the yaml config through the templating system.
//  4. Apply the yaml config to the cluster.
//
func (r *riggingImpl) Install(config map[string]string, name ...string) error {

	// 1. Create testing client and testing namespace.

	if err := r.createEnvironment(); err != nil {
		return err
	}

	// 2. Produce all images registered.

	cfg, err := r.updateConfig(config)
	if err != nil {
		return err
	}

	// 3. Pass the yaml config through the templating system.

	_ = cfg

	// 4. Apply yaml.
	if err := r.manifest.ApplyAll(); err != nil {
		return err
	}

	return nil
}

func (r *riggingImpl) createEnvironment() error {
	if r.namespace == "" {
		baseName := helpers.AppendRandomString(r.name)
		r.namespace = helpers.MakeK8sNamePrefix(baseName)
	}

	client, err := lifecycle.NewClient(
		kntest.Flags.Kubeconfig,
		kntest.Flags.Cluster,
		r.namespace)
	if err != nil {
		return fmt.Errorf("could not initialize clients: %v", err)
	}
	r.client = client

	if err := r.client.CreateNamespaceIfNeeded(); err != nil {
		return err
	}
	return nil
}

func (r *riggingImpl) updateConfig(config map[string]string) (map[string]interface{}, error) {
	ic, err := installer.ProduceImages()
	if err != nil {
		return nil, err
	}

	cfg := make(map[string]interface{})
	for k, v := range config {
		cfg[k] = v
	}
	// Implement template contract for Rigging:
	cfg["images"] = ic
	cfg["namespace"] = r.Namespace()
	return cfg, nil
}

// Uninstall implements Rigging.Uninstall
func (r *riggingImpl) Uninstall() error {

	// X. Delete yaml.
	if err := r.manifest.DeleteAll(); err != nil {
		return err
	}

	// Y. Delete namespace.
	if err := r.client.DeleteNamespaceIfNeeded(); err != nil {
		return err
	}
	return nil
}

// Objects implements Rigging.Objects
func (r *riggingImpl) Objects() []v1.ObjectReference {
	return nil
}

// Namespace implements Rigging.Namespace
func (r *riggingImpl) Namespace() string {
	return r.namespace
}
