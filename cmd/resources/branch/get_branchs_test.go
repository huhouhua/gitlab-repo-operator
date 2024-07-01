// Copyright 2024 The huhouhua Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package branch

import (
	"bytes"
	"fmt"
	cmdtesting "github.com/huhouhua/gl/cmd/testing"
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"strings"
	"testing"
)

func TestGetBranch(t *testing.T) {
	tests := []struct {
		name           string
		optionsFunc    func(opt *ListOptions)
		args           []string
		expectedOutput string
	}{{
		name:           "project name is an empty string",
		args:           []string{""},
		expectedOutput: fmt.Sprintf("error from server (NotFound): project %s not found", ""),
	}}
	ioStreams := cmdutil.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewGetBranchesCmd(factory, ioStreams)
			cmdOptions := NewListOptions(ioStreams)
			if tc.optionsFunc != nil {
				tc.optionsFunc(cmdOptions)
			}
			out := cmdtesting.RunTestForStdout(func() {
				var err error
				if err = cmdOptions.Complete(factory, cmd, tc.args); err != nil {
					fmt.Print(err)
					return
				}
				if err = cmdOptions.Validate(cmd, tc.args); err != nil {
					fmt.Print(err)
					return
				}
				if err = cmdOptions.Run(tc.args); err != nil {
					fmt.Print(err)
					return
				}
			})
			cmdtesting.TInfo(out)
			if tc.expectedOutput == "" {
				t.Errorf("%s: Invalid test case. Specify expected result.\n", tc.name)
			}
			if !strings.Contains(out, tc.expectedOutput) {
				t.Errorf("%s: Unexpected output! Expected\n%s\ngot\n%s", tc.name, tc.expectedOutput, out)
			}
		})
	}
}

func TestRunGetBranch(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		flags          map[string]string
		expectedOutput string
	}{{
		name:           "list all branch with project",
		args:           []string{"devops-olympus/loki"},
		flags:          map[string]string{"all": "true"},
		expectedOutput: "main",
	}, {
		name:           "list all branch with project id",
		args:           []string{"220"},
		flags:          map[string]string{"all": "true"},
		expectedOutput: "main",
	}, {
		name:           "list all branch default",
		args:           []string{"220"},
		flags:          map[string]string{},
		expectedOutput: "main",
	}}
	ioStreams := cmdutil.NewTestIOStreamsDiscard()
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for i, arg := range tc.args {
				cmdtesting.TInfo(fmt.Sprintf("(%d) %s", i, arg))
			}
			buf := new(bytes.Buffer)
			cmd := NewGetBranchesCmd(factory, ioStreams)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			for flag, value := range tc.flags {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					t.Errorf("set %s flag error", err.Error())
					return
				}
			}
			out := cmdtesting.RunTestForStdout(func() {
				cmd.Run(cmd, tc.args)
			})
			cmdtesting.TInfo(out)
			if tc.expectedOutput == "" {
				t.Errorf("%s: Invalid test case. Specify expected result.\n", tc.name)
			}
			if !strings.Contains(out, tc.expectedOutput) {
				t.Errorf("%s: Unexpected output! Expected\n%s\ngot\n%s", tc.name, tc.expectedOutput, out)
			}
		})
	}

}
