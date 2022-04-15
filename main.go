// Command inventory prints project latest release information
package main

import (
	"log"
	"os"

	"github.com/gregoryv/cmdline"
)

func main() {
	var (
		cli  = cmdline.NewBasicParser()
		s    = cli.Flag("-s, --skip-untagged")
		m    = cli.Flag("-m, --show-modified-date")
		args = cli.Args()
	)

	u := cli.Usage()
	u.Preface("List projects and release information")
	u.Example(
		"List all your projects",
		"$ inventory",
	)
	u.Example(
		"List specific projects",
		"$ inventory $HOME/src/github.com/YOURS/*",
	)
	cli.Parse()

	var cmd System
	cmd.SetSkipUntagged(s)
	cmd.SetShowModifiedDate(m)
	cmd.SetPaths(args)
	cmd.SetOutput(os.Stdout)
	cmd.SetRoot(os.Getenv("HOME"))

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
