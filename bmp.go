package bmpnurse

import (
	"fmt"
	"unsafe"
)

const (
	// WinBmImageLe little-endian repr of image type field
	// is used le format
	WinBmImageLe   int16 = 0x4D42
	WinBmImageBe   int16 = 0x424D
	OS2BmImageLe   int16 = 0x4142
	OS2BmImageBe   int16 = 0x4241
	OS2BmImageCI   int16 = 0x4943
	OS2BmImageCIBe int16 = 0x4349
	OS2BmImageCP   int16 = 0x5043
	OS2BmImageCPBe int16 = 0x4350
	OS2BmImageIC   int16 = 0x4349
	OS2BmImageICBe int16 = 0x4943
	OS2BmImagePT   int16 = 0x5450
	OS2BmImagePTBe int16 = 0x5054

	// BITMAPCOREHEADER have the same size = 12 bytes and
	// the same structure as OS21XBITMAPHEADER
	BITMAPCOREHEADER    int32 = 12
	OS22XBITMAPHEADER64 int32 = 64
	BITMAPINFOHEADER    int32 = 40
	BITMAPV4HEADER      int32 = 108
	BITMAPV5HEADER      int32 = 124
)

var DibHeaderTypeMap = map[int32]string{
	BITMAPCOREHEADER:    "BITMAPCOREHEADER 12 bytes",
	OS22XBITMAPHEADER64: "OS22XBITMAPHEADER64 64 bytes",
	BITMAPINFOHEADER:    "BITMAPINFOHEADER 40 bytes",
	BITMAPV4HEADER:      "BITMAPV4HEADER 108 bytes",
	BITMAPV5HEADER:      "BITMAPV5HEADER 124 bytes",
}

var BfTypesMapLe = map[int16]string{
	WinBmImageLe: "BM",
	OS2BmImageLe: "BA",
	OS2BmImageCI: "CI",
	OS2BmImageCP: "CP",
	OS2BmImageIC: "IC",
	OS2BmImagePT: "PT",
}

var BfTypesMapBe = map[int16]string{
	WinBmImageBe:   "BM",
	OS2BmImageBe:   "BA",
	OS2BmImageCIBe: "CI",
	OS2BmImageCPBe: "CP",
	OS2BmImageICBe: "IC",
	OS2BmImagePTBe: "PT",
}

// BmpHeader header of bmp image
type BmpHeader struct {
	// type of bitmap. BM from Windows
	bfType int16

	// image size bytes
	bfSizeBytes int32

	// reserved fields as described in bmp spec
	bfReserved1 int16
	bfReserved2 int16

	// byte array address (offset)
	bfOffBits int32
}

// BmpReport provide base info about bmp image
type BmpReport struct {
	// bmp protocol version
	BfType string

	// headers type
	DibHeadersType string
	ByteOrder      string

	// how many bits per pixel is
	BitsPerPixel int16
	HaveHeader   bool

	// flag is true when DeclaredSize == ActualSize
	HaveValidSize bool
	DeclaredSize  int32
	ActualSize    int32
}

// InspectBmpImage is user to inspect received .bmp file
func InspectBmpImage(rawData []byte) (report BmpReport, err error) {
	var dibHdrType int32

	// we have to work with int32 bcs of bmp documentation
	dataSize := int32(len(rawData))

	// we may got empty or invalid byte array, and we got
	// panic when we try got any value by pointer, so we
	// have to recover and return error
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("gobmp recover: '%+v'", v)
		}
	}()

	bmpHeader := getBmpHeader(rawData)
	if ok := validateBmpHeaderLeOrder(bmpHeader); ok {
		report.HaveHeader = ok
		report.ByteOrder = "little-endian"
		report.BfType = BfTypesMapLe[bmpHeader.bfType]
		// check file size

		dibHdrType = detectDibHeader(rawData)
	} else {
		// maybe we have not BmpHeader, let`s check DIB
		if ok = validateBmpHeaderBeOrder(bmpHeader); ok {
			report.HaveHeader = ok
			report.ByteOrder = "big-endian"
			report.BfType = BfTypesMapBe[bmpHeader.bfType]

			// detect DIB after initial header
			dibHdrType = detectDibHeader(rawData)
		} else {

			// check file size
			dibHdrType = detectDibHeaderAsFirst(rawData)
		}
	}

	report.HaveValidSize = false
	report.DeclaredSize = bmpHeader.bfSizeBytes
	report.ActualSize = dataSize

	// compare actual and expected image size (in bytes)
	if bmpHeader.bfSizeBytes != 0 && report.ActualSize == report.DeclaredSize {
		report.HaveValidSize = true
	}

	if info, ok := DibHeaderTypeMap[dibHdrType]; ok {
		report.DibHeadersType = info
	}

	report.BitsPerPixel = setBitsPerPixel(rawData, dibHdrType)
	return report, err
}

// IsValidSize compare expected and actual image size (in bytes)
func IsValidSize(rawData []byte) (broken bool, err error) {
	dataSize := int32(len(rawData))

	// we may got empty or invalid byte array, and we got
	// panic when we try got any value by pointer, so we
	// have to recover and return error
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("gobmp recover: '%+v'", v)
		}
	}()

	// we use optimization (we need only size field)
	size := getSizeBytes(rawData)
	return size != dataSize, err
}

func validateBmpHeaderLeOrder(hdr BmpHeader) bool {
	switch hdr.bfType {
	case WinBmImageLe, OS2BmImageLe, OS2BmImageCI,
		OS2BmImageIC, OS2BmImageCP, OS2BmImagePT:
		return true
	default:
		return false
	}
}

func validateBmpHeaderBeOrder(hdr BmpHeader) bool {
	switch hdr.bfType {
	case WinBmImageBe, OS2BmImageBe, OS2BmImageCIBe,
		OS2BmImageICBe, OS2BmImageCPBe, OS2BmImagePTBe:
		return true
	default:
		return false
	}
}

// TODO: refactor. Learn unsafe
func getBmpHeader(rd []byte) (hdr BmpHeader) {
	hdr.bfType = *(*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&rd[0])) + 0))
	hdr.bfSizeBytes = *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&rd[0])) + 2))
	hdr.bfReserved1 = *(*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&rd[0])) + 6))
	hdr.bfReserved2 = *(*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&rd[0])) + 8))
	hdr.bfOffBits = *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&rd[0])) + 10))
	return hdr
}

// getSizeBytes is an optimization when we need only declared image size
func getSizeBytes(rd []byte) int32 {
	return *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&rd[0])) + 2))
}

func detectDibHeader(rd []byte) int32 {
	// let`s check DIB header size using BmpHeader offset = 14 bytes
	return *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&rd[0])) + 14))
}

func detectDibHeaderAsFirst(rd []byte) int32 {
	// no offset is set.
	// we hope that DibHeader is first (from 0 pos)
	// in case we have .bmp in memory
	//
	// Size of DIB header is always 1st parameter
	return *(*int32)(unsafe.Pointer(&rd[0]))
}

func setBitsPerPixel(rd []byte, dibSize int32) int16 {
	switch dibSize {
	case BITMAPCOREHEADER, OS22XBITMAPHEADER64:
		// this header have offset to wished field = 24 bytes
		return *(*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&rd[0])) + 24))
	case BITMAPINFOHEADER, BITMAPV4HEADER, BITMAPV5HEADER:
		// all versions alter than BITMAPINFOHEADER add new
		// headers at the end of previous version
		return *(*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&rd[0])) + 28))
	default:
		return -1
	}
}
