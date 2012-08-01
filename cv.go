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
	"fmt"
	goimage "image"
	"image/color"
	"unsafe"
)

/************************************************
 * Elementary types
 ************************************************/
type Size struct {
	Width, Height int
}

type Image struct {
	// Private attributes
	iplImage    *C.IplImage
	ptr         unsafe.Pointer
	colorModel  color.Model
	initialized bool

	// Exported attributes
	Size     Size
	Data     unsafe.Pointer
	Step     int
	Channels int
	Depth    int
}

func (img *Image) Release() {
	if img.iplImage != nil {
		C.cvReleaseImage(&img.iplImage)
	}
}

func (img *Image) InitializeAs(other *Image) *Image {
	if img.initialized && (img.Size == other.Size) && (img.Channels == other.Channels) && (img.Depth == other.Depth) {
		// Do nothing if they are of the same size and type
		return img
	}

	img.Release()
	tmp := CreateImage(other.Size, other.Depth, other.Channels)
	*img = *tmp
	return img
}

func ImageFromIplImage(iplImage *C.IplImage) (*Image, error) {
	image := new(Image)
	image.ptr = unsafe.Pointer(iplImage)
	image.iplImage = iplImage
	image.initialized = true

	size := C.cvGetSize(image.ptr)
	image.Size = Size{int(size.width), int(size.height)}
	image.Depth = int(iplImage.depth)
	image.Channels = int(iplImage.nChannels)
	image.Data = unsafe.Pointer(iplImage.imageData)
	image.Step = int(iplImage.widthStep)

	if image.Depth != C.IPL_DEPTH_8U {
		return nil, errors.New(fmt.Sprintf("Unsupported image depth (%v) - Not 8 bit/channel", image.Depth))
	}

	if image.Channels == 1 {
		image.colorModel = color.GrayModel
	} else if image.Channels == 3 {
		image.colorModel = color.RGBAModel
	} else {
		panic("Unsupported image type - Unsupported number of channels")
	}

	return image, nil
}

/******* CreateImage *******/
const (
	CV_8U        = C.CV_8U
	CV_8S        = C.CV_8S
	CV_16S       = C.CV_16S
	CV_32S       = C.CV_32S
	CV_32F       = C.CV_32F
	CV_64F       = C.CV_64F
	IPL_DEPTH_8U = C.IPL_DEPTH_8U
)

func CreateImage(size Size, depth int, nChannels int) (image *Image) {
	iplImage := C.cvCreateImage(C.CvSize{C.int(size.Width), C.int(size.Height)}, C.int(depth), C.int(nChannels))
	image, _ = ImageFromIplImage(iplImage)
	image.initialized = true
	return
}

func NewImage() (image *Image) {
	image = new(Image)
	image.initialized = false
	image.Channels = 0
	image.Size = Size{0, 0}
	return
}

func (image *Image) checkInitialized() {
	if !image.initialized {
		panic("Tried to process non-initialized image")
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
	if !dest.initialized || dest.Size != size || dest.Channels != img.Channels || dest.Channels != img.Channels {
		dest.Release()
		*dest = *CreateImage(size, img.Depth, img.Channels)
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

	center := C.CvPoint2D32f{C.float(img.Size.Width / 2), C.float(img.Size.Height / 2)}
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
func (img *Image) GaussianBlur(radius int) (res *Image) {
	res = new(Image)
	res.InitializeAs(img)

	C.cvSmooth(img.ptr, res.ptr, C.CV_GAUSSIAN, C.int(radius), C.int(radius), C.double(0.0), C.double(0.0))

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
	
	cvlut := CreateImage(Size{256, 1}, img.Depth, 1)
	defer cvlut.Release()
	for i := 0; i < 256; i++ {
		C.cvSetReal2D(cvlut.ptr, C.int(0), C.int(i), C.double((*lut)[0][i]))
	}
	
	C.cvLUT(img.ptr, res.ptr, cvlut.ptr)
}

/******* CvtColor *******/
/*func (img *Image) CvtColorTo(dest *Image) {
	dest.InitializeAs(img)\
}*/

/************************************************
 * Implementation of Go's Image interface
 ************************************************/
func (image *Image) ColorModel() color.Model {
	return image.colorModel
}

func (image *Image) At(x, y int) (res color.Color) {
	scalar := C.cvGet2D(image.ptr, C.int(y), C.int(x))
	if image.Channels == 1 {
		res = color.Gray{uint8(scalar.val[0])}
	} else {
		// While OpenCV represents images as BGR, we ensure that any interaction with images results in an RGB image
		// (for example in cv.LoadImage)
		res = color.NRGBA{uint8(scalar.val[0]), uint8(scalar.val[1]), uint8(scalar.val[2]), 255}
	}

	return
}

func (image *Image) Bounds() goimage.Rectangle {
	p0 := goimage.Point{0, 0}
	p1 := goimage.Point{image.Size.Width, image.Size.Height}
	return goimage.Rectangle{p0, p1}
}
