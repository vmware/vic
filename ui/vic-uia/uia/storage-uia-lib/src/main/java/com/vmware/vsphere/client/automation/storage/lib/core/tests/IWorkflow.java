package com.vmware.vsphere.client.automation.storage.lib.core.tests;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.common.WorkflowStep;

/**
 * Class to represent the flow of steps
 */
public interface IWorkflow {

   /**
    * Append new step to the flow of steps</br>
    *
    * NOTE: The white list is by step type. </br>
    * Example: </br>
    * <code>Red specRed1;</code></br>
    * <code>Red specRed2;</code></br>
    * <code>Blue specBlue1;</code></br>
    * <code>Blue specBlue2;</code></br>
    * <code>append(title, step, specRed2) // Available specs for the step - specRed2, specBlue1, specBlue2</code>
    * </br>
    * <code>append(title, step, specRed2, specBlue1) // Available specs for the step - specRed2, specBlue1</code>
    * </br>
    *
    * @param title
    *           of the step
    * @param step
    * @param specsWhitelist white list of specs, which will be available for the newly added step
    */
   void append(String title, WorkflowStep step, BaseSpec... specsWhitelist);
}
