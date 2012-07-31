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
	resized := image.Resize(Size{320, 200}, INTER_NN)

	if resized.Size.Width != 320 || resized.Size.Height != 200 {
		t.Error("Resized image is not 320x200 pixels")
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
