package example

import (
	"fmt"
	"testing"

	"github.com/buraksezer/consistent"
)

// In your distributed system, you probably have a custom data type
// for your cluster members. Just add a String function to implement
// consistent.Member interface.

func Test_Sample(t *testing.T) {
	// Create a new consistent instance
	cfg := consistent.Config{
		PartitionCount:    7,
		ReplicationFactor: 20,
		Load:              1.25,
		Hasher:            hasher{},
	}
	c := consistent.New(nil, cfg)

	// Add some members to the consistent hash table.
	// Add function calculates average load and distributes partitions over members
	node1 := Member("node1.olricmq.com")
	c.Add(node1)

	node2 := Member("node2.olricmq.com")
	c.Add(node2)

	key := []byte("my-key")
	// calculates partition id for the given key
	// partID := hash(key) % partitionCount
	// the partitions is already distributed among members by Add function.
	owner := c.LocateKey(key)
	fmt.Println(owner.String())
	// Prints node2.olricmq.com
}
