package cv

/*
#cgo pkg-config: opencv
#include <opencv/cv.h>
#include <opencv/highgui.h>
#include <opencv2/imgproc/types_c.h>
*/
import "C"
import (
	"errors"
	//"fmt"
	goimage "image"
	"image/color"
	"unsafe"
)

/******* Image initialization methods *******/

func (img *Image) InitializeAs(other *Image) *Image {
	if img.Initialized && img.Size() == other.Size() && img.Type() == other.Type() {
		// Do nothing if they are of the same size and type
		return img
	}

	img.Release()
	tmp := CreateImage(other.Size(), other.imtype)
	*img = *tmp
	return img
}

func depthFromType(matType MatType, imageType bool) C.int {
	if matType.Depth == 8 && matType.ElemType == Unsigned {
		if imageType {
			return C.IPL_DEPTH_8U
		}
	}
	panic("TODO: depthFromType - implement this!")
	return C.int(0)
}

func typeFromDepthAndChannels(depth, channels C.int) MatType {
	if depth == C.IPL_DEPTH_8U && channels == 1 {
		return CV_8UC1
	} else if depth == C.IPL_DEPTH_8U && channels == 3 {
		return CV_8UC3
	}

	panic("TODO: typeFromDepthAndChannels - implement this!")
	return CV_8UC1
}

func ImageFromIplImage(iplImage *C.IplImage) (*Image, error) {
	image := new(Image)
	image.ptr = unsafe.Pointer(iplImage)
	image.iplImage = iplImage
	image.Initialized = true

	size := C.cvGetSize(image.ptr)
	image.size = Size{int(size.width), int(size.height)}
	image.imtype = typeFromDepthAndChannels(iplImage.depth, iplImage.nChannels)

	if image.imtype.NumChannels == 1 {
		image.colorModel = color.GrayModel
	} else if image.imtype.NumChannels == 3 {
		image.colorModel = color.RGBAModel
	} else {
		panic("Unsupported image type - Unsupported number of channels")
	}

	return image, nil
}

func CreateImage(size Size, imtype MatType) (image *Image) {
	var iplImage *C.IplImage

	if imtype.Depth == 8 && imtype.ElemType == Unsigned {
		iplImage = C.cvCreateImage(C.CvSize{C.int(size.Width), C.int(size.Height)}, C.IPL_DEPTH_8U, C.int(imtype.NumChannels))
	}

	image, _ = ImageFromIplImage(iplImage)
	return
}

func (img *Image) Initialize(size Size, imtype MatType) {
	if img.Initialized && img.Size() == size && img.imtype == imtype {
		return
	}

	if img.Initialized {
		img.Release()
	}

	*img = *CreateImage(size, imtype)
}

func NewImage() (image *Image) {
	image = new(Image)
	image.Initialized = false
	image.size = Size{0, 0}
	return
}

func (image *Image) checkInitialized() {
	if !image.Initialized {
		panic("Tried to process non-Initialized image")
	}
}

/******* Resize *******/
type InterpolationType int

const (
	INTER_NN       = C.CV_INTER_NN
	INTER_LINEAR   = C.CV_INTER_LINEAR
	INTER_AREA     = C.CV_INTER_AREA
	INTER_CUBIC    = C.CV_INTER_CUBIC
	INTER_LANCZOS4 = C.CV_INTER_LANCZOS4
)

func (img *Image) ResizeTo(dest *Image, size Size, interp InterpolationType) {
	if !dest.Initialized || dest.Size() != size || dest.imtype != img.imtype {
		dest.Release()
		*dest = *CreateImage(size, img.imtype)
	}

	C.cvResize(img.ptr, dest.ptr, C.int(interp))
}

func (img *Image) Resize(size Size, interp InterpolationType) (res *Image) {
	res = new(Image)
	img.ResizeTo(res, size, interp)
	return
}

