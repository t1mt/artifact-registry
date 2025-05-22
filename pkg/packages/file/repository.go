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
	"encoding/json"

	"go.linka.cloud/artifact-registry/pkg/codec"
	"go.linka.cloud/artifact-registry/pkg/crypt/openpgp"
	"go.linka.cloud/artifact-registry/pkg/storage"
)

const (
	RepositoryPublicKey  = "repository.key"
	RepositoryPrivateKey = "private.key"
)

var _ storage.Repository = (*repo)(nil)

type repo struct{}

func (r *repo) Index(ctx context.Context, key string, artifacts ...storage.Artifact) ([]storage.Artifact, error) {
	// no index for files
	return nil, nil
}

func (r *repo) GenerateKeypair() (string, string, error) {
	return openpgp.GenerateKeypair("Artifact Registry", "File", "")
}

func (r *repo) KeyNames() (string, string) {
	return RepositoryPublicKey, RepositoryPrivateKey
}

func (r *repo) Codec() storage.Codec {
	return codec.Funcs[storage.Artifact]{
		Format: "json",
		EncodeFunc: func(v storage.Artifact) ([]byte, error) {
			return json.Marshal(v)
		},
		DecodeFunc: func(b []byte) (storage.Artifact, error) {
			var a Package
			return &a, json.Unmarshal(b, &a)
		},
	}
}

func (r *repo) Name() string {
	return Name
}
