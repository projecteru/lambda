package utils

import (
	log "github.com/Sirupsen/logrus"
	"gitlab.ricebook.net/platform/lambda/types"
	"gopkg.in/urfave/cli.v2"
)

func GetParams(c *cli.Context) types.RunParams {

	runParams := types.RunParams{
		Name:       c.String("name"),
		Command:    c.String("command"),
		Network:    c.String("network"),
		Workingdir: c.String("working-dir"),
		Image:      c.String("image"),
		CPU:        c.Float64("cpu"),
		Mem:        c.Int64("mem"),
		Count:      c.Int("count"),
		Timeout:    c.Int("timeout"),
		Envs:       c.StringSlice("env"),
		Volumes:    c.StringSlice("volume"),
		OpenStdin:  c.Bool("interactive"),
	}

	if runParams.Name == "" {
		log.Fatal("Name missing")
	}

	if runParams.Command == "" {
		log.Fatal("Command missing")
	}

	return runParams
}

func DefaultString(runParams string, defaultConfig string) string {
	if runParams == "" {
		return defaultConfig
	}
	return runParams
}

func DefaultFloat64(runParams float64, defaultConfig float64) float64 {
	if runParams == 0.0 {
		return defaultConfig
	}
	return runParams
}

func DefaultInt(runParams int, defaultConfig int) int {
	if runParams == 0 {
		return defaultConfig
	}
	return runParams
}

func DefaultInt64(runParams int64, defaultConfig int64) int64 {
	if runParams == 0 {
		return defaultConfig
	}
	return runParams
}

func DefaultBool(runParams, defaultConfig bool) bool {
	if runParams == false {
		return defaultConfig
	}
	return runParams
}

func RebuildParams(runParams types.RunParams, defaultConfig types.DefaultConfig) types.RunParams {
	runParams.Pod = DefaultString(runParams.Pod, defaultConfig.Pod)
	runParams.Network = DefaultString(runParams.Network, defaultConfig.Network)
	runParams.Workingdir = DefaultString(runParams.Workingdir, defaultConfig.WorkingDir)
	runParams.Image = DefaultString(runParams.Image, defaultConfig.Image)
	runParams.CPU = DefaultFloat64(runParams.CPU, defaultConfig.Cpu)
	runParams.Mem = DefaultInt64(runParams.Mem, defaultConfig.Memory)
	runParams.Timeout = DefaultInt(runParams.Timeout, defaultConfig.Timeout)
	runParams.OpenStdin = DefaultBool(runParams.OpenStdin, defaultConfig.OpenStdin)
	return runParams
}
