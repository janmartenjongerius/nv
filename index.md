# Introduction

More than ever, modern development relies on environment variables. To easily debug the local environment or export it,
a chain of commands specific to your operating system would do. However, `env` wants to solve this in a modern way,
cross platform.

# Features

- Get an environment variable, ensuring the environment variable exists.
- Search for environment variables, interactively and programmatically.
- Export a list of required environment variables to a DotEnv file.
- Set, update and unset environment variables programmatically.

# Download

Current version: [![Stable version](https://img.shields.io/github/v/release/johmanx10/env?include_prereleases&label=)](https://github.com/johmanx10/env/releases/latest)

## <img alt="Debian linux" src=https://simpleicons.org/icons/debian.svg width=20 /> Debian linux 

| Package                       | Download |
|:------------------------------|:---------|
| Env      | [![DEB 32 Bit](https://img.shields.io/badge/-32--Bit-c60036?logo=debian)](https://github.com/johmanx10/env/releases/latest/download/env_386.deb) [![DEB 64 Bit](https://img.shields.io/badge/-64--Bit-c60036?logo=debian)](https://github.com/johmanx10/env/releases/latest/download/env_amd64.deb) |
| Plugins  | [![DEB 32 Bit](https://img.shields.io/badge/-32--Bit-c60036?logo=debian)](https://github.com/johmanx10/env/releases/latest/download/plugins_386.deb) [![DEB 64 Bit](https://img.shields.io/badge/-64--Bit-c60036?logo=debian)](https://github.com/johmanx10/env/releases/latest/download/plugins_amd64.deb) |

> The application is installed in `/usr/local/bin/env` so as not to conflict with `/usr/bin/env`.

> Plugins are installed in `/usr/local/lib/env`.

## <img alt="Redhat linux" src=https://simpleicons.org/icons/redhat.svg width=20 /> Redhat linux

Coming soon

## <img alt="Microsoft Windows" src=https://simpleicons.org/icons/windows.svg width=20 /> Microsoft Windows

Coming soon

## <img alt="Golang" src=https://simpleicons.org/icons/go.svg width=20 /> Install from source

[![Source](https://img.shields.io/badge/dynamic/json.svg?label=Source&url=https://api.github.com/repos/johmanx10/env&query=$.default_branch&logo=go&color=00acd7&logoColor=7fd5ea)](https://github.com/johmanx10/env/archive/main.zip)

```
go install janmarten.name/env
```

For builds from source, please refer to
[Go platform specific information](https://github.com/golang/go/wiki#platform-specific-information).

<img alt="Linux" src=https://simpleicons.org/icons/linux.svg width=20 />
<img alt="ChromeOS" src=https://simpleicons.org/icons/googlechrome.svg width=20 />
<img alt="Darwin" src=https://simpleicons.org/icons/apple.svg width=20 />
<img alt="FreeBSD" src=https://simpleicons.org/icons/freebsd.svg width=20 />
and more.
