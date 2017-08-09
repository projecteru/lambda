package rpc

import (
	"fmt"
	"net"
	"testing"
	"time"

	coreTypes "github.com/projecteru2/core/types"
	pb "github.com/projecteru2/lambda/rpc/mock_grpc"
	"github.com/projecteru2/lambda/types"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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

type mockCoreSrv struct{}

func (m *mockCoreSrv) RunAndWait(s pb.CoreRPC_RunAndWaitServer) error {
	opts, err := s.Recv()
	if err != nil {
		return err
	}

	spec, err := coreTypes.LoadSpecs(opts.DeployOptions.Specs)
	if err != nil {
		return err
	}
	e := spec.Entrypoints["hello"]
	timeout := e.RunAndWaitTimeout
	cmd := e.Command
	for i := 0; i < 5 && timeout > 0; i++ {
		t := time.Now().Second()
		r := fmt.Sprintf("run %v output: %d", cmd, t)
		s.Send(&pb.RunAndWaitMessage{
			ContainerId: "helloworld",
			Data:        []byte(r),
		})
		time.Sleep(1 * time.Second)
		timeout--
	}
	return nil
}

func TestRunAndWait(t *testing.T) {
	s, err := net.Listen("tcp", "127.0.0.1:8866")
	assert.NoError(t, err)
	grpcServer := grpc.NewServer()
	srv := &mockCoreSrv{}
	pb.RegisterCoreRPCServer(grpcServer, srv)
	go grpcServer.Serve(s)

	server := "127.0.0.1:8866"
	runParams := types.RunParams{
		Command: "date",
		Name:    "hello",
		Timeout: 5,
	}
	code := RunAndWait(server, runParams)
	assert.Equal(t, 0, code)
	grpcServer.GracefulStop()
}

//func TestRunAndWaitTimeout(t *testing.T) {
//	s, err := net.Listen("tcp", "127.0.0.1:8866")
//	assert.NoError(t, err)
//	grpcServer := grpc.NewServer()
//	srv := &mockCoreSrv{}
//	pb.RegisterCoreRPCServer(grpcServer, srv)
//	go grpcServer.Serve(s)
//
//	timeout := 1
//	timer := time.NewTimer(time.Duration(timeout+1) * time.Second)
//	runCh := make(chan interface{})
//	server := "127.0.0.1:8866"
//	runParams := types.RunParams{
//		Command: "date",
//		Name:    "hello",
//		Timeout: timeout,
//	}
//	go func() {
//		code := RunAndWait(server, runParams)
//		assert.Equal(t, 0, code)
//		runCh <- code
//	}()
//	select {
//	case <-timer.C:
//		t.Error("timeout failed!")
//	case <-runCh:
//		grpcServer.GracefulStop()
//	}
//}
