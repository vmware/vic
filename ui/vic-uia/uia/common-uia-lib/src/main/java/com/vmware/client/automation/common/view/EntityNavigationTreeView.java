/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.common.view;

import com.vmware.client.automation.vcuilib.commoncode.IDConstants;
import com.vmware.client.automation.vcuilib.commoncode.ObjNav;

/**
 * Implements the navigation tree view on the left
 */
public class EntityNavigationTreeView extends BaseView {

   /**
    * Retrieves focused entity name.
    *
    * @return
    */
   public String getFocusedEntityName() {
      return UI.component.value.get(IDConstants.ID_LABEL_FOCUSED_NODE);
   }

   /**
    * Verifies if a node is visible in Navigation pane
    *
    * @param nodeName
    * @throws Exception
    */
   public void verifyNodeNameVisible(String nodeName) throws Exception {
      ObjNav.verifyNodeNameVisible(nodeName);
   }
}
