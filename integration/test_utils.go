/*
Copyright 2017 The Kubernetes Authors.

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

package integration

import (
	"errors"
	"time"

	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/kubelet/apis/cri"
	"k8s.io/kubernetes/pkg/kubelet/apis/cri/v1alpha1/runtime"
	"k8s.io/kubernetes/pkg/kubelet/remote"

	"github.com/kubernetes-incubator/cri-containerd/pkg/util"
)

const (
	sock       = "/var/run/cri-containerd.sock"
	timeout    = 1 * time.Minute
	pauseImage = "gcr.io/google_containers/pause:3.0"
)

var (
	runtimeService cri.RuntimeService
	imageService   cri.ImageManagerService
)

func init() {
	var err error
	runtimeService, err = remote.NewRemoteRuntimeService(sock, timeout)
	if err != nil {
		glog.Exitf("Failed to create runtime service: %v", err)
	}
	imageService, err = remote.NewRemoteImageService(sock, timeout)
	if err != nil {
		glog.Exitf("Failed to create image service: %v", err)
	}
}

// Opts sets specific information in pod sandbox config.
type PodSandboxOpts func(*runtime.PodSandboxConfig)

// PodSandboxConfig generates a pod sandbox config for test.
func PodSandboxConfig(name, ns string, opts ...PodSandboxOpts) *runtime.PodSandboxConfig {
	config := &runtime.PodSandboxConfig{
		Metadata: &runtime.PodSandboxMetadata{
			Name: name,
			// Using random id as uuid is good enough for local
			// integration test.
			Uid:       util.GenerateID(),
			Namespace: ns,
		},
		Linux: &runtime.LinuxPodSandboxConfig{},
	}
	for _, opt := range opts {
		opt(config)
	}
	return config
}

// ContainerOpts to set any specific attribute like labels,
// annotations, metadata etc
type ContainerOpts func(*runtime.ContainerConfig)

func WithTestLabels() ContainerOpts {
	return func(cf *runtime.ContainerConfig) {
		cf.Labels = map[string]string{"key": "value"}
	}
}

func WithTestAnnotations() ContainerOpts {
	return func(cf *runtime.ContainerConfig) {
		cf.Annotations = map[string]string{"a.b.c": "test"}
	}
}

// ContainerConfig creates a container config given a name and image name
// and additional container config options
func ContainerConfig(name, image string, opts ...ContainerOpts) *runtime.ContainerConfig {
	cConfig := &runtime.ContainerConfig{
		Metadata: &runtime.ContainerMetadata{
			Name: name,
		},
		Image: &runtime.ImageSpec{Image: image},
	}
	for _, opt := range opts {
		opt(cConfig)
	}
	return cConfig
}

// CheckFunc is the function used to check a condition is true/false.
type CheckFunc func() (bool, error)

// Eventually waits for f to return true, it checks every period, and
// returns error if timeout exceeds. If f returns error, Eventually
// will return the same error immediately.
func Eventually(f CheckFunc, period, timeout time.Duration) error {
	start := time.Now()
	for {
		done, err := f()
		if done {
			return nil
		}
		if err != nil {
			return err
		}
		if time.Since(start) >= timeout {
			return errors.New("timeout exceeded")
		}
		time.Sleep(period)
	}
}
