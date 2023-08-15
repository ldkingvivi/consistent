package example

import (
	"fmt"
	"testing"

	"github.com/buraksezer/consistent"
)

func Test_AddAndList(t *testing.T) {
	members := []consistent.Member{}
	for i := 0; i < 8; i++ {
		member := Member(fmt.Sprintf("node%d.olricmq", i))
		members = append(members, member)
	}
	cfg := consistent.Config{
		PartitionCount:    71,
		ReplicationFactor: 20,
		Load:              1.25,
		Hasher:            hasher{},
	}

	c := consistent.New(members, cfg)
	owners := make(map[string]int)
	for partID := 0; partID < cfg.PartitionCount; partID++ {
		owner := c.GetPartitionOwner(partID)
		owners[owner.String()]++
	}
	fmt.Println("average load:", c.AverageLoad())

	for o, n := range owners {
		fmt.Printf("%+v:%d\n", o, n)
	}
}
