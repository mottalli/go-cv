package cv

import (
	"testing"
	"image/color"
)

const (
	imagePath string = "lena.jpg"
)

func TestLoadImageColorOK(t *testing.T) {
	image, err := LoadImage(imagePath, true)
	if err != nil {
		t.Error("Could not open file " + imagePath)
	}
	
	if image.Size.Width != 512 || image.Size.Height != 512 {
		t.Error("Wrong image size")
	}
	
	if image.ColorModel() != color.RGBAModel {
		t.Error("Color model is not RGBA")
	}
}

func TestLoadImageBWOK(t *testing.T) {
	image, err := LoadImage(imagePath, false)
	if err != nil {
		t.Error("Could not open file " + imagePath)
	}
	
	if image.ColorModel() != color.GrayModel {
		t.Error("Color model is not grayscale")
	}
}

func TestLoadImageBad(t *testing.T) {
	image, err := LoadImage("**inexistent file name**", true)
	if err == nil {
		t.Error("Did not return error opening invalid file")
	}
	
	if image != nil {
		t.Error("Did not return nil image when opening invalid file")
	}
}
