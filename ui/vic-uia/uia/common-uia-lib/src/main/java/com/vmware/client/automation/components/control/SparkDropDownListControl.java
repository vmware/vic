/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import java.util.ArrayList;
import java.util.List;

import com.vmware.flexui.componentframework.controls.spark.SparkDropDownList;
import com.vmware.flexui.selenium.BrowserUtil;

/**
 * This class handles the work with StackBlock component
 */
public class SparkDropDownListControl {

   private final static String DATAPROVIDER = "dataProvider.";

   /**
    * Method that returns the number of elements in a dropdown
    * @param id - id of the dropdown
    * @return int number of the items in the dropdown
    */
   public static int getNumberOfItems(String id) {
      return Integer.parseInt(new SparkDropDownList(id,BrowserUtil.flashSelenium).getNumElements());
   }

   /**
    * Method that returns a list with the names of the items in the dropdown
    * @param id - id of the dropdown
    * @return list with the names of the items in the dropdown
    */
   public static List<String> getNamesOfItems(String id) {
      SparkDropDownList dropDown = new SparkDropDownList(id,BrowserUtil.flashSelenium);
      int size = getNumberOfItems(id);
      ArrayList<String> result = new ArrayList<String>();
      for (int i = 0; i < size; i++) {
         result.add(dropDown.getProperty(DATAPROVIDER + i));
      }

      return result;
   }
}
