// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package main

import (
	"log"
	"os"
)

// enableSystemWideFIPS enables Mariner and Azure Linux 3 process-wide FIPS mode.
func enableSystemWideFIPS() (restore func(), err error) {
	// FIPS mode is enabled if OPENSSL_FORCE_FIPS_MODE is set, regardless of the value.
	_, ok := os.LookupEnv("OPENSSL_FORCE_FIPS_MODE")
	if ok {
		log.Println("FIPS mode already enabled.")
		return nil, nil
	}

	env("OPENSSL_FORCE_FIPS_MODE", "1")
	log.Println("Enabled Mariner and Azure Linux 3 FIPS mode.")

	return func() {
		if ok {
			err := os.Unsetenv("OPENSSL_FORCE_FIPS_MODE")
			if err != nil {
				log.Printf("Unable to unset OPENSSL_FORCE_FIPS_MODE: %v\n", err)
				return
			}
			log.Println("Successfully unset OPENSSL_FORCE_FIPS_MODE.")
		}
	}, nil
}
