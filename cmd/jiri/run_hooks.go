// Copyright 2017 The Fuchsia Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/btwiuse/jiri"
	"github.com/btwiuse/jiri/cmdline"
	"github.com/btwiuse/jiri/project"
)

var runHooksFlags struct {
	localManifest bool
	hookTimeout   uint
	attempts      uint
	fetchPackages bool
}

var cmdRunHooks = &cmdline.Command{
	Runner: jiri.RunnerFunc(runHooks),
	Name:   "run-hooks",
	Short:  "Run hooks using local manifest",
	Long: `
Run hooks using local manifest JIRI_HEAD version if -local-manifest flag is
false, else it runs hooks using current manifest checkout version.
`,
}

func init() {
	cmdRunHooks.Flags.BoolVar(&runHooksFlags.localManifest, "local-manifest", false, "Use local checked out manifest.")
	cmdRunHooks.Flags.UintVar(&runHooksFlags.hookTimeout, "hook-timeout", project.DefaultHookTimeout, "Timeout in minutes for running the hooks operation.")
	cmdRunHooks.Flags.UintVar(&runHooksFlags.attempts, "attempts", 1, "Number of attempts before failing.")
	cmdRunHooks.Flags.BoolVar(&runHooksFlags.fetchPackages, "fetch-packages", true, "Use fetching packages using jiri.")
}

func runHooks(jirix *jiri.X, args []string) (err error) {
	localProjects, err := project.LocalProjects(jirix, project.FastScan)
	if err != nil {
		return err
	}
	if runHooksFlags.attempts < 1 {
		return jirix.UsageErrorf("Number of attempts should be >= 1")
	}
	jirix.Attempts = runHooksFlags.attempts

	// Get hooks.
	var projs project.Projects
	var hooks project.Hooks
	var pkgs project.Packages
	if !runHooksFlags.localManifest {
		projs, hooks, pkgs, err = project.LoadUpdatedManifest(jirix, localProjects, runHooksFlags.localManifest)
	} else {
		projs, hooks, pkgs, err = project.LoadManifestFile(jirix, jirix.JiriManifestFile(), localProjects, runHooksFlags.localManifest)
	}
	if err != nil {
		return err
	}
	if err = project.RunHooks(jirix, hooks, runHooksFlags.hookTimeout); err != nil {
		return err
	}
	if err := project.FilterOptionalProjectsPackages(jirix, jirix.FetchingAttrs, nil, pkgs); err != nil {
		return err
	}
	// Get packages if the fetchPackages is true
	if runHooksFlags.fetchPackages && len(pkgs) > 0 {
		// Extend timeout for packages to be 5 times the timeout of a single hook.
		return project.FetchPackages(jirix, projs, pkgs, runHooksFlags.hookTimeout*5)
	}
	return nil
}
