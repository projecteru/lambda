package rpc

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"gitlab.ricebook.net/platform/core/rpc/gen"
	"gitlab.ricebook.net/platform/lambda/types"
	"google.golang.org/grpc"
)

const (
	appTmpl = `appname: "lambda"
entrypoints:
  %s:
    cmd: "%s"
    working_dir: "%s"
    run_and_wait_timeout: %d
`
)

// [exitcode] bytes
var EXIT_CODE = []byte{91, 101, 120, 105, 116, 99, 111, 100, 101, 93, 32}
var ENTER = []byte{10}
var SPLIT = []byte{62, 32}

func RunAndWait(server string, runParams types.RunParams) (code int) {
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[RunAndWait] did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewCoreRPCClient(conn)
	opts := generateOpts(runParams)

	resp, err := c.RunAndWait(context.Background())
	if err != nil {
		log.Fatalf("[RunAndWait] Run failed %v", err)
	}

	if resp.Send(&pb.RunAndWaitOptions{DeployOptions: opts}) != nil {
		log.Fatalf("[RunAndWait] Send options failed %v", err)
	}

	if runParams.OpenStdin {
		go func() {
			// 获得输入
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				command := scanner.Bytes()
				log.Debugf("input: %s", command)
				command = append(command, ENTER...)
				if err = resp.Send(&pb.RunAndWaitOptions{Cmd: command}); err != nil {
					log.Errorf("[RunAndWait] Send command %s error: %v", command, err)
				}
			}
			if err := scanner.Err(); err != nil {
				log.Errorf("[RunAndWait] Parse log failed, %v", err)
			}
		}()
	}

	for {
		msg, err := resp.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("[RunAndWait] Message invalid %v", err)
		}

		if bytes.HasPrefix(msg.Data, EXIT_CODE) {
			ret := string(bytes.TrimLeft(msg.Data, string(EXIT_CODE)))
			code, err = strconv.Atoi(ret)
			if err != nil {
				log.Fatalf("[RunAndWait] exit with %s", ret)
			}
			continue
		}
		data := msg.Data
		id := msg.ContainerId[:7]
		if !bytes.HasSuffix(data, SPLIT) {
			data = append(data, ENTER...)
		}
		fmt.Printf("[%s]: %s", id, data)
	}
	return
}

func generateOpts(rp types.RunParams) *pb.DeployOptions {
	for i, env := range rp.Envs {
		rp.Envs[i] = fmt.Sprintf("LAMBDA_%s", env)
	}
	spaces := generateSpecs(rp.Name, rp.Command, rp.Workingdir, rp.Volumes, rp.Timeout)
	opts := &pb.DeployOptions{
		Specs:      spaces,
		Appname:    "lambda",
		Image:      rp.Image,
		Podname:    rp.Pod,
		Entrypoint: rp.Name,
		CpuQuota:   rp.CPU,
		Memory:     rp.Mem,
		Count:      int32(rp.Count),
		Networks:   map[string]string{rp.Network: ""},
		Env:        rp.Envs,
		OpenStdin:  rp.OpenStdin,
	}

	// check opts
	if opts.Count < 0 || (opts.OpenStdin && opts.Count != 1) {
		log.Fatalf("[RunAndWait] Parameter error")
	}

	return opts
}

func generateSpecs(name, command, workingDir string, volumes []string, timeout int) string {
	specs := fmt.Sprintf(appTmpl, name, command, workingDir, timeout)
	if len(volumes) > 0 {
		vol := map[string][]string{}
		vol["volumes"] = volumes
		out, err := yaml.Marshal(vol)
		if err != nil {
			log.Fatalf("Parse failed %v", err)
		}
		specs = fmt.Sprintf("%s%s", specs, out)
	}
	return specs
}
