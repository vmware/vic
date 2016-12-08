/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.menu;

import java.util.ArrayList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;

/**
 * Class to represent the menu node hierarchical structure.
 */
public class MenuNode extends UIComponent {

   private static final Logger _logger = LoggerFactory.getLogger(MenuNode.class);

   protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   // Timeout for expand menu timeout
   private static final long EXPAND_TIMEOUT = 3000;

   private final MenuNode parentNode;
   private final String actionId;
   private final String label;
   private final String type;
   private final int childrenCount;
   private final List<MenuNode> children;

   /**
    * Create new instance of Menu node
    *
    * @param parent The nodes parent element
    * @param actionId The action id of menu node
    * @param label
    * @param type
    * @param childrenCount
    */
   public MenuNode(MenuNode parent, String actionId, String label, String type,
         String childrenCount) {
      super(buildId(parent, actionId), BrowserUtil.flashSelenium);
      this.parentNode = parent;
      this.actionId = actionId;
      this.label = label;
      this.type = type;
      this.children = new ArrayList<MenuNode>();
      this.childrenCount = parseInt(childrenCount);
   }

   /**
    * Get child of this node with a specific action id
    *
    * @param actionId
    * @return
    */
   public MenuNode getChildByActionId(String actionId) {
      MenuNode directChild = getDirectChildByActionId(actionId);
      if (directChild != null) {
         return directChild;
      }

      for (MenuNode child : children) {
         MenuNode menuNode = child.getChildByActionId(actionId);
         if (menuNode != null) {
            return menuNode;
         }
      }

      return null;
   }

   /**
    * Get child of this node with a specific label
    *
    * @param actionId
    * @return
    */
   public MenuNode getChildByActionLabel(String actionLabel) {
      MenuNode directChild = getDirectChildByActionLabel(actionLabel);
      if (directChild != null) {
         return directChild;
      }

      for (MenuNode child : children) {
         MenuNode menuNode = child.getChildByActionLabel(actionLabel);
         if (menuNode != null) {
            return menuNode;
         }
      }

      return null;
   }

   private MenuNode getDirectChildByActionId(String actionId) {
      for (MenuNode child : children) {
         if (child.actionId.equals(actionId)) {
            return child;
         }
      }

      return null;
   }

   private MenuNode getDirectChildByActionLabel(String actionLabel) {
      for (MenuNode child : children) {
         if (child.label.equals(actionLabel)) {
            return child;
         }
      }

      return null;
   }

   /**
    * Expands the menu hierarchy to current node
    *
    * @throws Exception
    */
   public void expandTo() throws Exception {
      expandTo(EXPAND_TIMEOUT);
   }

   /**
    * Expands the menu hierarchy to current node
    *
    * @param timeout
    * @throws Exception
    */
   public void expandTo(long timeout) throws Exception {
      if (parentNode != null) {
         parentNode.expandTo(timeout);
      }

      _logger.info(
            String.format(
                  "Navigating to menu node with label: %s and actionId: %s",
                  label,
                  actionId
               )
         );

      waitForElementExistence(timeout);
      visibleMouseOver();
   }

   /**
    * Generates direct menu action ID.
    * The method traverses the three and collects all parent IDs.
    */
   public String getDirectId() {
      if (parentNode != null) {
         return parentNode.getDirectId() + "." +  actionId;
      }

      return actionId;
   }

   /**
    * Get node's actionId
    *
    * @return
    */
   public String getActionId() {
      return actionId;
   }

   /**
    * Get node's parent. Null if this is the root
    *
    * @return
    */
   public MenuNode getParentNode() {
      return parentNode;
   }

   /**
    * Get node's type
    *
    * @return
    */
   public String getType() {
      return type;
   }

   /**
    * Get amount of children for the node
    *
    * @return
    */
   public int getChildrenCount() {
      return childrenCount;
   }

   /**
    * Add child to this node's children collection.
    *
    * @param menuNode
    */
   public void addChild(MenuNode menuNode) {
      children.add(menuNode);

   }

   /**
    * Get this node's children
    *
    * @return
    */
   public List<MenuNode> getChildren() {
      return new ArrayList<MenuNode>(children);
   }

   @Override
   public String toString() {
      StringBuilder childsStringBuilder = new StringBuilder();
      childsStringBuilder.append('[');
      if (children != null && children.size() > 0) {
         for (MenuNode child : children) {
            childsStringBuilder.append(child.toString());
            childsStringBuilder.append(',');
         }
         childsStringBuilder.deleteCharAt(childsStringBuilder.length() - 1);
      }
      childsStringBuilder.append(']');

      return String.format(
            "{uid: %S , type: %s, label: %s, children count: %s, children: %s}",
            actionId,
            type,
            label,
            childrenCount,
            childsStringBuilder.toString());
   }

   private void waitForElementExistence(long timeoutDuration) throws Exception {
      final long startTime = System.currentTimeMillis();

      _logger.info("Waiting for element existence: " + getUniqueId());

      long elapsedTime = 0;
      boolean notCompleted = true;
      while (notCompleted && elapsedTime < timeoutDuration) {
         // TODO rreymer: use SUITA tools instead
         Thread.sleep(50);
         notCompleted = !isComponentExisting();
         elapsedTime = System.currentTimeMillis() - startTime;
      }
   }

   private static String buildId(MenuNode parent, String id) {
      if (parent == null) {
         return id;
      }

      return parent.uniqueId + '.' + id;
   }

   private int parseInt(String integer) {
      try {
         return Integer.parseInt(integer);
      } catch (Exception e) {
         return 0;
      }
   }
}
