// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/goki/gi/gi"
	"github.com/goki/gi/gimain"
	"github.com/goki/gi/giv"
	"github.com/goki/gi/svg"
	"github.com/goki/ki/ki"
	"github.com/goki/mat32"
)

const ( // Width and height of the window, and size of the graph
	width     = 1920
	height    = 1080
	GraphSize = 800
)

var Vp *gi.Viewport2D
var EqTable *giv.TableView
var ParamsEdit *giv.StructView
var SvgGraph *svg.SVG
var SvgLines *svg.Group
var SvgMarbles *svg.Group
var SvgCoords *svg.Group
var gmin, gmax, gsz, ginc mat32.Vec2
var statusBar *gi.Frame
var fpsText *gi.Label
var errorText *gi.Label
var versionText *gi.Label
var mfr *gi.Frame
var gstru *giv.StructView
var lns *giv.TableView
var problemWithEval = false
var viewSettingsButton *gi.Button

func main() {
	gimain.Main(func() {
		mainrun()
	})
}

func mainrun() {
	TheSettings.Get()
	Gr.Defaults()
	InitEquationChangeSlice()
	rec := ki.Node{} // receiver for events
	rec.InitName(&rec, "rec")

	gi.SetAppName("marblesApp")
	gi.SetAppAbout("marbles allows you to enter equations, which are graphed, and then marbles are dropped down on the resulting lines, and bounce around in very entertaining ways!")

	win := gi.NewMainWindow("marblesApp", "Marbles", width, height)

	Vp = win.WinViewport2D()
	updt := Vp.UpdateStart()

	mfr = win.SetMainFrame()
	// the StructView will also show the Graph Toolbar which is main actions..
	gstru = giv.AddNewStructView(mfr, "gstru")
	gstru.Viewport = Vp // needs vp early for toolbar
	gstru.SetProp("height", "4.5em")
	gstru.SetStruct(&Gr)
	ParamsEdit = gstru
	lns = giv.AddNewTableView(mfr, "lns")
	lns.Viewport = Vp
	lns.SetSlice(&Gr.Lines)
	EqTable = lns

	frame := gi.AddNewFrame(mfr, "frame", gi.LayoutHoriz)

	SvgGraph = svg.AddNewSVG(frame, "graph")
	SvgGraph.SetProp("min-width", GraphSize)
	SvgGraph.SetProp("min-height", GraphSize)
	SvgGraph.SetStretchMaxWidth()
	SvgGraph.SetStretchMaxHeight()

	SvgLines = svg.AddNewGroup(SvgGraph, "SvgLines")
	SvgMarbles = svg.AddNewGroup(SvgGraph, "SvgMarbles")
	SvgCoords = svg.AddNewGroup(SvgGraph, "SvgCoords")

	gmin = mat32.Vec2{X: -10, Y: -10}
	gmax = mat32.Vec2{X: 10, Y: 10}
	gsz = gmax.Sub(gmin)
	ginc = gsz.DivScalar(GraphSize)

	SvgGraph.ViewBox.Min = gmin
	SvgGraph.ViewBox.Size = gsz
	SvgGraph.Norm = true
	SvgGraph.InvertY = true
	SvgGraph.Fill = true
	SvgGraph.SetProp("background-color", "white")
	SvgGraph.SetProp("stroke-width", ".2pct")

	statusBar = gi.AddNewFrame(mfr, "statusBar", gi.LayoutHoriz)
	statusBar.SetStretchMaxWidth()
	fpsText = gi.AddNewLabel(statusBar, "fpsText", "FPS: ")
	fpsText.SetProp("font-weight", "bold")
	fpsText.SetStretchMaxWidth()
	fpsText.Redrawable = true
	errorText = gi.AddNewLabel(statusBar, "errorText", "")
	errorText.SetProp("font-weight", "bold")
	errorText.SetStretchMaxWidth()
	errorText.Redrawable = true
	versionText = gi.AddNewLabel(statusBar, "versionText", "")
	versionText.SetProp("font-weight", "bold")
	versionText.SetStretchMaxWidth()
	versionText.SetText("Running version " + GetVersion())
	viewSettingsButton = gi.AddNewButton(statusBar, "viewSettingsButton")
	viewSettingsButton.SetText("Settings")
	viewSettingsButton.OnClicked(func() {
		pSettings := TheSettings
		giv.StructViewDialog(Vp, &TheSettings, giv.DlgOpts{Title: "Settings", Ok: true, Cancel: true}, rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			if sig == int64(gi.DialogAccepted) {
				TheSettings.Save()
				Gr.Params.Defaults()
				Gr.Graph()
				UpdateColors()
				ResetMarbles()
			} else if sig == int64(gi.DialogCanceled) {
				TheSettings = pSettings
			}
		})
	})
	// treeview := giv.AddNewTreeView(statusBar, "treeview")
	// treeview.SetRootNode(lns)
	lns.ChildByName("toolbar", -1).Delete(true)

	InitCoords()
	ResetMarbles()
	Gr.CompileExprs()
	Gr.Lines.Graph()
	UpdateColors()

	// Main Menu

	appnm := gi.AppName()
	mmen := win.MainMenu
	mmen.ConfigMenus([]string{appnm, "Edit", "Window"})

	amen := win.MainMenu.ChildByName(appnm, 0).(*gi.Action)
	amen.Menu = make(gi.Menu, 0, 10)
	amen.Menu.AddAppMenu(win)

	emen := win.MainMenu.ChildByName("Edit", 1).(*gi.Action)
	emen.Menu = make(gi.Menu, 0, 10)
	emen.Menu.AddCopyCutPaste(win)

	gi.SetQuitCleanFunc(func() {
		Gr.Stop()
		gi.Quit()
	})

	win.MainMenuUpdated()
	Vp.UpdateEndNoSig(updt)
	win.StartEventLoop()
}

