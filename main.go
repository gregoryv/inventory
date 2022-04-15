// Command inventory prints project latest release information
package main

import (
	"os"

	"github.com/gregoryv/cmdline"
)

func main() {
	var (
		cli  = cmdline.NewBasicParser()
		s    = cli.Flag("-s, --skip-untagged")
		f    = cli.Flag("-f, --show-full-path")
		m    = cli.Flag("-m, --show-modified-date")
		o    = cli.Option("-o, --order-by").Enum("releaseDate", "path", "releaseDate")
		i    = cli.Flag("-i, --include-vendor")
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
	cmd.SetShowFullPath(f)
	cmd.SetOrderBy(o)
	cmd.SetIncludeVendor(i)
	cmd.SetPaths(args)
	cmd.SetOutput(os.Stdout)
	cmd.SetRoot(os.Getenv("HOME"))
	cmd.Run()
}
