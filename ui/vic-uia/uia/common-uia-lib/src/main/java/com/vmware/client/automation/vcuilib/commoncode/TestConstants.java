/*
 * ************************************************************************
 *
 * Copyright 2009 VMware, Inc. All rights reserved. -- VMware Confidential
 *
 * ************************************************************************
 */
package com.vmware.client.automation.vcuilib.commoncode;

import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_CLUSTER_SETTINGS_EVC_EDIT_AMD_RADIO;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_CLUSTER_SETTINGS_EVC_EDIT_DISABLE_EVC_RADIO;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_CLUSTER_SETTINGS_EVC_EDIT_INTEL_RADIO;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_ADMINISTRATION_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_LICENSES_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_LICENSE_REPORTS_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_MANAGE_MONITOR_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_ROLE_MANAGER_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_RULES_PROFILE_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_SAVEDSEARCH_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_SEARCH_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_SOLUTION_PLUGIN_MANAGER_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_SSO_CONFIGURATION_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_SSO_USERS_GROUPS_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_TAG_MANAGER_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_TASKS_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.OBJ_NAV_VIRTUAL_INFRASTRUCTURE_NODE_ITEM;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.MAIN_NAVIGATION_TABS.MAIN_TABS_GETTING_STARTED;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.MAIN_NAVIGATION_TABS.MAIN_TABS_MANAGE;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.MAIN_NAVIGATION_TABS.MAIN_TABS_MONITOR;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.MAIN_NAVIGATION_TABS.MAIN_TABS_RELATED_ITEMS;
import static com.vmware.client.automation.vcuilib.commoncode.TestConstants.MAIN_NAVIGATION_TABS.MAIN_TABS_SUMMARY;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Define the constants to be used across the classes.
 *
 * NOTE: this class is a copy of the one from VCUI-QE-LIB
 */
public class TestConstants extends TestConstantsKey {

   private static final Logger logger = LoggerFactory.getLogger(TestConstants.class);

   /*
    *
    * ENUM SECTION All ENUMs go here
    */

   /**
    * Use this enum for any timeouts instead of using below constants: public
    * static final int DEFAULT_TIMEOUT_IMMEDIATE_INT_VALUE = 500; public static
    * final long DEFAULT_TIMEOUT_IMMEDIATE_LONG_VALUE = new
    * Long(DEFAULT_TIMEOUT_IMMEDIATE).longValue(); // one second public static
    * final int DEFAULT_TIMEOUT_ONE_SECOND_INT_VALUE = 1000; public static final
    * long DEFAULT_TIMEOUT_ONE_SECOND_LONG_VALUE = new
    * Long(DEFAULT_TIMEOUT_ONE_SECOND).longValue(); // five seconds public
    * static final int DEFAULT_TIMEOUT_FIVE_SECOND_INT_VALUE = 5000; // ten
    * seconds public static final long DEFAULT_TIMEOUT_TEN_SECONDS_LONG_VALUE =
    * new Long(DEFAULT_TIMEOUT_TEN_SECONDS).longValue(); public static final int
    * DEFAULT_WAIT_ITERATIONS_THREE_MILLION = 30000000; As of now I don't see
    * many uses of int constants and you can always use long instead of int, so
    * use of below implementation will be easy to maintain. If you need string
    * equivalents, place them in properties files, like for below constants.
    * public static final String DEFAULT_TIMEOUT_IMMEDIATE = "1"; public static
    * final String DEFAULT_TIMEOUT_ONE_SECOND = "1000"; public static final
    * String DEFAULT_TIMEOUT_TEN_SECONDS = "10000";
    *
    * @author Administrator
    */

   public enum Timeout {
      NO_TIMEOUT(0), ONE_MILLIS(1), ONE_HUNDRED_MILLIS(100), FIVE_HUNDRED_MILLIS(500), ONE_SECOND(
            1000), TWO_SECONDS(2000), THREE_SECONDS(3000), FIVE_SECONDS(5 * 1000), TEN_SECONDS(
            10 * 1000), ONE_MINUTE(60 * 1000), TWO_MINUTES(2 * 60 * 1000), THREE_MINUTES(
            3 * 60 * 1000), FIVE_MINUTES(5 * 60 * 1000), TEN_MINUTES(10 * 60 * 1000), THIRTY_SECONDS(
            30 * 1000);
      /**
       * Duration of the timeout in millis
       */
      private final long duration;

      private Timeout(long duration) {
         this.duration = duration;
      }

      /**
       * Returns the duration of the timeout
       *
       * @return
       */
      public long getDuration() {
         return duration;
      }

      /**
       * It will make the current thread to suspend execution
       */
      public void consume() {
         if (duration > 0) {
            try {
               logger.debug("Starting the timeout: " + this);
               Thread.sleep(duration);
               logger.debug("Finished the timeout: " + this);
            } catch (InterruptedException e) {
               logger.debug("The thread was interruped while performing the timeout: "
                     + this);
            }
         }
      }

      @Override
      public String toString() {
         StringBuilder output = new StringBuilder();
         output.append("Timeout: ").append(duration).append(" millis");
         return output.toString();
      }

   }

   public static enum SEARCH_MODIFIERS {
      BEGINS_WITH, CONTAINS, ENDS_WITH, IS
   };

   public static enum STANDARD_SWITCH_OBJECT_TYPE {
      VSWITCH, PORTGROUP, VNIC, PNIC
   }

   public static enum NUM_PORT_LIST {
      EIGHT("8"), TWENTY_FOUR("24"), FIFTY_SIX("56"), ONE_HUNDRED_TWENTY("120"), TWO_HUNDRED_FORTY_EIGHT(
            "248"), FIVE_HUNDRES_FOUR("504"), ONE_THOUSAND_SIXTEEN("1016"), TWO_THOUSAND_FORTY(
            "2040"), FOUR_THOUSAND_EIGHTY_EIGHT("4088");

      private final String numPort;

      NUM_PORT_LIST(String numPortValue) {
         numPort = numPortValue;
      }

      public String getNumPortValue() {
         return numPort;
      }
   }

   public static enum SortOrder {
      ASCENDING(false), DESCENDING(true);

      private final boolean sortOrder;

      SortOrder(boolean ascending) {
         sortOrder = ascending;
      }

      public boolean getSortOrder() {
         return sortOrder;
      }
   }

   public static enum CheckboxValue {
      CHECKED("checked"), UNCHECKED("unchecked");

      private final String checkBoxValue;

      CheckboxValue(String checkValue) {
         checkBoxValue = checkValue;
      }

      public String getCheckboxValue() {
         return checkBoxValue;
      }
   }

   public static enum HOST_NETWORKING_PAGE {
      GENERAL("0"), DNS_AND_ROUTING("1"), PHYSICAL_NETWORK_ADAPTERS("2"), VIRTUAL_NETWORK_ADAPTERS(
            "3"), VIRTUAL_SWITCHES("4");

      private final String pageIndex;

      HOST_NETWORKING_PAGE(String pageValue) {
         pageIndex = pageValue;
      }

      public String getPageIndex() {
         return pageIndex;
      }
   }

   // Resource Pool Reservation Drop down item position
   public static enum RESERVATION_DROPDOWN_POS {
      POS_RES_CURRENT("0"), POS_RES_MINIMUM("1"), POS_RES_MAXIMUM("2");
      private final String resPos;

      RESERVATION_DROPDOWN_POS(String rPos) {
         resPos = rPos;
      }

      public String getResPosition() {
         return resPos;
      }
   }

   // Resource Pool Limit Drop down item position
   public static enum LIMIT_DROPDOWN_POS {
      POS_LIM_CURRENT("0"), POS_LIM_MINIMUM("1"), POS_LIM_MAXIMUM("2"), POS_LIM_UNLIMITED(
            "3");
      private final String limPos;

      LIMIT_DROPDOWN_POS(String lPos) {
         limPos = lPos;
      }

      public String getLimPosition() {
         return limPos;
      }
   }

   public static enum CONJOINER {
      any, all
   }

   public static enum CRITERIA_COMPARATOR {
      Op_Contains("contains"), Op_Equals("is"), Op_Unequals("is not"), Op_GreaterThan(
            "is greater than"), Op_LessThan("is less than");

      private final String criteriaValue;

      CRITERIA_COMPARATOR(String cValue) {
         criteriaValue = cValue;
      }

      public String getCriteriaValue() {
         return criteriaValue;
      }
   }

   /**
    * Enum defined for VM's power states
    */
   public static enum VM_POWER_STATE {
      poweredOn, poweredOff, reset, suspended, rebootGuest, standbyGuest
   };

   public static enum INVOKE_METHODS {
      IM_CONTEXT_MENU, IM_ACTION_BAR, IM_MORE_ACTIONS, IM_BUTTON_ON_PORTLET, IM_GETTING_STARTED_LINK, IM_VM_TAB, IM_NETWORK_FOLDER, IM_TOOL_BAR
   }

   public static enum GUEST_MEMORY_PORTLET_VALUES {
      ACTIVE_GUEST_MEMORY, PRIVATE_MEMORY, SHARED_MEMORY, BALLOONED_MEMORY, COMPRESSED_MEMORY, SWAPPED_MEMORY, UNACCESSED_MEMORY
   }

   public static enum HOST_CPU_PORTLET_VALUES {
      HOST_CPU_CONSUMED, HOST_CPU_ACTIVE
   }

   public static enum HOST_MEMORY_PORTLET_VALUES {
      HOST_MEMORY_CONSUMED, HOST_MEMORY_OVERHEAD
   }

   public static enum HOST_TYPE {
      STANDALONE, CLUSTERED, LEGACY
   }

   public static enum HOST_STATE {
      HS_CONNECTED(STATE_CONNECTED), HS_DISCONNECTED(STATE_DISCONNECTED), HS_MAINTENANCE_MODE(
            STATE_MAINTENANCE_MODE), HS_NOTRESPONDING(STATE_NOTRESPONDING);
      private final String name;

      HOST_STATE(String value) {
         name = value;
      }

      public String getValue() {
         return name;
      }
   }

   public static enum VDS_STATUS {
      UP(ENTITY_STATUS_UP), OUT_OF_SYNC(ENTITY_STATUS_OUT_OF_SYNC);
      private final String name;

      VDS_STATUS(String value) {
         name = value;
      }

      public String getValue() {
         return name;
      }
   }

   public static enum HEALTH_STATUS {
      NORMAL(ENTITY_STATUS_NORMAL), WARNING(ENTITY_STATUS_WARNING), UNKNOWN(
            ENTITY_STATUS_UNKNOWN);
      private final String name;

      HEALTH_STATUS(String value) {
         name = value;
      }

      public String getValue() {
         return name;
      }
   }

   public static enum VMTAB_HEADER {
      VM_NAME_HEADER("Name"), VM_STATE_HEADER("State"), VM_STATUS_HEADER("Status"), VM_MEMO_SIZE_HEADER(
            "Memory Size (MB)"), VM_HOST_HEADER("Host"), VM_PROVISIONED_SPACE_HEADER(
            "Provisioned Space"), VM_USED_SPACE_HEADER("Used Space"), VM_HOST_CPU_HEADER(
            "Host CPU (MHz)"), VM_HOST_MEMO_HEADER("Host Mem (MB)"), VM_GUEST_MEMO_HEADER(
            "Guest Mem - %"), VM_GUEST_OS_HEADER("Guest OS"), VM_CPU_COUNT_HEADER(
            "CPU Count"), VM_NIC_COUNT_HEADER("NIC Count"), VM_ALARM_ACTION_HEADER(
            "Alarm Actions"), VM_VERSION_HEADER("VM Version");

      private final String columnName;

      VMTAB_HEADER(String colType) {
         columnName = colType;
      }

      public String getColumnName() {
         return columnName;
      }

      public static enum VMTAB_NAME {
         VM_NAME_HEADER, VM_STATE_HEADER, VM_STATUS_HEADER, VM_MEMO_SIZE_HEADER, VM_HOST_HEADER, VM_HOST_CPU_HEADER, VM_HOST_MEMO_HEADER, VM_GUEST_MEMO_HEADER, VM_GUEST_OS_HEADER, VM_CPU_COUNT_HEADER, VM_NIC_COUNT_HEADER, VM_ALARM_ACTION_HEADER
      }
   }

   public static enum ENUM_CLUSTER {

      STATUS_HA("HA"), STATUS_DRS("DRS");
      private final String statusType;

      ENUM_CLUSTER(String morType) {
         statusType = morType;
      }

      public String getStatusType() {
         return statusType;
      }
   }

   public static enum MENU_STATE {
      MENU_OPEN, MENU_CLOSED
   }

   public static enum POWER_OPS {
      POWER_ON, POWER_OFF, POWER_SUSPEND, POWER_RESET
   }

   public static enum INVENTORY_VIEW_OPTION {
      FULL_INVENTORY_VIEW, HOST_AND_CLUSTERS_VIEW, RESOURCE_POOL_VIEW, DATASTORES_VIEW, NETWORKS_VIEW, VMS_AND_TEMPLATES_VIEW
   }

   public static enum TASK_STATE {
      TASK_RUNNING, TASK_FAILED, TASK_COMPLETED, TASK_STATUS
   }

   public static enum VM_STATE {
      VM_POWERED_ON, VM_POWERED_OFF, VM_SUSPENDED
   }

   public static enum RENAME_NEG_SCENARIO {
      GREATER_THAN_MAX_CHARS, WHITE_SPACES, DUPLICATE_NAME, REMOVE_NODE_FROM_SERVER
   }

   public static enum _RESOURCE_ALLOCATION_TABS {
      _CPU_ALLOC_RESOURCE, _MEMORY_ALLOC_RESOURCE, _STORAGE_ALLOC_RESOURCE,
   }

   public static enum _PORTLET_ACTION {
      _PORTLET_MAXIMIZE, _PORTLET_RESTORE, _PORTLET_CLOSE, _PORTLET_MINIMIZE
   }

   public static enum NEW_RESOURCEPOOL {
      NEW_RP_CLUSTER, NEW_RP_HOST, NEW_RP_vAPP, NEW_RP_RESOURCEPOOL
   }

   public static enum FT_PANEL {
      GOS_VM, PRIMARY_HOST, SECONDARY_HOST
   }

   public static enum CLUSTER_TABS {
      VMS_TAB, RESOURCE_MGMT_TAB
   }

   public static enum _RESOURCE_ALLOCATION_VALUES {
      _CPU_TOTAL_CAPACITY, _CPU_RESERVED_CAPACITY, _CPU_AVAILABLE_CAPACITY, _CPU_RESERVATION_TYPE, _CPU_ENTITY_NAME, _CPU_ENTITY_SHARES, _CPU_ENTITY_SHARE_VALUE, _CPU_ENITY_PERCENT_SHARES, _CPU_ENTITY_RESERVATIONS, _CPU_ENTITY_RES_TYPE, _CPU_ENTITY_LIMIT, _MEMORY_TOTAL_CAPACITY, _MEMORY_RESERVED_CAPACITY, _MEMORY_AVAILABLE_CAPACITY, _MEMORY_RESERVATION_TYPE, _MEMORY_ENTITY_NAME, _MEMORY_ENTITY_SHARES, _MEMORY_ENTITY_SHARE_VALUE, _MEMORY_ENITY_PERCENT_SHARES, _MEMORY_ENTITY_RESERVATIONS, _MEMORY_ENTITY_RES_TYPE, _MEMORY_ENITY_LIMIT, _STORAGE_ENTITY_NAME, _STORAGE_ENTITY_DISK, _STORAGE_ENTITY_DATASTORE, _STORAGE_ENTITY_LIMIT, _STORAGE_ENTITY_SHARES, _STORAGE_ENTITY_SHARE_VALUE, _STORAGE_TOTAL_CAPACITY, _STORAGE_AVAILABLE_CAPACITY
   }

   public static enum DIALOG_ACTION {
      DA_OK, DA_CANCEL, DA_NEXT, DA_BACK, DA_FINISH, DA_ADD, DA_MOVE_UP, DA_MOVE_DOWN, DA_REMOVE
   }

   public static enum MSG_BOX_BUTTONS {
      MBB_YES, MBB_NO
   }

