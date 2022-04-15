// Command inventory prints project latest release information
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/gregoryv/cmdline"
)

func main() {
	var (
		cli          = cmdline.NewBasicParser()
		skipUntagged = cli.Flag("-s, --skip-untagged")
		paths        = cli.Args()
	)
	u := cli.Usage()
	u.Preface("List projects and release information")

	u.Example(
		"List all your projects",
		"$ inventory $HOME/src/github.com/YOURS/*",
	)
	cli.Parse()

	var cmd InventoryCmd
	cmd.SetSkipUntagged(skipUntagged)
	cmd.SetPaths(paths)
	cmd.SetOutput(os.Stdout)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

type InventoryCmd struct {
	skipUntagged bool
	paths        []string
	out          io.Writer
}

func (me *InventoryCmd) Run() error {
	var i int
	for _, dir := range me.Paths() {
		version, err := latestVersion(dir)
		if errors.Is(ErrNoTags, err) && me.skipUntagged {
			continue
		}
		i++
		date := latestCommitDate(dir)
		me.print(i, date, dir, version)
	}
	return nil
}

func (me *InventoryCmd) print(i int, date, dir, version string) {
	ref := filepath.Base(dir)
	fmt.Fprintf(me.out, "%v %s %s %s\n", i, date, ref, version)
}

func (me *InventoryCmd) SetSkipUntagged(v bool) { me.skipUntagged = v }
func (me *InventoryCmd) SetPaths(v []string)    { me.paths = v }
func (me *InventoryCmd) SetOutput(v io.Writer)  { me.out = v }

func (me *InventoryCmd) Paths() []string { return me.paths }

func latestCommitDate(repodir string) string {
	date, err := exec.Command("git", "-C", repodir, "log", "-1", "--format=%ct").Output()
	if err != nil {
		log.Println(repodir, err)
		return ""
	}
	date = bytes.TrimRight(date, "\n")
	sec, _ := strconv.Atoi(string(date))
	time := time.Unix(int64(sec), 0)
	return time.Format("2006-01-02")
}

func latestVersion(repodir string) (string, error) {
	tags := tags(repodir)
	if len(tags) == 0 {
		return "v0.0.0", ErrNoTags
	}
	sort.Sort(Tags(tags))
	return tags[0], nil
}

var ErrNoTags = errors.New("no tags")

func tags(repodir string) []string {
	res := make([]string, 0)
	tags, err := exec.Command("git", "-C", repodir, "tag").Output()
	if err != nil {
		log.Println(repodir, err)
		return res
	}
	for _, tag := range bytes.Split(tags, []byte("\n")) {
		if len(tag) == 0 {
			continue
		}
		res = append(res, string(tag))
	}
	return res
}

type Tags []string

func (s Tags) Len() int      { return len(s) }
func (s Tags) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Tags) Less(i, j int) bool {
	a := semver.New(parseVersion(s[i]))
	b := semver.New(parseVersion(s[j]))
	return b.LessThan(*a)
}

func parseVersion(v string) string {
	if len(v) == 0 {
		return ""
	}
	var (
		major int
		minor int
		patch int
	)
	fmt.Sscanf(v, "v%v.%v.%v", &major, &minor, &patch)
	return fmt.Sprintf("%v.%v.%v", major, minor, patch)
}
