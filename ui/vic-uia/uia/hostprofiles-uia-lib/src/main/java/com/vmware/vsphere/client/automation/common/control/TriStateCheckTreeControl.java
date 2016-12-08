/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.control;

import java.util.ArrayList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.delay.Delay;
import com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey;
import com.vmware.flexui.componentframework.controls.common.custom.TriStateChechBox;
import com.vmware.flexui.selenium.BrowserUtil;

/**
 * This is a control that describes the TriStateCheckTree that is found in Host Profiles
 * wizards, for example. It is a tree with multiple root nodes, that has a checkbox for
 * some nodes, that can be expanded/collapsed and checked.
 */
public class TriStateCheckTreeControl {

   private static final Logger _logger = LoggerFactory
         .getLogger(TriStateCheckTreeControl.class);

   private final TriStateChechBox _triStateCheckTree;

   private static final String DATAPROVIDER = "dataProvider.0";
   private static final String DATAPROVIDER_CHILDREN = ".children.";
   private static final String DATAPROVIDER_CHILDREN_LENGTH = "length";
   private static final String DATAPROVIDER_NODE_LABEL = "%s.label";
   private static final String DATAPROVIDER_NODE_VISIBLE = "%s.visible";
   private static final String DATAPROVIDER_NODE_ICC = "%s.ignoreComplianceCheckState";


   public TriStateCheckTreeControl(String id) {
      _triStateCheckTree = new TriStateChechBox(id, BrowserUtil.flashSelenium);
   }

   /**
    * Gets a List of the root nodes that have visible property set to true, i.e. are
    * displayed in UI, through the Dataprovider
    *
    * @return List of Strings with labels or list with single empty string, if no category
    *         is displayed
    */
   public List<String> getVisibleCategories() {
      String sizeProp =
            DATAPROVIDER + DATAPROVIDER_CHILDREN + DATAPROVIDER_CHILDREN_LENGTH;
      String condition = Boolean.TRUE.toString().toLowerCase();
      String filterProperty =
            DATAPROVIDER + DATAPROVIDER_CHILDREN + DATAPROVIDER_NODE_VISIBLE;
      String resultProperty =
            DATAPROVIDER + DATAPROVIDER_CHILDREN + DATAPROVIDER_NODE_LABEL;
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
            DATAPROVIDER + DATAPROVIDER_CHILDREN + DATAPROVIDER_CHILDREN_LENGTH;
      String condition = enableIccState;
      String filterProperty =
            DATAPROVIDER + DATAPROVIDER_CHILDREN + DATAPROVIDER_NODE_ICC;
      String resultProperty =
            DATAPROVIDER + DATAPROVIDER_CHILDREN + DATAPROVIDER_NODE_LABEL;
      return getFilteredNodes(sizeProp, condition, filterProperty, resultProperty);
   }

   /**
    * Method that lists the names of the direct visible children of the node that is a
    * parameter. The method opens by itself the tree to the node.
    *
    * @param nodeNamesInPath - path to the node whose children have to be checked for
    *           visibility
    * @return List of the names of the visible children, or empty list
    * @throws Exception if the node is inexistent
    */
   public List<String> getVisibleDirectChildren(List<String> nodeNamesInPath)
         throws Exception {
      String dpPrefix = getAndExpandNode(nodeNamesInPath, false).getDpPrefix();
      String sizeProp = dpPrefix + DATAPROVIDER_CHILDREN + DATAPROVIDER_CHILDREN_LENGTH;
      String condition = Boolean.TRUE.toString().toLowerCase();
      String filterProperty =
            dpPrefix + DATAPROVIDER_CHILDREN + DATAPROVIDER_NODE_VISIBLE;
      String resultProperty = dpPrefix + DATAPROVIDER_CHILDREN + DATAPROVIDER_NODE_LABEL;
      return getFilteredNodes(sizeProp, condition, filterProperty, resultProperty);
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
      String dpPrefix = getAndExpandNode(nodeNamesInPath, false).getDpPrefix();
      String sizeProp = dpPrefix + DATAPROVIDER_CHILDREN + DATAPROVIDER_CHILDREN_LENGTH;
      String condition = enableIccState;
      String filterProperty = dpPrefix + DATAPROVIDER_CHILDREN + DATAPROVIDER_NODE_ICC;
      String resultProperty = dpPrefix + DATAPROVIDER_CHILDREN + DATAPROVIDER_NODE_LABEL;
      return getFilteredNodes(sizeProp, condition, filterProperty, resultProperty);
   }

