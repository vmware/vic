/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.model;


/**
 * Item renderer that is used in the grid on Copy Settings to Host Profiles -> Select
 * target profiles page
 */
public class HorzBoxColumnExItemRenderer extends HorzBoxColumnItemRenderer {

   public HorzBoxColumnExItemRenderer(String parentId, String itemRendererData) {
      super(parentId, itemRendererData);
   }

   @Override
   protected String getLabelId() {
      return "/className=LabelEx";
   }
}
