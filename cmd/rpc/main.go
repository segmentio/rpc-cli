package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/segmentio/rpc-cli/internal/rpc"
	"github.com/tj/docopt"
)

const version = ""
const usage = `
  Usage:
    rpc <addr> <method> <args>...
    rpc <addr> <method>
    rpc <addr>

    rpc -h | --help
    rpc -v | --version

  Options:
    -h, --help          show help information
    -v, --version       show version information

`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	check(err)

	addr := args["<addr>"].(string)

	var input io.Reader

	if s, _ := os.Stdin.Stat(); s.Mode()&os.ModeCharDevice == 0 {
		input = os.Stdin
	}

	cmd := rpc.New()
	cmd.HTTP = strings.HasPrefix(addr, "http")
	cmd.Input = input
	cmd.Output = os.Stdout
	cmd.Addr = addr
	cmd.Method, _ = args["<method>"].(string)
	cmd.Args, _ = args["<args>"].([]string)

	check(cmd.Run())
}

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "rpc: %s\n", err)
	}
}
