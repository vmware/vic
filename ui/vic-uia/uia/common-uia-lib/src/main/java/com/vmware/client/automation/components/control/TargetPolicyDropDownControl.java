/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import com.vmware.flexui.componentframework.controls.spark.SparkDropDownList;
import com.vmware.flexui.selenium.BrowserUtil;

/**
 * This class handles the work with TargetPolicyDropDown component
 */
public class TargetPolicyDropDownControl {

   /**
    * Method to select and item by its name
    * @param id - id of the UI component
    * @param name - String name to select
    * @return - true if name's indexc >=0 and selected, else false
    */
   public static boolean selectItemByName(String id, String name) {
      SparkDropDownList dropDown = new SparkDropDownList(id, BrowserUtil.flashSelenium);
      boolean result;

      String index = dropDown.getIndexForValue(name);
      result = Integer.parseInt(index) >=0;
      if (result) {
         dropDown.selectItemByIndex(index);
      }
      return result;
   }

   /**
    * Method that gets the current preselection in the specified dropdown
    * @param id - id of the dropdown
    * @return String value of the displayed text from the preselection
    */
   public static String getSelectedItem(String id) {
       SparkDropDownList dropDown = new SparkDropDownList(id, BrowserUtil.flashSelenium);

       return dropDown.getSelectedItemText();
   }
}
