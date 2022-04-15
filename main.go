// Command inventory prints project latest release information
package main

import (
	"bytes"
	"fmt"
	"log"
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
		cli = cmdline.NewBasicParser()
	)
	cli.Parse()
	repos, _ := filepath.Glob("/home/gregory/src/github.com/gregoryv/*")
	var i int
	for _, repodir := range repos {
		i++
		v := latestVersion(repodir)
		date := latestCommitDate(repodir)
		fmt.Printf("%s %s/%s %s\n", date, "gregoryv", filepath.Base(repodir), v)
	}

	repos, _ = filepath.Glob("/home/gregory/src/xwing.7de.se/*")
	for _, repodir := range repos {
		i++
		v := latestVersion(repodir)
		date := latestCommitDate(repodir)
		fmt.Printf("%s %s/%s %s\n", date, "7de.se", filepath.Base(repodir), v)
	}
}

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

func latestVersion(repodir string) string {
	tags := tags(repodir)
	if len(tags) == 0 {
		return "v0.0.0"
	}
	sort.Sort(Tags(tags))
	return tags[0]
}

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
