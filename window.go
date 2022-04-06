package main

import (
	"github.com/goki/gi/gi"
	"github.com/goki/gi/giv"
	"github.com/goki/gi/svg"
	"github.com/goki/ki/ki"
	"github.com/goki/mat32"
)

func makeBasicElements() {

	// the StructView will also show the Graph Toolbar which is main actions..
	gstru = giv.AddNewStructView(mfr, "gstru")
	gstru.Viewport = vp // needs vp early for toolbar
	gstru.SetProp("height", "1em")
	gstru.SetStruct(&TheGraph)

	split := gi.AddNewSplitView(mfr, "split")
	split.SetProp("min-height", TheSettings.GraphSize)
	sidesplit := gi.AddNewSplitView(split, "sidesplit")
	sidesplit.Dim = mat32.Y
	lns = giv.AddNewTableView(sidesplit, "lns")
	lns.Viewport = vp
	lns.SetProp("index", false)
	lns.NoAdd = true
	lns.SetSlice(&TheGraph.Lines)
	lns.StyleFunc = func(tv *giv.TableView, slice interface{}, widg gi.Node2D, row, col int, vv giv.ValueView) {
		if col == 0 {
			newLabel := "<i><b>y=</b></i>"
			if row < len(functionNames) {
				newLabel = "<i><b>" + functionNames[row] + "(x)=</b></i>"
			}
			lbl := widg.(*giv.StructViewInline).Parts.Child(0).(*gi.Label)
			lbl.SetText(newLabel)
			lbl.SetProp("background-color", "yellow")
		}
		if col == 3 {
			clr := TheGraph.Lines[row].Colors.Color
			widg.SetProp("background-color", clr)
			widg.SetProp("color", clr)
			widg.(*gi.Action).Text = "LColors"
		}
		if col < 3 {
			edit := widg.(*giv.StructViewInline).Parts.Child(1).(*gi.TextField)
			edit.TextFieldSig.Connect(edit.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
				if col == 0 {
					TheGraph.Lines[row].Expr.Expr = string(edit.EditTxt)
				}
				if col == 1 {
					TheGraph.Lines[row].GraphIf.Expr = string(edit.EditTxt)
				}
				if col == 2 {
					TheGraph.Lines[row].Bounce.Expr = string(edit.EditTxt)
				}
				TheGraph.AutoGraph()
			})
			widg.SetProp("font-size", TheSettings.LineFontSize)
		}
	}

	params := giv.AddNewStructView(sidesplit, "params")
	params.SetStruct(&TheGraph.Params)

	sidesplit.SetSplits(6, 4)

	frame := gi.AddNewFrame(split, "frame", gi.LayoutHoriz)

	svgGraph = svg.AddNewSVG(frame, "graph")
	svgGraph.SetProp("min-width", TheSettings.GraphSize)
	svgGraph.SetProp("min-height", TheSettings.GraphSize)
	svgLines = svg.AddNewGroup(svgGraph, "SvgLines")
	svgMarbles = svg.AddNewGroup(svgGraph, "SvgMarbles")
	svgCoords = svg.AddNewGroup(svgGraph, "SvgCoords")
	svgTrackingLines = svg.AddNewGroup(svgGraph, "SvgTrackingLines")
	split.SetSplits(float32(width-TheSettings.GraphSize), float32(TheSettings.GraphSize)*7/8)
	gmin = mat32.Vec2{X: -10, Y: -10}
	gmax = mat32.Vec2{X: 10, Y: 10}
	gsz = gmax.Sub(gmin)
	var n float32 = 1.0 / float32(TheSettings.GraphInc)
	ginc = mat32.Vec2{X: n, Y: n}

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
	errorText.SetText("Graphed successfully")
	currentFileText = gi.AddNewLabel(statusBar, "currentFileText", "untitled.json")
	currentFileText.SetProp("font-weight", "bold")
	currentFileText.SetStretchMaxWidth()
	currentFileText.Redrawable = true
	versionText = gi.AddNewLabel(statusBar, "versionText", "")
	versionText.SetProp("font-weight", "bold")
	versionText.SetStretchMaxWidth()
	versionText.SetText("Running version " + GetVersion())
	lns.ChildByName("toolbar", -1).Delete(true)
	params.ChildByName("toolbar", -1).Delete(true)
	gstru.ChildByName("toolbar", -1).ChildByName("UpdtView", -1).Delete(true)
	gstru.SetProp("overflow", "hidden")

}

