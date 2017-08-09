package rpc

import (
	"testing"

	"github.com/projecteru2/lambda/types"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSpecs(t *testing.T) {
	name := "name"
	command := "date"
	workingDir := "/temp"
	volumes := []string{"/share"}
	timeout := 60
	specs := generateSpecs(name, command, workingDir, volumes, timeout)
	t.Log(specs)
}

func TestGenerateOpts(t *testing.T) {
	rps := types.RunParams{
		Image: "testImage",
	}
	pbDeployOpts := generateOpts(rps)
	assert.Equal(t, rps.Image, pbDeployOpts.Image)
}
