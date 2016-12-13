/**
 * Copyright 2012 VMware, Inc.  All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.components.navigator;

import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.LABEL_APP_LOGOUT;
import static com.vmware.flexui.selenium.BrowserUtil.flashSelenium;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.thoughtworks.selenium.FlashSelenium;
import com.vmware.client.automation.components.menu.ContextMenuBuilder;
import com.vmware.client.automation.components.menu.MenuNode;
import com.vmware.client.automation.vcuilib.commoncode.IDConstants;
import com.vmware.client.automation.vcuilib.commoncode.TestBaseUI;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.mx.Button;
import com.vmware.flexui.componentframework.controls.mx.Menu;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.SubToolAudit;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;

/**
 * This class provided ability to work with "Action" menu. That includes the
 * main action menu, toolbar menu, object navigator context menu and data
 * grid context menu.
 *
 * There are two type of methods: the first type expect that the menu is open
 * and the second type open the action menu and start to work with it. The first
 * type of methods is planned to be used with context menus for which the
 * automation need to execute mouse right-click on an object.
 * The second type of methods are for main and toolbar menus.
 *
 * @example ActionNavigator.invokeFromToolbarMenu(OrgPage.MI_CREATE_VDC); -
 *          will open the toolbar action menu and will click on the create new
 *          vDC menu item.
 *
 *          ActionNavigator.invokeFromActionsMenu(OrgPage.MI_ENABLE_ORG); -
 *          will open main action menu, go to the any sub menu in needed and
 *          click on the enable organization menu item.
 *
 *          ActionNavigator.invokeMenuAction(IDGroup.toIDGroup(
 *          "vsphere.core.vm.powerOnAction")); -
 *          expect that the menu is open, go to the designated sub menu item,
 *          and will click on the "Power On" action.
 *
 */
public class ActionNavigator {

   private static final Logger _logger = LoggerFactory.getLogger(ActionNavigator.class);

   private static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   // IDs of the sub menus.
   public static final IDGroup MI_ALL_VCD = IDGroup.toIDGroup("afContextMenu.vCloud");
   public static final IDGroup MI_ALL_VC = IDGroup.toIDGroup("afContextMenu.vCenter");
   public static final IDGroup MI_POWER = IDGroup.toIDGroup("afContextMenu.power");
   public static final IDGroup MI_ALARMS = IDGroup.toIDGroup("afContextMenu.alarms");
   public static final String ID_LOGOUT_MENU_ITEM = "Logout";
   private static final String ID_USER_MENU = "userMenu";

   // Action Navigator timeout constants.
   // TODO rreymer: find a better way managing such constants
   private static final long AN_TIMEOUT_ONE_MINUTE = 60000;
   private static final long AN_TIMEOUT_ONE_SECOND = 1000;

   /**
    * Opens the "Actions" menu on the top of the main content pane and clicks on
    * the menu item specified by the ID provided with the menuItemID parameter.
    * The method will open the sub menus specified by the subMenuIDs parameter.
    *
    * @param menuItemID    id of the menu item to be clicked.
    * @param subMenuIDs    path to the sub menu containing the menu item to click.
    *
    * @deprecated          use {@link #invokeFromActionsMenu(IDGroup)}
    */
   @Deprecated
   public static void invokeFromActionsMenu(IDGroup menuItemID, IDGroup... subMenuIDs) {
      ActionNavigator.invokeFromActionsMenu(menuItemID);
   }

   /**
    * Opens the "Actions" menu on the top of the main content pane and clicks on
    * the menu item specified by the ID provided with the menuItemID parameter.
    * The method will open the sub menus specified by the subMenuIDs parameter.
    *
    * @param menuItemID    id of the menu item to be clicked.
    * @param subMenuIDs    path to the sub menu containing the menu item to click.
    */
   public static void invokeFromActionsMenu(IDGroup menuItemID) {
      try {
         ActionNavigator.openMoreActions();
      } catch (Exception e) {
         SUITA.Factory.UI_AUTOMATION_TOOL.assertFor.isTrue(
               "Opening actions menu throw exception: " + e.getMessage(),
               false
               );
      }

      ActionNavigator.invokeMenuAction(menuItemID);
   }