func makeMainMenu() {
	appnm := gi.AppName()
	mmen := win.MainMenu
	mmen.ConfigMenus([]string{appnm, "File", "Edit"})

	fmen := win.MainMenu.ChildByName("File", 0).(*gi.Action)
	fmen.Menu = make(gi.Menu, 0, 10)
	fmen.Menu.AddAction(gi.ActOpts{Label: "New", ShortcutKey: gi.KeyFunMenuNew}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.Reset()
	})
	fmen.Menu.AddSeparator("sep0")
	fmen.Menu.AddAction(gi.ActOpts{Label: "Open", ShortcutKey: gi.KeyFunMenuOpen}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		giv.FileViewDialog(vp, "savedGraphs/", ".json", giv.DlgOpts{}, nil, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			if sig == int64(gi.DialogAccepted) {
				dlg := send.Embed(gi.KiT_Dialog).(*gi.Dialog)
				TheGraph.OpenJSON(gi.FileName(giv.FileViewDialogValue(dlg)))
			}
		})
	})
	fmen.Menu.AddAction(gi.ActOpts{Label: "Open Autosave", ShortcutKey: gi.KeyFunMenuOpenAlt1}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.OpenAutoSave()
	})
	fmen.Menu.AddSeparator("sep1")
	fmen.Menu.AddAction(gi.ActOpts{Label: "Save", ShortcutKey: gi.KeyFunMenuSave}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if TheGraph.State.File != "" {
			TheGraph.SaveLast()
		} else {
			giv.FileViewDialog(vp, "savedGraphs/", ".json", giv.DlgOpts{}, nil, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
				if sig == int64(gi.DialogAccepted) {
					dlg := send.Embed(gi.KiT_Dialog).(*gi.Dialog)
					TheGraph.SaveJSON(gi.FileName(giv.FileViewDialogValue(dlg)))
				}
			})
		}
	})
	fmen.Menu.AddAction(gi.ActOpts{Label: "Save as", ShortcutKey: gi.KeyFunMenuSaveAs}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		giv.FileViewDialog(vp, "savedGraphs/", ".json", giv.DlgOpts{}, nil, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			if sig == int64(gi.DialogAccepted) {
				dlg := send.Embed(gi.KiT_Dialog).(*gi.Dialog)
				TheGraph.SaveJSON(gi.FileName(giv.FileViewDialogValue(dlg)))
			}
		})
	})
	fmen.Menu.AddSeparator("sep2")
	fmen.Menu.AddAction(gi.ActOpts{Label: "Save as PNG", Shortcut: "Control+Alt+C"}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		img := TheGraph.Capture()
		giv.FileViewDialog(vp, "", ".png", giv.DlgOpts{}, nil, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			if sig == int64(gi.DialogAccepted) {
				dlg := send.Embed(gi.KiT_Dialog).(*gi.Dialog)
				SaveImageToFile(img, giv.FileViewDialogValue(dlg))
			}
		})
	})
	fmen.Menu.AddAction(gi.ActOpts{Label: "Copy PNG", Shortcut: "Shift+Control+C"}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.CopyGraphImage()
	})

	fmen.Menu.AddSeparator("sep3")
	fmen.Menu.AddAction(gi.ActOpts{Label: "Upload Graph", Shortcut: "Control+U"}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		gi.StringPromptDialog(vp, "", "", gi.DlgOpts{Title: "Upload Graph", Prompt: "Upload your graph for anyone else to see. Enter a name for your graph:"}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			if sig == int64(gi.DialogAccepted) {
				dlg := send.Embed(gi.KiT_Dialog).(*gi.Dialog)
				TheGraph.Upload(gi.StringPromptDialogValue(dlg))
			}
		})
	})
	fmen.Menu.AddAction(gi.ActOpts{Label: "Download Graph", Shortcut: "Control+D"}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.Download()
	})
	fmen.Menu.AddSeparator("sep4")
	fmen.Menu.AddAction(gi.ActOpts{Label: "Settings", ShortcutKey: gi.KeyFunMenuSaveAlt}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		pSettings := TheSettings
		giv.StructViewDialog(vp, &TheSettings, giv.DlgOpts{Title: "Settings", Ok: true, Cancel: true}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			if sig == int64(gi.DialogAccepted) {
				TheSettings.Save()
				svgGraph.SetProp("min-width", TheSettings.GraphSize)
				svgGraph.SetProp("min-height", TheSettings.GraphSize)
				var n float32 = 1.0 / float32(TheSettings.GraphInc)
				ginc = mat32.Vec2{X: n, Y: n}
				UpdateColors()
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

}
