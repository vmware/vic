/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.migrate.view;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.util.UiDelay;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.apl.Property;

/**
 * UI model of the Summary Page view located in: vCenter -> Virtual Machines ->
 * your VM -> Summary
 */
public class VmSummaryPageView extends BaseView {

   private static final String ID_LBL_HOST = "summary_hostName_valueLink[0]";

   /**
    * Get the Host label
    */
   public String getHost() {
      return UI.component.property.get(Property.TEXT, ID_LBL_HOST);
   }

   /**
    * Returns if the host label is visible
    * @return true if the host label is visible
    */
   public boolean isHostLabelVisible() {
      return new UIComponent(ID_LBL_HOST, BrowserUtil.flashSelenium).isVisibleOnPath();
   }

   /**
    * Check if given host is found in the VmSummaryPage
    *
    * @param expectedHost
    *           - expected name of the Host
    *
    * @return true if the expected host is found, false otherwise
    */
   public boolean isHostFound(String expectedHost) {

      long endTime = System.currentTimeMillis()
            + UiDelay.UI_OPERATION_TIMEOUT.getDuration() / 3;
      boolean expectedHostIsFound = expectedHost.equals(this.getHost());
      while (!expectedHostIsFound) {
         if (System.currentTimeMillis() > endTime) {
            _logger.error("Stop waiting for host " + expectedHost);
            break;
         }
         expectedHostIsFound = expectedHost.equals(this.getHost());
      }
      return expectedHostIsFound;
   }
}
