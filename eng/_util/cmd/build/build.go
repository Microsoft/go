// Copyright (c) Microsoft Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/microsoft/go-infra/json2junit"
	"github.com/microsoft/go-infra/patch"
	"github.com/microsoft/go-infra/submodule"
	"github.com/microsoft/go/_util/buildutil"
)

const description = `
This command builds Go, optionally running tests and packing an archive file.

Use this script to build Go on your local machine in the way the Microsoft
infrastructure builds it. "eng/run.ps1 build" automatically downloads a copy of
the Go compiler (required to build Go) then starts the build. This script is
also capable of running tests and packing an archive file: see Usage, above.

To build and test Go without the Microsoft infrastructure, use the Bash scripts
in 'src' such as 'src/run.bash' instead of this script.

Example: Build Go, run tests, and produce an archive file:

  eng/run.ps1 build -test -packbuild
`

func main() {
	var help = flag.Bool("h", false, "Print this help message.")
	o := &options{}

	flag.BoolVar(&o.SkipBuild, "skipbuild", false, "Disable building Go.")
	flag.BoolVar(&o.Test, "test", false, "Enable running tests.")
	flag.BoolVar(&o.PackBuild, "packbuild", false, "Enable creating an archive of this build using upstream 'distpack' and placing it in eng/artifacts/bin.")
	flag.BoolVar(&o.PackSource, "packsource", false, "Enable creating a source archive using upstream 'distpack' and placing it in eng/artifacts/bin.")
	flag.BoolVar(&o.CreatePDB, "pdb", false, "Create PDB files for all the PE binaries in the bin and tool directories. The PE files are modified in place and PDBs are placed in eng/artifacts/symbols.")

	flag.BoolVar(
		&o.Refresh, "refresh", false,
		"Refresh Go submodule: clean untracked files, reset tracked files, and apply patches before building.\n"+
			"For more refresh options, use the top level 'submodule-refresh' command instead of 'build'.")

	flag.StringVar(&o.Experiment, "experiment", "", "Include this string in GOEXPERIMENT.")
	flag.StringVar(&o.JUnitOutFile, "junitout", "", "Write the test output to this path as a JUnit file if this builder runs tests.")

	o.MaxMakeAttempts = buildutil.MaxMakeRetryAttemptsOrExit()

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

	// If build returns an error, handle it here with panic. Having build return an error makes it
	// easier to adapt build in the future to somewhere else in the module to use it as an API. (For
	// example, "build" could be changed to "Build" and run-builder could use it.)
	if err := build(o); err != nil {
		panic(err)
	}
}

type options struct {
	SkipBuild    bool
	Test         bool
	PackBuild    bool
	PackSource   bool
	CreatePDB    bool
	Refresh      bool
	Experiment   string
	JUnitOutFile string

	MaxMakeAttempts int
}

