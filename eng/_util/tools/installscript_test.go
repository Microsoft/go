// Copyright (c) Microsoft Corporation.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutil

import (
	"os/exec"
	"testing"
)

func TestInstallScriptUpToDate(t *testing.T) {
	cmd := exec.Command("go", "run", "github.com/microsoft/go-infra/goinstallscript", "-check")
	cmd.Dir = ".."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("goinstallscript is not up to date: %v, %v", string(out), err)
		t.Errorf("To update, in eng/_util, run: go run github.com/microsoft/go-infra/goinstallscript")
		t.Fail()
	}
}
