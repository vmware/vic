/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb.spec;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Physical NFS storage provider spec
 */
public class PhysicalNfsStorageProvisionerSpec extends EntitySpec{

   /**
    * storage user
    */
   public DataProperty<String> username;

   /**
    * storage password
    */
   public DataProperty<String> password;

   /**
    * Ip of the physical nfs storage
    */
   public DataProperty<String> ip;

   /**
    * Folder on physical nfs storage that is used as a
    * storage
    */
   public DataProperty<String> folder;

   /**
    * The name of the storage
    */
   public DataProperty<String> name;
}
