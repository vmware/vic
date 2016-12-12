package com.vmware.client.automation.components.navigator.navstep;

import com.vmware.client.automation.components.control.GridControl;
import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.client.automation.vcuilib.commoncode.AdvancedDataGrid;
import com.vmware.suitaf.apl.IDGroup;

/**
 * Navigation step that selects and double clicks (i.e. opens) an entity which has
 * a child-parent relationship with another entity. For example content library-vm template
 * or VDC-VM.
 *
 * The <code>ChildEntityNavStep</code> should be registered with a
 * <code>Navigator</code>. It will be automatically created by the
 * navigator when needed.
 */
public class ChildEntityNavStep extends BaseNavStep {
   private static final String GRID_ID = "list";
   private static final String COLUMN_NAME = "Name";
   private static final String NAV_ENTITY_ACTION = "childEntityAction";

   private final String childEntityName;
   private String gridId = GRID_ID;

   public ChildEntityNavStep(String childEntityName) {
      super(NAV_ENTITY_ACTION);
      this.childEntityName = childEntityName;
   }

   /**
    * Alternative constructor with grid ID parameter used to uniquely identify
    * the UI grid containing the child entity object
    *
    * @param gridId ID of the grid widget containing the child entity as a row item
    * @param childEntityName name of the child entity
    */
   public ChildEntityNavStep(String gridId, String childEntityName) {
      this(childEntityName);
      this.gridId = gridId;
   }

   @Override
   public void doNavigate(LocationSpec locationSpec) throws Exception {
      AdvancedDataGrid dataGrid = GridControl.findGrid(IDGroup.toIDGroup(gridId));
      Integer rowIndex = dataGrid.findItemByName(childEntityName);

      if (rowIndex == null || rowIndex < 0) {
         throw new IllegalArgumentException(
               "Couldn't find template with name '" + childEntityName + "'");
      }

      dataGrid.doubleClickCell(rowIndex, COLUMN_NAME);
   }
}
