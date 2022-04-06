package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/goki/gi/gi"
)

// SaveLast saves to the last opened or saved file
func (gr *Graph) SaveLast() {
	if gr.State.File == "" {
		errorText.SetText("no file has been opened or saved")
	} else {
		TheGraph.SaveJSON(gi.FileName(gr.State.File))
	}
}

// OpenJSON opens a graph from a JSON file
func (gr *Graph) OpenJSON(filename gi.FileName) error {
	b, err := os.ReadFile(string(filename))
	if HandleError(err) {
		return err
	}
	err = json.Unmarshal(b, gr)
	if HandleError(err) {
		return err
	}
	gr.State.File = string(filename)
	UpdateCurrentFileText()
	return err
}

// OpenAutoSave opens the last graphed graph, stays between sessions of the app
func (gr *Graph) OpenAutoSave() error {
	b, err := os.ReadFile("localData/autosave.json")
	if HandleError(err) {
		return err
	}
	err = json.Unmarshal(b, gr)
	if HandleError(err) {
		return err
	}
	return err
}

// SaveJSON saves a graph to a JSON file
func (gr *Graph) SaveJSON(filename gi.FileName) error {
	var b []byte
	var err error
	if TheSettings.PrettyJSON {
		b, err = json.MarshalIndent(gr, "", "  ")
	} else {
		b, err = json.Marshal(gr)
	}
	if HandleError(err) {
		return err
	}
	err = os.WriteFile(string(filename), b, 0644)
	HandleError(err)
	gr.State.File = string(filename)
	UpdateCurrentFileText()
	return err
}

// AutoSave saves the graph to autosave.json, called automatically
func (gr *Graph) AutoSave() error {
	var b []byte
	var err error
	if TheSettings.PrettyJSON {
		b, err = json.MarshalIndent(gr, "", "  ")
	} else {
		b, err = json.Marshal(gr)
	}
	if HandleError(err) {
		return err
	}
	err = os.WriteFile("localData/autosave.json", b, 0644)
	HandleError(err)
	return err
}

// Upload uploads the graph to the global database
func (gr *Graph) Upload(name string) {
	b, err := json.Marshal(gr)
	if HandleError(err) {
		return
	}
	UploadGraph(name, string(b))
}

// Download downloads a graph from the database
func (gr *Graph) Download() {
	cgwin := gi.NewMainWindow("cgwin", "Choose a graph", width, height)
	cgvp := cgwin.WinViewport2D()
	updt := cgvp.UpdateStart()
	cgmfr := cgwin.SetMainFrame()
	titleText := gi.AddNewLabel(cgmfr, "titleText", "Pick a graph to download")
	titleText.SetProp("font-size", "x-large")
	gi.AddNewSeparator(cgmfr, "TitleSeparator", true)
	graphs := GetGraphs()
	for k, d := range graphs {
		year, month, day := d.Date.Date()
		nameText := gi.AddNewLabel(cgmfr, fmt.Sprintf("Graph%vNameText", k), "<b>Graph Name:</b> "+d.Name)
		nameText.SetProp("font-size", "large")
		dateText := gi.AddNewLabel(cgmfr, fmt.Sprintf("Graph%vDateText", k), fmt.Sprintf("Published On %v %v, %v", month.String(), day, year))
		dateText.SetProp("font-size", "large")
		chooseButton := gi.AddNewButton(cgmfr, fmt.Sprintf("Graph%vButton", k))
		chooseButton.SetText("Open this graph")
		graphData := d.Graph
		chooseButton.OnClicked(func() {
			cgwin.Close()
			gr.OpenGraphFromString(graphData)
		})
		gi.AddNewSeparator(cgmfr, fmt.Sprintf("Graph%vSeparator", k), true)
	}
	cgvp.UpdateEndNoSig(updt)
	cgwin.StartEventLoop()
}

// OpenGraphFromString opens a graph given its json string
func (gr *Graph) OpenGraphFromString(data string) {
	err := json.Unmarshal([]byte(data), gr)
	if HandleError(err) {
		return
	}
}
