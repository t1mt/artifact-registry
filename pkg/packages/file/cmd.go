package file

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
		NewClient: func(params []string, opts []hclient.Option) (packages.Client, error) {
			if len(params) != 2 {
				return nil, fmt.Errorf("invalid number of arguments")
			}
			registry := params[0]
			repository := params[1]
			return NewClient(registry, repository, opts...)
		},
	}

	return cmd
}

// NewPull implements packages.CmdProvider.
func (c clientProvider) NewPull(ctx context.Context) *packages.Cmd {
	cmd := &packages.Cmd{
		NewClient: func(params []string, opts []hclient.Option) (packages.Client, error) {
			if len(params) != 2 {
				return nil, fmt.Errorf("invalid number of arguments")
			}
			registry := params[0]
			repository := params[1]
			return NewClient(registry, repository, opts...)
		},
	}

	return cmd
}

// NewPush implements packages.CmdProvider.
func (c clientProvider) NewPush(ctx context.Context) *packages.Cmd {
	cmd := &packages.Cmd{
		Usage:   fmt.Sprintf("push [repository] [path]"),
		ArgsLen: 2,
		NewClient: func(params []string, opts []hclient.Option) (packages.Client, error) {
			if len(params) != 2 {
				return nil, fmt.Errorf("invalid number of arguments")
			}
			registry := params[0]
			repository := params[1]
			return NewClient(registry, repository, opts...)
		},
	}

	return cmd
}

// NewSetup implements packages.CmdProvider.
func (c clientProvider) NewSetup(ctx context.Context) *packages.Cmd {
	// ignore setup cmd
	return nil
}

func (c clientProvider) MakePackages(r io.Reader) ([]storage.Artifact, error) {
	var p []*Package
	err := json.NewDecoder(r).Decode(&p)
	return storage.AsArtifact(p), err
}
