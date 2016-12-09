package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view;

import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore.Nfs41DatastoreSpec.Nfs41AuthenticationMode;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.NewDatastoreMessages;
import com.vmware.vsphere.client.test.i18n.I18n;

/**
 * WizardNavigator implementation for:
 *
 * New datastore wizard -> Configure Kerberos authentication
 */
public class ConfigureKerberosAuthenticationPage extends WizardNavigator {

   protected final static NewDatastoreMessages localizedMessages = I18n
         .get(NewDatastoreMessages.class);

   public static final String KERBEROS_SELECTION_HOLDER_COMPONENT = "kerberosUiPane";

   private static final String enableKerberosCheckboxSelector = String.format(
         "kerberosUiPane/label=%s",
         localizedMessages.enableDisableKrbAuthentication());

   private static final String krb5AuthenticationOptionSelector = String
         .format("kerberosUiPane/label=%s", localizedMessages.useKrb5());

   private static final String krb5iAuthenticationOptionSelector = String
         .format("kerberosUiPane/label=%s", localizedMessages.useKrb5i());

   /**
    * Interface for a nfs authentication mode transition
    */
   private static class Nfs41AuthenticationModeStateTransition {

      private final boolean isEnabled;

      public Nfs41AuthenticationModeStateTransition(boolean isEnabled) {
         this.isEnabled = isEnabled;
      }

      void transitToState() {
         if (UI.component.property.getBoolean(Property.VALUE,
               enableKerberosCheckboxSelector) != this.isEnabled) {
            UI.component.click(enableKerberosCheckboxSelector);
         }
      }
   }

   /**
    * Nfs41AuthenticationModeStateTransition implementation for krb5
    * authentication mode
    */
   private static final class KRB5AuthenticationModeStateTransition extends
         Nfs41AuthenticationModeStateTransition {

      public KRB5AuthenticationModeStateTransition() {
         super(true);
      }

      @Override
      void transitToState() {
         super.transitToState();
         UI.component.click(krb5AuthenticationOptionSelector);
      }

   }

   /**
    * Nfs41AuthenticationModeStateTransition implementation for krb5i
    * authentation mode
    */
   private static final class KRB5IAuthenticationModeStateTransition extends
         Nfs41AuthenticationModeStateTransition {

      public KRB5IAuthenticationModeStateTransition() {
         super(true);
      }

      @Override
      void transitToState() {
         super.transitToState();
         UI.component.click(krb5iAuthenticationOptionSelector);
      }

   }

   private static enum Nfs41AuthenticationModeInteractor {
      DISABLED(Nfs41AuthenticationMode.DISABLED,
            new Nfs41AuthenticationModeStateTransition(false)), KRB5(
            Nfs41AuthenticationMode.KRB5,
            new KRB5AuthenticationModeStateTransition()), KRB5I(
            Nfs41AuthenticationMode.KRB5I,
            new KRB5IAuthenticationModeStateTransition());

      private final Nfs41AuthenticationMode authenticationMode;
      private final Nfs41AuthenticationModeStateTransition stateTransition;

      Nfs41AuthenticationModeInteractor(
            Nfs41AuthenticationMode authenticationMode,
            Nfs41AuthenticationModeStateTransition stateTransition) {
         this.authenticationMode = authenticationMode;
         this.stateTransition = stateTransition;
      }

      public static Nfs41AuthenticationModeInteractor fromNfs41AuthenticationMode(
            Nfs41AuthenticationMode mode) {
         for (Nfs41AuthenticationModeInteractor translation : values()) {
            if (translation.authenticationMode == mode) {
               return translation;
            }
         }

         throw new RuntimeException(
               "Unknown translations for NFS authentication mode of " + mode);
      }
   }

   /**
    * Set the Kerberos authentication mode
    *
    * @param autenticationMode
    *           the desired authentication mode for Kerberos
    */
   public void setAuthenticationMode(Nfs41AuthenticationMode autenticationMode) {
      Nfs41AuthenticationModeStateTransition nfsStateTransition = Nfs41AuthenticationModeInteractor
            .fromNfs41AuthenticationMode(autenticationMode).stateTransition;

      nfsStateTransition.transitToState();
   }

   /**
    * Get the info message displayed for usage of Kerberos-based authentication
    *
    * @return
    */
   public String getUseKerberosInfo() {
      return UI.component.property.get(Property.TEXT,
            "_NfsKerberosSettingsPage_Text2");
   }

   /**
    * Checks if the enable Kerberos-based authentication checkbox is enabled
    *
    * @return
    */
   public boolean isEnableKerberosCheckboxEnabled() {
      boolean isKerberosSelectionParentEnabled = UI.component.property
            .getBoolean(Property.ENABLED, KERBEROS_SELECTION_HOLDER_COMPONENT);

      return isKerberosSelectionParentEnabled
            && UI.component.property.getBoolean(Property.ENABLED,
                  enableKerberosCheckboxSelector);
   }

   /**
    * Checks if the use kerberos for authentication only (krb5) radio option is
    * enabled
    *
    * @return
    */
   public boolean isUseKrb5Enabled() {
      boolean isKerberosSelectionParentEnabled = UI.component.property
            .getBoolean(Property.ENABLED, KERBEROS_SELECTION_HOLDER_COMPONENT);
      return isKerberosSelectionParentEnabled
            && UI.component.property.getBoolean(Property.ENABLED,
                  krb5AuthenticationOptionSelector);
   }

   /**
    * Checks if the use kerberos for authentication and data integrity (krb5i)
    * is enabled
    *
    * @return
    */
   public boolean isUseKrb5iEnabled() {
      boolean isKerberosSelectionParentEnabled = UI.component.property
            .getBoolean(Property.ENABLED, KERBEROS_SELECTION_HOLDER_COMPONENT);
      return isKerberosSelectionParentEnabled
            && UI.component.property.getBoolean(Property.ENABLED,
                  krb5iAuthenticationOptionSelector);
   }

}
