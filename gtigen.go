// Code generated by "goki generate"; DO NOT EDIT.

package main

import (
	"goki.dev/gti"
	"goki.dev/ordmap"
)

var _ = gti.AddType(&gti.Type{
	Name:      "main.Graph",
	ShortName: "main.Graph",
	IDName:    "graph",
	Doc:       "Graph contains the lines and parameters of a graph",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Params", &gti.Field{Name: "Params", Type: "github.com/kkoreilly/marbles.Params", LocalType: "Params", Doc: "the parameters for updating the marbles", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"Lines", &gti.Field{Name: "Lines", Type: "github.com/kkoreilly/marbles.Lines", LocalType: "Lines", Doc: "the lines of the graph -- can have any number", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"Marbles", &gti.Field{Name: "Marbles", Type: "[]*github.com/kkoreilly/marbles.Marble", LocalType: "[]*Marble", Doc: "", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\""}},
		{"State", &gti.Field{Name: "State", Type: "github.com/kkoreilly/marbles.State", LocalType: "State", Doc: "", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\""}},
		{"Functions", &gti.Field{Name: "Functions", Type: "github.com/kkoreilly/marbles.Functions", LocalType: "Functions", Doc: "", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\""}},
		{"Vectors", &gti.Field{Name: "Vectors", Type: "github.com/kkoreilly/marbles.Vectors", LocalType: "Vectors", Doc: "", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\""}},
		{"Objects", &gti.Field{Name: "Objects", Type: "github.com/kkoreilly/marbles.Objects", LocalType: "Objects", Doc: "", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\""}},
	}),
	Embeds: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{
		{"Graph", &gti.Method{Name: "Graph", Doc: "Graph updates graph for current equations, and resets marbles too", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"Run", &gti.Method{Name: "Run", Doc: "Run runs the marbles for NSteps", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"Stop", &gti.Method{Name: "Stop", Doc: "Stop stops the marbles", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"Step", &gti.Method{Name: "Step", Doc: "Step does one step update of marbles", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"StopSelecting", &gti.Method{Name: "StopSelecting", Doc: "StopSelecting stops selecting current marble", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"TrackSelectedMarble", &gti.Method{Name: "TrackSelectedMarble", Doc: "TrackSelectedMarble toggles track for the currently selected marble", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"SelectNextMarble", &gti.Method{Name: "SelectNextMarble", Doc: "SelectNextMarble selects the next marble in the viewbox", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
	}),
})

var _ = gti.AddType(&gti.Type{
	Name:      "main.Params",
	ShortName: "main.Params",
	IDName:    "params",
	Doc:       "Params are the parameters of the graph",
	Directives: gti.Directives{
		&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"NMarbles", &gti.Field{Name: "NMarbles", Type: "int", LocalType: "int", Doc: "Number of marbles", Directives: gti.Directives{}, Tag: "min:\"1\" max:\"10000\" step:\"10\" label:\"Number of marbles\""}},
		{"MarbleStartX", &gti.Field{Name: "MarbleStartX", Type: "github.com/kkoreilly/marbles.Expr", LocalType: "Expr", Doc: "Marble horizontal start position", Directives: gti.Directives{}, Tag: ""}},
		{"MarbleStartY", &gti.Field{Name: "MarbleStartY", Type: "github.com/kkoreilly/marbles.Expr", LocalType: "Expr", Doc: "Marble vertical start position", Directives: gti.Directives{}, Tag: ""}},
		{"StartVelY", &gti.Field{Name: "StartVelY", Type: "github.com/kkoreilly/marbles.Param", LocalType: "Param", Doc: "Starting horizontal velocity of the marbles", Directives: gti.Directives{}, Tag: "label:\"Starting velocity y\""}},
		{"StartVelX", &gti.Field{Name: "StartVelX", Type: "github.com/kkoreilly/marbles.Param", LocalType: "Param", Doc: "Starting vertical velocity of the marbles", Directives: gti.Directives{}, Tag: "label:\"Starting velocity x\""}},
		{"UpdateRate", &gti.Field{Name: "UpdateRate", Type: "github.com/kkoreilly/marbles.Param", LocalType: "Param", Doc: "how fast to move along velocity vector -- lower = smoother, more slow-mo", Directives: gti.Directives{}, Tag: ""}},
		{"TimeStep", &gti.Field{Name: "TimeStep", Type: "github.com/kkoreilly/marbles.Param", LocalType: "Param", Doc: "how fast time increases", Directives: gti.Directives{}, Tag: ""}},
		{"YForce", &gti.Field{Name: "YForce", Type: "github.com/kkoreilly/marbles.Param", LocalType: "Param", Doc: "how fast it accelerates down", Directives: gti.Directives{}, Tag: "label:\"Y force (Gravity)\""}},
		{"XForce", &gti.Field{Name: "XForce", Type: "github.com/kkoreilly/marbles.Param", LocalType: "Param", Doc: "how fast the marbles move side to side without collisions, set to 0 for no movement", Directives: gti.Directives{}, Tag: "label:\"X force (Wind)\""}},
		{"CenterX", &gti.Field{Name: "CenterX", Type: "github.com/kkoreilly/marbles.Param", LocalType: "Param", Doc: "the center point of the graph, x", Directives: gti.Directives{}, Tag: "label:\"Graph center x\""}},
		{"CenterY", &gti.Field{Name: "CenterY", Type: "github.com/kkoreilly/marbles.Param", LocalType: "Param", Doc: "the center point of the graph, y", Directives: gti.Directives{}, Tag: "label:\"Graph center y\""}},
		{"TrackingSettings", &gti.Field{Name: "TrackingSettings", Type: "github.com/kkoreilly/marbles.TrackingSettings", LocalType: "TrackingSettings", Doc: "", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})