/******* Rotate *******/
func (img *Image) RotateTo(dest *Image, angle float64) {
	dest.InitializeAs(img)

	center := C.CvPoint2D32f{C.float(img.Size().Width / 2), C.float(img.Size().Height / 2)}
	var rot *C.CvMat = C.cvCreateMat(2, 3, C.CV_32F)
	defer C.cvReleaseMat(&rot)
	C.cv2DRotationMatrix(center, C.double(angle), C.double(1.0), rot)

	C.cvWarpAffine(img.ptr, dest.ptr, rot, C.CV_INTER_LINEAR+C.CV_WARP_FILL_OUTLIERS, C.cvScalarAll(C.double(0.0)))
}

func (img *Image) Rotate(angle float64) (res *Image) {
	res = new(Image)
	img.RotateTo(res, angle)
	return
}

/******* Blur *******/
func (img *Image) GaussianBlurTo(dest *Image, radius int) {
	dest.InitializeAs(img)
	C.cvSmooth(img.ptr, dest.ptr, C.CV_GAUSSIAN, C.int(radius), C.int(radius), C.double(0.0), C.double(0.0))
}

func (img *Image) GaussianBlur(radius int) (res *Image) {
	res = new(Image)
	img.GaussianBlurTo(res, radius)

	return
}

/******* Copy / Clone *******/
func (img *Image) CopyTo(dest *Image) {
	dest.InitializeAs(img)
	C.cvCopy(img.ptr, dest.ptr, nil)
}

func (img *Image) Clone() (res *Image) {
	res, _ = ImageFromIplImage(C.cvCloneImage(img.iplImage))
	return
}

/******* LUT *******/
type LUT [][256]uint8

func (img *Image) LUTTo(res *Image, lut *LUT) {
	res.InitializeAs(img)

	cvlut := CreateImage(Size{256, 1}, CV_8UC1)
	defer cvlut.Release()
	for i := 0; i < 256; i++ {
		C.cvSetReal2D(cvlut.ptr, C.int(0), C.int(i), C.double((*lut)[0][i]))
	}

	C.cvLUT(img.ptr, res.ptr, cvlut.ptr)
}

/******* Split *******/
type Channels []Image

func (channels *Channels) Release() {
	for _, channel := range *channels {
		channel.Release()
	}
}

func (img *Image) SplitTo(channels *Channels) error {
	img.checkInitialized()

	n := img.imtype.NumChannels

	if n != len(*channels) {
		return errors.New("Trying to split image to wrong number of channels")
	}

	if n == 1 {
		C.cvSplit(img.ptr, (*channels)[0].ptr, nil, nil, nil)
	} else if n == 2 {
		C.cvSplit(img.ptr, (*channels)[0].ptr, (*channels)[1].ptr, nil, nil)
	} else if n == 3 {
		C.cvSplit(img.ptr, (*channels)[0].ptr, (*channels)[1].ptr, (*channels)[2].ptr, nil)
	} else if n == 4 {
		C.cvSplit(img.ptr, (*channels)[0].ptr, (*channels)[1].ptr, (*channels)[2].ptr, (*channels)[3].ptr)
	}

	return nil
}

func (img *Image) Split() (channels Channels) {
	img.checkInitialized()
	channels = make(Channels, img.imtype.NumChannels)

	for i := 0; i < img.imtype.NumChannels; i++ {
		channels[i].Initialize(img.Size(), MatType{img.imtype.Depth, img.imtype.ElemType, 1})
	}

	img.SplitTo(&channels)

	return
}

/******* CvtColor *******/

/************************************************
 * Implementation of Go's Image interface
 ************************************************/
func (image *Image) ColorModel() color.Model {
	return image.colorModel
}

func (image *Image) At(x, y int) (res color.Color) {
	scalar := image.ScalarAt(y, x)
	if image.imtype.NumChannels == 1 {
		res = color.Gray{uint8(scalar[0])}
	} else {
		// While OpenCV represents images as BGR, we ensure that any interaction with images results in an RGB image
		// (for example in cv.LoadImage)
		res = color.NRGBA{uint8(scalar[0]), uint8(scalar[1]), uint8(scalar[2]), 255}
	}

	return
}

func (image *Image) Bounds() goimage.Rectangle {
	p0 := goimage.Point{0, 0}
	p1 := goimage.Point{image.Size().Width, image.Size().Height}
	return goimage.Rectangle{p0, p1}
}
