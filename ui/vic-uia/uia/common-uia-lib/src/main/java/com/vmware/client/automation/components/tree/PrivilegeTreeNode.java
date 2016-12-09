/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.tree;

import java.util.ArrayList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.collect.Lists;
import com.vmware.client.automation.components.control.VerticalScrollBar;
import com.vmware.client.automation.delay.Delay;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;

/**
 * The class defines a privilege tree node. It is used to manipulate and validate the privilege tree. It can be used in
 * both edit and view mode of the privilege tree.
 * The UI control representing the privilege tree node consist of three items - collapse/expand control, check box and
 * label.
 */
public class PrivilegeTreeNode {

   private static final Logger _logger = LoggerFactory.getLogger(PrivilegeTreeNode.class);
   private static final String CHECKBOX_CONTROL_CLASS_NAME = "TriStateCheckBox";
   private static final String EXPAND_COLLAPSE_CONTROL_CLASS_NAME = "SpriteAsset";
   private static final String SCROLL_BAR_ID = "/className=VScrollBar";

   private static final String CB_STATE_CHECKED = "checked";
//   private static final String CB_STATE_TRANSIENT = "transient";
   private static final String CB_STATE_UNCHECKED = "unchecked";

   // The expand/collapse state is detected by the image name of the expand control. If it is closed its name is
   // something similar to _class_embed_css_assets_images_closedRightArrow_png_254702973. respectively when it is
   // expanded is contains Open instead of closed.
   private static final String EXPANDED_NODE_STATE_ID = "Open";
   private static final String COLLPASED_NODE_STATE_ID = "closed";

   private PrivilegeTreeNode _parent;
   // data id -
   private String _id;
   // Node ui control id based on the data provider data
   private String _uiId;
   private String _displayName;
   private int _childrenCount;
   private List<PrivilegeTreeNode> _children;
   // tree id
   private String _treeUid;

   /**
    * Create tree node.
    * @param treeUiId
    * @param parent
    * @param id
    * @param displayName
    * @param childrenCount
    */
   protected PrivilegeTreeNode(String treeUiId, PrivilegeTreeNode parent, String id, String displayName, String childrenCount) {
      this._parent = parent;
      this._displayName = displayName;
      this._id = id;
      this._uiId = buildNodeId();
      this._treeUid = treeUiId;
      this._childrenCount = parseInt(childrenCount);
      _children = new ArrayList<PrivilegeTreeNode>();
   }

   /**
    * Find child node of the current node by id. The id corresponds to the privilege id and it is unique for the
    * privileges tree. The search is done in the tree data provider data and the used algorithm is depth-first search.
    * @param nodeId id of the node.
    * @return
    */
   public PrivilegeTreeNode findNodeById(String nodeId) {
      return findNodeById(Lists.newArrayList(this), nodeId);
   }

   /**
    * Find child node of the current node by id. The id corresponds to the privilege name and it is not unique for the
    * privileges tree.
    * The method returns the first found node for that name. The search is done in the tree data provider data and the
    * used algorithm is depth-first search..
    * @param nodeLabel
    * @return
    */
   public PrivilegeTreeNode findNodeByLabel(String nodeLabel) {
      return findNodeByName(Lists.newArrayList(this), nodeLabel);
   }

   /**
    * Check or un-check tree node. If the node is already checked the method does nothing.
    * @param isSelected when true select the node otherwise un-check it.
    * @return
    */
   public boolean selectNode(boolean isSelected) {
      UIComponent checkboxTreeItemRenderer = findUiTreeItemRendered();
      return selectNode(checkboxTreeItemRenderer, isSelected);
   }

   /**
    * Return the state of the node. Possible values are: checked, unchecked or transient. Transient means that at least
    * one of the children nodes is selected, but not all of them.
    * @return String representing the selection state.
    */
   public String getSelectionState() {
      UIComponent checkboxTreeItemRenderer = findUiTreeItemRendered();
      return getSelectionState(checkboxTreeItemRenderer);
   }

   /**
    * Expand the node. If already expanded does nothing.
    * @return true if the node is expanded.
    */
   public boolean expandNode() {
      _logger.info("Expand node:" + _displayName);
      UIComponent checkboxTreeItemRenderer = findUiTreeItemRendered();
      return expandNode(checkboxTreeItemRenderer, true);
   }

   /**
    * Collapse the node. If already collapsed does nothing.
    * @return true if the node is collapsed.
    */
   public boolean collapseNode() {
      UIComponent checkboxTreeItemRenderer = findUiTreeItemRendered();
      return expandNode(checkboxTreeItemRenderer, false);
   }

   /**
    * Expand the whole path to the node.
    */
   public void expandTo() {
      List<PrivilegeTreeNode> path = getPathToRoot(this);
      for (PrivilegeTreeNode node : path) {
         if(!node.expandNode()) {
            throw new RuntimeException(String.format("Failed to expand %s tree node!", node._displayName));
         }
      }
   }