   public static enum USB_TYPE {
      USB11, USB20, USB30
   }

   public static enum USB_CONTROLLER_TYPE {
      NEW_USB_CONTROLLER, USB_CONTROLLER, USB_XHCI_CONTROLLER
   }

   public static enum MEMORY_VALUE_TYPE {
      CURRENT_VALUE, MINUMUM, DEFAULT, MAXIMUM, MAXIMUM_FOR_BEST_PERFORMANCE
   }

   // enum for Hard Disk and Memory Size Units
   public static enum SIZE_UNITS {
      MB, GB, TB
   }

   public enum GRID_MENU_ITEM {
      HIDE_COLUMN_NAME("1"), SIZE_TO_FIT_COLUMN_NAME("2"), SIZE_TO_FIT_ALL_COLUMNS("3"), LOCK_FIRST_COLUMN(
            "4"), SHOW_HIDE_COLUMNS("5"), SHOW_TOOLBAR("6"), ;

      private final String index;

      GRID_MENU_ITEM(String val) {
         index = val;
      }

      public String getIndex() {
         return index;
      }
   };

   /**
    * This is the Full List of Currently Supported Entities The nodeType value
    * matches the VC MOR Types (i.e. mor.getType()) Use getNodeType() to get the
    * Value
    *
    * @author akothari
    */
   public static enum NODE_TYPE {
      NT_VC("ServiceInstance"), NT_DATACENTER("Datacenter"), NT_HOST("HostSystem"), NT_CLUSTER(
            "ClusterComputeResource"), NT_VM("VirtualMachine"), NT_RESOURCE_POOL(
            "ResourcePool"), NT_DATASTORE("Datastore"), NT_FOLDER("Folder"), NT_VAPP(
            "VirtualApp"), NT_STANDARD_PORTGROUP("Network"), NT_STANDARD_SWITCH(
            "Network"), NT_DV_SWITCH("VmwareDistributedVirtualSwitch"), NT_DV_PORTGROUP(
            "DistributedVirtualPortgroup"), NT_DV_UPLINK(
            "DistributedVirtualUplinkPortgroup"), NT_TEMPLATE("VirtualMachineTemplate"), NT_STORAGE_POD(
            "StoragePod"), NT_VMFOLDER("VmFolder"), NT_HOST_PROFILE("HostProfile"), NT_COMPUTE_RESOURCE(
            "ComputeResource"),
      // Not sure if the node type is for storage profile is the correct one
      NT_VM_STORAGE_PROFILE("sps:SpsStorageProfile");

      private final String nodeType;

      NODE_TYPE(String morType) {
         nodeType = morType;
      }

      public String getNodeType() {
         return nodeType;
      }
   }

   /**
    * This is the Full List of Currently Supported Actions on all Entities
    * actionId variable is used to store the ID of the Supported Action Use
    * getActionId() to get the ID value from the Enum
    *
    * @author akothari
    */
   public static enum ACTION_TYPE {
      AT_POWER_ON(IDConstants.ID_ACTION_POWER_ON), AT_POWER_OFF(
            IDConstants.ID_ACTION_POWER_OFF), AT_SUSPEND(IDConstants.ID_ACTION_SUSPEND), AT_RESET(
            IDConstants.ID_ACTION_RESET), AT_MIGRATE(IDConstants.ID_ACTION_MIGRATE), AT_SNAPSHOT_MANAGER(
            IDConstants.ID_ACTION_SNAPSHOT_MANAGER), AT_REMOVE(
            IDConstants.ID_ACTION_REMOVE), AT_DELETE(IDConstants.ID_ACTION_DELETE), AT_RENAME(
            IDConstants.ID_ACTION_RENAME), AT_INSTALL_UPGRADE_TOOLS(
            IDConstants.ID_ACTION_INSTALL_UPGRADE_TOOLS), AT_UNMOUNT_TOOLS(
            IDConstants.ID_ACTION_UNMOUNT_TOOLS), AT_EDIT_SETTINGS(
            IDConstants.ID_ACTION_EDIT_SETTINGS), AT_EDIT_SETTINGS_VAPP(
            IDConstants.ID_ACTION_EDIT_SETTINGS_VAPP), AT_UNREGISTER(
            IDConstants.ID_ACTION_UNREGISTER), AT_NEW_FOLDER(
            IDConstants.ID_ACTION_NEW_FOLDER), AT_NEW_NETWORK_FOLDER(
            IDConstants.ID_ACTION_NETWORK_FOLDER), AT_NEW_VMFOLDER(
            IDConstants.ID_ACTION_NEW_VMFOLDER), AT_EDIT_SETTINGS_RP(
            IDConstants.ID_ACTION_EDIT_SETTINGS_RP), AT_NEW_RP(
            IDConstants.ID_ACTION_NEW_RP), AT_REMOVE_RP(IDConstants.ID_ACTION_REMOVE_RP), AT_REBOOT(
            IDConstants.ID_ACTION_REBOOT), AT_SHUTDOWN_HOST(
            IDConstants.ID_ACTION_SHUTDOWN_HOST), AT_ENTER_MAINTENANCE_MODE(
            IDConstants.ID_ACTION_ENTER_MAINTENANCE_MODE), AT_EXIT_MAINTENANCE_MODE(
            IDConstants.ID_ACTION_EXIT_MAINTENANCE_MODE), AT_ENTER_SDRS_MAINTENANCE_MODE(
            IDConstants.ID_ACTION_SDRS_ENTER_MAINTENANCE_MODE), AT_EXIT_SDRS_MAINTENANCE_MODE(
            IDConstants.ID_ACTION_SDRS_EXIT_MAINTENANCE_MODE), AT_MOVE_OUT_OF_DATASTORE_CLUSTER(
            IDConstants.ID_ACTION_MOVE_OUT_OF_DATASTORE_CLUSTER), AT_SHUTDOWN_GUEST(
            IDConstants.ID_ACTION_SHUTDOWN_GUEST), AT_RESTART_GUEST(
            IDConstants.ID_ACTION_RESTART_GUEST), AT_CLONE(IDConstants.ID_ACTION_CLONE), AT_TAKE_SNAPSHOT(
            IDConstants.ID_ACTION_TAKE_SNAPSHOT), AT_REVERT_CURRENT_SNAPSHOT(
            IDConstants.ID_ACTION_REVERT_CURRENT_SNAPSHOT), AT_CONVERT_TO_TEMPLATE(
            IDConstants.ID_ACTION_CONVERT_TO_TEMPLATE), AT_CREATE_VM(
            IDConstants.ID_ACTION_CREATE_VM), AT_CREATE_VM_DATASTORE(
            IDConstants.ID_ACTION_CREATE_VM_DATASTORE), AT_NEW_DC(
            IDConstants.ID_ACTION_NEW_DC), AT_UPGRADE_VIRTUAL_HW(
            IDConstants.ID_ACTION_UPGRADE_VIRTUAL_HW), AT_SHUTDOWN(
            IDConstants.ID_ACTION_SHUTDOWN), AT_PLUGIN_DISABLE(
            IDConstants.ID_PLUGIN_DISABLE), AT_PLUGIN_ENABLE(
            IDConstants.ID_PLUGIN_ENABLE), AT_POWER_ON_VAPP(IDConstants.ID_VAPP_POWER_ON), AT_POWER_OFF_VAPP(
            IDConstants.ID_VAPP_POWER_OFF), AT_SUSPEND_VAPP(IDConstants.ID_VAPP_SUSPEND), AT_NEW_VAPP(
            IDConstants.ID_ACTION_NEW_VAPP), AT_ANNOTATION(
            IDConstants.ID_ACTION_ANNOTATION), AT_REGISTER_VM(
            IDConstants.ID_ACTION_REGISTER_VM), AT_DEPLOY_VIRTUAL_MACHINE(
            IDConstants.ID_ACTION_DEPLOY_VIRTUAL_MACHINE), AT_SHUTDOWN_SINGLE_GUEST(
            IDConstants.ID_ACTION_SHUTDOWN_SINGLE_GUEST), AT_CREATE_VM_vmfolder(
            IDConstants.ID_ACTION_CREATE_VM_VMFolder), AT_TURN_FT_ON(
            IDConstants.ID_ACTION_TURN_FT_ON), AT_TURN_FT_OFF(
            IDConstants.ID_ACTION_TURN_FT_OFF), AT_DISABLE_FT(
            IDConstants.ID_ACTION_DISABLE_FT), AT_ENABLE_FT(
            IDConstants.ID_ACTION_ENABLE_FT), AT_MIGRATE_SECONDARY(
            IDConstants.ID_ACTION_MIGRATE_SECONDARY), AT_TEST_FAILOVER(
            IDConstants.ID_ACTION_TEST_FAILOVER), AT_TEST_RESTART_SECONDARY(
            IDConstants.ID_ACTION_TEST_RESTART_SECONDARY), AT_RENAME_CLUSTER(
            IDConstants.ID_ACTION_RENAME), AT_PROVISION_DATASTORE(
            IDConstants.ID_ACTION_PROVISION_DATASTORE), AT_REMOVE_DATASTORE(
            IDConstants.ID_ACTION_REMOVE_DATASTORE), AT_RENAME_DATASTORE(
            IDConstants.ID_ACTION_RENAME_DATASTORE), AT_RESCAN_DATASTORE(
            IDConstants.ID_ACTION_RESCAN_DATASTORE), AT_ADD_DIAGNOSTIC_PARTITION(
            IDConstants.ID_ACTION_ADD_DIAGNOSTIC_PARTITION), AT_MOUNT_DATASTORE(
            IDConstants.ID_ACTION_MOUNT_DATASTORE), AT_UNMOUNT_DATASTORE(
            IDConstants.ID_ACTION_UNMOUNT_DATASTORE), AT_UNMOUNT_NFS_DATASTORE(
            IDConstants.ID_ACTION_UNMOUNT_NFS_DATASTORE), AT_CONFIGURE_STORAGE_IO_CONTROL(
            IDConstants.ID_ACTION_CONFIGURE_STORAGE_IO_CONTROL), AT_INCREASE_DATASTORE_CAPACITY(
            IDConstants.ID_ACTION_INCREASE_DATASTORE_CAPACITY), AT_VM_CLONE(
            IDConstants.ID_ACTION_VM_CLONE), AT_CONSOLIDATE(
            IDConstants.ID_ACTION_CONSOLIDATE), AT_OPEN_CONSOLE(
            IDConstants.ID_ACTION_OPEN_CONSOLE), AT_TEMPLATE_CLONE(
            IDConstants.ID_ACTION_TEMPLATE_CLONE), AT_CONVERT_TEMPLATE_TO_VM(
            IDConstants.ID_ACTION_CONVERT_TEMPLATE_TO_VM), AT_CREATE_VDS(
            IDConstants.ID_ACTION_CREATE_VDS), AT_EDIT_SETTINGS_VDS(
            IDConstants.ID_ACTION_EDIT_SETTINGS_VDS), AT_EDIT_NETFLOW_VDS(
            IDConstants.ID_ACTION_EDIT_NETFLOW_VDS), AT_VDS_MANAGE_HOSTS(
            IDConstants.ID_ACTION_VDS_MANAGE_HOST), AT_DISCONNECT(
            IDConstants.ID_ACTION_DISCONNECT), AT_ADDHOST(IDConstants.ID_ACTION_ADDHOST), AT_CONNECT(
            IDConstants.ID_ACTION_RECONNECT), AT_CREATE_CLUSTER(
            IDConstants.ID_ACTION_CREATE_CLUSTER), AT_ADD_STORAGE(
            IDConstants.ID_ACTION_ADD_STORAGE), AT_RENAME_DATASTORE_CLUSTER(
            IDConstants.ID_ACTION_RENAME_DS_CLUSTER), AT_CREATE_DATASTORE_CLUSTER(
            IDConstants.ID_ACTION_NEW_DS_CLUSTER), AT_REMOVE_DATASTORE_CLUSTER(
            IDConstants.ID_ACTION_REMOVE_DS_CLUSTER), AT_EXPORT_OVF(
            IDConstants.ID_ACTION_EXPORT_OVF), AT_ACTION_CLONE_VM_TO_TEMPLATE(
            IDConstants.ID_ACTION_VM_CLONE_TO_TEMPLATE), AT_EDIT_GENERAL_SETTINGS(
            IDConstants.ID_ACTION_EDIT_GENERAL_SETTINGS), AT_ADD_NETWORKING(
            IDConstants.ID_ACTION_ADD_NETWORKING), AT_MOVE_NETWORK(
            IDConstants.ID_ACTION_MOVE_NETWORK), AT_EDIT_VM_START_SHUTDOWN_CONFIG(
            IDConstants.ID_EDIT_VM_START_SHUTDOWN_CONFIG), AT_CREATE_DATACENTER(
            IDConstants.ID_NEW_DATACENTER_BUTTON), AT_REMOVE_DATACENTER(
            IDConstants.ID_REMOVE_DATACENTER), AT_MANAGE_PORTGROUPS_VDS(
            IDConstants.ID_ACTION_MANAGE_PORTGROUPS_VDS), AT_STORAGE_IO_CONTROL(
            IDConstants.ID_HOST_STORAGE_IO_CONTROL), AT_CREATE_PORT_GROUP(
            IDConstants.ID_ACTION_CREATE_PORT_GROUP), AT_EDIT_DVPORTGROUP(
            IDConstants.ID_ACTION_EDIT_SETTINGS_DVPORTGROUP), AT_CREATE_DVPORTGROUP(
            IDConstants.ID_ACTION_CREATE_DVPORTGROUP), AT_EDIT_PORT_GROUP(
            IDConstants.ID_ACTION_EDIT_SETTINGS_DVPORTGROUP), AT_ADD_ALARM(
            IDConstants.ID_ACTION_ADD_ALARM), AT_DISABLE_ALARM(
            IDConstants.ID_ACTION_DISABLE_ALARM), AT_EDIT_CPU_SETTINGS_RP(
            IDConstants.ID_ACTION_EDIT_CPU_SETTINGS_RP), AT_EDIT_MEMORY_SETTINGS_RP(
            IDConstants.ID_ACTION_EDIT_MEMORY_SETTINGS_RP), AT_CREATE_HOST_PROFILE(
            IDConstants.ID_ACTION_CREATE_HOST_PROFILE_MENU), AT_RESET_HOST_CUSTOMIZATIONS(
            IDConstants.ID_ACTION_RESET_HOST_CUSTOMIZATIONS), AT_CREATE_HOST_PROFILE_BUTTON(
            IDConstants.ID_ACTION_CREATE_HOST_PROFILE), AT_EDIT_HOST_PROFILE(
            IDConstants.ID_ACTION_EDIT_HOST_PROFILE_MENU), AT_EDIT_HOST_PROFILE_BUTTON(
            IDConstants.ID_ACTION_EDIT_HOST_PROFILE), AT_DELETE_HOST_PROFILE_BUTTON(
            IDConstants.ID_ACTION_DELETE_HOST_PROFILE), AT_DELETE_HOST_PROFILE(
            IDConstants.ID_ACTION_DELETE_HOST_PROFILE_MENU), AT_CHECK_COMPLIANCE(
            IDConstants.ID_ACTION_CHECK_COMPLIANCE), AT_CHECK_HOST_PROFILE_COMPLIANCE(
            IDConstants.ID_ACTION_CHECK_HOST_PROFILE_COMPLIANCE), AT_RENAME_HOST_PROFILE(
            IDConstants.ID_ACTION_HOST_PROFILE_RENAME), AT_RUN_SCHEDULED_TASK(
            IDConstants.ID_ACTION_RUN_SCHEDULED_TASK), AT_COPY_SETTINGS_FROM_HOST(
            IDConstants.ID_ACTION_COPY_SETTINGS_FROM_HOST), AT_RESTORE_RESOURCEPOOL(
            IDConstants.ID_ACTION_RESTORE_RESOURCEPOOL), AT_MIGRATE_VM_NETWORKING(
            IDConstants.ID_ACTION_MIGRATE_VM_NETWORKING), AT_RENAME_DC_FOLDER(
            IDConstants.ID_ACTION_DC_FOLDERS_RENAME), AT_HOST_PROFILE_REMEDIATE(
            IDConstants.ID_ACTION_REMEDIATE_HOST), AT_UNREGISTER_DC(
            IDConstants.ID_ACTION_UNREGISTER_DC), AT_HOST_PROFILE_DETACH(
            IDConstants.ID_ACTION_DETACH_HOST), AT_HOST_PROFILE_DETACH_MULTIPLE(
            IDConstants.ID_ACTION_DETACH_HOST), AT_CHANGE_HOST_PROFILE(
            IDConstants.ID_ACTION_CHANGE_HOST_PROFILE), AT_DUPLICATE_HOST_PROFILE(
            IDConstants.ID_ACTION_DUPLICATE_HOST_PROFILE), AT_ATTACH_HOSTS(
            IDConstants.ID_ACTION_ATTACH_HOST), AT_UPGRADE_VDS(
            IDConstants.ID_ACTION_UPGRADE_VDS), AT_MOVE(IDConstants.ID_ACTION_MOVE), AT_ASSIGN_STORAGE_PROFILE(
            IDConstants.ID_ACTION_ASSIGN_STORAGE_PROFILE), AT_EDIT_STORAGE_PROFILE(
            IDConstants.ID_ACTION_EDIT_STORAGE_PROFILE), AT_REMOVE_STORAGE_PROFILE(
            IDConstants.ID_ACTION_REMOVE_STORAGE_PROFILE), AT_ENTER_STANDBY_MODE(
            IDConstants.ID_ACTION_ENTER_STANDBY_MODE_HOST), AT_DEPLOY_OVF(
            IDConstants.ID_ACTION_DEPLOY_OVF), AT_CLONE_VAPP(
            IDConstants.ID_ACTION_CLONE_VAPP), AT_EXIT_STANDBY_MODE(
            IDConstants.ID_ACTION_POWER_ON_HOST), AT_EDIT_DEFAULT_VM_COMPATIBILITY(
            IDConstants.ID_CONTEXTMENU_EDIT_DEFAULT_VM_COMPATIBILITY), AT_CANCEL_SCHEDULED_VM_UPGRADE(
            IDConstants.ID_CONTEXTMENU_CANCEL_SCHEDULED_VM_UPGRADE), AT_UPGRADE_TO_VMFS5_DATASTORE(
            IDConstants.ID_ACTION_UPGRADE_VMFS5_DATASTORE), AT_ACTION_ASSIGN_TAG(
            IDConstants.ID_ACTION_ASSIGN_TAG), AT_EDIT_ALARM(
            IDConstants.ID_ACTION_EDIT_ALARM), AT_ACTION_FILE_BROWSER(
            IDConstants.ID_ACTION_FILE_BROWSER), AT_ENTER_SDRS_MAINTANENCE_MODE(
            "vsphere.core.datastore.enterMaintenanceAction"), AT_EXIT_SDRS_MAINTANENCE_MODE(
            "vsphere.core.datastore.exitMaintenanceAction"), AT_ACTION_EDIT_RESOURCE_SETTING(
            IDConstants.ID_ACTION_EDIT_RESOURCE_SETTING), AT_NEW_HOST_CLUSTER_FOLDER(
            IDConstants.ID_ACTION_HOSTANDCLUSTER_FOLDER), AT_NEW_STORAGE_FOLDER(
            IDConstants.ID_ACTION_STORAGE_FOLDER), AT_NEW_VMFOLDER_SUBFOLDER(
            IDConstants.ID_ACTION_NEW_VMSUBFOLDER), AT_RECONFIGURE_HA(
            IDConstants.ID_ACTION_RECONFIGURE_HA);

