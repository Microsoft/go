# Developer Guide

This document is a guide for developers who want to contribute to the Microsoft Go repository.
It explains how to build the repository, how to work with the Go submodule, and how to use the different tools
that help maintain the repository.

## Setting up the repository

### Step 0: Contributor License Agreement

Most contributions require you to agree to a Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us the rights to use your contribution.
For details, visit https://cla.opensource.microsoft.com.

### Step 1: Install a Go toolchain

A preexisting Go toolchain is required to bootstrap the build process.
You can use your system's package manager to install Go, or you can download it from the [official Go website](https://golang.org/dl/).
The only requirement is that the Go version is high enough for the bootstrap process.
If the version is too low, the bootstrap process will fail and ask you to install a newer version.

### Step 2: Install the git-go-patch command

The `git-go-patch` command is a tool that helps you manage the patches in the `go` submodule.

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
> The program still works if you call it with its real name, but we think it's easier to remember and type something that looks like a `git` subcommand.

### Step 2: Initialize the submodule and apply patches

The repository uses a `go` submodule to store the Go source code.
All the patches that modify the Go source code are stored in the `patches` directory.

To initialize the submodule and apply the patches, run the following command:

```
git go-patch apply
```

### Step 3: Build the Go toolchain

You now can edit the `go/src` directory as you would any other Go project.

Note that the `go/src/go.mod` file uses a `go` directive with a version that is not available as a prebuilt binary.
This is because that module contains the future version of Go, which is not yet released.
In order to build the standard library packages located in `go/src` you will first need to build to Go toolchain from the `go/src` directory itself using the following command:

```
cd go/src
./make.bash # or make.bat on Windows
```

After building the Go toolchain, you can use it to develop and test changes in the `go/src` directory.
You will need to add the `go/bin` directory to your `PATH` to use the new Go toolchain.

### Step 4: Test that your environment is set up correctly

To test that your environment is set up correctly, run the following command:

```
cd go/src
go version
go test -short ./...
```

## Making changes to go/src

TODO