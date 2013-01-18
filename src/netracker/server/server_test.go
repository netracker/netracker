package server

import (
	"github.com/bmizerany/assert"
	"github.com/drewolson/testflight"
	"testing"
)

func TestRoot(t *testing.T) {
	s := New()

	testflight.WithServer(s.handler(), func(r *testflight.Requester) {
		response := r.Get("/")

		assert.Equal(t, 200, response.StatusCode)
	})
}
