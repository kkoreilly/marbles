package main

import (
	"encoding/json"
	"os"

	"github.com/goki/gi/gist"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
)

type Settings struct {
	LineDefaults   LineDefaults   `view:"inline" label:"Line Defaults"`
	GraphDefaults  Params         `view:"inline" label:"Graph Parameter Defaults"`
	MarbleSettings MarbleSettings `view:"inline" label:"Marble Settings"`
	ColorSettings  ColorSettings  `view:"no-inline" label:"Color Settings"`
}
type ColorSettings struct {
	BackgroundColor      gist.Color `label:"Background Color"`
	GraphColor           gist.Color `label:"Graph Background Color"`
	AxisColor            gist.Color `label:"Graph Axis Color"`
	StatusBarColor       gist.Color `label:"Status Bar Background Color"`
	ButtonColor          gist.Color `label:"Button Color"`
	StatusTextColor      gist.Color `label:"Status Bar Text Color"`
	GraphTextColor       gist.Color `label:"Graph Controls Text Color"`
	LineTextColor        gist.Color `label:"Line Text Color"`
	ToolBarColor         gist.Color `label:"Toolbar Background Color"`
	ToolBarButtonColor   gist.Color `label:"Toolbar Button Color"`
	GraphParamsColor     gist.Color `label:"Graph Parameters Background Color"`
	LinesBackgroundColor gist.Color `label:"Lines Background Color"`
}
type MarbleSettings struct {
	MarbleColor string
	MarbleSize  float64
}
type LineDefaults struct {
	Expr       string
	MinX       string
	MaxX       string
	MinY       string
	MaxY       string
	Bounce     string
	LineColors LineColors
}

var TheSettings Settings

var SettingProps = ki.Props{
	"ToolBar": ki.PropSlice{
		{Name: "Reset", Value: ki.Props{
			"label": "Reset Settings",
		},
		},
	}}

var KiT_TheSettings = kit.Types.AddType(&TheSettings, SettingProps)

func (se *Settings) Reset() {
	se.Defaults()
	se.Save()
}

func (se *Settings) Get() {
	b, err := os.ReadFile("localData/settings.json")
	if err != nil {
		se.Defaults()
		se.Save()
		return
	}
	err = json.Unmarshal(b, se)
	if err != nil {
		se.Defaults()
		se.Save()
		return
	}
	if se.LineDefaults.Expr == "" {
		se.LineDefaults.BasicDefaults()
		se.Save()
	}
	if se.GraphDefaults.MinSize.X == 0 {
		se.GraphDefaults.BasicDefaults()
		se.Save()
	}
	if se.MarbleSettings.MarbleColor == "" {
		se.MarbleSettings.Defaults()
		se.Save()
	}
	if se.ColorSettings.BackgroundColor == gist.NilColor {
		se.ColorSettings.Defaults()
		se.Save()
	}

}

func (se *Settings) Save() {
	b, err := json.MarshalIndent(se, "", "  ")
	if HandleError(err) {
		return
	}
	err = os.WriteFile("localData/settings.json", b, 0644)
	HandleError(err)
}
func (se *Settings) Defaults() {
	se.LineDefaults.BasicDefaults()
	se.GraphDefaults.BasicDefaults()
	se.MarbleSettings.Defaults()
	se.ColorSettings.Defaults()
}

func (ln *LineDefaults) BasicDefaults() {
	ln.Expr = "x"
	ln.LineColors.Color = gist.White
	ln.Bounce = "0.95"
	ln.MinX = "-10"
	ln.MaxX = "10"
	ln.MinY = "-10"
	ln.MaxY = "10"
	ln.LineColors.ColorSwitch = gist.White
}

func (ms *MarbleSettings) Defaults() {
	ms.MarbleColor = "default"
	ms.MarbleSize = 0.1
}

func (cs *ColorSettings) Defaults() {
	grey, _ := gist.ColorFromName("grey")
	lightblue, _ := gist.ColorFromName("lightblue")
	cs.BackgroundColor = gist.White
	cs.GraphColor = gist.White
	cs.AxisColor = grey
	cs.StatusBarColor = lightblue
	cs.ButtonColor = gist.White
	cs.StatusTextColor = gist.Black
	cs.GraphTextColor = gist.Black
	cs.LineTextColor = gist.Black
	cs.ToolBarColor = gist.White
	cs.ToolBarButtonColor = gist.White
	cs.GraphParamsColor = gist.White
	cs.LinesBackgroundColor = gist.White
}
