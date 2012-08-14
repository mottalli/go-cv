package cv

/*
#cgo pkg-config: opencv
#include <opencv/cv.h>
#include <opencv/highgui.h>
#include <opencv2/imgproc/types_c.h>
*/
import "C"
import (
	"image/color"
	"unsafe"
)

/************************************************
 * Elementary types
 ************************************************/
type Size struct {
	Width, Height int
}

type Scalar [4]float64

const (
	Signed = iota
	Unsigned
	Float
)

type MatType struct {
	Depth, ElemType, NumChannels int
}

var (
	CV_8UC1 MatType = MatType{8, Unsigned, 1}
	CV_8UC2 MatType = MatType{8, Unsigned, 2}
	CV_8UC3 MatType = MatType{8, Unsigned, 3}
	CV_8UC4 MatType = MatType{8, Unsigned, 4}
	CV_8SC1 MatType = MatType{8, Signed, 1}
	CV_8SC2 MatType = MatType{8, Signed, 2}
	CV_8SC3 MatType = MatType{8, Signed, 3}
	CV_8SC4 MatType = MatType{8, Signed, 4}

	CV_16UC1 MatType = MatType{16, Unsigned, 1}
	CV_16UC2 MatType = MatType{16, Unsigned, 2}
	CV_16UC3 MatType = MatType{16, Unsigned, 3}
	CV_16UC4 MatType = MatType{16, Unsigned, 4}
	CV_16SC1 MatType = MatType{16, Signed, 1}
	CV_16SC2 MatType = MatType{16, Signed, 2}
	CV_16SC3 MatType = MatType{16, Signed, 3}
	CV_16SC4 MatType = MatType{16, Signed, 4}

	CV_32UC1 MatType = MatType{32, Unsigned, 1}
	CV_32UC2 MatType = MatType{32, Unsigned, 2}
	CV_32UC3 MatType = MatType{32, Unsigned, 3}
	CV_32UC4 MatType = MatType{32, Unsigned, 4}
	CV_32SC1 MatType = MatType{32, Signed, 1}
	CV_32SC2 MatType = MatType{32, Signed, 2}
	CV_32SC3 MatType = MatType{32, Signed, 3}
	CV_32SC4 MatType = MatType{32, Signed, 4}
	CV_32FC1 MatType = MatType{32, Float, 1}
	CV_32FC2 MatType = MatType{32, Float, 2}
	CV_32FC3 MatType = MatType{32, Float, 3}
	CV_32FC4 MatType = MatType{32, Float, 4}

	CV_64FC1 MatType = MatType{64, Float, 1}
	CV_64FC2 MatType = MatType{64, Float, 2}
	CV_64FC3 MatType = MatType{64, Float, 3}
	CV_64FC4 MatType = MatType{64, Float, 4}
)

type Mat interface {
	NativePointer() unsafe.Pointer
	Release()
	ScalarAt(pos ...int) Scalar
	Size() Size
	Type() MatType
}

/**************************************************
 * Implementation of the IplImage structure
 **************************************************/
type Image struct {
	// Private attributes
	iplImage   *C.IplImage
	ptr        unsafe.Pointer
	colorModel color.Model
	imtype     MatType
	size       Size

	// Exported attributes
	Initialized bool
}

func (img *Image) NativePointer() unsafe.Pointer {
	return img.ptr
}

func (img *Image) ScalarAt(pos ...int) Scalar {
	n := len(pos)
	var s C.CvScalar
	if n == 1 {
		s = C.cvGet1D(img.ptr, C.int(pos[0]))
	} else if n == 2 {
		s = C.cvGet2D(img.ptr, C.int(pos[0]), C.int(pos[1]))
	} else if n == 3 {
		s = C.cvGet3D(img.ptr, C.int(pos[0]), C.int(pos[1]), C.int(pos[2]))
	}

	return Scalar{float64(s.val[0]), float64(s.val[1]), float64(s.val[2]), float64(s.val[3])}
}

func (img *Image) Release() {
	if img.iplImage != nil {
		C.cvReleaseImage(&img.iplImage)
		img.ptr = nil
		img.iplImage = nil
	}

	img.Initialized = false
}

func (img *Image) Size() Size {
	return img.size
}

func (img *Image) Type() MatType {
	return img.imtype
}
