/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import com.vmware.flexui.componentframework.controls.mx.NavBar;
import com.vmware.flexui.selenium.BrowserUtil;

/**
 * This class handles the work with ToggleButtonBar component
 */
public class ToggleButtonBarControl {

   /**
    * Method that clicks the Toggle Button Bar at the tab represented by the index
    * @param id - id of the toggle button bar
    * @param index - index of tab to click
    */
   public static void clickItemAtIndex(String id, String index) {
      NavBar toggleButtonBar = new NavBar(id, BrowserUtil.flashSelenium);

      toggleButtonBar.selectItemByIndex(index);
   }
}
