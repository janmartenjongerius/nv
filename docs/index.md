# Introduction

More than ever, modern development relies on environment variables. To easily debug the local environment or export it,
a chain of commands specific to your operating system would do. However, Nv wants to solve this in a modern way,
cross platform.

# Download

Current version: [![Stable version](https://img.shields.io/github/v/release/johmanx10/nv?include_prereleases&label=)](https://github.com/johmanx10/nv/releases/latest)

## Debian linux 

[![DEB 64-bit](https://img.shields.io/badge/Nv-64--Bit-c60036?logo=debian)](https://github.com/johmanx10/nv/releases/latest/download/nv_amd64.deb)

### Plugins
- [![DEB 64 Bit](https://img.shields.io/badge/Encoding:%20XML-64--Bit-c60036?logo=debian)](https://github.com/johmanx10/nv/releases/latest/download/nv_ext_encoding-xml.deb)

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
