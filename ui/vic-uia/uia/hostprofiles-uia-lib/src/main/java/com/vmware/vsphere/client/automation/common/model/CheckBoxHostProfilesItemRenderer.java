package com.vmware.vsphere.client.automation.common.model;

import com.vmware.flexui.componentframework.controls.custom.ItemRenderer;
import com.vmware.flexui.componentframework.controls.mx.CheckBox;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;

public class CheckBoxHostProfilesItemRenderer extends ItemRenderer {

   private static final String CLASSNAME_CB = "/className=CheckBox";

   public CheckBoxHostProfilesItemRenderer(String parentId, String itemRendererData) {
      super(parentId, itemRendererData);
   }

   @Override
   public String getComponentValue() {
      CheckBox checkBox =
            new CheckBox(getItemRendererUniqueId() + CLASSNAME_CB, BrowserUtil.flashSelenium);
      return String.valueOf(checkBox.isChecked());
   }

   @Override
   public void setComponentValue(String value) {
      CheckBox checkBox =
            new CheckBox(getItemRendererUniqueId() + CLASSNAME_CB, BrowserUtil.flashSelenium);
      checkBox.checkUncheckListCheckBox(value, SUITA.Environment.getUIOperationTimeout());
   }

}
