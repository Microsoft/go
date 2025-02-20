// Copyright (c) Microsoft Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package buildutil

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Retry runs f until it succeeds or the attempt limit is reached.
func Retry(attempts int, f func() error) error {
	var i = 0
	for ; i < attempts; i++ {
		if attempts > 1 {
			fmt.Printf("---- Running attempt %v of %v...\n", i+1, attempts)
		}
		err := f()
		if err != nil {
			if i+1 < attempts {
				fmt.Printf("---- Attempt failed with error: %v\n", err)
				continue
			}
			fmt.Printf("---- Final attempt failed.\n")
			return err
		}
		break
	}
	fmt.Printf("---- Successful on attempt %v of %v.\n", i+1, attempts)
	return nil
}

// MaxMakeRetryAttemptsOrExit returns max retry attempts for the Go build according to an env var.
func MaxMakeRetryAttemptsOrExit() int {
	return maxAttemptsOrExit("GO_MAKE_MAX_RETRY_ATTEMPTS")
}

func maxAttemptsOrExit(varName string) int {
	attempts, err := getEnvIntOrDefault(varName, 1)
	if err != nil {
		log.Fatal(err)
	}
	if attempts <= 0 {
		log.Fatalf("Expected positive integer for environment variable %q, but found: %v\n", varName, attempts)
	}
	return attempts
}

func getEnvIntOrDefault(varName string, defaultValue int) (int, error) {
	a, err := GetEnvOrDefault(varName, strconv.Itoa(defaultValue))
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(a)
	if err != nil {
		return 0, fmt.Errorf("env var %q is not an int: %w", varName, err)
	}
	return i, nil
}

// GetEnvOrDefault find an environment variable with name varName and returns its value. If the env
// var is not set, returns defaultValue.
//
// If the env var is found and its value is empty string, returns an error. This can't happen on
// Windows because setting an env var to empty string deletes it. However, on Linux, it is possible.
// It's likely a mistake, so we let the user know what happened with an error. For example, the env
// var might be empty string because it was set by "example=$(someCommand)" and someCommand
// encountered an error and didn't send any output to stdout.
func GetEnvOrDefault(varName, defaultValue string) (string, error) {
	v, ok := os.LookupEnv(varName)
	if !ok {
		return defaultValue, nil
	}
	if v == "" {
		return "", fmt.Errorf(
			"env var %q is empty, not a valid string. To use the default string %v, unset the env var",
			varName, defaultValue)
	}
	return v, nil
}

// AppendExperimentEnv sets the GOEXPERIMENT env var to the given value, or if GOEXPERIMENT is
// already set, appends a comma separator and then the given value.
func AppendExperimentEnv(experiment string) {
	// If the experiment enables a crypto backend, allow fallback to Go crypto. Go turns off cgo
	// and/or cross-builds in various situations during the build/tests, so we need to allow for it.
	if strings.Contains(experiment, "opensslcrypto") ||
		strings.Contains(experiment, "cngcrypto") ||
		strings.Contains(experiment, "boringcrypto") ||
		strings.Contains(experiment, "darwincrypto") ||
		strings.Contains(experiment, "systemcrypto") {

		experiment += ",allowcryptofallback"
	}
	if v, ok := os.LookupEnv("GOEXPERIMENT"); ok {
		experiment = v + "," + experiment
	}
	fmt.Printf("Setting GOEXPERIMENT: %v\n", experiment)
	if err := os.Setenv("GOEXPERIMENT", experiment); err != nil {
		panic(err)
	}
}

// UnassignGOROOT unsets the GOROOT env var if it is set.
//
// Setting GOROOT explicitly in the environment has not been necessary since Go
// 1.9 (https://go.dev/doc/go1.9#goroot), but a dev or build machine may still
// have it set. It interferes with attempts to run the built Go (such as when
// building the race runtime), so remove the explicit GOROOT if set.
func UnassignGOROOT() error {
	if explicitRoot, ok := os.LookupEnv("GOROOT"); ok {
		fmt.Printf("---- Removing explicit GOROOT from environment: %v\n", explicitRoot)
		if err := os.Unsetenv("GOROOT"); err != nil {
			return err
		}
	}
	return nil
}

// RunCmdMultiWriter runs a command and outputs the stdout to multiple [io.Writer].
// The writers are closed after the command completes.
func RunCmdMultiWriter(cmdline []string, stdout ...io.Writer) (err error) {
	c := exec.Command(cmdline[0], cmdline[1:]...)
	c.Stdout = io.MultiWriter(stdout...)
	c.Stderr = os.Stderr

	fmt.Printf("---- Running command: %v\n", c.Args)
	return c.Run()
}