      private final String actionId;

      ACTION_TYPE(String aID) {
         actionId = aID;
      }

      public String getActionId() {
         return actionId;
      }

      public static enum NODE_TYPE {
         NT_VC, NT_HOST, NT_CLUSTER, NT_VM, NT_DATACENTER, NT_RESOURCE_POOL, NT_VAPP, NT_STANDARD_PORTGROUP, NT_DV_SWITCH, NT_DV_PORTGROUP, NT_DV_UPLINK, NT_DATASTORE, NT_FOLDER, NT_TEMPLATE

      }
   }

   public static enum DATAGRID_ACTION_TYPE {
      DATAGRID_ACTION_FILTER(IDConstants.ID_SEARCHCONTROL_FILTERCONTROL), DATAGRID_ACTION_RESCAN_ADAPTER(
            IDConstants.ID_ACTION_RESCAN_ADAPTER), DATAGRID_ACTION_REFRESH_ADAPTER(
            IDConstants.ID_ACTION_REFRESH_ADAPTER), DATAGRID_ACTION_RESCAN_STORAGE(
            IDConstants.ID_ACTION_RESCAN_STORAGE), DATAGRID_ACTION_ADD_ISCSI(
            IDConstants.ID_BUTTON_ADD_ISCSI), DATAGRID_ACTION_FIND(
            IDConstants.ID_IMAGE_FIND);

      private final String dataGridActionId;

      DATAGRID_ACTION_TYPE(String aID) {
         dataGridActionId = aID;
      }

      public String getActionId() {
         return dataGridActionId;
      }

   }

   public static enum VAPP_STATE {
      VAPP_POWERED_ON, VAPP_POWERED_OFF, VAPP_SUSPENDED
   }

   /**
    * Enum defined for setting Shares Level for vApp
    */
   public static enum SHARES_LEVEL {
      LOW, HIGH, NORMAL, CUSTOM
   }

   public static enum VLAN_TYPE {
      VT_NO_VLAN, VT_VLAN, VT_VLAN_TRUNKING, VT_PRIVATE_VLAN
   }

   // Add more tabs here
   public static enum MAIN_NAVIGATION_TABS {
      MAIN_TABS_GETTING_STARTED, MAIN_TABS_SUMMARY, MAIN_TABS_MONITOR, MAIN_TABS_MANAGE, MAIN_TABS_RELATED_ITEMS
   };

   public static enum PRIVILEGE_POWEROPS {
      PRIVILEGES_POWER_ON, PRIVILEGES_POWER_OFF, PRIVILEGES_POWER_ONOFFSUSPENDRESET, PRIVILEGES_ONLY_POWER_ON, PRIVILEGES_POWER_SUSPEND, PRIVILEGES_POWER_RESET, PRIVILEGE_ONLY_POWERONOFF, PRIVILEGE_POWERONOFF
   }

   public static enum DISK_MODE {
      DEPENDENT, INDEPENDENT_PERSISTENT, INDEPENDENT_NON_PERSISTENT
   }

   public static enum CLUSTER_CONFIGURATION_LIST_OPTIONS {
      SERVICES("0"), DRS("1"), HA("2"), CONFIGURATION("3"), GENERAL("4"), VMWARE_EVC("5"), GROUPS(
            "6"), RULES("7"), VM_OVERRIDES("8"), HOST_OPTIONS("9");

      private final String clusterOption;

      CLUSTER_CONFIGURATION_LIST_OPTIONS(String option) {
         clusterOption = option;
      }

      public String getClusterOptionIndex() {
         return clusterOption;
      }
   };


   public static enum CLUSTER_DRS_LIST_OPTIONS {
      RECOMMENDATIONS("1"), FAULTS("2"), HISTORY("3"), CPU_UTILIZATION("4"), MEMORY_UTILIZATION(
            "5");


      private final String drsOption;


      CLUSTER_DRS_LIST_OPTIONS(String option) {
         drsOption = option;
      }


      public String getDRSOptionIndex() {
         return drsOption;

      }

   };


   public static enum ClusterRuleType {
      RULE_TYPE_KEEP_VMS_TOGETHER("Keep Virtual Machines Together"), RULE_TYPE_SEPARATE_VMS(
            "Separate Virtual Machines"), RULE_TYPE_VM_TO_HOSTS(
            "Virtual Machines to Hosts");
      private final String ruleType;

      ClusterRuleType(String ruleId) {
         ruleType = ruleId;
      }

      public String getRuleType() {
         return ruleType;
      }
   };

   // Datastore > Manage > Settings page specific
   public static enum DATASTORE_MANAGE_SETTING_PAGE_LIST_OPTIONS {
      GENERAL("General"), DEVICE_BACKING("Device Backing"), CONNECTIVITY_AND_MULTIPATHING(
            "Connectivity and Multipathing");
      private final String dsConfig;

      DATASTORE_MANAGE_SETTING_PAGE_LIST_OPTIONS(String dsConfigVal) {
         dsConfig = dsConfigVal;
      }

      public String getDatastoreManageSettingsPageVal() {
         return dsConfig;
      }
   }

   // NFS Datastore > Manage > Settings page specific
   public static enum NFS_MANAGE_SETTING_PAGE_LIST {
      GENERAL("General"), DEVICE_BACKING("Device Backing"), CONNECTIVITY_WITH_HOSTS(
            "Connectivity with Hosts");
      private final String pageValue;

      NFS_MANAGE_SETTING_PAGE_LIST(String pageSetting) {
         pageValue = pageSetting;
      }

      public String getValue() {
         return pageValue;
      }
   }

   public static enum RULE_VM_TO_HOSTS_RELATION_TYPES {
      HARD_AFFINITY, SOFT_AFFINITY, HARD_ANTI_AFFINITY, SOFT_ANTI_AFFINITY
   };

   public enum RULE {
      DRS_CLUSTER_AFFINITY, DRS_CLUSTER_ANTI_AFFINITY
   };

   public enum CLUSTER {
      DRS_CLUSTER, DRS_CLUSTER_HOST_ONE, DRS_CLUSTER_HOST_TWO
   }

   public enum FAULTS {
      DRS_FAULT_AFFINITY_ANTI_AFFINITY, DRS_FAULT_PINNED, DRS_FAULT_INCOMPATIBILITY, DRS_FAULT_STANDBY_AFFINITY_ANTIAFFINITY, DRS_FAULT_HA, DRS_FAULT_NO_ACTIVE_HOST
   }

   public enum RuleType {
      VMDK_AFFINITY("VIRTUAL_DISK_AFFINITY_RULE_TYPE", "VMDK affinity"), VMDK_ANTI_AFFINITY(
            "VIRTUAL_DISK_ANTI_AFFINITY_RULE_TYPE", "VMDK anti-affinity"), VM_ANTI_AFFINITY(
            "VM_ANTIAFFINITY_RULE_TYPE", "VM anti-affinity");

      private final String value;
      private final String label;

      private RuleType(String value, String label) {
         this.value = value;
         this.label = label;
      }

      public String getValue() {
         return value;
      }

      public String getLabel() {
         return label;
      }
   }

   public static enum VmOverridesAutomationLevel {
      FULLY_AUTOMATED("Fully Automated"), PARTIALLY_AUTOMATED("Partially Automated"), MANUAL(
            "Manual"), CLUSTER_DEFAULT("Use Cluster Settings"), DISABLED("Disabled");

      private final String autoLevel;

      VmOverridesAutomationLevel(String autoId) {
         autoLevel = autoId;
      }

      public String getAutoLevel() {
         return autoLevel;
      }
   }

   // VMOverrides Restart Priority
   public static enum VmOverridesRestartPriority {
      DISABLED(TestConstantsKey.RESTART_PRIORITY_DISABLED), LOW(
            TestConstantsKey.RESTART_PRIORITY_LOW), MEDIUM(
            TestConstantsKey.RESTART_PRIORITY_MEDIUM), HIGH(
            TestConstantsKey.RESTART_PRIORITY_HIGH), CLUSTER_DEFAULT(
            TestConstantsKey.RESTART_PRIORITY_CLUSTER_DEFAULT);

      private final String restartPriority;

      VmOverridesRestartPriority(String restartPriority) {
         this.restartPriority = restartPriority;
      }

      public String get() {
         return restartPriority;
      }
   }

   // VMOverrides Isolation Response
   public static enum VmOverridesIsolationResponse {
      LEAVEPOWEREDON(TestConstantsKey.ISOLATION_RESPONSE_LEAVE_POWERED_ON), POWEROFF(
            TestConstantsKey.ISOLATION_RESPONSE_POWER_OFF), SHUTDOWN(
            TestConstantsKey.ISOLATION_RESPONSE_SHUT_DOWN), CLUSTER_DEFAULT(
            TestConstantsKey.ISOLATION_RESPONSE_CLUSTER_DEFAULT);

      private final String isolationResponse;

      VmOverridesIsolationResponse(String isolationResponse) {
         this.isolationResponse = isolationResponse;
      }

      public String get() {
         return isolationResponse;
      }
   }

   // VM Monitoring
   public static enum VmMonitoring {
      DISABLED(TestConstantsKey.VM_MONITORING_DISABLED), VM_MONITORING_ONLY(
            TestConstantsKey.VM_MONITORING_ONLY), APPLICATIONS(
            TestConstantsKey.VM_MONITORING_AND_APP);

      private final String vmMonitoring;

      VmMonitoring(String vmMonitoring) {
         this.vmMonitoring = vmMonitoring;
      }

      public String get() {
         return vmMonitoring;
      }
   }

   // EVC Selection modes
   public static enum EvcSelectionModes {
      EVC_DISABLED(ID_CLUSTER_SETTINGS_EVC_EDIT_DISABLE_EVC_RADIO), EVC_AMD(
            ID_CLUSTER_SETTINGS_EVC_EDIT_AMD_RADIO), EVC_INTEL(
            ID_CLUSTER_SETTINGS_EVC_EDIT_INTEL_RADIO);

      private final String evcMode;

      EvcSelectionModes(String evcMode) {
         this.evcMode = evcMode;
      }

      public String getUiComponent() {
         return evcMode;
      }
   }

   /**
    * DRS VM Overrides ADG columns and indexes Keep the sequence in sync with
    * the actual default UI column sequence
    */
   public static enum DrsVmOverridesTableColumns {
      VM_OVERRIDES_NAME_COLUMN, VM_OVERRIDES_AUTO_LEVEL_COLUMN, VM_OVERRIDES_RESTART_PRIORITY_COLUMN, VM_OVERRIDES_HOST_ISOLATION_RESPONSE_COLUMN, VM_OVERRIDES_VM_MONITORING_COLUMN, VM_OVERRIDES_MONITORING_SENSITIVITY_COLUMN;
   }

   public static enum ISCSI_MATRIX {
      ROW_0, COLUMN_0, ROW_1, COLUMN_1
   };

   public static final HashMap<ISCSI_MATRIX, String> ISCSI_MATRIX_MAP =
         new HashMap<ISCSI_MATRIX, String>();
   static {
      ISCSI_MATRIX_MAP.put(ISCSI_MATRIX.ROW_0, "0");
      ISCSI_MATRIX_MAP.put(ISCSI_MATRIX.COLUMN_0, "0");
      ISCSI_MATRIX_MAP.put(ISCSI_MATRIX.ROW_1, "1");
      ISCSI_MATRIX_MAP.put(ISCSI_MATRIX.COLUMN_1, "1");
   }

   // Object Navigator Category Node Views
   public static enum OBJ_NAV_APP_VIEW {
      VIEW_LAUNCHER, VIEW_MANAGE_MONITOR, VIEW_VI_HOME, VIEW_RULES_PROFILE, VIEW_ADMINISTRATION, VIEW_LICENSES, VIEW_LICENSE_REPORTS, VIEW_VCS_EXTENSIONS, VIEW_SEARCH, VIEW_SAVED_SEARCHES, VIEW_TASKS, VIEW_PLUGIN_MANAGEMENT, VIEW_SSO_CONFIGURATION, VIEW_USERS_GROUPS, VIEW_ROLE_MANAGER, VIEW_TAGS
   }

