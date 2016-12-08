package com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore;

import org.apache.commons.lang.ArrayUtils;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vim.binding.vim.host.NasVolume;
import com.vmware.vim.binding.vim.host.NasVolume.SecurityType;

/**
 * Nfs datastore spec implementation for NFS 4.1 datastore
 */
public class Nfs41DatastoreSpec extends NfsDatastoreSpec {

   /**
    * Supported authentication mode for the nfs 4.1 datastore
    */
   public static enum Nfs41AuthenticationMode {
      DISABLED(localizedMessages.datastoreAutheticationModeDisabled(),
            NasVolume.SecurityType.AUTH_SYS), KRB5(localizedMessages
            .datastoreAutheticationModeKrb(), NasVolume.SecurityType.SEC_KRB5), KRB5I(
            localizedMessages.datastoreAutheticationModeKrbi(),
            NasVolume.SecurityType.SEC_KRB5I);
      /**
       * Localized string displayed to users
       */
      public final String localizedDisplayName;
      /**
       * The value of {@link SecurityType} associated with the current
       * {@link Nfs41AuthenticationMode}
       */
      public final NasVolume.SecurityType apiAuthentication;

      Nfs41AuthenticationMode(String localizedDisplayName,
            NasVolume.SecurityType apiAuthentication) {
         this.localizedDisplayName = localizedDisplayName;
         this.apiAuthentication = apiAuthentication;
      }

      /**
       * Retrieve {@link Nfs41AuthenticationMode} associated with specific
       * {@link SecurityType}
       *
       * @param desiredSecurityType
       * @return
       */
      public static Nfs41AuthenticationMode fromApiAuthenticationMode(
            SecurityType desiredSecurityType) {
         for (Nfs41AuthenticationMode nfsAuthenticationMode : values()) {
            if (nfsAuthenticationMode.apiAuthentication
                  .equals(desiredSecurityType)) {
               return nfsAuthenticationMode;
            }
         }

         throw new RuntimeException(
               String.format(
                     "The desired api security type of '%s' was not found in enumerated values '%s'",
                     desiredSecurityType, ArrayUtils.toString(values())));
      }
   }

   /**
    * The authentication mode of the nfs datastore
    */
   public DataProperty<Nfs41AuthenticationMode> authenticationMode;

   /**
    * The addresses of the remote nfs server
    */
   public DataProperty<String[]> remoteHostAdresses;

   /**
    * Initializes new instance of Nfs41DatastoreSpec
    */
   public Nfs41DatastoreSpec() {
      super(NfsVersion.NFS41);
   }
}
