/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import java.util.HashMap;
import java.util.Map;

import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;

/**
 * This class handles the work with SettingsBlock component
 */
public class SettingsBlockControl {

   protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;
   private static final String UID = "uid";

   // Component IDs
   private String ID_TITLE_LABEL = "/titleLabel";
   private String ID_SUBTITLE_LABEL = "/className=Scroller/className=Label";
   private String ID_STACKBLOCK = "/className=StackBlock";

   private String id;

   /**
    * In order the constructor to initialize correctly, it has to be navigated to ui
    * component.
    *
    * @param id - some identifier of the settings block
    */
   public SettingsBlockControl(String id) {
      UIComponent sttngsBlock = new UIComponent(id, BrowserUtil.flashSelenium);
      String uid = sttngsBlock.getProperty(UID);
      this.id = UID + "=" + uid;

   }

   /**
    * Method that gets its title
    *
    * @return - the title
    */
   public String getTitle() {
      return UI.component.property.get(
            Property.TEXT,
            IDGroup.toIDGroup(this.id + ID_TITLE_LABEL));
   }

   /**
    * Method that gets its subtitle
    *
    * @return - the subtitle
    */
   public String getSubtitle() {
      return UI.component.property.get(
            Property.TEXT,
            IDGroup.toIDGroup(this.id + ID_SUBTITLE_LABEL));
   }

   /**
    * Method that gets all stackblock control's key/value pairs information
    *
    * @return - the key/value pairs of the stackblock controls
    */
   public Map<String, String> getData() {
      Map<String, String> stackBlockData = new HashMap<String, String>();
      String stackBlockId = this.id + ID_STACKBLOCK;
      int numStackBlockControls =
            UI.component.existingCount(IDGroup.toIDGroup(stackBlockId));

      for (int i = 0; i < numStackBlockControls; i++) {
         UIComponent stackBlCtrlUI =
               new UIComponent(stackBlockId + "[" + i + "]", BrowserUtil.flashSelenium);
         StackBlockControl stackBlCtrl =
               new StackBlockControl(UID + "=" + stackBlCtrlUI.getProperty(UID));
         stackBlockData.put(stackBlCtrl.getKey(), stackBlCtrl.getValue());
      }
      return stackBlockData;
   }

}
