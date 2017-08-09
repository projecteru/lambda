package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/projecteru2/lambda/rpc"
	"github.com/projecteru2/lambda/types"
	"github.com/projecteru2/lambda/utils"
	"github.com/projecteru2/lambda/versioninfo"
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

	runParams := utils.GetParams(c)
	if admin {
		runParams.Pod = config.Default.AdminPod
		for _, v := range config.Default.AdminVolumes {
			runParams.Volumes = append(runParams.Volumes, v)
		}
	}

	if runParams.Count > config.Concurrency {
		log.Fatalf("Max concurrency limit %d", config.Concurrency)
	}

	runParams = utils.RebuildParams(runParams, config.Default)

	server := utils.PickServer(config.Servers)
	code := rpc.RunAndWait(server, runParams)
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
		&cli.BoolFlag{
			Name:    "stdin",
			Usage:   "open stdin for container",
			Aliases: []string{"s"},
			Value:   false,
		},
	}
	app.Action = runLambda
	app.Run(os.Args)
}
