/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.view;

import static com.vmware.flexui.selenium.BrowserUtil.flashSelenium;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.ActionNavigator;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.client.automation.vcuilib.commoncode.ContextHeader;
import com.vmware.client.automation.vcuilib.commoncode.IDConstants;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.custom.LabelItemRenderer;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Property;
import com.vmware.vsphere.client.automation.common.CommonUtil;
import com.vmware.vsphere.client.automation.common.VmStatus;

/**
 * The class represents the VM list view. The list can be found in vCenter - >
 * Virtual Machines.
 */
public class VmListView extends BaseView {
	private static final String VM_GRID_NAME_COLUMN = CommonUtil
			.getLocalizedString("vmList.grid.nameColumn");
	private static final String VM_GRID_STATUS_COLUMN = CommonUtil
			.getLocalizedString("vmList.grid.statusColumn");
	private static final String VM_GRID_STATUS_VALUE_NORMAL = CommonUtil
			.getLocalizedString("vmList.grid.statusValueNormal");
	private static final String DOWN_ARROW_BTN = "contextHeader/list/name=downArrowSkin";
	private static final IDGroup ID_GRID_VM = IDGroup
			.toIDGroup("VirtualMachine/list");
	private static final String VM_GRID_OVERALL_COMPLIANCE_COL = CommonUtil
			.getLocalizedString("vmGrid.vmOverallComplianceColumn");
	private static final String OVERALL_COMPLIANCE_COLUMN_AUTO_ID = "automationName=Overall Compliance";
   private static final IDGroup ID_SHOW_HIDE_SCROLLBAR = IDGroup
         .toIDGroup("contextHeader/list/className=VScrollBar");
   private static final IDGroup ID_SHOW_HIDE_SCROLLBAR_DOWN_BTN = IDGroup
         .toIDGroup("contextHeader/list/className=Button");
	private static final String FACET_BUTTON_ID = "facetIconBar/button[1]"; //TODO: istambolieva use a static id when issue 1285282 is fixed.
	private static final IDGroup OVERALL_COMPLIANCE_FACET_ID = IDGroup.toIDGroup("label=Overall Compliance/className=FacetExpandCollapsePanelSkin"); //TODO: istambolieva use a static id when issue 1293761 is fixed.
	private AdvancedDataGrid dataGridId;

	/**
	 * A mapping between a column index and its renderer in the placement policy
	 * grid.
	 */
	private static final Map<Integer, String> COL_SPEC = new HashMap<Integer, String>(
			1);

	static {
		// Mapping the correct renderer to columns in order to be able to get
		// column contents.
		COL_SPEC.put(2, LabelItemRenderer.class.getCanonicalName());
	}

	/**
	 * Checks whether the specified VM is listed in the VM list view.
	 *
	 * @param vmName
	 *            the name of the VM to be searched for
	 */
	public boolean isFoundInGrid(String vmName) {
		Integer rowIndex = getGrid().findItemByName(vmName);
		return (rowIndex != null) ? true : false;
	}

	/**
	 * Right-clicks on given VM in VM list
	 *
	 * @param vmName
	 *            the name of the VM to be right-clicked on
	 * @return true if the right-click is successful, false otherwise e.g. if
	 *         the VM is not found in the VM list.
	 */
	public boolean rightClickVm(String vmName) {
		return GridControl.rightClickEntity(getGrid(), VM_GRID_NAME_COLUMN,
				vmName);
	}

	/**
	 * Right-clicks on given VMs in VMs list
	 *
	 * @param vmNames
	 *            the names of VMs to be right-clicked on
	 * @return true if the right-click is successful, false otherwise, e.g. if
	 *         the VM is not found in the VMs list.
	 */
	public boolean rightClickVms(String... vmNames) {
		GridControl.rightClickEntity(getGrid(), VM_GRID_NAME_COLUMN, vmNames);

		try {
			ActionNavigator.waitForContextMenu(BrowserUtil.flashSelenium, true);
		} catch (Exception e) {
			_logger.error("Context menu was not opened correcly on VMs list");
			return false;
		}

		return true;
	}

	/**
	 * Select given VM in VM list
	 *
	 * @param vmName
	 *            the name of the VM to be selected
	 * @return true if the selection was successful, false otherwise
	 */
	public boolean selectVmByName(String vmName) {
		return GridControl.selectEntity(getGrid(), VM_GRID_NAME_COLUMN, vmName);
	}

	public String getActionsButtonId() {
		return IDConstants.ID_ADVANCE_DATAGRID_ALLACTIONS;
	}

	/**
	 * Gets the status of a VM visible in the vCenter>Virtual Machines grid.
	 *
	 * @param vmName
	 *            name of the VM which needs its Status gotten.
	 * @return A constant representing the status.
	 */
	public VmStatus getVmStatus(String vmName) {
		Integer rowIndex = getGrid().findItemByName(vmName);
		String statusValue = getGrid().getCellValue(rowIndex,
				VM_GRID_STATUS_COLUMN);
		if (statusValue.equals(VM_GRID_STATUS_VALUE_NORMAL)) {
			return VmStatus.NORMAL;
		}
		return VmStatus.UNKNOWN;
	}

