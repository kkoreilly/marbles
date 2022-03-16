// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/goki/gi/gi"
	"github.com/goki/gi/gimain"
	"github.com/goki/gi/giv"
	"github.com/goki/gi/svg"
	"github.com/goki/ki/ki"
	"github.com/goki/mat32"
	"github.com/inconshreveable/go-update"
)

const ( // Width and height of the window, and size of the graph
	width  = 1920
	height = 1080
)

var (
	vp                                                *gi.Viewport2D
	win                                               *gi.Window
	eqTable, lns                                      *giv.TableView
	paramsEdit, gstru                                 *giv.StructView
	svgGraph                                          *svg.SVG
	svgLines, svgMarbles, svgCoords, svgTrackingLines *svg.Group
	gmin, gmax, gsz, ginc                             mat32.Vec2
	mfr, statusBar                                    *gi.Frame
	fpsText, errorText, versionText, currentFileText  *gi.Label
	problemWithEval, problemWithCompile               bool
)

func main() {
	err := updateApp()
	if err != nil {
		panic(err)
	}
	gimain.Main(func() {
		mainrun()
	})
}

func updateApp() error {
	fmt.Println("Loading executable file... ")
	resp, err := http.Get("https://github.com/kplat1/marblesInfo/releases/download/v0.0-beta/marblesLinux")
	if err != nil {
		return err
	}
	fmt.Println("Applying update... ")
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		return err
	}
	fmt.Println("Finished!")
	return err
}

func mainrun() {
	TheSettings.Get()
	Gr.Defaults()
	InitEquationChangeSlice()
	InitBasicFunctionList()
	rec := ki.Node{} // receiver for events
	rec.InitName(&rec, "rec")

	gi.SetAppName("marbles")
	gi.SetAppAbout("marbles allows you to enter equations, which are graphed, and then marbles are dropped down on the resulting lines, and bounce around in very entertaining ways!")

	win = gi.NewMainWindow("marbles", "Marbles", width, height)

	vp = win.WinViewport2D()
	updt := vp.UpdateStart()

	mfr = win.SetMainFrame()
	makeBasicElements()
	InitCoords()
	ResetMarbles()
	Gr.CompileExprs()
	Gr.Lines.Graph(false)
	UpdateColors()

	InitDB()

	inClosePrompt := false
	win.SetCloseReqFunc(func(w *gi.Window) {
		if inClosePrompt {
			return
		}
		Gr.Stop()
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

// HandleError checks if there is an error. If there is, it sets the error text to the error, and returns true. Otherwise returns false.
func HandleError(err error) bool {
	if err != nil {
		errorText.SetText("Error: " + err.Error())
		return true
	}
	return false
}

// GetVersion finds the locally installed version and returns it
func GetVersion() string {
	b, err := os.ReadFile("localData/version.txt")
	if HandleError(err) {
		return "Error getting version"
	}
	return string(b)
}

// UpdateColors sets the colors of the app as specified in settings
func UpdateColors() {
	// Set the background color of the app
	mfr.SetProp("background-color", TheSettings.ColorSettings.BackgroundColor)
	lns.SetProp("background-color", TheSettings.ColorSettings.BackgroundColor)
	gstru.SetProp("background-color", TheSettings.ColorSettings.BackgroundColor)
	// Set the background color of the status bar and graph
	statusBar.SetProp("background-color", TheSettings.ColorSettings.StatusBarColor)
	errorText.CurBgColor = TheSettings.ColorSettings.StatusBarColor
	fpsText.CurBgColor = TheSettings.ColorSettings.StatusBarColor
	currentFileText.CurBgColor = TheSettings.ColorSettings.StatusBarColor
	svgGraph.SetProp("background-color", TheSettings.ColorSettings.GraphColor)
	// Set the text color of the status bar
	statusBar.SetProp("color", TheSettings.ColorSettings.StatusTextColor)
	// Set the color of the graph axis
	xAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)
	yAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)
	// Set the text color of the graph and line controls
	lns.SetProp("color", TheSettings.ColorSettings.LineTextColor)
	gstru.SetProp("color", TheSettings.ColorSettings.GraphTextColor)
	// Set the background color and button color for the toolbar
	tb := gstru.ToolBar()
	tb.SetProp("background-color", TheSettings.ColorSettings.ToolBarColor)
	children := tb.Children()
	for _, d := range []ki.Ki(*children) {
		d.SetProp("background-color", TheSettings.ColorSettings.ToolBarButtonColor)
	}
	// Set the background color for the graph parameters
	gstru.StructGrid().SetProp("background-color", TheSettings.ColorSettings.GraphParamsColor)
	// Set the background color for the lines
	lFrame := lns.ChildByName("frame", -1)
	lFrame.SetProp("background-color", TheSettings.ColorSettings.LinesBackgroundColor)
	lFrame.ChildByName("header", -1).SetProp("background-color", TheSettings.ColorSettings.LinesBackgroundColor)
	lFrame.ChildByName("grid-lay", -1).ChildByName("grid", -1).SetProp("background-color", TheSettings.ColorSettings.LinesBackgroundColor)

}

// UpdateCurrentFileText updates the current file text
func UpdateCurrentFileText() {
	updt := statusBar.UpdateStart()
	if currentFile == "" {
		currentFileText.SetText("untitled.json")
	}
	strs := strings.Split(currentFile, "savedGraphs")
	for k, d := range strs {
		if k != 1 {
			continue
		}
		d = strings.ReplaceAll(d, `\`, "")
		d = strings.ReplaceAll(d, `/`, "")
		currentFileText.SetText(d)
	}
	statusBar.UpdateEnd(updt)

}
