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

func imageFromIplImage(iplImage *C.IplImage) (*Image, error) {
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
