// Command yup-find is the CLI wrapper around github.com/gloo-foo/cmd-find.
package main

import (
	clix "github.com/gloo-foo/cli"
	command "github.com/gloo-foo/cmd-find"
	urf "github.com/urfave/cli/v3"
)

// version is the build version. It defaults to "dev" for local builds and is
// overridden at release time via the linker: -ldflags "-X main.version=<v>".
var version = "dev"

const (
	name         = "find"
	flagType     = "type"
	flagMaxDepth = "maxdepth"
)

// synopsis is the multi-line --help usage block; urfave/cli indents it three
// spaces, so the lines stay flush-left.
const synopsis = `find [PATH] [OPTIONS]

search for files in the directory hierarchy rooted at PATH (default: .).`

// spec declares the find wrapper. find is a source command: it walks the
// injected filesystem and produces its listing directly, so build returns it as
// the whole pipeline (a nil filter).
var spec = clix.Spec{
	Name:     name,
	Summary:  "search for files in a directory hierarchy",
	Synopsis: synopsis,
	Build:    build,
	Flags:    flags(),
}

// flags returns a fresh set of the wrapper's flags. Each call yields new flag
// values, so parsing one invocation never leaks urfave/cli's per-flag "was set"
// state into another (which IsSet reads).
func flags() []urf.Flag {
	return []urf.Flag{
		&urf.StringFlag{Name: flagType, Usage: "file is of type TYPE (f=file, d=directory)"},
		&urf.IntFlag{Name: flagMaxDepth, Usage: "descend at most LEVELS (a non-negative integer) levels"},
	}
}

// build maps the invocation to find's pipeline: the root operand and flags
// produce the walk source, with no filter.
func build(inv clix.Invocation) (clix.Source, clix.Command, error) {
	return command.Find(clix.File(root(inv.Args)), options(inv)...), nil, nil
}

// root is the single path operand, defaulting to the current directory.
func root(c *urf.Command) string {
	if c.NArg() == 0 {
		return "."
	}
	return c.Args().First()
}

// options folds the injected filesystem and the parsed flags into find's option
// values. FindFs is always applied.
func options(inv clix.Invocation) []any {
	opts := []any{command.FindFs{Fs: inv.Fs}}
	if inv.Args.IsSet(flagType) {
		opts = append(opts, command.FindType(inv.Args.String(flagType)))
	}
	if inv.Args.IsSet(flagMaxDepth) {
		opts = append(opts, command.FindMaxDepth(inv.Args.Int(flagMaxDepth)))
	}
	return opts
}

// runMain is an indirection seam so main's wiring is testable without spawning
// the process; a test swaps it and restores it.
var runMain = clix.Main

func main() { runMain(spec, version) }
