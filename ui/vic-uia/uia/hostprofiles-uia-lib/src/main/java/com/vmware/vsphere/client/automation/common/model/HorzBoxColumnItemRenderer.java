/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.model;

import com.vmware.flexui.componentframework.controls.custom.ItemRenderer;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;

/**
 * Item renderer that is used in the grid on Monitor -> Compliance tab of host profiles.
 */
public class HorzBoxColumnItemRenderer extends ItemRenderer {

   // IDs
   private static final String ID_CB = "/className=CheckBox";
   // Automation tool
   private static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   public HorzBoxColumnItemRenderer(String parentId, String itemRendererData) {
      super(parentId, itemRendererData);
      setComponentClassName("HorzBoxColumnRenderer");
   }

   @Override
   public String getDisplayedText() {
      return UI.component.property.get(
            Property.TEXT,
            IDGroup.toIDGroup(getItemRendererUniqueId() + getLabelId()));
   }

   @Override
   public void setComponentValue(String value) {
      UI.component.value
            .set(value, IDGroup.toIDGroup(getItemRendererUniqueId() + ID_CB));
   }

   protected String getLabelId() {
      return "/className=Label";
   }

}