func build(o *options) (err error) {

	scriptExtension := ".bash"
	executableExtension := ""
	archiveExtension := ".tar.gz"
	shellPrefix := []string{"bash"}

	if runtime.GOOS == "windows" {
		scriptExtension = ".bat"
		executableExtension = ".exe"
		archiveExtension = ".zip"
		shellPrefix = []string{"cmd.exe", "/c"}
	}

	// eng/run.ps1 guarantees that the current working directory is the root of the Go repo (our
	// GOROOT). Keep track of this so we can optionally pack it up later.
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if o.Refresh {
		config, err := patch.FindAncestorConfig(rootDir)
		if err != nil {
			return err
		}
		if err := submodule.Reset(rootDir, filepath.Join(config.RootDir, config.SubmoduleDir), true); err != nil {
			return err
		}
		if err := patch.Apply(config, patch.ApplyModeIndex); err != nil {
			return err
		}
	}

	// Get the target platform information. If the environment variable is different from the
	// runtime value, this means we're doing a cross-compiled build. These values are used for
	// capability checks and to make sure that if Pack is enabled, the output archive is formatted
	// correctly and uses the right filename.
	targetOS, err := buildutil.GetEnvOrDefault("GOOS", runtime.GOOS)
	if err != nil {
		return err
	}
	targetArch, err := buildutil.GetEnvOrDefault("GOARCH", runtime.GOARCH)
	if err != nil {
		return err
	}
	fmt.Printf("---- Target platform: %v_%v\n", targetOS, targetArch)

	if err := buildutil.UnassignGOROOT(); err != nil {
		return err
	}

	// The upstream build scripts in {repo-root}/src require your working directory to be src, or
	// they instantly fail. Change the current process dir so that we can run them.
	if err := os.Chdir("go/src"); err != nil {
		return err
	}

	if o.Experiment != "" {
		buildutil.AppendExperimentEnv(o.Experiment)
	}

	if !o.SkipBuild {
		// If we have a stage 0 copy of Go in an env variable (as set by run.ps1), use it in the
		// build command by setting GOROOT_BOOTSTRAP. The upstream build script "make.bash" uses
		// this env variable to find the copy of Go to use to build.
		//
		// Forcing the build script to use our stage 0 avoids uncertainty that could occur if we
		// allowed it to use arbitrary versions of Go from the build machine PATH.
		//
		// To avoid this behavior and use an ambiently installed version of Go from PATH, run
		// "make.bash" manually instead of using this tool.
		if stage0Goroot := os.Getenv("STAGE_0_GOROOT"); stage0Goroot != "" {
			if err := os.Setenv("GOROOT_BOOTSTRAP", stage0Goroot); err != nil {
				return err
			}
		}

		// Set GOBUILDEXIT so 'make.bat' exits with exit code upon failure. The ordinary behavior of
		// 'make.bat' is to always end with 0 exit code even if an error occurred, so 'all.bat' can
		// handle the error. See https://github.com/golang/go/issues/7806.
		if err := os.Setenv("GOBUILDEXIT", "1"); err != nil {
			return err
		}

		buildCommandLine := append(shellPrefix, "make"+scriptExtension)

		if err := buildutil.Retry(o.MaxMakeAttempts, func() error {
			return runCommandLine(buildCommandLine...)
		}); err != nil {
			return err
		}

		// The race runtime requires cgo.
		// It isn't supported on arm or 386.
		// It's supported on arm64, but the official linux-arm64 distribution doesn't include it.
		if os.Getenv("CGO_ENABLED") != "0" && targetArch != "arm" && targetArch != "arm64" && targetArch != "386" {
			fmt.Println("---- Building race runtime...")
			err := runCommandLine(
				filepath.Join("..", "bin", "go"+executableExtension),
				"install", "-race", "-a", "std",
			)
			if err != nil {
				return err
			}
		}
	}

	if o.Test {
		// Normally, use the dev script to build.
		testCommandLine := append(
			shellPrefix,
			[]string{
				"run" + scriptExtension,
				"--no-rebuild",
			}...,
		)

		if o.JUnitOutFile != "" {
			testCommandLine = append(testCommandLine, "-json")
			f, err := os.Create(o.JUnitOutFile)
			if err != nil {
				return err
			}
			conv := json2junit.NewConverter(f)
			defer func() {
				if closeErr := conv.Close(); err == nil {
					err = closeErr
				}
				if closeErr := f.Close(); err == nil {
					err = closeErr
				}
			}()
			if err := buildutil.RunCmdMultiWriter(testCommandLine, conv, buildutil.NewStripTestJSONWriter(os.Stdout)); err != nil {
				return err
			}
		} else {
			if err := buildutil.RunCmdMultiWriter(testCommandLine, os.Stdout); err != nil {
				return err
			}
		}
	}

	goRootDir := filepath.Join(rootDir, "go")
	if o.CreatePDB {
		if _, err := exec.LookPath("gopdb"); err != nil {
			return fmt.Errorf("gopdb not found in PATH: %v", err)
		}
		// Print the version of gopdb to the console.
		cmd := exec.Command("gopdb", "-version")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := runCmd(cmd); err != nil {
			return fmt.Errorf("gopdb failed: %v", err)
		}

		// Traverse the bin and tool directories to find all the binaries to generate PDBs for.
		binDir := filepath.Join(goRootDir, "bin")
		toolsDir := filepath.Join(goRootDir, "pkg", "tool", targetOS+"_"+targetArch)
		artifactsPDBDir := filepath.Join(rootDir, "eng", "artifacts", "symbols")

		if err := os.MkdirAll(artifactsPDBDir, os.ModePerm); err != nil {
			return err
		}

		var bins []string
		for _, dir := range []string{binDir, toolsDir} {
			entries, err := os.ReadDir(dir)
			if err != nil {
				return err
			}
			for _, entry := range entries {
				if !entry.Type().IsRegular() {
					continue
				}
				bins = append(bins, filepath.Join(dir, entry.Name()))
			}
		}

		// Generate PDBs for all the binaries.
		for _, bin := range bins {
			out := filepath.Join(artifactsPDBDir, filepath.Base(bin)+"."+targetOS+"-"+targetArch+".pdb")
			cmd := exec.Command("gopdb", "-o", out, bin)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := runCmd(cmd); err != nil {
				return fmt.Errorf("gopdb failed: %v", err)
			}
		}
	}

	if o.PackBuild || o.PackSource {
		// Find the host version of distpack. (Not the target version, which might not run.)
		toolsDir := filepath.Join(goRootDir, "pkg", "tool", runtime.GOOS+"_"+runtime.GOARCH)
		// distpack needs a VERSION file to run. If we're on the main branch, we don't have one, so
		// use dist's version calculation to create a temp dev version and put it in VERSION.
		var version string
		if data, err := os.ReadFile(filepath.Join(goRootDir, "VERSION")); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if version, err = writeDevelVersionFile(goRootDir, toolsDir); err != nil {
					return fmt.Errorf("unable to pack: failed writing development VERSION file: %v", err)
				}
				// Best effort: clean up the VERSION file when we're done. This is just for dev
				// workflows: the temp VERSION file should never be checked in.
				defer os.Remove(filepath.Join(goRootDir, "VERSION"))
			} else {
				return fmt.Errorf("unable to pack: VERSION file in unexpected state: %v", err)
			}
		} else {
			version, _, _ = strings.Cut(string(data), "\n")
		}
		cmd := exec.Command(filepath.Join(toolsDir, "distpack"+executableExtension))
		cmd.Env = append(os.Environ(), "GOROOT="+goRootDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := runCmd(cmd); err != nil {
			return fmt.Errorf("distpack failed: %v", err)
		}
		// distpack creates some files we don't need. Recreate the naming logic here to pick out the
		// files we want and copy them to our artifacts dir.
		distPackDir := filepath.Join(goRootDir, "pkg", "distpack")
		artifactsBinDir := filepath.Join(rootDir, "eng", "artifacts", "bin")
		type packCopy struct{ src, dst string }
		var packs []packCopy
		// Insert the build ID to make sure the archive filename is unique. We might change
		// patches but build the same submodule commit multiple times.
		buildID := getBuildID()
		if o.PackBuild {
			// distpack calls GOARCH=arm "arm" in its tar.gz filename, but the upstream release
			// process changes it to "armv6l" on https://go.dev/dl/ to match the historical name.
			// Do the same here.
			brandingTargetArch := targetArch
			if brandingTargetArch == "arm" {
				brandingTargetArch = "armv6l"
			}
			packs = append(packs, packCopy{
				src: filepath.Join(distPackDir, version+"."+targetOS+"-"+targetArch+archiveExtension),
				dst: filepath.Join(artifactsBinDir, version+"-"+buildID+"."+targetOS+"-"+brandingTargetArch+archiveExtension),
			})
		}
		if o.PackSource {
			packs = append(packs, packCopy{
				src: filepath.Join(distPackDir, version+".src.tar.gz"),
				dst: filepath.Join(artifactsBinDir, version+"-"+buildID+".src.tar.gz"),
			})
		}
		fmt.Printf("---- Copying distpack output to artifacts dir %v\n", artifactsBinDir)
		for _, p := range packs {
			fmt.Printf("---- Copying %q to %q...\n", p.src, p.dst)
			if err := copyFile(p.dst, p.src); err != nil {
				return err
			}
		}
	}

	fmt.Printf("---- Build command complete.\n")
	return nil
}

