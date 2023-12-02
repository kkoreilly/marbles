// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/gimain"
	"goki.dev/grr"
)

func main() { gimain.Run(app) }

func app() {
	gi.SetAppName("marbles")
	gi.SetAppAbout("marbles allows you to enter equations, which are graphed, and then marbles are dropped down on the resulting lines, and bounce around in very entertaining ways!")

	b := gi.NewBody()

	TheSettings.Defaults()
	TheGraph.Init(b)

	b.NewWindow().Run().Wait()
}

// TODO(kai/marbles): better error handling
func HandleError(err error) bool {
	return grr.Log(err) != nil
}
