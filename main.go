// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goki/gi/gi"
	"github.com/goki/gi/gimain"
	"github.com/goki/gi/giv"
	"github.com/goki/ki/ki"
)

const (
	width   = 1920
	height  = 1080
	version = "v0.0.0-dev"
)

var (
	vp                                                          *gi.Viewport2D
	win                                                         *gi.Window
	lns                                                         *giv.TableView
	gstru, params                                               *giv.StructView
	mfr, statusBar                                              *gi.Frame
	fpsText, valueText, errorText, versionText, currentFileText *gi.Label
)

func main() {
	CheckForFolders()
	gimain.Main(func() {
		mainrun()
	})
}

func mainrun() {
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

	mfr = win.SetMainFrame()
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

// HandleError checks if there is an error. If there is, it sets the error text to the error, and returns true. Otherwise returns false.
func HandleError(err error) bool {
	if err != nil {
		TheGraph.State.Error = err
		if statusBar != nil {
			updt := statusBar.UpdateStart()
			errorText.SetText("Error: " + TheGraph.State.Error.Error())
			statusBar.UpdateEnd(updt)
		}
		return true
	}
	return false
}

// GetVersion finds the locally installed version and returns it
func GetVersion() string {
	b, err := os.ReadFile(filepath.Join(GetMarblesFolder(), "localData/version.txt"))
	if err != nil {
		os.WriteFile(filepath.Join(GetMarblesFolder(), "localData/version.txt"), []byte(version), os.ModePerm)
		return version
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
	TheGraph.Objects.Graph.SetProp("background-color", TheSettings.ColorSettings.GraphColor)
	// Set the text color of the status bar
	statusBar.SetProp("color", TheSettings.ColorSettings.StatusTextColor)
	// Set the color of the graph axis
	TheGraph.Objects.XAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)
	TheGraph.Objects.YAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)
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
	if TheGraph.State.File == "" {
		currentFileText.SetText("untitled.json")
	}
	strs := strings.Split(TheGraph.State.File, "savedGraphs")
	for k, d := range strs {
		if k != 1 {
			continue
		}
		d = filepath.Base(d)
		currentFileText.SetText(d)
	}
	statusBar.UpdateEnd(updt)

}

// CheckForFolders makes sure that all needed folders exist
func CheckForFolders() {
	_, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	os.MkdirAll(filepath.Join(GetMarblesFolder(), "localData"), 0777)
	savedGraphsPath := filepath.Join(GetMarblesFolder(), "savedGraphs")
	os.MkdirAll(savedGraphsPath, 0777)
	// If the directory is empty, add a blank file to prevent the app from crashing due to no files.
	if CheckIfEmpty(savedGraphsPath) {
		os.WriteFile(filepath.Join(savedGraphsPath, "blank.json"), []byte(""), os.ModePerm)
	}
	// If a file has been added, we no longer need the blank file
	if CheckIfDirHas(savedGraphsPath, 2) {
		os.Remove(filepath.Join(savedGraphsPath, "blank.json"))
	}
}

// GetMarblesFolder returns the folder where marbles data is located
func GetMarblesFolder() string {
	d, _ := os.UserConfigDir()
	return filepath.Join(d, "Marbles")
}

// CheckIfEmpty checks if a directory is empty
func CheckIfEmpty(name string) bool {
	return !CheckIfDirHas(name, 1)
}

// CheckIfDirHas checks if a directory has at least n files
func CheckIfDirHas(name string, n int) bool {
	files, _ := os.ReadDir(name)
	if len(files) >= n {
		return true
	}
	return false
}