   // Object Navigator Tree Node Views
   public static enum OBJ_NAV_TREE_NODE_VIEW {
      VIEW_VIRTUAL_CENTERS, VIEW_DATACENTERS, VIEW_COMPUTE, VIEW_NETWORKING, VIEW_STORAGE, VIEW_VMS_TEMPLATES, VIEW_HOST_PROFILES, VIEW_VMS_STORAGE, VIEW_TAG_MANAGER, VIEW_SOLUTION_MANAGER, VIEW_HOSTS, VIEW_CLUSTERS, VIEW_DATASTORES, VIEW_DATASTORE_CLUSTERS, VIEW_STANDARD_NETWORKS, VIEW_DISTRIBUTED_SWITCHES, VIEW_DISTRIBUTED_PORT_GROUPS, VIEW_DISTRIBUTED_UPLINK_PORT_GROUPS, VIEW_VIRTUAL_MACHINES, VIEW_VAPPS, VIEW_HOST, VIEW_HOSTS_AND_CLUSTERS_TREE, VIEW_VMS_AND_TEMPLATES_TREE, VIEW_STORAGE_TREE, VIEW_NETWORKING_TREE, VIEW_DATACENTER, VIEW_DATASTORE, VIEW_VM_NETWORK, VIEW_RESOURCE_POOLS, VIEW_SAVED_SEARCHES, VIEW_VIRTUAL_CENTER, VIEW_DISTRIBUTED_SWITCH, VIEW_ORGANIZATIONS, VIEW_VIRTUAL_DATACENTERS, VIEW_CLOUD_RESOURCE_POOLS
   }

   // Cluster DRS automation level values
   public static enum CLUSTER_DRS_AUTO_LEVEL_VALUE {
      FULLY_AUTOMATED_VALUE("Fully Automated"), PARTIALLY_AUTOMATED_VALUE(
            "Partially Automated"), MANUAL_VALUE("Manual"), FULLY_AUTOMATED_INDEX("2"), PARTIALLY_AUTOMATED_INDEX(
            "1"), MANUAL_INDEX("0");

      private final String autoValue;

      CLUSTER_DRS_AUTO_LEVEL_VALUE(String autoLevelValue) {
         autoValue = autoLevelValue;
      }

      public String getDrsAutoValue() {
         return autoValue;
      }
   }

   // Cluster DPM automation level values
   public static enum ClusterDpmAutoLevel {
      OFF("Off"), MANUAL("Manual"), AUTOMATIC("Automatic"), OFF_INDEX("0"), MANUAL_INDEX(
            "1"), AUTOMATIC_INDEX("2");

      private final String autoValue;

      ClusterDpmAutoLevel(String autoLevelValue) {
         autoValue = autoLevelValue;
      }

      public String getDpmAutoValue() {
         return autoValue;
      }
   }

   // DRS Service List
   public static enum ClusterServicesListOptions {
      SERVICES("0"), DRS("1"), HA("2");

      private final String pos;

      ClusterServicesListOptions(String pos) {
         this.pos = pos;
      }

      public String getPosition() {
         return pos;
      }
   };

   public static enum LATENCY_LEVEL {
      LATENCY_LEVEL_CUSTOM, LATENCY_LEVEL_HIGH, LATENCY_LEVEL_MEDIUM, LATENCY_LEVEL_NORMAL, LATENCY_LEVEL_LOW
   }

   public static enum SCSI_CONTROLLER_TYPE {
      LSI_LOGIC_SAS(VM_DEVICE_LSI_SAS_LOGIC), BUS_LOGIC_PARALLEL(VM_DEVICE_BUS_LOGIC), VMWARE_PARAVIRTUAL(
            VM_DEVICE_PVSCSI), LSI_LOGIC_PARALLEL(VM_DEVICE_LSI_LOGIC_PARALLEL);

      private final String SCSICtrlType;

      SCSI_CONTROLLER_TYPE(String SCSIControllerType) {
         SCSICtrlType = SCSIControllerType;
      }

      public String getSCSIControllerType() {
         return SCSICtrlType;
      }
   }

   public static enum VM_PROVISIONING_TYPE {
      CREATE_NEW_VM(CREATE_FROM_SCRATCH), DEPLOY_TEMPLATE(DEPLOY_FROM_TEMPLATE), CLONE_VM_TO_VM(
            CLONE_EXISTING_VM), CLONE_VM_TO_TEMPLATE(CLONE_TEMPLATE_FROM_VM), CLONE_TEMPLATE_TO_TEMPLATE(
            CLONE_EXISTING_TEMPLATE), CONVERT_TEMPLATE_IN_VM(CONVERT_TEMPLATE_TO_VM);

      private final String vmProvisionType;

      VM_PROVISIONING_TYPE(String VMProvType) {
         vmProvisionType = VMProvType;
      }

      public String getVMProvisioningType() {
         return vmProvisionType;
      }
   }

   public static enum FLOPPY_DEVICE_TYPE {
      CLIENT_DEVICE(FLOPPY_TYPE_CLIENT_DEVICE), EXISITING_FILE(FLOPPY_TYPE_EXISTING_FILE), NEW_FLOPPY_IMAGE(
            FLOPPY_TYPE_NEW_IMAGE);

      private final String floppyDevice;

      FLOPPY_DEVICE_TYPE(String floppyDeviceType) {
         floppyDevice = floppyDeviceType;
      }

      public String getFloppyDeviceType() {
         return floppyDevice;
      }
   }

   public static enum GOS_FAMILY {
      Windows, Linux, Other;
   }

   public static enum GOS_VERSION {

      GOS_VERSION_WIN_2008_R2_64("Microsoft Windows Server 2008 R2 (64-bit)"), GOS_VERSION_WIN_2008_32(
            "Microsoft Windows Server 2008 (32-bit)"), GOS_VERSION_WIN_2008_64(
            "Microsoft Windows Server 2008 (64-bit)"), GOS_VERSION_WIN_7_32(
            "Microsoft Windows 7 (32-bit)"), GOS_VERSION_WIN_7_64(
            "Microsoft Windows 7 (64-bit)"), GOS_VERSION_WIN_VISTA_32(
            "Microsoft Windows Vista (32-bit)"), GOS_VERSION_WIN_VISTA_64(
            "Microsoft Windows Vista (64-bit)"), GOS_VERSION_WIN_SERVER2003_32(
            "Microsoft Windows Server 2003 (32-bit)"), GOS_VERSION_WIN_SERVER2003_64(
            "Microsoft Windows Server 2003 (64-bit)"), GOS_VERSION_WIN_SERVER2003ENTERPRISE_32(
            "Microsoft Windows Server 2003 Enterprise (32-bit)"), GOS_VERSION_WIN_SERVER2003ENTERPRISE_64(
            "Microsoft Windows Server 2003 Enterprise (64-bit)"), GOS_VERSION_WIN_SERVER2003DATACENTER_32(
            "Microsoft Windows Server 2003 Datacenter (32-bit)"), GOS_VERSION_WIN_SERVER2003DATACENTER_64(
            "Microsoft Windows Server 2003 Datacenter (64-bit)"), GOS_VERSION_WIN_SERVER2003STANDARD_32(
            "Microsoft Windows Server 2003 Standard (32-bit)"), GOS_VERSION_WIN_SERVER2003STANDARD_64(
            "Microsoft Windows Server 2003 Standard (64-bit)"), GOS_VERSION_WIN_SERVER2003WEB_32(
            "Microsoft Windows Server 2003 Web Edition (32-bit)"), GOS_VERSION_WIN_SBS2003(
            "Microsoft Small Business Server 2003"), GOS_VERSION_WIN_XPPRO_32(
            "Microsoft Windows XP Professional (32-bit)"), GOS_VERSION_WIN_XPPRO_64(
            "Microsoft Windows XP Professional (64-bit)"), GOS_VERSION_WIN_2000(
            "Microsoft Windows 2000"), GOS_VERSION_WIN_NT("Microsoft Windows NT"), GOS_VERSION_WIN_98(
            "Microsoft Windows 98"), GOS_VERSION_WIN_95("Microsoft Windows 95"), GOS_VERSION_WIN_3_1(
            "Microsoft Windows 3.1"), GOS_VERSION_LINUX_RED_HAT_ENT6_32(
            "Red Hat Enterprise Linux 6 (32-bit)"), GOS_VERSION_LINUX_RED_HAT_ENT6_64(
            "Red Hat Enterprise Linux 6 (64-bit)"), GOS_VERSION_LINUX_RED_HAT_ENT5_32(
            "Red Hat Enterprise Linux 5 (32-bit)"), GOS_VERSION_LINUX_RED_HAT_ENT5_64(
            "Red Hat Enterprise Linux 5 (64-bit)"), GOS_VERSION_LINUX_RED_HAT_ENT4_32(
            "Red Hat Enterprise Linux 4 (32-bit)"), GOS_VERSION_LINUX_RED_HAT_ENT4_64(
            "Red Hat Enterprise Linux 4 (64-bit)"), GOS_VERSION_LINUX_RED_HAT_ENT3_32(
            "Red Hat Enterprise Linux 3 (32-bit)"), GOS_VERSION_LINUX_RED_HAT_ENT3_64(
            "Red Hat Enterprise Linux 3 (64-bit)"), GOS_VERSION_LINUX_RED_HAT_ENT21(
            "Red Hat Enterprise Linux 2.1"), GOS_VERSION_CENT_OS_4_5_6_32(
            "CentOS 4/5/6 (32-bit)"), GOS_VERSION_CENT_OS_4_5_6_64(
            "CentOS 4/5/6 (64-bit)"), GOS_VERSION_ORACLE_ENT_LINUX4_5_32(
            "Oracle Linux 4/5/6 (32-bit)"), GOS_VERSION_ORACLE_ENT_LINUX4_5_64(
            "Oracle Linux 4/5/6 (64-bit)"), GOS_VERSION_NOVELL_SUSE_LINUX_ENT11_32(
            "Novell SUSE Linux Enterprise 11 (32-bit)"), GOS_VERSION_NOVELL_SUSE_LINUX_ENT11_64(
            "Novell SUSE Linux Enterprise 11 (64-bit)"), GOS_VERSION_NOVELL_SUSE_LINUX_ENT10_32(
            "Novell SUSE Linux Enterprise 10 (32-bit)"), GOS_VERSION_NOVELL_SUSE_LINUX_ENT10_64(
            "Novell SUSE Linux Enterprise 10 (64-bit)"), GOS_VERSION_NOVELL_SUSE_LINUX_ENT8_9_32(
            "Novell SUSE Linux Enterprise 8/9 (32-bit)"), GOS_VERSION_NOVELL_SUSE_LINUX_ENT8_9_64(
            "Novell SUSE Linux Enterprise 8/9 (64-bit)"), GOS_VERSION_NOVELL_OPEN_ENT_SERVER(
            "Novell Open Enterprise Server"), GOS_VERSION_ASIANUX3_32(
            "Asianux 3 (32-bit)"), GOS_VERSION_ASIANUX3_64("Asianux 3 (64-bit)"), GOS_VERSION_DEBIAN_GNULINUX5_32(
            "Debian GNU/Linux 5 (32-bit)"), GOS_VERSION_DEBIAN_GNULINUX5_64(
            "Debian GNU/Linux 5 (64-bit)"), GOS_VERSION_DEBIAN_GNULINUX4_32(
            "Debian GNU/Linux 4 (32-bit)"), GOS_VERSION_DEBIAN_GNULINUX4_64(
            "Debian GNU/Linux 4 (64-bit)"), GOS_VERSION_UBUNTU_LINUX_32(
            "Ubuntu Linux (32-bit)"), GOS_VERSION_UBUNTU_LINUX_64(
            "Ubuntu Linux (64-bit)"), GOS_VERSION_OTHER_26x_LINUX_32(
            "Other 2.6.x Linux (32-bit)"), GOS_VERSION_OTHER_26x_LINUX_64(
            "Other 2.6.x Linux (64-bit)"), GOS_VERSION_OTHER_24x_LINUX_32(
            "Other 2.4.x Linux (32-bit)"), GOS_VERSION_OTHER_24x_LINUX_64(
            "Other 2.4.x Linux (64-bit)"), GOS_VERSION_OTHER_LINUX_32(
            "Other Linux (32-bit)"), GOS_VERSION_OTHER_LINUX_64("Other Linux (64-bit)"), GOS_VERSION_NOVELL_NETWARE6x(
            "Novell NetWare 6.x"), GOS_VERSION_NOVELL_NETWARE5_1("Novell NetWare 5.1"), GOS_VERSION_ORACLE_SOLARIS10_32(
            "Oracle Solaris 10 (32-bit)"), GOS_VERSION_ORACLE_SOLARIS10_64(
            "Oracle Solaris 10 (64-bit)"), GOS_VERSION_SUN_SOLARIS9_EXPERIMENTAL(
            "Sun Microsystems Solaris 9"), GOS_VERSION_SUN_SOLARIS8_EXPERIMENTAL(
            "Sun Microsystems Solaris 8"), GOS_VERSION_FREEBSD_32("FreeBSD (32-bit)"), GOS_VERSION_FREEBSD_64(
            "FreeBSD (64-bit)"), GOS_VERSION_APPLE_MACOSXSERVER10_6_32(
            "Apple Mac OS X Server 10.6 (32-bit)"), GOS_VERSION_APPLE_MACOSXSERVER10_6_64(
            "Apple Mac OS X Server 10.6 (64-bit)"), GOS_VERSION_APPLE_MACOSXSERVER10_5_32(
            "Apple Mac OS X Server 10.5 (32-bit)"), GOS_VERSION_APPLE_MACOSXSERVER10_5_64(
            "Apple Mac OS X Server 10.5 (64-bit)"), GOS_VERSION_SERENITY_ECOMSTATION(
            "Serenity Systems eComStation 1"), GOS_VERSION_IBM_OS2("IBM OS/2"), GOS_VERSION_SCO_OPENSERVER5(
            "SCO OpenServer 5"), GOS_VERSION_SCO_UNIXWARE7("SCO UnixWare 7"), GOS_VERSION_MICROSOFT_DOS(
            "Microsoft MS-DOS"), GOS_VERSION_OTHER_32("Other (32-bit)");

      private final String gosVersion;

      GOS_VERSION(String gosVersionToSelect) {
         gosVersion = gosVersionToSelect;
      }

      public String getGOSVersion() {
         return gosVersion;
      }
   }

   public static enum NETWORK_ADAPTER_TYPE {
      E1000(NIC_TYPE_E1000), VMXNET2(NIC_TYPE_VirtualVmxnet2), VMXNET3(
            NIC_TYPE_VirtualVmxnet3), VMXNET(NIC_TYPE_VirtualVmxnet), FLEXIBLE(
            NIC_TYPE_FLEXIBLE), E1000E(NIC_TYPE_E1000E);

      private final String nwAdapterType;

      NETWORK_ADAPTER_TYPE(String adapterTypeToSelect) {
         nwAdapterType = adapterTypeToSelect;
      }

      public String getNetworkAdapterType() {
         return nwAdapterType;
      }
   }

   public static enum NETWORK_MAC_ADDRESS_TYPE {
      Manual, Automatic;
   }

   public static enum SCSI_BUS_SHARING {
      None, Virtual, Physical;
   }

   public static enum VIDEO_CARD_DISPLAY {
      one("1"), two("2"), three("3"), four("4");

      private final String videocardDisplay;

      VIDEO_CARD_DISPLAY(String cardDisplayToSelect) {
         videocardDisplay = cardDisplayToSelect;
      }

      public String getVideoCardDisplay() {
         return videocardDisplay;
      }
   }

   public static enum CDDVD_DEVICE_MODE {
      EMULATE("Emulate IDE"), PASSTHROUGH("Passthrough IDE");

      private final String deviceMode;

      CDDVD_DEVICE_MODE(String deviceModeToSelect) {
         deviceMode = deviceModeToSelect;
      }

      public String getCDDVDDeviceMode() {
         return deviceMode;
      }
   }

   public static enum VIDEO_CARD_SETTING {
      SPECIFY_CUSTOM("Specify custom settings"), AUTO_DETECT("Auto-detect settings");

      private final String videoSetting;

      VIDEO_CARD_SETTING(String cardSettingToSelect) {
         videoSetting = cardSettingToSelect;
      }

      public String getVideoCardSetting() {
         return videoSetting;
      }
   }

   public static enum VMTOOLS_RESTART {
      DEFAULT(DEFAULT_VMTOOLS), RESETVMTOOLS(RESET), RESTART(RESTARTGOS);

      private final String restartVMtools;

      VMTOOLS_RESTART(String restartToSelet) {
         restartVMtools = restartToSelet;
      }

      public String getVMToolsRestartMode() {
         return restartVMtools;
      }
   }

