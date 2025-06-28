package apk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	hclient "go.linka.cloud/artifact-registry/pkg/http/client"
	"go.linka.cloud/artifact-registry/pkg/packages"
	"go.linka.cloud/artifact-registry/pkg/storage"
)

var _ packages.CmdProvider = (*clientProvider)(nil)

func init() {
	packages.AddCmdProvider(Name, clientProvider{})
}

type clientProvider struct{}

// NewDelete implements packages.CmdProvider.
func (c clientProvider) NewDelete(ctx context.Context) *packages.Cmd {
	cmd := &packages.Cmd{
		NewClient: func(params packages.CmdParmas, opts []hclient.Option) (packages.Client, error) {
			registry := params.Registry
			repository := params.Repository
			return NewClient(registry, repository, "", "", opts...)
		},
	}

	return cmd
}

// NewPull implements packages.CmdProvider.
func (c clientProvider) NewPull(ctx context.Context) *packages.Cmd {
	cmd := &packages.Cmd{
		NewClient: func(params packages.CmdParmas, opts []hclient.Option) (packages.Client, error) {
			registry := params.Registry
			repository := params.Repository
			return NewClient(registry, repository, "", "", opts...)
		},
	}

	return cmd
}

// NewPush implements packages.CmdProvider.
func (c clientProvider) NewPush(ctx context.Context) *packages.Cmd {
	cmd := &packages.Cmd{
		Usage:   fmt.Sprintf("push [repository] [branch] [apk-repository] [path]"),
		ArgsLen: 4,
		NewClient: func(params packages.CmdParmas, opts []hclient.Option) (packages.Client, error) {
			registry := params.Registry
			repository := params.Repository
			branch := params.ExtraArgs[0]
			repo := params.ExtraArgs[1]
			return NewClient(registry, repository, branch, repo, opts...)
		},
	}

	return cmd
}

// NewSetup implements packages.CmdProvider.
func (c clientProvider) NewSetup(ctx context.Context) *packages.Cmd {
	cmd := &packages.Cmd{
		Usage:   fmt.Sprintf("setup [repository] [branch] [apk-repository]"),
		ArgsLen: 3,
		NewClient: func(params packages.CmdParmas, opts []hclient.Option) (packages.Client, error) {
			registry := params.Registry
			repository := params.Repository
			branch := params.ExtraArgs[0]
			repo := params.ExtraArgs[1]
			return NewClient(registry, repository, branch, repo, opts...)
		},
	}
	return cmd
}

func (c clientProvider) MakePackages(r io.Reader) ([]storage.Artifact, error) {
	var p []*Package
	err := json.NewDecoder(r).Decode(&p)
	return storage.AsArtifact(p), err
}
