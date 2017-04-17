package rpc

import (
	"context"
	"io"
	"sync"

	log "github.com/Sirupsen/logrus"
	"gitlab.ricebook.net/platform/core/rpc/gen"

	"google.golang.org/grpc"
)

type CIDs struct {
	sync.Mutex
	ids map[string]interface{}
}

func NewCIDs() *CIDs {
	c := &CIDs{}
	c.ids = map[string]interface{}{}
	return c
}

func (c *CIDs) Add(id string) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.ids[id]; !ok {
		c.ids[id] = new(interface{})
	}
}

func (c *CIDs) GetList() *pb.ContainerIDs {
	c.Lock()
	defer c.Unlock()
	ids := &pb.ContainerIDs{}
	ids.Ids = []*pb.ContainerID{}
	for id, _ := range c.ids {
		ids.Ids = append(ids.Ids, &pb.ContainerID{id})
	}
	return ids
}

func Remove(server string, cids *CIDs) {
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewCoreRPCClient(conn)

	resp, err := c.RemoveContainer(context.Background(), cids.GetList())
	if err != nil {
		log.Fatalf("Run failed %v", err)
	}

	for {
		msg, err := resp.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Message invalid %v", err)
		}
		log.Warnf("[TIMEOUT] %s remove %s", msg.Id[:7], msg.Message)
	}
}
