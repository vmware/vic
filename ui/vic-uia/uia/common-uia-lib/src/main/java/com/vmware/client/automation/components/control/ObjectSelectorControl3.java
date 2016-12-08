/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import java.util.ArrayList;
import java.util.List;

import org.apache.commons.lang.NotImplementedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.common.custom.ButtonScrollingButtonBar;
import com.vmware.flexui.componentframework.controls.mx.custom.InventoryTree;
import com.vmware.flexui.componentframework.controls.mx.custom.ViClientPermanentTabBar;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ManagedEntityUtil;
import com.vmware.vsphere.client.automation.srv.exception.ObjectNotFoundException;

/**
 * This class handles the work with object selector version 3 UI control
 */
// TODO: rkovachev re-factor the class to follow the approach similar to PrivilegeTreeNode 
public class ObjectSelectorControl3 {

   private static final Logger _logger =
         LoggerFactory.getLogger(ObjectSelectorControl3.class);

   private static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;
   private static final String DATA_ID = "/data.id=";
   private static final IDGroup DATAGRID_ID = IDGroup.toIDGroup("tiwoDialog/list");
   private static final IDGroup LOADING_PROGRESS_BAR_ID = IDGroup
         .toIDGroup("loadingProgressBar");

   // This enum is used for the three main tabs - Browse, Filter and Selected
   // Items - that are
   // used to switch between the Tree view and the List view
   public static enum Tab {
      FILTER("0"), BROWSE("1"), SELECTED_ITEMS("2");

      private final String _index;

      private Tab(String index) {
         this._index = index;
      }

      public String getIndex() {
         return this._index;
      }

   }

   /**
    * Method to click on a tab to switch between List View (Filter tab) and Tree
    * view (Browse tab)
    *
    * @param tabBarId
    *            ID of the tabBar that hold the view buttons
    * @param tab
    *            Tab to select
    */
   public static void selectView(String tabBarId, Tab tab) {
      ViClientPermanentTabBar tabBar =
            new ViClientPermanentTabBar(tabBarId, BrowserUtil.flashSelenium);

      tabBar.selectTabAt(tab.getIndex());
   }

   /**
    * Method that returns the selected view in the object selector, i.e. Filter,
    * Browse or Selected Objects
    *
    * @param tabBarId
    *            the ID of the tabBar that holds the view buttons
    * @return Tab element that specifies what view is selected
    */
   public static Tab getSelectedView(String tabBarId) {
      ViClientPermanentTabBar tabBar =
            new ViClientPermanentTabBar(tabBarId, BrowserUtil.flashSelenium);

      return Tab.valueOf(tabBar.getSelectedIndex().toString());
   }

   /**
    * Method to click on a List view tab to switch between the different tables
    * with objects
    *
    * @param tabBarId
    *            ID of the tabBar that hold the list view buttons
    * @param subTabIndex
    *            int value that corresponds to the zero-based index if the subtab to select.
    */
   public static void selectFilterViewTab(String tabBarId, int subTabIndex)
         throws Exception {
      ButtonScrollingButtonBar buttonBar =
            new ButtonScrollingButtonBar(tabBarId, BrowserUtil.flashSelenium);

      buttonBar.selectItemViClient(Integer.toString(subTabIndex));
   }

   /**
    * Method to click on a List view tab to switch between the different tables
    * with objects
    *
    * @param tabBarId
    *            ID of the tabBar that hold the list view buttons
    * @param subTabName
    *            String value that corresponds to the name of the subtab to select
    */
   public static void selectFilterViewTab(String tabBarId, final String subTabName)
         throws Exception {
      final ButtonScrollingButtonBar buttonBar =
            new ButtonScrollingButtonBar(tabBarId, BrowserUtil.flashSelenium);

      Object evaluator = new Object() {
         @Override
         public boolean equals(Object other) {
            int subTabIndex = buttonBar.getButtonIndexByName(subTabName);
            if (subTabIndex >= 0) {
               buttonBar.selectItemViClient(Integer.toString(subTabIndex));
               return true;
            } else {
               return false;
            }
         }
      };

      if (!UI.condition.isTrue(evaluator).await(SUITA.Environment.getPageLoadTimeout())) {
         throw new RuntimeException("Can not select the sub tab: " + subTabName);
      }
   }

