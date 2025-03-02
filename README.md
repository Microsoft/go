# The Go Programming Language

Go is an open source programming language that makes it easy to build simple,
reliable, and efficient software.

This repository, [microsoft/go](https://github.com/microsoft/go), contains the
infrastructure Microsoft uses to build Go. The submodule named `go` contains the
Go source code. By default, the submodule's remote URL is the official GitHub
mirror of Go, [golang/go](https://github.com/golang/go).  The canonical Git
repository for Go source code is located at https://go.googlesource.com/go.

This project is not involved in producing the [official binary distributions
of Go](https://go.dev/dl/).

Unless otherwise noted, the Go source files are distributed under the
BSD-style license found in the LICENSE file.

## Why does this fork exist?

This repository produces a modified version of Go that can be used to build FIPS
140-2 compliant applications. Our goal is to share this implementation with
others in the Go community who have the same requirement, and to merge this
capability into upstream Go as soon as possible. See
[eng/doc/fips](eng/doc/fips) for more information about this feature and the
history of FIPS 140-2 compliance in Go.

The binaries produced by this repository are also intended for general use
within Microsoft instead of the official binary distribution of Go.

We call this repository a fork even though it isn't a traditional Git fork. Its
branches do not share Git ancestry with the Go repository. However, the
repository serves the same purpose as a Git fork: maintaining a modified version
of the Go source code over time.

## Support

This project follows the upstream Go
[Release Policy](https://go.dev/doc/devel/release#policy).
This means we support each major release (1.X) until there are two newer major
releases. A new Go major version is
[released every six months](https://github.com/golang/go/wiki/Go-Release-Cycle),
so each Go major version is supported for about one year.

When upstream Go releases a new minor version (1.X.Y), we release a
corresponding microsoft/go version that may also include fork-specific changes.
This normally happens once a month. At any time, we may release a new revision
(1.X.Y-Z) to fix an issue without waiting for the next upstream minor release.
Revision releases are uncommon.

Each microsoft/go release is announced at the
[Microsoft for Go Developers](https://devblogs.microsoft.com/go/) blog.
Check out the upstream [golang-announce mailing list](https://groups.google.com/g/golang-announce)
for a summary of the changes in each Go version.

There are a few additional support resources internal to Microsoft:
* [Languages at Microsoft: Go](https://eng.ms/docs/more/languages-at-microsoft/go/articles/overview).
* [A Microsoft-internal email distribution list 📧 (instant join link)](https://idwebelements.microsoft.com/GroupManagement.aspx?Group=golang-announce&Operation=join)
  for release announcements.

## Download and install

We build the forked Go toolset with this list of OS/Arch combinations. To use a
prebuilt copy of Go while targeting a platform that is not on this list,
cross-compilation may be necessary.

* `linux_amd64`
* `linux_armv6l`
* `linux_arm64`
* `windows_amd64`

The following sections list the ways to get a build of the Microsoft fork of Go.

> [!NOTE]
> Don't see an option that works for you? Let us know!  
> File a GitHub issue, or comment on an existing issue in this tag:
  [![](https://img.shields.io/github/labels/microsoft/go/Area-Acquisition)](https://github.com/microsoft/go/labels/Area-Acquisition)

### Docker Container Images

**[microsoft/go-images](https://github.com/microsoft/go-images)** maintains and
documents container images that are available on Microsoft Artifact Registry.

### Azure Linux

The **[Azure Linux](https://github.com/microsoft/azurelinux)** distribution
includes builds of this Go fork.

* In Azure Linux 2.0, the package `msft-golang` installs this fork.
* In Azure Linux 3.0, the `golang` package installs this fork.

### Binary archive

[Signed builds of Go](https://github.com/microsoft/go/blob/microsoft/main/eng/doc/Downloads.md)
for several platforms are available as `zip` and `tar.gz` files.

### Build from source

#### Pre-patched source tarball

[The microsoft/go GitHub releases](https://github.com/microsoft/go/releases)
include a source tarball file ending in `.src.tar.gz`. After downloading and
extracting the tar.gz file, build it using the
[upstream instructions](https://go.dev/doc/install/source).

#### Clone and build

The first step is to clone this repository using Git and check out the desired
tag or commit. The `zip` file that GitHub offers for download is incomplete: it
doesn't include the `go` submodule.

If you want to contribute to the Microsoft for Go developers project, read the [Developer
Guide](eng/doc/DeveloperGuide.md). It lists the steps we recommend to set up a
Microsoft build of Go development environment, execute your first build, run
the standard library test suite, and contribute a PR.

If you just want to build on your own machine, you may find it more
convenient to use the tools provided by the `eng/run.ps1` script. We use this
script for CI builds. See [eng/README.md](eng/README.md) for more details about
`eng/run.ps1` and other repository infrastructure.

Once built, the Microsoft build of Go binary is found at `go/bin/go`.

## Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

Please read the [Developer Guide](eng/doc/DeveloperGuide.md) for more information about contributing to this project.

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft 
trademarks or logos is subject to and must follow 
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.
