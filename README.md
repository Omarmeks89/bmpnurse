# bmpnurse
small library for working with bmp images

You may check that `bmp` image is not broken:
```go
// some func that load bmp image
import "github.com/Omarmeks89/bmpnurse"

func HandleBmp(fName string) (bmpData []byte, err error){
	var valid bool
	
	bmpData, err = os.ReadFile(fname)
	if err != nil { 
		// ... 
	}
	
	// check that bmp size is valid (image is not broken)
	valid, err = bmpnurse.IsValidSize(bmpData)
	if !valid {
		return bmpData, err
	}
	// do sth useful
	// ...
}
```

or get more information creating `report` about image:
```go
import (
	"github.com/Omarmeks89/bmpnurse"
	"fmt"
)

func PrintBmpByteOrder(fName string) error {
	var report bmpnurse.BmpReport
	
	bmpData, err := os.ReadFile(fname)
	if err != nil { 
		// ... 
	}
	
	report, err = bmpnurse.InspectBmpImage(bmpData)
	if err != nil {
		return err		
	}
	
	fmt.Printf("'%v'\n", report.ByteOrder)	
}
```
as result we got:
```bash
'little-endian'
```

`BmpReport` looks like:
```go
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
```

## benchmark

Below benchmark results:
```bash
goos: linux
goarch: amd64
pkg: gobmp
cpu: Intel(R) Core(TM) i5-7300HQ CPU @ 2.50GHz
BenchmarkInspectBmpImage
BenchmarkInspectBmpImage-4      26053531                46.15 ns/op
BenchmarkIsBrokenBmp
BenchmarkIsBrokenBmp-4          179866705                6.692 ns/op
```
