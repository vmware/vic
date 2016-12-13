package com.vmware.vsphere.client.automation.common.annotations;

import java.util.UUID;

import com.vmware.client.automation.common.spec.EntitySpec;

/**
 * Spec implementation which do overrides the Objetc.hashCode and Object.equals
 * and can be safely compared after a deep clone
 */
public class MockedSpec extends EntitySpec {

   /**
    * GUID of current instance
    */
   private final UUID instanceId = UUID.randomUUID();

   @Override
   public boolean equals(Object obj) {
      if (obj instanceof MockedSpec) {
         return this.instanceId.equals(((MockedSpec) obj).instanceId);
      }

      return super.equals(obj);
   }

   @Override
   public int hashCode() {
      return instanceId.hashCode();
   }
}
