package packages

import (
	"context"
	"io"

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
	NewClient func(args []string, opts []hclient.Option) (Client, error)
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

func NewCmdProvider(name string) (CmdProvider, error) {
	p, ok := cmds[name]
	if !ok {
		return nil, ErrUnknownProvider
	}
	return p, nil
}

/* func NewCmd(ctx context.Context, name string, cmd CmdType) (*Cmd, error) {
	p, ok := cmds[name]
	if !ok {
		return nil, ErrUnknownProvider
	}
	switch cmd {
	case CmdSetup:
		return p.NewSetup(ctx)
	case CmdPull:
		return p.NewPull(ctx)
	case CmdPush:
		return p.NewPush(ctx)
	case CmdDelete:
		return p.NewDelete(ctx)
	default:
		return nil, fmt.Errorf("unsupported package type: %s", name)
	}
} */
