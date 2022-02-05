package main

import (
	"encoding/json"
	"os"

	"github.com/goki/gi/gist"
)

type Settings struct {
	LineDefaults   LineDefaults   `view:"inline" label:"Line Defaults"`
	GraphDefaults  Params         `view:"inline" label:"Graph Parameter Defaults"`
	MarbleSettings MarbleSettings `view:"inline" label:"Marble Settings"`
	ColorSettings  ColorSettings  `view:"no-inline" label:"Color Settings"`
}
type ColorSettings struct {
	BackgroundColor    gist.Color `label:"Background Color"`
	GraphColor         gist.Color `label:"Graph Background Color"`
	AxisColor          gist.Color `label:"Graph Axis Color"`
	StatusBarColor     gist.Color `label:"Status Bar Background Color"`
	ButtonColor        gist.Color `label:"Button Color"`
	StatusTextColor    gist.Color `label:"Status Bar Text Color"`
	GraphTextColor     gist.Color `label:"Graph Controls Text Color"`
	LineTextColor      gist.Color `label:"Line Text Color"`
	ToolBarColor       gist.Color `label:"Toolbar Background Color"`
	ToolBarButtonColor gist.Color `label:"Toolbar Button Color"`
	GraphParamsColor   gist.Color `label:"Graph Parameters Background Color"`
}
type MarbleSettings struct {
	MarbleColor string
	MarbleSize  float64
}
type LineDefaults struct {
	Expr        string
	MinX        string
	MaxX        string
	MinY        string
	MaxY        string
	Bounce      string
	Color       string
	ColorSwitch string
}

var TheSettings Settings

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
	ln.Color = "default"
	ln.Bounce = "0.95"
	ln.MinX = "-10"
	ln.MaxX = "10"
	ln.MinY = "-10"
	ln.MaxY = "10"
	ln.ColorSwitch = "none"
}

func (ms *MarbleSettings) Defaults() {
	ms.MarbleColor = "default"
	ms.MarbleSize = 0.1
}

func (cs *ColorSettings) Defaults() {
	white, _ := gist.ColorFromName("white")
	grey, _ := gist.ColorFromName("grey")
	black, _ := gist.ColorFromName("black")
	lightblue, _ := gist.ColorFromName("lightblue")
	cs.BackgroundColor = white
	cs.GraphColor = white
	cs.AxisColor = grey
	cs.StatusBarColor = lightblue
	cs.ButtonColor = white
	cs.StatusTextColor = black
	cs.GraphTextColor = black
	cs.LineTextColor = black
	cs.ToolBarColor = white
	cs.ToolBarButtonColor = white
	cs.GraphParamsColor = white
}
