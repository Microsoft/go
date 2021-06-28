// Copyright (c) Microsoft Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/microsoft/go/_core/pack"
)

const description = `
This command packs a built Go directory into an archive file and produces a
checksum file for the archive. It filters out the files that aren't necessary.
`

func main() {
	repoRootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	source := flag.String("source", repoRootDir, "The path of the Go directory to archive.")
	output := flag.String("o", "", "The path of the archive file to create. Format depends on extension. Default: a GOOS/GOARCH-dependent archive file in 'eng/artifacts/bin'.")

	var help = flag.Bool("h", false, "Print this help message.")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n", description)
	}

	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	if err := pack.Archive(*source, *output); err != nil {
		panic(err)
	}
}
