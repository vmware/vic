/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.menu;

import static com.vmware.flexui.selenium.BrowserUtil.flashSelenium;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.apache.commons.lang.StringUtils;

import com.vmware.client.automation.vcuilib.commoncode.IDConstants;
import com.vmware.flexui.componentframework.UIComponent;

/**
 * CLass used to encapsulate the parsing and building the context menu hierarchy
 */
public class ContextMenuBuilder {

   private final String UID_PROPERTY = "uid";
   private final String LABEL_PROPERTY = "label";
   private final String TYPE_PROPERTY = "type";
   private final String CHILDREN_PROPERTY = "children";
   private final String CHILDREN_LENGTH_PROPERTY = "children.length";

   private static final String ARRAY_DELIMITER = ";;;";
   private static final String LEVELS_DELIMITER = "!!!";
   private static final String PROPERTIES_QUERY_DELIMITER = ",";
   private static final String KEY_VALUE_DELIMITER = ":";
   private static final String CHILD_ELEMENTS_DELIMITER = ",";
   private static final String ELEMENT_ROWS_DELIMITER = "\n";

   /**
    * Build and return the menu's root element
    *
    * @return
    */
   public MenuNode getRootMenuNode() {
      return parseMenu();
   }

   private MenuNode parseMenu() {
      int MAX_LEVEL = 2;
      // List containing queries to get menu elements on each level
      List<String> query = new ArrayList<>();
      for (int i = 0; i <= MAX_LEVEL; i++) {
         query.add(getQueryOnLevel(i));
      }

      MenuNode rootNode = buildMenuHierarchy(MAX_LEVEL, query);

      return rootNode;
   }

   private MenuNode buildMenuHierarchy(int MAX_LEVEL, List<String> query) {
      //Get the menu elements for every level of the context menu
      Map<Integer, String[]> contextextMenuContents = ContextMenuBuilder
            .getProperties(IDConstants.ID_CONTEXT_MENU,
                  query.toArray(new String[0]));

      int numChildrenRoot = contextextMenuContents.get(0).length;
      MenuNode rootNode = new MenuNode(null, IDConstants.ID_CONTEXT_MENU,
            "root", "root", Integer.toString(numChildrenRoot));

      List<MenuNode> nodes = Arrays.asList(rootNode);
      for (int i = 0; i <= MAX_LEVEL; i++) {
         buildNodes(nodes, contextextMenuContents.get(i));
         nodes = replaceNodeWithItsChildren(nodes);
      }
      return rootNode;
   }

   private Map<String, String> parseMenuElement(String menuElement) {
      Map<String, String> menuElementMap = new HashMap<String, String>();

      String[] pairs = menuElement.split(ELEMENT_ROWS_DELIMITER);
      for (int i = 0; i < pairs.length; i++) {
         String pair = pairs[i];
         String[] keyValue = pair.split(KEY_VALUE_DELIMITER, -1);
         menuElementMap
               .put(keyValue[0], keyValue[1].equals("null")?null:keyValue[1]);
      }

      if (menuElementMap.get(CHILDREN_PROPERTY) != null) {
         menuElementMap.put(
               CHILDREN_LENGTH_PROPERTY,
               Integer.toString(menuElementMap.get(CHILDREN_PROPERTY).split(
                     CHILD_ELEMENTS_DELIMITER).length));
      }
      return menuElementMap;
   }

   private void buildNodes(List<MenuNode> parentNodes,
         String[] contextextMenuContents) {

      int currentElementIndex = 0;

      for (MenuNode parentNode : parentNodes) {
         int childrenCount = parentNode.getChildrenCount();
         if (childrenCount == 0) {
            currentElementIndex++;
            continue;
         }

         for (int i = 0; i < childrenCount; i++) {
            Map<String, String> element = parseMenuElement(contextextMenuContents[currentElementIndex]);
            MenuNode node = new MenuNode(parentNode, element.get(UID_PROPERTY),
                  element.get(LABEL_PROPERTY), element.get(TYPE_PROPERTY),
                  element.get(CHILDREN_LENGTH_PROPERTY));
            parentNode.addChild(node);
            currentElementIndex++;
         }
      }
   }

   private List<MenuNode> replaceNodeWithItsChildren(List<MenuNode> nodes) {
      List<MenuNode> result = new ArrayList<>();
      for (MenuNode menuNode : nodes) {
         if (menuNode.getChildrenCount() > 0) {
            result.addAll(menuNode.getChildren());
         } else {
            result.add(menuNode);
         }
      }
      return result;
   }

   private String getQueryOnLevel(int level) {
      StringBuilder result = new StringBuilder("dataProvider.source.[].");
      for (int i = 0; i < level; i++) {
         result.append("children.[].");
      }

      result.append("{all}");
      return result.toString();
   }

   /**
    * Get properties for component with specific ID.
    *
    * NOTE: this method was taken from VCUI-QE-LIB: PropertiesUtil.
    *
    * @param componentId
    * @param queries
    * @return Map<LevelName, values> - if the the levels has single value, the
    *         array will have one element
    */
   private static Map<Integer, String[]> getProperties(String componentId,
         String... queries) {
      UIComponent component = new UIComponent(componentId, flashSelenium);
      String query = StringUtils.join(queries, PROPERTIES_QUERY_DELIMITER);

      String rawResult = component.getProperties(query, ARRAY_DELIMITER,
            LEVELS_DELIMITER);
      String[] levelValues = rawResult.split(LEVELS_DELIMITER, -1);
      if (levelValues.length != queries.length) {
         return null;
      }

      Map<Integer, String[]> result = new HashMap<Integer, String[]>();
      for (int i = 0; i < queries.length; i++) {
         result.put(i, levelValues[i].split(ARRAY_DELIMITER, -1));
      }

      return result;
   }
}