   /**
    * Opens the data grid "Actions" menu on the top of the main content pane and
    * clicks on the menu item specified by the ID provided with the menuItemID
    * parameter. The method will open the sub menus specified by the subMenuIDs
    * parameter.
    *
    * @param menuItemID    id of the menu item to be clicked.
    * @param subMenuIDs    path to the sub menu containing the menu item to click.
    */
   public static void invokeFromToolbarMenu(IDGroup menuItemID) {
      ActionNavigator.invokeDebugFromToolbarMenu(false, menuItemID);
   }

   /**
    * Opens the data grid "Actions" menu on the top of the main content pane and
    * clicks on the menu item specified by the ID provided with the menuItemID
    * parameter. The method will open the sub menus specified by the subMenuIDs
    * parameter and will capture screenshots of the menu actions
    *
    * @param menuItemID    id of the menu item to be clicked.
    */
   public static void invokeDebugFromToolbarMenu(IDGroup menuItemID) {
      ActionNavigator.invokeDebugFromToolbarMenu(true, menuItemID);
   }

   /**
    * Opens a specific data grid "Actions" menu on the top of the main content
    * pane and clicks on the menu item specified by the ID provided with the
    * menuItemID parameter. The method will open the sub menus specified by the
    * subMenuIDs parameter.
    *
    * @param gridActionsButton   specifies the view in which the 'Actions' button is
    * @param menuItemID          id of the menu item to be clicked
    * @param subMenuIDs          path to the sub menu containing the menu item to click
    *
    * @deprecated                use {@link #invokeFromToolbarMenu(IDGroup)}
    */
   @Deprecated
   public static void invokeFromToolbarMenu(String actionsButtonId,
         IDGroup menuItemID, IDGroup... subMenuIDs) {
      ActionNavigator.invokeFromToolbarMenu(actionsButtonId, menuItemID);
   }

   /**
    * Opens a specific data grid "Actions" menu on the top of the main content
    * pane and clicks on the menu item specified by the ID provided with the
    * menuItemID parameter. The method will open the sub menus specified by the
    * subMenuIDs parameter.
    *
    * @param gridActionsButton   specifies the view in which the 'Actions' button is
    * @param menuItemID          id of the menu item to be clicked
    */
   public static void invokeFromToolbarMenu(String actionsButtonId, IDGroup menuItemID) {
      try {
         ActionNavigator.openMenuItem(actionsButtonId, false);
      } catch (Exception e) {
         SUITA.Factory.UI_AUTOMATION_TOOL.assertFor.isTrue(
               "Opening DataGrid actions menu throw exception: " + e.getMessage(),
               false
               );
      }

      ActionNavigator.invokeMenuAction(menuItemID);
   }

   /**
    * The method clicks on the menu item specified by the ID provided with the
    * menuItemID parameter. The method expects that the menu is opened. The
    * method will open the sub menu described by the subMenuIDs parameter.
    *
    * NOTE: This method expect that the first level menu is open and loaded.
    * Before invoking this method the developer have to validate it. This method
    * is recommended to be used with context menu for which the test need to
    * execute right-click mouse action.
    *
    * @param menuItemID    id of the menu item to be clicked.
    * @param subMenuIDs    path to the sub menu containing the menu item to click.
    *
    * @deprecated          use {@link #invokeMenuAction(IDGroup)}
    */
   @Deprecated
   public static void invokeMenuAction(IDGroup menuItemID, IDGroup... subMenuIDs) {
      ActionNavigator.invokeMenuAction(false, menuItemID);
   }

   /**
    * The method clicks on the menu item specified by the ID provided with the
    * menuItemID parameter. The method expects that the menu is opened. The
    * method will open the sub menu described by the subMenuIDs parameter.
    *
    * NOTE: This method expect that the first level menu is open and loaded.
    * Before invoking this method the developer have to validate it. This method
    * is recommended to be used with context menu for which the test need to
    * execute right-click mouse action.
    *
    * @param menuItemID    id of the menu item to be clicked.
    */
   public static void invokeMenuAction(IDGroup menuItemID) {
      ActionNavigator.invokeMenuAction(false, menuItemID);
   }

