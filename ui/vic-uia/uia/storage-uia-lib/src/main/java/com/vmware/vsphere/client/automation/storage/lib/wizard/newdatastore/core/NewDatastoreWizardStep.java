/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.storage.lib.wizard.newdatastore.core;

import com.vmware.vsphere.client.automation.srv.common.spec.DatastoreSpec;
import com.vmware.vsphere.client.automation.common.annotations.UsesSpec;
import com.vmware.vsphere.client.automation.storage.lib.wizard.core.step.WizardStep;


public abstract class NewDatastoreWizardStep extends WizardStep {

   @UsesSpec()
   protected DatastoreSpec datastoreSpec;

}