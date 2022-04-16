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

	sys := NewSystem()
	sys.SetSkipUntagged(s)
	sys.SetShowModifiedDate(m)
	sys.SetShowFullPath(f)
	sys.SetOrderBy(o)
	sys.SetIncludeVendor(i)
	sys.SetPaths(args)
	sys.SetOutput(os.Stdout)
	sys.SetRoot(os.Getenv("HOME"))
	sys.Run()
}