   /**
    * Get the list of node children.
    * @return
    */
   protected List<PrivilegeTreeNode> getChildren() {
     return  new ArrayList<PrivilegeTreeNode>(this._children);
   }

   /**
    * Get the count of node children.
    * @return
    */
   protected int getChildrenCount() {
      return this._childrenCount;
   }

   /**
    * Add node child.
    * @param node
    */
   protected void addChild(PrivilegeTreeNode node) {
      this._children.add(node);
   }

   /**
    * Return list of nodes defining the path from the root to the node.
    * @param node
    * @return
    */
   private List<PrivilegeTreeNode> getPathToRoot(PrivilegeTreeNode node) {
      List<PrivilegeTreeNode> path = new ArrayList<PrivilegeTreeNode>();
      do {
         path.add(0, node);
         node = node._parent;
      } while(node != null);
      return path;
   }

   /**
    * Find node by node id. That id is the one provided from the data provider.
    * @param nodes
    * @param nodeId
    * @return
    */
   private PrivilegeTreeNode findNodeById(List<PrivilegeTreeNode> nodes, String nodeId) {
      for (PrivilegeTreeNode treeNode : nodes) {
         if(treeNode._id.equals(nodeId)) {
            return treeNode;
         }
         PrivilegeTreeNode node = findNodeById(treeNode.getChildren(), nodeId);
         if(node != null) {
            return node;
         }
      }
      return null;
   }

   /**
    * Find node by name - the displayed on the screen name.
    * @param nodes
    * @param nodeName
    * @return
    */
   private PrivilegeTreeNode findNodeByName(List<PrivilegeTreeNode> nodes, String nodeName) {
      for (PrivilegeTreeNode treeNode : nodes) {
         if(treeNode._displayName.equals(nodeName)) {
            return treeNode;
         }
         PrivilegeTreeNode node = findNodeById(treeNode.getChildren(), nodeName);
         if(node != null) {
            return node;
         }
      }
      return null;
   }

   private int parseInt(String integer) {
      try {
         return Integer.parseInt(integer);
      } catch (Exception e) {
         return 0;
      }
   }

   /**
    * Retrieve the checkbox control id.
    * @param treeItemRenderer
    * @return
    */
   private String getSelectionBoxId(UIComponent treeItemRenderer) {
      return getNodeItemsIds(treeItemRenderer, CHECKBOX_CONTROL_CLASS_NAME);
   }

   /**
    * Retrieve the expand/collapse button control id.
    * @param treeItemRenderer
    * @return
    */
   private String getExpandControlId(UIComponent treeItemRenderer) {
      return getNodeItemsIds(treeItemRenderer, EXPAND_COLLAPSE_CONTROL_CLASS_NAME);
   }

   /**
    * Retrieve the node item control children names.
    * @param treeItemRenderer
    * @return
    */
   private String getNodeItemsIds(UIComponent treeItemRenderer, String itemType) {
      String idString = treeItemRenderer.getProperties("[].name");
      String[] names = idString.split(",");
      for (String name : names) {
         if(name.contains(itemType)) {
            _logger.info("Node item id: name=" + name);
            return "name=" + name;
         }
      }
      throw new RuntimeException("Failed to find " + itemType);
   }

   /**
    * Find the UIComponent representing the node if expanded in the tree.
    * @return
    */
   private UIComponent findUiTreeItemRendered() {
      _logger.info("Find tree item: " + _displayName + " with uiId: " + _uiId);
      UIComponent checkboxTreeItemRenderer = new UIComponent(_uiId, BrowserUtil.flashSelenium);
      if(!checkboxTreeItemRenderer.isVisibleOnPath()) {
         VerticalScrollBar scrollBar = new VerticalScrollBar(_treeUid + SCROLL_BAR_ID);
         if(!scrollBar.isVisible()) {
            _logger.info("No schrollbar is shown");
            return null;
         }
         // Move up to top of the tree pane
         if(this._parent == null) {
            _logger.info("Scroll to top of the tree");
            scrollBar.goToTop();
         } else  {
            while(scrollBar.moveDown()) {
               if(checkboxTreeItemRenderer.isVisibleOnPath()) {
                  return checkboxTreeItemRenderer;
               }
            }
            _logger.warn("Node for " + _displayName + " is not visible for uiId " + _uiId);
            return null;
         }
      }
      return checkboxTreeItemRenderer;
   }

   private String buildNodeId() {
      String nodeIdFormat = "data.id=%s";
      if(this._parent == null) {
         nodeIdFormat = "data.id=%s[1]";
      }
      return String.format(nodeIdFormat, this._id);
   }

