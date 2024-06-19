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

package require

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// RequiresMaxArgs returns an error if there is not at most max args
func RequiresMaxArgs(max int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) <= max {
			return nil
		}
		return errors.Errorf(
			"%q requires at most %d %s.\nSee '%s --help'.\n\nUsage:  %s\n\n%s",
			cmd.CommandPath(),
			max,
			pluralize("argument", max),
			cmd.CommandPath(),
			cmd.UseLine(),
			cmd.Short,
		)
	}
}

//nolint:unparam
func pluralize(word string, number int) string {
	if number == 1 {
		return word
	}
	return word + "s"
}
