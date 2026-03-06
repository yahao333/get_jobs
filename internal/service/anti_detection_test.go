package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomDuration(t *testing.T) {
	ad := NewAntiDetection()
	min := 100 * time.Millisecond
	max := 200 * time.Millisecond

	// Test multiple calls to ensure shared random source works
	for i := 0; i < 100; i++ {
		d := ad.randomDuration(min, max)
		assert.True(t, d >= min, "duration should be >= min")
		assert.True(t, d < max, "duration should be < max")
	}
}
