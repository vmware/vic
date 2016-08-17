// Copyright 2016 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package index

import (
	"container/list"
	"fmt"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type Node interface {
	// Returns the string representation (usually the ID) of the node
	String() string

	// Deep copy of the node
	Copy() Node
}

type node struct {
	Node
	parent   *node
	children []*node
	mask     uint32
}

func (n *node) addChild(child *node) {
	n.children = append(n.children, child)
}

type Index struct {
	root        *node
	lookupTable map[string]*node
	l           sync.Mutex
}

func NewIndex() *Index {
	i := &Index{
		lookupTable: make(map[string]*node),
	}

	return i
}

// Insert inserts a copy of the given node to the tree under the given parent.
func (i *Index) Insert(parent string, n Node) error {
	i.l.Lock()
	defer i.l.Unlock()

	_, ok := i.lookupTable[n.String()]
	if ok {
		return fmt.Errorf("node %s already exists in index", n.String())
	}

	newNode := &node{
		Node: n.Copy(),
	}

	if parent == n.String() {
		if i.root != nil {
			return fmt.Errorf("node cannot point to self unless it's root")
		}

		// set root
		i.root = newNode

	} else {
		p, ok := i.lookupTable[parent]
		if !ok {
			return fmt.Errorf("Can't find parent %s", parent)
		}
		newNode.parent = p
		p.addChild(newNode)
	}

	i.lookupTable[n.String()] = newNode
	return nil
}

// Get returns a Copy of the named node.
func (i *Index) Get(nodeId string) (Node, error) {
	i.l.Lock()
	defer i.l.Unlock()

	n, ok := i.lookupTable[nodeId]
	if !ok {
		return nil, fmt.Errorf("Node %s not found", nodeId)
	}

	return n.Copy(), nil
}

// Delete deletes a leaf node
func (i *Index) Delete(nodeId string) (Node, error) {
	i.l.Lock()
	defer i.l.Unlock()

	n, ok := i.lookupTable[nodeId]
	if !ok {
		return nil, fmt.Errorf("Node %s not found", nodeId)
	}

	if len(n.children) != 0 {
		return nil, fmt.Errorf("Node %s has children %#q", nodeId, n.children)
	}

	// remove the reference to the node from its parent
	var deleted bool
	err := i.bfsworker(i.root, func(needle *node) (iterflag, error) {
		for idx, child := range needle.children {
			if child.String() == nodeId {
				// remove the child
				needle.children = append(needle.children[:idx], needle.children[idx+1:]...)
				log.Debugf("Removing %s from parent (%s children : %#q)", nodeId, needle.String(), needle.children)
				deleted = true
				return STOP, nil
			}
		}

		// continue iterating
		return NOOP, nil
	})

	if err != nil {
		return nil, err
	}

	if !deleted {
		err = fmt.Errorf("%s not found in tree", nodeId)
		log.Errorf("%s", err)
		return nil, err
	}

	delete(i.lookupTable, nodeId)

	return n.Node, nil
}

type iterflag int

const (
	NOOP iterflag = iota
	STOP
)

type visitor func(Node) (iterflag, error)

func (i *Index) bfs(root *node, visitFunc visitor) error {
	i.l.Lock()
	defer i.l.Unlock()

	// XXX Look into parallelizing this without breaking API boundaries.
	return i.bfsworker(root, func(n *node) (iterflag, error) { return visitFunc(n.Node) })
}

func (i *Index) bfsworker(root *node, visitFunc func(*node) (iterflag, error)) error {

	queue := list.New()
	queue.PushBack(root)

	for queue.Len() > 0 {
		n := queue.Remove(queue.Front()).(*node)

		flag, err := visitFunc(n)
		if err != nil {
			return err
		}

		if flag == STOP {
			return nil
		}

		for _, child := range n.children {
			queue.PushBack(child)
		}
	}

	return nil
}
