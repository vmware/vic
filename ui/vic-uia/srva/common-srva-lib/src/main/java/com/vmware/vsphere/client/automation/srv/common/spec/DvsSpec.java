/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Container class for distributed virtual switch properties.
 * The parent of DVS should be datacenter.
 */
public class DvsSpec extends ManagedEntitySpec {

   public enum NetworkResourceControlVersion {
      VERSION_1("version1"), VERSION_2("version2"), VERSION_3("version3");

      private String value;

      public String getValue() {
         return value;
      }

      private NetworkResourceControlVersion(String value) {
         this.value = value;
      }
   }

   /**
    * The NIOC version: version2, version3, etc.
    */
   public DataProperty<NetworkResourceControlVersion> networkResourceControlVersion;

   /**
    * Host to be attached to the DVS.
    */
   public DataProperty<HostSpec> host;

   /**
    * Number of dvs uplinks
    */
   public DataProperty<String> uplinksNum;

   /**
    * Version of the DVS. If not assigned will be used default for the release.
    */
   public DataProperty<String> dvsVersion;
}
