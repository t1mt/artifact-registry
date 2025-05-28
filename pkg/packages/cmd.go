package packages

import (
	"context"
	"io"
	"strings"

	hclient "go.linka.cloud/artifact-registry/pkg/http/client"
	"go.linka.cloud/artifact-registry/pkg/storage"
)

type CmdType int

const (
	CmdSetup CmdType = 1 + iota
	CmdPull
	CmdPush
	CmdDelete
)

type Cmd struct {
	Usage     string
	ArgsLen   int
	NewClient func(params CmdParmas, opts []hclient.Option) (Client, error)
}

type CmdParmas struct {
	Registry   string
	Repository string
	ExtraArgs  []string
}

func NewParams(args []string) CmdParmas {
	arg := args[0]
	var registry, repository string
	parts := strings.SplitN(arg, "/", 2)
	registry = parts[0]
	if len(parts) > 1 {
		repository = parts[1]
	}
	return CmdParmas{
		Registry:   registry,
		Repository: repository,
		ExtraArgs:  args[1:],
	}
}

type CmdProvider interface {
	NewSetup(ctx context.Context) *Cmd
	NewPull(ctx context.Context) *Cmd
	NewPush(ctx context.Context) *Cmd
	NewDelete(ctx context.Context) *Cmd

	MakePackages(r io.Reader) ([]storage.Artifact, error)
}

var cmds = map[string]CmdProvider{}

func AddCmdProvider(name string, provider CmdProvider) {
	cmds[name] = provider
}

func NewCmd(name string) (CmdProvider, error) {
	p, ok := cmds[name]
	if !ok {
		return nil, ErrUnknownProvider
	}
	return p, nil
}
