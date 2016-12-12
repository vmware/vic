/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Container class for Datastore cluster (Storage POD) properties.
 */
public class DatastoreClusterSpec extends ManagedEntitySpec {

   /**
    * Property that shows whether SDRS is enabled
    */
   public DataProperty<Boolean> sdrsEnabled;

   /**
    * Property that shows whether I/O load balancing is enabled
    */
   public DataProperty<Boolean> ioBalancingEnabled;

   /**
    * Property that shows SDRS automation behavior - fullyAutomated or manual
    */
   public DataProperty<SdrsBehavior> sdrsBehavior;
}
