/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.spec;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Base class for Virtual Center Objects specs.
 *
 */
public class ManagedEntitySpec extends EntitySpec {


   @Name
   public DataProperty<String> name;

   public DataProperty<ManagedEntitySpec> parent;

   public enum ReconfigurableConfigSpecVersion {
      OLD("oldConfigVersion"), NEW("newConfigVersion");

      private String value;

      public String getValue() {
         return value;
      }

      private ReconfigurableConfigSpecVersion(String value) {
         this.value = value;
      }
   }

   /**
    * 
    */
   public DataProperty<ReconfigurableConfigSpecVersion> reconfigurableConfigVersion;

   @Override
   public String toString() {
      String parent = "NONE";
      if(this.parent.isAssigned()) {
         if(this.parent.get().name.isAssigned()) {
            parent = this.parent.get().name.get();
         } else {
            parent = "NAME NOT ASSIGNED";
         }
      }
      return String.format(
            "%s: name->%s, parent->%s",
            this.getClass().getSimpleName(),
            this.name.isAssigned() ? this.name.get() : "NONE",
                  parent);
   }
}
