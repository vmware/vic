package com.vmware.vsphere.client.automation.srv.common.spec;

/**
 * Represents supported datastore types.
 */
public enum DatastoreType {
   VMFS {
      @Override
      public String toString() {
         return "VMFS";

      }
   },
   NFS {
      @Override
      public String toString() {
         return "NFS";

      }
   },
   VSAN {
      @Override
      public String toString() {
         return "VSAN";

      }
   }
}
