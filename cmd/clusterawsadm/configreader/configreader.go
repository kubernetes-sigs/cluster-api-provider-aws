/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package configreader

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	bootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/bootstrap/v1alpha1"
	bootstrapschemev1 "sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/api/bootstrap/v1alpha1/scheme"
)

type errEmptyBootstrapConfig string

func (e errEmptyBootstrapConfig) Error() string {
	return fmt.Sprintf("bootstrap config file %q was empty", string(e))
}

// LoadConfigFile loads a YAML file representing a bootstrapv1.AWSIAMConfiguration.
func LoadConfigFile(name string) (*bootstrapv1.AWSIAMConfiguration, error) {
	// compute absolute path based on current working dir
	iamConfigFile, err := filepath.Abs(name)
	if err != nil {
		return nil, fmt.Errorf("failed to convert IAM config path into absolute path %s, error: %w", name, err)
	}
	loader, err := newFsLoader(iamConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize filesystem loader: %w", err)
	}
	return loader.Load()
}

// Loader loads configuration from a storage layer.
type loader interface {
	// Load loads and returns the AWSIAMConfiguration from the storage layer, or an error if a configuration could not be loaded.
	Load() (*bootstrapv1.AWSIAMConfiguration, error)
}

// fsLoader loads configuration from `configDir`..
type fsLoader struct {

	// bootstrapCodecs is the scheme used to decode config files
	bootstrapCodecs *serializer.CodecFactory
	// bootstrapFile is an absolute path to the file containing a serialized KubeletConfiguration
	bootstrapFile string
}

// ReadFile reads a file.
func (fsLoader) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Clean(filename))
}

// NewFsLoader returns a Loader that loads a AWSIAMConfiguration from the `config file`.
func newFsLoader(bootstrapFile string) (loader, error) {
	_, bootstrapCodecs, err := bootstrapschemev1.NewSchemeAndCodecs()
	if err != nil {
		return nil, err
	}

	return &fsLoader{
		bootstrapCodecs: bootstrapCodecs,
		bootstrapFile:   bootstrapFile,
	}, nil
}

func (loader *fsLoader) Load() (*bootstrapv1.AWSIAMConfiguration, error) {
	data, err := loader.ReadFile(loader.bootstrapFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read bootstrap config file %q, error: %w", loader.bootstrapFile, err)
	}

	// no configuration is an error, some parameters are required
	if len(data) == 0 {
		return nil, errEmptyBootstrapConfig(loader.bootstrapFile)
	}

	kc, err := DecodeBootstrapConfiguration(loader.bootstrapCodecs, data)
	if err != nil {
		return nil, err
	}

	fileDir := filepath.Dir(loader.bootstrapFile)

	// make all paths absolute
	resolveRelativePaths([]*string{&fileDir}, "")
	return kc, nil
}

// resolveRelativePaths makes relative paths absolute by resolving them against `root`.
func resolveRelativePaths(paths []*string, root string) {
	for _, path := range paths {
		// leave empty paths alone, "no path" is a valid input
		// do not attempt to resolve paths that are already absolute
		if len(*path) > 0 && !filepath.IsAbs(*path) {
			*path = filepath.Join(root, *path)
		}
	}
}

// DecodeBootstrapConfiguration decodes a serialized AWSIAMConfiguration to the internal type.
func DecodeBootstrapConfiguration(bootstrapCodecs *serializer.CodecFactory, data []byte) (*bootstrapv1.AWSIAMConfiguration, error) {
	obj := &bootstrapv1.AWSIAMConfiguration{}

	if err := runtime.DecodeInto(bootstrapCodecs.UniversalDecoder(), data, obj); err != nil {
		return nil, errors.Wrap(err, "error decoding metadata.yaml")
	}

	bootstrapv1.SetDefaults_AWSIAMConfiguration(obj)

	return obj, nil
}
