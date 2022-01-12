// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux
// +build linux

package openssl

import (
	"fmt"
	"os"
	"testing"
)

// Test that func init does not panic.
func TestInit(t *testing.T) {}

func TestMain(m *testing.M) {
	if !Enabled() {
		fmt.Sprintln("skipping for non-FIPS enabled machines")
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// Test that Unreachable panics.
func TestUnreachable(t *testing.T) {
	defer func() {
		if Enabled() {
			if err := recover(); err == nil {
				t.Fatal("expected Unreachable to panic")
			}
		} else {
			if err := recover(); err != nil {
				t.Fatalf("expected Unreachable to be a no-op")
			}
		}
	}()
	Unreachable()
}

// Test that UnreachableExceptTests does not panic (this is a test).
func TestUnreachableExceptTests(t *testing.T) {
	UnreachableExceptTests()
}