   /**
    * Method that selects a node in the tree, needs a path in starting from the root node
    * of the node to select and containing all parents and the node to select as a last
    * element
    *
    * @param nodeNamesInPath - path as described above
    * @return true if successful, false otherwise
    * @throws Exception
    */
   public boolean selectItem(List<String> nodeNamesInPath) {
      try {
         Node node = getAndExpandNode(nodeNamesInPath, true);
         List<String> returnUidList = node.getReturnUidList();
         _triStateCheckTree.selectItem(returnUidList.get(returnUidList.size() - 1));
         return true;
      } catch (Exception e) {
         _logger.info(e.getStackTrace().toString());
      }

      return false;
   }

   /**
    * Method that expands a node in the tree, needs a path in starting from the root node
    * of the node to select and containing all parents and the node to select as a last
    * element
    *
    * @param nodeNamesInPath - path as described above
    * @return true if successful, false otherwise
    * @throws Exception
    */
   public boolean expandNode(List<String> nodeNamesInPath) {
      try {
         getAndExpandNode(nodeNamesInPath, true);
         return true;
      } catch (Exception e) {
         _logger.info(e.getStackTrace().toString());
      }

      return false;
   }


   /**
    * Method that returns the value from the dataprovider of the tree for a specified
    * property
    *
    * @param nodeNamesInPath - the path to the node
    * @param propName - the property whose value is needed
    * @return string value of the property
    * @throws Exception
    */
   public String getDataproviderPropertyValue(List<String> nodeNamesInPath,
         String propName) throws Exception {
      Node node = getAndExpandNode(nodeNamesInPath, false);
      return _triStateCheckTree.getProperty(node.getDpPrefix() + "." + propName);
   }

   // private helper methods
   /**
    * This class represents a node with its dataprovider uid of type 0.3.8. etc. and
    * dataprovider prefix pf type dataprovider.children.0.children.3 etc.
    *
    * Uid is used, for check or expand operations and dpprefix is needed for verification
    * of property values per node
    */
   private class Node {
      // dataprovider uid of type 0.3.8. etc.
      private List<String> returnUidList;
      // dataprovider prefix pf type dataprovider.children.0.children.3 etc.
      private String dpPrefix;

      public Node(List<String> returnUidList, String dpPrefix) {
         this.returnUidList = returnUidList;
         this.dpPrefix = dpPrefix;
      }

      public List<String> getReturnUidList() {
         return this.returnUidList;
      }

      public String getDpPrefix() {
         return this.dpPrefix;
      }
   }

   private Node getAndExpandNode(List<String> nodeNamesInPath, boolean expand)
         throws Exception {
      _triStateCheckTree
            .waitForElementEnable(Delay.timeout.forSeconds(30).getDuration());
      String dpPrefix = DATAPROVIDER;
      String currName = "";
      String currUID = "";
      List<String> returnUidList = new ArrayList<String>();
      _logger.debug("search Nodes [" + nodeNamesInPath + "]\n");
      for (String nodeName : nodeNamesInPath) {
         dpPrefix = dpPrefix + DATAPROVIDER_CHILDREN;
         int length =
               Integer.parseInt(_triStateCheckTree.getProperty(dpPrefix
                     + DATAPROVIDER_CHILDREN_LENGTH));
         boolean foundNode = false;
         for (int i = 0; i < length; i++) {
            currName =
                  _triStateCheckTree.getProperty(dpPrefix + i
                        + TestConstantsKey.DATA_PROVIDER_LIST_LABEL);
            if (nodeName.equals(currName)) {
               dpPrefix = dpPrefix + i;
               currUID = _triStateCheckTree.getProperty(dpPrefix + ".uid");
               returnUidList.add(currUID);
               if (expand) {
                  _triStateCheckTree.expandItem(currUID);
                  _triStateCheckTree.waitForElementEnable(Delay.timeout.forMinutes(1)
                        .getDuration());
                  _logger.debug("expanding... UID=[" + currUID + "]\t\tfor Label=["
                        + currName + "]\n");
               }
               foundNode = true;
               break;
            }
         }
         if (foundNode) {
            _logger.info("node Name Found [" + nodeName + "] at Prefix [" + dpPrefix
                  + "]");
         } else {
            _logger.info("node Name NOT Found [" + nodeName + "] at Prefix [" + dpPrefix
                  + "]");
            throw new Exception("Node not found!");
         }
      }

      return new Node(returnUidList, dpPrefix);
   }

   private List<String> getFilteredNodes(String sizeProperty, String condition,
         String filterProperty, String resultProperty) {
      List<String> result = new ArrayList<String>();
      int dataProviderSize =
            Integer.parseInt(_triStateCheckTree.getProperty(sizeProperty));

      for (int i = 0; i < dataProviderSize; i++) {
         if (condition.equals(_triStateCheckTree.getProperty(String.format(
               filterProperty,
               i)))) {
            result.add(_triStateCheckTree.getProperty(String.format(resultProperty, i)));
         }
      }

      return result;
   }
}
