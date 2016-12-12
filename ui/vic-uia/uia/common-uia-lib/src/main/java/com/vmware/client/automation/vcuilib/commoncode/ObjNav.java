/*
 * ************************************************************************
 *
 * Copyright 2011 VMware, Inc.  All rights reserved. -- VMware Confidential
 *
 * ************************************************************************
 */

package com.vmware.client.automation.vcuilib.commoncode;

import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.SELECTED_SET_DATAGRID;
import static com.vmware.client.automation.vcuilib.commoncode.TestBaseUI.verifySafely;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.DEFAULT_TIMEOUT_ONE_SECOND_INT_VALUE;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.DEFAULT_TIMEOUT_TEN_SECONDS_LONG_VALUE;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.Timeout.ONE_SECOND;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.EQUALS_SIGN;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.INV_TREE_NODE_ID_NAME_PROPERTY;
import static com.vmware.flexui.selenium.BrowserUtil.flashSelenium;

import java.util.ArrayList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.thoughtworks.selenium.SeleniumException;
import com.vmware.client.automation.vcuilib.commoncode.TestConstants.OBJ_NAV_TREE_NODE_VIEW;
import com.vmware.flexui.componentframework.controls.mx.Label;
import com.vmware.flexui.componentframework.controls.mx.custom.ObjectNavigator;

/**
 * Object Navigator implementation.
 *
 * NOTE: this class is a partial copy of the one from VCUI-QE-LIB
 */
public class ObjNav {

   private static final Logger logger = LoggerFactory.getLogger(ObjNav.class);

   /**
    * Method to verify if the node with the provided name is visible in the
    * selected sets datagrid. This is the lower datagrid in the Object Navigator
    *
    * @param nodeName String Node name
    * @throws Exception
    */
   public static void verifyNodeNameVisible(String nodeName) throws Exception {
      verifyNodeNameVisible(nodeName, true);
   }

   /**
    * Method to verify if the node with the provided name is visible or not
    * in the selected sets datagrid. This is the lower datagrid in the
    * Object Navigator
    *
    * @param nodeName String Node name
    * @param isNodeVisible Boolean If node to find should be visible
    * @throws Exception
    */
   public static void verifyNodeNameVisible(String nodeName, Boolean isNodeVisible)
         throws Exception {

      ObjectNavigator selectedSetsDatagrid =
            new ObjectNavigator(SELECTED_SET_DATAGRID, flashSelenium);
      int rowIndex = -1;
      Boolean bFound = false;
      rowIndex = ObjNav.getRowIndexForNodeName(selectedSetsDatagrid, nodeName);
      bFound = (rowIndex > -1) ? true : false;

      verifySafely(bFound, isNodeVisible, "Node: " + nodeName + " visible in the "
            + "Object Navigator");
   }

   /**
    * Method to Select Node in the Object Navigator using name
    *
    * @param objNavView OBJ_NAV_TREE_NODE_VIEW Tree node view in Object
    *           Navigator
    * @param nodeToNavigateMor MOR of the node to navigate
    * @return boolean isNavigated true or false based on the node selection
    * @throws Exception
    */
   // TODO: rkovachvev the method is full of workarounds and not proper check of internal
   // flex properties. PR 1299768
   public static boolean selectNodeByName(OBJ_NAV_TREE_NODE_VIEW objNavView,
         String nodeName) throws Exception {
      Boolean isNavigated = false;

      String nodeId = INV_TREE_NODE_ID_NAME_PROPERTY + EQUALS_SIGN + nodeName;

      // Select the node in the selectedSetDatagrid
      Label labelObject1 =
            new Label(SELECTED_SET_DATAGRID + "/" + nodeId, flashSelenium);
      Label labelObject2 = labelObject1;

      // Workaround: The last element appears twice and one of the instances is not visible.
      // If this is the one - just use the second instance.
      // We still have to use the first instance for mouse clicks.
      // PR 1299768: Why do we check for labelObject2 but at the end we use labelObject1?
      if (!labelObject2.getVisible()) {
         logger.info(nodeId + " is not visible try with index [1]");
         if (!labelObject2.isVisibleOnPath()) {
            logger.info(nodeId + " is not visible on path");
         } else {
            logger.info(nodeId + " is visible on path");
         }
         labelObject2 =
               new Label(SELECTED_SET_DATAGRID + "/" + nodeId + "[1]", flashSelenium);
      }

      try {
         // Wait for the element to be visible(10 secs)
         labelObject2
               .waitForElementVisibleOnPath(DEFAULT_TIMEOUT_TEN_SECONDS_LONG_VALUE);
      } catch (Throwable t) {
         logger.error("The node navigation is not visible on path. "
               + "Though will try to click on it as there is a cases "
               + "when the visible property is not populated on slow environments!");
         t.printStackTrace();
      }

      // We can still only click on the first instance
      labelObject1.visibleMouseDownUp();
      labelObject1.visibleClick(DEFAULT_TIMEOUT_ONE_SECOND_INT_VALUE);
      logger.info("Navigated successfully to: " + nodeId + " from " + "Object Navigator");
      // PR 1299768: Useless check for labelObject2 status as the SUITA wrapper never use the returned value.
      // isNavigated = Boolean.parseBoolean(labelObject2.getProperty("selected"));

      return isNavigated;
   }

   /**
    * Method to get the row Index for Node in a datagrid
    *
    * @param nodeSetsDataGrid ObjectNavigator
    * @param nodeName String Node name
    * @return int Row Index for the node
    */
   public static int getRowIndexForNodeName(ObjectNavigator nodeSetsDataGrid,
         String nodeName) {

      int numElements = retryIfNodeListNotEntirelyLoaded(nodeSetsDataGrid);
      List<String> itemsList = new ArrayList<String>(1);
      for (int i = 0; i < numElements; i++) {
         itemsList.add(nodeSetsDataGrid.getProperty("dataProvider." + i + ".name"));
      }

      return itemsList.indexOf(nodeName);
   }

   /**
    * Helper method to handle problems when data provider length of the nodes
    * datagrid is not entirely loaded. Returns number of elements of the
    * datagrid.
    *
    * @param edgeSetsDataGrid
    * @return
    */
   private static int retryIfNodeListNotEntirelyLoaded(ObjectNavigator edgeSetsDataGrid) {
      int numElements = 0;
      try {
         // Sometimes on slow environments the dataProvider property is not
         // fully loaded and this method throws Selenium Exception for NPObject
         numElements = new Integer(edgeSetsDataGrid.getProperty("dataProvider.length"));
         // The catch will give another second in case of NPObject exception and
         // retry again to get the number of elements in the datagrid
      } catch (SeleniumException e) {
         e.printStackTrace();
         logger.info("Object navigator data provider length is null at "
               + "this point, trying once again after 1 sec...");
         ONE_SECOND.consume();
         numElements = new Integer(edgeSetsDataGrid.getProperty("dataProvider.length"));
      }
      return numElements;
   }
}
