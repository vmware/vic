package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.step;

import com.vmware.client.automation.assertions.EqualsAssertion;
import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.storage.lib.wizard.core.step.WizardStep;
import com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.view.ConfigureKerberosAuthenticationPage;

/**
 * {@link WizardStep} implementation for verification of Configure Kerberos
 * authentication page
 */
public class VerifyKerberosAuthenticationStep extends WizardStep {

   /**
    * {@link BaseSpec} implementation for the state of the Configure Kerberos
    * authetication page of the new datastore wizard
    */
   public static class KerborosAuthenticationStateSpec extends BaseSpec {

      /**
       * The use kerberos authentication info message
       */
      public DataProperty<String> infoMessage;

      /**
       * The state of the Enable Kerberos-based authentication checkbox
       */
      public DataProperty<Boolean> isEnableKerberosAuthenticationEnabled;

      /**
       * The state of "Use Kerberos for authentication only (krb5)" option
       */
      public DataProperty<Boolean> isKrb5Enabled;

      /**
       * The state of
       * "Use Kerberos for authentication and data integrity (krb5i)" option
       */
      public DataProperty<Boolean> isKrb5iEnabled;
   }

   @UsesSpec
   private KerborosAuthenticationStateSpec verificationSpec;

   @Override
   protected void executeWizardOperation(WizardNavigator wizardNavigator)
         throws Exception {
      ConfigureKerberosAuthenticationPage page = new ConfigureKerberosAuthenticationPage();

      verifySafely(new EqualsAssertion(page.getUseKerberosInfo(),
            verificationSpec.infoMessage.get(), "Use keberos info"));

      verifySafely(new EqualsAssertion(page.isEnableKerberosCheckboxEnabled(),
            verificationSpec.isEnableKerberosAuthenticationEnabled.get(),
            "Enable Kerberos-based authentication checkbox state"));

      verifySafely(new EqualsAssertion(page.isUseKrb5Enabled(),
            verificationSpec.isKrb5Enabled.get(),
            "Use Kerberos for authetication only radio button state"));

      verifySafely(new EqualsAssertion(page.isUseKrb5Enabled(),
            verificationSpec.isKrb5Enabled.get(),
            "Use Kerberos for authetication and data integrity radio button state"));
   }
}
