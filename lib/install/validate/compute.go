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

package validate

import (
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"context"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

func (v *Validator) compute(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// Compute

	// compute resource looks like <toplevel>/<sub/path>
	// this should map to /datacenter-name/host/<toplevel>/Resources/<sub/path>
	// we need to validate that <toplevel> exists and then that the combined path exists.

	pool, err := v.ResourcePoolHelper(ctx, input.ComputeResourcePath)
	v.NoteIssue(err)
	if pool == nil {
		return
	}

	// stash the pool for later use
	v.ResourcePoolPath = pool.InventoryPath

	// some hoops for while we're still using session package
	v.Session.Pool = pool
	v.Session.PoolPath = pool.InventoryPath
	v.Session.ClusterPath = v.inventoryPathToCluster(pool.InventoryPath)

	clusters, err := v.Session.Finder.ComputeResourceList(v.Context, v.Session.ClusterPath)
	if err != nil {
		log.Errorf("Unable to acquire reference to cluster %q: %s", path.Base(v.Session.ClusterPath), err)
		v.NoteIssue(err)
		return
	}

	if len(clusters) != 1 {
		err := fmt.Errorf("Unable to acquire unambiguous reference to cluster %q", path.Base(v.Session.ClusterPath))
		log.Error(err)
		v.NoteIssue(err)
	}

	v.Session.Cluster = clusters[0]

	// TODO: for vApp creation assert that the name doesn't exist
	// TODO: for RP creation assert whatever we decide about the pool - most likely that it's empty
}

func (v *Validator) ResourcePoolHelper(ctx context.Context, path string) (*object.ResourcePool, error) {
	defer trace.End(trace.Begin(path))

	// if compute-resource is unspecified is there a default
	if path == "" || path == "/" {
		if v.Session.Pool != nil {
			log.Debugf("Using default resource pool for compute resource: %q", v.Session.Pool.InventoryPath)
			return v.Session.Pool, nil
		}

		// if no path specified and no default available the show all
		v.suggestComputeResource("*")
		return nil, errors.New("No unambiguous default compute resource available: --compute-resource must be specified")
	}

	ipath := v.computePathToInventoryPath(path)
	log.Debugf("Converted original path %q to %q", path, ipath)

	// first try the path directly without any processing
	pools, err := v.Session.Finder.ResourcePoolList(ctx, path)
	if err != nil {
		log.Debugf("Failed to look up compute resource as absolute path %q: %s", path, err)
		if _, ok := err.(*find.NotFoundError); !ok {
			// we return err directly here so we can check the type
			return nil, err
		}

		// if it starts with datacenter then we know it's absolute and invalid
		if strings.HasPrefix(path, "/"+v.Session.DatacenterPath) {
			v.suggestComputeResource(path)
			return nil, err
		}
	}

	if len(pools) == 0 {
		// assume it's a cluster specifier - that's the formal case, e.g. /cluster/resource/pool
		// not /cluster/Resources/resource/pool
		// everything from now on will use this assumption

		pools, err = v.Session.Finder.ResourcePoolList(ctx, ipath)
		if err != nil {
			log.Debugf("failed to look up compute resource as cluster path %q: %s", ipath, err)
			if _, ok := err.(*find.NotFoundError); !ok {
				// we return err directly here so we can check the type
				return nil, err
			}
		}
	}

	if len(pools) == 1 {
		log.Debugf("Selected compute resource %q", pools[0].InventoryPath)
		return pools[0], nil
	}

	// both cases we want to suggest options
	v.suggestComputeResource(ipath)

	if len(pools) == 0 {
		log.Debugf("no such compute resource %q", path)
		// we return err directly here so we can check the type
		return nil, err
	}

	// TODO: error about required disabmiguation and list entries in nets
	return nil, errors.New("ambiguous compute resource " + path)
}

func (v *Validator) suggestComputeResource(path string) {
	defer trace.End(trace.Begin(path))

	log.Infof("Suggesting valid values for --compute-resource based on %q", path)

	// allow us to work on inventory paths
	path = v.computePathToInventoryPath(path)

	var matches []string
	for matches = nil; matches == nil; matches = v.findValidPool(path) {
		// back up the path until we find a pool
		newpath := filepath.Dir(path)
		if newpath == "." {
			// filepath.Dir returns . which has no meaning for us
			newpath = "/"
		}
		if newpath == path {
			break
		}
		path = newpath
	}

	if matches == nil {
		// Backing all the way up didn't help
		log.Info("Failed to find resource pool in the provided path, showing all top level resource pools.")
		matches = v.findValidPool("*")
	}

	if matches != nil {
		// we've collected recommendations - displayname
		log.Info("Suggested values for --compute-resource:")
		for _, p := range matches {
			log.Infof("  %q", v.inventoryPathToComputePath(p))
		}
		return
	}

	log.Info("No resource pools found")
}

func (v *Validator) findValidPool(path string) []string {
	defer trace.End(trace.Begin(path))

	// list pools in path
	matches := v.listResourcePools(path)
	if matches != nil {
		sort.Strings(matches)
		return matches
	}

	// no pools in path, but if path is cluster, list pools in cluster
	clusters, err := v.Session.Finder.ComputeResourceList(v.Context, path)
	if len(clusters) == 0 {
		// not a cluster
		log.Debugf("Path %q does not identify a cluster (or clusters) or the list could not be obtained: %s", path, err)
		return nil
	}

	if len(clusters) > 1 {
		log.Debugf("Suggesting clusters as there are multiple matches")
		matches = make([]string, len(clusters))
		for i, c := range clusters {
			matches[i] = c.InventoryPath
		}
		sort.Strings(matches)
		return matches
	}

	log.Debugf("Suggesting pools for cluster %q", clusters[0].Name())
	matches = v.listResourcePools(fmt.Sprintf("%s/Resources/*", clusters[0].InventoryPath))
	if matches == nil {
		// no child pools so recommend cluster directly
		return []string{clusters[0].InventoryPath}
	}

	return matches
}

func (v *Validator) listResourcePools(path string) []string {
	defer trace.End(trace.Begin(path))

	pools, err := v.Session.Finder.ResourcePoolList(v.Context, path+"/*")
	if err != nil {
		log.Debugf("Unable to list pools for %q: %s", path, err)
		return nil
	}

	if len(pools) == 0 {
		return nil
	}

	matches := make([]string, len(pools))
	for i, p := range pools {
		matches[i] = p.InventoryPath
	}

	return matches
}

func (v *Validator) GetResourcePool(input *data.Data) (*object.ResourcePool, error) {
	defer trace.End(trace.Begin(""))

	return v.ResourcePoolHelper(v.Context, input.ComputeResourcePath)
}

func (v *Validator) ValidateCompute(ctx context.Context, input *data.Data, computeRequired bool) (*config.VirtualContainerHostConfigSpec, error) {
	defer trace.End(trace.Begin(""))
	conf := &config.VirtualContainerHostConfigSpec{}

	if input.ComputeResourcePath == "" && !computeRequired {
		return conf, nil
	}

	log.Infof("Validating compute resource")
	v.compute(ctx, input, conf)
	return conf, v.ListIssues()
}
