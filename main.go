// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/gimain"
	"goki.dev/gi/v2/giv"
	"goki.dev/ki/v2"
)

var (
	vp                                                          *gi.Viewport2D
	win                                                         *gi.Window
	lns                                                         *giv.TableView
	params                                                      *giv.StructView
	mainSplit                                                   *gi.SplitView
	graphToolbar                                                *gi.ToolBar
	mfr, statusBar                                              *gi.Frame
	fpsText, valueText, errorText, versionText, currentFileText *gi.Label
)

func main() { gimain.Run(app) }

func app() {
	TheSettings.Get()
	TheGraph.Defaults()
	InitBasicFunctionList()
	rec := ki.Node{} // receiver for events
	rec.InitName(&rec, "rec")

	gi.SetAppName("marbles")
	gi.SetAppAbout("marbles allows you to enter equations, which are graphed, and then marbles are dropped down on the resulting lines, and bounce around in very entertaining ways!")

	win = gi.NewMainWindow("marbles", "Marbles", width, height)

	vp = win.WinViewport2D()
	updt := vp.UpdateStart()

	makeBasicElements()
	TheGraph.SetFunctionsTo(DefaultFunctions)
	InitCoords()
	TheGraph.CompileExprs()
	TheGraph.Lines.Graph()
	ResetMarbles()
	SetCompleteWords(TheGraph.Functions)
	UpdateColors()

	InitClipboard()
	go InitDB()

	inClosePrompt := false
	win.SetCloseReqFunc(func(w *gi.Window) {
		if inClosePrompt {
			return
		}
		TheGraph.Stop()
		TheGraph.AutoSave()
		if !TheSettings.ConfirmQuit {
			gi.Quit()
		}
		gi.PromptDialog(vp, gi.DlgOpts{Title: "Close", Prompt: "Close marbles app?"}, gi.AddOk, gi.AddCancel, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			if sig == int64(gi.DialogAccepted) {
				gi.Quit()
			} else {
				inClosePrompt = false
			}
		})
	})
	makeMainMenu()
	win.MainMenuUpdated()
	vp.UpdateEndNoSig(updt)
	win.StartEventLoop()
}
