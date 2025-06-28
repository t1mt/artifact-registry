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
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.linka.cloud/artifact-registry/pkg/packages"
)

func newPkgPullCmd(typ string) *cobra.Command {
	var output string
	cmd := &cobra.Command{
		Use:     fmt.Sprintf("pull [repository] [path]"),
		Short:   fmt.Sprintf("Download %s package from the repository", typ),
		Aliases: []string{"dl", "get", "read", "download"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			path := args[1]
			if output == "" {
				output = filepath.Base(path)
			}
			if _, err := os.Stat(filepath.Dir(output)); err != nil {
				return err
			}
			var (
				c   packages.Puller
				err error
			)
			cli, err := packages.NewCmd(typ)
			if err != nil {
				return err
			}
			c, err = cli.NewPull(cmd.Context()).NewClient(packages.NewParams(args), opts)
			if err != nil {
				return err
			}
			r, size, err := c.Pull(ctx, path)
			if err != nil {
				return err
			}
			f, err := os.Create(output)
			if err != nil {
				return err
			}
			defer f.Close()
			pw := newProgressReader(r, size)
			defer pw.Close()
			go pw.Run(ctx)
			if _, err := io.Copy(f, pw); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file")
	return cmd
}
