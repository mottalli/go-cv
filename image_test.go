package cv

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestColor(t *testing.T) {
	image, _ := LoadImage(imagePath, true)
	defer image.Release()

	r, g, b, a := image.At(30, 20).RGBA()

	var er, eg, eb, ea uint32
	er, eg, eb, ea = 224*0x101, 134*0x101, 110*0x101, 255*0x101

	if r != er {
		t.Error(fmt.Sprintf("Expected red value of %v, got %v", er, r))
	}

	if g != eg {
		t.Error(fmt.Sprintf("Expected green value of %v, got %v", eg, g))
	}

	if b != eb {
		t.Error(fmt.Sprintf("Expected blue value of %v, got %v", eb, b))
	}

	if a != ea {
		t.Error(fmt.Sprintf("Expected alpha value of %v, got %v", ea, a))
	}
}

func TestGray(t *testing.T) {
	image, _ := LoadImage(imagePath, false)
	defer image.Release()

	r, g, b, a := image.At(100, 100).RGBA()

	if r != g || g != b {
		t.Error("Obtained non-grayscale image")
	}

	if a != 255*0x101 {
		t.Error("Grayscale alpha channel invalid")
	}
}

func TestResize(t *testing.T) {
	image, _ := LoadImage(imagePath, true)
	defer image.Release()

	newsize := Size{320, 200}
	resized := image.Resize(newsize, INTER_NN)
	defer resized.Release()

	if resized.Size() != newsize {
		t.Error("Resized image is not 320x200 pixels")
	}
}

func TestInitialize(t *testing.T) {
	img := NewImage()
	defer img.Release()

	if img.Initialized {
		t.Error("Empty image is initialized")
	}

	size := Size{320, 200}

	img.Initialize(size, CV_8UC3)

	if !img.Initialized {
		t.Error("Initialized image not marked as initialized, or image was not initialized successfully")
	}

	if img.Size() != size {
		t.Error("Wrong image size while intialising")
	}

	if img.Type() != CV_8UC3 {
		t.Error("Wrong image type while initialising")
	}
}

func TestReinitialize(t *testing.T) {
	img := NewImage()
	defer img.Release()

	img.Initialize(Size{10, 20}, CV_8UC1)
	// Initialize again with different parameters
	img.Initialize(Size{30, 40}, CV_8UC3)

	if img.Size().Width != 30 || img.Size().Height != 40 || img.Type() != CV_8UC3 {
		t.Error("Invalid format when re-initialising image")
	}
}

func TestSplit(t *testing.T) {
	image, _ := LoadImage(imagePath, true)
	defer image.Release()

	channels := image.Split()
	defer channels.Release()

	if len(channels) != 3 {
		t.Error("Could not split image into 3 channels")
	}

	err := image.SplitTo(&channels)
	if err != nil {
		t.Error("Error while splitting image into existing channel list")
	}
}

func BenchmarkRandomImageAccess(b *testing.B) {
	b.StopTimer()
	image, _ := LoadImage(imagePath, true)

	b.StartTimer()
	bounds := image.Bounds()

	for c := 0; c < b.N; c++ {
		image.At(rand.Intn(bounds.Max.X), rand.Intn(bounds.Max.Y))
	}
}

func BenchmarkSequentialImageAccess(b *testing.B) {
	b.StopTimer()
	image, _ := LoadImage(imagePath, true)

	b.StartTimer()
	bounds := image.Bounds()

	for c := 0; c < b.N; {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				image.At(x, y)
				c++
				if c == b.N {
					break
				}
			}
		}
	}
}
