// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

//go:build linux

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// enableSystemWideFIPS fallback is a no-op because the current platform either doesn't support or
// doesn't require system-wide FIPS to be enabled to run tests.
func enableSystemWideFIPS() (restore func(), err error) {
	cmd := exec.Command("openssl", "version", "-a")
	log.Printf("---- Running command: %v\n", cmd.Args)
	out, err := cmd.CombinedOutput()
	sout := string(out)
	if err != nil {
		return nil, fmt.Errorf("failed to check openssl version: %v, %v", err, string(out))
	}
	log.Print(sout)

	lines := strings.Split(sout, "\n")
	if !strings.Contains(sout, "OpenSSL 1.") {
		// Only OpenSSL 1 needs special handling for FIPS mode,
		// at least on the platforms we test on.
		log.Println("Using fallback (no-op) for enableSystemWideFIPS. It either isn't supported on this platform or isn't necessary.")
		//return nil, nil
	}

	// Search for the OPENSSLDIR path in the output.
	var ossldir string
	for _, line := range lines {
		var found bool
		if ossldir, found = strings.CutPrefix(line, "OPENSSLDIR: "); found {
			break
		}
	}
	if ossldir == "" {
		return nil, fmt.Errorf("failed to find OPENSSLDIR in openssl version output")
	}
	ossldir = strings.Trim(ossldir, `"`)

	// Append the FIPS configuration to the openssl.cnf file.
	// OpenSSL will merge duplicated sections, so we don't need
	// to check if the section already exists.
	opensslcnf := filepath.Join(ossldir, "openssl.cnf")
	prevContent, err := os.ReadFile(opensslcnf)
	if err != nil {
		return nil, fmt.Errorf("failed to read openssl.cnf file: %v", err)
	}
	err = os.WriteFile(opensslcnf, append(prevContent, []byte("\n\n[evp_sect]\nfips_mode = yes\n")...), 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write to openssl.cnf file: %v", err)
	}

	log.Println("Enabled FIPS mode.")

	return func() {
		err := os.WriteFile(opensslcnf, prevContent, 0644)
		if err != nil {
			log.Printf("Unable to restore openssl.cnf file: %v\n", err)
			return
		}
		log.Println("Successfully restored openssl.cnf file.")
	}, nil
}
