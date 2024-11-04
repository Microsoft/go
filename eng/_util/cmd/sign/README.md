# `sign` and the Microsoft Go signing infrastructure

Most of the logic for signing (extracting files, repackaging, creating checksums) is implemented by this `sign` command.

The [`/eng/signing`](/eng/signing) directory contains the MSBuild project that `sign` invokes to run real signing.
The MSBuild project uses [MicroBuild Signing](https://dev.azure.com/devdiv/DevDiv/_wiki/wikis/DevDiv.wiki/650/MicroBuild-Signing) (internal Microsoft wiki link).

To see signing in action, go to [`/eng/pipeline/README.md`](/eng/pipeline/README.md) and follow the link for `microsoft-go`.

## Dry run

1. Create the directory `/eng/signing/tosign` and add the `.tar.gz` and `.zip` artifacts to sign.
    * Download artifacts from the `microsoft-go` pipeline, for example.
    * It's ok to skip downloading some artifacts. The signing process doesn't require all platforms to be present.
    * If you specify `-files`, you can use your own directory.
1. From the root of the repository, run `pwsh eng/run.ps1 sign -n`

The `-n` argument makes it a dry run: it extracts/repacks files in the same way it would if it were signing them, but no signing is done.
This doesn't involve .NET/MSBuild, so this is a good way for a developer to test changes to the signing logic.

See `pwsh eng/run.ps1 sign -h` for more options.

## Test signing

> [!NOTE]
> Test signing has not been observed to work.
> It has been documented for completeness, in case someone wants to try.

### Prerequisites

* Windows
* .NET Core SDK 8.0 or later.
    * [Download](https://dot.net/download)
* The signing plugin.
    1. Download the latest NuGet Package: https://devdiv.visualstudio.com/DevDiv/_artifacts/feed/MicroBuildToolset/NuGet/MicroBuild.Plugins.Signing
    1. Extract its contents (the file is a zip) to `%userprofile%\.nuget\packages\microbuild.plugins.signing\1.1.900`.
        * Optionally make the versioned dir's name match the version of the package you downloaded. It will be discovered dynamically, as a plugin, whether or not the version matches.

### Test signing run

1. Set up `tosign` as described in the dry run section.
1. From the root of the repository, run `pwsh eng/run.ps1 sign`

## Real signing

This can't be done from a dev machine.
It occurs in the `microsoft-go` pipeline, on a Windows machine.
See [`/eng/pipeline/README.md`](/eng/pipeline/README.md).

The invocation of `sign` can be found in [`/eng/pipeline/stages/sign-stage.yml`](/eng/pipeline/stages/sign-stage.yml).