   /**
    * Method that is used to select an object in a table under the Object
    * Selector's List view
    *
    * @param name
    *            name of the entry to select
    * @return True if selection was successful, otherwise false
    */
   public static boolean selectFilterViewItem(ManagedEntitySpec spec) throws Exception {

      // TODO in single-selection mode radiobuttons are shown and in multiple
      // selection mode, checkboxes are shown
      // implement method body
      throw new NotImplementedException();
   }

   /**
    * Method that is used to get the selected item in the List view of the
    * object selector
    *
    * @param columnName - the name of the column whose value to get
    * @return String the name of the selected item or null if nothing is
    *         selected
    */
   public static String getSelectedFilterViewItem(String columnName) throws Exception {
      return GridControl.getSelectedEntityColumnValue(DATAGRID_ID, columnName);
   }

   /**
    * Method that is used to get the contents of a column in a table in
    * a Filter view subtab
    *
    * @param columnName - name of the column
    * @return List of stirngs with column contents
    */
   public static List<String> getFilterViewTableContents(String datagridId,
         String columnName) {
      return GridControl.getColumnContents(
            GridControl.findGrid(IDGroup.toIDGroup(datagridId)),
            columnName);
   }

   /**
    * Method that is used to get the selected subaTab in the List view of the
    * object selector
    *
    * @return String the name of the selected subTab or null if nothing is
    *         selected
    */
   public static String getSelectedFilterViewSubTab() throws Exception {

      // implement method body
      throw new NotImplementedException();
   }

   /**
    * Method that expands the inventory tree if needed and selects the specified
    * node
    *
    * @param invTreeId
    *            - String id of the inventory tree
    * @param spec
    *            - the ManagedEntitySpec of the item to select
    * @return True if node was successfully selected, false otherwise
    * @throws Exception
    *             Throws an exception if the action script calls to the Inventory
    *             tree fail
    */
   public static boolean selectBrowseViewItem(String invTreeId, ManagedEntitySpec spec)
         throws Exception {
      InventoryTree inventoryTree =
            new InventoryTree(invTreeId, BrowserUtil.flashSelenium);
      try {
         expandBrowseViewNode(invTreeId, spec);
         String morNodeId = constructNodeId(spec);
         try {
            waitForNodeReady(
                  invTreeId,
                  morNodeId,
                  (int) SUITA.Environment.getUIOperationTimeout());
            inventoryTree.selectMultiViewInvNode(morNodeId);
         } catch (AssertionError ae) {
            throw new Exception("Assertion Error has occured while navigating tree: "
                  + ae.getMessage());
         }

         return spec.name.get().equals(getSelectedBrowseViewItem(invTreeId));
      } catch (ObjectNotFoundException oe) {
         // this exception is caught here, b/c there are tests that try to verify a specified object is no longer
         // seen in obj controller
         _logger.info("Object is not present in inventory!");
      }
      return false;
   }

   /**
    * Method that expands the whole inventory tree
    *
    * @return True if tree was successfully expanded, otherwise false
    */
   public static boolean expandBrowseView() throws Exception {

      // implement method body
      throw new NotImplementedException();
   }

