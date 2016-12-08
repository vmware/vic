/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import com.google.common.base.Strings;
import com.vmware.client.automation.exception.TocTreeException;
import com.vmware.client.automation.util.UIUtil;
import com.vmware.flexui.componentframework.controls.mx.List;
import com.vmware.flexui.componentframework.controls.mx.Tree;
import com.vmware.flexui.selenium.BrowserUtil;

/**
 * TocTree utility methods.
 */
public final class TocTreeUtil {

   /** This index is returned when Toc Item is not found. */
   private static final int ITEM_NOT_FOUND_INDEX = -1;
   private static final String PROPERTY_VISIBLE = "visible";
   private static final String PROPERTY_AUTOMATION_NAME = "automationName";
   private static final String PROPERTY_CLASS_NAME = "className";
   private static final String CLASS_TOC_TREE_ITEM = "TocTreeItemRenderer";
   private static final String CLASS_LIST_BASE_CONTENT = "ListBaseContentHolder";

   /**
    * Loops through all TOC tree items in search for item with specific name.
    * TODO mdzhokanov: optimize this method to be faster
    *
    * @param tocItemName
    *           the name of the target Toc item
    * @return
    * @throws TocTreeException
    *            in case Toc item with such name is not found
    */
   public static int getTocItemIndexByName(Tree tocTree, String tocItemName) {
      // Traverse the TocTree in search of "ListBaseContentHolder" component
      List listBaseContent = getTreeContentHolder(tocTree);

      // Sanity check
      if (listBaseContent == null) {
         throw new TocTreeException("Could not find ListBase child of TocTree '"
               + tocTree.getUniqueId() + "'");
      }

      // Loop the listBaseContent in search of specific Toc Item
      int tocItemIndex =
            TocTreeUtil.findTocItemIndexByName(listBaseContent, tocItemName);

      // Sanity check
      if (tocItemIndex == -1) {
         throw new TocTreeException("Could not find Toc Item with such name '"
               + tocItemName + "' in Toc Tree '" + tocTree.getUniqueId() + "'");
      }

      return tocItemIndex;
   }

   /**
    * Finds the relevant <code>ListBaseContentHolder</code> component. Each Toc Tree
    * includes such content holder.
    *
    * NOTE: return null in case the content holder is not found.
    */
   private static List getTreeContentHolder(Tree tocTree) {
      int childrenCount = tocTree.getNumChildren();
      // Loop through all tocTree ui children in search of ListBaseContentHolder
      for (int i = 0; i < childrenCount; i++) {
         String childClassName =
               TocTreeUtil.getListChildProperty(tocTree, i, PROPERTY_CLASS_NAME);
         if (childClassName.equals(CLASS_LIST_BASE_CONTENT)) {
            String id = tocTree.getChildIdAtIndex(String.valueOf(i));
            return new List(id, BrowserUtil.flashSelenium);
         }
      }

      // ListBaseContentHolder child is not found
      return null;
   }

   /**
    * Loops through <link>ListBaseContentHolder</link> children in search for TocTreeItem
    * that has specific name. When the item is found current method returns its index.
    *
    * @param list
    *           reference to <link>ListBaseContentHolder</link> component
    * @param tocItemName
    *           name of the target toc item
    * @return index of the item in the toc tree
    */
   public static int findTocItemIndexByName(List list, String tocItemName) {
      // Sanity check
      if (Strings.isNullOrEmpty(tocItemName)) {
         return ITEM_NOT_FOUND_INDEX;
      }

      int childrenCount = list.getNumChildren();
      int itemsCounter = 0;

      // Loop through all items in the List
      for (int i = 0; i < childrenCount; i++) {
         String childClassName = getListChildProperty(list, i, PROPERTY_CLASS_NAME);
         String visible = getListChildProperty(list, i, PROPERTY_VISIBLE);

         // When TocTreeItem is found (hidden items are ignored), then check its name
         if (childClassName.equals(CLASS_TOC_TREE_ITEM) && visible.equals("true")) {
            String automationName =
                  getListChildProperty(list, i, PROPERTY_AUTOMATION_NAME);
            if (automationName.equals(tocItemName)) {
               return itemsCounter;
            }
            // increment items counter (it is incremented only on TocTreeItems visit)
            itemsCounter++;
         }
      }

      // item not found
      return ITEM_NOT_FOUND_INDEX;
   }

   /**
    * Returns specific property for a
    * <link>com.vmware.flexui.componentframework.controls.mx.List</link> child.
    *
    * @param list
    *           target list
    * @param childIndex
    *           sequence number of the target child
    * @param propertyKey
    *           key for the target property
    * @return property value
    */
   public static String getListChildProperty(List list, int childIndex,
         String propertyKey) {
      String propertyRaw =
            list.getChildPropertyAtIndex(String.valueOf(childIndex), propertyKey);
      return UIUtil.getPropertyValue(propertyRaw);
   }
}
