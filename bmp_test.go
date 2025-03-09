package bmpnurse

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestBuildBmpHeaderLe(t *testing.T) {
	// 14 bytes of bmp header
	header := []byte{66, 77, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	hdr := getBmpHeader(header)
	tp := BfTypesMapLe[hdr.bfType]

	require.Equal(t, "BM", tp)
}

func TestValidateBmpHeader(t *testing.T) {
	header := []byte{66, 77, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	hdr := getBmpHeader(header)
	valid := validateBmpHeaderLeOrder(hdr)
	require.Equal(t, true, valid)
}

func TestBuildBmpHeaderBe(t *testing.T) {
	// 14 bytes of bmp header
	// check big-endian option
	header := []byte{77, 66, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	hdr := getBmpHeader(header)

	tp := BfTypesMapBe[hdr.bfType]
	require.Equal(t, "BM", tp)
}

func TestDetectDIBHeader(t *testing.T) {
	header := []byte{66, 77, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12, 0, 0, 0}
	info := "BITMAPCOREHEADER 12 bytes"

	dibSize := detectDibHeader(header)
	require.Equal(t, int32(12), dibSize)

	DIBType, ok := DibHeaderTypeMap[dibSize]
	require.Equal(t, true, ok)
	require.Equal(t, info, DIBType)
}

func TestSetBitsPerPixel(t *testing.T) {
	// 8 = 24 bit (little-endian)
	header := []byte{
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		8,
		0,
	}

	pixels := setBitsPerPixel(header, int32(12))
	require.Equal(t, int16(8), pixels)
}

func TestIsBrokenBmpRecoverPanic(t *testing.T) {
	// we got panic because of empty data
	var data []byte
	_, err := IsValidSize(data)
	require.Error(t, err)
}

func TestInspectBmpImageRecoverPanic(t *testing.T) {
	// we got panic because of empty data
	var data []byte
	_, err := InspectBmpImage(data)
	require.Error(t, err)
}

// test with real valid bmp file
func TestInspectBmpImage(t *testing.T) {
	data, _ := os.ReadFile("img.bmp")

	report, _ := InspectBmpImage(data)
	require.Equal(t, true, report.HaveValidSize)
}

func TestIsBrokenBmp(t *testing.T) {
	data, _ := os.ReadFile("img.bmp")

	broken, _ := IsValidSize(data)
	require.Equal(t, false, broken)
}