   /**
    * Method that expands a specified node in the inventory tree
    *
    * @param invTreeId
    *            String id of the inventory tree
    * @param spec
    *            The ManagedEntitySpec of the item to expand
    *
    */
   public static void expandBrowseViewNode(String invTreeId, ManagedEntitySpec spec)
         throws Exception {
      InventoryTree inventoryTree =
            new InventoryTree(invTreeId, BrowserUtil.flashSelenium);
      List<String> path = buildInventoryPath(spec);
      try {
         boolean isRootNode = true;
         for (String node : path) {
            if(waitForNodeReady(
                  invTreeId,
                  node,
                  (int) SUITA.Environment.getUIOperationTimeout())) {
               inventoryTree.expandMultiViewInvNode(node);
            } else if (isRootNode) {
               _logger.warn("The Object selectore does NOT show the VC as root node!");
            } else {
               throw new Exception("Node with ID is not visible:" + node + " for " + spec.toString());
            }
            isRootNode = false;
         }
      } catch (AssertionError ae) {
         throw new Exception("Assertion Error has occured while navigating tree: "
               + ae.getMessage(), ae);
      }
   }

   /**
    * Method that returns the Name of the selected node in the tree
    *
    * @param invTreeId
    *            String id of the inventory tree
    * @return String - the name of the selected element or null if nothing
    *         is selected
    */
   public static String getSelectedBrowseViewItem(String invTreeId) throws Exception {
      UIComponent navTree = new UIComponent(invTreeId, BrowserUtil.flashSelenium);
      return navTree.getProperty("selectedItems.length").equals("0") ? null : navTree
            .getProperty("selectedItem.label");
   }

   /**
    * Method that returns the labels of all browse view items that should be visible
    * in UI
    *
    * @param invTreeId
    *            String id of the inventory tree
    * @return List<String> - list of all teh labels of present Browse View items
    */
   public static List<String> getAllVisibleBrowseViewItems(String invTreeId)
         throws Exception {
      ArrayList<String> visibleItems = new ArrayList<String>();
      final String dataProvider0 = "dataProvider.0";
      UIComponent navTree = new UIComponent(invTreeId, BrowserUtil.flashSelenium);
      int length = Integer.parseInt(navTree.getProperty("dataProvider.length"));
      if (length != 0) {
         visibleItems.add(navTree.getProperty(dataProvider0 + ".label"));
         boolean nodeOpen =
               navTree.getProperty(dataProvider0 + ".nodeOpen").equalsIgnoreCase("true");
         int childrenLength =
               Integer.parseInt(navTree.getProperty(dataProvider0 + ".children.length"));
         visibleItems.addAll(getOpenItemsNames(
               dataProvider0 + ".children",
               nodeOpen,
               childrenLength,
               navTree));

      }
      return visibleItems;
   }

   /**
    * Method that return the number of root elements in the Object Selector
    * Browse view, i.e. number of VC servers or number of datacenters
    *
    * @param invTreeId
    * @return number of items
    */
   public static int getNumberOfFirstLevelItemsInBrowseView(String invTreeId) {
      UIComponent navTree = new UIComponent(invTreeId, BrowserUtil.flashSelenium);
      return Integer.parseInt(navTree.getProperty("dataProvider.length"));
   }

   /**
    * Method that gets the label value of the first root element in Browse view,
    * i.e. first VC server, first datacenter, etc.
    *
    * @param invTreeId
    * @return Name of first item or null if no such
    */
   public static String getNameOfFirstRootItemInBrowseView(String invTreeId) {
      UIComponent navTree = new UIComponent(invTreeId, BrowserUtil.flashSelenium);

      if (Integer.parseInt(navTree.getProperty("dataProvider.length")) > 0) {
         return navTree.getProperty("dataProvider.0.label");
      }

      return null;
   }

   /**
    * Checks if loading progress bar is available.
    * It does not check if it is visible or not.
    *
    * @return if the loading progress bar is present on the page
    */
   public static boolean loadingProgressBarPresent() {
      return UI.condition.isFound(LOADING_PROGRESS_BAR_ID).await(
            SUITA.Environment.getUIOperationTimeout() / 2);
   }

   /**
    * Waits loading progress bar of the object selector to disappear.
    *
    * @throws AssertionError if the loading bar is still visible after the configured
    *             timeout has passed
    */
   public static void waitForObjectSelectorToLoad() {
      Object evaluator = new Object() {
         @Override
         public boolean equals(Object other) {
            return UI.component.property.getBoolean(Property.SHOW_PROGRESS, DATAGRID_ID)
                  .equals(other);
         }
      };

      UI.condition.isTrue(evaluator).await(SUITA.Environment.getPageLoadTimeout());
   }

