// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

//go:generate core generate

import (
	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/core"
)

func app() {
	TheSettings.Defaults()

	b := core.NewBody("Marbles")
	b.AddTopBar(func(bar *core.Frame) {
		core.NewToolbar(bar).Maker(TheGraph.MakeToolbar)
	})

	TheGraph.Init(b)

	b.RunMainWindow()
}

// TODO(kai/marbles): better error handling
func HandleError(err error) bool {
	return errors.Log(err) != nil
}
