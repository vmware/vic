/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.verification.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Spec used to validate specific page from a wizard
 */
public class WizardPageVerificationSpec extends BaseSpec {

   /**
    * Property that stores left navigation page title.
    */
   public DataProperty<String> leftNavPageTitle;

   /**
    * Property that stores wizard page header that is located
    * on top of the page content area.
    */
   public DataProperty<String> pageHeader;

   /**
    * Property that stores wizard page header description that is located
    * right below the page header.
    */
   public DataProperty<String> pageHeaderDesription;

   /**
    * Property that stores the expected state of back button(enabled/disabled).
    */
   public DataProperty<Boolean> backButtonEnabled;

   /**
    * Property that stores the expected state of next button(enabled/disabled).
    */
   public DataProperty<Boolean> nextButtonEnabled;

   /**
    * Property that stores the expected state of finish button(enabled/disabled).
    */
   public DataProperty<Boolean> finishButtonEnabled;
}
