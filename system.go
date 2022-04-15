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
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type System struct {
	skipUntagged     bool
	showModifiedDate bool

	root  string
	paths []string
	out   io.Writer
}

func (me *System) SetSkipUntagged(v bool)     { me.skipUntagged = v }
func (me *System) SetShowModifiedDate(v bool) { me.showModifiedDate = v }
func (me *System) SetPaths(v []string)        { me.paths = v }
func (me *System) SetOutput(v io.Writer)      { me.out = v }
func (me *System) SetRoot(v string)           { me.root = v }

func (me *System) Run() error {
	if len(me.paths) == 0 {
		// find all project directories
		filepath.Walk(me.root, func(pth string, f os.FileInfo, err error) error {
			if f == nil || !f.IsDir() {
				return nil
			}
			if path.Base(pth) == "pkg" {
				return filepath.SkipDir
			}
			if f.Name() == ".git" {
				me.paths = append(me.paths, filepath.Dir(pth))
				return filepath.SkipDir
			}
			return nil
		})
	}

	result := make([]Project, 0)
	for _, dir := range me.paths {
		tag, err := latestTag(dir)
		if errors.Is(ErrNoTags, err) && me.skipUntagged {
			continue
		}
		date := latestCommitDate(dir)
		var p Project
		p.SetLastModified(date)
		p.SetPath(dir)
		p.SetLatest(tag)
		result = append(result, p)
	}

	me.format(result)
	return nil
}

func (me *System) format(result []Project) {
	w := me.out
	fmt.Fprintf(w, "# Showing %v of %v projects\n", len(result), len(me.paths))
	fmt.Fprintln(w, me.Header())

	for _, p := range result {
		path := filepath.Base(p.Path())
		parts := []string{p.ReleaseDate(), path, p.Version()}
		if me.showModifiedDate {
			parts = append(parts, p.LastModified())
		}
		fmt.Fprintln(w, strings.Join(parts, " "))
	}
}

func (me *System) Header() string {
	var buf bytes.Buffer
	buf.WriteString("# Released Version Project")
	if me.showModifiedDate {
		buf.WriteString(" Modified")
	}
	return buf.String()
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

func latestTag(repodir string) (Tag, error) {
	tags := tags(repodir)
	if len(tags) == 0 {
		return NoTag, ErrNoTags
	}
	sort.Sort(Tags(tags))
	return tags[0], nil
}

var NoTag = Tag{
	name: "v0.0.0",
	date: "0000-00-00",
}
var ErrNoTags = errors.New("no tags")

func tags(repodir string) []Tag {
	res := make([]Tag, 0)
	tags, err := exec.Command("git", "-C", repodir,
		"for-each-ref", "--format=%(tag) %(taggerdate:short)", "refs/tags",
	).Output()
	if err != nil {
		log.Println(repodir, err)
		return res
	}
	for _, entry := range bytes.Split(tags, []byte("\n")) {
		if len(entry) == 0 {
			continue
		}
		fields := strings.Fields(string(entry))
		if len(fields) != 2 {
			continue
		}
		t := Tag{
			name: fields[0],
			date: fields[1],
		}
		res = append(res, t)
	}
	return res
}