   /**
    * Expand/collapse tree node
    * @param checkboxTreeItemRenderer
    * @param expand
    * @return true if the node is expanded/collapsed
    */
   private boolean expandNode(final UIComponent checkboxTreeItemRenderer, boolean expand) {
      if(this.getChildren().size() == 0) {
         // It is a leaf - nothing to do
         _logger.info("Nothing to do the node is a leaf.");
         return true;
      }
      String expandControlId = getExpandControlId(checkboxTreeItemRenderer);
      final UIComponent expandControl = new UIComponent(expandControlId, BrowserUtil.flashSelenium);
      _logger.debug("Expand control id: " + expandControlId);
      String expectedState = EXPANDED_NODE_STATE_ID;
      String targetState = COLLPASED_NODE_STATE_ID;
      if(expand) {
         expectedState = COLLPASED_NODE_STATE_ID;
         targetState = EXPANDED_NODE_STATE_ID;
      }
      if(isStateExpected(expandControl, expectedState)) {
         // expand/collapse it
         expandControl.leftMouseClick();
         // wait for the expand/collapse control to get invisible after clicking on it.
         long timeout = Delay.timeout.forSeconds(5).getDuration();
         SUITA.Factory.UI_AUTOMATION_TOOL.condition.isTrue(
               new Object() {
                  public boolean equals(Object obj) {
                     return !expandControl.isVisibleOnPath();
                  }
               }).await(timeout);

         if(!expandControl.isVisibleOnPath()) {
            waitForNodeToExpand(checkboxTreeItemRenderer);
            expandControlId = getExpandControlId(checkboxTreeItemRenderer);
            _logger.debug("Expand control id after click: " + expandControlId);
            UIComponent collapsedControl = new UIComponent(expandControlId, BrowserUtil.flashSelenium);
            return isStateExpected(collapsedControl, targetState);
         }
         return false;
      } else if(isStateExpected(expandControl, targetState)) {
         _logger.info("Nothing to do the node is already " + (expand ? "expanded" : "colapsed"));
         return true;
      } else {
         throw new RuntimeException("Unable to detect tree node state!");
      }
   }

   /**
    * Validate the the expand state of the UI node is the same as the state provided by the expandedState parameter.
    * The expand/collapse state is detected by the image name of the expand control. If it is closed its name is
    * something similar to _class_embed_css_assets_images_closedRightArrow_png_254702973. respectively when it is
    * expanded is contains Open instead of closed.
    * @param expandControl
    * @param expandedState
    * @return
    */
   private boolean isStateExpected(UIComponent expandControl, String expandedState) {
      String buttonName = expandControl.getChildPropertyAtIndex("0", "name");
      return buttonName.contains(expandedState);
   }

   /**
    * Check/uncheck tree node.
    * @param checkboxTreeItemRenderer
    * @param check
    * @return
    */
   private boolean selectNode(UIComponent checkboxTreeItemRenderer, boolean check) {
      final UIComponent cbControl = new UIComponent(getSelectionBoxId(checkboxTreeItemRenderer), BrowserUtil.flashSelenium);
      String targetState = CB_STATE_UNCHECKED;
      if(check) {
         targetState = CB_STATE_CHECKED;
      }
      String initialState = cbControl.getProperty("state");
      if(initialState.equals(targetState)) {
         _logger.info("Nothing to it is already checked!");
         return true;
      } else {
         cbControl.leftMouseClick();
         String newState = cbControl.getProperty("state");
         return newState.equals(targetState);
      }
   }

   /**
    * After a node is expanded/collapsed the node item is redrawn several times. As a result the
    * expand/collapsed control id returned from first invocation of the getExpandControlId might not be the
    * final one and the usage of it will throw exception.
    * To be sure the that the node is visualized properly the automation does invocation of the getExpandControlId
    * till it return two consecutive same values.
    * @param checkboxTreeItemRenderer
    */
   private void waitForNodeToExpand(final UIComponent checkboxTreeItemRenderer) {
      Long timeout = Delay.timeout.forSeconds(5).getDuration();
      SUITA.Factory.UI_AUTOMATION_TOOL.condition.isTrue(
            new Object() {
               public boolean equals(Object obj) {
                  String collapseControlBeforeSleep = getExpandControlId(checkboxTreeItemRenderer);
                  // wait for half second
                  Delay.sleep.forMillis(1000).consume();
                  String collapseControlAfterSleep = getExpandControlId(checkboxTreeItemRenderer);
                  _logger.info("Colapse control before and after sleep: " + collapseControlBeforeSleep + " " + collapseControlAfterSleep);
                  return collapseControlBeforeSleep.equals(collapseControlAfterSleep);
               }
            }).await(timeout);
   }

   /**
    * Retrieve UI node selection state based.
    * @param checkboxTreeItemRenderer
    * @return
    */
   private String getSelectionState(UIComponent checkboxTreeItemRenderer) {
      final UIComponent cbControl = new UIComponent(getSelectionBoxId(checkboxTreeItemRenderer), BrowserUtil.flashSelenium);
      return cbControl.getProperty("state");
   }
}