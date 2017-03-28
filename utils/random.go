package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func PickServer(servers []string) string {
	l := len(servers)
	i := rand.Int() % l
	return servers[i]
}