   /**
    * Check if menu action is present. The method expects that the menu is
    * opened. The method will open the sub menu described by the subMenuIDs
    * parameter.
    *
    * @param menuItemID    id of the menu item to be clicked
    * @param subMenuIds    path to the sub menu containing the menu item to click
    *
    * @return              if the menu action is present
    *
    * @deprecated          use {@link #isMenuActionPresent(IDGroup)}
    */
   @Deprecated
   public static boolean isMenuActionPresent(IDGroup menuItemID, IDGroup... subMenuIds) {
      return ActionNavigator.isMenuActionPresent(menuItemID);
   }

   /**
    * Check if menu action is present. The method expects that the menu is
    * opened. The method will open the sub menu described by the subMenuIDs
    * parameter.
    *
    * @param menuItemID    id of the menu item to be clicked
    *
    * @return              if the menu action is present
    */
   public static boolean isMenuActionPresent(IDGroup menuItemID) {
      try {
         String direcotMenuItemId = menuItemID.getValue(Property.DIRECT_ID);

         // Expand the menu to the specified action
         ActionNavigator.expandAction(
               direcotMenuItemId,
               false,
               flashSelenium,
               SUITA.Environment.getUIOperationTimeout()
               );

         if (!direcotMenuItemId.startsWith(IDConstants.ID_CONTEXT_MENU)) {
            direcotMenuItemId = ActionNavigator.getDirectMenuActionId(direcotMenuItemId);
         }

         // check if component is present before checking its properties
         if (!UI.condition.isFound(direcotMenuItemId)
               .await(SUITA.Environment.getUIOperationTimeout())) {
            return false;
         }

         if (!UI.component.property.getBoolean(Property.VISIBLE, direcotMenuItemId)) {
            return false;
         }

         return UI.component.property.getBoolean(
               Property.ENABLED,
               direcotMenuItemId + "/className=UITextField"
               );
      } catch (IllegalArgumentException iae) {
         // if in expandAction method, the action is not found, then we return action not present, i.e. false
         return false;
      } catch (Exception e) {
         SUITA.Factory.UI_AUTOMATION_TOOL.assertFor.isTrue(
               "Invocation of menu item throw exception: " + e.getMessage(),
               false
               );
      }

      return false;
   }

   /**
    * The method verifies if a menu item is enabled. The method expects that the
    * menu is opened. The method will open the sub menu described by the
    * subMenuIDs parameter.
    *
    * NOTE: This method expect that the first level menu is open and loaded.
    * Before invoking this method the developer have to validate it. This method
    * is recommended to be used with context menu for which the test need to
    * execute right-click mouse action.
    *
    * @param menuItemID    id of the menu item to be clicked
    * @param subMenuIDs    path to the sub menu containing the menu item to click
    *
    * @return              if the action is enabled or not
    *
    * @deprecated          {@link #isMenuActionEnabled(IDGroup)}
    */
   @Deprecated
   public static boolean isMenuActionEnabled(IDGroup menuItemID, IDGroup... subMenuIDs) {
      return ActionNavigator.isMenuActionEnabled(menuItemID);
   }

   /**
    * The method verifies if a menu item is enabled. The method expects that the
    * menu is opened. The method will open the sub menu described by the
    * subMenuIDs parameter.
    *
    * NOTE: This method expect that the first level menu is open and loaded.
    * Before invoking this method the developer have to validate it. This method
    * is recommended to be used with context menu for which the test need to
    * execute right-click mouse action.
    *
    * @param menuItemID    id of the menu item to be clicked
    *
    * @return              if the action is enabled or not
    */
   public static boolean isMenuActionEnabled(IDGroup menuItemID) {
      try {
         return ActionNavigator.isActionEnabledById(
               menuItemID.getValue(Property.DIRECT_ID)
               );
      } catch (Exception e) {
         SUITA.Factory.UI_AUTOMATION_TOOL.assertFor.isTrue(
               "Invocation of menu item throw exception: " + e.getMessage(),
               false);
         return false;
      }
   }


   /**
    * Method that clicks the logout menu item.
    */
   public static void invokeLogoutMenuItem() {
      Menu userMenu = new Menu(ID_USER_MENU, flashSelenium);
      userMenu.waitForElementEnable(SUITA.Environment.getUIOperationTimeout(), 1000);
      userMenu.selectMenuItemByName(
            LABEL_APP_LOGOUT,
            AN_TIMEOUT_ONE_SECOND
            );
   }

