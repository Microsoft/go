# MSBuild signing infrastructure

This directory contains a component of the Microsoft Go signing infrastructure written using MSBuild.
`Sign.csproj` is the interface between the Go signing command [`/eng/_util/cmd/sign`][sign] and MicroBuild, an internal Microsoft toolset written to primarily support .NET projects that use MSBuild.

See [`/eng/_util/cmd/sign`][sign] for more information about the signing infrastructure.

[sign]: /eng/_util/cmd/sign