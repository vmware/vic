/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.ft.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.util.UiDelay;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.srv.common.spec.FaultToleranceSpec.FaultToleranceStatus;

/**
 * Represents the fault tolerance portlet, as seen on the Summary tab of a VM.
 */
public class FtPortletView extends BaseView {
   private static final IDGroup FT_STATUS = IDGroup.toIDGroup("ftStatusLabel");
   private static final IDGroup FT_STATE = IDGroup.toIDGroup("ftStateLabel");
   private static final IDGroup SECONDARY_VM_LOCATION = IDGroup
         .toIDGroup("secondaryLocationLinkBtn");

   /**
    * Returns the Fault Tolerance status. Possible statuses are Not protected,
    * Protected, etc.
    *
    * @return the fault tolerance status
    */
   public String getFaultToleranceStatus() {
      return UI.component.property.get(Property.TEXT, FT_STATUS);
   }

   /**
    * Returns the Fault Tolerance state. Possible states are Starting,
    * Suspended, etc.
    *
    * @return the fault tolerance state
    */
   public String getFaultToleranceState() {
      return UI.component.property.get(Property.TEXT, FT_STATE);
   }

   /**
    * Returns the Secondary VM Location in the form of a host IP.
    *
    * @return the secondary vm location
    */
   public String getSecondaryVmLocation() {
      return UI.component.property.get(Property.TEXT, SECONDARY_VM_LOCATION);
   }

   /**
    * Wait and refresh the page until the expected status has been reached.
    *
    * @param expectedStatus
    *           - the status we expect
    * @return true if the expected status is reached, false otherwise
    * @throws InterruptedException
    */
   public boolean waitForStatus(FaultToleranceStatus expectedStatus)
         throws InterruptedException {
      long endTime = System.currentTimeMillis()
            + UiDelay.UI_OPERATION_TIMEOUT.getDuration();

      while (!expectedStatus.getValue().equals(getFaultToleranceStatus())) {
         this.refreshPage();
         if (System.currentTimeMillis() > endTime) {
            _logger.error("Stop waiting for status %s", expectedStatus);
            break;
         }
      }
      return expectedStatus.getValue().equals(getFaultToleranceStatus());
   }
}