   public static enum VMTOOLS_SHUTDOWN {
      DEFAULT(DEFAULT_VMTOOLS), POWER_OFF_TOOLS(POWER_OFF), SHUT_DOWN(SHUTDOWNGOS);

      private final String shutDownVMTools;

      VMTOOLS_SHUTDOWN(String shutDownToSelect) {
         shutDownVMTools = shutDownToSelect;
      }

      public String getVMToolsShutdownMode() {
         return shutDownVMTools;
      }
   }

   public static enum VMTOOLS_SUSPEND {
      Default, Suspend;
   }

   public static enum VM_VERSION {
      VM_VERSION_4("vmx-04"), VM_VERSION_7("vmx-07"), VM_VERSION_8("vmx-08"), VM_VERSION_9(
            "vmx-09");

      private final String vmVersion;

      VM_VERSION(String version) {
         vmVersion = version;
      }

      public String getVmVersion() {
         return vmVersion;
      }
   }

   public static enum CPUID_MASK {
      HIDE("Hide the NX/XD flag from guest"), EXPOSE("Expose the NX/XD flag to guest");

      private final String idMask;

      CPUID_MASK(String idMask) {
         this.idMask = idMask;
      }

      public String getMask() {
         return idMask;
      }
   }

   public static enum CPU_HT_SHARING {
      Any, None, Internal;
   }

   public static enum CPU_MMU_VIRTUALIZATION {
      AUTOMATIC("Automatic"), SW_CPU_AND_MMU("Software CPU and MMU"), HW_CPU_SW_MMU(
            "Hardware CPU, Software MMU"), HW_CPU_AND_MMU("Hardware CPU and MMU");

      private final String cpummuVirtualization;

      CPU_MMU_VIRTUALIZATION(String cpummuVirtualization) {
         this.cpummuVirtualization = cpummuVirtualization;
      }

      public String getVirtualization() {
         return cpummuVirtualization;
      }
   }

   public static enum HERTZ_UNIT {
      GHz, MHz;
   }

   public static enum DEVICE_NODE {
      SCSI, IDE;
   }

   public static enum NEW_HDD_PROVISIONING {
      THIN, FLAT;
   }

   public static enum CD_DVD_DEVICE_TYPE {
      CD_DVD_DEVICE_TYPE_CLIENT_DEVICE("Client Device"), CD_DVD_DEVICE_TYPE_HOST_DEVICE(
            "Host Device"), CD_DVD_DEVICE_TYPE_DATASTORE_ISO_FILE("Datastore ISO File");

      private final String cddvdDeviceType;

      CD_DVD_DEVICE_TYPE(String cddvdDevice) {
         cddvdDeviceType = cddvdDevice;
      }

      public String getCDDVDDeviceType() {
         return cddvdDeviceType;
      }
   }

   public static enum CD_DVD_DEVICE_MODE {
      EMULATE("Emulate IDE"), PASSTHROUGH("Passthrough IDE");

      private final String CDDVDDeviceMode;

      CD_DVD_DEVICE_MODE(String cddvcdeviceMode) {
         CDDVDDeviceMode = cddvcdeviceMode;
      }

      public String getMode() {
         return CDDVDDeviceMode;
      }
   }

   public static enum DEBUG_STATISTICS {
      RUN_NORMALLY("Run normally"), RECORD_DEBUG("Record Debugging Information"), RECORD_STATISTICS(
            "Record Statistics");

      private final String debugstats;

      DEBUG_STATISTICS(String debugStatistics) {
         debugstats = debugStatistics;
      }

      public String getDebugStats() {
         return debugstats;
      }
   }

   public static enum SWAPFILE_LOCATION {
      defaultOption, vmOption, hostOption;
   }

   public static enum NEW_HDD_PROVISIONING_TYPE {
      THIN("Thin Provision"), THICK_LAZYZEROED("Thick Provision Lazy Zeroed"), THICK_EAGERZEROED(
            "Thick Provision Eager Zeroed");

      private final String provType;

      NEW_HDD_PROVISIONING_TYPE(String provHDDType) {
         provType = provHDDType;
      }

      public String getdiskType() {
         return provType;
      }
   }

   // Networking

   // Related Items network table columns

   public static enum RI_NETWORKS_TABLE_COLUMN {
      NAME_ID(SIMPLE_SEARCH_NAME_COLUMN), NUMBER_OF_PORTS_ID(NUMBER_OF_PORTS_COLUMN);
      private final String strValue;

      private RI_NETWORKS_TABLE_COLUMN(String strValue) {
         this.strValue = strValue;
      }

      public String getValue() {
         return strValue;
      }
   }

   // DRM
   public enum DatastoreShow {
      CONNECT_TO_ALL_HOSTS("Show datastores connected to all hosts"), CONNECT_TO_SOME_HOSTS(
            "Show datastores connected to some hosts"), ALL("Show all datastores");

      private final String text;

      private DatastoreShow(String text) {
         this.text = text;
      }

      public String getText() {
         return text;
      }
   }

   public enum Time {
      MINUTES("Minutes", "unitMinutes"), HOURS("Hours", "unitHours"), DAYS("Days",
            "unitDays");

      private final String text;
      private final String unitString;

      private Time(String text, String unitString) {
         this.text = text;
         this.unitString = unitString;
      }

      public String getText() {
         return text;
      }

      public String getUnitString() {
         return unitString;
      }
   }

   public enum UiComponentProperty {

      CURRENT_STATE("currentState", "State_Expanded"), VALUE("value", ""), ERROR_STRING(
            "errorString", ""), TEXT("text", ""), LABEL("label", ""), RESIZABLE(
            "resizable", ""), ENABLED("enabled", ""), EDITABLE("editable", ""), VISIBLE(
            "visible", ""), ADAPTER_STATUS_BLOCK("adapterStatusBlock.baselinePosition",
            ""), ADAPTER_STATUS_BLOCK_LABEL("adapterStatusBlock.titleLabel.text", ""), ISCSI_GENERAL_BLOCK(
            "generalBlock.baselinePosition", ""), ISCSI_GENERAL_BLOCK_LABEL(
            "generalBlock.titleLabel.text", ""), ISCSI_AUTH_BLOCK(
            "authenticationBlock.baselinePosition", ""), ISCSI_AUTH_BLOCK_LABEL(
            "authenticationBlock.titleLabel.text", ""), ISCSI_NETWORK_INTERFACE_BLOCK(
            "networkInterfaceBlock.baselinePosition", ""), ISCSI_NETWORK_INTERFACE_BLOCK_LABEL(
            "networkInterfaceBlock.titleLabel.text", ""), ISCSI_IP_DNS_BLOCK(
            "ipAndDnsBlock.baselinePosition", ""), ISCSI_IP_DNS_BLOCK_LABEL(
            "ipAndDnsBlock.titleLabel.text", ""), NAME("Name", ""), SELECTED_INDEX(
            "selectedIndex", ""), DATA_SELECTED("data.selected", ""), TITLE("title", ""), SELECTED_ITEM_NAME(
            "selectedItem.name", ""), LOCKED_COLUMN_COUNT("lockedColumnCount", ""), SELECTED_TEXT(
            "selectedText", "");

      /**
       * Name of the property.
       */
      private final String name;

      /**
       * Value of the property.
       */
      private final String value;

      private UiComponentProperty(String name, String value) {
         this.name = name;
         this.value = value;
      }

      public String getName() {
         return name;
      }

      public String getValue() {
         return value;
      }
   }

   // Order of the items in this enum should be the same as default order of
   // columns in in Ports table
   public static enum PORTS_TABLE_COLUMNS {
      PORT_ID(PORT_ID_COLUMN_NAME), NAME(NAME_COLUMN_NAME), CONNECTEE(
            CONNECTEE_COLUMN_NAME), RUNTIME_MAC_ADDRESS(RUNTIME_MAC_ADDRESS_COLUMN_NAME), PORT_GROUP(
            PORT_GROUP_COLUMN_NAME), DIRECT_PATH_IO(DIRECT_PATH_IO_COLUMN_NAME), STATE(
            STATE_COLUMN_NAME), VLAN_ID(VLAN_ID_COLUMN_NAME), BROADCAST_INGRESS_TRAFFIC(
            BROADCAST_INGRESS_TRAFFIC_COLUMN_NAME), BROADCAST_EGRESS_TRAFFIC(
            BROADCAST_EGRESS_TRAFFIC_COLUMN_NAME);
      private final String strValue;

      private PORTS_TABLE_COLUMNS(String strValue) {
         this.strValue = strValue;
      }

      public String getValue() {
         return strValue;
      }
   }

   public static enum TOC_EDIT_DVUPLINK_PORTGROUP_SETTINGS {
      GENERAL, ADVANCED, VLAN, LACP, MONITORING, MISCELLANEOUS
   }

   public static enum TOC_EDIT_DVUPLINK_PORT_SETTINGS {
      PROPERTIES, VLAN, LACP, MONITORING, MISCELLANEOUS
   }

   public static enum TOC_EDIT_DVPORTGROUP_SETTINGS {
      GENERAL, ADVANCED, SECURITY, TRAFFIC_SHAPING, VLAN, TEAMING_AND_FAILOVER, MONITORING, MISCELLANEOUS
   }

   public static enum TOC_EDIT_PORT_SETTINGS {
      PROPERTIES, SECURITY, TRAFFIC_SHAPING, VLAN, TEAMING_AND_FAILOVER, MONITORING, MISCELLANEOUS
   }

   // Edit Host Profile - Networking enums
   public static enum IPV4_ADDRESS_DROPDOWN {
      APPLY_SPECIFIED_IPV4_CONFIG, USER_SPECIFIED_IPV4, USE_DHCP, PROMPT_THE_USER_IF_NO_DEFAULT
   }

   public static enum STATIC_IPV6_ADDRESS_DROPDOWN {
      USER_SPECIFIED_IPV6, PROMPT_THE_USER_IF_NO_DEFAULT, USER_MUST_CHOOSE_POLICY_OPTION
   }

   public static enum LICENSE_VIEW {
      LV_LICENSES(LIC_VIEW_LICENSES), LV_REPORTS(LIC_VIEW_LICENSE_REPORTS);

      private final String title;

      private LICENSE_VIEW(String title) {
         this.title = title;
      }

      public String getTitle() {
         return title;
      }
   }

   /**
    * Enum for tabs in Licenses Management view
    */
   public static enum LICENSES_TABS {
      // The order should be same as in the UI since we're using ordinal() for
      // tab index
      LT_LICENSE_KEYS, LT_PRODUCTS, LT_VCS, LT_HOSTS, LT_SOLUTUIONS;

      private final String tabIndex;

      private LICENSES_TABS() {
         this.tabIndex = "" + ordinal();
      }

      public String getTabIndex() {
         return tabIndex;
      }
   }

   /**
    * Enum for License Keys list columns
    */
   public static enum LICENSE_KEYS_LIST_COLUMN {
      LICENSE_KEY(LIC_COL_LICENSE_KEY), PRODUCT(LIC_COL_PRODUCT), USAGE(LIC_COL_USAGE), CAPACITY(
            LIC_COL_CAPACITY), SECOND_USAGE(LIC_COL_SECOND_USAGE), SECOND_CAPACITY(
            LIC_COL_SECOND_CAPACITY), LABEL(LIC_LICENSE_KEYS_COL_LABEL), EXPIRES(
            LIC_COL_EXPIRES), ASSIGNED(LIC_LICENSE_KEYS_COL_ASSIGNED);

      private final String columnName;

      private LICENSE_KEYS_LIST_COLUMN(String columnName) {
         this.columnName = columnName;
      }

      /**
       * @return the columnName
       */
      public String getColumnName() {
         return columnName;
      }

      @Override
      public String toString() {
         return getColumnName();
      }

      /**
       * Get LICENSE_KEYS_LIST_COLUMN object by name
       *
       * @param columnName String
       * @return LICENSE_KEYS_LIST_COLUMN object with matching name or null
       */
      public static LICENSE_KEYS_LIST_COLUMN getColumnByName(String columnName) {
         for (LICENSE_KEYS_LIST_COLUMN col : values()) {
            if (columnName.equals(col.getColumnName())) {
               return col;
            }
         }
         return null;
      }
   }

   /**
    * Enum for License Product list columns
    */
   public static enum LICENSE_PRODUCT_LIST_COLUMN {
      NAME(LIC_LICENSE_PRODUCTS_COL_NAME), USAGE(LIC_COL_USAGE), CAPACITY(
            LIC_COL_CAPACITY), SECOND_USAGE(LIC_COL_SECOND_USAGE), SECOND_CAPACITY(
            LIC_COL_SECOND_CAPACITY);

      private final String columnName;

      private LICENSE_PRODUCT_LIST_COLUMN(String columnName) {
         this.columnName = columnName;
      }

      /**
       * @return the columnName
       */
      public String getColumnName() {
         return columnName;
      }

      @Override
      public String toString() {
         return getColumnName();
      }

      /**
       * Get LICENSE_PRODUCT_LIST_COLUMN object by name
       *
       * @param columnName String
       * @return LICENSE_PRODUCT_LIST_COLUMN object with matching name or null
       */
      public static LICENSE_PRODUCT_LIST_COLUMN getColumnByName(String columnName) {
         for (LICENSE_PRODUCT_LIST_COLUMN col : values()) {
            if (columnName.equals(col.getColumnName())) {
               return col;
            }
         }
         return null;
      }
   }

   /**
    * Enum for License Hosts list columns
    */
   public static enum LICENSE_HOSTS_LIST_COLUMN {
      HOST(LIC_HOSTS_COL_HOST), STATE(LIC_HOSTS_COL_STATE), USAGE(LIC_COL_USAGE), SECOND_USAGE(
            LIC_COL_SECOND_USAGE), PRODUCT(LIC_COL_PRODUCT), LICENSE_KEY(
            LIC_COL_LICENSE_KEY), EXPIRES(LIC_COL_EXPIRES), CPUS(LIC_HOSTS_COL_CPUS), CORES(
            LIC_HOSTS_COL_CORES), VMS(LIC_HOSTS_COL_VMS);

      private final String columnName;

      private LICENSE_HOSTS_LIST_COLUMN(String columnName) {
         this.columnName = columnName;
      }

      /**
       * @return the columnName
       */
      public String getColumnName() {
         return columnName;
      }

      @Override
      public String toString() {
         return getColumnName();
      };

      /**
       * Get LICENSE_HOSTS_LIST_COLUMN object by name
       *
       * @param columnName String
       * @return LICENSE_HOSTS_LIST_COLUMN object with matching name or null
       */
      public static LICENSE_HOSTS_LIST_COLUMN getColumnByName(String columnName) {
         for (LICENSE_HOSTS_LIST_COLUMN col : values()) {
            if (columnName.equals(col.getColumnName())) {
               return col;
            }
         }
         return null;
      }
   }

   /**
    * Enum for License Hosts list columns
    */
   public static enum LICENSE_VC_LIST_COLUMN {
      VC_INSTANCE(LIC_VCS_COL_VC_INSTANCE), USAGE(LIC_COL_USAGE), PRODUCT(
            LIC_COL_PRODUCT), LICENSE_KEY(LIC_COL_LICENSE_KEY), EXPIRES(LIC_COL_EXPIRES);

      private final String columnName;

      private LICENSE_VC_LIST_COLUMN(String columnName) {
         this.columnName = columnName;
      }

      /**
       * @return the columnName
       */
      public String getColumnName() {
         return columnName;
      }

      @Override
      public String toString() {
         return getColumnName();
      };

      /**
       * Get LICENSE_VC_LIST_COLUMN object by name
       *
       * @param columnName String
       * @return LICENSE_VC_LIST_COLUMN object with matching name or null
       */
      public static LICENSE_VC_LIST_COLUMN getColumnByName(String columnName) {
         for (LICENSE_VC_LIST_COLUMN col : values()) {
            if (columnName.equals(col.getColumnName())) {
               return col;
            }
         }
         return null;
      }
   }


   /**
    * Enum for License Key Details view
    */
   public enum LICENSE_KEY_DETAILS {
      KEY, VCS, PRODUCT, COST_UNIT, VRAM_PER_CPU, LABEL, EXPIRATION;
   }

