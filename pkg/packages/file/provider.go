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

package file

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.linka.cloud/artifact-registry/pkg/packages"
	"go.linka.cloud/artifact-registry/pkg/storage"
)

const Name = "file"

var _ packages.Provider = (*provider)(nil)

func init() {
	packages.Register(Name, newProvider)
}

func newProvider(_ context.Context) (packages.Provider, error) {
	return &provider{}, nil
}

type provider struct{}

func (p *provider) Repository() storage.Repository {
	return &repo{}
}

func (p *provider) Routes() []*packages.Route {
	return []*packages.Route{
		{
			Path:   "", // ?filename={}
			Method: http.MethodPut,
			Handler: packages.Push(func(r *http.Request, reader io.Reader, key string) (storage.Artifact, error) {
				query := r.URL.Query()
				path := query.Get("filename")
				// path := mux.Vars(r)["filename"]
				return NewPackage(reader, path)
			}),
		},
		{
			Path:   "", // ?filename={}
			Method: http.MethodDelete,
			Handler: packages.Delete(func(r *http.Request) string {
				query := r.URL.Query()
				path := query.Get("filename")
				// path := mux.Vars(r)["filename"]
				return fmt.Sprintf("%s", path)
			}),
		},
		{
			Path:   "", // ?filename={}
			Method: http.MethodGet,
			Handler: packages.Pull(func(r *http.Request) string {
				query := r.URL.Query()
				filename := query.Get("filename")
				return strings.TrimPrefix(filename, "/")
				// return mux.Vars(r)["filename"]
			}),
		},
	}
}
