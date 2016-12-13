/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Container class for cluster High Availability Admission Control properties.
 */
public class AdmissionControlSpec extends ManagedEntitySpec {

   /**
    * Property that shows the host failures cluster tolerates
    */
   public DataProperty<Integer> hostFailuresTolerates;

   /**
    * Property that shows the defined failover policy
    */
   public DataProperty<AdmissionControlFailoverPolicy> failoverPolicy;

   /**
    * Property that shows the VM resource reduction event threshold in percent
    */
   public DataProperty<Integer> resourceReductionThreshold;
}