   /**
    * Invokes menu action.
    *
    * The menu has to be opened beforehand.
    *
    * @param actionId      the ID of the action that will be invoked
    * @param flashSelenium
    * @param timeout       the timeout
    * @throws Exception    if the action cannot be found in the menu
    */
   public static void invokeAction(String actionId, FlashSelenium flashSelenium,
         long timeout) throws Exception {
      expandAction(actionId, true, flashSelenium, timeout);
   }

   /**
    * Expands the menu to the action that was specified by an ID.
    *
    * The menu action can be invoked or only navigated to it.
    *
    * @param actionId      the ID of the action that will be invoked
    * @param invoke        if the action has to be invoked or not
    * @param flashSelenium
    * @param timeout       the timeout
    * @throws Exception    if the action cannot be found in the menu
    */
   public static void expandAction(String actionId, boolean invoke,
         FlashSelenium flashSelenium, long timeout) throws Exception {

      ActionNavigator.waitForContextMenu(flashSelenium, true);
      ActionNavigator.waitForRefreshButtonEnable(AN_TIMEOUT_ONE_MINUTE);

      MenuNode rootNode = new ContextMenuBuilder().getRootMenuNode();
      MenuNode menuNode = rootNode.getChildByActionId(actionId);

      if (menuNode != null) {
         menuNode.expandTo();

         if (invoke) {
            menuNode.leftMouseClick();

            // Close the context menu
            ActionNavigator.closeContextMenuSafely();
         }
      } else {
         _logger.error("Unable to navigate to action, actionId : " + actionId);
         throw new IllegalArgumentException(
               "Unable to navigate to action, actionId : " + actionId
               );
      }
   }

   /**
    * Generates menu action direct ID that comprises of all parent menus item IDs.
    *
    * @param actionId   the id of the menu action
    * @return           generated direct menu action ID
    */
   public static String getDirectMenuActionId(String actionId) {
      MenuNode rootNode = new ContextMenuBuilder().getRootMenuNode();
      MenuNode menuNode = rootNode.getChildByActionId(actionId);

      return menuNode.getDirectId();
   }

   /**
    * Method waits for context menu to be visible, enabled and fullyLoaded
    *
    * NOTE: This method has been taken from VCUI-QE-LIB: ActionFunction.
    *
    * @param flashSelenium
    * @throws Exception
    */
   public static void waitForContextMenu(FlashSelenium flashSelenium, boolean ready)
         throws Exception {
      UIComponent contextMenu =
            new UIComponent(IDConstants.ID_CONTEXT_MENU, flashSelenium);
      contextMenu.waitForElementEnable(AN_TIMEOUT_ONE_SECOND);
      if (ready) {
         // Wait for ContextMenu.fullyLoaded = true
         _logger.info("Waiting for ContextMenu.fullyLoaded = true");

         if (contextMenu.getPropertyExists("fullyLoaded")) {
            contextMenu.waitForComponentPropertyValue(
                  "fullyLoaded",
                  "true",
                  AN_TIMEOUT_ONE_SECOND,
                  60
                  );
         }
      }
   }

   // ---------------------------------------------------------------------------
   // Private methods

   /**
    * Opens the More Actions menu.
    *
    * NOTE: This method has been taken from VCUI-QE-LIB: ActionFunction.
    *
    * @throws Exception   if the actions menu is not shown on screen.
    */
   public static void openMoreActions() throws Exception {
      UIComponent moreActions = new UIComponent(IDConstants.ID_MORE_ACTIONS_ICON, flashSelenium);
      moreActions.waitForElementEnable(AN_TIMEOUT_ONE_MINUTE,
            AN_TIMEOUT_ONE_SECOND);
      moreActions.waitForElementVisibleOnPath(AN_TIMEOUT_ONE_MINUTE,
            AN_TIMEOUT_ONE_SECOND);
      moreActions.click();
   }

