/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.control;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Stack;

import org.apache.commons.collections.CollectionUtils;

import com.vmware.client.automation.delay.Delay;
import com.vmware.flexui.componentframework.controls.mx.Tree;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;

/**
 * Class that in teh specified tree allows selection of a node by the path to the node
 * through finding nodes indices in the dataprovider
 */
public class HostProfileTreeControl extends Tree {
   private final static String DATAPROVIDER = "dataProvider.0";
   private final static String CHILDREN = ".children";
   private final static String LENGTH = ".length";
   private final static String LABEL = ".label";
   private final static String VISIBLE = ".visible";
   private final static String DOT = ".";
   private final static String NODE_ICC  = ".%s.ignoreComplianceCheckState";

   private final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;
   private String filterId;

   public HostProfileTreeControl(String id, String filterId) {
      super(id, BrowserUtil.flashSelenium);
      this.filterId = filterId;
   }

   /**
    * Method that selects an item and returns the last item's index
    *
    * @param path - path to the item including all ndoes to it
    */
   public void selectItem(java.util.List<String> path) {
      if (path == null) {
         throw new RuntimeException("There is no path passed!");
      }
      int lastElementIndex = path.size() - 1;
      selectItem(String.valueOf(lastElementIndex));
   }

   /**
    * Method that executes the filter on the view with supplied string
    * @param filter - string to filter by
    */
   public void filter(String filter) {
      UI.component.value.set(filter, filterId);

      // wait for filter to execute
      Delay.sleep.forSeconds(5).consume();

      UI.component.setFocus(filterId);
   }

   /**
    * Method that returns only the visible categories
    *
    * @return - the currently visible categories
    */
   public List<String> getVisibleCategories() {
      return getVisibleNodes(DATAPROVIDER);
   }

   /**
    * Method to get all visible children of a node
    *
    * @param path
    * @return
    */
   public List<String> getVisibleDirectChildren(List<String> path) {
      String dataProviderPath = getNodeDataproviderPrefix(path);
      return getVisibleNodes(dataProviderPath);
   }

   /**
    * Method that gets the direct children of the node described by the path, that have
    * the specified ignore check compliance state
    *
    * @param nodeNamesInPath - names of the parent nodes that form teh path to th enode
    * @param enableIccState - specified state: ignoreComplianceCheck,
    *           transientIgnoreComplianceCheck and nonIgnoreComplianceCheck
    * @return list of the names of node that respect the criteria or an empty arraylist
    * @throws Exception if the node doesn't exist
    */
   public List<String> getEnabledIccStateDirectChildren(List<String> nodeNamesInPath,
                                                        String enableIccState) throws Exception {
      String dpPrefix = getNodeDataproviderPrefix(nodeNamesInPath);
      String sizeProp = dpPrefix + CHILDREN + LENGTH;
      String condition = enableIccState;
      String filterProperty = dpPrefix + CHILDREN + NODE_ICC;
      String resultProperty = dpPrefix + CHILDREN + ".%s" + LABEL;
      return getFilteredNodes(sizeProp, condition, filterProperty, resultProperty);
   }

   /**
    * Method that gets the list of root nodes that have the specified ignore check
    * compliance state
    *
    * @param enableIccState - specified state: ignoreComplianceCheck,
    *           transientIgnoreComplianceCheck and nonIgnoreComplianceCheck
    * @return list of the names of node that respect the criteria or an empty arraylist
    * @throws Exception if the node doesn't exist
    */
   public List<String> getEnabledIccStateCategories(String enableIccState)
           throws Exception {
      String sizeProp =
              DATAPROVIDER +CHILDREN + LENGTH;
      String condition = enableIccState;
      String filterProperty =
              DATAPROVIDER + CHILDREN + NODE_ICC;
      String resultProperty =
              DATAPROVIDER + CHILDREN + ".%s"+ LABEL;
      return getFilteredNodes(sizeProp, condition, filterProperty, resultProperty);
   }

   /**
    * Method that gets the property value of a node
    *
    * @param path - list of the names of the nodes preceding the node to get and including
    *           the node itself
    * @param propName - name of property of the node to get
    * @return - value of the proeprty
    * @throws Exception
    */
   public String getDataproviderPropertyValue(List<String> path, String propName)
         throws Exception {
      String dataProviderPath = getNodeDataproviderPrefix(path);
      return getProperty(dataProviderPath + "." + propName);
   }

   // Helper methods
   private List<String> getVisibleNodes(String dpPrefix) {
      List<String> result = new ArrayList<String>();

      dpPrefix += CHILDREN;
      int numOfChildren = Integer.valueOf(getProperty(dpPrefix + LENGTH)).intValue();

      for (int i = 0; i < numOfChildren; i++) {
         if (Boolean.TRUE.toString().equalsIgnoreCase(
               getProperty(dpPrefix + DOT + i + VISIBLE))) {
            result.add(getProperty(dpPrefix + DOT + i + LABEL));
         }
      }

      return result;
   }

   private String getNodeDataproviderPrefix(List<String> path) {
      if (CollectionUtils.isEmpty(path)) {
         throw new RuntimeException("No path supplied!");
      }

      StringBuilder sbNodeDp = new StringBuilder(DATAPROVIDER + CHILDREN);

      Collections.reverse(path);
      Stack<String> stackPath = new Stack<String>();
      stackPath.addAll(path);

      return getDpPrefix(stackPath, sbNodeDp);
   }

   private String getDpPrefix(Stack<String> path, StringBuilder sb) {
      boolean isFound = false;
      int length = Integer.valueOf(getProperty(sb + LENGTH)).intValue();

      String name = path.pop();
      for (int i = 0; i < length; i++) {
         if (name.equals(getProperty(sb + DOT + i + LABEL))) {
            sb.append("." + i);
            isFound = true;
            break;
         }
      }
      if (path.empty()) {
         if (sb.length() != 0 && isFound) {
            return sb.toString();
         } else {
            throw new RuntimeException("Node's dataProvider id not found!");
         }
      } else {
         sb.append(CHILDREN);
         return getDpPrefix(path, sb);
      }
   }

   private List<String> getFilteredNodes(String sizeProperty, String condition,
                                         String filterProperty, String resultProperty) {
      List<String> result = new ArrayList<String>();
      int dataProviderSize =
              Integer.parseInt(this.getProperty(sizeProperty));

      for (int i = 0; i < dataProviderSize; i++) {
         if (condition.equals(this.getProperty(String.format(
                 filterProperty,
                 i)))) {
            result.add(this.getProperty(String.format(resultProperty, i)));
         }
      }

      return result;
   }

}
