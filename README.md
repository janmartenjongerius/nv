# Introduction

Nv is an alternative lookup tool for Environment Variables.
It can be used to:

- Get an environment variable, ensuring the environment variable exists.
- Search for environment variables, interactively and programmatically.
- Export a list of required environment variables to a DotEnv file.

Nv is written in Go and builds into a native binary.

* See the [user documentation](http://janmarten.name/nv) for more information and pre-built downloads.
* Use the [code documentation](https://pkg.go.dev/janmarten.name/nv) for specifics about code symbols.

# Installation

Either install the application through Go:

```
go install janmarten.name/nv
```

Or clone the repository and install from source:

```
git clone git@github.com:johmanx10/nv.git
cd nv
make bin/nv
```

For builds from source, please refer to
[Go platform specific information](https://github.com/golang/go/wiki#platform-specific-information).


# Project status

| Statistic        | State |
|:-----------------|:------|
| Downloads        | ![GitHub Releases](https://img.shields.io/github/downloads/johmanx10/nv/latest/total?label=Latest%20release) [![Downloads](https://img.shields.io/github/downloads/johmanx10/nv/total?label=All%20releases)](https://github.com/johmanx10/nv/releases) |
| Workflows        | [![Build](https://github.com/johmanx10/nv/workflows/Build/badge.svg)](https://github.com/johmanx10/nv/actions?query=workflow%3ABuild) [![Release](https://github.com/johmanx10/nv/workflows/Release/badge.svg)](https://github.com/johmanx10/nv/actions?query=workflow%3ARelease) [![Documentation](https://github.com/johmanx10/nv/workflows/Documentation/badge.svg)](https://github.com/johmanx10/nv/actions?query=workflow%3ADocumentation) [![GitHub deployments](https://img.shields.io/github/deployments/johmanx10/nv/github-pages?label=Github%20pages&logo=jekyll)](https://github.com/johmanx10/nv/deployments/activity_log?environment=github-pages) |
| Current version  | [![Stable version](https://img.shields.io/github/v/release/johmanx10/nv?include_prereleases&label=)](https://github.com/johmanx10/nv/releases/latest) ![GitHub Release Date](https://img.shields.io/github/release-date/johmanx10/nv?label=Released) ![GitHub commits since latest release (by date)](https://img.shields.io/github/commits-since/johmanx10/nv/latest?label=Commits%20behind&logo=git&logoColor=white) |
| Go version       | [![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/johmanx10/nv?logo=go&label=&logoColor=white)](https://golang.org/doc/go1.14) |
| Source code size | [![Code size in bytes](https://img.shields.io/github/languages/code-size/johmanx10/nv?label=&logo=git&logoColor=white)](https://github.com/johmanx10/nv) |
| Code quality     | [![Go Report Card](https://goreportcard.com/badge/github.com/johmanx10/nv)](https://goreportcard.com/report/github.com/johmanx10/nv) [![Codecov](https://img.shields.io/codecov/c/github/johmanx10/nv?label=Coverage&logo=codecov&logoColor=white)](https://codecov.io/gh/johmanx10/nv) |
| License          | [![License](https://img.shields.io/github/license/johmanx10/nv?label=&color=blue)](https://github.com/johmanx10/nv/blob/main/LICENSE) |

# Special thanks

| Supporter | Support |
|----:|:---:|
| [![JetBrains logo](docs/img/jetbrains.png)](https://www.jetbrains.com/?from=NV) |  This project is supported by [JetBrains](https://www.jetbrains.com/?from=NV) through providing a license for their [Go IDE, Goland](https://www.jetbrains.com/go/?from=NV). <br/><br/>[![Goland IDE](docs/img/logo-goland.png)](https://www.jetbrains.com/go/?from=NV) |

