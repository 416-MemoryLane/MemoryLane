package tests

import (
	"fmt"
	"memory-lane/app/papaya"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Reaching consensus", func() {
	var g *papaya.Gallery

	When("the app progresses through a happy path", func() {
		BeforeEach(func() {
			clearTestGallery()
		})

		It("", func() {
			g = initGallery()
			fmt.Printf("testing with: %v", g)
		})
	})

	When("the app receives messages out of order", func() {
		BeforeEach(func() {
			clearTestGallery()
		})

		It("", func() {
			g = initGallery()
			fmt.Printf("testing with: %v", g)
		})
	})

})