   /**
    * Enum for tabs in License Product Types
    */
   public static enum LICENSE_PRODUCT_TYPE {
      LPT_VC, LPT_HOST, LPT_SOLUTUION;
   }

   /* END OF ENUM SECTION */
   /*
    * --------------------------------------------------------------------------
    * -
    */

   /*
    *
    * Hashmap<STRING,STRING> All string maps go here
    */

   public static final HashMap<String, String> GOS_FAMILY_MAP =
         new HashMap<String, String>();
   static {
      GOS_FAMILY_MAP.put("windowsGuest", "Windows");
      GOS_FAMILY_MAP.put("linuxGuest", "Linux");
      GOS_FAMILY_MAP.put("netwareGuest", "Other");
      GOS_FAMILY_MAP.put("solarisGuest", "Other");
      GOS_FAMILY_MAP.put("otherGuestFamily", "Other");
      GOS_FAMILY_MAP.put("darwinGuestFamily", "Other");
   }

   public static final HashMap<String, String> LIC_API_TO_UI_COST_UNIT_MAP =
         new HashMap<String, String>();
   static {
      LIC_API_TO_UI_COST_UNIT_MAP
            .put(LIC_API_COST_UNIT_INSTANCE, LIC_COST_UNIT_INSTANCE);
      LIC_API_TO_UI_COST_UNIT_MAP.put(LIC_API_COST_UNIT_VM, LIC_COST_UNIT_VM);
      LIC_API_TO_UI_COST_UNIT_MAP.put(LIC_API_COST_UNIT_VRAM, LIC_COST_UNIT_VRAM);
      LIC_API_TO_UI_COST_UNIT_MAP.put(LIC_API_COST_UNIT_CPU, LIC_COST_UNIT_CPU);
      LIC_API_TO_UI_COST_UNIT_MAP.put(LIC_API_COST_UNIT_CPU6CORE, LIC_COST_UNIT_CPU);
      LIC_API_TO_UI_COST_UNIT_MAP.put(LIC_API_COST_UNIT_CPU12CORE, LIC_COST_UNIT_CPU);
   }

   public static final HashMap<String, String> LIC_PRODUCTS_CHART_COST_UNIT_MAP =
         new HashMap<String, String>();
   static {
      LIC_PRODUCTS_CHART_COST_UNIT_MAP.put(LIC_API_COST_UNIT_INSTANCE, "");
      LIC_PRODUCTS_CHART_COST_UNIT_MAP.put(LIC_API_COST_UNIT_VM, " (VM)");
      LIC_PRODUCTS_CHART_COST_UNIT_MAP.put(LIC_API_COST_UNIT_VRAM, " (vRAM)");
      LIC_PRODUCTS_CHART_COST_UNIT_MAP.put(
            LIC_API_COST_UNIT_CPU,
            " (unlimited cores per CPU)");
      LIC_PRODUCTS_CHART_COST_UNIT_MAP.put(
            LIC_API_COST_UNIT_CPU6CORE,
            " (1-6 cores per CPU)");
      LIC_PRODUCTS_CHART_COST_UNIT_MAP.put(
            LIC_API_COST_UNIT_CPU12CORE,
            " (1-12 cores per CPU)");
   }

   public static final HashMap<String, String> LIC_PRODUCT_DETAILS_COST_UNIT_MAP =
         new HashMap<String, String>();
   static {
      LIC_PRODUCT_DETAILS_COST_UNIT_MAP.put(
            LIC_API_COST_UNIT_INSTANCE,
            LIC_COST_UNIT_INSTANCES);
      LIC_PRODUCT_DETAILS_COST_UNIT_MAP.put(LIC_API_COST_UNIT_VM, LIC_COST_UNIT_VMS);
      LIC_PRODUCT_DETAILS_COST_UNIT_MAP
            .put(LIC_API_COST_UNIT_VRAM, LIC_COST_UNIT_GBVRAM);
      LIC_PRODUCT_DETAILS_COST_UNIT_MAP.put(LIC_API_COST_UNIT_CPU, LIC_COST_UNIT_CPUS);
      LIC_PRODUCT_DETAILS_COST_UNIT_MAP.put(
            LIC_API_COST_UNIT_CPU6CORE,
            LIC_COST_UNIT_CPUS);
      LIC_PRODUCT_DETAILS_COST_UNIT_MAP.put(
            LIC_API_COST_UNIT_CPU12CORE,
            LIC_COST_UNIT_CPUS);
   }

   // Mapping for vApp and RP - tabs and indexes
   public static final HashMap<String, String> RES_GROUPS_TAB_INDEXES_MAP =
         new HashMap<String, String>();
   static {
      RES_GROUPS_TAB_INDEXES_MAP.put(GETTING_STARTED_TAB, "0");
      RES_GROUPS_TAB_INDEXES_MAP.put(SUMMARY_TAB, "1");
      RES_GROUPS_TAB_INDEXES_MAP.put(MONITOR_TAB, "2");
      RES_GROUPS_TAB_INDEXES_MAP.put(MANAGE_TAB, "3");
      RES_GROUPS_TAB_INDEXES_MAP.put(RELATED_ITEMS_TAB, "4");
   }

   // Mapping for VM Lists - columns and indexes
   public static final HashMap<String, String> VM_LIST_INDEXES_MAP =
         new HashMap<String, String>();
   static {
      VM_LIST_INDEXES_MAP.put(VM_NAME_HEADER, "0");
      VM_LIST_INDEXES_MAP.put(VM_STATE_HEADER, "1");
      VM_LIST_INDEXES_MAP.put(VM_STATUS_HEADER, "2");
      VM_LIST_INDEXES_MAP.put(VM_PROVISIONED_SPACE, "3");
      VM_LIST_INDEXES_MAP.put(VM_USED_SPACE, "4");
      VM_LIST_INDEXES_MAP.put(VM_HOST_CPU_HEADER, "5");
      VM_LIST_INDEXES_MAP.put(VM_HOST_MEMO_HEADER, "6");
      VM_LIST_INDEXES_MAP.put(VM_GUEST_MEMO_HEADER, "7");
   }

   // Map between API power states and UI representation
   public static final HashMap<String, String> API_NGC_UI_POWER_STATES_MAP =
         new HashMap<String, String>();
   static {
      API_NGC_UI_POWER_STATES_MAP.put(POWEREDOFF, POWERED_OFF);
      API_NGC_UI_POWER_STATES_MAP.put(POWEREDON, POWERED_ON);
      API_NGC_UI_POWER_STATES_MAP.put(POWERSUSPENDED, SUSPENDED);
   }

   public static final HashMap<String, MAIN_NAVIGATION_TABS> MAIN_TAB_NAMES_MAP =
         new HashMap<String, MAIN_NAVIGATION_TABS>();
   static {
      MAIN_TAB_NAMES_MAP.put(GETTING_STARTED_TAB, MAIN_TABS_GETTING_STARTED);
      MAIN_TAB_NAMES_MAP.put(SUMMARY_TAB, MAIN_TABS_SUMMARY);
      MAIN_TAB_NAMES_MAP.put(MANAGE_TAB, MAIN_TABS_MANAGE);
      MAIN_TAB_NAMES_MAP.put(MONITOR_TAB, MAIN_TABS_MONITOR);
      MAIN_TAB_NAMES_MAP.put(RELATED_ITEMS_TAB, MAIN_TABS_RELATED_ITEMS);
   }

   // HashMap for Views based on objects
   public static final HashMap<String, OBJ_NAV_TREE_NODE_VIEW> NODE_VIEW_BASED_ON_OBJECT =
         new HashMap<String, OBJ_NAV_TREE_NODE_VIEW>();
   static {
      // Nodes directly under Virtual Centers
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_FOLDER.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_VIRTUAL_CENTERS);

