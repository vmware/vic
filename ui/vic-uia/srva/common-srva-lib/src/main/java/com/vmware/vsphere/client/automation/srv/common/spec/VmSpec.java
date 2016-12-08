/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Container class for VM properties. The parent of a VM should be a
 * <code>HostSpec</code>.
 */
public class VmSpec extends ManagedEntitySpec {

   /**
    * Property that shows VM's guest OS id.
    */
   public DataProperty<String> guestId;

   /**
    * <code>true</code> indicates that this VM has a network adapter.
    */
   public DataProperty<NicSpec> nicList;

   /**
    * Property that specifies VM hardware version it should in format vmx-07
    */
   public DataProperty<String> hardwareVersion;

   /**
    * Property that specifies concrete datastore on which the VM would be
    * placed. If both the 'parent' and 'datastore' properties of the VmSpec are
    * used, it's up to the programmer to watch for consistency between these two
    * properties.
    */
   public DataProperty<DatastoreSpec> datastore;

   /**
    * VM Folder in which the VM will be created. If the property stays
    * unassigned, the default VM folder will be used
    */
   public DataProperty<FolderSpec> vmFolder;

   /**
    * The compute resource - drs cluster, resource pool or vApp. If not set the
    * default resource pool will be used.
    */
   public DataProperty<ManagedEntitySpec> computeResource;

   /**
    * Property that specifies the memory in MB for the VM
    */
   public DataProperty<Long> memoryInMB;

   /**
    * Property that specifies the number of CPUs for the VM
    */
   public DataProperty<Integer> numCPUs;

   /**
    * Property that specifies the fault tolerance role of the VM
    */
   public DataProperty<FaultToleranceVmRoles> ftRole;

   /**
    * Enumeration for Fault Tolerance VM roles. This enumeration is used when
    * working with the API - roles are represented by integers in the backend.
    */
   public enum FaultToleranceVmRoles {
      PRIMARY(1), SECONDARY(2);

      private final Integer ftRole;

      private FaultToleranceVmRoles(Integer role) {
         ftRole = role;
      }

      public Integer getRole() {
         return ftRole;
      }
   }

   /*
    * Property that specifies Vm Storage Policy which would be assigned to the
    * VM.
    */
   public DataProperty<StoragePolicySpec> profile;

   /**
    * Virtual disks of the vm.
    */
   public DataProperty<VirtualDiskSpec> diskSpec;

}
