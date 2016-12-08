/**
 * ************************************************************************
 *
 * Copyright 2009 VMware, Inc.  All rights reserved. -- VMware Confidential
 *
 * ************************************************************************
 */
package com.vmware.client.automation.vcuilib.commoncode;

import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.DETERMINE_NICS_TO_USE_LABEL;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.IPV4_ADDRESS_LABEL;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.MOVE_DOWN_UPLINK_BUTTON;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.MOVE_UP_UPLINK_BUTTON;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.STATIC_IPV6_ADDRESS_LABEL;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.UPGRADE_VDS_VERSION_41;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.UPGRADE_VDS_VERSION_50;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstantsKey.UPGRADE_VDS_VERSION_51;

/**
 * ID constants class. Naming convention is the constant name should start with
 * ID if the element is identified by the id, or LABEL if it is identified by
 * it's label etc.,
 *
 * NOTE: this class is a partial copy of the one from VCUI-QE-LIB
 */
public class IDConstants {

   // Certificate Warning Dialog
   public static final String ID_CERTIFICATE_WARNING_IGNORE_BUTTON = "buttonContinue";
   public static final String ID_CERTIFICATE_WARNING_CANCEL_BUTTON = "buttonCancel";
   public static final String ID_CERTIFICATE_WARNING_INSTALL_CERT_CHECKBOX =
         "checkBoxInstallCertificate";
   public static final String ID_PORTLET_SUFFIX = ".chrome";
   public static final String ID_CERTIFICATE_WARNING_ERROR_TEXT = "errWarningText";

   // Entities
   public static final String EXTENSION_ENTITY_VM = "vm.";
   public static final String EXTENSION_ENTITY_HOST = "host.";
   public static final String EXTENSION_ENTITY_CLUSTER = "cluster.";
   public static final String EXTENSION_ENTITY_RESOURCE_POOL = "resourcePool.";
   public static final String EXTENSION_ENTITY_VAPP = "vApp.";
   public static final String EXTENSION_ENTITY_DATACENTER = "datacenter.";
   public static final String EXTENSION_ENTITY_STANDARD_NETWORK = "network.";
   public static final String EXTENSION_ENTITY_DV_SWITCH = "dvs.";
   public static final String EXTENSION_ENTITY_DV_PORTGROUP = "dvPortgroup.";
   public static final String EXTENSION_ENTITY_DV_UPLINK = "dvPortgroup.";
   public static final String EXTENSION_ENTITY_DATASTORE = "datastore.";
   public static final String EXTENSION_ENTITY_FOLDER = "folder.";
   public static final String EXTENSION_ENTITY_VC = "folder.";
   public static final String EXTENSION_ENTITY_VCENTER = "vc.";
   public static final String EXTENSION_ENTITY_STORAGE_POD = "dscluster.";
   public static final String EXTENSION_ENTITY_TEMPLATE = "template.";
   public static final String EXTENSION_ENTITY_OVF = "ovf.";
   public static final String EXTENSION_ENTITY_HOST_PROFILE = "hostprofile.";
   public static final String EXTENSION_ENTITY_DATASTORE_EXPLORER =
         "datastore.explorer.";
   public static final String EXTENSION_ENTITY_HP = "hp.";
   public static final String EXTENSION_ENTITY_PROFILE = "profile.";
   public static final String EXTENSION_ENTITY_SPBM = "spbm.";

   // Extension points for each entity
   public static final String EXTENSION_PREFIX = "vsphere.core.";
   public static final String EXTENSION_ENTITY = "<ENTITY_TYPE>";
   public static final String MAIN_TAB_PREFIX = "<MAIN_TAB>";
   public static final String ID_OBJECT_VIEWS = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("views").toString();

   // VC selector
   public static final String ID_VC_SERVER_COMBOBOX = "serverComboBox";

   // Advanced datagrid tool bar IDs
   public static final String ID_SEARCHCONTROL_FILTERCONTROL = "filterControl";
   public static final String ID_TOOLBAR_FILTERCONTROL_TEXT_INPUT =
         "filterControl/textInput";
   public static final String ID_TOOLBAR_DATAGRID = "dataGridToolbar";
   public static final String ID_TOOBALR_DATAGRID_STATUS = "dataGridStatusbar";
   public static final String ID_CONTAINER_DATAGRID_TOOL = "toolContainer";
   public static final String ID_CONTAINER_DATAGRID_STATUS = "container";
   public static final String ID_IMAGE_FIND = "find";
   public static final String ID_SEARCHCONTROL_SEARCHCONTROL = "searchControl";
   public static final String ID_MENU_SEARCH_CONTEXT = "searchContextMenu";
   public static final String ID_ACTION_RESCAN_ADAPTER = "rescanAdapterButton";
   public static final String ID_ACTION_REFRESH_ADAPTER =
         "vsphere.core.host.refreshStorageSystem";
   public static final String ID_ACTION_REFRESH = "refresh";
   public static final String ID_ACTION_RESCAN_STORAGE = "vsphere.core.host.rescanHost";
   public static final String ID_IMAGE_COLLAPSEALL = "collapseAll";
   public static final String ID_IMAGE_EXPANDALL = "expandAll";
   public static final String ID_IMAGE_ADD = "add";
   public static final String ID_IMAGE_DELETE = "delete";
   public static final String ID_BUTTON_ADD_ISCSI = "addAdapterButton";
   public static final String ID_BUTTON_BUTTON_NOSLASH = "button";
   public static final String ID_GENERIC_CONTENTGROUP = "contentGroup";
   public static final String ID_DEFAULT_PROPERTY_VIEW = "defaultPropertyView";
   public static final String ID_TOOLBAR_DATAGRID_REFRESH_BUTTON = ID_TOOLBAR_DATAGRID
         + "/" + ID_CONTAINER_DATAGRID_TOOL + "/refresh";
   public static final String ID_STORAGE_TREE = "tocTree";

   public static final String ID_LABEL_MTU = "mtuLabel";
   // Solution IDs
   public static final String ID_SOLUTION_VSPHERE = "Virtual Infrastructure";
   public static final String ID_MONITOR_MAIN_TAB_CONTAINER = "monitor";
   public static final String ID_MANAGE_MAIN_TAB_CONTAINER = "manage";
   public static final String ID_RELATED_ITEMS_MAIN_TAB_CONTAINER = "related";

   // Sub tab views extension points
   public static final String ID_SUB_TAB_MONITOR_VIEWS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("monitorViews").toString();
   public static final String ID_SUB_TAB_MANAGE_VIEWS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("manageViews").toString();
   public static final String ID_SUB_TAB_RELATED_ITEMS_VIEWS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("relatedViews").toString();

   // Action Related IDs : Bug 521186 regarding [1]
   public static final String ID_CONTEXT_MENU = "afContextMenu";
   public static final String ID_NO_SUB_MENU = "noSubMenu";
   public static final String CLASS_NAME = "className=";
   public static final String AUTOMATIONNAME = "automationName=";
   public static final String AUTOMATIONVALUE = "automationValue";
   public static final String SHOW_ON_BAR_SUFFIX = ".showOnBar";
   // Action -- Categories
   public static final String ID_CATEGORY_ALL_VCENTER = "vCenter";
   public static final String ID_CATEGORY_POWER_OPS =
         "vsphere.core.powerOpsActionCategory";
   public static final String ID_CATEGORY_RESOURCE_MGMT =
         "vsphere.core.resourceManagementActionCategory";
   public static final String ID_CATEGORY_UNCATEGORIZED =
         "vsphere.core.uncategorizedActionCategory";
   public static final String ID_CATEGORY_INVENTORY =
         "vsphere.core.inventoryActionCategory";
   public static final String ID_CATEGORY_CONFIG = "vsphere.core.configActionCategory";
   public static final String ID_CATEGORY_SNAPSHOT =
         "vsphere.core.snapshotsActionCategory";
   public static final String ID_CATEGORY_FAULT_TOLERANCE =
         "vsphere.core.faultToleranceActionCategory";
   public static final String ID_CATEGORY_CONSOLE = "vsphere.console.consoleCategory";
   public static final String ID_CATEGORY_GUEST_OS = ID_CONTEXT_MENU + ".guestOs";
   // TODO:Need to change the name later, few packages using it
   public static final String ID_ACTIONS_BUTTONS = "actionBar";
   public static final String ID_ADVANCED_GRID_MENU = "AdvancedMenu";
   public static final String ID_ACTIONS_BAR = ID_ACTIONS_BUTTONS;
   public static final String ID_TITLE_BAR_GROUP = "titleBarGroup";
   public static final String ID_BUTTON = "/button";
   public static final String ID_MORE_ACTIONS_ICON = "actionButton";
   public static final String ID_BUTTON_LICENSE_OK =
         "_LicenseExpirationNotifyView_Button1";

   public static final String ID_LABEL_NUM_OF_PROCESSORS =
         "summary_numProcessors_valueLbl";
   public static final String ID_LABEL_NUM_OF_VMOTION =
         "summary_numVMotionMigrations_valueLbl";
   public static final String ID_LABEL_RESOURCE_USAGE_GRID_ROW_ID_0 =
         "resourceUsageGridRowId_0/";
   public static final String ID_LABEL_RESOURCE_USAGE_GRID_ROW_ID_1 =
         "resourceUsageGridRowId_1/";
   public static final String ID_LABEL_RESOURCE_USAGE_GRID_ROW_ID_2 =
         "resourceUsageGridRowId_2/";
   public static final String ID_LABEL_RES_CONSUME_GRID_VIEW_LIMIT_LBL =
         "/resConsumeGridView/limitLbl";
   public static final String ID_LABEL_RES_CONSUME_GRID_VIEW_RESERVATION_LBL =
         "/resConsumeGridView/reservationLbl";

   // Action -- Actions
   public static final String ID_ACTION_POWER_ON = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("powerOnAction").toString();
   public static final String ID_ACTION_POWER_OFF = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("powerOffAction").toString();
   public static final String ID_ACTION_SUSPEND = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("suspendAction").toString();
   public static final String ID_ACTION_RESET = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("resetAction").toString();
   public static final String ID_ACTION_PROVISION_DATASTORE = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DATASTORE).append("addAction")
         .toString();
   public static final String ID_ACTION_UNMOUNT_DATASTORE = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DATASTORE).append("unmountDatastore")
         .toString();
   public static final String ID_ACTION_MOUNT_DATASTORE = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DATASTORE).append("mountDatastore")
         .toString();
   public static final String ID_ACTION_REMOVE_DATASTORE = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DATASTORE)
         .append("deleteVmfsDatastoreAction").toString();
   public static final String ID_ACTION_RENAME_DATASTORE = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DATASTORE).append("renameAction")
         .toString();
   public static final String ID_ACTION_UNMOUNT_NFS_DATASTORE = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DATASTORE)
         .append("unmountNfsDatastoreAction").toString();
   public static final String ID_ACTION_CONFIGURE_STORAGE_IO_CONTROL = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DATASTORE)
         .append("configureStorageIOControl").toString();
   public static final String ID_ACTION_INCREASE_DATASTORE_CAPACITY = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DATASTORE).append("increaseAction")
         .toString();
   public static final String ID_ACTION_UPGRADE_VMFS5_DATASTORE = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DATASTORE)
         .append("upgradeVmfsDatastore").toString();
   public static final String ID_ACTION_DC_FOLDERS_RENAME = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("renameDcAction").toString();
   public static final String ID_ACTION_RESCAN_DATASTORE = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_HOST).append("rescanForDatastores")
         .toString();
   public static final String ID_ACTION_ADD_DIAGNOSTIC_PARTITION = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_HOST)
         .append("addDiagnosticPartitionAction").toString();
   public static final String ID_ACTION_MIGRATE = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("migrateAction").toString();
   public static final String ID_ACTION_DEPLOY_VIRTUAL_MACHINE = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("provisioning.")
         .append("cloneTemplateToVmAction").toString();
   public static final String ID_ACTION_UNREGISTER = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("unregisterAction").toString();
   public static final String ID_ACTION_UNREGISTER_DC = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("unregisterDcAction")
         .toString();
   public static final String ID_ACTION_REMOVE = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("removeAction").toString();
   public static final String ID_ACTION_DELETE = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("deleteAction").toString();
   public static final String ID_ACTION_RENAME = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("renameAction").toString();
   public static final String ID_ACTION_INSTALL_UPGRADE_TOOLS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("installUpgradeToolsAction")
         .toString();
   public static final String ID_ACTION_UNMOUNT_TOOLS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY)
         .append("unmountToolsInstallerAction").toString();

   public static final String ID_ACTION_EDIT_SETTINGS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("provisioning.")
         .append("editAction").toString();
   public static final String ID_ACTION_EDIT_SETTINGS_VM = "btn_" + EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "provisioning.editAction";
   public static final String ID_HOST_STORAGE_IO_CONTROL = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("configureStorageIOControl")
         .toString();
   public static final String ID_ACTION_NEW_FOLDER = EXTENSION_PREFIX
         + EXTENSION_ENTITY_FOLDER + "createFolderAction";
   public static final String ID_ACTION_NEW_VMFOLDER = EXTENSION_PREFIX
         + EXTENSION_ENTITY_DATACENTER + "createVmFolderAction";
   public static final String ID_ACTION_NEW_VMSUBFOLDER = EXTENSION_PREFIX
         + EXTENSION_ENTITY_FOLDER + "createVmFolderAction";
   public static final String ID_ACTION_NETWORK_FOLDER = EXTENSION_PREFIX
         + EXTENSION_ENTITY_DATACENTER + "createNetworkFolderAction";
   public static final String ID_ACTION_STORAGE_FOLDER = EXTENSION_PREFIX
         + EXTENSION_ENTITY_DATACENTER + "createDatastoreFolderAction";
   public static final String ID_ACTION_HOSTANDCLUSTER_FOLDER = EXTENSION_PREFIX
         + EXTENSION_ENTITY_DATACENTER + "createHostFolderAction";
   public static final String ID_ACTION_EDIT_SETTINGS_VAPP = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VAPP + "editSettingsAction";
   public static final String ID_ACTION_EDIT_SETTINGS_RP = new StringBuffer(
         EXTENSION_PREFIX).append("resourcePool.editAction").toString();
   public static final String ID_ACTION_EDIT_CPU_SETTINGS_RP = new StringBuffer(
         EXTENSION_PREFIX).append("resourcePool.editCpuAction").toString();
   public static final String ID_BUTTON_EDIT_CPU_SETTINGS_RP = new StringBuffer("btn_")
         .append(ID_ACTION_EDIT_CPU_SETTINGS_RP).toString();
   public static final String ID_ACTION_EDIT_MEMORY_SETTINGS_RP = new StringBuffer(
         EXTENSION_PREFIX).append("resourcePool.editMemoryAction").toString();
   public static final String ID_BUTTON_EDIT_MEMORY_SETTINGS_RP = new StringBuffer(
         "btn_").append(ID_ACTION_EDIT_MEMORY_SETTINGS_RP).toString();
   public static final String ID_EDIT_VM_START_SHUTDOWN_CONFIG = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("editVmStartupAction")
         .toString();
   public static final String ID_ACTION_NEW_RP = new StringBuffer(EXTENSION_PREFIX)
         .append("resourcePool.createAction").append(SHOW_ON_BAR_SUFFIX).toString();
   public static final String ID_ACTION_NEW_VAPP = new StringBuffer(EXTENSION_PREFIX)
         .append("vApp.createAction").append(SHOW_ON_BAR_SUFFIX).toString();
   public static final String ID_ACTION_REMOVE_RP = new StringBuffer(EXTENSION_PREFIX)
         .append("resourcePool.removeAction").toString();
   public static final String ID_ACTION_SHUTDOWN_GUEST = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("shutdownAction").toString();
   public static final String ID_ACTION_SHUTDOWN_SINGLE_GUEST = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("shutdownActionSingleVm")
         .toString();
   public static final String ID_ACTION_CLONE = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("provisioning.cloneVmToVmAction").toString();
   public static final String ID_ACTION_RESTART_GUEST = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("rebootAction").toString();
   public static final String ID_ACTION_REVERT_CURRENT_SNAPSHOT = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("revertToCurrentSnapshot")
         .toString();
   public static final String ID_ACTION_CONSOLIDATE = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_VM).append("consolidateSnapshots").toString();
   public static final String ID_ACTION_CONVERT_TO_TEMPLATE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "markAsTemplateAction";
   public static final String ID_ACTION_CREATE_VM = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "provisioning" + ".createVmAction";
   public static final String ID_ACTION_CREATE_VM_DATASTORE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "provisioning.createVmAction.pinStorage";
   public static final String ID_ACTION_CREATE_VM_VMFolder = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "provisioning.createVmAction.pinVmFolder";
   public static final String ID_ACTION_VM_CLONE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "provisioning" + ".cloneVmToVmAction";
   public static final String ID_ACTION_VM_CLONE_TO_TEMPLATE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "provisioning" + ".cloneVmToTemplateAction";
   public static final String ID_ACTION_CREATE_CLUSTER = EXTENSION_PREFIX
         + EXTENSION_ENTITY_CLUSTER + "createAction";
   public static final String ID_ACTION_EXPORT_OVF = EXTENSION_PREFIX
         + EXTENSION_ENTITY_OVF + "exportAction";
   public static final String ID_ACTION_DEPLOY_OVF = EXTENSION_PREFIX
         + EXTENSION_ENTITY_OVF + "deployAction";
   public static final String ID_ACTION_NEW_DC = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("createAction").toString();
   public static final String ID_ACTION_SHUTDOWN = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("shutdownAction").toString();
   public static final String ID_PLUGIN_DISABLE = new StringBuffer(EXTENSION_PREFIX)
         .append("admin.disablePluginAction").toString();
   public static final String ID_PLUGIN_ENABLE = new StringBuffer(EXTENSION_PREFIX)
         .append("admin.enablePluginAction").toString();
   public static final String ID_ACTION_TURN_FT_ON = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("ftTurnOnAction").toString();
   public static final String ID_ACTION_TURN_FT_OFF = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("ftTurnOffAction").toString();
   public static final String ID_ACTION_DISABLE_FT = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("ftDisableAction").toString();
   public static final String ID_ACTION_ENABLE_FT = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("ftEnableAction").toString();
   public static final String ID_ACTION_MIGRATE_SECONDARY = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("ftMigrateSecondaryAction")
         .toString();
   public static final String ID_ACTION_TEST_FAILOVER = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("ftTestFailoverAction")
         .toString();
   public static final String ID_ACTION_TEST_RESTART_SECONDARY = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY)
         .append("ftTestRestartSecondaryAction").toString();
   public static final String ID_ACTION_OPEN_CONSOLE =
         "vsphere.console.openConsoleAction";
   public static final String ID_ACTION_RECONFIGURE_HA =
         "vsphere.core.host.reconfigureDasAction";

   public static final String ID_ACTION_CREATE_VDS = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_DV_SWITCH).append("createDvsAction").toString();

   public static final String ID_ACTION_CREATE_PORT_GROUP = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH).append("createPortgroup")
         .toString();

   public static final String ID_ACTION_ADD_STORAGE = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("moveDatastoresIntoAction").toString();

   public static String ID_ACTION_RENAME_DS_CLUSTER = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_STORAGE_POD).append("renameAction").toString();

   public static String ID_ACTION_NEW_DS_CLUSTER = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_STORAGE_POD).append("createAction").toString();

   public static String ID_ACTION_REMOVE_DS_CLUSTER = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_STORAGE_POD).append("removeAction").toString();

   public static final String ID_ACTION_STARTED_TAB__CREATE_VM =
         "automationName=Getting Started/" + ID_ACTION_CREATE_VM;

   public static final String ID_ACTION_EDIT_SETTINGS_VDS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
         .append("editDvsSettingsAction").toString();

   public static final String ID_ACTION_EDIT_NETFLOW_VDS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
         .append("editNetFlowAction").toString();

   public static final String ID_ACTION_EDIT_HEALTH_CHECK = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
         .append("editHealthCheckAction").toString();

   public static final String ID_ACTION_EDIT_PRIVATE_VLAN = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
         .append("editPrivateVlanAction").toString();

   public static final String ID_ACTION_VDS_MANAGE_HOST = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
         .append("manageHostOnDvsAction").toString();

   public static final String ID_ACTION_RESTORE_CONFIG_VDS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
         .append("restoreConfigAction").toString();

   public static final String ID_ACTION_EXPORT_CONFIG_VDS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
         .append("exportConfigAction").toString();

   public static final String ID_ACTION_MANAGE_PORTGROUPS_VDS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH).append("managePortgroups")
         .toString();

   public static final String ID_ACTION_EDIT_SETTINGS_DVPORTGROUP = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_PORTGROUP)
         .append("editSettingsAction").toString();

   public static final String ID_ACTION_RENAME_DVPORTGROUP = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_PORTGROUP).append("renameAction")
         .toString();

   public static final String ID_ACTION_RESTORE_CONFIG_DVPORTGROUP = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_PORTGROUP)
         .append("restoreConfigAction").toString();

   public static final String ID_ACTION_EXPORT_CONFIG_DVPORTGROUP = new StringBuffer(
         EXTENSION_PREFIX).append("dvpg.").append("exportConfigAction").toString();

   public static final String ID_ACTION_REMOVE_DVPORTGROUP = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_PORTGROUP).append("removeAction")
         .toString();

   public static final String ID_ACTION_REMOVE_VDS = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_DV_SWITCH).append("removeAction").toString();

   public static final String ID_ACTION_CREATE_DVPORTGROUP = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH).append("createPortgroup")
         .toString();

   public static final String ID_ACTION_ADDHOST_TO_VDS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
         .append("addHostToDvsAction").toString();

   public static final String ID_ACTION_RENAME_VDS = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_DV_SWITCH).append("renameAction").toString();

   public static final String ID_ACTION_MOVE_VDS = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_DV_SWITCH).append("moveAction").toString();

   public static final String ID_ACTION_TEMPLATE_CLONE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "provisioning" + ".cloneTemplateToTemplateAction";
   public static final String ID_ACTION_CONVERT_TEMPLATE_TO_VM = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "provisioning" + ".convertTemplateToVmAction";
   public static final String ID_ACTION_POWER_ON_VAPP = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VAPP + "powerOnAction";
   public static final String ID_ACTION_CLONE_VAPP = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VAPP + "cloneAction";
   public static final String ID_ACTION_UPGRADE_VIRTUAL_HW = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "upgradeVirtualHardwareAction";
   public static final String ID_ACTION_TAKE_SNAPSHOT = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "takeSnapshotAction";
   public static final String ID_ACTION_SNAPSHOT_MANAGER = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "manageSnapshotAction";
   public static final String ID_ACTION_REBOOT = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + "rebootAction";
   public static final String ID_ACTION_SHUTDOWN_HOST = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + "shutdownAction";
   public static final String ID_ACTION_ENTER_MAINTENANCE_MODE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + "enterMaintenanceAction";
   public static final String ID_ACTION_EXIT_MAINTENANCE_MODE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + "exitMaintenanceAction";
   public static final String ID_ACTION_SDRS_ENTER_MAINTENANCE_MODE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_DATASTORE + "enterMaintenanceAction";
   public static final String ID_ACTION_SDRS_EXIT_MAINTENANCE_MODE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_DATASTORE + "exitMaintenanceAction";
   public static final String ID_ACTION_MOVE_OUT_OF_DATASTORE_CLUSTER = EXTENSION_PREFIX
         + EXTENSION_ENTITY_DATASTORE + "moveOutFromDatastoreClusterAction";
   public static final String ID_ACTION_DISCONNECT = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + "disconnectAction";
   public static final String ID_ACTION_RECONNECT = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + "reconnectAction";
   public static final String ID_ACTION_ADDHOST = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + "addAction";
   public static final String ID_ACTION_RESTORE_RESOURCEPOOL = EXTENSION_PREFIX
         + EXTENSION_ENTITY_CLUSTER + "restoreRpTreeAction";
   public static final String ID_ACTION_MIGRATE_VM_NETWORKING = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH).append("vmMigrateAction")
         .toString();
   public static final String ID_ACTION_UPGRADE_VDS = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_DV_SWITCH).append("upgradeDvsAction").toString();
   public static final String ID_ACTION_MOVE = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("moveAction").toString();
   public static final String ID_ACTION_ENTER_STANDBY_MODE_HOST = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + "enterStandbyAction";
   public static final String ID_ACTION_POWER_ON_HOST = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + "exitStandbyAction";
   // Added constants for add alarm & disable alarm in context menu
   public static final String ID_ACTION_ADD_ALARM = "vsphere.opsmgmt.alarms.addAlarm";
   public static final String ID_ACTION_EDIT_ALARM = "vsphere.opsmgmt.alarms.edit";
   public static final String ID_ACTION_DISABLE_ALARM =
         "vsphere.opsmgmt.alarms.disableAlarmActions";
   public static final String ID_VC_SERVER_ADDRESS = "serverComboBox";
   public static final String ID_VC_SERVER_USER = "usernameTextInput";
   public static final String ID_VC_SERVER_PASSWORD = "passwordTextInput";
   public static final String ID_LOGIN_BUTTON = "loginButton";
   public static final String ID_LOGIN_ERROR_LABEL = "_ErrorComponent_Text1";
   public static final String ID_LOGIN_CLIENT_INFO_LABEL = "clientVersion";
   public static final String ID_LOGIN_INTEGRATION_PLUGIN_INFO_LABEL = "pluginVersion";
   public static final String ID_LOGIN_QUESTION_MARK = "vmrcHelpButton";
   public static final String ID_SIGN_POST_POPUP = "vmrcHelpButton.signpostpopup";
   public static final String ID_SIGN_POST_POPUP_CLOSE_BUTTON = new StringBuffer(
         ID_SIGN_POST_POPUP).append("/")
         .append("automationName=vmrcHelpButton.signpostpopupClosebutton").toString();
   public static final String ID_SIGN_POST_POPUP_MORE_INFORMATION_LINK =
         new StringBuffer(ID_SIGN_POST_POPUP).append("/").append("contextLink")
               .toString();
   public static final String ID_SIGN_POST_POPUP_LABEL = new StringBuffer(
         ID_SIGN_POST_POPUP).append("/").append("displayText").toString();
   public static final String ID_OBJ_NAV_BACK_BUTTON = "appObjectNavigatorBackButton";
   public static final String ID_REMEMBER_ME_CHECKBOX = "rememberMeCheckBox";
   public static final String ID_REFRESH_BUTTON = "refreshButton";
   public static final String ID_REFRESH_BUTTON_DATA_GRID =
         "advancedSearch/searchResults/resultsBox/resultDetailsPanel/refreshButton";
   public static final String ID_APPLICATIONS_BUTTON = "appsButton";
   public static final String ID_APPLICATION_LABEL = "currentAppLabel";
   public static final String ID_APPLICATIONS_MENU = "launchMenu";
   public static final String ID_BUTTON_MINIMIZE_TIWO = "minimizeButton";
   public static final String ID_APPS_HELP_BUTTON = "helpButton";
   public static final String ID_APPS_HELP_MENU = "helpMenu";
   public static final String LABEL_APPS_HELPS_ABOUT = "About VMware vSphere";
   public static final String ID_ABOUT_DIALOG_BOX = "aboutDialog";
   public static final String ID_ABOUT_DIALOG_OK_BUTTON = new StringBuffer(
         ID_ABOUT_DIALOG_BOX).append("/okButton").toString();
   public static final String ID_ABOUT_DIALOG_CLIENT_ABOUT_LABEL = new StringBuffer(
         ID_ABOUT_DIALOG_BOX).append("/vsAbout").toString();
   public static final String ID_ABOUT_DIALOG_VC_ABOUT_LABEL = new StringBuffer(
         ID_ABOUT_DIALOG_BOX).append("/vcAbout").toString();
   public static final String ID_ABOUT_DIALOG_INTEGRATION_PLUGIN_ABOUT_LABEL =
         new StringBuffer(ID_ABOUT_DIALOG_BOX).append("/vsPluginAbout").toString();
   public static final String ID_NAV_TREE = "navTree";
   public static final String ID_NAV_TREE_NODE = ID_NAV_TREE + "/data.id=";
   public static final String ID_TREE = "tree";
   public static final String ID_NAV_TREE_SIDEBAR = "navigatorSidebar";
   public static final String ID_NAVIGATOR_VIEW_SELECTOR = "navigatorViewSelector/";
   public static final String ID_HOSTS_AND_CLUSTERS_VIEW_BUTTON = new StringBuffer(
         ID_NAVIGATOR_VIEW_SELECTOR).append("vsphere.core.physicalInventorySpec")
         .toString();
   public static final String ID_VMS_AND_TEMPLATES_VIEW_BUTTON = new StringBuffer(
         ID_NAVIGATOR_VIEW_SELECTOR).append("vsphere.core.virtualInventorySpec")
         .toString();
   public static final String ID_DATASTORES_VIEW_BUTTON = new StringBuffer(
         ID_NAVIGATOR_VIEW_SELECTOR).append("vsphere.core.storageInventorySpec")
         .toString();
   public static final String ID_NETWORKING_VIEW_BUTTON = new StringBuffer(
         ID_NAVIGATOR_VIEW_SELECTOR).append("vsphere.core.networkingInventorySpec")
         .toString();
   public static final String ID_HEALTH_PROFILES_VIEW_BUTTON = new StringBuffer(
         ID_NAVIGATOR_VIEW_SELECTOR).append("automationName=Health profiles").toString();
   // TODO file a PR to get the Id of the "messageSpinne" -> Inv Tree spinner
   public static final String ID_INV_TREE_SPINNER = new StringBuffer(
         "automationName=messageSpinner").toString();
   public static final String ID_FLEX_APP = "container_app";
   public static final String ID_FLEX_ADMIN_APP = "admin-app";
   public static final String ID_FLEX_CSHARP_APP = "csharp-app";
   public static final String ID_SHORTCUTS_BAR = "shortcutsBar";
   public static final String ID_CURRENT_APPLICATION = "currentAppLabel";
   public static final String LABEL_VCENTER_MANAGEMENT = "vCenter Management";
   public static final String LABEL_APP_TASK_CONSOLE = "Task Console";
   public static final String LABEL_APP_EVENT_CONSOLE = "Event Console";
   public static final String ID_LOGGED_IN_USER_POPUP_BUTTON = "loggedInUser";
   public static final String ID_USER_MENU = "userMenu";
   public static final String LABEL_APP_LOGOUT = "Logout";
   public static final String ID_MENU_ITEM_LOGOUT =
         "afContextMenu.SyntheticAction.000000.Logout";
   public static final String LABEL_APP_LOGOUT_USERNAME = new StringBuffer(
         LABEL_APP_LOGOUT).append(" <USER_NAME>...").toString();
   public static final String ID_APP_MONITORING_TASKCONSOLE = EXTENSION_PREFIX
         + "tasks.application";
   public static final String LABEL_APP_SEARCH = "Search";
   public static final String LABEL_APP_REPORTING = "Reporting";
   public static final String LABEL_APP_SYSTEM_ADMINISTRATION = "System Administration";
   public static final String LABEL_APP_PLUGIN_MANAGEMENT = "Plug-In Management";
   public static final String LABEL_APP_SAMPLE_APPLICATION = "Sample Application";
   public static final String LABEL_APP_MONITORING = "Monitoring";
   public static final String LABEL_APP_SMASH_TEST = "Smash Test";
   public static final String LABEL_HEALTH_MANAGEMENT = "Health Management";
   public static final String LABEL_CONTROL_CENTER = "Control Center";
   public static final String LABEL_LICENSING = "Licensing";
   public static final String ID_GETTING_STARTED_ACTION_PREFIX =
         "automationName=Getting Started/";
   public static final String LABEL_APP_CHANGE_PASSWORD = "Change Password...";

   // Change Password Dialog
   public static final String ID_TEXTINPUT_CHANGEPASSWORD_OLDPASSWORD = "oldPassword";
   public static final String ID_TEXTINPUT_CHANGEPASSWORD_NEWPASSWORD = "newPassword";
   public static final String ID_TEXTINPUT_CHANGEPASSWORD_CONFIRMPASSWORD =
         "newPasswordConfirm";

   // Message of the day banner
   public static final String ID_TOP_NOTICE_BANNER = "noticeBannerView";
   public static final String ID_TOP_NOTICE_BANNER_CLOSE_BUTTON = "bannerCloseButton";

   // TODO: ID_EXPAND_COLLAPSE_BUTTON has to be removed once everyone has
   // fixed their implementation using ID_SLIDE_PANEL
   public static final String ID_EXPAND_COLLAPSE_BUTTON = "expandCollapseImage";

   // Property of ID_SLIDE_PANEL
   public static final String ID_SIDEBAR_ARROW_BUTTON = "arrowButton";
   public static final String ID_SIDEBAR_PORTLET_EXPAND_COLLAPSE_BUTTON =
         "expandCollapseButton";
   public static final String ID_APP_SIDEBAR_BUTTON = "appSidebar";
   public static final String ID_RECENT_TASKS_ICON_ON_SIDE_BAR = ID_APP_SIDEBAR_BUTTON
         + "/vsphere.core.tasks.recentTasksView/iconPlaceholder/iconComponent";
   public static final String ID_SIDEBAR_PUSHPIN_ICON = ID_APP_SIDEBAR_BUTTON
         + "/pushPinButton";

   public static final String ID_ALERTS_BUTTON =
         "com.vmware.vsphere.client.samplePlugin.samplePortlet1";
   public static final String ID_INVENTORY_BUTTON = "vsphere.core.inventory.application";
   public static final String ID_RECENT_TASKS_BUTTON =
         "com.vmware.vsphere.client.samplePlugin.samplePortlet2";
   public static final String ID_VIEW_COMBO_BOX = "viewSelectCbx";
   public static final String ID_MIGRATE_BUTTON = "vsphere.core.vm.migrateAction";
   public static final String ID_NAVIGATOR = "navigator";
   public static final String ID_WORK_IN_PROGRESS_PORTLET =
         "vsphere.core.tiwo.sidebarView.portletChrome";
   public static final String ID_WORK_IN_PROGRESS_ICON = "vsphere.core.tiwo.sidebarView";
   public static final String ID_WORK_IN_PROGRESS_TIWO_LIST = "tiwoListView";
   public static final String ID_WORK_IN_PROGRESS_TIWO_TASK_LIST = "itemDescription";
   public static final String ID_WORK_IN_PROGRESS_TIWO_TASK_LIST_SUFFIX = "itemSuffix";

   public static final String ID_SEARCH_BUTTON = "vsphere.core.search.application";

   // Move To wizard
   public static final String ID_TABNAVIGATOR_MOVETO_OPTION =
         "objectSelectorDialog/_mainTabNavigator";
   public static final String ID_TABNAVIGATOR_MOVETO_ENTITY_TYPE =
         ID_TABNAVIGATOR_MOVETO_OPTION + "/"
               + "filterObjectsContainer/filterViewToggleButtonBar";
   public static final String ID_ADVANCEDATAGRID_MOVETO_RP =
         "ResourcePool_filterGridView";
   public static final String ID_DIALOG_MOVE_TO = "selectionWidgetDialog";
   public static final String ID_TREE_MOVE_TO = ID_DIALOG_MOVE_TO
         + "/navTreeView/navTree";

   // TODO:Get the id of Sample Portlet - Alerts
   public static final String ID_PORTLET_ALARMS =
         "automationName=Sample Portlet - Alerts";

   // TODO:Get the id of Sample Portlet - Recent Tasks
   public static final String ID_PORTLET_RECENT_TASKS =
         "automationName=Sample Portlet - Recent Tasks";

   // Register VC items
   public static final String ID_REGISTER_VC_LINK = "registerVcLinkButton";
   public static final String ID_UNREGISTER_VC_LINK = "unRegisterVcLinkButton";
   public static final String ID_REGISTERED_VC_TABLE = "listVcServers";
   public static final String ID_REGISTER_VC_CLIENT_URL = "clientUrlLink";

   public static final String ID_REGISTER_VC_DIALOG_BOX = "addVcView";
   public static final String ID_REGISTER_VC_SERVER_TEXTFIELD = new StringBuffer(
         ID_REGISTER_VC_DIALOG_BOX).append("/").append("serverTextInput").toString();
   public static final String ID_REGISTER_VC_USERNAME_TEXTFIELD = new StringBuffer(
         ID_REGISTER_VC_DIALOG_BOX).append("/").append("usernameTextInput").toString();
   public static final String ID_REGISTER_VC_PASSWORD_TEXTFIELD = new StringBuffer(
         ID_REGISTER_VC_DIALOG_BOX).append("/").append("passwordTextInput").toString();
   public static final String ID_REGISTER_WEB_CLIENT_SERVER_TEXTFIELD =
         new StringBuffer(ID_REGISTER_VC_DIALOG_BOX).append("/")
               .append("clientIpTextInput").toString();
   public static final String ID_REGISTER_VC_REGISTER_BUTTON = new StringBuffer(
         ID_REGISTER_VC_DIALOG_BOX).append("/").append("registerVcButton").toString();
   public static final String ID_REGISTER_VC_CANCEL_BUTTON = new StringBuffer(
         ID_REGISTER_VC_DIALOG_BOX).append("/").append("cancelButton").toString();
   public static final String ID_REGISTER_VC_ERROR_LABEL = new StringBuffer(
         ID_REGISTER_VC_DIALOG_BOX).append("/").append("error").append("/")
         .append("errWarningText").toString();

   public static final String ID_UNREGISTER_VC_DIALOG_BOX = "addVcView";
   public static final String ID_UNREGISTER_VC_SERVER_TEXTFIELD = new StringBuffer(
         ID_UNREGISTER_VC_DIALOG_BOX).append("/").append("serverTextInput").toString();
   public static final String ID_UNREGISTER_VC_USERNAME_TEXTFIELD = new StringBuffer(
         ID_UNREGISTER_VC_DIALOG_BOX).append("/").append("usernameTextInput").toString();
   public static final String ID_UNREGISTER_VC_PASSWORD_TEXTFIELD = new StringBuffer(
         ID_UNREGISTER_VC_DIALOG_BOX).append("/").append("passwordTextInput").toString();
   public static final String ID_UNREGISTER_VC_UNREGISTER_BUTTON = new StringBuffer(
         ID_UNREGISTER_VC_DIALOG_BOX).append("/").append("registerVcButton").toString();
   public static final String ID_UNREGISTER_VC_CANCEL_BUTTON = new StringBuffer(
         ID_UNREGISTER_VC_DIALOG_BOX).append("/").append("cancelButton").toString();
   public static final String ID_UNREGISTER_VC_ERROR_LABEL = new StringBuffer(
         ID_UNREGISTER_VC_DIALOG_BOX).append("/").append("error").append("/")
         .append("errWarningText").toString();

   // TODO: jmak
   // Confirm Unregister VC Dialog Box
   // public static final String ID_CONFIRM_UNREGISTER_VC_DIALOG_BOX =
   // "automationName=Confirm Unregister vCenter Server";
   // public static final String ID_CONFIRM_UNREGISTER_VC_YES_BUTTON =
   // new StringBuffer(ID_CONFIRM_UNREGISTER_VC_DIALOG_BOX).
   // append("/automationName=Yes").toString();
   // public static final String ID_CONFIRM_UNREGISTER_VC_NO_BUTTON =
   // new StringBuffer(ID_CONFIRM_UNREGISTER_VC_DIALOG_BOX).
   // append("/automationName=No").toString();

   // VM Hardware Portlet
   public static final String ID_INVENTORY_TREE_VM_NODE = "automationName=vm";
   public static final String ID_VM_SUMMARY_TAB = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_VM).append("summary.").append("container").toString();
   public static final String ID_VM_SUMMARY_TAB_PORTLETS_PANEL = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_VM).append("summaryView").toString()
         + "/rootChrome";
   public static final String ID_LABEL_VM_HARDWARE_OTHER_DEVICES =
         "automationName=Other Devices";
   public static final String ID_VM_HARDWARE_CPU_LABEL = "titleLabel.viewCpu";
   public static final String ID_VM_HARDWARE_MEMORY_LABEL = "titleLabel.viewMem";
   public static final String ID_EDIT_HW_MEMORY_CONTROL = "memoryControl";
   public static final String ID_EDIT_HW_MEMORY_VIDEO_CARD_SECTION =
         "hardwareStack/Video card /bodySection";
   public static final String ID_SCSI_CONTROLLER_LABEL_GENERIC = "titleLabel.scsi_?";
   public static final String ID_SCSI_CONTROLLER_GENERIC_COMBOBOX = "typeSetCombo";
   public static final String ID_SCSI_DEVICE_GENERIC_COMBOBOX =
         "titleLabel.SCSI device 1";
   public static final String ID_VM_HARDWARE_FLOPPY_LABEL_PREFIX =
         "titleLabel.view_Floppy_1";
   public static final String ID_VM_HARDWARE_CDDVD_LABEL_PREFIX =
         "titleLabel.view_Cdrom_1";
   public static final String ID_VM_HARDWARE_FLOPPY_LABEL_SUFFIX =
         "view_Floppy_1/headerLabel";
   public static final String ID_VM_HARDWARE_CDDVD_LABEL_SUFFIX =
         "view_Cdrom_1/headerLabel";
   public static final String ID_VM_HARDWARE_DEVICE_LABEL_SUFFIX =
         "/labelBox/deviceLabel";
   public static final String ID_VM_HARDWARE_CPU_UTILIZATION_LABEL =
         "cpuUtilization/_VmCpuView_Text1";
   public static final String ID_VM_HARDWARE_CPU_SHARES_LABEL = "cpuShares";
   public static final String ID_VM_HARDWARE_CPU_RESERVATION_LABEL = "cpuReservation";
   public static final String ID_VM_HARDWARE_CPU_LIMIT_LABEL = "cpuLimit";
   public static final String ID_VM_HARDWARE_CPU_HTSHARING_LABEL = "htSharing";
   public static final String ID_VM_HARDWARE_MEMORY_UTILIZATION_LABEL =
         "memUtilization/_VmMemView_Text1";
   public static final String ID_VM_HARDWARE_MEMORY_SHARES_LABEL = "memShares";
   public static final String ID_VM_HARDWARE_MEMORY_RESERVATION_LABEL = "memReservation";
   public static final String ID_VM_HARDWARE_MEMORY_LIMIT_LABEL = "memLimit";
   public static final String ID_VM_HARDWARE_MEMORY_HOSTOVERHEAD_LABEL = "memOverhead";
   public static final String ID_VM_HARDWARE_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "hardwareSummaryView.chrome";
   public static final String ID_VM_HARDWARE_CPU_VPMC_LABEL = "vPMCEnabled";
   public static final String ID_CPU_PHYSICAL_SOCKETS = "sockets";
   public static final String ID_WAKE_ON_LAN_LABEL = "rowWakeOnLAN";
   public static final String ID_VMOPTION_STANDBY_LABEL = "titleLabel.standbyPage";
   public static final String ID_SCSI_NOT_RECOMMENDED_WARN =
         "scsi_?/_SCSIControllerPage_Text1";
   public static final String ID_BTN_EDIT_DEFAULT_VM_COMPATIBILITY = new StringBuffer(
         "btn_").append(EXTENSION_PREFIX).append(EXTENSION_ENTITY_VM)
         .append("editDefaultHwVersionAction").toString();
   public static final String ID_CONTEXTMENU_EDIT_DEFAULT_VM_COMPATIBILITY =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_VM)
               .append("editDefaultHwVersionAction").toString();
   public static final String ID_CONTEXTMENU_CANCEL_SCHEDULED_VM_UPGRADE =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_VM)
               .append("upgradeVirtualHardwareActionCancel").toString();

   // edit button is part of Action Bar buttons, not part of VM Hardware
   public static final String ID_EDIT_HARDWARE_BUTTON_LEFT_PORTLET = ID_ACTIONS_BUTTONS
         + "/" + ID_ACTIONS_BAR + "/" + EXTENSION_PREFIX + EXTENSION_ENTITY_VM
         + "editAction";
   public static final String ID_VM_HARDWARE_NETWORKADAPTER_LABEL =
         "titleLabel.viewNIC_1";
   public static final String ID_VM_HARDWARE_NETWORKADAPTER = "titleLabel.viewNIC_{0}";
   public static final String ID_VM_HARDWARE_NETWORKADAPTER_MACADDRESS_LABEL =
         "macAddress";
   public static final String ID_VM_HARDWARE_NETWORKADAPTER_NETWORK_LINK =
         "vmNetwork/_VmNetView_LinkButton1";
   public static final String ID_VM_HARDWARE_NETWORKADAPTER_NETWORK_CONNECTION =
         "vmNetwork/_VmNetBlock_Label3";
   public static final String ID_VM_HARDWARE_HARDDISK_UTILIZATION_LABEL =
         "diskCapacity/_VmDiskBlock_Text1";
   public static final String ID_VM_HARDWARE_HARDDISK_LOCATION_STORAGE_LINK =
         "datastoreLink";
   public static final String ID_VM_HARDWARE_HARDDISK_LOCATION_STORAGE_SIZE =
         "diskLocation/_VmDiskBlock_Text2";
   public static final String ID_VM_HARDWARE_HARDDISK_LABEL = "titleLabel.viewDisk_1";

   public static final String ID_VM_HARDWARE_OTHERHARDWARE_LABEL =
         "titleLabel.viewOtherHardware";
   public static final String ID_VM_HARDWARE_OTHERHARDWARE_SCSI_LIST = "scsiAdapterList";
   public static final String ID_SCSI_CONTROLLER_CHANGE_TYPE_BUTTON =
         "scsi_?/changeTypeButton";
   public static final String ID_SCSI_CONTROLLER_DONT_CHANGE_TYPE_BUTTON =
         "scsi_?/dontChangeTypeButton";
   public static final String ID_LABEL_SCSI_CONTROLLER_DEVICE_WILL_NOT_BE_REMOVED =
         "scsi_?/_SCSIControllerPage_Label1";
   public static final String ID_LABEL_SCSI_CONTROLLER_CHANGE_TYPE_MESSAGE =
         "scsi_?/_SCSIControllerPage_Text2";
   public static final String ID_VM_HARDWARE_OTHERHARDWARE_CONTROLLER_LIST =
         "controllerList";
   public static final String ID_VM_HARDWARE_OTHERHARDWARE_INPUTDEVICE_LIST =
         "inputDeviceList";
   public static final String ID_VM_HARDWARE_HWVERSION_LABEL = "hwVersionLabel";
   public static final String ID_BUTTON_VM_POWER_ON =
         "vsphere.core.vm.powerOnAction/button";
   public static final String ID_BUTTON_VM_SHUTDOWN =
         "vsphere.core.vm.shutdownAction/button";
   public static final String ID_EDIT_VM_POPUP_OKBUTTON = "okButton";
   public static final String ID_POWEROPS_PROGRESSBAR = "powerCommand/spinner";
   public static final String ID_EDIT_HW_MEMORY_LABEL = "titleLabel.memoryPage";
   public static final String ID_EDIT_HW_MEMORY_COMBO_COLLAPSED =
         "memoryControlHeader/topRow/comboBox";
   public static final String ID_EDIT_HW_MEMORY_COMBO = "memoryControl/topRow/comboBox";
   public static final String ID_EDIT_HARDWARE_POPUP = "tiwoDialog";
   public static final String ID_EDIT_HW_OK_BUTTON = "vmConfigForm/okButton";
   public static final String ID_EDIT_HW_LIMIT_COMBO =
         "memoryPage/limitControl/topRow/comboBox";
   public static final String ID_EDIT_HW_RESERVATION_COMBO =
         "memoryPage/reservationControl/topRow/comboBox";
   public static final String ID_EDIT_HW_SHARES_COMBO =
         "memoryPage/sharesControl/levels";
   public static final String ID_EDIT_HW_SHARES_STEPPER =
         "memoryPage/sharesControl/numShares";
   public static final String ID_POWER_STATE_LABEL = "_powerStateLabel";
   public static final String ID_POWERED_OFF_STATE = "_poweredOffLabel";
   public static final String ID_POWER_COMMAND_LINK = "powerCommand/link";
   public static final String ID_POWER_ON_ICON = EXTENSION_PREFIX + EXTENSION_ENTITY_VM
         + "powerOnAction";
   public static final String ID_POWER_OFF_ICON = EXTENSION_PREFIX + EXTENSION_ENTITY_VM
         + "powerOffAction";
   public static final String ID_SUSPEND_ICON = EXTENSION_PREFIX + EXTENSION_ENTITY_VM
         + "suspendAction";
   public static final String ID_EDIT_HARDWARE_SHORT_ICON =
         "banner/outerBox/iconHolder/theIcon";
   public static final String ID_EDIT_HARDWARE_SHORT_LABEL =
         "outerBox/contentsBox/shortLabel";
   public static final String ID_EDIT_HARDWARE_MORE_LINK =
         "contentsBox/topLineBox/moreLink";
   public static final String ID_EDIT_HARDWARE_ERROR_DETAILS_LABEL =
         "contentsBox/details";
   public static final String ID_EDIT_HARDWARE_ERROR_CLOSE_BUTTON =
         "contentsBox/topLineBox/xButton";
   public static final String ID_ACTIONS_DIALOG_BOX = "appActions";
   public static final String ID_EDIT_HW_CPU_CPUIDMASK_COMBO = "cpuIdCombo";
   public static final String ID_EDIT_HW_CPU_CPUMMU_COMBO = "cpuMmuCombo";
   public static final String ID_VMCI_LABEL = "titleLabel.VMCI device";
   public static final String ID_VMCI_CHECK_BOX = "enableVMCICheckBox";
   public static final String ID_LABEL_VMCI_DEPRICATED_MESSAGE = "_VMCIDevicePage_Text2";
   public static final String ID_LABEL_VMCI_OBSOLETE = "automationName=Obsolete";
   public static final String ID_LABEL_ITEMLABEL = "labelItem";
   public static final String ID_ENT_BIOS_CHECKBOX = "enterBIOSSetupCB";
   public static final String ID_BOOT_RETRY_CHECKBOX = "bootRetryEnabledCB";
   public static final String ID_BOOT_RECOVERY_STEPPER = "bootRetryDelayNum";
   public static final String ID_EDIT_HW_MEMORY_HOTADD_CHECK_BOX = "memoryHotAdd";
   public static final String ID_EDIT_HW_MEMORY_UNIT = "memoryControl/topRow/units";
   public static final String ID_EDIT_HW_MEMORY_VIDEO_CARD_LABEL =
         "titleLabel.Video card ";
   public static final String ID_BUTTON_EDIT_HW_VIDEO_MEMORY_CALC =
         "memoryCalculatorBtn";
   public static final String ID_DIALOG_VIDEO_MEMORY_CALC =
         "className=VideoMemoryCalculator";
   public static final String ID_COMBO_VIDEO_CARD_DISPLAY_NUM =
         "className=VideoMemoryCalculator/numberOfDisplays/className=ComboBox";
   public static final String ID_COMBO_VIDEO_CARD_RESOLUTION =
         "className=VideoMemoryCalculator/resolution/className=ComboBox";
   public static final String ID_LABEL_VIDEO_CARD_TOTAL_MEMORY =
         "className=VideoMemoryCalculator/totalVideoMemory/className=Label";
   public static final String ID_BUTTON_VIDEO_MEMCALC_CANCEL =
         "className=VideoMemoryCalculator/cancel";
   public static final String ID_BUTTON_VIDEO_MEMCALC_OK =
         "className=VideoMemoryCalculator/ok";
   public static final String ID_EDIT_HW_MEMORY_VIDEO_SETTING_COMBO_BOX =
         "autoDetectCombo";
   public static final String ID_EDIT_HW_VIDEO_MEMORY_CALCULATOR_BUTTON =
         "memoryCalculatorBtn";
   public static final String ID_EDIT_HW_VIDEO_MEMORY_CALCULATOR_CONTENTS =
         "contentGroup";
   public static final String ID_EDIT_HW_VIDEO_MEMORY_INCREASED_TEXT = "autoIncrease";
   public static final String ID_VIDEO_3D_MEMORY_CALCULATOR_WARNING =
         "automationName=Extra 64MB plus 16MB per additional display was added to accommodate 3D graphics";
   public static final String ID_VIDEO_MEMORY_CALCULATOR_DISPLAY_COMBO_BOX =
         "accessibilityName=Number of displays:";
   public static final String ID_VIDEO_MEMORY_CALCULATOR_RESOLUTION_COMBO_BOX =
         "accessibilityName=Resolution:";
   public static final String ID_VIDEO_MEMORY_CALCULATOR_OKBUTTON = "automationName=OK";
   public static final String ID_VIDEO_MEMORY_CALCULATOR_CANCELBUTTON =
         "automationName=Cancel";
   public static final String ID_EDIT_HW_VIDEO_3D_RENDERER_COMBO_BOX = "use3DRenderer";
   public static final String ID_EDIT_HW_VIDEO_3D_GRAPHICS = "enable3DSupport_";
   public static final String ID_EDIT_HW_VIDEO_MEMORY_SIZE = "videoRamSizeInKB_";
   public static final String ID_EDIT_HW_NUM_DISPLAY_COMBO_BOX = "numDisplaysCombo";
   public static final String ID_EDIT_HW_VIDEO_RAM_TEXT_BOX = "videoRamInput";
   public static final String ID_LABEL_EDIT_HW_VIDEO_MEMORY_UNIT =
         "_VideoCardPage_Label1";
   public static final String ID_EDIT_HW_VIDEO_DISPLAY_TEXT = new StringBuffer(
         ID_EDIT_HW_VIDEO_RAM_TEXT_BOX).append("textDisplay").toString();
   public static final String ID_EDIT_HW_ENABLE_3D_CHECK_BOX = "enable3D";
   public static final String ID_VM_CDDVD_CONNECTION_BUTTON = "view_Cdrom_1/menuButton";
   public static final String ID_VM_FLOPPY_CONNECTION_BUTTON =
         "view_Floppy_1/menuButton";
   public static final String ID_VM_CDDVD_CONNECTED_TO_HOST_CDROM_DEVICE =
         "automationName=Host CD-ROM device";
   public static final String ID_VM_CDDVD_OPTION_HOST_CD_DEVICE =
         "automationName=Connect to host CD device...";
   public static final String ID_VM_CDDVD_OPTION_DATASTORE_ISO =
         "automationName=Connect to CD/DVD image on a datastore...";
   public static final String ID_VM_FLOPPY_OPTION_DATASTORE_ISO =
         "automationName=Connect to floppy image on a datastore...";
   public static final String ID_VM_FLOPPY_OPTION_HOST_FLOPPY =
         "automationName=Connect to host floppy device...";
   public static final String ID_VM_CDDVD_CONNECTION_OPTIONS_MENU = "menu_view_Cdrom_1";
   public static final String ID_VM_FLOPPY_CONNECTION_OPTIONS_MENU =
         "menu_view_Floppy_1";
   public static final String ID_EDIT_HW_NO_PERMISSON_MESSAGE = "_message";
   public static final String ID_SCSI_CONTROLLER_TYPE = "controllerType";
   public static final String ID_LINK_VM_HARDWARE_NETWORK_ADAPTER = "headerNetworkLink";
   public static final String ID_PORTLET_LINK_VM_HARDWARE_NETWORK_ADAPTER =
         "viewNIC_{0}/headerNetworkLink";
   public static final String ID_PORTLET_VM_HARDWARE_NETWORK_ADAPTER_NETWORK =
         "viewNIC_{0}/bodyNetworkLink";
   public static final String ID_PORTLET_VM_HARDWARE_NETWORK_ADAPTER_STATE =
         "viewNIC_{0}/_VmNetBlock_Label3";
   public static final String ID_LINK_VM_HARDWARE_NETWORK_ADAPTER_NETWORK =
         "bodyNetworkLink";

   public static final String ID_NETWORK_ADAPTER = "Network adapter ";
   // Annotations portlet
   public static final String ID_PORTLET_MINMAX_BUTTON = "minMaxButton";
   public static final String ID_VM_ANNOTATION_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "AnnotationsView";
   public static final String ID_VM_ANNOTATION_EDIT_LINK = "editButton";
   public static final String ID_VM_EDIT_ANNOTATIONS_PANEL = "tiwoDialog";
   public static final String ID_VM_EDIT_ANNOTATIONS_VALUE_GRID = "attributesGrid";
   public static final String ID_VM_EDIT_ANNOTATIONS_TEXTAREA = "notesTextArea";
   public static final String ID_EDIT_ANNOTATIONS_POPUP_ADDLINK = "addButton";
   public static final String ID_EDIT_ANNOTATIONS_POPUP_REMOVELINK = "removeButton";
   public static final String ID_EDIT_ANNOTATIONS_POPUP_OKBUTTON = "okButton";
   public static final String ID_EDIT_ANNOTATIONS_POPUP_CANCELBUTTON = "cancelButton";

   // Related Items portlet
   public static final String ID_VM_RELATED_ITEMS_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "relatedItemsView.chrome";
   public static final String ID_VM_SUMMARY_VAPP_DETAILS_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "vAppDetailsView.chrome";
   public static final String ID_VM_SUMMARY_VAPP_PRODUCTNAME_LABEL =
         "productUrlLabel/plainLabel";
   public static final String ID_VM_SUMMARY_VAPP_VENDOR_LABEL =
         "vendorUrlLabel/plainLabel";
   public static final String ID_VM_SUMMARY_VAPP_VERSION_LABEL = "versionText";

   // Fault Tolerance Portlet
   public static final String ID_GOS_VM_FT_SUMMARY_PORTLET =
         "vsphere.core.vm.ftView.chrome";
   public static final String ID_HOST_FT_PORTLET =
         "vsphere.core.host.summary.ft.hostFTView.chrome";
   public static final String ID_HOST_FT_VERSION_LABEL = "ftVersionLabel";
   public static final String ID_HOST_FT_PORTLET_POWERED_ON_PRIMARY_VM_COUNT_LABEL =
         "poweredOnPrimaryVMsLabel";
   public static final String ID_HOST_FT_PORTLET_POWERED_ON_SECONDARY_VM_COUNT_LABEL =
         "poweredOnSecondaryVMsLabel";
   public static final String ID_HOST_FT_PORTLET_SECONDARY_VM_COUNT_LABEL =
         "secondaryVMsCountLabel";
   public static final String ID_HOST_FT_PORTLET_PRIMARY_VM_COUNT_LABEL =
         "primaryVMsCountLabel";
   public static final String ID_HOST_FT_PORTLET_VERSION_NUMBER = "ftVersionLabel";
   public static final String ID_GOS_VM_FT_PORTLET_STATUS_LABEL = "ftStatusLabel";

   public static final String ID_HOST_FT_PORTLET_POWERED_ON_PRIMARY_VM_COUNT_TEXT =
         "automationName=Powered On Primary VMs";
   public static final String ID_HOST_FT_PORTLET_POWERED_ON_SECONDARY_VM_COUNT_TEXT =
         "automationName=Powered On Secondary VMs";
   public static final String ID_HOST_FT_PORTLET_SECONDARY_VM_COUNT_TEXT =
         "automationName=Total Secondary VMs";
   public static final String ID_HOST_FT_PORTLET_PRIMARY_VM_COUNT_TEXT =
         "automationName=Total Primary VMs";
   public static final String ID_HOST_FT_PORTLET_VERSION_NUMBER_TEXT =
         "automationName=Fault Tolerance Version";
   public static final String ID_GOS_VM_FT_PORTLET_POWERED_ON_PRIMARY_VM_LABEL =
         "poweredOnPrimaryVMsLabel";
   public static final String ID_GOS_VM_FT_PORTLET_STATE_LABEL = "ftStateLabel";
   public static final String ID_GOS_VM_FT_PORTLET_SECONDARY_LOCATION =
         "secondaryLocationLinkBtn";
   public static final String ID_GOS_VM_FT_PORTLET_SECONDARY_CPU =
         "totalSecondaryCpuLabel";
   public static final String ID_GOS_VM_FT_PORTLET_SECONDARY_MEMORY =
         "totalSecondaryMemoryLabel";
   public static final String ID_GOS_VM_FT_PORTLET_SECONDARY_LATENCY =
         "secondaryLatencyLabel";
   public static final String ID_GOS_VM_FT_PORTLET_LOG_BANDWIDTH = "logBandwidthLabel";
   public static final String ID_TASK_FAILURE_FT_POPUP_LIST = "list";

   // Custom Attributes Portlet
   public static final String ID_VM_CUSTOM_ATRRIBUTES_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "summary.customAttributesView.chrome";

   // Notes Portlet
   public static final String ID_VM_NOTES_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "summary.annotationsNotesView.chrome";

   // Storage Profiles Portlet
   public static final String ID_VM_STORAGE_PROFILE_PORTLET = EXTENSION_PREFIX + "spbm."
         + EXTENSION_ENTITY_VM + "storageProfilesView.chrome";

   // Advanced Configuration Portlet
   public static final String ID_VM_ADVANCED_CONFIG_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "advancedConfigurationView.chrome";

   // Tags Portlet
   public static final String ID_SUMMARY_TAGS_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY + "summary.tagsView.chrome";

   // Notes Portlet
   public static final String ID_SUMMARY_NOTES_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY + "summary.annotationsView.chrome";

   // Details Portlet
   public static final String ID_SUMMARY_DETAILS_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY + "summary.detailsView.chrome";

   // Policies Portlet
   public static final String ID_POLICIES_PORTLET = EXTENSION_PREFIX + EXTENSION_ENTITY
         + "summary.policyView.chrome";

   // Confirm Power Off Dialog Box
   // BUG 503689: Confirm Power Off dialog box's items don't have id
   public static final String ID_CONFIRM_POWER_OFF_YES_BUTTON =
         "automationName=Confirm Power Off/automationName=Yes";
   public static final String ID_CONFIRM_POWER_OFF_NO_BUTTON =
         "automationName=Confirm Power Off/automationName=No";

   // Confirm Suspend Dialog Box
   public static final String ID_CONFIRM_SUSPEND_YES_BUTTON =
         "automationName=Confirm Suspend/automationName=Yes";
   public static final String ID_CONFIRM_SUSPEND_NO_BUTTON =
         "automationName=Confirm Suspend/automationName=No";

   // Confirm Reset dialog box
   public static final String ID_CONFIRM_RESET_YES_BUTTON =
         "automationName=Confirm Reset/automationName=Yes";
   public static final String ID_CONFIRM_RESET_NO_BUTTON =
         "automationName=Confirm Reset/automationName=No";

   // Confirm Logout dialog box
   public static final String ID_CONFIRM_LOGOUT_YES_BUTTON =
         "automationName=Confirm Logout/automationName=Yes";
   public static final String ID_CONFIRM_LOGOUT_NO_BUTTON =
         "automationName=Confirm Logout/automationName=No";

   // Confirm Upgrade HW dialog box
   public static final String ID_CONFIRM_VM_UPGRADE_YES_BUTTON =
         "automationName=Confirm Virtual Machine Upgrade/automationName=Yes";
   public static final String ID_CONFIRM_VM_UPGRADE_NO_BUTTON =
         "automationName=Confirm Virtual Machine Upgrade/automationName=No";

   // Confirm Upgrade HW dialog box
   public static final String ID_CONFIRM_RESTART_YES_BUTTON =
         "automationName=Confirm Restart/automationName=Yes";
   public static final String ID_CONFIRM_RESTART_NO_BUTTON =
         "automationName=Confirm Restart/automationName=No";
   public static final String ID_CONFIRM_SHUTDOWN_YES_BUTTON =
         "automationName=Confirm Shut Down/automationName=Yes";
   public static final String ID_CONFIRM_SHUTDOWN_NO_BUTTON =
         "automationName=Confirm Shut Down/automationName=No";
   public static final String ID_BUTTON_VM_COMPATIBILITY = "comboBox";

   // Confirm Host Disconnect dialog box
   public static final String ID_CONFIRM_HOST_DISCONNECT_YES_BUTTON =
         "automationName=Confirm Disconnect/automationName=Yes";

   // Confirm Host Reconnect dialog box
   public static final String ID_CONFIRM_HOST_RECONNECT_YES_BUTTON =
         "automationName=Reconnect host/automationName=Yes";

   // Confirm Host Remove dialog box
   public static final String ID_CONFIRM_HOST_REMOVE_YES_BUTTON = "automationName=Yes";
   public static final String ID_CONFIRM_HOST_ADD_YES_BUTTON =
         "YesNoDialog/automationName=Yes";
   public static final String ID_CONFIRM_HOST_ADD_NO_BUTTON =
         "YesNoDialog/automationName=No";
   public static final String ID_CONFIRM_HOST_ADD_ALERT_DIALOG = "YesNoDialog";
   public static final String ID_HOST_SUMMARY_STATUS_LABEL =
         "summary_stateProp_valueLbl";

   // Host State Summary Tab
   public static final String ID_HOST_STATE_STATUS = "summary_stateProp_valueLbl";
   public static final String ID_HOST_STATE_CONNECTED = "Connected";
   public static final String ID_HOST_STATE_DISCONNECTED = "Disconnected";
   public static final String ID_HOST_STATE_MAINTENANCE_MODE = "Maintenance Mode";
   public static final String ID_HOST_STATE_NOT_RESPONDING = "Not responding";

   // Reboot Host Variables
   public static final String ID_REBOOTHOST_DIALOG_BOX = "wizardContent";
   public static final String ID_LABEL_REBOOTHOST_DIALOG_BOX_WARNING = "warning";
   public static final String ID_LABEL_REBOOTHOST_DIALOG_BOX_WARNING_TEXT =
         "This host is not in maintenance mode.";
   public static final String ID_LABEL_REBOOTHOST_DIALOG_BOX_CONFIRMATION =
         "confirmation";
   public static final String ID_LABEL_REBOOTHOST_DIALOG_BOX_CONFIRMATION_TEXT =
         "Are you sure you want to reboot the selected host?";
   public static final String ID_LABEL_REBOOTHOST_DIALOG_BOX_LOGREASON = "logReason";
   public static final String ID_LABEL_REBOOTHOST_DIALOG_BOX_LOGREASON_TEXT =
         "Log a reason for this reboot operation:";
   public static final String ID_TEXTINPUT_REBOOTHOST_DIALOG_BOX_TESTDISPLAY = "reason";
   public static final String ID_TEXTINPUT_REBOOTHOST_DIALOG_BOX_TESTDISPLAY_TEXT =
         "Enter a reason here";

   // Add Host Variables
   public static final String ID_ADDHOST_CONNECTION_PAGE = "hostConnectionPage";
   public static final String ID_LABEL_ADDHOST_INPUT_LABEL = "inputLabel";
   public static final String ID_LABEL_ADDHOST_RESOURCE_POOL =
         "drsClusterRpPage/headerLabel";
   public static final String ID_ADDHOST_HOSTNAME_TEXTINPUT =
         "hostNameAndLocationPage/inputLabel";
   public static final String ID_ADDHOST_USERNAME_TEXTINPUT = "userName";
   public static final String ID_ADDHOST_PASSWORD_TEXTINPUT = "password";
   public static final String ID_ADDHOST_NEXT_BUTTON = "next";
   public static final String ID_ADDHOST_FINISH_BUTTON = "finish";
   public static final String ID_ASSIGN_LICENSE_HOST_DETAILS =
         "_AssignLicenseKeyDetailsView_Label1";
   public static final String ID_ASSIGN_LICENSE_KEY_LIST = "existingKeysList";
   public static final String ID_ASSIGN_LICENSE_COMBO = "openButton";
   public static final String ID_ENABLE_LOCKDOWN_MODE = "enableLockDown";
   public static final String ID_VM_LOCATION_SEARCH_INPUT = "searchInput";
   public static final String ID_ADD_HOST_WARNING =
         "hostSummaryPage/alreadyManagedMessage";
   public static final String ID_ADD_HOST_FAILED_MESSAGE =
         "validationMessageContainer/_message";
   public static final String ID_REVIEW_SETTINGS_VENDER_NAME = "txtVendor";
   public static final String ID_REVIEW_SETTINGS_MODEL_NAME = "txtModel";
   public static final String ID_REVIEW_SETTINGS_HOST_VERSION_NAME =
         "_HostSummaryPage_PropertyGridRow4/txtHostVersion";
   public static final String ID_REVIEW_SETTINGS_HOST_VERSION_NAME1 =
         "_ConfirmationPage_PropertyGridRow2/txtHostVersion";
   public static final String ID_REVIEW_SETTINGS_HOST_NAME =
         "_HostSummaryPage_PropertyGridRow1/txtHostName";
   public static final String ID_REVIEW_SETTINGS_HOST_NAME1 =
         "_ConfirmationPage_PropertyGridRow1/txtHostName";
   public static final String ID_RADIO_BUTTON_ROOT_RP = "rbRootRp";
   public static final String ID_RADIO_BUTTON_NEW_RP = "rbNewRp";
   public static final String ID_TEXT_INPUT_NEW_RP_NAME = "newRpName";
   public static final String ID_LABEL_HOST_SUMMARY_TEXT = "txtHostName";
   public static final String ID_TEXT_RESOURCES_DESTINATION_TEXT = "txtRp";
   public static final String ID_LABEL_HOST_SUMMARY_TOC = "step_hostSummaryPage";
   public static final String ID_LABEL_CONNECTION_SETTINGS_TOC =
         "step_hostConnectionPage";
   public static final String ID_LABEL_CONNECTION_SETTINGS_TEXT =
         "_HostConnectionPage_Text1";
   public static final String ID_LABEL_SECURITY_PROFILE_LOCKDOWN_MODE =
         "vsphere.core.host.manage.settingsView/tocTree/automationName=Security Profile";
   public static final String ID_NAVTREE_VMLOCATIONPAGE = "vmLocationPage/navTree";
   public static final String ID_LABEL_NATIVE_HOST_NAME = "localhost";
   public static final String ID_ADDHOST_CONFIRMATION_PAGE = "confirmationPage";

   // Remove host variables
   public static final String ID_ALERT_FORM_TEXTFIELD =
         "className=AlertForm/className=UITextField";
   public static final String ID_ALERT_FORM_OK_BUTTON =
         "className=AlertForm/automationName=OK";

   // Host Info
   public static final String ID_HOST_STATUS_PORTLET =
         "vsphere.core.host.summary.statusView.chrome";
   public static final String ID_HOST_HARDWARE_PORTLET =
         "vsphere.core.host.summary.hardwareView.chrome";
   public static final String ID_HOST_CONFIGURATION_PORTLET =
         "vsphere.core.host.summary.configurationView.chrome";
   public static final String ID_HOST_ANNOTATION_PORTLET =
         "vsphere.core.host.summary.annotationsView.chrome";
   public static final String ID_HOST_STATUS_GENERAL = "General";
   public static final String ID_HOST_ANNOTATION_LABEL =
         "vsphere.core.host.summary.annotationsView/stackEditor/titleLabel";
   public static final String ID_HOST_ANNOTATIONS_EDIT_BUTTON =
         "vsphere.core.host.summary.annotationsView/editButton";
   public static final String ID_HOST_ANNOTATIONS_ADD_BUTTON =
         "editAnnotationView/addButton";
   public static final String ID_HOST_ANNOTATIONS_REMOVE_BUTTON =
         "editAnnotationView/removeButton";
   public static final String ID_HOST_ANNOTATIONS_OK_BUTTON =
         "editAnnotationView/okButton";
   public static final String ID_HOST_ANNOTATIONS_CANCEL_BUTTON =
         "editAnnotationView/cancelButton";
   public static final String ID_HOST_ANNOTATIONS_ATTRIBUTES_GRID =
         "editAnnotationView/attributesGrid";
   public static final String ID_HOST_HARDWARE_MANUFACTURER =
         "_HostHardwareView_StackBlock1/titleLabel";
   public static final String ID_HOST_HARDWARE_MODEL =
         "_HostHardwareView_StackBlock2/titleLabel";
   public static final String ID_HOST_HARDWARE_CPU =
         "titleLabel._HostHardwareView_StackBlock3";
   public static final String ID_HOST_HARDWARE_MEMORY =
         "titleLabel._HostHardwareView_StackBlock4";
   public static final String ID_HOST_HARDWARE_NETWORKING = "titleLabel.networkBlock";
   public static final String ID_HOST_HARDWARE_STORAGE = "titleLabel.storageBlock";
   public static final String ID_HOST_HARDWARE_CPU_CPUCORES = "CPU Cores";
   public static final String ID_HOST_HARDWARE_CPU_PROCESSORTYPE = "Processor Type";
   public static final String ID_HOST_HARDWARE_CPU_SOCKETS = "Sockets";
   public static final String ID_HOST_HARDWARE_CPU_CORESPERSOCKET = "Cores per Socket";
   public static final String ID_HOST_HARDWARE_CPU_LOGICALPROCESSORS =
         "Logical Processors";
   public static final String ID_HOST_HARDWARE_CPU_HYPERTHREADING = "Hyperthreading";
   public static final String ID_HOST_HARDWARE_MEMORY_SYSTEM = "System";
   public static final String ID_HOST_HARDWARE_MEMORY_VIRTUALMACHINES =
         "Virtual Machines";
   public static final String ID_HOST_HARDWARE_MEMORY_CAPACITY = "Capacity";
   public static final String ID_HOST_HARDWARE_NETWORKING_ROUTING = "Routing";
   public static final String ID_HOST_HARDWARE_NETWORKING_NETWORKS = "Networks";
   public static final String ID_HOST_HARDWARE_NETWORKING_PHYSICALADAPTERS =
         "Physical Adapters";
   public static final String ID_HOST_HARDWARE_STORAGE_DATASTORES = "Datastores";
   public static final String ID_HOST_HARDWARE_STORAGE_PHYSICALADAPTERS =
         "Physical Adapters";
   public static final String ID_HOST_CONFIGURATION_ESXVERSION = "ESX Version";
   public static final String ID_HOST_CONFIGURATION_VMOTIONENABLED = "VMotion Enabled";
   public static final String ID_HOST_CONFIGURATION_CONFIGUREDFORFT =
         "_HostConfigurationView_PropertyGridRow3";
   public static final String ID_VM_SUMMARY_TAB_NETWORK =
         "vsphere.core.vm.summaryView.container";
   public static final String ID_HOST_CPU_VIEW_PERFORMANCE = "cpuLink";
   public static final String ID_HOST_CPU_VIEW_PERFORMANCE_ALERT =
         "text=Cpu -> View Performance...";
   public static final String ID_HOST_MEMORY_VIEW_PERFORMANCE = "memLink";
   public static final String ID_HOST_MEMORY_VIEW_PERFORMANCE_ALERT =
         "text=Memory -> View Performance...";
   public static final String ID_HOST_STORAGE_STORAGE_PODS = "spods";
   public static final String ID_HOST_VIRTUAL_SWITCHES_VIEW =
         "vsphere.core.host.manage.settings.virtualSwitchesView";
   public static final String ID_HOST_VIRTUAL_SWITCH_MANAGE_DIAGRAM = "hpsDiagram";
   public static final String ID_HOST_MANAGE_SETTINGS_NETWORKING_TOC = "tocTree";
   public static final String ID_HOST_VDS_STACK_BLOCK_GENERAL =
         "_HostProxySwitchDetailsView_StackBlock1";
   public static final String ID_HOST_VDS_PROPERTIES_NAME =
         "_HostProxySwitchDetailsView_LabelEx1";
   public static final String ID_HOST_VDS_PROPERTIES_UPLINKS =
         "_HostProxySwitchDetailsView_LabelEx5";
   public static final String ID_HOST_VDS_PROPERTIES_NETIOCONTROL =
         "_HostProxySwitchDetailsView_LabelEx7";
   public static final String ID_HOST_VDS_PROPERTIES_MAX_MTU =
         "_HostProxySwitchDetailsView_LabelEx8";
   public static final String ID_HOST_VDS_STACK_BLOCK_ADVANCED =
         "_HostProxySwitchDetailsView_StackBlock2";
   public static final String ID_SEARCH_STACK = "searchStack";
   // Status Portlet
   public static final String ID_HEALTH_STATE_LABEL = "entityStatus/statusLabel";
   public static final String ID_VM_SUMMARY_STATUS_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "statusView.chrome";
   public static final String ID_NODE_LABEL = "title";
   public static final String ID_VM_SUMMARY_RELATEDITEMS_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "relatedItemsView.chrome";
   public static final String ID_VM_SUMMARY_RELATEDITEMS_DATASTORE_LIST =
         "RelatedItemsServerObjectList3";
   public static final String ID_VM_SUMMARY_RELATEDITEMS_NETWORK_LIST =
         "RelatedItemsServerObjectList2";
   public static final String ID_VM_SUMMARY_RELATEDITEMS_HOST_LIST =
         "RelatedItemsServerObjectList0";
   public static final String ID_VM_SUMMARY_RELATEDITEMS_RP_LIST =
         "RelatedItemsServerObjectList1";

   public static final String ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID =
         "headerContainer/wizardPageNavigator/step.id=";

   // Standard Switch
   public static final String ID_STANDARD_SWITCH_STATUS_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_STANDARD_NETWORK + "statusView.chrome";
   public static final String ID_STANDARD_SWITCH_NETWORK_DETAILS_PORTLET =
         EXTENSION_PREFIX + EXTENSION_ENTITY_STANDARD_NETWORK + MAIN_TAB_PREFIX
               + ".detailsView.chrome";
   public static final String ID_STANDARD_SWITCH_STATUS_OVERALL_LABEL =
         "entityStatus/statusLabel";
   public static final String ID_STANDARD_SWITCH_NETWORK_DETAILS_ACCESSIBLE =
         "accessible";
   public static final String ID_STANDARD_SWITCH_NETWORK_DETAILS_IP_POOL = "ipPool";
   public static final String ID_STANDARD_SWITCH_NETWORK_DETAILS_VMS = "vms";
   public static final String ID_STANDARD_SWITCH_NETWORK_DETAILS_HOSTS = "hosts";
   public static final String ID_ACTION_EDIT_GENERAL_SETTINGS = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + EXTENSION_ENTITY_STANDARD_NETWORK
         + "editGeneralSettingsAction";
   public static String ID_STANDARD_SWITCH_EDIT_ADVANCED_SETTINGS_BUTTON = "btn_"
         + EXTENSION_PREFIX + EXTENSION_ENTITY_HOST + EXTENSION_ENTITY_STANDARD_NETWORK
         + "editAdvancedSettingsAction";
   public static final String ID_ACTION_ADD_NETWORKING = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + EXTENSION_ENTITY_STANDARD_NETWORK
         + "addNetworkingAction";
   public static final String ID_ACTION_MOVE_NETWORK = EXTENSION_PREFIX
         + EXTENSION_ENTITY_STANDARD_NETWORK + "moveAction";
   public static final String ID_HOST_NETWORKING_LIST = "tocTree";

   // Standard Switch Add Networking wizard
   public static final String ID_STANDARD_SWITCH_VMKERNEL_TYPE_RADIO =
         "vmkernelTypeRadioBtn";
   public static final String ID_STANDARD_SWITCH_CONSOLE_TYPE_RADIO =
         "consoleTypeRadioBtn";
   public static final String ID_STANDARD_SWITCH_PNIC_TYPE_RADIO = "pnicTypeRadioBtn";
   public static final String ID_STANDARD_SWITCH_VM_TYPE_RADIO = "vmTypeRadioBtn";
   public static final String ID_STANDARD_SWITCH_NEW_SWITCH_RADIO = "newSwitchRadioBtn";
   public static final String ID_STANDARD_SWITCH_EXISTING_SWITCH_RADIO =
         "existingSwitchRadioBtn";
   public static final String ID_BUTTON_ASSIGN_BUTTON = "addButton";
   public static final String ID_STANDARD_SWITCH_ADD_NETWORKING_WIZARD_CONNECTION_TYPE_PAGE =
         "connectionTypePage";
   public static final String ID_STANDARD_SWITCH_ADD_NETWORKING_WIZARD_CONNECTION_TARGET_PAGE =
         "connectionTargetPage";
   public static final String ID_STANDARD_SWITCH_ADD_NETWORKING_WIZARD_NEW_SWITCH_PAGE =
         "newSwitchPage";
   public static final String ID_STANDARD_SWITCH_ADD_NETWORKING_WIZARD_VMKERNEL_PROPERTIES_PAGE =
         "portPropertiesPage";
   public static final String ID_STANDARD_SWITCH_ADD_NETWORKING_WIZARD_VMPORTGROUP_PROPERTIES_PAGE =
         "vmNetworkingPage";
   public static final String ID_STANDARD_SWITCH_NETWORK_LABEL = "networkLabelInput";
   public static final String ID_STANDARD_SWITCH_VLANID = "vlanIdList";
   public static final String ID_STANDARD_SWITCH_VMOTION_CHECKBOX = "vMotionChk";
   public static final String ID_STANDARD_SWITCH_FT_CHECKBOX = "ftChk";
   public static final String ID_STANDARD_SWITCH_MANAGEMENT_TRAFFIC_CHECKBOX =
         "mngTrafficCkb";
   public static final String ID_NETWORKING_DIAGRAM_SELECT_SWITCH_COMBO_BOX =
         "virtualSwitchList";
   public static final String ID_STANDARD_SWITCH_DIAGRAM = "vssDiagram";
   public static final String ID_NETWORKING_DIAGRAM_SELECT_ITEM =
         "_VirtualSwitchesView_Text1";
   public static final String ID_NETWORKING_MANAGEMENT_TOCTREE =
         "vsphere.core.host.manage.networkingView/tocTree";
   public static final String ID_EDIT_VSWITCH_BUTTON = "editVswitchSettingsAction";
   public static final String ID_EDIT_PORTGROUP_BUTTON = "editPortgroupSettingsAction";
   public static final String ID_EDIT_VNIC_BUTTON = "editVnicSettingsAction";
   public static final String ID_EDIT_PNIC_BUTTON = "editPnicSettingsAction";
   public static final String ID_REMOVE_VSWITCH_BUTTON = "removeVswitchAction";
   public static final String ID_REMOVE_PORTGROUP_BUTTON = "removePortgroupAction";
   public static final String ID_REMOVE_VNIC_BUTTON = "removeVnicAction";
   public static final String ID_NETWORK_LABEL = "networkLabel";
   public static final String ID_NUM_PORTS_LIST = "numPortsList";
   public static final String ID_MTU_STEPPER = "mtuNumStepper";
   public static final String ID_NUM_PORTS_LABEL = "numPortsLabel";
   public static final String ID_MTU_LABEL = "mtuLabel";
   public static final String ID_BUTTON_MANAGE_PHYSICAL_ADAPTER =
         "_SwitchButtonBar_ActionButton1";
   public static final String ID_BUTTON_MIGRATE_NETWORKING =
         "_SwitchButtonBar_ActionButton3";
   public static final String ID_REMOVE_VSWITCH_DIALOG = "YesNoDialog";
   public static final String ID_REMOVE_YES_BUTTON = ID_REMOVE_VSWITCH_DIALOG
         + "/automationName=Yes";
   public static final String ID_REMOVE_NO_BUTTON = ID_REMOVE_VSWITCH_DIALOG
         + "/automationName=No";
   public static final String ID_LABEL_SELECTION_CONNECTION_STEP_ID =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "step_connectionTypePage";
   public static final String ID_LABEL_CONNECTION_TARGET_STEP_ID =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "step_connectionTargetPage";
   public static final String ID_LABEL_NEW_SWITCH_STEP_ID =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "step_newSwitchPage";
   public static final String ID_LABEL_SUMMARY_STEP_ID =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "step_defaultSummaryPage";
   public static final String ID_LABEL_SELECT_PNIC = "_FailoverOrderEditor_Text1";
   public static final String ID_BUTTON_MOVE_UP = "automationName=Move Up";
   public static final String ID_BUTTON_MOVE_DOWN = "automationName=Move Down";
   public static final String ID_BUTTON_REMOVE = "automationName=Remove";
   public static final String ID_LIST_ASSIGNED_PNIC = "list";
   public static final String ID_LIST_UNCLAIMED_PNIC = "unclaimedList";
   public static final String ID_DIALOG_ASSIGNED_PNIC = "assignDialog";
   public static final String ID_LIST_FAILOVER_GROUP = "failoverGroup";
   public static final String ID_DATA_PROVIDER_SOURCE_SOURCE =
         "dataProvider.source.source.";
   public static final String ID_DATA_PROVIDER_PNIC_LIST_ADAPTERS =
         ID_DATA_PROVIDER_SOURCE_SOURCE + "0.adapters.";
   public static final String ID_DATA_PROVIDER_UPLINKS_LISTS =
         ID_DATA_PROVIDER_SOURCE_SOURCE + "0.items";
   public static final String ID_DATA_PROVIDER_VIRTUAL_SWITCH = ".virtualSwitch";
   public static final String ID_DATA_PROVIDER_DEVICE = ".device";
   public static final String ID_LIST_PNIC = "pnicList";
   public static final String ID_HOST_NETWORKING_PNIC_GRID = ID_LIST_PNIC + "/list";
   public static final String ID_HOST_NETWORKING_PNIC_FAILOVER_GRID = "failoverOrder";
   public static final String ID_HOST_NETWORKING_PNIC_SPEED_COMBO = "speedCombo";
   public static final String ID_HOST_NETWORKING_PNIC_EDIT_SPEED_OK_BUTTON = "okBtn";
   public static final String ID_HOST_NETWORKING_FAILOVER_PNIC_REMOVE =
         "toolTip=Remove selected";
   public static final String ID_HOST_NETWORKING_IP_STATUS_COMBO =
         "confIpv6SupportOnHostDropDown";

   public static final String ID_PAGE_EDIT_NIC_SETTINGS =
         "wizardPageNavigator/step_nicSettingsView";
   public static final String ID_TREE_NETWORKING_VIEW =
         "vsphere.core.host.manage.networkingView/tocTree";

   public static final String ID_CDP_ENABLED_STATUS_LABEL = "cdpHeaderCollapsedLabel";
   public static final String ID_CDP_DISABLED_STATUS_LABEL = "_PnicDetailsView_Text11";
   public static final String ID_LLDP_ENABLED_STATUS_LABEL = "lldpHeaderCollapsedLabel";
   public static final String ID_LLDP_DISABLED_STATUS_LABEL = "_PnicDetailsView_Text13";

   public static final String ID_MANAGE_PHYSICAL_ADAPTERS_LIST = "mappingList";
   public static final String ID_MIGRATE_NETWORKING_LIST = "hostTreeList";
   public static final String ID_PORTGROUP_SETTINGS_NAME = "networkLabelInput";
   public static final String ID_HOST_NETWORKING_VNIC_LIST = "vnicList";
   public static final String ID_HOST_NETWORKING_VNIC_GRID =
         ID_HOST_NETWORKING_VNIC_LIST + "/list";
   public static final String ID_HOST_NETWORKING_VNIC_EDIT_BUTTON =
         "vsphere.core.host.network.editVnicSettingsAction/button";
   public static final String ID_HOST_NETWORKING_VNIC_REMOVE_BUTTON =
         "vsphere.core.host.network.removeVnicAction/button";
   public static final String ID_HOST_NETWORKING_PNIC_EDIT_BUTTON =
         "vsphere.core.host.network.editPnicSettingsAction/button";
   public static final String ID_HOST_NETWORKING_PNIC_REMOVE_BUTTON =
         "vsphere.core.host.network.editPnicSettingsAction/button";
   public static final String ID_HOST_NETWORKING_WARNING_TEXT = "errWarningText";

   public static final String ID_HOST_PHYSICAL_ADAPTERS_PROPERTIES_VIEW = "propView";

   // VC Summary
   public static final String ID_VC_SUMMARY_TAB = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_VC).append("summary.").append("container").toString();
   public static final String ID_VC_STATUS_PORTLET = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_VC).append("statusView").append(ID_PORTLET_SUFFIX)
         .toString();
   public static final String ID_VC_DETAILS_PORTLET = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_VC).append("summary.detailsView")
         .append(ID_PORTLET_SUFFIX).toString();
   public static final String ID_VC_LICENCE_PORTLET = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_VCENTER).append("summary.licensingView")
         .append(ID_PORTLET_SUFFIX).toString();
   public static final String ID_VC_STATUS_ACTIVE_TASKS = "tasks";
   public static final String ID_VC_DETAILS_VMS = "summary_vmCountProp_valueLbl";
   public static final String ID_VC_DETAILS_HOSTS = "summary_hostCountProp_valueLbl";

   // DC Summary
   public static final String ID_DC_SUMMARY =
         "com.vmware.vsphere.client.dcui.subView.summary";
   public static final String ID_DC_STATUS_PORTLET = EXTENSION_PREFIX
         + "<ENTITY_TYPE>summary.statusView.chrome";
   public static final String ID_DC_DETAILS_PORTLET = EXTENSION_PREFIX
         + "<ENTITY_TYPE>summary.detailsView.chrome";
   public static final String ID_DC_NETWORK_PORTLET =
         "com.vmware.vsphere.client.networkui.portlet.dcNetworks";
   public static final String ID_DC_NETWORK_CLOSE_BUTTON =
         "com.vmware.vsphere.client.networkui.portlet.dcNetworks/closeButton";
   public static final String ID_DC_NETWORK_MINMAX_BUTTON =
         "com.vmware.vsphere.client.networkui.portlet.dcNetworks/"
               + ID_PORTLET_MINMAX_BUTTON;
   public static final String ID_DC_STATUS_CLOSE_BUTTON =
         "com.vmware.vsphere.client.dcui.portlet.status/closeButton";
   public static final String ID_DC_STATUS_MINMAX_BUTTON =
         "com.vmware.vsphere.client.dcui.portlet.status/" + ID_PORTLET_MINMAX_BUTTON;
   public static final String ID_DC_DETAILS_CLOSE_BUTTON =
         "com.vmware.vsphere.client.dcui.portlet.details/closeButton";
   public static final String ID_DC_DETAILS_MINMAX_BUTTON =
         "com.vmware.vsphere.client.dcui.portlet.details/" + ID_PORTLET_MINMAX_BUTTON;
   public static final String ID_TAB_NAVIGATOR = "tabNavigator";
   public static final String ID_SUB_TABS_RELATED_ITEMS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append(MAIN_TAB_PREFIX).toString();
   public static final String ID_SUB_TABS_TOGGLE_BUTTON_BAR = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append(MAIN_TAB_PREFIX)
         .append("/toggleButtonBar").toString();
   public static final String ID_SUB_TABS_TAB_BAR = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append(MAIN_TAB_PREFIX).append("/tabBar").toString();
   public static final String ID_NAV_BOX = "appBody";
   public static final String ID_LABEL_MAIN_VIEW_TITLE = "appBody/titleText";
   public static final String ID_DC_STATUS_LABEL = "statusLabel";
   public static final String ID_DC_NETWORKS = "dcNetworks";
   public static final String ID_DC_VMS = "dcVms";
   public static final String ID_DC_HOSTS = "dcHosts";
   public static final String ID_DVS_NETWORK_DETAILS =
         "com.vmware.vsphere.client.networkui.portlet.dvsNetworks";
   public static final String ID_DC_NETWORK_GRID = EXTENSION_PREFIX
         + "<ENTITY_TYPE>networksView";
   public static final String ID_STATUS_ICON = "statusIcon";

   // Host Control
   public static final String ID_HOST_STATUS_GENERAL_VALUE = "entityStatus/statusLabel";
   public static final String ID_HOST_HARDWARE_MANUFACTURER_VALUE =
         "_HostHardwareView_StackBlock1/manufacturer";
   public static final String ID_HOST_HARDWARE_MODEL_VALUE =
         "_HostHardwareView_StackBlock2/model";
   public static final String ID_HOST_HARDWARE_CPU_VALUE = "cpu";
   public static final String ID_HOST_HARDWARE_CPU_CPUCORES_VALUE = "cpuCores";
   public static final String ID_HOST_HARDWARE_CPU_PROCESSORTYPE_VALUE = "cpuType";
   public static final String ID_HOST_HARDWARE_CPU_SOCKETS_VALUE = "cpuSockets";
   public static final String ID_HOST_HARDWARE_CPU_CORESPERSOCKET_VALUE =
         "cpuCoresPerSocket";
   public static final String ID_HOST_HARDWARE_CPU_LOGICALPROCESSORS_VALUE =
         "cpuLogical";
   public static final String ID_HOST_HARDWARE_CPU_HYPERTHREADING_VALUE =
         "cpuHyperthreading";
   public static final String ID_HOST_HARDWARE_MEMORY_VALUE = "memory";
   public static final String ID_HOST_HARDWARE_MEMORY_SYSTEM_VALUE = "memSystem";
   public static final String ID_HOST_HARDWARE_MEMORY_VIRTUALMACHINES_VALUE = "memVMs";
   public static final String ID_HOST_HARDWARE_MEMORY_SERVICECONSOLE_VALUE =
         "memService";
   public static final String ID_HOST_HARDWARE_MEMORY_CAPACITY_VALUE = "memCapacity";
   public static final String ID_HOST_HARDWARE_NETWORKING_VALUE = "networking";
   public static final String ID_HOST_HARDWARE_NETWORKING_ROUTING_VALUE = "routing";
   public static final String ID_HOST_HARDWARE_NETWORKING_NETWORKS_VALUE = "numNetworks";
   public static final String ID_HOST_HARDWARE_NETWORKING_PHYSICALADAPTERS_VALUE =
         "numNics";
   public static final String ID_HOST_HARDWARE_STORAGE_VALUE = "storage";
   public static final String ID_HOST_HARDWARE_STORAGE_DATASTORES_VALUE = "datastores";
   public static final String ID_HOST_HARDWARE_STORAGE_PHYSICALADAPTERS_VALUE =
         "numHbas";
   public static final String ID_HOST_CONFIGURATION_ESXVERSION_VALUE = "version";
   public static final String ID_HOST_CONFIGURATION_VMOTIONENABLED_VALUE = "vmotion";
   public static final String ID_HOST_STORAGE_DATASTORES_PORTLET = "title=Datastores";
   public static final String ID_HOST_TAB_NAVIGATOR = ID_NAV_BOX + "/tabNavigator";
   public static final String ID_HOST_SUMMARY_TAB =
         "vsphere.core.host.summary.container";
   public static final String ID_HOST_STORAGE_TAB =
         "vsphere.core.host.storageView.container";
   public static final String ID_HOST_ANNOTATION_OK_BUTTON =
         "editAnnotationView/okButton";
   public static final String ID_HOST_ANNOTATION_CANCEL_BUTTON =
         "editAnnotationView/cancelButton";
   public static final String ID_HOST_CONFIG_FTENABLED = "ftEnabled";
   public static final String ID_HOST_SUMMARY_FT_MISCONFIGURED_ARROW_IMAGE =
         "arrowImageftStack";
   public static final String ID_HOST_SUMMARY_HA_STATE_LABEL = "haEnabled";
   public static final String ID_HOST_NETWORKS_TAB =
         "vsphere.core.host.networksView.container";

   // Host Operations
   public static final String ID_HOST_MAINTENANCE_CONFIRM_DIALOG =
         "confirmMaintenanceForm";
   public static final String ID_HOST_MAINTENANCE_CONFIRM_DIALOG_YESBUTTON = "yesButton";
   public static final String ID_HOST_MAINTENANCE_CONFIRM_DIALOG_NOBUTTON = "noButton";
   public static final String ID_HOST_ERROR_DESCRIPTION_TEXT = "descriptionText";
   public static final String ID_MM_EVACUATE_VM_CHECKBOX = "evacuateVMsCheckBox";
   public static final String ID_PANEL_CONFIRM_ENTER_STANDBY_MODE = "confirmStandbyForm";
   public static final String ID_BUTTON_ENTER_STANDBY_MODE_YES = "yesButton";
   public static final String ID_BUTTON_ENTER_STANDBY_MODE_NO = "noButton";

   // VM Summary View
   public static final String ID_GUEST_OS_PORTLET_AUTOMATION = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "guestOSView.chrome";
   public static final String ID_GUEST_OS_DETAILS_GOS_NAME =
         "summary_guestFullName_valueLbl";
   public static final String ID_GUEST_OS_DETAILS_GOS_IP =
         "summary_ipAddressProp_valueLbl";
   public static final String ID_GUEST_OS_DETAILS_GOS_DNS =
         "summary_dnsNameProp_valueLbl";
   public static final String ID_GUEST_OS_DETAILS_GOS_TOOLS_RUNNING =
         "summary_vmwareTools_valueLbl";
   public static final String ID_GUEST_OS_DETAILS_PLUGIN_LINK = "vmConsoleLink";
   public static final String ID_GUEST_OS_DETAILS_PLUGIN_IMAGE = "shineImage";
   public static final String ID_VM_SUMMAY_HOST_NAME_LABEL =
         "summary_hostName_valueLink";
   public static final String ID_VM_COMPATIBILITY_LABEL = "summary_vmVersion_valueLbl";
   public static final String ID_VM_SUMMARY_ALARM_DISPLAY_OBJECT_CONTAINER =
         "issueContainer";

   // DRS Groups
   public static final String ID_BUTTON_ADD_DRS_GROUP = "addDrsGroup";
   public static final String ID_TEXTINPUT_TEXT_DISPLAY = "groupName";
   public static final String ID_DROPDOWN_DRS_GROUP_TYPE = "groupType";
   public static final String ID_BUTTON_ADD_DRS_GROUP_MEMBER = "addGroupMember";
   public static final String ID_BUTTON_ADD_DRS_GROUP_SEARCH = "searchButton";
   public static final String ID_ADVANCED_DATAGRID_DRS_GROUP = "groupsList";
   public static final String ID_CHECKBOX_DRS_GROUP =
         "_CheckBoxColumnRenderer_CheckBox1";
   public static final String ID_ADVANCED_DATAGRID_DRS_GROUP_VMFILTER =
         "VirtualMachine_filterGridView";
   public static final String ID_ADVANCED_DATAGRID_DRS_GROUP_HOSTFILTER =
         "HostSystem_filterGridView";
   public static final String ID_BUTTON_REMOVE_DRS_GROUPS = "removeDrsGroup";
   public static final String ID_BUTTON_REMOVE_DRS_GROUP_MEMBER = "removeGroupMember";
   public static final String ID_BUTTON_EDIT_DRS_GROUP = "editDrsGroup";
   public static final String ID_ADVANCED_DATAGRID_DRS_GROUPMEM = "membersList";
   public static final String ID_REMOVE_BUTTON_DRS_GROUPS = "removeDrsGroup";
   public static final String ID_DROPDOWN_DRS_RULES_VM_GROUP = "vmGroup";
   public static final String ID_DROPDOWN_DRS_RULES_RELATION_TYPE = "relationType";
   public static final String ID_DROPDOWN_DRS_RULES_HOST_GROUP = "hostGroup";
   public static final String ID_ADVANCED_DATAGRID_DRS_GROUP_SELECTEDOBJECTS =
         "VirtualMachine_selectedObjectGridView";

   public static final String ID_CLUSTER_DRS_VM_GROUP = "0";
   public static final String ID_CLUSTER_DRS_HOST_GROUP = "1";
   public static final String ID_CLUSTER_GROUP_NAME_INDEX = "1";
   public static final String ID_ADVANCED_DATAGRID_FAULT_LIST = "faultList";
   public static final String ID_CHECKBOX_EVACUATE_VMS = "evacuateVMsCheckBox";

   // Cluster IPMI
   public static final String ID_CLUSTER_SETTINGS_EDIT_BUTTON = "editDrsButton";
   public static final String ID_CLUSTER_EDIT_SETTINGS_PANEL = "tiwoDialog";
   public static final String ID_CLUSTER_SERVICES_EDIT_PANEL_POWER_MANAGEMENT_MANUAL_SPARK_RADIO_BUTTON =
         "_DrsConfigForm_RadioButton5";
   public static final String ID_CLUSTER_SERVICES_EDIT_PANEL_POWER_MANAGEMENT_AUTOMATIC_SPARK_RADIO_BUTTON =
         "_DrsConfigForm_RadioButton6";
   public static final String ID_CLUSTER_SERVICES_EDIT_PANEL_POWER_MANAGEMENT_OFF_SPARK_RADIO_BUTTON =
         "_DrsConfigForm_RadioButton4";
   public static final String ID_CLUSTER_SERVICES_EDIT_PANEL_EXPAND_POWER_MANAGEMENT_AUTOMATION_IMAGE =
         "arrowImagedpmBlock";
   public static final String ID_CLUSTER_EDIT_SETTINGS_PANEL_OK_BUTTON =
         new StringBuffer(ID_CLUSTER_EDIT_SETTINGS_PANEL).append("/okButton").toString();
   public static final String ID_BUTTON_CANCEL_CLUSTER_EDIT_SETTINGS_PANEL =
         new StringBuffer(ID_CLUSTER_EDIT_SETTINGS_PANEL).append("/cancelButton")
               .toString();
   public static final String ID_CLUSTER_MANAGE_SETTING_POWER_MANAGEMENT_LABEL_VALUE =
         "_DrsConfigView_Label2";
   public static final String ID_CLUSTER_SUMMARY_TAB_POWER_MANAGEMENT_LABEL_VALUE =
         "_DrsDetailsView_Label6";
   public static final String ID_CLUSTER_TURNON_DPM_RECOMMENDATION_DIALOG_OK =
         "automationName=OK";
   public static final String ID_CLUSTER_TURNON_DPM_RECOMMENDATION_DIALOG_CANCEL =
         "automationName=Cancel";
   public static final String ID_CLUSTER_POWER_MANAGEMENT_RECOMMENDATION_DIALOG =
         "className=AlertForm";
   public static final String ID_SLIDER_DRS_MIGRATION_THRESHOLD =
         "_DrsConfigForm_HSlider1";
   public static final String ID_SLIDER_DPM__THRESHOLD = "_DrsConfigForm_HSlider2";
   public static final String ID_IMAGE_DRS_AUTOMATION_ARROW =
         "arrowImagedrsAutomationBlock";
   public static final String ID_IMAGE_POWER_MANAGEMENT_ARROW = "arrowImagedpmBlock";
   public static final String ID_LABEL_DRS_AUTOMATION_LEVEL = "_DrsConfigView_Text1";
   public static final String ID_LABEL_DRS_MIGRATION_THRESHOLD = "_DrsConfigView_Text2";
   public static final String ID_LABEL_VIRTUAL_MACHINE_AUTOMATION =
         "_DrsConfigView_Text3";
   public static final String ID_LABEL_DPM_AUTOMATION_LEVEL = "_DrsConfigView_Text4";
   public static final String ID_LABEL_DPM_THRESHOLD = "_DrsConfigView_Text5";
   public static final String ID_LISTBOX_DPM_AUTOMATION_LEVEL = "dpmAutLevelDropDown";
   public static final String ID_VGROUP_DPM_THRESHOLD = "_DrsConfigForm_VGroup6";
   public static final String ID_STACKEDITOR_DRS_CONFIG_BLOCK =
         "_DrsConfigForm_StackEditor1";

   // Cluster
   public static final String ID_CLUSTER_TITLE = "vsphere.core.cluster.view/title";
   public static final String ID_CONFIRM_YES_NO_DIALOG = "YesNoDialog";
   public static final String ID_REMOVE_CONFIRM_YES_BUTTON = ID_CONFIRM_YES_NO_DIALOG
         + "/automationName=Yes";
   public static final String ID_REMOVE_CONFIRM_NO_BUTTON = ID_CONFIRM_YES_NO_DIALOG
         + "/automationName=No";
   public static final String ID_LABEL_SUMARRY_CLUSTER_COUNT =
         "summary_clusterCountProp_valueLbl";
   // Cluster Services list
   public static final String ID_CLUSTER_SETTINGS_SERVICES_SPARK_LIST =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER)
               .append("manage.").append("settingsView").append("/tocTree").toString();

   // Cluster HA - High Availability
   public static final String ID_CLUSTER_SERVICES_EDIT_PANEL_EXPAND_HEARTBEAT_DATASTORE_IMAGE =
         new StringBuffer(ID_CLUSTER_EDIT_SETTINGS_PANEL).append(
               "/arrowImagedsHeartbeatBlock").toString();
   public static final String ID_CLUSTER_SERVICES_EDIT_PANEL_PREFERED_DATASTORE_LABEL =
         "_DatastoresForHeartbeatingList_Label1";
   public static final String ID_CLUSTER_SERVICES_EDIT_PANEL_HOST_MOUNTING_DATASTORE_LABEL =
         "_DatastoresForHeartbeatingList_Label2";

   // Cluster Summary - Portlet IDs
   public static final String ID_CLUSTER_SUMMARY_RESOURCES_PORTLET =
         "vsphere.core.cluster.summary.resourcesView.chrome";
   public static final String ID_CLUSTER_SUMMARY_RESOURCES_COMPUTE_EXPAND =
         "arrowImage_ClusterResourceView_StackBlock1";
   public static final String ID_CLUSTER_SUMMARY_RESOURCES_HOSTS =
         "automationName=Hosts";
   public static final String ID_CLUSTER_SUMMARY_RESOURCES_PROCESSORS =
         "automationName=Total Processors";
   public static final String ID_CLUSTER_SUMMARY_RESOURCES_CPU =
         "automationName=Total CPU Resources";
   public static final String ID_CLUSTER_SUMMARY_RESOURCES_MEMORY =
         "automationName=Total Memory";
   public static final String ID_CLUSTER_SUMMARY_RESOURCES_EVC =
         "automationName=EVC Mode";

   public static final String ID_CLUSTER_SUMMARY_HA_PROTECTED = "_HaDetailsView_Label1";
   public static final String ID_CLUSTER_SUMMARY_HA_CONFIGURED_FAILOVER =
         "_HaDetailsView_Label9";
   public static final String ID_CLUSTER_SUMMARY_HA_CURRENT_FAILOVER =
         "_HaDetailsView_Label11";
   public static final String ID_CLUSTER_SUMMARY_HA_HOST_MONITORING =
         "_HaDetailsView_Label15";
   public static final String ID_CLUSTER_SUMMARY_HA_VM_MONITORING =
         "_HaDetailsView_Label17";

   public static final String ID_CLUSTER_SUMMARY_CONSUMERS_PORTLET =
         "vsphere.core.cluster.summary.consumersView.chrome";
   public static final String ID_CLUSTER_SUMMARY_CONSUMERS_VM_EXPAND =
         "arrowImagevmBlock";
   public static final String ID_CLUSTER_SUMMARY_CONSUMERS_RESOURCE_POOLS =
         "automationName=Resource Pools";
   public static final String ID_CLUSTER_SUMMARY_CONSUMERS_VAPPS =
         "automationName=vApps";
   public static final String ID_CLUSTER_SUMMARY_CONSUMERS_VMS =
         "automationName=Virtual Machines";
   public static final String ID_CLUSTER_SUMMARY_CONSUMERS_POWERED_ON =
         "automationName=Powered-on";
   public static final String ID_CLUSTER_SUMMARY_CONSUMERS_POWERED_OFF =
         "automationName=Powered-off";
   public static final String ID_CLUSTER_SUMMARY_CONSUMERS_TOTAL =
         "automationName=Total";

   public static final String ID_CLUSTER_SUMMARY_TOTAL_PROCESSORS =
         "summary_numProcessors_label";
   public static final String ID_CLUSTER_SUMMARY_TOTAL_MIGRATIONS =
         "summary_numVMotionMigrations_label";
   public static final String ID_CLUSTER_SUMMARY_DRS_BADGE = "drsBadge";

   // Cluster DRS Automation
   public static final String ID_CLUSTER_SETTINGS_DRS_AUTOMATION_VALUE_LABEL =
         "_DrsConfigView_Label1";
   public static final String ID_CLUSTER_SETTINGS_EXPAND_DRS_AUTOMATION_IMAGE =
         "arrowImagedrsAutomationBlock";
   public static final String ID_CLUSTER_SETTINGS_EXPAND_DRS_VM_AUTOMATION_VALUE_LABEL =
         "_DrsConfigView_Text3";
   public static final String ID_CLUSTER_SUMMARY_MIGRATION_AUTOMATION_LABEL =
         "_DrsDetailsView_Label3";
   public static final String ID_CLUSTER_SERVICES_EDIT_PANEL_EXPAND_DRS_AUTOMATION_IMAGE =
         new StringBuffer(ID_CLUSTER_EDIT_SETTINGS_PANEL).append(
               "/arrowImagedrsAutomationBlock").toString();
   public static final String ID_CLUSTER_SERVICES_EDIT_PANEL_DRS_FULL_SPARK_RADIO_BUTTON =
         new StringBuffer(ID_CLUSTER_EDIT_SETTINGS_PANEL).append(
               "/_DrsConfigForm_RadioButton3").toString();
   public static final String ID_CLUSTER_SERVICES_EDIT_PANEL_DRS_PARTIAL_SPARK_RADIO_BUTTON =
         new StringBuffer(ID_CLUSTER_EDIT_SETTINGS_PANEL).append(
               "/_DrsConfigForm_RadioButton2").toString();
   public static final String ID_CLUSTER_SERVICES_EDIT_PANEL_DRS_MANUAL_SPARK_RADIO_BUTTON =
         new StringBuffer(ID_CLUSTER_EDIT_SETTINGS_PANEL).append(
               "/_DrsConfigForm_RadioButton1").toString();
   public static final String ID_CLUSTER_SERVICES_EDIT_PANEL_VM_AUTOMATION =
         new StringBuffer(ID_CLUSTER_EDIT_SETTINGS_PANEL).append(
               "/_DrsConfigForm_CheckBox2").toString();
   public static final String ID_DROPDOWN_AUTOMATION_LEVEL = "drsBehavior";
   public static final String ID_DROPDOWN_MIGRATION_THRESHOLD_LEVEL =
         "drsMigrationThreshold";
   public static final String ID_CLUSTER_SETTINGS_GROUPS_NAME_LABEL =
         "tiwoDialog/groupName";
   public static final String ID_SLIDER_DRS_MIGRATION_THRESHOLD_VALUE =
         "drsMigrationThreshold";
   public static final String ID_DROP_DOWN_LIST_DRS_BEHAVIOR = "drsBehavior";
   public static final String ID_BUTTON_CLUSTER_BROWSE = "browseBtn";
   public static final String ID_BUTTON_CLUSTER_DC_CLUSTER =
         "clustersForDatacenter.button";
   public static final String ID_BUTTON_CLUSTER_DC_INV_FOLDERS =
         "invFoldersForDatacenter.button";
   public static final String ID_BUTTON_CLUSTER_HOST_AND_CLUSTER =
         "clustersForHAndC.button";
   public static final String ID_BUTTON_CLUSTER_CREATE_GLOBAL =
         "vsphere.core.cluster.createActionGlobal/button";
   public static final String ID_BUTTON_CLUSTER_BROWSE_SEARCH = "searchButton";
   public static final String ID_LABEL_CLUSTER_BROWSE_INVALID = "messageLabel";
   public static final String ID_LISTBOX_DRS_AUTOMATION_LEVEL = "drsAutLevelDropDown";
   public static final String ID_BUTTON_SCHEDULE_DRS = "scheduleDrsButton";
   public static final String ID_IMAGE_RELEVANT_CLUSTER_SETTINGS =
         "arrowImage_VmOverridesConfigForm_StackBlock1";
   public static final String ID_IMAGE_VSPHERE_DRS =
         "arrowImage_VmOverridesConfigForm_StackBlock2";
   public static final String ID_TEXT_VALIDATION_MESSAGE = "_message";

   // Cluster General Edit settings
   public static String ID_CLUSTER_EDIT_GENERAL_SETTINGS_BUTTON =
         "btn_vsphere.core.cluster.editGeneralSettingsAction";

   // Cluster Swap File options
   public static String ID_CLUSTER_SWAP_FILE_VM = "_GeneralConfigForm_RadioButton1";
   public static String ID_CLUSTER_SWAP_FILE_HOST = "_GeneralConfigForm_RadioButton2";
   public static String ID_CLUSTER_SETTINGS_CONFIGURATION_GENERAL_SWAP_LOC =
         "_GeneralConfigView_Label1";
   public static String ID_LABEL_DEFAULT_SWAPFILE_LOCATION = "location";

   // Host Swap File options
   public static String ID_BUTTON_HOST_VM_SWAPFILE_LOC_EDIT_BUTTON =
         "btn_vsphere.core.host.editVmSwapfileLocationAction";
   public static String ID_RADIO_HOST_VM_DIR_OPTION = "vmdirOption";
   public static String ID_RADIO_HOST_DATASTORE_OPTION = "datastoreOption";
   public static final String ID_BUTTON_ADD_DRS_GROUP_MEMBER_SECTION =
         "addDrsGroupMember";
   public static final String ID_BUTTON_REMOVE_DRS_GROUP_MEMBER_SECTION =
         "removeDrsGroupMember";

   // Cluster Summary
   public static final String ID_CLUSTER_STATUS_PORTLET = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER)
         .append("summary.statusView.chrome").toString();
   public static final String ID_CLUSTER_RESOURCES_PORTLET = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER)
         .append("summary.resourcesView.chrome").toString();
   public static final String ID_CLUSTER_CONSUMERS_PORTLET = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER)
         .append("summary.consumersView.chrome").toString();
   public static final String ID_PORTLET_CLUSTER_VSPHERE_DRS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER)
         .append("summary.drsDetailsView.chrome").toString();
   public static final String ID_PORTLET_CLUSTER_CUSTOM_ATTRIBUTES = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER)
         .append("summary.customAttributesView.chrome").toString();
   public static final String ID_CLUSTER_SERVICES_PORTLET = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER)
         .append("summary.servicesView.chrome").toString();
   public static final String ID_HOST_NETWORK_NETWORKS_PORTLET =
         "com.vmware.vsphere.client.networkui.portlet.hostNetworks";
   public static final String ID_CLUSTER_HOSTS_GRID =
         "com.vmware.vsphere.client.hostui.subView.Hosts/hostList";
   public static final String ID_CLUSTER_GUEST_MEMORY_PORTLET =
         "vsphere.core.cluster.guestMemoryView.chrome";
   public static final String ID_CLUSTER_HOST_MEMORY_PORTLET =
         "vsphere.core.cluster.hostMemoryView.chrome";
   public static final String ID_CLUSTER_HOST_CPU_PORTLET =
         "vsphere.core.cluster.hostCPUView.chrome";
   public static final String ID_DRS_CLUSTER_LOAD_BALANCE = "_DrsDetailsView_Label1";
   public static final String ID_PORTLET_CLUSTER_VSPHERE_HA =
         "vsphere.core.cluster.summary.haPortletView.chrome";

   // Utilization tab labels
   public static final String ID_GUEST_MEMORY_PORTLET = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("monitor.")
         .append("utilization.").append("guestMemoryView.chrome").toString();
   public static final String ID_PRIVATE_MEMORY = ID_GUEST_MEMORY_PORTLET
         + "/privateMemoryLbl";
   public static final String ID_ACTIVE_GUEST_MEMORY = ID_GUEST_MEMORY_PORTLET
         + "/activeMemoryLbl";
   public static final String ID_BALLOONED_MEMORY = ID_GUEST_MEMORY_PORTLET
         + "/balloonedMemoryLbl";
   public static final String ID_SHARED_MEMORY = ID_GUEST_MEMORY_PORTLET
         + "/sharedMemoryLbl";
   public static final String ID_UNACCESSED_MEMORY = ID_GUEST_MEMORY_PORTLET
         + "/unaccessedMemoryLbl";
   public static final String ID_SWAPPED_MEMORY = ID_GUEST_MEMORY_PORTLET
         + "/swappedMemoryLbl";
   public static final String ID_COMPRESSED_MEMORY = ID_GUEST_MEMORY_PORTLET
         + "/compressedMemoryLbl";

   public static final String ID_HOST_MEMORY_PORTLET =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("monitor.")
               .append("utilization.").append("hostMemoryView.chrome").toString();
   public static final String ID_HOST_CPU_PORTLET = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("monitor.").append("utilization.")
         .append("hostCPUView.chrome").toString();
   public static final String ID_HOST_CPU_ACTIVE = ID_HOST_CPU_PORTLET
         + "/activeOverheadResourceLbl";
   public static final String ID_HOST_CPU_CONSUMED = ID_HOST_CPU_PORTLET
         + "/consumedResourceLbl";
   public static final String ID_HOST_MEMORY_CONSUMED = ID_HOST_MEMORY_PORTLET
         + "/consumedResourceLbl";
   public static final String ID_HOST_MEMORY_OVERHEAD = ID_HOST_MEMORY_PORTLET
         + "/activeOverheadResourceLbl";

   // Cluster Summary - Cluster Resources portlet
   public static final String ID_LABEL_CLUSTER_COMPUTE =
         "titleLabel._ClusterResourceView_StackBlock1";
   public static final String ID_CLUSTER_HOSTS = "numHostsMain";
   public static final String ID_CLUSTER_NUM_HOSTS = "numHosts";
   public static final String ID_CLUSTER_TOTAL_PROCESSORS = "numCores";
   public static final String ID_CLUSTER_TOTAL_CPU_RESOURCES = "totalCpu";
   public static final String ID_CLUSTER_TOTAL_MEMORY = "totalMem";
   public static final String ID_CLUSTER_EVC = "evcStatus";
   public static final String ID_CLUSTER_VIRTUAL_MACHINES = "titleLabel.vmBlock";
   public static final String ID_CLUSTER_POWERED_ON_VMS = "numPoweredOnVms";
   public static final String ID_CLUSTER_POWERED_OFF_VMS = "numPoweredOffVms";
   public static final String ID_LABEL_TOTAL_NUM_VMS = "numTotalVms";

   // Cluster Summary - Cluster Consumers portlet
   public static final String ID_CLUSTER_RESOURCEPOOLS = "numResourcePools";
   public static final String ID_CLUSTERS_VMS = "numTotalVms";

   // Cluster Summary - Status portlet
   public static final String ID_CLUSTER_STATUS = "entityStatus/statusLabel";
   public static final String ID_CLUSTER_STATUS_ICON = "statusIcon";

   // Cluster Summary - Services portlet
   public static final String ID_CLUSTER_VMWARE_HA =
         "vsphere.core.cluster.servicesView.chrome/titleLabel.haBlock";
   public static final String ID_CLUSTER_FAILOVER_THRESHOLD = "failoverThreshold";
   public static final String ID_CLUSTER_HOST_MONITORING = "hostMonitoring";
   public static final String ID_CLUSTER_HOST_MONITORING_BLOCK =
         "titleLabel.haHostMonitoringBlock";
   public static final String ID_CLUSTER_HOST_MONITORING_RESTART_PRIORITY_DROPDOWN =
         "restartPriorityStateDropDown";
   public static final String ID_CLUSTER_HOST_MONITORING_ISOLATION_RESPONSE_DROPDOWN =
         "isolationResponseStateDropDown";
   public static final String ID_CLUSTER_ADMISSION_CONTROL =
         "admissionControlStateChkBox";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_CAPACITY =
         "failoverLevelInput";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_CAPACITY_INCREMENT =
         "failoverLevelInput/incrementButton";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_CAPACITY_DECREMENT =
         "failoverLevelInput/decrementButton";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_CAPACITY_RADIO_BUTTON =
         "failoverLevelPolicyRadioButton";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_PERCENTAGE_RADIO_BUTTON =
         "reservedCapacityPolicyRadioButton";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_PERCENTAGE_CPU_STEPPER =
         "failoverReservedCpu";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_PERCENTAGE_MEM_STEPPER =
         "failoverReservedMem";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_HOSTS_RADIO_BUTTON =
         "failoverHostPolicyRadioButton";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_HOSTS_ADD_BUTTON =
         "addFailoverHostButton";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_HOSTS_DELETE_BUTTON =
         "removeFailoverHostsButton";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_HOSTS_GRID =
         "failoverHostsAdvancedDataGridEx";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_BLOCK =
         "titleLabel.haHostAdmissionControlBlock";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_DISABLE_RADIO_BUTTON =
         "admissionControlDisabledRadioButton";
   public static final String ID_CLUSTER_VM_MONITORING = "vmMonitoringStateDropDown";
   public static final String ID_CLUSTER_VM_MONITORING_BLOCK =
         "titleLabel.haVMMonitoringBlock";
   public static final String ID_CLUSTER_VM_MONITORING_SENSITIVITY_SLIDER =
         "sensitivityPresetSlider";
   public static final String ID_CLUSTER_ADVANCED_SETTINGS_BLOCK =
         "titleLabel.haAdvOptionsBlock";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_PERCENTAGE_TEXT =
         "admissionControlGridRow/_HaConfigView_Text5";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_PERCENTAGE_CPU_TEXT =
         "admissionControlDescriptionGridRow/_HaConfigView_Text6";
   public static final String ID_CLUSTER_ADMISSION_CONTROL_PERCENTAGE_MEM_TEXT =
         "admissionControlAdditionalDescriptionGridRow/_HaConfigView_Text7";

   // Cluster Summary - vSphere DRS portlet
   public static final String ID_LABEL_CLUSTER_DRS_SETTINGS = "step_drsConfigForm";
   public static final String ID_LABEL_CLUSTER_DRS_BALANCED = "_DrsDetailsView_Label1";
   public static final String ID_LABEL_CLUSTER_MIGRATION_AUTOMATION =
         "_DrsDetailsView_Label3";
   public static final String ID_LABEL_CLUSTER_MIGRATION_THRESHOLD =
         "_DrsDetailsView_Text1";
   public static final String ID_LABEL_CLUSTER_POWER_MANAGEMENT =
         "_DrsDetailsView_Label6";
   public static final String ID_LABEL_CLUSTER_DRS_RECOMMENDATIONS =
         "_DrsDetailsView_Label8";
   public static final String ID_LABEL_CLUSTER_DRS_FAULTS = "_DrsDetailsView_Label10";
   public static final String ID_LABEL_CLUSTER_DRS_RECOMMENDATIONS_COUNT =
         "_DrsDetailsView_Label8";

   // Cluster Summary - Custom Attributes portlet
   public static final String ID_LABEL_CLUSTER_CUSTOM_ATTRIBUTES =
         "vsphere.core.cluster.summary.customAttributesView.chrome";

   // Cluster Monitor sub tab
   public static final String ID_CLUSTER_MONITOR_HA_SPARK_LIST = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER).append("monitor.")
         .append("haView").append("/tocTree").toString();
   public static final String ID_BUTTON_CLUSTER_MONITOR_DRS_VIEW = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER)
         .append(ID_MONITOR_MAIN_TAB_CONTAINER).append(".drsView.").append("button")
         .toString();
   public static final String ID_CLUSTER_MONITOR_HA_HOSTS_PORTLET =
         "vsphere.core.cluster.monitor.ha.summaryView.hostsRuntimeInfo.chrome";
   public static final String ID_CLUSTER_MONITOR_HA_HOSTS_MASTER_TEXT =
         "_HostsRuntimeInfoView_Label1";
   public static final String ID_CLUSTER_MONITOR_HA_HOSTS_CONNECTED_TO_MASTER_TEXT =
         "_HostsRuntimeInfoView_Label2";
   public static final String ID_CLUSTER_MONITOR_HA_SUMMARY_TAB = "1";
   public static final String ID_CLUSTER_MONITOR_HA_HEARTBEAT_TAB = "2";
   public static final String ID_CLUSTER_MONITOR_HA_CONFIG_ISSUES_TAB = "3";
   public static final String ID_CLUSTER_MONITOR_HA_BADGE = "haBadge";
   public static final String ID_CLUSTER_MONITOR_FAULT_LIST = "faultList";
   public static final String ID_CLUSTER_MONITOR_FAULT_DETAILS_LIST = "faultDetailList";
   public static final String ID_CLUSTER_MONITOR_HA_VMS_PORTLET =
         "vsphere.core.cluster.monitor.ha.summaryView.vmsRuntimeInfo.chrome";
   public static final String ID_BUTTON_CLUSTER_MONITOR_HA_VIEW =
         "vsphere.core.cluster.monitor.haView.button";
   public static final String ID_ADV_DATAGRID_CLUSTER_MONITOR_HA_CONFIG_ISSUES =
         "entityDataGrid";

   // Cluster Virtual Machine sub tab
   public static final String ID_VMS_LIST = "list";
   public static final String ID_STANDALONE_HOST_RELATED_ITEMS_VM_LIST =
         "vmsForStandaloneHost/list";
   public static final String ID_STANDALONE_HOST_RELATED_ITEMS_TEMPLATES_LIST =
         "vmTemplatesForStandAloneHost/list";
   public static final String ID_CLUSTERED_HOST_RELATED_ITEMS_VM_LIST =
         "vmsForClusteredHost/list";
   public static final String ID_CLUSTER_RELATED_ITEMS_VM_LIST =
         "vsphere.core.cluster.related/vmsForCluster/list";

   // Cluster Manage sub tab
   public static final String ID_CLUSTER_SETTINGS_CONFIGURATION_SPARK_LIST =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER)
               .append("manage.").append("settingsView").append("/tocTree").toString();
   public static final String ID_CLUSTER_SETTINGS_VM_OVERRIDES_ADVANCED_DATAGRID =
         "vmOverridesGrid";
   public static final String ID_CLUSTER_SETTINGS_CONFIGURATION_GENERAL_NAME =
         "_GeneralConfigView_Label1";
   // TODO: the id and automationName of both the "services" label and HA text
   // label are same, hence using the text. Will modify once it gets updated.
   public static final String ID_CLUSTER_SETTINGS_HA_STATUS_TEXT =
         "text=vSphere HA is Turned ON";
   public static final String ID_CLUSTER_SETTINGS_DRS_STATUS_TEXT =
         "_DrsConfigView_Label1";
   public static final String ID_CLUSTER_SETTINGS_HA_ISOLATION_ADDR =
         "das.isolationaddress1";
   public static final String ID_CLUSTER_SETTINGS_HA_USE_ISOLATION =
         "das.usedefaultisolationaddress";

   public static final String ID_CLUSTER_SETTINGS_EVC_EDIT_BUTTON =
         "btn_vsphere.core.cluster.evc.configureEvcAction";
   public static final String ID_CLUSTER_SETTINGS_EVC_EDIT_INTEL_RADIO =
         "intelVendorButton";
   public static final String ID_CLUSTER_SETTINGS_EVC_EDIT_AMD_RADIO = "amdVendorButton";
   public static final String ID_CLUSTER_SETTINGS_EVC_EDIT_DISABLE_EVC_RADIO =
         "disableEvcButton";
   public static final String ID_CLUSTER_SETTINGS_EVC_COMPATIBILITY_CHECK_LABEL =
         "compatibilityPanel/itemMessage";

   // Cluster Monitor sub tab
   public static final String ID_CLUSTER_SETTINGS_DRS_SPARK_LIST = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER).append("monitor.")
         .append("drsView").append("/tocTree").toString();

   // Create Cluster
   public static final String ID_TEXT_CREATE_CLUSTER_CLUSTER_NAME = "nameTextInput";
   public static final String ID_CHECKBOX_CREATE_CLUSTER_DRS = "drsTurnOn";
   public static final String ID_CREATE_CLUSTER_HA_CHECKBOX = "haTurnOn";
   public static final String ID_DIALOG_OK_BUTTON = "okButton";
   // HA cluster configuration
   public static final String ID_CLUSTER_EDIT_HA_BUTTON = "editHaSettingsBtn";
   public static final String ID_ENABLE_HA_CHECKBOX = "enableHACheckBox";
   public static final String ID_CLUSTER_HA_HOST_MONITORING_LABEL =
         "_HaConfigView_Label2";
   // HA cluster block
   public static final String ID_CREATE_CLUSTER_HA_BLOCK =
         "tiwoDialog/titleLabel.haStackBlock";
   public static final String ID_CREATE_CLUSTER_HA_BLOCK_HOST_MONITORING =
         "tiwoDialog/enableHostMonitoring";
   public static final String ID_CREATE_CLUSTER_HA_BLOCK_ADMISSION_CONTROL =
         "tiwoDialog/admissionControlExpandedCheckBox";
   public static final String ID_CREATE_CLUSTER_HA_BLOCK_HOST_FAILURES_RADIO =
         "tiwoDialog/failoverLevel";
   public static final String ID_CREATE_CLUSTER_HA_BLOCK_RESERVED_CAPACITY_RADIO =
         "tiwoDialog/reservedCapacity";
   public static final String ID_CREATE_CLUSTER_HA_BLOCK_VM_MONITORING =
         "tiwoDialog/vmMonitoringExpanded";

   // Virtual Machines Tab
   public static final String ID_VMS_ADVANCE_DATAGRID = "vmList";
   public static final String ID_VAPP_VMS_ADVANCE_DATAGRID = "vmsForVApp/list";
   public static final String ID_VMS_ADVANCE_DATAGRID_PROGRESS_BAR =
         ID_VMS_ADVANCE_DATAGRID + "/loadingProgressBar";
   public static final String ID_ADVANCED_DATAGRID_INLINE_EDIT_UI_COMPONENT =
         "uicEditor";

   public static final String MAX_HORIZONTAL_SCROLL_POSITION =
         "maxHorizontalScrollPosition";
   public static final String MAX_VERTICAL_SCROLL_POSITION = "maxVerticalScrollPosition";

   public static final String SPARK_MAX_HORIZONTAL_SCROLL_POSITION = "contentWidth";
   public static final String SPARK_MAX_VERTICAL_SCROLL_POSITION = "contentHeight";

   public static final String ID_DATA_PROVIDER_POWERSTATE = "runtime.powerState";
   public static final String ID_DATA_PROVIDER_STATUS = "summary.overallStatus";
   public static final String ID_DATA_PROVIDER_ALARM_ACTION = "alarmActionsEnabled";
   public static final String ID_DATA_PROVIDER_VM_VERSION = "config.version";
   public static final String ID_BUTTON_SHUTDOWN_SINGLE_VM = new StringBuffer(
         "vsphere.core.vm.shutdownActionSingleVm").append(ID_BUTTON).toString();
   public static final String ID_BUTTON_REBOOT_VM = new StringBuffer(
         "vsphere.core.vm.rebootAction").append(ID_BUTTON).toString();

   // Getting Started Tab
   public static final String ID_GETTING_STARTED_TAB_NAME = "Getting Started";
   // Hosts Tab
   public static final String ID_HOSTS_TAB = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("related.hostsView").append(".container")
         .toString();
   public static final String ID_HOSTS_TAB_STORAGE = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("related.hostsView").append(".container")
         .toString();
   public static final String ID_HOSTS_ADVANCE_DATAGRID = "hostList";
   public static final String ID_VC_HOSTS_LIST = "hostsForVCenter/list";
   public static final String ID_DC_HOSTS_LIST = "hostsForDatacenter/list";
   public static final String ID_LIST_CLUSTERED_HOSTS = "hostsForCluster/list";
   public static final String ID_HOST_STATE_ICON_PREFIX = "hostList/Image_";
   public static final String ID_HOST_STATE_ICON_INDEX = "15";
   public static final String ID_HOST_STATE_ICON_DIFFER = "10";
   public static final String ID_HOSTS_TAB_BUTTON = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("related.hostsView").append(".button")
         .toString();
   public static final String ID_HOST_SETTINGS_VIRTUAL_MACHINE_SPARK_LIST =
         new StringBuffer(IDConstants.EXTENSION_PREFIX)
               .append(IDConstants.EXTENSION_ENTITY_HOST).append("manage.")
               .append("settingsView").append("/tocListVirtual Machines").toString();

   // Related Items Tab
   public static final String ID_RELATED_ITEMS_TAB = "Related Objects";
   public static final String ID_VIRTUAL_MACHINE_TAB = EXTENSION_PREFIX
         + "<ENTITY_TYPE>related.vmsView";
   public static final String ID_VIRTUAL_MACHINES_BUTTON = new StringBuffer(
         ID_VIRTUAL_MACHINE_TAB).append(".").append("button").toString();
   public static final String ID_VIRTUAL_MACHINES_VIEW = new StringBuffer(
         ID_VIRTUAL_MACHINE_TAB).append(".").append("relationshipsView").toString();

   // Manage Tab
   public static final String ID_MANAGE_TAB = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("manage.").append("container").toString();
   public static final String ID_SETTINGS_BUTTON = new StringBuffer(ID_MANAGE_TAB)
         .append(".").append("button").toString();
   public static final String ID_MANAGE_TAB_NAME = "Manage";
   public static final String ID_SETTINGS_BUTTON_VM_MANAGETAB = new StringBuffer(
         EXTENSION_PREFIX).append("vm.manage.settings.").append("button").toString();

   // Settings Tab
   public static final String ID_TAB_SETTINGS = "Settings";
   public static final String ID_SYSTEM_LIST = "tocTree";

   // Monitor Tab
   public static final String ID_MONITOR_TAB = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("monitor.").append("container").toString();
   public static final String ID_MONITOR_TAB_NAME = "Monitor";
   public static final String ID_PROFILE_COMPLIANCE = "Profile Compliance";
   public static final String ID_RES_MGMT_TAB = EXTENSION_PREFIX
         + "<ENTITY_TYPE>monitor.resMgmtView";
   //Related Items Tab
   public static final String ID_RELATED_TAB = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("related.").append("container").toString();
   // Summary Tab
   public static final String ID_SUMMARY_TAB = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("summary.").append("container").toString();
   public static final String ID_STATUS_PORTLET = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("summary.statusView.chrome").toString();
   public static final String ID_HOSTS_GRID = new StringBuffer(ID_HOSTS_TAB).append(
         "/hostList").toString();
   public static final String ID_SUMMARY_TAB_NAME = "Summary";

   // Getting Started Tab
   public static final String ID_GETTING_STATRTED_TAB_NAME = "Getting Started";

   // Tasks Tab
   public static final String ID_TASKS_TAB = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("taskView").append(".container").toString();
   public static final String ID_TASKS_TAB_ADVANCE_DATAGRID = "taskGrid";
   public static final String ID_TASKS_TAB_TASK_NAME = "_TaskDetailsPane_Text1";
   public static final String ID_TASKS_TAB_TASK_STATUS = "statusLabel";
   public static final String ID_TASKS_TAB_TASK_TARGET = "taskEntityLabel";
   public static final String ID_TASKS_TAB_TASK_INITIATOR = "initiatorLabel";
   public static final String ID_TASKS_TAB_TASK_VCENTER = "vCenterLabel";
   public static final String ID_TASKS_BUTTON = "vsphere.core.vm.taskView.button";
   public static final String ID_TASKTAB_BUTTON = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("monitor.").append("tasks.").append("button")
         .toString();
   public static final String ID_TASKS_TAB_VIEW_PANE =
         "vsphere.core.cluster.monitor.tasksView";
   public static final String ID_TASKS_TAB_DETAILS_PANE = new StringBuffer(
         ID_TASKS_TAB_VIEW_PANE).append("/detailsPane").toString();
   public static final String ID_TASKS_TAB_RELATED_EVENTS_PANE =
         "vsphere.opsmgmt.event.relatedEventsView";

   // Events tab
   public static final String ID_EVENTS_TAB_VIEW_PANE =
         "vsphere.opsmgmt.event.cluster.eventView";
   public static final String ID_EVENTS_TAB_DETAILS_PANE = new StringBuffer(
         ID_EVENTS_TAB_VIEW_PANE).append("/eventDetailsView").toString();

   // Issues tab
   public static final String ID_ISSUES_TAB = "vsphere.core.cluster.monitor.issues";
   public static final String ID_ISSUES_TAB_ISSUE_LIST = new StringBuffer(ID_ISSUES_TAB)
         .append("/issueList").toString();
   public static final String ID_ISSUES_TAB_TOC_TREE = new StringBuffer(ID_ISSUES_TAB)
         .append("/tocTree").toString();

   // Permissions Tab
   public static final String ID_PERMISSIONS_TAB_GRID =
         "vsphere.core.cluster.manage.permissions/permissionGrid";

   // Profile Compliance
   public static final String ID_PROFILE_COMPLIANCE_TAB =
         "vsphere.core.cluster.monitor.profileCompliance.hostsCompliance";
   public static final String ID_PROFILE_COMPLIANCE_TAB_SETTINGS = new StringBuffer(
         ID_PROFILE_COMPLIANCE_TAB).append("/_HostListComplianceView_SettingsBlock1")
         .toString();
   public static final String ID_PROFILE_COMPLIANCE_TAB_CHECK_BUTTON = new StringBuffer(
         ID_PROFILE_COMPLIANCE_TAB).append("/checkComplianceActionButton").toString();
   public static final String ID_PROFILE_COMPLIANCE_TAB_HOST_LIST = new StringBuffer(
         ID_PROFILE_COMPLIANCE_TAB).append("/hostList").toString();

   // Monitor view
   public static final String ID_MONITOR_VIEW = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("monitorView").toString();
   public static final String ID_MONITOR_VIEW_TAB_NAVIGATOR_ID = ID_MONITOR_VIEW
         .replace("<ENTITY_TYPE>", EXTENSION_ENTITY_VM) + "/" + ID_TAB_NAVIGATOR;

   // Virtual Switches Diagram
   public static final String ID_NETWORK_DIAGRAM = "diagram";
   public static final String ID_NETWORK_DIAGRAM_VIEW = "switchPanelView";
   public static final String ID_NETWORK_DIAGRAM_ID_PROPERTY = "name";
   public static final String ID_NETWORK_DIAGRAM_OBJECT_CHILDREN_PROPERTY =
         "items.source";
   public static final String ID_NETWORK_DIAGRAM_HEADER_TEXT = "headerText";
   public static final String ID_NETWORK_DIAGRAM_INFO_BUTTON = "infoButton";
   public static final String ID_NETWORK_DIAGRAM_VLAN_TEXT = "vlanIdText";
   public static final String ID_NETWORK_DIAGRAM_PORTS_IMAGE = "numPortsImage";
   public static final String ID_NETWORK_DIAGRAM_EXPANDABLE_BUTTON = "expandableButton";
   public static final String ID_NETWORK_DIAGRAM_NO_ADAPTERS_LABEL =
         "noPhysicalAdaptersLabel";
   public static final String ID_NETWORK_DIAGRAM_PROPERTY_ARRAY_SEPARATOR = "!!?!!";
   public static final String ID_NETWORK_DIAGRAM_SEPARATOR = ",";
   public static final String ID_NETWORK_DIAGRAM_EXPAND_COLAPSE_IMAGE =
         "expandCollapseImage";
   public static final String ID_NETWORK_DIAGRAM_TITLE_TEXT = "titleText";
   public static final String ID_NETWORK_DIAGRAM_PORT_TITLE_TEXT = "vmNameText";
   public static final String ID_NETWORK_DIAGRAM_PORT_MAC_TEXT = "macAddressText";
   public static final String ID_NETWORK_DIAGRAM_PORT_VM_STATUS_IMAGE = "vmStatusImage";
   public static final String ID_NETWORK_DIAGRAM_PORT_STATUS_IMAGE = "portStatusImage";
   public static final String ID_SWITCH_DIAGRAM_SWITCH_LABEL = new StringBuffer(
         ID_STANDARD_SWITCH_DIAGRAM).append("/switchNameLabel").toString();

   // DVSwitch - Summary tab
   public static final String ID_DVSWITCH_SUMMARY_STATUS_PORTLET = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
         .append("summary.statusView").append(ID_PORTLET_SUFFIX).toString();
   public static final String ID_DVSWITCH_SUMMARY_DETAILS_PORTLET = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
         .append("summary.detailsView").append(ID_PORTLET_SUFFIX).toString();
   public static final String ID_DVSWITCH_SUMMARY_PVLANSETTINGS_PORTLET =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
               .append("summary.pvlanView").append(ID_PORTLET_SUFFIX).toString();
   public static final String ID_DVSWITCH_STATUS_LABEL =
         ID_DVSWITCH_SUMMARY_STATUS_PORTLET + "/statusLabel";
   public static final String ID_DVSWITCH_ACTIVETASKS_LABEL =
         ID_DVSWITCH_SUMMARY_STATUS_PORTLET + "/tasks";
   public static final String ID_DVSWITCH_SWITCHPANEL_LABEL =
         ID_DVSWITCH_SUMMARY_DETAILS_PORTLET + "/titleLabel.switchBlock";
   public static final String ID_DVSWITCH_MANUFACTURER_LABEL =
         ID_DVSWITCH_SUMMARY_DETAILS_PORTLET + "/manufacturer";
   public static final String ID_DVSWITCH_VERSION_LABEL =
         ID_DVSWITCH_SUMMARY_DETAILS_PORTLET + "/version";
   public static final String ID_DVSWITCH_PORTSPANEL_LABEL =
         ID_DVSWITCH_SUMMARY_DETAILS_PORTLET + "/titleLabel.portsBlock";
   public static final String ID_DVSWITCH_AVAILABLEPORTS_LABEL =
         ID_DVSWITCH_SUMMARY_DETAILS_PORTLET + "/availablePorts";
   public static final String ID_DVSWITCH_TOTALPORTS_LABEL =
         ID_DVSWITCH_SUMMARY_DETAILS_PORTLET + "/totalPorts";
   public static final String ID_DVSWITCH_NETWORKS_LABEL =
         ID_DVSWITCH_SUMMARY_DETAILS_PORTLET + "/networks";
   public static final String ID_DVSWITCH_HOSTS_LABEL =
         ID_DVSWITCH_SUMMARY_DETAILS_PORTLET + "/hosts";
   public static final String ID_DVSWITCH_VIRTUALMACHINES_LABEL =
         ID_DVSWITCH_SUMMARY_DETAILS_PORTLET + "/vms";
   public static final String ID_DVSWITCH_PVLANSETTINGS_ADVANCEDDATAGRID =
         EXTENSION_PREFIX + EXTENSION_ENTITY_DV_SWITCH + "pvlanView";

   // DVSwitch - ResourceManagement tab
   public static final String ID_RESOURCE_MANAGEMENT_STATUS =
         "resourceManagementStatusLbl";
   public static final String ID_RESOURCE_MANAGEMENT_POOLLIST = "resourcePoolList";
   public static final String ID_TOTAL_PHYSICAL_ADAPTERS = "totalPhysicalAdaptersLbl";
   public static final String ID_TOTAL_BANDWIDTH = "totalBandwidthLbl";

   // dvPortgroup and uplink group - summary tab
   public static final String ID_PORTGROUP_EXTENSION_PREFIX = new StringBuffer(
         EXTENSION_PREFIX).append("<ENTITY_TYPE>").toString();
   public static final String ID_PORTGROUP_SUMMARY_STATUS_PORTLET = new StringBuffer(
         ID_PORTGROUP_EXTENSION_PREFIX).append("summary.statusView").toString();
   public static final String ID_PORTGROUP_STATUS_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_STATUS_PORTLET).append(
         "/statusDisplay/entityStatus/statusLabel").toString();
   public static final String ID_PORTGROUP_ACTIVETASKS_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_STATUS_PORTLET).append("/tasks").toString();
   public static final String ID_PORTGROUP_SUMMARY_DETAILS_PORTLET = new StringBuffer(
         ID_PORTGROUP_EXTENSION_PREFIX).append("summary.detailsView.chrome").toString();
   public static final String ID_PORTGROUP_PORTSPANEL_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_DETAILS_PORTLET).append("/titleLabel.portsBlock")
         .toString();
   public static final String ID_PORTGROUP_AVAILABLEPORTS_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_DETAILS_PORTLET).append("/availablePorts").toString();
   public static final String ID_PORTGROUP_TOTALPORTS_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_DETAILS_PORTLET).append("/totalPorts").toString();
   public static final String ID_PORTGROUP_PORTBINDING_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_DETAILS_PORTLET).append("/pgType").toString();
   public static final String ID_PORTGROUP_SWITCHNAME_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_DETAILS_PORTLET).append("/switchName").toString();
   public static final String ID_PORTGROUP_IPPOOL_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_DETAILS_PORTLET).append("/ipPoolName").toString();
   public static final String ID_PORTGROUP_HOSTS_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_DETAILS_PORTLET).append("/hosts").toString();
   public static final String ID_PORTGROUP_VIRTUALMACHINES_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_DETAILS_PORTLET).append("/vms").toString();
   public static final String ID_PORTGROUP_SUMMARY_POLICIES_PORTLET = new StringBuffer(
         ID_PORTGROUP_EXTENSION_PREFIX).append("summary.policyView.chrome").toString();
   public static final String ID_PORTGROUP_SECURITYPOLICYPANEL_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET)
         .append("/titleLabel.securityPolicyBlock").toString();
   public static final String ID_PORTGROUP_PROMISCUOUSMODE_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/promiscuousMode").toString();
   public static final String ID_PORTGROUP_MACADDRESSCHANGE_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/macAddressChanges").toString();
   public static final String ID_PORTGROUP_FORGEDTRANSMITS_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/forgedTransmits").toString();
   public static final String ID_PORTGROUP_VLANPOLICYPANEL_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/titleLabel.vlanBlock")
         .toString();
   public static final String ID_PORTGROUP_VLANTYPE_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/vlanType").toString();
   public static final String ID_PORTGROUP_VLANID_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/vlanId").toString();
   public static final String ID_PORTGROUP_INSHAPINGPOLICYPANEL_LABEL =
         new StringBuffer(ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append(
               "/titleLabel.inShapingBlock").toString();
   public static final String ID_PORTGROUP_INSHAPINGSTATUS_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/inShapingStatus").toString();
   public static final String ID_PORTGROUP_INSHAPINGAVGBW_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/inShapingAvgBW").toString();
   public static final String ID_PORTGROUP_INSHAPINGPEAKGBW_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/inShapingPeakBW").toString();
   public static final String ID_PORTGROUP_INSHAPINGBURSTSIZE_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/inShapingBurstSize").toString();
   public static final String ID_PORTGROUP_OUTSHAPINGPOLICYPANEL_LABEL =
         new StringBuffer(ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append(
               "/titleLabel.outShapingBlock").toString();
   public static final String ID_PORTGROUP_OUTSHAPINGSTATUS_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/outShapingStatus").toString();
   public static final String ID_PORTGROUP_OUTSHAPINGAVGBW_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/outShapingAvgBW").toString();
   public static final String ID_PORTGROUP_OUTSHAPINGPEAKGBW_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/outShapingPeakBW").toString();
   public static final String ID_PORTGROUP_OUTSHAPINGBURSTSIZE_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/outShapingBurstSize")
         .toString();
   public static final String ID_PORTGROUP_TEAMINGPOLICYPANEL_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/titleLabel.teamingBlock")
         .toString();
   public static final String ID_PORTGROUP_TEAMINGLOADBALANCING_LABEL =
         new StringBuffer(ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append(
               "/teamingLoadBalancing").toString();
   public static final String ID_PORTGROUP_TEAMINGNWFAILOVER_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/teamingNWFailoverDetecetion")
         .toString();
   public static final String ID_PORTGROUP_TEAMINGNOTIFYSWITCHES_LABEL =
         new StringBuffer(ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append(
               "/teamingNotifySwitches").toString();
   public static final String ID_PORTGROUP_TEAMINGFAILBACK_LABEL = new StringBuffer(
         ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append("/teamingFailback").toString();
   public static final String ID_PORTGROUP_TEAMINGACTIVEUPLINKS_LABEL =
         new StringBuffer(ID_PORTGROUP_SUMMARY_POLICIES_PORTLET)
               .append("/activeUplinks").toString();
   public static final String ID_PORTGROUP_TEAMINGSTANDBYUPLINKS_LABEL =
         new StringBuffer(ID_PORTGROUP_SUMMARY_POLICIES_PORTLET).append(
               "/standbyUplinks").toString();
   public static final String ID_DVUPLINKS_STATUS_PORTLET =
         "vsphere.core.dvPortgroup.statusView.chrome";

   // Resource Management portlet
   public static final String ID_RESOURCEMGMT_PORTLET = EXTENSION_PREFIX
         + "<ENTITY_TYPE>monitor.resMgmtView.container";
   public static final String ID_RESOURCE_ALLOCATION_PORTLET = EXTENSION_PREFIX
         + "<ENTITY_TYPE>resAllocationView.chrome";

   // ResourceManagement Tab
   public static final String ID_RESOURCE_MANAGEMENT =
         "com.vmware.vsphere.client.resmgmtui.subView.RMtab";
   public static final String ID_RESOURCE_MANAGEMENT_TOC = "tocListOther";
   public static final String ID_RESOURCE_MANAGEMENT_CPU_TAB =
         "vsphere.core.resMgmt.cpuAllocationView.container";
   public static final String ID_RESOURCE_MANAGEMENT_MEM_TAB =
         "vsphere.core.resMgmt.memoryAllocationView.container";
   public static final String ID_RESOURCE_MANAGEMENT_STORAGE_TAB =
         "vsphere.core.resMgmt.storageAllocationView.container";
   public static final String ID_TOTAL_CAPACITY_VALUE = "totalCapacityValue";
   public static final String ID_CPU_TOTAL_CAPACITY = ID_RESOURCE_MANAGEMENT_CPU_TAB
         + "/" + ID_TOTAL_CAPACITY_VALUE;
   public static final String ID_MEM_TOTAL_CAPACITY = ID_RESOURCE_MANAGEMENT_MEM_TAB
         + "/" + ID_TOTAL_CAPACITY_VALUE;
   public static final String ID_RESERVED_CAPACITY_VALUE = "reservedCapacityValue";
   public static final String ID_CPU_RESERVED_CAPACITY = ID_RESOURCE_MANAGEMENT_CPU_TAB
         + "/" + ID_RESERVED_CAPACITY_VALUE;
   public static final String ID_MEM_RESERVED_CAPACITY = ID_RESOURCE_MANAGEMENT_MEM_TAB
         + "/" + ID_RESERVED_CAPACITY_VALUE;
   public static final String ID_AVAILABLE_CAPACITY_VALUE = "availableCapacityValue";
   public static final String ID_CPU_AVAILABLE_CAPACITY = ID_RESOURCE_MANAGEMENT_CPU_TAB
         + "/" + ID_AVAILABLE_CAPACITY_VALUE;
   public static final String ID_MEM_AVAILABLE_CAPACITY = ID_RESOURCE_MANAGEMENT_MEM_TAB
         + "/" + ID_AVAILABLE_CAPACITY_VALUE;
   public static final String ID_CPU_RESERVATION_TYPE = ID_RESOURCE_MANAGEMENT_CPU_TAB
         + "/reservationTypeLbl";
   public static final String ID_MEM_RESERVATION_TYPE = ID_RESOURCE_MANAGEMENT_MEM_TAB
         + "/reservationTypeLbl";
   public static final String RESOURCE_CPU_EXPECTED = "CPU";
   public static final String RESOURCE_MEM_EXPECTED = "Memory";
   public static final String ID_CPU_RESOURCE_GRID = ID_RESOURCE_MANAGEMENT_CPU_TAB
         + "/resourceGridView";
   public static final String ID_MEM_RESOURCE_GRID = ID_RESOURCE_MANAGEMENT_MEM_TAB
         + "/resourceGridView";
   public static final String ID_STORAGE_RESOURCE_GRID =
         ID_RESOURCE_MANAGEMENT_STORAGE_TAB + "/storageGrid";
   public static final String ID_RESOURCE_MGMT_TAB_NAME = "Resource Management";
   public static final String ID_VM_IMG = ID_CPU_RESOURCE_GRID + "/Image_8";
   public static final String ID_STORAGE_TAB_VM_IMG = ID_STORAGE_RESOURCE_GRID
         + "/Image_8";
   public static final String ID_STORAGE_RESOURCE_GRID_TOTAL_CAPACITY =
         ID_RESOURCE_MANAGEMENT_STORAGE_TAB + "/totalCapacityLbl";
   public static final String ID_STORAGE_RESOURCE_GRID_AVAILABLE_CAPACITY =
         ID_RESOURCE_MANAGEMENT_STORAGE_TAB + "/availableCapacityLbl";
   public static final String ID_RESOURCEMGMT_TAB_NAVIGATOR = "label="
         + ID_RESOURCE_MGMT_TAB_NAME + "/tabNavigator";
   public static final String ID_ALLOCATION_METER = "allocationMeter";
   public static final String ID_RESOURCE_MANAGEMENT_CPU_TAB_GRAPH =
         ID_RESOURCE_MANAGEMENT_CPU_TAB + "/" + ID_ALLOCATION_METER;
   public static final String ID_RESOURCE_MANAGEMENT_MEM_TAB_GRAPH =
         ID_RESOURCE_MANAGEMENT_MEM_TAB + "/" + ID_ALLOCATION_METER;
   public static final String ID_RESOURCE_MANAGEMENT_STORAGE_TAB_GRAPH =
         ID_RESOURCE_MANAGEMENT_STORAGE_TAB + "/" + ID_ALLOCATION_METER;
   public static final String ID_RESOURCE_MANAGEMENT_TAB = EXTENSION_PREFIX
         + "<ENTITY_TYPE>monitor.resMgmtView";
   public static final String ID_TASKS_TAB_VIEW = EXTENSION_PREFIX
         + "<ENTITY_TYPE>monitor.tasks";
   public static final String ID_EVENTS_TAB_VIEW = EXTENSION_PREFIX
         + "<ENTITY_TYPE>monitor.events";
   public static final String ID_UTILIZATION_TAB_VIEW = EXTENSION_PREFIX
         + "<ENTITY_TYPE>utilizationView";
   public static final String ID_RESOURCE_MANAGEMENT_BUTTON = new StringBuffer(
         ID_RESOURCE_MANAGEMENT_TAB).append(".").append("button").toString();
   public static final String ID_TASKS_TAB_BUTTON = new StringBuffer(ID_TASKS_TAB_VIEW)
         .append(".").append("button").toString();
   public static final String ID_EVENTS_TAB_BUTTON =
         new StringBuffer(ID_EVENTS_TAB_VIEW).append(".").append("button").toString();
   public static final String ID_UTILIZATION_TAB_BUTTON = new StringBuffer(
         ID_UTILIZATION_TAB_VIEW).append(".").append("button").toString();

   public static final String ID_VIEW_EVENT_DETAILS = "eventDetailsView";
   public static final String ID_LABEL_EVENT_CREATION_TIME = "eventCreationTimeLabel";
   public static final String ID_LABEL_EVENT_USER = "eventUserLabel";
   public static final String ID_LABEL_DESCRIPTION = "descriptionLabel";
   public static final String ID_LABEL_EVENT_TARGET = "eventRelatedTargetLabel";
   public static final String ID_LABEL_EVENT_TYPE = "eventTypeLabel";
   public static final String ID_LABEL_EVENT_TASK = "eventRelatedTaskLabel";

   public static final String ID_UID_TASK_VIEW =
         "vsphere.core.cluster.monitor.tasksView/detailsBox/detailsPane";
   public static final String ID_UID_MONITOR_TASKS =
         "vsphere.core.cluster.monitor.tasks";

   public static final String ID_RESOURCE_MANAGEMENT_UTIL_TAB_GUEST_MEMORY_FIRST_GRAPH =
         EXTENSION_PREFIX + EXTENSION_ENTITY
               + "monitor.utilization.guestMemoryView/guestMemoryMeterGroup/firstMeter";
   public static final String ID_RESOURCE_MANAGEMENT_UTIL_TAB_GUEST_MEMORY_SECOND_GRAPH =
         EXTENSION_PREFIX + EXTENSION_ENTITY
               + "monitor.utilization.guestMemoryView/guestMemoryMeterGroup/secondMeter";
   public static final String ID_RESOURCE_MANAGEMENT_UTIL_TAB_HOST_MEMORY_CONSUMED_GRAPH =
         EXTENSION_PREFIX
               + EXTENSION_ENTITY
               + "monitor.utilization.hostMemoryView/resConsumeMeterGroup/resConsumedMeter";
   public static final String ID_RESOURCE_MANAGEMENT_UTIL_TAB_HOST_MEMORY_OVERHEAD_GRAPH =
         EXTENSION_PREFIX
               + EXTENSION_ENTITY
               + "monitor.utilization.hostMemoryView/resConsumeMeterGroup/resActiveOverheadMeter";
   public static final String ID_RESOURCE_MANAGEMENT_UTIL_TAB_HOST_CPU_CONSUMED_GRAPH =
         EXTENSION_PREFIX + EXTENSION_ENTITY
               + "monitor.utilization.hostCPUView/resConsumeMeterGroup/resConsumedMeter";
   public static final String ID_RESOURCE_MANAGEMENT_UTIL_TAB_HOST_CPU_ACTIVE_GRAPH =
         EXTENSION_PREFIX
               + EXTENSION_ENTITY
               + "monitor.utilization.hostCPUView/resConsumeMeterGroup/resActiveOverheadMeter";
   public static final String ID_ROOT_PANEL_RESMGMT = ID_RESOURCEMGMT_PORTLET
         + "/rootChrome";
   public static final String ID_VM_HOST_MEMORY_PORTLET = ID_HOST_MEMORY_PORTLET
         .replace("utilization.", "resMgmt.");
   public static final String ID_VM_HOST_CPU_PORTLET = ID_HOST_CPU_PORTLET.replace(
         "utilization.",
         "resMgmt.");
   public static final String ID_VM_HOST_MEMORY_CONSUMED =
         "vsphere.core.vm.monitor.resMgmt.hostMemoryView/consumedResourceLbl";
   public static final String ID_VM_HOST_MEMORY_OVERHEAD =
         "vsphere.core.vm.monitor.resMgmt.hostMemoryView/activeOverheadResourceLbl";
   public static final String ID_VM_HOST_CPU_CONSUMED =
         "vsphere.core.vm.monitor.resMgmt.hostCPUView/consumedResourceLbl";
   public static final String ID_VM_HOST_CPU_ACTIVE =
         "vsphere.core.vm.monitor.resMgmt.hostCPUView/activeOverheadResourceLbl";

   // Resource Pool Summary Tab
   public static final String STATUS_PORTLET_EXPECTED = "Status";
   public static final String RESOURCE_SETTING_EXPECTED = "Resource Settings";
   public static final String RESOURCE_CONSUMER_EXPECTED = "Resource Consumers";
   public static final String GUEST_MEMORY_EXPECTED = "Guest Memory";
   public static final String HOST_CPU_EXPECTED = "Host CPU";
   public static final String HOST_MEMORY_EXPECTED = "Host Memory";
   public static final String ID_RESOURCE_POOL_SUMMARY_TAB_NAME = "Summary";
   public static final String ID_RESOURCE_POOL_SUMMARY_TAB = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_RESOURCE_POOL)
         .append("summary.container").toString();
   public static final String ID_RESOURCE_POOL_RESOURCE_MANAGEMENT_TAB =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_RESOURCE_POOL)
               .append("monitor.resMgmtView.container").toString();
   public static final String ID_RESOURCE_POOL_STATUS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_RESOURCE_POOL)
         .append("summary.statusView.chrome").toString();
   public static final String ID_RESOURCE_POOL_RESOURCE_SETTING = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_RESOURCE_POOL)
         .append("summary.rsrcSettingsView.chrome").toString();
   public static final String ID_RESOURCE_POOL_RESOURCE_CONSUMER = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_RESOURCE_POOL)
         .append("summary.rsrcConsumersView.chrome").toString();
   public static final String ID_RESOURCE_POOL_GUEST_MEMORY = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_RESOURCE_POOL)
         .append("monitor.utilization.guestMemoryView.chrome").toString();
   public static final String ID_RESOURCE_POOL_HOST_CPU = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_RESOURCE_POOL)
         .append("monitor.utilization.hostCPUView.chrome").toString();
   public static final String ID_CLUSTER_MONITOR_GUEST_MEMORY = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER)
         .append("monitor.utilization.guestMemoryView.chrome").toString();
   public static final String ID_CLUSTER_MONITOR_HOST_CPU = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_CLUSTER)
         .append("monitor.utilization.hostCPUView.chrome").toString();
   public static final String ID_CLUSTER_MONITOR_HOST_MEMORY = new StringBuffer(
         EXTENSION_PREFIX + EXTENSION_ENTITY_CLUSTER
               + "monitor.utilization.hostMemoryView.chrome").toString();
   public static final String ID_RESOURCE_POOL_HOST_MEMORY = new StringBuffer(
         EXTENSION_PREFIX + EXTENSION_ENTITY_RESOURCE_POOL
               + "monitor.utilization.hostMemoryView.chrome").toString();
   public static final String ID_RESOURCE_POOL_HOST_MEMORY_CONFIGURED =
         new StringBuffer(IDConstants.ID_RESOURCE_POOL_HOST_MEMORY).append(
               "/resConsumeGridView/configuredMemLabel").toString();
   public static final String ID_RESOURCE_POOL_HOST_MEMORY_LIMIT = new StringBuffer(
         IDConstants.ID_RESOURCE_POOL_HOST_MEMORY)
         .append("/resConsumeGridView/limitLbl").toString();
   public static final String ID_RESOURCE_POOL_HOST_MEMORY_RESERVATION =
         new StringBuffer(IDConstants.ID_RESOURCE_POOL_HOST_MEMORY).append(
               "/resConsumeGridView/reservationLbl").toString();
   public static final String ID_RESOURCE_POOL_HOST_MEMORY_SHARES = new StringBuffer(
         IDConstants.ID_RESOURCE_POOL_HOST_MEMORY).append(
         "/resConsumeGridView/sharesLabel").toString();
   public static final String ID_RESOURCE_POOL_HOST_CPU_LIMIT = new StringBuffer(
         IDConstants.ID_RESOURCE_POOL_HOST_CPU).append("/resConsumeGridView/limitLbl")
         .toString();
   public static final String ID_RESOURCE_POOL_HOST_CPU_RESERVATION = new StringBuffer(
         IDConstants.ID_RESOURCE_POOL_HOST_CPU).append(
         "/resConsumeGridView/reservationLbl").toString();
   public static final String ID_RESOURCE_POOL_HOST_CPU_SHARES = new StringBuffer(
         IDConstants.ID_RESOURCE_POOL_HOST_CPU)
         .append("/resConsumeGridView/sharesLabel").toString();
   public static final String ID_RESOURCE_POOL_STATUS_OVERALL =
         "entityStatus/statusLabel";
   public static final String ID_RESOURCE_POOL_RESOURCE_SETTING_CPU_SHARES = "cpuShares";
   public static final String ID_RESOURCE_POOL_RESOURCE_SETTING_CPU_RESERVE =
         "cpuReserve";
   public static final String ID_RESOURCE_POOL_RESOURCE_SETTING_CPU_LIMIT = "cpuLimit";
   public static final String ID_RESOURCE_POOL_RESOURCE_SETTING_CPU_WORSTCASE =
         "cpuWorstCase";
   public static final String ID_RESOURCE_POOL_RESOURCE_SETTING_MEM_SHARES = "memShares";
   public static final String ID_RESOURCE_POOL_RESOURCE_SETTING_MEM_RESERVE =
         "memReserve";
   public static final String ID_RESOURCE_POOL_RESOURCE_SETTING_MEM_LIMIT = "memLimit";
   public static final String ID_RESOURCE_POOL_RESOURCE_SETTING_MEM_WORSTCASE =
         "memWorstCase";
   public static final String ID_RESOURCE_POOL_RESOURCE_SETTING_CPU_LABEL =
         "titleLabel._ResourcePoolSettingsView_StackBlock1";
   public static final String ID_RESOURCE_POOL_RESOURCE_SETTING_MEM_LABEL =
         "titleLabel._ResourcePoolSettingsView_StackBlock2";
   public static final String ID_RESOURCE_POOL_RESOURCE_CONSUMERS_VM_TEMPLET_LABEL =
         "titleLabel._ResourcePoolConsumersView_StackBlock1";
   public static final String ID_RESOURCE_POOL_RESOURCE_CONSUMERS_POWERED_ON_VM_LABEL =
         "titleLabel._ResourcePoolConsumersView_StackBlock2";
   public static final String ID_RESOURCE_POOL_RESOURCE_CONSUMERS_CHILD_RP_LABEL =
         "titleLabel._ResourcePoolConsumersView_StackBlock3";
   public static final String ID_RESOURCE_POOL_RESOURCE_CONSUMERS_VM_TEMPLET_THIS_RP =
         "vmDirectCount";
   public static final String ID_RESOURCE_POOL_RESOURCE_CONSUMERS_VM_TEMPLET_TOTAL_DESCEND =
         "vmTotalCount";
   public static final String ID_RESOURCE_POOL_RESOURCE_CONSUMERS_POWERED_ON_VM_THIS_RP =
         "vmPoweredOnDirectCount";
   public static final String ID_RESOURCE_POOL_RESOURCE_CONSUMERS_POWERED_ON_VM_TOTAL_DESCEND =
         "vmPoweredOnTotalCount";
   public static final String ID_RESOURCE_POOL_RESOURCE_CONSUMERS_CHILD_RP_THIS_RP =
         "rpDirectCount";
   public static final String ID_RESOURCE_POOL_RESOURCE_CONSUMERS_CHILD_RP_TOTAL_DESCEND =
         "rpTotalCount";
   public static final String ID_RESOURCE_POOL_GUEST_MEMORY_ACTIVE_GUEST_MEMORY =
         "activeMemoryLbl";
   public static final String ID_RESOURCE_POOL_GUEST_MEMORY_PRIVATE = "privateMemoryLbl";
   public static final String ID_RESOURCE_POOL_GUEST_MEMORY_SHARED = "sharedMemoryLbl";
   public static final String ID_RESOURCE_POOL_GUEST_MEMORY_BALLOONED =
         "balloonedMemoryLbl";
   public static final String ID_RESOURCE_POOL_GUEST_MEMORY_COMPRESSED =
         "compressedMemoryLbl";
   public static final String ID_RESOURCE_POOL_GUEST_MEMORY_SWAPPED = "swappedMemoryLbl";
   public static final String ID_RESOURCE_POOL_GUEST_MEMORY_UNACCESSED =
         "unaccessedMemoryLbl";
   public static final String ID_RESOURCE_POOL_HOST_MEMORY_CONSUMED = new StringBuffer(
         ID_RESOURCE_POOL_HOST_MEMORY).append("/resConsumeGridView/consumedResourceLbl")
         .toString();
   public static final String ID_RESOURCE_POOL_HOST_MEMORY_OVERHEAD = new StringBuffer(
         ID_RESOURCE_POOL_HOST_MEMORY).append(
         "/resConsumeGridView/activeOverheadResourceLbl").toString();
   public static final String ID_RESOURCE_POOL_HOST_CPU_CONSUMED = new StringBuffer(
         ID_RESOURCE_POOL_HOST_CPU).append("/resConsumeGridView/consumedResourceLbl")
         .toString();
   public static final String ID_RESOURCE_POOL_HOST_CPU_ACTIVE = new StringBuffer(
         ID_RESOURCE_POOL_HOST_CPU).append(
         "/resConsumeGridView/activeOverheadResourceLbl").toString();
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_LINK =
         "editSettingsBtn";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL =
         "tiwoDialog";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_OK_BUTTON =
         "tiwoDialog/okButton";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CANCEL_BUTTON =
         "tiwoDialog/cancelButton";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CPU_RESERVATION_COMBO =
         "cpuConfigControl/reservationCombo/topRow/comboBox";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CPU_LIMIT_COMBO =
         "cpuConfigControl/limitCombo/topRow/comboBox";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_MEMORY_RESERVATION_COMBO =
         "memoryConfigControl/reservationCombo/topRow/comboBox";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_MEMORY_LIMIT_COMBO =
         "memoryConfigControl/limitCombo/topRow/comboBox";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CPU_SHARES =
         "cpuConfigControl/sharesControl/levels";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CPU_SHARE_NUMERIC_STEPPER =
         "cpuConfigControl/numShares";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CPU_RESERVATION_UNIT =
         "cpuConfigControl/reservationCombo/topRow/units";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CPU_LIMIT_UNIT =
         "cpuConfigControl/limitCombo/topRow/units";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_MEMORY_SHARES =
         "memoryConfigControl/sharesControl/levels";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_MEMORY_SHARE_NUMERIC_STEPPER =
         "memoryConfigControl/numShares";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_MEMORY_RESERVATION_UNIT =
         "memoryConfigControl/reservationCombo/topRow/units";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_MEMORY_LIMIT_UNIT =
         "memoryConfigControl/limitCombo/topRow/units";
   public static final String ID_RESOURCE_POOL_RESOURCE_MANAGEMENT_ADVANCE_DATAGRID_CPU =
         ID_RESOURCE_MANAGEMENT_CPU_TAB + "/resourceGridView";
   public static final String ID_RESOURCE_POOL_RESOURCE_MANAGEMENT_ADVANCE_DATAGRID_MEMORY =
         ID_RESOURCE_MANAGEMENT_MEM_TAB + "/resourceGridView";;
   public static final String ID_RESOURCE_POOL_RESOURCE_MANAGEMENT_MEMORY_TAB =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_RESOURCE_POOL)
               .append("resAllocationView/resAllocationTabNavigator").toString();
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CPU_MAX_RESERVATION_LABEL =
         "cpuConfigControl/reservation/reservationCombo/bottomRow/label2";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CPU_MAX_LIMIT_LABEL =
         "limitCombo/bottomRow/label2";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_MEMORY_MAX_RESERVATION_LABEL =
         "memoryConfigControl/reservation/reservationCombo/bottomRow/label2";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_MEMORY_MAX_LIMIT_LABEL =
         "memoryConfigControl/limit/limitCombo/bottomRow/label2";
   public static final String ID_RESOURCE_POOL_SUMMARY = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_RESOURCE_POOL)
         .append("summaryView/rootChrome").toString();
   public static final String ID_NEW_RESOURCE_POOL_NAME_TEXT_INPUT =
         "nameStackBlok/nameTextInput";
   public static final String ID_NEW_RESOURCE_POOL_PANEL_OKBUTTON =
         "tiwoDialog/okButton";
   public static final String ID_NEW_RESOURCE_POOL_PANEL = "tiwoDialog";
   public static final String ID_NEW_RESOURCE_POOL_PANEL_CANCELBUTTON =
         "tiwoDialog/cancelButton";
   public static final String ID_RP_TITLE_LABEL =
         "vsphere.core.resourcePool.objectView/title";
   public static final String ID_LABEL_MANAGE_SETTINGS_CPU_SHARES =
         "cpuResourceViewControl/shares/_ResourceAllocationInfoControl_Label1";
   public static final String ID_LABEL_MANAGE_SETTINGS_CPU_RESERVATION =
         "cpuResourceViewControl/reservation/_ResourceAllocationInfoControl_Label2";
   public static final String ID_LABEL_MANAGE_SETTINGS_CPU_LIMIT =
         "cpuResourceViewControl/limit/_ResourceAllocationInfoControl_Label3";
   public static final String ID_LABEL_MANAGE_SETTINGS_MEMORY_SHARES =
         "memoryResourceViewControl/shares/_ResourceAllocationInfoControl_Label1";
   public static final String ID_LABEL_MANAGE_SETTINGS_MEMORY_RESERVATION =
         "memoryResourceViewControl/reservation/_ResourceAllocationInfoControl_Label2";
   public static final String ID_LABEL_MANAGE_SETTINGS_MEMORY_LIMIT =
         "memoryResourceViewControl/limit/_ResourceAllocationInfoControl_Label3";
   public static final String ID_LABEL_SUMMARY_VM_TEMPLATES_COUNT =
         "summary_vmsAndTemplatesCount_valueLbl";
   public static final String ID_LABEL_SUMMARY_POWERED_ON_COUNT =
         "summary_poweredOnVmsCount_valueLbl";
   public static final String ID_LABEL_SUMMARY_CHILD_RP_COUNT =
         "summary_childResPoolsCount_valueLbl";
   public static final String ID_LABEL_SUMMARY_CHILD_VAPP_COUNT =
         "summary_childVAppsCount_valueLbl";
   public static final String ID_SUMMARY_USAGE_GRID_ROW_CPU = "resourceUsageGridRowId_0";
   public static final String ID_SUMMARY_USAGE_GRID_ROW_MEMORY =
         "resourceUsageGridRowId_1";
   public static final String ID_LABEL_RP_SUMMARY_USAGE = "usageLabel";
   public static final String ID_LABEL_RP_SUMMARY_FREE_CAPACITY = "freeCapacityLabel";

   // Resourcepool ResourceManagement tab
   public static final String ID_LABEL_RESOURCEALLOCATION_CPU_AVAILABLE_RESERVATION =
         "vsphere.core.resMgmt.cpuAllocationView/availableCapacityValue";
   public static final String ID_LABEL_RESOURCEALLOCATION_MEMORY_AVAILABLE_RESERVATION =
         "vsphere.core.resMgmt.memoryAllocationView/availableCapacityValue";
   public static final String ID_RESOURCE_MANAGEMENT_TOCTREE = "tocTree";
   public static final String ID_RESOURCE_MGMT_TOCTREE =
         "vsphere.core.resourcePool.monitor.resMgmtView.wrapper/tocTree";
   public static final String ID_RP_MANAGE_SETTINGS_TOCTREE =
         "vsphere.core.resourcePool.manage.settings.wrapper/tocTree";
   public static final String ID_CLUSTER_RESOURCE_MGMT_TOCTREE =
         "vsphere.core.cluster.monitor.resMgmtView.wrapper/tocTree";
   public static final String ID_HOST_RESOURCE_MGMT_TOCTREE =
         "vsphere.core.host.monitor.resMgmtView.wrapper/tocTree";

   // BUG 503689: Because of this bug we won't need confirm removal dialogue
   // id.
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CONFIRM_REMOVAL =
         "automationName=Confirm removal";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CONFIRM_REMOVAL_YES_BUTTON =
         "automationName=Confirm removal/automationName=Yes";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CONFIRM_REMOVAL_NO_BUTTON =
         "automationName=Confirm removal/automationName=No";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_CPU_LABEL =
         "editSettingsSE/cpuStackBlock/titleLabel.cpuStackBlock";
   public static final String ID_RESOURCE_POOL_EDIT_RESOURCE_SETTING_PANEL_MEMORY_LABEL =
         "editSettingsSE/memoryStackBlock/titleLabel.memoryStackBlock";
   public static final String ID_RESOURCE_POOL_RESOURCE_MANAGEMENT_ADVANCE_DATAGRID_CPU_RPNAME =
         "cpuAllocationTabView/resourceGridView/Label_9";
   public static final String ID_RESOURCE_POOL_RESOURCE_MANAGEMENT_ADVANCE_DATAGRID_MEMORY_RPNAME =
         "memoryAllocationTabView/resourceGridView/Label_9";
   public static final String ID_RESOURCE_POOL_REMOVE_YES_BUTTON =
         "automationName=Remove Resource Pool/automationName=Yes";

   // Datastore Summary Tab
   public static final String ID_DATASTORE_SUMMARY =
         "com.vmware.vsphere.client.dsui.subView.summary";
   public static final String ID_DATASTORE_STATUS_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_DATASTORE + "summary.statusView.chrome";
   public static final String ID_DATASTORE_DETAILS_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_DATASTORE + "summary.detailsView.chrome";
   public static final String ID_DATASTORE_CAPACITY_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_DATASTORE + "summary.capacityView.chrome";
   public static final String ID_DATASTORE_OVERALLSTATUS = "entityStatus/statusLabel";
   public static final String ID_DATASTORE_URL = "dsUrl";
   public static final String ID_DATASTORE_TYPE = "dsType";
   public static final String ID_DATASTORE_HOSTS = "dsHostCount";
   public static final String ID_DATASTORE_VMS = "dsVmCount";
   public static final String ID_DATASTORE_CAPACITY = "dsCapacity";
   public static final String ID_DATASTORE_PROVISIONED = "dsProvisioned";
   public static final String ID_DATASTORE_FREE = "dsFreeSpace";
   public static final String ID_DATASTORE_VMS_LIST = "vmList";
   // Datastore Summary tab - new capabilities, custom attributes portlets
   public static final String ID_DATASTORE_CAPABILITES_PORTLET = EXTENSION_PREFIX
         + "spbm.storage.storageCapabilitiesPortlet.chrome";
   public static final String ID_DATASTORE_CUSTOMATTRIBUTES_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_DATASTORE + "summary.customAttributesView.chrome";
   // Datastore files
   public static String ID_DATASTORE_FILES_TREE = "tree";
   public static String ID_DATASTORE_FILES_DATAGRID = "datagrid";

   // Vapp Summary Tab
   public static final String ID_VAPP_STATUS_PORTLET =
         "vsphere.core.vApp.summary.statusView";
   public static final String ID_VAPP_DETAILS_PORTLET =
         "vsphere.core.vApp.summary.detailsView.chrome";
   public static final String ID_VAPP_ANNOTATION_PORTLET =
         "vsphere.core.vApp.summary.annotationsNotesView.chrome";
   public static final String ID_VAPP_RELATED_ITEMS_PORTLET =
         "vsphere.core.vApp.summary.relatedItemsView.chrome";
   public static final String ID_VAPP_SUMMARY_RELATED_ITEMS_DATASTORE_LIST =
         ID_VAPP_RELATED_ITEMS_PORTLET + "/RelatedItemsServerObjectList0";
   public static final String ID_VAPP_SUMMARY_RELATED_ITEMS_NETWORK_LIST =
         ID_VAPP_RELATED_ITEMS_PORTLET + "/RelatedItemsServerObjectList1";
   public static final String ID_VAPP_HEALTH_STATUS =
         "vsphere.core.vApp.summary.statusView/status/statusLabel";
   public static final String ID_VAPP_STATUS = "vAppStateUrlLabel/plainLabel";
   public static final String ID_VAPP_VERSION = "versionText";
   public static final String ID_VAPP_PRODUCT = "productUrlLabel/plainLabel";
   public static final String ID_VAPP_VENDOR = "vendorUrlLabel/plainLabel";
   public static final String ID_VAPP_VMS_LIST =
         "com.vmware.vsphere.client.views.vm.list.VMList";
   public static final String ID_VAPP_ANNOTATION_EDIT_LINK = "editButton";
   public static final String ID_VAPP_ANNOTATION_TEXTAREA = "annotationTextArea";
   public static final String ID_VAPP_NOTES_TEXTAREA = "notes";
   public static final String ID_VAPP_PRODUCT_URL_LINK = "productUrlLabel/linkButton";
   // Vapp Power Ops
   public static final String ID_VAPP_POWER_ON = "vsphere.core.vApp.powerOnAction";
   public static final String ID_VAPP_POWER_OFF = "vsphere.core.vApp.shutdownAction"; // "vsphere.core.vApp.powerOffAction"
   // ;
   public static final String ID_VAPP_SUSPEND = "vsphere.core.vApp.suspendAction";
   public static final String ID_CONFIRM_SUSPEND_YES = "automationName=Yes";
   public static final String ID_CONFIRM_SUSPEND_NO = "automationName=No";

   // Edit vApp Annotation popup dialog
   public static final String ID_VAPP_EDIT_ANNOTATION_TEXTAREA = "notesTextArea";
   public static final String ID_VAPP_EDIT_ANNOTATION_OKBUTTON = "okButton";
   // Edit vApp Panel
   public static final String ID_VAPP_EDIT_PANEL = "tiwoDialog";
   public static final String ID_VAPP_EDIT_CPU_LABEL = "titleLabel.editCpuResources";
   public static final String ID_VAPP_EDIT_MEMORY_LABEL =
         "titleLabel.editMemoryResources";
   public static final String ID_VAPP_EDIT_UNRECOG_OVF_LABEL =
         "titleLabel.viewUnrecognizedOvfSections";
   public static final String ID_VAPP_EDIT_PRODUCT_LABEL = "titleLabel.editProduct";
   public static final String ID_VAPP_EDIT_PRODUCT_NAME_HEADER_TEXTAREA =
         "headerProductName";
   public static final String ID_VAPP_EDIT_PRODUCT_NAME_TEXTAREA = "productName";
   public static final String ID_VAPP_EDIT_PRODUCT_VERSION_TEXTAREA = "version";
   public static final String ID_VAPP_EDIT_PRODUCT_FULLVERSION_TEXTAREA = "fullVersion";
   public static final String ID_VAPP_EDIT_PRODUCT_URL_TEXTAREA = "productUrl";
   public static final String ID_VAPP_EDIT_PRODUCT_VENDOR_TEXTAREA = "vendor";
   public static final String ID_VAPP_EDIT_PRODUCT_VENDORURL_TEXTAREA = "vendorUrl";
   public static final String ID_VAPP_EDIT_PRODUCT_APPURL_TEXTAREA = "appUrl";
   public static final String ID_VAPP_EDIT_IP_ALLOCATION_LABEL =
         "titleLabel.editIpAllocation";
   public static final String ID_VAPP_EDIT_IP_ALLOCATION_FIXED_RADIO = "fixedButton";
   public static final String ID_VAPP_EDIT_IP_ALLOCATION_TRANSIENT_RADIO =
         "transientButton";
   public static final String ID_VAPP_EDIT_IP_ALLOCATION_DHCP_RADIO = "dhcpButton";
   public static final String ID_VAPP_EDIT_PROPERTIES_BUTTON =
         "titleLabel.editProperties";
   public static final String ID_VAPP_EDIT_PROPERTIES_CATEGORY_LABEL = "categoryLabel";
   public static final String ID_VAPP_EDIT_PROPERTIES_LABEL_LABEL = "labelLabel";
   public static final String ID_VAPP_EDIT_PROPERTIES_LABEL_LABEL1 =
         "propertyStackBlock0PropertyGrid/propertyStackBlock0EditPropertyControl0Row/labelItem/PropertyGridLabel";
   public static final String ID_VAPP_EDIT_PROPERTIES_DESCRIPTION_LABEL =
         "descriptionLabel";
   public static final String ID_VAPP_EDIT_PROPERTIES_DESCRIPTION_LABEL1 =
         "propertyStackBlock0EditPropertyControl0DescriptionLabel";
   public static final String ID_VAPP_EDIT_PROPERTIES_CONFIG_VALUE_TEXTAREA =
         "appPropertiesPage/textDisplay";
   public static final String ID_VAPP_EDIT_PROPERTIES_CONFIG_VALUE_PASSWORD_SPARK_TEXT_INPUT =
         "appPropertiesPage/propertyStackBlock0EditPropertyControl0PasswordGridInput/textDisplay";
   public static final String ID_VAPP_EDIT_PROPERTIES_CONFIG_VALUE_CONFIRM_PASSWORD_SPARK_TEXT_INPUT =
         "appPropertiesPage/propertyStackBlock0EditPropertyControl0PasswordGridConfirmation/textDisplay";
   public static final String ID_VAPP_ADVANCED_PROPERTIES_BUTTON =
         "titleLabel.advancedPropertyListView";
   public static final String ID_VAPP_ADVANCED_PROPERTIES_GRID = "grid";
   public static final String ID_VAPP_ADVANCED_PROPERTIES_NEW_BUTTON = "newButton";
   public static final String ID_VAPP_ADVANCED_PROPERTIES_EDIT_BUTTON = "editButton";
   public static final String ID_VAPP_ADVANCED_PROPERTIES_DELETE_BUTTON = "deleteButton";
   public static final String ID_VAPP_EDIT_MEMORY_LIMIT_COMBOBOX =
         "editMemoryResources/limit/limitCombo/topRow/comboBox";
   public static final String ID_VAPP_EDIT_AUTHORING_PANEL =
         "authoringPanel/propertyGrid";
   public static final String ID_VAPP_EDIT_IP_ALLOCATION_PANEL =
         "deploymentPanel/editIpAllocation/ipAllocationControl";
   public static final String ID_VAPP_EDIT_CPU_RESOURCE_ALLOCATION_PANEL =
         "deploymentPanel/editCpuResources/configControl";
   public static final String ID_VAPP_EDIT_MEMORY_RESOURCE_ALLOCATION_PANEL =
         "deploymentPanel/editMemoryResources/configControl";
   public static final String ID_VAPP_CPU_RESOURCE_LABEL = "titleLabel.editCpuResources";
   public static final String ID_VAPP_MEMORY_RESOURCE_LABEL =
         "titleLabel.editMemoryResources";
   public static final String ID_VAPP_PRODUCT_NAME_HEADER = "headerProductName";
   public static final String ID_VAPP_PRODUCT_HEADER = "titleLabel.editProduct";
   public static final String ID_VAPP_EDIT_DHCP_CHECKBOX = "dhcpCheckBox";
   public static final String ID_VAPP_EDIT_OVF_ENV_CHECKBOX = "ovfEnvCheckBox";

   // vApp Edit Property Settings
   public static final String ID_VAPP_EDIT_PROPERTY_DIALOG =
         "advancedPropertyEditDialog";
   public static final String ID_VAPP_PROPERTY_CATEGORY_COMBOBOX = "categoryComboBox";
   public static final String ID_VAPP_PROPERTY_LABEL_TEXTINPUT = "labelInput";
   public static final String ID_VAPP_PROPERTY_KEY_CLASSID_TEXTINPUT = "classIdInput";
   public static final String ID_VAPP_PROPERTY_KEY_ID_TEXTINPUT = "idInput";
   public static final String ID_VAPP_PROPERTY_KEY_INSTANCEID_TEXTINPUT =
         "instanceIdInput";
   public static final String ID_VAPP_PROPERTY_KEY_DESCRIPTION_TEXTAREA =
         "descriptionTextArea";
   public static final String ID_VAPP_PROPERTY_STATIC_RADIO = "staticRadioButton";
   public static final String ID_VAPP_PROPERTY_TYPE_DROPDOWN = "typeDropDown";
   public static final String ID_VAPP_PROPERTY_MIN_LENGTH_TEXTINPUT = "minLengthInput";
   public static final String ID_VAPP_PROPERTY_MAX_LENGTH_TEXTINPUT = "maxLengthInput";
   public static final String ID_VAPP_PROPERTY_DEFAULT_VALUE_TEXTINPUT =
         "defaultValueInput";
   public static final String ID_VAPP_PROPERTY_CHOICE_INPUT = "stringChoiceInput";
   public static final String ID_VAPP_PROPERTY_USER_CONFIGURABLE_CHECKBOX =
         "userConfigurableCheckBox";
   public static final String ID_VAPP_PROPERTY_DYNAMIC_RADIO = "dynamicRadioButton";
   public static final String ID_VAPP_PROPERTY_MACRO_DROPDOWN = "macroDropDown";
   public static final String ID_VAPP_PROPERTY_NETWORK_DROPDOWN = "networkDropDown";
   public static final String ID_VAPP_EDIT_PROPERTY_DIALOG_OK_BUTTON = "okButton";
   public static final String ID_VAPP_EDIT_PROPERTY_DIALOG_CANCEL_BUTTON =
         "cancelButton";
   public static final String ID_VAPP_EDIT_PROPERTY_RANGE_LABEL = "rangeLabel";
   public static final String ID_VAPP_EDIT_PROPERTY_STRING_CHOICE_LABEL =
         "stringChoiceLabel";
   public static final String ID_VAPP_EDIT_PROPERTY_DEFAULT_VALUE_LABEL =
         "defaultValueLabel";
   public static final String ID_VAPP_EDIT_PROPERTY_TRUE_RADIO = "booleanTrueButton";
   public static final String ID_VAPP_EDIT_PROPERTY_FALSE_RADIO = "booleanFalseButton";
   public static final String ID_VAPP_AUTH_IP_PROTOCOL_DROPDOWN = "protocolDropDown";
   public static final String ID_VAPP_DEPLOY_IP_PROTOCOL_DROPDOWN =
         "ipAllocationControl/protocolDropDown";
   public static final String ID_VAPP_DEPLOY_IP_ALLOCATION_DROPDOWN =
         "ipAllocationControl/policyDropDown";
   public static final String ID_VAPP_NO_PROPERTIES_LABEL = "noPropertiesLabel";

   // New vApp
   public static final String ID_NEW_VAPP_SELECT_NAME_FOLDER_PAGE = "nameFolderPage";
   public static final String ID_NEW_VAPP_SELECT_NAME_FOLDER_PAGE_TITLE =
         "pageStack/tiwoDialog/nameFolderPage/pageHeaderTitleElement";
   public static final String ID_DEPLOY_OVF_NAME_LOCATION_PAGE = "nameAndLocationPage";
   public static final String ID_NEW_VAPP_SELECT_NAME_DESTINATION_PAGE =
         "destinationPage";
   public static final String ID_NEW_VAPP_SELECT_RESOURCES_PAGE = "resourcesPage";
   public static final String ID_NEW_VAPP_SUMMARY_PAGE = "defaultSummaryPage";
   public static final String ID_NEW_VAPP_SELECT_RESOURCES_PAGE_TITLE =
         "tiwoDialog/resourcesPage/pageHeaderTitleElement";
   public static final String ID_NEW_VAPP_SUMMARY_PAGE_TITLE =
         "tiwoDialog/defaultSummaryPage/pageHeaderTitleElement";
   public static final String ID_NEW_VAPP_NAME_TEXT_INPUT =
         ID_NEW_VAPP_SELECT_NAME_FOLDER_PAGE + "/nameTextInput";
   public static final String ID_NEW_VAPP_CPU_SHARES_COMBO =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE + "/editCpuResources/sharesControl/levels";
   public static final String ID_NEW_VAPP_CPU_SHARES_NUMERIC_STEPPER =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE + "/editCpuResources/sharesControl/numShares";
   public static final String ID_NEW_VAPP_MEMORY_SHARES_COMBO =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE + "/editMemoryResources/sharesControl/levels";
   public static final String ID_NEW_VAPP_MEMORY_SHARES_NUMERIC_STEPPER =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE
               + "/editMemoryResources/sharesControl/numShares";
   public static final String ID_NEW_VAPP_CPU_RESERVATION_COMBO =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE
               + "/editCpuResources/reservationCombo/comboBox";
   public static final String ID_NEW_VAPP_CPU_RESERVATION_LIMIT_LABEL =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE + "/editCpuResources/reservationCombo/label2";
   public static final String ID_NEW_VAPP_MEMORY_RESERVATION_COMBO =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE
               + "/editMemoryResources/reservationCombo/comboBox";
   public static final String ID_NEW_VAPP_MEMORY_RESERVATION_LIMIT_LABEL =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE
               + "/editMemoryResources/reservationCombo/label2";
   public static final String ID_NEW_VAPP_CPU_RESERVATION_TYPE_CHECKBOX =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE + "/editCpuResources/expReservationCheck";
   public static final String ID_NEW_VAPP_MEMORY_RESERVATION_TYPE_CHECKBOX =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE + "/editMemoryResources/expReservationCheck";
   public static final String ID_NEW_VAPP_CPU_LIMIT_COMBO =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE + "/editCpuResources/limitCombo/comboBox";
   public static final String ID_NEW_VAPP_CPU_LIMIT_MAX_LABEL =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE + "/editCpuResources/limitCombo/label2";
   public static final String ID_NEW_VAPP_MEMORY_LIMIT_COMBO =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE + "/editMemoryResources/limitCombo/comboBox";
   public static final String ID_NEW_VAPP_MEMORY_LIMIT_MAX_LABEL =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE + "/editMemoryResources/limitCombo/label2";
   public static final String ID_VAPP_CPU_SHARES_NUMERIC_STEPPER_INCREMENT_BUTTON =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE
               + "/editCpuResources/sharesControl/numShares/incrementButton";
   public static final String ID_VAPP_CPU_SHARES_NUMERIC_STEPPER_DECREMENT_BUTTON =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE
               + "/editCpuResources/sharesControl/numShares/decrementButton";
   public static final String ID_VAPP_MEMORY_SHARES_NUMERIC_STEPPER_INCREMENT_BUTTON =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE
               + "/editMemoryResources/sharesControl/numShares/incrementButton";
   public static final String ID_VAPP_MEMORY_SHARES_NUMERIC_STEPPER_DECREMENT_BUTTON =
         ID_NEW_VAPP_SELECT_RESOURCES_PAGE
               + "/editMemoryResources/sharesControl/numShares/decrementButton";
   public static final String ID_VAPP_CREATION_TYPE_LIST = "optionsList";

   // Common vApp
   public static final String ID_VAPP_COMMON_CPU_SHARES_LEVEL =
         "/editCpuResources/sharesControl/levels";
   public static final String ID_VAPP_COMMON_CPU_SHARES_NUM_SHARES =
         "/editCpuResources/sharesControl/numShares";
   public static final String ID_VAPP_COMMON_CPU_RESERVATION =
         "/editCpuResources/reservationCombo/comboBox";
   public static final String ID_VAPP_COMMON_CPU_RESERVATION_TYPE =
         "/editCpuResources/expReservationCheck";
   public static final String ID_VAPP_COMMON_CPU_LIMIT =
         "/editCpuResources/limitCombo/comboBox";

   public static final String ID_VAPP_COMMON_MEM_SHARES_LEVEL =
         "/editMemoryResources/sharesControl/levels";
   public static final String ID_VAPP_COMMON_MEM_SHARES_NUM_SHARES =
         "/editMemoryResources/sharesControl/numShares";
   public static final String ID_VAPP_COMMON_MEM_RESERVATION =
         "/editMemoryResources/reservationCombo/comboBox";
   public static final String ID_VAPP_COMMON_MEM_RESERVATION_TYPE =
         "/editMemoryResources/expReservationCheck";
   public static final String ID_VAPP_COMMON_MEM_LIMIT =
         "/editMemoryResources/limitCombo/comboBox";

   // Edit vApp
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK =
         "tiwoDialog/optionsContainer";

   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_CPU_SHARES_LEVEL =
         ID_EDIT_VAPP_DIALOG_OPTIONS_STACK + ID_VAPP_COMMON_CPU_SHARES_LEVEL;
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_CPU_SHARES_VALUE =
         ID_EDIT_VAPP_DIALOG_OPTIONS_STACK + ID_VAPP_COMMON_CPU_SHARES_NUM_SHARES;
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_CPU_RESERVATION =
         ID_EDIT_VAPP_DIALOG_OPTIONS_STACK + ID_VAPP_COMMON_CPU_RESERVATION;
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_CPU_RESERVATION_TYPE =
         ID_EDIT_VAPP_DIALOG_OPTIONS_STACK + ID_VAPP_COMMON_CPU_RESERVATION_TYPE;
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_CPU_LIMIT =
         ID_EDIT_VAPP_DIALOG_OPTIONS_STACK + ID_VAPP_COMMON_CPU_LIMIT;
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_CPU_MAX_LIMIT =
         ID_EDIT_VAPP_DIALOG_OPTIONS_STACK + "/editCpuResources/reservationCombo/label2";

   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_MEM_SHARES_LEVEL =
         ID_EDIT_VAPP_DIALOG_OPTIONS_STACK + ID_VAPP_COMMON_MEM_SHARES_LEVEL;
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_MEM_SHARES_VALUE =
         ID_EDIT_VAPP_DIALOG_OPTIONS_STACK + ID_VAPP_COMMON_MEM_SHARES_NUM_SHARES;
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_MEM_RESERVATION =
         ID_EDIT_VAPP_DIALOG_OPTIONS_STACK + ID_VAPP_COMMON_MEM_RESERVATION;
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_MEM_RESERVATION_TYPE =
         ID_EDIT_VAPP_DIALOG_OPTIONS_STACK + ID_VAPP_COMMON_MEM_RESERVATION_TYPE;
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_MEM_LIMIT =
         ID_EDIT_VAPP_DIALOG_OPTIONS_STACK + ID_VAPP_COMMON_MEM_LIMIT;
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_MEM_MAX_LIMIT =
         ID_EDIT_VAPP_DIALOG_OPTIONS_STACK
               + "/editMemoryResources/reservationCombo/label2";
   public static final String ID_NEW_VAPP_CREATION_TYPE_PAGE = "creationTypePage";

   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_CPU_SHARES_NUM_STEPPER_INC =
         IDConstants.ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_CPU_SHARES_VALUE
               + "/incrementButton";
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_CPU_SHARES_NUM_STEPPER_DEC =
         IDConstants.ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_CPU_SHARES_VALUE
               + "/decrementButton";
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_MEM_SHARES_NUM_STEPPER_INC =
         IDConstants.ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_MEM_SHARES_VALUE
               + "/incrementButton";
   public static final String ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_MEM_SHARES_NUM_STEPPER_DEC =
         IDConstants.ID_EDIT_VAPP_DIALOG_OPTIONS_STACK_MEM_SHARES_VALUE
               + "/decrementButton";

   public static final String ID_EDIT_VAPP_DIALOG_PRODUCT_PRODUCT_URL_BUTTON =
         "viewProductUrlButton";
   public static final String ID_EDIT_VAPP_DIALOG_PRODUCT_VENDOR_URL_BUTTON =
         "viewVendorUrlButton";


   // Migration Wizard
   public static final String ID_MIGRATION_WIZARD_PANEL = "tiwoDialog";
   public static final String ID_MIGRATION_WIZARD_NEXT_BUTTON =
         ID_MIGRATION_WIZARD_PANEL + "/next";
   public static final String ID_MIGRATION_WIZARD_PREVIOUS_BUTTON =
         ID_MIGRATION_WIZARD_PANEL + "/back";
   public static final String ID_MIGRATION_WIZARD_FINISH_BUTTON =
         ID_MIGRATION_WIZARD_PANEL + "/finish";
   public static final String ID_MIGRATION_WIZARD_CANCEL_BUTTON =
         ID_MIGRATION_WIZARD_PANEL + "/cancel";
   public static final String ID_MIGRATION_WIZARD_RP_PAGE_NAVTREE =
         ID_MIGRATION_WIZARD_PANEL + "/selectResourcePoolPage/navTreeView/navTree";
   public static final String ID_MIGRATION_WIZARD_PAGE_TITLE = ID_MIGRATION_WIZARD_PANEL
         + "/title";
   public static final String ID_MIGRATION_CHANGE_HOST_DESCRIPTION_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/changeHostDesc";
   public static final String ID_MIGRATION_COMPATIBILITY_CHECK_SUCCESS_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectResourcePoolPage/complianceSuccess";
   public static final String ID_MIGRATION_COMPATIBILITY_CHECK_SUCCESS_SELECT_HOST_PAGE_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectHostPage/complianceSuccess";
   public static final String ID_MIGRATION_COMPATIBILITY_CHECK_SUCCESS_SELECT_DATASTORE_PAGE_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectDatastorePage/complianceSuccess";
   public static final String ID_MIGRATION_CHANGEHOSTDATASTORE_ERROR_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL
               + "/migrationTypePage/changeHostDatastoreErrors/errWarningText";
   public static final String ID_LABEL_MIGRATION_CHANGE_HOST_ERR_WARNING =
         "changeHostErrors/errWarningText";
   public static final String ID_MIGRATION_DELETED_VM_ERROR_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/migrationTypePage/errWarningText";
   public static final String ID_MIGRATION_SELECTHOSTPAGE_ERROR_TEXT_GENERAL_LABEL =
         ID_MIGRATION_WIZARD_PANEL
               + "/selectHostPage/compatibilityView/vm1_label3_error";
   public static final String ID_MIGRATION_SELECTHOSTPAGE_VM_ERROR_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectHostPage/compatibilityView/vm" + "vmIndex"
               + "_label" + "labelIndex" + "_vm";
   public static final String ID_MIGRATION_SELECTHOSTPAGE_HOST_ERROR_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectHostPage/compatibilityView/vm" + "vmIndex"
               + "_label" + "labelIndex" + "_host";
   public static final String ID_MIGRATION_SELECTHOSTPAGE_ERROR_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectHostPage/compatibilityView/vm" + "vmIndex"
               + "_label" + "labelIndex" + "_error";
   public static final String ID_MIGRATION_SELECTHOSTPAGE_WARNING_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectHostPage/compatibilityView/vm" + "vmIndex"
               + "_label" + "labelIndex" + "_warning";
   public static final String ID_MIGRATION_SELECTRPPAGE_VM_ERROR_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectResourcePoolPage/compatibilityView/vm"
               + "vmIndex" + "_label" + "labelIndex" + "_vm";
   public static final String ID_MIGRATION_SELECTRPPAGE_HOST_ERROR_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectResourcePoolPage/compatibilityView/vm"
               + "vmIndex" + "_label" + "labelIndex" + "_host";
   public static final String ID_MIGRATION_SELECTRPPAGE_ERROR_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectResourcePoolPage/compatibilityView/vm"
               + "vmIndex" + "_label" + "labelIndex" + "_error";
   public static final String ID_MIGRATION_SELECTRPPAGE_WARNING_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectResourcePoolPage/compatibilityView/vm"
               + "vmIndex" + "_label" + "labelIndex" + "_warning";
   public static final String ID_MIGRATION_SELECRPPAGE_HOST_NOT_CONNECTED_ERROR =
         ID_MIGRATION_WIZARD_PANEL
               + "/selectResourcePoolPage/compatibilityView/error_resourcePoolPage_hostNotConnected";
   public static final String ID_MIGRATION_SELECTDATASTOREPAGE_VM_ERROR_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectDatastorePage/compatibilityView/vm"
               + "vmIndex" + "_label" + "labelIndex" + "_vm";
   public static final String ID_MIGRATION_SELECTDATASTOREPAGE_HOST_ERROR_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectDatastorePage/compatibilityView/vm"
               + "vmIndex" + "_label" + "labelIndex" + "_host";
   public static final String ID_MIGRATION_SELECTDATASTOREPAGE_ERROR_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectDatastorePage/compatibilityView/vm"
               + "vmIndex" + "_label" + "labelIndex" + "_error";
   public static final String ID_MIGRATION_SELECTDATASTOREPAGE_WARNING_TEXT_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/selectDatastorePage/compatibilityView/vm"
               + "vmIndex" + "_label" + "labelIndex" + "_warning";
   public static final String ID_MIGRATION_NO_PERMISSION_ON_DATASTORE_ERROR_TEXT_LABEL =
         "error_datastorePage_allocateSpacePermission";
   public static final String ID_MIGRATION_DEST_HOST_LABEL = ID_MIGRATION_WIZARD_PANEL
         + "/host";
   public static final String ID_MIGRATION_VMS_TO_MIGRATE_LABEL =
         ID_MIGRATION_WIZARD_PANEL + "/vms";
   public static final String ID_MIGRATION_CHANGE_HOST_RADIO = ID_MIGRATION_WIZARD_PANEL
         + "/radioBtnChangeHost";
   public static final String ID_MIGRATION_CHANGE_DATASTORE_RADIO =
         ID_MIGRATION_WIZARD_PANEL + "/radioBtnChangeDatastore";
   public static final String ID_MIGRATION_CHANGE_HOST_DATASTORE_RADIO =
         ID_MIGRATION_WIZARD_PANEL + "/radioBtnChangeHostAndDatastore";
   public static final String ID_MIGRATION_VMOTION_PRIORITY_HIGH_RADIO =
         ID_MIGRATION_WIZARD_PANEL + "/radioBtnHighPriority";
   public static final String ID_MIGRATION_VMOTION_PRIORITY_LOW_RADIO =
         ID_MIGRATION_WIZARD_PANEL + "/radioBtnLowPriority";
   public static final String ID_ALLOW_HOST_SELECTION_CHECKBOX =
         ID_MIGRATION_WIZARD_PANEL + "/allowHostSelectionChkBox";
   public static final String ID_MIGRATION_SPINNER = ID_MIGRATION_WIZARD_PANEL
         + "/progressBar";
   public static final String ID_MIGRATION_DISK_FORMAT_SAME = ID_MIGRATION_WIZARD_PANEL
         + "/selectDiskFormatPage/radioBtnSameFormat";
   public static final String ID_MIGRATION_DISK_FORMAT_THIN = ID_MIGRATION_WIZARD_PANEL
         + "/selectDiskFormatPage/radioBtnThinFormat";
   public static final String ID_MIGRATION_DISK_FORMAT_THICK = ID_MIGRATION_WIZARD_PANEL
         + "/selectDiskFormatPage/radioBtnThickFormat";
   public static final String ID_MIGRATION_SELECT_DATASTORES_LIST =
         ID_MIGRATION_WIZARD_PANEL + "/storageList";
   public static final String ID_MIGRATION_SELECT_HOST_LIST = ID_MIGRATION_WIZARD_PANEL
         + "/hostList";
   public static final String ID_MIGRATION_PAGE_NAVIGATOR_MIGRATION_TYPE_STEP =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "step_migrationTypePage";
   public static final String ID_MIGRATION_PAGE_NAVIGATOR_SELECT_RP_STEP =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "step_selectResourcePoolPage";
   public static final String ID_MIGRATION_PAGE_NAVIGATOR_SELECT_HOST_STEP =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "step_selectHostPage";
   public static final String ID_MIGRATION_PAGE_NAVIGATOR_SELECT_DATASTORE_STEP =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "step_selectDatastorePage";
   public static final String ID_MIGRATION_PAGE_NAVIGATOR_SELECT_DISK_FORMAT_STEP =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "step_selectDiskFormatPage";
   public static final String ID_MIGRATION_PAGE_NAVIGATOR_SELECT_VMOTION_PRIORITY_STEP =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "step_selectMovePriorityPage";
   public static final String ID_MIGRATION_PAGE_NAVIGATOR_REVIEW_SELECTIONS_STEP =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "step_migrationSummaryPage";
   public static final String ID_MIGRATION_PAGE_NAVIGATOR_NEXT_BUTTON =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "/panoramaNextBtn";
   public static final String ID_MIGRATION_PAGE_NAVIGATOR_PREVIOUS_BUTTON =
         ID_LABEL_WIZARD_PAGE_NAVIGATOR_STEP_ID + "/panoramaPreviousBtn";
   public static final String ID_MIGRATION_MIGRATION_TYPE_PAGE = "migrationTypePage";
   public static final String ID_MIGRATION_SELECT_RP_PAGE = "selectResourcePoolPage";
   public static final String ID_MIGRATION_SELECT_HOST_PAGE = "selectHostPage";
   public static final String ID_MIGRATION_SELECT_DATASTORE_PAGE = "selectDatastorePage";
   public static final String ID_MIGRATION_SELECT_DISK_FORMAT_PAGE =
         "selectDiskFormatPage";
   public static final String ID_MIGRATION_SELECT_MOVE_PRIORITY_PAGE =
         "selectMovePriorityPage";
   public static final String ID_MIGRATION_SUMMARY_PAGE = "migrationSummaryPage";
   public static final String ID_MIGRATION_SELECT_DATASTORE_MODE_BUTTON = "modeBtn";
   public static final String ID_SELECT_DATASTORE_ADVANCED_VIEW_ADVANCED_DATAGRID =
         "disksList";
   public static final String ID_ADVANCED_DISK_SELECTOR_COMBOBOX = "uicEditor";
   public static final String ID_MIGRATION_DISK_SELECTOR_COMBOBOX =
         ID_MIGRATION_SELECT_DATASTORE_PAGE + "/diskFormatSelector";
   public static final String ID_MIGRATION_WARNING_LABEL = "vm1_label3_warning";
   public static final String ID_MIGRATION_ERROR_LABEL = "vm1_label3_error";
   public static final String ID_LABEL_MIGRATION_HOST_NOT_CONNECTED_ERROR =
         "error_resourcePoolPage_hostNotConnected";
   public static final String ID_LABEL_MIGRATION_HOST_NOT_FOUND_ERROR =
         "error_resourcePoolPage_noHostsFound";
   public static final String ID_LABEL_MIGRATION_SELCT_HOST_ERROR =
         "error_resourcePoolPage_selectRPOnly";
   public static final String ID_DATAGRID_CLUSTER_TAB_NAME =
         "ClusterComputeResource_filterGridView";
   public static final String ID_HARDWARE_VERSION_LABEL = "hardwareVersion";
   public static final String ID_MIGRATION_CHANGEHOSTDATASTORE_ERROR_WARN_LABEL =
         "changeHostDatastoreErrors/errWarningText";

   // General Wizard
   public static final String ID_WIZARD_TIWO_DIALOG = "tiwoDialog";
   public static final String ID_WIZARD_CONTAINER = "wizardContainer";
   public static final String ID_WIZARD_NEXT_BUTTON = "next";
   public static final String ID_WIZARD_PREVIOUS_BUTTON = "back";
   public static final String ID_WIZARD_FINISH_BUTTON = "finish";
   public static final String ID_WIZARD_CANCEL_BUTTON = "cancel";
   public static final String ID_WIZARD_OK_BUTTON = "ok";
   public static final String ID_WIZARD_LOADING_PROGRESSBAR = "progressBar";
   public static final String ID_WIZARD_VALIDATING_DATA_PROGRESSBAR =
         "animatedLoadingProgressBar";
   public static final String ID_WIZARD_SHOW_BUTTON = "showButton";
   public static final String ID_BUTTON_NEXT = "next";
   public static final String ID_VAPP_NEXT_BUTTON =
         "Create vApp/bottomBar/wizardButtonsContainer/rightButtonsContainer/next";
   public static final String ID_DIALOG_POP_UP = "dialogPopup";
   public static final String ID_SINGLE_PAGE_DIALOG_ID = "SinglePageDialog_ID";

   // General TIWO Dialog
   public static final String ID_TIWO_DIALOG = "tiwoDialog";
   public static final String ID_TIWO_DIALOG_OK_BUTTON = "okButton";
   public static final String ID_TIWO_DIALOG_CANCEL_BUTTON = "cancelButton";
   public static final String ID_RENAME_ENTITY_TEXT_INPUT = "inputLabel";
   public static final String ID_EXPORT_OVF_PROGRESS = "ovfProgressDialog0";
   public static final String ID_EXPORT_PROGRESSBAR = "progressBar";
   public static final String ID_TIWO_CONTENTS = "contents";
   public static final String ID_CLOSE_BUTTON_ON_POPUP_VALIDATION_MESSAGE =
         ID_TIWO_DIALOG + "/" + ID_TIWO_CONTENTS + "/" + "_closeButton";

   // Provision Datastore Dialog
   public static final String ID_PROVISION_DATASTORE_NAME =
         "datastoreNameAndLocationPage/inputLabel";
   public static final String ID_PROVISION_DATASTORE_VMFS_TYPE = "vmfsRadio";
   public static final String ID_PROVISION_DATASTORE_NFS_TYPE = "nfsRadio";
   public static final String ID_PROVISION_DATASTORE_DEVICE_SELECTION_LIST =
         "deviceList";
   public static final String ID_PROVISION_DATASTORE_VMFS_VERSION_TYPE_VMFS5 =
         "vmfs5Radio";
   public static final String ID_PROVISION_DATASTORE_VMFS_VERSION_TYPE_VMFS3 =
         "vmfs3Radio";
   public static final String ID_PROVISION_DATASTORE_PARTITION_CONFIGURATION_LIST =
         "partitionConfigList";
   public static final String ID_PROVISION_DATASTORE_NFS_SERVER = "nfsServer";
   public static final String ID_PROVISION_DATASTORE_NFS_FOLDER = "nfsFolder";
   public static final String ID_VMFS_DATASTORE_PROPERTY_NAME = "vmfsDsName";
   public static final String ID_NFS_DATASTORE_PROPERTY_NAME = "dsName";
   public static final String ID_LIST_HOST_ACCESSIBILITY = "hosts";
   public static final String ID_PROVISION_DATASTORE_ASSIGN_NEW_SIGNATURE_RADIO =
         "newSignatureRadio";
   public static final String ID_PROVISION_DATASTORE_KEEP_EXISTING_SIGNATURE_RADIO =
         "keepSignatureRadio";
   public static final String ID_PROVISION_DATASTORE_FORMAT_DISK_RADIO =
         "formatDiskRadio";
   public static final String ID_PROVISION_DATASTORE_SELECT_HOST_DROPDOWN_LIST =
         "hostSelector";
   public static final String ID_LABEL_DATASTORE_NAME_AND_LOCATION_PAGE =
         "step_datastoreNameAndLocationPage";
   public static final String ID_LABEL_DATASTORE_TYPE_PAGE = "step_datastoreTypePage";
   public static final String ID_LABEL_VMFS_DEVICE_SELECTION_PAGE =
         "step_vmfsSelectDevicePage";
   public static final String ID_LABEL_VMFS_VERSION_PAGE = "step_vmfsSelectVersionPage";
   public static final String ID_LABEL_PARTITION_CONFIGURATION_PAGE =
         "step_vmfsPartitionConfigurationPage";
   public static final String ID_LABEL_READY_TO_COMPLETE_PAGE =
         "step_vmfsReadyToCompletePage";
   public static final String ID_LABEL_NFS_CONFIGURATION_PAGE = "step_nfsSettingsPage";
   public static final String ID_LABEL_NFS_READY_TO_COMPLETE_PAGE =
         "step_nfsReadyToCompletePage";
   public static final String ID_LIST_MAXIMUM_FILE_SIZE = "maxFileSizeList";
   public static final String ID_LABEL_DATASTORE_NAME_TEXT = "txtDatastoreName";
   public static final String ID_LABEL_DATASTORE_TYPE_TEXT = "txtDatastoreType";
   public static final String ID_LABEL_DATASTORE_PARTITION_FORMAT_TEXT =
         "txtPartitionFormat";
   public static final String ID_LABEL_DATASTORE_VMFS_VERSION_TEXT = "txtVmfsVersion";
   public static final String ID_LABEL_DATASTORE_MAX_FILE_SIZE_TEXT = "txtMaxFileSize";
   public static final String ID_LABEL_DATASTORE_BLOCK_SIZE_TEXT = "txtBlockSize";
   public static final String ID_PAGE_DATASTORE_NAME_AND_LOCATION_PAGE =
         "datastoreNameAndLocationPage";
   public static final String ID_TEXT_DATASTORE_NAME_AND_LOCATION_SEARCH_INPUT =
         "datastoreNameAndLocationPage/searchInput";
   public static final String ID_LABEL_NEW_DATASTORE_SEARCH_TEXT = "Label_2";
   public static final String ID_LIST_DATASTORES_FOR_STORAGE_FOLDER =
         "datastoresForStorageFolder/list";
   public static final String ID_BUTTON_SEARCH_DATASTORE_NAME_AND_LOCATION =
         "datastoreNameAndLocationPage/searchButton";
   public static final String ID_PAGE_DATASTORE_TYPE_PAGE = "datastoreTypePage";
   public static final String ID_PAGE_NFS_CONFIGURATION_PAGE = "nfsSettingsPage";

   // NFS related
   public static final String ID_TEXT_UNMOUNT_NFS = "selectHostsTextArea";
   public static final String ID_BUTTON_DATASTORE_GETTING_STARTED_LINK_1 =
         "vsphere.core.datastore.gettingStarted/gettingStartedHelpLink_1";
   public static final String ID_BUTTON_DATASTORE_GETTING_STARTED_LINK_2 =
         "vsphere.core.datastore.gettingStarted/gettingStartedHelpLink_2";
   public static final String ID_BUTTON_DATASTORE_GETTING_STARTED_LINK_3 =
         "vsphere.core.datastore.gettingStarted/gettingStartedHelpLink_3";
   public static final String ID_TEXT_GENERAL_PROPERTIES_NAME =
         "_DatastorePropertiesView_PropertyGridRow1";
   public static final String ID_TEXT_GENERAL_PROPERTIES_TYPE =
         "_DatastorePropertiesView_PropertyGridRow2";
   public static final String ID_LABEL_GENERAL_CAPACITY =
         "titleLabel.capacityStackBlock";
   public static final String ID_TEXT_GENERAL_CAPACITY_TOTAL_CAPACITY =
         "_DatastoreCapacityView_PropertyGridRow2";
   public static final String ID_TEXT_GENERAL_CAPACITY_PROVISIONED_SPACE =
         "_DatastoreCapacityView_PropertyGridRow3";
   public static final String ID_TEXT_GENERAL_CAPACITY_FREE_SPACE =
         "_DatastoreCapacityView_PropertyGridRow4";
   public static final String ID_LABEL_GENERAL_STORAGE_IO_CONTROL =
         "titleLabel.storageIORMStackBlock";
   public static final String ID_TEXT_GENERAL_STORAGE_IO_CONTROL_STATUS =
         "_DatastoreCapabilitiesView_PropertyGridRow3";
   public static final String ID_TEXT_GENERAL_STORAGE_IO_CONTROL_MODE =
         "_DatastoreCapabilitiesView_PropertyGridRow4";
   public static final String ID_TEXT_GENERAL_STORAGE_IO_CONTROL_DRS_IOM =
         "_DatastoreCapabilitiesView_PropertyGridRow6";
   public static final String ID_TEXT_DEVICE_BACKING_SERVER =
         "_NfsDatastoreBackingView_PropertyGridRow1";
   public static final String ID_TEXT_DEVICE_BACKING_FOLDER =
         "_NfsDatastoreBackingView_PropertyGridRow2";
   public static final String ID_GRID_HOSTS_FOR_DATASTORE = "hostsForDatastore/list";
   public static final String ID_CHECK_BOX_NFS_READ_ONLY = "nfsMountReadOnly";
   public static final String ID_BUTTON_SELECT_ALL_HOST = "btnSelectAll";
   public static final String ID_LABEL_DS_READY_TO_COMPLETE_NAME = "txtDatastoreName";
   public static final String ID_LABEL_DS_READY_TO_COMPLETE_TYPE = "txtDatastoreType";
   public static final String ID_LABEL_DS_READY_TO_COMPLETE_NFS_SERVER = "txtNfsServer";
   public static final String ID_LABEL_DS_READY_TO_COMPLETE_NFS_FOLDER = "txtNfsFolder";
   public static final String ID_LABEL_DS_READY_TO_COMPLETE_NFS_ACCESS_MODE =
         "txtNfsAccessMode";
   public static final String ID_LABEL_DS_READY_TO_COMPLETE_NFS_GENERAL =
         "_NfsReadyToCompletePage_Label1";
   public static final String ID_LABEL_DS_READY_TO_COMPLETE_NFS_SHARE =
         "_NfsReadyToCompletePage_Label2";
   public static final String ID_LABEL_DS_READY_TO_COMPLETE_NFS_HOST_ACCESS =
         "_NfsReadyToCompletePage_Label3";
   public static final String ID_LIST_DS_READY_TO_COMPLETE_NFS_HOST_ACCESS =
         "nfsReadyToCompletePage/selectedHosts";
   public static final String ID_GRID_DATASTORES_FOR_HOST =
         "datastoresForStandaloneHost/list";
   public static final String ID_GRID_DATASTORES_LIST = "Datastore/list";

   // Confirm Remove Datastore Dialog
   public static final String ID_REMOVE_DATASTORE_SELECT_ALL_BUTTON = "selectAllButton";

   public static final String ID_BUTTON_RECOMMENDATIONS_APPLY = "buttonApply";
   public static final String ID_DIALOG_MIGRATE_RECOMMENDATIONS = "dsSdrsRecsDialog";

   public static final String ID_LABEL_RP_RESTORE = "_RestoreRpTreeView_Label1";

   public static final String ID_GRID_FILES_DATAGRID =
         "vsphere.core.datastore.manage.filesView/datagrid";

   // Configure Storage I/O dialog
   public static final String ID_STORAGE_IO_CONTROL_STATUS = "storageIOControlStatus";
   public static final String ID_CONFIGURE_STORAGE_IO_CONTROL_EDIT_BUTTON =
         "btn_vsphere.core.datastore.configureStorageIOControl";
   public static final String ID_ENABLE_STORAGE_IO_CONTROL_CHECKBOX =
         "enableIORMCheckBox";
   public static final String ID_CONFIGURE_STORAGE_IO_CONTROL_OK_BUTTON = "okButton";
   public static final String ID_BUTTON_REFRESH_DATASTORE =
         "btn_vsphere.core.datastore.refreshAction";
   public static final String ID_BUTTON_UPGRADE_VMFS_DATASTORE =
         "btn_vsphere.core.datastore.upgradeVmfsDatastore";
   public static final String ID_LABEL_CAPACITY_TITLE =
         "vsphere.core.datastore.manage.settings.general.capacityView/titleLabel";
   public static final String ID_LABEL_PROPERTIES_TITLE =
         "vsphere.core.datastore.manage.settings.general.propertiesView/titleLabel";
   public static final String ID_LABEL_DATASTORE_CAPABILITIES_TITLE =
         "vsphere.core.datastore.manage.settings.general.capabilitiesView/titleLabel";
   public static final String ID_LABEL_SELECT_DEVICE_PAGE = "step_selectDevicePage";
   public static final String ID_LABEL_SPECIFY_CONFIGURATION = "step_configurationPage";
   public static final String ID_LABEL_READY_TO_COMPLETE = "step_readyToCompletePage";
   public static final String ID_CHECKBOX_DISABLE_STORAGE_IO_STATISTICS_COLLECTION =
         "disableSIOStatCollCheckBox";
   public static final String ID_CHECKBOX_EXCUDE_STORAGE_IO_STATISTICS_FROM_SDRS =
         "excludeIOStatCollCheckBox";
   public static final String ID_HBOX_CONGESTION_THRESHOLD_SECTION =
         "_EditDatastoreCapabilitiesForm_HGroup1";
   public static final String ID_LABEL_HEADER_TEXT = "headerText";
   public static final String ID_DATASTORE_DEVICE_SELECTION_LIST =
         "datastoresForDatacenter/list";
   public static final String ID_LABEL_DEVICE_INFO = "txtDeviceInfo";
   public static final String ID_LABEL_READY_TO_COMPLETE_PAGE_GENERAL_SECTION =
         "_ReadyToCompletePage_Label1";
   public static final String ID_LABEL_READY_TO_COMPLETE_PAGE_DEVICE_AND_FORMATTING_SECTION =
         "_ReadyToCompletePage_Label2";
   public static final String ID_RADIO_CONGESTION_THRESHOLD_MANUAL =
         "manualThresholdRadioButton";
   public static final String ID_RADIO_CONGESTION_THRESHOLD_PERCENTAGE_OF_PEAK_THROUGHPUT =
         "percentPeakThroughputRadioButton";
   public static final String ID_TEXT_PERCENTAGE_OF_PEAK_THROUGHPUT =
         "percentPeakThroughput";
   public static final String ID_LABEL_DEVICE_SELECTION_PAGE =
         "step_deviceSelectionPage";
   public static final String ID_LABEL_PARTITION_TYPE_PAGE = "step_partitionTypePage";
   public static final String ID_LABEL_PARTITION_TYPE_TEXT = "txtPartitionType";
   public static final String ID_LABEL_DEVICE_TYPE_TEXT = "txtDisk";
   public static final String ID_LABEL_LUN_TEXT = "txtLun";
   public static final String ID_LABEL_CAPACITY_TEXT = "txtCapacity";

   // Datastore Manage sub tab
   public static final String ID_DATASTORE_MANAGE_SETTINGS_PAGE_LIST =
         "vsphere.core.datastore.manage.settings/tocTree";

   // Increase Datastore Capacity
   public static final String ID_BUTTON_INCREASE_DATASTORE_CAPACITY =
         "btn_vsphere.core.datastore.increaseAction";
   public static final String ID_TEXTINPUT_INCREASE_SIZE_BY = "textDisplay";
   public static final String ID_IMAGE_CAPACITY = "arrowImagecapacityStackBlock";
   public static final String ID_LABEL_TOTAL_CAPACITY = "totalSize";
   public static final String ID_IMAGE_STORAGE_IO_CONTROL =
         "arrowImagestorageIORMStackBlock";
   public static final String ID_IMAGE_FILE_SYSTEM_BLOCK =
         "arrowImagefileSystemStackBlock";
   public static final String ID_LABEL_FILE_SYSTEM_TYPE = "vmfsDsFileSystemType";
   public static final String ID_LABEL_MAX_FILE_SIZE = "vmfsDsMaxFileSize";
   public static final String ID_LABEL_BLOCK_SIZE = "vmfsDsBlockSize";
   public static final String ID_LABEL_DRIVE_TYPE = "vmfsDsDriveType";
   public static final String ID_LABEL_PROVISIONED_SPACE = "provisionedSpace";
   public static final String ID_LABEL_FREE_SPACE = "freeSpace";
   public static final String ID_LABEL_STORAGE_IO_CONTROL_MODE = "storageIOControlMode";
   public static final String ID_LABEL_SDRS_IO_METRICS = "statsSdrsAggregation";
   public static final String ID_IMAGE_HARDWARE_ACCELERATION_BLOCK =
         "arrowImagehardwareAccelerationStackBlock";
   public static final String ID_LIST_CONNECTIVITY_AND_MULTIPATHING_HOST_LIST =
         "vsphere.core.datastore.manage.settings.connectivityView/hostList";
   public static final String ID_BUTTON_EDIT_MULTIPATHING_POLICIES =
         "editMultipathingPolicyAction";
   public static final String ID_LABEL_PATH_SELECTION_POLICY = "pathSelectionPolicy";
   public static final String ID_COMBOBOX_PATH_SELECTION_POLICY =
         "tiwoDialog/policyCombo";
   public static final String ID_LIST_PATH_SELECTION_LIST = "tiwoDialog/pathsList";
   public static final String ID_TEXT_INCREASE_DS_SIZE_BY = "increaseSizeStepper";
   public static final String ID_BUTTON_DATASTORE_REFRESH =
         "vsphere.core.datastore.explorer.refreshDatastoreAction/button";
   public static final String ID_LABEL_DEVICE_BACKING_TITLE =
         "vsphere.core.datastore.manage.settings.backing.extentListView/titleLabel";
   public static final String ID_LABEL_DEVICE_DETAILS_TITLE =
         "vsphere.core.datastore.manage.settings.backing.extentDetailsView/titleLabel";
   public static final String ID_LIST_DEVICE_BACKING =
         "vsphere.core.datastore.manage.settings.backing.extentListView/_ExtentsView_ExtentsListView1";
   public static final String ID_LIST_PRIMARY_PARTITION =
         "vsphere.core.datastore.manage.settings.backing.extentDetailsView/primaryPartitions";
   public static final String ID_LIST_LOGICAL_PARTITION =
         "vsphere.core.datastore.manage.settings.backing.extentDetailsView/logicalPartitions";
   public static final String ID_IMAGE_PATH_SELECTION_POLICY_BLOCK =
         "arrowImage_MultipathingPoliciesView_StackBlock1";

   // VM Provisioning Wizard
   public static final String ID_CONTAINER_VM_PROV_PAGES = "vmConfigPages";
   public static final String ID_BUTTON_VMPROVISION_PAGE_SEARCH =
         "searchControl/searchButton";
   public static final String ID_VMPROVISIONING_SELECT_RP_PAGE = "selectResourcePage";
   public static final String ID_TEXT_VMPROVISION_RPPAGE_COMPLIANCE_SUCCESS =
         ID_VMPROVISIONING_SELECT_RP_PAGE + "/complianceSuccess";
   public static final String ID_TEXT_VMPROVISION_SEARCH_RESULT_TIMESTAMP =
         "resultRetrievalTimestamp";
   public static final String ID_VMPROVISIONING_CREATION_TYPE_PAGE =
         "provisioningTypePage";
   public static final String ID_VMPROVISIONING_SELECT_TEMPLATE_PAGE =
         "selectTemplatePage";
   public static final String ID_VMPROVISIONING_SELECT_SOURCE_OBJECT_PAGE =
         "selectSourceObjectPage";
   public static final String ID_VMPROVISIONING_SELECT_NAME_FOLDER_PAGE =
         "selectNameLocationPage";
   public static final String ID_VMPROVISIONING_SELECT_DATASTORE_PAGE =
         "selectStoragePage";
   public static final String ID_VMPROVISIONING_SELECT_VM_VERSION_PAGE =
         "selectVmVersionPage";
   public static final String ID_VMPROVISIONING_SELECT_GUEST_OS_PAGE =
         "selectGuestOsPage";
   public static final String ID_VMPROVISIONING_CUSTOMIZE_GUEST_OS_PAGE =
         "customizeGuestOsPage";
   public static final String ID_VMPROVISIONING_DEPLOY_SELECT_NAME_FOLDER_PAGE =
         "selectNameFolderPage";
   public static final String ID_VMPROVISIONING_CUSTOMIZE_HARDWARE_PAGE =
         "customizeHardwarePage";
   public static final String ID_VMPROVISIONING_SUMMARY_PAGE = "summaryPage";
   public static final String ID_VMPROVISIONING_TYPE_LIST = "optionsList";
   public static final String ID_VMPROVISIONING_SEARCH_TEMPLATE_TEXT_INPUT =
         ID_VMPROVISIONING_SELECT_TEMPLATE_PAGE + "/searchInput";
   public static final String ID_VMPROVISIONING_SEARCH_BUTTON =
         ID_VMPROVISIONING_SELECT_TEMPLATE_PAGE + "/searchButton";
   public static final String ID_VMPROVISIONING_TEMPLATE_LIST_ADV_DATAGRID =
         "vsphere.core.search.searchResultsListView";
   public static final String ID_VMPROVISIONING_VM_NAME_TEXT_INPUT =
         ID_VMPROVISIONING_SELECT_NAME_FOLDER_PAGE + "/nameTextInput";
   public static final String ID_ADVDATAGRID_VMPROVISIONING_FILTER_DATACENTERS =
         "Datacenter_filterGridView";
   public static final String ID_ADVDATAGRID_CLUSTERSFILTER =
         "ClusterComputeResource_filterGridView";
   public static final String ID_ADVDATAGRID_VMPROVISIONING_FILTER_FOLDERS =
         "Folder_filterGridView";
   public static final String ID_ADVDATAGRID_VMPROVISIONING_SELECTED_OBJECTS_DATACENTERS =
         "Datacenter_selectedObjectGridView";
   public static final String ID_ADVDATAGRID_VMPROVISIONING_SELECTED_OBJECTS_FOLDERS =
         "Folder_selectedObjectGridView";
   public static final String ID_ADVDATAGRID_VMPROVISIONING_FILTER_CLUSTERS =
         "ClusterComputeResource_filterGridView";
   public static final String ID_ADVDATAGRID_VMPROVISIONING_FILTER_HOST =
         "HostSystem_filterGridView";
   public static final String ID_ADVDATAGRID_VMPROVISIONING_FILTER_RESOURCEPOOL =
         "ResourcePool_filterGridView";
   public static final String ID_ADVDATAGRID_VMPROVISIONING_FILTER_VAPP =
         "VirtualApp_filterGridView";
   public static final String ID_ADVDATAGRID_VMPROVISIONING_SELECTED_OBJECTS_CLUSTERS =
         "ClusterComputeResource_selectedObjectGridView";
   public static final String ID_ADVDATAGRID_VMPROVISIONING_SELECTED_OBJECTS_HOST =
         "HostSystem_selectedObjectGridView";
   public static final String ID_ADVDATAGRID_VMPROVISIONING_SELECTED_OBJECTS_RESOURCEPOOL =
         "ResourcePool_selectedObjectGridView";
   public static final String ID_ADVDATAGRID_VMPROVISIONING_SELECTED_OBJECTS_VAPP =
         "VirtualApp_selectedObjectGridView";
   public static final String ID_BUTTONBAR_VMPROVISIONING_FILTER_VIEW =
         "filterViewToggleButtonBar";
   public static final String ID_BUTTONBAR_VMPROVISIONING_SELECTED_OBJECTS_VIEW =
         "selectedObjectsViewToggleButtonBar";
   public static final String ID_TABNAVIGATOR_VMPROVISIONING = "_mainTabNavigator";
   public static final String ID_CHECKBOX_DATACENTERS =
         "_CheckBoxColumnRenderer_CheckBox1";
   public static final String ID_VMPROVISIONING_DATASTORE_LIST_ADV_DATAGRID =
         ID_VMPROVISIONING_SELECT_DATASTORE_PAGE + "/datastoreList";
   public static final String ID_VMPROVISIONING_DATASTORE_ADV_DATAGRID =
         ID_VMPROVISIONING_SELECT_DATASTORE_PAGE + "/storageList";
   public static final String ID_VMPROVISIONING_GOS_FAMILY_COMBO =
         ID_VMPROVISIONING_SELECT_GUEST_OS_PAGE + "/gosFamilyCombo";
   public static final String ID_VMPROVISIONING_GOS_VERSION_COMBO =
         ID_VMPROVISIONING_SELECT_GUEST_OS_PAGE + "/gosVersionCombo";
   public static final String ID_VMPROVISIONING_GOS_VERSION_TEXT_AREA =
         ID_VMPROVISIONING_SELECT_GUEST_OS_PAGE + "/alternateName";
   public static final String ID_VMPROVISIONING_CUSTOMIZE_GOS_CHECKBOX =
         "customizeGOSCheckBox";
   public static final String ID_VMPROVISIONING_CUSTOMIZE_HW_CHECKBOX =
         "customizeHardwareCheckBox";
   public static final String ID_VMPROVISIONING_NUMBER_VMS_TO_CREATE_NUM_STEPPER =
         ID_VMPROVISIONING_SELECT_TEMPLATE_PAGE + "/deployVmCountNumericStepper";
   public static final String ID_VMPROVISIONING_HDD = "disk_";
   public static final String ID_TEXT_BOX_VMPROVISIONING_HDD_SIZE = "textDisplay";
   public static final String ID_VMPROVISIONING_HDD_THIN_PROVISION_CHECK_BOX = "thin";
   public static final String ID_VMPROVISIONING_PROV_TYPE_VALUE_SUMMARY_PAGE =
         "provisioningTypeValue";
   public static final String ID_VMPROVISIONING_DESTINATION_SUMMARY_PAGE = "hostValue";
   public static final String ID_VMPROVISIONING_VMNAME_VALUE_SUMMARY_PAGE = "vmValue";
   public static final String ID_VMPROVISIONING_TEMPLATE_VALUE_SUMMARY_PAGE =
         "sourceVmValue";
   public static final String ID_VMPROVISIONING_RP_VALUE_SUMMARY_PAGE = "rpValue";
   public static final String ID_VMPROVISIONING_DATASTORE_VALUE_SUMMARY_PAGE =
         "datastoreValue";
   public static final String ID_VMPROVISIONING_FOLDER_VALUE_SUMMARY_PAGE =
         "folderValue";
   public static final String ID_VMPROVISIONING_STORAGE_VALUE_SUMMARY_PAGE =
         "diskFormatValue";
   public static final String ID_VMPROVISIONING_STORAGE_DATASTORE_CONF_SUMMARY_PAGE =
         "editedDisk_0_Value";
   public static final String ID_VMPROVISIONING_POWER_ON_CHECKBOX = "powerOnVM";
   public static final String ID_VMPROVISIONING_SELECT_RP_ONLY_ERROR_TEXT_LABEL =
         "error_resourcePoolPage_selectRPOnly";
   public static final String ID_VMPROVISIONING_NON_DRS_CLUSTER_ERROR_TEXT_LABEL =
         "error_resourcePoolPage_nonDrsCluster";
   public static final String ID_VMPROVISIONING_NO_HOSTS_CLUSTER_ERROR_TEXT_LABEL =
         "error_resourcePoolPage_noHostsFound";
   public static final String ID_VMPROVISIONING_HOST_NOT_CONNECTED_OR_IN_MM_ERROR_TEXT_LABEL =
         "error_resourcePoolPage_hostNotConnected";
   public static final String ID_LABEL_VMPROVISIONING_CLUSTEREDHOST_NOT_CONNECTED_OR_IN_MM_ERROR =
         "error_resourcePoolPage_noConnectedHostsFound";
   public static final String ID_VMPROVISIONING_NO_PERMISSION_ON_RP_ERROR_TEXT_LABEL =
         "error_resourcePoolPage_noPermissionOnRP0";
   public static final String ID_VMPROVISIONING_NO_PERMISSION_ON_DATASTORE_ERROR_TEXT_LABEL =
         "error_datastoreProvisioningPage_noPermissionOnDatastore";
   public static final String ID_VMPROVISIONING_VM_VERSION_SEVEN = "radButtonVmVersion";
   public static final String ID_VMPROVISIONING_DISK_SIZE_COMBO_BOX = "scale";
   public static final String ID_COMBO_BOX_HARD_DISK_UNIT = "disk_?/diskSize/scale";
   public static final String ID_VMPROVISIONING_NEW_HARD_DISK_LABEL =
         "titleLabel.disk_0";
   public static final String ID_VMPROVISIONING_NEW_HARD_DISK_LABEL_GENERIC =
         "titleLabel.disk_?";
   public static final String ID_VMPROVISIONING_CAPACITY_STEPPER = "capacity";
   public static final String ID_VMPROVISIONING_VALIDATION_MESSAGE = "validationMessage";
   public static final String ID_VMPROVISIONING_CONTAINER = "_container";
   public static final String ID_VMPROVISIONING_MESSAGE = "_message";
   public static final String ID_VMPROVISIONING_MESSAGE_TEXT = "messageText";
   public static final String ID_VMPROVISIONING_LABEL_CREATEVM_ERROR =
         ID_VMPROVISIONING_SELECT_NAME_FOLDER_PAGE + "/issueDescription";
   public static final String ID_LABEL_VALIDATION_ERRORS_MESSAGE = "_errorLabel";
   public static final String ID_LABEL_CLOSE_ERRORS_MESSAGE = "_closeButton";
   public static final String ID_VALIDATION_ERRORS_CLOSE_BUTTON = "_closeButton";
   public static final String ID_VALIDATION_ERRORS_PANE = "_errorDataGroup";
   public static final String ID_LABEL_DATA_STATE = "data.pageInfo.state";
   public static final String ID_LABEL_DATA_TITLE = "data.pageInfo.title";
   public static final String ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR =
         "tiwoDialog/wizardPageNavigator/dataGroup";
   public static final String ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR + "/step.id=";

   public static final String ID_VMPROVISIONING_SELECT_CREATION_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "step_creationTypePage";
   public static final String ID_VMPROVISIONING_SELECT_A_VM_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "step_selectTemplatePage";
   public static final String ID_VMPROVISIONING_EDIT_SETTINGS_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "rootSettingsPage";
   public static final String ID_VMPROVISIONING_SELECT_NAME_FOLDER_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "step_selectNameFolderPage";
   public static final String ID_VMPROVISIONING_SELECT_A_RESOURCE_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX
               + "step_selectResourcePoolPage";
   public static final String ID_VMPROVISIONING_SELECT_A_STORAGE_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "step_selectDatastorePage";
   public static final String ID_VMPROVISIONING_SELECT_VM_VERSION_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "step_selectVmVersionPage";
   public static final String ID_VMPROVISIONING_SELECT_GUEST_OS_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "step_selectGuestOsPage";
   public static final String ID_VMPROVISIONING_CUSTOMIZE_HARDWARE_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX
               + "step_customizeHardwarePage";
   public static final String ID_VMPROVISIONING_REVIEW_SETTINGS_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "step_summaryPage";

   public static final String ID_VMPROVISIONING_SELECT_TEMPLATE_INVENTORY_TREE =
         ID_WIZARD_TIWO_DIALOG + "/" + ID_VMPROVISIONING_SELECT_TEMPLATE_PAGE + "/"
               + ID_WIZARD_SHOW_BUTTON;
   public static final String ID_VMOPTIONS_VMTOOLS_RESTART_COMBO = "cmbRestart";
   public static final String ID_VMOPTIONS_VMTOOLS_SHUTDOWN_COMBO = "cmbShutdown";
   public static final String ID_VMOPTIONS_VMTOOLS_SUSPEND_COMBO = "cmbSuspend";
   public static final String ID_VMPROVISIONING_DRS_DISABLED_DATASTORE_ADV_GRID =
         "datastorePerPodBox/datastoreList";
   public static final String ID_VMPROVISIONING_CLUSTERED_DATASTORE_ADV_DATAGRID =
         new StringBuffer(ID_VMPROVISIONING_SELECT_DATASTORE_PAGE).append(
               ID_VMPROVISIONING_DRS_DISABLED_DATASTORE_ADV_GRID).toString();
   public static final String ID_VMPROVISIONING_DATASTORE_CLUSTER_GRID =
         ID_VMPROVISIONING_SELECT_DATASTORE_PAGE + "/storageBox/storageList";
   public static final String ID_SELECT_NAME_FOLDER_TREE = new StringBuffer(
         ID_WIZARD_TIWO_DIALOG).append("/")
         .append(ID_VMPROVISIONING_SELECT_NAME_FOLDER_PAGE).append("/")
         .append(ID_NAV_TREE).toString();
   public static final String ID_SELECT_RESOURCE_TREE = new StringBuffer(
         ID_WIZARD_TIWO_DIALOG).append("/").append(ID_VMPROVISIONING_SELECT_RP_PAGE)
         .append("/").append(ID_NAV_TREE).toString();
   public static final String ID_VM_VERSION_COMBO = "vmVersionSelector/comboBox";
   public static final String ID_VMPROVISIONING_DISKPROV_EAGER_ZEROED_RADIO =
         "eageryScurb";
   public static final String ID_VMPROVISIONING_DISKPROV_THIN_RADIO = "thin";
   public static final String ID_VMPROVISIONING_DISKPROV_FLAT_RADIO = "flat";
   public static final String ID_LINK_CPU_ADVANCED = "cpuPage/cpuIdMask/cpuIdAdvanced";
   public static final String ID_DIALOG_DRS_RECOMMENDATION = "vmSdrsRecsDialog";
   public static final String ID_TAB_DRS_RECOMMENDATION =
         "vmSdrsRecsDialog/tabs/automationName=Recommendations";
   public static final String ID_TAB_DRS_FAULTS =
         "vmSdrsRecsDialog/tabs/automationName=Faults";
   // Constant of VM -> Manage -> Settings tab
   public static final String ID_VM_MANAGE_SETTINGS_WRAPPER = "manage.settings.wrapper";
   public static final String ID_VM_MANAGE_SETTINGS_VMHARDWARE_LABEL =
         "vsphere.core.vm.manage.settings.wrapper/tocTree/automationName=VM Hardware";
   public static final String ID_VM_MANAGE_SETTINGS_VMHARDWARE_VIEW =
         "vsphere.core.vm.manage.settings.vmHardwareView";
   public static final String ID_VM_MANAGE_SETTINGS_VMOPTIONS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_VM)
         .append(ID_VM_MANAGE_SETTINGS_WRAPPER).append("/").append("text=VM Options")
         .toString();
   public static final String ID_VMPROVISIONING_MEMORY_RESERVATION_CHECK =
         "memoryPage/reservation/reserveAll";
   public static final String ID_VIRTUAL_DEVICE_NODE_NEW_HDD = " New Hard disk";
   public static final String ID_VIRTUAL_DEVICE_NODE_NEW_CD_DVD = " New CD/DVD Drive";
   public static final String ID_LABEL_VMOPTIONS_ADV_ARROW_IMAGE = "arrowImageadvanced";
   public static final String ID_LABEL_VMOPTIONS_SWAP_FILE_LOC = "swapfileText";
   public static String ID_LABEL_VMPROVISIONING_GOS_SUMMARY_PAGE = "gosNameValue";
   public static String ID_LABEL_VMPROVISIONING_NICCOUNT_SUMMARY_PAGE = "nicCountValue";
   public static String ID_LABEL_VMPROVISIONING_HDDTYPE_SUMMARY_PAGE = "newDisk_0_Value";
   public static String ID_LABEL_VMPROVISIONING_NEW_HDD = "newDiskLocation_?_Value";
   public static String ID_FORM_ALERT = "className=Alert";

   // Migrate VM Wizard
   public static String ID_RADIO_BUTTON_CHANGE_HOST = "radioBtnChangeHost";
   public static String ID_RADIO_BUTTON_CHANGE_DATASTORE = "radioBtnChangeDatastore";
   public static String ID_RADIO_BUTTON_CHANGE_HOST_AND_DATASTORE =
         "radioBtnChangeHostAndDatastore";

   //Getting started Image
   public static final String ID_GETTING_STARTED_IMAGE = "_GettingStartedView_Image2";

   // Recent Tasks Portlet
   public static final String ID_RECENT_TASKS_PORTLET = EXTENSION_PREFIX
         + "tasks.recentTasksView.portletChrome";
   public static final String ID_TASKS_LIST = ID_RECENT_TASKS_PORTLET
         + "/recentTasksList";
   public static final String ID_RECENT_TASKS_VIEW_COMBO_BOX = ID_RECENT_TASKS_PORTLET
         + "/filterComboBox";
   public static final String ID_RECENT_TASKS_STATE_BUTTON_BAR = ID_RECENT_TASKS_PORTLET
         + "/stateButtonBar";
   public static final String ID_RECENT_TASKS_ALL_BUTTON =
         ID_RECENT_TASKS_STATE_BUTTON_BAR + "/automationName=All";
   public static final String ID_RECENT_TASKS_RUNNING_BUTTON =
         ID_RECENT_TASKS_STATE_BUTTON_BAR + "/automationName=Running";
   public static final String ID_RECENT_TASKS_FAILED_BUTTON =
         ID_RECENT_TASKS_STATE_BUTTON_BAR + "/automationName=Failed";
   public static final String ID_RECENT_TASKS_GOTOTASKCONSOLE_LINK =
         "taskConsoleLinkButton";
   public static final String ID_TASKS_GRID = "recentTasksGridView";

   // Task Console Related
   public static final String ID_TASK_CONSOLE_TASK_GRID = "taskGrid";
   public static final String ID_TASK_CONSOLE_DETAIL_TASKNAME_LABEL =
         "compositeDescriptionBox";
   public static final String ID_TASK_CONSOLE_DETAIL_TASKSTATUS_LABEL = "statusData";
   public static final String ID_TASK_CONSOLE_DETAIL_TASKINITATIOR_LABEL =
         "initiatorData";
   public static final String ID_TASK_CONSOLE_DETAIL_TASKTARGET_BUTTON = "entityName";
   public static final String ID_TASK_CONSOLE_DETAIL_RELATEDEVENTS_LIST =
         "relatedEventDataGrid";
   public static final String ID_TASK_CONSOLE_DETAIL_VCENTERSERVER_LABEL = "vCenterData";
   public static final String ID_MORE_ACTIONS_DATAGRID = "rhs";
   public static final String ID_TASK_NOTIFICATION_WINDOW =
         "className=com.vmware.ui.views.notification.NoticePopupView";
   public static final String ID_TASK_NOTIFICATION_TASK_NAME_LABEL =
         "vsphere.core.operationNotificationView/taskNameBox/taskName";
   public static final String ID_TASK_CONSOLE_DETAIL_FT_POPUP = "description";
   public static final String ID_TASK_CONSOLE_DETAILS_COLUMN_CONTENTS =
         "dataProvider.list.source.0.details";
   public static final String ID_TASK_CONSOLE_DETAILS_LINK_BUTTON = "detailsLink_7";
   public static final String ID_TASK_CONSOLE_ERROR_STACK = "stackLabel0";
   public static final String ID_TASK_CONSOLE_SUBMIT_ERROR_LINK =
         "submitErrorReportLink";
   public static final String ID_EVENT_RELATED_TARGET_LINK =
         "eventRelatedTargetLinkButton";

   // VM Question Related
   public static final String ID_VM_QUESTION_DESCRIPTION_LABEL = "questionDescription";
   public static final String ID_VM_QUESTION_ANSWEROK_BUTTON = "answerQuestion";

   // Datacenter Add related
   public static final String ID_MORE_ACTIONS_BUTTON = "actionButton";
   public static final String ID_NEW_DATACENTER_BUTTON =
         "vsphere.core.datacenter.createAction";
   public static final String ID_NEW_FOLDER_BUTTON = "vsphere.core.folder.createAction";
   public static final String ID_NEW_VMFOLDER_BUTTON =
         "vsphere.core.datacenter.createVmFolderAction";
   public static final String ID_REMOVE_DATACENTER =
         "vsphere.core.datacenter.removeAction";
   public static final String ID_DC_NAME = "SinglePageDialog_ID/inputLabel";
   public static final String ID_DC_OBJECT_ADVANCE_DATAGRID = "Datacenter/list";

   // Datacenter > Related Objects tab > Virtual Machines subtab
   public static final String ID_DC_VM_MORE_ACTIONS_BUTTON = "vmsForDatacenter" + "/"
         + ID_CONTAINER_DATAGRID_TOOL + "/" + ID_MORE_ACTIONS_BUTTON;

   // Edit Vm Constants
   public static final String ID_ADDING_DISK_LABEL = "";
   public static final String ID_EDIT_ADD_DEVICE_COMBO = "addHardware/popUpButton";
   public static final String ID_EDIT_ADD_DEVICE_MENU = "menu";
   public static final String ID_HARD_DRIVE_UTILIZATION_GENERIC =
         "com.vmware.vsphere.client.vmui.portlet.hardware/viewDisk_?/diskUtilization/_VmDiskView_Text1";
   public static final String ID_EDIT_HW_STACK = "hardwareStack";
   public static final String ID_EDIT_ADD_DEVICE_BUTTON = "addDeviceButton";
   public static final String ID_EDIT_VM_DEVICE_CANNOT_ADD_MESSAGE_LABEL = "message";
   public static final String ID_EDIT_SETTING_VPMC_CHECK = "cpuPage/counters";
   public static final String ID_LABEL_EDIT_SETTING_VPMC =
         "cpuPage/countersRow/className=PropertyGridLabel";

   // Add Dependency constants
   public static final String ID_DATAGRID_VSERVICES = "dataGrid";
   public static final String ID_BUTTON_REMOVE_DEPENDENCY =
         "removeVServiceDependencyButton";
   public static final String ID_BUTTON_ADD_DEPENDENCY = "addVServiceDependencyButton";
   public static final String ID_DROPDOWN_DEPENDENCY = "providerNameSelect";
   public static final String ID_TEXTINPUT_DEPENDENCY = "nameTextInput";
   public static final String ID_BUTTON_EDIT_DEPENDENCY = "editVServiceDependencyButton";
   public static final String ID_CHECK_BOX_REQUIRED = "dependencyRequiredCheckBox";
   public static final String ID_COMBO_BOX_DEPENDENCY_LIST = "bindProviderInput";
   public static final String ID_DEPENDENCY_NAME = "nameLabel";
   public static final String ID_DEPENDENCY_PROVIDER = "providerLabel";
   public static final String ID_DEPENDENCY_REQUIRED = "requiredLabel";
   public static final String ID_DEPENDENCY_TYPE = "typeLabel";
   public static final String ID_RADIO_BINDING_YES = "requiredYes";
   public static final String ID_RADIO_REQUIRED_YES = "bindingYes";
   public static final String ID_TEXTINPUT_EDIT_DEPENDENCY_DESCRIPTION =
         "dependencyDescriptionArea";
   public static final String ID_TEXTINPUT_EDIT_DEPENDENCY_NAME = "dependencyNameInput";
   public static final String ID_RADIO_LOCAL_FILE = "localRadioButton";
   public static final String ID_LABEL_VM_NAME = "plainLabel";
   public static final String ID_LABEL_PROVIDER_VALUE = "providerValue";
   public static final String ID_DROPDOWN_PROVIDER = "providerDropDown";

   // Added New Device. Forget not to replace the question mark with the actual
   // no. of device added.
   // so if you want to verify whether the second device is added successfully,
   // replace the ? with 2 and so on.
   public static final String ID_NEW_HARD_DISK_CAPACITY_STEPPER = "disk_?/capacity";
   public static final String ID_NEW_HARD_DISK_SHARES_COMBO_BOX = "disk_?/sharesLevel";
   public static final String ID_NEW_HARD_DISK_CUSTOM_SHARE_TEXT_BOX =
         "disk_?/customShares";
   public static final String ID_NEW_HARD_DISK_THIN_PROVISION_CHECK_BOX = "disk_?/thin";
   public static final String ID_NEW_HARD_DISK_FLAT_PROVISION_CHECK_BOX =
         "disk_?/eageryScurb";
   public static final String ID_RADIO_NEW_HARD_DISK_THIN_PROVISION = "disk_?/thin";
   public static final String ID_RADIO_NEW_HARD_DISK_THICK_EAGER_ZEROED_PROVISION =
         "disk_?/eageryScurb";
   public static final String ID_RADIO_NEW_HARD_DISK_THICK_LAZY_ZEROED_PROVISION =
         "disk_?/flat";
   public static final String ID_NEW_HARD_DISK_DATASTORE_COMBO_BOX = "disk_?/location";
   public static final String ID_HARD_DRIVE_DRIVE_TYPE = "disk_?/diskTypeLabel";
   public static final String ID_NEW_HARD_DISK_SELECT_DATASTORE =
         "automationName=Select a datastore cluster or datastore";
   public static final String ID_NEW_HARD_DISK_DATASTORE_LOCATION = "storageList";
   public static final String ID_NEW_HARD_DISK_ACTUALDATASTORE_LOCATION =
         "datastoreList";
   public static final String ID_NEW_HARD_DISK_RECOMMENDATION_LOCATION = "recList";
   public static final String ID_NEW_HARD_DISK_SCSI_TYPE_1 = "SCSI(1:1) SCSI device 1";
   public static final String ID_SDRS_CHECKBOX = "sdrsChkBox";
   public static final String ID_NEW_HARD_DISK_DATASTORE_LOCATION_OK_BUTTON = "buttonOk";
   public static final String ID_NEW_HARD_DISK_DATASTORE_LOCATION_CANCEL_BUTTON =
         "buttonCancel";
   public static final String ID_NEW_HARD_DISK_SCSI_NODE_COMBO_BOX = "disk_?/nodeSCSI";
   public static final String ID_NEW_HARD_DISK_SCSIDEVICE_NODE_COMBO_BOX = "nodeSCSI";
   public static final String ID_NEW_HARD_DISK_DISK_MODE_DEPENDENT_RADIO_BUTTON =
         "disk_?/radioButtonDependent";
   public static final String ID_NEW_HARD_DISK_SCSI_RADIO = "disk_?/scsi";
   public static final String ID_NEW_HARD_DISK_MODE = "disk_?/diskMode";
   public static final String ID_NEW_HARD_DISK_IDE_RADIO = "disk_?/ide";
   public static final String ID_NEW_HARD_DISK_IDEDEVICE_NODE_COMBO_BOX = "nodeIDE";
   public static final String ID_RDM_HARD_DISK_SHARES_COMBO_BOX = "sharesLevel";
   public static final String SELECT_TARGET_LUN_CANCEL_BUTTON = "btnCancel";
   public static final String ID_LABEL_HARD_DISK_DEVICE_WILL_BE_REMOVED =
         "_DiskPage_Label1";
   public static final String ID_SCSI_DEVICE_FINISHPAGE = "scsiCtrl_?_Value";
   public static final String ID_ALERT_FORM_OCCUPIED_NODE_MESSAGE =
         "/automationName=This bus node is already occupied by another device";

   // TO DO: ID_NEW_HARD_DISK_DISK_MODE_INDEPENDENT_PERSISTENT_RADIO_BUTTON may
   // be removed in future
   public static final String ID_NEW_HARD_DISK_DISK_MODE_INDEPENDENT_PERSISTENT_RADIO_BUTTON =
         "disk_?/radioButtonIP";

   // Radio button for RDM disk
   public static final String ID_HARD_DISK_DISK_MODE_INDEPENDENT_PERSISTENT_RADIO_BUTTON =
         "radioButtonIP";
   public static final String ID_NEW_HARD_DISK_DISK_MODE_INDEPENDENT_NON_PERSISTENT_RADIO_BUTTON =
         "disk_?/radioButtonIN";
   public static final String ID_NEW_SCSI_DEVICE_LABEL = "titleLabel.New SCSI device";
   public static final String ID_SCSI_DEVICE_LABEL = "titleLabel.SCSI device ?";
   public static final String ID_SCSI_CONTROLLER_TYPE_COMBO_BOX = "scsi_?/typeSetCombo";
   public static final String ID_SCSI_CONTROLLER_SHARE_TYPE_COMBO_BOX =
         "scsi_?/busSharingCombo";
   public static final String ID_DATASTORE_BROWSER_OK_BUTTON = "dsBrowserDlg/buttonOk";
   public static final String ID_SCSI_CONTROLLER_CHANGE_BUTTON = "changeTypeButton";
   public static final String ID_PARALLEL_OUTPUT_FILE_TEXT_BOX_GENERIC =
         "parallelPort_?/file";
   public static final String ID_DATASTORE_BROWSER_CANCEL_BUTTON =
         "dsBrowserDlg/buttonCancel";
   public static final String ID_HARD_DRIVE_LABEL_GENERIC = "titleLabel.viewDisk_?";
   public static final String ID_HARD_DRIVE_LABEL_EDIT_VM_GENERIC = "titleLabel.disk_?";
   public static final String ID_HARD_DRIVE_CAPACITY_TEXT_BOX =
         "disk_?/editing/diskSize/text";
   public static final String ID_CD_DVD_LABEL_GENERIC = "titleLabel.cdrom_?";
   public static final String ID_CD_DVD_ISO_FILE_NAME_TEXT_INPUT = "cdrom_?/file";
   public static final String ID_FLOPPY_ISO_FILE_NAME_TEXT_INPUT = "fileNameInput";
   public static final String ID_CD_DVD_DATASTORE_BROWSE_BUTTON = "cdrom_?/browse";
   public static final String ID_FLOPPY_DATASTORE_BROWSE_BUTTON = "floppy_?/browse";
   public static final String ID_CD_DVD_VIRTUAL_DEVICE_NODE_COMBO_BOX = "cdrom_?/node";
   public static final String ID_CD_DVD_HOST_DEVICE_FILE_NAME_COMBO_BOX =
         "cdrom_?/devices";
   public static final String ID_CD_DVD_STATUS_GENERIC_CHECK_BOX = "cdrom_?/connected";
   public static final String ID_CD_DVD_START_CONNECTED_GENERIC_CHECK_BOX =
         "cdrom_?/startConnected";
   public static final String ID_FLOPPY_STATUS_GENERIC_CHECK_BOX = "floppy_?/connected";
   public static final String ID_FLOPPY_FILE_NAME_TEXT_INPUT = "floppy_?/file";
   public static final String ID_FLOPPY_DEVICE_BROWSE_BUTTON = "floppy_?/browse";
   public static final String ID_HDD_CONTROL_PREFIX = "disk_?/";
   public static final String ID_SCSI_CONTROLLER_PREFIX = "scsi_?/";
   public static final String ID_TEXT_INPUT_VM_PROV_CPU_RESERVATION =
         "cpuPage/reservation/className=TextInput";
   public static final String ID_TEXT_INPUT_VM_PROV_CPU_LIMIT =
         "cpuPage/limit/className=TextInput";
   public static final String ID_TEXT_INPUT_DISK_IOLIMIT = "disk_?/ioLimit/textDisplay";
   public static final String ID_COMBO_NEW_SCSI_DEVICE_NODE =
         "New SCSI device/nodeCombo";

   // This constant will be removed later
   // public static final String ID_FLOPPY_START_CONNECTED_GENERIC_CHECK_BOX =
   // "floppy_?/startConnected";
   public static final String ID_FLOPPY_START_CONNECTED_GENERIC_CHECK_BOX =
         "floppy_?/startConnected";
   public static final String ID_FLOPPY_DEVICE_DRIVE_GENERIC = "floppy_?/drives";
   public static final String ID_PARALLEL_PORT_TYPE_GENERIC_COMBO_BOX =
         "parallelPort_?/type";
   public static final String ID_PARALLEL_PORT_DEVICE_GENERIC_COMBO_BOX =
         "parallelPort_?/device";
   public static final String ID_SERIAL_PORT_DEVICE_GENERIC_COMBO_BOX =
         "serialPort_?/device";
   public static final String ID_SERIAL_OUTPUT_FILE_TEXT_BOX_GENERIC =
         "serialPort_?/file";
   public static final String ID_SERIAL_PORT_TYPE_GENERIC_COMBO_BOX =
         "serialPort_?/type";
   public static final String ID_SERIAL_PORT_PHYSICAL_DEVICE_GENERIC =
         "serialPort_?/device";
   public static final String ID_SERIAL_PORT_YIELD_CPU_GENERIC =
         "serialPort_?/bodySection/_SerialPortPage_PropertyGridRow4/yield";
   public static final String ID_SERIAL_PORT_OUTPUTFILE_GENERIC_TEXTBOX =
         "serialPort_?/bodySection/_SerialPortPage_PropertyGridRow3/backing/file";
   public static final String ID_SERIAL_PORT_URI_DIRECTION_GENERIC =
         "serialPort_?/uridirection";
   public static final String ID_SERIAL_PORT_URI_GENERIC = "serialPort_?/portURI";
   public static final String ID_SERIAL_PORT_USE_VSPC_GENERIC = "serialPort_?/useVSPC";
   public static final String ID_SERIAL_PORT_VSPC_GENERIC = "serialPort_?/vspcURI";
   public static final String ID_SERIAL_PORT_BROWSE_BUTTON_GENERIC =
         "serialPort_?/browse";
   public static final String ID_PARALLEL_PORT_BROWSE_BUTTON_GENERIC =
         "parallelPort_?/browse";
   public static final String ID_SERIAL_PORT_PIPE_NAME_TEXT_BOX_GENERIC =
         "serialPort_?/pipeName";
   public static final String ID_SERIAL_PORT_NEAR_END_COMBO_BOX_GENERIC =
         "serialPort_?/endpoint";
   public static final String ID_SERIAL_PORT_FAR_END_COMBO_BOX_GENERIC =
         "serialPort_?/noRxLoss";
   public static final String ID_NIC_STATUS_GENERIC_CHECK_BOX =
         "ethernetCard_?/connected";
   public static final String ID_MAC_ADDRESS_COMBO_BOX_GENERIC =
         "ethernetCard_?/addressType";
   public static final String ID_MAC_ADDRESS_TEXT_BOX_GENERIC =
         "ethernetCard_?/macAddress";
   public static final String ID_NIC_START_CONNECTED_GENERIC_CHECK_BOX =
         "ethernetCard_?/startConnected";
   public static final String ID_PARALLEL_PORT_STATUS_GENERIC_CHECK_BOX =
         "parallelPort_?/connected";
   public static final String ID_PARALLEL_PORT_START_CONNECTED_GENERIC_CHECK_BOX =
         "parallelPort_?/startConnected";
   public static final String ID_SERIAL_PORT_STATUS_GENERIC_CHECK_BOX =
         "serialPort_?/connected";
   public static final String ID_SERIAL_PORT_START_CONNECTED_GENERIC_CHECK_BOX =
         "serialPort_?/startConnected";
   public static final String ID_SCSI_STATUS_GENERIC_CHECK_BOX = "scsi_?/connected";
   public static final String ID_USB_CONTROLLER_STATUS_CHECK_BOX =
         "viewStack/USB controller /connected";
   public static final String ID_USB_CONTROLLER_LABEL = "titleLabel.USB controller ";
   public static final String ID_USB_DEVICE_LABEL = "titleLabel.usb_?";
   public static final String ID_USB_CONTROLLER_TYPE_COMBO_BOX =
         "New USB Controller/usbType";
   public static final String ID_USB_DEVICE_GRIDLABEL = "titleLabel.view_UsbDevices";
   public static final String ID_USB_TYPE_2 = "USB 2.0";
   public static final String ID_USB_TYPE_3 = "USB 3.0";
   public static final String ID_USB_DEVICE = "USB Devices";
   public static final String ID_CD_DVD_DEVICE_TYPE_GENERIC = "cdrom_?/type";
   public static final String ID_FLOPPY_DEVICE_TYPE_GENERIC = "floppy_?/type";
   public static final String ID_NETWORK_TYPE_GENERIC_COMBO_BOX = "ethernetCard_?/type";
   public static final String ID_NETWORK_CONNECTION_GENERIC_COMBO_BOX =
         "ethernetCard_?/connection";
   public static final String ID_CD_DVD_DEVICE_MODE_GENERIC = "cdrom_?/mode";
   public static final String ID_CD_DVD_DEVICE_MEDIA = "cdrom_?/devices";
   public static final String ID_FLOPPY_LABEL_GENERIC = "titleLabel.floppy_?";
   public static final String ID_FLOPPY_LABEL_STATUS_GENERIC = "Floppy_?";
   public static final String ID_NIC_LABEL_GENERIC = "titleLabel.ethernetCard_?";
   public static final String ID_NIC_MAC_VALUE_GENERIC =
         "com.vmware.vsphere.client.vmui.portlet.hardware/viewNIC_?/macAddress";
   public static final String ID_PARALLEL_PORT_LABEL_GENERIC =
         "titleLabel.parallelPort_?";
   public static final String ID_SERIAL_PORT_LABEL_GENERIC = "titleLabel.serialPort_?";
   public static final String ID_SERIAL_PORT_MODE_GENERIC =
         "serialPort_?/headerPropertyGrid/type";
   public static final String ID_HDD_REMOVE_BUTTON_GENERIC =
         "disk_?/headerPropertyGrid/remove";
   public static final String ID_HDD_DELETE_FILE_FROM_DS_GENERIC =
         "disk_?/headerPropertyGrid/purge";
   public static final String ID_CD_DVD_REMOVE_BUTTON_GENERIC = "cdrom_?/remove";
   public static final String ID_FLOPPY_REMOVE_BUTTON_GENERIC = "floppy_?/remove";
   public static final String ID_NIC_REMOVE_BUTTON_GENERIC = "ethernetCard_?/remove";
   public static final String ID_PARALLEL_PORT_REMOVE_BUTTON_GENERIC =
         "parallelPort_?/remove";
   public static final String ID_PARALLEL_PORT_DEVICE_NOT_REMOVED_LABEL =
         "parallelPort_?/_ParallelPortPage_Label1";
   public static final String ID_PARALLEL_PORT_RESTORE_BUTTON_GENERIC =
         "parallelPort_?/remove";
   public static final String ID_SERIAL_PORT_REMOVE_BUTTON_GENERIC =
         "serialPort_?/remove";
   public static final String ID_SERIAL_PORT_DEVICE_NOT_REMOVED_LABEL =
         "serialPort_?/_SerialPortPage_Label1";
   public static final String ID_SCSI_REMOVE_BUTTON_GENERIC = "scsi_?/remove";
   public static final String ID_SCSI_DEVICE_REMOVE_BUTTON_GENERIC =
         "SCSI device ?/remove";
   public static final String ID_SCSI_DEVICE = "SCSI device ?";
   public static final String ID_SCSI_DEVICE_VIRTUAL_NODE = "nodeCombo";
   public static final String ID_USB_CONTROLLER_REMOVE_BUTTON =
         "viewStack/USB controller /remove";
   public static final String ID_USB_XHCI_CONTROLLER_REMOVE_BUTTON =
         "hardwareStack/USB xHCI controller /remove";
   public static final String ID_USB_DEVICE_REMOVE_BUTTON =
         "usb_?/headerPropertyGrid/remove";
   public static final String ID_LABEL_USB_CONTROLLER_UNSUPPORTED =
         "_USBControllerPage_Text1";
   public static final String ID_LABEL_USB_DEVICE_WILL_BE_REMOVED_WARNING =
         "_USBControllerPage_Label1";
   public static final String ID_LABEL_SCSI_DEVICE_WILL_BE_REMOVED_WARNING =
         "_ScsiDevicePage_Label1";
   public static final String ADD_DEVICE_ERROR_OK_BUTTON = "label=OK";
   public static final String ID_DATASTORE_BROWSER_DIALOG = "dsBrowserDlg";
   public static final String ID_DATASTORE_BROWSER_TREE = "dsBrowserDlg/tree";
   public static final String ID_EDIT_HW_MAX_MEMORY_LABEL = "memoryControl/label1";
   public static final String ID_EDIT_HW_MAX_MEMORY_VALUE_LABEL = "memoryControl/label2";
   public static final String ID_EDIT_HW_RESERVATION_MEMORY_LABEL =
         "memoryPage/reservationControl/bottomRow/label1";
   public static final String ID_EDIT_HW_RESERVATION_MEMORY_VALUE_LABEL =
         "memoryPage/reservationControl/bottomRow/label2";
   public static final String ID_EDIT_HW_LIMIT_MEMORY_LABEL =
         "memoryPage/limitControl/bottomRow/label1";
   public static final String ID_EDIT_HW_LIMIT_MEMORY_VALUE_LABEL =
         "memoryPage/limitControl/bottomRow/label2";
   public static final String ID_EDIT_HW_RESERVATION_LABEL =
         "reservation/labelReservationTitle";
   public static final String ID_EDIT_HW_CANCEL_BUTTON = "cancelButton";
   public static final String ID_EDIT_MEMORY_NO_AFFINITY = "radNoAffinity";
   public static final String ID_EDIT_MEMORY_USE_AFFINITY = "radAffinity";
   public static final String ID_EDIT_MEMORY_AFFINITY_CHECKBOX =
         "memoryAffinityContainer/_VmMemoryPage_HBox1/chkNodes1[?]";
   public static final String ID_LABEL_DEVICE_USB = "usb_?";
   public static final String ID_COMBO_BOX_DEVICE = "deviceCombo";
   public static final String ID_LABEL_CONNECTION_STATUS = "connection";
   public static final String ID_LABEL_USB_UNIQUE_ID = "uuid";
   public static final String ID_LABEL_USB_DEVICE_PRESENCE = "presence";
   public static final String ID_LABEL_USB_VMOTION_SUPPORT = "vmotion";
   public static final String ID_LABEL_USB_DEVICE_WILL_BE_REMOVED =
         "_USBDevicePage_Text1";
   public static final String ID_LABEL_SELECT_HOST_DEVICE_DETAILS = "details";
   public static final String ID_TREE_VIEW_VM_SETTINGS = "tocTree";
   public static final String ID_BUTTON_EDIT_IN_MANAGE_SETTINGS_TAB =
         "vsphere.core.vm.manage.settings.vmOptionsView/btn_vsphere.core.vm.provisioning.editAction";

   // EDIT VM CPU related ID constants
   public static final String ID_EDIT_HW_CPU_LABEL = "titleLabel.cpuPage";
   public static final String ID_EDIT_HW_CPU_COMBO =
         "viewStack/hardwareStack/cpuPage/cpuCombo";
   public static final String ID_LABEL_EDIT_HW_SOCKETS = "sockets";
   public static final String ID_BUTTON_CPU_HELP = "coresPerSocketHelp";
   public static final String ID_TEXT_HELP = "displayText";
   public static final String ID_LABEL_CPU_VALUE = "cpuValue";
   public static final String ID_LABEL_VM_HARDWARE_CPU = "cpuMeterLabel";
   public static final String ID_EDIT_HW_CPU_RESERVATION_COMBO =
         "cpuPage/reservationControl/topRow/comboBox";
   public static final String ID_EDIT_HW_CPU_CORES_COMBO =
         "viewStack/hardwareStack/cpuPage/coresPerSocketCombo";
   public static final String ID_EDIT_HW_CPU_HOT_REMOVE_CHECK_BOX =
         "cpuPage/cpuHotPlugGroup/cpuHotRemove";
   public static final String ID_EDIT_HW_CPU_LIMIT_COMBO =
         "cpuPage/limitControl/topRow/comboBox";
   public static final String ID_EDIT_HW_CPU_SHARES_COMBO =
         "cpuPage/sharesControl/levels";
   public static final String ID_EDIT_HW_CPU_SHARES_STEPPER =
         "cpuPage/sharesControl/numShares";
   public static final String ID_EDIT_HW_CPU_HT_SHARE_COMBO =
         "cpuPage/htSharingControl/htSharingCombo";
   public static final String ID_EDIT_HW_CPU_MAX_RESERVATION_TEXT_LABEL =
         "reservationControl/bottomRow/label1";
   public static final String ID_EDIT_HW_CPU_MAX_RESERVATION_VALUE_LABEL =
         "reservationControl/bottomRow/label2";
   public static final String ID_EDIT_HW_MEMORY_MAX_RESERVATION_VALUE_LABEL =
         "memoryControl/bottomRow/label2";
   public static final String ID_EDIT_HW_MEMORY_COMBOBOX =
         "memoryControl/topRow/comboBox";
   public static final String ID_EDIT_HW_MEMORY_UNITS_COMBOBOX =
         "memoryControl/topRow/units";
   public static final String ID_EDIT_HW_CPU_MAX_LIMIT_TEXT_LABEL =
         "limitControl/bottomRow/label1";
   public static final String ID_EDIT_HW_CPU_MAX_LIMIT_VALUE_LABEL =
         "limitControl/bottomRow/label2";
   public static final String ID_EDIT_HW_CPU_AFFINITY_TEXT_BOX =
         "cpuPage/affinityControl/txtAffinity";
   public static final String ID_EDIT_HW_CPU_HOT_ADD_CHECK_BOX =
         "cpuPage/cpuHotPlugGroup/cpuHotAdd";
   public static final String ID_EDIT_HW_CPU_COMBO_COLLAPSED = "cpuCombo";
   public static final String ID_EDIT_VIDEO_CARD_LABEL = "titleLabel.Video card ";
   public static final String ID_EDIT_ENABLE_3D_SUPPORT_CHECKBOX = "enable3D";
   public static final String ID_EDIT_USB_CONTROLLER_LABEL =
         "titleLabel.New USB Controller";
   public static final String ID_EDIT_USB_SETTINGS_COMBOBOX = "usbType";
   public static final String ID_EDIT_ADDED_USB_XHCI_CONTROLLER_LABEL =
         "titleLabel.USB xHCI controller ";
   public static final String ID_EDIT_ADDED_USB_CONTROLLER_LABEL =
         "titleLabel.USB controller ";
   public static final String ID_EDIT_USB_TYPE_LABEL = "summaryLabel";
   public static final String ID_LABEL_VM_HARDWARE_MEMORY = "memMeterLabel";
   public static final String ID_LABEL_HYPERTHREADING_STATUS = "htStatus";
   public static final String ID_LABEL_AVAILABLE_CPUS = "cpus";
   public static final String ID_NESTED_HYPERVISOR_CHECK_BOX = "nestedHV";
   // MHz/GHz selection combo box

   public static final String ID_EDIT_HW_CPU_RESERVATION_UNIT =
         "cpuPage/reservationControl/topRow/units";
   public static final String ID_EDIT_HW_CPU_LIMIT_UNIT =
         "cpuPage/limitControl/topRow/units";
   public static final String ID_EDIT_HW_MEMORY_RESERVATION_UNIT =
         "memoryPage/reservationControl/topRow/units";
   public static final String ID_EDIT_HW_MEMORY_LIMIT_UNIT =
         "memoryPage/limitControl/topRow/units";

   // Edit VM HARD-DISk related ID constants
   public static final String ID_COMBO_BOX_EDIT_HW_HARD_DISK_LIMIT_IOPS = "ioLimit";
   public static final String ID_TEXT_INPUT_EDIT_HW_HARD_DISK_LIMIT_IOPS =
         "automationName=textInput";

   // Edit VM - vApp Options constants
   public static final String ID_EDIT_VM_VAPP_OPTIONS_BUTTON = "vAppOptionsStackButton";
   public static final String ID_EDIT_VM_VAPP_OPTIONS_PRODUCT_PAGE =
         "authoringPanel/productPage";
   public static final String ID_EDIT_VM_VAPP_OPTIONS_PRODUCT_LABEL =
         ID_EDIT_VM_VAPP_OPTIONS_PRODUCT_PAGE + "/titleLabel.productPage";
   public static final String ID_EDIT_VM_VAPP_ENABLE_CHECKBOX =
         "toggleVAppOptionsCheckBox";
   public static final String ID_EDIT_VM_VAPP_OPTIONS_PROPERTIES_LABEL =
         "titleLabel.vAppPropertiesPage";
   public static final String ID_EDIT_VM_VAPP_OPTIONS_PROPERTIES_LABEL1 =
         "titleLabel.propertyStackBlock0";
   public static final String ID_EDIT_VM_VAPP_IPALLOCATION_LABEL =
         "titleLabel.ipAllocationPage";
   public static final String ID_EDIT_VM_VAPP_OVFSECTIONS_LABEL =
         "titleLabel.unrecognizedOvfSectionsView";
   public static final String ID_EDIT_VM_VAPP_PRODUCT_LABEL = "titleLabel.productPage";
   public static final String ID_EDIT_VM_VAPP_AUTHPROPERTIES_LABEL =
         "titleLabel.advancedPropertyListView";
   public static final String ID_EDIT_VM_VAPP_AUTHIPALLOCATION_LABEL =
         "titleLabel.supportedIpAllocationView";
   public static final String ID_EDIT_VM_VAPP_AUTHOVFSETTINGS_LABEL =
         "titleLabel.ovfSettingsPage";
   public static final String ID_EDIT_VM_VAPP_VIEW_OVF_BUTTON =
         "viewOvfEnvironmentButton";
   public static final String ID_EDIT_VM_VAPP_VIEW_OVF_DIALOG = "ovfEnvironmentDialog";
   public static final String ID_EDIT_VM_VAPP_VIEW_OVG_CLOSE_BUTTON = "closeButton2";
   public static final String ID_EDIT_VM_VAPP_OVF_ISO_CHECK_OPTION = "isoCheckBox";
   public static final String ID_EDIT_VM_VAPP_OVF_TOOLS_CHECK = "toolsCheckBox";
   public static final String ID_EDIT_VM_VAPP_OVF_INSTALL_BOOT_CHECK =
         "installBootEnableCheckBox";
   public static final String ID_EDIT_VM_VAPP_OVF_INSTALL_BOOT_GROUP =
         "installBootDelayGroup";
   public static final String ID_VAPP_PROPERTIES_DESC_COMBO =
         "propertyStackBlock0EditPropertyControl0InputControl1";
   public static final String ID_VAPP_PROPERTIES_DESC_CHECKBOX =
         "propertyStackBlock0EditPropertyControl0InputControl1";
   public static final String ID_VAPP_PROPERTIES_DESC_TEXT_INPUT =
         "propertyStackBlock0EditPropertyControl0InputControl1";

   // Edit VM -> Virtual Hardware
   public static final String ID_COMBOBOX_EDIT_VM_NETWORK_ADAPTER = "connection";
   public static final String ID_LABEL_EDIT_VM_NETWORK_ADAPTER =
         "titleLabel.ethernetCard_0";
   public static final String ID_TEXTINPUT_EDIT_VM_NETWORK_ADAPTER_PORT_ID = "port";

   // Virtual machine Monitor tab
   public static final String ID_VM_MONITOR_TAB_ISSUES_TOCTREE = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY)
         .append("monitor.issues.wrapper/tocTree").toString();

   // Host AdvancedDataGrid
   public static final String ID_HOST_STORAGE_ADVDATAGRID = "datastoreList";
   public static final String ID_STORAGE_TAB_BUTTON = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("related.datastoresView").append(".button")
         .toString();

   // Host State
   public static final String ID_STATE_HOST_STATE = "state";

   // Permissions Portlet
   public static final String ID_PERMISSIONS_PORTLET = EXTENSION_PREFIX
         + "<ENTITY_TYPE>permissionView.chrome";
   public static final String ID_ADD_PERMISSION_BUTTON =
         "vsphere.core.permission.add/button";
   public static final String ID_ADD_USER_BUTTON_ON_ADD_PERMISSION_PORTLET = "addUser";
   public static final String ID_USERS_AND_GROUPS_ADVANCED_DATAGRID = "usersDataGrid"; // Users
   // and
   // Groups
   // advance
   // datagrid
   // in
   // Add
   // permission
   // portlet
   public static final String ID_ADD_USER_BUTTON_USERS_AND_GROUP_PORTLET =
         "addUserButton";
   public static final String ID_ROLE_LIST_ON_ADD_PERMISSION_PORTLET =
         "SinglePageDialog_ID/roles";
   public static final String ID_DOMAIN_ON_ADD_PERMISSION_PORTLET = "domainsComboBox";
   public static final String ID_CHECKBOX_PROPAGATE_PERMISSION = "propogate";
   public static final String ID_TRI_STATE_CHECKBOX_PRIVILEGE_TREE = "privGrid";
   public static final String ID_ADD_PERMISSION_DIALOG_VIEW_CHILDREN_LINK =
         "viewChildren";
   public static final String ID_ADD_PERMISSION_DIALOG_CHILDREN_VIEW = "childrenView";
   public static final String ID_ADD_PERMISSION_DIALOG_TEXT2_FOR_CURRENT_OBJ =
         "_EditPermissionListView_Text2";
   public static final String ID_ADD_PERMISSION_DIALOG_TEXT3_FOR_CURRENT_OBJ =
         "_EditPermissionListView_Text3";


   // Ports Tab
   public static final String ID_PORTS_ADV_DATA_GRID = "portList";
   public static final String ID_MONITORING_PORT_STATE_BUTTON = "monitoringButton";

   public static final String ID_PERMISSIONS_TABLE = new StringBuffer(
         ID_PERMISSIONS_PORTLET).append("/").append(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("permissionView").toString();

   // Networks Tab
   public static final String ID_VC_NETWORK_ADV_DATA_GRID = "dvpgForVCenter/list";
   public static final String ID_DC_NETWORK_ADV_DATA_GRID = "dvpgForDatacenter/list";
   public static final String ID_DC_NETWORK_ADV_DATA_GRID_UPLINKS =
         "uplinksForDatacenter/list";
   public static final String ID_DC_NETWORK_ADV_DATA_GRID_SWITCHES =
         "dvsForDatacenter/list";
   public static final String ID_DC_NETWORK_ADV_DATA_GRID_STANDARD_NETWORKS =
         "standardnetworksForDatacenter/list";
   public static final String ID_HOST_NETWORK_ADV_DATA_GRID =
         "stdnwksForStandaloneHost/list";
   public static final String ID_HOST_NETWORK_ADV_DATA_GRID_SWITCHES =
         "dvsForStandaloneHost/list";
   public static final String ID_HOST_NETWORK_ADV_DATA_GRID_UPLINKS =
         "uplinksForStandaloneHost/list";
   public static final String ID_CLUSTER_NETWORK_ADV_DATA_GRID =
         "stdnwksForCluster/list";
   public static final String ID_CLUSTER_NETWORK_ADV_DATA_GRID_SWITCHES =
         "dvsForCluster/list";
   public static final String ID_HOST_NETWORK_ADV_DATA_GRID_STDNETWORKS =
         "stdnwksForStandaloneHost/list";
   public static final String ID_DVSWITCH_ADV_DATA_GRID = "pgsForVMWDVS/list";
   public static final String ID_DVSWITCH_ADV_DATA_GRID_UPLINK = "uplinksForVMWDVS/list";
   public static final String ID_DVSWITCH_ADV_DATA_GRID_HOSTS = "hostsForVMWDVS/list";
   public static final String ID_DVSWITCH_ADV_DATA_GRID_VIRTUAL_MACHINES =
         "vmsForVMWDVS/list";
   public static final String ID_PORT_GROUP_ADV_DATA_GRID_HOSTS = "hostsForDVPG/list";
   public static final String ID_PORT_GROUP_ADV_DATA_GRID_VIRTUAL_MACHINES =
         "vmsForDVPG/list";
   public static final String ID_PORT_GROUP_ADV_DATA_GRID_OBJECTS =
         "DistributedVirtualPortgroup/list";
   public static final String ID_TASK_DVSWITCH_ADV_DATA_GRID =
         "vsphere.core.monitor.taskView";
   public static final String ID_DC_RELATED_ITEMS_NETWORK_GRID =
         "vsphere.core.datacenter.related.networksView";
   public static final String ID_VM_RELATED_ITEMS_NETWORK_GRID = "networksForVm/list";
   public static final String ID_VAPP_RELATED_ITEMS_NETWORK_GRID =
         "networksForVApp/list";

   // Virtual Machine Portlet
   public static final String ID_RELATED_OBJECTS_PORTLET = EXTENSION_PREFIX
         + "<ENTITY_TYPE>related";
   public static final String ID_VC_VMS_DATAGRID =
         "vsphere.core.folder.related/vmsForVCenter/list";
   public static final String ID_DATASTORE_VMS_DATAGRID =
         "vsphere.core.datastore.related/vmsForDatastore/list";
   public static final String ID_VIRTUAL_MACHINE_PORTLET = EXTENSION_PREFIX
         + "<ENTITY_TYPE>related.vmsView";
   public static final String ID_VIRTUAL_MACHINE_DATAGRID = new StringBuffer(
         ID_VIRTUAL_MACHINE_PORTLET).append("/").append("vmList").toString();

   public static final String ID_RESOURCEPOOL_VIRTUAL_MACHINE_DATAGRID =
         new StringBuffer(ID_VIRTUAL_MACHINE_PORTLET).append("/").append("vmList")
               .toString();

   // Sub button
   public static final String ID_BUTTON_TASKS =
         "vsphere.core.folder.monitor.tasks.button";
   public static final String ID_BUTTON_VC_HOST =
         "vsphere.core.folder.related/hostsForVCenter.button";
   public static final String ID_BUTTON_DATASTORE_HOST =
         "vsphere.core.datastore.related/hostsForDatastore.button";
   public static final String ID_BUTTON_VC_VMS =
         "vsphere.core.folder.related/vmsForVCenter.button";
   public static final String ID_BUTTON_DATASTORE_VMS =
         "vsphere.core.datastore.related/vmsForDatastore.button";
   public static final String ID_BUTTON_VC_DATASTORE =
         "vsphere.core.folder.related/datastoresForVCenter.button";
   public static final String ID_BUTTON_DATASTORE_TASKS =
         "vsphere.core.datastore.monitor.tasks.button";

   // VMOption dialog box
   public static final String ID_AFTER_POWERON = "afterPowerOnCB";
   public static final String ID_AFTER_RESUME = "afterResumeCB";
   public static final String ID_BEFORE_SUSPENDING = "beforeSuspendingCB";
   public static final String ID_BEFORE_SHUTDOWN = "beforeShutdownCB";
   public static final String ID_UPGRADE_POLICY = "upgradePolicyCB";
   public static final String ID_SYNC_TIME = "synchTimeCB";
   public static final String ID_OK = "okButton";
   public static final String ID_CANCEL = "cancelButton";
   public static final String ID_TOOLS_STATUS = "toolsRunningStatus";
   public static final String ID_VMTOOLS_LABEL =
         "viewStack/optionsStack/toolsPage/titleLabel";
   public static final String ID_VMOPTIONS_LINK = "optionsStackButton";
   public static final String ID_BUTTON_VMHARDWARE = "hardwareStackButton";
   public static final String ID_VM_OPTIONS_BOOTDELAY_STEPPER = "bootDelayNum";
   public static final String ID_VM_OPTIONS_BUTTON =
         "vmConfigForm/viewStackSwitcherBox/optionsStackButton";
   public static final String ID_VM_OPTIONS_BIOS_LABEL = "BIOSPage/titleLabel";
   public static final String ID_VM_OPTIONS_CONFIGURATION_BUTTON = "configurationButton";
   public static final String ID_BUTTON_DEBUGGING_AND_STATISTICS = "openButton";
   public static final String ID_VM_OPTIONS_EDITCONFIG_PANEL = "configParametersForm";
   public static final String ID_VM_OPTIONS_CONFIGPARAM_GRID = "configParmGrid";
   public static final String ID_VM_OPTIONS_CONFIGPARAM_ADDROW_BUTTON = "btnAddRow";
   public static final String ID_VM_OPTIONS_CONFIGPARAM_OK_BUTTON = "btnOk";
   public static final String ID_VM_OPTIONS_CONFIGPARAM_CANCEL_BUTTON = "btnCancel";
   // TIWO related
   // TODO Once we have an ID we should change that
   public static final String ID_CONFIRM_REMOVAL_ALERT =
         "automationName=Confirm removal";
   public static final String ID_TIWO_SIDEBAR_BUTTON = "vsphere.core.tiwo.sidebarView";
   public static final String ID_TIWO_TASKS_LIST = "tiwoListView";
   public static final String ID_TIWO_BUTTON = "vsphere.core.tiwo.sidebarView";
   public static final String ID_TIWO_PANEL =
         "vsphere.core.tiwo.sidebarView.portletChrome";
   public static final String ID_EDIT_VM_CPU_SCHED_AFFINITY =
         "affinityControl/txtAffinity";
   public static final String ID_TIWO_DELETE_ICON = "deleteIcon";
   public static final String ID_TIWO_ITEM_DESCRIPTION = "itemDescription";
   public static final String ID_CONFIRM_REMOVAL_YES_BUTTON =
         "automationName=Confirm removal/automationName=Yes";
   public static final String ID_CONFIRM_REMOVAL_NO_BUTTON =
         "automationName=Confirm removal/automationName=No";

   // Cluster - Services Tab
   public static final String ID_CLUSTER_HISTORY_ADVDATAGRID = "historyDataGrid";
   public static final String ID_CLUSTER_DRS_FAULTS_ADVDATAGRID = "faultDataGrid";
   public static final String ID_CLUSTER_DRS_RECOMMENDATIONS_ADVDATAGRID =
         "recommendationDataGrid";
   public static final String ID_CLUSTER_DRS_RECOMMENDATIONS_BUTTON =
         "refreshRecommendationsAction";
   public static final String ID_LABEL_CLUSTER_DRS_RECOMMENDATIONS_EMPTY_LIST =
         "recommendationDataGrid/emptyListIndicator";

   public static final String ID_CHECKBOX_OVEERIDE_RECOMMENDATION = "overrideSwitchedOn";

   // Quick Search
   public static final String ID_SEARCH_CONTAINER = "searchButtonContainer";
   public static final String ID_QUICK_SEARCH_BOX = "searchInput";
   public static final String ID_QUICK_SEARCH_ICON = ID_SEARCH_CONTAINER
         + "/searchButton";

   public static final String ID_QUICK_SEARCH_VIEW = "className=QuickSearchView";
   public static final String ID_QUICK_SEARCH_RESULT_LIST = "quickSearchResults";
   public static final String ID_QUICK_SEARCH_RESULT_ITEMS = ID_QUICK_SEARCH_RESULT_LIST
         + "/resultsBox";
   public static final String ID_RESULTS_GRID = "resultsGrid";
   public static final String ID_CATEGORY_RESULTS = "categoryResultsVBox";
   public static final String ID_NO_RESULTS_FOUND_LABEL = "noResultsFound";
   public static final String ID_SHOW_ALL_LINK = "showAllLink";
   public static final String ID_ADVANCED_SEARCH_PAGE = "advancedSearch";
   public static final String ID_NAV_OBJECT = "title";
   public static final String ID_LOCATING_NODE = "messageLabelBox";

   // Simple Search
   public static final String ID_SIMPLE_SEARCH = "SimpleSearch";
   public static final String ID_SIMPLE_SEARCH_INPUT_BOX =
         "simpleSearch/searchBox/_SearchControl_HBox4/searchInput";
   public static final String ID_SIMPLE_SEARCH_BUTTON = "simpleSearch/searchButton";
   public static final String ID_SIMPLE_SEARCH_RESULTS_GRID = "resultsHost";
   public static final String ID_SIMPLE_SEARCH_PREVIEW = "previewName";
   public static final String ID_SIMPLE_SEARCH_NO_RESULTS =
         "simpleSearch/searchResults/messageBox/messageLabel";
   public static final String ID_SIMPLE_SEARCH_SAVE_IMAGE = "saveButton";
   public static final String ID_SEARCH_NAME_LABEL = "searchNameInput";
   public static final String ID_LABEL_OPEN = "openSearch";
   public static final String ID_LABEL_DELETE = "deleteSearch";
   public static final String ID_LABEL_DELETE_CHECK = "deleteSearch[1]";
   public static final String ID_STRING_SAVE = "Save";
   public static final String ID_STRING_OVERWRITE = "Overwrite";
   public static final String SAVE_SEARCH = "=Save Search";
   public static final String SAVE_SEARCH_DIALOG_STATUS =
         "_SaveSearchDialog_Grid1/statusLabel";
   public static final String ID_LABEL_SIMPLE_SEARCH_RESULTS_TAB =
         "simpleSearch/resultTabs";
   public static final String ID_ADV_SRCH_FULLLABEL = ".fullLabel";

   // Remove Confirmation ID's
   public static final String ID_CONFIRM_REMOVE_YES = "label=Yes";
   public static final String ID_CONFIRM_REMOVE_NO = "label=No";

   // VMOptions
   public static final String ID_VMOPTIONS_GENERAL_VMNAME = "vmNameHeaderText";
   public static final String ID_VMOPTIONS_GENERAL_CONFIG_FILE = "configFileText";
   public static final String ID_VMOPTIONS_GENERAL_WORKING_LOC = "workingLocText";
   public static final String ID_VMOPTIONS_GENERAL_GOS_FAMILY =
         "gosFamilyRow/gosFamilyCombo";
   public static final String ID_VMOPTIONS_GENERAL_GOS_VERSION =
         "gosVersionRow/gosVersionCombo";
   public static final String ID_VMOPTIONS_OK_BUTTON =
         "tiwoDialog/vmConfigForm/buttonContainer/okButton";
   public static final String ID_VMOPTIONS_ADVANCED_DISABLE_ACCELARATION =
         "disableAccelerationCB";
   public static final String ID_VMOPTIONS_ADVANCED_ENABLE_LOGGING = "enableLoggingCB";
   public static final String ID_VMOPTIONS_ADVANCED_DEBUG_STATISTICS = "monitorCombo";
   public static final String ID_VMOPTIONS_ADVANCED_SWAPLOC_DEFAULT = "defaultOption";
   public static final String ID_VMOPTIONS_ADVANCED_SWAPLOC_VMOPTION = "vmOption";
   public static final String ID_VMOPTIONS_ADVANCED_SWAPLOC_HOSTOPTION = "hostOption";
   public static final String ID_VMOPTIONS_BUTTON =
         "vmConfigControl/viewStackSwitcherBox/optionsStackButton";
   public static final String ID_VMOPTIONS_GENERAL_LABEL =
         "optionsStack/generalPage/titleLabel.generalPage";
   public static final String ID_VMOPTIONS_REMOTECONSOLE_LABEL = "titleLabel.vmrcPage";
   public static final String ID_VMOPTIONS_VMTOOLS_LABEL = "titleLabel.toolsPage";
   public static final String ID_VMOPTIONS_REMOTECONSOLE_LOCK_CHECKBOX =
         "gosAutoLockBodyCB";
   public static final String ID_VMOPTIONS_REMOTECONSOLE_LIMIT_CONNECTIONS_CHECKBOX =
         "connectionsCB";
   public static final String ID_VMOPTIONS_REMOTECONSOLE_LIMIT_CONNECTIONS_INPUT =
         "connectionsNum";
   public static final String ID_VMOPTIONS_ADVANCED_LABEL = "titleLabel.advancedPage";
   public static final String ID_VMOPTIONS_BOOT_LABEL = "titleLabel.BIOSPage";
   public static final String ID_VMOPTIONS_BOOT_OPTIONS_COMBO = "firmwareCombo";
   public static final String ID_VMOPTIONS_BOOT_OPTIONS_WARNING = "_BIOSPage_Text3";
   public static final String ID_VMOPTIONS_BOOT_OPTIONS_DELAY_STEPPER = "bootDelayNum";
   public static final String ID_VMOPTIONS_NPIV_LABEL =
         "titleLabel.fibreChannelNpivPage";
   public static final String ID_VM_OPTIONS_NPIV_TEMPDISABLE_CHECKBOX =
         "temporarilyDisableNpiv";
   public static final String ID_VM_OPTIONS_NPIV_LEAVEUNCHANGED_RADIO =
         "wwnLeaveUnchanged";
   public static final String ID_VM_OPTIONS_NPIV_WWNGENERATE_RADIO = "wwnGenerate";
   public static final String ID_VM_OPTIONS_NPIV_WWNREMOVE_RADIO = "wwnRemove";
   public static final String ID_VM_OPTIONS_NPIV_NEWASSIGNMENTS_TEXTAREA =
         "wwnAssignments";
   public static final String ID_VM_OPTIONS_NPIV_WWNN_COMBO = "numberOfNodeWwnsCombo";
   public static final String ID_VM_OPTIONS_NPIV_WWPN_COMBO = "numberOfPortWwnsCombo";
   public static final String ID_VM_OPTIONS_LATENCY_DROPDOWN = "latencyLevel";
   public static final String ID_VM_OPTIONS_LATENCY_TEXT_BOX = "latencyValue";
   public static final String ID_VM_OPTIONS_LATENCY_UNITS = "latencyUnits";
   public static final String ID_VM_OPTIONS_HEADER_VMNAME = "vmNameHeaderText";
   public static final String ID_VM_OPTIONS_DEBUG_STATS_TEXT = "monitorText";
   public static final String ID_VMOPTIONS_ADVANCED_LABEL_MANAGETAB =
         "titleLabel.advanced";
   public static final String ID_LABEL_SETINGS_TAB_VM_NAME = "vmNameHeaderText";
   public static final String ID_LABEL_FORCE_BIOS_SETUP = AUTOMATIONNAME
         + "Force BIOS setup";
   public static final String ID_LABEL_FORCE_EFI_SETUP = AUTOMATIONNAME
         + "Force EFI setup";

   public static final String ID_LINK_MORE_RECOMMENDATIONS = "files_link_2";
   public static final String ID_LABEL_RECOMMENDATION_DATASTORE = "files_value_2";

   // Declaration for latency value
   public static char Mu = (char) 956;
   public static final String ID_LATENCY_UNITS_VAL = Mu + "s";
   public static final String ID_LATENCY_UNITS_MS_VAL = "ms";

   public static final String ID_LOADING_SPINNER = "loadingSpinner";
   // Cluster Power Operations
   public static final String ID_APPLY_RECOMMENDATIONS_BUTTON = "applyButton";
   public static final String ID_APPLY_RECOMMENDATIONS_LIST = "recList";
   public static final String ID_CLUSTER_POWER_ON_RECOMMENDATIONS =
         "powerOnVmResultsForm";
   public static final String ID_CLUSTER_RECOMMENDATIONS_CHECKBOX = "CheckBox_6";
   public static final String ID_CLUSTER_APPLY_RECOMMENDATIONS_CHECKBOX = "CheckBox_11";
   public static final String ID_APPLY_RECOMMENDATIONS_CANCEL_BUTTON = "cancelOKButton";
   public static final String ID_CLUSTER_VMS_LIST = "vmList";
   public static final String ID_FAILED_POWER_ON = "powerOnVmResultsForm";
   public static final String ID_FAULTS_LIST = "faultsList";
   public static final String CONTENT = "content";

   public static final String ID_ADD_VC_BUTTON = "addVcButton";
   public static final String ID_ADD_VC_CONFIRMATION_OK_BUTTON = "AutomationName=OK";

   // Edit vapp
   public static final String ID_EDIT_VAPP_SETTINGS =
         "vsphere.core.vApp.editSettingsAction";
   public static final String ID_EDIT_VAPP_SETTINGS_PRODUCT = "editProduct/titleLabel";
   public static final String ID_EDIT_VAPP_PRODUCT_URL = "productUrl";

   /**
    * public static final String ID_ADD_CIS_BUTTON = "addCIsButton"; public
    * static final String LABEL_APPLY_TEMPLATE = "Apply Template"; public static
    * final String[] ALL_MAIN_APP_AREAS = {LABEL_DASHBOARDS_TAB,
    * LABEL_MANAGEMENT_TAB, LABEL_REPORTING_TAB, LABEL_VCONTROL_TAB,
    * LABEL_VPERSPECTIVE_TAB};
    */

   // Edit VM - vApp Options
   public static final String ID_EDITVM_VAPPOPTIONS_BUTTON = "vAppOptionsStackButton";
   public static final String ID_EDITVM_VAPPOPTIONS_IPALLOCATION_LABEL =
         "titleLabel.ipAllocationPage";
   public static final String ID_EDITVM_VAPPOPTIONS_IPALLOCATION_FIXED_RADIO =
         "fixedButton";
   public static final String ID_VAPP_EDIT_IPALLOCATION_POLICY_TYPE_BUTTON =
         "policyDropDown/openButton";
   public static final String ID_VAPP_EDIT_IPALLOCATION_POLICY_TYPE_LABEL =
         "labelDisplay";
   public static final String ID_VAPP_VM_EDIT_IPALLOCATION_LABEL = "policyLabel";
   public static final String ID_VAPP_VM_EDIT_IPALLOCATION_PROTOCOL_LABEL =
         "protocolDescription";
   public static final String ID_VAPP_VM_EDIT_IPALLOCATION_POLICY_LABEL =
         "policyDescription";
   public static final String ID_EDITVM_VAPPOPTIONS_IPALLOCATION_TRANSIENT_RADIO =
         "transientButton";
   public static final String ID_EDITVM_VAPPOPTIONS_IPALLOCATION_DHCP_RADIO =
         "dhcpButton";

   // Host Maintenance
   public static final String ID_ROOT_HOST_TAB =
         "vsphere.core.folder.hostsView.container";
   public static final String ID_DC_HOST_TAB =
         "vsphere.core.datacenter.hostsView.container";
   public static final String ID_HOST_ENTER_MAINTENANCE_MODE =
         "vsphere.core.host.enterMaintenanceAction";
   public static final String ID_HOST_EXIT_MAINTENANCE_MODE =
         "vsphere.core.host.exitMaintenanceAction";
   public static final String ID_HOST_CONFIRM_MM_YES_BUTTON = "yesButton";
   public static final String ID_HOST_CONFIRM_MM_NO_BUTTON = "noButton";
   public static final String ID_HOST_CONFIRM_MM_PANEL =
         "automationName=Confirm Maintenance Mode";
   public static final String ID_HOST_CONFIRM_MM_ALERT = "automationName=Warning";
   public static final String ID_HOST_CONFIRM_MM_ALERT_OK_BUTTON = "automationName=OK";

   // VMRC
   public static final String ID_PLUGIN_DESRIPTION_TEXT = "pluginDescription";
   public static final String ID_PLUGIN_DOWNLOAD_LINKBUTTON = "downloadLink";
   public static final String ID_CLIENT_INTEGR_PLUGIN_TEXT = "_VmrcPluginView_Text3";
   public static final String ID_PLUGIN_VERSION_TEXT = "pluginVersion";
   public static final String ID_VMRC_FIRST_POPUP_WINDOW = "vmrc_0";
   public static final String ID_VMRC_FLEX_APP = "vmrc-app";
   public static final String ID_VMRC_CAD_BUTTON = "cadButton";
   public static final String ID_VMRC_OFF_BUTTON = "btnPowerOff";
   public static final String ID_VMRC_SUSPEND_BUTTON = "btnSuspend";
   public static final String ID_VMRC_ON_BUTTON = "btnPowerOn";
   public static final String ID_VMRC_RESET_BUTTON = "btnReset";
   public static final String ID_VMRC_FULL_SCR_BUTTON = "fullscreenButton";
   public static final String ID_VMRC_CONSOLE_BUTTONBARBUTTON = "view_0";
   public static final String ID_VMRC_SETTINGS_BUTTONBARBUTTON = "view_1";
   public static final String ID_CONTEXTMENU_OPEN_CONSOLE = new StringBuffer(
         ID_CONTEXT_MENU).append(".").append(ID_ACTION_OPEN_CONSOLE).toString();
   public static final String ID_CONTEXTMENU_GUEST_OS = "guestOs";
   public static final String ID_GUEST_OS_OPEN_CONSOLE = new StringBuffer(
         ID_CONTEXT_MENU).append(".").append(ID_CONTEXTMENU_GUEST_OS).append(".")
         .append(ID_ACTION_OPEN_CONSOLE).toString();
   public static final String ID_TOOLBAR_DATAGRID_OPEN_CONSOLE = ID_ACTION_OPEN_CONSOLE
         + "/button";
   public static final String ID_ACTION_OPEN_CONSOLE_BUTTON = ID_ACTION_OPEN_CONSOLE
         + "/button";

   // More Actions
   public static final String ID_MOREACTIONS_LHS = "lhs";
   public static final String ID_MOREACTIONS_FILTER_TXT = "filterText";
   public static final String ID_MOREACTIONS_ONLY_AVAILABLE = "onlyAvailable";
   public static final String ID_MOREACTIONS_FILTER_LABEL = "filterLabel";
   public static final String ID_MOREACTIONS_VCENTER_MGMT = "toc_app";
   public static final String ID_MOREACTIONS_ALLACTIONS = "toc_all";
   public static final String ID_MOREACTIONS_PANEL = "moreActionsAnchoredDialog";
   public static final String ID_ADVANCE_DATAGRID_ALLACTIONS = "allActions/button";

   // VM Upgrade HW
   public static final String ID_VM_CONFIRM_HW_UPG =
         "automationName=Confirm Virtual Machine Upgrade";
   public static final String ID_VM_CONFIRM_HW_UPG_YES_BUTTON = "YES";
   public static final String ID_VM_CONFIRM_HW_UPG_NO_BUTTON = "NO";
   public static final String ID_VM_CONFIRM_HW_UPG_ANSWER_RADIO = "vmQuestions/choice1";
   public static final String ID_VM_CONFIRM_HW_UPG_ANSWER_RADIO_ANSWER =
         "answerQuestion";
   public static final String ID_EDIT_VM_ERROR_LABEL = "addHw/addHardware/message";
   public static final String ID_EDIT_VM_UPGRADE = "automationName=Upgrade";
   public static final String ID_VM_UPGRADE_HW_ON_NEXT_POWERON = "scheduled";
   public static final String ID_VERSION_LABEL = "labelDisplay";
   public static final String ID_VM_VERSION_SUMMARY = "summary_vmVersion_valueLbl";
   public static final String ID_VM_COMPATIBILITY_DROPDOWN =
         "SinglePageDialog_ID/comboBoxHelp";
   public static final String ID_VM_VERSION_STATUS_SUMMARYTAB = "hwVersionUpgrade";
   public static final String ID_UPGRADE_AFTER_SHUTDOWN = "onSoftPowerOff";
   public static final String ID_HOST_DEFAULT_COMPATIBLITY = "textDesc";
   // HW version Displayed
   public static final String ID_VM_SUMMARY_COMPATIBILITY_PREFIX = "VM Compatibility: ";
   public static final String ID_VM_HW_VERSION_7 = "ESX 4.0 and later";
   public static final String ID_VM_HW_VERSION_8 = "ESX 5.0 and later";
   public static final String ID_VM_HW_VERSION_9 = "ESXi 5.1 and later";
   public static final String ID_VM_VERSION_DROPDOWN = "vmVersionSelector/comboBox";

   // Advanced Search
   public static final String ID_OPERATOR_COMBOBOX = "operatorComboBox";
   public static final String ID_VALUE_COMBOBOX = "valueComboBox";
   public static final String ID_VALUE_INPUT = "valueInput";
   public static final String ID_TEXT_INPUT = "textInput";
   public static final String ID_COMBOBOX = "ComboBox";
   public static final String ID_TEXTINPUT = "TextInput";
   public static final String ID_TEXTINPUTEX = "TextInputEx";
   public static final String ID_MULTILINETOKENSCONTAINER_TOKEN =
         "multiLineTokensContainer/token_";
   public static final String ID_ADV_SRCH_REFRESH_BUTTON =
         "advancedSearch/searchResults/resultsBox/resultDetailsPanel/refreshButton";
   public static final String ID_ADV_SEARCH_LINK = "simpleSearch/advancedSearchLink";
   public static final String ID_ADV_SEARCH_MESSAGE = "messageLabel";
   public static final String ID_SIMPLE_SEARCH_LINK = "simpleSearchLink";
   public static final String ID_ADV_SEARCH = "_AdvancedSearchView_Label1";
   public static final String ID_SEARCH_BY_ENTITY = "typeSelector";
   public static final String ID_PROPERTY_SELECTOR = "propertySelector";
   public static final String ID_CONJOINER_SELECTOR = "conjoinerSelector";
   public static final String ID_ADVANCED_SEARCH_BUILDER = "advancedSearchBuilder_";
   public static final String ID_CONTENTROWS_ADVANCEDSEARCHROW =
         "contentRows/advancedSearchRow_";
   public static final String ID_ADD_CRITERIA_LINK = "addRowLink";
   public static final String ID_ADVANCEDSEARCHROW = "advancedSearchRow_";
   public static final String ID_CRITERIA_COMPARATOR = "token_0";
   public static final String ID_CRITERIA_VALUE = "token_1";
   public static final String ID_TEXT_INPUT_CRITERIA_VALUE = "TextInputTokenEditor";
   public static final String ID_COMBO_BOX_CRITERIA_VALUE = "ComboboxTokenEditor";
   public static final String ID_NUMERIC_CRITERIA_VALUE = "NumericTokenEditor";
   public static final String ID_RICHTXT_CRITERIA_VALUE = "RichEditableText";
   public static final String ID_NUMBERINPUT_CRITERIA_VALUE = "NumberInput";
   public static final String ID_LABEL_CRITERIA_VALUE = "LabelTokenEditor";
   public static final String ID_TEXT = "text";
   public static final String ID_NAME = "name";
   public static final String ID_BUTTON_SEARCH = "advancedSearch/searchButton";
   public static final String ID_BUTTON_SAVE = "advancedSearch/saveButton";
   public static final String ID_TYPE_SELECTOR = "typeSelector";
   public static final String ID_CUSTOM_ATTRIBUTE_VALUE_TEXT = "token_5";
   public static final String ID_CUSTOM_ATTRIBUTE_VALUE_CONDITION = "token_4";
   public static final String ID_CUSTOM_ATTRIBUTE_NAME_TEXT = "token_2";
   public static final String ID_ADV_SEARCH_RESULTS_GRID = "resultsHost";
   public static final String ID_NEW_SEARCHTYPE_LINK = "addBuilderLink";
   public static final String ID_REMOVE_CRITERIA_LINK = "removeRowLink";
   public static final String ID_REMOVE_SEARCHTYPE_LINK = "removeBuilderLink";
   public static final String ID_QUICKSEARCH_SAVEDLIST =
         "savedSearchesPopup/contentStack/savedSearchesContainer/savedList";
   public static final String ID_BUTTON_ADDITIONAL =
         "mainControlBar/toolbarViewsBar/searchControl/additionalButton";
   public static final String ID_FLEXBOX_QUICKSEARCH_SEARCH_CONTAINER =
         "contentStack/savedSearchesContainer";
   public static final String ID_LABEL_PREVIEWTITLE =
         "appBody/searchViews/advancedSearch/previewName";
   public static final String ID_LABEL_INVENTORY_PATH = "previewInventoryPath";
   public static final String ID_TEXT_SUMMARYVIEW_CONTAINER = "summaryView.container";
   public static final String ID_TEXT_GETTING_STARTED_CONTAINER =
         "gettingStarted.container";
   public static final String ID_TEXT_PREVIEW_DOT_SPACE = "Preview: ";
   public static final String ID_BUTTON_SAVESEARCHDIALOG_CANCEL = "cancelButton";
   public static final String ID_BUTTON_SAVESEARCHDIALOG_SAVE = "saveButton";
   public static final String ID_SAVEDLIST = "savedList";
   public static final String ID_LABEL_SAVESEARCHDIALOG_INPUTLABEL = "searchNameInput";

   // CIM
   public static final String ID_TREE_VIEW_CIM_SENSOR = "sensorsTreeView";
   public static final String ID_VIEW_CIM_EVENT_LOG = "eventsLogView";
   public static final String ID_BUTTON_BAR_CIM_ACTION_FUNCTION =
         "_SensorView_ActionfwButtonBar1";
   public static final String ID_LABEL_CIM_SENSORS = "_TocListItemRenderer_Label1";
   public static final String ID_LIST_CIM_TOCLIST = "tocList";
   public static final String ID_BUTTON_CIM_UPDATE =
         "vsphere.core.cimmonitor.updateHardwareData/button";
   public static final String ID_BUTTON_CIM_RESET_SENSORS =
         "vsphere.core.cimmonitor.resetHostSensors/button";
   public static final String ID_BUTTON_CIM_EXPORT_DATA =
         "vsphere.core.cimmonitor.exportData/button";
   public static final String ID_BUTTON_CIM_COPY_TO_CLIPBOARD = "copyToClipboard";
   public static final String ID_TEXT_FIELD_CIM_EXPORT_HARDWARRE_DETAILS = "xmlText";
   public static final String ID_BUTTON_CIM_EXPORT_SAVE_AS = "_ExportDataDialog_Button1";
   public static final String ID_BUTTON_CIM_EXPORT_CLOSE = "_ExportDataDialog_Button2";
   public static final String ID_BUTTON_CIM_RESET_EVENT_LOG =
         "vsphere.core.cimmonitor.resetHostEventsLog/button";
   public static final String ID_LABEL_CIM_CIMDATA_TITLE = "titleLabel";
   public static final String ID_LABEL_CIM_ALERTWARNING_TITLE = "titleGroup/titleLabel";
   public static final String ID_LABEL_CIM_UPDATED_TIME_INFORMATION =
         "_HardwareStatusView_Label1";
   public static final String ID_LABEL_CIM_BIOS_VIEW_INFORMATION = "_BiosView_Label?";
   public static final String ID_LABEL_CIM_SYSTEM_BOARD_INFORMATION =
         "_SystemBoardView_Label?";
   public static final String ID_LABEL_CIM_ALERTS_AND_WARNINGS_INFORMATION =
         "normalStatusLabel";
   public static final String ID_LABEL_NO_ALERTS_WARNINGS = "emptyListIndicator";
   public static final String ID_LABEL_NO_HARDWARE_MONITERING_SERVICE_WARNING =
         "_ErrorsView_Label1";
   public static final String ID_LABEL_NO_HOST_DATA_WARNING = "emptyListIndicator";

   // User World Swap.
   public static final String ID_LABEL_USER_WORLD_SWAP_LINK =
         "_TocListItemRenderer_Label1";
   public static final String ID_BUTTON_USER_WORLD_SWAP_EDIT =
         "btn_vsphere.core.host.editUserWorldSwapSettings";
   public static final String ID_CHECK_BOX_USER_WORLD_SWAP_ENABLE = "enabledCheckBox";
   public static final String ID_CHECK_BOX_USER_WORLD_SWAP_CAN_USE_DATASTORE =
         "userSpecifiedDatastoreCheckBox";
   public static final String ID_CHECK_BOX_USER_WORLD_SWAP_CAN_USE_HOST_CACHE =
         "hostCacheCheckBox";
   public static final String ID_CHECK_BOX_USER_WORLD_SWAP_CAN_USE_VM_SWAPFILE_LOCATION =
         "localDatastoreCheckBox";
   public static final String ID_DROP_DOWN_LIST_USER_WORLD_SWAP_DATASTORE =
         "datastoreNameList";
   public static final String ID_DROP_DOWN_BUTTON_USER_WORLD_SWAP_DATASTORE =
         "openButton";
   public static final String ID_BUTTON_USER_WORLD_SWAP_HELP = "helpButton";
   public static final String ID_LABEL_USER_WORLD_SWAP = "labelItem";
   public static final String ID_LABEL_USER_WORLD_SWAP_DISABLED = "disabled";
   public static final String ID_LABEL_USER_WORLD_SWAP_USE_HOST_CACHE = "hostCache";
   public static final String ID_LABEL_USER_WORLD_SWAP_USE_DATASTORE =
         "userSpecifiedDatastore";
   public static final String ID_LABEL_USER_WORLD_SWAP_USE_VM_SWAPFILE_LOCATION =
         "hostLocalSwapDatastore";
   public static final String ID_LABEL_USER_WORLD_SWAP_DESCRIPTION = "descriptionLabel";
   public static final String ID_PROPERTY_GRID = "propertyGrid";

   // NetDump
   public static final String ID_CHECKBOX_NETDUMP_ENABLED = "Enabled";
   public static final String ID_TEXT_INPUT_HOSTvNIC_TO_USE = "HostVNic";
   public static final String ID_TEXT_INPUT_NETDUMP_SERVER_IP = "NetworkServerIP";
   public static final String ID_TEXT_INPUT_NETDUMP_SERVER_PORT =
         "NetworkServerPort/className=TextInputEx";
   public static final String ID_UID_HOSTPROFILE_NETWORKING_CONFIGURATION = "0.3.";
   public static final String ID_UID_HOSTPROFILE_NETWORK_COREDUMP_SETTINGS = "0.3.10.";
   public static final String ID_IMAGE_REMEDIATION_RESULTS_ARROW_IMAGE =
         "remediationResults/arrowImage";
   public static final String ID_LABEL_ERROR_STACK_MESSAGE = "stackLabel1";
   public static final String ID_LABEL_ERROR_REMEDIATE_HOST_FAILED =
         "error/_container/_message";
   public static final String ID_LABEL_ERROR_REMEDIATE_HOST_FAILED_MESSAGE =
         "remediationResults/gridRowLabelId";
   public static final String ID_LABEL_SERVER_IP_ADDRESS = "automationName=<HOST_NAME>";
   public static final String ID_LABEL_LISTENING_ON_PORT =
         "automationName=<NETDUMP_SERVER_PORT>";
   public static final String ID_TOC_TREE_ITEM_ESXI_DUMP_COLLECTOR = "5";

   // Stateless Esxi
   public static final String ID_TREE_HOST_MANAGE_SETTINGS_TOC_TREE =
         "vsphere.core.host.manage.settingsView/tocTree";
   public static final String ID_RADIO_BUTTON_USE_AUTHENTICATION =
         "AuthJoinDomainForm_RadioButton1";
   public static final String ID_RADIO_BUTTON_USE_CAM_SERVER = "camServerRb";
   public static final String ID_BUTTON_IMPORT_CERTIFICATE =
         "btn_vsphere.core.host.importCertificate";
   public static final String ID_TEXT_DISPLAY_CERTIFICATE_PATH = "certPath/textDisplay";
   public static final String ID_TEXT_DISPLAY_AUTHENTCATION_PROXY_SERVER_IP_ADDRESS =
         "ipAddressInputTop";
   public static final String ID_TEXT_INPUT_AUTHENTICATION_SERVICES_DOMAIN_NAME =
         "domain/textDisplay";
   public static final String ID_CHECK_BOX_FIREWALL_SETTINGS_ALLOW_ALL_IP = "allIp";
   public static final String ID_LABEL_POWER_MANAGE_TECHNOLOGY =
         "HostSystem:powerConfigHardware.technology";
   public static final String ID_LABEL_POWER_SYSTEM_CUSTOM_POLICY =
         "HostSystem:powerConfigHardware.powerSystemCurrentPolicy";
   public static final String ID_BUTTON_EDIT_POWER_MANAGEMENT = "btnProperties";
   public static final String ID_BUTTON_ASSIGN_LICENSEKEY =
         "btn_vsphere.license.management.host.assignLicenseKeyAction";
   public static final String ID_CHECKBOX_VMOTION_TRAFFIC = "mngTrafficCkb1";
   public static final String ID_LABEL_TOTAL_PHYSICAL_MEMORY_LABEL =
         "titleLabel.HostSystem:memoryConfig.totalPhysicalMemoryTitle";
   public static final String ID_LABEL_TOTAL_SYSTEM_MEMORY_LABEL =
         "titleLabel.HostSystem:memoryConfig.systemPhysicalMemoryTitle";
   public static final String ID_LABEL_TOTAL_VM_MEMORY_LABEL =
         "titleLabel.HostSystem:memoryConfig.vmPhysicalMemoryTitle";
   public static final String ID_BUTTON_EDIT_PROCESSOR_SETTINGS =
         "Edit_HostSystem:processorGeneralConfig";
   public static final String ID_LIST_NETWORK_VIRTUAL_SWITCH_LIST = "switchList";


   /*
    * public static final String ID_SEARCH_NAME_LABEL = "searchNameInput";
    * public static String ID_LABEL_OPEN = "openSearch"; public static String
    * ID_LABEL_DELETE = "deleteSearch"; public static String
    * ID_LABEL_DELETE_CHECK = "deleteSearch[1]"; public static String
    * ID_STRING_SAVE = "Save"; public static String ID_STRING_OVERWRITE =
    * "Overwrite"; public static String SAVE_SEARCH = "=Save Search";
    */

   // Plugins Management
   public static final String ID_PLUGINS_MANAGEMENTUI_SAMPLE_PORTLET = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_VM).append("sampleView1")
         .append(ID_PORTLET_SUFFIX).toString();
   public static final String ID_PLUGINS_MANAGEMENTUI_SAMPLE_VM_VIEW_TAB =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_VM)
               .append("sampleView2").append(".container").toString();
   public static final String ID_PLUGINS_MANAGEMENTUI_SAMPLE_PORTLET_ALERT =
         new StringBuffer(EXTENSION_PREFIX).append("sample.sampleSideBarView")
               .toString();
   public static final String ID_PLUGINS_MANAGEMENTUI_SIDEBAR_SAMPLE_PORTLET =
         new StringBuffer(EXTENSION_PREFIX).append(
               "sample.sampleSideBarView.portletChrome").toString();
   public static final String ID_PLUGINS_MANAGEMENTUI_DATAGRID =
         "vsphere.core.admin.pluginApplication";
   public static final String ID_PLUGINS_MANAGEMENTUI_ENABLE_BUTTON = new StringBuffer(
         EXTENSION_PREFIX).append("admin.enablePluginAction").toString();
   public static final String ID_PLUGINS_MANAGEMENTUI_DISABLE_BUTTON = new StringBuffer(
         EXTENSION_PREFIX).append("admin.disablePluginAction").toString();

   // Snapshots
   public static final String ID_TAKE_SNAPSHOT_PANEL = "tiwoDialog";
   public static final String ID_TAKE_SNAPSHOT_NAME = "snapshotName";
   public static final String ID_TAKE_SNAPSHOT_DESC = "snapshotDescription";
   public static final String ID_TAKE_SNAPSHOT_MEM = "snapshotView/snapshotMemory";
   public static final String ID_TAKE_SNAPSHOT_QUIESCE = "snapshotQuiesce";
   public static final String ID_SNAPSHOT_MNGR_PANEL = "tiwoDialog";
   public static final String ID_CLOSE_BUTTON = "tiwoDialog/exitButton";
   public static final String ID_DELETE_ALL_BUTTON = "tiwoDialog/deleteAllButton";
   public static final String ID_SNAPSHOT_TREE = "snapshotTree";
   public static final String ID_EDIT_BUTTON = "tiwoDialog/editButton";
   public static final String ID_DELETE_BUTTON = "tiwoDialog/deleteButton";
   public static final String ID_GOTO_BUTTON = "tiwoDialog/revertButton";
   public static final String ID_DELETE_CONFIRM_ALERT = "automationName=Confirm Delete";
   public static final String ID_DELETE_CONFIRM_ALERT_YES = "automationName=Yes";
   public static final String ID_DELETE_CONFIRM_ALERT_NO = "automationName=No";
   public static final String ID_TAKE_SNAPSHOT_CANCEL_BUTTON = "tiwoDialog/cancelButton";
   public static final String ID_GOTO_CONFIRM_ALERT =
         "automationName=Confirm go to Snapshot";
   public static final String ID_GOTO_CONFIRM_ALERT_YES = "automationName=Yes";
   public static final String ID_GOTO_CONFIRM_ALERT_NO = "automationName=No";
   public static final String ID_REVERT_TO_CURRENT_CONFIRM_ALERT =
         "automationName=Confirm go to Current Snapshot";
   public static final String ID_REVERT_TO_CURRENT_CONFIRM_ALERT_YES =
         "automationName=Yes";
   public static final String ID_REVERT_TO_CURRENT_CONFIRM_ALERT_NO =
         "automationName=No";
   public static final String ID_SNAPSHOT_MANAGER_DISK_CAPACITY =
         "viewGrid/diskUsageLabel/diskUsage";
   public static final String ID_CONSOLIDATE_SINGLE_LABEL = "issueDescription0";
   public static final String ID_CONSOLIDATE_DOUBLE_LABEL = "issueDescription1";
   // vApp StartOrder
   public static final String ID_VAPP_START_ORDER_PANEL = "tiwoDialog";
   public static final String ID_VAPP_START_ORDER_STARTUP_PANEL = "startupPanel";
   public static final String ID_VAPP_START_ORDER_SHUTDOWN_PANEL = "shutdownPanel";
   public static final String ID_VAPP_START_ORDER_DOWN_BUTTON =
         "startOrder/buttonVBox/downButton";
   public static final String ID_VAPP_START_ORDER_UP_BUTTON =
         "startOrder/buttonVBox/upButton";
   public static final String ID_VAPP_START_ORDER_LIST =
         "startOrder/listVBox/startOrderList";
   public static final String ID_VAPP_START_ORDER_START_UP_DELAY_STEPPER =
         ID_VAPP_START_ORDER_STARTUP_PANEL + "/startupDelayHBox/startupDelay";
   public static final String ID_VAPP_SHUTDOWN_DELAY_STEPPER =
         ID_VAPP_START_ORDER_SHUTDOWN_PANEL + "/shutdownDelay";
   public static final String ID_VAPP_START_ORDER_LABEL = "titleLabel.startOrder";
   public static final String ID_VAPP_STARTUP_COMBO = ID_VAPP_START_ORDER_STARTUP_PANEL
         + "/startupOperationHBox/startupOperation";
   public static final String ID_VAPP_SHUTDOWN_COMBO =
         ID_VAPP_START_ORDER_SHUTDOWN_PANEL + "/shutdownOperationHBox/shutdownOperation";
   public static final String ID_VAPP_STARTUP_TOOLS_CHECK_BOX =
         ID_VAPP_START_ORDER_STARTUP_PANEL + "/startupToolsReady";

   // Host Storage
   public static final String ID_HOST_STORAGE_CAPACITY_PANEL =
         "vsphere.core.hosts.datastore.capacityView.chrome";
   public static final String ID_HOST_STORAGE_DETAILS_PANEL =
         "vsphere.core.hosts.datastore.detailsView.chrome";
   public static final String ID_HOST_STORAGE_PATHS_PANEL =
         "vsphere.core.hosts.datastore.pathsView.chrome";
   public static final String ID_HOST_STORAGE_CAPACITY_TOTAL =
         "dsCapacityPropertyGrid/_DatastoreCapacityView_PropertyGridRow1/dsCapacity";
   public static final String ID_HOST_STORAGE_CAPACITY_PROVISIONED =
         "dsCapacityPropertyGrid/_DatastoreCapacityView_PropertyGridRow2/dsProvisioned";
   public static final String ID_HOST_STORAGE_CAPACITY_FREE =
         "dsCapacityPropertyGrid/_DatastoreCapacityView_PropertyGridRow3/dsFreeSpace";
   public static final String ID_HOST_DETAILS_LOCATION =
         "dsGeneralPropertyGrid/_DatastoreDetailsView_PropertyGridRow1/dsUrl";
   public static final String ID_HOST_DETAILS_FILESYS_TYPE =
         "dsGeneralPropertyGrid/_DatastoreDetailsView_PropertyGridRow2/dsType";
   public static final String ID_HOST_DETAILS_NO_OF_HOSTS =
         "dsGeneralPropertyGrid/_DatastoreDetailsView_PropertyGridRow3/dsHostCount";
   public static final String ID_HOST_DETAILS_VMS_AND_TEMPLATES =
         "dsGeneralPropertyGrid/_DatastoreDetailsView_PropertyGridRow4/dsVmCount";
   public static final String ID_HOST_DETAILS_MAX_FILE_SIZE =
         "dsGeneralPropertyGrid/_DatastoreDetailsView_PropertyGridRow5/dsMaxFileSize";
   public static final String ID_HOST_DETAILS_CONGESTION_MGT =
         "_DatastoreDetailsView_PropertyGridRow6/dsCongestionMgmt";
   public static final String ID_HOST_DETAILS_HW_ACCELARATION =
         "dsGeneralPropertyGrid/pgrVStorage/dsVStorage";
   public static final String ID_HOST_PATHS_TOTAL =
         "viewPaths/_DatastorePathCountView_PropertyGridRow1/dsTotalPaths";
   public static final String ID_HOST_PATHS_BROKEN =
         "viewPaths/_DatastorePathCountView_PropertyGridRow2/dsBrokenPaths";
   public static final String ID_HOST_PATHS_DISABLED =
         "viewPaths/_DatastorePathCountView_PropertyGridRow3/dsDisabledPaths";

   // Host AutoStart Config
   public static final String ID_AUTOSTARTSTOP_EDIT_BTN =
         "btn_vsphere.core.host.editVmStartupAction";
   public static final String ID_AUTOSTARTSTOP_CANCEL_BUTTON = "cancelButton";
   public static final String ID_AUTOSTARTSTOP_OK_BUTTON = "okButton";
   public static final String ID_GRID_AUTOSTART = "gridAutoStart";
   public static final String ID_AUTOSTARTSTOP_MOVE_DOWN_BUTTON =
         "tiwoDialog/moveDown/button";
   public static final String ID_AUTOSTARTSTOP_MOVE_UP_BUTTON =
         "tiwoDialog/moveUp/button";
   public static final String ID_AUTOSTARTSTOP_DATA_PROVIDER =
         "dataProvider.source.source.";

   // Host Power Management
   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_SPARK_LIST =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_HOST)
               .append("manage.").append("settingsView").append("/tocTree").toString();
   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_PROPERTIES_BUTTON =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_HOST)
               .append("manage.").append("settingsView").append("/btnProperties")
               .toString();
   public static final String ID_HOST_SETTINGS_SYSTEM_MANAGEMENT_SPARK_LIST =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_HOST)
               .append("manage.").append("settingsView").append("/tocListSystem")
               .toString();

   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_MANAGE_TITLE_LABEL =
         "titleLabel";

   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_MANAGE_TECHNOLOGY_LABEL =
         "titleLabel.HostSystem:powerConfigHardware.technologyTitle";
   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_MANAGE_ACTIVE_POLICY_LABEL =
         "titleLabel.HostSystem:powerConfigHardware.powerSystemCurrentPolicyTitle";
   public static final String ID_LABEL_HOST_SETTINGS_POWER_MANAGEMENT_TECHBOLOGY_VALUE =
         "HostSystem:powerConfigHardware.technology";
   public static final String ID_LABEL_HOST_SETTINGS_POWER_MANAGEMENT_ACTIVE_POLICY_VALUE =
         "HostSystem:powerConfigHardware.powerSystemCurrentPolicy";

   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_EDIT_PROPERTIES_DIALOG_PANEL =
         "tiwoDialog";
   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_EDIT_PROPERTIES_DIALOG_PANEL_HIGH_PERFORMANCE_RADIO_BUTTON =
         "automationName=High performance";
   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_EDIT_PROPERTIES_DIALOG_PANEL_BALANCED_RADIO_BUTTON =
         "automationName=Balanced";
   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_EDIT_PROPERTIES_DIALOG_PANEL_LOW_POWER_RADIO_BUTTON =
         "automationName=Low power";
   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_EDIT_PROPERTIES_DIALOG_PANEL_CUSTOM_RADIO_BUTTON =
         "automationName=Custom";

   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_EDIT_PROPERTIES_DIALOG_PANEL_HELP_BUTTON =
         "helpButton";
   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_EDIT_PROPERTIES_DIALOG_PANEL_OK_BUTTON =
         "okButton";
   public static final String ID_HOST_SETTINGS_POWER_MANAGEMENT_EDIT_PROPERTIES_DIALOG_PANEL_CANCEL_BUTTON =
         "cancelButton";

   // Datacenter IP Pools
   public static final String ID_DC_IP_POOL_TAB =
         "vsphere.core.datacenter.ipPoolsView.container";
   public static final String ID_DC_IP_POOL_LIST = "ipPoolList";
   public static final String ID_DC_IPPOOL_IPV4_POOL =
         "detailsStackEditor/ipv4ConfigControl/titleLabel";
   public static final String ID_DC_IPPOOL_IPV6_POOL =
         "detailsStackEditor/ipv6ConfigControl/titleLabel";
   public static final String ID_DC_IPPOOL_OTHER_NW_PROPERTIES =
         "detailsStackEditor/otherNetworkProperties/titleLabel";
   public static final String ID_DC_IPPOOL_IPV4_POOL_SUBNET =
         "ipv4ConfigControl/propertyGrid/subnetGridRow/subnetLabel";
   public static final String ID_DC_IPPOOL_IPV4_POOL_NETMASK =
         "ipv4ConfigControl/propertyGrid/netmaskGridRow/netmaskLabel";
   public static final String ID_DC_IPPOOL_IPV4_POOL_GATEWAY =
         "ipv4ConfigControl/propertyGrid/gatewayGridRow/gatewayLabel";
   public static final String ID_DC_IPPOOL_IPV4_POOL_IPPOOL =
         "ipv4ConfigControl/propertyGrid/ipPoolGridRow/ipPoolTextArea";
   public static final String ID_DC_IPPOOL_IPV4_POOL_DHCP_PRESENT =
         "ipv4ConfigControl/propertyGrid/dhcpGridRow/dhcpLabel";
   public static final String ID_DC_IPPOOL_IPV4_POOL_DNS_SERVERS =
         "ipv4ConfigControl/propertyGrid/dnsGridRow/dnsTextArea";
   public static final String ID_DC_IPPOOL_IPV6_POOL_SUBNET =
         "ipv6ConfigControl/propertyGrid/subnetGridRow/subnetLabel";
   public static final String ID_DC_IPPOOL_IPV6_POOL_NETMASK =
         "ipv6ConfigControl/propertyGrid/netmaskGridRow/netmaskLabel";
   public static final String ID_DC_IPPOOL_IPV6_POOL_GATEWAY =
         "ipv6ConfigControl/propertyGrid/gatewayGridRow/gatewayLabel";
   public static final String ID_DC_IPPOOL_IPV6_POOL_IPPOOL =
         "ipv6ConfigControl/propertyGrid/ipPoolGridRow/ipPoolTextArea";
   public static final String ID_DC_IPPOOL_IPV6_POOL_DHCP_PRESENT =
         "ipv6ConfigControl/propertyGrid/dhcpGridRow/dhcpLabel";
   public static final String ID_DC_IPPOOL_IPV6_POOL_DNS_SERVERS =
         "ipv6ConfigControl/propertyGrid/dnsGridRow/dnsTextArea";
   public static final String ID_DC_IPPOOL_OTHER_NW_PROPERTIES_DNS_DOMAIN =
         "otherNetworkProperties/otherPropertyGrid/dnsDomainGridRow/dnsDomainLabel";
   public static final String ID_DC_IPPOOL_OTHER_NW_PROPERTIES_HOST_PREFIX =
         "otherNetworkProperties/otherPropertyGrid/hostPrefixGridRow/hostPrefixLabel";
   public static final String ID_DC_IPPOOL_OTHER_NW_PROPERTIES_DNS_SERVERS =
         "otherNetworkProperties/otherPropertyGrid/dnsSearchPathGridRow/dnsSearchPathLabel";
   public static final String ID_DC_IPPOOL_OTHER_NW_PROPERTIES_HTTP_PROXY =
         "otherNetworkProperties/otherPropertyGrid/httpProxyGridRow/httpProxyLabel";
   public static final String ID_DC_IPPOOL_OTHER_NW_PROPERTIES_ASSOCIATED_NWS =
         "otherNetworkProperties/otherPropertyGrid/networksGridRow/networksTextArea";

   //Network Protocol Profile
   public static final String ID_NETWORK_PROFILE_TAB =
         "_ScrollingButtonBarButtonRenderer_Label1";
   public static final String ID_NETWORK_PROTOCOL_PROFILE_ADD = "ipPoolAddAction";
   public static final String ID_NETWORK_PROTOCOL_PROFILE_EDIT = "ipPoolEditAction";
   public static final String ID_NETWORK_PROTOCOL_PROFILE_DELETE = "ipPoolDeleteAction";
   public static final String ID_NETWORK_PROTOCOL_PROFILE_ADD_PANEL = "tiwoDialog";
   public static final String ID_NETWORK_PROTOCOL_PROFILE_EDIT_PANEL = "tiwoDialog";
   public static final String ID_NETWORK_PROTOCOL_PROFILE_TEXT_INPUT = "nameTextInput";
   public static final String ID_NETWORK_PROTOCOL_PROFILE_BUTTON_FINISH = "finish";
   public static final String ID_NETWORK_PROTOCOL_PROFILE_BUTTON_CANCEL = "cancelButton";
   public static final String ID_NETWORK_PROTOCOL_PROFILE_BUTTON_NEXT = "next";
   public static final String ID_NETWORK_PROTOCOL_PROFILE_GRID_LIST = "list";
   public static final String ID_NETWORK_PROTOCOL_PROFILE_GRID_COLUMN_NAME = "Name";
   public static final String ID_NETWORK_PROTOCOL_PROFILE_NET_LABEL = "netLabel";


   // Folder Summary
   public static final String ID_FOLDER_SUMMARY_TAB = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_FOLDER).append(ID_TEXT_SUMMARYVIEW_CONTAINER)
         .toString();
   public static final String ID_FOLDER_STATUS_PORTLET = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_FOLDER).append("statusView")
         .append(ID_PORTLET_SUFFIX).toString();
   public static final String ID_FOLDER_DETAILS_PORTLET = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_FOLDER).append("detailsView")
         .append(ID_PORTLET_SUFFIX).toString();
   public static final String ID_FOLDER_STATUS_ACTIVE_TASKS = "tasks";
   public static final String ID_FOLDER_DETAILS_VMS = "dcVms";
   public static final String ID_FOLDER_DETAILS_HOSTS = "dcHosts";

   // Advance Search
   public static final String ADVS_MESSAGE_LABEL =
         "advancedSearch/messageBox/messageLabel";

   // Reporting Automation Ids
   public static final String ID_LAUNCH_POPUP_BUTTON = "appsButton";
   public static final String ID_LAUNCH_MENU = "launchMenu";

   // VM Resource Management Tab Ids
   public static final String ID_VM_RES_MGMT_HOSTCPU = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "monitor.resMgmt.hostCPUView.chrome";
   public static final String ID_VM_RES_MGMT_HOSTMEM = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "monitor.resMgmt.hostMemoryView.chrome";
   public static final String ID_VM_RES_MGMT_GUESTMEM = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "monitor.resMgmt.guestMemoryView.chrome";
   public static final String ID_VM_RES_MGMT_CONSUMED_HOST_CPU = ID_VM_RES_MGMT_HOSTCPU
         + "/consumedResourceLbl";
   public static final String ID_VM_RES_MGMT_ACTIVE_HOST_CPU = ID_VM_RES_MGMT_HOSTCPU
         + "/activeOverheadResourceLbl";
   public static final String ID_VM_RES_MGMT_CONSUMED_HOST_MEM = ID_VM_RES_MGMT_HOSTMEM
         + "/consumedResourceLbl";
   public static final String ID_VM_RES_MGMT_ACTIVE_HOST_MEM = ID_VM_RES_MGMT_HOSTMEM
         + "/activeOverheadResourceLbl";
   public static final String ID_VM_RES_MGMT_ACTIVE_GUEST_MEM = ID_VM_RES_MGMT_GUESTMEM
         + "/activeMemoryLbl";
   public static final String ID_VM_RES_MGMT_PRIVATE_GUEST_MEM = ID_VM_RES_MGMT_GUESTMEM
         + "/privateMemoryLbl";
   public static final String ID_VM_RES_MGMT_SHARED_GUEST_MEM = ID_VM_RES_MGMT_GUESTMEM
         + "/sharedMemoryLbl";
   public static final String ID_VM_RES_MGMT_BALOONED_GUEST_MEM =
         ID_VM_RES_MGMT_GUESTMEM + "/balloonedMemoryLbl";
   public static final String ID_VM_RES_MGMT_COMPRESSED_GUEST_MEM =
         ID_VM_RES_MGMT_GUESTMEM + "/compressedMemoryLbl";
   public static final String ID_VM_RES_MGMT_SWAPPED_GUEST_MEM = ID_VM_RES_MGMT_GUESTMEM
         + "/swappedMemoryLbl";
   public static final String ID_VM_RES_MGMT_UNACCESSED_GUEST_MEM =
         ID_VM_RES_MGMT_GUESTMEM + "/unaccessedMemoryLbl";

   // Reporting Automation ID
   // TODO: replace the automation name once the ID is available
   public static final String ID_REPORTING_REPORTS_TAB = "automationName=Reports";
   public static final String ID_REPORTING_ADMINISTRATION_TAB =
         "automationName=Administration";
   public static final String ID_REPORTING_CREATEREPORT_BUTTON = "createReportButton";
   public static final String ID_REPORTING_RUNREPORT_BUTTON = "runReportButton";
   public static final String ID_REPORTING_DELETEREPORT_BUTTON = "deleteReportButton";
   public static final String ID_REPORTING_DELETEDOCUMENTBUTTON = "deleteDocumentButton";
   public static final String ID_REPORTING_CREATEREPORT_DIALOG = "tiwoDialog";
   public static final String ID_REPORTING_CREATEREPORT_DIALOG_CREATE_BUTTON =
         "tiwoDialog/saveOrCreateButton";
   public static final String ID_REPORTING_CREATEREPORT_DIALOG_CANCEL_BUTTON =
         "tiwoDialog/cancelButton";
   public static final String ID_REPORTING_CREATEREPORT_PLUS_BUTTON = "createButton";
   public static final String ID_TIME_STAMP_LABEL = "timestampLabel";
   public static final String ID_REPORTING_CREATEREPORT_DIALOG_REPORTNAME_TEXTBOX =
         "metadesc_textbox_nameBox";
   public static final String ID_REPORTING_CREATEREPORT_DIALOG_FORMATLIST_COMBOBOX =
         "metadesc_menulist_formatList";
   public static final String ID_REPORTING_CREATEREPORT_DIALOG__FIELD_LISTBOX =
         "metadesc_listbox_null";
   public static final String ID_REPORTING_CREATEDREPORT_LIST = "objectList";
   public static final String ID_REPORTING_DELETEREPORT_ALERT_YES_BUTTON =
         "automationName=Yes";
   // for template popup in the create report dialog
   public static final String ID_REPORTING_CREATEREPORT_DIALOG_TEMPLATE_POPUP =
         "templateMenu";
   public static final String ID_REPORTING_CREATEREPORT_DIALOG_TARGETOBJECT_TEXTBOX =
         "metadesc_textbox_targetBox";

   // New Resource Pool
   public static final String ID_RP_CPU_SHARES_LABEL =
         "cpuConfigControl/shares/sharesControl/levels";
   public static final String ID_RP_CPU_RESERVATION_LABEL =
         "cpuConfigControl/reservation/reservationCombo/topRow/comboBox";
   public static final String ID_RP_CPU_RESERVATION_TYPE_LABEL =
         "cpuConfigControl/expandableReservation/expReservationCheck";
   public static final String ID_RP_CPU_LIMIT_LABEL =
         "cpuConfigControl/limit/limitCombo/topRow/comboBox";
   public static final String ID_RP_CPU_SHARES_NUMERIC_STEPPER =
         "cpuConfigControl/shares/sharesControl/numShares";
   public static final String ID_RP_CPU_RESERVATION_UNITS =
         "cpuConfigControl/reservation/reservationCombo/topRow/units";
   public static final String ID_RP_CPU_LIMIT_UNITS =
         "cpuConfigControl/limit/limitCombo/topRow/units";
   public static final String ID_RP_CPU_MAX_RESERVATION_LABEL =
         "cpuConfigControl/reservation/reservationCombo/bottomRow/label2";
   public static final String ID_RP_CPU_MAX_LIMIT_LABEL =
         "cpuConfigControl/limit/limitCombo/bottomRow/label2";
   public static final String ID_RP_CPU_RESERVATION_CHECKBOX =
         "cpuConfigControl/expandableReservation/expReservationCheck";
   public static final String ID_RP_MEMORY_SHARES_LABEL =
         "memoryConfigControl/shares/sharesControl/levels";
   public static final String ID_RP_MEMORY_RESERVATION_LABEL =
         "memoryConfigControl/reservation/reservationCombo/topRow/comboBox";
   public static final String ID_RP_MEMORY_RESERVATION_TYPE_LABEL =
         "memoryConfigControl/expandableReservation/expReservationCheck";
   public static final String ID_RP_MEMORY_LIMIT_LABEL =
         "memoryConfigControl/limit/limitCombo/topRow/comboBox";
   public static final String ID_RP_MEMORY_SHARES_NUMERIC_STEPPER =
         "memoryConfigControl/shares/sharesControl/numShares";
   public static final String ID_RP_MEMORY_RESERVATION_UNITS =
         "memoryConfigControl/reservation/reservationCombo/topRow/units";
   public static final String ID_RP_MEMORY_LIMIT_UNITS =
         "memoryConfigControl/limit/limitCombo/topRow/units";
   public static final String ID_RP_MEMORY_MAX_RESERVATION_LABEL =
         "memoryConfigControl/reservation/reservationCombo/bottomRow/label2";
   public static final String ID_RP_MEMORY_MAX_LIMIT_LABEL =
         "memoryConfigControl/limit/limitCombo/bottomRow/label2";
   public static final String ID_RP_MEMORY_RESERVATION_CHECKBOX =
         "memoryConfigControl/expandableReservation/expReservationCheck";
   public static final String ID_RP_OK_BUTTON = "tiwoDialog/mainPane/okButton";
   public static final String ID_RP_CANCEL_BUTTON = "tiwoDialog/mainPane/cancelButton";
   public static final String ID_RP_NAME_TEXTINPUT =
         "tiwoDialog/mainPane/editSettingsSE/nameStackBlok/nameTextInput";
   public static final String ID_RP_PANEL_CPU_LABEL = "titleLabel.cpuStackBlock";
   public static final String ID_RP_PANEL_MEM_LABEL = "titleLabel.memoryStackBlock";
   public static final String ID_RP_ARROW_CPU_STACK = "arrowImagecpuStackBlock";
   public static final String ID_RP_ARROW_MEMORY_STACK = "arrowImagememoryStackBlock";
   public static final String ID_IMAGE_COLLAPSE_RP_SUMMARY_RESOURCE_SETTINGS_CPU =
         "arrowImage_ResourcePoolSettingsView_StackBlock1";
   public static final String ID_IMAGE_COLLAPSE_RP_SUMMARY_RESOURCE_SETTINGS_MEMORY =
         "arrowImage_ResourcePoolSettingsView_StackBlock2";
   public static final String ID_BUTTON_RP_GETTING_STARTED_EDIT =
         "vsphere.core.resourcePool.editAction";
   public static final String ID_LABEL_CPU_RESOURCE_SHARES =
         "cpuResourceViewControl/shares";
   public static final String ID_LABEL_CPU_RESOURCE_RESERVATION =
         "cpuResourceViewControl/reservation";
   public static final String ID_LABEL_CPU_RESOURCE_LIMIT =
         "cpuResourceViewControl/limit";
   public static final String ID_LABEL_MEMORY_RESOURCE_SHARES =
         "memoryResourceViewControl/shares";
   public static final String ID_LABEL_MEMORY_RESOURCE_RESERVATION =
         "memoryResourceViewControl/reservation";
   public static final String ID_LABEL_MEMORY_RESOURCE_LIMIT =
         "memoryResourceViewControl/limit";
   public static final String ID_PORTLET_RESOURCEPOOL_SUMMARY_RESOURCESETTING =
         "vsphere.core.resourcePool.summary.rsrcSettingsView.chrome";
   public static final String ID_PORTLET_RESOURCEPOOL_SUMMARY_RESOURCECONSUMERS =
         "vsphere.core.resourcePool.summary.rsrcConsumersView.chrome";

   // Vm Edit Resource Setting
   public static final String ID_ACTION_EDIT_RESOURCE_SETTING =
         "vsphere.core.vm.provisioning.editResourceAction";
   public static final String ID_COMBOBOX_VM_EDIT_RESOURCE_SETTING_CPU_SHARES =
         "cpuShares/cpuSharesControl/levels";
   public static final String ID_COMBOBOX_VM_EDIT_RESOURCE_SETTING_CPU_RESERVATION =
         "cpuReservation/comboBox";
   public static final String ID_COMBOBOX_VM_EDIT_RESOURCE_SETTING_CPU_LIMIT =
         "cpuLimit/comboBox";
   public static final String ID_COMBOBOX_VM_EDIT_RESOURCE_SETTING_MEMORY_SHARES =
         "memoryShares/memorySharesControl/levels";
   public static final String ID_COMBOBOX_VM_EDIT_RESOURCE_SETTING_MEMORY_RESERVATION =
         "memoryReservation/comboBox";
   public static final String ID_COMBOBOX_VM_EDIT_RESOURCE_SETTING_MEMORY_LIMIT =
         "memoryLimit/comboBox";
   public static final String ID_BUTTON_EDIT_RESOURCE_SETTING_SCHEDULE =
         ID_ACTION_EDIT_RESOURCE_SETTING + "/button";
   public static final String ID_BUTTON_VM_SCHEDULE_NEW_TASK =
         "contextObjectScheduledActions/button";
   public static final String ID_COMBOBOX_SCHEDULED_TASK_VM_EDIT_RESOURCE_SETTING_CPU_SHARES =
         "cpuShares/cpuSharesControl/levels";
   public static final String ID_COMBOBOX_SCHEDULED_TASK_VM_EDIT_RESOURCE_SETTING_MEMORY_SHARES =
         ID_COMBOBOX_VM_EDIT_RESOURCE_SETTING_MEMORY_SHARES + "/levels";;
   public static final String ID_MENU_SCHEDULE_EDIT_VM_RESOURCE_SETTING =
         "automationName=Edit Resource Settings...";
   public static final String ID_LABEL_VM_HOST_CPU_WORST_CASE_ALLOCATION =
         ID_VM_HOST_CPU_PORTLET + "/worstCaseAllocationLabel";
   public static final String ID_LABEL_VM_HOST_MEMORY_WORST_CASE_ALLOCATION =
         ID_VM_HOST_MEMORY_PORTLET + "/worstCaseAllocationLabel";
   public static final String ID_LABEL_VM_HOST_CONFIGURED = ID_VM_HOST_MEMORY_PORTLET
         + "/configuredMemLabel";

   public static final String ID_FAULTS_LIST_CANCEL_OK_BUTTON = "cancelOKButton";
   public static final String ID_PORTLET_VM_MONITOR_RESOURCEALLOCATION_HOST_MEMORY =
         "vsphere.core.vm.monitor.resMgmt.hostMemoryView.chrome";
   public static final String ID_PORTLET_VM_MONITOR_RESOURCEALLOCATION_HOST_CPU =
         "vsphere.core.vm.monitor.resMgmt.hostCPUView.chrome";
   public static final String ID_PORTLET_VM_MONITOR_RESOURCEALLOCATION_GUEST_MEMORY =
         "vsphere.core.vm.monitor.resMgmt.guestMemoryView.chrome";
   // Template ResourceManagement
   public static final String ID_RESOURCE_MGT_CPU_PANEL =
         "vsphere.core.vm.hostCPUView.chrome";
   public static final String ID_RESOURCE_MGT_MEMORY_PANEL =
         "vsphere.core.vm.hostMemoryView.chrome";
   public static final String ID_RESOURCE_MGT_GUEST_MEMORY_PANEL =
         "vsphere.core.vm.guestMemoryView.chrome";
   public static final String ID_RESOURCE_MGT_CPU_CONSUMED =
         "vsphere.core.vm.hostCPUView/resConsumeGridView/consumedResourceLbl";
   public static final String ID_RESOURCE_MGT_CPU_ACTIVE =
         "vsphere.core.vm.hostCPUView/resConsumeGridView/activeOverheadResourceLbl";
   public static final String ID_RESOURCE_MGT_MEMORY_CONSUMED =
         "vsphere.core.vm.hostMemoryView/resConsumeGridView/consumedResourceLbl";
   public static final String ID_RESOURCE_MGT_MEMORY_ACTIVE =
         "vsphere.core.vm.hostMemoryView/resConsumeGridView/activeOverheadResourceLbl";
   public static final String ID_RESOURCE_MGT_GUEST_MEMORY =
         "vsphere.core.vm.guestMemoryView/memoryGridView/activeMemoryLbl";
   public static final String ID_RESOURCE_MGT_GUEST_PRIVATE =
         "vsphere.core.vm.guestMemoryView/memoryGridView/privateMemoryLbl";
   public static final String ID_RESOURCE_MGT_GUEST_SHARED =
         "vsphere.core.vm.guestMemoryView/memoryGridView/sharedMemoryLbl";
   public static final String ID_RESOURCE_MGT_GUEST_BALLOONED =
         "vsphere.core.vm.guestMemoryView/memoryGridView/balloonedMemoryLbl";
   public static final String ID_RESOURCE_MGT_GUEST_COMPRESSED =
         "vsphere.core.vm.guestMemoryView/memoryGridView/compressedMemoryLbl";
   public static final String ID_RESOURCE_MGT_SWAPPED_MEMORY =
         "vsphere.core.vm.guestMemoryView/memoryGridView/swappedMemoryLbl";
   public static final String ID_RESOURCE_MGT_GUEST_UNACCESSED =
         "vsphere.core.vm.guestMemoryView/memoryGridView/unaccessedMemoryLbl";

   // Licensing Reporting Ids
   public static final String ID_LIC_LOADING_DATA_PROGRESSBAR =
         "animatedLoadingProgressBar";
   public static final String ID_LIC_REPORTING_VIEW =
         "vsphere.license.licenseReportView";
   public static final String ID_LIC_REPORTING = "vsphere.license.reporting"; // Added by wgao, replace "ID_LIC_REPORTING_VIEW"
   public static final String ID_LIC_REPORTING_WARNING_BOX = "warningBox";
   public static final String ID_LIC_VC_SELECTOR = "vcSelector";
   public static final String ID_LIC_VC_COMOBOX = "vcComboBox";
   //ublic static final String ID_LIC_SERVER_COMOBOX= "serverComboBox"; // For vsphere5x branch, added by wgao
   public static final String ID_LIC_SINGLEVC_CHECKBOX = "singleVcCheckBox"; // For "Show data only for selected VC, added by wgao"
   public static final String ID_LIC_TIME_PERIOD_COMOBOX = "selectTimeComboBox";
   public static final String ID_LIC_CUSTOM_PERIOD_RECALCULATE_BUTTON =
         "buttonRecalculate";
   public static final String ID_LIC_DATE_SELECTOR = "dateSelector";
   public static final String ID_LIC_CUSTOM_PERIOD_RANGEFIELD = "dateRangeField";
   public static final String ID_LIC_CUSTOM_PERIOD_FROM_DATEFIELD =
         ID_LIC_CUSTOM_PERIOD_RANGEFIELD + "/dateFieldFrom";
   public static final String ID_LIC_CUSTOM_PERIOD_TO_DATEFIELD =
         ID_LIC_CUSTOM_PERIOD_RANGEFIELD + "/dateFieldTo";
   public static final String ID_LIC_BASE_BARCHART = "barChart";
   public static final String ID_LIC_BARCHART = ID_LIC_BASE_BARCHART + "/barChart";
   public static final String ID_LIC_BARCHART_VSCROLL = ID_LIC_BASE_BARCHART
         + "/verticalScrollBar";
   public static final String ID_LIC_PRODUCTASSETVIEWSTACK = "productAssetViewStack";
   public static final String ID_LIC_ERROR_LABEL = ID_LIC_PRODUCTASSETVIEWSTACK
         + "/issueLabel";
   public static final String ID_LIC_GENERATED_AT_LABEL = ID_LIC_PRODUCTASSETVIEWSTACK
         + "/generatedLabel";
   public static final String ID_LIC_LICENSE_USAGE_ASSETS_PANEL =
         ID_LIC_PRODUCTASSETVIEWSTACK + "/assetView";
   //public static final String ID_LIC_PRODUCT_DETAILS_PANEL = ID_LIC_REPORTING_VIEW
   //    + "/productDetailsPanel";
   public static final String ID_LIC_PRODUCT_DETAILS_PANEL = ID_LIC_REPORTING
         + "/productDetailsPanel"; //Replace "ID_LIC_REPORTING_VIEW" with "ID_LIC_REPORTING" for MN.Next web client
   public static final String ID_LIC_PRODUCT_DETAILS = ID_LIC_REPORTING
         + "/productDetails"; //Add for MN.Next web client
   public static final String ID_ASSETS_VIEW_ASSET_LIST = ID_LIC_PRODUCTASSETVIEWSTACK
         + "/assetView/assetList";
   public static final String ID_LIC_PD_LICENSE_LIST = ID_LIC_PRODUCT_DETAILS_PANEL
         + "/licenseList";
   public static final String ID_LIC_PD_PRODUCT_NAME_LABEL =
         ID_LIC_PRODUCT_DETAILS_PANEL + "/productLabel";
   public static final String ID_LIC_PD_AVERAGE_USAGE = ID_LIC_PRODUCT_DETAILS_PANEL
         + "/avgUsageLabel";
   public static final String ID_LIC_PD_AVERAGE_CAPACITY = ID_LIC_PRODUCT_DETAILS_PANEL
         + "/avgCapacityLabel";
   public static final String ID_LIC_PD_AVERAGE_USAGE_PERCENT =
         ID_LIC_PRODUCT_DETAILS_PANEL + "/avgPercentUsageLabel";
   public static final String ID_LIC_PD_CURRENT_USAGE = ID_LIC_PRODUCT_DETAILS_PANEL
         + "/currentUsageLabel";
   public static final String ID_LIC_PD_CURRENT_CAPACITY = ID_LIC_PRODUCT_DETAILS_PANEL
         + "/currentCapabityLabel";
   public static final String ID_LIC_PD_CURRENT_USAGE_PERCENT =
         ID_LIC_PRODUCT_DETAILS_PANEL + "/currentPercentUsageLabel";
   public static final String ID_LIC_PD_THRESHOLD = ID_LIC_PRODUCT_DETAILS_PANEL
         + "/thresholdLabel";
   public static final String ID_LIC_PD_EDIT_THRESHOLD_LINKBUTTON =
         ID_LIC_PRODUCT_DETAILS_PANEL + "/editThresholdLinkButton";
   public static final String ID_LIC_PD_WARNING_ICON = ID_LIC_PRODUCT_DETAILS_PANEL
         + "/_ProductDetailsPanel_Image1";
   public static final String ID_LIC_PD_FORM_HEADING = ID_LIC_PRODUCT_DETAILS_PANEL
         + "/selectDetailsHeading";
   public static final String ID_LIC_STATIC_WARNING_LABEL = ID_LIC_REPORTING_WARNING_BOX
         + "/staticWarningLabel";
   public static final String ID_LIC_WARNING_LINK_BUTTON = ID_LIC_REPORTING_WARNING_BOX
         + "/warningLinkButton0";

   public static final String ID_LIC_EXPORT_LINKBUTTON = "exportButton";
   public static final String ID_LIC_EXPORT_DIALOG = "exportDialog";
   public static final String ID_LIC_EXPORT_DIALOG_EXPORT_BUTTON = ID_LIC_EXPORT_DIALOG
         + "/buttonExport";
   public static final String ID_LIC_EXPORT_DIALOG_CANCEL_BUTTON = ID_LIC_EXPORT_DIALOG
         + "/buttonCancel";
   // TODO: Using the title to identify the dialog as a workaround because
   // FlexAlert dialogs don't have ids
   public static final String ID_LIC_SAVE_EXPORT_DIALOG =
         TestConstantsKey.TITLE_PROPERTY + TestConstantsKey.EQUALS_SIGN
               + "Save Export File";
   // TODO: Using the button label as a workaround because FlexAlert buttons
   // don't have ids
   public static final String ID_LIC_SAVE_EXPORT_YES = TestConstantsKey.LABEL_PROPERTY
         + TestConstantsKey.EQUALS_SIGN + TestConstantsKey.STRING_YES;
   public static final String ID_LIC_SAVE_EXPORT_NO = TestConstantsKey.LABEL_PROPERTY
         + TestConstantsKey.EQUALS_SIGN + TestConstantsKey.STRING_NO;

   public static final String ID_LIC_EDIT_THRESHOLD_DIALOG = "editThresholdDialog";
   public static final String ID_LIC_EDIT_THRESHOLD_DIALOG_VC_LABEL =
         ID_LIC_EDIT_THRESHOLD_DIALOG + "/vcLabel";
   public static final String ID_LIC_EDIT_THRESHOLD_DIALOG_PRODUCT_LABEL =
         ID_LIC_EDIT_THRESHOLD_DIALOG + "/productLabel";
   public static final String ID_LIC_EDIT_THRESHOLD_DIALOG_THRESHOLD_NUM_STEPPER =
         ID_LIC_EDIT_THRESHOLD_DIALOG + "/thresholdNumericStepper";
   public static final String ID_LIC_EDIT_THRESHOLD_DIALOG_OK_BUTTON =
         ID_LIC_EDIT_THRESHOLD_DIALOG + "/buttonOK";
   public static final String ID_LIC_EDIT_THRESHOLD_DIALOG_CANCEL_BUTTON =
         ID_LIC_EDIT_THRESHOLD_DIALOG + "/buttonCancel";
   public static final String ID_LIC_EDIT_THRESHOLD_DIALOG_NOTE_LABEL =
         ID_LIC_EDIT_THRESHOLD_DIALOG + "/multiVcAssetsText";
   // TODO: Using the title to identify the dialog as a workaround because
   // FlexAlert dialogs don't have ids
   public static final String ID_LIC_ADD_THRESHOLD_ALERT =
         TestConstantsKey.TITLE_PROPERTY + TestConstantsKey.EQUALS_SIGN
               + TestConstantsKey.LIC_TITLE_ADD_THRESHOLD_ALERT;
   public static final String ID_LIC_MODIFY_THRESHOLD_ALERT =
         TestConstantsKey.TITLE_PROPERTY + TestConstantsKey.EQUALS_SIGN
               + TestConstantsKey.LIC_TITLE_MODIFY_THRESHOLD_ALERT;
   public static final String ID_LIC_REMOVE_THRESHOLD_ALERT =
         TestConstantsKey.TITLE_PROPERTY + TestConstantsKey.EQUALS_SIGN
               + TestConstantsKey.LIC_TITLE_REMOVE_THRESHOLD_ALERT;
   // TODO: Using the button label as a workaround because FlexAlert buttons
   // don't have ids
   public static final String ID_LIC_FLEX_ALERT_OK_BUTTON =
         TestConstantsKey.LABEL_PROPERTY + TestConstantsKey.EQUALS_SIGN
               + TestConstantsKey.STRING_OK;

   public static final String ID_LIC_INFO_LABEL = ID_LIC_PRODUCTASSETVIEWSTACK
         + "/infoLabel";
   public static final String ID_ASSETS_VIEW_BREADCRUMB_CONTROL = ID_LIC_REPORTING
         + "/breadcrumbControl"; // Replace "ID_LIC_REPORTING_VIEW" with "ID_LIC_REPORTING" for MN.Next web client
   public static final String ID_ASSETS_VIEW_VC_LINK_BUTTON =
         ID_ASSETS_VIEW_BREADCRUMB_CONTROL + "/button0";
   public static final String ID_ASSETS_VIEW_LICENSE_KEY_LINK_BUTTON =
         ID_ASSETS_VIEW_BREADCRUMB_CONTROL + "/button2";
   public static final String ID_LIC_PD_AVERAGE_LINK_BUTTON = "usageAggregation";
   public static final String ID_LIC_PD_AVERAGE_EXPLANATION_TITLE_WINDOW =
         "detailsPopup";

   public static final String ID_ACTION_CREATE_VAPP = new StringBuffer(EXTENSION_PREFIX)
         .append("vApp.createAction").toString();

   public static final String ID_ACTION_ANNOTATION = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append("editAnnotationAction").toString();
   // UIComponent Property Names
   public static final String ID_UI_COMPONENT_PROPERTY_X = "x";
   public static final String ID_UI_COMPONENT_PROPERTY_Y = "y";
   public static final String ID_UI_COMPONENT_PROPERTY_TEXT = "text";
   public static final String ID_UI_COMPONENT_PROPERTY_VALUE = "value";
   public static final String ID_UI_COMPONENT_PROPERTY_SELECTED = "selected";
   public static final String ID_UI_COMPONENT_PROPERTY_SELECTED_ITEM = "selectedItem";
   public static final String ID_UI_COMPONENT_PROPERTY_ENABLED = "enabled";
   public static final String ID_UI_COMPONENT_PROPERTY_TABENABLED = "tabEnabled";


   // Hardware
   public static final String MHz = "MHz";
   public static final String MB = "MB";

   // Task Tab Data Provider Properties
   public final static String ID_TASK_TAB_DATA_PROVIDER_COMPLETE_TIME =
         "completedStartTime.time";
   public final static String ID_TASK_TAB_DATA_PROVIDER_START_TIME =
         "actualStartTime.time";
   public final static String ID_TASK_TAB_DATA_PROVIDER_QUEUED_TIME =
         "requestedStartTime.time";
   public final static String ID_TASK_TAB_DATA_PROVIDER_TARGET_NAME = "targetName";
   public final static String ID_TASK_TAB_DATA_PROVIDER_TASK_STATUS = "status";
   public final static String ID_TASK_TAB_DATA_PROVIDER_TASK_STATUS_COMPLETE =
         "Completed";
   public final static String ID_TASK_TAB_DATA_PROVIDER_TASK_STATUS_FAILED = "Failed";
   public final static String ID_TASK_TAB_DATA_PROVIDER_LIST_LENGTH =
         "dataProvider.list.length";
   public final static String ID_TASK_TAB_DATA_PROVIDER_LIST_SOURCE =
         "dataProvider.list.source.";
   public final static String ID_TASK_TAB_CLUSTER = "taskGrid";

   // Event Tab Data Provider Properties
   public final static String ID_EVENT_TAB_CLUSTER = "eventsDataGrid";

   // Register VM Wizard
   public static final String ID_REGISTERVM_CANCEL_BUTTON = "buttonCancel";
   public static final String ID_REGISTERVM_DATAFIELD_NAME = "name";
   public static final String ID_REGISTERVM_DATAFIELD_VALUE = "value";
   public static final String ID_REGISTERVM_NAMELOCATION_PAGE = "nameLocationPage";
   public static final String ID_REGISTERVM_RP_PAGE = "resourcePoolPage";
   public static final String ID_REGISTERVM_HC_PAGE = "hostClusterPage";
   public static final String ID_REGISTERVM_SPECIFIC_HOST_PAGE = "specificHostPage";
   public static final String ID_REGISTERVM_READY_TO_COMPLETE = "summaryPage";
   public static final String ID_CLUSTER_HOSTS_LIST = "clusterHostList";
   public static final String ID_DATASTORE_BROWSER_FILELIST = "fileList";
   public static final String ID_REGISTERVM_WIZARD_OK_BUTTON = "buttonOk";
   public static final String ID_ACTION_REGISTER_VM = EXTENSION_PREFIX
         + EXTENSION_ENTITY_DATASTORE + "registerVmAction";
   public static final String ID_SEARCH_BUTTON_NAME_LOCATION =
         "tiwoDialog/nameLocationPage/folderSelectWidget/searchButton";
   public static final String ID_SEARCH_TEXT_INPUT_NAME_LOCATION =
         "folderSelectWidget/searchControl/_SearchControl_HBox4/searchInput";
   public static final String ID_SEARCH_BUTTON_HOST_CLUSTER =
         "hostSelectWidget/searchControl/searchStack/searchButtonContainer/searchButton";
   public static final String ID_SEARCH_TEXT_INPUT_HOST_CLUSTER =
         "hostSelectWidget/searchControl/_SearchControl_HBox4/searchInput";
   public static final String ID_SEARCH_BUTTON_RESOURCE_POOL =
         "poolSelectWidget/searchControl/searchStack/searchButtonContainer/searchButton";
   public static final String ID_SEARCH_TEXT_INPUT_RESOURCE_POOL =
         "poolSelectWidget/searchControl/_SearchControl_HBox4/searchInput";
   public static final String ID_REGISTER_VM_VALIDATION_CONTAINER =
         "validationMessageContainer";
   public static final String ID_REGISTER_VM_CONTAINER = "_container";
   public static final String ID_REGISTER_VM_MESSAGE = "_message";
   public static final String ID_REGISTER_VM_SEARCH_NO_RESULTS =
         "contentStack/searchResultsTable/messageBox/messageLabel";
   public static final String ID_REGISTER_VM_NO_PRIVILEGES_WARNING = "errWarningText";
   public static final String ID_REGISTER_VM_DATACENTER_PRIVILEGES =
         "You do not have privilege \"Register\" on the selected Datacenter.";
   public static final String ID_BUTTON_CLUSTERREGISTER = "automationName=Clusters";
   public static final String ID_NAME_LOCATION_PAGE_NAV_TREE =
         "tiwoDialog/nameLocationPage/navTreeView/navTree";
   public static final String ID_DATASTORE_NAME_LOCATION_PAGE_NAV_TREE =
         "tiwoDialog/datastoreNameAndLocationPage/navTreeView/navTree";
   public static final String ID_HOST_CLUSTER_PAGE_NAV_TREE =
         "tiwoDialog/hostClusterPage/navTreeView/navTree";
   public static final String ID_ERROR_WARNING_LABEL = "errWarningText";

   // Error, Warning, Confirmation Dialog & overlay error IDs
   public static final String ID_ERROR_DIALOG = "ErrorDialog";
   public static final String ID_ERROR_DIALOG_YESNO = "YesNoErrorDialog";
   public static final String ID_ERROR_OK_LABEL = "label=OK";
   public static final String ID_WARNING_DIALOG = "WarningDialog";
   public static final String ID_WARNING_OK_LABEL = "label=OK";
   public static final String ID_CONFIRMATION_DIALOG = "confirmationDialog";
   public static final String ID_YES_NO_DIALOG = "YesNoDialog";
   public static final String ID_AUTOMATION_NAME_YES = "automationName=Yes";
   public static final String ID_AUTOMATION_NAME_NO = "automationName=No";
   public static final String ID_AUTOMATION_NAME_OK = "automationName=OK";
   public static final String ID_AUTOMATION_NAME_CANCEL = "automationName=Cancel";
   public static final String ID_YES_BUTTON = ID_YES_NO_DIALOG + "/"
         + ID_AUTOMATION_NAME_YES;
   public static final String ID_NO_BUTTON = ID_YES_NO_DIALOG + "/"
         + ID_AUTOMATION_NAME_NO;
   public static final String ID_ERROR_WARNING_OK = "OK";
   public static final String ID_CONFIRM_YES_BUTTON = ID_CONFIRMATION_DIALOG + "/"
         + ID_AUTOMATION_NAME_YES;
   public static final String ID_CONFIRM_NO_BUTTON = ID_CONFIRMATION_DIALOG + "/"
         + ID_AUTOMATION_NAME_NO;
   public static final String ID_CONFIRM_YES_LABEL = "label=Yes";
   public static final String ID_CONFIRM_NO_LABEL = "label=No";
   public static final String ID_OVERLAY_ERROR_CONTAINER = "errorContainer";
   public static final String ID_INFO_DIALOG = "YesNoDialog";

   // file a bug to change the ID to
   // public static final String ID_OVERLAY_ERROR_LABEL_CONTAINER =
   // "error_1_box";
   public static final String ID_OVERLAY_ERROR_LABEL = "error_";
   public static final String LABEL_ERROR_0 = "error_0";
   public static final String LABEL_LOGIN_ERROR = "error_text";

   public static final String ID_ADV_SEARCH_NO_RESULTS =
         "advancedSearch/searchResults/messageBox/messageLabel";
   public static final String ID_LABEL_SIMPLE_SEARCH_MESSAGELABEL = "/messageLabel";

   // Advanced DataGrid
   public static final String ID_HIDE_COLUMN_NAME = "1";
   public static final String ID_SIZE_TO_FIT_COLUMN_NAME = "2";
   public static final String ID_SIZE_TO_FIT_ALL_COLUMNS = "3";
   public static final String ID_LOCK_FIRST_COLUMN = "4";
   public static final String ID_SHOW_HIDE_COLUMNS = "5";
   public static final String ID_SELECT_ALL_BUTTON = "selectAllLinkButton";
   public static final String ID_SHOW_HIDE_COLUMNS_ADVDATAGRID = "list";
   public static final String ID_SHOW_HIDE_COLUMNS_CLOSE_BUTTON = "closeButton";
   public static final String ID_DATACENTER_STORAGE_TAB =
         "vsphere.core.datacenter.datastoresView.container";
   public static final String ID_F2_TEXT_INPUT = "uicEditor";
   public static final String ID_VC_STORAGE_TAB =
         "vsphere.core.folder.datastoresView.container";
   public static final String ID_ADVDATAGRID_CONTEXT_MENU = "headerContextMenu";

   public static final String ID_MENU_DATAGRID_CONTEXT_HEADER = "contextHeader";
   public static final String ID_DATAGRID_CONTEXT_HEADER = "contextHeader/list";
   public static final String ID_DATAGRID_CONTEXT_HEADER_SCROLLER =
         ID_DATAGRID_CONTEXT_HEADER + "/className=VScrollBar";
   public static final String ID_DATAGRID_CONTEXT_HEADER_SCROLLER_DOWN_BUTTON =
         ID_DATAGRID_CONTEXT_HEADER + "/name=downArrowSkin";
   public static final String ID_DATAGRID_CONTEXT_HEADER_SCROLL_POSITION =
         "scrollPosition";
   public static final String ID_DATA_STORE_CLUSTER_HOSTS_SHOW_HIDE_COLUMNS_ADVDATAGRID =
         "standaloneHostsForDatastoreCluster";
   public static final String ID_DATA_STORE_CLUSTER_CLUSTER_SHOW_HIDE_COLUMNS_ADVDATAGRID =
         "clustersForDatastoreCluster";
   public static final String ID_DATA_STORE_CLUSTER_DATASTORES_SHOW_HIDE_COLUMNS_ADVDATAGRID =
         "dssForDatastoreCluster";
   public static final String ID_DATA_STORE_CLUSTER_VM_SHOW_HIDE_COLUMNS_ADVDATAGRID =
         "vmsForDatastoreCluster";

   // Startup welcome page dialog
   public static final String ID_STARTUP_WELCOME_PAGE_DIALOG = "startupPage";
   public static final String ID_DO_NOT_SHOW_WELCOME_PAGE_ON_STARTUP_CHECKBOX =
         "doNotShowOnStartup";
   public static final String ID_CLOSE_STARTUP_WELCOME_PAGE_IMAGE = "closeStartupPage";

   // Edit VM Startup and shutdown dialog
   public static final String ID_EDIT_VM_STARTUP_SHUTDOWN_CHECKBOX =
         "cbStartAutomatically";
   public static final String ID_STARTUP_DELAY_FLEXNUMERICSTEPPER = "startupDelayText";
   public static final String ID_SHUTDOWN_ACTION_COMBOBOX = "cbShutdownAction";
   public static final String ID_SHUTDOWN_DELAY_FLEXNUMERICSTEPPER = "shutdownDelayText";
   public static final String ID_CONTINUE_IMMEDIATELY_CHECKBOX = "cbContinueImmediately";

   public static final String ID_DATASTORE = "vmxLocation";
   public static final String ID_MORE_RECOMMENDATION_LINK = "recommendationLink";
   public static final String ID_STORAGE_DRS_CHECKBOX = "sdrsChkBox";
   public static final String ID_STORAGE_LIST = "storageList";
   public static final String ID_DATASTORE_LIST = "datastoreList";
   public static final String ID_VC_DATASTORE_LIST = "datastoresForVCenter/list";
   public static final String ID_NEW_DEVICE = "popUpButton";
   public static final String ID_NEW_DEVICE_MENU = "menu";
   public static final String LABEL_NEW_HARD_DISK = "New Hard Disk";
   public static final String ID_LABEL_EXISTING_HARD_DISK = "Existing Hard Disk";
   public static final String LABEL_NONE = "None";
   public static final String ID_ADD_NEW_DEVICE_BUTTON = "addDeviceButton";
   public static final String LABEL_BROWSE = "Browse...";
   public static final String ID_DATASTORE_LOCATION_COMBOBOX = "location";
   public static final String ID_OK_BUTTON = "buttonOk";
   public static final String ID_CANCEL_BUTTON = "buttonCancel";
   public static final String ID_NEXT_BUTTON = "next";
   public static final String ID_FINISH_BUTTON = "finish";
   public static final String ID_VIEW_FAULTS_TAB = "_DatastoreRecommendationDialog_Box2";
   public static final String ID_VIEW_RECOMMENDATIONS_TAB =
         "_DatastoreRecommendationDialog_Box1";
   public static final String ID_TOOLBAR_TABS = "tabs";
   public static final String ID_CANCEL_DOWNLOAD = "cancelDownload";
   public static final String ID_DOWNLOAD_COMPLETE = "downloadCompleteLabel";

   public static final String ID_ADVANCED_BASIC_BUTTON = "modeBtn";
   public static final String ID_CUSTOMIZE_HARDWARE_CHECKBOX =
         "customizeHardwareCheckBox";
   public static final String LABEL_ADVANCED = "Advanced";
   public static final String LABEL_BASIC = "Basic";

   // Pod summary constants
   public static final String POD_SERVICES_PORTLET =
         "vsphere.core.dscluster.summary.servicesView.chrome";
   public static final String POD_DATASTORE_CLUSTER_CONSUMERS =
         "vsphere.core.dscluster.summary.consumersView.chrome";
   public static final String POD_DATASTORE_CLUSTER_RESOURCES =
         "vsphere.core.dscluster.summary.resourcesView.chrome";
   public static final String POD_SERVICES_VMWARE_SDRS = "arrowImagedrsBlock";
   public static final String POD_IO_METRICS = "iometrics";
   public static final String POD_IO_AUTOMATION_LEVEL = "automationLevel";
   public static final String POD_SPACE_THRESHOLD = "spaceThreshold";
   public static final String POD_IO_THRESHOLD = "ioThreshold";
   public static final String POD_VIRTUALMACHINES_COUNT = "spVmCount";
   public static final String POD_USED_SPACE = "spProvisioned";
   public static final String POD_FREE_SPACE = "spFreeSpace";
   public static final String POD_TOTAL_CAPACITY = "spCapacity";
   public static final String POD_TOTAL_DATASTORES = "spDatatores";
   public static final String ID_LABEL_SNAPSHOTS_NUMBER_FOR_DATASTORE_CLUSTER =
         "summary_snapshotCount_valueLbl";


   // Constant used in drm.cluster.powerops
   public static final String ID_RECOMMENDATION_TEXT = "recommendationText";

   public static final String ID_EDIT_SETTING_VM_BUTTON =
         "automationName=Getting Started/" + "vsphere.core.vm.provisioning.editAction";
   public static final String ID_EDIT_SETTING_VM_BUTTON_SUMMARY_TAB =
         "automationName=Summary/editHardwareLink";
   public static final String ID_PARALLEL_PORT_DATASTORE_TREE = "tree";
   public static final String ID_SAVE_FILE_AS_TEXTFIELD = "fileNameInput";
   public static final String ID_CONTENTS_LIST = "dbListView/" + "fileList";
   public static final String ID_BROWSE_BUTTON = "browse";
   public static final String ID_SUMMARY_EDITLINK = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY_VM).append("provisioning.editAction").toString();
   public static final String ID_SUMMARY_CLICK = "actionsBlock";
   public static final String ID_SUMMARY_EDITLINK_PARRENT = new StringBuffer(
         ID_SUMMARY_CLICK).append('/').append(ID_SUMMARY_EDITLINK).toString();

   // iSCSI Storage constants
   public static final String ID_STORAGE_ADAPTERS = "tocListStorage";
   public static final String ID_STORAGE_ADAPTERS_GRID =
         "vsphere.core.host.configuration.storageAdaptersList/list";
   public static final String STORAGE_ADAPTERS_GROUPING_PROP =
         "dataProvider.source.grouping.fields.0.name";
   public static final String ID_BUTTON_STORAGE =
         "vsphere.core.host.manage.storage.button";
   public static final String ID_VIEW_PROPERTIES =
         "vsphere.core.storage.adaptersDetails.views.properties";
   public static final String ID_LABEL_ADAPTERSTATUS = "adapterStatusBlockLabel";
   public static final String ID_TIWO_DIALOG_SELECTION = "tiwoDialog/selection";
   public static final String ID_BUTTON_GENERALBLOCKEDIT = "generalBlockEditButton";
   public static final String ID_TEXT_INPUT_ISCSIALIAS = "iScsiAlias";
   public static final String ID_DIALOG_ADD_ISSCI_CONFIRM = "addIScsiConfirmationDialog";
   public static final String ID_DYNAMIC_STATIC_TARGETS_SOURCE = "dataProvider.source.";
   public static final String ID_DYNAMIC_STATIC_TARGETS_ADDRESS = ".address.value";
   public static final String ID_DYNAMIC_STATIC_TARGETS_ISCSI_NAME = ".iScsiName.value";

   // General Storage constants
   public static final String ID_TEXT_STORAGE_VAAI = "vStorageHeader";
   public static final String ID_GRID_STORAGE_VAAI = "vStorageGrid";
   public static final String ID_TEXT_STORAGE_PROVIDER_NAME = "providerName";
   public static final String ID_TEXT_STORAGE_PROVIDER_STATUS = "providerStatus";
   public static final String ID_STORAGE_SETTING_GROUP = "tocListContainer";
   public static final String ID_STORAGE_SETTING_GENERAl = "1";
   public static final String ID_ADD_STORAGE_ADAPTER_MENU = "addStorageAdapterMenu";
   public static final String ID_ADD_STORAGE_ADAPTER_BUTTON = "addAdapterButton";
   public static final String ID_ADAPTER_STATUS_BUTTON = "adapterStatusButton";
   public static final String ID_STORAGE_ADAPTER_PROPERTIES_TAB =
         "vsphere.core.storage.adaptersDetails.views.properties.container";
   public static final String ID_REFRESH_STORAGE_ADAPTERS =
         "vsphere.core.host.refreshStorageSystem";
   public static final String ID_TEXT_DOCUMENT = "document.text";
   public static final String ID_DROPDOWN_PNIC = "pNicDropDown";
   public static final String ID_STEPPER_VLANID = "vlanIdStepper";
   public static final String ID_STEPPER_PRIORITY = "priorityClassStepper";
   public static final String ID_TEXT_MAC_ADDRESS = "macAddressText";
   public static final String ID_TEXT_DEVICE_MODEL = "storage.adapters.deviceModel.text";
   public static final String ID_LABEL_NETWORK_ADAPTER = "networkAdapterLabel";
   public static final String ID_FORM_FCOE_PNIC = "_AddFcoeForm_FormItem1";
   public static final String ID_FORM_FCOE_VLAN_ID = "_AddFcoeForm_FormItem2";
   public static final String ID_FORM_FCOE_PCLASS = "_AddFcoeForm_FormItem3";
   public static final String ID_FORM_FCOE_MAC_ADDRESS = "_AddFcoeForm_FormItem4";
   public static final String ID_TEXT_FCOE_VLAN_ID = "_AddFcoeForm_Text2";
   public static final String ID_TEXT_FCOE_PCLASS = "_AddFcoeForm_Text3";

   public static final String ID_LIST_STORAGE = "tocListOther";
   public static final String ID_TREE_STORAGE =
         "vsphere.core.host.manage.storage/tocTree";
   public static final String ID_TREE_DATASTORE_CLUSTER_SETTINGS =
         "vsphere.core.dscluster.manage.settingsView/tocTree";
   public static final String ID_LIST_HOST_CACHE_CONFIGURATION = "ssdDatastoreList";
   public static final String ID_BUTTON_EDIT = "editButton/button";
   public static final String ID_BUTTON_RESCAN_ADAPTER = "rescanAdapterButton/button";
   public static final String ID_BUTTON_REFRESH = "refreshButton";
   public static final String ID_TIWO_DIALOG_PANEL = "tiwoDialog";
   public static final String ID_CHECKBOX_HOST_CACHE_CONFIGURATION =
         "enabledCacheCheckBox";
   public static final String ID_RADIO_HOST_CACHE_CONFIGURATION_MAX_SPACE =
         "useMaximumAvailableSpaceRadio";
   public static final String ID_RADIO_HOST_CACHE_CONFIGURATION_CUSTOM_SIZE =
         "useCustomSizeRadio";
   public static final String ID_NUMERIC_STEPPER_SSD_SIZE = "swapSizeStepper";
   public static final String ID_SLIDER_SSD_SIZE = "swapSizeSlider";
   public static final String ID_BUTTON_RESCAN_HOST =
         "vsphere.core.host.rescanHost/button";
   public static final String ID_CHECKBOX_HBA = "rescanAllHBAs";
   public static final String ID_CHECKBOX_VMFS = "rescanVMFS";
   public static final String ID_BUTTON_UNMOUNT =
         "vsphere.core.storage.devices.unmount/button";
   public static final String ID_BUTTON_UNMOUNT_YES =
         "automationName=Unmount Device/automationName=Yes";
   public static final String ID_BUTTON_MOUNT =
         "vsphere.core.storage.devices.mount/button";
   public static final String ID_STORAGE_DEVICES_GRID = "list";
   public static final String ID_LABEL_HOST_CACHE_CONFIGURATION_MAX_SPACE =
         "_HostCacheConfigurationForm_VGroup2/labelDisplay";
   public static final String ID_TEXT_HOST_CACHE_CONFIGURATION_FILTER = "textInput";
   public static final String ID_LABEL_HOST_CACHE_CONFIGURATION_SEARCH_INFO =
         "infoPanel1";
   public static final String ID_LABEL_HOST_CACHE_CONFIGURATION_LIST_INDICATOR =
         "emptyListIndicator";
   public static final String ID_LIST_HOST_STORAGE = "hostList";
   public static final String ID_BUTTON_HOST_UNMOUNT = "unmountButton";
   public static final String ID_BUTTON_HOST_MOUNT = "mountButton";
   public static final String ID_BUTTON_REFRESH_HOST_STORAGE_SYSTEM =
         ID_ACTION_REFRESH_ADAPTER + ID_BUTTON;
   public static final String ID_LIST_ADAPTERS_DETAILS_DEVICES =
         "vsphere.core.storage.adaptersDetails.views.devices/list";
   public static final String ID_BUTTON_RENAME =
         "vsphere.core.storage.devices.rename/button";
   public static final String ID_LIST_STORAGE_DEVICES =
         "vsphere.core.host.configuration.storageDevicesList/list";
   public static final String ID_BUTTON_STORAGE_DEVICES_COPYTOCLIPBOARD =
         "vsphere.core.host.configuration.storageDevicesList/copyToClipboard";
   public static final String ID_LIST_STORAGE_ADAPTERS =
         "vsphere.core.host.configuration.storageAdaptersList/list";
   public static final String ID_VIEW_DEVICE_DETAILS =
         "vsphere.core.host.configuration.storageDevicesDetails";
   public static final String ID_LABEL_CONNECTIVITY_MULTIPATHING =
         "titleGroup/titleLabel";
   public static final String ID_LABEL_NO_DETAILS = "noDetailsAvailableLabel";
   public static final String ID_TEXT_PSP_VALUE = "pathSelectionPolicy";
   public static final String ID_LABEL_FILE_SYSTEM = "titleLabel.fileSystemStackBlock";
   public static final String ID_TREE_STORAGE_TOC =
         "vsphere.core.host.manage.storage/tocTree";
   public static final String ID_BUTTON_BAR_STORAGE_DEVICES =
         "_StorageDevicesList_ActionfwButtonBar1/";
   public static final String ID_STORAGE_DEVICES_TAB_NAVIGATOR =
         "_StorageDevicesDetailsView_TabbedView1/tabNavigator";
   public static final String ID_BUTTON_STORAGE_DEVICES_EDIT_MULTIPATH =
         "editMultipathingPolicyAction";
   public static final String ID_COMBO_BOX_PATH_SELECTION_POLICY = "policyCombo";
   public static final String ID_LIST_EDIT_MULTIPATH = "pathsList";
   public static final String ID_TEXT_PATH_SELECT_POLICY = "pathSelectionPolicy";
   public static final String ID_TEXT_PREFERRED_PATH = "preferredPath";
   public static final String ID_LIST_STORAGE_DEVICES_PATH =
         "vsphere.core.storage.devicesDetails.views.paths/list";
   public static final String ID_BUTTON_ENABLE = "enableButton";
   public static final String ID_BUTTON_DISABLE = "disableButton";
   public static final String ID_LIST_STORAGE_ADAPTERS_PATH =
         "vsphere.core.storage.adaptersDetails.views.paths/list";
   public static final String ID_TEXT_STORAGE_DEVICE_HEADER = "deviceHeader";
   public static final String ID_ARROW_IMAGE_PATHS =
         "arrowImage_MultipathingPathsView_StackBlock2";
   public static final String ID_LIST_CONNECTIVITY_MULTIPATH =
         "vsphere.core.datastore.manage.settings.connectivity.multipathingView/pathsList";
   public static final String ID_CHECKBOX_RUNTIME_NAME = "TriStateCheckBox_2";
   public static final String ID_CHECKBOX_STATUS = "TriStateCheckBox_3";
   public static final String ID_CHECKBOX_ADAPTER = "TriStateCheckBox_4";
   public static final String ID_CHECKBOX_SDADAPTER = "TriStateCheckBox_10";
   public static final String ID_CHECKBOX_OWNER = "TriStateCheckBox_12";
   public static final String ID_TEXT_SSD_INFO = "ssdInfoText";
   public static final String ID_GROUP_EDIT_MULTIPATH =
         "_MultipathingPolicyForm_VGroup2";

   // Object Navigator constants
   public static final String APPLICATION_NAVIGATOR = "vsphere.core.navigator";
   public static final String LICENSE_APPLICATION_NAVIGATOR =
         "vsphere.license.navigator";
   public static final String CATEGORY_NODE_VIEW = "CategoryNodeView_";
   public static final String CATEGORY_NODE_LIST_VIEW = "CategoryNodeListView_";
   public static final String CATEGORY_NODE_TREE_VIEW = "CategoryNodeTreeView_";
   public static final String TREE_NODE_ITEM = "TreeNodeItem_";
   public static final String RELATED_ITEMS_NODE_VIEW = "RelatedItemsNodeView_";
   public static final String APP_OBJ_NAVIGATOR = "appObjectNavigator";
   public static final String OBJ_NAV_LAUNCHER_VIEW = new StringBuffer(
         CATEGORY_NODE_LIST_VIEW).append(APPLICATION_NAVIGATOR).append(".launcher")
         .toString();
   public static final String OBJ_NAV_VI_VIEW =
         new StringBuffer(CATEGORY_NODE_TREE_VIEW).append(APPLICATION_NAVIGATOR)
               .append(".virtualInfrastructure").toString();
   public static final String OBJ_NAV_RULES_PROFILES_VIEW = new StringBuffer(
         CATEGORY_NODE_TREE_VIEW).append(APPLICATION_NAVIGATOR)
         .append(".rulesAndProfiles").toString();
   public static final String OBJ_NAV_ADMINISTRATION_VIEW = new StringBuffer(
         CATEGORY_NODE_TREE_VIEW).append(APPLICATION_NAVIGATOR)
         .append(".administration").toString();
   public static final String OBJ_NAV_LICENSING_VIEW = new StringBuffer(
         CATEGORY_NODE_TREE_VIEW).append(LICENSE_APPLICATION_NAVIGATOR)
         .append(".licensing").toString();
   public static final String OBJ_NAV_HOME_NODE_ITEM = new StringBuffer(TREE_NODE_ITEM)
         .append(APPLICATION_NAVIGATOR).append(".home").toString();
   public static final String OBJ_NAV_MANAGE_MONITOR_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".tasks").toString();
   public static final String OBJ_NAV_ADMINISTRATION_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".administration")
         .toString();
   public static final String OBJ_NAV_SEARCH_NODE_ITEM =
         new StringBuffer(TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
               .append(".search").toString();
   public static final String OBJ_NAV_SAVEDSEARCH_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".savedSearches")
         .toString();
   public static final String OBJ_NAV_VIRTUAL_INFRASTRUCTURE_NODE_ITEM =
         new StringBuffer(TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
               .append(".virtualInfrastructure").toString();
   public static final String OBJ_NAV_RULES_PROFILE_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".rulesAndProfiles")
         .toString();
   public static final String OBJ_NAV_TASKS_NODE_ITEM = new StringBuffer(TREE_NODE_ITEM)
         .append(APPLICATION_NAVIGATOR).append(".tasks").toString();
   public static final String OBJ_NAV_EVENTS_NODE_ITEM =
         new StringBuffer(TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
               .append(".events").toString();
   public static final String OBJ_NAV_VI_HOME_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".viHome").toString();
   public static final String OBJ_NAV_VI_TREE_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".viInventoryTree")
         .toString();
   public static final String OBJ_NAV_VIRTUAL_CENTERS_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
         .append(".vilist.Folder_isRootFolder").toString();
   public static final String OBJ_NAV_DATACENTERS_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".vilist.Datacenter_")
         .toString();
   public static final String OBJ_NAV_HOSTS_NODE_ITEM = new StringBuffer(TREE_NODE_ITEM)
         .append(APPLICATION_NAVIGATOR).append(".vilist.HostSystem_").toString();
   public static final String OBJ_NAV_CLUSTERS_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
         .append(".vilist.ClusterComputeResource_").toString();
   public static final String OBJ_NAV_DATASTORES_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".vilist.Datastore_")
         .toString();
   public static final String OBJ_NAV_DATASTORE_CLUSTERS_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".vilist.StoragePod_")
         .toString();
   public static final String OBJ_NAV_RESOURCE_POOLS_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
         .append(".vilist.ResourcePool_isNonRootRP").toString();
   public static final String OBJ_NAV_STANDARD_NETWORKS_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
         .append(".vilist.Network_isStandardNetwork").toString();
   public static final String OBJ_NAV_DISTRIBUTED_SWITCHES_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
         .append(".vilist.DistributedVirtualSwitch_").toString();
   public static final String OBJ_NAV_DISTRIBUTED_PORT_GROUPS_NODE_ITEM =
         new StringBuffer(TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
               .append(".vilist.DistributedVirtualPortgroup_!isUplinkPortgroup")
               .toString();
   public static final String OBJ_NAV_DISTRIBUTED_UPLINK_PORT_GROUPS_NODE_ITEM =
         new StringBuffer(TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
               .append(".vilist.DistributedVirtualPortgroup_isUplinkPortgroup")
               .toString();
   public static final String OBJ_NAV_VIRTUAL_MACHINES_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
         .append(".vilist.VirtualMachine_isNormalVMOrPrimaryFTVM").toString();
   //   public static final String OBJ_NAV_VIRTUAL_MACHINES_NODE_ITEM = new StringBuffer(
   //	         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(
   //	         ".vilist.VirtualMachine_!config.template").toString();
   public static final String OBJ_NAV_VAPPS_NODE_ITEM = new StringBuffer(TREE_NODE_ITEM)
         .append(APPLICATION_NAVIGATOR).append(".vilist.VirtualApp_").toString();
   public static final String OBJ_NAV_VM_TEMPLATES_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
         .append(".vilist.VirtualMachine_config.template").toString();
   public static final String OBJ_NAV_COMPUTE_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".viCompute").toString();
   public static final String OBJ_NAV_NETWORKING_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".viNetworking")
         .toString();
   public static final String OBJ_NAV_STORAGE_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".viStorage").toString();
   public static final String OBJ_NAV_VMS_TEMPLATES_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".viVmsAndTemplates")
         .toString();
   public static final String OBJ_NAV_HOST_PROFILES_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".hostProfiles")
         .toString();
   public static final String OBJ_NAV_VM_STORAGE_PROFILES_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".storageProfiles")
         .toString();
   public static final String OBJ_NAV_TAG_MANAGER_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".tagsManager").toString();
   public static final String OBJ_NAV_SOLUTION_EXTENSION_NODE_ITEM =
         "TreeNodeItem_vsphere.core.navigator.solutionVCAndExtensions";

   public static final String OBJ_NAV_ROLE_MANAGER_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".accessRoleManager")
         .toString();
   public static final String OBJ_NAV_LICENSES_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(LICENSE_APPLICATION_NAVIGATOR).append(".management")
         .toString();
   public static final String OBJ_NAV_LICENSE_REPORTS_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(LICENSE_APPLICATION_NAVIGATOR).append(".reporting")
         .toString();
   public static final String OBJ_NAV_SSO_USERS_GROUPS_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".ssoUsersAndGroups")
         .toString();
   public static final String OBJ_NAV_SSO_VC_REGISTTRATION_TOOL_NODE_ITEM =
         new StringBuffer(TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
               .append(".ssoVcRegistrationTool").toString();
   public static final String OBJ_NAV_SSO_CONFIGURATION_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".ssoConfiguration")
         .toString();
   public static final String OBJ_NAV_SOLUTION_PLUGIN_MANAGER_NODE_ITEM =
         new StringBuffer(TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
               .append(".solutionPluginManager").toString();
   public static final String OBJ_NAV_SOLUTION_VCS_EXTENSIONS_NODE_ITEM =
         new StringBuffer(TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
               .append(".solutionVCAndExtensions").toString();
   public static final String OBJ_NAV_VI_TREE_NODE_VIEW = new StringBuffer(
         RELATED_ITEMS_NODE_VIEW).append(APPLICATION_NAVIGATOR)
         .append(".viVCenterServers").toString();
   public static final String OBJ_NAV_DATACENTERS_TREE_NODE_VIEW = new StringBuffer(
         RELATED_ITEMS_NODE_VIEW).append(APPLICATION_NAVIGATOR).append(".viDatacenters")
         .toString();
   public static final String OBJ_NAV_COMPUTE_TREE_NODE_VIEW = new StringBuffer(
         RELATED_ITEMS_NODE_VIEW).append(APPLICATION_NAVIGATOR).append(".viCompute")
         .toString();
   public static final String OBJ_NAV_NETWORKING_TREE_NODE_VIEW = new StringBuffer(
         RELATED_ITEMS_NODE_VIEW).append(APPLICATION_NAVIGATOR).append(".viNetworking")
         .toString();
   public static final String OBJ_NAV_STORAGE_TREE_NODE_VIEW = new StringBuffer(
         RELATED_ITEMS_NODE_VIEW).append(APPLICATION_NAVIGATOR).append(".viStorage")
         .toString();
   public static final String OBJ_NAV_VMS_TEMPLATES_TREE_NODE_VIEW = new StringBuffer(
         RELATED_ITEMS_NODE_VIEW).append(APPLICATION_NAVIGATOR)
         .append(".viVmsAndTemplates").toString();
   public static final String OBJ_NAV_HOST_PROFILES_TREE_NODE_VIEW = new StringBuffer(
         RELATED_ITEMS_NODE_VIEW).append(APPLICATION_NAVIGATOR).append(".hostProfiles")
         .toString();
   public static final String OBJ_NAV_VM_STORAGE_PROFILES_TREE_NODE_VIEW =
         new StringBuffer(RELATED_ITEMS_NODE_VIEW).append(APPLICATION_NAVIGATOR)
               .append(".storageProfiles").toString();
   // Obj Nav Inventory Trees Nodes
   public static final String OBJ_NAV_HOSTS_AND_CLUSTERS_TREE_NODE_ITEM =
         new StringBuffer(TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
               .append(".hostsClustersTree").toString();
   public static final String OBJ_NAV_VMS_AND_TEMPLATES_TREE_NODE_ITEM =
         new StringBuffer(TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR)
               .append(".vmTemplatesTree").toString();
   public static final String OBJ_NAV_STORAGE_TREE_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".storageTree").toString();
   public static final String OBJ_NAV_NETWORKING_TREE_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(APPLICATION_NAVIGATOR).append(".networkingTree")
         .toString();
   public static final String OBJ_NAV_DATASTORE_NODE_VIEW = new StringBuffer(
         RELATED_ITEMS_NODE_VIEW).append(APPLICATION_NAVIGATOR)
         .append(".vilist.Datastore_").toString();
   public static final String OBJ_NAV_VIRTUAL_CENTER_NODE_VIEW = new StringBuffer(
         RELATED_ITEMS_NODE_VIEW).append(APPLICATION_NAVIGATOR)
         .append(".vilist.Folder_isRootFolder").toString();
   public static final String EDGE_SETS_DATAGRID = "edgeSetsDataGrid";
   public static final String SELECTED_SET_DATAGRID = "selectedSetDataGrid";
   public static final String VC_DATAPROVIDER_ID = "Folder";
   public static final String CONTENTS_DATAPROVIDER_ID = "contentsForVCenter";
   public static final String DATACENTERS_DATAPROVIDER_ID = "dcsForVCenter";
   public static final String HOSTS_DATAPROVIDER_ID = "hostsForVCenter";
   public static final String DATASTORE_HOSTS_DATAPROVIDER_ID = "hostsForDatastore";

   public static final String CLUSTERS_DATAPROVIDER_ID = "clustersForVCenter";
   public static final String VMS_DATAPROVIDER_ID = "vmsForVCenter";
   public static final String HOST_PROFILES_DATAPROVIDER_ID = "hpsForVCenter";
   public static final String RESOURCE_POOL_DATAPROVIDER_ID = "resourcePoolsForHAndC";
   public static final String NETWORKS_DATAPROVIDER_ID = "networksForNetworking";
   public static final String APP_OBJECT_NAVIGATOR_HOME_BUTTON = "homeButton";
   public static final String ID_LABEL_FOCUSED_NODE = "focusedObjLabel";
   public static final String DATACENTERS_FOLDER_CONTENTS_DATAPROVIDER_ID =
         "contentsForDatacenterFolder";
   public static final String ID_LABEL_NODE_ITEM = "titleText";
   public static final String ID_SELECTED_ITEM_ID = "selectedItem.id";
   public static final String ID_BUTTON_OBJ_NAV_HISTORY = "historyDropdownButton";
   public static final String ID_SPARK_LIST_OBJ_NAV_HISTORY = "historyList";

   // Resource Pool Data provide IDs
   public static final String HOST_RP_CONTENTS_DATAPROVIDER_ID =
         "rootRpContentsForStandAloneHost";
   public static final String CLUSTER_RP_CONTENTS_DATAPROVIDER_ID =
         "rootRpContentsForCluster";
   public static final String RP_CONTENTS_DATAPROVIDER_ID = "contentsForRp";

   public static final String ID_UID_CLUSTER_EVENTS =
         "vsphere.core.cluster.monitor.events";

   // vApp data provide IDs
   public static final String HOST_VAPP_CONTENTS_DATAPROVIDER_ID =
         "vappsForStandaloneHost";
   public static final String HOST_VAPP_CONTENTS_DATAPROVIDER_ID2 = "vApps";
   public static final String CLUSTER_VAPP_CONTENTS_DATAPROVIDER_ID = "vappsForCluster";
   public static final String VAPP_CONTENTS_DATAPROVIDER_ID = "contentsForVApp";

   // Datacenter view ids
   public static final String DATACENTERS_DATAPROVIDER_ID_DATACENTERS_VIEW =
         "Datacenter";
   public static final String DATACENTERS_CONTENTS_DATAPROVIDER_ID =
         "childObjectsForDatacenter";

   // Hosts view ids
   public static final String HOSTS_DATAPROVIDER_ID_HOSTS_VIEW = "HostSystem";

   // Resource Pools view ids
   public static final String RESOURCE_POOLS_DATAPROVIDER_ID_RESOURCE_POOLS_VIEW =
         "ResourcePool";

   // Clusters view ids
   public static final String CLUSTERS_DATAPROVIDER_ID_CLUSTERS_VIEW =
         "ClusterComputeResource";

   // Datastore view ids
   public static final String DATASTORES_DATAPROVIDER_ID_DATASTORES_VIEW = "Datastore";

   // Datastore clusters view ids
   public static final String DATASTORE_CLUSTERS_DATAPROVIDER_ID_DATASTORE_CLUSTERS_VIEW =
         "StoragePod";
   public static final String ID_LABEL_POD_VM_OVERRIDE_SETTING =
         "vsphere.core.dscluster.manage.wrapper/tocTree/automationName=VM Overrides";

   // Standard Networks view ids
   public static final String STANDARD_NETWORKS_ID_STANDARD_NETWORKS_VIEW = "Network";

   // Distributed switches view ids
   public static final String DISTRIBUTED_SWITCHES_ID_DISTRIBUTED_SWITCHES_VIEW =
         "DistributedVirtualSwitch";

   // Distributed port groups view ids
   public static final String DISTRIBUTED_PORT_GROUPS_ID_DISTRIBUTED_PORT_GROUPS_VIEW =
         "DistributedVirtualPortgroup";

   // VMs view ids
   public static final String VMS_ID_VMS_VIEW = "VirtualMachine";

   // Standard Networks view ids
   public static final String VAPPS_ID_VAPPS_VIEW = "VirtualApp";

   // Compute View ids
   public static final String HOSTS_DATAPROVIDER_ID_COMPUTE_VIEW = "hostsForHAndC";
   public static final String CLUSTERS_DATAPROVIDER_ID_COMPUTE_VIEW = "clustersForHAndC";
   public static final String RES_POOLS_DATAPROVIDER_ID_COMPUTE_VIEW =
         "resourcePoolsForHAndC";

   // Networking view ids
   public static final String STANDARD_NETWORKS_DATAPROVIDER_ID_NETWORKING_VIEW =
         "Network";
   public static final String DISTRIBUTED_SWITCHES_DATAPROVIDER_ID_NETWORKING_VIEW =
         "vdsForNetworking";
   public static final String DISTRIBUTED_PORT_GROUPS_DATAPROVIDER_ID_NETWORKING_VIEW =
         "portgroupsForNetworking";
   public static final String UPLINK_PORT_GROUPS_DATAPROVIDER_ID_NETWORKING_VIEW =
         "uplinksForNetworking";
   public static final String HOSTS_DATAPROVIDER_ID_NETWORKING_VIEW =
         "hostsForNetworking";
   public static final String NETWORK_FOLDER_CONTENTS_DATAPROVIDER_ID =
         "contentsForNetworkFolder";

   // Host view ids
   public static final String CHILD_RPS_DATAPROVIDER_ID_HOST_VIEW =
         "rootRpContentsForStandAloneHost";
   public static final String VMS_DATAPROVIDER_ID_HOST_VIEW = "vmsForStandaloneHost";
   public static final String DATASTORES_DATAPROVIDER_ID_HOST_VIEW =
         "datastoresForStandaloneHost";

   // Storage view ids
   public static final String DATASTORES_DATAPROVIDER_ID_STORAGE_VIEW =
         "datastoresForStorage";
   public static final String DATASTORE_CLUSTERS_DATAPROVIDER_ID_STORAGE_VIEW =
         "datastoreClustersForStorage";
   public static final String STORAGE_FOLDER_CONTENTS_DATAPROVIDER_ID =
         "contentsForStorageFolder";

   // VMS and Templates view ids
   public static final String VMS_DATAPROVIDER_ID_STORAGE_VIEW = "vmsForVmsAndTemplates";
   public static final String VM_TEMPLATES_DATAPROVIDER_ID_STORAGE_VIEW =
         "templatesForVmsAndTemplates";
   public static final String VAPPS_DATAPROVIDER_ID_STORAGE_VIEW =
         "vAppsForVmsAndTemplates";
   public static final String VM_FOLDER_CONTENTS_DATAPROVIDER_ID =
         "contentsForVMandTemplFolder";
   public static final String VM_TEMPLATES_ID_VM_TEMPLATES_VIEW = "VirtualMachine";
   public static final String ID_VMFOLDER_CONTENTS_ADVANCED_DATAGRID =
         VM_FOLDER_CONTENTS_DATAPROVIDER_ID + "/list";
   public static final String ID_VM_DATASTORES_CONTENTS_ADVANCED_DATAGRID =
         "datastoresForVm" + "/list";

   // Host Profile view ids
   public static final String HOST_PROFILES_ID_HOST_PROFILES_VIEW = "hostProfiles";

   // Datacenter node view ids
   public static final String VM_ID_DATACENTER_VIEW = "vmsForDatacenter";

   // Datastore node view ids
   public static final String VM_ID_DATASTORE_VIEW = "vmsForDatastore";

   // VM Network node view ids
   public static final String VMS_ID_VM_NETWORK_VIEW = "vmsForNetwork";
   public static final String HOSTS_ID_VM_NETWORK_VIEW = "hostsForNetwork";

   // VC node view ids
   public static final String VMS_ID_VIRTUAL_CENTER_VIEW = "vmsForVCenter";
   public static final String VM_TEMPLATE_ID_VIRTUAL_CENTER_VIEW =
         "vmTemplatesForVCenter";
   public static final String VAPPS_ID_VIRTUAL_CENTER_VIEW = "vappsForVCenter";
   public static final String DV_PORT_GROUP_ID_VIRTUAL_CENTER_VIEW = "dvpgForVCenter";

   // vDS node view ids
   public static final String UPLINK_PORTGROUP_ID_VDS_VIEW = "uplinksForVMWDVS";

   // Host Profile ids
   public static final String ATTACH_HOST_ACTION = "attachAction";
   public static final String CREATE_PROFILE_ACTION = "createProfileFromAHostAction";
   public static final String CREATE_ACTION_GLOBAL = "createActionGlobal";
   public static final String CHECK_COMPLIANCE_ACTION = "checkComplianceAction";
   public static final String CHANGE_HOST_PROFILE_ACTION = "changeHostProfileAction";
   public static final String DETACH_HOST_ACTION = "detachHostClusterAction";
   public static final String DELETE_PROFILE_ACTION = "deleteAction";
   public static final String RENAME_PROFILE_ACTION = "renameAction";
   public static final String DUPLICATE_HOST_PROFILE_ACTION = "duplicateProfileAction";
   public static final String EDIT_ACTION_GLOBAL = "editHostProfileAction";
   public static final String REMEDIATE_HOST_ACTION = "remediateHostAction";
   public static final String CHANGE_REF_HOST_ACTION = "changeRefHostAction";
   public static final String RESET_HOST_CUSTOMIZATIONS_ACTION =
         "resetHostCustomizationsAction";
   public static final String ID_MONITOR_COMPLIANCE_VIEW = "monitor.complianceView";
   public static final String ID_MANAGE_SETTINGS_VIEW = "manage.settingsView";
   public static final String ID_PROFILE_COMPLIANCE_EXTENSION_VIEW =
         "monitor.profileComplianceMultiExtentsionView";
   public static final String ID_MANAGE_PROFILE_COMPLIANCE = "manage.profileCompliance.";
   public static final String ID_MONITOR_PROFILE_COMPLIANCE =
         "monitor.profileCompliance.";
   public static final String ID_IMAGE_ICON_CLUSTER_ATTACH_HOST_WIZARD =
         "Image_CollapseExpand_0_";
   public static final String ID_HOST_PROFILE_SELECTION_PAGE =
         "hostProfileSelectionPage";
   public static final String ID_HOST_PROFILE_PROFILE_NAME = ID_TIWO_DIALOG
         + "/profileName";
   public static final String MONITOR_COMPLIANCE_HOST = "monitor.profileCompliance.host";
   public static final String ID_BUTTON_CHECK_COMPLIANCE_CLUSTER_MONITOR =
         "checkComplianceActionButton";
   public static final String ID_TOC_LIST_CLUSTER_SETTINGS =
         "vsphere.core.cluster.manage.settingsView/tocTree";
   public static final String ID_PROFILE_COMPLIANCE_PAGE =
         "vsphere.core.cluster.monitor.profileCompliance";
   public static final String ID_PROFILE_COMPLIANCE_HOST_PAGE =
         ID_PROFILE_COMPLIANCE_PAGE + ".hostsCompliance";
   public static final String ID_ADVANCED_DATAGRID_CLUSTER_COMPLIANCE =
         ID_PROFILE_COMPLIANCE_HOST_PAGE + "/" + ID_HOSTS_ADVANCE_DATAGRID;
   public static final String ID_COLLAPSED_TEXT_CLUSTER_REQUIREMENTS =
         ID_PROFILE_COMPLIANCE_PAGE + ".cluster/collapsedText";
   public static final String ID_COLLAPSED_TEXT_HOST_REQUIREMENTS =
         ID_PROFILE_COMPLIANCE_PAGE + ".host/collapsedText";
   public static final String ID_LABEL_HOST_PROFILE_STATUS =
         ID_PROFILE_COMPLIANCE_HOST_PAGE + "/" + "hostProfileStatusLabel";
   public static final String ID_IMAGE_COMPLIANCE_FAILURES = "failures/arrowImagenull";
   public static final String ID_LABEL_COMPLIANCE_TITLE = "failures/titleLabel";
   public static final String ID_IMAGE_ICON_CLUSTER_REQUIREMENTS =
         "vsphere.core.cluster.monitor."
               + "profileCompliance.cluster/detailsBlock/arrowImagedetailsBlock";
   public static final String ID_IMAGE_ICON_HOST_PROFILE =
         "vsphere.core.cluster.monitor."
               + "profileCompliance.host/detailsBlock/arrowImagedetailsBlock";
   public static final String ID_LABEL_CLUSTER_REQUIREMENTS = "vsphere."
         + "core.cluster.monitor.profileCompliance.cluster/detailsBlock/detailsText1";
   public static final String ID_LABEL_HOST_PROFILE_TEXT_1 = "vsphere."
         + "core.cluster.monitor.profileCompliance.host/detailsBlock/detailsText1";
   public static final String ID_LABEL_HOST_PROFILE_TEXT_3 = "vsphere."
         + "core.cluster.monitor.profileCompliance.host/detailsBlock/detailsText3";
   public static final String ID_MANAGE_PROFILE_COMPLIANCE_HOST =
         ID_MANAGE_PROFILE_COMPLIANCE + "host";
   public static final String ID_HOSTS_COMPLIANCE = "hostsCompliance";
   public static final String ID_COMPLIANCE_LIST = "complianceList";
   public static final String ID_TEXT_DISPLAY = "textDisplay";
   public static final String ID_ACTION_CREATE_HOST_PROFILE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HP + CREATE_ACTION_GLOBAL;
   public static final String ID_ACTION_EDIT_HOST_PROFILE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HP + EDIT_ACTION_GLOBAL;
   public static final String ID_ACTION_DELETE_HOST_PROFILE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HP + DELETE_PROFILE_ACTION;
   public static final String ID_ACTION_DELETE_HOST_PROFILE_MENU = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST_PROFILE + DELETE_PROFILE_ACTION;
   public static final String ID_ACTION_CREATE_HOST_PROFILE_MENU = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST_PROFILE + CREATE_PROFILE_ACTION;
   public static final String ID_ACTION_RESET_HOST_CUSTOMIZATIONS = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST_PROFILE + RESET_HOST_CUSTOMIZATIONS_ACTION;
   public static final String ID_ACTION_EDIT_HOST_PROFILE_MENU = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST_PROFILE + EDIT_ACTION_GLOBAL;
   public static final String ID_ACTION_CHECK_COMPLIANCE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + CHECK_COMPLIANCE_ACTION;
   public static final String ID_ACTION_CHECK_HOST_PROFILE_COMPLIANCE = EXTENSION_PREFIX
         + "hostProfile." + CHECK_COMPLIANCE_ACTION;
   public static final String ID_ACTION_DETACH_HOST = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST_PROFILE + DETACH_HOST_ACTION;
   public static final String ID_ACTION_HOST_PROFILE_RENAME = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST_PROFILE + RENAME_PROFILE_ACTION;
   public static final String ID_ACTION_COPY_SETTINGS_FROM_HOST = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST_PROFILE + CHANGE_REF_HOST_ACTION;
   public static final String ID_ACTION_RUN_SCHEDULED_TASK =
         "vsphere.client.scheduling.runScheduledOperationAction";
   public static final String ID_ACTION_REMEDIATE_HOST = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST_PROFILE + REMEDIATE_HOST_ACTION;
   public static final String ID_ACTION_ATTACH_HOST = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST_PROFILE + ATTACH_HOST_ACTION;
   public static final String ID_COMPLIANCE_VIEW_BUTTON = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HP + ID_MONITOR_COMPLIANCE_VIEW + ".button";
   public static final String ID_MANAGE_SETTINGS_VIEW_BUTTON = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HP + ID_MANAGE_SETTINGS_VIEW + ".button";
   public static final String ID_ENTITY_NAME = "entityName";
   public static final String ID_BUTTON_PROFILE_COMPLIANCE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_CLUSTER + ID_PROFILE_COMPLIANCE_EXTENSION_VIEW + ".button";
   public static final String ID_BUTTON_ATTACH_PROFILE = "btn_" + EXTENSION_PREFIX
         + EXTENSION_ENTITY_CLUSTER + EXTENSION_ENTITY_PROFILE + EXTENSION_ENTITY_HOST
         + EXTENSION_ENTITY_PROFILE + "attach";
   public static final String ID_BUTTON_LINK_LEARN_ABOUT_HP = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HP + "gettingStartedView/gettingStartedHelpLink_1/link";
   public static final String ID_BUTTON_LINK_LEARN_USE_HP = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HP + "gettingStartedView/gettingStartedHelpLink_2/link";
   public static final String ID_BUTTON_LINK_LEARN_ABOUT_HOSTS = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HP + "gettingStartedView/gettingStartedHelpLink_3/link";
   public static final String ID_BUTTON_SETTING_VIEW_HP = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HP + "manage.settingsView.button";
   public static final String ID_BUTTON_SCHEDULED_TASKS_HP = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HP + "manage.scheduledOpsView.button";
   public static final String ID_SETTINGS_COMPONENT = "settingsComponent";
   public static final String ID_SELECTION_CONTROL = "selectionControl";
   public static final String ID_COLLAPSED_TEXT = "collapsedText";
   public static final String ID_LABEL_NTP_CLIENT = "lblNtpClient";
   public static final String ID_DATAGRID_COMPLIANCE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_CLUSTER + ID_MONITOR_PROFILE_COMPLIANCE
         + ID_HOSTS_COMPLIANCE + "/" + "hostList";
   public static final String ID_BUTTON_EDIT_TIME_CONFIG = "btn_" + EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + "editTimeConfigAction";
   public static final String ID_RADIO_BUTTON_CONFIG_NTP = "rbConfigWithNTP";
   public static final String ID_RADIO_BUTTON_CONFIG_NTP_MANUALLY = "rbConfigManually";
   public static final String ID_LABEL_SELECT_REFERENCE_HOST =
         "_SelectReferenceHostPage_Label1";
   public static final String ID_LABEL_NAME_DESCRIPTION =
         "_NameAndDescriptionPage_Label1";
   public static final String ID_ADVANCED_DATAGRID_CHANGE_PROFILE = "selectionControl";
   public static final String ID_ACTION_CHANGE_HOST_PROFILE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST_PROFILE + CHANGE_HOST_PROFILE_ACTION;
   public static final String ID_PORTLET_HOST_PROFILE_COMPLIANCE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST + "summary.hostProfileCompliancePortlet";
   public static final String ID_ACTION_DUPLICATE_HOST_PROFILE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_HOST_PROFILE + DUPLICATE_HOST_PROFILE_ACTION;
   public static final String ID_LABEL_PROFILE_NAME = "profileNameLabel";
   public static final String ID_ADVANCED_DATAGRID_CREATE_HOST_PROFILE =
         "HostSystem_filterGridView";
   public static final String ID_DEFAULT_SUMMARY_PAGE = "defaultSummaryPage";
   public static final String ID_PAGE_DESCRIPTION_ELEMENT =
         "pageHeaderDescriptionElement";
   public static final String ID_REFERENCE_HOST_PAGE = "referenceHostPage";
   public static final String ID_NAME_DESCRIPTION_PAGE = "nameAndDescriptionPage";
   public static final String ID_ADVANCED_DATAGRID_AVAILABLE_VIEW = "availableView";
   public static final String ID_ADVANCED_DATAGRID_SELECTED_VIEW = "selectedView";
   public static final String ID_BUTTON_ATTACH_SELECTED = "attachSlected";
   public static final String ID_COLUMN_HEADER_NAME = "name";
   public static final String ID_HOST_COMPLIANCE_STATUS = "hostComplianceStatus";
   public static final String ID_LABEL_SUMMARY_PROFILE_NAME =
         "summary_profileName_valueLbl";
   public static final String ID_LABEL_SUMMARY_PROFILE_DESCRIPTION =
         "summary_profileDescripion_valueLbl";
   public static final String ID_LABEL_SUMMARY_PROFILE_CREATION_DATE =
         "summary_profileCreationDate_valueLbl";
   public static final String ID_LABEL_SUMMARY_PROFILE_EDIT_DATE =
         "summary_profileModificationDate_valueLbl";
   public static final String ID_PAGE_ID_NAME_DESCRIPTION = "pageIdNameAndDescription";
   public static final String ID_DESCRIPTION = "description";
   public static final String ID_PAGE_EDIT = "pageIdEdit";
   public static final String ID_TRI_STATE_CHECKBOX_TREE_EDIT_PROFILE =
         "hostProfileInternalTree";
   public static final String ID_LIST_SYSTEM_IMAGE_CACHE =
         "tiwoDialog/pageIdEdit/className=DropDownList";
   public static final String ID_CHECKBOX_VMFS_OVERWRITE =
         "tiwoDialog/pageIdEdit/overwriteVmfs";
   public static final String ID_TEXT_INPUT_EDIT_PROFILE =
         "tiwoDialog/pageIdEdit/cpuMinlimit";
   public static final String ID_TEXT_INPUT_EDIT_PROFILE_FIRST_DISK =
         "tiwoDialog/pageIdEdit/firstDisk";
   public static final String ID_LABEL_PROFILE_EDIT_PROFILE =
         "tiwoDialog/pageIdEdit/titleLabel";
   public static final String ID_LABEL_ERROR_EDIT_PROFILE =
         "tiwoDialog/pageIdEdit/_errorLabel";
   public static final String ID_TEXT_TARGET_NAME = "targetNameText";
   public static final String ID_TEXT_OPERATION_NAME = "operationText";
   public static final String ID_MESSAGE_TEXT = "messageText";
   public static final String ID_BUTTON_NO_REMOVE_TASK = "noBtn";
   public static final String ID_BUTTON_DETACH_ATTACH_HOST_WIZARD = "detachSelected";
   public static final String ID_LABEL_ATTACH_HOST_WARNING =
         "hostSelectionPage/_message";
   public static final String ID_BUTTON_ISSUES_SUB_TAB =
         "vsphere.core.hp.monitor.issues.button";
   public static final String ID_BUTTON_NETWORKING_SUB_TAB =
         "vsphere.core.host.manage.networkingView.button";
   public static final String ID_TOC_TREE_NETWORKING_VIEW =
         "vsphere.core.host.manage.networkingView/tocTree";
   public static final String ID_LIST_DNS_ROUTING =
         "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_BUTTON_EDIT_PASSWORD = "tiwoDialog/editPasswordButton";
   public static final String ID_TEXT_DISPLAY_PASSWORD = "password/textDisplay";
   public static final String ID_TEXT_DISPLAY_CONFIRM_PASSWORD =
         "confirmPassword/textDisplay";
   public static final String ID_BUTTON_UPDATE_PASSWORD = "updatePasswordButton";
   public static final String ID_DATAGRID_TASKS_LIST =
         "vsphere.core.hp.manage.scheduledOpsView.defaultListView/list";
   public static final String ID_BUTTON_REMOVE_SCHEDULED_TASK =
         "vsphere.client.scheduling.removeScheduledOperationAction/button";
   public static final String ID_BUTTON_RUN_SCHEDULED_TASK =
         "vsphere.client.scheduling.runScheduledOperationAction/button";
   public static final String ID_DATAGRID_HISTORY_TASKS =
         "vsphere.core.hp.vm.manage.scheduledOpsView.defaultDetailsView/historyList";

   public static final String ID_CLASSNAME_ANCHORED_DIALOG = "className=AnchoredDialog";

   public static final String ID_TEXT_INPUT_SNAPSHOT_NAME = "snapshotTextInput";

   // Edit Host Profile - Networking
   public static final String ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_NETWORKING_CONFIGURATION =
         "0.3.";
   public static final String ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_VSPHERE_DISTRIBUTED_SWITCH =
         ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_NETWORKING_CONFIGURATION + ".7.";
   public static final String ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_DVSWITCH =
         ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_VSPHERE_DISTRIBUTED_SWITCH + ".0.";
   public static final String ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_UPLINK_PORT_CONFIGURATION =
         ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_DVSWITCH + ".0.";
   public static final String ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_UPLINK_PORT_CONFIGURATION_SUBITEM =
         ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_UPLINK_PORT_CONFIGURATION + ".0.";
   public static final String ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_HOST_PORT_GROUP =
         ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_NETWORKING_CONFIGURATION + ".2.";
   public static final String ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_VMKERNEL =
         ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_HOST_PORT_GROUP + ".0.";
   public static final String ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_IPADDRESS =
         ID_CHECKBOX_PROFILE_TREE_ITEM_RENDERER_VMKERNEL + ".0.";
   // TODO: Need unique IDs from dev - PR 847053
   public static final String ID_SPARK_TEXTINPUT_NAMES_OF_NICS_TO_ATTACH = "label="
         + DETERMINE_NICS_TO_USE_LABEL + ":/className=RichEditableText";
   public static final String ID_SPARK_DROPDOWNLIST_IPV4_ADDRESS = "label="
         + IPV4_ADDRESS_LABEL + ":/className=DropDownList";
   public static final String ID_SPARK_DROPDOWNLIST_STATIC_IPV6_ADDRESS = "label="
         + STATIC_IPV6_ADDRESS_LABEL + ":/className=DropDownList";

   // 'Remediate Hosts' dialog - 'Customize Hosts' page
   // TODO: Getting unique id is block due to PR, to be filed by Stanko.
   public static final String ID_SPARK_TEXTINPUT_HOST_IPV6_ADDRESS = "";
   public static final String ID_SPARK_TEXTINPUT_HOST_IPV6_ADDRESS_PREFIX = "";

   // DRS Cluster
   public static final String ID_ADD_CLUSTER_RULE_BUTTON = "addClusterRule";
   public static final String ID_BUTTON_EDIT_CLUSTER_RULE = "editClusterRule";
   public static final String ID_CLUSTER_RULE_TYPE_DROPDOWN_LIST = "ruleType";
   public static final String ID_CLUSTER_RULE_VM_TO_HOSTS_RELATION_TYPE_DROPDOWN_LIST =
         "relationType";
   public static final String ID_CLUSTER_RULE_NAME_SPARK_TEXT_INPUT = "ruleName";
   public static final String ID_CREATE_DRS_RULE_PANEL = "tiwoDialog";
   public static final String ID_CREATE_DRS_RULE_CANCEL_BUTTON =
         ID_CREATE_DRS_RULE_PANEL + "/cancelButton";
   public static final String ID_CREATE_DRS_RULE_OK_BUTTON = ID_CREATE_DRS_RULE_PANEL
         + "/okButton";
   public static final String ID_DRS_RULE_MEMBER_BUTTON = "addRuleMember";
   public static final String ID_ADD_DRS_RULE_MEMBER_OK_BUTTON = "okButton";
   public static final String ID_BUTTON_ADD_DRS_RULE_MEMBER = "addRuleMember";
   public static final String ID_ADD_DRS_RULE_MEMBER_CANCEL_BUTTON = "buttonCancel";
   public static final String ID_CLUSTER_CONFIGURATION_TOC = "tocListConfiguration";
   public static final String ID_DRS_RULES_ADVANCE_DATAGRID = "rulesList";
   public static final String ID_DRS_RULES_REMOVE_BUTTON = "removeClusterRule";
   public static final String ID_CONFIRM_DIALOG = "confirmationDialog";
   public static final String ID_CONFIRMATION_YES_BUTTON = ID_CONFIRM_DIALOG
         + "/automationName=Yes";
   public static final String ID_CONFIRMATION_NO_BUTTON = ID_CONFIRM_DIALOG
         + "/automationName=No";
   public static final String ID_CHECKBOX_RULE_MEMBER =
         "_CheckBoxColumnRenderer_CheckBox1";
   public static final String ID_ALERT_TURN_OFF_DRS = "YesNoDialog";
   public static final String ID_CHECKBOX_ENABLE_DRS = "enableDrsCheckBox";
   public static final String ID_BUTTON_BROWSE = "browseButton";
   public static final String ID_BUTTON_CLOSE = "_closeButton";
   public static final String ID_LIST_CONFLICT = "conflictList";
   public static final String ID_BUTTON_RULE_DETAILS = "clusterRuleMemberDetails";
   public static final String ID_LABEL_CONFLICT = "_VmRuleMemberDetailsView_Label1";
   public static final String ID_LABEL_CONFLICT_DETAILS =
         "_VmRuleMemberDetailsView_Label1";
   public static final String ID_LABEL_MEMBER_OF = "_VmRuleMemberDetailsView_Label2";
   public static final String ID_LIST_RULE_MEMBER_CONFLICT = "conflictsList";
   public static final String ID_LIST_RULE_MEMBER = "rulesList";
   public static final String ID_BUTTTON_GO_TO = "goToBtn";
   public static final String ID_BUTTON_CLOSE_RULE_MEMBER = "closeBtn";
   public static final String ID_GROUP_RULE_MEMBER = "contents";
   public static final String ID_REMOVE_DRS_RULE_MEMBER = "removeClusterRuleMember";
   public static final String ID_DRS_RULE_ENABLED_CHECKBOX = "ruleEnabled";
   public static final String ID_CHECKBOX_OVERRIDE_RECOMMENDATIONS =
         "_CheckBoxSkin_Group1";
   public static final String ID_LIST_RULE_MEMBER_CONTROL =
         "ruleConfigControlPlaceholder";
   public static final String ID_BUTTON_ADD_RULE_MEMBER = "addClusterRuleMember";
   public static final String ID_DIALOG_OBJECT_SELECTOR = "objectSelectorDialog";
   public static final String ID_BUTTON_REMOVE_RULE_MEMBER = "removeRuleMember";
   public static final String ID_IMAGE_DRS_ADV_OPTIONS_BLOCK =
         "arrowImagedrsAdvOptionsBlock";
   public static final String ID_BUTON_DELETE_OPTION = "deleteOptionBtn";
   public static final String ID_ADVANCED_DATA_GRID_VM_RULE_DETAILS = "vmList";
   public static final String ID_PORTLET_DATASTORE_CLUSTER_SUMMARY_STORAGE_CAPABILITIES =
         "vsphere.core.spbm.storage.storageCapabilitiesPortlet.chrome";
   public static final String ID_PORTLET_DATASTORE_CLUSTER_SUMMARY_SERVICES =
         "vsphere.core.dscluster.summary.servicesView.chrome";
   public static final String ID_PORTLET_DATASTORE_CLUSTER_SUMMARY_CONSUMERS =
         "vsphere.core.dscluster.summary.consumersView.chrome";
   public static final String ID_PORTLET_DATASTORE_CLUSTER_SUMMARY_RESOURCES =
         "vsphere.core.dscluster.summary.resourcesView.chrome";
   public static final String ID_ADVANCED_DATA_GRID_STORAGE_DRS_FAULT =
         "vsphere.core.dscluster.monitor.sdrs.faultView/faultList";
   public static final String ID_ADVANCED_DATA_GRID_STORAGE_DRS_FAULT_DETAILS =
         "vsphere.core.dscluster.monitor.sdrs.faultView/faultDetailList";
   // DRS VM Overrides
   public static final String ID_ADD_BTN = "addBtn";
   public static final String ID_CLUSTER_SELECT_VM_DIALOG = "objectSelectorDialog";
   public static final String ID_CLUSTER_VM_OVERRIDES_DELETE_BUTTON = "deleteBtn";
   public static final String ID_CLUSTER_VM_OVERRIDES_EDIT_BUTTON = "editBtn";
   public static final String ID_CLUSTER_VM_OVERRIDES_PANEL_DESELECT_BUTTON =
         new StringBuffer(ID_CLUSTER_EDIT_SETTINGS_PANEL).append("/deselectVmButton")
               .toString();
   public static final String ID_CLUSTER_VM_OVERRIDES_PANEL_SELECT_BUTTON =
         new StringBuffer(ID_CLUSTER_EDIT_SETTINGS_PANEL).append("/selectVmButton")
               .toString();
   public static final String ID_CLUSTER_VM_OVERRIDES_PANEL_OK_BUTTON =
         new StringBuffer(ID_CLUSTER_EDIT_SETTINGS_PANEL).append("/okButton").toString();
   public static final String ID_CLUSTER_VM_OVERRIDES_PANEL_CANCEL_BUTTON =
         new StringBuffer(ID_CLUSTER_EDIT_SETTINGS_PANEL).append("/cancelButton")
               .toString();
   public static final String ID_CLUSTER_VM_OVERRIDES_VM_SELECTION_CANCEL_BUTTON =
         new StringBuffer(ID_CLUSTER_SELECT_VM_DIALOG).append("/cancelButton")
               .toString();
   public static final String ID_CLUSTER_VM_OVERRIDES_PANEL_AUTO_LEVEL_COMBOBOX =
         "automationLevelInput";
   public static final String ID_CLUSTER_VM_OVERRIDES_PANEL_RESTART_PRIORITY_COMBOBOX =
         "_VmOverridesConfigForm_DropDownList2";
   public static final String ID_CLUSTER_VM_OVERRIDES_PANEL_ISOLATION_RESPONSE_COMBOBOX =
         "_VmOverridesConfigForm_DropDownList3";
   public static final String ID_CLUSTER_VM_OVERRIDES_PANEL_MONITORING_SENSITIVITY_COMBOBOX =
         "_VmOverridesConfigForm_DropDownList4";
   public static final String ID_VM_OVERRIDE_SELECT_VM_DATAGRID =
         IDConstants.ID_CLUSTER_SELECT_VM_DIALOG + "/VirtualMachine_filterGridView";
   public static final String ID_CLUSTER_ADD_VM_OVERRIDES_PANEL_AUTO_LEVEL_COMBOBOX =
         "automationLevelInput";
   public static final String ID_CLUSTER_ADD_VM_OVERRIDES_PANEL_RESTART_PRIORITY_COMBOBOX =
         "restartPriorityInput";
   public static final String ID_CLUSTER_ADD_VM_OVERRIDES_PANEL_ISOLATION_RESPONSE_COMBOBOX =
         "isolationResponseInput";
   public static final String ID_CLUSTER_ADD_VM_OVERRIDES_PANEL_VM_MONITORING_COMBOBOX =
         "vmMonitoringInput";
   public static final String ID_CLUSTER_ADD_VM_OVERRIDES_PANEL_AUTO_LEVEL_COMBOBOX_LABEL =
         ID_CLUSTER_ADD_VM_OVERRIDES_PANEL_AUTO_LEVEL_COMBOBOX + "/labelDisplay";
   public static final String ID_CLUSTER_ADD_VM_OVERRIDES_FT_VMS_WARNING_TEXT =
         "drsAutoLevForFtSecVms/errWarningText";
   public static final String ID_CLUSTER_ADD_VM_OVERRIDES_FT_VMS_WARNING_TEXT_EVC_DISABLED =
         "drsAutoLev/errWarningText";
   public static final String ID_CLUSTER_ADD_VM_OVERRIDES_VM_NAME_DATA_PROVIDER_PROPERTY =
         "vmName";
   public static final String ID_LABEL_VM_OVERRIDES_VM_ATUOMATION =
         "_VmOverridesConfigForm_Label15";

   // Storage extension - VAAI
   public static final String ID_LABEL_HARDWARE_ACCERLERATION =
         "titleLabel.hardwareAccelerationStackBlock";

   // Storage extension - VASA
   public static final String ID_STORAGE_DEVICE = "2";
   public static final String ADD_STORAGE_PROVIDER_BUTTON =
         "vsphere.core.storage.storageProviders.addNew/button";
   public static final String ID_BUTTON_SYNC_STORAGE_PROVIDER =
         "vsphere.core.storage.storageProviders.sync/button";
   public static final String REMOVE_STORAGE_PROVIDER_BUTTON =
         "vsphere.core.storage.storageProviders.remove/button";
   public static final String STORAGE_PROVIDER_LIST = "providersList";
   public static final String NAME_STORAGE_PROVIDER_TEXT = "providerName";
   public static final String URL_STORAGE_PROVIDER_TEXT = "providerUrl";
   public static final String USER_STORAGE_PROVIDER_TEXT = "username";
   public static final String PASSWD_STORAGE_PROVIDER_TEXT = "password";
   public static final String ID_CHECKBOX_USE_CERT = "useCustomCertificate";
   public static final String ID_BUTTON_HOST_SCSI_VOLUME_SYNC =
         "hostScsiVolumeSyncButton";
   public static final String ID_BUTTON_HOST_DATASTORE_SYNC = "hostDatastoreSyncButton";
   public static final String ID_BUTTON_HOST_VM_SYNC = "hostVmSyncButton";
   public static final String ID_BUTTON_DATASTORE_VM_SYNC = "datastoreVmSyncButton";
   public static final String ID_BUTTON_DATASTORE_SCSI_VOLUME_SYNC =
         "datastoreScsiVolumeSyncButton";
   public static final String ID_BUTTON_VM_DATASTORE_SYNC = "vmDatastoreSyncButton";
   public static final String ID_BUTTON_VM_SCSI_VOLUME_SYNC = "vmScsiVolumeSyncButton";
   public static final String ID_TEXT_DS_SUMMARY_USER_CAPABILITY = "userCapabilityName";

   // Storage extension - SPBM
   public static final String ID_ACTION_ASSIGN_STORAGE_PROFILE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_SPBM + "assignStorageCapabilitiesAction";
   public static final String ID_ACTION_EDIT_STORAGE_PROFILE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_SPBM + "editStorageProfileAction";
   public static final String ID_ACTION_REMOVE_STORAGE_PROFILE = EXTENSION_PREFIX
         + EXTENSION_ENTITY_SPBM + "removeStorageProfileAction";
   public static final String ID_ALERT_REMOVE_STORAGE_CAPABILITY =
         "removeStorageCapabilityQuestion";
   public static final String ID_BUTTON_ADD_STORAGE_PROFILE =
         "vsphere.core.spbm.profileList.createStorageProfileAction/button";
   public static final String ID_BUTTON_EDIT_STORAGE_PROFILE =
         "vsphere.core.spbm.profileList.editStorageProfileAction/button";
   public static final String ID_BUTTON_ENABLE_STORAGE_PROFILE =
         "vsphere.core.spbm.enableStorageProfilesAction/button";
   public static final String ID_BUTTON_MANAGE_STORAGE_CAPABILITIES =
         "vsphere.core.spbm.manageStorageCapabilitiesAction/button";
   public static final String ID_BUTTON_REMOVE_STORAGE_PROFILE =
         "vsphere.core.spbm.profileList.removeStorageProfileAction/button";
   public static final String ID_BUTTON_ADD_STORAGE_CAPABILITIES =
         "vsphere.core.spbm.createStorageCapabilitiesAction/button";
   public static final String ID_BUTTON_EDIT_STORAGE_CAPABILITIES =
         "vsphere.core.spbm.editStorageCapabilitiesAction/button";
   public static final String ID_BUTTON_REMOVE_STORAGE_CAPABILITIES =
         "vsphere.core.spbm.removeStorageCapabilitiesAction/button";
   public static final String ID_BUTTON_SELECT_ALL_CAPABILITIES = "selectAllButton";
   public static final String ID_BUTTON_CLOSE_MANAGE_CAPABILITIES = "closeFormButton";
   public static final String ID_FORM_MANAGE_STORAGE_CAPABILITIES =
         "automationName=Create Storage Capability";
   public static final String ID_LABEL_SUMMARY_STORAGE_PROFILE_NAME =
         "summary_profileName_valueLbl";
   public static final String ID_LABEL_SUMMARY_STORAGE_PROFILE_DESCRIPTION =
         "summary_profileDescripion_valueLbl";
   public static final String ID_LABEL_SUMMARY_HEADER_TOTAL_ENTITIES =
         "summaryHeader/_ProfileSummaryControl_Label3";
   public static final String ID_LIST_STORAGE_PROFILE_LIST = "profileList";
   public static final String ID_LIST_STORAGE_CAPABILITIES = "storageCapabilities";
   public static final String ID_LIST_CAPABILITY_LIST = "list";
   public static final String ID_LIST_SUMMARY_STORAGE_CAPABILITIES_GRID =
         "capabilitiesGrid";
   public static final String ID_TEXT_STORAGE_PROFILE_NAME = "profileName";
   public static final String ID_TEXT_STORAGE_PROFILE_DESCRIPTION = "profileDescription";
   public static final String ID_TEXT_STORAGE_CAPABILITY_NAME = "capabilityName";
   public static final String ID_TEXT_STORAGE_CAPABILITY_DESCRIPTION =
         "capabilityDescription";
   public static final String ID_LINK_EDIT_STORAGE_PROFILE =
         "vsphere.core.spbm.editStorageProfileAction";
   public static final String ID_LINK_GETTING_STARTED_LEARN_MORE_1 =
         "vsphere.core.sp.gettingStarted/gettingStartedHelpLink_1";
   public static final String ID_LINK_GETTING_STARTED_LEARN_MORE_2 =
         "vsphere.core.sp.gettingStarted/gettingStartedHelpLink_2";
   public static final String ID_VIEW_MANAGE_STORAGE_CAPABILITIES =
         "automationName=Manage Storage Capabilities";
   public static final String ID_DROP_DOWN_LIST_CAPABILITY_NAME = "capabilitiesDropdown";
   public static final String ID_BUTTON_CREATE_NEW_CAPABILITY =
         "createNewCapabilityButton";
   public static final String ID_FORM_CREATE_STORAGE_CAPABILITY =
         "_CreateStorageCapabilityForm_FormItem1";
   public static final String ID_TEXT_USER_CAPABILITY_NAME = "userCapabilityName";
   public static final String ID_BUTTON_MANAGE_STORAGE_PROFILE =
         "btn_vsphere.core.spbm.manageVmStorageProfilesAction";
   public static final String ID_BUTTON_PROPAGATE_TO_DISK = "propagateButton";
   public static final String ID_LIST_STORAGE_PROFILE = "spDropdown";
   public static final String ID_LIST_VM_LIST = "list";
   public static final String ID_BUTTON_RELATED_OBJECTS_VM =
         "vmsForStorageProfiles.button";
   public static final String ID_BUTTON_RELATED_OBJECTS_VM_TEMPLATES =
         "vmtemplsForStorageProfiles.button";
   public static final String ID_BUTTON_CHECK_STORAGE_PROFILE_COMPLIANCE =
         "vsphere.core.spbm.checkStorageProfileComplianceAction/button";

   // Storage extension - vm Provision with storage profile
   public static final String ID_LIST_STORAGE_PROFILE_SELECTOR =
         "storageProfileSelector";
   public static final String ID_TEXT_VM_SUMMARY_COMPLIANCE_STATE = "complianceText";
   public static final String ID_BUTTON_MANAGE_PROFILES = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "manage.profilesView.button";
   public static final String ID_TEXT_HOME_VM_PROFILE =
         "viewVmProfileAssignments.homeProfileLabel.text";
   public static final String ID_TEXT_HARD_DISK_1_PROFILE = "Hard disk 1.gridRow.text";
   public static final String ID_TEXT_HARD_DISK_2_PROFILE = "Hard disk 2.gridRow.text";
   public static final String ID_LABEL_VALUE_LABEL = "valueLabel";
   public static final String ID_LIST_VM_STORAGE_PROFILE_COMPLIANCE_LIST =
         "complianceList";
   public static final String ID_LABEL_ROLLUP_COMPLIANCE = "ProfileSummaryItem2" + "/"
         + ID_LABEL_VALUE_LABEL;
   public static final String ID_LABEL_ROLLUP_NON_COMPLIANCE = "ProfileSummaryItem0"
         + "/" + ID_LABEL_VALUE_LABEL;
   public static final String ID_DATAGRID_DISKS_GRID = "disksGrid";
   public static final String ID_TEXT_VM_SUMMARY_STORAGE_PROFILES = "profilesText";
   public static final String ID_LABEL_STORAGE_PROFILES_VALUE = "storageProfilesValue";
   public static final String ID_LABEL_STORAGE_PROFILES_VALUE_UNDER_NEW_DISK =
         "newDiskStorageProfile_0_Value";

   // Storage extension - storage host profile
   public static final String ID_TEXT_PSP_NAME = "pspName";
   public static final String ID_TEXT_BYTES = "bytesBeforeSwitchingToNextPath";
   public static final String ID_TEXT_RULE_NUMBER = "ruleNumber";
   public static final String ID_TEXT_RULE_CLASS = "claimruleClass";
   public static final String ID_DROP_DOWN_LIST_CLAIM_TYPE = "PolicyOptionDropDown-1";
   public static final String ID_TEXT_VENDOR_NAME = "vendorName";
   public static final String ID_TEXT_DEVICE_NAME = "deviceName";
   public static final String ID_TEXT_DRIVER_NAME = "driverName";
   public static final String ID_TEXT_WWNN = "wwnn";
   public static final String ID_TEXT_WWPN = "wwpn";
   public static final String ID_TEXT_LUN = "lun";
   public static final String ID_TEXT_IQN = "iqn";
   public static final String ID_TEXT_ADAPTER_NAME = "adapterName";
   public static final String ID_TEXT_CHANNEL = "channel";
   public static final String ID_TEXT_TARGET = "target";
   public static final String ID_TEXT_TRANSPORT_NAME = "transportName";

   // SDRS
   public static final String SDRS_INDEX = "1";
   public static String ID_TEXT_INPUT_LABEL = "inputLabel";
   public static String ID_HOST_AND_CLUSTER_LIST =
         "_DsClusterHostsAndClustersPage_List1";
   public static String ID_HOST_LIST = "HostSystem_filterGridView";
   public static String ID_DATAGRID_CLUSTER_LIST =
         "ClusterComputeResource_filterGridView";
   public static String ID_TURN_ON_STORAGE_DRS = "sdrsEnabled";
   public static String ID_SELECT_VM_OVERRIDE_BUTTON = "selectVmButton";
   public static String ID_BUTTON_OK = "buttonOk";
   public static String ID_AUTOMATION_LEVEL_DROPDOWN = "sdrsAutomationDropDown";
   public static String ID_COMBOBOX_INTRA_VM_AFFINITY = "intraVmAffinityDropDown";
   public static String ID_VM_LIST = "vmList";
   public static String ID_LIST_VMDKS = "vmdksList";
   public static String ID_ADD_BUTTON = "addBtn";
   public static String ID_DELETE_BTN = "deleteBtn";
   public static String ID_VM_OVERRIDES_GRID = "vmOverridesGrid";
   public static String ID_CONFIGURATION_LIST = "tocListConfiguration";
   public static String ID_TOC_LIST_OTHER = "tocListOther";
   public static String ID_REFRESH_RECOMMENDATION_ACTION =
         "refreshRecommendationsAction";
   public static String ID_APPLY_RECOMMENDATION_ACTION = "applyRecommendationsAction";
   public static String ID_BUTTON_APPLY = "buttonApply";
   public static String ID_CHECKBOX_OVERRIDE_RECOMMND = "overrideSwitchedOn";
   public static String ID_SDRS_RECOMMENTDATION_LIST = "sdrsRecList";
   public static final String ID_GRID_VMS_GRID =
         "vsphere.core.datastore.related/tabViews/vmsForDatastore/list";
   public static final String ID_GRID_VMS_CLUSTER_GRID =
         "vsphere.core.datastore.related/tabViews/vmsForDatastoreCluster/list";
   public static String ID_CHECKBOX_CHECK_UNCHECK_RECOMMND = "chbApply";
   public static final String ID_DS_CLUSTER_MANAGE_SETTINGS_SERVICES = "tocListServices";
   public static final String ID_DS_CLUSTER_EDIT_TURN_SDRS_ON = "sdrsTurnedOnChkb";
   public static final String ID_EDIT_DS_CLUSTER_BUTTON =
         "btn_vsphere.core.dscluster.configureDatastoreClusterAction";
   public static final String ID_SDRS_LABEL = "_SdrsConfigView_Label1";
   public static final String ID_DATAGRID_VIRTUAL_MACHINE_OVERRIDE =
         "VirtualMachine_filterGridView";
   public static final String ID_RADIO_BUTTON_AUTOMATION_LEVEL_FULLY_AUTOMATED =
         "_DsClusterAutomationLevelPage_RadioButton2";
   public static final String ID_LABEL_SDRS_IS_TURNED_ON = "_SdrsConfigView_Label1";
   public static final String ID_LABEL_DATASTORE_CLUSTER_NAME = "promptLabel";
   public static final String ID_LABEL_WIZARD_HOST = "promptLabel";
   public static final String ID_LABEL_WIZARD_CREATE_VM =
         "_SelectCreationTypeProvisioningPage_Label1";

   public static final String ID_TOCTREE_SDRS_RECOMMENDATIONS =
         "vsphere.core.dscluster.monitor.sdrsView/tocTree";
   public static final String ID_TOCTREE_SDRS_RULES =
         "vsphere.core.dscluster.manage.settingsView/tocTree";
   public static final String ID_LABEL_LOCATION = "locationLabel";
   public static final String ID_LABEL_ENTITY_NAME = "entityName";
   public static String ID_LABEL_SDRS_TURNED_ON = ID_TURN_ON_STORAGE_DRS
         + "/labelDisplay";
   public static final String ID_LABEL_CREATE_DATASTORE_CLUSTER_GENERAL =
         "_DsClusterNameAndLocationPage_Text1";
   public static final String ID_LABEL_SDRS_AUTOMATION =
         "titleLabel.sdrsAutomationBlock";
   public static final String ID_RADIO_BUTTON_AUTOMATION_LEVL_NO_AUTOMATION =
         "_DsClusterAutomationLevelPage_RadioButton1";
   public static String ID_LABEL_NO_AUTOMATION =
         ID_RADIO_BUTTON_AUTOMATION_LEVL_NO_AUTOMATION + "/labelDisplay";
   public static String ID_LABEL_FULLY_AUTOMATED =
         ID_RADIO_BUTTON_AUTOMATION_LEVEL_FULLY_AUTOMATED + "/labelDisplay";
   public static final String ID_LABEL_NO_AUTOMATION_EXPLANATION =
         "_DsClusterAutomationLevelPage_Text1";
   public static final String ID_LABEL_FULLY_AUTOMATED_EXPLANATION =
         "_DsClusterAutomationLevelPage_Text2";
   public static final String ID_LABEL_NO_ADVANCED_OPTIONS =
         "_DsClusterAutomationLevelPage_Label1";
   public static final String ID_LABEL_ADVANCED_OPTIONS =
         "titleLabel.sdrsAdvOptionsBlock";
   public static final String ID_CHECKBOX_IO_LOAD_BALANCE = "ioLoadBalanceEnabled";
   public static final String ID_STACK_EDITOR_ADVANCED_OPTION =
         "_DsClusterRuntimeSettingsPage_StackBlock1";

   public static final String ID_ARROW_IMAGE =
         "arrowImage_DsClusterRuntimeSettingsPage_StackBlock1";
   public static final String ID_NUMERIC_STEPPER_SPACE_UTILIZATION_THRESHOLD =
         "spaceUtilizationThreshold";
   public static final String ID_NUMERIC_STEPPER_IO_LATENCY_THRESHOLD =
         "ioLatencyThreshold";
   public static final String ID_NUMERIC_STEPPER_MIN_SPACE_UTILIZATION_DIFFERENCE =
         "minSpaceUtilizationDifference";
   public static final String ID_CLUSTER_CONNECTIVITY = "crList";

   public static final String ID_LABEL_IO_METRIC_INCLUSION =
         "_DsClusterRuntimeSettingsPage_SettingsBlock1/titleGroup/titleLabel";
   public static final String ID_LABEL_IO_METRIC_INCLUSION_EXPLANATION =
         "_DsClusterRuntimeSettingsPage_Label1";
   public static final String ID_LABEL_IO_ENABLE_IO_METRIC_EXPLANATION =
         "_DsClusterRuntimeSettingsPage_Text1";
   public static final String ID_LABEL_SDRS_THRESHOLD =
         "_DsClusterRuntimeSettingsPage_SettingsBlock2/titleGroup/titleLabel";
   public static final String ID_LABEL_SDRS_THRESHOLD_EXPLANATION =
         "_DsClusterRuntimeSettingsPage_Label2";

   /* Id's on the Ready to Complete page of Creare New Datastore Cluster */
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_GENERAL =
         "_DsClusterReadyToCompletePage_Label1";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_DATASTORE_CLUSTRE_NAME =
         "_DsClusterReadyToCompletePage_Label2";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_NAME_VALUE =
         "_DsClusterReadyToCompletePage_Label3";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_STORAGE_DRS =
         "_DsClusterReadyToCompletePage_Label4";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_STORAGE_DRS_VALUE =
         "_DsClusterReadyToCompletePage_Label5";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_SDRS_AUTOMATION =
         "_DsClusterReadyToCompletePage_Label6";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_AUTOMATION_LEVEL =
         "_DsClusterReadyToCompletePage_Label7";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_AUTOMATION_LEVEL_VALUE =
         "_DsClusterReadyToCompletePage_Label8";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_SDRS_RUNTIME_SETTINGS =
         "_DsClusterReadyToCompletePage_Label9";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_STORAGE_IO =
         "_DsClusterReadyToCompletePage_Label10";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_STORAGE_IO_VALUE =
         "_DsClusterReadyToCompletePage_Label11";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_UTILIZED_SPACE =
         "_DsClusterReadyToCompletePage_Label12";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_UTILIZED_SPACE_VALUE =
         "_DsClusterReadyToCompletePage_Label13";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_IO_LATENCY =
         "_DsClusterReadyToCompletePage_Label14";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_IO_LATENCY_VALUE =
         "_DsClusterReadyToCompletePage_Label15";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_ADVANCED_OPTIONS =
         "_DsClusterReadyToCompletePage_Label16";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_UTILIZATION_DIFFERENCE =
         "_DsClusterReadyToCompletePage_Label19";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_UTILIZATION_DIFFERENCE_VALUE =
         "_DsClusterReadyToCompletePage_Label20";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_CHECK_IMBALANCE =
         "_DsClusterReadyToCompletePage_Label21";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_CHECK_IMBALANCE_VALUE =
         "_DsClusterReadyToCompletePage_Label22";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_IO_THRESHOLD =
         "_DsClusterReadyToCompletePage_Label23";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_IO_THRESHOLD_VALUE =
         "_DsClusterReadyToCompletePage_Label24";
   public static final String ID_LABEL_SDRS_READY_TO_COMPLETE_PAGE_DATASTORES =
         "_DsClusterReadyToCompletePage_Label25";
   public static final String ID_GRID_ADG_DATASTORES = "adgDatastores";
   public static final String ID_GRID_ADG_HOSTS_AND_CLUSTERS = "adgClustersAndHosts";

   public static final String ID_COMBOBOX_HOST_CONNECTION_STATUS_CONTROL =
         "hostConnectionStatusControl";
   public static final String ID_TOOGLE_BAR_BUTONN_FILTER_VIEW =
         "filterViewToggleButtonBar";
   public static final String ID_BUTTON_HOSTS = ID_TOOGLE_BAR_BUTONN_FILTER_VIEW
         + "/automationName=Standalone Hosts";
   public static final String ID_BUTTON_CLUSTERS = ID_TOOGLE_BAR_BUTONN_FILTER_VIEW
         + "/automationName=Clusters";

   public static final String ID_TEXT_WARNING_MESSAGE = "warningArea/_message";
   public static final String ID_LABEL_DETAILS_HINT = "detailsHint";
   public static final String ID_LABEL_PREVIEW_CONTAINER_TITLE_LABEL =
         "previewContainer/titleLabel";
   public static final String ID_LABEL_LIST_CONTAINER_TITLE_LABEL =
         "listContainer/titleLabel";
   public static final String ID_LABEL_DATASTORES_SELECTED =
         "datastoreList/dataGridStatusbar/infoPanel2";

   public static final String ID_BUTTON_DROP_DOWN_ARROW = "dropdownArrowButton";
   public static final String ID_LIST_FILTER_CONTEXT_MENU = "filterContextMenu";
   public static final String ID_COMPONENT_SEARCH_CONTROL = "searchControl";
   public static final String ID_CONTEXT_HEADER = "contextHeader";
   public static final String ID_CONTEXT_HEADER_LIST = "/list";
   public static final String ID_CONTEXT_HEADER_ALL_COLUMN_NAMES_PROPERTY =
         "dataProvider.source.source.source.[].name";
   public static final String ID_FILTER_CONTEXT_HEADER = "filterContextHeader";
   public static final String ID_STACK_BLOCK_SDRS_ADV_OPTIONS = "sdrsAdvOptionsBlock";
   public static final String ID_STACK_BLOCK_SRDS_IO_METRICS = "sdrsIOMetricBlock";
   public static final String ID_STACK_BLOCK_STORAGE_DRS_AUTOMATION =
         "sdrsAutomationBlock";
   public static final String ID_STACK_BLOCK_AUTOMATION_LEVEL = "autLevelBlock";
   public static final String ID_STACK_BLOCK_IO_METRICS = "ioMetricsBlock";
   public static final String ID_STACK_BLOCK_SDRS_CONFIG_3 =
         "_SdrsConfigForm_StackBlock3";
   public static final String ID_IMAGE_SDRS_ADV_OPTIONS_BLOCK =
         "arrowImagesdrsAdvOptionsBlock";
   public static final String ID_BUTTON_ADD_OPTION = "addOptionBtn";
   public static final String ID_GRID_ADVANCED_OPTIONS = "advancedOptionsGrid";
   public static final String ID_GRID_HORIZONTAL_SCROLL = "className=HScrollBar";
   public static final String ID_GRID_HORIZONTAL_SCROLL_RIGHT_BUTTON =
         ID_GRID_HORIZONTAL_SCROLL + "/name=downArrowSkin";
   public static final String ID_GRID_HORIZONTAL_SCROLL_LEFT_BUTTON =
         ID_GRID_HORIZONTAL_SCROLL + "/name=upArrowSkin";
   public static final String ID_CREATE_VDS_NAME_AND_LOCATION_PAGE =
         "tiwoDialog/nameAndLocationPage";
   public static final String ID_NAME_AND_LOCATION_PAGE_SEARCH_GRID =
         ID_CREATE_VDS_NAME_AND_LOCATION_PAGE + "/" + ID_SIMPLE_SEARCH_RESULTS_GRID;
   public static final String ID_BUTTON_DS_CLLUSTER_CREATE =
         "vsphere.core.dscluster.createActionGlobal/button";
   public static final String ID_TEXT_SEARCH_DATASTORE_CLUSTER_LOCATION =
         "nameAndLocationPage/entityLocationSelector/searchControl/_SearchControl_HBox4/searchInput";
   public static final String ID_BUTTON_SEARCH_DATASTORE_CLUSTER_LOCATION =
         "nameAndLocationPage/entityLocationSelector/searchControl/searchStack/searchButtonContainer/searchButton";
   public static final String ID_LIST_STORAGE_POD_VIEW = "StoragePod/list";
   public static final String ID_BUTTON_DS_CLUSTER_CREATE =
         "vsphere.core.dscluster.createAction/button";

   // IDs for SDRS page and for Edit cluster page
   public static final String ID_IMAGE_ARROW_AUT_LEVEL_BLOCK = "arrowImageautLevelBlock";
   public static final String ID_IMAGE_ARROW_IO_METRICS_BLOCK =
         "arrowImageioMetricsBlock";
   public static final String ID_IMAGE_ARROW_ADVANCED_OPTIONS_BLOCK =
         "arrowImage_SdrsConfigForm_StackBlock3";
   public static final String ID_LISTBOX_SDRS_AUTOMATION_LEVEL = "sdrsAutLevelDropDown";
   public static final String ID_CHECKBOX_IO_METRIC = "_SdrsConfigForm_CheckBox2";
   public static final String ID_CHECKBOX_IO_METRIC_SECTION_EXPANDED =
         "_SdrsConfigForm_CheckBox3";
   public static final String ID_NUMERIC_STEPPER_EDIT_UTILIZED_SPACE_THRESHOLD =
         "_SdrsConfigForm_NumericStepper1";
   public static final String ID_NUMERIC_STEPPER_EDIT_IO_LATENCY =
         "_SdrsConfigForm_NumericStepper2";
   public static final String ID_NUMERIC_STEPPER_EDIT_MIN_UTILIZATION_DIFFERENCE =
         "_SdrsConfigForm_NumericStepper3";
   public static final String ID_LISTBOX_LOAD_BALANCE_INTERVAL_CONTROL =
         "loadBalanceIntervalControl";
   public static final String ID_LISTBOX_LOAD_BALANCE_INTERVAL_UNIT_CONTROL =
         "loadBalanceIntervalUnitControl";
   public static final String ID_LABEL_SDRS_CONFIG_VIEW2 = "_SdrsConfigView_Label2";
   public static final String ID_IMAGE_ARROW_DRS_AUTOMATION_BLOCK =
         "arrowImagesdrsAutomationBlock";
   public static final String ID_IMAGE_ARROW_DRS_IO_METRIC_BLOCK =
         "arrowImagesdrsIOMetricBlock";
   public static final String ID_LABEL_SDRS_CONFIG_VIEW3 = "_SdrsConfigView_Text3";
   public static final String ID_LABEL_SDRS_CONFIG_VIEW_TEXT2 = "_SdrsConfigView_Text2";
   public static final String ID_LABEL_SDRS_CONFIG_VIEW_TEXT4 = "_SdrsConfigView_Text4";
   public static final String ID_LABEL_SDRS_CONFIG_VIEW_TEXT5 = "_SdrsConfigView_Text5";
   public static final String ID_LABEL_SDRS_CONFIG_VIEW_TEXT6 = "_SdrsConfigView_Text6";
   public static final String ID_LABEL_SDRS_CONFIG_VIEW_TEXT7 = "_SdrsConfigView_Text7";
   public static final String ID_LABEL_SDRS_CONFIG_VIEW_TEXT8 = "_SdrsConfigView_Text8";
   public static final String ID_LABEL_SDRS_CONFIG_VIEW_TEXT9 = "_SdrsConfigView_Text9";
   public static final String ID_SLIDER_SDRS_CONFIG_FORM3 = "_SdrsConfigForm_HSlider3";

   //sdrs scheduled task
   public static final String ID_BUTTON_DATASTORE_CLUSTER_SCHEDULED_TASK =
         "btn_vsphere.core.dscluster.scheduleConfigureDatastoreClusterAction";
   public static final String ID_TEXT_NAME = "nameTxt";
   public static final String ID_BUTTON_SCHEDULING_OPTIONS = "schedulingOptionsBtn";
   public static final String ID_PAGE_SCHEDULE_OPTION = "step_wizardScheduleOptionsPage";

   //SDRS faults
   public static final String ID_TEXT_FAULTS_HEADER = "faultsHeaderText";

   //edit vm settins
   public static final String ID_BUTTON_SDRS_RULES_STACK = "sdrsRulesStackButton";
   public static final String ID_BUTTON_ADD_RULE = "addRule";
   public static final String ID_BUTTON_EDIT_RULE = "editRule";
   public static final String ID_BUTTON_REMOVE_RULE = "removeRule";

   public static final String ID_BUTTON_ADD_VM_RULE_MEMBERS = "addVmRuleMembers";
   public static final String ID_BUTTON_REMOVE_VM_RULE_MEMBERS = "removeVmRuleMembers";
   public static final String ID_DATAGRID_VM_RULE_DETAILS = "vmRuleDetails";
   public static final String ID_DATAGRID_VMDK_RULE_DETAILS = "vmdkRuleDetails";


   // Host Authentication Services
   public static String ID_HOST_JOIN_DOMAIN = "btn_vsphere.core.host.authJoinDomain";
   public static String ID_HOST_LEAVE_DOMAIN = "btn_vsphere.core.host.authLeaveDomain";
   public static String ID_HOST_JOIN_DOMAIN_NAME = "domain";
   public static String ID_HOST_JOIN_DOMAIN_USERNAME = "username";
   public static String ID_HOST_JOIN_DOMAIN_PASSWORD = "password";
   public static String ID_HOST_LEAVE_DOMAIN_CONFIRMATION_YES_BUTTON =
         "confirmationDialog/automationName=Yes";
   public static String ID_HOST_LEAVE_DOMAIN_CONFIRMATION_NO_BUTTON =
         "confirmationDialog/automationName=No";
   public static final String HOSTS_FOLDER_CONTENTS_DATAPROVIDER_ID =
         "contentsForHandCFolder";
   public static String ID_RB_USING_CREDENTIALS = "_AuthJoinDomainForm_RadioButton1";
   public static String ID_RB_USING_PROXY = "camServerRb";
   public static String ID_ADDRESS_INPUT_PROXY_SERVER_IP = "camServerIp";
   public static String ID_ADD_INPUT_IMPORT_CERT_INPUT_SERVER_IP = "serverIp";

   // iSCSI Static and Dynamic Discovery constants
   public static final String ISCSI_TARGETS_ADD_BUTTON = "addButton";
   public static final String ISCSI_TARGETS_REMOVE_BUTTON = "removeButton";
   public static final String ISCSI_TARGETS_AUTHENTICATION_BUTTON = "authButton";
   public static final String ISCSI_TARGETS_ADVANCED_BUTTON = "advancedButton";
   public static final String ISCSI_TARGETS_RESCAN_ADAPTER_BUTTON =
         "rescanAdapterButton";
   public static final String ISCSI_SERVER_NAME = "iScsiServer";
   public static final String ISCSI_TARGET_NAME = "iScsiTarget";
   public static final String ISCSI_DYNAMIC_DISCOVERY_PORT = "port";
   public static final String ISCSI_ADD_DYNAMIC_TARGET_DIALOG = "tiwoDialog";
   //public static final String ID_PATH_ISCSI_STORAGE = "vsphere.core.host.manage.storage/";
   public static final String ISCSI_DYNAMIC_DISCOVERY_INHERITSETTINGS_CHECK =
         "inheritCheckBox";
   public static final String ISCSI_ADAPTER_DETAILS =
         "vsphere.core.host.configuration.storageAdaptersDetails";
   public static final String ISCSI_TARGETS_ID =
         "vsphere.core.storage.adaptersDetails.views.swIscsi.targets.container";
   public static final String ISCSI_PROPERTIES_ID =
         "vsphere.core.storage.adaptersDetails.views.properties.container";
   public static final String ISCSI_DEVICES_ID =
         "vsphere.core.storage.adaptersDetails.views.devices.container";
   public static final String ISCSI_DEVICES_TAB = "automationName=Devices";
   public static final String ISCSI_RESCAN_ADAPTERS_MESSAGE = "_message";
   public static final String ID_MESSAGE_CONTAINER_ERRORS = "messageArea";
   public static final String ISCSI_RESCAN_CLOSE = "_button";
   public static final String ISCSI_TAB_NAVIGATOR =
         "_StorageAdaptersDetailsView_TabbedView1/tabNavigator";
   public static final String ISCSI_TAB_NAVIGATOR_DYNAMIC_DISCOVERY =
         "targetsButtonBar/label=Dynamic Discovery";
   public static final String ISCSI_TAB_NAVIGATOR_STATIC_DISCOVERY =
         "targetsButtonBar/label=Static Discovery";
   public static final String ISCSI_STATIC_TARGETS = "staticTargetsList";
   public static final String ISCSI_DYNAMIC_TARGETS = "dynamicTargetsList";
   public static final String ISCSI_TARGET_ADDRESS_DETAILS =
         "dataProvider.source.0.address.value";
   public static final String ISCSI_TARGET_PORT_DETAILS =
         "dataProvider.source.0.port.value";
   public static final String ISCSI_TARGET_IQN_DETAILS =
         "dataProvider.source.0.iScsiName.value";
   public static final String ISCSI_REMOVE_TARGETS_DIALOG = "YesNoDialog";
   public static final String ISCSI_REMOVE_TARGETS_CONFIRM_YES = "automationName=Yes";

   public static final String ID_LABEL_ISCSI_ADAPTERS_STATUS =
         "_StorageAdaptersPropertiesView_PropertyGridRow1";
   public static final String ID_LABEL_ISCSI_ADAPTERS_DEVICE =
         "storage.adapters.deviceName";
   public static final String ID_LABEL_ISCSI_ADAPTERS_DEVICE_VALUE =
         "storage.adapters.deviceName.text";
   public static final String ID_LABEL_ISCSI_ADAPTERS_MODEL_LABEL =
         "storage.adapters.deviceModel";
   public static final String ID_LABEL_ISCSI_ADAPTERS_MODEL_VALUE =
         "storage.adapters.deviceModel.text";
   public static final String ID_LABEL_ISCSI_ADAPTERS_NAME_LABEL =
         "storage.adapters.properties.general.iscsiName";
   public static final String ID_LABEL_ISCSI_ADAPTERS_NAME_VALUE =
         "storage.adapters.properties.general.iscsiName.text";
   public static final String ID_LABEL_ISCSI_ADAPTERS_ALIAS_LABEL =
         "storage.adapters.properties.general.iscsiAlias";
   public static final String ID_LABEL_ISCSI_ADAPTERS_ALIAS_VALUE =
         "storage.adapters.properties.general.iscsiAlias.text";
   public static final String ID_LABEL_ISCSI_ADAPTERS_DISCOVERY_LABEL =
         "storage.adapters.properties.general.iscsiTargets";
   public static final String ID_LABEL_ISCSI_ADAPTERS_DISCOVERY_VALUE =
         "storage.adapters.properties.general.iscsiTargets.text";
   public static final String ID_LABEL_ISCSI_ADAPTERS_AUTH_LABEL =
         "_StorageAdaptersPropertiesView_PropertyGridRow12";
   public static final String ID_LABEL_ISCSI_ADAPTERS_AUTH_VALUE =
         "authenticationBlockLabel";
   public static final String ID_LABEL_ISCSI_ADAPTERS_MAC_LABEL =
         "_StorageAdaptersPropertiesView_PropertyGridRow3";
   public static final String ID_LABEL_ISCSI_ADAPTERS_MAC_VALUE = "macAddressLabel";
   public static final String ID_LABEL_ISCSI_ADAPTERS_IP_LABEL =
         "_StorageAdaptersPropertiesView_PropertyGridRow4";
   public static final String ID_LABEL_ISCSI_ADAPTERS_IP_VALUE = "ipAddressLabel";//"storage.adapters.deviceName.text";
   public static final String ID_LABEL_ISCSI_ADAPTERS_SUBNET_LABEL =
         "_StorageAdaptersPropertiesView_PropertyGridRow5";
   public static final String ID_LABEL_ISCSI_ADAPTERS_SUBNET_VALUE = "subnetMaskLabel";
   public static final String ID_LABEL_ISCSI_ADAPTERS_DEFAULTGATEWAY_LABEL =
         "_StorageAdaptersPropertiesView_PropertyGridRow6";
   public static final String ID_LABEL_ISCSI_ADAPTERS_DEFAULTGATEWAY_VALUE =
         "defaultGatewayLabel";
   public static final String ID_LABEL_ISCSI_ADAPTERS_IPV6_LABEL =
         "_StorageAdaptersPropertiesView_PropertyGridRow7";
   public static final String ID_LABEL_ISCSI_ADAPTERS_IPV6_VALUE = "ipv6AddressLabel";
   public static final String ID_LABEL_ISCSI_ADAPTERS_SUBNETIPV6_LABEL =
         "_StorageAdaptersPropertiesView_PropertyGridRow8";
   public static final String ID_LABEL_ISCSI_ADAPTERS_SUBNETIPV6_VALUE =
         "subnetMaskv6Label";
   public static final String ID_LABEL_ISCSI_ADAPTERS_DGIPV6_LABEL =
         "_StorageAdaptersPropertiesView_PropertyGridRow9";
   public static final String ID_LABEL_ISCSI_ADAPTERS_DGIPV6_VALUE =
         "defaultGatewayv6Label";
   public static final String ID_LABEL_ISCSI_ADAPTERS_PREF_DNS_LABEL =
         "_StorageAdaptersPropertiesView_PropertyGridRow10";
   public static final String ID_LABEL_ISCSI_ADAPTERS_PREF_DNS_VALUE =
         "preferredDnsLabel";
   public static final String ID_LABEL_ISCSI_ADAPTERS_ALT_DNS_LABEL =
         "_StorageAdaptersPropertiesView_PropertyGridRow11";
   public static final String ID_LABEL_ISCSI_ADAPTERS_ALT_DNS_VALUE =
         "alternateDnsLabel";

   public static final String ID_BUTTON_ADD_ADAPTER = "addAdapterButton/button";
   public static final String ID_BUTTON_GENERAL_BLOCK = "generalBlockEditButton";
   public static final String ID_PATH_STORAGECONTAINER =
         "vsphere.core.storage.adaptersDetails.views.paths.container";
   public static final String ID_CONTAINER_ISCSI_ADVANCEDOPT =
         "vsphere.core.storage.adaptersDetails.views.swIscsi.advancedOptions.container";
   public static final String ID_CONATINER_ISCSI_PORTBINDING =
         "vsphere.core.storage.adaptersDetails.views.swIscsi.networkPortBinding.container";
   public static final String ID_DATAGRID_ISCSI_ADV_OPTIONS_PREFIX =
         "HostSystem:key-vim.host.InternetScsiHba-";
   public static final String ID_BUTTON_EDIT_ISCSI_ADV_OPTIONS_PREFIX =
         "Edit_HostSystem:key-vim.host.InternetScsiHba-";
   public static final String ID_DATAGRID_ISCSI_ADV_OPTIONS_POSTFIX =
         ".internetScsiAdvancedSettings";
   public static final String ID_BUTTON_ADV_OPTIONS_EDIT =
         "vsphere.core.host.manage.storage/actionId=com.vmware.resourcefw.defaultEditAction";
   public static final String ID_DATAGRID_ISCSI_TARGETS_ADVANCED_OPTIONS =
         "HostSystem:internetScsiTargetAdvancedSettings[object InternetScsiTargetsAdvancedSettingsRetrievalSpec]";
   public static final String ID_TEXT_ISCSI_NAME = "iScsiName";
   public static final String ID_TEXT_ISCSI_ALIAS = "iScsiAlias";

   // iSCSI Network Port Binding Dialog constants
   public static final String ID_DIALOG_VIEW_DETAILS =
         "storageAdaptersNetworkPortBindingDetailsView";
   public static final String ID_DIALOG_BIND_VNIC = "storageAdaptersAddForm";
   public static final String ID_TAB_NAVIGATOR_VNIC_SETTINGS = "vnicSettingsTabView";
   public static final String ID_BUTTON_ISCSI_NETWORK_PORT_BINDING_CLOSE_VIEW_DETAILS =
         "closeDetailsButton";
   public static final String ID_LIST_ISCSI_ADD_VMKERNEL_ADAPTER =
         "storageAdaptersAddList";
   public static final String ID_LABEL_VNIC_LIST_HEADER_INFO =
         "_StorageAdaptersAddForm_Text1";
   public static final String ID_LABEL_VNIC_EMPTY_DETAILS_AREA_INFO =
         "_StorageAdaptersDetailsTabs_Text1";
   public static final String ID_IMAGE_VNIC_EMPTY_DETAILS_AREA_INFO =
         "_StorageAdaptersDetailsTabs_BitmapImage1";
   public static final String ID_BUTTON_SHOWHIDE_CONTEXT_MENU_OK_BUTTON =
         "contextHeader/detailsPanel/commitPanel/okButton";
   public static final String ID_BUTTON_SHOWHIDE_CONTEXT_MENU_CLOSE_BUTTON =
         "contextHeader/caption/closeButton";

   public static final String ID_ISCSI_BIND_VMKERNEL_ADAPTER = "storageAdaptersAddForm";
   public static final String ID_LABEL_STATUS_ERROR_MESSAGE =
         "storageAdaptersNetworkPortBindingDetailsView/complianceStatusNavContent/propViewMessageText";

   public static final String ID_TAB_VNIC_SETTINGS_STATUS = "0";

   // iSCSI network port binding View constants
   public static final String ID_BUTTON_ISCSI_NETWORK_PORT_BINDING_ADD =
         "_StorageAdaptersNetworkPortBindingList_ActionfwButtonBar1/button[0]";
   public static final String ID_BUTTON_ISCSI_NETWORK_PORT_BINDING_REMOVE =
         "_StorageAdaptersNetworkPortBindingList_ActionfwButtonBar1/button[1]";
   public static final String ID_BUTTON_ISCSI_NETWORK_PORT_BINDING_VIEW_DETAILS =
         "_StorageAdaptersNetworkPortBindingList_ActionfwButtonBar1/button[2]";
   public static final String ID_GRID_ISCSI_NETWORK_PORT_BINDING =
         "storageAdaptersNetworkPortBindingList";
   public static final String ID_LABEL_ISCSI_NETWORK_PORT_BINDING_EMPTY_LIST_INDICATOR =
         "storageAdaptersNetworkPortBindingList/emptyListIndicator";

   // Edit Security Profile
   public static final String ID_CHECKBOX_SECURITY_PROFILE_SOFTWARE_ISCSI_CLIENT =
         "automationName=Software iSCSI Client";
   public static final String ID_LABEL_SECURITY_PROFILE_SOFTWARE_ISCSI_CLIENT =
         "outStackBlock/firewallOutPropGrid/automationName=Software iSCSI Client";
   public static final String ID_LIST_SECURITY_PROFILE_EDIT_FIREWALL =
         "editSecurityProfile/firewallList";

   public static final String ID_BLOCK_ISCSI_IP =
         "_EditIScsiIPSettingsForm_SettingsBlock1/titleLabel";
   public static final String ID_BLOCK_ISCSI_IPV6 =
         "_EditIScsiIPSettingsForm_SettingsBlock2/titleLabel";
   public static final String ID_BLOCK_ISCSI_DNS =
         "_EditIScsiIPSettingsForm_SettingsBlock3/titleLabel";

   public static final String ID_LABEL_AUTH_TEXT =
         "_EditIScsiAuthenticationSettingsForm_Text1";
   public static final String ID_LABEL_SUBNET_MASK_IPV6_VALUE =
         "subnetMaskV6NotSupported";
   public static final String ID_LABEL_DG_IPV6_VALUE = "defaultGatewayV6NotSupported";
   public static final String ID_BLOCK_OUTGOING_CHAP =
         "_InternetScsiAuthenticationForm_SettingsBlock1/titleLabel";
   public static final String ID_LABEL_OUTGOING_CHAP_NAME =
         "_InternetScsiAuthenticationForm_FormItem2";
   public static final String ID_LABEL_OUTGOING_CHAP_SECRET =
         "_InternetScsiAuthenticationForm_FormItem3";

   public static final String ID_BLOCK_INCOMING_CHAP =
         "_InternetScsiAuthenticationForm_SettingsBlock2/titleLabel";
   public static final String ID_LABEL_INCOMING_CHAP_NAME =
         "_InternetScsiAuthenticationForm_FormItem4";
   public static final String ID_LABEL_INCOMING_CHAP_SECRET =
         "_InternetScsiAuthenticationForm_FormItem5";


   public static final String ID_FORM_ISCSI_IP = "_EditIScsiIPSettingsForm_FormItem1";
   public static final String ID_FORM_ISCSI_SUBNET =
         "_EditIScsiIPSettingsForm_FormItem2";
   public static final String ID_FORM_ISCSI_DEFAULTGATEWAY =
         "_EditIScsiIPSettingsForm_FormItem3";
   public static final String ID_FORM_ISCSI_IPV6 = "_EditIScsiIPSettingsForm_FormItem4";
   public static final String ID_FORM_ISCSI_SUBNETIPV6 =
         "_EditIScsiIPSettingsForm_FormItem5";
   public static final String ID_FORM_ISCSI_DGIPV6 =
         "_EditIScsiIPSettingsForm_FormItem6";
   public static final String ID_FORM_ISCSI_PRE_DNS =
         "_EditIScsiIPSettingsForm_FormItem7";
   public static final String ID_FORM_ISCSI_ALT_DNS =
         "_EditIScsiIPSettingsForm_FormItem8";

   public static final String ID_LABEL_CURRENT_MAX_SPEED =
         "_StorageAdaptersPropertiesView_PropertyGridRow2";
   public static final String ID_LABEL_CURRENT_MAX_SPEED_VALUE = "currentMaxSpeedLabel";

   public static final String ID_BLOCK_NETWORK_INTERFACE = "networkInterfaceBlock";
   public static final String ID_BUTTON_EDIT_IP_ISCSI = "editIpSettingsButton";
   public static final String ID_TEXT_IP_ADDRESS = "ipAddress";
   public static final String ID_TEXT_SUBNET_MASK = "subnetMask";
   public static final String ID_TEXT_DEFAULT_GATEWAY = "defaultGateway";
   public static final String ID_LABEL_IPV4_ADDRESS = "ipAddressLabel";
   public static final String ID_LABEL_SUBNET_MASK_V4 = "subnetMaskLabel";
   public static final String ID_LABEL_DEFAULT_GATEWAY_V4 = "defaultGatewayLabel";
   public static final String ID_LABEL_IPV6_ADDRESS = "ipv6AddressLabel";
   public static final String ID_LABEL_SUBNET_MASK_V6 = "subnetMaskv6Label";

   public static final String ID_COMPONENT_DROP_SHADOW = "dropShadow";

   public static final String ID_LABEL_ISCSI_NAME =
         "storage.adapters.properties.general.iscsiName.text";
   public static final String ID_LABEL_ISCSI_ALIAS =
         "storage.adapters.properties.general.iscsiAlias.text";

   public static final String ID_FORM_ISCSI_NAME =
         "_EditIScsiGeneralSettingsForm_FormItem1";
   public static final String ID_FORM_ISCSI_ALIAS =
         "_EditIScsiGeneralSettingsForm_FormItem2";

   public static final String ID_GRID_ISCSI_PATHS =
         "vsphere.core.storage.adaptersDetails.views.paths/list";

   public static final String ID_FORM_ISCSI_AUTH_OUTGOING_NAME =
         "_InternetScsiAuthenticationForm_FormItem2";
   public static final String ID_FORM_ISCSI_AUTH_OUTGOING_SECRET =
         "_InternetScsiAuthenticationForm_FormItem3";
   public static final String ID_FORM_ISCSI_AUTH_INCOMING_NAME =
         "_InternetScsiAuthenticationForm_FormItem4";
   public static final String ID_FORM_ISCSI_AUTH_INCOMING_SECRET =
         "_InternetScsiAuthenticationForm_FormItem5";

   public static final String ID_TEXT_ISCSI_AUTH_OUTGOING = "outgoingName";
   public static final String ID_TEXT_ISCSI_AUTH_INCOMING = "incomingName";

   public static final String ID_TEXT_ISCSI_AUTH_OUTGOING_NAME =
         ID_FORM_ISCSI_AUTH_OUTGOING_NAME + "/" + ID_TEXT_ISCSI_AUTH_OUTGOING
               + "/textDisplay";
   public static final String ID_TEXT_ISCSI_AUTH_OUTGOING_SECRET =
         ID_FORM_ISCSI_AUTH_OUTGOING_SECRET + "/outgoingSecret/textDisplay";
   public static final String ID_TEXT_ISCSI_AUTH_INCOMING_NAME =
         ID_FORM_ISCSI_AUTH_INCOMING_NAME + "/" + ID_TEXT_ISCSI_AUTH_INCOMING
               + "/textDisplay";
   public static final String ID_TEXT_ISCSI_AUTH_INCOMING_SECRET =
         ID_FORM_ISCSI_AUTH_INCOMING_SECRET + "/incomingSecret/textDisplay";

   public static final String ID_LABEL_ISCSI_AUTH_BLOCK = "authenticationBlockLabel";
   public static final String ID_BUTTON_ISCSI_AUTH_EDIT = "authenticationButton";
   public static final String ID_LIST_ISCSI_AUTH_METHODS = "authMethods";
   public static final String ID_BLOCK_ISCSI_RESCAN_MESSAGE_AREA = "messageArea";
   public static final String ID_CHECKBOX_ISCSI_INITIATOR_OUTGOING = "outgoingCheckBox";
   public static final String ID_CHECKBOX_ISCSI_INITIATOR_INCOMING = "incomingCheckBox";

   // OVF constants
   public static final String ID_OVF_NAME_INPUT = "nameInput";
   public static final String ID_EXPORT_OVF_DIRECTORY = "chooseDirectoryButton";
   public static final String ID_EXPORT_OVF_FORMAT = "formatInput";
   public static final String ID_EXPORT_OVF_OVERWRITE = "overwriteCheckBox";
   public static final String ID_EXPORT_OVF_PROGRESS_CANCELBTN = "cancelBtn";
   public static final String ID_EXPORT_OVF_PROGRESS_CLOSEBTN = "closeBtn";
   public static final String ID_EXPORT_OVF_FILE_LOCATION = "directoryLabel";
   public static final String ID_CHECKBOX_ENABLE_ADVANCED_OPTIONS =
         "enableAdvancedCheckBox";
   public static final String ID_CHECKBOX_INCLUDE_BIOS_UUID = "biosUuidCheckBox";
   public static final String ID_CHECKBOX_INCLUDE_MAC_ADDRESS = "macAddressesCheckBox";
   public static final String ID_CHECKBOX_INCLUDE_EXTRA_CONFIG = "extraConfigCheckBox";
   public static final String ID_TEXT_INPUT_ANNOTATION = "annotationInput";
   public static final String ID_CHECKBOX_INCLUDE_IMAGE_FILES =
         "includeImageFilesCheckBox";
   public static final String ID_CHECKBOX_CDROM_START_CONNECTED =
         "_CdromPage_HBox1/startConnected";
   public static final String ID_CHECKBOX_FLOPPY_START_CONNECTED =
         "_FloppyPage_HBox1/startConnected";
   public static final String ID_TEXT_ADVANCED_WARNING = "advancedWarning/_message";
   public static final String ID_TEXT_EULA_HEADER = "eulaPage/pageHeaderTitleElement";
   public static final String ID_TEXT_EULA_DESCRIPTION =
         "eulaPage/pageHeaderDescriptionElement";
   public static final String ID_BUTTON_ACCEPT_EULA = "eulaPage/acceptButton";
   public static final String ID_DROPDOWNLIST_POLICY = "networkPage/policyDropDown";
   public static final String ID_TEXT_OVF_FILE = "tiwoDialog/propViewValueText";
   public static final String ID_DEPLOY_OVF_SELECT_SOURCE_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "step_sourcePage";
   public static final String ID_DEPLOY_OVF_REVIEW_DETAILS_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "step_detailsPage";
   public static final String ID_DEPLOY_OVF_ACCEPT_EULA_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "step_eulaPage";
   public static final String ID_DEPLOY_OVF_NAME_LOCATION_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "step_nameAndLocationPage";
   public static final String ID_DEPLOY_OVF_CONFIGURATION_STEP =
         ID_VMPROVISIONING_WIZARDPAGE_NAVIGATOR_STEP_INDEX + "step_configPage";
   public static final String ID_LABEL_SELECTED_FILE = "tiwoDialog/localFileStatusLabel";
   public static final String ID_STORAGE_PAGE = "storagePage";


   // Global Actions
   public static final String ID_BUTTON_VM_CREATE_GLOBAL = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_VM).append("createVmActionGlobal")
         .append(ID_BUTTON).toString();
   public static final String ID_BUTTON_VAPP_CREATE_GLOBAL = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_VAPP).append(CREATE_ACTION_GLOBAL)
         .append(ID_BUTTON).toString();
   public static final String ID_BUTTON_ADD_HOST_GLOBAL = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_HOST).append("addActionGlobal")
         .append(ID_BUTTON).toString();
   public static final String ID_BUTTON_CREATE_ACTION_GLOBAL = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("createActionGlobal")
         .append(ID_BUTTON).toString();
   public static final String ID_BUTTON_ADD_ACTION_GLOBAL = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("addActionGlobal")
         .append(ID_BUTTON).toString();
   public static final String ID_BUTTON_DVS_CREATE_ACTION_GLOBAL = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
         .append("createDvsActionGlobal").append(ID_BUTTON).toString();
   public static final String ID_BUTTON_DVS_CREATE_ACTION = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH).append("createDvsAction")
         .append(ID_BUTTON).toString();
   public static final String ID_BUTTON_DEPLOY_OVF_GLOBAL = new StringBuffer(
         ID_ACTION_DEPLOY_OVF).append(".global").append(ID_BUTTON).toString();

   //Deploy ovf pages
   public static final String ID_WIZARD_OVF_SOURCE_PAGE_TITLE =
         "tiwoDialog/sourcePage/pageHeaderTitleElement";
   public static final String ID_WIZARD_OVF_DETAILS_PAGE_TITLE =
         "tiwoDialog/detailsPage/pageHeaderTitleElement";
   public static final String ID_WIZARD_OVF_NAME_FOLDER_PAGE_TITLE =
         "tiwoDialog/nameAndLocationPage/pageHeaderTitleElement";
   public static final String ID_WIZARD_OVF_SUMMARY_PAGE_TITLE =
         "tiwoDialog/defaultSummaryPage/pageHeaderTitleElement";
   public static final String ID_WIZARD_PAGE_HEADER_TITLE = "pageHeaderTitleElement";
   public static final String ID_WIZARD_CREATION_PAGE_TITLE =
         "tiwoDialog/creationTypePage/pageHeaderTitleElement";
   public static final String ID_WIZARD_DESTINATION_PAGE_TITLE =
         "tiwoDialog/destinationPage/pageHeaderTitleElement";
   public static final String ID_WIZARD_NAMEFOLDER_PAGE_TITLE =
         "tiwoDialog/nameFolderPage/pageHeaderTitleElement";
   public static final String ID_WIZARD_STORAGE_PAGE_TITLE =
         "tiwoDialog/storagePage/pageHeaderTitleElement";
   public static final String ID_WIZARD_COMPLETE_PAGE_TITLE =
         "tiwoDialog/defaultSummaryPage/pageHeaderTitleElement";
   public static final String ID_WIZARD_NETWORK_PAGE_TITLE =
         "tiwoDialog/networkPage/pageHeaderTitleElement";
   public static final String ID_WIZARD_EULA_PAGE_TITLE =
         "tiwoDialog/eulaPage/pageHeaderTitleElement";
   public static final String ID_WIZARD_PROPERTIES_PAGE_TITLE =
         "tiwoDialog/propertiesPage/pageHeaderTitleElement";
   public static final String ID_WIZARD_CONFIGURATION_PAGE_TITLE =
         "tiwoDialog/configPage/pageHeaderTitleElement";
   public static final String ID_WIZARD_PROPERTIES_TEXT_TWO =
         "propertyStackBlock0EditPropertyControl2Row/textDisplay";
   public static final String ID_WIZARD_PROPERTIES_TEXT_THREE =
         "propertyStackBlock0EditPropertyControl3Row/textDisplay";

   //"next"; ID_WIZARD_NEXT_BUTTON
   public static final String ID_OVF_DEPLOY_URL_RADIO = "urlRadioButton";
   public static final String ID_OVF_DEPLOY_URL_TEXT = "textDisplay";
   public static final String ID_OVF_DEPLOY_REVIEWDETAILS_PAGE = "step_detailsPage";
   public static final String ID_OVF_DEPLOY_SOURCE_PAGE = "step_sourcePage";
   public static final String ID_OVF_DEPLOY_EULA_PAGE = "step_eulaPage";
   public static final String ID_OVF_DEPLOY_ACCEPT_BUTTON = "acceptButton";
   public static final String ID_OVF_DEPLOY_NAMELOCATION_PAGE =
         "step_nameAndLocationPage";
   public static final String ID_OVF_DEPLOY_NAME_FIELD = "textDisplay"; //- name text
   public static final String ID_OVF_DEPLOY_CONFIG_PAGE = "step_configPage";
   public static final String ID_OVF_DEPLOY_STORAGE_PAGE = "step_storagePage";
   public static final String ID_OVF_DEPLOY_STORAGE_TYPE_LABEL = "labelDisplay";
   public static final String ID_OVF_DEPLOY_STORAGE_LIST = "storageList";
   public static final String ID_OVF_DEPLOY_NETWORK_PAGE = "step_networkPage";
   public static final String ID_OVF_DEPLOY_NETWORK_IP_TYPE = "labelDisplay"; // - for ip type dhcp or static
   public static final String ID_OVF_DEPLOY_PROPERTIES_PAGE = "step_propertiesPage";
   public static final String ID_OVF_DEPLOY_SUMMARY_PAGE = "step_defaultSummaryPage";

   //clone OVF pages
   public static final String ID_CLONE_OVF_DESTINATION_PAGE = "step_destinationPage";

   //clone VApp
   public static final String ID_CLONE_VAPP_SELECT_CREATION_TYPE =
         "step_creationTypePage";
   public static final String ID_CLONE_VAPP_SELECTION_TYPE_PAGE = "dataGroup";


   // Permission constants
   public static final String ID_REMOVE_PERMISSION_BUTTON =
         "vsphere.core.permission.remove/button";
   public static final String ID_CHANGE_PERMISSION_BUTTON =
         "vsphere.core.permission.change/button";
   public static final String ID_PERMISSION_GRID = "permissionGrid";
   public static final String ID_ADD_USER_BUTTON = "addUser";
   public static final String ID_USERNAME_TEXTINPUT = "usersSelectionTextInput";
   public static final String ID_USER_OK_BUTTON = "okButton";

   // RDM related settings
   public static final String SELECT_TARGET_LUN_GRID = "grid";
   public static final String EDIT_VM_OK_BUTTON = "okButton";
   public static final String SELECT_TARGET_LUN_OK_BUTTON = "btnOK";
   public static final String ID_RDM_DISK = "RDM Disk";

   // compatibility
   public static final String ID_EDIT_VM_RDM_COMPATIBILITY_COMBO_BOX = "compatibility";
   // disk file
   public static final String ID_EDIT_VM_RDM_FILE_TEXT = "file";
   // disk mode
   public static final String ID_EDIT_VM_RDM_DISKMODE_APPEND = "radioButtonAppend";
   public static final String ID_EDIT_VM_RDM_DISKMODE_DEPENDENT = "radioButtonDependent";
   // independent nonpersistent
   public static final String ID_EDIT_VM_RDM__DISKMODE_IN = "radioButtonIN";
   // independent persistent
   public static final String ID_EDIT_VM_RDM_DISKMODE_IP = "radioButtonIP";
   public static final String ID_EDIT_VM_RDM_DISKMODE_CONTROL =
         "radioButtonNonpersistent";
   public static final String ID_EDIT_VM_RDM_DISKMODE_UNDOABLE = "radioButtonUndoable";
   // disk size
   public static final String ID_EDIT_VM_RDM_DISK_SIZE_COMBO_BOX =
         "diskSize/textDisplay";
   public static final String ID_EDIT_VM_RDM_DISKSIZE_LABEL = "maxSize";
   // disk type
   public static final String ID_EDIT_VM_RDM_DISKTYPE_LABEL = "diskTypeLabel";
   // disk limit
   public static final String ID_EDIT_VM_RDM_DISK_IOLIMIT_COMBO_BOX = "ioLimit";
   // disk location
   public static final String ID_EDIT_VM_RDM_DISK_LOCATION_COMBO_BOX = "location";
   // disk lun
   public static final String ID_EDIT_VM_RDM_DISK_LUN_LABEL = "lun";
   // disk provisioning
   public static final String ID_EDIT_VM_RDM_DISK_FLAT_PREINI_CHECKBOX = "eageryScurb";
   public static final String ID_EDIT_VM_RDM_DISK_ALLOC_CHECKBOX = "thin";
   // disk shares
   public static final String ID_EDIT_VM_RDM_DISK_TEXT_INPUT = "customShares";
   public static final String ID_EDIT_VM_RDM_DISK_SHARES_COMBOBOX = "sharesLevel";
   // disk write through
   public static final String ID_EDIT_VM_RDM_DISK_WRITE_THROUGH = "writeThrough";

   public static final String ID_EDIT_VM_RDM_REMOVE_LINK = "remove";
   public static final String ID_EDIT_VM_RDM_PURGE_CHECK_BOX = "purge";
   public static final String ID_ROLES_COMBO_BOX = "roles";

   // Host Memory CPU
   public static final String ID_HOST_PROCESSOR_MANUFACTURER =
         "HostSystem:processorSystemConfig.manufacturer";
   public static final String ID_HOST_PROCESSOR_MODEL =
         "HostSystem:processorSystemConfig.model";
   public static final String ID_HOST_PROCESSOR_BOOT_DEVICE =
         "HostSystem:processorSystemConfig.bootDevice";
   public static final String ID_HOST_PROCESSOR_BIOS_VERSION =
         "HostSystem:processorSystemConfig.biosVersion";
   public static final String ID_HOST_PROCESSOR_RELEASE_DATE =
         "HostSystem:processorSystemConfig.releaseDate";
   public static final String ID_HOST_PROCESSOR_ASSET_TAG =
         "HostSystem:processorSystemConfig.assetTag";
   public static final String ID_HOST_PROCESSOR_SERVICE_TAG =
         "HostSystem:processorSystemConfig.serviceTag";

   // Events related
   public static final String ID_EVENTS_ADVANCED_DATAGRID = "eventsDataGrid";
   public static final String ID_EVENTS_TYPE_DESCRIPTION_LABEL = "eventTypeDescLabel";
   public static final String ID_EVENTS_TYPE_DESCRIPTION_TEXTINPUT = "eventTypeDesc";
   public static final String ID_EVENTS_TYPE_CAUSE_DESCRIPTION_LABEL = "causeDesc";
   public static final String ID_EVENTS_TYPE_ACTION_DESCRIPTION_LABEL = "actionDesc";
   public static final String ID_ALARM_RECONFIGURE_EVENT =
         "An alarm has been reconfigured";
   public static final String ID_EMAIL_SENT_FAILED_EVENT =
         "An error occurred while sending email notification of a triggered alarm";
   public static final String ID_VM_DEPLOY_FROM_TEMPLATE_EVENT =
         "A virtual machine has been created from the specified template";
   public static final String ID_ALARM_SCRIPT_FAILURE_EVENT =
         "The vCenter Server logs this event if an error occurs while running a script when an alarm triggers.";
   public static final String ID_SNMP_TRAP_FAILURE_EVENT =
         "The vCenter Server logs this event if an error occurs while sending an SNMP trap when an alarm triggers.";
   public static final String ID_WWNS_CHANGED_EVENT = "WWNs are changed";
   public static final String ID_WWNS_CHANGED =
         "The WWN (World Wide Name) assigned to the virtual machine was changed";

   public static final String ID_SNMP_TRAP_FAILURE_CAUSE =
         "An SNMP trap could not be sent for a triggered alarm";
   public static final String ID_ALARM_RECONFIGURED_CAUSE =
         "A user has reconfigured an alarm";
   public static final String ID_EMAIL_SENT_FAILED_CAUSE =
         "Failed to send email for a triggered alarm";
   public static final String ID_ALARM_SCRIPT_FAILURE_CAUSE =
         "There was an error running the script";
   public static final String ID_ALARM_TEMPLATE_DEPLOYMENT_CAUSE_ONE =
         "A user action caused a virtual machine to be created from the template";
   public static final String ID_ALARM_TEMPLATE_DEPLOYMENT_TWO =
         "A scheduled task caused a virtual machine to be created from the template";
   public static final String ID_ALARM_WWNS_CHANGED_CAUSE =
         "The virtual machine was assigned a new WWN, possibly due to a conflict caused by another virtual machine being assigned the same WWN";

   public static final String ID_SNMP_TRAP_FAILURE_ACTION =
         "Action: Check the vCenter Server SNMP settings. Make sure that the vCenter Server network can handle SNMP packets.";
   public static final String ID_EMAIL_SENT_FAILED_ACTION =
         "Action: Check the vCenter Server SMTP settings for sending email notifications";
   public static final String ID_ALARM_SCRIPT_FAILURE_ACTION =
         "Action: Fix the script or failure condition";

   // SIOC
   public static final String MANUAL_RADIO_BUTTON = "manualThresholdRadioButton";
   public static final String STORAGE_IO_CONTROL_IMAGE =
         "arrowImagestorageIORMStackBlock";
   public static final String STORAGE_IO_CONTROL_STATUS = "storageIOControlStatus";
   public static final String STORAGE_CONGESTION_THRESHOLD_VALUE =
         "storageIOControlThreshold";
   public static final String ENABLE_STORAGE_IO_CONTROL_CHECK_BOX = "enableIORMCheckBox";
   public static final String THRESHOLD_VALUE = "manualThreshold";
   public static final String RESET_TO_DEFAULTS_BUTTON = "btnReset";
   public static final String THRESHOLD_VALUE_EXTENSION = "className=TextInput";
   public static final String BAR_BUTTON = "className=Button";

   // Tags Related
   public static final String ID_TAGS_BUTTON =
         "vsphere.core.tagging.application.dashboard.list.button";
   public static final String ID_CATEGORIES_BUTTON =
         "vsphere.core.tagging.application.dashboard.categories.button";
   public static final String ID_SINGLE_CARDINALITY_RADIOBUTTON =
         "_NewCategoryForm_RadioButton1";
   public static final String ID_MUTLIPLE_CARDINALITY_RADIO_BUTTON =
         "_NewCategoryForm_RadioButton2";
   public static final String ID_CREATE_TAGS_BUTTON =
         "vsphere.core.tagging.application.dashboard.list/dataGridToolbar/button";
   public static final String ID_CREATE_CATEGORY_BUTTON =
         "vsphere.core.tagging.application.dashboard.categories/dataGridToolbar/vsphere.core.tagging.createCategoryManageAction/button";
   public static final String ID_DELETE_CATEGORY_BUTTON =
         "vsphere.core.tagging.application.dashboard.categories/dataGridToolbar/vsphere.core.tagging.deleteCategoryAction/button";
   public static final String ID_EDIT_CATEGORY_BUTTON =
         "vsphere.core.tagging.application.dashboard.categories/dataGridToolbar/vsphere.core.tagging.editCategoryAction/button";
   public static final String ID_ASSIGN_TAG_BUTTON =
         "vsphere.core.tagging.assignTagAction/button";
   public static final String ID_AVAILABLE_CATEGORIES_GRID = "categoriesList";
   public static final String ID_AVAILABLE_TAGS_GRID = "tagsList";
   public static final String ID_TAGGING_TEXTINPUT = "textDisplay";
   public static final String ID_CATEGORY_DROPDOWN = "categoryDropDown";
   public static final String ID_ERROR_POPUP = "contents";
   public static final String TAGGING_BAR = "vsphere.core.tagging.";
   public static final String BUTTON_SUFFIX = "button";
   public static final String ID_NEW_TAG_BUTTON = TAGGING_BAR + "createTagAction" + "/"
         + BUTTON_SUFFIX;
   public static final String ID_DETACH_TAG_BUTTON = TAGGING_BAR + "detachTagAction"
         + "/" + BUTTON_SUFFIX;
   public static final String ID_TAG_NAME_TEXTFIELD = "nameInput";
   public static final String ID_TAG_DESCRIPTION_TEXTFIELD = "descInput";;
   public static final String ID_CATEGORY_NAME_TEXTFIELD = "categoryNameInput";
   public static final String ID_BUTTON_EDIT_TAGS =
         "vsphere.core.tagging.editTagAction/button";
   public static final String ID_BUTTON_DELETE_TAGS =
         "vsphere.core.tagging.deleteTagAction/button";
   public static final String ID_BUTTON_CREATE_TAGS_GSPAGE =
         "vsphere.core.tagging.createTagManageAction";
   public static final String ID_BUTTON_CREATE_CATEGORY_GSPAGE =
         "vsphere.core.tagging.createCategoryManageAction";
   public static final String ID_TEXTFIELD_TAG_DESCRIPTION = "categoryDescInput";
   public static final String ID_LIST_CATEGORY_ASSOCIABLEOBJECTS = "associableList";
   public static final String ID_SCROLLINGLIST_TAGS_CATEGORY_BUTTONCONTAINER =
         "scrollingList";
   public static final String ID_CHECKBOX_ASSOCIABLEOBJECTS_ALLOBJECTS =
         ID_LIST_CATEGORY_ASSOCIABLEOBJECTS + "/" + AUTOMATIONNAME + "All objects";
   public static final String ID_CHECKBOX_ASSOCIABLEOBJECTS_OBJECTS =
         ID_LIST_CATEGORY_ASSOCIABLEOBJECTS + "/" + AUTOMATIONNAME;
   public static final String ID_BUTTON_CONVERT_CA =
         "vsphere.core.tagging.createCAsAction";
   public static final String ID_LIST_CALIST = "caList";
   public static final String ID_LIST_GRID_CATEGORY =
         "vsphere.core.tagging.application.dashboard.list/className=DropDownList";
   public static final String ID_PORTLET_DCCLUSTER_TAG =
         "vsphere.core.dscluster.summary.tagsView.chrome";
   public static final String ID_PORTLET_HOST_TAG =
         "vsphere.core.host.summary.tagsView.chrome";
   public static final String ID_PORTLET_DS_TAG =
         "vsphere.core.datastore.summary.tagsView.chrome";
   public static final String ID_PORTLET_VM_TAG =
         "vsphere.core.vm.summary.tagsView.chrome";
   public static final String ID_PORTLET_DVS_TAG =
         "vsphere.core.dvs.summary.tagsView.chrome";
   public static final String ID_PORTLET_CLUSTER_TAG =
         "vsphere.core.cluster.summary.tagsView.chrome";
   public static final String ID_PORTLET_DVPG_TAG =
         "vsphere.core.dvPortgroup.summary.tagsView.chrome";
   public static final String ID_PORTLET_RP_TAG =
         "vsphere.core.resourcePool.summary.tagsView.chrome";
   public static final String ID_PORTLET_FOLDER_TAG =
         "vsphere.core.folder.summary.tagsView.chrome";
   public static final String ID_PORTLET_DC_TAG =
         "vsphere.core.datacenter.summary.tagsView.chrome";
   public static final String ID_PORTLET_NW_TAG =
         "vsphere.core.network.summary.tagsView.chrome";
   public static final String ID_PORTLET_VAPP_TAG =
         "vsphere.core.vApp.summary.tagsView.chrome";
   public static final String ID_TEXTINPUT_CATEGORY_INPUTBOX =
         "vsphere.core.tagging.application.dashboard.categories/filterControl/textInput";
   public static final String ID_BUTTON_CATEGORY_SEARCH =
         "vsphere.core.tagging.application.dashboard.categories/filterControl/searchButton";
   public static final String ID_BUTTON_CLUSTER_EVENTS_INPUTBOX =
         "vsphere.opsmgmt.event.cluster.eventView/filterControl/textInput";
   public static final String ID_BUTTON_CLUSTER_EVENTS_SEARCH =
         "vsphere.opsmgmt.event.cluster.eventView/filterControl/searchButton";
   public static final String ID_TEXTINPUT_TAG_INPUTBOX =
         "vsphere.core.tagging.application.dashboard.list/filterControl/textInput";
   public static final String ID_BUTTON_TAG_SEARCH =
         "vsphere.core.tagging.application.dashboard.list/filterControl/searchButton";
   public static final String ID_LIST_TAGSGRID_CATEGORYLIST =
         "dataGridToolbar/className=DropDownList";
   public static final String ID_TEXTINPUT_ASSIGNTAG_INPUTBOX =
         "tagsList/dataGridToolbar/filterControl/textInput";
   public static final String ID_BUTTON__ASSIGNTAG_SEARCH =
         "tagsList/dataGridToolbar/filterControl/searchButton";
   public static final String ID_ACTION_EDIT_CATEGORY =
         "afContextMenu.vsphere.core.tagging.editCategoryAction";
   public static final String ID_ACTION_DELETE_CATEGORY =
         "afContextMenu.vsphere.core.tagging.deleteCategoryAction";
   public static final String ID_ACTION_EDIT_TAG =
         "afContextMenu.vsphere.core.tagging.editTagAction";
   public static final String ID_ACTION_DELETE_TAG =
         "afContextMenu.vsphere.core.tagging.deleteTagAction";
   public static final String ID_ACTION_DETACH_TAG =
         "afContextMenu.vsphere.core.tagging.detachTagAction";

   // Storage reports
   public static final String ID_LABEL_INPUT = "inputLabel";
   public static final String ID_REPORTS_UPDATED_BUTTON = "smsUiSyncButton";
   public static final String ID_DATACENTER_ENTITY = "datacenter";
   public static final String ID_HOST_ENTITY = "host";
   public static final String ID_VM_ENTITY = "Vm";
   public static final String ID_CLUSTER_ENTITY = "cluster";
   public static final String ID_DATASTORE_ENTITY = "Datastore";
   public static final String ID_GRID_HOSTS_VM = "hostVmGrid";
   public static final String ID_GRID_TORAGE_REPORTS_DATACENTER_CLUSTER =
         new StringBuffer(ID_DATACENTER_ENTITY).append(ID_CLUSTER_ENTITY).append("Grid")
               .toString();
   public static final String ID_MENU_VM_FILE_REPORT = new StringBuffer("vm")
         .append(ID_VM_ENTITY).append("FileReportMenu").toString();
   public static final String ID_MENU_HOST_VM_REPORT = new StringBuffer(ID_HOST_ENTITY)
         .append(ID_VM_ENTITY).append("ReportMenu").toString();
   public static final String ID_MENU_DC_VM_REPORT = new StringBuffer(
         ID_DATACENTER_ENTITY).append(ID_VM_ENTITY).append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DC_VM = new StringBuffer(
         ID_DATACENTER_ENTITY).append(ID_VM_ENTITY).append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_HOST_DATASTORE = new StringBuffer(
         ID_HOST_ENTITY).append(ID_DATASTORE_ENTITY).append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_HOSTS_VM = new StringBuffer(
         ID_HOST_ENTITY).append(ID_VM_ENTITY).append("Grid").toString();
   public static final String ID_BUTTON_STORAGE_REPORTS_DC_VM_SYNC = new StringBuffer(
         ID_DATACENTER_ENTITY).append(ID_VM_ENTITY).append("SyncButton").toString();
   public static final String ID_BUTTON_STORAGE_REPORTS_HOST_VM_UPDATE =
         new StringBuffer(ID_HOST_ENTITY).append(ID_VM_ENTITY).append("SyncButton")
               .toString();
   public static final String ID_BUTTON_STORAGE_REPORTS_VIEW = new StringBuffer(
         EXTENSION_PREFIX).append(ID_HOST_ENTITY).append(".monitor.smsView.button")
         .toString();
   public static final String ID_BUTTON_DATACENTER_REPORTS_VIEW = new StringBuffer(
         EXTENSION_PREFIX).append(ID_DATACENTER_ENTITY)
         .append(".monitor.smsView.button").toString();
   public static final String ID_RP_ENTITY = "resourcePool";
   public static final String ID_GRID_STORAGE_REPORTS_DATACENTER_CLUSTER =
         new StringBuffer(ID_DATACENTER_ENTITY).append("Cluster").append("Grid")
               .toString();
   public static final String ID_GRID_STORAGE_REPORTS_DC_DATASTORE = new StringBuffer(
         ID_DATACENTER_ENTITY).append(ID_DATASTORE_ENTITY).append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DATASTORE_VM = new StringBuffer(
         "datastore").append(ID_VM_ENTITY).append("Grid").toString();
   public static final String ID_MENU_DC_DATASTORE_REPORT = new StringBuffer(
         ID_DATACENTER_ENTITY).append(ID_DATASTORE_ENTITY).append("ReportMenu")
         .toString();
   public static final String ID_BUTTON_STORAGE_REPORTS_VM_DATASTORE = new StringBuffer(

   EXTENSION_PREFIX).append("vm").append(".monitor.smsView.button").toString();

   public static final String ID_BUTTON_STORAGE_REPORTS_DATASTORE = new StringBuffer(
         EXTENSION_PREFIX).append("datastore").append(".monitor.smsView.button")
         .toString();
   public static final String ID_GRID_STORAGE_REPORTS_CLUSTER_VM = new StringBuffer(
         ID_CLUSTER_ENTITY).append(ID_VM_ENTITY).append("Grid").toString();

   public static final String ID_GRID_STORAGE_REPORTS_VM_DATASTORE = new StringBuffer(
         "vm").append(ID_DATASTORE_ENTITY).append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DATASTORE_HOST = new StringBuffer(
         "datastore").append("Host").append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_CLUSTER_DATASTORE =
         new StringBuffer(ID_CLUSTER_ENTITY).append(ID_DATASTORE_ENTITY).append("Grid")
               .toString();
   public static final String ID_GRID_STORAGE_REPORTS_CLUSTER_HOST = new StringBuffer(
         ID_CLUSTER_ENTITY).append("Host").append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_CLUSTER_RP = new StringBuffer(
         ID_CLUSTER_ENTITY).append("ResourcePool").append("Grid").toString();
   public static final String ID_MENU_DATASTORE_VM_REPORT =
         new StringBuffer("datastore").append(ID_VM_ENTITY).append("ReportMenu")
               .toString();
   public static final String ID_MENU_CLUSTER_VM_REPORT = new StringBuffer(
         ID_CLUSTER_ENTITY).append(ID_VM_ENTITY).append("ReportMenu").toString();
   public static final String ID_MENU_CLUSTER_HOST_REPORT = new StringBuffer(
         ID_CLUSTER_ENTITY).append("Host").append("ReportMenu").toString();
   public static final String ID_MENU_CLUSTER_RP_REPORT = new StringBuffer(
         ID_CLUSTER_ENTITY).append("ResourcePool").append("ReportMenu").toString();
   public static final String ID_MENU_CLUSTER_DATASTORE_REPORT = new StringBuffer(
         ID_CLUSTER_ENTITY).append("Datastore").append("ReportMenu").toString();
   public static final String ID_MENU_CLUSTER_SCSI_VOLUME_REPORT = new StringBuffer(
         ID_CLUSTER_ENTITY).append("ScsiVolume").append("ReportMenu").toString();
   public static final String ID_MENU_VM_DATASTORE_REPORT = new StringBuffer("vm")
         .append(ID_DATASTORE_ENTITY).append("ReportMenu").toString();
   public static final String ID_MENU_VM_SCSI_VOLUME_REPORT = new StringBuffer("vm")
         .append("ScsiVolume").append("ReportMenu").toString();
   public static final String ID_MENU_RP_VM_REPORT = new StringBuffer(ID_RP_ENTITY)
         .append(ID_VM_ENTITY).append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_RP_DATASTORE = new StringBuffer(
         ID_RP_ENTITY).append(ID_DATASTORE_ENTITY).append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_HOST_SCSI_VOLUMES =
         new StringBuffer(ID_HOST_ENTITY).append("ScsiVolume").append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_CLUSTER_SCSI_VOLUMES =
         new StringBuffer(ID_CLUSTER_ENTITY).append("ScsiVolume").append("Grid")
               .toString();
   public static final String ID_GRID_STORAGE_REPORTS_VM_SCSI_VOLUMES =
         new StringBuffer("vm").append("ScsiVolume").append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DATASTORE_SCSI_VOLUMES =
         new StringBuffer("datastore").append("ScsiVolume").append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_VM_FILES = new StringBuffer("vm")
         .append("VmFile").append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_RP_VM = new StringBuffer(
         ID_RP_ENTITY).append(ID_VM_ENTITY).append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_HOST_VM = new StringBuffer(
         ID_HOST_ENTITY).append(ID_VM_ENTITY).append("Grid").toString();
   public static final String ID_BUTTON_STORAGE_REPORTS_RP_VM_UPDATE = new StringBuffer(
         ID_RP_ENTITY).append(ID_VM_ENTITY).append("SyncButton").toString();
   public static final String ID_BUTTON_STORAGE_REPORTS_CLUSTER_VM_UPDATE =
         new StringBuffer(ID_CLUSTER_ENTITY).append(ID_VM_ENTITY).append("SyncButton")
               .toString();
   public static final String ID_BUTTON_STORAGE_REPORTS_VM_DATASTORE_UPDATE =
         new StringBuffer("vm").append(ID_DATASTORE_ENTITY).append("SyncButton")
               .toString();
   public static final String ID_BUTTON_STORAGE_REPORTS_DATASTORE_VM_UPDATE =
         new StringBuffer("datastore").append("Vm").append("SyncButton").toString();
   public static final String ID_BUTTON_CLUSTER_REPORTS_VIEW = new StringBuffer(
         EXTENSION_PREFIX).append("cluster").append(".monitor.smsView.button")
         .toString();
   public static final String ID_BUTTON_RESOURCEPOOL_REPORTS_VIEW = new StringBuffer(
         EXTENSION_PREFIX).append(ID_RP_ENTITY).append(".monitor.smsView.button")
         .toString();
   public static final String ID_MENU_DC_HOST_REPORT = new StringBuffer(
         ID_DATACENTER_ENTITY).append(ID_HOST_ENTITY).append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DC_HOST = new StringBuffer(
         ID_DATACENTER_ENTITY).append(ID_HOST_ENTITY).append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DC_CLUSTER = new StringBuffer(
         ID_DATACENTER_ENTITY).append(ID_CLUSTER_ENTITY).append("Grid").toString();
   public static final String ID_MENU_DC_CLUSTER_REPORT = new StringBuffer(
         ID_DATACENTER_ENTITY).append(ID_CLUSTER_ENTITY).append("ReportMenu").toString();
   public static final String ID_BUTTON_STORAGE_REPORTS_DC_DATASTORE_UPDATE =
         new StringBuffer(ID_DATACENTER_ENTITY).append(ID_DATASTORE_ENTITY)
               .append("SyncButton").toString();
   public static final String ID_MENU_HOST_DATASTORE_REPORT = new StringBuffer(
         ID_HOST_ENTITY).append(ID_DATASTORE_ENTITY).append("ReportMenu").toString();
   public static final String ID_MENU_HOST_SCSI_REPORT =
         new StringBuffer(ID_HOST_ENTITY).append("ScsiVolume").append("ReportMenu")
               .toString();
   public static final String ID_MENU_DC_SCSI_PATH_REPORT = new StringBuffer(
         ID_DATACENTER_ENTITY).append("ScsiPath").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_HOST_SCSI_PATH = new StringBuffer(
         ID_HOST_ENTITY).append("ScsiPath").append("Grid").toString();
   public static final String ID_MENU_HOST_SCSI_PATH_REPORT = new StringBuffer(
         ID_HOST_ENTITY).append("ScsiPath").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_HOST_SCSI_ADAPTER =
         new StringBuffer(ID_HOST_ENTITY).append("ScsiAdapter").append("Grid")
               .toString();
   public static final String ID_MENU_HOST_SCSI_ADAPTER_REPORT = new StringBuffer(
         ID_HOST_ENTITY).append("ScsiAdapter").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_HOST_SCSI_TARGET =
         new StringBuffer(ID_HOST_ENTITY).append("ScsiTarget").append("Grid").toString();
   public static final String ID_MENU_HOST_SCSI_TARGET_REPORT = new StringBuffer(
         ID_HOST_ENTITY).append("ScsiTarget").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DC_SCSI_PATH = new StringBuffer(
         ID_DATACENTER_ENTITY).append("ScsiPath").append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DC_SCSI_ADAPTER =
         new StringBuffer(ID_DATACENTER_ENTITY).append("ScsiAdapter").append("Grid")
               .toString();
   public static final String ID_MENU_DC_SCSI_ADAPTER_REPORT = new StringBuffer(
         ID_DATACENTER_ENTITY).append("ScsiAdapter").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DC_SCSI_TARGET = new StringBuffer(
         ID_DATACENTER_ENTITY).append("ScsiTarget").append("Grid").toString();
   public static final String ID_MENU_DC_SCSI_TARGET_REPORT = new StringBuffer(
         ID_DATACENTER_ENTITY).append("ScsiTarget").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DC_NAS_MOUNT = new StringBuffer(
         ID_DATACENTER_ENTITY).append("NasMount").append("Grid").toString();
   public static final String ID_MENU_DC_NAS_MOUNT_REPORT = new StringBuffer(
         ID_DATACENTER_ENTITY).append("NasMount").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_HOST_NAS_MOUNT = new StringBuffer(
         ID_HOST_ENTITY).append("NasMount").append("Grid").toString();
   public static final String ID_MENU_HOST_NAS_MOUNT_REPORT = new StringBuffer(
         ID_HOST_ENTITY).append("NasMount").append("ReportMenu").toString();
   public static final String ID_MENU_RP_DATASTORES_REPORT = new StringBuffer(
         ID_RP_ENTITY).append(ID_DATASTORE_ENTITY).append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_RP_SCSI_VOLUME = new StringBuffer(
         ID_RP_ENTITY).append("ScsiVolume").append("Grid").toString();
   public static final String ID_MENU_RP_SCSI_VOLUME_REPORT = new StringBuffer(
         ID_RP_ENTITY).append("ScsiVolume").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_RP_SCSI_PATH = new StringBuffer(
         ID_RP_ENTITY).append("ScsiPath").append("Grid").toString();
   public static final String ID_MENU_RP_SCSI_PATH_REPORT = new StringBuffer(
         ID_RP_ENTITY).append("ScsiPath").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_RP_SCSI_ADAPTER =
         new StringBuffer(ID_RP_ENTITY).append("ScsiAdapter").append("Grid").toString();
   public static final String ID_MENU_RP_SCSI_ADAPTER_REPORT = new StringBuffer(
         ID_RP_ENTITY).append("ScsiAdapter").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_RP_SCSI_TARGET = new StringBuffer(
         ID_RP_ENTITY).append("ScsiTarget").append("Grid").toString();
   public static final String ID_MENU_RP_SCSI_TARGET_REPORT = new StringBuffer(
         ID_RP_ENTITY).append("ScsiTarget").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_RP_NAS_MOUNT = new StringBuffer(
         ID_RP_ENTITY).append("NasMount").append("Grid").toString();
   public static final String ID_MENU_RP_NAS_MOUNT_REPORT = new StringBuffer(
         ID_RP_ENTITY).append("NasMount").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_CLUSTER_SCSI_VOLUME =
         new StringBuffer(ID_CLUSTER_ENTITY).append("ScsiVolume").append("Grid")
               .toString();
   public static final String ID_GRID_STORAGE_REPORTS_CLUSTER_SCSI_PATH =
         new StringBuffer(ID_CLUSTER_ENTITY).append("ScsiPath").append("Grid")
               .toString();
   public static final String ID_MENU_CLUSTER_SCSI_PATH_REPORT = new StringBuffer(
         ID_CLUSTER_ENTITY).append("ScsiPath").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_VM_SCSI_VOLUME = new StringBuffer(
         "vm").append("ScsiVolume").append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_VM_SCSI_PATH = new StringBuffer(
         "vm").append("ScsiPath").append("Grid").toString();
   public static final String ID_MENU_VM_SCSI_PATH_REPORT = new StringBuffer("vm")
         .append("ScsiPath").append("ReportMenu").toString();
   public static final String ID_MENU_VM_SCSI_ADAPTERS_REPORT = new StringBuffer("vm")
         .append("ScsiAdapter").append("ReportMenu").toString();
   public static final String ID_MENU_VM_SCSI_TARGET_REPORT = new StringBuffer("vm")
         .append("ScsiTarget").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_VM_SCSI_ADAPTER =
         new StringBuffer("vm").append("ScsiAdapter").append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_VM_SCSI_TARGET = new StringBuffer(
         "vm").append("ScsiTarget").append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_VM_NAS_MOUNT = new StringBuffer(
         "vm").append("NasMount").append("Grid").toString();
   public static final String ID_MENU_VM_NAS_MOUNT_REPORT = new StringBuffer("vm")
         .append("NasMount").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DATASTORE_VM_FILES =
         new StringBuffer("datastore").append("VmFile").append("Grid").toString();
   public static final String ID_MENU_DATASTORE_VM_FILES_REPORT = new StringBuffer(
         "datastore").append("VmFile").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DATASTORE_SCSI_VOLUME =
         new StringBuffer("datastore").append("ScsiVolume").append("Grid").toString();
   public static final String ID_MENU_DATASTORE_SCSI_VOLUMES_REPORT = new StringBuffer(
         "datastore").append("ScsiVolume").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DATASTORE_SCSI_ADAPTER =
         new StringBuffer("datastore").append("ScsiAdapter").append("Grid").toString();
   public static final String ID_MENU_DATASTORE_SCSI_ADAPTERS_REPORT = new StringBuffer(
         "datastore").append("ScsiAdapter").append("Grid").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DATASTORE_SCSI_PATH =
         new StringBuffer("datastore").append("ScsiPath").append("Grid").toString();
   public static final String ID_MENU_DATASTORE_SCSI_PATH_REPORT = new StringBuffer(
         "datastore").append("ScsiPath").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DATASTORE_SCSI_TARGET =
         new StringBuffer("datastore").append("ScsiTarget").append("Grid").toString();
   public static final String ID_MENU_DATASTORE_SCSI_TARGETS_REPORT = new StringBuffer(
         "datastore").append("ScsiTarget").append("ReportMenu").toString();
   public static final String ID_GRID_STORAGE_REPORTS_DATASTORE_NAS_MOUNT =
         new StringBuffer("datastore").append("NasMount").append("Grid").toString();
   public static final String ID_MENU_DATASTORE_NAS_MOUNT_REPORT = new StringBuffer(
         "datastore").append("NasMount").append("ReportMenu").toString();


   // Added constants for RDM
   public static final String ID_POWERED_OFF = "PoweredOff";
   public static final String ID_POWERED_ON = "PoweredOn";

   public static final String ID_HA_PROTECTED_IMAGE = "haProtectedBadge";
   public static final String ID_HOST_CONFIGURATION_FT_REQUIREMENTS =
         "ftReasonPrerequisitesText";

   // Added constants for networking

   public static final String ID_INCREMENT_BUTTON = "incrementButton";
   public static final String ID_DECREMENT_BUTTON = "decrementButton";

   // IDs of VLAN policy
   public static final String ID_POLICY_VLAN_ID_LABEL = "_DvPortVlanPolicyView_LabelEx2";

   // IDs for Edit Port Group Wizard
   public static final String ID_WIZARD_PAGE_NAVIGATOR = "wizardPageNavigator";

   // ToC List items` ids
   public static final String ID_GENERAL_EDIT_PORT_GROUP_MENU_ITEM =
         "step_portgeoupGeneralProperties";
   public static final String ID_ADVANCED_EDIT_PORT_GROUP_MENU_ITEM =
         "step_portgeoupAdvancedProperties";
   public static final String ID_SECURITY_EDIT_PORT_GROUP_MENU_ITEM = "step_security";
   public static final String ID_TRAFFIC_SHAPING_EDIT_PORT_GROUP_MENU_ITEM =
         "step_trafficShaping";
   public static final String ID_VLAN_EDIT_PORT_GROUP_MENU_ITEM = "step_vlan";
   public static final String ID_TEAMING_AND_FAILOVER_EDIT_PORT_GROUP_MENU_ITEM =
         "step_5";
   public static final String ID_MONITORING_EDIT_PORT_GROUP_MENU_ITEM =
         "step_monitoring";
   public static final String ID_MISCELLANEOUS_EDIT_PORT_GROUP_MENU_ITEM = "step_misc";
   public static final String ID_TEAMING_AND_FAILOVER_EDIT_PORT_MENU_ITEM = "step_4";

   // General page
   public static final String ID_GENERAL_PAGE_EDIT_PORT_GROUP = "generalPropertiesForm";
   public static final String ID_PORT_GROUP_NAME_EDIT_PORT_GROUP = "portgroupName";
   public static final String ID_NUMBER_OF_PORTS_EDIT_PORT_GROUP =
         ID_GENERAL_PAGE_EDIT_PORT_GROUP + "/"
               + "_DvPortgroupGeneralPropertiesPage_NumericStepper1";
   public static final String ID_NUMBER_OF_PORTS_EDIT_PORT_GROUP_INCREMENT =
         ID_NUMBER_OF_PORTS_EDIT_PORT_GROUP + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_NUMBER_OF_PORTS_EDIT_PORT_GROUP_DECREMENT =
         ID_NUMBER_OF_PORTS_EDIT_PORT_GROUP + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_PORT_BINDING_EDIT_PORT_GROUP =
         ID_GENERAL_PAGE_EDIT_PORT_GROUP + "/" + "portBinding";
   public static final String ID_PORT_ALLOCATION_EDIT_PORT_GROUP =
         ID_GENERAL_PAGE_EDIT_PORT_GROUP + "/" + "portAllocDrowDown";
   public static final String ID_NETWORK_RESOUCE_POOL_EDIT_PORT_GROUP =
         ID_GENERAL_PAGE_EDIT_PORT_GROUP + "/" + "poolSelector";
   public static final String ID_DYNAMIC_BINDING_WARNING_EDIT_PORT_GROUP =
         ID_GENERAL_PAGE_EDIT_PORT_GROUP + "/" + "portBindingWarningText";
   public static final String ID_EDIT_PORT_GROUP_DESCRIPTION =
         "_DvPortgroupGeneralPropertiesPage_TextAreaEx1";
   public static final String ID_EDIT_PORT_GROUP_NETWORK_RESOURCE_POOL =
         ID_GENERAL_PAGE_EDIT_PORT_GROUP + "/" + "poolSelector";

   // Security page
   public static final String ID_PROMISCUOUS_MODE_DROPDOWN_EDIT_PORT_GROUP =
         "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList1" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_MAC_ADDRESS_CHANGES_DROPDOWN_EDIT_PORT_GROUP =
         "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList2" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_FORGED_TRANSMITS_DROPDOWN_EDIT_PORT_GROUP =
         "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList3" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";

   // Advanced page
   public static final String ID_EDIT_DVPG_ADVANCED_STACK_BLOCK =
         "_DvPortgroupAdvancedPropertiesPage_StackBlock1";
   public static final String ID_ADVANCED_DISABLE_OVERRIDE =
         "_DvPortgroupPolicyOverrideAllowedRadioSelection_RadioButton2";
   public static final String ID_ADVANCED_ALLOW_OVERRIDE =
         "_DvPortgroupPolicyOverrideAllowedRadioSelection_RadioButton1";
   public static final String ID_CONFIGURE_RESET_AT_DISCONNECT_DROP_DOWN =
         "_DvPortgroupAdvancedPropertiesPage_BooleanDropDownList1";

   public static final String ID_ADVANCED_BLOCK_PORTS_ALLOW_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection1"
               + "/" + ID_ADVANCED_ALLOW_OVERRIDE;
   public static final String ID_ADVANCED_BLOCK_PORTS_DISABLE_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection1"
               + "/" + ID_ADVANCED_DISABLE_OVERRIDE;

   public static final String ID_ADVANCED_TRAFFIC_SHAPING_ALLOW_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection2"
               + "/" + ID_ADVANCED_ALLOW_OVERRIDE;
   public static final String ID_ADVANCED_TRAFFIC_SHAPING_DISABLE_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection2"
               + "/" + ID_ADVANCED_DISABLE_OVERRIDE;

   public static final String ID_ADVANCED_VENDOR_ALLOW_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection3"
               + "/" + ID_ADVANCED_ALLOW_OVERRIDE;
   public static final String ID_ADVANCED_VENDOR_DISABLE_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection3"
               + "/" + ID_ADVANCED_DISABLE_OVERRIDE;

   public static final String ID_ADVANCED_VLAN_ALLOW_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection4"
               + "/" + ID_ADVANCED_ALLOW_OVERRIDE;
   public static final String ID_ADVANCED_VLAN_DISABLE_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection4"
               + "/" + ID_ADVANCED_DISABLE_OVERRIDE;

   public static final String ID_ADVANCED_UPLINK_TEAMING_ALLOW_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection5"
               + "/" + ID_ADVANCED_ALLOW_OVERRIDE;
   public static final String ID_ADVANCED_UPLINK_TEAMING_DISABLE_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection5"
               + "/" + ID_ADVANCED_DISABLE_OVERRIDE;

   public static final String ID_ADVANCED_NETWORK_IO_ALLOW_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection6"
               + "/" + ID_ADVANCED_ALLOW_OVERRIDE;
   public static final String ID_ADVANCED_NETWORK_IO_DISABLE_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection6"
               + "/" + ID_ADVANCED_DISABLE_OVERRIDE;

   public static final String ID_ADVANCED_SECURITY_ALLOW_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection7"
               + "/" + ID_ADVANCED_ALLOW_OVERRIDE;
   public static final String ID_ADVANCED_SECURITY_DISABLE_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection7"
               + "/" + ID_ADVANCED_DISABLE_OVERRIDE;

   public static final String ID_ADVANCED_NETFLOW_ALLOW_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection8"
               + "/" + ID_ADVANCED_ALLOW_OVERRIDE;
   public static final String ID_ADVANCED_NETFLOW_DISABLE_RADIOBUTTON_EDIT_DVPG =
         "_DvPortgroupAdvancedPropertiesPage_DvPortgroupPolicyOverrideAllowedRadioSelection8"
               + "/" + ID_ADVANCED_DISABLE_OVERRIDE;


   // Traffic shapping page
   public static final String ID_TRAFFIC_SHAPPING_PAGE_EDIT_PORT_GROUP =
         "trafficShaping";

   public static final String ID_EGRESS_TRAFFIC_SHAPING_STACK_BLOCK_EDIT_PORT_GROUP =
         "_DvPortTrafficShapingPolicyListPage_StackBlock2";
   public static final String ID_EGRESS_TRAFFIC_SHAPING_SETTINGS_BLOCK_EDIT_PORT_GROUP =
         "_DvPortTrafficShapingPolicyListPage_SettingsBlock2";
   public static final String ID_EGRESS_STATUS_EDIT_PORT_GROUP =
         ID_EGRESS_TRAFFIC_SHAPING_SETTINGS_BLOCK_EDIT_PORT_GROUP + "/"
               + "_DvPortTrafficShapingPolicyPage_BooleanDropDownList1";
   public static final String ID_EGRESS_AVERAGE_EDIT_PORT_GROUP =
         ID_EGRESS_TRAFFIC_SHAPING_SETTINGS_BLOCK_EDIT_PORT_GROUP + "/"
               + "averageNumStepper";
   public static final String ID_EGRESS_PEAK_EDIT_PORT_GROUP =
         ID_EGRESS_TRAFFIC_SHAPING_SETTINGS_BLOCK_EDIT_PORT_GROUP + "/"
               + "peakNumStepper";
   public static final String ID_EGRESS_BURST_EDIT_PORT_GROUP =
         ID_EGRESS_TRAFFIC_SHAPING_SETTINGS_BLOCK_EDIT_PORT_GROUP + "/"
               + "burstNumStepper";

   public static final String ID_INGRESS_TRAFFIC_SHAPING_STACK_BLOCK_EDIT_PORT_GROUP =
         "_DvPortTrafficShapingPolicyListPage_StackBlock1";
   public static final String ID_INGRESS_TRAFFIC_SHAPING_SETTINGS_BLOCK_EDIT_PORT_GROUP =
         "_DvPortTrafficShapingPolicyListPage_SettingsBlock1";
   public static final String ID_INGRESS_TRAFFIC_SHAPING_STATUS_EDIT_PORT_GROUP =
         ID_INGRESS_TRAFFIC_SHAPING_SETTINGS_BLOCK_EDIT_PORT_GROUP + "/"
               + "_DvPortTrafficShapingPolicyPage_BooleanDropDownList1";
   public static final String ID_INGRESS_AVERAGE_EDIT_PORT_GROUP =
         ID_INGRESS_TRAFFIC_SHAPING_SETTINGS_BLOCK_EDIT_PORT_GROUP + "/"
               + "averageNumStepper";
   public static final String ID_INGRESS_PEAK_EDIT_PORT_GROUP =
         ID_INGRESS_TRAFFIC_SHAPING_SETTINGS_BLOCK_EDIT_PORT_GROUP + "/"
               + "peakNumStepper";
   public static final String ID_INGRESS_BURST_EDIT_PORT_GROUP =
         ID_INGRESS_TRAFFIC_SHAPING_SETTINGS_BLOCK_EDIT_PORT_GROUP + "/"
               + "burstNumStepper";

   public static final String ID_EGRESS_AVERAGE_EDIT_PORT_GROUP_INCREMENT_BUTTON =
         ID_EGRESS_AVERAGE_EDIT_PORT_GROUP + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_EGRESS_AVERAGE_EDIT_PORT_GROUP_DECREMENT_BUTTON =
         ID_EGRESS_AVERAGE_EDIT_PORT_GROUP + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_EGRESS_PEAK_EDIT_PORT_GROUP_INCREMENT_BUTTON =
         ID_EGRESS_PEAK_EDIT_PORT_GROUP + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_EGRESS_PEAK_EDIT_PORT_GROUP_DECREMENT_BUTTON =
         ID_EGRESS_PEAK_EDIT_PORT_GROUP + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_EGRESS_BURST_EDIT_PORT_GROUP_INCREMENT_BUTTON =
         ID_EGRESS_BURST_EDIT_PORT_GROUP + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_EGRESS_BURST_EDIT_PORT_GROUP_DECREMENT_BUTTON =
         ID_EGRESS_BURST_EDIT_PORT_GROUP + "/" + ID_DECREMENT_BUTTON;

   public static final String ID_INGRESS_AVERAGE_EDIT_PORT_GROUP_INCREMENT_BUTTON =
         ID_INGRESS_AVERAGE_EDIT_PORT_GROUP + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_INGRESS_AVERAGE_EDIT_PORT_GROUP_DECREMENT_BUTTON =
         ID_INGRESS_AVERAGE_EDIT_PORT_GROUP + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_INGRESS_PEAK_EDIT_PORT_GROUP_INCREMENT_BUTTON =
         ID_INGRESS_PEAK_EDIT_PORT_GROUP + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_INGRESS_PEAK_EDIT_PORT_GROUP_DECREMENT_BUTTON =
         ID_INGRESS_PEAK_EDIT_PORT_GROUP + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_INGRESS_BURST_EDIT_PORT_GROUP_INCREMENT_BUTTON =
         ID_INGRESS_BURST_EDIT_PORT_GROUP + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_INGRESS_BURST_EDIT_PORT_GROUP_DECREMENT_BUTTON =
         ID_INGRESS_BURST_EDIT_PORT_GROUP + "/" + ID_DECREMENT_BUTTON;

   // VLAN page
   public static final String ID_VLAN_ID_PAGE_EDIT_PORT_GROUP = "vlan";
   public static final String ID_VLAN_TYPE_DROPDOWN_EDIT_PORT_GROUP = "vlanTypeSelector";
   public static final String ID_DROPDOWN_PVLAN_ID_EDIT_PORT_GROUP =
         "privateVlanIdSelector";
   public static final String ID_VLAN_ID_EDIT_PORT_GROUP =
         ID_VLAN_ID_PAGE_EDIT_PORT_GROUP + "/" + "_DvPortVlanPolicyPage_NumericStepper1";
   public static final String ID_VLAN_ID_EDIT_PORT_GROUP_INCREMET_BUTTON =
         ID_VLAN_ID_EDIT_PORT_GROUP + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_VLAN_ID_EDIT_PORT_GROUP_DECREMET_BUTTON =
         ID_VLAN_ID_EDIT_PORT_GROUP + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_VLAN_TRUNK_RANGE_EDIT_PORT_GROUP = "vlanIdRanges";
   public static final String ID_VLAN_VLANID_EDIT_SETTINGS = "openButton";

   // Teaming and failover edit page
   public static final String ID_LOAD_BALANCING_EDIT_PORT_GROUP =
         "loadBalancingSelector";
   public static final String ID_NETWORK_FAILURE_DETECTION_EDIT_PORT_GROUP =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList1/_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_NOTIFY_SWITCHES_EDIT_PORT_GROUP =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList2/_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_FAILBACK_EDIT_PORT_GROUP =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList3/_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_LOAD_BALANCING_DETAILS = "loadBalancingDetails";
   public static final String MOVE_DOWN_BUTTON_TOOLTIP = "/toolTip="
         + MOVE_DOWN_UPLINK_BUTTON;
   public static final String MOVE_UP_BUTTON_TOOLTIP = "/toolTip="
         + MOVE_UP_UPLINK_BUTTON;

   // Monitoring page
   public static final String ID_NET_FLOW_EDIT_PORT_GROUP =
         "_DvPortMonitoringPolicyPage_OverridableBooleanDropDownList1" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_NET_FLOW_EDIT_UPLINK_PORT_GROUP = "monitoring" + "/"
         + "_OverridableBooleanDropDownList_BooleanDropDownList1";

   // Miscellaneous page
   public static final String ID_BLOCK_ALL_PORTS_EDIT_PORT_GROUP =
         "_DvPortMiscPoliciesPage_OverridableBooleanDropDownList1" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";

   public static final String ID_EDIT_VM_NETWORK_CONNECTION = "connection";

   // IDs in Create vDS wizard
   public static final String ID_CREATE_VDS_NAME_FIELD =
         "nameAndLocationPage/inputLabel";
   public static final String ID_CREATE_VDS_UPLINKS_FIELD = "uplinkPortsNumber";
   public static final String ID_CREATE_VDS_NETIOCONTROL_DROPDOWN = "resMgmtEnableState";
   public static final String ID_CREATE_VDS_DEFAULT_PORTGROUP_CHECKBOX =
         "defaultPortGroupCreateState";
   public static final String ID_CREATE_VDS_DEFAULT_PORTGROUP_FIELD =
         "defaultPortGroupName";

   public static final String ID_CREATE_VDS_VERSION_DESCRIPTION = "descriptionText";

   public static final String ID_CREATE_VDS_NAME_AND_LOCATION_STEP =
         "step_nameAndLocationPage";
   public static final String ID_CREATE_VDS_CONFIGURE_DISTRIBUTED_SWITCH_STEP =
         "step_dvsGeneralConfigPage";
   public static final String ID_CREATE_VDS_SELECT_DISTRIBUTED_SWITCH_VERSION_STEP =
         "step_selectDvsVersionPage";
   public static final String ID_CREATE_VDS_READY_TO_COMPLETE_STEP =
         "step_defaultSummaryPage";

   public static final String ID_CREATE_VDS_CONFIGURE_DISTRIBUTED_SWITCH_PAGE =
         ID_WIZARD_TIWO_DIALOG + "/" + "dvsGeneralConfigPage";
   public static final String ID_CREATE_VDS_SELECT_DISTRIBUTED_SWITCH_VERSION_PAGE =
         ID_WIZARD_TIWO_DIALOG + "/" + "selectDvsVersionPage";
   public static final String ID_CREATE_VDS_READY_TO_COMPLETE_PAGE =
         ID_WIZARD_TIWO_DIALOG + "/" + "defaultSummaryPage";

   // vDS alarms IDs
   public static final String VDS_ALARAM_TYPE = "vim.dvs.VmwareDistributedVirtualSwitch";
   public static final String ID_ALARM_PREFIX = "vim.event.";
   public static final String ID_VDS_ALARM_HOST_LEFT = ID_ALARM_PREFIX
         + "DvsHostLeftEvent";
   public static final String ID_VDS_ALARM_PORT_LEAVE_PORTGROUP = ID_ALARM_PREFIX
         + "DvsPortLeavePortgroupEvent";
   public static final String ID_VDS_ALARM_PORT_JOIN_PORTGROUP = ID_ALARM_PREFIX
         + "DvsPortJoinPortgroupEvent";
   public static final String ID_VDS_ALARM_UPGRADE_AVAILABLE = ID_ALARM_PREFIX
         + "DvsUpgradeAvailableEvent";
   public static final String ID_VDS_ALARM_UPGRADE_IN_PROGRESS = ID_ALARM_PREFIX
         + "DvsUpgradeInProgressEvent";
   public static final String ID_VDS_ALARM_UPGRADE_REJECTED = ID_ALARM_PREFIX
         + "DvsUpgradeRejectedEvent";
   public static final String ID_VDS_ALARM_PORT_BLOCKED = ID_ALARM_PREFIX
         + "DvsPortBlockedEvent";
   public static final String ID_VDS_ALARM_PORT_CONNECTED = ID_ALARM_PREFIX
         + "DvsPortConnectedEvent";
   public static final String ID_VDS_ALARM_PORT_CREATED = ID_ALARM_PREFIX
         + "DvsPortCreatedEvent";
   public static final String ID_VDS_ALARM_PORT_DELETED = ID_ALARM_PREFIX
         + "DvsPortDeletedEvent";
   public static final String ID_VDS_ALARM_PORT_DISCONNECTED = ID_ALARM_PREFIX
         + "DvsPortDisconnectedEvent";
   public static final String ID_VDS_ALARM_PORT_ENTERED_PASSTHRU = ID_ALARM_PREFIX
         + "DvsPortEnteredPassthruEvent";
   public static final String ID_VDS_ALARM_PORT_LINK_DOWN = ID_ALARM_PREFIX
         + "DvsPortLinkDownEvent";
   public static final String ID_VDS_ALARM_PORT_LINK_UP = ID_ALARM_PREFIX
         + "DvsPortLinkUpEvent";
   public static final String ID_VDS_ALARM_PORT_EXITED_PASSTHRU = ID_ALARM_PREFIX
         + "DvsPortExitedPassthruEvent";
   public static final String ID_VDS_ALARM_PORT_RECONFIGURED = ID_ALARM_PREFIX
         + "DvsPortReconfiguredEvent";
   public static final String ID_VDS_ALARM_PORT_UNBLOCKED = ID_ALARM_PREFIX
         + "DvsPortUnblockedEvent";
   public static final String ID_VDS_ALARM_HEALTH_STATUS_CHANGE = ID_ALARM_PREFIX
         + "DvsHealthStatusChangeEvent";
   public static final String ID_VDS_ALARM_HOST_JOINED = ID_ALARM_PREFIX
         + "DvsHostJoinedEvent";
   public static final String ID_VDS_ALARM_HOST_STATUS_UPDATE = ID_ALARM_PREFIX
         + "DvsHostStatusUpdated";
   public static final String ID_VDS_ALARM_HOST_OUT_OF_SYNC = ID_ALARM_PREFIX
         + "DvsHostWentOutOfSyncEvent";
   public static final String ID_VDS_ALARM_HOST_IN_SYNC = ID_ALARM_PREFIX
         + "DvsHostBackInSyncEvent";
   public static final String ID_VDS_ALARM_UPGRADED = ID_ALARM_PREFIX
         + "DvsUpgradedEvent";
   public static final String ID_VDS_ALARM_RECONFIGURED = ID_ALARM_PREFIX
         + "DvsReconfiguredEvent";
   public static final String ID_VDS_ALARM_EXPRESSION_ATTRIBUTE_DESCRIPTION =
         "configSpec.description";
   public static final String ID_VDS_ALARM_EXPRESSION_ATTRIBUTE_PORTS =
         "configSpec.numStandalonePorts";
   public static final String ID_VDS_ALARM_EXPRESSION_ATTRIBUTE_NAME = "configSpec.name";

   // DvPortGroup alarms IDs
   public static final String ID_DVPG_ALARM_RECONFIGURED = ID_ALARM_PREFIX
         + "DVPortgroupReconfiguredEvent";
   public static final String ID_DVPG_ALARM_RENAMED = ID_ALARM_PREFIX
         + "DVPortgroupRenamedEvent";

   // IDs in Edit vDS NetFlow wizard
   public static final String ID_EDIT_VDS_NETFLOW_COLLECTOR_IP =
         "collectorIpAddressInput";
   public static final String ID_EDIT_VDS_NETFLOW_PORT = "portNumberStepper";
   public static final String ID_EDIT_VDS_NETFLOW_VDS_IP = "dvsIpAddressInput";
   public static final String ID_EDIT_VDS_NETFLOW_ACTIVE_FLOW =
         "activeFlowExportTimeoutStepper";
   public static final String ID_EDIT_VDS_NETFLOW_IDLE_FLOW =
         "idleFlowExportTimeoutStepper";
   public static final String ID_EDIT_VDS_NETFLOW_SAMPLING_RATE = "samplingRateStepper";
   public static final String ID_EDIT_VDS_NETFLOW_PROCESS_INTERNAL_FLOWS =
         "_NetFlowPage_BooleanDropDownList1";
   public static final String ID_EDIT_VDS_NETFLOW_VDS_SETTINGS_INFO = "dvsSettingsInfo";
   public static final String ID_VDS_NETFLOW_PROPERTY_VIEW = "netFlowPropertyView";

   // IDs in Remove vDS confirmation dialog
   public static final String ID_REMOVE_VDS_CONFIRMATION_DIALOG = "YesNoDialog";
   public static final String ID_REMOVE_VDS_CONFIRM_YES_BUTTON =
         ID_REMOVE_VDS_CONFIRMATION_DIALOG + "/automationName=Yes";
   public static final String ID_REMOVE_VDS_CONFIRM_NO_BUTTON =
         ID_REMOVE_VDS_CONFIRMATION_DIALOG + "/automationName=No";

   // IDs in vDS > Manage > Settings > Properties
   public static final String ID_VDS_PROPERTIES_VIEW = "settingsComponent/propView";

   // vDS > Manage > Settings tab
   public static final String ID_VDS_MANAGE_SETTINGS_TOC = "tocTree";
   public static final String ID_VDS_MANAGE_SETTINGS_EDIT_BUTTON =
         "btn_vsphere.core.dvs.editDvsSettingsAction";

   // IDs in vDS Edit Settings dialog
   public static final String ID_VDS_EDIT_SETTINGS_PAGE = "wizardContent";
   public static final String ID_VDS_EDIT_SETTINGS_CONTENT_PAGE =
         ID_VDS_EDIT_SETTINGS_PAGE + "/" + "pageStack";
   public static final String ID_VDS_EDIT_SETTINGS_DISCOVERY_PROTOCOL_SECTION =
         ID_VDS_EDIT_SETTINGS_PAGE + "/" + "ldpBlock";
   public static final String ID_VDS_EDIT_SETTINGS_ADMIN_CONTACTS_SECTION =
         ID_VDS_EDIT_SETTINGS_PAGE + "/" + "adminContactBlock";
   public static final String ID_VDS_EDIT_SETTINGS_NAME_FIELD =
         ID_VDS_EDIT_SETTINGS_CONTENT_PAGE + "/" + "switchName";
   public static final String ID_VDS_EDIT_SETTINGS_UPLINKS_FIELD =
         ID_VDS_EDIT_SETTINGS_CONTENT_PAGE + "/"
               + "_DvsConfigGeneralPage_NumericStepper1";
   public static final String ID_VDS_EDIT_SETTINGS_UPLINKS_INCREMENT =
         ID_VDS_EDIT_SETTINGS_UPLINKS_FIELD + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_VDS_EDIT_SETTINGS_UPLINKS_DECREMENT =
         ID_VDS_EDIT_SETTINGS_UPLINKS_FIELD + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_VDS_EDIT_SETTINGS_EDIT_UPLINK_NAMES_LINK =
         ID_VDS_EDIT_SETTINGS_CONTENT_PAGE + "/" + "_DvsConfigGeneralPage_LinkButton1";
   public static final String ID_VDS_EDIT_SETTINGS_DESCRIPTION_FIELD =
         ID_VDS_EDIT_SETTINGS_CONTENT_PAGE + "/" + "_DvsConfigGeneralPage_TextAreaEx1";
   public static final String ID_VDS_EDIT_SETTINGS_RESOURCE_ALLOCATION_DROPDOWN =
         ID_VDS_EDIT_SETTINGS_PAGE + "/" + "iormList";
   public static final String ID_VDS_EDIT_SETTINGS_MAXIMUM_MTU_FIELD =
         ID_VDS_EDIT_SETTINGS_CONTENT_PAGE + "/"
               + "_DvsConfigAdvancedPage_NumericStepper1";
   public static final String ID_VDS_EDIT_SETTINGS_DEFAULT_MAX_PORTS_PER_HOST =
         ID_VDS_EDIT_SETTINGS_CONTENT_PAGE + "/"
               + "_DvsConfigAdvancedPage_NumericStepper2";
   public static final String ID_VDS_EDIT_SETTINGS_MAXIMUM_MTU_INCREMENT =
         ID_VDS_EDIT_SETTINGS_MAXIMUM_MTU_FIELD + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_VDS_EDIT_SETTINGS_MAXIMUM_MTU_DECREMENT =
         ID_VDS_EDIT_SETTINGS_MAXIMUM_MTU_FIELD + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_VDS_EDIT_SETTINGS_DP_TYPE_DROPDOWN =
         ID_VDS_EDIT_SETTINGS_CONTENT_PAGE + "/" + "discoveryProtocolType";
   public static final String ID_VDS_EDIT_SETTINGS_DP_OPERATION_DROPDOWN =
         ID_VDS_EDIT_SETTINGS_CONTENT_PAGE + "/" + "ldpOperationTypeList";
   public static final String ID_VDS_EDIT_SETTINGS_ADMIN_NAME_FIELD =
         ID_VDS_EDIT_SETTINGS_CONTENT_PAGE + "/" + "contactNameInput";
   public static final String ID_VDS_EDIT_SETTINGS_OTHER_DETAILS_FIELD =
         ID_VDS_EDIT_SETTINGS_CONTENT_PAGE + "/" + "contactDetailsInput";

   public static final String ID_VDS_EDIT_UPLINK_NAMES_CONTAINER = "uplinkContainer";

   public static final String ID_VDS_MANAGE_PORTS_VIEW =
         "vsphere.core.dvs.manage.portsView";
   public static final String ID_REFRESH_VDS_PORTS_BUTTON =
         "vsphere.core.dvPortgroup.manage.portsView/toolContainer/filterControlContainer/refresh";

   public static final String ID_FILTER_TEXT_INPUT = "textInput";
   public static final String ID_FILTER_SEARCH_BUTTON = "searchButton";

   // IDs in vDS Edit Uplink Names dialog
   public static final String ID_VDS_EDIT_UPLINK_NAMES_DIALOG = "namesForm";
   public static final String ID_UPLINKS_RELATED_ITEMS = "hostsForUplinkPG/list";

   // IDs in vDS Manage Settings Resource Allocation
   public static final String ID_VDS_CREATE_RESOURCE_POOL_BUTTON =
         "vsphere.core.dvs.createResPool/id=button";
   public static final String ID_VDS_EDIT_RESOURCE_POOL_BUTTON =
         "vsphere.core.dvs.editResPool/id=button";
   public static final String ID_VDS_REMOVE_RESOURCE_POOL_BUTTON =
         "vsphere.core.dvs.removeResPool/id=button";

   public static final String ID_YES_CONFIRM_BUTTON = "confirmationDialog/name=YES";
   public static final String ID_NO_CONFIRM_BUTTON = "confirmationDialog/name=NO";

   public static final String ID_NETWORK_RESOURCE_POOL_LIST_COUNT =
         ID_RESOURCE_MANAGEMENT_POOLLIST + "/infoPanel1";

   public static final String ID_NRP_DETAILS_TAB_CONTAINER = "detailsTabContainer";
   public static final String ID_NRP_DVPORTGROUP_MAPPING = "resPoolDvPortgroupMapping";
   public static final String ID_NRP_DVPORTGROUP_COUNT =
         "dvPortgroupMapppingTab/infoPanel1";

   // IDs in vDS New Network Resource Pool
   public static final String ID_VDS_NEW_NRP_NAME_FIELD =
         "_ResPoolPropertiesPage_TextInputEx1";
   public static final String ID_VDS_NEW_NRP_ORIGIN_LABEL =
         "_ResPoolPropertiesPage_Text2";
   public static final String ID_VDS_NEW_NRP_DESCRIPTION_FIELD =
         "_ResPoolPropertiesPage_TextAreaEx1";
   public static final String ID_VDS_NEW_NRP_HOST_LIMIT_NUMERIC_STEPPER =
         "_ResPoolPropertiesPage_NumericStepper1";
   public static final String ID_VDS_NEW_NRP_HOST_LIMIT_CHECKBOX =
         "_ResPoolPropertiesPage_CheckBox1";
   public static final String ID_VDS_NEW_NRP_PHYSICAL_ADAPTER_SHARES_DROPDOWN =
         "_ResPoolPropertiesPage_DropDownList1";
   public static final String ID_VDS_NEW_NRP_PHYSICAL_ADAPTER_SHARES_FIELD =
         "_ResPoolPropertiesPage_NumericStepper2";
   public static final String ID_VDS_NEW_NRP_QOS_TAG_DROPDOWN =
         "_ResPoolPropertiesPage_DropDownList2";

   // IDs in vDS Edit Network Resource Pool
   public static final String ID_VDS_EDIT_NRP_NAME_FIELD =
         "_ResPoolPropertiesPage_Text1";
   public static final String ID_VDS_EDIT_NRP_DESCRIPTION_FIELD =
         "_ResPoolPropertiesPage_Text3";

   public static final String ID_NETWORK_RESOURCE_POOL_DETAILS_LABEL =
         "_ResAllocView_Text1";

   public static final String ID_VDS_RESOURCE_POOL_DETAILS = "detailsContainer";
   public static final String ID_VDS_RESOURCE_POOL_DETAILS_NAME =
         "_ResPoolDetails_Text1";
   public static final String ID_VDS_RESOURCE_POOL_DETAILS_ORIGIN =
         "_ResPoolDetails_Text2";
   public static final String ID_VDS_RESOURCE_POOL_DETAILS_DESCRIPTION =
         "_ResPoolDetails_Text3";
   public static final String ID_VDS_RESOURCE_POOL_DETAILS_HOST_LIMIT =
         "_ResPoolDetails_Text4";
   public static final String ID_VDS_RESOURCE_POOL_DETAILS_ADAPTER_SHARES =
         "_ResPoolDetails_Text5";
   public static final String ID_VDS_RESOURCE_POOL_DETAILS_QOS_TAG =
         "_ResPoolDetails_Text6";

   // IDs in vDS Manage Hosts wizard
   public static final String ID_VDS_MANAGE_HOSTS_MEMBERS_LIST = "hostsList";
   public static final String ID_VDS_MANAGE_HOSTS_MEMBERS_BUTTON =
         "memberHostImageButton";

   // IDs in vDS Add and Manage Hosts wizard
   public static final String ID_LIST_VDS_MANAGE_HOSTS_MEMBERS = "hostsList";
   public static final String ID_LIST_TOOLBAR_VDS_MANAGE_HOSTS_MEMBERS =
         "dataGridToolbar";
   public static final String ID_LIST_VDS_MANAGE_HOSTS_ASSIGN_UPLINK =
         "uplinkPortSelector";
   public static final String ID_LIST_VDS_MANAGE_HOSTS_ASSIGN_PORTGROUP =
         "portGroupsList";
   public static final String ID_LIST_VDS_MANAGE_HOSTS_ASSIGN_PORTGROUP_NETWORKS =
         "networkList";
   public static final String ID_LIST_VDS_MANAGE_HOSTS_HOST_TREE = "hostTreeList";
   public static final String ID_LIST_VDS_MANAGE_HOSTS_VIRTUAL_ADAPTERS =
         "virtualAdaptersList";
   public static final String ID_LIST_VDS_MANAGE_HOSTS_VALIDATION = "validationList";

   public static final String ID_BUTTON_VDS_MANAGE_HOSTS_NEW_HOSTS =
         "newHostImageButton";
   public static final String ID_BUTTON_VDS_MANAGE_HOSTS_MEMBER_HOSTS =
         "memberHostImageButton";
   public static final String ID_BUTTON_VDS_MANAGE_HOSTS_ASSIGN_UPLINK =
         "assignUplinkButtoon";
   public static final String ID_BUTTON_VDS_MANAGE_HOSTS_ASSIGN_PORT_GROUP =
         "assignPortgroupButton";
   public static final String ID_BUTTON_VDS_MANAGE_HOSTS_UNASSIGN_PORT_GROUP =
         "unassignPortgroupButton";
   public static final String ID_BUTTON_VDS_MANAGE_HOSTS_REMOVE_HOST =
         "removeHostImageButton";
   public static final String ID_BUTTON_VDS_MANAGE_HOSTS_VIEW_ADAPTER_SETTINGS =
         "showSelectedObjDetailsButton";

   public static final String ID_RADIO_VDS_MANAGE_HOSTS_ADD_HOSTS =
         "addHostsOnly/radioBtn";
   public static final String ID_RADIO_VDS_MANAGE_HOSTS_MIGRATE_HOSTS =
         "mngHostsOnly/radioBtn";
   public static final String ID_RADIO_VDS_MANAGE_HOSTS_REMOVE_HOSTS =
         "removeHosts/radioBtn";
   public static final String ID_RADIO_VDS_MANAGE_HOSTS_ADD_AND_MIGRATE_HOSTS =
         "addAndMngHosts/radioBtn";

   public static final String ID_LABEL_VDS_MANAGE_HOSTS_VALIDATION_STATUS =
         "overallStatusText";
   public static final String ID_LABEL_VDS_MANAGE_HOSTS_VALIDATION_MESSAGE =
         "propViewMessageText";

   public static final String ID_PAGE_VDS_MANAGE_HOSTS_VALIDATE_CHANGES =
         "viewVnicDependenciesPage";
   public static final String ID_PAGE_VDS_MANAGE_HOSTS_SELECT_VNICS =
         "selectVirtualAdaptersPage";
   public static final String ID_PAGE_VDS_MANAGE_HOSTS_SELECT_PNICS =
         "selectPhysicalAdaptersPage";
   public static final String ID_PAGE_VDS_MANAGE_HOSTS_SELECT_TASK =
         "selectOperationPage";
   public static final String ID_PAGE_VDS_MANAGE_HOSTS_SELECT_HOST = "selectHostsPage";
   public static final String ID_PAGE_VDS_MANAGE_HOSTS_SELECT_VM_NETWORKING =
         "selectVirtualMachinesPage";
   public static final String ID_PAGE_VDS_MANAGE_HOSTS_SELECT_DEFAULT_SUMMARY =
         "defaultSummaryPage";

   public static final String ID_DIALOG_VDS_MANAGE_HOSTS_ASSIGN_PORTGROUP =
         "className=NetworkSelectorDialog";
   public static final String ID_DIALOG_VDS_MANAGE_HOSTS_ASSIGN_UPLINK =
         "assignUplinkDialog";
   public static final String ID_DIALOG_VDS_MANAGE_HOSTS_SELECT_HOSTS =
         "className=DvsSelectHostsDialog";

   public static final String ID_ALERT_VDS_MANAGE_HOSTS_OK_CANCEL = "OkCancelDialog";

   // IDs in Host > Manage > Networking > Virtual switches >> Manage Pnics Dialog for VSS/VDS
   public static final String ID_LIST_ASSIGN_ADAPTER = "connecteeList/list";
   public static final String ID_LIST_ASSIGN_ADAPTER_FAILOVER = "failoverOrder/list";
   public static final String ID_BUTTON_VSS_ADD_ADAPTER = "automationName=Add adapters";
   public static final String ID_BUTTON_VSS_REMOVE_ADAPTER =
         "automationName=Remove selected";
   public static final String ID_BUTTON_VDS_ADD_ADAPTER = "automationName=Add adapter";
   public static final String ID_BUTTON_VDS_REMOVE_ADAPTER =
         "automationName=Remove selected adapters";
   public static final String ID_BUTTON_MOVE_UP_ADAPTER = "automationName=Move Up";
   public static final String ID_BUTTON_MOVE_DOWN_ADAPTER = "automationName=Move Down";
   public static final String ID_BUTTON_INVOKE_MANAGE_PNIC =
         "vsphere.core.host.network.managePhysicalNicsAction/button";

   // IDs in Host > Manage > Networking > Virtual switches >> Migrate Vnic to VSS Wizard
   public static final String ID_PAGE_SELECT_VNIC = "nicListPage";
   public static final String ID_PAGE_CONFIGURE_PORTGROUP = "configPage";
   public static final String ID_LIST_VNIC = ID_PAGE_SELECT_VNIC + "/list";
   public static final String ID_TEXTBOX_NETWORK_LABEL = "configPage/networkLabel";

   // IDs in Host > Manage > Networking > Virtual adapters >> Remove Vnic Dialog
   public static final String ID_DIALOG_REMOVE_VNIC =
         "className=VmkernelDeleteConfirmDialog";
   public static final String ID_DIALOG_REMOVE_VNIC_VALIDATE_CHANGES =
         "className=DependencyReviewDialog";
   public static final String ID_PAGE_REMOVE_VNIC = ID_DIALOG_REMOVE_VNIC + "/content";
   public static final String ID_PAGE_REMOVE_VNIC_VALIDATE_CHANGES =
         ID_DIALOG_REMOVE_VNIC_VALIDATE_CHANGES + "/contentVGroup";
   public static final String ID_LABEL_INFO_MESSAGE = "adapterTextField";
   public static final String ID_BUTTON_VALIDATE_CHANGES = "showDependenciesButton";

   // Related Objects -> Subtabs button ids
   public static final String ID_RELATED_OBJECTS_HOSTS_TAB_BUTTON =
         "hostsForDVPG.button";
   public static final String ID_RELATED_OBJECTS_VIRTUAL_MACHINES_TAB_BUTTON =
         "vmsForDVPG.button";
   public static final String ID_RELATED_OBJECTS_VM_TEMPLATES_TAB_BUTTON =
         "vmTemplatesForDVPG.button";

   // IDs of Manage Distributed Port Groups Wizard (MDPGW)

   // Wizard Page Ids
   public static final String ID_MDPGW_SELECT_PORT_GROUP_POLICIES_PAGE_ID =
         "policyTypePage";
   public static final String ID_MDPGW_SELECT_PORT_GROUPS_PAGE_ID = "pgListPage";
   public static final String ID_MDPGW_SECURITY_PAGE_ID = "securityPolicyPage";
   public static final String ID_MDPGW_TRAFFIC_SHAPING_PAGE_ID =
         "trafficShapingPolicyPage";
   public static final String ID_MDPGW_VLAN_PAGE_ID = "vlanPolicyPage";
   public static final String ID_MDPGW_TEAMING_AND_FAILOVER_PAGE_ID =
         "failoverPolicyPage";
   public static final String ID_MDPGW_RESOURCE_ALLOCATION_PAGE_ID = "netIormPolicyPage";
   public static final String ID_MDPGW_MONITORING_PAGE_ID = "monitoringPolicyPage";
   public static final String ID_MDPGW_MISCELLANEOUS_PAGE_ID = "miscPoliciesPage";
   public static final String ID_MDPGW_READY_TO_COMPLETE_PAGE_ID = "defaultSummaryPage";
   public static final String ID_MDPGW_FAILOVER_ORDER_LIST = "list";

   // ToC Ids
   public static final String ID_MDPGW_SELECT_PORT_GROUP_POLICIES_TOC =
         "step_policyTypePage";
   public static final String ID_MDPGW_SELECT_PORT_GROUPS_TOC = "step_pgListPage";
   public static final String ID_MDPGW_CONFIGURE_POLICIES_TOC = "step_2";
   public static final String ID_MDPGW_SECURITY_TOC = "step_securityPolicyPage";
   public static final String ID_MDPGW_TRAFFIC_SHAPING_TOC =
         "step_trafficShapingPolicyPage";
   public static final String ID_MDPGW_VLAN_TOC = "step_vlanPolicyPage";
   public static final String ID_MDPGW_TEAMING_AND_FAILOVER_TOC =
         "step_failoverPolicyPage";
   public static final String ID_MDPGW_RESOURCE_ALLOCATION_TOC =
         "step_netIormPolicyPage";
   public static final String ID_MDPGW_MONITORING_TOC = "step_monitoringPolicyPage";
   public static final String ID_MDPGW_MISCELLANEOUS_TOC = "step_miscPoliciesPage";
   public static final String ID_MDPGW_READY_TO_COMPLETE_TOC = "step_defaultSummaryPage";

   // Security page
   public static final String ID_MDPGW_SECURITY_PROMISCUOUS_MODE_DROPDOWN_LIST =
         "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList1" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_MDPGW_SECURITY_MAC_ADDRESS_CHANGES_DROPDOWN_LIST =
         "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList2" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_MDPGW_SECURITY_FORGED_TRANSMITS_DROPDOWN_LIST =
         "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList3" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   // Traffic shaping page
   public static final String ID_NUMMERIC_STEPPER_TEXT_FIELD = "textDisplay";
   public static final String ID_MDPGW_INGRESS_POLICY_PAGE = "ingressPage";
   public static final String ID_MDPGW_INGRESS_STATUS_DROP_DOWN =
         ID_MDPGW_INGRESS_POLICY_PAGE + "/"
               + "_DvPortTrafficShapingPolicyPage_BooleanDropDownList1";
   public static final String ID_MDPGW_INGRESS_AVERAGE_BANDWIDTH_FIELD =
         ID_MDPGW_INGRESS_POLICY_PAGE + "/" + "averageNumStepper";
   public static final String ID_MDPGW_INGRESS_PEAK_BANDWIDTH_FIELD =
         ID_MDPGW_INGRESS_POLICY_PAGE + "/" + "peakNumStepper";
   public static final String ID_MDPGW_INGRESS_BURST_SIZE_FIELD =
         ID_MDPGW_INGRESS_POLICY_PAGE + "/" + "burstNumStepper";
   public static final String ID_MDPGW_EGRESS_POLICY_PAGE = "egressPage";
   public static final String ID_MDPGW_EGRESS_STATUS_DROP_DOWN =
         ID_MDPGW_EGRESS_POLICY_PAGE + "/"
               + "_DvPortTrafficShapingPolicyPage_BooleanDropDownList1";
   public static final String ID_MDPGW_EGRESS_AVERAGE_BANDWIDTH_FIELD =
         ID_MDPGW_EGRESS_POLICY_PAGE + "/" + "averageNumStepper";
   public static final String ID_MDPGW_EGRESS_PEAK_BANDWIDTH_FIELD =
         ID_MDPGW_EGRESS_POLICY_PAGE + "/" + "peakNumStepper";
   public static final String ID_MDPGW_EGRESS_BURST_SIZE_FIELD =
         ID_MDPGW_EGRESS_POLICY_PAGE + "/" + "burstNumStepper";
   public static final String ID_MDPGW_UNACCEPTABLE_PEAK_AVERAGE_VALIDATION_ERROR =
         "_errorLabel";
   public static final String ID_MDPGW_LOAD_BALANCING_ICON = "loadBalancingDetails";
   public static final String ID_MDPGW_TEAMING_AND_FAILOVER_WARNING_ICON =
         "propViewIconImage";
   // Vlan page
   public static final String ID_MDPGW_VLAN_POLICY_PAGE = "vlanPolicyPage";
   public static final String ID_MDPGW_VLAN_MODE_DROPDOWN_LIST = "vlanTypeSelector";
   public static final String ID_MDPGW_VLAN_POLICY_PAGE_VLAN_ID_FIELD =
         ID_MDPGW_VLAN_POLICY_PAGE + "/" + "_DvPortVlanPolicyPage_NumericStepper1";
   public static final String ID_MDPGW_PORT_VLAN_POLICY_PAGE_TRUNK_RANGE_TEXT_INPUT =
         "vlanIdRanges";
   // Teaming and failover page
   public static final String ID_MDPGW_LOAD_BALACING_DROPDOWN_LIST =
         "loadBalancingSelector";
   public static final String ID_MDPGW_NETWORK_FAILURE_DETECTION_DROPDOWN_LIST =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList1" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_MDPGW_NOTIFY_SWITCHES_DROPDOWN_LIST =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList2" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_MDPGW_FAILBACK_DROPDOWN_LIST =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList3" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_MDPGW_FAILOVER_MOVE_DOWN_BUTTON =
         "_DvPortFailoverPolicyPage_Button2";
   public static final String ID_MDPGW_FAILOVER_MOVE_UP_BUTTON =
         "_DvPortFailoverPolicyPage_Button1";
   public static final String ID_MDPGW_FAILOVER_ORDER_ADV_GRID = "failoverOrderEditor";
   public static final String ID_MESSAGE_TEXT_ID = "propViewMessageText";
   public static final String ID_ADDITIONAL_TIP_ID = "additionalTip";
   public static final String ID_CONFIRM_DIALOG_BUTTON_YES = "confirmButton";
   public static final String ID_CONFIRM_DIALOG_BUTTON_NO = "rejectButton";
   // Resource Allocation page
   public static final String ID_MDPGW_NETWORK_RESOUCRE_POOL_DROPDOWN_LIST =
         "poolSelector";
   // Monitoring page
   public static final String ID_MDPGW_NET_FLOW_DROPDOWN =
         "_DvPortMonitoringPolicyPage_OverridableBooleanDropDownList1";
   public static final String ID_MDPGW_NET_FLOW_DROPDOWN_LIST =
         ID_MDPGW_NET_FLOW_DROPDOWN + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_MDPGW_NET_FLOW_CHECKBOX = ID_MDPGW_NET_FLOW_DROPDOWN
         + "/" + "_OverridableBooleanDropDownList_CheckBox1";
   // Miscellaneous page
   public static final String ID_MDPGW_BLOCK_ALL_PORTS_DROPDOWN_LIST =
         "_DvPortMiscPoliciesPage_OverridableBooleanDropDownList1" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_MDPGW_BLOCK_ALL_PORTS_TEXT =
         "_DvPortMiscPoliciesPage_RichEditableText1";
   // Forged Transmits
   public static final String ID_MDPGW_CONFIGURE_POLICY_STEP_NUMBER = "3";
   public static final String ID_MDPGW_STEP_DISPLAY_LABEL = "substepLabelDisplay";
   public static final String ID_MDPGW_STEP_INDEX_LABEL = "subStepIndexDisplay";
   public static final String ID_MDPGW_STEP_LETTER_LABEL = "subStepLetterDisplay";
   public static final String ID_MDPGW_DATA_GROUP_PAGE = "dataGroup";
   public static final String ID_MDPGW_PORT_GROUPS_GRID = "list";
   public static final String ID_MDPGW_NAME_GRID_COLUMN = "Name";

   // IDs of UI components present on dvPortGroup -> Getting Started tab
   public static final String ID_DISTRIBUTED_PORT_GROUP_QUESTION_LABEL = "titleLabel";
   public static final String ID_DISTRIBUTED_PORT_GROUP_DESCRIPTION_LABEL =
         "descriptionLabel";
   public static final String ID_LEARN_MORE_ABOUT_DISTRIBUTED_PORT_GROUP_LINK =
         "gettingStartedHelpLink_1" + "/" + "link";
   public static final String ID_LEARN_HOW_TO_SETUP_A_NETWORK_LINK =
         "gettingStartedHelpLink_2" + "/" + "link";
   public static final String ID_DISTRIBUTED_PORT_GROUP_IMAGE = ID_GETTING_STARTED_IMAGE;

   // IDs of Policy Portlet
   public static final String ID_POLICY_PORTLET_SECURITY_STACK_BLOCK =
         "securityPolicyBlock";
   public static final String ID_POLICY_PORTLET_SECURITY_PROMISCUOUS_MODE_LABEL =
         "promiscuousMode";
   public static final String ID_POLICY_PORTLET_SECURITY_MAC_ADDRESS_CHANGES_LABEL =
         "macAddressChanges";
   public static final String ID_POLICY_PORTLET_SECURITY_FORGED_TRANSMITS_LABEL =
         "forgedTransmits";

   // Portlet policy VLAN
   public static final String ID_POLICY_PORTLET_VLAN_TYPE_LABEL =
         "summary_vlanProp_valueLbl";
   public static final String ID_POLICY_PORTLET_VLAN_STACK_BLOCK = "vlanBlock";

   // Portlet policy Egress
   public static final String ID_POLICY_PORTLET_EGRESS_TRAFFIC_SHAPING_STACK_BLOCK =
         "outShapingBlock";
   public static final String ID_POLICY_PORTLET_EGRESS_TRAFFIC_SHAPING_PEAK_BANDWIDTH_LABEL =
         "outShapingPeakBW";
   public static final String ID_POLICY_PORTLET_EGRESS_TRAFFIC_BURST_SIZE_LABEL =
         "outShapingBurstSize";
   public static final String ID_POLICY_PORTLET_EGRESS_TRAFFIC_SHAPING_AVERAGE_BANDWIDTH_LABEL =
         "outShapingAvgBW";
   public static final String ID_POLICY_PORTLET_EGRESS_TRAFFIC_SHAPING_STATUS_LABEL =
         "outShapingStatus";

   // Portlet policy Ingress
   public static final String ID_POLICY_PORTLET_INGRESS_TRAFFIC_SHAPING_STACK_BLOCK =
         "inShapingBlock";
   public static final String ID_POLICY_PORTLET_INGRESS_TRAFFIC_SHAPING_PEAK_BANDWIDTH_LABEL =
         "inShapingPeakBW";
   public static final String ID_POLICY_PORTLET_INGRESS_TRAFFIC_BURST_SIZE_LABEL =
         "inShapingBurstSize";
   public static final String ID_POLICY_PORTLET_INGRESS_TRAFFIC_SHAPING_AVERAGE_BANDWIDTH_LABEL =
         "inShapingAvgBW";
   public static final String ID_POLICY_PORTLET_INGRESS_TRAFFIC_SHAPING_STATUS_LABEL =
         "inShapingStatus";

   // Portlet policy Failover
   public static final String ID_POLICY_PORTLET_TEAMING_AND_FAILOVER_ACTIVE_UPLINKS_LABEL =
         "activeUplinks";
   public static final String ID_POLICY_PORTLET_TEAMING_AND_FAILOVER_STANDBY_UPLINKS_LABEL =
         "standbyUplinks";
   public static final String ID_POLICY_PORTLET_TEAMING_AND_FAILOVER_UNUSED_UPLINKS_LABEL =
         "unusedUplinks";
   public static final String ID_POLICY_PORTLET_TEAMING_AND_FAILOVER_LOAD_BALANCING =
         "teamingLoadBalancing";
   public static final String ID_POLICY_PORTLET_NETWORK_FAILURE_DETECTION =
         "teamingNWFailoverDetecetion";
   public static final String ID_POLICY_PORTLET_NOTIFY_SWITCHES =
         "teamingNotifySwitches";
   public static final String ID_POLICY_PORTLET_FAILBACK = "teamingFailback";
   public static final String ID_POLICY_PORTLET_TEAMING_AND_FAILOVER_STACK_BLOCK =
         "teamingBlock";

   // IDs of DvPortSettingsView and DvPortgroupSettingsView components

   // Id of Ports View which contains DvPortSettingsView component
   public static final String ID_DVPORT_GROUP_MANAGE_PORTS_SUB_TAB_VIEW =
         "vsphere.core.dvPortgroup.manage.portsView";
   public static final String ID_PORTS_VIEW_PORT_SETTINGS_POLICY_STACK_BLOCK =
         "_DvPortSettingsView_StackBlock4";
   public static final String ID_PORTS_VIEW_PORT_SETTINGS_PROPERTIES_STACK_BLOCK =
         "_DvPortSettingsView_StackBlock1";
   public static final String ID_PORTS_VIEW_PORT_SETTINGS_PROPERTIES_NAME_LABEL =
         "propViewValueText";
   public static final String ID_PORTS_VIEW_PORT_SETTINGS_PROPERTIES_DESCRIPTION_LABEL =
         "_DvPortSettingsView_Text2";
   public static final String ID_PORTS_VIEW_PORT_SETTINGS_PROPERTIES_RESOURCE_POOL_LABEL =
         "_DvPortSettingsView_Text3";

   // Id of Settings View which contain DvPortgroupSettingsView component
   public static final String ID_DVPOTRT_GROUP_MANAGE_SETTINGS_TAB_VIEW =
         "vsphere.core.dvPortgroup.manage.settingsView";
   public static final String ID_SETTINGS_VIEW_POLICY_STACK_BLOCK =
         "_DvPortgroupSettings_StackBlock2";

   // IDs of General Properties
   public static final String ID_DVPG_MANAGE_SETTINGS_PROPERTIES_VIEW =
         "vsphere.core.dvPortgroup.manage.propertiesView/settingsComponent/portgroupPropertyView";
   public static final String ID_DVPG_MANAGE_SETTINGS_POLICIES_VIEW =
         "vsphere.core.dvPortgroup.manage.policiesView/settingsComponent/portgroupPropertyView";
   public static final String ID_MANAGE_PORTS_PROPERTY_VIEW =
         "_PortsView_Scroller1/propView";
   public static final String ID_EDIT_DVPORTGROUP_SETTINGS_BUTTON_MANAGE_TAB = "btn_"
         + ID_ACTION_EDIT_SETTINGS_DVPORTGROUP;
   public static final String ID_VIEW_PROPERTIES_STACK_BLOCK =
         "_DvPortgroupSettings_StackBlock1";
   public static final String ID_VIEW_GENERAL_STACK_BLOCK =
         "_DvPortgroupPropertiesView_DvPortgroupGeneralPropertiesView1";
   public static final String ID_PORT_GROUP_NAME_LABEL =
         "_DvPortgroupGeneralPropertiesView_Text1";
   public static final String ID_PORT_BINDING_LABEL =
         "_DvPortgroupGeneralPropertiesView_Text2";
   public static final String ID_NUMBER_OF_PORTS_LABEL =
         "_DvPortgroupGeneralPropertiesView_Text4";
   public static final String ID_PORT_ALLOCATION_LABEL =
         "_DvPortgroupGeneralPropertiesView_Text3";
   public static final String ID_NETWORK_RESOURCE_POOL_LABEL =
         "_DvPortgroupGeneralPropertiesView_Text5";
   public static final String ID_VIEW_GENERAL_ADVANCED_BLOCK =
         "_DvPortgroupPropertiesView_DvPortgroupAdvancedPropertiesView1";
   public static final String ID_PORT_GROUP_DESCRIPTION_LABEL =
         "_DvPortgroupGeneralPropertiesView_Text6";
   public static final String ID_UPLINK_PORT_GROUP_STATE_LABEL =
         "_DvPortSettingsView_LabelEx1";
   public static final String ID_PORTGROUP_POLICIES_TOC_LIST =
         "vsphere.core.dvPortgroup.manage.settings.wrapper/tocTree";
   public static final String ID_DVPG_MANAGE_PORTS_PROPERTIES_VIEW =
         "vsphere.core.dvPortgroup.manage.portsView/portSettings/propView";


   // vDS Switch Detils portlet values
   public static final String ID_VDS_SUMMARY_TAB_DETAILS_PORTLET = ID_DC_DETAILS_PORTLET
         .replaceFirst(EXTENSION_ENTITY, EXTENSION_ENTITY_DV_SWITCH);
   public static final String ID_VDS_SUMMARY_TAB_SWITCH_MANUFACTURER =
         "summary_manifactureProp_valueLbl";
   public static final String ID_VDS_SUMMARY_TAB_SWITCH_VERSION =
         "summary_versionProp_valueLbl";
   public static final String ID_VDS_SUMMARY_TAB_DETAILS_PORTS =
         ID_VDS_SUMMARY_TAB_DETAILS_PORTLET + "/ports";
   public static final String ID_VDS_SUMMARY_TAB_DETAILS_AVAILABLE_PORTS =
         ID_VDS_SUMMARY_TAB_DETAILS_PORTLET + "/availablePorts";
   public static final String ID_VDS_SUMMARY_TAB_DETAILS_TOTAL_PORTS =
         ID_VDS_SUMMARY_TAB_DETAILS_PORTLET + "/totalPorts";
   public static final String ID_VDS_SUMMARY_TAB_DETAILS_NETWORKS =
         ID_VDS_SUMMARY_TAB_DETAILS_PORTLET + "/networks";
   public static final String ID_VDS_SUMMARY_TAB_DETAILS_HOSTS =
         ID_VDS_SUMMARY_TAB_DETAILS_PORTLET + "/hosts";
   public static final String ID_VDS_SUMMARY_TAB_DETAILS_VMS =
         ID_VDS_SUMMARY_TAB_DETAILS_PORTLET + "/vms";

   public static final String ID_VDS_RESOURCE_ALLOCATION_TAB =
         "vsphere.core.dvs.manage.resAllocationView";
   public static final String ID_VDS_RESOURCE_ALLOCATION_TAB_NETWORK_IO_STATUS =
         ID_VDS_RESOURCE_ALLOCATION_TAB + "/resourceManagementStatusLbl";

   // vDS Health check
   public static final String ID_HEALTH_CHECK_ACTION_BUTTON =
         "btn_vsphere.core.dvs.editHealthCheckAction";
   public static final String ID_HEALTH_CHECK_VLAN_MTU_DROPDOWN = "vlanMtuComboBox";
   public static final String ID_HEALTH_CHECK_TEAMING_FAILOVER_DROPDOWN =
         "teamingFailoverComboBox";
   public static final String ID_HEALTH_CHECK_VLAN_MTU_LABEL = "vlanMtuLabel";
   public static final String ID_HEALTH_CHECK_TEAMING_FAILOVER_LABEL =
         "teamingFailoverLabel";
   public static final String ID_HEALTH_CHECK_HOST_STATUS_GRID = "hostListView" + "/"
         + "datagrid";
   public static final String ID_HEALTH_CHECK_VIEW = "vsphere.core.dvs.healthCheckView";
   public static final String ID_HEALTH_CHECK_SETTINGS_PROPERTIES_VIEW =
         ID_HEALTH_CHECK_VIEW + "/propView";
   public static final String HEALTH_CHECK_MTU_STATUS = ".mtuStatus";
   public static final String HEALTH_CHECK_VLAN_STATUS = ".vlanStatus";
   public static final String HEALTH_CHECK_VDS_STATUS = ".vdsLocalizedStatus";
   public static final String HEALTH_CHECK_TEAMING_STATUS = ".teamingStatus";
   public static final String HEALTH_CHECK_HOST_STATE = ".hostState";
   public static final String HEALTH_CHECK_DATA_PROVIDER = "dataProvider.list.source.";
   public static final String HEALTH_CHECK_HOST_NAME = ".runtime.hostName.value";
   public static final String ID_HEALTH_CHECK_MONITORING_VIEW =
         "vsphere.core.dvs.monitor.healthCheckView";
   public static final String HEALTH_CHECK_REFRESH_BUTTON =
         ID_HEALTH_CHECK_HOST_STATUS_GRID + "/" + ID_TOOLBAR_DATAGRID_REFRESH_BUTTON;
   public static final String HEALTH_CHECK_DETAILS_TAB_NAVIGATOR =
         ID_HEALTH_CHECK_MONITORING_VIEW + "/advTabBar/tabBar";
   public static final String HEALTH_CHECK_VLAN_GRID = "vlanListView/datagrid";
   public static final String HEALTH_CHECK_TEAMING_PROPERTY_VIEW =
         "teamingFailoverView/propertyView";
   public static final String HEALTH_CHECK_MTU_GRID = "mtuListView/datagrid";

   // IDs for "Create Distributed Port Group" wizard
   public static final String ID_CREATE_PORT_GROUP_SELECT_NAME_AND_LOCATION_PAGE =
         "nameAndLocationPage";
   public static final String ID_CREATE_PORT_GROUP_PROPERTIES_PAGE = "propertiesPage";
   public static final String ID_CREATE_PORT_GROUP_SECURITY_PAGE = "securityPolicyPage";
   public static final String ID_CREATE_PORT_GROUP_TRAFFIC_SHAPING_PAGE =
         "trafficShapingPolicyPage";
   public static final String ID_CREATE_PORT_GROUP_TEAMING_AND_FAILOVER_PAGE =
         "failoverPolicyPage";
   public static final String ID_CREATE_PORT_GROUP_MONITORING_PAGE =
         "monitoringPolicyPage";
   public static final String ID_CREATE_PORT_GROUP_MISCELLENEOUS_PAGE =
         "miscPoliciesPage";
   public static final String ID_CREATE_PORT_GROUP_ADVANCED_PROPERTIES_PAGE =
         "advancedPropertiesPage";
   public static final String ID_CREATE_PORT_GROUP_FINAL_PAGE = "finalPage";
   public static final String ID_CREATE_PORT_GROUP_CUSTOMIZE_PAGE = "policiesPage";
   public static final String ID_FAILOVER_ORDER =
         "_DvPortFailoverPolicyPage_StackBlock5";

   public static final String ID_CREATE_PORT_GROUP_READY_TO_COMPLETE =
         IDConstants.ID_CREATE_PORT_GROUP_FINAL_PAGE + "/className=PropertyView";

   // vDS Features portlet support statuses
   public static final String ID_VDS_SUMMARY_TAB_FEATURES_PORTLET =
         "vsphere.core.dvs.summary.featuresView.chrome";
   public static final String ID_VDS_SUMMARY_TAB_FEATURES_PROPERTY_GRID =
         ID_VDS_SUMMARY_TAB_FEATURES_PORTLET + "/dvsFeaturesGrid";
   public static final String ID_VDS_SUMMARY_TAB_FEATURES_NETIOCTRL_STATUS =
         ID_VDS_SUMMARY_TAB_FEATURES_PROPERTY_GRID + "/netiorm";
   public static final String ID_VDS_SUMMARY_TAB_FEATURES_USERDEFRP_STATUS =
         ID_VDS_SUMMARY_TAB_FEATURES_PROPERTY_GRID + "/userDefResourcePools";
   public static final String ID_VDS_SUMMARY_TAB_FEATURES_DIRECTPATHIO_STATUS =
         ID_VDS_SUMMARY_TAB_FEATURES_PROPERTY_GRID + "/npt";
   public static final String ID_VDS_SUMMARY_TAB_FEATURES_NETFLOW_STATUS =
         ID_VDS_SUMMARY_TAB_FEATURES_PROPERTY_GRID + "/ipfix";
   public static final String ID_VDS_SUMMARY_TAB_FEATURES_LLDP_STATUS =
         ID_VDS_SUMMARY_TAB_FEATURES_PROPERTY_GRID + "/lldp";
   public static final String ID_VDS_SUMMARY_TAB_FEATURES_PORTMIRRORING_STATUS =
         ID_VDS_SUMMARY_TAB_FEATURES_PROPERTY_GRID + "/vspan";

   public static final String ID_VDS_SUMMARY_TAB_ISSIUE_DESCRIPTION_LABEL =
         "vsphere.core.dvs.summaryView/issueSummary/issueDescription0";

   // Performance Chart
   public static final String ID_LINK_CHART_OPTIONS = "chartOptionsLink";
   public static final String ID_ADG_COUNTER_GROUP = "counterGroupsGrid";
   public static final String ID_ADG_OBJECT_SELECT = "selectObjectGrid";
   public static final String ID_ADG_COUNTER_SELECT = "selectCountersGrid";
   public static final String ID_LIST_TIME_SELECT = "timeIntervalSelect";
   public static final String ID_LIST_CHART_SELECT = "chartTypeSelect";
   public static final String ID_BUTTON_SUBMIT = "submitBtn";
   public static final String ID_BUTTON_CANCEL = "cancelBtn";
   public static final String ID_BUTTON_COUNTER_ALL = "selectCountersAll";
   public static final String ID_BUTTON_COUNTER_NONE = "selectCountersNone";
   public static final String ID_BUTTON_OBJECT_ALL = "selectObjectsAll";
   public static final String ID_BUTTON_OBJECT_NONE = "selectObjectsNone";
   public static final String ID_LEGEND_PERF_CHART = "legend";
   public static final String ID_PERFORMANCE_CHART_OPTIONS_DIALOG =
         "advPerfChartsConfigDlg";
   public static final String ID_PERFORMANCE_CHART_OPTIONS_HELP =
         ID_PERFORMANCE_CHART_OPTIONS_DIALOG + "/helpBtn";
   public static final String ID_PERFORMANCE_CHART_OPTIONS_LAST_RADIO =
         ID_PERFORMANCE_CHART_OPTIONS_DIALOG + "/customLast";
   public static final String ID_PERFORMANCE_CHART_OPTIONS_FROM_TO_RADIO =
         ID_PERFORMANCE_CHART_OPTIONS_DIALOG + "/customFromTo";
   public static final String ID_PERFORMANCE_CHART_OPTIONS_LAST_NUMERIC_STEPPER =
         ID_PERFORMANCE_CHART_OPTIONS_DIALOG + "/customLastNumStep";
   public static final String ID_PERFORMANCE_CHART_ADVANCED = "dateTimeChart";
   public static final String ID_EXPORT_MENU = "exportMenu";
   public static final String ID_PERFORMANCE_CHART_ADVANCED_TYPE_DROPDOWN =
         "availableCounterGroupsDropDown";
   public static final String ID_TOC_TREE_PERFORMANCE_CHARTS = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY).append("monitor.performance")
         .append("/tocTree").toString();

   // Performance Charts - Overview
   public static final String ID_PERFORMANCE_CHARTS_MAIN_CONTAINER =
         "vsphere.opsmgmt.performance." + EXTENSION_ENTITY + "performanceView";
   public static final String ID_PERFORMANCE_CHARTS_MAIN_SCROLLABLE_CONTAINER =
         "vsphere.core." + EXTENSION_ENTITY
               + "monitor.performance.overviewView.container";
   public static final String ID_PERFORMANCE_CHARTS_SCROLLABLE_CONTAINER =
         "scrollerGroup";
   public static final String ID_PERF_CHARTS_VIEW_COMBO_BOX = "modeList";
   public static final String ID_PERF_CHARTS_TIME_RANGE_COMBO_BOX = "timeRangeList";
   public static final String ID_PERF_CHARTS_THUMBNAILS_CONTAINER = "thumbnails";
   public static final String ID_PERF_CHARTS_OVERVIEW_TITLE =
         ID_PERFORMANCE_CHARTS_MAIN_CONTAINER + "/title";
   public static final String ID_PERF_CHARTS_OVERVIEW_SECOND_TITLE =
         ID_PERFORMANCE_CHARTS_MAIN_CONTAINER + "/secondTitle";
   public static final String ID_PERF_CHARTS_APPLY_TIME_RANGE_BUTTON =
         "btnApplyTimeRange";
   public static final String ID_PERF_CHARTS_OVERVIEW_CUSTOM_TIME_RANGE_DATE_FROM =
         "fieldFrom/dateField";
   public static final String ID_PERF_CHARTS_OVERVIEW_CUSTOM_TIME_RANGE_DATE_TO =
         "fieldTo/dateField";
   public static final String ID_PERF_CHARTS_CUSTOM_TIME_RANGE_TIME_FROM =
         "fieldFrom/timeField";
   public static final String ID_PERF_CHARTS_CUSTOM_TIME_RANGE_TIME_TO =
         "fieldTo/timeField";
   public static final String ID_PERF_CHARTS_THUMBNAIL_NAVIGATOR_PAGE_INFO_LABEL =
         "pageInfo";
   public static final String ID_PERF_CHARTS_THUMBNAIL_CONTAINER = "thumbnails";
   public static final String ID_PERF_CHARTS_THUMBNAIL_NAVIGATOR_FIRST_PAGE_BUTTON =
         "firstPage";
   public static final String ID_PERF_CHARTS_THUMBNAIL_NAVIGATOR_PREVIOUS_PAGE_BUTTON =
         "prevPage";
   public static final String ID_PERF_CHARTS_THUMBNAIL_NAVIGATOR_LAST_PAGE_BUTTON =
         "lastPage";
   public static final String ID_PERF_CHARTS_THUMBNAIL_NAVIGATOR_NEXT_PAGE_BUTTON =
         "nextPage";
   public static final String ID_PERF_CHARTS_THUMBNAIL_FIRST_ROW_LINK_BUTTON =
         "row2_entityLink";
   public static final String ID_PERF_CHARTS_THUMBNAIL_SECOND_ROW_LINK_BUTTON =
         "row4_entityLink";

   // Performance Charts - Advanced
   public static final String ID_PERF_CHARTS_CUSTOM_TIME_RANGE_DATE_FROM =
         "customFromDate/dateField";
   public static final String ID_PERF_CHARTS_CUSTOM_TIME_RANGE_DATE_TO =
         "customToDate/dateField";
   public static final String ID_PERFORMANCE_CHART_ADVANCED_OPTIONS_SAVED_DROPDOWN =
         "savedChartSettingsSelect";
   public static final String ID_PERFORMANCE_CHART_ADVANCED_OPTIONS_SAVE_AS_BUTTON =
         "saveAsChartSettingsBtn";
   public static final String ID_PERFORMANCE_CHART_ADVANCED_OPTIONS_DELETE_BUTTON =
         "deleteSavedChartSettingsBtn";
   public static final String ID_PERFORMANCE_CHART_ADVANCED_OPTIONS_LOAD_STARTUP_CHECKBOX =
         "alwaysLoadChartSettingsChk";
   public static final String ID_CONFIRM_DIALOG_OK_BUTTON = "OkCancelDialog/label=OK";
   public static final String ID_SAVE_SETTINGS_DIALOG = "advPerfChartsSaveSettingsDlg";
   public static final String ID_SAVE_SETTINGS_DIALOG_TEXT_INPUT =
         ID_SAVE_SETTINGS_DIALOG + "/settingsNameInput";

   // IDs "Edit Distributed Port Group" wizard
   public static final String ID_DVPORTGROUP_NAME_FIELD =
         "_DvPortgroupGeneralPropertiesPage_TextInput1";

   // IDs for "Create Distributed Port Group" wizard pages

   // Wizard Page Ids
   public static final String ID_CREATE_PG_SELECT_NAME_AND_LOCATION_PAGE_ID =
         "step_nameAndLocationPage";
   public static final String ID_CREATE_PG_CONFIGURE_SETTINGS_PAGE_ID =
         "step_propertiesPage";
   public static final String ID_CREATE_PG_READY_TO_COMPLETE_PAGE_ID = "step_finalPage";

   public static final String ID_CREATE_PG_WIZARD = "Create a Distributed Port Group";
   public static final String ID_CREATE_PG_CONFIGURE_SETTINGS_PAGE = "propertiesPage";
   public static final String ID_CREATE_PG_READY_TO_COMPLETE_PAGE = "finalPage";
   public static final String ID_CREATE_PG_CUSTOMIZE_PAGE = "policiesPage";
   public static final String ID_CREATE_PG_SECURITY_PAGE = "securityPolicyPage";
   public static final String ID_CREATE_PG_EGRESS_PAGE = "egressPage";
   public static final String ID_CREATE_PG_INGRESS_PAGE = "ingressPage";

   // IDs for Configure settings page at "Create Distributed Port Group" wizard
   public static final String ID_CREATE_PG_PORTGROUP_NAME_FIELD = "inputLabel";
   public static final String ID_CREATE_PG_NUMBER_OF_PORTS_FIELD =
         ID_CREATE_PG_CONFIGURE_SETTINGS_PAGE + "/"
               + "_DvPortgroupGeneralPropertiesPage_NumericStepper1";
   public static final String ID_CREATE_PG_NUMBER_OF_PORTS_FIELD_INCREMENT_BUTTON =
         ID_CREATE_PG_NUMBER_OF_PORTS_FIELD + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_CREATE_PG_NUMBER_OF_PORTS_FIELD_DECREMENT_BUTTON =
         ID_CREATE_PG_NUMBER_OF_PORTS_FIELD + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_CREATE_PG_PORT_BINDING_DROPDOWN = "portBinding";
   public static final String ID_CREATE_PG_NETWORK_RESOURCE_POOL_DROPDOWN =
         ID_CREATE_PG_CONFIGURE_SETTINGS_PAGE + "/" + "poolSelector";
   public static final String ID_CREATE_PG_PORT_ALLOCATION_TYPES_DROPDOWN =
         "portAllocDrowDown";
   public static final String ID_CREATE_PG_VLAN_TYPE_DROPDOWN =
         ID_CREATE_PG_CONFIGURE_SETTINGS_PAGE + "/" + "vlanTypeSelector";
   public static final String ID_CREATE_PG_VLAN_ID =
         ID_CREATE_PG_CONFIGURE_SETTINGS_PAGE + "/"
               + "_DvPortVlanPolicyPage_NumericStepper1";
   public static final String ID_CREATE_PG_VLAN_ID_INCREMENT_BUTTON =
         ID_CREATE_PG_VLAN_ID + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_CREATE_PG_VLAN_ID_DECREMENT_BUTTON =
         ID_CREATE_PG_VLAN_ID + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_CREATE_PG_VLAN_TRUNK_RANGE = "vlanIdRanges";
   public static final String ID_CREATE_PG_CHECKBOX = "editAdvancedSelector";
   public static final String ID_CREATE_PG_TABLE_OF_CONTENTS_CUSTOMIZE =
         "_WizardTableOfContentsItemRenderer_Group2";
   public static final String ID_TITLE_LABLE = "titleLabel";
   public static final String ID_CREATE_PG_VLAN_LABEL = "_PropertiesPage_SettingsBlock1"
         + "/" + ID_TITLE_LABLE;
   public static final String ID_CREATE_PG_ADVANCED_LABEL =
         "_PropertiesPage_SettingsBlock2" + "/" + ID_TITLE_LABLE;
   public static final String ID_CREATE_PG_CONFIGURE_SETTINGS_LABEL =
         "pageHeaderTitleElement";
   public static final String ID_CREATE_PG_PRIVATE_VLAN_LABEL =
         "_DvPortVlanPolicyPage_LabelEx3";
   public static final String ID_CREATE_PG_PRIVATE_VLAN_POLICY_LABEL =
         "_DvPortVlanPolicyPage_Text1";
   // IDs for Security Policies page at "Create Distributed Port Group" wizard
   public static final String ID_CREATE_PG_PROMISCUOUS_MODE_DROPDOWN =
         ID_CREATE_PG_SECURITY_PAGE
               + "/"
               + "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList1/_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_CREATE_PG_MAC_ADDRESS_CHANGES_DROPDOWN =
         ID_CREATE_PG_SECURITY_PAGE
               + "/"
               + "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList2/_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_CREATE_PG_FORGED_TRANSMITS_DROPDOWN =
         ID_CREATE_PG_SECURITY_PAGE
               + "/"
               + "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList3/_OverridableBooleanDropDownList_BooleanDropDownList1";

   // IDs for Traffic Shaping page at "Create Distributed Port Group" wizard
   public static final String ID_CREATE_PG_INGRESS_SECTION =
         "_DvPortTrafficShapingPolicyListPage_DvPortTrafficShapingPolicyPage1";
   public static final String ID_CREATE_PG_INGRESS_TRAFFIC_SHAPING_STATUS_DROPDOWN =
         ID_CREATE_PG_INGRESS_PAGE + "/"
               + "_DvPortTrafficShapingPolicyPage_BooleanDropDownList1";
   public static final String ID_CREATE_PG_INGRESS_AVERAGE_BANDWIDTH_FIELD =
         ID_CREATE_PG_INGRESS_PAGE + "/" + "averageNumStepper";
   public static final String ID_CREATE_PG_INGRESS_AVERAGE_BANDWIDTH_FIELD_INCREMENT_BUTTON =
         ID_CREATE_PG_INGRESS_AVERAGE_BANDWIDTH_FIELD + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_CREATE_PG_INGRESS_AVERAGE_BANDWIDTH_FIELD_DECREMENT_BUTTON =
         ID_CREATE_PG_INGRESS_AVERAGE_BANDWIDTH_FIELD + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_CREATE_PG_INGRESS_PEAK_BANDWIDTH_FIELD =
         ID_CREATE_PG_INGRESS_PAGE + "/" + "peakNumStepper";
   public static final String ID_CREATE_PG_INGRESS_PEAK_BANDWIDTH_FIELD_INCREMENT_BUTTON =
         ID_CREATE_PG_INGRESS_PEAK_BANDWIDTH_FIELD + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_CREATE_PG_INGRESS_PEAK_BANDWIDTH_FIELD_DECREMENT_BUTTON =
         ID_CREATE_PG_INGRESS_PEAK_BANDWIDTH_FIELD + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_CREATE_PG_INGRESS_BURST_SIZE_FIELD =
         ID_CREATE_PG_INGRESS_PAGE + "/" + "burstNumStepper";
   public static final String ID_CREATE_PG_INGRESS_BURST_SIZE_FIELD_INCREMENT_BUTTON =
         ID_CREATE_PG_INGRESS_BURST_SIZE_FIELD + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_CREATE_PG_INGRESS_BURST_SIZE_FIELD_DECREMENT_BUTTON =
         ID_CREATE_PG_INGRESS_BURST_SIZE_FIELD + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_CREATE_PG_EGRESS_SECTION =
         "_DvPortTrafficShapingPolicyListPage_DvPortTrafficShapingPolicyPage2";
   public static final String ID_CREATE_PG_EGRESS_TRAFFIC_SHAPING_STATUS_DROPDOWN =
         ID_CREATE_PG_EGRESS_PAGE + "/"
               + "_DvPortTrafficShapingPolicyPage_BooleanDropDownList1";
   public static final String ID_CREATE_PG_EGRESS_AVERAGE_BANDWIDTH_FIELD =
         ID_CREATE_PG_EGRESS_PAGE + "/" + "averageNumStepper";
   public static final String ID_CREATE_PG_EGRESS_AVERAGE_BANDWIDTH_FIELD_INCREMENT_BUTTON =
         ID_CREATE_PG_EGRESS_AVERAGE_BANDWIDTH_FIELD + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_CREATE_PG_EGRESS_AVERAGE_BANDWIDTH_FIELD_DECREMENT_BUTTON =
         ID_CREATE_PG_EGRESS_AVERAGE_BANDWIDTH_FIELD + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_CREATE_PG_EGRESS_PEAK_BANDWIDTH_FIELD =
         ID_CREATE_PG_EGRESS_PAGE + "/" + "peakNumStepper";
   public static final String ID_CREATE_PG_EGRESS_PEAK_BANDWIDTH_FIELD_INCREMENT_BUTTON =
         ID_CREATE_PG_EGRESS_PEAK_BANDWIDTH_FIELD + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_CREATE_PG_EGRESS_PEAK_BANDWIDTH_FIELD_DECREMENT_BUTTON =
         ID_CREATE_PG_EGRESS_PEAK_BANDWIDTH_FIELD + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_CREATE_PG_EGRESS_BURST_SIZE_FIELD =
         ID_CREATE_PG_EGRESS_PAGE + "/" + "burstNumStepper";
   public static final String ID_CREATE_PG_EGRESS_BURST_SIZE_FIELD_INCREMENT_BUTTON =
         ID_CREATE_PG_EGRESS_BURST_SIZE_FIELD + "/" + ID_INCREMENT_BUTTON;
   public static final String ID_CREATE_PG_EGRESS_BURST_SIZE_FIELD_DECREMENT_BUTTON =
         ID_CREATE_PG_EGRESS_BURST_SIZE_FIELD + "/" + ID_DECREMENT_BUTTON;
   public static final String ID_CREATE_PG_TRAFFIC_SHAPING_LABEL_TEXT =
         "_DvPortTrafficShapingPolicyListPage_Text1";
   // IDs for Teaming and Failover page at "Create Distributed Port Group"
   // wizard
   public static final String ID_TEAMING_AND_FAILOVER_STACK_BLOCK =
         "_CustomizationPage_StackBlock3";
   public static final String ID_LOAD_BALANCING_DROP_DOWN = "loadBalancingSelector";
   public static final String ID_NETWORK_FAILURE_DETECTION_DROP_DOWN =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList1" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_NOTIFY_SWITCHES_DROP_DOWN =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList2" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_FAILBACK_DROP_DOWN =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList3" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";

   // IDs for Monitoring page at "Create Distributed Port Group" wizard
   public static final String ID_CREATE_PG_MONITORING_PAGE = "monitoringPolicyPage";
   public static final String ID_CREATE_PG_MONITORING_NET_FLOW_DROPDOWN =
         ID_CREATE_PG_MONITORING_PAGE
               + "/"
               + "_DvPortMonitoringPolicyPage_OverridableBooleanDropDownList1/_OverridableBooleanDropDownList_BooleanDropDownList1";
   // IDs for Miscellaneous at "Create Distributed Port Group" wizard
   public static final String ID_CREATE_PG_MISCELLANEOUS_PAGE = "miscPoliciesPage";
   public static final String ID_CREATE_PG_MISCELLANEOUS_BLOCK_ALL_PORTS_DROPDOWN =
         ID_CREATE_PG_MISCELLANEOUS_PAGE
               + "/"
               + "_DvPortMiscPoliciesPage_OverridableBooleanDropDownList1/_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_CREATE_PG_BLOCK_ALL_PORTS_TEXT =
         "_DvPortMiscPoliciesPage_Text1";
   // IDs for Additional properties at "Create Distributed Port Group" wizard
   public static final String ID_CREATE_PG_PORTGROUP_DESCRIPTION_FIELD =
         "descriptionText";
   public static final String ID_ADDITIONAL_PROPERTIES_STACK_BLOCK =
         "_CustomizationPage_StackBlock6";
   public static String ID_ACTION_BAR_TITLE = ID_NAV_BOX + "/contentGroup/title";
   public static final String ID_TRAFFIC_SHAPING_INGRESS_EDIT_PG =
         "_DvPortTrafficShapingPolicyListPage_StackBlock1";
   public static final String ID_TRAFFIC_SHAPING_EGRESS_EDIT_PG =
         "_DvPortTrafficShapingPolicyListPage_StackBlock2";
   public static final String ID_TRAFFIC_SHAPING_STATUS_DROPDOWN =
         "_DvPortTrafficShapingPolicyPage_BooleanDropDownList1";
   public static final String ID_TRAFFIC_SHAPING_AVERAGE =
         "_DvPortTrafficShapingPolicyPage_NumericStepper1";
   public static final String ID_TRAFFIC_SHAPING_PEAK =
         "_DvPortTrafficShapingPolicyPage_NumericStepper2";
   public static final String ID_TRAFFIC_SHAPING_BURST =
         "_DvPortTrafficShapingPolicyPage_NumericStepper3";
   public static final String ID_CREATE_PG_PORT_GROUP_CONFIGURE_AT_RESET =
         "_DvPortgroupAdvancedPropertiesPage_BooleanDropDownList1";
   public static final String ID_PORT_GROUP_OBJECTS_NEW_PORT_GROUP_BUTTON_ID =
         "vsphere.core.dvs.createPortgroupGlobal" + "/" + ID_BUTTON_BUTTON_NOSLASH;

   // Host->Manage->Settings->Networking|General
   public static final String ID_BUTTON_EDIT_GENERAL_SETTINGS =
         "btn_vsphere.core.host.network.editGeneralSettingsAction";
   public static final String ID_HOST_NETWORKING_GENERAL_IP_STATUS =
         "ipv6SupportedHeader";

   // Host->Manage->Settings->Networking|DNS and Routing
   public static final String ID_LABEL_HOST_NAME = "hostNameLabel";
   public static final String ID_LABEL_HOST_DOMAIN = "hostDomainLabel";
   public static final String ID_LABEL_DNS_METHOD = "dnsMethodLabel";
   public static final String ID_LABEL_VMK_NETWORK_ADAPTER = "vnicLabel";
   public static final String ID_LABEL_PREFERRED_DNS_SERVER = "preferredDnsServerLabel";
   public static final String ID_LABEL_ALTERNATE_DNS_SERVER = "alternateDnsServerLabel";
   public static final String ID_LABEL_SEARCH_DOMAINS = "searchDomainLabel";
   public static final String ID_LABEL_IPV4_VMKERNEL_GATEWAY =
         "ipv4VmkernelGatewayLabel";
   public static final String ID_LABEL_IPV6_VMKERNEL_GATEWAY =
         "ipv6VmkernelGatewayLabel";
   public static final String ID_BUTTON_EDIT_DNS_AND_ROUTING =
         "btn_vsphere.core.host.network.editDnsAndRoutingAction";

   public static final String ID_LABEL_DNS_AND_ROUTING_TITLE = "titleLabel";
   public static final String ID_PROPERTYGRIDROW_HOST_NAME = "hostNameRow";
   public static final String ID_PROPERTYGRIDROW_DOMAIN_NAME = "hostDomainRow";
   public static final String ID_PROPERTYGRIDROW_DNS_METHOD = "dnsMethodRow";
   public static final String ID_PROPERTYGRIDROW_VMK_ADAPTER = "vmkernelAdapterRow";
   public static final String ID_PROPERTYGRIDROW_PREFERRED_DNS_SERVER =
         "preferredDnsServerRow";
   public static final String ID_PROPERTYGRIDROW_ALTERNATE_DNS_SERVER =
         "alternateDnsServerRow";
   public static final String ID_PROPERTYGRIDROW_SEARCH_DOMAINS = "searchDomainRow";
   public static final String ID_LABEL_DEFAULT_GATEWAYS_TITLE = "defaultGatewayLabel";
   public static final String ID_PROPERTYGRIDROW_IPV4_VMKERNEL_GATEWAY =
         "ipv4VmkernelGatewayRow";
   public static final String ID_PROPERTYGRIDROW_IPV6_VMKERNEL_GATEWAY =
         "ipv6VmkernelGatewayRow";
   public static final String ID_PROPERTYGRIDROW_IPV4_CONSOLE_GATEWAY =
         "ipv4ConsoleGatewayRow";
   public static final String ID_PROPERTYGRIDROW_IPV6_CONSOLE_GATEWAY =
         "ipv6ConsoleGatewayRow";

   // 'Edit DNS and Routing Configuration' dialog
   public static final String ID_TEXTINPUT_HOST_NAME = "hostNameInput";
   public static final String ID_TEXTINPUT_HOST_DOMAIN = "hostDomainInput";
   public static final String ID_RADIO_BUTTON_DNS_METHOD_DYNAMIC =
         "_DnsSettingsPage_RadioButton1";
   public static final String ID_RADIO_BUTTON_DNS_METHOD_STATIC =
         "_DnsSettingsPage_RadioButton2";
   public static final String ID_SPARK_DROPDOWN_LIST_VMK_ADAPTER =
         "vmkernelAdaptersList";
   public static final String ID_TEXTINPUT_PREFERRED_DNS_SERVER =
         "preferredDnsServerInput";
   public static final String ID_TEXTINPUT_ALTERNATE_DNS_SERVER =
         "alternateDnsServerInput";
   public static final String ID_TEXTINPUT_SEARCH_DOMAINS = "searchDomainInput";
   public static final String ID_TEXTINPUT_IPV4_VMKERNEL_GATEWAY =
         "ipv4VmkernelGatewayInput";
   public static final String ID_SPARK_LABEL_DNS_CONFIGURATION =
         "step_editDnsSettingsView";
   public static final String ID_SPARK_LABEL_ROUTING = "step_gatewaysSettingsView";

   // Summary tab for port group
   public static final String ID_NETWORK_DETAILS_PORTLET =
         "vsphere.core.dvPortgroup.summary.detailsView.chrome";
   public static final String ID_VIRTUAL_MACHINES = "vms";
   public static final String ID_HOSTS = "hosts";
   public static final String ID_NETWORK_PROTOCOL_PROFILE = "ipPoolName";
   public static final String ID_NUMBER_OF_PORTS_STACK_BLOCK = "portsBlock";
   public static final String ID_NUMBER_OF_PORTS_ALL = "ports";
   public static final String ID_NUMBER_OF_PORTS_AVAILABLE = "freeCapacityLabel";
   public static final String ID_NUMBER_OF_PORTS_USED = "usageLabel";
   public static final String ID_NUMBER_OF_PORTS_TOTAL = "totalCapacityLabel";
   public static final String ID_PORT_BINDING_PORTGROUP_SUMMARY =
         "summary_portBindProp_valueLbl";
   public static final String ID_ASSOCIATED_SWITCH = "switchButton";

   // dvPortGroup Ports tab IDs
   public static final String ID_EDIT_PORT_SETTINGS_BUTTON = new StringBuffer(
         EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
         .append("editPortSettingsAction/button").toString();

   // Edit Port settings dialog IDs
   public static final String ID_OVERRIDE_CHECKBOX_EDIT_PORT =
         "_OverridableBooleanDropDownList_CheckBox1";

   public static final String ID_PORT_NAME_EDIT_PORT =
         "_DvPortPropertiesPage_TextInputEx1";
   public static final String ID_PORT_DESCRIPTION_EDIT_PORT =
         "_DvPortPropertiesPage_TextAreaEx1";

   public static final String ID_PROMISCUOUS_MODE_DROP_DOWN_EDIT_PORT =
         "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList1" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_MAC_ADDRESS_CHANGES_DROP_DOWN_EDIT_PORT =
         "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList2" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_FORGED_TRANSMITS_DROP_DOWN_EDIT_PORT =
         "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList3" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_OVERRIDE_PROMISCUOUS_MODE_CHECKBOX_EDIT_PORT =
         "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList1" + "/"
               + ID_OVERRIDE_CHECKBOX_EDIT_PORT;
   public static final String ID_OVERRIDE_MAC_ADDRESS_CHANGES_CHECKBOX_EDIT_PORT =
         "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList2" + "/"
               + ID_OVERRIDE_CHECKBOX_EDIT_PORT;
   public static final String ID_OVERRIDE_FORGED_TRANSMITS_CHECKBOX_EDIT_PORT =
         "_DvPortSecurityPolicyPage_OverridableBooleanDropDownList3" + "/"
               + ID_OVERRIDE_CHECKBOX_EDIT_PORT;

   public static final String ID_INGRESS_TRAFFIC_POLICY_PAGE_EDIT_PORT = "ingressPage";
   public static final String ID_EGRESS_TRAFFIC_POLICY_PAGE_EDIT_PORT = "egressPage";
   public static final String ID_INGRESS_TRAFFIC_AVERAGE_EDIT_PORT =
         ID_INGRESS_TRAFFIC_POLICY_PAGE_EDIT_PORT + "/" + "averageNumStepper";
   public static final String ID_INGRESS_TRAFFIC_PEAK_EDIT_PORT =
         ID_INGRESS_TRAFFIC_POLICY_PAGE_EDIT_PORT + "/" + "peakNumStepper";
   public static final String ID_INGRESS_TRAFFIC_BURST_SIZE_EDIT_PORT =
         ID_INGRESS_TRAFFIC_POLICY_PAGE_EDIT_PORT + "/" + "burstNumStepper";
   public static final String ID_EGRESS_TRAFFIC_AVERAGE_EDIT_PORT =
         ID_EGRESS_TRAFFIC_POLICY_PAGE_EDIT_PORT + "/" + "averageNumStepper";
   public static final String ID_EGRESS_TRAFFIC_PEAK_EDIT_PORT =
         ID_EGRESS_TRAFFIC_POLICY_PAGE_EDIT_PORT + "/" + "peakNumStepper";
   public static final String ID_EGRESS_TRAFFIC_BURST_SIZE_EDIT_PORT =
         ID_EGRESS_TRAFFIC_POLICY_PAGE_EDIT_PORT + "/" + "burstNumStepper";
   public static final String ID_INGRESS_TRAFFIC_DROP_DOWN_EDIT_PORT =
         ID_INGRESS_TRAFFIC_POLICY_PAGE_EDIT_PORT + "/"
               + "_DvPortTrafficShapingPolicyPage_BooleanDropDownList1";
   public static final String ID_EGRESS_TRAFFIC_DROP_DOWN_EDIT_PORT =
         ID_EGRESS_TRAFFIC_POLICY_PAGE_EDIT_PORT + "/"
               + "_DvPortTrafficShapingPolicyPage_BooleanDropDownList1";


   public static final String ID_VLAN_COMBOBOX_EDIT_PORT = "vlanTypeSelector";
   public static final String ID_VLANID_NUMSTEPPER_EDIT_PORT =
         "_DvPortVlanPolicyPage_NumericStepper1";
   public static final String ID_VLAN_TRUNK_RANGE_INPUT_EDIT_PORT = "vlanIdRanges";
   public static final String VLAN_TRUNK_RANGE_ERROR_MESSAGE_PROPERTY = "errorString";

   public static final String ID_OVERRIDE_LOAD_BALANCING_CHECKBOX_EDIT_PORT =
         "_DvPortFailoverPolicyPage_HGroup1" + "/"
               + "_DvPortFailoverPolicyPage_CheckBox1";
   public static final String ID_OVERRIDE_NETWORK_FAILURE_CHECKBOX_EDIT_PORT =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList1" + "/"
               + ID_OVERRIDE_CHECKBOX_EDIT_PORT;
   public static final String ID_OVERRIDE_NOTIFY_CHECKBOX_EDIT_PORT =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList2" + "/"
               + ID_OVERRIDE_CHECKBOX_EDIT_PORT;
   public static final String ID_OVERRIDE_FAILBACK_CHECKBOX_EDIT_PORT =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList3" + "/"
               + ID_OVERRIDE_CHECKBOX_EDIT_PORT;
   public static final String ID_UPLINK_FAILOVER_STACK_BLOCK_EDIT_PORT =
         "_DvPortFailoverPolicyPage_StackBlock5";
   public static final String ID_UPLINK_FAILOVER_OVERRIDE_CHECKBOX_EDIT_PORT =
         "_DvPortFailoverPolicyPage_CheckBox2";
   public static final String ID_UPLINK_FAILOVER_MOVE_UP_BUTTON =
         "_DvPortFailoverPolicyPage_Button1";

   public static final String ID_LOAD_BALANCING_DROP_DOWN_EDIT_PORT =
         "_DvPortFailoverPolicyPage_HGroup1" + "/" + "loadBalancingSelector";
   public static final String ID_NETWORK_FAILURE_DROP_DOWN_EDIT_PORT =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList1" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_FAILBACK_DROP_DOWN_EDIT_PORT =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList3" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";
   public static final String ID_NOTIFY_SWITCHES_DROP_DOWN_EDIT_PORT =
         "_DvPortFailoverPolicyPage_OverridableBooleanDropDownList2" + "/"
               + "_OverridableBooleanDropDownList_BooleanDropDownList1";

   public static final String ID_OVERRIDE_CHECKBOX_NETWORK_RP_EDIT_PORT =
         "_DvPortPropertiesPage_DvPortNetIormPolicyPage1" + "/"
               + "_DvPortNetIormPolicyPage_CheckBox1";
   public static final String ID_NETWORK_RESOURCE_POOL_EDIT_PORT_DROPDOWN =
         "portProperties" + "/" + "poolSelector";
   public static final String ID_OVERRIDE_CHECKBOX_BLOCK_PORT_EDIT_PORT = "misc" + "/"
         + ID_OVERRIDE_CHECKBOX_EDIT_PORT;
   public static final String ID_BLOCK_PORT_DROP_DOWN_EDIT_PORT = "misc" + "/"
         + "_OverridableBooleanDropDownList_BooleanDropDownList1";

   public static final String ID_OVERRIDE_INGRESS_CHECK_BOX = "ingressPage" + "/"
         + "overrideCbx";
   public static final String ID_OVERRIDE_CHECKBOX_EGRESS_TRAFFIC_EDIT_PORT =
         "egressPage" + "/" + "overrideCbx";

   public static final String ID_OVERRIDE_CHECKBOX_SECURITY_PROMISCUOUS_EDIT_PORT =
         "security" + "/" + ID_OVERRIDE_CHECKBOX_EDIT_PORT;
   public static final String ID_OVERRIDE_CHECKBOX_SECURITY_MAC_ADD_EDIT_PORT =
         "security" + "/" + ID_OVERRIDE_CHECKBOX_EDIT_PORT;
   public static final String ID_OVERRIDE_CHECKBOX_SECURITY_FORGED_TRANDSMIT_EDIT_PORT =
         "security" + "/" + ID_OVERRIDE_CHECKBOX_EDIT_PORT;

   public static final String ID_OVERRIDE_VLAN_CHECK_BOX = "vlan" + "/" + "overrideCkb";
   public static final String ID_OVERRIDE_MONITORING_CHECKBOX_EDIT_PORT = "monitoring"
         + "/" + ID_OVERRIDE_CHECKBOX_EDIT_PORT;

   // Migrate Virtual Machine Networking wizard
   public static final String ID_ADVANCE_DATAGRID_MIGRATE_VM_NETWORKING_SOURCE_NETWORK =
         "sourceNetworks";
   public static final String ID_ADVANCE_DATAGRID_MIGRATE_VM_NETWORKING_DEST_NETWORK =
         "destinationNetworks";
   public static final String ID_ADVANCE_DATAGRID_MIGRATE_VM_NETWORKING_VMS_TO_MIGRATE =
         "vmList";

   // SSO IDs
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_BUTTON =
         "sso.actions.addIds/button";
   public static final String ID_SSO_CONFIGURATION_DELETE_IDENTITYSOURCE_BUTTON =
         "sso.actions.deleteIds/button";
   public static final String ID_SSO_CONFIGURATION_EDIT_IDENTITYSOURCE_BUTTON =
         "sso.actions.editIds/button";
   public static final String ID_SSO_CONFIGURATION_ADDDEFAULTDOMAIN_IDENTITYSOURCE_BUTTON =
         "automationName=Add to Default Domains";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_AD_RADIO =
         "adRadioButton";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_OPENLDAP_RADIO =
         "openLdapRadioButton";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_LOCALOS_RADIO =
         "localRadioButton";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_NAME_TEXT =
         "domainFriendlyNameTextInput";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_PRIMARYLDAP_TEXT =
         "primaryUrlTextInput";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_SECONDARYLDAP_TEXT =
         "secondaryUrlTextInput";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_USERCN_TEXT =
         "userCnTextInput";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_DOMAINNAME_TEXT =
         "domainNameTextInput";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_DOMAINALIAS_TEXT =
         "domainAliasTextInput";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_GROUPCN_TEXT =
         "groupDnTextInput";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_AUTHENTICATIONTYPE_TEXT =
         "authTypeComboBox";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_USERNAME_TEXT =
         "usernameTextInput";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_PASSWORD_TEXT =
         "passwordTextInput";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_TESTCONNECTION_TEXT =
         "probeConnectionButton";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_ALERT_DIALOG =
         "AlertDialog";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_ALERT_DIALOG_OK_BUTTON =
         "automationName=OK";
   public static final String ID_SSO_CONFIGURATION_ADD_IDENTITYSOURCE_DIALOG_OK_BUTTON =
         "okButton";
   public static final String ID_SSO_CONFIGURATION_IDENTITYSOURCES_GRID = "idsTable";
   public static final String ID_SSO_CONFIGURATION_DELETE_IDENTITYSOURCE_CONFIRMATION_DIALOG =
         "confirmationDialog";
   public static final String ID_SSO_CONFIGURATION_DELETE_USER_CONFIRMATION_DIALOG =
         "confirmationDialog";
   public static final String ID_SSO_CONFIGURATION_DELETE_IDENTITYSOURCE_CONFIRMATION_DIALOG_YESBUTTON =
         "automationName=Yes";
   public static final String ID_SSO_CONFIGURATION_DELETE_USER_CONFIRMATION_DIALOG_YESBUTTON =
         "automationName=Yes";
   public static final String ID_SSO_USERSANDGROUPS_LOCALUSERS_GRID = "usersTable";
   public static final String ID_SSO_USERSANDGROUPS_APPLICATIONUSERS_GRID =
         "solutionUsersTable";
   public static final String ID_SSO_CONFIGURATION_ADDUSER_DIALOG_OK_BUTTON = "okButton";
   public static final String ID_SSO_CONFIGURATION_ADD_USER_BUTTON =
         "sso.actions.createUser/button";
   public static final String ID_SSO_CONFIGURATION_DELETE_USER_BUTTON =
         "sso.actions.deleteUser/button";
   public static final String ID_SSO_CONFIGURATION_DELETE_APPUSER_BUTTON =
         "sso.actions.deleteAppUser/button";
   public static final String ID_SSO_CONFIGURATION_GUESTUSER_RADIO = "guestUserButton";
   public static final String ID_SSO_CONFIGURATION_ADDUSER_USERNAME_TEXT = "username";
   public static final String ID_SSO_CONFIGURATION_ADDUSER_PASSWORD_TEXT = "password";
   public static final String ID_SSO_CONFIGURATION_ADDUSER_CONFIRMPASSWORD_TEXT =
         "passwordConfirm";
   public static final String ID_SSO_CONFIGURATION_ADDUSER_FIRSTNAME_TEXT = "firstName";
   public static final String ID_SSO_CONFIGURATION_ADDUSER_LASTNAME_TEXT = "lastName";
   public static final String ID_SSO_CONFIGURATION_ADDUSER_EMAIL_TEXT = "email";
   public static final String ID_SSO_CONFIGURATION_ADDUSER_DIALOG_CANCEL_BUTTON =
         "cancelButton";
   public static final String ID_SSO_CONFIGURATION_LOCK_USER_BUTTON =
         "sso.actions.lockUser/button";
   public static final String ID_SSO_CONFIGURATION_UNLOCK_USER_BUTTON =
         "sso.actions.unlockUser/button";
   public static final String ID_SSO_CONFIGURATION_DISABLE_USER_BUTTON =
         "sso.actions.disableUser/button";
   public static final String ID_SSO_CONFIGURATION_ENABLE_USER_BUTTON =
         "sso.actions.enableUser/button";
   public static final String ID_SSO_CONFIGURATION_EDIT_USER_BUTTON =
         "sso.actions.editUser/button";
   public static final String ID_SSO_CONFIGURATION_ADDUSER_DESCRIPTION_TEXT =
         "description";
   public static final String ID_SSO_USERSANDGROUPS_ADD_GROUP_BUTTON =
         "sso.actions.createGroup/button";
   public static final String ID_SSO_USERSANDGROUPS_EDIT_GROUP_BUTTON =
         "sso.actions.editGroup/button";
   public static final String ID_SSO_USERSANDGROUPS_ADDGROUP_GROUPNAME_TEXT =
         "groupname";
   public static final String ID_SSO_USERSANDGROUPS_ADDGROUP_GROUPDESCRIPTION_TEXT =
         "description";
   public static final String ID_SSO_USERSANDGROUPS_ADD_TIWODIALOG = "tiwoDialog";
   public static final String ID_SSO_USERSANDGROUPS_ADDGROUP_DIALOG = "tiwoDialog";
   public static final String ID_SSO_USERSANDGROUPS_ADDGROUP_DIALOG_OK_BUTTON =
         "okButton";
   public static final String ID_SSO_USERSANDGROUPS_DELETE_GROUP_BUTTON =
         "sso.actions.deleteGroup/button";
   public static final String ID_SSO_USERSANDGROUPS_DELETE_GROUP_PRINCIPAL_CONFIRMATION_DIALOG =
         "confirmationDialog";
   public static final String ID_SSO_USERSANDGROUPS_DELETE_GROUP_PRINCIPAL_CONFIRMATION_DIALOG_YESBUTTON =
         "automationName=Yes";
   public static final String ID_SSO_USERSANDGROUPS_DELETE_PRINCIPAL_BUTTON =
         "sso.actions.removePrincipalFromGroup/button";
   public static final String ID_BUTTON_SSO_USERSANDGROUPS_ADD_PRINCIPAL_BUTTON =
         "sso.actions.addPrincipalToGroup/button";
   public static final String ID_RADIOBUTTON_SSO_CONFIGURATION_REGULARUSER =
         "regularUserButton";
   public static final String ID_RADIOBUTTON_SSO_CONFIGURATION_ADMINISTRATORUSER =
         "adminUserButton";
   public static final String ID_ADVDATAGRID_SSO_USERSANDGROUPS_LOCALGROUPS =
         "groupsTable";
   public static final String ID_ADVDATAGRID_SSO_USERSANDGROUPS_USERSINGROUPS_GRID =
         "usersInGroupTable";
   public static final String ID_COMBOBOX_SSO_USERSANDGROUPS_ADDPRINCIPAL =
         "identitySources";
   public static final String ID_BUTTON_SSO_CONFIGURATION_USERACTION_CONFIRMATION_DIALOG_YES =
         "automationName=Yes";
   public static final String ID_DIALOG_SSO_CONFIGURATION_USERACTION_CONFIRMATION_DIALOG =
         "confirmationDialog";
   public static final String ID_LABEL_SSO_ADDUSER_EMAIL_ERROR = "_errorLabel";
   public static final String ID_LABEL_SSO_LOCKEDUSER_LOGINPAGE_MESSAGE =
         "_ErrorComponent_Text1";
   public static final String ID_TEXTINPUT_SSO_FILTER_INPUTBOX =
         "filterControl/textInput";
   public static final String ID_BUTTON_SSO_FILTER = "filterControl/searchButton";
   public static final String ID_LABEL_SSO_ADDIDS_VALIDATION_ERROR_TITLE = "_titleLabel";
   public static final String ID_LABEL_SSO_ADDIDS_VALIDATION_ERROR_TEXT = "_errorLabel";
   public static final String ID_LABEL_SSO_ADDIDS_ERROR_TEXT = "_message";
   public static final String ID_DIALOG_SSO_ADDIDS_TIWODIALOG = "tiwoDialog";
   public static final String ID_DIALOG_SSO_TIWODIALOG = "tiwoDialog";
   public static final String ID_BUTTON_SSO_CONFIGURATION_TIWODIALOG_OK = "okButton";
   public static final String ID_BUTTON_SSO_CONFIGURATION_TIWODIALOG_CANCEL =
         "cancelButton";
   public static final String ID_BUTTONBAR_SSO_CONFIGURATION_POLICIES =
         "sso.admin.policies.extension/toggleButtonBar";
   public static final String ID_BUTTON_SSO_CONFIGURATION_EDIT_PASSWORDPOLICY =
         "editPasswordPolicyBtn";
   public static final String ID_BUTTON_SSO_CONFIGURATION_EDIT_LOCKOUTPOLICY =
         "editLockoutPolicyBtn";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_PASSWORDPOLICYDESCRIPTION =
         "description";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_MAXIMUMLIFETIME =
         "passwordLifetimeDays";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_RESTRICTEDREUSE =
         "prohibitedPreviousPasswordsCount";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_MAXIMUMLENGTH = "maxLength";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_MINIMUMLENGTH = "minLength";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_SPECIALCHARACTERS =
         "minSpecialCharCount";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_ALPHABETICCHARACTERS =
         "minAlphabeticCount";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_UPPERCASECHARACTERS =
         "minUppercaseCount";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_LOWERCASECHARACTERS =
         "minLowercaseCount";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_NUMERICCHARACTERS =
         "minNumericCount";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_MINIDENTICALADJECENTCHARACTERS =
         "minIdenticalAdjacentCharacters";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_MAXIDENTICALADJECENTCHARACTERS =
         "maxIdenticalAdjacentCharacters";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_LOCKOUTPOLICYDESCRIPTION =
         "description";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_FAILEDATTEMPTS =
         "maxFailedAttempts";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_TIMEINTERVAL =
         "failedAttemptIntervalSec";
   public static final String ID_TEXTINPUT_SSO_CONFIGURATION_UNLOCKTIME =
         "autoUnlockIntervalSec";
   public static final String ID_LABEL_SSO_CONFIGURATION_LOCKOUTPOLICYDESCRIPTION =
         "_SsoLockoutPolicyView_Label3";
   public static final String ID_LABEL_SSO_CONFIGURATION_FAILEDATTEMPTS =
         "_SsoLockoutPolicyView_Label6";
   public static final String ID_LABEL_SSO_CONFIGURATION_TIMEINTERVAL =
         "_SsoLockoutPolicyView_Label8";
   public static final String ID_LABEL_SSO_CONFIGURATION_UNLOCKTIME =
         "_SsoLockoutPolicyView_Label10";
   public static final String ID_LABEL_SSO_CONFIGURATION_PASSWORDPOLICYDESCRIPTION =
         "_SsoPasswordPoliciesView_Label3";
   public static final String ID_LABEL_SSO_CONFIGURATION_MAXIMUMLIFETIME =
         "_SsoPasswordPoliciesView_Label6";
   public static final String ID_LABEL_SSO_CONFIGURATION_RESTRICTEDREUSE =
         "_SsoPasswordPoliciesView_Label8";
   public static final String ID_LABEL_SSO_CONFIGURATION_MAXIMUMLENGTH =
         "_SsoPasswordPoliciesView_Label11";
   public static final String ID_LABEL_SSO_CONFIGURATION_MINIMUMLENGTH =
         "_SsoPasswordPoliciesView_Label13";
   public static final String ID_LABEL_SSO_CONFIGURATION_SPECIALCHARACTERS =
         "_SsoPasswordPoliciesView_Label16";
   public static final String ID_LABEL_SSO_CONFIGURATION_ALPHABETICCHARACTERS =
         "_SsoPasswordPoliciesView_Label15";
   public static final String ID_LABEL_SSO_CONFIGURATION_UPPERCASECHARACTERS =
         "_SsoPasswordPoliciesView_Label17";
   public static final String ID_LABEL_SSO_CONFIGURATION_LOWERCASECHARACTERS =
         "_SsoPasswordPoliciesView_Label18";
   public static final String ID_LABEL_SSO_CONFIGURATION_NUMERICCHARACTERS =
         "_SsoPasswordPoliciesView_Label19";
   public static final String ID_LABEL_SSO_CONFIGURATION_IDENTICALADJECENTCHARACTERS =
         "_SsoPasswordPoliciesView_Label20";
   public static final String ID_LABEL_ERROR_BANNER = "noticeText";
   public static final String ID_CHECKBOX_SSO_SSPI_LOGIN = "sspiCheckBox";


   public static final String ID_HOST_BOOT_OPTIONS_BUTTON = "btnBootDevice";

   // VC Configuration
   public static final String ID_VC_CONFIG_VCNAME =
         "Folder:vcServerSettings.runtimeSettings.VirtualCenter.InstanceName";
   public static final String ID_VC_CONFIG_VCMANAGEDIP =
         "Folder:vcServerSettings.runtimeSettings.VirtualCenter.ManagedIP";
   public static final String ID_VC_CONFIG_VCUNIQUEID =
         "Folder:vcServerSettings.runtimeSettings.instance.id";
   public static final String ID_VC_CONFIG_USER_DIRECTORY_TIMEOUT = "adTimeoutLabel";
   public static final String ID_VC_CONFIG_QUERY_LIMIT = "adQueryLimitLabel";
   public static final String ID_VC_CONFIG_VALIDATION_PERIOD =
         "Folder:vcServerSettings.activeDirectorySettings.ads.checkIntervalEnabled.Text";
   public static final String ID_VC_CONFIG_MAIL_SERVER =
         "Folder:vcServerSettings.mailSettings.mail.smtp.server";
   public static final String ID_VC_CONFIG_MAIL_SENDER =
         "Folder:vcServerSettings.mailSettings.mail.sender";
   public static final String ID_VC_CONFIG_SNMP_PRIMARY_RECIEVER_NAME =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.1.name";
   public static final String ID_VC_CONFIG_SNMP_PRIMARY_RECIEVER_ENABLED =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.1.enabled";
   public static final String ID_VC_CONFIG_SNMP_PRIMARY_RECIEVER_PORT =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.1.port";
   public static final String ID_VC_CONFIG_SNMP_PRIMARY_RECIEVER_COMMUNITY =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.1.community";
   public static final String ID_VC_CONFIG_SNMP_SECONDARY_RECIEVER_NAME =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.2.name";
   public static final String ID_VC_CONFIG_SNMP_SECONDARY_RECIEVER_ENABLED =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.2.enabled";
   public static final String ID_VC_CONFIG_SNMP_SECONDARY_RECIEVER_PORT =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.2.port";
   public static final String ID_VC_CONFIG_SNMP_SECONDARY_RECIEVER_COMMUNITY =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.2.community";
   public static final String ID_VC_CONFIG_SNMP_THIRD_RECIEVER_NAME =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.3.name";
   public static final String ID_VC_CONFIG_SNMP_THIRD_RECIEVER_ENABLED =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.3.enabled";
   public static final String ID_VC_CONFIG_SNMP_THIRD_RECIEVER_PORT =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.3.port";
   public static final String ID_VC_CONFIG_SNMP_THIRD_RECIEVER_COMMUNITY =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.3.community";
   public static final String ID_VC_CONFIG_SNMP_FOURTH_RECIEVER_NAME =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.4.name";
   public static final String ID_VC_CONFIG_SNMP_FOURTH_RECIEVER_ENABLED =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.4.enabled";
   public static final String ID_VC_CONFIG_SNMP_FOURTH_RECIEVER_PORT =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.4.port";
   public static final String ID_VC_CONFIG_SNMP_FOURTH_RECIEVER_COMMUNITY =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.4.community";
   public static final String ID_VC_CONFIG_HTTP_PORT = "httpPortLabel";
   public static final String ID_VC_CONFIG_HTTPS_PORT = "httpsPortLabel";
   public static final String ID_VC_CONFIG_NORMAL = "normalTimeoutLabel";
   public static final String ID_VC_CONFIG_LONG = "longTimeoutLabel";
   public static final String ID_VC_CONFIG_LOGLEVEL = "logLevelLabel";
   public static final String ID_VC_CONFIG_MAXIMUM_CONNECTIONS =
         "Folder:vcServerSettings.databaseSettings.VirtualCenter.MaxDBConnection";
   public static final String ID_VC_CONFIG_TASK_RETENTION = "dbTaskRetentionLabel";
   public static final String ID_VC_CONFIG_EVENT_RETENTION = "dbEventRetentionLabel";
   public static final String ID_VC_CONFIG_SSL_SETTINGS = "sslSettingsHeaderText";
   public static final String ID_VC_CONFIG_EDIT_STATISTICS_DATAGRID =
         "statisticsIntervals";
   public static final String ID_VC_CONFIG_EDIT_VCNAME =
         "Folder:vcServerSettings.runtimeSettings.VirtualCenter.InstanceName";
   public static final String ID_VC_CONFIG_EDIT_VCMANAGEDIP =
         "Folder:vcServerSettings.runtimeSettings.VirtualCenter.ManagedIP";
   public static final String ID_VC_CONFIG_EDIT_VCUNIQUEID =
         "Folder:vcServerSettings.runtimeSettings.instance.id";
   public static final String ID_VC_CONFIG_EDIT_USER_DIRECTORY_TIMEOUT =
         "Folder:vcServerSettings.activeDirectorySettings.ads.timeout";
   public static final String ID_VC_CONFIG_EDIT_QUERY_LIMIT =
         "Folder:vcServerSettings.activeDirectorySettings.ads.maxFetch";
   public static final String ID_VC_CONFIG_EDIT_VALIDATION_PERIOD =
         "Folder:vcServerSettings.activeDirectorySettings.ads.checkIntervalEnabled";
   public static final String ID_VC_CONFIG_EDIT_MAIL_SERVER =
         "Folder:vcServerSettings.mailSettings.mail.smtp.server";
   public static final String ID_VC_CONFIG_EDIT_MAIL_SENDER =
         "Folder:vcServerSettings.mailSettings.mail.sender";
   public static final String ID_VC_CONFIG_EDIT_HTTP_PORT =
         "Folder:vcServerSettings.portsSettings.WebService.Ports.http";
   public static final String ID_VC_CONFIG_EDIT_HTTPS_PORT =
         "Folder:vcServerSettings.portsSettings.WebService.Ports.https";
   public static final String ID_VC_CONFIG_EDIT_NORMAL =
         "Folder:vcServerSettings.timeoutSettings.client.timeout.normal";
   public static final String ID_VC_CONFIG_EDIT_LONG =
         "Folder:vcServerSettings.timeoutSettings.client.timeout.long";
   public static final String ID_VC_CONFIG_EDIT_LOGLEVEL =
         "Folder:vcServerSettings.logSettings.log.level";
   public static final String ID_VC_CONFIG_EDIT_MAXIMUM_CONNECTIONS =
         "Folder:vcServerSettings.databaseSettings.VirtualCenter.MaxDBConnection";
   public static final String ID_VC_CONFIG_EDIT_TASK_RETENTION =
         "Folder:vcServerSettings.databaseSettings.task.maxAge";
   public static final String ID_VC_CONFIG_EDIT_EVENT_RETENTION =
         "Folder:vcServerSettings.databaseSettings.event.maxAge";
   public static final String ID_VC_CONFIG_EDIT_SSL_SETTINGS = "verifySslCertificates";
   public static final String ID_VC_CONFIG_STATISTICS_SECTION =
         "titleLabel.statisticsSettingsStackBlock";
   public static final String ID_VC_CONFIG_RUNTIME_SETTINGS_SECTION =
         "titleLabel.Folder:vcServerSettings.runtimeSettingsTitle";
   public static final String ID_VC_CONFIG_USER_DIRECTORY_SECTION =
         "titleLabel.Folder:vcServerSettings.activeDirectorySettingsTitle";
   public static final String ID_VC_CONFIG_DATABASE_SECTION =
         "titleLabel.Folder:vcServerSettings.databaseSettingsTitle";
   public static final String ID_VC_CONFIG_MAIL_SECTION =
         "titleLabel.Folder:vcServerSettings.mailSettingsTitle";
   public static final String ID_VC_CONFIG_SNMP_RECIEVERS_SECTION =
         "titleLabel.Folder:vcServerSettings.snmpSettingsTitle";
   public static final String ID_VC_CONFIG_EDIT_STATISTICS_SECTION =
         "step_vcStatisticsConfigPage";
   public static final String ID_VC_CONFIG_EDIT_RUNTIME_SETTINGS_SECTION =
         "step_vcRuntimeConfigPage";
   public static final String ID_VC_CONFIG_EDIT_USER_DIRECTORY_SECTION =
         "step_vcActiveDirectoryConfigPage";
   public static final String ID_VC_CONFIG_EDIT_DATABASE_SECTION =
         "step_vcDatabaseConfigPage";
   public static final String ID_VC_CONFIG_EDIT_MAIL_SECTION = "step_vcMailConfigPage";
   public static final String ID_VC_CONFIG_EDIT_PORTS_SECTION = "step_vcPortsConfigPage";
   public static final String ID_VC_CONFIG_EDIT_TIMEOUT_SETTINGS_SECTION =
         "step_vcTimeoutConfigPage";
   public static final String ID_VC_CONFIG_EDIT_LOGGING_SETTINGS_SECTION =
         "step_vcLogConfigPage";
   public static final String ID_VC_CONFIG_EDIT_SSL_SETTINGS_SECTION =
         "step_vcSslConfigPage";
   public static final String ID_VC_CONFIG_EDIT_SNMP_RECIEVERS_SECTION =
         "step_vcSnmpConfigPage";
   public static final String ID_VC_CONFIG_EDIT_BUTTON =
         "btn_vsphere.core.folder.editVcSettingsAction";
   public static final String ID_VCCONFIG_SETTINGS_SPARK_LIST = "tocListOther";
   public static final String ID_VCCONFIG_STATISTICS_HEADER = "statisticsHeaderText";
   public static final String ID_VCCONFIG_STATISTICS_ESTIMATED_DBSIZE =
         "estimatedDbSize";
   public static final String ID_VCCONFIG_STATISTICS_NO_OF_HOSTS =
         "estimatedNumberOfHosts";
   public static final String ID_VCCONFIG_STATISTICS_NO_OF_VMS = "estimatedNumberOfVms";
   public static final String ID_VCCONFIG_STATISTICS_ESTIMATED_DBSIZE_EXPANDED =
         "estimatedDbSpaceText";
   public static final String ID_VCCONFIG_ADV_SETTINGS_ADVDATAGRID =
         "Folder:vcServerSettings.advancedSettings";
   public static final String ID_VC_CONFIG_ADV_SETTINGS_EDIT_BUTTON =
         "btn_vsphere.core.folder.editVcAdvancedSettingsAction";
   public static final String ID_VC_CONFIG_ADV_SETTINGS_CONFIG_KEY = "settingKey";
   public static final String ID_VC_CONFIG_ADV_SETTINGS_CONFIG_VALUE = "settingValue";
   public static final String ID_VC_CONFIG_ADV_SETTINGS_ADD_BUTTON = "addSettingBtn";

   // Constants for Alarms
   public static final String ID_ALARMS_ISSUES_TAB_ADVGRID = "alarmList";
   public static final String ID_ALARMS_LIST = "alarmDefTable";
   public static final String ID_DATASTORE_VM_LIST = "datastoreVmList";
   public static final String ID_BUTTON_NEW_ALARM_WIZARD_CONTAINER_TRIGGER_ADD =
         "alarmTriggersPage/toolbarControl/toolContainer/add";
   public static final String ID_BUTTON_NEW_ALARM_WIZARD_CONTAINER__TRIGGER_DELETE =
         "alarmTriggersPage/toolbarControl/toolContainer/delete";
   public static final String ID_BUTTON_NEW_ALARM_WIZARD_CONTAINER_ACTION_ADD =
         "alarmActionsPage/toolbarControl/toolContainer/add";
   public static final String ID_BUTTON_NEW_ALARM_WIZARD_CONTAINER_ACTION_DELETE =
         "alarmActionsPage/toolbarControl/toolContainer/delete";
   public static final String ID_BUTTON_NEW_ALARM_WIZARD_CONTAINER_TRIGGER_EVENT_ADD =
         "alarmTriggersPage/eventConditions/toolbarControl/toolContainer/add";
   public static final String ID_CONDITION_GRID = "eventConditions/conditionsGrid";
   public static final String ID_COMBOBOX_ARGUMENT = "argumentCombo";
   public static final String ID_COMBOBOX_OPERATOR = "operatorCombo";
   public static final String ID_TEXTFIELD_CONDITION = "argumentCondition";
   public static final String ID_CONFIGURATION_TEXT = "configurationText";
   public static final String ID_BUTTON_NEW_ALARM_WIZARD_CONTAINER_TRIGGER_EVENT_DELETE =
         "alarmTriggersPage/eventConditions/toolbarControl/toolContainer/delete";
   public static final String ID_BUTTON_ALARMS_LIST_ADD =
         "alarmDefTable/dataGridToolbar/toolContainer/add";
   public static final String ID_BUTTON_ALARMS_LIST_DELETE =
         "alarmDefTable/dataGridToolbar/toolContainer/delete";
   public static final String ID_ALARM_DEFINITIONS_TAB = "Alarm Definitions";
   public static final String ID_LABEL_ALARM_NAME = "alarmName";
   public static final String ID_LABEL_ALARM_DESCRIPTION = "description";
   public static final String ID_COMBOBOX_ALARM_TYPE = "monitorType";
   public static final String ID_RADIOBUTTON_ALARMTYPE_STATE = "stateMonitor";
   public static final String ID_RADIOBUTTON_ALARMTYPE_EVENT = "eventMonitor";
   public static final String ID_CHECKBOX_ALARM_ENABLE = "enabledCheckbox";
   public static final String ID_DATAGRID_ALARM_TRIGGERS_PAGE =
         "_TriggersPage_VGroup2/grid";
   public static final String ID_SPARKLIST_EVENT_TRIGGERS = "eventTriggers";
   public static final String ID_SPARKLIST_EVENT_STATUS = "eventStatus";
   public static final String ID_SPARKLIST_STATE_TRIGGERS =
         "alarmTriggersPage/grid/className=TriggerTypeEditor";
   public static final String ID_SPARKLIST_ALARM_ANYALL_DROPDOWN =
         "alarmTriggersPage/anyAllDropDown";
   public static final String ID_SPARKLIST_STATE_TRIGGERS_CONDITION =
         "alarmTriggersPage/stateConditions";
   public static final String ID_SPARKLIST_STATE_TRIGGERS_CONDITION_VALUES =
         "alarmTriggersPage/stateValues";
   public static final String ID_DATAGRID_ALARM_ACTIONS_PAGE = "alarmActionsPage/grid";
   public static final String ID_SPARKLIST_ALARM_ACTIONS =
         "alarmActionsPage/actionWrappersDropDown";
   public static final String ID_SPARKLIST_ALARM_ACTION_REPITITIONS =
         "alarmActionsPage/repetition";
   public static final String ID_BUTTON_MONITOR_TAB_ALARMS_VIEW =
         "vsphere.opsmgmt.alarms.alarmsView.button";
   public static final String ID_SPARKLIST_STATE_TRIGGERS_METRIC_CONDITIONS =
         "alarmTriggersPage/metricConditions";
   public static final String ID_SPARKLIST_STATE_TRIGGERS_METRIC_VALUES = "metricValues";
   public static final String ID_SPARKLIST_STATE_TRIGGERS_METRIC_DURATION =
         "metricDuration";
   public static final String ID_LABEL_NEW_ALARM_WIZARD_ERROR_STRING = "_errorLabel";
   public static final String ID_TEXTINPUT_STATE_TRIGGERS_METRIC_VALUES =
         "metricValues/automationName=textInput";
   public static final String ID_NEW_ALARM_DEFINITION_WIZARD_TITLE =
         "New Alarm Definition";
   public static final String ID_TRIGGER_HOST_CONNECTION_STATE = "Host Connection State";
   public static final String ID_CONFIGURATION_MIGRATEVM_WIZARD =
         "automationName=Migrate Action Wizard";
   public static final String ID_BUTTON_CONFIG = "configButton";
   public static final String ID_STATE_TRIGGER_VM_STATE = "Host Connection State";
   public static final String ID_EVENT_TRIGGER_VM_POWERED_ON = "VM powered on";
   public static final String ID_EVENT_TRIGGER_VM_POWERED_OFF = "VM powered off";
   public static final String ID_STATE_TRIGGER_HOST_NETWORK_USAGE = "Host Network Usage";
   public static final String ID_ALARM_ACTION_ENTER_MAINTENANCE_MODE =
         "Enter maintenance mode";
   public static final String ID_ALARM_ACTION_ENTERING_MAINTENANCE_MODE =
         "Entering maintenance mode";
   public static final String ID_ALARM_ACTION_EXIT_MAINTENANCE_MODE =
         "Exit maintenance mode";
   public static final String ID_ALARM_TRIGGER_DATASTORE_DISK_USUAGE =
         "Datastore Disk Usage";
   public static final String ID_ALARM_TRIGGER_DATASTORE_DISK_PROVISIONED =
         "Datastore Disk Provisioned";
   public static final String ID_ALARM_ACTION_RUN_A_COMMAND = "Run a command";
   // Alarms grid settings
   public static final String ID_DISPLAY_TEXT_FIELD = "displayTextField";
   public static final String ID_ALARM_NAME_LABEL = "alarmDetailsName" + "/"
         + ID_DISPLAY_TEXT_FIELD;
   public static final String ID_DEFINED_IN_SECTION = "alarmDetailsDefinedIn";
   public static final String ID_ALARM_DEFINED_IN_LABEL = ID_DEFINED_IN_SECTION + "/"
         + ID_DISPLAY_TEXT_FIELD;
   public static final String ID_ALARM_DEFINED_IN_ICON = ID_DEFINED_IN_SECTION + "/"
         + "displayIcon";
   public static final String ID_ALARM_DEESCRIPTION_LABEL = "alarmDetailsDescription"
         + "/" + ID_DISPLAY_TEXT_FIELD;
   public static final String ID_ALARM_MONITOR_TYPE_LABEL = "alarmDetailsMonitorType"
         + "/" + ID_DISPLAY_TEXT_FIELD;
   public static final String ID_ALARM_ENABLED_LABEL = "alarmDetailsEnabled" + "/"
         + ID_DISPLAY_TEXT_FIELD;
   public static final String ID_TRIGGER_STACK_BLOCK = "triggersBlock";
   public static final String ID_TRIGGER_EVENTS_LABEL = "lblAllAny";
   public static final String ID_ALARM_ACTIONS_LABEL = "lblActionsHeaderCollapsed";
   public static final String ID_ALARMS_GENERAL_PAGE = "step_alarmGeneralPage";
   public static final String ID_ALARMS_TRIGGERS_PAGE = "step_alarmTriggersPage";
   public static final String ID_ALARMS_ACTIONS_PAGE = "step_alarmActionsPage";
   public static final String ID_BUTTON_NEW_ALARM_WIZARD_CONTAINER_ACTION_CONDITION_ADD =
         "eventConditions/toolbarControl/toolContainer/add";
   public static final String ID_TRIGGERS_BLOCK_ARROW_BUTTON = "arrowImagetriggersBlock";
   public static final String ID_ACTIONS_BLOCK_ARROW_BUTTON = "arrowImageactionsBlock";
   public static final String ID_ALARM_DETAILS_NAME = "alarmDetailsName";

   // Client_Settings
   public static final String ID_LABEL_USER_MENU_REMOVE_STORED_DATA =
         "Remove Stored Data...";
   public static final String ID_LABEL_USER_MENU_RESET_TO_FACTORY_DEF =
         "Reset To Factory Defaults";
   public static final String ID_CHECKBOX_REMOVE_STORED_DATA =
         "com.vmware.usersettings.tiwo";
   public static final String ID_CHECKBOX_GETTING_STARTED_PREFERENCES =
         "com.vmware.usersettings.gettingStarted";
   public static final String ID_CHECKBOX_SAVED_SEARCHES =
         "com.vmware.usersetting.savedsearches";
   public static final String ID_LIST_WORK_IN_PROGRESS = "tiwoListView";
   public static final String ID_ICON_SAVE_SEARCH =
         "TreeNodeItem_vsphere.core.navigator.savedSearches";
   public static final String ID_LABEL_SAVE_SEARCH_COUNT = "countLabel";
   public static final String ID_MENU_LIST_REMOVE_STORED_DATA = "removeStoredDataForm";
   public static final String ID_IMAGE_ICON_CLOSE = "closeIcon";
   public static final String ID_BUTTON_OK_CLOSE_TAB = "okBtn";

   // Advanced System Settings
   public static final String ID_GRID_PARAMETER_LIST = "resourcesGrid";
   public static final String ID_TEXT_PARAMETER_VALUE_TYPE = "textEdit";
   public static final String ID_PARAMETER_VALUE_TYPE_NUMBER = "numberEdit";
   public static final String ID_BUTTON_ADV_SET_EDIT =
         "vsphere.core.host.editAdvancedHostSettingsAction/button";
   public static final String ID_COMBOBOX_ENABLE_DISABLE = "choiceEdit";
   public static final String ID_BUTTON_ADV_SETTINGS_CANCEL = "cancelButton";
   public static final String ID_TEXTFIELD_PARAMETER_EDIT = "textDisplay";


   // Upgrade vSphere Distributed Switch wizard
   public static final String ID_RADIO_UPGRADE_VDS_VERSION_51 = "label="
         + UPGRADE_VDS_VERSION_51;
   public static final String ID_RADIO_UPGRADE_VDS_VERSION_50 = "label="
         + UPGRADE_VDS_VERSION_50;
   public static final String ID_RADIO_UPGRADE_VDS_VERSION_41 = "label="
         + UPGRADE_VDS_VERSION_41;

   // Host SecurityProfile
   public static final String ID_SPARK_LIST_HOST_SECURITYPROFILE = "tocTree";
   public static final String ID_BUTTON_EDIT_HOSTSERVICES_SETTINGS =
         "btn_vsphere.core.host.editHostServices";
   public static final String ID_BUTTON_EDIT_HOSTFIREWALL_SETTINGS =
         "btn_vsphere.core.host.editHostFirewallInfo";
   public static final String ID_HEADER_HOSTSERVICE_NAME = "Name";
   public static final String ID_BUTTON__HOSTSERVICE_STARTBUTTON =
         "serviceDetail/automationName=Start";
   public static final String ID_BUTTON__HOSTSERVICE_RESTARTBUTTON =
         "serviceDetail/automationName=Restart";
   public static final String ID_BUTTON__HOSTSERVICE_STOPBUTTON =
         "serviceDetail/automationName=Stop";
   public static final String ID_LABEL_HOSTSERVICE_RUNNING_STATUS =
         "_ServiceDetailView_LabelEx2";
   public static final String ID_ADVANCEDDATAGRID_HOST_SERVICES = "serviceList";
   public static final String ID_STARTUP_POLICY_COMBOBOX = "policy";
   public static final String ID_POLICY_DESCRIPTION_LABEL = "policyDesc";
   public static final String ID_ADVANCEDDATAGRID_HOST_FIREWALL = "firewallList";
   public static final String ID_PROPERTYGRID_HOST_FIREWALL_INPUT = "firewallInPropGrid";
   public static final String ID_PROPERTYGRID_HOST_FIREWALL_OUTPUT =
         "firewallOutPropGrid";
   public static final String ID_PROPERTYGRIDROW_HOST_FIREWALL_INPUT =
         "firewallInPropGrid/_FirewallView_PropertyGridRow1";
   public static final String ID_PROPERTYGRIDROW_HOST_FIREWALL_OUTPUT =
         "firewallOutPropGrid/_FirewallView_PropertyGridRow2";
   public static final String ID_VIEW_HOST_SECURITYPROFILE_FIREWALLVIEW =
         "vsphere.core.host.manage.settings.securityProfile.firewallView";
   public static final String ID_VIEW_HOST_SECURITYPROFILE_SERVICESVIEW =
         "vsphere.core.host.manage.settings.securityProfile.servicesView";
   public static final String ID_VIEW_HOST_SECURITYPROFILELOCKDOWNMODEVIEW =
         "vsphere.core.host.manage.settings.securityProfile.lockdownView";
   public static final String ID_VIEW_HOST_SECURITYPROFILE_HIPALVIEW =
         "vsphere.core.host.settings.securityProfile.imageProfileAcceptanceLevelView";
   public static final String ID_VIEW_HOST_SECURITYPROFILE_EDITFIREWALLVIEW =
         "editSecurityProfile";
   public static final String ID_VIEW_HOST_SECURITYPROFILE_EDITLOCKDOWNVIEW =
         "editLockdownView";
   public static final String ID_VIEW_HOST_SECURITYPROFILE_EDITHIPALVIEW =
         "editImageProfileView";
   public static final String ID_BUTTON_EDIT_HOSTACCEPTANCE_LEVEL =
         "btn_vsphere.core.host.editImageProfile";
   public static final String ID_BUTTON_EDIT_HOSTLOCKDOWN_MODE =
         "btn_vsphere.core.host.editLockdown";
   public static final String ID_BUTTON_EDIT_FIREWALL_SETTING =
         "btn_vsphere.core.host.editHostFirewallInfo";
   public static final String ID_COMBOBOX_HOSTACCEPTANCE_LEVEL = "acceptLevel";
   public static final String ID_LABEL_HOSTACCEPTANCE_LEVEL = "acceptanceLevelValue";
   public static final String ID_CHECKBOX_HOST_LOCKDOWNMODE = "enableLockdown";
   public static final String ID_LABEL_HOSTLOCKDOWN_MODE = "lockdownValue";
   public static final String ID_TEXTFIELD_IP_EDIT = "textDisplay";
   public static final String ID_CHECKBOX_ALLOWED_IP = "allIp";

   // Schedule Taks ID
   public static final String ID_SCHEDULE_TASK_OPTION_WINDOW =
         "dataGroup/step_wizardScheduleOptionsPage";
   public static final String ID_CHANGE_SCHEDULE_TASK_OPTION =
         "wizardScheduleOptionsPage/_SchedulingOptionsPage_PropertyGridRow3/schedulingOptionsBtn";
   public static final String ID_SCHEDULE_TASK_RUN_NOW_OPTION = "runNowRadio";
   public static final String ID_SCHEDULE_TASK_TASK_NAME = "nameTxt";
   public static final String ID_SCHEDULE_TASK_TASK_DESCRIPTION = "descrTxt";
   public static final String ID_TEXTFIELD_SCHEDULE_TASK_EMAIL = "emailTxt";
   public static final String ID_RADIOBUTTON_RUN_AFTER_STARTUP = "runAfterStartupRadio";
   public static final String ID_RADIOBUTTON_RUN_LATER = "runLaterRadio";
   public static final String ID_RADIOBUTTON_RUN_RECURRENT = "runRecurrentRadio";
   public static final String ID_RADIOBUTTON_HOURLY = "hourlyRadio";
   public static final String ID_RADIOBUTTON_DAILY = "dailyRadio";
   public static final String ID_RADIOBUTTON_WEEKLY = "weeklyRadio";
   public static final String ID_RADIOBUTTON_MONTHLY = "monthlyRadio";
   public static final String ID_RADIOBUTTON_WEEKlY = "_WeeklyRecurrenceView_Label1";
   public static final String ID_CHEKBOX_MONDAY = "monChk";
   public static final String ID_CHEKBOX_TUESDAY = "tueChk";
   public static final String ID_CHEKBOX_WEDNESDAY = "wedChk";
   public static final String ID_CHEKBOX_THRUSDAY = "thuChk";
   public static final String ID_CHEKBOX_FRIDAY = "friChk";
   public static final String ID_CHEKBOX_SATURDAY = "satChk";
   public static final String ID_CHEKBOX_SUNDAY = "sunChk";
   public static final String ID_RADIOBUTTON_DAY_OF_MONTH = "dayOfMonthRadio";
   public static final String ID_RADIOBUTTON_WEEK_OF_MONTH = "dayOfWeekRadio";
   public static final String ID_TEXTBOX_AFTER_STARTUP = "minutesTxt/textDisplay";
   public static final String ID_TEXTBOX_RECURRING_EVERY_HOUR =
         "hourlyRecurrenceView/intervalTxt/textDisplay";
   public static final String ID_TEXTBOX_RECURRING_EVERY_MINUTE =
         "hourlyRecurrenceView/minutesTxt/textDisplay";
   public static final String ID_TEXTBOX_RECURRING_INTERVAL =
         "recurrentSchedulerGroup/dailyRecurrenceView/intervalTxt/textDisplay";
   public static final String ID_TEXTBOX_RECURRING_TIME =
         "recurrentSchedulerGroup/perDayGroup/occurOnceTime";
   public static final String ID_RADIOBUTTON_OCCUR_ONCE = "occurOnceRadio";
   public static final String ID_TEXTBOX_WEEKLY_RECURRENCE =
         "weeklyRecurrenceView/textDisplay";
   public static final String ID_TEXTBOX_DAY_OF_MONTH =
         "recurrentSchedulerGroup/monthlyRecurrenceView/dayOfMonthTxt/textDisplay";
   public static final String ID_TEXTBOX_INTERVAL_MONTH =
         "monthlyRecurrenceView/intervalTxt/textDisplay";
   public static final String ID_LABEL_MONTHLY_RECURRENCE =
         "monthlyRecurrenceView/weekOfMonthLst/labelDisplay";
   public static final String ID_LABEL_WEEKLEY_RECURRENCE =
         "monthlyRecurrenceView/dayOfWeekLst/labelDisplay";
   public static final String ID_TEXTBOX_MONTHLY_TIME_OCCURANCE =
         "recurrenceOptions/recurrentSchedulerGroup/perDayGroup/occurOnceTime";
   public static final String ID_TEXTBOX_DAYOFWEEK =
         "dayOfWeekIntervalTxt/textDisplay/textDisplay";
   public static final String ID_BUTTON_WEEKOF_MONTH = "weekOfMonthLst/openButton";
   public static final String ID_BUTTON_DAY_OF_WEEK = "dayOfWeekLst/openButton";
   public static final String ID_TEXTBOX_DAY_OF_WEEK =
         "dayOfWeekIntervalTxt/textDisplay";
   public static final String ID_TEXTBOX_CURRENT_SCHEDULER = "currentScheduler";
   public static final String ID_TEXTBOX_USER_EMAIL = "emailTxt/textDisplay";
   public static final String ID_BUTTON_RUN_ACTION =
         "vsphere.client.scheduling.runScheduledOperationAction/button";
   public static final String ID_BUTTON_EDIT_ACTION =
         "vsphere.client.scheduling.editScheduledOperationAction/button";
   public static final String ID_BUTTON_REMOVE_ACTION =
         "vsphere.client.scheduling.removeScheduledOperationAction/button";
   public static final String ID_LIST_SCHEDULING_OPS = "list";
   public static final String ID_TEXTBOX_RUN_LATER_DATE =
         "runLaterDateTimePicker/dateField";
   public static final String ID_TEXTBOX_RUN_LATER_TIME =
         "runLaterDateTimePicker/timeField";
   public static final String ID_TEXTBOX_TASK = "task";
   public static final String ID_TEXTBOX_SCHEDULED = "schedule";
   public static final String ID_TEXTBOX_NEXT_RUN = "nextRun";
   public static final String ID_TEXTBOX_SCHEDULE_DETAILS =
         "ScheduledOperationDetailsView_Text7";
   public static final String ID_TEXTBOX_LAST_MODIFIED = "lastModified";
   public static final String ID_TEXTBOX_MODIFIED_BY = "modifiedBy";
   public static final String ID_BUTTON_REMOVE_TASK = "yesBtn";
   public static final String ID_BUTTON_SHUTDOWN_GUEST = "confirmationDialog/YES";
   public static final String ID_BUTTON_TASK_REMOVE_CONFIRM_YES =
         "suppressRemoveScheduledTaskPrompt/yesBtn";
   public static final String ID_BUTTON_TASK_REMOVE_CONFIRM_NO =
         "suppressRemoveScheduledTaskPrompt/noBtn";
   public static final String ID_LABEL_CURRENT_SCHEDULED = "schedule";
   public static final String ID_LABEL_ERROR_FOR_MISSING_NAME = "_errorLabel";
   public static final String ID_PANEL_NO_TASK_NAME =
         "automationName=No scheduled task name is specified";
   public static final String ID_BUTTON_CANCEL_TASK = "cancelButton";
   public static final String ID_LABEL_TASK_NAME = "automationName=Task name:";
   public static final String ID_LABEL_TASK_DESCRIPTION =
         "automationName=Task description:";
   public static final String ID_LABEL_CONFIGURED_SCHEDULER =
         "automationName=Configured Scheduler:";
   public static final String ID_LABEL_EMAIL_TO_USER =
         "automationName=Send email to the following addresses when the task is complete:";
   public static final String ID_TEXT_OCURRENCE_ONCE_TIME = "occurOnceTime";
   public static final String ID_COMBOBOX_WEEKOF_MONTH = "weekOfMonthLst";
   public static final String ID_COMBOBOX_DAY_OF_WEEK = "dayOfWeekLst";
   public static final String ID_RADIOBUTTON_DAYOFWEEK = "dayOfWeekRadio";
   public static final String ID_RADIOBUTTON_DAYOFMONTH = "dayOfMonthRadio";
   public static final String ID_TEXTBOX_INTERVAL = "intervalTxt";
   public static final String ID_TEXTBOX_DAYOFMONTH = "dayOfMonthTxt";
   public static final String ID_TEXTBOX_DAY_OF_WEEK_INTERVAL = "dayOfWeekIntervalTxt";
   public static final String ID_BUTTON_NEW_SCHEDULED_TASK =
         "automationName=Schedule New Task";
   public static final String ID_BUTTON_RUN_DPM_SCHEDULED_TASK =
         ID_ACTION_RUN_SCHEDULED_TASK + "/button";
   public static final String ID_BUTTON_DRS_SCHEDULE = "scheduleDrsButton";

   // Datastore Tab

   public static final String ID_SEARCH_FILE_FOLDER =
         "searchBox/automationName=searchInput";
   public static final String ID_SEARCH_BUTTON_FILE =
         "searchBox/searchButtonContainer/automationName=searchButton";
   public static final String ID_DATASTORE_SEARCH_VIEW =
         "datastoreExplorerControl/datagrid";
   public static final String ID_DATASTORE_SEARCH_RESULT_VIEW =
         "datastoreExplorerControl/searchResults";
   public static final String ID_DATASTORE_TREE_VIEW = "datastoreExplorerControl/tree";
   public static final String ID_DATASTORE_FOLDER =
         "vsphere.core.datastore.explorer.createFolderAction/button";
   public static final String ID_TEXT_BOX_DATASTORE = "folderName";
   public static final String ID_BUTTON_DATASTORE_CREATE = "createButton";
   public static final String ID_BUTTON_CREATE_FOLDER = "automationName=New Folder";
   public static final String ID_BUTTON_DATASTORE_CANCEL = "cancelButton";
   public static final String ID_BUTTON__BROWSE_DATASTORE_REFRESH =
         "vsphere.core.datastore.explorer.refreshDatastoreAction/button";
   public static final String ID_BUTTON_DATASTORE_UPLOAD =
         "automationName=Install the Client Integration Plugin-in to enable file transfers";
   public static final String ID_DATASTORE_DELETE_FOLDER =
         "afContextMenu.vsphere.core.datastore.explorer.deleteFileAction";
   public static final String ID_RETURN_TO_FILE_BROWSER = "showButton";
   public static final String ID_REGISTER_VM =
         "afContextMenu.vsphere.core.datastore.registerVmAction";
   public static final String ID_DOWNLOAD_DATASTORE =
         "afContextMenu.vsphere.core.datastore.explorer.downloadFromDatastoreAction";
   public static final String ID_CONFIRMATION_TEXT =
         "automationName=This action is not available for any of the selected objects at this time.";
   public static final String ID_DELETE_FILE_OK_BUTTON = "automationName=OK";
   public static final String ID_DATASTORE_DELETE_FOLDER_ADMIN = new StringBuffer(
         ID_DATASTORE_DELETE_FOLDER).append("/").append("automationName=Delete File")
         .toString();
   public static final String ID_TEXT_BOX_REGISTER_VM_NAME = "regVmName";


   // Role Manager IDs
   public static final String ID_ROLE_CREATE = "create";
   public static final String ID_ROLE_REMOVE = "remove";
   public static final String ID_ROLE_CLONE = "clone";
   public static final String ID_ROLE_EDIT = "edit";
   public static final String ID_ROLE_NAME = "SinglePageDialog_ID/roleName";
   public static final String ID_ROLES_LIST = "rolesView";
   public static final String ID_TEXT_MESSAGE = "_message";

   // Time Configuration
   public static final String ID_BUTTON_EDIT_TIMECONFIG =
         "btn_vsphere.core.host.editTimeConfigAction";
   public static final String ID_RADIOBUTTON_CONFIGURE_DATE_TIME = "rbConfigManually";
   public static final String ID_RADIOBUTTON_NTP = "rbConfigWithNTP";
   public static final String ID_TEXT_AREA_NTP = "areaNtpServers";
   public static final String ID_LABEL_NTP_SERVER = "lblNtpServers";
   public static final String ID_BUTTON_START = "btnStart";
   public static final String ID_BUTTONSTOP = "btnStop";
   public static final String ID_BUTTON_RESTART = "btnRestart";
   public static final String ID_COMBOBOX_ACTIVATION_POLICY = "comboStartupPolicy";
   public static final String ID_LABEL_NTP_STATUS = "lblNtpStatus";
   public static final String ID_LABEL_DATE_AND_TIME = "lblDateAndTime";
   public static final String ID_TEXT_FIELD_DATE = "dateField";
   public static final String ID_TEXT_FIELD_TIME = "timeField";
   public static final String ID_TREE_VIEW_HOST_MANAGE_SETTINGS_SYSTEM_TIME_CONFIG =
         "tocTree";

   // Service Health Management
   public static final String ID_TAB_SERVICE_HEALTH = "Service Health";
   public static final String ID_IMAGE_BUTTON_EXPAND_ALL_ = "expandAll";
   public static final String ID_IMAGE_BUTTON_COLLAPSE_ALL = "collapseAll";
   public static final String ID_TREE_VIEW_HEALTH_STATUS =
         "_ServiceHealthView_ServiceHealthTreeView1";
   public static final String ID_PANEL_ITEM_COUNT_INFO = "infoPanel1";
   public static final String ID_HEALTH_WINDOW_TITLE = "titleText";
   public static final String ID_WINDOW_CONTAINER_HEALTH_STATUS_ =
         "vsphere.core.folder.monitor.healthView.container";
   public static final String ID_LINK_HEALTH_SERVICE =
         "TreeNodeItem_vsphere.core.applicationnavigator.serviceStatus";
   public static final String ID_LABEL_HEALTH_SUMMARY = "healthSummaryLabel";

   // Sessions
   public static final String ID_TAB_SESSIONS = "Sessions";
   public static final String ID_LIST_USER_SESSIONS = "userSessionsList";
   public static final String ID_LABEL_ACTIVE_SESSIONS = "_VcSessionsView_Label1";
   public static final String ID_BUTTON_TERMINATE_SESSION = "terminateSessionBtn";
   public static final String ID_BUTTON_EDIT_MSG_OF_THE_DAY =
         "btn_vsphere.core.folder.editMessageOfTheDayAction";

   // Image_customization
   public static final String ID_LIST_SPEC = "specList";
   public static final String ID_COMBOBOX_SPEC_GOS_SELECTOR = "gosSelector";
   public static final String ID_TEXT_SPEC_NAME = "custSpecName";
   public static final String ID_TEXT_SPEC_DESCRIPTION = "custSpecDesc";
   public static final String ID_LIST_PAGE_TITLE = "dataGroup";
   public static final String ID_CHECKBOX_SYSPREP_ANSWER_FILE = "useSysPrepChkBox";
   public static final String ID_RADIOBUTTON_IMPORT_ANSWER_FILE = "importSysprepFile";
   public static final String ID_RADIOBUTTON_CREATE_ANSWER_FILE = "createSysprepFile";
   public static final String ID_TEXT_AREA_CREATE_SYSPREP = "sysprepTextContent";
   public static final String ID_TEXT_SPEC_ORG_NAME = "organizationName";
   public static final String ID_TEXT_SPEC_OWNER_NAME = "ownerName";
   public static final String ID_RADIOBUTTON_SPEC_COMP_NAME =
         "computerNamePage/enterName";
   public static final String ID_CHECKBOX_APPEND_NUMERIC_VALUE = "includeUniqueNumber";
   public static final String ID_RADIOBUTTON_SPEC_VM_NAME = "useVmName";
   public static final String ID_RADIOBUTTON_SPEC_DEPLOY_WIZARD =
         "computerNamePage/enterNameInDeployWizard";
   public static final String ID_TEXT_SPEC_COMP_NAME = "computerName";
   public static final String ID_TEXT_SPEC_PROD_KEY = "productKey";
   public static final String ID_TEXT_ADMIN_PASSWD = "adminPwd";
   public static final String ID_TEXT_CONFIRM_ADMIN_PASSWD = "adminPwdConfirm";
   public static final String ID_LIST_SPEC_TIME_ZONE_SELECTOR = "timeZoneSelector";
   public static final String ID_TEXT_SPEC_COMMAND_ENTRY = "entryInput";
   public static final String ID_BUTTON_SPEC_COMMAND_ADD = "addButton";
   public static final String ID_RADIOBUTTON_SPEC_STD_NETWORK_SETTINGS =
         "standardSettings";
   public static final String ID_RADIOBUTTON_SPEC_CUSTOM_NETWORK_SETTINGS_RADIOBTN =
         "customSettings";
   public static final String ID_RADIOBUTTON_SPEC_WORKGROUP = "workgroupRadio";
   public static final String ID_RADIOBUTTON_DOMAIN = "domainRadio";
   public static final String ID_TEXT_DOMAIN_USERNAME = "username";
   public static final String ID_TEXT_DOMAIN_PASSWD = "password";
   public static final String ID_TEXT_DOMAIN_CONFIRM_PASSWD = "confirmPassword";
   public static final String ID_CHECKBOX_GENERATE_SID = "generateSid";
   public static final String ID_DIALOG_SPEC_YESNO = "YesNoDialog";
   public static final String ID_BUTTON_ADMIN_PASSWD_YES = "automationName=Yes";
   public static final String ID_BUTTON_ADMIN_PASSWD_NO = "automationName=No";
   public static final String ID_LABEL_TIME_ZONE_PAGE = "step_timeZonePage";
   public static final String ID_BUTTON_SPEC_CREATE_NEW =
         "vsphere.core.vm.customizeGos.createNewAction/button";
   public static final String ID_BUTTON_SPEC_DELETE =
         "vsphere.core.vm.customizeGos.removeAction/button";
   public static final String ID_TREE_NODE_SPEC =
         "TreeNodeItem_vsphere.core.navigator.customizationSpecManager";
   public static final String ID_TEXTFIELD_DOMAIN_NAME = "domainNameGroup";
   public static final String ID_LABEL_SPEC_NAME = "viewNameLabel";
   public static final String ID_LABEL_SPEC_DESC = "viewDescriptionLabel";
   public static final String ID_COMBOBOX_SERVER = "serverComboBox";
   public static final String ID_BUTTON_GENERATE_CUSTOM_NAME = "generateNameGroup";
   public static final String ID_BUTTON_FILTER_CONTROL = "filterControl";
   public static final String ID_BUTTON_SPEC_EDIT =
         "vsphere.core.vm.customizeGos.editAction/button";
   public static final String ID_BUTTON_SPEC_IMPORT =
         "vsphere.core.vm.customizeGos.importAction/button";
   public static final String ID_BUTTON_SPEC_EXPORT =
         "vsphere.core.vm.customizeGos.exportAction/button";
   public static final String ID_BUTTON_SPEC_DUPLICATE =
         "vsphere.core.vm.customizeGos.duplicateAction/button";
   public static final String ID_DIALOG_EDIT_SPEC = "tiwoDialog";
   public static final String ID_GRID_SPEC_PROPERTIES = "headerPropertyGrid";
   public static final String ID_BUTTON_BACK_SPEC = "back";
   public static final String ID_BUTTON_CANCEL_SPEC = "cancel";
   public static final String ID_LABEL_NIC_NAME = "automationName=NIC1";
   public static final String ID_BUTTON_ADD_NIC = "button";
   public static final String ID_BUTTON_EDIT_NIC = "automationName=Edit";
   public static final String ID_BUTTON_DELETE_NIC = "button";
   public static final String ID_CHECKBOX_AUTO_LOG_ON = "autoLogonChkBox";
   public static final String ID_NUMERIC_STEPPER_LOG_ON_TIMES = "logonTimes";
   public static final String ID_BUTTON_SPEC_COMMAND_DELETE = "deleteButton";
   public static final String ID_BUTTON_SPEC_COMMAND_MOVE_UP = "moveUpButton";
   public static final String ID_BUTTON_SPEC_COMMAND_MOVE_DOWN = "moveDownButton";
   public static final String ID_DIALOG_BLANK_ADMIN_PASSWORD =
         "automationName=Confirm Blank Administrator Password";
   public static final String ID_TEXT_NO_PRIVILEGE_ERROR_MSG = "error_0";

   public static final String ID_VMPROVISIONING_DSCLUSTER_GRID =
         "storageLocator/storageList";
   public static final String ID_DIALOG_LUN = "automationName=Select Target LUN";
   public static final String ID_LABEL_HARD_DISK = "Hard disk";
   public static final String ID_INFO_LIST = "infoGrid";
   public static final String ID_TOOL_TIP = "toolTip";

   // EVC CPU Specific
   public static final String ID_LABEL_EVC_STATUS_SUMMARY = "evcStatus";
   public static final String ID_LABEL_EVC_HOST_SUMMARY = "currentEvcModeText";
   public static final String ID_LABEL_SUPPORTED_EVC_MODE = "supportedEvcModeText";
   public static final String ID_BUTTON_EDIT_EVC_CLUSTER =
         "btn_vsphere.core.cluster.evc.configureEvcAction";
   public static final String ID_BUTTON_EVC_DISABLE = "disableEvcButton";
   public static final String ID_BUTTON_EVC_AMD = "amdVendorButton";
   public static final String ID_BUTTON_EVC_INTEL = "intelVendorButton";
   public static final String ID_LIST_VMWARE_EVC_MODES = "evcVendorSpecificModesList";
   public static final String ID_TEXTINPUT_COMPATIBILITY_MESSAGE = "itemMessage";
   public static final String ID_ADVANCEDDATAGRIDHEADER_DC_VM_LIST =
         "vmsForDatacenter/list";
   public static final String ID_ADVANCEDDATAGRIDHEADER_CLS_VM_LIST =
         "vmsForCluster/list";
   public static final String ID_LABEL_HOST_ADD_COMPATIBILITY_CHECK = "_errorLabel";
   public static final String ID_TEXTINPUT_ADDHOST = "inputLabel";
   public static final String ID_ADVANCEDDATAGRID_CLUSTER_NAME = "Label_4";
   public static final String ID_MOVEHOST_PAGE_NAVTREE =
         "/selectResourcePoolPage/navTreeView/navTree";
   public static final String ID_LABEL_EVC_MODE_EXPAND = "titleLabel.evcStack";

   // Memory Compression
   public static final String ID_LABEL_COMPRESSED_VM_MEMORY = "compressedMemoryLbl";

   // Host Connected USB IDS
   public static final String ID_BUTTON_EDIT_ALL_VM_SETTINGS =
         "btn_vsphere.core.vm.provisioning.editAction";
   public static final String ID_COMBOBOX_HARDWARE_VERSION = "version";

   // constants for Solution Manager
   public static final String ID_EXTENSION_BUTTON =
         "tabBar/automationName=Extension Types";
   public static final String ID_SAMPLE_SOLUTION_LABEL =
         "automationName=EAM Sample Solution";
   public static final String ID_CLUSTER_HOST_VMS_LIST = "list";
   public static final String ID_GRID_AGENT_SUMMARY_VIEW = "relatedGrid";
   public static final String ID_LABEL_RELATED_ITEM = "RelatedItemsPropertyGridRow1";
   public static final String ID_LABEL_AGENT_RESOURCE_POOL =
         "automationName=Resource pool";
   public static final String ID_RELATED_ITEM_AGENT_VM_PORTLET = EXTENSION_PREFIX
         + EXTENSION_ENTITY_VM + "relatedItemsView.chrome";
   public static final String ID_LABEL_RESOURCE_POOL =
         "relatedItemsView/relatedGrid/labelItem";
   public static final String ID_LIST_DC_TOP_LEVEL_ITEMS =
         "childObjectsForDatacenter/list";
   public static final String ID_BUTTON_ESX_AGENT_FOLDER = "ESX Agents";
   public static final String ID_VC_EXTENSION_LIST = "extensionsOfVirtualCenters/list";
   public static final String ID_LABEL_ESX_AGENT_MANAGER =
         "automationName=vSphere ESX Agent Manager";

   //constants for Log Browser
   public static final String ID_LABEL_LOG_TYPE = "automationName=Type:";
   public static final String ID_LABEL_TYPE_DISPLAY = "labelDisplay";
   public static final String ID_LABEL_ADJACENT_VIEW = "automationName=Adjacent:";
   public static final String ID_BUTTON_ACTION = "Actions";
   public static final String ID_TEXT_FILTER =
         "mainTable/topControlBar/noObjectSelector/textInput";
   public static final String ID_LABEL_MAIN_TITLE =
         "mainTable/mainDataArea/refreshScreen/mainLabel";
   public static final String ID_BUTTON_REFRESH_LOG_BROWSER = "Refresh";
   public static final String ID_BUTTON_RETRIEVE_LOGS = "retrieveButton";


   // Client Connected USB
   public static final String ID_BUTTON_USB_DEVICES = "view_UsbDevices/menuButton";
   public static final String ID_LABEL_USB_DEVICE_ATTACHED =
         "vmHardwareView/view_Usb_1/headerLabel";
   public static final String ID_MENU_USB_DEVICE_NAME = "menu_view_UsbDevices";
   public static final String ID_LABEL_USB_DEVICE_REMOVE = "view_Usb_1/menuButton";
   public static final String ID_LABEL_USB_DEVICE_DISCONNECT = "menu_view_Usb_1";
   public static final String ID_LABEL_USB_DEVICE_STATUS =
         "vmHardwareView/view_UsbDevices/headerLabel";
   public static final String ID_PROPERTY_GRID_VM_DATASTORE =
         "RelatedItemsPropertyGridRow3";
   public static final String ID_BUTTON_USB_DEVICE_REMOVE_LINK_BUTTON = "remove";
   public static final String ID_CHECKBOX_HOST_USB_VMOTION_SUPPORT = "remoteHostCheck";

   // Cluster Manage->Setting-Configuration
   public static final String ID_ADVANCED_DATAGRID_HOST_OPTIONS = "hostOptionsGrid";
   public static final String ID_DROPDOWNLIST_HOST_POWER = "ddlDpmBehaviour";

   // iSCSI Static and Dynamic Discovery
   public static final String ID_CHECK_BOX_INHERIT_SETTINGS = "inheritCheckBox";
   public static final String ID_LIST_AUTH_METHODS = "authMethods";
   public static final String ID_LABEL_AUTH_METHODS = "authMethods/labelDisplay";
   public static final String ID_TEXT_INPUT_OUTGOING_NAME = "outgoingName";
   public static final String ID_TEXT_INPUT_INCOMING_NAME = "incomingName";
   public static final String ID_CHECK_BOX_OUTGOING = "outgoingCheckBox";
   public static final String ID_CHECK_BOX_INCOMING = "incomingCheckBox";
   public static final String ID_TEXT_INPUT_OUTGOING_SECRET = "outgoingSecret";
   public static final String ID_TEXT_INPUT_INCOMING_SECRET = "incomingSecret";
   public static final String ID_LABEL_VALIDATION_ERROR = "_titleLabel";
   public static final String ID_TEXT_LABEL_VALIDATION_ERROR = "_errorLabel";
   // Permissions constants
   public static final String ID_WIZARD_CONTENT = "wizardContent";
   public static final String ID_BUTTON_HELP = "helpButton";
   public static final String ID_BUTTON_CHECKNAMES = "checkNamesButton";
   public static final String ID_TEXTINPUT_USERSELECTION = "usersSelectionTextInput";
   public static final String ID_TEXTINPUT_GROUPSELECTION = "groupsSelectionTextInput";
   public static final String ID_COMBOBOX_DOMAINS = "domainsComboBox";
   public static final String ID_COMBOBOX_SORTTYPES = "sortTypes";
   public static final String ID_TEXTINPUT_SEARCH = "searchInput";
   public static final String ID_WIZARD_SEARCH_INPUT = ID_WIZARD_TIWO_DIALOG + "/"
         + ID_SEARCHCONTROL_SEARCHCONTROL + "/" + ID_TEXTINPUT_SEARCH;
   public static final String ID_WIZARD_QUICK_SEARCH_ICON = ID_WIZARD_TIWO_DIALOG + "/"
         + ID_SEARCHCONTROL_SEARCHCONTROL + "/" + ID_SEARCH_STACK + "/"
         + ID_QUICK_SEARCH_ICON;
   public static final String ID_BUTTON_DMESEARCH = "searchButton";
   public static final String ID_BUTTON_SEL_USERGROUP_ADD_USER = "addUserButton";
   public static final String ID_LINK_DEFINED_IN = "This object and its children";
   public static final String ID_LINK_DEFINED_IN_NON_PROPAGATE = "This object";
   public static final String ID_BUTTON_REMOVE_USER = "removeUser";
   public static final String ID_CD_DVD_MEDIA_TEXT = "cdrom_?/client";
   public static final String ID_FLOPPY_MEDIA_TEXT = "floppy_?/client";
   public static final String ID_CHANGE_PERMISSION_WARNING = "_message";
   public static final String ID_LOADING_PROGRESS_BAR = "automationName=Loading...";

   //Edit IPMI Setting
   public static final String ID_TEXT_INPUT_IPMI_LOGIN = "ipmiLoginTxi";
   public static final String ID_TEXT_INPUT_IPMI_PASSWORD = "ipmiPasswordTxi";
   public static final String ID_TEXT_INPUT_IPMI_IP = "ipmiIpAddressCtrl";
   public static final String ID_TEXT_INPUT_IPMI_MAC = "ipmiMacAddressTxi";
   public static final String ID_BUTTON_EDIT_HOST_OPTIONS_OK_BUTTON = "okBtn";
   public static final String ID_BUTTON_EDIT_HOST_OPTIONS_CANCEL_BUTTON = "cancelBtn";
   public static final String ID_BUTTON_EDIT_HOST_OPTIONS_IPMI = "editBtn";
   public static final String ID_BUTTON_EDIT_IPMI_STANDALONE_HOST =
         "Edit_HostSystem:powerConfigSoftware";

   // vDS->Manage->Settings->Private VLAN
   public static final String ID_BUTTON_PVLAN_EDIT =
         "btn_vsphere.core.dvs.editPrivateVlanAction";
   public static final String ID_ADVDATAGRID_PVLAN = "list";

   // 'Edit Private VLAN Settings' dialog
   public static final String ID_BUTTON_ADD_PRIMARY_VLAN_ID = "_PrivateVlanPage_Button1";
   public static final String ID_BUTTON_REMOVE_PRIMARY_VLAN_ID =
         "_PrivateVlanPage_Button2";
   public static final String ID_BUTTON_ADD_SECONDARY_VLAN_ID =
         "_PrivateVlanPage_Button3";
   public static final String ID_BUTTON_REMOVE_SECONDARY_VLAN_ID =
         "_PrivateVlanPage_Button4";
   public static final String ID_ADVDATAGRID_PRIMARY_VLAN_ID = "primaryEntries";
   public static final String ID_ADVDATAGRID_SECONDARY_VLAN_ID = "secondaryEntries";
   public static final String ID_BUTTON_VLAN_TYPE_DROPDOWN = "openButton";
   public static final String ID_LABEL_VLAN_TYPE_ISOLATED = "automationName=Isolated";

   // 'Removing Primary VLAN ID' dialog
   public static final String ID_BUTTON_CONFIRM_REMOVE_PRIMARY_VLAN_ID_YES =
         "automationName=Yes";

   public static final String ID_TEXT_INPUT_ERROR_MSG = "_message";
   public static final String ID_WARNING_MSG = "warningMessage";

   // Network Protocol Profiles
   public static final String ID_CREATE_NPPS_NAME = "nameTextInput";
   public static final String ID_CREATE_NPPS_NETW_ASSOC = "networkAssociationList";
   public static final String ID_CREATE_NPPS_NETW_ASSOC_PORTGROUPS_LABEL = "netLabel";
   public static final String ID_CREATE_NPPS_NETW_ASSOC_PORTGROUPS_CHECKBOX =
         "netSelect";
   public static final String ID_IPV4INPUT = "ipv4Input";
   public static final String ID_IPV6INPUT = "ipv6Input";

   public static final String ID_CREATE_NPPS_NETW_BITS = "networkBitsTextInput";
   public static final String ID_CREATE_NPPS_IPV4_SUBNET = "subnetAddressInput/ipv4";
   public static final String ID_CREATE_NPPS_IPV4_SUBNET_BITS = ID_IPV4INPUT + "/"
         + ID_CREATE_NPPS_NETW_BITS;
   public static final String ID_CREATE_NPPS_IPV4_GATEWAY = "gatewayAddressInput/ipv4";
   public static final String ID_CREATE_NPPS_IPV4_DHCP = ID_IPV4INPUT + "/dhcpCheckbox";
   public static final String ID_CREATE_NPPS_IPV4_DNS = ID_IPV4INPUT
         + "/dnsServersTextInput";
   public static final String ID_CREATE_NPPS_IPV4_ENABLE_NPP = ID_IPV4INPUT
         + "/ipPoolCheckbox";
   public static final String ID_CREATE_NPPS_IPV4_NPP_RANGE = ID_IPV4INPUT
         + "/ipPoolRangeTextInput";

   public static final String ID_CREATE_NPPS_IPV6_SUBNET = ID_IPV6INPUT
         + "/subnetAddressInput/ipv6";
   public static final String ID_CREATE_NPPS_IPV6_SUBNET_BITS = ID_IPV6INPUT + "/"
         + ID_CREATE_NPPS_NETW_BITS;
   public static final String ID_CREATE_NPPS_IPV6_GATEWAY = ID_IPV6INPUT
         + "/gatewayAddressInput/ipv6";
   public static final String ID_CREATE_NPPS_IPV6_DHCP = ID_IPV6INPUT + "/dhcpCheckbox";
   public static final String ID_CREATE_NPPS_IPV6_DNS = ID_IPV6INPUT
         + "/dnsServersTextInput";
   public static final String ID_CREATE_NPPS_IPV6_ENABLE_NPP = ID_IPV6INPUT
         + "/ipPoolCheckbox";
   public static final String ID_CREATE_NPPS_IPV6_NPP_RANGE = ID_IPV6INPUT
         + "/ipPoolRangeTextInput";

   public static final String ID_NPP_SUBTABS_ASSOC_NETW = "networksList";

   public static final String ID_IPV4CONFIG = "ipv4Config";
   public static final String ID_IPV6CONFIG = "ipv6Config";

   public static final String ID_EDIT_NPP_IPV4_RANGE_BTN = ID_IPV4INPUT
         + "/ipPoolRangeViewToggle";
   public static final String ID_EDIT_NPP_IPV6_RANGE_BTN = ID_IPV6INPUT
         + "/ipPoolRangeViewToggle";

   public static final String ID_EDIT_NPP_IPV4_RANGE_LIST = ID_IPV4INPUT
         + "/ipPoolRangeView";
   public static final String ID_EDIT_NPP_IPV6_RANGE_LIST = ID_IPV6INPUT
         + "/ipPoolRangeView";

   public static final String ID_NPP_SUBTABS_IPV4_SUBNET = ID_IPV4CONFIG
         + "/subnetLabel";
   public static final String ID_NPP_SUBTABS_IPV4_SUBNET_BITS = ID_IPV4CONFIG
         + "/netmaskLabel";
   public static final String ID_NPP_SUBTABS_IPV4_GATEWAY = ID_IPV4CONFIG
         + "/gatewayLabel";
   public static final String ID_NPP_SUBTABS_IPV4_DHCP = ID_IPV4CONFIG + "/dhcpLabel";
   public static final String ID_NPP_SUBTABS_IPV4_DNS = ID_IPV4CONFIG
         + "/dnsServersTextArea";
   public static final String ID_NPP_SUBTABS_IPV4_ENABLE_NPP = ID_IPV4CONFIG
         + "/ipPoolLabel";

   public static final String ID_NPP_SUBTABS_IPV6_SUBNET = ID_IPV6CONFIG
         + "/subnetLabel";
   public static final String ID_NPP_SUBTABS_IPV6_SUBNET_BITS = ID_IPV6CONFIG
         + "/netmaskLabel";
   public static final String ID_NPP_SUBTABS_IPV6_GATEWAY = ID_IPV6CONFIG
         + "/gatewayLabel";
   public static final String ID_NPP_SUBTABS_IPV6_DHCP = ID_IPV6CONFIG + "/dhcpLabel";
   public static final String ID_NPP_SUBTABS_IPV6_DNS = ID_IPV6CONFIG
         + "/dnsServersTextArea";
   public static final String ID_NPP_SUBTABS_IPV6_ENABLE_NPP = ID_IPV6CONFIG
         + "/ipPoolLabel";

   public static final String ID_NPP_SUBTABS_OTHER_NETW_DNS = "dnsDomainLabel";
   public static final String ID_NPP_SUBTABS_OTHER_NETW_HOST_PREFIX = "hostPrefixLabel";
   public static final String ID_NPP_SUBTABS_OTHER_NETW_DNS_SEARCH_PATH =
         "dnsSearchPathLabel";
   public static final String ID_NPP_SUBTABS_OTHER_NETW_HTTP_PROXY = "httpProxyLabel";

   //Address range in Add/Edit/Associate IP Pool wizards
   public static final String ID_NPP_IPV6_ADDRESS_RANGE = ID_IPV6INPUT
         + "/addressRangeLabel";
   public static final String ID_NPP_IPV4_ADDRESS_RANGE = ID_IPV4INPUT
         + "/addressRangeLabel";

   public static final String ID_NPP_CREATE_NPPS_IPV4_SUBNET_BITS = ID_IPV4INPUT
         + "/netmaskLabel";
   public static final String ID_NPP_CREATE_NPPS_IPV6_SUBNET_BITS = ID_IPV6INPUT
         + "/netmaskLabel";

   public static final String ID_CREATE_NPPS_OTHER_NETW_DNS = "dnsDomainTextInput";
   public static final String ID_CREATE_NPPS_OTHER_NETW_HOST_PREFIX =
         "hostPrefixTextInput";
   public static final String ID_CREATE_NPPS_OTHER_NETW_DNS_SEARCH_PATH =
         "dnsSearchPathTextInput";
   public static final String ID_CREATE_NPPS_OTHER_NETW_HTTP_PROXY =
         "httpProxyTextInput";

   public static final String ID_NPPS_ADVANCED_DATAGRID = "list";

   public static final String ID_IMAGE_EDIT = "inlineEdit";
   public static final String ID_NPP_ASSOCIATE_CREATE_NEW_NPP_RADIO_BTN =
         "workflowCreateNew";
   public static final String ID_NPP_ASSOCIATE_USE_EXISTING_NPP_RADIO_BTN =
         "workflowAssociate";

   public static final String ID_NPP_IMAGE_EDIT = "ipPoolEditAction";
   public static final String ID_NPP_IMAGE_ADD = "ipPoolAddAction";
   public static final String ID_NPP_IMAGE_DELETE = "ipPoolDeleteAction";

   public static final String ID_NPP_ADD_ASSOC_BTN = "addAssociationButton";
   public static final String ID_NPP_REMOVE_ASSOC_BTN = "removeAssociationButton";
   public static final String ID_NPP_EDIT_ASSOC_BTN = "ipPoolBlock/editButton";

   public static final String ID_NPPS_EDIT_BTN = "editButton";

   public static final String ID_NPPS_ADVANCED_DATAGRID_COLUMN_NAME = "Name";

   public static final String ID_NPPS_TAB_BUTTON = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append(MAIN_TAB_PREFIX).append(".ipPoolsView.button")
         .toString();

   public static final String ID_NPP_TAB_BUTTON = new StringBuffer(EXTENSION_PREFIX)
         .append(EXTENSION_ENTITY).append(MAIN_TAB_PREFIX).append(".ipPoolView.button")
         .toString();

   public static final String ID_DVPORTGROUP_NETWORK_DETAILS_IP_POOL = "ipPoolName";

   public static final String ID_NPPS_TAB_NAV = "ipPoolTabNav";
   public static final String ID_NPPS_TAB_NAV_OTH_NETW_CONF = "otherNetworkConfig";
   public static final String ID_NPPS_TITLE_LABEL = "titleLabel";
   public static final String ID_NPPS_TITLE_LABEL_DETAILS_PANE = "detailsSB" + "/"
         + ID_NPPS_TITLE_LABEL;
   public static final String ID_NPPS_TITLE_LABEL_DVPORTGROUP = "ipPoolBlock" + "/"
         + ID_NPPS_TITLE_LABEL;
   public static final String ID_NPP_TITLE_LABEL_NO_NPP_ASSOCIATED =
         "networksNoIpPoolAssociated";
   public static final String ID_NPP_TITLE_LABEL_NO_NPP_SELECTED =
         "networksNoIpPoolSelected";

   public static final String ID_NPPS_NETW_ASSOC_PORTGROUPS_LABEL_ALREADY_ASSOCIATED =
         "netIpPoolLabel";
   public static final String ID_NPP_NETW_ALREADY_ASSOC_ERR_MSG = "errorMessage";

   public static final String ID_NPP_WIZARD_PAGE_NAME_NETW_ASSOC =
         "ipPoolNameAndNetworkPage";
   public static final String ID_NPP_WIZARD_PAGE_IPV4 = "ipPoolIpv4ConfigPage";
   public static final String ID_NPP_WIZARD_PAGE_IPV6 = "ipPoolIpv6ConfigPage";
   public static final String ID_NPP_WIZARD_PAGE_OTH_NETW_CONF =
         "ipPoolOtherNetworkConfigPage";
   public static final String ID_NPP_WIZARD_PAGE_SET_ASSOC = "ipPoolSelectWorkflowPage";
   public static final String ID_NPP_WIZARD_PAGE_EXISTING_ASSOC =
         "ipPoolChooseExistingPage";
   public static final String ID_NPP_WIZARD_PAGE_READY_TO_COMPLETE =
         "defaultSummaryPage";


   // System logs control ID
   public static final String ID_EXPORT_SYSTEM_LOGS = "exportLogFiles";
   public static final String ID_COMBO_LOGITEMS = "logItems";
   public static final String ID_BUTTON_EXPORT_DATA = "exportDataButton";
   public static final String ID_LABEL_SYSTEMLOGS_COUNTS = "logContentDetails";
   public static final String ID_TEXTINPUT_FILTER = "filterControl";
   public static final String ID_GRID_HOSTSYSTEM_FILTER = "HostSystem_filterGridView";
   public static final String ID_CHECKBOX_INCLUDE_VCLOGS = "vcLogsCheckbox";
   public static final String ID_BUTTON_NEXT_2000_LINES = "showNext";
   public static final String ID_BUTTON_CANCEL_ALL_LINES = "automationName=Cancel";
   public static final String ID_BUTTON_SHOW_ALL_LINES = "showAll";
   public static final String ID_CHECKBOX_SHOW_LINE_NUMBERS = "showLinesNumber";
   public static final String ID_BUTTON_CANCEL_SYSTEMLOGS = "cancelButton";
   public static final String ID_BUTTON_OK_SYSTEMLOGS = "okButton";
   public static final String ID_GRID_SYSTEMLOGS_OBJECTDATA =
         "HostSystem_selectedObjectGridView";
   public static final String ID_IMAGE_WORKINGPROGRESS = "_TiwoListItemRenderer_Image1";
   public static final String ID_BUTTON_GENERATE_LOGS = "logBundleButton";
   public static final String ID_CHECKBOX_HOSTSELECTION =
         "_CheckBoxColumnRenderer_CheckBox1[1]";
   public static final String ID_LABEL_SELECT_COLUMNS =
         "automationName=Select Columns...";

   public static final String ID_ALARMS_PORTLET =
         "vsphere.opsmgmt.alarms.sidebarView.portletChrome";
   public static final String ID_ALARMS_PORTLET_SIDEBARGRID =
         "vsphere.opsmgmt.alarms.sidebarView.portletChrome/allAlarmsSideBarGrid";
   public static final String ID_ALARMS_PORTLET_COLLAPSE_BUTTON =
         "portletCollapseButton";
   public static final String ID_ALARM_PORTLET_EXPAND_BUTTON = "portletExpandButton";
   public static final String ID_ALARMS_STYLE_NAME = "styleName";
   public static final String ID_ISSUES_LIST = "issueList";
   public static final String ID_TRIGGEREDALARMS_LIST = "alarmList";
   public static final String ID_TEXT_CHANGE_PERMISSION_WARNING = "tiwoDialog/_message";
   public static final String ID_TRIGGERS_HOST_DISK_USAGE = "Host Disk Usage";
   public static final String ID_TRIGGERS_HOST_NETWORK_USAGE = "Host Network Usage";
   public static final String ID_TRIGGERS_HOST_MEMORY_USAGE = "Host Memory Usage";
   public static final String ID_TRIGGERS_HOST_CPU_USUAGE = "Host CPU Usage";
   public static final String ID_ALARM_ACTION = "ALARM";
   // Flex UI components classes constants
   public static final String FLEX_CLASS_ROLLOVER_IMAGE = "RolloverImage";
   public static final String FLEX_CLASS_TEXT_INPUT = "TextInput";
   public static final String FLEX_CLASS_TEXT_FIELD = "UITextField";
   public static final String FLEX_CLASS_TOC_TREE_ITEM_RENDERER = "TocTreeItemRenderer";
   public static final String ID_TEXT_ERROR = "errorText";

   public static final String ID_ITEM_ID = "itemId";

   //Storage 2TB+
   public static final String ID_STORAGE_CAPACITY_IMAGE = "arrowImagecapacityStackBlock";
   public static final String ID_DATASTORE_TOTAL_CAPACITY = "totalSize";
   public static final String ID_PROVISION_DATASTORE_PARTITION_FORMAT =
         "tiwoDialog/txtPartitionFormat";
   public static final String ID_PROVISION_DATASTORE_FUTURE_DATASTORE_SIZE =
         "tiwoDialog/txtFutureDatastoreSize";
   public static final String ID_PROVISION_DATASTORE_SIZE_STEPPER =
         "datastoreSizeStepper";

   public static final String ID_LIST_STORAGE_DISK_FORMAT_SELECTOR =
         "diskFormatSelector";

   // Cluster Getting started Links
   public static final String ID_CLUSTER_GETTING_STARTED_CLUSTERS =
         "gettingStartedHelpLink_1";
   public static final String ID_CLUSTER_GETTING_STARTED_RESOURCE_POOLS =
         "gettingStartedHelpLink_2";

   // Getting started Tabs
   public static final String ID_FUE_CREATE_FOLDER_BUTTON =
         "vsphere.core.folder.createFolderAction";
   public static final String ID_FUE_CREATE_ROOT_DATACENTER_BUTTON =
         "vsphere.core.datacenter.createAction";
   public static final String ID_FUE_ADD_HOST_BUTTON = "vsphere.core.host.addAction";
   public static final String ID_FUE_CREATE_CLUSTER_BUTTON =
         "vsphere.core.cluster.createAction";
   public static final String ID_FUE_CREATE_VM_ONDC_BUTTON =
         "vsphere.core.vm.provisioning.createVmAction";
   public static final String ID_FUE_CREATE_DATASTORE_BUTTON =
         "vsphere.core.datastore.addAction";
   public static final String ID_FUE_CREATE_DVPORTGROUP_BUTTON =
         "vsphere.core.dvs.createDvsAction";
   public static final String ID_FUE_ADD_CLUSTERED_HOST_BUTTON =
         "vsphere.core.host.addAction";
   public static final String ID_FUE_CREATE_VM_BUTTON =
         "vsphere.core.vm.provisioning.createVmAction";
   public static final String ID_FUE_CREATE_VM_HOST_BUTTON =
         "vsphere.core.vm.provisioning.createVmAction";
   public static final String ID_FUE_CREATE_DATACENTER_BUTTON =
         "vsphere.core.datacenter.createAction";
   public static final String ID_FUE_CREATE_VM_RESOURCEPOOL_BUTTON =
         "vsphere.core.vm.provisioning.createVmAction";
   public static final String ID_FUE_EDIT_RESOURCEPOOL_BUTTON =
         "vsphere.core.resourcePool.editAction";
   public static final String ID_FUE_CREATE_RESOURCEPOOL_BUTTON =
         "vsphere.core.resourcePool.createAction.showOnBar";
   public static final String ID_FUE_POWERON_VAPP_BUTTON =
         "vsphere.core.vApp.powerOnAction";
   public static final String ID_FUE_EDIT_VAPP_BUTTON =
         "vsphere.core.vApp.editSettingsAction";
   public static final String ID_FUE_POWERON_BUTTON = "vsphere.core.vm.powerOnAction";
   public static final String ID_FUE_FUE_POWEROFF_BUTTON =
         "vsphere.core.vm.powerOffAction";
   public static final String ID_FUE_SUSPEND_VM_BUTTON = "vsphere.core.vm.suspendAction";
   public static final String ID_FUE_EDIT_VM_BUTTON =
         "vsphere.core.vm.provisioning.editAction";
   public static final String ID_FUE_CLONE_VMTOTEMPLATE_BUTTON =
         "vsphere.core.vm.provisioning.cloneTemplateToVmAction";
   public static final String ID_FUE_CONVERT_TEMPLATETOVM_BUTTON =
         "vsphere.core.vm.provisioning.convertTemplateToVmAction";
   public static final String ID_FUE_MOVE_DATASTORE_INCLUSTER_BUTTON =
         "vsphere.core.dscluster.moveDatastoresIntoAction";
   public static final String ID_DVSETTINGS_BUTTON =
         "vsphere.core.dvs.editDvsSettingsAction";
   public static final String ID_FUE_CREATEDVS_PORTGROUP_BUTTON =
         "vsphere.core.dvs.createPortgroup";
   public static final String ID_FUE_DVPORTGROUP_EDIT_BUTTON =
         "vsphere.core.dvPortgroup.editSettingsAction";
   public static final String ID_DVPORTGROUP_EDIT_BUTTON_SUMMARY_TAB =
         "btn_vsphere.core.dvPortgroup.editSettingsAction";
   public static final String ID_FUE_EDIT_UPLINK_PORTGROUP_BUTTON =
         "vsphere.core.dvPortgroup.editSettingsAction";
   public static final String ID_FUE_NETORK_FOLDER_NEW_PORT_GROUP_BUTTON_ID =
         ID_FUE_CREATEDVS_PORTGROUP_BUTTON + "/" + ID_BUTTON_BUTTON_NOSLASH;
   public static final String ID_FUE_ACTION_BAR_NEW_PORT_GROUP_BUTTON_ID =
         "vsphere.core.dvs.createPortgroupGlobal" + "/" + ID_BUTTON_BUTTON_NOSLASH;

   // vDS Port Mirroring
   public static final String ID_ACTION_EDIT_PORTMIRRORING_SESSION_VDS =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
               .append("editSpanSessionAction").toString();
   public static final String ID_ACTION_ADD_PORTMIRRORING_SESSION_VDS =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
               .append("addSpanSessionAction").toString();
   public static final String ID_ACTION_REMOVE_PORTMIRRORING_SESSION_VDS =
         new StringBuffer(EXTENSION_PREFIX).append(EXTENSION_ENTITY_DV_SWITCH)
               .append("removeSpanSessionAction").toString();
   public static final String ID_PORTMIRRORING_SESSION_LIST = "list";
   public static final String ID_ADD_PORTMIRRORING_DESCRIPRION_PAGE = "typePage";
   public static final String ID_INVOKE_PORTMIRRORING_WIZARD_PREFIX =
         "sessionList/list/dataGridToolbar/toolContainer/_SpanSessionList_ActionfwButtonBar1/";
   public static final String ID_EDIT_PORTMIRRORING_ENCAP_VLANID =
         "encapsulationVlanRow/encapsulationVlanIdStepper2";
   public static final String ID_EDIT_PORTMIRRORING_PRESERVE_VLAN =
         "preserveVlanIdCheckBox2";
   public static final String ID_EDIT_PORTMIRRORING_NORMALIO_STATUS =
         "normalTrafficAllowList";
   public static final String ID_EDIT_PORTMIRRORING_PACKETLEN_CHECKBOX =
         "mirroredPacketLengthCheckBox";
   public static final String ID_EDIT_PORTMIRRORING_PACKETLEN =
         "mirroredPacketLengthStepper";
   public static final String ID_PORTMIRRORING_ADD_WIZARD_TAB_NAVIGATOR =
         "wizardPageNavigator";
   public static final String ID_PORTMIRRORING_ADD_WIZARD_PROPERTY_PAGE_ID =
         "propertyPage";
   public static final String ID_PORTMIRRORING_ADD_WIZARD_SOURCE_PAGE_ID = "sourcePage";
   public static final String ID_PORTMIRRORING_ADD_WIZARD_DEST_PAGE_ID =
         "destinationPage";
   public static final String ID_PORTMIRRORING_ADD_WIZARD_SUMMARY_PAGE_ID =
         "summaryPage";
   public static final String ID_EDIT_PORTMIRRORING_SESSION_NAME = "nameInput";
   public static final String ID_EDIT_PORTMIRRORING_SESSION_STATUS = "statusList";
   public static final String ID_EDIT_PORTMIRRORING_SAMPLINGRATE = "samplingRateStepper";
   public static final String ID_EDIT_PORTMIRRORING_VLANID = "vlanId";
   public static final String ID_EDIT_PORTMIRRORING_IPADRESS = "ipAddressInput";
   public static final String ID_EDIT_PORTMIRRORING_ADD_VLANID_BUTTON =
         "vlanList/vlanIdList/dataGridToolbar/toolContainer/toolTip=Add VLAN ID/button";
   public static final String ID_EDIT_PORTMIRRORING_DELETE_VLANID_BUTTON =
         "vlanList/vlanIdList/dataGridToolbar/toolContainer/toolTip=Remove VLAN ID/button";
   public static final String ID_EDIT_PORTMIRRORING_ADD_IPADDRESS_BUTTON =
         "ipAddressList/dataGridToolbar/toolContainer/toolTip=Add IP address/button";
   public static final String ID_EDIT_PORTMIRRORING_DELETE_IPADDRESS_BUTTON =
         "ipAddressList/dataGridToolbar/toolContainer/toolTip=Remove IP address/button";
   public static final String ID_EDIT_PORTMIRRORING_SESSION_DESC = "descriptionText";
   public static final String ID_EDIT_PORTMIRRORING_ENCAP_VLAN_LEGACY =
         "encapsulationVlanIdCheckBox";
   public static final String ID_EDIT_PORTMIRRORING_ENCAP_VLANID_LEGACY =
         "encapsulationVlanIdStepper";
   public static final String ID_EDIT_PORTMIRRORING_PRESERVE_VLAN_LEGACY =
         "preserveVlanIdCheckBox";
   public static final String ID_PORTMIRRORING_ADD_DVPORTS_INFO_LABEL =
         "descriptionLabel";
   public static final String ID_PORTMIRRORING_INGRESS_DVPORTS_BUTTON = "ingressButton";
   public static final String ID_PORTMIRRORING_EGRESS_DVPORTS_BUTTON = "egressButton";
   public static final String ID_PORTMIRRORING_INGRESS_EGRESS_DVPORTS_BUTTON =
         "ingressEgressButton";
   public static final String ID_PORTMIRRORING_DELETE_DSTDVPORT_BUTTON = "removeButton";
   public static final String ID_PORTMIRRORING_DELETE_DVPORTS_BUTTON =
         "removePortButton";
   public static final String ID_PORTMIRRORING_ADD_DVPORTS_BUTTON = "addPortButton";
   public static final String ID_PORTMIRRORING_SELECT_DVPORTS_BUTTON =
         "addPortBySelectorButton";
   public static final String ID_PORTMIRRORING_SELECT_DVPORTS_LIST =
         "portBrowserList/datagrid";
   public static final String ID_PORTMIRRORING_ADD_DVPORTS_DIALOG = "portSelectionView";
   public static final String ID_PORTMIRRORING_ADD_DVPORTS_TEXT = "portList";
   public static final String ID_PORTMIRRORING_ADD_DVPORTS_OKBUTTON = "okButton";
   public static final String ID_PORTMIRRORING_ADD_VLANID_OR_IP_OKBUTTON = "okBtn";
   public static final String ID_PORTMIRRORING_ADD_VLANID_OR_IP_CANCELBUTTON =
         "cancelBtn";
   public static final String ID_PORTMIRRORING_ADD_DVPORTS_CANCELBUTTON = "cancelButton";
   public static final String ID_PORTMIRRORING_ADD_UPLINKS_LIST_ID =
         "availableUplinksList";
   public static final String ID_PORTMIRRORING_DELETE_UPLINKS_LIST_ID =
         "selectedUplinksList";
   public static final String ID_PORTMIRRORING_ADD_UPLINKS_LIST_FIELD = "key";
   public static final String ID_PORTMIRRORING_ADD_UPLINKS_LIST_ID_FOR_MIXEDDEST =
         "uplinkList";
   public static final String ID_PORTMIRRORING_ADD_UPLINKS_BUTTON_ID = "addUplinkBtn";
   public static final String ID_PORTMIRRORING_DELETE_UPLINKS_BUTTON_ID =
         "removeUplinkBtn";
   public static final String ID_PORTMIRRORING_ADD_UPLINKS_BUTTON_ID_FOR_MIXEDDEST =
         "addUplinkButton";
   public static final String ID_PORTMIRRORING_DETAILS_TAB = "detailsView";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_DESTINATION =
         "dvPortDestList";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_DSTPORTS_LIST =
         ID_PORTMIRRORING_DETAILS_TAB_DESTINATION + "/destinationList";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_LEGACYDST_LIST =
         "mixedDestList/destinationList";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_TYPE_LABEL = "typeLabel";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_STATUS_LABEL = "statusLabel";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_NORMALIO_LABEL =
         "normalTrafficAllowedLabel";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_PKTLEN_LABEL =
         "mirroredPacketLengthLabel";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_SAMPLING_LABEL =
         "samplingRateLabel";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_DESC_LABEL =
         "descriptionRow/descriptionLabel";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_SRCPORTS_LIST = "vlanIdList";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_SRCPORTSTYPE2_LIST =
         "portIdList";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_IP_LIST =
         "ipAddressList/ipAddressList";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_ENCAPVLAN_LABEL =
         "encapsulationVlanIdLabel";
   public static final String ID_PORTMIRRORING_DETAILS_TAB_KEEPORIGVLAN_LABEL =
         "stripOriginalVlanLabel";
   public static final String ID_ADD_PORTMIRRORING_SESSION_TYPE_RADIOBTN_PREFIX =
         "/label=";
   public static final String ID_EDIT_PORTMIRRORING_SESSION_CANCLE_BUTTON =
         "cancelButton";
   public static final String ID_EDIT_PORTMIRRORING_SESSION_OK_BUTTON = "okButton";
   public static final String ID_WIZARD_ADD_SESSION_PAGE_UNMARKED_STATUS = "normal";
   public static final String ID_WIZARD_ADD_SESSION_PAGE_MARKED_STATUS =
         "normalAndComplete";
   public static final String ID_WIZARD_ADD_SESSION_PAGE_CURRENT_STATUS = "currentState";

   //Host Memory CPU
   public static final String ID_COMBOBOX_HOST_BOOT_OPTIONS = "bootDeviceComboBox";
   public static final String ID_LABEL_HOST_MEMORY_TOTAL_MEMORY =
         "HostSystem:memoryConfig.totalPhysicalMemory";
   public static final String ID_LABEL_HOST_MEMORY_SYSTEM_MEMORY =
         "HostSystem:memoryConfig.systemPhysicalMemory";
   public static final String ID_LABEL_HOST_MEMORY_VM_PHYSICAL_MEMORY =
         "HostSystem:memoryConfig.vmPhysicalMemory";
   public static final String ID_LABEL_HOST_MEMORY_CONSOLE_MEMORY =
         "HostSystem:memoryConfig.consolePhysicalMemory";
   public static final String ID_TEXTINPUT_HOST_MEMORY_CONSOLE_MEMORY =
         "HostSystem:memoryConfig.consolePhysicalMemory";
   public static final String ID_BUTTON_HOST_MEMORY_EDIT =
         "Edit_HostSystem:memoryConfig";
   public static final String ID_PATH_ISCSI_STORAGE =
         "vsphere.core.host.manage.storage/";
   public static final String ID_CHECKBOX_OUTGOING = "outgoingCheckBox";
   public static final String ID_CHECKBOX_INCOMING = "incomingCheckBox";

   //NESTED HV
   public static final String ID_EDIT_HW_CPU_NESTEDHV_CHECK_BOX = "nestedHV";

   //CONFIG RULE
   public static final String ID_FILE_BROWSER_EXT_TYPE = "fileTypeCombo";
   public static final String ID_FILE_BROWSER_CANCEL_BUTTON = "buttonCancel";
   public static final String ID_NETWORK_NIC_TYPE = "type";
   public static final String ID_NETWORK_NIC_STATUS_GENERIC_CHECK_BOX = "status";
   public static final String ID_NETWORK_NIC_START_CONNECTED_GENERIC_CHECK_BOX =
         "startConnected";
   public static final String ID_ACTION_ASSIGN_TAG =
         "vsphere.core.tagging.assignTagAction";
   public static final String ID_ACTION_FILE_BROWSER =
         "vsphere.core.datastore.browseDatastoreAction";
   public static final String ID_ADD_PERMISSION = "Add...";
   public static final String ID_REMOVE_PERMISSION = "Remove...";
   public static final String ID_CHANGE_PERMISSION = "Change Role...";

   //Export List
   public static final String ID_TAB_NAVIGATOR_TABVIEWS = ID_TAB_NAVIGATOR
         + "/tabViews/";
   public static final String ID_EXPORT_LIST_PARENT =
         "<DATAGRID_ID>/dataGridStatusbar/dropDownBox";
   public static final String ID_EXPORT_LIST = ID_EXPORT_LIST_PARENT + "/export";
   public static final String ID_IMAGE_EXPORT_LIST = ID_TAB_NAVIGATOR_TABVIEWS
         + ID_EXPORT_LIST;

   public static final String ID_VC_VMTAB_EXPORT_LIST =
         "tabNavigator/tabViews/vmsForVCenter/list/dataGridStatusbar/dropDownBox/export";
   public static final String ID_VC_HOSTTAB_EXPORT_LIST =
         "tabNavigator/tabViews/hostsForVCenter/list/dataGridStatusbar/dropDownBox/export";
   public static final String ID_VC_HOSTPROFILESTAB_EXPORT_LIST =
         "tabNavigator/tabViews/hpsForVCenter/list/dataGridStatusbar/dropDownBox/export";
   public static final String ID_HOST_DATASTORE_EXPORT_LIST =
         "tabNavigator/tabViews/datastoresForStandaloneHost/list/dataGridStatusbar/container/dropDownBox/export";
   public static final String ID_HOST_VIRTUALMACHINE_EXPORT_LIST =
         "tabNavigator/tabViews/vmsForStandaloneHost/list/dataGridStatusbar/container/dropDownBox/export";
   public static final String ID_CLUSTER_DATASTORE_EXPORT_LIST =
         "tabNavigator/tabViews/dssForCluster/list/dataGridStatusbar/container/dropDownBox/export";
   public static final String ID_DC_HOST_EXPORT_LIST =
         "tabNavigator/tabViews/hostsForDatacenter/list/dataGridStatusbar/dropDownBox/export";
   public static final String ID_EXPORTLIST_SAVEBUTTON = "okButton";
   public static final String ID_EXPORTLIST_GENRATECSVREPORT_BUTTON =
         "generateReportButton";
   public static final String ID_EXPORTLIST_SELECTALLCLOUMN =
         "selectAllColumnsLinkButton";
   public static final String ID_RADIO_EXPORTLIST_ALL_ROWS =
         "_ListExportConfigForm_RadioButton1";
   public static final String ID_RADIO_EXPORTLIST_SELECTED_ROWS =
         "_ListExportConfigForm_RadioButton2";
   public static final String ID_PROGRESS_BUSYOVERLAY = "busyOverlay";

   // Export Events
   public static final String ID_BUTTON_EXPORTEVENTS = "export";
   public static final String ID_EXPORT_COLUMNS = "exportColumns";
   public static final String ID_RADIOBUTTON_SYSTEM_TYPE = "systemTypeEvent";
   public static final String ID_RADIOBUTTON_USER_TYPE = "userTypeEvent";
   public static final String ID_RADIOBUTTON_LIMITS = "numberSelectionRadioButton";
   public static final String ID_RADIOBUTTON_ALLEVENTS = "allEventRadioButton";
   public static final String ID_RADIOBUTTON_SELECTED_USERS = "selectedUsers";
   public static final String ID_RADIOBUTTON_ALL_USERS = "allUsers";
   public static final String ID_LINKBUTTON_UNSELECTALL_COLUMNS =
         "selectAllColumnsLinkButton";
   public static final String ID_BUTTON_GENERATE_CSV_REPORTS = "loadButton";
   public static final String ID_CHECKBOX_SEVERITY_INFORMATION = "infoSeverityCheckBox";
   public static final String ID_CHECKBOX_SEVERITY_ERROR = "errorSeverityCheckBox";
   public static final String ID_CHECKBOX_SEVERITY_WARNING = "warningSeverityCheckBox";
   public static final String ID_NUMERICSTEPPER_LAST_PERIOD =
         "lastPeriodNummericStepper";
   public static final String ID_NUMERICSTEPPER_LIMIT_EVENT =
         "limiteEventNumericStepper";
   public static final String ID_SELECTED_USERS_TEXT = "selectedUsersText";
   public static final String ID_LOADED_EVENT_RESULT = "loadedEventResult";
   public static final String ID_BUTTON_SEARCH_USERS = "searchUserButton";
   public static final String ID_BUTTON_ADD_USER = "addUserButton";
   public static final String ID_RADIOBUTTON_LAST_PERIOD = "lastPeriodRadioButton";
   public static final String ID_RADIOBUTTON_FROM_PERIOD = "currentPeriodRadioButton";

   public static final String ID_PROPERTIES_VIEW_TAB = "className=TabBar/text=";
   public static final String ID_MIGRATE_VMS_NETWORKING = "migrateVmNetworking";

   // Upgrade VMFS datastore
   public static final String ID_CHECKBOX_DATASTORE = "itemCheckBox";
   public static final String ID_UPGRADE_VMFS_DATASTORE_LABEL =
         "_UpgradeVmfsDatastoreForm_Text1";
   public static final String ID_UPGRADE_VMFS_VALIDATION_LABEL = "textLabel";
   public static final String ID_BUTTON_CLEAR_ALL_BUTTON = "clearAllButton";

   //Unique LUN
   public static final String ID_TEXTINPUT_RENAME = "nameInputField";

   //Unmount datastore
   public static final String ID_EXPLANATION_TEXT_LABEL = "explanationTextArea";
   public static final String ID_REQUIREMENTS_TEXT_LABEL = "requirementsTextArea";
   public static final String ID_SELECT_HOST_TEXT_LABEL = "selectHostsTextArea";
   public static final String ID_SELECT_HOSTS_LIST = "hostsList";

   // Licensing Management
   public static final String ID_LIC_SERVER_COMOBOX = "serverComboBox";
   public static final String ID_LICENSE_MANAGEMENT_VIEW = "vsphere.license.management";

   public static final String ID_TABNAVIGATOR = "tabNavigator";
   public static final String ID_LICENSE_KEYS_LIST =
         ID_LICENSE_MANAGEMENT_VIEW
               + "/"
               + ID_TABNAVIGATOR
               + "/vsphere.license.licenseKey.list.container/vsphere.license.licenseKey.list/licenseKeyList";
   public static final String ID_LICENSE_ADDLICENSEKEYS_BUTTON = ID_LICENSE_KEYS_LIST
         + "/vsphere.license.management.addLicenseKeysAction";
   public static final String ID_LICENSE_CLIPBOARD_BUTTON = ID_LICENSE_KEYS_LIST
         + "/copyToClipboard";
   public static final String ID_ADD_LICENSE_DIALOG = "tiwoDialog";
   public static final String ID_ADD_LICENSE_DIALOG_LICENSE_INPUT =
         ID_ADD_LICENSE_DIALOG + "/" + "licenseKeysText";
   public static final String ID_ADD_LICENSE_DIALOG_FINISH_BUTTON =
         ID_ADD_LICENSE_DIALOG + "/" + "finish";
   public static final String ID_ADD_LICENSE_DIALOG_BACK_BUTTON = ID_ADD_LICENSE_DIALOG
         + "/" + "back";
   public static final String ID_ADD_LICENSE_DIALOG_NEXT_BUTTON = ID_ADD_LICENSE_DIALOG
         + "/" + "next";
   public static final String ID_ADD_LICENSE_DIALOG_CANCEL_BUTTON =
         ID_ADD_LICENSE_DIALOG + "/" + "cancel";
   public static final String YES_NO_DIALOG = "YesNoDialog";
   public static final String YES_NO_DIALOG_YES_BUTTON = YES_NO_DIALOG + "/" + "Yes";
   public static final String YES_NO_DIALOG_NO_BUTTON = YES_NO_DIALOG + "/" + "NO";
   public static final String ID_LICENSE_KEYS_TAB = new StringBuffer(
         ID_LICENSE_MANAGEMENT_VIEW).append(ID_TABNAVIGATOR).append(".License Keys")
         .toString();
   public static final String OBJ_NAV_LICENSING_MANAGE_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(LICENSE_APPLICATION_NAVIGATOR).append(".management")
         .toString();
   public static final String OBJ_NAV_LICENSING_REPORTS_NODE_ITEM = new StringBuffer(
         TREE_NODE_ITEM).append(LICENSE_APPLICATION_NAVIGATOR).append(".reporting")
         .toString();
   public static final String ID_NOTIFICATION_BANNER_DETAILS_BUTTON = "detailsBtn";
   public static final String ID_LICENSE_EXPIRATION_NOTIFY_VIEW =
         "licenseExpirationNotifyView";
   public static final String ID_LICENSE_EXPIRATION_GENERAL_INFO =
         "licenseExpirationNotifyView/container/_LicenseExpirationNotifyView_Text1";
   public static final String ID_LICENSE_EXPIRATION_DETAIL_INFO =
         "licenseExpirationNotifyView/container/expiringVcs/textDisplay";
   public static final String ID_LICENSE_EXPIRATION_NOTIFY_CLOSEBUTTON =
         "licenseExpirationNotifyView/container/okBtn";
   public static final String ID_LICENSE_KYES_TABLE =
         "vsphere.license.management/extensionsView/tabNavigator/vsphere.license.licenseKey.list.container/vsphere.license.licenseKey.list/licenseKeyList";
   public static final String ID_ASSIGN_VC_LICENSE_EDIT_BUTTON =
         "btn_vsphere.license.management.vc.assignLicenseKeyAction";
   public static final String ID_ASSIGN_HOST_LICENSE_EDIT_BUTTON =
         "btn_vsphere.license.management.host.assignLicenseKeyAction";
   public static final String ID_DOWNGRADE_LICENSE_ERROR_PANEL =
         "SinglePageDialog_ID/errorPanel";
   public static final String ID_DOWNGRADE_LICENSE_ERROR_PANEL_TEXT1 =
         "_AssignLicenseKeyErrorPanel_Label1";
   public static final String ID_DOWNGRADE_LICENSE_ERROR_PANEL_TEXT2 =
         "_AssignLicenseKeyErrorPanel_Label2";
   public static final String ID_DOWNGRADE_LICENSE_ERROR_PANEL_DETAILBUTTON =
         "detailsBtn";
   public static final String ID_DOWNGRADE_LICENSE_DETAIL_VIEW_OK_BUTTON =
         "container/okBtn";
   public static final String ID_DOWNGRADE_LICENSE_DETAIL_VIEW_TEXT =
         "container/warningDetails";
   public static final String ID_ASSIGNLICENSE_WARNING_VIEW_TEXT1 =
         "_AssignLicenseKeyWarningDetailsView_Text1";
   public static final String ID_ASSIGNLICENSE_WARNING_VIEW_TEXT2 =
         "_AssignLicenseKeyWarningDetailsView_Text2";
   public static final String ID_DECODE_LICENSE_BUTTON =
         "SinglePageDialog_ID/details/decodeLicenseBtn";
   public static final String ID_ADDHOST_FAIL_IN_PROGRESS_LABEL =
         "vsphere.core.tiwo.sidebarView.portletChrome/tiwoListView/itemDescription";
   public static final String ID_ADDHOST_FAIL_IN_PROGRESS =
         "vsphere.core.tiwo.sidebarView.portletChrome/tiwoListView";
   public static final String ID_ADDHOST_CANCEL_BUTTON = "cancel";
   public static final String ID_ADDHOST_ERROR = "_container/_message";
   public static final String ID_LICENSEKEY_LIST_DATAGRID =
         "vsphere.license.licenseKey.list/licenseKeyList";
   public static final String ID_PRODUCTS_LIST_DATAGRID =
         "vsphere.license.product.list/productList";
   public static final String ID_VCINSTANCE_LIST_DATAGRID =
         "vsphere.license.vc.list/vcList";
   public static final String ID_HOSTS_LIST_DATAGRID =
         "vsphere.license.host.list/hostList";
   public static final String ID_SOLUTIONS_LIST_DATAGRID =
         "vsphere.license.solution.list/solutionList";
   public static final String ID_MOVE_HOST_TO_DIALOG = "objectSelectorDialog";
   public static final String ID_MOVE_HOST_CLUSTER_TAB =
         "objectSelectorDialog/_mainTabNavigator/filterObjectsContainer/filterViewToggleButtonBar/automationName=Clusters";
   public static final String ID_ADVANCE_DATAGRID_SELECT_CLUSTER =
         "objectSelectorDialog/_mainTabNavigator/filterObjectsContainer/ClusterComputeResource_filteringObjectsContainer/ClusterComputeResource_filterGridView";
   public static final String ID_RECENT_TASK_ERROR_LABEL =
         "vsphere.core.tasks.recentTasksView.portletChrome/recentTasksList/viewStack/statusHBox/statusLabel";
   public static final String ID_ASSIGN_NEWKEY_DURING_ADDHOST =
         "assignLicensePage/assignLicenseKeyForm/stateOptions";
   public static final String ID_ASSIGN_NEWKEY_ASSETUI =
         "SinglePageDialog_ID/stateOptions";
   public static final String ID_HOSTPROFILE_NAMEINPUT = "profileName";
   public static final String ID_DROPSHADOW_ERROR = "automationName=dropShadow";
   public static final String ID_HOSTPROFILE_ADVDATAGRID =
         "vsphere.core.hostProfiles.itemsView/tabViews/hostProfiles/list";
   public static final String ID_ATTACH_HP_BUTTON =
         "vsphere.core.hostprofile.attachAction/button";
   public static final String ID_COMPLIANCECHECK_HP_BUTTON =
         "vsphere.core.hp.checkComplianceAction/button";
   public static final String ID_CREATE_HOSTPROFILE_BUTTON =
         "vsphere.core.hp.createActionGlobal/button";
   public static final String ID_ATTACH_SELECTED_BUTTON =
         "selectionControl/buttons/attachSlected";
   public static final String ID_SELECTHOST_HP_ADVGRID =
         "selectionControl/availableViewContainer/availableView";
   public static final String ID_CANCELBUTTON_HOSTPROFILE = "cancel";
   public static final String ID_CREATEHOSTPROFILE_SELECTHOST =
         "HostSystem_filteringObjectsContainer/HostSystem_filterGridView";
   public static final String ID_CREATE_DVS_INPUT = "inputLabel/textDisplay";
   public static final String ID_DVS_ADDHOST_BUTTON =
         "vsphere.core.dvs.manageHostOnDvsAction/button";
   public static final String ID_DVS_INCOMPATIBLE_BUTTON = "incompatHostsImageButton";
   public static final String ID_ADVDATAGRID_INCOMP_HOSTLIST =
         "incompatibleHostsList/hostList";
   public static final String ID_COMPATIBLE_ISSUE_LABEL =
         "compatibilityIssues/propViewMessageText";
   public static final String ID_COMPATIBLE_ISSUE_CLOSEBUTTON = "btnClose";
   public static final String ID_COMPATIBLE_HOST_ADVGRID = "hostsList";
   public static final String ID_ADDNEWHOST_INTO_DVS_BUTTON =
         "hostsList/dataGridToolbar/toolContainer/newHostImageButton";
   public static final String ID_ASSIGNNEWKEY_TOASSET_OKBUTTON = "okButton";
   public static final String ID_ASSIGNNEWKEY_TOASSET_CANCELBUTTON = "cancelButton";
   public static final String ID_DVS_LIST_ADVDATAGRID =
         "vsphere.core.viDvs.itemsView/tabViews/DistributedVirtualSwitch/list";
   public static final String ID_SUMMARY_ASSIGNHOST_LICENSE_KEY_BUTTON =
         "vsphere.core.host.summary.licensingView.chrome/vsphere.core.host.summary.licensingView/assignLicenseBtn";
   public static final String ID_SUMMARY_ASSIGNVC_LICENSEKEY_BUTTON =
         "vsphere.core.vc.summary.licensingView.chromevsphere.core.vc.summary.licensingView/assignLicenseBtn";
   public static final String ID_HOSTNOTCOMPATIBLE_ALARM =
         "vsphere.opsmgmt.alarms.sidebarView.portletChrome/tabsViewStack/allTab/allAlarmsSideBarGrid/alarmName";
   public static final String ID_TRIGGERED_ALARM_TAB =
         "vsphere.opsmgmt.alarms.alarmsView.button";
   public static final String ID_TRIGGERED_ALARM_ADVGRID =
         "vsphere.opsmgmt.alarms.triggeredAlarmsView/alarmList";
   public static final String ID_ISSUES_TOC_TREE = "vsphere.core." + EXTENSION_ENTITY
         + "monitor.issues.container/tocTree";
   public static final String ID_LICENSE_INPUTTEXT = "licenseKeyText";
   public static final String LABEL_APP_REMOVE_STOREDDATA = "Remove Stored Data...";
   public static final String WORK_IN_PROGRESS_CHECKBOX =
         "removeStoredDataForm/removeDataContainer/com.vmware.usersettings.tiwo";
   public static final String ID_SVM_ERROR_LABEL =
         "migrationTypePage/viewStack/contents/changeDatastoreErrors/errWarningText";
   public static final String ID_EDIT_VM_BUTTON =
         "btn_vsphere.core.vm.provisioning.editAction";
   public static final String ID_MEMORY_SET_LABEL = "memoryControl/topRow/comboBox";
   public static final String ID_CPU_SET_LABEL = "cpuComboContainer/cpuCombo";
   public static final String ID_EDIT_VM_OKBUTTON = "okButton";
   public static final String ID_EDIT_VM_CANCELBUTTON = "cancelButton";
   public static final String ID_EDIT_DATASTORE_BUTTON =
         "btn_vsphere.core.datastore.configureStorageIOControl";
   public static final String ID_CHECKBOX_STORAGEIO =
         "SinglePageDialog_ID/StorageDialogId/enableIORMCheckBox";
   public static final String ID_VIRTUAL_NICS_ADVGRID =
         "vsphere.core.host.manage.settings.vnicView/vnicList";
   public static final String ID_EDIT_VNICS_BUTTON =
         "vsphere.core.host.network.editVnicSettingsAction/button";
   public static final String ID_VMOTION_CHECKBOX = "vMotionChk";
   public static final String ID_FT_CHECKBOX = "ftChk";
   public static final String ID_SPARKLIST_OF_CLUSTER =
         "vsphere.core.cluster.manage.settingsView/tocTreeContainer/tocTreeMainContainer/tocTree";
   public static final String ID_SPARKLIST_OF_VC_ALARM =
         "vsphere.core.folder.monitor.issues.wrapper/tocTreeContainer/tocTreeMainContainer/tocTree";
   public static final String ID_SPARKLIST_OF_VC_SETTING =
         "vsphere.core.folder.manage.settingsView/tocTreeContainer/tocTreeMainContainer/tocTree";
   public static final String ID_SPARKLIST_OF_HOST_SETTING =
         "vsphere.core.host.manage.settingsView/tocTreeContainer/tocTreeMainContainer/tocTree";
   public static final String ID_SPARKLIST_OF_HOST_ALARM =
         "vsphere.core.host.monitor.issues.wrapper/tocTreeContainer/tocTreeMainContainer/tocTree";
   public static final String ID_SPARKLIST_OF_HOST_NETWORKING =
         "vsphere.core.host.manage.networkingView.container/vsphere.core.host.manage.networkingView/tocTreeContainer/tocTreeMainContainer/tocTree";
   public static final String ID_PRODUCT_TEXT = "settingsComponent/productText";
   public static final String ID_LICENSEKEY_TEXT = "settingsComponent/licenseKeyText";
   public static final String ID_EXPIRATION_TEXT = "settingsComponent/expiresText";
   public static final String ID_FEATURES_TEXT = "settingsComponent/featuresText";
   public static final String ID_ADDON_TEXT = "settingsComponent/addOnsText";
   public static final String ID_TRIGGERED_ISSUES =
         "vsphere.core.host.monitor.issues.commonView/issueList";
   public static final String ID_TASK_DETAILS = "statusData";
   public static final String ID_LICENSEKEY_COLUMN_NAME = "License Key";
   public static final String ID_PRODUCT_COLUMN_NAME = "Name";
   public static final String ID_VC_COLUMN_NAME = "vCenter Server Instance";
   public static final String ID_HOST_COLUMN_NAME = "Host";
   public static final String ID_SOLUTIONS_COLUMN_NAME = "Solution";
   public static final String ID_REMOVELICENSEKEYS_BUTTON =
         "vsphere.license.management.removeLicenseKeyAction";
   public static final String ID_REMOVELICENSEKEYS_CONFIRMATION_DIALOG =
         "confirmationDialog";
   public static final String ID_ASSIGNVCLICENSEKEYS_BUTTON =
         "vsphere.license.management.vc.assignLicenseKeyAction";
   public static final String ID_ASSIGNHOSTLICENSEKEYS_BUTTON =
         "vsphere.license.management.host.assignLicenseKeyAction";
   public static final String ID_ADDNEWKEYCHECK_LIST = "readyToCompletePage/licenseKeys";
   public static final String ID_ASSIGNKEY_DETAILS_KEY =
         "tiwoDialog/details/licenseKeyHeader";
   public static final String ID_ASSIGNKEY_DETAILS_PRODUCT =
         "tiwoDialog/details/product";
   public static final String ID_ASSIGNKEY_DETAILS_CAPACITY =
         "tiwoDialog/details/capacityWithAvailable";
   public static final String ID_ASSIGNKEY_DETAILS_ADDITIONALCAPACITY =
         "tiwoDialog/details/additionalCapacityWithAvailable";
   public static final String ID_ASSIGNKEY_DETAILS_EXPIRES =
         "tiwoDialog/details/expires";
   public static final String ID_ASSIGNKEY_DETAILS_VRAMPERCPU =
         "tiwoDialog/details/ramPerCpuEntitlement";
   public static final String ID_LICENSE_ERROR_LABEL = "_errorLabel";
   public static final String ID_HA_SPARK_LIST_SEQUENCE = "2";
   public static final String ID_DRS_SPARK_LIST_SEQUENCE = "1";
   public static final String ID_LICENSE_SPARK_LIST_SEQUENCE = "5";
   public static final String ID_OTHER_IN_LICNESE_SPARK_LIST_SEQUENCE = "6";
   public static final String ID_HOST_TEAMING_HEALTH_SUMMARY_LABEL = "summaryItemLabel";
   public static final String ID_HOST_TEAMING_HEALTH_DETAILS_LABEL = "statusItemLabel";
   public static final String ID_VDS_MANAGE_SETTINGS_VIEW =
         "vsphere.core.dvs.manage.settings";

   // Licensing Ids
   public static final String ID_LICENSE_REPORTS_VIEW = "vsphere.license.reporting";
   public static final String ID_LICENSE_LICENSE_KEYS_LIST = new StringBuffer(
         ID_LICENSE_MANAGEMENT_VIEW).append("/vsphere.license.licenseKey.list")
         .toString();
   public static final String ID_LICENSE_ADD_LICENSE_KEYS_BUTTON = new StringBuffer(
         ID_LICENSE_LICENSE_KEYS_LIST).append(
         "/vsphere.license.management.addLicenseKeysAction").toString();
   public static final String ID_LICENSE_UPDATE_LABEL_BUTTON = new StringBuffer(
         ID_LICENSE_LICENSE_KEYS_LIST).append(
         "/vsphere.license.management.updateLicenseKeyAction").toString();
   public static final String ID_LICENSE_REMOVE_LICENSE_KEY_BUTTON = new StringBuffer(
         ID_LICENSE_LICENSE_KEYS_LIST).append(
         "/vsphere.license.management.removeLicenseKeyAction").toString();
   public static final String ID_ACTION_LICENSE_UPDATE_LABEL_BUTTON =
         "vsphere.license.management.updateLicenseKeyAction";
   public static final String ID_LICENSE_UPDATE_LABEL_TEXTINPUT =
         "tiwoDialog/SinglePageDialog_ID/labelText";
   public static final String ID_LICENSE_ADD_KEYS_ENTER_KEYS_PAGE = "addKeysPage";
   public static final String ID_LICENSE_ADD_KEYS_ENTER_KEYS_TEXTINPUT =
         "licenseKeysText";
   public static final String ID_LICENSE_ADD_KEYS_LABEL_TEXTINPUT =
         "optionalLabelText/textDisplay";
   public static final String ID_LICENSE_ADD_KEYS_READY_TO_COMPLETE_PAGE =
         "readyToCompletePage";
   public static final String ID_LICENSE_ADD_KEYS_READY_TO_COMPLETE_KEYS_LIST =
         "licenseKeys";
   public static final String ID_LICENSE_HOSTS_TAB_LIST = new StringBuffer(
         ID_LICENSE_MANAGEMENT_VIEW).append("/").append(ID_TAB_NAVIGATOR)
         .append("/vsphere.license.host.list/hostList").toString();
   public static final String ID_LICENSE_ASSIGN_LICENSE_KEY_BUTTON = new StringBuffer(
         ID_LICENSE_HOSTS_TAB_LIST).append(
         "/vsphere.license.management.host.assignLicenseKeyAction").toString();
   public static final String ID_LICENSE_ASSIGN_NEW_LICENSE_KEY_TEXTINPUT =
         "licenseKeyText";
   public static final String ID_LICENSE_ASSIGN_NEW_LICENSE_KEY_LABEL_TEXTINPUT =
         "labelText";
   public static final String ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW = "details";
   public static final String ID_LICENSE_ASSIGN_NEW_LICENSE_DECODE_BUTTON =
         new StringBuffer(ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW).append(
               "/decodeLicenseBtn").toString();
   public static final String ID_LICENSE_ASSIGN_LICENSE_EXISING_KEYS_LIST =
         "existingKeysList";
   public static final String ID_LICENSE_ASSIGN_EXISTING_OR_NEW_KEY_DROPDOWN =
         "stateOptions";
   public static final String ID_LICENSE_ASSIGN_LICENSE_HOST_TOTALCPUS_LABEL =
         new StringBuffer(ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW).append("/totalCpus")
               .toString();
   public static final String ID_LICENSE_ASSIGN_LICENSE_HOST_TOTALVMS_LABEL =
         new StringBuffer(ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW).append("/totalVms")
               .toString();
   public static final String ID_LICENSE_ASSIGN_LICENSE_HOST_MAXCORESPERCPU_LABEL =
         new StringBuffer(ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW).append(
               "/maxCoresPerCpu").toString();
   public static final String ID_LICENSE_ASSIGN_LICENSE_HOST_MAXCPUSPERASSET_LABEL =
         new StringBuffer(ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW).append(
               "/maxCpusPerAsset").toString();
   public static final String ID_LICENSE_ASSIGN_LICENSE_KEY_HEADER_LABEL =
         new StringBuffer(ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW).append(
               "/licenseKeyHeader").toString();
   public static final String ID_LICENSE_ASSIGN_LICENSE_LICENSE_PRODUCT_LABEL =
         new StringBuffer(ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW).append("/product")
               .toString();
   public static final String ID_LICENSE_ASSIGN_LICENSE_LICENSE_CAPACITY_LABEL =
         new StringBuffer(ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW).append(
               "/capacityWithAvailable").toString();
   public static final String ID_LICENSE_ASSIGN_LICENSE_LICENSE_SECOND_CAPACITY_LABEL =
         new StringBuffer(ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW).append(
               "/additionalCapacityWithAvailable").toString();
   public static final String ID_LICENSE_ASSIGN_LICENSE_LICENSE_EXPIRATION_LABEL =
         new StringBuffer(ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW).append("/expires")
               .toString();
   public static final String ID_LICENSE_ASSIGN_LICENSE_LICENSE_KEYLABEL_LABEL =
         new StringBuffer(ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW).append("/keyLabel")
               .toString();
   public static final String ID_LICENSE_ASSIGN_LICENSE_VRAM_PER_CPU_LABEL =
         new StringBuffer(ID_LICENSE_ASSIGN_LICENSE_DETAILS_VIEW).append(
               "/ramPerCpuEntitlement").toString();

   public static final String ID_LICENSE_KEY_SUMMARY_VIEW =
         "tabNavigator/vsphere.license.licenseKey.summary/vsphere.license.licenseKey.summaryView";
   public static final String ID_LICENSE_KEY_SUMMARY_HEADER = new StringBuffer(
         ID_LICENSE_KEY_SUMMARY_VIEW).append("/summaryHeader").toString();
   public static final String ID_LICENSE_KEY_SUMMARY_KEY_LABEL = "objName";
   public static final String ID_LICENSE_KEY_RELATED_VCS_LIST =
         "tabNavigator/vsphere.license.licenseKey.related/vcsForVcLicenseKey/vcList";
   public static final String ID_LICENSE_KEY_RELATED_HOSTS_LIST =
         "tabNavigator/vsphere.license.licenseKey.related/hostsForHostLicenseKey/hostList";
   public static final String ID_LICENSE_KEY_RELATED_PRODUCTS_LIST =
         "tabNavigator/vsphere.license.licenseKey.related/productsForLicenseKey/productList";

   // add by Juan
   public static final String ID_SUMMARY_SELECTED = "vgPropsAndBadges/objName";
   public static final String ID_SUMMARY_PRODUCT = "summary_productProp_valueLbl";
   public static final String ID_SUMMARY_VC = "summary_vcNameProp_valueLbl";
   public static final String ID_SUMMARY_KEYSNUMBER =
         "summary_numberOfKeysProp_valueLbl";
   public static final String ID_SUMMARY_COSTUNIT = "summary_costUnitProp_valueLbl";
   public static final String ID_SUMMARY_VRAMPERCPU =
         "summary_vRamEntitlementProp_valueLbl";
   public static final String ID_SUMMARY_LABEL = "summary_labelProp_valueLbl";
   public static final String ID_SUMMARY_EXPIRATION =
         "summary_expirationDateProp_valueLbl";
   public static final String ID_ASSET_SUMMARY_USAGE = "assignedText";
   public static final String ID_ASSET_SUMMARY_PRODUCT_LINK = "productLink";
   public static final String ID_ASSET_SUMMARY_EXPIRATION = "expirationDateText";
   public static final String ID_HOST_FILTER = "licenseFilterControl/filterOptions";

   public static final String ID_NO_PHYSICAL_ADAPTERS_ATTACHED_ERROR_TEXT_ID =
         "noPhysicalAdaptersAttached/errWarningText";
   public static final String ID_DATA_NOT_RETRIEVED_ERROR_TEXT_ID =
         "dataNotRetrievedErrorComponent/errWarningText";
   public static final String ID_BUTTON_EDIT_HOST_SETTINGS = "editHostHpsSettingsButton";
   //DME constants
   public static final String ID_RELATED_ITEMS_DATASTORE_ADVANCE_DATAGRID = "list";
   public static final String ID_RELATED_ITEMS_NETWORK_ADVANCE_DATAGRID = "list";
   public static final String ID_MOVE_TO_NAV_TREE = "navTreeView/navTree";
   public static final String ID_RELATED_ITEMS_TOP_LEVEL_OBJECTS_DATAGRID = "list";
   public static final String ID_TASK_CONSOLE_TITLE = "compositeDescriptionBox";
   //MemCPU
   public static final String ID_LABEL_VMPROVISIONING_MEMORY_SUMMARY_PAGE =
         "memoryValue";
   public static final String ID_LABEL_VMPROVISIONING_CPU_SUMMARY_PAGE = "cpuValue";
   public static final String ID_VMOPTIONS_BOOT_OPTIONS_UPDATE_WARNING_TEXT =
         "_BIOSPage_Text3";

   // FT specific constants
   public static final String ID_POWER_ON_MIXED_VM_CONFIRMATION_DIALOG =
         "powerOnVmConfirmationDialog";
   public static final String ID_POWER_ON_MIXED_VMS_WARNING_YES_BUTTON =
         ID_POWER_ON_MIXED_VM_CONFIRMATION_DIALOG + "/automationName=Yes";
   public static final String ID_POWER_ON_MIXED_VMS_WARNING_NO_BUTTON =
         ID_POWER_ON_MIXED_VM_CONFIRMATION_DIALOG + "/automationName=No";
   // SysResourceAlllocation:Host:CPUMemory config
   public static final String ID_BUTTON_SYSCONFIG_TAB_SETTINGS =
         "vsphere.core.host.manage.settingsView.button";
   public static final String ID_BUTTON_SYS_RESERVATION_EDIT = "systemReservationEdit";
   public static final String ID_BUTTON_SYS_RESERVATION_ADVANCED_EDIT = "detailEdit";
   public static final String ID_BUTTON_SYS_ALLOCATION_LIST = new StringBuffer(
         IDConstants.EXTENSION_PREFIX).append(IDConstants.EXTENSION_ENTITY_HOST)
         .append("manage.").append("settingsView").append("/tocTree").toString();
   public static final String ID_LIST_ITEM_HOST_MANAGE_SETTINGS_SYSTEM_RESOURCE_ALLOCATION =
         new StringBuffer(ID_BUTTON_SYS_ALLOCATION_LIST).append(
               "/automationName=System Resource Allocation").toString();
   public static final String ID_LABEL_SYS_ALLOCATION_SIMPLE_LINK = "simpleLink";
   public static final String ID_LABEL_SYS_ALLOCATION_ADVANCED_LINK = "advancedLink";
   public static final String ID_LABEL_SYS_ALLOCATION_ADVANCED_MEMORY =
         "_SimpleResourceAllocationDetailsView_LabelEx1";
   public static final String ID_LABEL_SYS_ALLOCATION_ADVANCED_CPU =
         "_SimpleResourceAllocationDetailsView_LabelEx2";

   // CPU related ID's
   public static final String ID_TEXT_SYS_ALLOCATION_CPU_SHARES_LEVEL =
         "cpuConfigControl/shares/sharesControl/levels";
   public static final String ID_TEXT_SYS_ALLOCATION_CPU_SHARES_VALUE =
         "cpuConfigControl/shares/sharesControl/numShares";
   public static final String ID_TEXT_SYS_ALLOCATION_CPU_RESERVATION =
         "cpuConfigControl/reservation/reservationCombo/topRow/comboBox";
   public static final String ID_TEXT_SYS_ALLOCATION_CPU_RESERVATION_UNIT =
         "cpuConfigControl/reservation/reservationCombo/topRow/units";
   public static final String ID_CHECKBOX_SYS_ALLOCATION_CPU_RESERVATION_TYPE =
         "cpuConfigControl/expandableReservation/expReservationCheck";
   public static final String ID_LABEL_SYS_ALLOCATION_CPU_RESERVATION_MAX =
         "cpuConfigControl/reservation/reservationCombo/label2";
   public static final String ID_COMBO_SYS_ALLOCATION_CPU_LIMIT =
         "cpuConfigControl/limit/limitCombo/topRow/comboBox";
   public static final String ID_COMBO_SYS_ALLOCATION_CPU_LIMIT_UNIT =
         "cpuConfigControl/limit/limitCombo/topRow/units";
   public static final String ID_LABEL_SYS_ALLOCATION_CPU_LIMIT_MAX =
         "cpuConfigControl/limit/limitCombo/label2";
   public static final String ID_LABEL_SYS_ALLOCATION_CPU_ADVANCED_LINK_SHARE_VALUE =
         "_AdvancedResourceAllocationDetailsView_LabelEx1";
   public static final String ID_LABEL_SYS_ALLOCATION_CPU_ADVANCED_LINK_RESERVATION_VALUE =
         "_AdvancedResourceAllocationDetailsView_LabelEx2";
   public static final String ID_LABEL_SYS_ALLOCATION_CPU_ADVANCED_LINK_LIMIT_VALUE =
         "_AdvancedResourceAllocationDetailsView_LabelEx3";
   public static final String ID_CHECKBOX_SYS_ALLOCATION_CPU_ADVANCED_LINK_EXPANDABLE_VALUE =
         "_AdvancedResourceAllocationDetailsView_CheckBox1";
   public static final String ID_CHECKBOX_SYS_ALLOCATION_CPU_ADVANCED_LINK_LIMIT_VALUE =
         "_AdvancedResourceAllocationDetailsView_LabelEx3";

   // Memory related ID's
   public static final String ID_TEXT_SYS_ALLOCATION_MEMORY_SHARES_LEVEL =
         "memoryConfigControl/shares/sharesControl/levels";
   public static final String ID_TEXT_SYS_ALLOCATION_MEMORY_SHARES_VALUE =
         "memoryConfigControl/shares/sharesControl/numShares";
   public static final String ID_TEXT_SYS_ALLOCATION_MEMORY_RESERVATION =
         "memoryConfigControl/reservation/reservationCombo/topRow/comboBox";
   public static final String ID_TEXT_SYS_ALLOCATION_MEMORY_RESERVATION_UNIT =
         "memoryConfigControl/reservation/reservationCombo/topRow/units";
   public static final String ID_CHECKBOX_SYS_ALLOCATION_MEMORY_RESERVATION_TYPE =
         "memoryConfigControl/expandableReservation/expReservationCheck";
   public static final String ID_LABEL_SYS_ALLOCATION_MEMORY_RESERVATION_MAX =
         "memoryConfigControl/reservation/reservationCombo/label2";
   public static final String ID_COMBO_SYS_ALLOCATION_MEMORY_LIMIT =
         "memoryConfigControl/limit/limitCombo/topRow/comboBox";
   public static final String ID_COMBO_SYS_ALLOCATION_MEMORY_LIMIT_UNIT =
         "memoryConfigControl/limit/limitCombo/topRow/units";
   public static final String ID_LABEL_SYS_ALLOCATION_MEMORY_LIMIT_MAX =
         "memoryConfigControl/limit/limitCombo/label2";
   public static final String ID_BUTTON_SYS_ALLOCATION_MEMORY_CPU_EDIT_SETTING_OK =
         "okButton";
   public static final String ID_BUTTON_SYS_ALLOCATION_MEMORY_CPU_EDIT_SETTING_CANCEL =
         "cancelButton";
   public static final String ID_LABEL_SYS_ALLOCATION_MEMORY_ADVANCED_LINK_SHARE_VALUE =
         "_AdvancedResourceAllocationDetailsView_LabelEx4";
   public static final String ID_LABEL_SYS_ALLOCATION_MEMORY_ADVANCED_LINK_RESERVATION_VALUE =
         "_AdvancedResourceAllocationDetailsView_LabelEx5";
   public static final String ID_LABEL_SYS_ALLOCATION_MEMORY_ADVANCED_LINK_LIMIT_VALUE =
         "_AdvancedResourceAllocationDetailsView_LabelEx6";
   public static final String ID_CHECKBOX_SYS_ALLOCATION_MEMORY_ADVANCED_LINK_EXPANDABLE_VALUE =
         "_AdvancedResourceAllocationDetailsView_CheckBox2";
   public static final String ID_CHECKBOX_SYS_ALLOCATION_MEMORY_ADVANCED_LINK_LIMIT_VALUE =
         "_AdvancedResourceAllocationDetailsView_LabelEx6";
   public static final String ID_DATAGRID_SYS_ALLOCATION_RESOURCEPOOL_TREE =
         "resourcesGrid";
   public static final String ID_DATAPROVIDER_SOURCE = "dataProvider.";
   public static final String ID_DATAPROVIDER_LABEL = ".label";
   public static final String ID_DATAPROVIDER_DATA = ".data";
   public static final String DATAPROVIDER_NAME_PAGE_TITLE = ".pageInfo.title";
   public static final String ID_CHECKBOX_INCLUDE_SERVER_LICE_INFO =
         "includeServerLicenseInfo";
   public static final String ID_RADIO_WINDOWS_LIC_PER_SEAT =
         "windowsLicensePage/perSeat";
   public static final String ID_RADIO_WINDOWS_LIC_PER_SERVER =
         "windowsLicensePage/perServer";
   public static final String ID_TEXT_MAX_CONNECTION =
         "windowsLicensePage/maxConnectionsPerServer";
   public static final String ID_DATAGRID_LOCAL_TIMING = "hwClockSelector";
   public static final String ID_LIST_AREA_INDEX = "areaSelector";
   public static final String ID_IPV4ADDRESSINPUT_PRIMARY_DNS = "primaryDns";
   public static final String ID_IPV4ADDRESSINPUT_SECONDARY_DNS = "secondaryDns";
   public static final String ID_IPV4ADDRESSINPUT_TETRIARY_DNS = "tetriaryDns";
   public static final String ID_LABEL_NW_GENERAL_PROPERTIES = "step_properties";
   public static final String ID_LABEL_NW_DNS_PROPERTIES = "step_dnsProperties";
   public static final String ID_LABEL_NW_WINS_PROPERTIES = "step_winsProperties";
   public static final String ID_RADIO_NW_USE_DHCP = "useDHCP";
   public static final String ID_RADIO_NW_RPOMPT_USER = "promptUser";
   public static final String ID_RADIO_NW_USE_VC_GENERATED_IP =
         "useGenerated/generateIPAddress";
   public static final String ID_RADIO_NW_USE_FOLLOWING_IP_ADDRESS = "useIPAddress";
   public static final String ID_RADIO_NW_USE_DHCP_DNS_ENTRY = "useDhcpDns";
   public static final String ID_RADIO_NW_USE_MANUAL_DNS_ENTRY = "useManualDns";
   public static final String ID_TEXT_ALTERNATE_GATEWAY = "alternateGateway";
   public static final String ID_TEXT_PREFERRED_DNS_SERVER = "preferredDnsServer";
   public static final String ID_TEXT_ALTERNATE_DNS_SERVER = "alternateDnsServer";
   public static final String ID_UI_COMPONENT_ANCHORED_DIALOG =
         "className=AnchoredDialog";
   public static final String ID_BUTTON_DNS_PROP_ADD = "dnsProperties/addButton";
   public static final String ID_BUTTON_DNS_PROP_DELETE = "dnsProperties/deleteButton";
   public static final String ID_LABEL_NW_ALTERNATE_GATEWAY_VALUE =
         "alternateGatewayLabel";
   public static final String ID_UI_COMPONENT_DNS_SERVER_FIELDS =
         "_NetworkGeneralProperties_VGroup6";
   public static final String ID_LIST_DNS_SUFFIX =
         "dnsProperties/dnsSuffixList/entryList";
   public static final String ID_TEXTINPUT_VM_PROV_SEARCH =
         "selectResourcePage/searchInput";
   public static final String ID_SCSICONTROLLER_CHANGETYPE = "changeTypeButton";
   public static final String ID_SCSICONTROLLER_DONT_CHANGETYPE = "dontChangeTypeButton";
   public static final String LIST_DNS_SUFFIX_ENTRYLIST =
         "dnsProperties/dnsSuffixList/entryList";
   public static final String LIST_COMMAND_ENTRYLIST = "entryList";
   public static final String LIST_NETWORKSETTINGS_LIST = "networkSettings";
   public static final String TEXTINPUT_NET_BIOS_NAME =
         "computerNameBox/netBiosNameTextInput";
   public static final String TEXT_USE_GENERATED_NAME = "useGenerated/textDisplay";

   //multiVC
   public static final String ID_COMBOBOX_STATISTIC_INTERVAL = "uicEditor";
   public static final String ID_PAGE_VC_SETTINGS = "wizardPageNavigator";
   public static final String ID_TEXT_INPUT_STATISTICS_INTERVAL_SAVE_DURATION =
         "uicEditor";
   public static final String ID_NUM_STEPPER_ESTIMATED_NO_OF_VMS =
         "estimatedNumberOfVms";
   public static final String ID_NUM_STEPPER_ESTIMATED_NO_OF_HOSTS =
         "estimatedNumberOfHosts";
   public static final String ID_LIST_VCENTER_SERVERS = "selectedSetDataGrid";
   public static final String ID_LABEL_STATISTICS_SETTINGS =
         "step_vcStatisticsConfigPage";
   public static final String ID_TEXT_RUNTIME_SETTINGS_UNIQUE_ID =
         "Folder:vcServerSettings.runtimeSettings.instance.id";
   public static final String ID_TEXT_RUNTIME_SETTINGS_IP_ADDRESS =
         "Folder:vcServerSettings.runtimeSettings.VirtualCenter.ManagedIP";
   public static final String ID_TEXT_RUNTIME_SETTINGS_VC_NAME =
         "Folder:vcServerSettings.runtimeSettings.VirtualCenter.InstanceName";
   public static final String ID_TEXT_AD_SETTINGS_TIMEOUT =
         "Folder:vcServerSettings.activeDirectorySettings.ads.timeout";
   public static final String ID_CHECKBOX_AD_SETTINGS_QUERY_LIMIT =
         "Folder:vcServerSettings.activeDirectorySettings.ads.maxFetchEnabled";
   public static final String ID_TEXT_AD_SETTINGS_QUERY_LIMIT =
         "Folder:vcServerSettings.activeDirectorySettings.ads.maxFetch";
   public static final String ID_CHECKBOX_AD_SETTINGS_VALIDATION_PERIOD =
         "Folder:vcServerSettings.activeDirectorySettings.ads.checkIntervalEnabled";
   public static final String ID_TEXT_AD_SETTINGS_VALIDATION_PERIOD =
         "Folder:vcServerSettings.activeDirectorySettings.ads.checkInterval";
   public static final String ID_LABEL_RUNTIME_SETTINGS = "step_vcRuntimeConfigPage";
   public static final String ID_LABEL_ACTIVE_DIRECTORY_SETTINGS =
         "step_vcActiveDirectoryConfigPage";
   public static final String ID_LABEL_MAIL_SETTINGS = "step_vcMailConfigPage";
   public static final String ID_LABEL_SSL_SETTINGS = "step_vcSslConfigPage";
   public static final String ID_TEXT_MAIL_SETTINGS_SERVER =
         "Folder:vcServerSettings.mailSettings.mail.smtp.server";
   public static final String ID_TEXT_MAIL_SETTINGS_SENDER =
         "Folder:vcServerSettings.mailSettings.mail.sender";
   public static final String ID_LABEL_SNMP_RECV_SETTINGS = "step_vcSnmpConfigPage";
   public static final String ID_TEXT_SNMP_SETTINGS_RECEIVER =
         "Folder:vcServerSettings.snmpSettings.snmp.receiver.1.name";
   public static final String ID_LABEL_PORTS_SETTINGS = "step_vcPortsConfigPage";
   public static final String ID_LABEL_PORTS_PAGE_TITLE = "vcPortsConfigPage";
   public static final String ID_LABEL_TIMEOUT_SETTINGS = "step_vcTimeoutConfigPage";
   public static final String ID_TEXT_TIMEOUT_SETTINGS_NORMAL_OP =
         "Folder:vcServerSettings.timeoutSettings.client.timeout.normal";
   public static final String ID_TEXT_TIMEOUT_SETTINGS_LONG_OP =
         "Folder:vcServerSettings.timeoutSettings.client.timeout.long";
   public static final String ID_LABEL_LOG_SETTINGS = "step_vcLogConfigPage";
   public static final String ID_COMBOBOX_LOGGING_OPTIONS =
         "Folder:vcServerSettings.logSettings.log.level";
   public static final String ID_LABEL_DB_SETTINGS = "step_vcDatabaseConfigPage";
   public static final String ID_TEXT_DB_SETTINGS_MAX_CONNECTIONS =
         "Folder:vcServerSettings.databaseSettings.VirtualCenter.MaxDBConnection";
   public static final String ID_CHECKBOX_DB_SETTINGS_TASK_CLEANUP =
         "Folder:vcServerSettings.databaseSettings.task.maxAgeEnabled";
   public static final String ID_TEXT_DB_SETTINGS_TASK_RETENTION =
         "Folder:vcServerSettings.databaseSettings.task.maxAge";
   public static final String ID_CHECKBOX_DB_SETTINGS_EVENT_CLEANUP =
         "Folder:vcServerSettings.databaseSettings.event.maxAgeEnabled";
   public static final String ID_TEXT_DB_SETTINGS_EVENT_RETENTION =
         "Folder:vcServerSettings.databaseSettings.event.maxAge";
   public static final String ID_LIST_ADVANCED_SETTINGS =
         "Folder:vcServerSettings.advancedSettings";
   public static final String ID_BUTTON_VC_CONFIG_ADV_SET_EDIT =
         "btn_vsphere.core.folder.editVcAdvancedSettingsAction";
   public static final String ID_LIST_VCENTER_VMS = "vmsForVCenter/list";
   public static final String ID_GRID_PERMISSION = "permissionGrid";
   public static final String ID_TAB_PERMISSION = "Permissions";
   public static final String ID_TAB_SYSTEM_LOGS = "System Logs";
   public static final String ID_LABEL_MANGE_SETTINGS_LICENSING =
         "vsphere.core.folder.manage.settings.licensing";
   public static final String ID_TAB_EVENTS = "Events";
   public static final String ID_TAB_SCHEDULED_TASKS = "Scheduled Tasks";
   public static final String ID_LABEL_SCHEDULED_TASKS = "howToCreateLabel";
   public static final String ID_TEXT_SSL_CERT_INFO = "verifySslCertificatesInfo";
   public static final String ID_DIALOG_MSG_ALERT = "AlertDialog";
   public static final String ID_TEXT_MSG = "messageTextArea";
   public static final String ID_LABEL_VC_COUNT =
         "TreeNodeItem_vsphere.core.navigator.vilist.Folder_isRootFolder";
   public static final String ID_LABEL_DC_COUNT =
         "TreeNodeItem_vsphere.core.navigator.vilist.Datacenter_";
   public static final String ID_LABEL_CLUSTER_COUNT =
         "TreeNodeItem_vsphere.core.navigator.vilist.ClusterComputeResource_";
   public static final String ID_TAB_LINKED_SERVER = "Linked vCenter Server systems";
   public static final String ID_GRID_LINKED_SERVER = "linkedVCentersForVCenter";
   public static final String ID_LABEL_RTC_NAME = "propViewValueText";


   //Login
   public static final String ID_USER_MENU_CHANGE_PASSWD_OPTION = "Change Password...";
   public static final String ID_USER_MENU_CHANGE_PASSWD_AUTOMATION_NAME =
         "automationName=Change Password...";
   public static final String ID_CHECKBOX_WIN_AUTHENTICATION = "sspiCheckBox";

   public static final String ID_DATACENTER_GS_TAB =
         "vsphere.core.datacenter.gettingStarted.container";
   public static final String ID_DATACENTER_SUMMARY_TAB =
         "vsphere.core.datacenter.summary.container";
   public static final String ID_DATACENTER_MANAGE_TAB =
         "vsphere.core.datacenter.manage.container";
   public static final String ID_DATACENTER_MONITOR_TAB =
         "vsphere.core.datacenter.monitor.container";
   public static final String ID_DATACENTER_RELATED_OBJECTS_TAB =
         "vsphere.core.datacenter.related.container";

   // Navigation
   public static final String ID_MONITOR_CONTAINER = "vsphere.core." + EXTENSION_ENTITY
         + ".monitor";
   public static final String ID_MANAGE_CONTAINER = "vsphere.core." + EXTENSION_ENTITY
         + ".manage";
   public static final String ID_RELATED_OBJECTS_CONTAINER = "vsphere.core."
         + EXTENSION_ENTITY + ".related";

}