	/**
	 * Selects and item and executes an action from it's context menu.
	 *
	 * @param entityName
	 *            - the name of the item to select
	 * @param addityionalEntityNames
	 *            - names of the additional items in case of multi-selection
	 * @param actionId
	 *            - the ID of the context menu action
	 */
	public void executeAction(String entityName,
			List<String> addityionalEntityNames, IDGroup actionId,
			IDGroup... submenuIds) {
		List<String> entityNames = new ArrayList<String>();
		entityNames.add(entityName);
		entityNames.addAll(addityionalEntityNames);

		// Get a reference to the advanced data grid.
		GridControl.rightClickEntity(getGrid(),
				CommonUtil.getLocalizedString("column.name"),
				entityNames.toArray(new String[] {}));
		ActionNavigator.invokeMenuAction(actionId, submenuIds);
	}

	/**
	 * Selects and item and executes an action from it's context menu.
	 *
	 * @param entityName
	 *            - the name of the item to select
	 * @param actionId
	 *            - the ID of the context menu action
	 * @param submenuIds
	 *            - IDs of sub menu items
	 */
	public void executeAction(String entityName, IDGroup actionId,
			IDGroup... submenuIds) {
		executeAction(entityName, Collections.<String> emptyList(), actionId,
				submenuIds);
	}

	/**
	 * Gets the index of a VM in the grid
	 *
	 * @param name
	 *            - the name of the VM
	 * @return the index of a VM in the grid or -1 if not found
	 */
	public int getVmIndex(String name) {
		AdvancedDataGrid grid = GridControl.findGrid(IDGroup
				.toIDGroup(ID_GRID_VM));
		return GridControl.getEntityIndex(grid,
				CommonUtil.getLocalizedString("vmGrid.vmNameColumnName"), name);
	}

   /**
    * Check if 'Overall compliance column exists in the datagrid.'
    */
   public boolean isOverallComplianceAvailable() {
      boolean isOverallCompliance = false;
      if (UI.condition.isFound(OVERALL_COMPLIANCE_COLUMN_AUTO_ID).await(
            SUITA.Environment.getUIOperationTimeout() / 5)) {
         isOverallCompliance = true;
      }
      return isOverallCompliance;
   }

   /**
    * Scroll the vertical scrollbar until "Overall compliance" is displayed.
    */
   public Boolean scrollToOverallComplianceItem() {
      AdvancedDataGrid grid = GridControl.findGrid(ID_GRID_VM);
      ContextHeader header = grid.openShowHideCoumnHeader(CommonUtil
            .getLocalizedString(VM_GRID_NAME_COLUMN));
      Boolean isFound = false;
      try {
         if (isOverallComplianceAvailable()) {
            isFound = true;
            header.selectDeselectColums(
                  new String[] { VM_GRID_OVERALL_COMPLIANCE_COL }, null);
            header.clickOK();
            return isFound;
         } else if (UI.condition.isFound(ID_SHOW_HIDE_SCROLLBAR).await(
               SUITA.Environment.getUIOperationTimeout() / 5)) {
            // the element is not found but scroll is available
            // scroll all the elements up and check if the desired element will
            // appear
            Integer scrollPosition = UI.component.property.getInteger(
                  Property.SCROLL_POSITION, ID_SHOW_HIDE_SCROLLBAR);
            Integer scrollMaxPosition = UI.component.property.getInteger(
                  Property.SCROLL_MAX_POS, ID_SHOW_HIDE_SCROLLBAR);
            UIComponent downButton = new UIComponent(DOWN_ARROW_BTN,
                  flashSelenium);
            while (scrollPosition < scrollMaxPosition) {
               scrollPosition++;
               downButton.leftMouseClick();

               // check if the element appeared
               if (isOverallComplianceAvailable()) {
                  isFound = true;
                  header.selectDeselectColums(
                        new String[] { VM_GRID_OVERALL_COMPLIANCE_COL }, null);
                  header.clickOK();
                  return isFound;
               }
            }
         }
      } catch (AssertionError ae) {
      }
      header.clickOK();
      return isFound;
   }


   /**
    *
    * Click Facet button
    */
   public void clickFacetButton() {
      if (UI.component.exists(FACET_BUTTON_ID)) {
         UI.component.click(FACET_BUTTON_ID);
      } else {
         _logger.info("The facet button is missing");
      }
   }

   /**
    * Check that Overall Compliance Filter exists in the VM Facet
    * @return <code>boolean</code> true if facet button is displayed.
    */
   public boolean IsOverallComplianceFacetPresent() {
      boolean isOverallFacetFilterFound = false;
      if (UI.condition.isFound(OVERALL_COMPLIANCE_FACET_ID).await(
            SUITA.Environment.getUIOperationTimeout() / 5)) {
         isOverallFacetFilterFound = true;
      }
      return isOverallFacetFilterFound;
   }

	// ---------------------------------------------------------------------------
	// Protected methods
	protected String getGridId() {
		return "vsphere.core.viVms.itemsView/list";
	}

	// ---------------------------------------------------------------------------
	// Private methods

	/**
	 * Finds and returns the advanced data grid on the VM list view.
	 *
	 * @return <code>AdvancedDataGrid</code> object, null if the current
	 *         location is not the VM view.
	 */
   private AdvancedDataGrid getGrid() {
      dataGridId = GridControl.findGrid(IDGroup.toIDGroup(getGridId()), COL_SPEC);
      return dataGridId;
   }

   /**
    * Clicks on the "Overall compliance" checkbox in the list of the hidden
    * columns.
    *
    */
   private boolean isOverallComplianceClicked(ContextHeader contextHeader, boolean isFoundOverallCompliance) {
      if (isOverallComplianceAvailable()) {
         isFoundOverallCompliance = true;
         contextHeader.selectDeselectColums(
               new String[] { VM_GRID_OVERALL_COMPLIANCE_COL }, null);
         contextHeader.clickOK();
      }
      return isFoundOverallCompliance;
   }
}