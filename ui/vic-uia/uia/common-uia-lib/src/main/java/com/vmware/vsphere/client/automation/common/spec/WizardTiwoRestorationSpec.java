/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.components.navigator.WizardNavigator;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Container class for TIWO restoration options for a wizard.
 */
public class WizardTiwoRestorationSpec extends BaseSpec {

   /**
    * Presents the expected view that should have been reached.
    */
   public DataProperty<Class<? extends WizardNavigator>> expectedView;

   /**
    * Expected error message
    */
   public DataProperty<String> expectedErrorMessage;
}
