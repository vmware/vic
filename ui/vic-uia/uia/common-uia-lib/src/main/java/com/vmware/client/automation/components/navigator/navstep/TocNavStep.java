/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.navigator.navstep;

import org.openqa.selenium.NotFoundException;

import com.google.common.base.Strings;
import com.vmware.client.automation.components.control.TocTreeControl;
import com.vmware.client.automation.components.control.VerticalScrollBar;
import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.client.automation.util.UiDelay;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.SubToolAudit;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;

/**
 * A <code>NavigationStep</code> for selecting a table of content item.
 */
public class TocNavStep extends BaseNavStep {

   private static final String DEFAULT_TOC_TREE_ID = "tocTree";
   private static final String SCROLL_BAR_ID = "/className=VScrollBar";

   private int _tocItemPosition = -1;
   private String _tocItemName = "";

   /**
    * The constructor defines a mapping between the navigation identifier and
    * the ToC item position.
    *
    * @param nId
    * @see <code>NavigationStep</code>.
    *
    * @param tocItemPosition
    *           Zero based position of the ToC item.
    */
   @Deprecated
   public TocNavStep(String nid, int tocItemPosition) {
      super(nid);

      _tocItemPosition = tocItemPosition;
   }

   /**
    * The constructor defines a mapping between the navigation identifier and
    * the ToC item name.
    *
    * @param nId
    * @see <code>NavigationStep</code>.
    *
    * @param tocItemName
    *           ToC item name.
    */
   public TocNavStep(String nid, String tocItemName) {
      super(nid);

      _tocItemName = tocItemName;
   }

   @Override
   public void doNavigate(LocationSpec locationSpec) throws Exception {
      // Take ToCTree UID. The visible one has priority.
      String uid = UI.component.property.get(Property.DIRECT_ID,
            IDGroup.toIDGroup(DEFAULT_TOC_TREE_ID));

      if (Strings.isNullOrEmpty(uid)) {
         throw new IllegalStateException(String.format(
               "The ToC tree cannot be found on the screen: %s",
               DEFAULT_TOC_TREE_ID));
      }

      TocTreeControl tocTree = new TocTreeControl(uid);

      SUITA.Factory.UI_AUTOMATION_TOOL.audit.snapshotAppScreen(
            SubToolAudit.getFPID(), "PRE_CLICK_TOC_ITEM_INDEX_"
                  + _tocItemPosition);
      if (_tocItemPosition >= 0 && Strings.isNullOrEmpty(_tocItemName)) {

         // Navigate by position.
         tocTree.selectItemByIndex(_tocItemPosition);

      } else if (_tocItemPosition < 0 && !Strings.isNullOrEmpty(_tocItemName)) {
         // Navigate by name.
         UIComponent treeItem = findUiTreeItem(_tocItemName);
         if (null != treeItem) {
            treeItem.leftMouseClick();
         } else {
            throw new NotFoundException("Toc item with name " + treeItem
                  + " was not found.");
         }
      } else {
         throw new IllegalArgumentException(
               "tocItemPosition and tocItemName cannot be both set at the same time.");
      }
   }

   /**
    * Searches for Toc item name within the view. If not found but scroller is
    * found the method will try to scroll from top to down and search for the
    * toc item.
    *
    * @return Returns UIComponents in case tocItem is visible and null otherwise
    */
   public static UIComponent findUiTreeItem(String tocItemName)
         throws Exception {
      String uiId = "automationName=" + tocItemName;
      _logger.info("Search for tree item: " + tocItemName + " with uiId: "
            + uiId);

      UIComponent treeItem = new UIComponent(uiId, BrowserUtil.flashSelenium);
      if (UI.condition.isFound(uiId).await(
            UiDelay.UI_OPERATION_TIMEOUT.getDuration() / 3)) {
         _logger.info("Found tree item: " + tocItemName);
         return treeItem;
      }

      // In case we didn't find toc tree item visible - search for scrollbar
      // Then scroll from top to bottom and search again
      String scrollBarId = DEFAULT_TOC_TREE_ID + SCROLL_BAR_ID;
      if (!UI.condition.isFound(scrollBarId).await(
            UiDelay.UI_OPERATION_TIMEOUT.getDuration() / 3)) {
         _logger.info("No scrollbar is shown and the searched item "
               + tocItemName + " is not visible!");
         return null;
      }

      VerticalScrollBar scrollBar = new VerticalScrollBar(scrollBarId);

      // Move up to top of the tree pane
      _logger.info("Scroll to top of the tree");
      scrollBar.goToTop();
      do {
         if (UI.condition.isFound(uiId).await(
               UiDelay.UI_OPERATION_TIMEOUT.getDuration() / 3)) {
            _logger.info("Found tree item: " + tocItemName);
            return treeItem;
         }
      } while (scrollBar.moveDown());

      _logger.info("Can not Find tree item: " + tocItemName);
      return null;
   }
}