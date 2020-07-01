/*
Copyright 2014 The Kubernetes Authors.

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

package hellokubernetes

import (
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdtesting "k8s.io/kubectl/pkg/cmd/testing"
)

func TestExtraArgsFail(t *testing.T) {
	cmdtesting.InitTestErrorHandler(t)

	f := cmdtesting.NewTestFactory()
	defer f.Cleanup()

	c := NewCmdHelloKubernetes(f, genericclioptions.NewTestIOStreamsDiscard())
	options := HelloKubernetesOptions{}
	if options.Validate(c, []string{"rc"}) == nil {
		t.Errorf("unexpected non-error")
	}
}

func TestHelloKubernetesObject(t *testing.T) {
	cmdtesting.InitTestErrorHandler(t)
	_, _, rc := cmdtesting.TestData()
	rc.Items[0].Name = "redis-master-controller"

	tf := cmdtesting.NewTestFactory().WithNamespace("test")
	defer tf.Cleanup()

	ioStreams, _, buf, _ := genericclioptions.NewTestIOStreams()
	cmd := NewCmdHelloKubernetes(tf, ioStreams)
	cmd.Flags().Set("filename", "../../../testdata/redis-master-controller.yaml")
	cmd.Flags().Set("output", "name")
	cmd.Run(cmd, []string{})

	// uses the name from the file, not the response
	if buf.String() != "Hello redis-master ReplicationController\n" {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestHelloKubernetesMultipleObject(t *testing.T) {
	cmdtesting.InitTestErrorHandler(t)

	tf := cmdtesting.NewTestFactory().WithNamespace("test")
	defer tf.Cleanup()

	ioStreams, _, buf, _ := genericclioptions.NewTestIOStreams()
	cmd := NewCmdHelloKubernetes(tf, ioStreams)
	cmd.Flags().Set("filename", "../../../testdata/redis-master-controller.yaml")
	cmd.Flags().Set("filename", "../../../testdata/frontend-service.yaml")
	cmd.Flags().Set("output", "name")
	cmd.Run(cmd, []string{})

	// Names should come from the REST response, NOT the files
	if buf.String() != "Hello redis-master ReplicationController\nHello frontend Service\n" {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestCreateDirectory(t *testing.T) {
	cmdtesting.InitTestErrorHandler(t)
	_, _, rc := cmdtesting.TestData()
	rc.Items[0].Name = "name"

	tf := cmdtesting.NewTestFactory().WithNamespace("test")
	defer tf.Cleanup()

	ioStreams, _, buf, _ := genericclioptions.NewTestIOStreams()
	cmd := NewCmdHelloKubernetes(tf, ioStreams)
	cmd.Flags().Set("filename", "../../../testdata/replace/legacy")
	cmd.Flags().Set("output", "name")
	cmd.Run(cmd, []string{})

	if buf.String() != "Hello frontend ReplicationController\nHello redis-master ReplicationController\nHello redis-slave ReplicationController\n" {
		t.Errorf("unexpected output: %s", buf.String())
	}
}
