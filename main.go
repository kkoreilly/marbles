// Copyright (c) 2020, kplat1. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math/rand"
	"os"
	"strings"
	"time"

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
	graphSize = 800
)

var (
	vp                                               *gi.Viewport2D
	eqTable, lns                                     *giv.TableView
	paramsEdit, gstru                                *giv.StructView
	svgGraph                                         *svg.SVG
	svgLines, svgMarbles, svgCoords                  *svg.Group
	gmin, gmax, gsz, ginc                            mat32.Vec2
	mfr, statusBar                                   *gi.Frame
	fpsText, errorText, versionText, currentFileText *gi.Label
	problemWithEval                                  bool
)

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

	vp = win.WinViewport2D()
	updt := vp.UpdateStart()

	mfr = win.SetMainFrame()
	// the StructView will also show the Graph Toolbar which is main actions..
	gstru = giv.AddNewStructView(mfr, "gstru")
	gstru.Viewport = vp // needs vp early for toolbar
	gstru.SetProp("height", "4.5em")
	gstru.SetStruct(&Gr)
	paramsEdit = gstru
	lns = giv.AddNewTableView(mfr, "lns")
	lns.Viewport = vp
	lns.SetSlice(&Gr.Lines)
	eqTable = lns

	frame := gi.AddNewFrame(mfr, "frame", gi.LayoutHoriz)

	svgGraph = svg.AddNewSVG(frame, "graph")
	svgGraph.SetProp("min-width", graphSize)
	svgGraph.SetProp("min-height", graphSize)
	svgGraph.SetStretchMaxWidth()
	svgGraph.SetStretchMaxHeight()

	svgLines = svg.AddNewGroup(svgGraph, "SvgLines")
	svgMarbles = svg.AddNewGroup(svgGraph, "SvgMarbles")
	svgCoords = svg.AddNewGroup(svgGraph, "SvgCoords")

	gmin = mat32.Vec2{X: -10, Y: -10}
	gmax = mat32.Vec2{X: 10, Y: 10}
	gsz = gmax.Sub(gmin)
	ginc = gsz.DivScalar(graphSize)

	svgGraph.ViewBox.Min = gmin
	svgGraph.ViewBox.Size = gsz
	svgGraph.Norm = true
	svgGraph.InvertY = true
	svgGraph.Fill = true
	svgGraph.SetProp("background-color", "white")
	svgGraph.SetProp("stroke-width", ".2pct")

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
	currentFileText = gi.AddNewLabel(statusBar, "currentFileText", "untitled.json")
	currentFileText.SetProp("font-weight", "bold")
	currentFileText.SetStretchMaxWidth()
	currentFileText.Redrawable = true
	versionText = gi.AddNewLabel(statusBar, "versionText", "")
	versionText.SetProp("font-weight", "bold")
	versionText.SetStretchMaxWidth()
	versionText.SetText("Running version " + GetVersion())
	// viewSettingsButton = gi.AddNewButton(statusBar, "viewSettingsButton")
	// viewSettingsButton.SetText("Settings")
	// viewSettingsButton.OnClicked(func() {
	// 	pSettings := TheSettings
	// 	giv.StructViewDialog(vp, &TheSettings, giv.DlgOpts{Title: "Settings", Ok: true, Cancel: true}, rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
	// 		if sig == int64(gi.DialogAccepted) {
	// 			TheSettings.Save()
	// 			Gr.Params.Defaults()
	// 			Gr.Graph()
	// 			UpdateColors()
	// 			ResetMarbles()
	// 		} else if sig == int64(gi.DialogCanceled) {
	// 			TheSettings = pSettings
	// 		}
	// 	})
	// })
	// treeview := giv.AddNewTreeView(statusBar, "treeview")
	// treeview.SetRootNode(gstru)
	lns.ChildByName("toolbar", -1).Delete(true)
	gstru.ChildByName("toolbar", -1).ChildByName("UpdtView", -1).Delete(true)

	InitCoords()
	ResetMarbles()
	Gr.CompileExprs()
	Gr.Lines.Graph(false)
	UpdateColors()

	InitDB()

	// Main Menu

	appnm := gi.AppName()
	mmen := win.MainMenu
	mmen.ConfigMenus([]string{appnm, "File", "Edit", "Window", "Objectives"})

	fmen := win.MainMenu.ChildByName("File", 0).(*gi.Action)
	fmen.Menu = make(gi.Menu, 0, 10)
	fmen.Menu.AddAction(gi.ActOpts{Label: "New", ShortcutKey: gi.KeyFunMenuNew}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		Gr.Reset()
	})
	fmen.Menu.AddSeparator("sep0")
	fmen.Menu.AddAction(gi.ActOpts{Label: "Open", ShortcutKey: gi.KeyFunMenuOpen}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		giv.FileViewDialog(vp, "savedGraphs/", ".json", giv.DlgOpts{}, nil, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			if sig == int64(gi.DialogAccepted) {
				dlg := send.Embed(gi.KiT_Dialog).(*gi.Dialog)
				Gr.OpenJSON(gi.FileName(giv.FileViewDialogValue(dlg)))
			}
		})
	})
	fmen.Menu.AddAction(gi.ActOpts{Label: "Open Autosave", ShortcutKey: gi.KeyFunMenuOpenAlt1}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		Gr.OpenAutoSave()
	})
	fmen.Menu.AddSeparator("sep1")
	fmen.Menu.AddAction(gi.ActOpts{Label: "Save", ShortcutKey: gi.KeyFunMenuSave}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if currentFile != "" {
			Gr.SaveLast()
		} else {
			giv.FileViewDialog(vp, "savedGraphs/", ".json", giv.DlgOpts{}, nil, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
				if sig == int64(gi.DialogAccepted) {
					dlg := send.Embed(gi.KiT_Dialog).(*gi.Dialog)
					Gr.SaveJSON(gi.FileName(giv.FileViewDialogValue(dlg)))
				}
			})
		}
	})
	fmen.Menu.AddAction(gi.ActOpts{Label: "Save as", ShortcutKey: gi.KeyFunMenuSaveAs}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		giv.FileViewDialog(vp, "savedGraphs/", ".json", giv.DlgOpts{}, nil, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			if sig == int64(gi.DialogAccepted) {
				dlg := send.Embed(gi.KiT_Dialog).(*gi.Dialog)
				Gr.SaveJSON(gi.FileName(giv.FileViewDialogValue(dlg)))
			}
		})
	})
	fmen.Menu.AddSeparator("sep2")
	fmen.Menu.AddAction(gi.ActOpts{Label: "Settings", ShortcutKey: gi.KeyFunMenuSaveAlt}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		pSettings := TheSettings
		giv.StructViewDialog(vp, &TheSettings, giv.DlgOpts{Title: "Settings", Ok: true, Cancel: true}, rec.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
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

	amen := win.MainMenu.ChildByName(appnm, 0).(*gi.Action)
	amen.Menu = make(gi.Menu, 0, 10)
	amen.Menu.AddAppMenu(win)

	emen := win.MainMenu.ChildByName("Edit", 1).(*gi.Action)
	emen.Menu = make(gi.Menu, 0, 10)
	emen.Menu.AddCopyCutPaste(win)

	omen := win.MainMenu.ChildByName("Objectives", 2).(*gi.Action)
	omen.Menu = make(gi.Menu, 0, 10)
	omen.Menu.AddAction(gi.ActOpts{Label: "Generate Objectives"}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		rand.Seed(time.Now().UnixNano())
		nObjs := rand.Intn(3) + 3
		for i := 0; i < nObjs; i++ {
			Gr.AddObjective()
		}

	})
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

}
