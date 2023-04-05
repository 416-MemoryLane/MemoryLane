package tests

import (
	"log"
	"memory-lane/app/papaya"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
)

func TestData(t *testing.T) {

}

const TEST_GALLERY_DIR = "./test-gallery"

var _ = BeforeSuite(func() {
	clearTestGallery()
	initGallery()
})

var _ = AfterSuite(func() {
	clearTestGallery()
})

var _ = BeforeEach(func() {
	clearTestGallery()
	initGallery()
})

func clearTestGallery() {
	err := os.Remove(TEST_GALLERY_DIR)
	if err != nil {
		panic(err)
	}
}

func initGallery() *papaya.Gallery {
	g, err := papaya.NewGallery(TEST_GALLERY_DIR, log.New(os.Stdout, "", log.Lshortfile|log.Ltime))
	if err != nil {
		panic(err)
	}

	return g
}