      // Nodes directly under Datacenters
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_DATACENTER.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_DATACENTERS);

      // Nodes directly under Hosts
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_HOST.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_HOSTS);

      // Nodes directly under Clusters
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_CLUSTER.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_CLUSTERS);

      // Nodes directly under Resource Pools
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_RESOURCE_POOL.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_RESOURCE_POOLS);

      // Nodes directly under Virtual Machines
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_VM.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_VIRTUAL_MACHINES);

      // Nodes directly under VM Templates
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_TEMPLATE.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_VMS_TEMPLATES);

      // Nodes directly under VAPPS
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_VAPP.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_VAPPS);

      // Nodes directly under Standard networks
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_STANDARD_SWITCH.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_STANDARD_NETWORKS);

      // Nodes directly under Distributed portgroups
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_DV_PORTGROUP.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_DISTRIBUTED_PORT_GROUPS);
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_DV_UPLINK.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_DISTRIBUTED_UPLINK_PORT_GROUPS);

      // Nodes directly under Distributed Switches
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_DV_SWITCH.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_DISTRIBUTED_SWITCHES);

      // Nodes directly under Datastores
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_DATASTORE.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_DATASTORES);

      // Nodes directly under Datastore clusters
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_STORAGE_POD.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_DATASTORE_CLUSTERS);

      // Nodes directly under Host Profiles
      NODE_VIEW_BASED_ON_OBJECT.put(
            NODE_TYPE.NT_HOST_PROFILE.getNodeType(),
            OBJ_NAV_TREE_NODE_VIEW.VIEW_HOST_PROFILES);

   }

   // List for mor to be excluded for object Navigation Object navigator
   public static final List<String> FOLDER_OBJECTNAV_EXCLUSION_LIST =
         new ArrayList<String>();
   static {
      FOLDER_OBJECTNAV_EXCLUSION_LIST.add("host");
      FOLDER_OBJECTNAV_EXCLUSION_LIST.add("datastore");
      FOLDER_OBJECTNAV_EXCLUSION_LIST.add("vm");
      FOLDER_OBJECTNAV_EXCLUSION_LIST.add("network");
   }

   // HashMap for Node Ids under the DataCenter view
   public static final HashMap<String, String> NODE_IDS_DATACENTERS_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_DATACENTERS_VIEW.put(
            NODE_TYPE.NT_DATACENTER.getNodeType(),
            IDConstants.DATACENTERS_DATAPROVIDER_ID_DATACENTERS_VIEW);
   }

   // HashMap for Node Ids under the Hosts view
   public static final HashMap<String, String> NODE_IDS_HOSTS_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_HOSTS_VIEW.put(
            NODE_TYPE.NT_HOST.getNodeType(),
            IDConstants.HOSTS_DATAPROVIDER_ID_HOSTS_VIEW);
   }

   // HashMap for Node Ids under the Virtual Centers view
   public static final HashMap<String, String> NODE_IDS_VIRTUAL_CENTERS_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_VIRTUAL_CENTERS_VIEW.put(
            NODE_TYPE.NT_FOLDER.getNodeType(),
            IDConstants.VC_DATAPROVIDER_ID);
      NODE_IDS_VIRTUAL_CENTERS_VIEW.put(
            NODE_TYPE.NT_RESOURCE_POOL.getNodeType(),
            IDConstants.VC_DATAPROVIDER_ID);
   }

   // HashMap for Node Ids under the Resource Pools view
   public static final HashMap<String, String> NODE_IDS_RESOURCE_POOLS_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_RESOURCE_POOLS_VIEW.put(
            NODE_TYPE.NT_RESOURCE_POOL.getNodeType(),
            IDConstants.RESOURCE_POOLS_DATAPROVIDER_ID_RESOURCE_POOLS_VIEW);
   }

   // HashMap for Node Ids under the Clusters view
   public static final HashMap<String, String> NODE_IDS_CLUSTERS_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_CLUSTERS_VIEW.put(
            NODE_TYPE.NT_CLUSTER.getNodeType(),
            IDConstants.CLUSTERS_DATAPROVIDER_ID_CLUSTERS_VIEW);
   }

   // HashMap for Node Ids under the Datastores view
   public static final HashMap<String, String> NODE_IDS_DATASTORES_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_DATASTORES_VIEW.put(
            NODE_TYPE.NT_DATASTORE.getNodeType(),
            IDConstants.DATASTORES_DATAPROVIDER_ID_DATASTORES_VIEW);
   }

   // HashMap for Node Ids under the Datastore Cluster view
   public static final HashMap<String, String> NODE_IDS_DATASTORE_CLUSTERS_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_DATASTORE_CLUSTERS_VIEW.put(
            NODE_TYPE.NT_STORAGE_POD.getNodeType(),
            IDConstants.DATASTORE_CLUSTERS_DATAPROVIDER_ID_DATASTORE_CLUSTERS_VIEW);
   }

   // HashMap for Node Ids under the Standard Networks view
   public static final HashMap<String, String> NODE_IDS_STANDARD_NETWORKS_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_STANDARD_NETWORKS_VIEW.put(
            NODE_TYPE.NT_STANDARD_PORTGROUP.getNodeType(),
            IDConstants.STANDARD_NETWORKS_DATAPROVIDER_ID_NETWORKING_VIEW);
   }

   // Hashmap for Node Ids under the VM Network node view
   public static final HashMap<String, String> NODE_IDS_VM_NETWORK_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_VM_NETWORK_VIEW.put(
            NODE_TYPE.NT_VM.getNodeType(),
            IDConstants.VMS_ID_VM_NETWORK_VIEW);
      NODE_IDS_VM_NETWORK_VIEW.put(
            NODE_TYPE.NT_HOST.getNodeType(),
            IDConstants.HOSTS_ID_VM_NETWORK_VIEW);
   }

   // HashMap for Node Ids under the Distributed switches view
   public static final HashMap<String, String> NODE_IDS_DISTRIBUTED_SWITCHES_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_DISTRIBUTED_SWITCHES_VIEW.put(
            NODE_TYPE.NT_DV_SWITCH.getNodeType(),
            IDConstants.DISTRIBUTED_SWITCHES_ID_DISTRIBUTED_SWITCHES_VIEW);
   }

   // HashMap for Node Ids under the Distributed port group view
   public static final HashMap<String, String> NODE_IDS_DISTRIBUTED_PORT_GROUPS_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_DISTRIBUTED_PORT_GROUPS_VIEW.put(
            NODE_TYPE.NT_DV_PORTGROUP.getNodeType(),
            IDConstants.DISTRIBUTED_PORT_GROUPS_ID_DISTRIBUTED_PORT_GROUPS_VIEW);
   }

   // HashMap for Node Ids under the Distributed uplinkport group view
   public static final HashMap<String, String> NODE_IDS_DISTRIBUTED_UPLINK_PORT_GROUPS_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_DISTRIBUTED_UPLINK_PORT_GROUPS_VIEW.put(
            NODE_TYPE.NT_DV_PORTGROUP.getNodeType(),
            IDConstants.DISTRIBUTED_PORT_GROUPS_ID_DISTRIBUTED_PORT_GROUPS_VIEW);
   }

   // HashMap for Node Ids under the VMs view
   public static final HashMap<String, String> NODE_IDS_VMS_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_VMS_VIEW.put(NODE_TYPE.NT_VM.getNodeType(), IDConstants.VMS_ID_VMS_VIEW);
   }

   // HashMap for Node Ids under the vApps view
   public static final HashMap<String, String> NODE_IDS_VAPPS_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_VAPPS_VIEW.put(
            NODE_TYPE.NT_VAPP.getNodeType(),
            IDConstants.VAPPS_ID_VAPPS_VIEW);
   }

   // HashMap for Node Ids under the Compute view
   public static final HashMap<String, String> NODE_IDS_COMPUTE_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_COMPUTE_VIEW.put(
            NODE_TYPE.NT_HOST.getNodeType(),
            IDConstants.HOSTS_DATAPROVIDER_ID_COMPUTE_VIEW);
      NODE_IDS_COMPUTE_VIEW.put(
            NODE_TYPE.NT_CLUSTER.getNodeType(),
            IDConstants.CLUSTERS_DATAPROVIDER_ID_COMPUTE_VIEW);
   }

   // HashMap for Node Ids under the Networking view
   public static final HashMap<String, String> NODE_IDS_NETWORKING_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_NETWORKING_VIEW.put(
            NODE_TYPE.NT_STANDARD_PORTGROUP.getNodeType(),
            IDConstants.STANDARD_NETWORKS_DATAPROVIDER_ID_NETWORKING_VIEW);
      NODE_IDS_NETWORKING_VIEW.put(
            NODE_TYPE.NT_DV_SWITCH.getNodeType(),
            IDConstants.DISTRIBUTED_SWITCHES_DATAPROVIDER_ID_NETWORKING_VIEW);
      NODE_IDS_NETWORKING_VIEW.put(
            NODE_TYPE.NT_DV_PORTGROUP.getNodeType(),
            IDConstants.DISTRIBUTED_PORT_GROUPS_DATAPROVIDER_ID_NETWORKING_VIEW);
      NODE_IDS_NETWORKING_VIEW.put(
            NODE_TYPE.NT_DV_UPLINK.getNodeType(),
            IDConstants.UPLINK_PORT_GROUPS_DATAPROVIDER_ID_NETWORKING_VIEW);
      NODE_IDS_NETWORKING_VIEW.put(
            NODE_TYPE.NT_HOST.getNodeType(),
            IDConstants.HOSTS_DATAPROVIDER_ID_NETWORKING_VIEW);
   }

   // HashMap for Node Ids under the Storage view
   public static final HashMap<String, String> NODE_IDS_STORAGE_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_STORAGE_VIEW.put(
            NODE_TYPE.NT_DATASTORE.getNodeType(),
            IDConstants.DATASTORES_DATAPROVIDER_ID_STORAGE_VIEW);
      NODE_IDS_STORAGE_VIEW.put(
            NODE_TYPE.NT_STORAGE_POD.getNodeType(),
            IDConstants.DATASTORE_CLUSTERS_DATAPROVIDER_ID_STORAGE_VIEW);
   }

   // HashMap for Node Ids under the VM Templates view
   public static final HashMap<String, String> NODE_IDS_VMS_TEMPLATES_VIEW =
         new HashMap<String, String>();
   static {

      NODE_IDS_VMS_TEMPLATES_VIEW.put(
            NODE_TYPE.NT_TEMPLATE.getNodeType(),
            IDConstants.VM_TEMPLATES_ID_VM_TEMPLATES_VIEW);

   }

   // HashMap for Node Ids under the Datacenter node view
   public static final HashMap<String, String> NODE_IDS_DATACENTER =
         new HashMap<String, String>();
   static {

      NODE_IDS_DATACENTER.put(
            NODE_TYPE.NT_VM.getNodeType(),
            IDConstants.VM_ID_DATACENTER_VIEW);

   }

   // HashMap for Node Ids under the Datastore node view
   public static final HashMap<String, String> NODE_IDS_DATASTORE =
         new HashMap<String, String>();
   static {

      NODE_IDS_DATASTORE.put(
            NODE_TYPE.NT_VM.getNodeType(),
            IDConstants.VM_ID_DATASTORE_VIEW);

   }

   // HashMap for Contents Node Ids under objects
   public static final HashMap<String, String> CONTENTS_NODE_IDS_UNDER_OBJECT =
         new HashMap<String, String>();
   static {
      CONTENTS_NODE_IDS_UNDER_OBJECT.put(
            NODE_TYPE.NT_FOLDER.getNodeType(),
            IDConstants.VMS_DATAPROVIDER_ID_STORAGE_VIEW);
      CONTENTS_NODE_IDS_UNDER_OBJECT.put(
            NODE_TYPE.NT_TEMPLATE.getNodeType(),
            IDConstants.VM_TEMPLATES_DATAPROVIDER_ID_STORAGE_VIEW);
      CONTENTS_NODE_IDS_UNDER_OBJECT.put(
            NODE_TYPE.NT_VAPP.getNodeType(),
            IDConstants.VAPPS_DATAPROVIDER_ID_STORAGE_VIEW);
   }

   // HashMap for Node Ids under the Host Profiles view
   public static final HashMap<String, String> NODE_IDS_HOST_PROFILES_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_HOST_PROFILES_VIEW.put(
            NODE_TYPE.NT_HOST_PROFILE.getNodeType(),
            IDConstants.HOST_PROFILES_ID_HOST_PROFILES_VIEW);
   }

   // HashMap for Node Ids under the Networking view
   public static final HashMap<String, String> NODE_IDS_HOST_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_HOST_VIEW.put(
            NODE_TYPE.NT_RESOURCE_POOL.getNodeType(),
            IDConstants.CHILD_RPS_DATAPROVIDER_ID_HOST_VIEW);
      NODE_IDS_HOST_VIEW.put(
            NODE_TYPE.NT_VM.getNodeType(),
            IDConstants.VMS_DATAPROVIDER_ID_HOST_VIEW);
      NODE_IDS_HOST_VIEW.put(
            NODE_TYPE.NT_DATASTORE.getNodeType(),
            IDConstants.DATASTORES_DATAPROVIDER_ID_HOST_VIEW);
   }

   // Hashmap for Node Ids under the Virtual Center node view
   public static final HashMap<String, String> NODE_IDS_VIRTUAL_CENTER_NODE_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_VIRTUAL_CENTER_NODE_VIEW.put(
            NODE_TYPE.NT_VM.getNodeType(),
            IDConstants.VMS_ID_VIRTUAL_CENTER_VIEW);
      NODE_IDS_VIRTUAL_CENTER_NODE_VIEW.put(
            NODE_TYPE.NT_TEMPLATE.getNodeType(),
            IDConstants.VM_TEMPLATE_ID_VIRTUAL_CENTER_VIEW);
      NODE_IDS_VIRTUAL_CENTER_NODE_VIEW.put(
            NODE_TYPE.NT_VAPP.getNodeType(),
            IDConstants.VAPPS_ID_VIRTUAL_CENTER_VIEW);
      NODE_IDS_VIRTUAL_CENTER_NODE_VIEW.put(
            NODE_TYPE.NT_DV_PORTGROUP.getNodeType(),
            IDConstants.DV_PORT_GROUP_ID_VIRTUAL_CENTER_VIEW);
   }

   // Hashmap for Node Ids under the Distributed Virtual Switch
   public static final HashMap<String, String> NODE_IDS_DISTRIBUTED_VIRTUAL_SWITCH_NODE_VIEW =
         new HashMap<String, String>();
   static {
      NODE_IDS_DISTRIBUTED_VIRTUAL_SWITCH_NODE_VIEW.put(
            NODE_TYPE.NT_DV_UPLINK.getNodeType(),
            IDConstants.UPLINK_PORTGROUP_ID_VDS_VIEW);
   }

   /* End of Hashmap<STRING,STRING> */
   /*
    * --------------------------------------------------------------------------
    * -
    */

   /*
    *
    * String array section
    */

   // Path ids
   public static String MANAGE_MONITOR_PATH_IDS[] = { OBJ_NAV_MANAGE_MONITOR_NODE_ITEM };
   public static String VI_HOME_PATH_IDS[] =
         { OBJ_NAV_VIRTUAL_INFRASTRUCTURE_NODE_ITEM };
   public static String RULES_PROFILES_PATH_IDS[] = { OBJ_NAV_RULES_PROFILE_NODE_ITEM };
   public static String ADMINISTRATION_PATH_IDS[] = { OBJ_NAV_ADMINISTRATION_NODE_ITEM };
   public static String LICENSING_PATH_IDS[] = { OBJ_NAV_ADMINISTRATION_NODE_ITEM };
   public static String LICENSES_PATH_IDS[] = { OBJ_NAV_LICENSES_NODE_ITEM };
   public static String LICENSE_REPORTS_PATH_IDS[] =
         { OBJ_NAV_LICENSE_REPORTS_NODE_ITEM };
   public static String SEARCH_PATH_IDS[] = { OBJ_NAV_SEARCH_NODE_ITEM };
   public static String SAVEDSEARCH_PATH_IDS[] = { OBJ_NAV_SAVEDSEARCH_NODE_ITEM };
   public static String TASKS_PATH_IDS[] = { OBJ_NAV_MANAGE_MONITOR_NODE_ITEM,
         OBJ_NAV_TASKS_NODE_ITEM };
   public static String PLUGIN_MANAGEMENT_PATH_IDS[] = {
         OBJ_NAV_ADMINISTRATION_NODE_ITEM, OBJ_NAV_SOLUTION_PLUGIN_MANAGER_NODE_ITEM };
   public static String SSO_CONFIGURATION_PATH_IDS[] = {
         OBJ_NAV_ADMINISTRATION_NODE_ITEM, OBJ_NAV_SSO_CONFIGURATION_NODE_ITEM };
   public static String USERS_GROUPS_PATH_IDS[] = { OBJ_NAV_ADMINISTRATION_NODE_ITEM,
         OBJ_NAV_SSO_USERS_GROUPS_NODE_ITEM };
   public static String VC_RELATED_ITEMS_IDS[] = { RELATED_ITEMS_CONTENTS,
         RELATED_ITEMS_DATACENTERS, RELATED_ITEMS_CLUSTERS, RELATED_ITEMS_HOSTS,
         RELATED_ITEMS_VIRTUAL_MACHINES, RELATED_ITEMS_VM_TEMPLATES,
         RELATED_ITEMS_VAPPS, RELATED_ITEMS_DATASTORES,
         RELATED_ITEMS_DATASTORE_CLUSTERS, RELATED_ITEMS_STANDARD_NETWORKS,
         RELATED_ITEMS_DISTRIBUTED_SWITCHES, RELATED_ITEMS_DISTRIBUTED_PORT_GROUPS,
         RELATED_ITEMS_EXTENSIONS };
   public static String DC_RELATED_ITEMS_IDS[] = { RELATED_ITEMS_CONTENTS,
         RELATED_ITEMS_CLUSTERS, RELATED_ITEMS_HOSTS, RELATED_ITEMS_VIRTUAL_MACHINES,
         RELATED_ITEMS_VM_TEMPLATES, RELATED_ITEMS_VAPPS, RELATED_ITEMS_DATASTORES,
         RELATED_ITEMS_DATASTORE_CLUSTERS, RELATED_ITEMS_STANDARD_NETWORKS,
         RELATED_ITEMS_DISTRIBUTED_SWITCHES, RELATED_ITEMS_DISTRIBUTED_PORT_GROUPS };
   public static String HOST_RELATED_ITEMS_IDS[] = {
         RELATED_ITEMS_CHILD_RESOURCE_POOLS_VMS, RELATED_ITEMS_VIRTUAL_MACHINES,
         RELATED_ITEMS_VAPPS, RELATED_ITEMS_DATASTORES, RELATED_ITEMS_NETWORKS,
         RELATED_ITEMS_DISTRIBUTED_SWITCHES };
   public static String CLUSTER_RELATED_ITEMS_IDS[] = {
         RELATED_ITEMS_CHILD_RESOURCE_POOLS_VMS, RELATED_ITEMS_HOSTS,
         RELATED_ITEMS_VIRTUAL_MACHINES, RELATED_ITEMS_VAPPS, RELATED_ITEMS_DATASTORES,
         RELATED_ITEMS_DATASTORE_CLUSTERS, RELATED_ITEMS_NETWORKS,
         RELATED_ITEMS_DISTRIBUTED_SWITCHES };
   public static String RP_RELATED_ITEMS_IDS[] = {
         RELATED_ITEMS_CHILD_RESOURCE_POOLS_VMS, RELATED_ITEMS_VIRTUAL_MACHINES,
         RELATED_ITEMS_VAPPS };
   public static String NETWORK_RELATED_ITEMS_IDS[] = { RELATED_ITEMS_VIRTUAL_MACHINES,
         RELATED_ITEMS_HOSTS };
   public static String VDS_RELATED_ITEMS_IDS[] = { RELATED_ITEMS_HOSTS,
         RELATED_ITEMS_VIRTUAL_MACHINES, RELATED_ITEMS_DISTRIBUTED_PORT_GROUPS,
         RELATED_ITEMS_UPLINK_PORT_GROUPS };
   public static String VDS_PORTGROUP_RELATED_ITEMS_IDS[] = { RELATED_ITEMS_HOSTS,
         RELATED_ITEMS_VIRTUAL_MACHINES };
   public static String VIRTUAL_MACHINE_RELATED_ITEMS_IDS[] = { RELATED_ITEMS_NETWORKS,
         RELATED_ITEMS_DATASTORES };
   public static String VDS_UPLINK_RELATED_ITEMS_IDS[] = { RELATED_ITEMS_HOSTS };
   public static String DATASTORES_RELATED_ITEMS_IDS[] = {
         RELATED_ITEMS_VIRTUAL_MACHINES, RELATED_ITEMS_HOSTS };
   public static String VM_TEMPLATE_RELATED_ITEMS_IDS[] = { RELATED_ITEMS_HOST };
   public static String VAPP_RELATED_ITEMS_IDS[] = {
         RELATED_ITEMS_CHILD_RESOURCE_POOLS_VMS, RELATED_ITEMS_VIRTUAL_MACHINES,
         RELATED_ITEMS_VAPPS, RELATED_ITEMS_DATASTORES, RELATED_ITEMS_NETWORKS };
   public static String HOST_PROFILE_RELATED_ITEMS_IDS[] = { RELATED_ITEMS_HOSTS };
   public static String ROLE_MANAGER_PATH_IDS[] = { OBJ_NAV_ADMINISTRATION_NODE_ITEM,
         OBJ_NAV_ROLE_MANAGER_NODE_ITEM };
   public static String TAG_MANAGER_PATH_IDS[] = { OBJ_NAV_TAG_MANAGER_NODE_ITEM };
   public static String VM_STORAGE_PROFILE_RELATED_ITEMS_IDS[] = {
         RELATED_ITEMS_VIRTUAL_MACHINES, RELATED_ITEMS_VM_TEMPLATES };

   public static final String[] POWER_SUBMENU = { "Power On", "Shut Down Guest",
         "Restart Guest", "Power Off", "Suspend", "Reset" };
   public static final String[] CONFIGURATION_SUBMENU = { "Rename",
         "Upgrade Virtual Hardware", "Edit Settings", "Install/Upgrade VMware Tools",
         "Unmount Tools Installer", };
   public static final String[] INVENTORY_SUBMENU = { "Migrate", "Convert to Template",
         "Clone", "Remove from Inventory", "Delete" };
   public static final String[] SNAPSHOTS_SUBMENU = { "Snapshot Manager...",
         "Take Snapshot...", "Revert to current snapshot" };

   // More Actions
   public static final String[] MORE_ACTION = { "Power On", "Suspend",
         "Shut Down Guest", "Power Off", "Restart Guest", "Reset", "Edit Settings",
         "Install/Upgrade Tools", "Unmount Tools Installer", "Migrate", "Rename",
         "Remove", "Delete", "Snapshot Manager...", "New Folder", "New Datacenter",
         "Shutdown" };

   /* End of String array section */
   /*
    * --------------------------------------------------------------------------
    * -
    */

   /*
    *
    * TODO Misc to be remediated section
    */

   // Constants used for DRS rebalance
   public static final int CPU_LOAD = 1;
   public static final int NO_LOAD = 4;
   public static final int MIGRATION_THRESHOLD_LEVEL_LOW = 1;
   public static final int MIGRATION_THRESHOLD_LEVEL_HIGH = 5;

   public static int MAX_CD_DRIVE_ALLOWED = 3;
   public static int MAX_FLOPPY_DRIVE_ALLOWED = 2;
   public static int MAX_SERIAL_PORT_ALLOWED = 4;
   public static int MAX_PARALLEL_PORT_ALLOWED = 3;
   public static int MAX_LSI_PORT_ALLOWED = 4;
   public static final int MAX_CHARS_NAME_LENGTH = 80;
   public static final int MAX_CHARS_NAME_LENGTH_ESX41 = 49;

   public static final boolean DEBUG = false;

   // Immediate
   public static final String DEFAULT_TIMEOUT_IMMEDIATE = "1";
   public static final int DEFAULT_TIMEOUT_IMMEDIATE_INT_VALUE = 500;
   public static final long DEFAULT_TIMEOUT_IMMEDIATE_LONG_VALUE = new Long(
         DEFAULT_TIMEOUT_IMMEDIATE).longValue();

   // one second
   public static final String DEFAULT_TIMEOUT_ONE_SECOND = "1000";
   public static final int DEFAULT_TIMEOUT_ONE_SECOND_INT_VALUE = 1000;
   public static final long DEFAULT_TIMEOUT_ONE_SECOND_LONG_VALUE = 1000L;

   // five seconds
   public static final int DEFAULT_TIMEOUT_FIVE_SECOND_INT_VALUE = 5000;

   // ten seconds
   public static final String DEFAULT_TIMEOUT_TEN_SECONDS = "10000";
   public static final long DEFAULT_TIMEOUT_TEN_SECONDS_LONG_VALUE = new Long(
         DEFAULT_TIMEOUT_TEN_SECONDS).longValue();
   public static final int DEFAULT_WAIT_ITERATIONS_THREE_MILLION = 30000000;

   public static final int DEFAULT_TIMEOUT_THREE_MINUTES = 180000;
   public static final int EDIT_OPTION_CPU = 1;
   public static final int EDIT_OPTION_MEMORY = 2;
   public static final int EDIT_OPTION_RESERVATION = 3;
   public static final int EDIT_OPTION_LIMIT = 4;
   public static final int EDIT_OPTION_SHARES = 5;
   public static final int EDIT_OPTION_SHARES_OPTIONS = 6;
   public static final int EDIT_OPTION_MEMORY_NEGATIVE = 7;
   public static final int EDIT_OPTION_LIMIT_NEGATIVE = 8;
   public static final int EDIT_OPTION_RESERVATION_NEGATIVE = 9;

   public static final int VM_STATE_SUSPEND = 2;

   public static final int DEFAULT_TIMEOUT_THREE_MINUTE = 180000;
   public static final int VM_STATE_POWER_OFF = 0;
   public static final int VM_STATE_POWER_ON = 1;
   public static final int VM_STATE_POWERED_OFF = 0;
   public static final int VM_STATE_POWERED_ON = 1;
   public static final int VM_STATE_SUSPENDED = 2;
   public static final int VM_STATE_RESET = 3;

   public static final int DEFAULT_TIMEOUT_ONE_MINUTE = 60000;
   public static final int DEFAULT_TIMEOUT_TWO_SECOND = 2000;
   public static final int POWER_STATE_FOR_VM_OFF = 0;
   public static final int POWER_STATE_FOR_VM_ON = 1;
   public static final int POWEREDOFF_VM_POWER_OFF_AGAIN = 2;
   public static final int POWEREDON_VM_POWER_ON_AGAIN = 3;
   public static final int PRIVILEGES_POWER_ON = 1;
   public static final int PRIVILEGES_POWER_OFF = 2;
   public static final int PRIVILEGES_POWER_ONOFFSUSPENDRESET = 3;
   public static final int PRIVILEGES_ONLY_POWER_ON = 4;
   public static final int PRIVILEGES_POWER_SUSPEND = 5;
   public static final int PRIVILEGES_POWER_RESET = 6;

   public static final int CAPACITY_GB = 1024;

   public static final String VSWITCH_EDIT_SETTINGS_TITLE = SPACE + DASH + SPACE
         + "Edit Settings";

   public static final String TURN_FT_ON_WARNING_DIALOG_WITH_DISK_WARNING =
         TURN_FT_ON_DISK_WARNING
               + "The memory reservation of this VM will be changed to the memory"
               + " size of the VM and maintained equal to it until Fault Tolerance"
               + " is turned off.\n\nDo you want to turn On Fault Tolerance?";
   public static final String TURN_FT_ON_WARNING_DIALOG_DRS_WITH_DISK_WARNING =
         TURN_FT_ON_DISK_WARNING
               + "The DRS automation level for this VM will change to disabled."
               + "\n\nThe memory reservation of this VM will be changed to the "
               + "memory size of the VM and maintained equal to it until Fault "
               + "Tolerance is turned off.\n\nDo you want to turn On Fault "
               + "Tolerance?";

   // Constant for default value of CPU Shares and Memory Shares
   public static final int DEFAULT_CPU_SHARES = 4000;
   public static final int DEFAULT_MEMORY_SHARES = 163840;

   public static final long SMALL_DISK_CAPACITY_IN_KB = 2048;

   public static final long BYTE_NUMBER = 1;
   public static final long KILO_BYTE_NUMBER = BYTE_NUMBER * 1024;
   public static final long MEGA_BYTE_NUMBER = KILO_BYTE_NUMBER * 1024;
   public static final long GIGA_BYTE_NUMBER = MEGA_BYTE_NUMBER * 1024;
   public static final long TERA_BYTE_NUMBER = GIGA_BYTE_NUMBER * 1024;

   public static final long HERTZ_NUMBER = 1;
   public static final long KILO_HERTZ_NUMBER = HERTZ_NUMBER * 1000;
   public static final long MEGA_HERTZ_NUMBER = KILO_HERTZ_NUMBER * 1000;
   public static final long GIGA_HERTZ_NUMBER = MEGA_HERTZ_NUMBER * 1000;
   public static final long TERA_HERTZ_NUMBER = GIGA_HERTZ_NUMBER * 1000;

   public static final double unitByte = 1024.0;
   public static final int limit = -1;

   // Host Security profile

   public static final String[] HOST_SERVICES_LIST = { "Direct Console UI",

   "ESXi", "SSH", "lbtd",

   "Local Security Authentication Server (Active Directory Service)",

   "I/O Redirector (Active Directory Service)",

   "Network Login Server (Active Directory Service)", "NTP Daemon",

   "CIM Server", "vSphere High Availability Agent", "vpxa" };

   public static final String[] HOST_FIREWALLSERVICES_LIST =
         { "SSH Client", "SSH Server", "CIM Server", "CIM Secure Server", "CIM SLP",
               "DHCPv6", "DVFilter", "DVSSync", "HBR", "IKED", "NFC", "WOL",
               "Active Directory All", "DHCP Client", "DNS Client", "Fault Tolerance",
               "FTP Client", "gdbserver", "httpClient", "Software iSCSI Client",
               "NFS Client", "NTP Client", "VM serial port connected over network",
               "SNMP Server", "syslog", "vCenter Update Manager", "vMotion",
               "VM serial port connected to vSPC", "vSphere Client", "vprobeServer",
               "VMware vCenter Agent", "vSphere Web Access" };

   public static final String[] HOST_IMAGEPROFILE_ACCEPTANCE_LEVELS = {
         "Community Supported", "Partner Supported", "VMware Accepted",
         "VMware Certified" };

   public static enum EDIT_HOST_SERVICE_OPERATIONS {
      OPERATION_START, OPERATION_RESTART, OPERATION_STOP
   };

   // Host Tab for all entities
   public static final int EXPECTED_COLUMN_COUNT = 12;

   // Host Tab for all entities
   public static final int TASK_TAB_EXPECTED_COLUMN_COUNT = 9;
   public static final long LIMIT_UNLIMITED_LONG = -1;
   public static final long DEFAULT_RESERVATION_LONG = 0;
   public static boolean DEFAULT_EXPANDABLE_BOOL = true;
   public static String DEFAULT_EXPANDABLE = Boolean.toString(DEFAULT_EXPANDABLE_BOOL);
   public static final String DEFAULT_RESERVATION = Long
         .toString(DEFAULT_RESERVATION_LONG);
   public static final int DEFAULT_CPU_SHARES_MULTIPLIER = 1000;
   public static final int DEFAULT_MEM_SHARES_MULTIPLIER = 10;
   public static final int DEFAULT_CPU_SHARES_LOW_INT =
         2 * DEFAULT_CPU_SHARES_MULTIPLIER;
   public static final String DEFAULT_CPU_SHARES_LOW = Integer
         .toString(DEFAULT_CPU_SHARES_LOW_INT);
   public static final int DEFAULT_CPU_SHARES_NORMAL_INT =
         4 * DEFAULT_CPU_SHARES_MULTIPLIER;
   public static final String DEFAULT_CPU_SHARES_NORMAL = Integer
         .toString(DEFAULT_CPU_SHARES_NORMAL_INT);
   public static final int DEFAULT_CPU_SHARES_HIGH_INT =
         8 * DEFAULT_CPU_SHARES_MULTIPLIER;
   public static final String DEFAULT_CPU_SHARES_HIGH = Integer
         .toString(DEFAULT_CPU_SHARES_HIGH_INT);
   public static final int DEFAULT_MEMORY_SHARES_LOW_INT = 81920;
   public static final String DEFAULT_MEMORY_SHARES_LOW = Integer
         .toString(DEFAULT_MEMORY_SHARES_LOW_INT);
   public static final int DEFAULT_MEMORY_SHARES_NORMAL_INT = 163840;
   public static final String DEFAULT_MEMORY_SHARES_NORMAL = Integer
         .toString(DEFAULT_MEMORY_SHARES_NORMAL_INT);
   public static final int DEFAULT_MEMORY_SHARES_HIGH_INT = 327680;
   public static final String DEFAULT_MEMORY_SHARES_HIGH = Integer
         .toString(DEFAULT_MEMORY_SHARES_HIGH_INT);

   public static final String VMPROV_DEFAULT_GOS_VERSION = GOS_VERSION_WIN_2008_R2_64;
   public static final int VMPROV_DEFAULT_NUMBER_OF_VMS_TO_CREATE = 1;

   public static final String EXTENSION_VC_SERVER = "<VC_SERVER>";
   public static final String LIC_LUD_MISSING = "Licensing usage data for "
         + EXTENSION_VC_SERVER + " is missing for the selected time period.";

   public static final String VIRTUALSWITCH_NAME = "vSwitch_"
         + System.currentTimeMillis();

   public static final String FT_HOST_NOT_CONFIGURED_LEGACY_VIRTUAL_HOST =
         FT_HOST_NOT_CONFIGURED_NO_FT_LOGGING_NIC + "\r"
               + FT_HOST_NOT_CONFIGURED_NO_VMOTION_NIC + "\r"
               + FT_HOST_NOT_CONFIGURED_CPU_NOT_FT_ENABLED;

   public static final String INVALID_IP_ADDRESS_ERROR_STRING = "Invalid IP address.";

   /*
    * Kept these here instead of moving to properties file since annotations
    * expect direct constants and were throwing error to get these from parent
    * TestConstantsKey.java class
    */
   public static final String TEST_GROUP_PRECHECKIN = "precheckin";
   public static final String TEST_GROUP_PRECHECKIN_NEW = "precheckin_new";
   public static final String TEST_GROUP_VIRTUAL = "virtual";
   public static final String TEST_GROUP_PHYSICAL = "physical";
   public static final String TEST_GROUP_P0_VIRTUAL = "p0_virtual";
   public static final String TEST_GROUP_P0_PHYSICAL = "p0_physical";
   public static final String TEST_GROUP_LEGACY_HOST = "legacy_host";
   public static final String TEST_GROUP_ESX51_HOST = "esx51_host";

   // EVC CPU Specific
   public enum CLUSTER_CONFIGURATIONS_LIST_OPTIONS {
      SERVICES("0"), VSPHERE_DRS("1"), VSPHERE_HA("2"), CONFIGURATION("3"), GENERAL("4"), VMWARE_EVC(
            "5");
      private final String clsConfig;

      CLUSTER_CONFIGURATIONS_LIST_OPTIONS(String clsConfigVal) {
         clsConfig = clsConfigVal;
      }

      public String getClsuerEVCVal() {
         return clsConfig;
      }
   }

   // Host Connected USB
   public enum VM_MANAGE_SETTINGS_HARDWARE_LIST_OPTIONS {
      VM_HARDWARE("0"), VM_OPTIONS("1"), VM_SDRS_RULES("2"), VAPP_OPTIONS("3");
      private final String vmHardwareOption;

      VM_MANAGE_SETTINGS_HARDWARE_LIST_OPTIONS(String vmHardwareOpt) {
         vmHardwareOption = vmHardwareOpt;
      }

      public String getVMHardwareOption() {
         return vmHardwareOption;
      }
   }

   // command to unlink the vswitch
   public static final String CMD_VSWITCH_UNLINK = "esxcfg-vswitch -U ";
   public static final String CMD_VSWITCH_LINK = "esxcfg-vswitch -L ";

   /**
    * Use StringUtil.uniqueName()
    */
   @Deprecated
   public static String uniqueName() {
      return "NewObject" + System.currentTimeMillis();
   }

   /* End of Misc to be remediated section */
   /*
    * --------------------------------------------------------------------------
    * -
    */

   // Time Configuration
   public static enum SERVICE_ACTIVATION_POLICY {
      SAP_PORT_USAGE, SAP_HOST, SAP_MANUAL
   };

   public static enum SERVICE_TYPE_BUTTON {
      BUTTON_START, BUTTON_STOP, BUTTON_RESTART
   };

   // Client_Settings
   public static enum USER_MENU_REMOVE_OPTIONS {
      REMOVE_STORED_DATA, RESET_TO_FACTORY_DEFAULTS
   };

   public static enum REMOVE_STORED_DATA_OPTIONS {
      WORK_IN_PROGRESS, GETTING_STARTED_PREFERENCES, SAVED_SEARCHES
   };

   // Host states
   public enum HOST_POWER_MODE {
      MM_MODE, STANDBY_MODE, CONNECTED_MODE;
   }

   public static final String NEW_VAPP_WIZARD_CLASS = "NewVAppWizard";
   public static final String ADD_LICENSE_KEYS_WIZARD_CLASS = "AddLicenseKeysWizard";
   public static final String ASSIGN_LICENSE_KEY_DIALOG_CLASS = "SinglePageDialog";

   public static final String NETW_DEFAULT_IPV4_SUBNET_MASK = "255.255.255.0";
   public static final String NETW_DEFAULT_IPV6_SUBNET_MASK =
         "ffff:ffff:ffff:ffff:ffff:ffff:0:0";
   public static final String NETW_DEFAULT_IPV6_SUBNET_MASK_32 = "ffff:ffff:0:0:0:0:0:0";

   // Host roles
   public static final String PRIV_ADD_STANDALONE_HOST =
         "Host.Inventory.AddHostToCluster";
   public static final String PRIV_ADD_CLUSTERED_HOST =
         "Host.Inventory.AddStandaloneHost";
   public static final String PRIV_REMOVE_HOST = "Host.Inventory.RemoveHostFromCluster";
   public static final String PRIV_HOST_MM = "Host.Config.Maintenance";
   public static final String PRIV_GLOBAL_LOGEVENT = "Global.LogEvent";

   public static final String SELECT_COLUMNS = "2";

   public enum DayOfWeek {
      MONDAY("Monday", 0), TUESDAY("Tuesday", 1), WEDNESDAY("Wednesday", 2), THURSDAY(
            "Thusday", 3), FRIDAY("Friday", 4), SATURDAY("Saturday", 5), SUNDAY(
            "Sunday", 6);

      private final String label;
      private int index;

      private DayOfWeek(String label, int index) {
         this.label = label;
      }

      public String getLabel() {
         return label;
      }

      public int getIndex() {
         return index;
      }
   }

   public enum WeekOrdinal {
      FIRST("first", 0), SECOND("second", 1), THIRD("third", 2), FOURTH("fourth", 3), LAST(
            "last", 4);

      private final String label;
      private final int index;

      private WeekOrdinal(String label, int index) {
         this.label = label;
         this.index = index;
      }

      public String getLabel() {
         return label;
      }

      public int getIndex() {
         return index;
      }
   }

   public enum HOST_COMPATIBLITY_STATUSES {
      COMPOTABILE_STATUS(COMPATIBLE), INCOMPATIBLE_STATUS(INCOMPATIBLE);
      private final String value;

      private HOST_COMPATIBLITY_STATUSES(String value) {
         this.value = value;
      }

      public String getValue() {

         return value;
      }

   }

   public static final String[] TAG_CATEGORY_ASSOCIABLE_OBJECT_LIST = { "All objects",
         "Cluster", "Datacenter", "Datastore", "Distributed Port Group",
         "Distributed Switch", "Folder", "Host", "Network", "Resource Pool",
         "Datastore Cluster", "vApp", "Virtual Machine" };

   public enum FOLDER_TYPE {
      FOLDER_HOSTANDCLUSTERS, FOLDER_VMANDTEMPLATES, FOLDER_STORAGE, FOLDER_NETWORK, ROOT, FOLDER_VMANDTEMPLATES_SUBFOLDER
   }
}
