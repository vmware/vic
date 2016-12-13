/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.vm.clone.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.CustomizeHwVmSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreClusterSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VmSpec;

/**
 * Spec for Clone VM
 *
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.clone.spec.CloneVmSpec}
 */
@Deprecated
public class CloneVmSpec extends VmSpec {

   /**
    * Destination compute resource for the VM
    */
   public DataProperty<ManagedEntitySpec> targetComputeResource;

   /**
    * Datastore cluster on which to place the VM
    */
   public DataProperty<DatastoreClusterSpec> targetDatastoreCluster;

   /**
    * Clone VM > Select clone options > Customize the operating system
    */
   public DataProperty<Boolean> customizeGos;

   /**
    * Clone VM > Select clone options > Customize this virtual machine's
    * hardware (Experimental)
    */
   public DataProperty<Boolean> customizeHw;

   /**
    * Clone VM > Select clone options > Power on virtual machine after creation
    */
   public DataProperty<Boolean> powerOnVm;


   /**
    * Contains all customize vm hardware data
    */
   public DataProperty<CustomizeHwVmSpec> customizeHwVmSpec;

}