package com.vmware.client.automation.vcuilib.commoncode;

import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.DEFAULT_TIMEOUT_ONE_MINUTE;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.thoughtworks.selenium.FlashSelenium;
import com.vmware.client.automation.vcuilib.commoncode.TestConstants.NODE_TYPE;
import com.vmware.flexui.componentframework.UIComponent;

/**
 * This Class is to invoke a specific Action(s) on a specified Object(s)
 * Functions related to the ActionFramework should go in this Static Class
 * Currently, there is 3 ways to invoke an Action : ContextMenu, ActionBar,
 * MoreActions
 *
 * NOTE: this class is partial copy of the one from VCUI-QE-LIB
 */
public class ActionFunction {

   private static final Logger logger = LoggerFactory.getLogger(ActionFunction.class);

   /**
    * Method waits for component to become visible
    *
    * @param component
    * @param timeout
    * @throws Exception
    * @return true if its visible
    */
   public static boolean waitForFirstVisibleElement(UIComponent component, long timeout,
         boolean failOnError) throws Exception {
      boolean result = true;
      String isVisible = component.getVisibleProperty("visible");

      if (!isVisible.equals("true")) {
         final long sleepTime = (timeout / 100) * 5; // Sleep time = 5% of
         // timeout
         final long startTime = System.currentTimeMillis();
         boolean suceeded = false;
         int retry = 1;

         // Wait until timeout is reached or find a visible element
         while (!suceeded && ((System.currentTimeMillis() - startTime) < timeout)) {
            logger.debug("ActionFunction.waitForFirstVisibleElement: Retry #"
                  + Integer.toString(retry));

            isVisible = component.getVisibleProperty("visible");

            suceeded = isVisible.equals("true");
            retry++;

            // Sleep some time in order to avoid massive flex java calls
            Thread.sleep(sleepTime);
         }

         if (!suceeded) {
            String err =
                  "ActionFunction.waitForFirstVisibleElement: Reached maximum timeout "
                        + Long.toString(timeout) + " milliseconds.";
            logger.error(err);
            if (failOnError) {
               logger.error("Fail On Error set to TRUE");
               TestBaseUI.verifySafely(suceeded, true, err);
            }
         }

         logger.info("Result = Succeed = " + suceeded);
         result = suceeded;
      }

      return result;
   }


   public static void invokeActionFromContextMenuForDataGrid(String id,
         FlashSelenium flashSelenium, long timeout) throws Exception {
      UIComponent action = new UIComponent(id, flashSelenium);
      waitForFirstVisibleElement(action, DEFAULT_TIMEOUT_ONE_MINUTE, true);
      action.visibleMouseDownUp(timeout);
   }

   /**
    * This function is to get the replaceValue for <ENTITY_TYPE> using morType
    *
    * @param morType
    * @return String retVal : This contains the replaceValue for Entity Type
    * @throws Exception
    */
   public static String convertMorTypeToExtensionEntity(String morType) throws Exception {
      String retVal = null;

      if (morType == NODE_TYPE.NT_VC.getNodeType()) {
         retVal = IDConstants.EXTENSION_ENTITY_VC;
      } else if (morType.equals(NODE_TYPE.NT_DATACENTER.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_DATACENTER;
      } else if (morType.equals(NODE_TYPE.NT_HOST.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_HOST;
      } else if (morType.equals(NODE_TYPE.NT_CLUSTER.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_CLUSTER;
      } else if (morType.equals(NODE_TYPE.NT_VM.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_VM;
      } else if (morType.equals(NODE_TYPE.NT_RESOURCE_POOL.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_RESOURCE_POOL;
      } else if (morType.equals(NODE_TYPE.NT_DATASTORE.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_DATASTORE;
      } else if (morType.equals(NODE_TYPE.NT_FOLDER.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_FOLDER;
      } else if (morType.equals(NODE_TYPE.NT_VMFOLDER.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_FOLDER;
      } else if (morType.equals(NODE_TYPE.NT_VAPP.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_VAPP;
      } else if (morType.equals(NODE_TYPE.NT_STANDARD_PORTGROUP.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_STANDARD_NETWORK;
      } else if (morType.equals(NODE_TYPE.NT_DV_SWITCH.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_DV_SWITCH;
      } else if (morType.equals(NODE_TYPE.NT_DV_PORTGROUP.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_DV_PORTGROUP;
      } else if (morType.equals(NODE_TYPE.NT_DV_UPLINK.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_DV_UPLINK;
      } else if (morType == NODE_TYPE.NT_TEMPLATE.getNodeType()) {
         retVal = IDConstants.EXTENSION_ENTITY_TEMPLATE;
      } else if (morType.equals(NODE_TYPE.NT_STORAGE_POD.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_STORAGE_POD;
      } else if (morType.equals(NODE_TYPE.NT_HOST_PROFILE.getNodeType())) {
         retVal = IDConstants.EXTENSION_ENTITY_HP;
      } else {
         throw new Exception("Unsupported Node Type: " + morType);
      }
      return retVal;
   }
}
