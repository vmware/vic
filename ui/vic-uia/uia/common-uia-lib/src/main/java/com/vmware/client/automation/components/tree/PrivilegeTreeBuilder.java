/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.tree;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.vmware.client.automation.components.control.PropertiesUtil;
import com.vmware.suitaf.SUITA;

/**
 * Class that parse the whole privilege tree data provider data and based on it builds {@link PrivilegeTreeNode}
 * root node.
 */
public class PrivilegeTreeBuilder {

   private String _privilegeTreeId = "privGrid";

   // data id property
   private final String DATA_ID_PROPERTY = "id";
   private final String DATA_LABEL_PROPERTY = "displayName";
   private final String CHILDREN_LENGTH_PROPERTY = "children.length";

   // statically limit the tree to 4 levels.
   private int MAX_TREE_LEVEL = 4;

   /**
    * Create privilege tree builder. As tree id is used the default "privGrid" id.
    */
   public PrivilegeTreeBuilder() {
   }

   /**
    * Create privilege tree builder by specifying tree id.
    * @param treeId
    */
   public PrivilegeTreeBuilder(String treeId) {
      this._privilegeTreeId = treeId;
   }

   /**
    * Build and return the {@link PrivilegeTreeNode} root of the privilege tree.
    * By default is used "privGrid" as tree id.
    * @return
    */
   public PrivilegeTreeNode getRootNode() {
      return parseTree();
   }

   /**
    * Parse data provider data to the MAX_TREE_LEVEL and build the privilege tree.
    * @return
    */
   private PrivilegeTreeNode parseTree() {
       List<Map<String, String>> properties = new ArrayList<>();
       List<String> query = new ArrayList<>();
       for (int i = 0; i <= MAX_TREE_LEVEL; i++) {
           Map<String, String> map = getChildrenPropertiesOnLevel(i);
           properties.add(map);
           query.addAll(map.values());
       }

       PrivilegeTreeNode rootNode = buildHierarchy(MAX_TREE_LEVEL, properties, query);
       return rootNode;
   }

   private Map<String, String> getChildrenPropertiesOnLevel(int level) {
      Map<String, String> properties = new HashMap<>();
      properties.put(DATA_ID_PROPERTY, getPropertyOnLevel(DATA_ID_PROPERTY, level));
      properties.put(DATA_LABEL_PROPERTY, getPropertyOnLevel(DATA_LABEL_PROPERTY, level));
      properties.put(CHILDREN_LENGTH_PROPERTY, getPropertyOnLevel(CHILDREN_LENGTH_PROPERTY, level));
      return properties;
   }

   /**
    * Create tree hierarchy using the provider data limit by the max level.
    * @param treeMaxLevel
    * @param properties
    * @param query
    * @return
    */
   private PrivilegeTreeNode buildHierarchy(int treeMaxLevel,
         List<Map<String, String>> properties, List<String> query) {

      if (!SUITA.Factory.UI_AUTOMATION_TOOL.condition.isFound(_privilegeTreeId)
            .await(SUITA.Environment.getUIOperationTimeout()))
         return null;

      Map<String, String[]> treeContents =
            PropertiesUtil.getProperties(_privilegeTreeId, query.toArray(new String[0]));

      // Builder root node
      String idPropertyForRoot = properties.get(0).get(DATA_ID_PROPERTY);
      String labelPropertyForRoot = properties.get(0).get(DATA_LABEL_PROPERTY);
      String childrenCountPropertyForRoot = properties.get(0).get(CHILDREN_LENGTH_PROPERTY);
      PrivilegeTreeNode rootNode = new PrivilegeTreeNode(_privilegeTreeId,
            null, treeContents.get(idPropertyForRoot)[0],
            treeContents.get(labelPropertyForRoot)[0], treeContents.get(childrenCountPropertyForRoot)[0]);

      List<PrivilegeTreeNode> nodes = Arrays.asList(rootNode);
      for (int i = 1; i <= treeMaxLevel; i++) {
         buildNodes(nodes, properties.get(i), treeContents);
         nodes = replaceNodeWithItsChildren(nodes);
      }
      return rootNode;
   }

   /**
    * Build tree nodes for the respective level.
    * @param parentNodes
    * @param properties
    * @param treeDataProviderContents
    */
   private void buildNodes(List<PrivilegeTreeNode> parentNodes, Map<String, String> properties,
         Map<String, String[]> treeDataProviderContents) {

      String[] dataIds = treeDataProviderContents.get(properties.get(DATA_ID_PROPERTY));
      String[] labels = treeDataProviderContents.get(properties.get(DATA_LABEL_PROPERTY));
      String[] childrenLength =
            treeDataProviderContents.get(properties.get(CHILDREN_LENGTH_PROPERTY));

      int currentElementIndex = 0;

      for (PrivilegeTreeNode parentNode : parentNodes) {
         int childrenCount = parentNode.getChildrenCount();
         if (childrenCount == 0) {
            currentElementIndex++;
            continue;
         }

         for (int i = 0; i < childrenCount; i++) {
            PrivilegeTreeNode node =
                  new PrivilegeTreeNode(
                        _privilegeTreeId,
                        parentNode,
                        dataIds[currentElementIndex],
                        labels[currentElementIndex],
                        childrenLength[currentElementIndex]);
            parentNode.addChild(node);
            currentElementIndex++;
         }
      }
   }

   private List<PrivilegeTreeNode> replaceNodeWithItsChildren(List<PrivilegeTreeNode> nodes) {
      List<PrivilegeTreeNode> result = new ArrayList<>();
      for (PrivilegeTreeNode node : nodes) {
          if (node.getChildrenCount() > 0) {
              result.addAll(node.getChildren());
          } else {
              result.add(node);
          }
      }
      return result ;
   }

   private String getPropertyOnLevel(String property, int level) {
      StringBuilder result = new StringBuilder("dataProvider.source.[].");
      for (int i = 0; i < level; i++) {
         result.append("children.source.[].");
      }

      result.append(property);
      return result.toString();
   }
}