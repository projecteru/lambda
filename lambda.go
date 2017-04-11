package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"gitlab.ricebook.net/platform/lambda/rpc"
	"gitlab.ricebook.net/platform/lambda/types"
	"gitlab.ricebook.net/platform/lambda/utils"
	"gitlab.ricebook.net/platform/lambda/versioninfo"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
)

var (
	debug bool
	admin bool
)

func setupLog(l string) error {
	level, err := log.ParseLevel(l)
	if err != nil {
		return err
	}
	log.SetLevel(level)

	formatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	log.SetFormatter(formatter)
	return nil
}

func initConfig(configPath string) (types.Config, error) {
	config := types.Config{}

	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return config, err
	}

	return config, nil
}

func runLambda(c *cli.Context) error {
	if debug {
		setupLog("DEBUG")
	} else {
		setupLog("INFO")
	}

	config, err := initConfig(c.String("config"))
	if err != nil {
		log.Fatal(err)
	}

	pod, image, name, command, network, workingDir, envs, volumes, cpu, mem, count, timeout := utils.GetParams(c)
	if admin {
		pod = config.Default.AdminPod
		for _, v := range config.Default.AdminVolumes {
			volumes = append(volumes, v)
		}
	}

	if count > config.Concurrency {
		log.Fatalf("Max concurrency limit %d", config.Concurrency)
	}

	pod = utils.DefaultString(pod, config.Default.Pod)
	network = utils.DefaultString(network, config.Default.Network)
	workingDir = utils.DefaultString(workingDir, config.Default.WorkingDir)
	image = utils.DefaultString(image, config.Default.Image)
	cpu = utils.DefaultFloat64(cpu, config.Default.Cpu)
	mem = utils.DefaultInt64(mem, config.Default.Memory)
	timeout = utils.DefaultInt(timeout, config.Default.Timeout)

	server := utils.PickServer(config.Servers)
	code := rpc.RunAndWait(server, pod, image, name, command,
		network, workingDir, envs, volumes, cpu, mem, count, timeout)
	if code == 0 {
		return nil
	}
	return cli.Exit("", code)
}

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Print(versioninfo.VersionString())
	}
}

func main() {
	app := cli.App{}
	app.Name = versioninfo.NAME
	app.Usage = "run code on eru"
	app.Version = versioninfo.VERSION
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Value: "/etc/eru/lambda.yaml",
			Usage: "config file path for lambda, in yaml",
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "name for this lambda",
		},
		&cli.StringFlag{
			Name:  "command",
			Usage: "how to run it",
		},
		&cli.StringSliceFlag{
			Name:  "env",
			Usage: "set env can use multiple times",
		},
		&cli.StringSliceFlag{
			Name:  "volume",
			Usage: "set volume can use multiple times",
		},
		&cli.StringFlag{
			Name:        "pod",
			Usage:       "where to run",
			DefaultText: "in config file",
		},
		&cli.StringFlag{
			Name:        "network",
			Usage:       "SDN name",
			DefaultText: "in config file",
		},
		&cli.StringFlag{
			Name:        "working-dir",
			Usage:       "use as current working dir",
			DefaultText: "in config file",
		},
		&cli.StringFlag{
			Name:        "image",
			Usage:       "base image for running",
			DefaultText: "in config file",
		},
		&cli.Float64Flag{
			Name:        "cpu",
			Usage:       "how many cpu",
			DefaultText: "in config file",
		},
		&cli.Int64Flag{
			Name:        "mem",
			Usage:       "how many memory in bytes",
			DefaultText: "in config file",
		},
		&cli.IntFlag{
			Name:        "timeout",
			Usage:       "when to interrupt",
			DefaultText: "in config file",
		},
		&cli.IntFlag{
			Name:  "count",
			Usage: "how many containers",
			Value: 1,
		},
		&cli.BoolFlag{
			Name:        "admin",
			Usage:       "admin or not",
			Value:       false,
			Destination: &admin,
		},
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "enable debug",
			Aliases:     []string{"d"},
			Value:       false,
			Destination: &debug,
		},
	}
	app.Action = runLambda
	app.Run(os.Args)
}
