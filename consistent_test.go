// Copyright (c) 2018 Burak Sezer
// All rights reserved.
//
// This code is licensed under the MIT License.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files(the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and / or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions :
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package consistent

import (
	"fmt"
	"hash/fnv"
	"testing"
)

func newConfig() *Config {
	return &Config{
		PartitionCount:    23,
		ReplicationFactor: 20,
		LoadFactor:        1.25,
		Hasher:            hasher{},
	}
}

type testMember string

func (tm testMember) Name() string {
	return string(tm)
}

type hasher struct{}

func (hs hasher) Sum64(data []byte) uint64 {
	h := fnv.New64()
	h.Write(data)
	return h.Sum64()
}

func TestConsistentAdd(t *testing.T) {
	cfg := newConfig()
	c := New(nil, cfg)
	members := make(map[string]struct{})
	for i := 0; i < 8; i++ {
		member := testMember(fmt.Sprintf("node%d.olricmq", i))
		members[member.Name()] = struct{}{}
		c.Add(member)
	}
	for member, _ := range members {
		found := false
		for _, mem := range c.GetMembers() {
			if member == mem.Name() {
				found = true
			}
		}
		if !found {
			t.Fatalf("%s could not be found", member)
		}
	}
}

func TestConsistentRemove(t *testing.T) {
	members := []Member{}
	for i := 0; i < 8; i++ {
		member := testMember(fmt.Sprintf("node%d.olricmq", i))
		members = append(members, member)
	}
	cfg := newConfig()
	c := New(members, cfg)
	if len(c.GetMembers()) != len(members) {
		t.Fatalf("inserted member count is different")
	}
	for _, member := range members {
		c.Remove(member.Name())
	}
	if len(c.GetMembers()) != 0 {
		t.Fatalf("member count should be zero")
	}
}

func TestConsistentLoad(t *testing.T) {
	members := []Member{}
	for i := 0; i < 8; i++ {
		member := testMember(fmt.Sprintf("node%d.olricmq", i))
		members = append(members, member)
	}
	cfg := newConfig()
	c := New(members, cfg)
	if len(c.GetMembers()) != len(members) {
		t.Fatalf("inserted member count is different")
	}
	maxLoad := c.AverageLoad()
	for member, load := range c.LoadDistribution() {
		if load > maxLoad {
			t.Fatalf("%s exceeds max load. Its load: %f, max load: %f", member, load, maxLoad)
		}
	}
}

func TestConsistentLocateKey(t *testing.T) {
	cfg := newConfig()
	c := New(nil, cfg)
	key := []byte("OlricMQ")
	res := c.LocateKey(key)
	if res != nil {
		t.Fatalf("This should be nil: %v", res)
	}
	members := make(map[string]struct{})
	for i := 0; i < 8; i++ {
		member := testMember(fmt.Sprintf("node%d.olricmq", i))
		members[member.Name()] = struct{}{}
		c.Add(member)
	}
	res = c.LocateKey(key)
	if res == nil {
		t.Fatalf("This shouldn't be nil: %v", res)
	}
}

func TestConsistentInsufficientMemberCount(t *testing.T) {
	members := []Member{}
	for i := 0; i < 8; i++ {
		member := testMember(fmt.Sprintf("node%d.olricmq", i))
		members = append(members, member)
	}
	cfg := newConfig()
	c := New(members, cfg)
	key := []byte("OlricMQ")
	partID := c.FindPartitionID(key)
	_, err := c.GetPartitionBackups(partID, 30)
	if err != ErrInsufficientMemberCount {
		t.Errorf("Expected ErrInsufficientMemberCount(%v), Got: %v", ErrInsufficientMemberCount, err)
	}
}

func TestConsistentBackupMembers(t *testing.T) {
	members := []Member{}
	for i := 0; i < 8; i++ {
		member := testMember(fmt.Sprintf("node%d.olricmq", i))
		members = append(members, member)
	}
	cfg := newConfig()
	c := New(members, cfg)
	key := []byte("OlricMQ")
	partID := c.FindPartitionID(key)
	backups, err := c.GetPartitionBackups(partID, 2)
	if err != nil {
		t.Errorf("Expected nil, Got: %v", err)
	}
	if len(backups) != 2 {
		t.Errorf("Expected backup count is 2. Got: %d", len(backups))
	}
	owner := c.GetPartitionOwner(partID)
	for _, backup := range backups {
		if backup.Name() == owner.Name() {
			t.Fatalf("Backup is equal the partition owner: %s", owner.Name())
		}
	}
}