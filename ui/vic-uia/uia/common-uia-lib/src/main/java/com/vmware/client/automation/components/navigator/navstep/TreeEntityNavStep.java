/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.navigator.navstep;

import java.util.List;

import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.flexui.componentframework.controls.mx.Tree.TreeNode;
import com.vmware.flexui.componentframework.controls.mx.custom.InventoryTree;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.vsphere.client.automation.components.tree.spec.TreeLocationSpec;

/**
 * The <code>TreeEntityNavStep</code> selects finds and focuses the node we want
 * using its name.
 *
 * The action can be used from any page that has the inventory tree visible.
 */
public class TreeEntityNavStep extends BaseNavStep {

   protected static InventoryTree tree;
   private static final String TREE_ID = "navTree";
   private static final String DATA_ID_PROPERTY = "data.id";

   public TreeEntityNavStep(String entityName) {
      super(entityName);
   }

   @Override
   public void doNavigate(LocationSpec locationSpec) throws Exception {
      boolean isNavigationDone = false;
      tree = new InventoryTree(TREE_ID, BrowserUtil.flashSelenium);
      tree.waitForNodesToFinishLoading();
      List<TreeNode> vcNodes = tree.getTreeContent();
      TreeLocationSpec treeLocationSpec = (TreeLocationSpec) locationSpec;
      for (TreeNode vcNode : vcNodes) {
         isNavigationDone = expandUntilChildNodeIsFound(vcNode.getNodeText(),
               getNId(), getEntityTypeFromNode(vcNode),
               treeLocationSpec.getEntityType());
         if (isNavigationDone) {
            break;
         }
      }
      focusNode(getNId());
   }

   /**
    * Expands tree nodes under the passed node until we find the node we want.
    * We compare nodes using the name and entity type.
    *
    * @param parentName
    *           the name of the parent node we want to look under
    * @param wantedNode
    *           the node we want to find
    * @param entityType
    *           the entity type of the parent node
    * @param wantedType
    *           the entity type we are looking for
    * @return True if we've reached the desired node, False otherwise
    */
   private static boolean expandUntilChildNodeIsFound(String parentName,
         String wantedNode, String entityType, String wantedType) {
      boolean isNodeFound = false;
      if (!parentName.equals(wantedNode) || !entityType.equals(wantedType)) {
         TreeNode parent = refreshTreeNode(parentName);
         if (parent != null) {
            List<TreeNode> children = parent.getChildren();
            for (int i = 0; i < children.size(); i++) {
               parent = refreshTreeNode(parentName);
               TreeNode child = parent.getChildren().get(i);
               String nodeId = child.getProperty(DATA_ID_PROPERTY);
               String nodeType = getEntityTypeFromNode(child);
               child.getNodeText();
               tree.expandMultiViewInvNode(nodeId);
               tree.waitForNodesToFinishLoading();
               isNodeFound = expandUntilChildNodeIsFound(child.getNodeText(),
                     wantedNode, nodeType, wantedType);
               if (isNodeFound) {
                  break;
               }
            }
         }
      } else {
         isNodeFound = true;
      }
      return isNodeFound;
   }

   /**
    * Retrieve the desired tree node from a list of all visible nodes
    *
    * @param visibleNodes
    *           a list of visible <code>TreeNode</code> elements
    * @param nodeName
    *           the node name
    * @return TreeNode the desired <code>TreeNode</code>
    */
   private static TreeNode getTreeNode(List<TreeNode> visibleNodes,
         String nodeName) {
      TreeNode resultNode = null;
      for (TreeNode childNode : visibleNodes) {
         if (childNode.getNodeText().equals(nodeName)) {
            resultNode = childNode;
            break;
         } else {
            resultNode = getTreeNode(childNode.getChildren(), nodeName);
            if (resultNode != null) {
               break;
            }
         }
      }
      return resultNode;
   }

   /**
    * Focus into a node from the inventory tree by node name
    *
    * @param nodeName
    *           the node name
    */
   public static void focusNode(String nodeName) {
      TreeNode node = refreshTreeNode(nodeName);
      tree.focusIntoInvNode(node.getProperty(DATA_ID_PROPERTY));
   }

   /**
    * Get the entity type from the parsed data of a tree node
    *
    * @param node
    *           the node whose type we want
    * @return A String representation of the entity type
    */
   private static String getEntityTypeFromNode(TreeNode node) {
      return node.getProperty(DATA_ID_PROPERTY).split(":")[1];
   }

   /**
    * Refreshes the node so that it reflects the current state of the tree
    *
    * @param nodeName
    *           the node name
    * @return The refreshed <code>TreeNode</code>, containing accurate data
    */
   private static TreeNode refreshTreeNode(String nodeName) {
      List<TreeNode> visibleNodes = new InventoryTree(TREE_ID,
            BrowserUtil.flashSelenium).getTreeContent();
      return getTreeNode(visibleNodes, nodeName);
   }
}