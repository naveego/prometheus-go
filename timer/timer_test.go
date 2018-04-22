package timer_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/naveego/prometheus-go/timer"
)

var _ = Describe("MemoryTimer", func() {
	It("Should track the duration of an operation", func() {
		t := &MemoryTimer{}
		Expect(t.Elapsed()).To(Equal(time.Duration(0)), "Expected Elapsed to be 0 after creation")

		t.Start()
		time.Sleep(time.Second * 1)
		t.Stop()
		Expect(t.Elapsed()).To(BeNumerically(">", 0))
	})

	Describe("Start", func() {
		It("Should return the timer instance", func() {
			t := &MemoryTimer{}
			Expect(t.Start()).To(Equal(t))
		})
	})
})
