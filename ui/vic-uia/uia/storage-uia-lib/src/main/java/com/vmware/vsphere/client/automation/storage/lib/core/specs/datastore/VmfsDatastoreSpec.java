package com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore;

import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;

/**
 * Datastore spec implementation for VMFS datastore
 */
public class VmfsDatastoreSpec extends DatastoreSpec {

   /**
    * Initializes new instance of VmfsDatastoreSpec
    */
   public VmfsDatastoreSpec() {
      super.type.set(DatastoreType.VMFS);
   }

}
