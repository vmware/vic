/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.step;

import java.util.List;
import java.util.concurrent.TimeUnit;

import com.vmware.client.automation.common.spec.RenameEntitySpec;
import com.vmware.client.automation.vcuilib.commoncode.GlobalFunction;
import com.vmware.client.automation.vcuilib.commoncode.TestConstants.TASK_STATE;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.TestScope;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA.Environment;
import com.vmware.vsphere.client.automation.common.CommonUtil;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.common.view.RenameEntityPage;

/**
 * Renames a managed entity. The rename dialog must be open upfront.
 */
public class RenameEntityStep extends BaseWorkflowStep {

   private String _oldName;
   private String _newName;
   private boolean _newNameExists;
   private String _taskName;

   @Override
   public void prepare() throws Exception {
      RenameEntitySpec renameSpec = getSpec().get(RenameEntitySpec.class);
      if (renameSpec != null) {
         if (!renameSpec.newName.isAssigned()) {
            throw new IllegalArgumentException("You must provide a newName property in the RenameEntitySpec");
         }

         _newName = renameSpec.newName.get();
         _newNameExists = false;
         if (renameSpec.isNewNameExisting.isAssigned()) {
            _newNameExists = renameSpec.isNewNameExisting.get();
         }

         if (_newNameExists) {
            if (!renameSpec.oldName.isAssigned()) {
               throw new IllegalArgumentException("You must provide an oldName property in the RenameEntitySpec");
            }
            _oldName = renameSpec.oldName.get();

            if (!renameSpec.taskName.isAssigned()) {
               throw new IllegalArgumentException("You must provide a taskName property in the RenameEntitySpec");
            }
            _taskName = renameSpec.taskName.get();
         }
      } else {
         ManagedEntitySpec entity = getSpec().get(ManagedEntitySpec.class);
         if (entity == null) {
            throw new IllegalArgumentException("The spec has no links to 'ManagedEntitySpec' instances");
         }

         _newName = entity.name.get();
         _newNameExists = false;
      }

   }

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      RenameEntitySpec renameSpec = filteredWorkflowSpec.get(RenameEntitySpec.class);
      if (renameSpec != null) {
         if (!renameSpec.newName.isAssigned()) {
            throw new IllegalArgumentException("You must provide a newName property in the RenameEntitySpec");
         }

         _newName = renameSpec.newName.get();
         _newNameExists = false;
         if (renameSpec.isNewNameExisting.isAssigned()) {
            _newNameExists = renameSpec.isNewNameExisting.get();
         }

         if (_newNameExists) {
            if (!renameSpec.oldName.isAssigned()) {
               throw new IllegalArgumentException("You must provide an oldName property in the RenameEntitySpec");
            }
            _oldName = renameSpec.oldName.get();

            if (!renameSpec.taskName.isAssigned()) {
               throw new IllegalArgumentException("You must provide a taskName property in the RenameEntitySpec");
            }
            _taskName = renameSpec.taskName.get();
         }
      } else {
         ManagedEntitySpec entity = getSpec().get(ManagedEntitySpec.class);
         if (entity == null) {
            throw new IllegalArgumentException("The spec has no links to 'ManagedEntitySpec' instances");
         }

         _newName = entity.name.get();
         _newNameExists = false;
      }

   }

   @Override
   public void execute() throws Exception {
      RenameEntityPage renameEntityPage = new RenameEntityPage();

      new RenameEntityPage().setName(_newName);
      boolean isSubmitted = renameEntityPage.clickOk();

      if (_newName.isEmpty()) {
         verifyFatal(TestScope.FULL,
               isSubmitted == false,
               "Verifying the dialog IS NOT submitted");

         List<String> valdiationErrors = renameEntityPage.getMessagesFromValidationPanel();
         verifyFatal(TestScope.FULL,
               valdiationErrors.size() == 1,
               "Verifying the error pane has 1 error messages");

         String expectedMessage = CommonUtil.getLocalizedString("message.errorNoName");
         verifyFatal(TestScope.FULL,
               valdiationErrors.get(0).equals(expectedMessage),
               "Verifying the error pane message says: " + expectedMessage);

         // Close the dialog
         renameEntityPage.cancel();
      } else if (_newNameExists) {
         verifyFatal(TestScope.FULL,
               isSubmitted == true,
               "Verifying the dialog is submitted");

         GlobalFunction.waitUntilTasksState(
               _taskName,
               _oldName,
               TASK_STATE.TASK_FAILED,
               Environment.getUIOperationTimeout(),
               TimeUnit.SECONDS.toMillis(3),
               BrowserUtil.flashSelenium);

         String[] tasks = GlobalFunction.getAllTasks(BrowserUtil.flashSelenium);

         String expectedMessage = _newName + " " + CommonUtil.getLocalizedString("message.errorAlreadyInUse");
         verifyFatal(TestScope.FULL,
               tasks[0].contains(expectedMessage),
               "Verifying the error message contains: " + expectedMessage);
      } else {
         verifyFatal(TestScope.FULL,
               isSubmitted == true,
               "Verifying the dialog is submitted");
      }
   }
}
