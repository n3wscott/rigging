/*
Copyright 2019 The Rigging Authors

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

package example

import (
	"encoding/json"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/n3wscott/rigging/pkg/installer"
	"github.com/n3wscott/rigging/pkg/runner"
)

func init() {
	installer.RegisterPackage("github.com/n3wscott/rigging/example/cmd/echo")
}

// EchoTestImpl a very simple example test implementation.
func EchoTestImpl(t *testing.T) {

	cfg := make(map[string]interface{})
	cfg["echo"] = "hello world"

	i, err := installer.NewInstaller(t, cfg, installer.EndToEndConfigYaml([]string{"echo"})...)
	if err != nil {
		t.Errorf("failed to create installer, %s", err)
		return
	}

	// Create the resources for the test.
	if err := i.Setup(); err != nil {
		t.Errorf("failed to create, %s", err)
		return
	}

	// Teardown deferred.
	defer func() {
		if err := i.Teardown(); err != nil {
			t.Errorf("failed to create, %s", err)
		}
		// Just chill for tick.
		time.Sleep(10 * time.Second)
	}()

	jobGVR := schema.GroupVersionResource{
		Group:    "batch",
		Version:  "v1",
		Resource: "jobs",
	}

	msg, err := i.WaitUntilJobDone("echo")
	if err != nil {
		t.Error(err)
		return
	}
	if msg == "" {
		t.Error("No terminating message from the pod")
		return
	} else {
		out := &runner.Output{}
		if err := json.Unmarshal([]byte(msg), out); err != nil {
			t.Error(err)
			return
		}
		if !out.Success {
			if logs, err := i.LogsFor(i.Namespace, "echo", jobGVR); err != nil {
				t.Error(err)
			} else {
				t.Logf("job: %s\n", logs)
			}
			return
		}
	}
}
