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
public class ComplianceColumnItemRenderer extends ItemRenderer {

   // IDs
   // TODO: add handling of clusters, if necessary private final static String ID_CLUSTER
   // = "/clusterComplianceContainer";
   private final static String ID_HOST = "/hostComplianceContainer/className=Label";
   // Automation tool
   private static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   public ComplianceColumnItemRenderer(String parentId, String itemRendererData) {
      super(parentId, itemRendererData);
      setComponentClassName("ComplianceColumnRenderer");
   }

   @Override
   public String getDisplayedText() {
      return UI.component.property.get(
            Property.TEXT,
            IDGroup.toIDGroup(getItemRendererUniqueId() + ID_HOST));
   }

}
