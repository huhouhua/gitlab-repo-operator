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
	"errors"
	"fmt"
	cmdtesting "github.com/huhouhua/gitlab-repo-operator/cmd/testing"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"strings"
	"testing"
)

func TestDeleteBranch(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		optionsFunc func() *DeleteOptions
		validate    func(opt *DeleteOptions, cmd *cobra.Command, args []string) error
		run         func(opt *DeleteOptions, args []string) error
		wantError   error
	}{{
		name: "delete by name",
		args: []string{"develop"},
		optionsFunc: func() *DeleteOptions {
			opt := NewDeleteOptions()
			opt.project = "huhouhua/gitlab-repo-branch"
			return opt
		},
		run: func(opt *DeleteOptions, args []string) error {
			var err error
			out := cmdtesting.RunTestForStdout(func() {
				err = opt.Run(args)
			})
			expectedOutput := fmt.Sprintf("Branch (%s) from project (%s) has been deleted", opt.branch, opt.project)
			if !strings.Contains(out, expectedOutput) {
				err = errors.New(fmt.Sprintf("delete by path : Unexpected output! Expected\n%s\ngot\n%s", expectedOutput, out))
			}
			return err
		},
		wantError: nil,
	}, {
		name: "branch not found",
		args: []string{"not-found"},
		optionsFunc: func() *DeleteOptions {
			opt := NewDeleteOptions()
			opt.project = "huhouhua/gitlab-repo-branch"
			return opt
		},
		run: func(opt *DeleteOptions, args []string) error {
			err := opt.Run(args)
			var repo *gitlab.ErrorResponse
			if errors.As(err, &repo) && repo.Message == "{message: 404 Branch Not Found}" {
				return nil
			}
			return err
		},
	}, {
		name: "project not found",
		args: []string{"not-found"},
		optionsFunc: func() *DeleteOptions {
			opt := NewDeleteOptions()
			opt.project = "not-found"
			return opt
		},
		run: func(opt *DeleteOptions, args []string) error {
			err := opt.Run(args)
			var repo *gitlab.ErrorResponse
			if errors.As(err, &repo) && repo.Message == "{message: 404 Project Not Found}" {
				return nil
			}
			return err
		},
	}, {
		name: "not definition branch",
		args: []string{},
		validate: func(opt *DeleteOptions, cmd *cobra.Command, args []string) error {
			err := opt.Validate(cmd, args)
			if err.Error() == "please enter branch" {
				return err
			}
			return nil
		},
	}}
	factory := cmdutil.NewFactory(cmdtesting.NewFakeRESTClientGetter())
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewDeleteBranchCmd(factory)
			var cmdOptions *DeleteOptions
			if tc.optionsFunc != nil {
				cmdOptions = tc.optionsFunc()
			} else {
				cmdOptions = NewDeleteOptions()
			}
			var err error
			if err = cmdOptions.Complete(factory, cmd, tc.args); err != nil && !errors.Is(err, tc.wantError) {
				t.Errorf("expected %v, got: '%v'", tc.wantError, err)
				return
			}
			if tc.validate != nil {
				err = tc.validate(cmdOptions, cmd, tc.args)
				if err != nil {
					return
				}
			} else {
				if err = cmdOptions.Validate(cmd, tc.args); err != nil && !errors.Is(err, tc.wantError) {
					t.Errorf("expected %v, got: '%v'", tc.wantError, err)
					return
				}
			}
			if tc.run != nil {
				err = tc.run(cmdOptions, tc.args)
				if err != nil {
					t.Error(err)
				}
				return
			}
			if err = cmdOptions.Run(tc.args); !errors.Is(err, tc.wantError) {
				t.Errorf("expected %v, got: '%v'", tc.wantError, err)
				return
			}
		})
	}
}
