package bmpnurse

import (
	"log"
	"os"
	"testing"
)

var testData []byte

// avoid bench optimization
var tReport BmpReport
var isBroken bool

func init() {
	data, err := os.ReadFile("img.bmp")
	if err != nil {
		log.Fatal(err)
	}
	testData = data
}

func BenchmarkInspectBmpImage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		report, err := InspectBmpImage(testData)
		if err != nil {
			b.Fatal(err)
		}
		tReport = report
	}
}

func BenchmarkIsBrokenBmp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result, err := IsValidSize(testData)
		if err != nil {
			b.Fatal(err)
		}
		isBroken = result
	}
}
