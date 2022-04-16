package main

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"os"

	"github.com/kbinani/screenshot"
	"golang.design/x/clipboard"
)

// CopyGraphImage captures a picture of the graph and copies it
func (gr *Graph) CopyGraphImage() {
	img := gr.Capture()
	CopyImage(img)
}

// Capture captures an image of the graph and returns it
func (gr *Graph) Capture() *image.RGBA {
	img, err := screenshot.CaptureRect(gr.Objects.Graph.BBox2D())
	HandleError(err)
	return img
}

// CopyImage copies an image to the clipboard
func CopyImage(img *image.RGBA) {
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	b, err := io.ReadAll(buf)
	if HandleError(err) {
		return
	}
	clipboard.Write(clipboard.FmtImage, b)
}

// SaveImageToFile saves the image to a given filename
func SaveImageToFile(img *image.RGBA, filename string) {
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	b, err := io.ReadAll(buf)
	if HandleError(err) {
		return
	}
	os.WriteFile(filename, b, os.ModePerm)
}

// InitClipboard inits the clipboard
func InitClipboard() {
	clipboard.Init()
}
