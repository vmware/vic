package com.vmware.client.automation.vcuilib.commoncode;

import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.EQUALS_SIGN;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.INV_TREE_NODE_ID_PROPERTY;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.mx.custom.InventoryTree;
import com.vmware.flexui.selenium.BrowserUtil;

/**
 * Inventory tree implementation.
 *
 * NOTE: this class is a partial copy of the one from VCUI-QE-LIB
 */
public class InvTree {

   private static final Logger logger = LoggerFactory.getLogger(InvTree.class);

   public static InvTree invTree = null;
   public static String invTreeId = null;
   private InventoryTree inventoryTree = null;

   public static enum NODE_RELATION {
      PARENT_OF, CHILD_OF;
   }

   private InvTree() {
      invTreeId = IDConstants.ID_NAV_TREE;
      inventoryTree = new InventoryTree(invTreeId, BrowserUtil.flashSelenium);
   }

   public static InvTree getInstance() {
      synchronized (InvTree.class) {
         if (invTree == null) {
            invTree = new InvTree();
         }
      }
      return invTree;
   }

   /**
    * This function will wait for the node to be Visible.
    * nodeId should be in the same format as String returned by convertMorToString()
    *
    * @param nodeId
    * @param maxRetry
    * @return true if node becomes visible within maxRetry(in seconds), else false
    */
   public boolean waitForNodeReady(String nodeId, int maxRetry) {

      if (InvTree.getInstance().getInventoryTree().getVisible()) {
         logger.info("NavTree is Visible");

         UIComponent node =
               new UIComponent(invTreeId + "/" + INV_TREE_NODE_ID_PROPERTY + EQUALS_SIGN
                     + nodeId, BrowserUtil.flashSelenium);
         int count = 0;
         while (count < maxRetry && !(node.getVisibleProperty("visible").equals("true"))) {
            count++;
            try {
               Thread.sleep(TestConstants.DEFAULT_TIMEOUT_ONE_SECOND_LONG_VALUE);
            } catch (InterruptedException e) {
               e.printStackTrace();
               logger.error(e.getMessage());
            }
         }
         if (count >= maxRetry) {
            logger.info("Fail - WaitForNodeReady Count is more than " + maxRetry + " :: "
                  + count);
            TestBaseUI.captureSnapShot("_NodeNotReady_");
            return false;
         } else {
            logger.info("Success - WaitForNodeReady Count is less than " + maxRetry + " :: "
                  + count);
            return true;
         }
      } else {
         logger.info("NavTree is Not Visible");
         return false;
      }
   }

   public InventoryTree getInventoryTree() {
      return inventoryTree;
   }
}
