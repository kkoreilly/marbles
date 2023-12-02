package main

import (
	"goki.dev/colors"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
	"goki.dev/goosi/events"
	"goki.dev/mat32/v2"
	"goki.dev/svg"
)

func (gr *Graph) MakeBasicElements(b *gi.Body) {
	sp := gi.NewSplits(b)

	lns := giv.NewTableView(sp).SetSlice(&gr.Lines)
	/*
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
	*/
	lns.OnChange(func(e events.Event) {
		gr.AutoGraph()
	})

	gr.Objects.Graph = gi.NewSVG(sp)
	gr.Objects.SVG = svg.NewSVG(500, 500)
	gr.Objects.Graph.SVG = gr.Objects.SVG

	gr.Objects.SVG.Norm = true
	gr.Objects.SVG.InvertY = true
	gr.Objects.SVG.Fill = true
	gr.Objects.SVG.BackgroundColor.SetSolid(colors.Scheme.Surface)

	gr.Vectors.Min = mat32.Vec2{X: -graphViewBoxSize, Y: -graphViewBoxSize}
	gr.Vectors.Max = mat32.Vec2{X: graphViewBoxSize, Y: graphViewBoxSize}
	gr.Vectors.Size = gr.Vectors.Max.Sub(gr.Vectors.Min)
	var n float32 = 1.0 / float32(TheSettings.GraphInc)
	gr.Vectors.Inc = mat32.Vec2{X: n, Y: n}

	gr.Objects.Root = &gr.Objects.SVG.Root
	gr.Objects.Root.ViewBox.Min = gr.Vectors.Min
	gr.Objects.Root.ViewBox.Size = gr.Vectors.Size

	svg.NewCircle(gr.Objects.Root).SetRadius(50)

	sp.SetSplits(0.3, 0.7)

	/*
		params = giv.AddNewStructView(sidesplit, "params")
		params.SetStruct(&TheGraph.Params)
		params.ViewSig.Connect(params.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			TheGraph.AutoGraph()
		})

		sidesplit.SetSplits(6, 4)
	*/

	// graphFrame := gi.AddNewFrame(sp, "graphFrame", gi.LayoutVert)

	// TheGraph.Objects.Graph = svg.AddNewSVG(graphFrame, "graph")
	// TheGraph.Objects.Graph.SetFixedHeight(units.NewDot(float32(TheSettings.GraphSize) - 20))
	// TheGraph.Objects.Graph.SetFixedWidth(units.NewDot(float32(TheSettings.GraphSize) - 20))
	// TheGraph.Objects.Graph.SetStretchMaxWidth()

	gr.Objects.Lines = svg.NewGroup(gr.Objects.Root, "lines")
	gr.Objects.Marbles = svg.NewGroup(gr.Objects.Root, "marbles")
	gr.Objects.Coords = svg.NewGroup(gr.Objects.Root, "coords")
	gr.Objects.TrackingLines = svg.NewGroup(gr.Objects.Root, "tracking-lines")

	// TheGraph.Objects.Graph.ViewBox.Min = TheGraph.Vectors.Min
	// TheGraph.Objects.Graph.ViewBox.Size = TheGraph.Vectors.Size
	// TheGraph.Objects.Graph.Norm = true
	// TheGraph.Objects.Graph.InvertY = true
	// TheGraph.Objects.Graph.Fill = true
	// TheGraph.Objects.Graph.SetProp("background-color", "white")
	// TheGraph.Objects.Graph.SetProp("stroke-width", ".2pct")

	// gp := float32(TheSettings.GraphSize) / float32(width)
	// sp.SetSplits(1-gp, gp)

	// gi.AddNewSeparator(graphFrame, "sep", true)

	// statusBar = gi.AddNewFrame(mfr, "statusBar", gi.LayoutHoriz)
	// statusBar.SetStretchMaxWidth()
	// fpsText = gi.AddNewLabel(statusBar, "fpsText", "FPS: ")
	// fpsText.SetProp("font-weight", "bold")
	// fpsText.SetStretchMaxWidth()
	// fpsText.Redrawable = true
	// valueText = gi.AddNewLabel(statusBar, "valueText", "f(0) ≈ ")
	// valueText.SetProp("font-weight", "bold")
	// valueText.SetStretchMaxWidth()
	// valueText.Redrawable = true
	// errorText = gi.AddNewLabel(statusBar, "errorText", "")
	// errorText.SetProp("font-weight", "bold")
	// errorText.SetStretchMaxWidth()
	// errorText.Redrawable = true
	// errorText.SetText("Graphed successfully")
	// currentFileText = gi.AddNewLabel(statusBar, "currentFileText", "untitled.json")
	// currentFileText.SetProp("font-weight", "bold")
	// currentFileText.SetStretchMaxWidth()
	// currentFileText.Redrawable = true
	// versionText = gi.AddNewLabel(statusBar, "versionText", "")
	// versionText.SetProp("font-weight", "bold")
	// versionText.SetStretchMaxWidth()
	// versionText.SetText("Running version " + GetVersion())
	// lns.ToolBar().Delete(true)
	// params.ToolBar().Delete(true)
}

