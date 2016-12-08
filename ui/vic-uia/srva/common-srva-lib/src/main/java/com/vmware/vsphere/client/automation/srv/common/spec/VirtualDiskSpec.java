/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vim.binding.vim.vm.device.VirtualDisk;

/**
 * A data model for a virtual disk device intended for use in the UI and
 * providing a (partial) view over the data modeled in {@link VirtualDisk}.
 */
public class VirtualDiskSpec extends BaseSpec {
   public enum ProvisioningType {
      THIN("Thin provision"), THICK_LAZY_ZEROED("Thick provision lazy zeroed"), THICK_EAGER_ZEROED(
            "Thick provision eager zeroed");
      private String name;

      ProvisioningType(String name) {
         this.name = name;
      }

      @Override
      public String toString() {
         return this.name;
      }
   }

   public enum Sharing {
      UNSPECIFIED("Unspecified"), NO_SHARING("No Sharing"), MULTI_WRITER(
            "Multi-writer");
      private String name;

      Sharing(String name) {
         this.name = name;
      }

      @Override
      public String toString() {
         return this.name;
      }
   }

   /**
    * Name of the virtual disk plus the 'vmdk' extension.
    */
   public DataProperty<String> name;

   /**
    * Path of the folder where this virtual disk is stored on its datastore
    * (without disk name). This path string must not start or end with '/'. An
    * empty string ('') is allowed value. By default, this path points to a
    * folder bearing the name of the vm in which the disk has been created at
    * first.
    */
   public DataProperty<String> path;

   /**
    * Parent datastore of this virtual disk.
    */
   public DataProperty<DatastoreSpec> parentDatastore;

   /**
    * Size in MB.
    */
   public DataProperty<Long> size;

   /**
    * One of {@link VirtualDiskSpec.ProvisioningType}.
    */
   public DataProperty<ProvisioningType> provisioning;

   /**
    * One of {@link VirtualDiskSpec.Sharing}
    */
   public DataProperty<Sharing> sharing;

   /**
    * Returns the absolute path of the virtual disk on its datastore. Examples:
    *
    * <pre>
    * [datastore1] vmName/vmName_1.vmdk
    * [datastore1] vmName/folder1/folder2/vmName_1.vmdk
    * </pre>
    *
    * @return
    */
   public String getAbsolutePath() {
      return String.format("[%s] %s/%s", parentDatastore.get().name.get(),
            path.get(), name.get());
   }
}
