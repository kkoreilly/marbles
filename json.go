package main

import (
	"encoding/json"
	"os"

	"github.com/goki/gi/gi"
)

// Save last saves to the last opened or saved file
func (gr *Graph) SaveLast() {
	if LastSavedFile == "" {
		errorText.SetText("no file has been opened or saved")
	} else {
		Gr.SaveJSON(gi.FileName(LastSavedFile))
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
	LastSavedFile = string(filename)
	gr.Graph()
	return err
}

// Open Autosave opens the last graphed graph, stays between sessions of the app
func (gr *Graph) OpenAutoSave() error {
	b, err := os.ReadFile("localData/autosave.json")
	if HandleError(err) {
		return err
	}
	err = json.Unmarshal(b, gr)
	if HandleError(err) {
		return err
	}
	gr.Graph()
	return err
}

// SaveJSON saves a graph to a JSON file
func (gr *Graph) SaveJSON(filename gi.FileName) error {
	b, err := json.MarshalIndent(gr, "", "  ")
	if HandleError(err) {
		return err
	}
	err = os.WriteFile(string(filename), b, 0644)
	HandleError(err)
	LastSavedFile = string(filename)
	return err
}

// Autosave saves the graph to autosave.json, called automatically
func (gr *Graph) AutoSave() error {
	b, err := json.MarshalIndent(gr, "", "  ")
	if HandleError(err) {
		return err
	}
	err = os.WriteFile("localData/autosave.json", b, 0644)
	HandleError(err)
	return err
}
