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

	pod, image, name, command, network, envs, cpu, mem, count, timeout := utils.GetParams(c)
	if admin {
		pod = config.AdminPod
	}
	if network == "" {
		network = config.DefaultSDN
	}
	if image == "" {
		image = config.BaseImage
	}

	server := utils.PickServer(config.Servers)
	code := rpc.RunAndWait(server, pod, image, name, command, network, envs, cpu, mem, count, timeout)
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
			Name:  "pod",
			Usage: "where to run",
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
			Usage: "set envs can use multiple times",
		},
		&cli.StringFlag{
			Name:  "network",
			Usage: "SDN name (default: define in config file)",
		},
		&cli.StringFlag{
			Name:  "image",
			Usage: "use image (default: define in config file)",
		},
		&cli.IntFlag{
			Name:  "timeout",
			Usage: "when to interrupt",
			Value: 10,
		},
		&cli.Float64Flag{
			Name:  "cpu",
			Usage: "how many cpu",
			Value: 1.0,
		},
		&cli.Int64Flag{
			Name:  "mem",
			Usage: "how many memory in bytes",
			Value: 536870912,
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
