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

package file

import (
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/huhouhua/gl/cmd/require"
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"github.com/huhouhua/gl/util/cli"
	"github.com/huhouhua/gl/util/progress"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"strings"
	"sync"
)

type ReplaceOptions struct {
	gitlabClient *gitlab.Client
	path         string
	branchList   *gitlab.ListBranchesOptions
	content      []byte
	Project      string
	Ref          string
	RefMatch     string
	FileName     string
	Force        bool
	ioStreams    cli.IOStreams
}

func NewReplaceOptions(ioStreams cli.IOStreams) *ReplaceOptions {
	return &ReplaceOptions{
		ioStreams: ioStreams,
		branchList: &gitlab.ListBranchesOptions{
			ListOptions: gitlab.ListOptions{
				Page:    1,
				PerPage: 10,
			},
		},
	}
}

var (
	replaceFileDesc = "replace file for project "

	replaceFileExample = `# edit file
gl replace files app/my.yml -p myproject --ref=main -f ./my.yml`
)

func NewReplaceFileCmd(f cmdutil.Factory, ioStreams cli.IOStreams) *cobra.Command {
	o := NewReplaceOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "file",
		Aliases:               []string{"f"},
		Short:                 replaceFileDesc,
		Example:               replaceFileExample,
		Args:                  require.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{"file"},
	}
	o.AddFlags(cmd)
	return cmd
}
func (o *ReplaceOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddProjectVarPFlag(cmd, &o.Project)
	f := cmd.Flags()
	f.StringVar(&o.Ref, "ref", o.Ref, "The name of a repository branch or tag or, if not given, the default branch.")
	f.StringVar(&o.RefMatch, "ref-match", o.RefMatch, "match repository branch or tag or, if not given, the use --ref matching branch.")
	f.StringVarP(&o.FileName, "filename", "f", "", "to use to replace the repository file .")
	f.BoolVar(&o.Force, "force", o.Force, "If true, immediately remove repository file from API and bypass graceful deletion. Note that immediate deletion of some  repository file may result in inconsistency or data loss and requires confirmation.")
	cmdutil.VerifyMarkFlagRequired(cmd, "project")
	cmdutil.VerifyMarkFlagRequired(cmd, "filename")
}

// Complete completes all the required options.
func (o *ReplaceOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) > 0 {
		o.path = args[0]
	}
	if strings.TrimSpace(o.Ref) != "" {
		o.branchList.Search = pointer.ToString(o.Ref)
	}
	if strings.TrimSpace(o.RefMatch) != "" {
		o.branchList.Regex = pointer.ToString(o.RefMatch)
	}
	o.gitlabClient, err = f.GitlabClient()
	o.content, err = cmdutil.ReadFile(o.FileName)
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *ReplaceOptions) Validate(cmd *cobra.Command, args []string) error {
	if strings.TrimSpace(o.Project) == "" {
		_ = cmd.Usage()
		return fmt.Errorf("please enter project name and id")
	}
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *ReplaceOptions) Run(args []string) error {
	for {
		branches, err := o.nextBranches()
		if err != nil {
			return err
		}
		if len(branches) == 0 {
			break
		}
		var wg = sync.WaitGroup{}
		for _, item := range branches {
			wg.Add(1)
			go func(branch *gitlab.Branch) {
				defer wg.Done()
				o.updateFile(branch)
			}(item)
		}
		wg.Wait()
	}
	return nil
}

func (o *ReplaceOptions) updateFile(branch *gitlab.Branch) {
	s := progress.CreatingEvent(false).WithText(fmt.Sprintf(" %s ...", branch.Name)).Start()
	defer s.Stop()
	_, _, err := o.gitlabClient.RepositoryFiles.UpdateFile(o.Project, o.path, &gitlab.UpdateFileOptions{
		Branch:        pointer.ToString(branch.Name),
		CommitMessage: pointer.ToString(fmt.Sprintf("update %s from gl command line", o.path)),
		Content:       pointer.ToString(string(o.content)),
	})
	if err != nil {
		s.Error("Error\n", err.Error())
		return
	}
	s.Success()
}

func (o *ReplaceOptions) nextBranches() ([]*gitlab.Branch, error) {
	s := progress.CreatingEvent(true).WithText(fmt.Sprintf(" pull branch on page %d", o.branchList.ListOptions.Page)).Start()
	defer s.Stop()
	branches, _, err := o.gitlabClient.Branches.ListBranches(o.Project, o.branchList)
	o.branchList.ListOptions.Page++
	if err != nil {
		s.Error("Error\n", err.Error())
		return nil, err
	}
	if len(branches) == 0 {
		s.Done()
		return branches, nil
	}
	s.Success()
	return branches, nil
}