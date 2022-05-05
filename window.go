package main

import (
	"path/filepath"
	"strconv"

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
	lns.SetProp("inact-key-nav", false)
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
			// lbl.SetProp("background-color", "yellow")
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
				if col == 0 && TheGraph.State.Error == nil {
					val := TheGraph.Lines[row].Expr.Eval(0, 0, 0)
					funcName := "y"
					if row < len(functionNames) {
						funcName = functionNames[row]
					}
					valueText.SetText(funcName + "(0) ≈ " + strconv.FormatFloat(val, 'f', 3, 64))
				}
			})
			widg.SetProp("font-size", TheSettings.LineFontSize)
			edit.SetCompleter(edit, ExprComplete, ExprCompleteEdit)
		}
	}
	lns.ViewSig.Connect(lns.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.AutoGraph()
	})

	params = giv.AddNewStructView(sidesplit, "params")
	params.SetStruct(&TheGraph.Params)
	params.ViewSig.Connect(params.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.AutoGraph()
	})

	sidesplit.SetSplits(6, 4)

	frame := gi.AddNewFrame(split, "frame", gi.LayoutHoriz)

	TheGraph.Objects.Graph = svg.AddNewSVG(frame, "graph")
	TheGraph.Objects.Graph.SetProp("min-width", TheSettings.GraphSize)
	TheGraph.Objects.Graph.SetProp("min-height", TheSettings.GraphSize)
	TheGraph.Objects.Lines = svg.AddNewGroup(TheGraph.Objects.Graph, "TheGraph.Objects.Lines")
	TheGraph.Objects.Marbles = svg.AddNewGroup(TheGraph.Objects.Graph, "TheGraph.Objects.Marbles")
	TheGraph.Objects.Coords = svg.AddNewGroup(TheGraph.Objects.Graph, "TheGraph.Objects.Coords")
	TheGraph.Objects.TrackingLines = svg.AddNewGroup(TheGraph.Objects.Graph, "TheGraph.Objects.TrackingLines")
	split.SetSplits(float32(width-TheSettings.GraphSize), float32(TheSettings.GraphSize)*7/8)
	TheGraph.Vectors.Min = mat32.Vec2{X: -graphViewBoxSize, Y: -graphViewBoxSize}
	TheGraph.Vectors.Max = mat32.Vec2{X: graphViewBoxSize, Y: graphViewBoxSize}
	TheGraph.Vectors.Size = TheGraph.Vectors.Max.Sub(TheGraph.Vectors.Min)
	var n float32 = 1.0 / float32(TheSettings.GraphInc)
	TheGraph.Vectors.Inc = mat32.Vec2{X: n, Y: n}

	TheGraph.Objects.Graph.ViewBox.Min = TheGraph.Vectors.Min
	TheGraph.Objects.Graph.ViewBox.Size = TheGraph.Vectors.Size
	TheGraph.Objects.Graph.Norm = true
	TheGraph.Objects.Graph.InvertY = true
	TheGraph.Objects.Graph.Fill = true
	TheGraph.Objects.Graph.SetProp("background-color", "white")
	TheGraph.Objects.Graph.SetProp("stroke-width", ".2pct")

	statusBar = gi.AddNewFrame(mfr, "statusBar", gi.LayoutHoriz)
	statusBar.SetStretchMaxWidth()
	fpsText = gi.AddNewLabel(statusBar, "fpsText", "FPS: ")
	fpsText.SetProp("font-weight", "bold")
	fpsText.SetStretchMaxWidth()
	fpsText.Redrawable = true
	valueText = gi.AddNewLabel(statusBar, "valueText", "f(0) ≈ ")
	valueText.SetProp("font-weight", "bold")
	valueText.SetStretchMaxWidth()
	valueText.Redrawable = true
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
	lns.ToolBar().Delete(true)
	params.ToolBar().Delete(true)
	gstru.ToolBar().Child(0).Delete(true)
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
		giv.FileViewDialog(vp, filepath.Join(GetMarblesFolder(), "savedGraphs")+"/", ".json", giv.DlgOpts{}, nil, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
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
			giv.FileViewDialog(vp, filepath.Join(GetMarblesFolder(), "savedGraphs")+"/", ".json", giv.DlgOpts{}, nil, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
				if sig == int64(gi.DialogAccepted) {
					dlg := send.Embed(gi.KiT_Dialog).(*gi.Dialog)
					TheGraph.SaveJSON(gi.FileName(giv.FileViewDialogValue(dlg)))
				}
			})
		}
	})
	fmen.Menu.AddAction(gi.ActOpts{Label: "Save as", ShortcutKey: gi.KeyFunMenuSaveAs}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		giv.FileViewDialog(vp, filepath.Join(GetMarblesFolder(), "savedGraphs")+"/", ".json", giv.DlgOpts{}, nil, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
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
				TheGraph.Objects.Graph.SetProp("min-width", TheSettings.GraphSize)
				TheGraph.Objects.Graph.SetProp("min-height", TheSettings.GraphSize)
				var n float32 = 1.0 / float32(TheSettings.GraphInc)
				TheGraph.Vectors.Inc = mat32.Vec2{X: n, Y: n}
				UpdateColors()
				TheGraph.AutoGraphAndUpdate()
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
