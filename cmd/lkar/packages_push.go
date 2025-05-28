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
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.linka.cloud/artifact-registry/pkg/packages"
	"go.linka.cloud/grpc-toolkit/logger"
)

func newPkgPushCmd(typ string) *cobra.Command {
	prvd, err := packages.NewCmd(typ)
	if err != nil {
		panic(err)
	}
	cli := prvd.NewPush(context.TODO())
	use := cli.Usage
	index := cli.ArgsLen
	return &cobra.Command{
		Use:     use,
		Short:   fmt.Sprintf("Push %s package to the repository", typ),
		Aliases: []string{"put", "create", "upload"},
		Args:    cobra.ExactArgs(index),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			f, err := os.Open(args[index-1])
			if err != nil {
				return err
			}
			defer f.Close()
			i, err := f.Stat()
			if err != nil {
				return err
			}
			c, err := cli.NewClient(packages.NewParams(args), opts)
			// c, err := client(args)
			if err != nil {
				return err
			}
			pw := newProgressReader(f, i.Size())
			go pw.Run(ctx)
			if err := c.Push(ctx, pw); err != nil {
				logger.C(ctx).Error(err)
				return err
			}
			return nil
		},
	}
}
