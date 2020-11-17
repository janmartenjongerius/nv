# Introduction

More than ever, modern development relies on environment variables. To easily debug the local environment or export it,
a chain of commands specific to your operating system would do. However, Nv wants to solve this in a modern way,
cross platform.

# Download

## 64-Bit

| Operating system | Package |
|:-----------------|:--------|
| Debian Linux     | [![Download](https://img.shields.io/github/v/release/johmanx10/nv?include_prereleases&label=&color=c60036&logo=debian)](https://github.com/johmanx10/nv/releases/latest/download/nv_amd64.deb) |

## Install from source

Either install the application through Go:

```
go install janmarten.name/nv
```

Or clone the repository and install from source:

```
git clone git@github.com:johmanx10/nv.git
cd nv
make install
```

For builds from source, please refer to
[Go platform specific information](https://github.com/golang/go/wiki#platform-specific-information).
