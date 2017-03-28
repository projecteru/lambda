package utils

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

func GetParams(c *cli.Context) (pod, image, repo, name, command, network string,
	envs []string, cpu float64, mem int64, count int) {

	pod = c.String("pod")
	if pod == "" {
		log.Fatal("Need specify a pod")
	}

	name = c.String("name")
	if name == "" {
		log.Fatal("Name missing")
	}

	command = c.String("command")
	if command == "" {
		log.Fatal("Command missing")
	}

	network = c.String("network")
	repo = c.String("repo")
	image = c.String("image")
	cpu = c.Float64("cpu")
	mem = c.Int64("mem")
	count = c.Int("count")
	envs = c.StringSlice("env")
	return
}
