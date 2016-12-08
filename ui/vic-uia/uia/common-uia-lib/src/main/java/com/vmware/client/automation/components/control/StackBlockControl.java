/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;

/**
 * This class handles the work with StackBlock component - it is a pair of key/value
 * elemnt
 */
public class StackBlockControl {

    private final static String CURRENT_STATE_PROPERTY = "currentState";
    private final static String EXPANDED_STATE_STATUS = "State_Expanded";
    private final static String ID_UI_TEXTFIELD = "/className=UITextField";
    private final static String ID_LABEL = "/className=PropertyGridLabel" + ID_UI_TEXTFIELD;
    private final static String ID_VALUE = "/className=LabelValueElement" + ID_UI_TEXTFIELD;
    private final static String ID_NUMBER_VALUE = "/className=NumberLabelValueElement" + ID_UI_TEXTFIELD;
    protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   private String id;

   public StackBlockControl(String id) {
      this.id = id;
   }

   /**
    * Method that returns the value of the key of the control
    *
    * @return - label value
    */
   public String getKey() {
      return UI.component.property.get(Property.TEXT, IDGroup.toIDGroup(id + ID_LABEL));
   }

    /**
     * Method that returns the value of value of the control
     *
     * @return - label value
     */
    public String getValue() {
        IDGroup labelId = IDGroup.toIDGroup(id + ID_VALUE);
        if (!UI.component.exists(labelId)) {
            labelId = IDGroup.toIDGroup(id + ID_NUMBER_VALUE);
            if (!UI.component.exists(labelId)) {
                throw new RuntimeException("Label not found!");
            }
        }
        return UI.component.property.get(Property.TEXT, labelId);
    }


   /**
    * This method is deprecated the SUITA interfaces should be used instead. Method that
    * gets whether a stack block is expanded or not
    *
    * @param id - id of the stack block
    * @return - true if expanded, false otherwise
    */
   @Deprecated
   public static boolean isExpanded(String id) {
      UIComponent component = new UIComponent(id, BrowserUtil.flashSelenium);
      try {
         String componentState = component.getProperty(CURRENT_STATE_PROPERTY);
         return componentState != null && componentState.equals(EXPANDED_STATE_STATUS);
      } catch (NullPointerException npe) {
         SUITA.Factory.UI_AUTOMATION_TOOL.logger.warn(String.format(
               "The searched property %s is not found for the object %s!",
               CURRENT_STATE_PROPERTY,
               id));
         return false;
      }
   }
}
