/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.navigator.navstep;

import org.openqa.selenium.NotFoundException;

import com.google.common.base.Strings;
import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.client.automation.util.UiDelay;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.common.custom.ButtonScrollingButtonBar;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.SubToolAudit;
import com.vmware.suitaf.apl.IDGroup;

/**
 * A <code>NavigationStep</code> for selecting a secondary tab such as Settings.
 */
public class SecondaryTabNavStep extends BaseNavStep {

   // Used by Manage and Monitor secondary tabs.
   private static String ID_TOGGLE_BUTTON_BAR = "toggleButtonBar";
   private static final String MANAGE_PATH_ID = "manage";

   // Used by Related Objects secondary tab.
   private static String ID_TAB_BAR = "className=RelatedItemsView/tabBar";

   private int _tabPosition = -1;
   private String _tabName = "";

   /**
    * The constructor defines a mapping between the navigation identifier and
    * the tab position.
    *
    * @param nId
    * @see <code>NavigationStep</code>.
    *
    * @param tabPosition
    *           Zero based position of the tab.
    */
   @Deprecated
   public SecondaryTabNavStep(String nid, int tabPosition) {
      super(nid);

      _tabPosition = tabPosition;
   }

   /**
    * The constructor defines a mapping between the navigation identifier and
    * the tab name.
    *
    * @param nId
    * @see <code>NavigationStep</code>.
    *
    * @param tabName
    *           Tab name.
    */
   public SecondaryTabNavStep(String nid, String tabName) {
      super(nid);

      _tabName = tabName;
   }

   @Override
   public void doNavigate(LocationSpec locationSpec) throws Exception {
      // Check which secondary tab bar is used.
      String subTabNavigatorId = ID_TOGGLE_BUTTON_BAR;
      if (!UI.component.exists(IDGroup.toIDGroup(subTabNavigatorId))) {
         subTabNavigatorId = ID_TAB_BAR;

         if (!UI.component.exists(IDGroup.toIDGroup(subTabNavigatorId))) {
            if (locationSpec.path.stringValue().contains(MANAGE_PATH_ID)) {
               // Try to see if we use flat tabs
               if (!Strings.isNullOrEmpty(_tabName)) {
                  // Navigate by name
                  UIComponent treeItem = TocNavStep.findUiTreeItem(_tabName);
                  if (null != treeItem) {
                     treeItem.leftMouseClick();
                  }
               }
               return;
            } else {
               throw new NotFoundException("Secondary tab object not found.");
            }
         }
      }

      // Init the ButtonScrollingButtonBar.
      ButtonScrollingButtonBar subTabsBar = new ButtonScrollingButtonBar(
            subTabNavigatorId, BrowserUtil.flashSelenium);

      // When we navigate to different secondary tab of the same object
      // we have different id of the toggleButtonBar
      // Example: first time we go to VC -> Monitor and the id is "toggleButtonBar"
      // second time we go to VC -> Manage the id is "toggleButtonBar[1]"
      // Below condition is a workaround for this issue
      // TODO: think of a better solution to the problem
      if (!subTabsBar.isVisibleOnPath()) {
         int retryIndex = 1;
         _logger.info("Element with id " + ID_TOGGLE_BUTTON_BAR
            + " is not visible, re-trying with index " + retryIndex);
         subTabsBar = new ButtonScrollingButtonBar(subTabNavigatorId + "["
            + retryIndex + "]", BrowserUtil.flashSelenium);
      }

      subTabsBar
            .waitForElementEnable(SUITA.Environment.getUIOperationTimeout());

      SUITA.Factory.UI_AUTOMATION_TOOL.audit.snapshotAppScreen(
            SubToolAudit.getFPID(), "PRE_CLICK_SECONDARY_TAB_INDEX_"
                  + _tabPosition);

      if (!Strings.isNullOrEmpty(_tabName)) {
         // Navigate by name.
         int tabIndex = subTabsBar.getButtonIndexByName(_tabName);
         if(tabIndex >= 0){
            subTabsBar.clickScrollingButtonBarItemAtIndex(String
                  .valueOf(tabIndex));
         }
         else{
            throw new NotFoundException("Tab with name "+ _tabName +" not found.");
         }
      } else if (_tabPosition >= 0) {
         // Navigate by position.
         subTabsBar.clickScrollingButtonBarItemAtIndex(String
               .valueOf(_tabPosition));

      } else {
         throw new IllegalArgumentException(
               "tabPosition and tabName cannot be both undefined at the same time.");
      }
   }
}
