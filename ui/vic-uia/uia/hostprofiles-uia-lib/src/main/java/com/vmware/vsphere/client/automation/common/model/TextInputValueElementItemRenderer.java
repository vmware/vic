/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.model;

import com.vmware.flexui.componentframework.controls.custom.ItemRenderer;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;

/**
 * Item renderer that is used in the grid on Customize pages of host profiles
 * wizards. It works with TextInputValueElementItemRenderer taht is similar
 * to TextInput element.
 */
public class TextInputValueElementItemRenderer extends ItemRenderer {

   // ID of text field
   private final static String ID_UITEXTFIELD = "/className=UITextField";
   // Automation tool
   private static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   public TextInputValueElementItemRenderer(String parentId, String itemRendererData) {
      super(parentId, itemRendererData);
      setComponentClassName("TextInputValueElement");
   }

   @Override
   public String getDisplayedText() {
      return UI.component.property.get(
            Property.TEXT,
            IDGroup.toIDGroup(getItemRendererUniqueId() + ID_UITEXTFIELD));
   }

   @Override
   public void setComponentValue(String value) {
      UI.component.value.set(
            value,
            IDGroup.toIDGroup(getItemRendererUniqueId() + ID_UITEXTFIELD));
      UI.component.click(IDGroup.toIDGroup(getItemRendererUniqueId() + ID_UITEXTFIELD));
   }

}
