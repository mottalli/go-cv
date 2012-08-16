package cv

/*
#cgo pkg-config: opencv
#include <opencv/cv.h>
#include <opencv/highgui.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

func LoadImage(filename string, loadColor bool) (*Image, error) {
	cname := C.CString(filename)
	defer C.free(unsafe.Pointer(cname))

	if !loadColor {
		iplgray := C.cvLoadImage(cname, C.int(0))
		if iplgray == nil {
			return nil, errors.New("Could not open file: " + filename)
		}
		return imageFromIplImage(iplgray)
	}

	iplbgr := C.cvLoadImage(cname, C.int(1))
	if iplbgr == nil {
		return nil, errors.New("Could not open file: " + filename)
	}

	bgr, _ := imageFromIplImage(iplbgr)
	defer bgr.Release()

	rgb := new(Image).InitializeAs(bgr)
	C.cvCvtColor(bgr.ptr, rgb.ptr, C.CV_BGR2RGB)

	return rgb, nil
}

func Show(image *Image, windowName string) {
	cname := C.CString(windowName)
	defer C.free(unsafe.Pointer(cname))

	C.cvShowImage(cname, image.ptr)
}

func WaitKey(miliseconds int) rune {
	key := C.cvWaitKey(C.int(miliseconds))
	return rune(key)
}

/**************************************************
 * Video capture
 **************************************************/

type Capture struct {
	handler *C.CvCapture
}

func CaptureFromCam(device int) (cap Capture, err error) {
	cap.handler = C.cvCaptureFromCAM(C.int(device))
	if cap.handler == nil {
		err = errors.New("Could not open capture device")
	}
	return
}

func (cap *Capture) QueryFrame() (*Image, error) {
	iplImage := C.cvQueryFrame(cap.handler)
	return imageFromIplImage(iplImage)
}

func (cap *Capture) Close() {
	C.cvReleaseCapture(&cap.handler)
}
