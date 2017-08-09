package utils

import (
	"testing"

	"github.com/projecteru2/lambda/types"
	"github.com/stretchr/testify/assert"
)

func TestRebuildParams(t *testing.T) {
	params := types.RunParams{}
	defaultConfig := types.DefaultConfig{
		Image:   "default:image",
		Timeout: 600,
	}
	rebuildParams := RebuildParams(params, defaultConfig)

	assert.Equal(t, rebuildParams.Image, defaultConfig.Image)
	assert.NotEqual(t, params.Image, rebuildParams.Image)

	assert.Equal(t, rebuildParams.Timeout, defaultConfig.Timeout)
	assert.NotEqual(t, params.Timeout, rebuildParams.Timeout)
}

func TestPickServer(t *testing.T) {
	servers := []string{"A", "B"}
	var A, B int
	for i := 0; i < 99; i++ {
		s := PickServer(servers)
		switch s {
		case "A":
			A++
		case "B":
			B++
		}
	}
	t.Logf("A: %d", A)
	t.Logf("B: %d", B)
	sub := A - B
	if B-A > 0 {
		sub = B - A
	}
	assert.True(t, sub < 30)
}
