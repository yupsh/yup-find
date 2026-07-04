package main

import (
	"context"
	"testing"

	clix "github.com/gloo-foo/cli"
	"github.com/spf13/afero"
	urf "github.com/urfave/cli/v3"
)

// parse runs args through a bare command carrying the wrapper's flags and
// returns the parsed accessor.
func parse(t *testing.T, args ...string) *urf.Command {
	t.Helper()
	var got *urf.Command
	app := &urf.Command{
		Name:   name,
		Flags:  flags(),
		Action: func(_ context.Context, c *urf.Command) error { got = c; return nil },
	}
	if err := app.Run(context.Background(), args); err != nil {
		t.Fatalf("parse: %v", err)
	}
	return got
}

func TestRoot(t *testing.T) {
	if got := root(parse(t, name)); got != "." {
		t.Fatalf("root=%q, want . for no operand", got)
	}
	if got := root(parse(t, name, "/etc")); got != "/etc" {
		t.Fatalf("root=%q, want /etc", got)
	}
}

func TestOptions(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want int // FindFs + one per set flag
	}{
		{"none", []string{name}, 1},
		{"type", []string{name, "--type", "f"}, 2},
		{"maxdepth", []string{name, "--maxdepth", "2"}, 2},
		{"both", []string{name, "--type", "d", "--maxdepth", "3"}, 3},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inv := clix.Invocation{Args: parse(t, tc.args...), Fs: afero.NewMemMapFs()}
			if got := len(options(inv)); got != tc.want {
				t.Fatalf("options len=%d, want %d", got, tc.want)
			}
		})
	}
}

func TestBuild(t *testing.T) {
	inv := clix.Invocation{Args: parse(t, name), Fs: afero.NewMemMapFs()}
	src, filter, err := build(inv)
	if err != nil || src == nil || filter != nil {
		t.Fatalf("build: src=%v filter=%v err=%v (want source, nil filter)", src, filter, err)
	}
}

func Test_main(t *testing.T) {
	orig := runMain
	t.Cleanup(func() { runMain = orig })
	var gotName clix.Name
	runMain = func(s clix.Spec, _ clix.Version) { gotName = s.Name }
	main()
	if gotName != name {
		t.Fatalf("main used spec %q, want %s", gotName, name)
	}
}