   /**
    * The method clicks on the menu item specified by the ID provided with the
    * menuItemID parameter. The method expects that the menu is opened. The
    * method will open the sub menu described by the subMenuIDs parameter.
    *
    * NOTE: This method expect that the first level menu is open and loaded.
    * Before invoking this method the developer have to validate it. This method
    * is recommended to be used with context menu for which the test need to
    * execute right-click mouse action.
    *
    * @param debugMode     debug info log on/off
    * @param menuItemID    id of the menu item to be clicked.
    */
   private static void invokeMenuAction(boolean debugMode, IDGroup menuItemID) {
      try {
         ActionNavigator.invokeActionById(
               debugMode,
               menuItemID.getValue(Property.DIRECT_ID)
               );

         // The Refresh button is replaced by the Loading spinner on action
         // invocation

      } catch (Exception e) {
         SUITA.Factory.UI_AUTOMATION_TOOL.assertFor.isTrue(
               "Invocation of menu item throw exception: " + e.getMessage(),
               false
               );
      }
   }

   /**
    * Opens the data grid "Actions" menu on the top of the main content pane and
    * clicks on the menu item specified by the ID provided with the menuItemID
    * parameter. The method will open the sub menus specified by the subMenuIDs
    * parameter and will capture screenshots of the menu actions
    *
    * @param debug         if screenshots should be made on every step
    * @param menuItemID    id of the menu item to be clicked.
    */
   private static void invokeDebugFromToolbarMenu(boolean debug, IDGroup menuItemID) {
      try {
         ActionNavigator.openMenuItem(IDConstants.ID_ADVANCE_DATAGRID_ALLACTIONS, debug);
      } catch (Exception e) {
         SUITA.Factory.UI_AUTOMATION_TOOL.assertFor.isTrue(
               "Opening DataGrid actions menu throw exception: " + e.getMessage(),
               false
               );
      }

      if (debug) {
         UI.audit.snapshotAppScreen(SubToolAudit.getFPID(), "DEBUG_OPEN_MENU");
      }

      ActionNavigator.invokeMenuAction(debug, menuItemID);

      if (debug) {
         UI.audit.snapshotAppScreen(SubToolAudit.getFPID(), "DEBUG_OPEN_MENU_ACTIONS");
      }
   }

   /**
    * Verifies if an action for objects is enabled. Searches by action id. Used
    * for objects which aren't ManagedObjectReference, e.g. LinceseKey,
    * HostProfile, etc. NOTE: The method expects that the Action menu is open.
    *
    * @param menuItemId    the id of the action that will be checked
    *
    * @return              if the action is enabled or not
    * @throws Exception    if the action can not be found in the menu
    */
   private static boolean isActionEnabledById(String menuItemId) throws Exception {

      // Expand the menu to the specified action
      ActionNavigator.expandAction(
            menuItemId,
            false,
            flashSelenium,
            SUITA.Environment.getUIOperationTimeout()
            );

      try {
         // We check the label object actually
         String actionLabelId = menuItemId;
         if (!actionLabelId.startsWith(IDConstants.ID_CONTEXT_MENU)) {
            actionLabelId = ActionNavigator.getDirectMenuActionId(actionLabelId);
         }

         actionLabelId += "/className=UITextField";

         UIComponent actionLabel = new UIComponent(actionLabelId, flashSelenium);
         ActionNavigator.waitForFirstVisibleElement(
               actionLabel,
               SUITA.Environment.getUIOperationTimeout(),
               false
               );

         ActionNavigator.waitForRefreshButtonEnable(
               SUITA.Environment.getPageLoadTimeout()
               );

         return actionLabel.getVisibleProperty("enabled").equals("true");
      } finally {
         // Close the context menu
         ActionNavigator.closeContextMenuSafely();
      }
   }

   /**
    * This method will invoke the specified by menuItemId action. Navigates to
    * the corresponding sub menu described by the subMenuPath parameter and
    * clicks on the menu item. After clicking on the menu item the method close
    * the opened menu. NOTE: The method expects that the Action menu is open and
    * loaded on the screen.
    *
    * @param debugMode     if screenshots has to be created
    * @param menuItemId    id of the action item
    *
    * @throws Exception    if the sub menu or menu item can not be found
    */
   private static void invokeActionById(boolean debugMode,
         String menuItemId) throws Exception {

      if (debugMode) {
         UI.audit.snapshotAppScreen(SubToolAudit.getFPID(), "DEBUG_SUBMENU");
      }

      // Do the click on the menu item.
      ActionNavigator.invokeAction(
            menuItemId,
            flashSelenium,
            SUITA.Environment.getUIOperationTimeout()
            );

      if (debugMode) {
         UI.audit.snapshotAppScreen(SubToolAudit.getFPID(), "SUBMENU_CLICKED");
      }
   }

