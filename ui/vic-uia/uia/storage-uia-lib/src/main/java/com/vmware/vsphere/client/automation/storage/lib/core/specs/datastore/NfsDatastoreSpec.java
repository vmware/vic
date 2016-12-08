package com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore;

import org.apache.commons.lang.ArrayUtils;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vim.binding.vim.host.MountInfo;
import com.vmware.vim.binding.vim.host.FileSystemVolume.FileSystemType;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreType;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * Datastore spec implementation for NFS datastore
 */
public abstract class NfsDatastoreSpec extends DatastoreSpec {

   protected final static DatastoreSpecMessages localizedMessages = I18n
         .get(DatastoreSpecMessages.class);

   /**
    * Enumeration of all available NFS version
    */
   public static enum NfsVersion {
      /**
       * Nfs version 3
       */
      NFS3(FileSystemType.NFS, localizedMessages.nfsVersion3()),

      /**
       * Nfs version 4.1
       */
      NFS41(FileSystemType.NFS41, localizedMessages.nfsVersion41());

      /**
       * Localized string displayed to users
       */
      public final String localizedNfsVersionString;

      /**
       * API enum value for {@link FileSystemType} enumeration associated with
       * the current {@link NfsVersion}
       */
      public final FileSystemType apiFileSystemEnum;

      NfsVersion(FileSystemType apiFileSystemEnum,
            String localizedNfsVersionString) {
         this.localizedNfsVersionString = localizedNfsVersionString;
         this.apiFileSystemEnum = apiFileSystemEnum;
      }
   }

   /**
    * Enumeration of available NFS datastore access modes
    */
   public static enum DatastoreAccessMode {
      READ_ONLY(MountInfo.AccessMode.readOnly, localizedMessages
            .datastoreAccessModeReadOnly()), READ_WRITE(
            MountInfo.AccessMode.readWrite, localizedMessages
                  .datastoreAccessModeReadWrite());
      /**
       * Localized string displayed to users
       */
      public final String localizedDisplayName;

      /**
       * API enum value for {@link MountInfo.AccessMode} enumeration associated
       * with the current {@link DatastoreAccessMode}
       */
      public final MountInfo.AccessMode apiDatastoreAccessMode;

      DatastoreAccessMode(final MountInfo.AccessMode apiDatastoreAccessMode,
            final String localizedDisplayName) {
         this.localizedDisplayName = localizedDisplayName;
         this.apiDatastoreAccessMode = apiDatastoreAccessMode;
      }

      public static final DatastoreAccessMode fromApiAccessMode(
            MountInfo.AccessMode desiredApiMode) {
         for (DatastoreAccessMode currentValue : values()) {
            if (currentValue.apiDatastoreAccessMode.equals(desiredApiMode)) {
               return currentValue;
            }
         }

         throw new RuntimeException(
               String.format(
                     "The desired api state '%s' was not found in enumerated values '%s'",
                     desiredApiMode, ArrayUtils.toString(values())));
      }
   }

   /**
    * The version of the nfs
    */
   public final NfsVersion nfsVersion;

   public DataProperty<DatastoreAccessMode> accessMode;

   /**
    * Initializes new instance of Nfs Datastore spec
    *
    * @param nfsVersion
    *           the nfs version to be associated with the current instance
    */
   protected NfsDatastoreSpec(NfsVersion nfsVersion) {
      this.nfsVersion = nfsVersion;
      super.type.set(DatastoreType.NFS);
   }
}
