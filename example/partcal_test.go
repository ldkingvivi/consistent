package example

import (
	"fmt"
	"github.com/buraksezer/consistent"
	"log"
	"testing"
)

func Test_PartitionMapping(t *testing.T) {
	// check across multiple runs, the owner and parition are stable exactly the same

	// create baseline
	members := []consistent.Member{}
	for i := 0; i < 108; i++ {
		member := Member(fmt.Sprintf("node%d.olricmq", i))
		members = append(members, member)
	}

	pCount := 6113

	cfg := consistent.Config{
		PartitionCount:    pCount,
		ReplicationFactor: 40,
		Load:              1.0000001,
		Hasher:            hasher{},
	}

	c := consistent.New(members, cfg)
	mapping := make(map[int]string)
	for p := 0; p < pCount; p++ {
		m := c.GetPartitionOwner(p)
		mapping[p] = m.String()
	}

	// try 100 times
	for i := 0; i < 100; i++ {
		newM := []consistent.Member{}
		for j := 0; j < 108; j++ {
			member := Member(fmt.Sprintf("node%d.olricmq", j))
			newM = append(newM, member)
		}

		newC := consistent.New(newM, cfg)

		for p := 0; p < pCount; p++ {
			m := newC.GetPartitionOwner(p)
			nm := mapping[p]
			if nm != m.String() {
				t.Errorf("not consistent")
			}
		}
		log.Printf("%d all match\n", i)
	}
}

func Test_MisMemberMapping(t *testing.T) {
	// new with []100 and remove 1, if this is the same as new with []99, pick a middle one to remove

	// create baseline
	members := []consistent.Member{}
	for i := 0; i < 108; i++ {
		if i == 5 || i == 13 || i == 98 {
			continue
		}

		member := Member(fmt.Sprintf("node%d.olricmq", i))
		members = append(members, member)
	}

	pCount := 6113

	cfg := consistent.Config{
		PartitionCount:    pCount,
		ReplicationFactor: 40,
		Load:              1.0000001,
		Hasher:            hasher{},
	}

	c := consistent.New(members, cfg)
	mapping := make(map[int]string)
	for p := 0; p < pCount; p++ {
		m := c.GetPartitionOwner(p)
		mapping[p] = m.String()
	}

	// create nw one
	newM := []consistent.Member{}
	for j := 0; j < 108; j++ {
		member := Member(fmt.Sprintf("node%d.olricmq", j))
		newM = append(newM, member)
	}
	newC := consistent.New(newM, cfg)
	newC.Remove("node13.olricmq")
	newC.Remove("node98.olricmq")
	newC.Remove("node5.olricmq")

	for p := 0; p < pCount; p++ {
		m := newC.GetPartitionOwner(p)
		nm := mapping[p]
		if nm != m.String() {
			t.Errorf("not consistent")
		}
	}
}