// Handle Error checks if there is an error. If there is, it sets the error text to the error, and returns true. Otherwise returns false.
func HandleError(err error) bool {
	if err != nil {
		errorText.SetText("Error: " + err.Error())
		return true
	}
	return false
}

// Get Vesion finds the locally installed version and returns it
func GetVersion() string {
	b, err := os.ReadFile("localData/version.txt")
	if HandleError(err) {
		return "Error getting version"
	}
	return string(b)
}

// Update Colors sets the colors of the app as specified in settings
func UpdateColors() {
	// Set the background color of the app
	mfr.SetProp("background-color", TheSettings.ColorSettings.BackgroundColor)
	lns.SetProp("background-color", TheSettings.ColorSettings.BackgroundColor)
	gstru.SetProp("background-color", TheSettings.ColorSettings.BackgroundColor)
	// Set the background color of the status bar and graph
	statusBar.SetProp("background-color", TheSettings.ColorSettings.StatusBarColor)
	errorText.CurBgColor = TheSettings.ColorSettings.StatusBarColor
	fpsText.CurBgColor = TheSettings.ColorSettings.StatusBarColor
	SvgGraph.SetProp("background-color", TheSettings.ColorSettings.GraphColor)
	// Set the text color of the status bar
	statusBar.SetProp("color", TheSettings.ColorSettings.StatusTextColor)
	// Set the color of the graph axis
	xAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)
	yAxis.SetProp("stroke", TheSettings.ColorSettings.AxisColor)
	// Set the text color of the graph and line controls
	lns.SetProp("color", TheSettings.ColorSettings.LineTextColor)
	gstru.SetProp("color", TheSettings.ColorSettings.GraphTextColor)
	//Set the color of the settings button
	viewSettingsButton.SetProp("background-color", TheSettings.ColorSettings.ButtonColor)
	// Set the background color and button color for the toolbar
	tb := gstru.ToolBar()
	tb.SetProp("background-color", TheSettings.ColorSettings.ToolBarColor)
	children := tb.Children()
	for _, d := range []ki.Ki(*children) {
		d.SetProp("background-color", TheSettings.ColorSettings.ToolBarButtonColor)
	}
	// Set the background color for the graph parameters
	gstru.StructGrid().SetProp("background-color", TheSettings.ColorSettings.GraphParamsColor)
}