func writeDevelVersionFile(goRootDir, toolsDir string) (string, error) {
	cmd := exec.Command(filepath.Join(toolsDir, "dist"), "version")
	cmd.Env = append(os.Environ(), "GOROOT="+goRootDir)
	vBytes, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("unable to get dist version: %v (%v)", err, string(vBytes))
	}
	fields := strings.Fields(string(vBytes))
	if len(fields) < 2 {
		return "", fmt.Errorf("expected at least 2 fields in dist version output, got %q in %q", len(fields), string(vBytes))
	}
	if fields[0] != "devel" {
		return "", fmt.Errorf("expected first field 'devel' in dist version, got %q", fields[0])
	}
	// The second field should be something like "go1.21-abcde1234", and the remaining fields are a
	// timestamp. Just using the second field as is: the full VERSION file string is placed into the
	// archive filename, so this keeps it simple and avoids special characters.
	if err := os.WriteFile(filepath.Join(goRootDir, "VERSION"), []byte(fields[1]), 0o666); err != nil {
		return "", err
	}
	return fields[1], nil
}

// copyFile copies src to dst, creating dst's directory if necessary. Handles errors robustly,
// see https://github.com/golang/go/blob/c3458e35f4/src/cmd/internal/archive/archive_test.go#L57
// Doesn't copy file permissions.
func copyFile(dst, src string) (err error) {
	err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}
	var s, d *os.File
	s, err = os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err = os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if e := d.Close(); err == nil {
			err = e
		}
	}()
	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}
	return nil
}

func runCommandLine(commandLine ...string) error {
	c := exec.Command(commandLine[0], commandLine[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return runCmd(c)
}

func runCmd(cmd *exec.Cmd) error {
	fmt.Printf("---- Running command: %v\n", cmd.Args)
	return cmd.Run()
}

// getBuildID returns BUILD_BUILDNUMBER if defined (e.g. a CI build). Otherwise, "dev".
func getBuildID() string {
	archiveVersion := os.Getenv("BUILD_BUILDNUMBER")
	if archiveVersion == "" {
		return "dev"
	}
	return archiveVersion
}
