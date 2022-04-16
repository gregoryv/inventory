package main

import (
	"github.com/coreos/go-semver/semver"
)

type Project struct {
	lastModified string
	path         string
	latest       Tag
}

func (me *Project) SetLastModified(v string) { me.lastModified = v }
func (me *Project) SetPath(v string)         { me.path = v }
func (me *Project) SetLatest(v Tag)          { me.latest = v }

func (me *Project) LastModified() string { return me.lastModified }
func (me *Project) Path() string         { return me.path }
func (me *Project) Version() string      { return me.latest.name }
func (me *Project) ReleaseDate() string  { return me.latest.date }

// ----------------------------------------

type byReleaseDate []Project

func (s byReleaseDate) Len() int      { return len(s) }
func (s byReleaseDate) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byReleaseDate) Less(i, j int) bool {
	return s[i].ReleaseDate() < s[j].ReleaseDate()
}

type byPath []Project

func (s byPath) Len() int      { return len(s) }
func (s byPath) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byPath) Less(i, j int) bool {
	return s[i].Path() < s[j].Path()
}

// ----------------------------------------

type Tags []Tag

func (s Tags) Len() int      { return len(s) }
func (s Tags) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Tags) Less(i, j int) bool {
	a, err := semver.NewVersion(s[i].name)
	if err != nil {
		return false
	}

	b, err := semver.NewVersion(s[j].name)
	if err != nil {
		return false
	}
	return b.LessThan(*a)
}

// ----------------------------------------

type Tag struct {
	name string
	date string
}
