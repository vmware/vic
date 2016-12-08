/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.tree;

/**
 * An enumeration for entity types in NGC.
 */
public enum EntityTypes {
   DATACENTER("Datacenter"), FOLDER("Folder"), HOST("HostSystem"), COMPUTE_RESOURCE(
         "ComputeResource"), CLUSTER("ClusterComputeResource"), RESOURCE_POOL(
         "ResourcePool"), VAPP("VirtualApp"), VM("VirtualMachine"), DATASTORE(
         "Datastore"), DVS("VmwareDistributedVirtualSwitch"), DVPORTGROUP(
         "DistributedVirtualPortgroup"), DS_CLUSTER("StoragePod"), NETWORK(
         "Network");
   private String entityType;

   EntityTypes(String value) {
      entityType = value;
   }

   public String getEntityType() {
      return entityType;
   }
}