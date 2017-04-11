package utils

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

func GetParams(c *cli.Context) (pod, image, name, command, network, workingDir string,
	envs, volumes []string, cpu float64, mem int64, count, timeout int) {

	name = c.String("name")
	if name == "" {
		log.Fatal("Name missing")
	}

	command = c.String("command")
	if command == "" {
		log.Fatal("Command missing")
	}

	network = c.String("network")
	workingDir = c.String("working-dir")
	image = c.String("image")
	cpu = c.Float64("cpu")
	mem = c.Int64("mem")
	count = c.Int("count")
	timeout = c.Int("timeout")
	envs = c.StringSlice("env")
	volumes = c.StringSlice("volume")
	return
}

func DefaultString(src string, def string) string {
	if src == "" {
		return def
	}
	return src
}

func DefaultFloat64(src float64, def float64) float64 {
	if src == 0.0 {
		return def
	}
	return src
}

func DefaultInt(src int, def int) int {
	if src == 0 {
		return def
	}
	return src
}

func DefaultInt64(src int64, def int64) int64 {
	if src == 0 {
		return def
	}
	return src
}