   private static void openMenuItem(String actionsButtonID,
         boolean debugMode) throws Exception {
      UIComponent moreActions = new UIComponent(actionsButtonID, flashSelenium);
      moreActions.waitForElementEnable(SUITA.Environment.getUIOperationTimeout(), 1000);
      moreActions.click();

      if (debugMode) {
         UI.audit.snapshotAppScreen(SubToolAudit.getFPID(), "OPEN_MENU_ITEM");
      }

      ActionNavigator.waitForContextMenu(flashSelenium, true);
      ActionNavigator.waitForRefreshButtonEnable(
            SUITA.Environment.getUIOperationTimeout()
            );

      if (debugMode) {
         UI.audit.snapshotAppScreen(SubToolAudit.getFPID(), "OPEN_MENU_ITEM_LOADED");
      }
   }

   /**
    * Wait for the refresh button to be visible and enabled
    *
    * NOTE: This method has been taken from VCUI-QE-LIB: GlobalFunction.
    *
    * @param timeout
    */
   private static void waitForRefreshButtonEnable(long timeout) {
      Button refreshButton = new Button(IDConstants.ID_REFRESH_BUTTON, flashSelenium);
      refreshButton.waitForElementEnable(timeout);
   }

   /**
    * Closes context menu if it is opened.
    *
    * NOTE: This method has been taken from VCUI-QE-LIB: ActionFunction.
    *
    * @throws Exception
    */
   private static void closeContextMenuSafely() throws Exception {
      // Verification of context menu is over and we need to close it because
      // menu is still visible the safest way to do that is just to click on
      // actions bar, its safe to click on search input text on the top right
      // corner of the screen.
      UIComponent searchInputText = new UIComponent(
            IDConstants.ID_QUICK_SEARCH_BOX,
            flashSelenium
            );
      ActionNavigator.waitForFirstVisibleElement(searchInputText, AN_TIMEOUT_ONE_MINUTE, true);
      searchInputText.visibleMouseDownUp(AN_TIMEOUT_ONE_SECOND);
   }

   /**
    * Method waits for component to become visible
    *
    * NOTE: This method has been taken from VCUI-QE-LIB: ActionFunction.
    *
    * @param component     UI component that should appear
    * @param timeout       how long to wait for element's appearance
    * @throws Exception    if the element could not be found
    * @return              true if its visible
    */
   private static boolean waitForFirstVisibleElement(UIComponent component,
         long timeout, boolean failOnError) throws Exception {
      boolean result = true;
      String isVisible = component.getVisibleProperty("visible");

      if (!isVisible.equals("true")) {
         final long sleepTime = (timeout / 100) * 5; // Sleep time = 5% of
         // timeout
         final long startTime = System.currentTimeMillis();
         boolean suceeded = false;
         int retry = 1;

         // Wait until timeout is reached or find a visible element
         while (!suceeded
               && ((System.currentTimeMillis() - startTime) < timeout)) {
            _logger.info(
                  String.format(
                        "ActionNavigator.waitForFirstVisibleElement: Retry #%d",
                        retry
                        )
                  );

            isVisible = component.getVisibleProperty("visible");

            suceeded = isVisible.equals("true");
            retry++;

            // Sleep some time in order to avoid massive flex java calls
            Thread.sleep(sleepTime);
         }

         if (!suceeded) {
            String err = "ActionNavigator.waitForFirstVisibleElement: Reached maximum timeout "
                  + Long.toString(timeout) + " milliseconds.";
            _logger.warn(err);
            if (failOnError) {
               _logger.error("Fail On Error set to TRUE");
               TestBaseUI.verifySafely(suceeded, true, err);
            }
         }

         _logger.error("Result = Succeed = " + suceeded);
         result = suceeded;
      }

      return result;
   }
}