   /**
    * Method that goes through the inventory tree in ObjectSelector:
    *
    * dataProvider.0.children.0.children.0, etc., while there are such
    * elements and their node is open
    *
    * @param dataProvider
    * @param nodeOpen
    * @param length
    * @param navTree
    * @return List with the labels of visible nodes
    */
   private static List<String> getOpenItemsNames(String dataProvider, boolean nodeOpen,
         int length, UIComponent navTree) {
      List<String> visibleItems = new ArrayList<String>();
      if (nodeOpen) {
         for (int i = 0; i < length; i++) {
            visibleItems.add(navTree.getProperty(dataProvider + "." + i + ".label"));
            nodeOpen =
                  navTree.getProperty(dataProvider + "." + i + ".nodeOpen")
                        .equalsIgnoreCase("true");
            int nodeLength =
                  Integer.parseInt(navTree.getProperty(dataProvider + "." + i
                        + ".children.length"));
            visibleItems.addAll(getOpenItemsNames(
                  dataProvider + "." + i + ".children",
                  nodeOpen,
                  nodeLength,
                  navTree));
         }
      }
      return visibleItems;
   }

   /**
    * This function waits for the node to be Visible. nodeId should be in the
    * same format as String returned by ManagedEntityUtil.getInvTreePathToEntity
    *
    * @param invTreeId
    *            String id of the inventory tree
    * @param nodeId
    *            the last element in the list returned by
    *            ManagedEntityUtil.getInvTreePathToEntity
    * @param timeout
    *            the number of seconds to wait for the object to become visible
    * @return true if node becomes visible within timeout(in seconds), else
    *         false
    */
   private static boolean waitForNodeReady(String invTreeId, String nodeId, int timeout) {
      String nodeComponentId = invTreeId + DATA_ID + nodeId;
      _logger.info("NODE_ID:" + nodeComponentId);
      return UI.condition.isFound(nodeComponentId).await(timeout);
   }

   /**
    * Constructs the node id of the object navigator tree.
    * The automation uses the internal data . id to navigate and node id is based on the entity href.
    * Currently the id format is - serverGUID:type:value
    * TODO: rkovachev use the platform provided mechanism to get the mor to string or better to use labels as id.
    * @param moRef
    *           The managed object reference to be converted
    * @return its string representation serverGUID:type:value
    * @throws Exception
    */
   private static String constructNodeId(ManagedEntitySpec spec) throws Exception{
      ManagedObjectReference moRef = ManagedEntityUtil.getManagedObject(spec)._getRef();
      VcService service = VcServiceUtil.getVcService(spec.service.get());
      // this is a workaround as VCService returns null for severGuid of moRefs
      String serverGuid = service.getServiceGuid();
      return serverGuid + ":" + moRef.getType() + ":" + moRef.getValue();
   }

   /**
    * Utility method that creates a List<String> path that contains all the
    * parent and child'sMO refs that lead to the entity in the spec. The last
    * element is the name of the entity, for which a path is built.
    *
    * @param spec
    *           The ManagedEntitySpec of the entity whose path we want to get
    * @param serviceSpec
    *           specification of the LDU connection details.
    * @return List<String> List of Strings in the form of moRefs
    *         (serverGUID:type:value), where the first element is the root and
    *         the last is the entity whose path is being built
    * @throws Exception
    *             In case that it cannot log to vc service, and exception will be
    *             thrown
    */
   private static List<String> buildInventoryPath(ManagedEntitySpec spec) throws Exception {
      List<String> path = new ArrayList<String>();
      path.add(0, constructNodeId(spec));

      ManagedEntitySpec tempSpec = spec;
      while (tempSpec.parent.isAssigned()) {
         tempSpec = tempSpec.parent.get();
         path.add(0, constructNodeId(tempSpec));
      }
      return path;
   }
}
