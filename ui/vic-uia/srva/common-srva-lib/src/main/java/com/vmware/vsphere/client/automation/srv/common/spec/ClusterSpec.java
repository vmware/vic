/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Container class for Virtual Center cluster properties. Properties necessary
 * for cluster creation are included.
 *
 */
public class ClusterSpec extends ManagedEntitySpec {
   /**
    * Property that shows whether DRS is enabled
    */
   public DataProperty<Boolean> drsEnabled;
   /**
    * Property that shows DRS automation behavior - fullyAutomated,
    * partiallyAutomated or manual
    */
   public DataProperty<DrsBehavior> drsBehavior;

   /**
    * Property that shows whether SPBM is enabled
    */
   public DataProperty<Boolean> isSpbmEnabled;

   /**
    * Property that shows whether vSAN is enabled.
    */
   public DataProperty<Boolean> vsanEnabled;

   /**
    * Property that shows whether vSAN automatically claims storage.
    */
   public DataProperty<Boolean> vsanAutoClaimStorage;

   /**
    * Property that shows whether HA is enabled
    */
   public DataProperty<Boolean> vsphereHA;

   /**
    * Property that contains Admission Control settings for HA
    */
   public DataProperty<AdmissionControlSpec> admissionControlSpec;

   /**
    * Property that shows the cluster automation level
    */
   public static enum DrsAutoLevel {
      PARTIAL_AUTO, FULL_AUTO, MANUAL
   }

}
