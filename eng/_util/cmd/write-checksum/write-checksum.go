// Copyright (c) Microsoft Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/microsoft/go/_util/internal/checksum"
)

const description = `
This command creates a SHA256 checksum file for the given files, in the same
location and with the same name as each given file but with ".sha256" added to
the end. Pass files as non-flag arguments.

Generated files are compatible with "sha256sum -c".
`

func main() {
	help := flag.Bool("h", false, "Print this help message.")

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
	if flag.NArg() == 0 {
		flag.Usage()
		log.Fatal("No files specified.")
	}
	for _, m := range flag.Args() {
		if err := checksum.WriteSHA256ChecksumFile(m); err != nil {
			log.Fatal(err)
		}
	}
}
