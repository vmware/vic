/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.components.navigator.navstep;

import static com.vmware.flexui.selenium.BrowserUtil.flashSelenium;

import com.google.common.base.Joiner;
import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.client.automation.vcuilib.commoncode.ObjNav;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.Property;

/**
 * The <code>EntityNavStep</code> is able to select a first class entity from an
 * inventory list.
 *
 * The action can be used from any inventory list displaying first class
 * entities.
 *
 * The <code>EntityNavStep</code> should be registered with a
 * <code>Navigator</code>. It will be automatically created by the navigator
 * when needed.
 */
public class EntityNavStep extends BaseNavStep {

   private static final String NAV_ENTITY_ACTION = "entityAction";
   private static final String INV_TREE_NODE_ID_NAME_PROPERTY = "data.name";
   private static final String SELECTED_SET_DATAGRID = "selectedSetDataGrid";
   private static final String SCROLL_BAR_ID = "className=VScrollBar";
   private static final String DROP_DOWN_ARROW_ID = "name=downArrowSkin";

   // Number of refresh retries in order to load the needed item in the
   // navigation tree. If the item is not found after that number of refreshes
   // the navigation will fail with item not found error.
   private static final int REFRESH_RETRY_COUNT = 5;

   private final String _entityName;

   /**
    * The constructor initializes the <code>EntityNavStep</code> with the name
    * of the entity to be selected.
    *
    * @param entityName
    *           The name of the entity to select.
    */
   public EntityNavStep(String entityName) {
      super(NAV_ENTITY_ACTION);
      _entityName = entityName;
   }

   @Override
   public void doNavigate(LocationSpec locationSpec) throws Exception {
      Joiner joiner = Joiner.on("/");
      final String scrollId = joiner.join(SELECTED_SET_DATAGRID, SCROLL_BAR_ID);

      boolean entityIsFound = false;
      String entityIdToNavigate = joiner.join(SELECTED_SET_DATAGRID,
            INV_TREE_NODE_ID_NAME_PROPERTY + "=" + _entityName);

      int refreshCount = 0;
      do {

         if (isEntityOnScreen(entityIdToNavigate)) {
            // the element was found no need of scrolling

            entityIsFound = true;
            try {
               ObjNav.selectNodeByName(null, _entityName);
            } catch (AssertionError e) {
               // TODO rreymer: remove this once the VCQE lib gets updated
               _logger.warn("ObjNav.initializeOpsObjects fails with Assertion error");
            }
         } else {
            _logger.info("Try to go through the scroller!");
            // Check if there is a vertical scroll bar
            if (UI.condition.isFound(scrollId).estimate()) {
               // the element is not found but scroll is available
               // scroll all the elements down and check if the desired element
               // will appear
               Integer scrollPosition = UI.component.property.getInteger(
                     Property.SCROLL_POSITION, scrollId);
               Integer scrollMaxPosition = UI.component.property.getInteger(
                     Property.SCROLL_MAX_POS, scrollId);
               while (scrollPosition < scrollMaxPosition) {
                  // Three clicks on the scroller down arrow to load the next items in the navigation pane.
                  clickScrollerDownButton(3);
                  scrollPosition = UI.component.property.getInteger(
                        Property.SCROLL_POSITION, scrollId);

                  // check if the element appeared
                  if (isEntityOnScreen(entityIdToNavigate)) {
                     entityIsFound = true;
                     // PR 1325918: Double click on the scroller down arrow to minimized the chance that the
                     // found entity is hidden by newly added items and autorefresh.
                     _logger.info("Click two times to assure the item is visible!");
                     clickScrollerDownButton(2);
                     // Select the found entity
                     ObjNav.selectNodeByName(null, _entityName);
                     break;
                  }
               }
            }

            // Entity is not found even after scrolling down - maybe it is not
            // loaded yet. Refresh the vsphere client.
            if (!entityIsFound) {
               refreshCount++;
               if (refreshCount >= REFRESH_RETRY_COUNT) {
                  throw new Exception("Unable to find the " + _entityName
                        + " in the vavigation tree!");
               }

               UI.audit.snapshotAppScreen("NO-ENTITY-" + _entityName, "");
               Thread.sleep(1000);
               new BaseView().refreshPage();
               UI.audit.snapshotAppScreen("NO-ENTITY-AFTER-REFRESH-"
                     + refreshCount + "COUNT" + _entityName, "");
            }
         }
      } while (!entityIsFound);
   }

   // ---------------------------------------------------------------------------
   // Private methods

   /**
    * Check and wait for the entity id to be shown in the navigation tree.
    *
    * @param entityIdToNavigate
    *           id of the entity to be waited for.
    * @return true if found on scree. If the id is not visible after timeout
    *         return false.
    */
   private boolean isEntityOnScreen(String entityIdToNavigate) {
      return UI.condition.isFound(entityIdToNavigate).await(
            SUITA.Environment.getUIOperationTimeout() / 5);
   }

   /**
    * Click on the navigation pane scroll bar down button.
    * @param positionsToMoveDown number of click operations to be executed.
    */
   private void clickScrollerDownButton(int positionsToMoveDown) {
      UIComponent downButton = new UIComponent(
            SELECTED_SET_DATAGRID + "/" + SCROLL_BAR_ID + "/"
                  + DROP_DOWN_ARROW_ID, flashSelenium);

      for(int i = 0; i < positionsToMoveDown; i++) {
         downButton.leftMouseClick();
      }
   }
}
