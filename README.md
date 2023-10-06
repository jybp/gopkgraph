# gopkgraph

⚠️ This is experimental, performance may be slow.

A simple tool to generate a dependency graph of Golang packages. 
The output has the same format as the [go mod graph](https://go.dev/ref/mod#go-mod-graph) command. It can be used alongside tools that support that format such as [modgraphviz](https://pkg.go.dev/golang.org/x/exp/cmd/modgraphviz).

![gopkgraph](docs/docker.png?raw=true)

# Installation

- modgraphviz: `$ go install golang.org/x/exp/cmd/modgraphviz`
- graphviz: https://graphviz.org/download/

```shell
$ go install github.com/jybp/gopkgraph@latest
```

# Usage

List all dependencies of the package in the current directory belonging to the same module:

```shell
$ gopkgraph | modgraphviz | dot -Tpng -o pkgs.png
```

List all dependencies of the package in the current directory including one level of dependencies from other modules belonging to `github.com/user`:

```shell
$ gopkgraph -mods=1 | grep " \github.com/user" | modgraphviz | dot -Tpng -o pkgs.png
```

## Flags

|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Flag&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;| Description |
| --- | --- |
| `-pkg` | Path to the package. Omit this flag to target the current directory. (optional) |
| `-mods` | Max depth for packages from other modules. (optional) |
| `-stdlib` | Max depth for packages from the stdlib. (optional) |
| `-help` | Print flags. (optional) |
