# VCH Deployment with a Shared NFS Datastore Fails with an Error About No Single Host Being Able to Access All Datastores #

Deploying a virtual container host (VCH) to a cluster, and specifying a shared NFS datastore as the image store, fails with the error `No single host can access all of the requested datastores.` 

## Problem ##

This error occurs even if all of the hosts in the cluster do appear to have access to the shared NFS datastore.

## Cause ##

VCHs require datastores to be writable. The shared NFS datastore is possibly mounted as read-only.

## Solution ##

To see whether a datastore is writable or read-only, consult  `mountInfo` in the Managed Object Browser (MOB) of the vCenter Server instance to which you are deploying the VCH. 

1. Go to https://<i>vcenter_server_address</i>/mob/.
2. Click **content**.
3. Click **group-xx (Datacenters)** in the `rootFolder` row.
4. Click the managed object reference (MoRef) of your datacenter in the `childEntity` row.
5. Click the MoRef of the shared NFS datastore in the `datastore` row.
6. Click the `DatastoreHostMount` link in the `host` row.
7. Click **`mountInfo`** and check the `accessMode` value.
8. If the `accessMode` value is `readOnly`, unmount the datastore from vCenter Server and remount it with `readWrite` permissions.