package integration

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/log"
	raftutils "github.com/docker/swarmkit/manager/state/raft/testutils"
	"golang.org/x/net/context"
)

const opsTimeout = 64 * time.Second

// Cluster is representation of cluster - connected nodes.
type testCluster struct {
	ctx        context.Context
	cancel     context.CancelFunc
	api        *dummyAPI
	nodes      map[string]*testNode
	nodesOrder map[string]int
	errs       chan error
	wg         sync.WaitGroup
	counter    int
}

// NewCluster creates new cluster to which nodes can be added.
// AcceptancePolicy is set to most permissive mode on first manager node added.
func newTestCluster() *testCluster {
	ctx, cancel := context.WithCancel(context.Background())
	c := &testCluster{
		ctx:        ctx,
		cancel:     cancel,
		nodes:      make(map[string]*testNode),
		nodesOrder: make(map[string]int),
		errs:       make(chan error, 1024),
	}
	c.api = &dummyAPI{c: c}
	return c
}

// Stop makes best effort to stop all nodes and close connections to them.
func (c *testCluster) Stop() error {
	c.cancel()
	for _, n := range c.nodes {
		if err := n.Stop(); err != nil {
			return err
		}
	}
	c.wg.Wait()
	close(c.errs)
	for err := range c.errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// RandomManager chooses random manager from cluster.
func (c *testCluster) RandomManager() *testNode {
	var managers []*testNode
	for _, n := range c.nodes {
		if n.IsManager() {
			managers = append(managers, n)
		}
	}
	idx := rand.Intn(len(managers))
	return managers[idx]
}

// AddManager adds node with Manager role(both agent and manager).
func (c *testCluster) AddManager() error {
	// first node
	var n *testNode
	if len(c.nodes) == 0 {
		node, err := newTestNode("", "")
		if err != nil {
			return err
		}
		n = node
	} else {
		joinAddr, err := c.RandomManager().node.RemoteAPIAddr()
		if err != nil {
			return err
		}
		clusterInfo, err := c.api.ListClusters(context.Background(), &api.ListClustersRequest{})
		if err != nil {
			return err
		}
		if len(clusterInfo.Clusters) == 0 {
			return fmt.Errorf("joining manager: there is no cluster created in storage")
		}
		node, err := newTestNode(joinAddr, clusterInfo.Clusters[0].RootCA.JoinTokens.Manager)
		if err != nil {
			return err
		}
		n = node
	}

	c.counter++
	ctx := log.WithLogger(c.ctx, log.L.WithField("testnode", c.counter))

	c.wg.Add(1)
	go func() {
		c.errs <- n.node.Start(ctx)
		c.wg.Done()
	}()

	select {
	case <-n.node.Ready():
	case <-time.After(opsTimeout):
		return fmt.Errorf("node did not ready in time")
	}

	c.nodes[n.node.NodeID()] = n
	c.nodesOrder[n.node.NodeID()] = c.counter
	return nil
}

// AddAgent adds node with Agent role(doesn't participate in raft cluster).
func (c *testCluster) AddAgent() error {
	// first node
	var n *testNode
	if len(c.nodes) == 0 {
		return fmt.Errorf("there is no manager nodes")
	}
	joinAddr, err := c.RandomManager().node.RemoteAPIAddr()
	if err != nil {
		return err
	}
	clusterInfo, err := c.api.ListClusters(context.Background(), &api.ListClustersRequest{})
	if err != nil {
		return err
	}
	if len(clusterInfo.Clusters) == 0 {
		return fmt.Errorf("joining agent: there is no cluster created in storage")
	}
	node, err := newTestNode(joinAddr, clusterInfo.Clusters[0].RootCA.JoinTokens.Worker)
	if err != nil {
		return err
	}
	n = node

	c.counter++
	ctx := log.WithLogger(c.ctx, log.L.WithField("testnode", c.counter))

	c.wg.Add(1)
	go func() {
		c.errs <- n.node.Start(ctx)
		c.wg.Done()
	}()

	select {
	case <-n.node.Ready():
	case <-time.After(opsTimeout):
		return fmt.Errorf("node is not ready in time")
	}
	c.nodes[n.node.NodeID()] = n
	c.nodesOrder[n.node.NodeID()] = c.counter
	return nil
}

// CreateService creates dummy service.
func (c *testCluster) CreateService(name string, instances int) (string, error) {
	spec := &api.ServiceSpec{
		Annotations: api.Annotations{Name: name},
		Mode: &api.ServiceSpec_Replicated{
			Replicated: &api.ReplicatedService{
				Replicas: uint64(instances),
			},
		},
		Task: api.TaskSpec{
			Runtime: &api.TaskSpec_Container{
				Container: &api.ContainerSpec{Image: "alpine", Command: []string{"sh"}},
			},
		},
	}

	resp, err := c.api.CreateService(context.Background(), &api.CreateServiceRequest{Spec: spec})
	if err != nil {
		return "", err
	}
	return resp.Service.ID, nil
}

// Leader returns TestNode for cluster leader.
func (c *testCluster) Leader() (*testNode, error) {
	resp, err := c.api.ListNodes(context.Background(), &api.ListNodesRequest{
		Filters: &api.ListNodesRequest_Filters{
			Roles: []api.NodeRole{api.NodeRoleManager},
		},
	})
	if err != nil {
		return nil, err
	}
	for _, n := range resp.Nodes {
		if n.ManagerStatus.Leader {
			tn, ok := c.nodes[n.ID]
			if !ok {
				return nil, fmt.Errorf("leader id is %s, but it isn't found in test cluster object", n.ID)
			}
			return tn, nil
		}
	}
	return nil, fmt.Errorf("cluster leader is not found in storage")
}

// RemoveNode removes node entirely. It tries to demote managers.
func (c *testCluster) RemoveNode(id string, graceful bool) error {
	node, ok := c.nodes[id]
	if !ok {
		return fmt.Errorf("remove node: node %s not found", id)
	}
	// demote before removal
	if node.IsManager() {
		if err := c.SetNodeRole(id, api.NodeRoleWorker); err != nil {
			return fmt.Errorf("demote manager: %v", err)
		}

	}
	if err := node.Stop(); err != nil {
		return err
	}
	delete(c.nodes, id)
	if graceful {
		if err := raftutils.PollFuncWithTimeout(nil, func() error {
			resp, err := c.api.GetNode(context.Background(), &api.GetNodeRequest{NodeID: id})
			if err != nil {
				return fmt.Errorf("get node: %v", err)
			}
			if resp.Node.Status.State != api.NodeStatus_DOWN {
				return fmt.Errorf("node %s is still not down", id)
			}
			return nil
		}, opsTimeout); err != nil {
			return err
		}
	}
	if _, err := c.api.RemoveNode(context.Background(), &api.RemoveNodeRequest{NodeID: id, Force: !graceful}); err != nil {
		return fmt.Errorf("remove node: %v", err)
	}
	return nil
}

// SetNodeRole sets role for node through control api.
func (c *testCluster) SetNodeRole(id string, role api.NodeRole) error {
	node, ok := c.nodes[id]
	if !ok {
		return fmt.Errorf("set node role: node %s not found", id)
	}
	if node.IsManager() && role == api.NodeRoleManager {
		return fmt.Errorf("node is already manager")
	}
	if !node.IsManager() && role == api.NodeRoleWorker {
		return fmt.Errorf("node is already worker")
	}

	var initialTimeout time.Duration
	// version might change between get and update, so retry
	for i := 0; i < 5; i++ {
		time.Sleep(initialTimeout)
		initialTimeout += 500 * time.Millisecond
		resp, err := c.api.GetNode(context.Background(), &api.GetNodeRequest{NodeID: id})
		if err != nil {
			return err
		}
		spec := resp.Node.Spec.Copy()
		spec.Role = role
		if _, err := c.api.UpdateNode(context.Background(), &api.UpdateNodeRequest{
			NodeID:      id,
			Spec:        spec,
			NodeVersion: &resp.Node.Meta.Version,
		}); err != nil {
			// there possible problems on calling update node because redirecting
			// node or leader might want to shut down
			if grpc.ErrorDesc(err) == "update out of sequence" {
				continue
			}
			return err
		}
		if role == api.NodeRoleManager {
			// wait to become manager
			return raftutils.PollFuncWithTimeout(nil, func() error {
				if !node.IsManager() {
					return fmt.Errorf("node is still not a manager")
				}
				return nil
			}, opsTimeout)
		}
		// wait to become worker
		return raftutils.PollFuncWithTimeout(nil, func() error {
			if node.IsManager() {
				return fmt.Errorf("node is still not a worker")
			}
			return nil
		}, opsTimeout)
	}
	return fmt.Errorf("set role %s for node %s, got sequence error 5 times", role, id)
}

// Starts a node from a stopped state
func (c *testCluster) StartNode(id string) error {
	n, ok := c.nodes[id]
	if !ok {
		return fmt.Errorf("set node role: node %s not found", id)
	}

	ctx := log.WithLogger(c.ctx, log.L.WithField("testnode", c.nodesOrder[id]))

	c.wg.Add(1)
	go func() {
		c.errs <- n.node.Start(ctx)
		c.wg.Done()
	}()

	select {
	case <-n.node.Ready():
	case <-time.After(opsTimeout):
		return fmt.Errorf("node did not ready in time")
	}
	if n.node.NodeID() != id {
		return fmt.Errorf("restarted node does not have have the same ID")
	}
	return nil
}
