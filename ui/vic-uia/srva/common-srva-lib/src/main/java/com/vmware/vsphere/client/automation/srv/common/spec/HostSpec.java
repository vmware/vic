/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Container class for Vritual Center host properties. Properties necessary
 * for host addition are included.
 *
 */
public class HostSpec extends ManagedEntitySpec {

   /**
    * Admin user of the host.
    */
   public DataProperty<String> userName;

   /**
    * Password of the above user.
    */
   public DataProperty<String> password;

   /**
    * Host port, default is 443.
    */
   public DataProperty<Integer> port;

   /**
    * iSCSI server IP.
    */
   public DataProperty<String> iscsiServerIp;

   /**
    * Number of the free physical NICs of the host
    */
   public DataProperty<Integer> numberFreePnics;
}
