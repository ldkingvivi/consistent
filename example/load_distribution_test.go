package example

import (
	"fmt"
	"github.com/buraksezer/consistent"
	"log"
	"math"
	"math/rand"
	"testing"
)

func Test_LoadDistribution(t *testing.T) {
	members := []consistent.Member{}
	for i := 0; i < 8; i++ {
		member := Member(fmt.Sprintf("node%d.olricmq", i))
		members = append(members, member)
	}
	cfg := consistent.Config{
		PartitionCount:    6113,
		ReplicationFactor: 40,
		Load:              1.0000001,
		Hasher:            hasher{},
	}
	c := consistent.New(members, cfg)

	keyCount := 1000000
	load := (c.AverageLoad() * float64(keyCount)) / float64(cfg.PartitionCount)
	fmt.Println("Maximum key count for a member should be around this: ", math.Ceil(load))
	distribution := make(map[string]int)
	key := make([]byte, 4)
	for i := 0; i < keyCount; i++ {
		rand.Read(key)
		member := c.LocateKey(key)
		distribution[member.String()]++
	}
	for member, count := range distribution {
		fmt.Printf("member: %s, key count: %d\n", member, count)
	}

	m := c.LocateKey([]byte("node1234567oaishfosdbfksaufhakfjbasjkldfiuafkjabsdfhasufbaskjdfbasdlfasdlfhasfsajkdf"))
	log.Printf("locate node to %s", m)
}
