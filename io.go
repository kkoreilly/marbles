package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"cogentcore.org/core/core"
)

// SaveLast saves to the last opened or saved file
func (gr *Graph) SaveLast() { //types:add
	if gr.State.File != "" {
		TheGraph.SaveJSON(gr.State.File)
	}
}

// OpenJSON opens a graph from a JSON file
func (gr *Graph) OpenJSON(filename core.Filename) error { //types:add
	b, err := os.ReadFile(string(filename))
	if HandleError(err) {
		return err
	}
	err = json.Unmarshal(b, gr)
	if HandleError(err) {
		return err
	}
	gr.State.File = filename
	// gr.AutoGraphAndUpdate()
	return err
}

// OpenAutoSave opens the last graphed graph, stays between sessions of the app
func (gr *Graph) OpenAutoSave() error {
	b, err := os.ReadFile(filepath.Join(core.TheApp.AppDataDir(), "autosave.json"))
	if HandleError(err) {
		return err
	}
	err = json.Unmarshal(b, gr)
	if HandleError(err) {
		return err
	}
	// gr.AutoGraphAndUpdate()
	return err
}

// SaveJSON saves a graph to a JSON file
func (gr *Graph) SaveJSON(filename core.Filename) error { //types:add
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
	err = os.WriteFile(string(filename), b, os.ModePerm)
	if HandleError(err) {
		return err
	}
	gr.State.File = filename
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
	err = os.WriteFile(filepath.Join(core.TheApp.AppDataDir(), "autosave.json"), b, 0666)
	HandleError(err)
	return err
}
