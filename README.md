<!-- generated by Test_generate_readme, DO NOT EDIT! -->

inventory - Command for listing projects and their versions

This command is used to find and display a set of git projects and
their release information.

## Quick start

```shell
$ go install github.com/gregoryv/inventory
$ inventory -h
Usage: ./inventory [OPTIONS]

List projects and release information

Options
    -s, --skip-untagged
    -f, --show-full-path
    -m, --show-modified-date
    -o, --order-by : "releaseDate" [path releaseDate]
    -i, --include-vendor
    -h, --help

Examples
    List all your projects
    $ inventory

    List specific projects
    $ inventory $HOME/src/github.com/YOURS/*

```
