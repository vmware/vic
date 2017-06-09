// Copyright 2016-2017 VMware, Inc. All Rights Reserved.
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
	"context"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/vic/lib/config"
	"github.com/vmware/vic/lib/install/data"
	"github.com/vmware/vic/pkg/errors"
	"github.com/vmware/vic/pkg/trace"
)

func (v *Validator) compute(ctx context.Context, input *data.Data, conf *config.VirtualContainerHostConfigSpec) {
	defer trace.End(trace.Begin(""))

	// ComputeResourcePath should resolve to a ComputeResource, ClusterComputeResource or ResourcePool

	pool, err := v.ResourcePoolHelper(ctx, input.ComputeResourcePath)
	v.NoteIssue(err)
	if pool == nil {
		return
	}

	// TODO: for vApp creation assert that the name doesn't exist
	// TODO: for RP creation assert whatever we decide about the pool - most likely that it's empty
}

func (v *Validator) inventoryPath(obj object.Reference) string {
	elt, err := v.Session.Finder.Element(v.Context, obj.Reference())
	if err != nil {
		log.Warnf("failed to get inventory path for %s: %s", obj.Reference(), err)
		return ""
	}

	return elt.Path
}

// ResourcePoolHelper finds a resource pool from the input compute path and shows
// suggestions if unable to do so when the path is ambiguous.
func (v *Validator) ResourcePoolHelper(ctx context.Context, path string) (*object.ResourcePool, error) {
	defer trace.End(trace.Begin(path))

	// if compute-resource is unspecified is there a default
	if path == "" {
		if v.Session.Pool != nil {
			log.Debugf("Using default resource pool for compute resource: %q", v.Session.Pool.InventoryPath)
			return v.Session.Pool, nil
		}

		// if no path specified and no default available the show all
		v.suggestComputeResource()
		return nil, errors.New("No unambiguous default compute resource available: --compute-resource must be specified")
	}

	pool, err := v.Session.Finder.ResourcePool(ctx, path)
	if err != nil {
		switch err.(type) {
		case *find.NotFoundError:
			// fall through to ComputeResource check
		case *find.MultipleFoundError:
			log.Errorf("Failed to use --compute-resource=%q as resource pool: %s", path, err)
			v.suggestResourcePool(path)
			return nil, err
		default:
			return nil, err
		}
	}

	var compute *object.ComputeResource

	if pool == nil {
		// check if its a ComputeResource or ClusterComputeResource
		compute, err = v.Session.Finder.ComputeResource(ctx, path)
		if err != nil {
			switch err.(type) {
			case *find.NotFoundError, *find.MultipleFoundError:
				v.suggestComputeResource()
			}

			return nil, err
		}

		// Use the default pool
		pool, err = compute.ResourcePool(ctx)
		if err != nil {
			return nil, err
		}
		pool.InventoryPath = v.inventoryPath(pool.Reference())
	} else {
		// TODO: add an object.ResourcePool.Owner method (see compute.ResourcePool.GetCluster)
		var p mo.ResourcePool

		if err = pool.Properties(ctx, pool.Reference(), []string{"owner"}, &p); err != nil {
			log.Errorf("unable to get cluster of resource pool %s: %s", pool.Name(), err)
			return nil, err
		}

		compute = object.NewComputeResource(pool.Client(), p.Owner)
		compute.InventoryPath = v.inventoryPath(compute.Reference())
	}

	// stash the pool for later use
	v.ResourcePoolPath = pool.InventoryPath

	// some hoops for while we're still using session package
	v.Session.Pool = pool
	v.Session.PoolPath = pool.InventoryPath

	v.Session.Cluster = compute
	v.Session.ClusterPath = compute.InventoryPath

	return pool, nil
}

func (v *Validator) ListComputeResource() ([]string, error) {
	compute, err := v.Session.Finder.ComputeResourceList(v.Context, "*")
	if err != nil {
		return nil, fmt.Errorf("unable to list compute resource: %s", err)
	}

	if len(compute) == 0 {
		return nil, nil
	}

	matches := make([]string, len(compute))
	for i, c := range compute {
		matches[i] = c.Name()
	}
	return matches, nil
}

func (v *Validator) suggestComputeResource() {
	defer trace.End(trace.Begin(""))

	compute, err := v.ListComputeResource()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Suggested values for --compute-resource:")
	for _, c := range compute {
		log.Infof("  %q", c)
	}
}

func (v *Validator) ListResourcePool(path string) ([]string, error) {
	pools, err := v.Session.Finder.ResourcePoolList(v.Context, path)
	if err != nil {
		return nil, fmt.Errorf("unable to list resource pool: %s", err)
	}

	if len(pools) == 0 {
		return nil, nil
	}

	matches := make([]string, len(pools))
	for i, p := range pools {
		matches[i] = p.InventoryPath
	}
	return matches, nil
}

func (v *Validator) suggestResourcePool(path string) {
	defer trace.End(trace.Begin(""))

	pools, err := v.ListResourcePool(path)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Suggested resource pool values for --compute-resource:")
	for _, c := range pools {
		p := strings.TrimPrefix(c, v.DatacenterPath+"/host/")
		log.Infof("  %q", p)
	}
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
