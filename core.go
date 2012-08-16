package cv

/*
#cgo pkg-config: opencv
#include <opencv/cv.h>
#include <opencv/highgui.h>
#include <opencv2/imgproc/types_c.h>
*/
import "C"
import (
//"unsafe"
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
