# Developer Guide

This document is a guide for developers who want to contribute to the Microsoft Go repository.
It explains how to build the repository, how to work with the Go submodule, and how to use the different tools that help maintain the repository.

This guide is primarily intended for developers working for the Go team at Microsoft, but it can also be useful for external contributors.

## Setting up the repository

### Contributor License Agreement

Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution.
For details, visit https://cla.opensource.microsoft.com.

### Install a Go toolchain

A preexisting Go toolchain is required to bootstrap the build process.
You can use your system's package manager to install Go, or you can download it from the [official Go website](https://golang.org/dl/).
The only requirement is that the Go version is high enough for the bootstrap process.
If the version is too low, the bootstrap process will fail and ask you to install a newer version.

This repository implements some scripts (provided by `eng/run.ps1`) to facilitate installing the correct bootstrapping Go version and also to build the Go toolchain from source, see the [`eng` Readme](../eng/README.md) for more information.
It is recommended that you get familiar with both the upstream Go build process and the scripts provided in this repository.

### Install git and the git-go-patch command

This repository heavily relies on advanced Git features to manage the Go submodule, so it is recommended to develop with a local clone of the repository rather than using the GitHub web interface.

You will need to have Git installed on your system, either from your system's package manager or from the [official Git website](https://git-scm.com/downloads).

The [`git-go-patch`](https://github.com/microsoft/go-infra/tree/main/cmd/git-go-patch) command is a tool that helps you manage the patches in the `go` submodule.

To install the `git-go-patch` command, run the following command:

```
go install github.com/microsoft/go-infra/cmd/go-patch@latest
```

> [!NOTE]
> Make sure `git-go-patch` is accessible in your shell's `PATH` variable.
> You may need to add `$GOPATH/bin` to your `PATH`. Use `go env GOPATH` to locate it.

Then, run the command to see the help documentation:

```
git go-patch -h
```

> [!NOTE]
> `git` detects that our `git-go-patch` executable starts with `git-` and makes it available as `git go-patch`.

### Initialize the submodule and apply patches

The repository uses a [Git submodule](https://git-scm.com/book/en/v2/Git-Tools-Submodules) named `go` to store the Go source code.
All the patches that modify the Go source code are stored in the [`patches`](../../patches) directory.

To initialize the submodule and apply the patches, run the following command:

```
git go-patch apply
```

### Build the Go toolchain

You now can edit the `go/src` directory as you would the upstream Go project.
[The upstream "Installing Go from source" instructions](https://go.dev/doc/install/source) apply to the `go` directory and can be used to build and test.

In order to make changes to the standard library packages located in `go/src` you will first need to build to Go toolchain from the `go/src` directory itself using the following command:

```
cd go/src
./make.bash # or make.bat on Windows
```

> [!NOTE]
> Rebuilding the Go toolchain from source is not necessary for changes in the Go standard library, they are immediately reflected in any future `go build`, `go test`, or `go run` commands.
> However, if you are making changes to the Go toolchain itself (any package under `go/src/cmd`), you will need to rebuild the Go toolchain.

The newly built Go toolchain will be available in the `go/bin` directory. From now one this guide will assume that any mention of the `go` command refers to the one in the `go/bin` directory.
There are different ways to use the new Go toolchain:
- Add `go/bin` to your `PATH`, although but it is not recommended because it will probably contain unstable features that may interfere with other Go projects.
- You can use the full path to the `go` command in the `go/bin` directory.
- You can instruct your IDE to use the `go` command in the `go/bin` directory (recommended approach). See the [IDE setup](#ide-setup) section for more information.

### Test that your environment is set up correctly

To test that your environment is set up correctly, run the following command:

```
cd go/src
go version
go test -short ./...
```

## IDE setup

### Visual Studio Code

Visual Studio Code (VS Code from now on) is a popular IDE for Go development. We recommend using the official Go extension for VS Code.
Please refer to the [Go extension documentation](https://code.visualstudio.com/docs/languages/go) for more information on how to set up VS Code for Go development.

#### Using the Go toolchain from the `go` submodule

You can use the Go toolchain from the `go` submodule in VS Code by following these steps:

1. In VS Code, open `Command Palette's Help` > `Show All Commands`. Or use the keyboard shortcut (`Ctrl+Shift+P`).
1. Search for `Go: Choose Go environment` then run the command from the pallet.
1. Select `Choose from file browser`.
1. Select the `go` command in the `go/bin` directory.
1. Save the file and restart VS Code.

## Making changes to go/src

TODO