/*
func makeToolbar() {
	graphToolbar = gi.AddNewToolBar(mfr, "graphToolbar")
	graphToolbar.AddAction(gi.ActOpts{Name: "Graph", Label: "Graph", Icon: "file-image", Tooltip: "graph the equations and reset the marbles"}, graphToolbar.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.Graph()
		vp.SetNeedsFullRender()
	})

	graphToolbar.AddAction(gi.ActOpts{Name: "Run", Label: "Run", Icon: "run", Tooltip: "runs the marbles"}, graphToolbar.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.Run()
	})

	graphToolbar.AddAction(gi.ActOpts{Name: "Stop", Label: "Stop", Icon: "stop", Tooltip: "stop the marbles"}, graphToolbar.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.Stop()
	})

	graphToolbar.AddAction(gi.ActOpts{Name: "Step", Label: "Step", Icon: "step-fwd", Tooltip: "steps the marbles for one step"}, graphToolbar.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.Step()
	})

	graphToolbar.AddSeparator("sep1")

	graphToolbar.AddAction(gi.ActOpts{Name: "NextMarble", Label: "Next Marble", Icon: "forward", ShortcutKey: gi.KeyFunFocusNext, Tooltip: "selects the next marble"}, graphToolbar.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.SelectNextMarble()
	})
	graphToolbar.AddAction(gi.ActOpts{Name: "Unselect", Label: "Unselect", Icon: "stop", Tooltip: "stops selecting the marble"}, graphToolbar.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.StopSelecting()
	})
	graphToolbar.AddAction(gi.ActOpts{Name: "Track", Label: "Track", Icon: "edit", ShortcutKey: gi.KeyFunTranspose, Tooltip: "toggles track for the currently selected marble"}, graphToolbar.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.TrackSelectedMarble()
	})

	graphToolbar.AddSeparator("sep2")

	graphToolbar.AddAction(gi.ActOpts{Name: "NewLine", Label: "New Line", Icon: "plus", Shortcut: "Command+M", Tooltip: "adds a new blank line"}, graphToolbar.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.AddLine()
		vp.SetNeedsFullRender()
	})

	for _, d := range *graphToolbar.Children() {
		d.SetProp("border-radius", units.NewPx(7))
	}
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
	fmen.Menu.AddAction(gi.ActOpts{Label: "Save as PNG", Shortcut: "Command+Alt+C"}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		img := TheGraph.Capture()
		giv.FileViewDialog(vp, "", ".png", giv.DlgOpts{}, nil, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			if sig == int64(gi.DialogAccepted) {
				dlg := send.Embed(gi.KiT_Dialog).(*gi.Dialog)
				SaveImageToFile(img, giv.FileViewDialogValue(dlg))
			}
		})
	})
	fmen.Menu.AddAction(gi.ActOpts{Label: "Copy PNG", Shortcut: "Shift+Command+C"}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		TheGraph.CopyGraphImage()
	})

	fmen.Menu.AddSeparator("sep3")
	fmen.Menu.AddAction(gi.ActOpts{Label: "Upload Graph", Shortcut: "Command+U"}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		gi.StringPromptDialog(vp, "", "", gi.DlgOpts{Title: "Upload Graph", Prompt: "Upload your graph for anyone else to see. Enter a name for your graph:"}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
			if sig == int64(gi.DialogAccepted) {
				dlg := send.Embed(gi.KiT_Dialog).(*gi.Dialog)
				TheGraph.Upload(gi.StringPromptDialogValue(dlg))
			}
		})
	})
	fmen.Menu.AddAction(gi.ActOpts{Label: "Download Graph", Shortcut: "Command+D"}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
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
				gp := float32(TheSettings.GraphSize) / float32(width)
				mainSplit.SetSplits(1-gp, gp)
				UpdateColors()
				TheGraph.Objects.Graph.SetFixedHeight(units.NewDot(float32(TheSettings.GraphSize) - 20))
				TheGraph.Objects.Graph.SetFixedWidth(units.NewDot(float32(TheSettings.GraphSize) - 20))
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
*/
