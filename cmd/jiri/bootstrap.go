// Copyright 2018 The Fuchsia Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"

	"github.com/btwiuse/jiri"
	"github.com/btwiuse/jiri/cipd"
	"github.com/btwiuse/jiri/cmdline"
)

var cmdBootstrap = &cmdline.Command{
	Runner: cmdline.RunnerFunc(runBootstrap),
	Name:   "bootstrap",
	Short:  "Bootstrap essential packages",
	Long: `
Bootstrap essential packages such as cipd.
`,
	ArgsName: "<package ...>",
	ArgsLong: "<package ...> is a list of packages that can be bootstrapped by jiri. If the list is empty, jiri will list supported packages.",
}

func runBootstrap(env *cmdline.Env, args []string) error {
	if len(args) == 0 {
		// Currently it only supports cipd. We may add more packages from buildtools in the future.
		fmt.Printf("Supported package(s):\n\tcipd\n")
		return nil
	}
	for _, v := range args {
		switch strings.ToLower(v) {
		case "cipd":
			jirix, err := jiri.NewX(env)
			if err != nil {
				return err
			}
			cipdPath, err := cipd.Bootstrap(jirix.CIPDPath())
			if err != nil {
				return err
			}
			fmt.Printf("cipd bootstrapped to path:%q\n", cipdPath)

		default:
			return fmt.Errorf("unsupported package %q", v)
		}
	}
	return nil
}
