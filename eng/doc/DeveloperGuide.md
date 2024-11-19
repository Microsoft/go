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
You can use your system's package manager to install Go, download Go from the [official Go website](https://golang.org/dl/), or download a prebuilt version of Microsoft Go itself.

The only requirement is that the Go version is high enough for the bootstrap process.
If you attempt to build Go while using a bootstrap Go with a version that is too low, the bootstrap process will fail and ask you to install a newer version.

> [!NOTE]
> The in-support versions of Go found on the [official Go website](https://golang.org/dl/) are always high enough to bootstrap the development branch.
> This is because:
> * The [last two major versions of Go are supported by the Go project](https://go.dev/s/release#release-maintenance). (Microsoft Go has the same policy.)
> * Go N can always be bootstrapped by [both N-1 and N-2](https://tip.golang.org/doc/install/source#go14).

> [!NOTE]
> This repository's `eng/run.ps1` PowerShell script is able to download a correct bootstrapping Go version automatically before building Microsoft Go from source.
> We recommend that Microsoft Go team members be familiar with this script because it is used by Microsoft Go CI.
> However, it isn't necessary to use the script for most work on the Microsoft Go patches.
> See the [`eng` Readme](/eng/README.md) for more information about `eng/run.ps1`.

### Install git and the git-go-patch command

This repository heavily relies on advanced Git features to manage the Go submodule, so it is recommended to develop with a local Git clone of the repository rather than other methods, e.g. using the GitHub web interface.

Make sure Git is installed on your system.
You can get Git from your system's package manager or the [official Git website](https://git-scm.com/downloads).

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
We recommend reading the upstream instructions, but we've included some minimal instructions here to get started.

First, use the following commands to build the Go toolchain using the source in the `go/src` directory:

- On Unix-like systems:
    ```bash
    cd go/src
    ./make.bash
    ```

- On Windows:
    ```bat
    cd go/src
    .\make.bat
    ```

The newly built Go toolchain will be available in the `go/bin` directory.
An app built by `go/bin/go` will use the standard library in `go/src`, so changes that you make to the standard library are reflected in the built app.

From now on, when this guide mentions the `go` command, it refers to executing the `go` binary in the `go/bin` directory.

> [!NOTE]
> Rebuilding the Go toolchain from source is not necessary for changes in the Go standard library: changes are immediately reflected in any `go build`, `go test`, or `go run` commands.
> However, if you make changes to the Go toolchain itself (any package under `go/src/cmd`), you do need to rebuild the Go toolchain.

There are different ways to use the new Go toolchain:

- Use the full path to the `go` command.
- Add the full path of `go/bin` to the start of `PATH`.
  - We only recommend setting `PATH` in a specific terminal session, not user-wide or system-wide. The development version of Go will probably contain unstable features that may interfere with your other Go projects.
- Instruct your IDE to use the `go` command. Recommended approach for most development work. See the [IDE setup](#ide-setup) section for more information.

### Test that your environment is set up correctly

To test that your environment is set up correctly, run the following commands, which work the same on all platforms:

```
cd go/src
go version
go test -short ./...
```

## IDE setup

### VS Code

[VS Code](https://code.visualstudio.com/) (Visual Studio Code) is a popular IDE for Go development.
We recommend using the official Go extension for VS Code.
Please refer to the [Go extension documentation](https://code.visualstudio.com/docs/languages/go) for more information on how to set up VS Code for Go development.

#### Using the Go toolchain from the `go` submodule

You can use your build of `go` in VS Code by following these steps:

1. In VS Code, open the command palette.
    - `View` > `Command Palette...`.
    - Default keyboard shortcut: `Ctrl+Shift+P`.
1. Search for `Go: Choose Go environment` and select it.
1. Select `Choose from file browser`.
1. Select the `go` executable in the `go/bin` directory. (On Windows, `go.exe`.)
1. Open the command palette.
1. Search for `Developer: Reload Window` and select it.

## Making changes to go/src

TODO