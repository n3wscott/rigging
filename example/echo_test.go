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
	"fmt"
	"github.com/n3wscott/rigging"
	"github.com/n3wscott/rigging/pkg/installer"
	"github.com/n3wscott/rigging/pkg/runner"
	"k8s.io/apimachinery/pkg/api/meta"
	"regexp"
	"testing"
	"time"
)

func init() {
	installer.RegisterPackage("github.com/n3wscott/rigging/example/cmd/echo")
}

func TestParse(t *testing.T) {
	s := "rigging-mekmiegt/echo (batch/v1, Kind=Job)"
	matches := regexp.MustCompile(`\[(.*?)\]`).FindAllStringSubmatch(s, -1)
	if matches == nil {
		fmt.Println("No matches found.")
		return
	}

	for i, match := range matches {
		full := match[0]
		submatches := match[1:len(match)]
		fmt.Printf("%v => \"%v\" from \"%v\"\n", i, submatches[0], full)
	}
}

// EchoTestImpl a very simple example test implementation.
func EchoTestImpl(t *testing.T) {
	opts := []rigging.Option{
		rigging.WithImages(map[string]string{
			"echo": "n3wscott.azurecr.io/echo-b301ec929b6c030bb4dd170136bb2fb3@sha256:ab79d01c4478302f21fd08071c6d78fbe5a0096ae017860ab382f51f327b73d5",
		}),
	}

	rig, err := rigging.NewInstall(opts, []string{"echo"}, map[string]string{"echo": "hello world"})
	if err != nil {
		t.Fatalf("failed to create rig, %s", err)
	}

	t.Logf("Created a new testing rig at namespace %s.", rig.Namespace())

	// Uninstall deferred.
	defer func() {
		if err := rig.Uninstall(); err != nil {
			t.Errorf("failed to uninstall, %s", err)
		}
	}()

	refs := rig.Objects()
	for _, r := range refs {
		k := r.GroupVersionKind()
		gvk, _ := meta.UnsafeGuessKindToResource(k)
		t.Log("UnsafeGuessKindToResource:")
		t.Log(gvk)

		msg, err := rig.WaitForReadyOrDone(r, 45*time.Second)
		if err != nil {
			t.Fatalf("failed to wait for ready or done, %s", err)
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
			//if !out.Success {
			//	if logs, err := i.LogsFor(i.Namespace, r.Name, gvk); err != nil {
			//		t.Error(err)
			//	} else {
			//		t.Logf("job: %s\n", logs)
			//	}
			//	return
			//}
		}
	}
}
