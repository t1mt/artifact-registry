// Copyright 2023 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.linka.cloud/artifact-registry/pkg/packages"
)

func newPkgDeleteCmd(typ string) *cobra.Command {
	return &cobra.Command{
		Use:     fmt.Sprintf("delete [repository] [path]"),
		Short:   fmt.Sprintf("Delete %s package from the repository", typ),
		Aliases: []string{"rm", "remove", "del"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				c   packages.Deleter
				err error
			)
			cli, err := packages.NewCmd(typ)
			if err != nil {
				return err
			}
			c, err = cli.NewDelete(cmd.Context()).NewClient(packages.NewParams(args), opts)
			if err != nil {
				return err
			}
			return c.Delete(cmd.Context(), args[1])
		},
	}
}
