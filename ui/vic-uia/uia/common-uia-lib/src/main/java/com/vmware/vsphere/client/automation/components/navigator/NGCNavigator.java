/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator;

import com.vmware.client.automation.components.navigator.Navigator;
import com.vmware.client.automation.components.navigator.navstep.EntityViewNavStep;
import com.vmware.client.automation.components.navigator.navstep.HomeViewNavStep;
import com.vmware.client.automation.components.navigator.navstep.PrimaryTabNavStep;
import com.vmware.client.automation.components.navigator.navstep.SecondaryTabNavStep;
import com.vmware.client.automation.components.navigator.navstep.TocNavStep;
import com.vmware.client.automation.components.navigator.navstep.TreeTabNavStep;
import com.vmware.vsphere.client.automation.common.CommonUtil;

/**
 * A <code>Navigator</code> with VCD specific navigation steps registered on it.
 */
public class NGCNavigator extends Navigator {

   private static NGCNavigator s_vcdNavigator = null;

   // Navigation identifiers

   // First Level
   public final static String NID_HOME_CONSUMERS = "home.consumers";

   public final static String NID_HOME_RESOURCES = "home.resources";

   public final static String NID_HOME_VCENTER = "home.vcenter";

   public final static String NID_HOME_LIBRARIES = "home.contentLibraries";

   public final static String NID_HOME_RULES_AND_PROFILES = "home.rules";

   public final static String NID_HOME_ADMINISTRATION = "home.administration";

   public final static String NID_HOME_TAGS = "home.tags";

   public final static String NID_HOME_NEW_SEARCH = "home.newSearch";

   public final static String NID_VM_STORAGE_POLICIES = "home.vmStoragePolicies";

   public final static String NID_HOME_HOSTS_AND_CLUSTERS_TAB = "home.hostsAndClustersTab";

   public final static String NID_HOME_VMS_AND_TEMPLATES_TAB = "home.vmsAndTemplatesTab";

   public final static String NID_HOME_STORAGE_TAB = "home.storageTab";

   public final static String NID_HOME_NETWORKING_TAB = "home.networkingTab";

   // Second level

   // Consumers
   public final static String NID_CONSUMERS_ORGS = "entity.orgs";

   public final static String NID_CONSUMERS_CLOUD_VAPPS = "entity.cloudvapps";

   // Administration
   public final static String NID_ADMINISTRATION_LICENSES = "entity.administration.licenses";

   public final static String NID_ADMINISTRATION_ROLES = "administration.roles";

   public final static String NID_HOME_CEIP = "home.ceip";

   public final static String NID_ADMINISTRATION_CLIENT_PLUGINS="administration.client.plugins";

   public final static String NID_ADMINISTRATION_VC_SERVER_EXTENSIONS="administration.vc.server.extensions";

   // Resources
   public final static String NID_RESOURCES_CRPS = "entity.crps";

   // vCenter
   public final static String NID_VCENTER_VCS = "entity.vcs";

   // Datacenters
   public final static String NID_VCENTER_DCS = "entity.dcs";

   // Distributed Switches
   public final static String NID_VCENTER_DVSES = "entity.dvs";

   // Hosts
   public final static String NID_VCENTER_HOSTS = "entity.hosts";

   // Datastores
   public final static String NID_VCENTER_DATASTORES = "entity.datastores";

   // Datastore clusters
   public final static String NID_VCENTER_DATASTORE_CLUSTERS = "entity.datastoreClusters";

   // Hosts and Clusters
   public final static String NID_VCENTER_HOSTSANDCLUSTERS = "entity.HostsAndClusters";

   // VMs and Templates
   public final static String NID_VCENTER_VMSANDTEMPLATES = "entity.VmsAndTemplates";

   // Storage
   public final static String NID_VCENTER_STORAGE = "entity.storage";

   // Clusters
   public final static String NID_VCENTER_CLUSTERS = "entity.clusters";

   // Resource Pools
   public final static String NID_VCENTER_RPS = "entity.rps";

   // vDCs
   public final static String NID_VCENTER_VDC = "entity.vdc";

   // VMs
   public final static String NID_VCENTER_VM = "entity.vm";

   // VM templates in folders
   public final static String NID_VCENTER_TEMPLATES_IN_FOLDERS = "entity.template.in.folder";

   // vApps
   public final static String NID_VCENTER_VAPP = "entity.vapp";

   // Tags
   public final static String NID_TAGS_GETTING_STARTED = "tags.gettingStarted";

   public final static String NID_TAGS_ITEMS = "tags.items";

   public final static String NID_BACKING_TAGS_ITEMS = "backingTags.items";

   public final static String NID_BACKING_CATEGORIES_ITEMS = "backingCategories.items";

   // Custom Attributes
   public final static String NID_CUSTOM_ATTRIBUTES = "customAttributes.items";

   // Host Profiles
   public final static String NID_HOST_PROFILES = "entity.hostProfiles";
   public final static String NID_HOST_PROFILES_II_TAB_COMPLIANCE = "entity.hostProfiles.monitor.compliance";
   public final static String NID_HOST_PROFILES_II_TAB_SETTINGS = "entity.hostProfiles.manage.settings";

   // Entity level tabs
   public final static String NID_ENTITY_PRIMARY_TAB_GETTING_STARTED = "entity.l1tab.gettingStarted";

   public final static String NID_ENTITY_PRIMARY_TAB_SUMMARY = "entity.l1tab.summary";

   public final static String NID_ENTITY_PRIMARY_TAB_MONITOR = "entity.l1tab.monitor";

   public final static String NID_ENTITY_PRIMARY_TAB_MANAGE = "entity.l1tab.manage";

   public final static String NID_ENTITY_PRIMARY_TAB_RELATED_OBJECTS = "entity.l1tab.relatedObjects";

   public final static String NID_ENTITY_PRIMARY_TAB_RELATED_OBJECTS_CL = "entity.l1tab.relatedObjectsContentLibrary";

   public final static String NID_ENTITY_PRIMARY_TAB_OBJECTS = "entity.l1tab.itemsView";

   public final static String NID_ENTITY_PRIMARY_TAB_DATASTORES = "entity.l1tab.datastores";

   public final static String NID_ENTITY_PRIMARY_TAB_HOSTS = "entity.l1tab.hosts";

   public final static String NID_ENTITY_PRIMARY_TAB_NETWORKS = "entity.l1tab.networks";

   public final static String NID_ENTITY_PRIMARY_TAB_VMS = "entity.l1tab.vms";

   public final static String NID_ENTITY_PRIMARY_TAB_RPS = "entity.l1tab.rps";

   // Secondary tabs
   public final static String NID_ENTITY_MANAGE_TAB_II_TAB_CUSTOM_ATTRIBUTES = "host.l2tab.manage.customAttributes";

   // PBM secondary tabs - Monitor
   public final static String NID_PBM_TAB_MONITOR_VMS = "pbm.l2tab.monitor.vms";

   public final static String NID_PBM_TAB_MONITOR_TAGS = "pbm.l2tab.monitor.policyTags";

   // PBM secondary tabs - Related Objects
   public final static String NID_PBM_TAB_RELATED_OBJECTS_CLUSTERS = "pbm.l2tab.reatedObjects.clusters";

   public final static String NID_PBM_TAB_RELATED_OBJECTS_HOSTS = "pbm.l2tab.reatedObjects.hosts";

   public final static String NID_PBM_TAB_RELATED_OBJECTS_POLICY_TAGS = "pbm.l2tab.reatedObjects.policyTags";

   public final static String NID_PBM_TAB_RELATED_OBJECTS_VMS = "pbm.l2tab.reatedObjects.vms";

   // Administration
   public final static String NID_ADMINISTRATION_PRIMARY_TAB_LICENSES = "entity.administration.l1tab.licenses";

   public final static String NID_ADMINISTRATION_PRIMARY_TAB_ASSETS = "entity.administration.l1tab.assets";

   // Administration secondary navigation for "VCs"
   public final static String NID_ADMINISTRATION_SECONDARY_TAB_ASSETS_VCS = "entity.administration.l2tab.assets.vcs";

   // Administration secondary navigation for "Hosts"
   public final static String NID_ADMINISTRATION_SECONDARY_TAB_ASSETS_HOSTS = "entity.administration.l2tab.assets.hosts";

   // VM primary navigation for "Summary"
   public final static String NID_VM_SUMMARY = "vm.l1tab.summary";

   // VM secondary navigation for "Monitor"
   public final static String NID_VM_MONITOR_POLICIES = "vm.l2tab.monitor.policies";

   // VM Monitor TOC navigation for "Policies"
   public final static String NID_VM_MONITOR_POLICIES_TOC_STORAGE = "vm.monitor.policies.toc.storage";

   // VM secondary navigation for "Manage"
   public final static String NID_VM_MANAGE_POLICIES = "vm.l2tab.manage.policies";

   // VM Manage TOC navigation for "Policies"
   public final static String NID_VM_MANAGE_POLICIES_TOC_STORAGE = "vm.manage.policies.toc.storage";

   // VC TOC Navigation for "Settings"
   public final static String NID_VC_MANAGE_SETTINGS_TOC_VCHA = "vc.manage.settings.toc.vcha";

   // VC TOC Navigation for Monitor-VCHA
   public final static String NID_VC_MONITOR_VCHA_TOC_HEALTH = "vcenter.monitor.vcha.toc.health";

   // ORG Secondary navigation level
   public final static String NID_ORG_TAB_MANAGE_SETTINGS = "org.l2tab.manage.settings";

   // ORG TOC navigation
   public final static String NID_ORG_TOC_MANAGE_SETTINGS_GENERAL = "org.toc.manage.settings.general";

   public final static String NID_ORG_TOC_MANAGE_SETTINGS_EMAIL = "org.toc.manage.settings.email";

   public final static String NID_ORG_TOC_MANAGE_SETTINGS_POLICIES = "org.toc.manage.settings.policies";

   // CRP Secondary navigation level
   public final static String NID_CRP_TAB_MANAGE_SETTINGS = "crp.l2tab.manage.settings";

   public final static String NID_CRP_TAB_RELATED_OBJECTS_VDCS = "crp.l2tab.relatedObjects.vdcs";

   public final static String NID_CRP_TAB_RELATED_OBJECTS_STORAGE_PROFILES = "crp.l2tab.relatedObjects.storageProfiles";

   public final static String NID_CRP_TAB_RELATED_OBJECTS_CLUSTERS = "crp.l2tab.relatedObjects.clusters";

   public final static String NID_CRP_TAB_RELATED_OBJECTS_VMS = "crp.l2tab.relatedObjects.vms";

   public final static String NID_CRP_TAB_RELATED_OBJECTS_HOSTS = "crp.l2tab.relatedObjects.hosts";

   public final static String NID_CRP_TAB_RELATED_OBJECTS_DATASTORES = "crp.l2tab.relatedObjects.datastores";

   // Cloud vApp Secondary navigation level
   public static final String NID_CLOUDVAPP_TAB_MANAGE_SETTINGS = "cloudvapp.l2tab.manage.settings";

   public static final String NID_CLOUDVAPP_TAB_MANAGE_START_STOP_VMS = "cloudvapp.l2tab.manage.startstopvms";

   // Cloud vApp TOC navigation
   public static final String NID_CLOUDVAPP_TOC_MANAGE_SETTINGS_LEASES = "org.toc.manage.settings.leases";

   // Cluster primary tabs navigation
   public final static String NID_CLUSTER_PRIMARY_TAB_GETTING_STARTED = "cluster.l1tab.gettingStarted";

   public final static String NID_CLUSTER_PRIMARY_TAB_SUMMARY = "cluster.l1tab.summary";

   public final static String NID_CLUSTER_PRIMARY_TAB_MONITOR = "cluster.l1tab.monitor";

   public final static String NID_CLUSTER_PRIMARY_TAB_MANAGE = "cluster.l1tab.manage";

   public final static String NID_CLUSTER_PRIMARY_TAB_RELATED_OBJECTS = "cluster.l1tab.relatedObjects";

   // Datastore primary tabs navigation
   public final static String NID_DATASTORE_PRIMARY_TAB_SUMMARY = "datastore.l1tab.summary";

   public final static String NID_DATASTORE_PRIMARY_TAB_MANAGE = "datastore.l1tab.manage";

   public final static String NID_DATASTORE_PRIMARY_TAB_VMS = "datastore.l1tab.vms";

   // Datastore secondary tabs navigation
   public final static String NID_DATASTORE_MANAGE_II_TAB_TAGS = "datastore.l2tab.manage.tags";

   public final static String NID_DATASTORE_VMS_II_TAB_VMS = "datastore.l2tab.vms.vms";

   public final static String NID_DATASTORE_TOC_CONNECTIVITY_WITH_HOSTS = "datastore.toc.connectivityWithHosts";

   // Datastore Cluster primary tabs navigation
   public final static String NID_DS_CLUSTER_PRIMARY_TAB_GETTING_STARTED = "dsCluster.l1tab.gettingStarted";

   public final static String NID_DS_CLUSTER_PRIMARY_TAB_SUMMARY = "dsCluster.l1tab.summary";

   public final static String NID_DS_CLUSTER_PRIMARY_TAB_MONITOR = "dsCluster.l1tab.monitor";

   public final static String NID_DS_CLUSTER_PRIMARY_TAB_MANAGE = "dsCluster.l1tab.manage";

   public final static String NID_DS_CLUSTER_PRIMARY_TAB_RELATED_OBJECTS = "dsCluster.l1tab.relatedObjects";

   // Datastore Cluster secondary tabs navigation
   public final static String NID_DS_CLUSTER_II_TAB_SDRS = "dsCluster.l2tab.monitor.sdrs";
   public final static String NID_DS_CLUSTER_MANAGE_II_TAB_SETTINGS = "dsCluster.l2tab.manage.settings";

   // Cluster secondary tabs navigation
   public final static String NID_CLUSTER_MONITOR_II_TAB_ISSUES = "cluster.l2tab.monitor.issues";

   public final static String NID_CLUSTER_MONITOR_II_TAB_PERFORMANCE = "cluster.l2tab.monitor.performance";

   public final static String NID_TOC_MONITOR_PERFORMANCE_ADVANCED = "cluster.toc.monitor.performance.advanced";

   public final static String NID_TOC_MONITOR_PERFORMANCE_OVERVIEW = "cluster.toc.monitor.performance.overview";

   public final static String NID_CLUSTER_MONITOR_II_TAB_PROFILE_COMPLIANCE = "cluster.l2tab.monitor.profileCompliance";

   public final static String NID_CLUSTER_MONITOR_II_TAB_TASKS = "cluster.l2tab.monitor.tasks";

   public final static String NID_CLUSTER_MONITOR_II_TAB_EVENTS = "cluster.l2tab.monitor.events";

   public final static String NID_CLUSTER_MONITOR_II_TAB_RESOURCE_ALLOCATION = "cluster.l2tab.monitor.events";

   public final static String NID_CLUSTER_MONITOR_II_TAB_UTILIZATION = "cluster.l2tab.monitor.utilization";

   public final static String NID_CLUSTER_MONITOR_II_TAB_VSPHERE_DRS = "cluster.l2tab.monitor.vsphere.drs";

   public final static String NID_CLUSTER_MONITOR_II_TAB_VSPHERE_HA = "cluster.l2tab.monitor.vsphere.ha";

   public final static String NID_CLUSTER_MONITOR_II_TAB_STORAGE_REPORTS = "cluster.l2tab.monitor.storageReports";

   public final static String NID_CLUSTER_MANAGE_II_TAB_SETTINGS = "cluster.l2tab.manage.settings";

   public final static String NID_CLUSTER_MANAGE_II_TAB_ALARM_DEFINITIONS = "cluster.l2tab.manage.alarmDefinitions";

   public final static String NID_CLUSTER_MANAGE_II_TAB_TAGS = "cluster.l2tab.manage.tags";

   public final static String NID_CLUSTER_MANAGE_II_TAB_PERMISSIONS = "cluster.l2tab.manage.permissions";

   public final static String NID_CLUSTER_MANAGE_II_TAB_SCHEDULED_TASKS = "cluster.l2tab.manage.scheduledTasks";

   public final static String NID_CLUSTER_RELATED_OBJECTS_II_TAB_TOP_OBJECTS = "cluster.l2tab.relatedObjects.topObjets";

   public final static String NID_CLUSTER_RELATED_OBJECTS_II_TAB_HOSTS = "cluster.l2tab.relatedObjects.hosts";

   public final static String NID_CLUSTER_RELATED_OBJECTS_II_TAB_VMS = "cluster.l2tab.relatedObjects.vms";

   public final static String NID_DATASTORE_RELATED_OBJECTS_II_TAB_VMS = "datastore.l2tab.relatedObjects.vms";

   public final static String NID_DS_CLUSTER_RELATED_OBJECTS_II_TAB_DATASTORES = "dscluster.l2tab.relatedObjects.datastores";

   public final static String NID_VCENTER_RELATED_OBJECTS_II_TAB_VMS = "vcenter.l2tab.relatedObjects.vms";

   public final static String NID_VCENTER_RELATED_OBJECTS_II_TAB_VAPPS = "vcenter.l2tab.relatedObjects.vapps";

   public final static String NID_VCENTER_VMS_II_TAB_VMS = "vcenter.l2tab.vms.vms";

   public final static String NID_VCENTER_VMS_II_TAB_VAPPS = "vcenter.l2tab.vms.vapps";

   public final static String NID_VCENTER_MANAGE_II_TAB_SETTINGS = "vcenter.l2tab.manage.settings";

   public final static String NID_VCENTER_MONITOR_II_TAB_VCHA = "vcenter.l2tab.monitor.vcha";

   public final static String NID_VCENTER_MONITOR_II_TAB_ISSUES = "vcenter.l2tab.monitor.issues";

   public final static String NID_VCENTER_MONITOR_II_TAB_EVENTS = "vcenter.l2tab.monitor.events";

   public final static String NID_VCENTER_MONITOR_TOC_TRIGGERED_ALARMS = "vcenter.toc.monitor.issues.triggeredAlarms";

   public final static String NID_HOSTSANDCLUSTERS_RELATED_OBJECTS_II_TAB_VMS = "hostsandclusters.l2tab.relatedObjects.vms";

   public final static String NID_DATACENTER_RELATED_OBJECTS_II_TAB_VMS = "datacenter.l2tab.relatedObjects.vms";

   public final static String NID_DATACENTER_VMS_II_TAB_VMS = "datacenter.l2tab.vms.vms";

   public final static String NID_DATACENTER_VMS_II_TAB_VM_FOLDERS = "datacenter.l2tab.vms.vmfolders";

   public final static String NID_DATACENTER_HOSTS_II_TAB_HOSTS = "datacenter.l2tab.hosts.hosts";

   public final static String NID_CLUSTER_RELATED_OBJECTS_II_TAB_VAPPS = "cluster.l2tab.relatedObjects.vapps";

   public final static String NID_CLUSTER_RELATED_OBJECTS_II_TAB_DATASTORES = "cluster.l2tab.relatedObjects.datastores";

   public final static String NID_CLUSTER_RELATED_OBJECTS_II_TAB_DATASTORE_CLUSTERS = "cluster.l2tab.relatedObjects.datastoreClusters";

   public final static String NID_CLUSTER_RELATED_OBJECTS_II_TAB_NETWORKS = "cluster.l2tab.relatedObjects.networks";

   public final static String NID_CLUSTER_RELATED_OBJECTS_II_TAB_DVSES = "cluster.l2tab.relatedObjects.distributedSwitches";

   public final static String NID_CLUSTER_HOSTS_II_HOSTS = "cluster.l2tab.hosts.hosts";

   public final static String NID_CLUSTER_VMS_II_VAPPS = "cluster.l2tab.vms.vapps";

   public final static String NID_CLUSTER_VMS_II_VMS = "cluster.l2tab.vms.vms";

   // Cluster manage setting toc navigation
   public final static String NID_CLUSTER_TOC_MANAGE_SETTINGS_GENERAL = "cluster.toc.manage.settings.general";

   public final static String NID_CLUSTER_TOC_MANAGE_SETTINGS_VSPHERE_DRS = "cluster.toc.manage.settings.vsphereDrs";

   public final static String NID_CLUSTER_TOC_MANAGE_SETTINGS_VSAN_GENERAL = "cluster.toc.manage.settings.vsan.general";

   public final static String NID_CLUSTER_TOC_MANAGE_SETTINGS_VSAN_DISK_MANAGEMENT = "cluster.toc.manage.settings.vsan.diskManagement";

   public final static String NID_CLUSTER_TOC_MANAGE_SETTINGS_CONFIGURATION_LICENSING = "cluster.toc.manage.settings.configuration.licensing";

   // Host primary tabs navigation
   public final static String NID_HOST_PRIMARY_TAB_MANAGE = "host.l1tab.manage";
   public final static String NID_HOST_PRIMARY_TAB_RELATED_OBJECTS = "host.l1tab.relatedObjects";

   // Host secondary tabs navigation
   public final static String NID_HOST_MANAGE_TAB_SETTINGS = "host.l2tab.manage.settings";
   public final static String NID_HOST_MONITOR_TAB_PERFORMANCE = "host.l1tab.monitor.performance";
   public final static String NID_HOST_MANAGE_TAB_NETWORKING = "host.l1tab.manage.networking";

   public final static String NID_HOST_RELATED_OBJECTS_II_TAB_VMS = "host.l2tab.relatedObjects.vms";
   public final static String NID_HOST_RELATED_OBJECTS_II_TAB_VAPPS = "host.l2tab.relatedObjects.vapps";
   public final static String NID_HOST_RELATED_OBJECTS_II_TAB_DATASTORES = "host.l2tab.relatedObjects.datastores";

   public final static String NID_HOST_RELATED_OBJECTS_II_TAB_NETWORKS = "host.l2tab.relatedObjects.networks";
   public final static String NID_HOST_NETWORKS_II_TAB_NETWORKS = "host.l2tab.netowkrs.networks";

   public final static String NID_HOST_VMS_II_TAB_VMS = "host.l2tab.vms.vms";
   public final static String NID_HOST_VMS_II_TAB_VAPPS = "host.l2tab.vms.vapps";

   // Configure Other TOC
   public final static String NID_TOC_CONFIGURE_OTHER_TAGS = "toc.configure.other.tags";

   // Host settings tabs navigation
   public final static String NID_HOST_TOC_MANAGE_SETTINGS_LICENSING = "host.toc.manage.settings.licensing";
   public final static String NID_HOST_TOC_MANAGE_SETTINGS_ADVANCED_SYSTEM_SETTINGS = "host.toc.manage.settings.advanced.system.settings";
   public final static String NID_HOST_TOC_MANAGE_SETTINGS_GRAPHICS = "host.toc.manage.settings.graphics";

   // Host manage networking toc navigation
   public final static String NID_HOST_TOC_MANAGE_NETWORKING_VIRTUAL_SWITCHES = "host.toc.manage.networking.virtualSwitches";
   public final static String NID_HOST_TOC_MANAGE_NETWORKING_VMKERNEL_ADAPTERS = "host.toc.manage.networking.vmkernelAdapters";

   // Distributed Switch primary tabs navigation
   public final static String NID_DVS_PRIMARY_TAB_MANAGE = "dvs.l1tab.manage";
   public final static String NID_DVS_PRIMARY_TAB_RELATED_OBJECTS = "dvs.l1tab.relatedObjects";

   // Distributed Switch secondary tabs navigation
   public final static String NID_DVS_MANAGE_TAB_SETTINGS = "dvs.l2tab.manage.settings";
   public final static String NID_DVS_MANAGE_TAB_PORTS = "dvs.l2tab.manage.ports";
   public final static String NID_DVS_MANAGE_TAB_RESOURCE_ALLOCATION = "dvs.l2tab.manage.resourceAllocation";
   public final static String NID_DVS_RELATED_OBJECTS_II_TAB_DVPGS = "dvs.l2tab.relatedObjects.dvgps";
   public final static String NID_DVS_NETWORKS_II_TAB_DVPGS = "dvs.l2tab.networks.dvgps";

   // Distributed Switch manage settings TOC
   public final static String NID_DVS_MANAGE_TAB_SETTINGS_TOC_TOPOLOGY = "dvs.l2tab.manage.settings.topology";

   // vApp primary tabs navigation
   public final static String NID_VAPP_PRIMARY_TAB_RELATED_OBJECTS = "vapp.l1tab.relatedObjects";
   public final static String NID_VAPP_PRIMARY_TAB_SUMMARY = "vapp.l1tab.summary";
   public final static String NID_VAPP_PRIMARY_TAB_NETWORKS = "vapp.l1tab.networks";

   // vApp secondary tabs navigation
   public final static String NID_VAPP_RELATED_OBJECTS_II_TAB_VMS = "vapp.l2tab.relatedObjects.vms";
   public final static String NID_VAPP_RELATED_OBJECTS_II_TAB_VAPPS = "vapp.l2tab.relatedObjects.vapps";
   public final static String NID_VAPP_RELATED_OBJECTS_II_TAB_NWS = "vapp.l2tab.relatedObjects.networks";

   public final static String NID_VAPP_VMS_II_TAB_VMS = "vapp.l2tab.vms.vms";
   public final static String NID_VAPP_VMS_II_TAB_VAPPS = "vapp.l2tab.vms.vapps";

   // Resource pool primary tabs navigation
   public final static String NID_RESPOOL_PRIMARY_TAB_RELATED_OBJECTS = "respool.l1tab.relatedObjects";
   public final static String NID_RESPOOL_PRIMARY_TAB_SUMMARY = "respool.l1tab.summary";

   public final static String NID_RESPOOL_VMS_II_TAB_VMS = "respool.l2tab.vms.vms";
   public final static String NID_RESPOOL_VMS_II_TAB_VAPPS = "respool.l2tab.vms.vapps";

   // Resource pool secondary tabs navigation
   public final static String NID_RESPOOL_RELATED_OBJECTS_II_TAB_VMS = "respool.l2tab.relatedObjects.vms";
   public final static String NID_RESPOOL_RELATED_OBJECTS_II_TAB_VAPPS = "respool.l2tab.relatedObjects.vapps";

   // Content Libraries
   public final static String NID_CL_TAB_SUMMARY = "library.l1tab.summary";
   public final static String NID_CL_TAB_MANAGE = "library.l1tab.manage";
   public final static String NID_CL_TAB_RELATED_OBJECTS = "library.l1tab.relatedObjects";
   public final static String NID_CL_TAB_TEMPLATES = "library.l1tab.ovfTemplateForLibrary";
   public final static String NID_CL_TAB_OTHER_TYPES = "library.l1tab.otherTypes";

   // Library secondary tabs navigation
   public final static String NID_CL_RELATED_OBJECTS_TEMPLATES = "library.l2tab.ovfTemplateForLibrary";
   public final static String NID_CL_RELATED_OBJECTS_OTHER_TYPES = "library.l2tab.otherTypes";
   public final static String NID_CL_RELATED_OBJECTS_DATASTORES = "library.l2tab.datastores";

   // -------------------------------------------------------------------------
   // Datacenter primary tabs navigation
   public final static String NID_DATACENTER_PRIMARY_TAB_GETTING_STARTED = "datacenter.l1tab.gettingStarted";

   public final static String NID_DATACENTER_PRIMARY_TAB_SUMMARY = "datacenter.l1tab.summary";

   public final static String NID_DATACENTER_PRIMARY_TAB_MONITOR = "datacenter.l1tab.monitor";

   public final static String NID_DATACENTER_PRIMARY_TAB_MANAGE = "datacenter.l1tab.manage";

   public final static String NID_DATACENTER_PRIMARY_TAB_RELATED_OBJECTS = "datacenter.l1tab.relatedObjects";

   // -------------------------------------------------------------------------
   // Datacenter secondary tabs navigation
   public final static String NID_DATACENTER_MONITOR_TASKS = "datacenter.l2tab.monitor.tasks";
   public final static String NID_DATACENTER_MONITOR_EVENTS = "datacenter.l2tab.monitor.events";
   public final static String NID_DATACENTER_RELATED_OBJECTS_II_TAB_TOP_LEVEL_OBJECTS = "datacenter.l2tab.relatedObjects.topLevelObjects";
   public final static String NID_DATACENTER_RELATED_OBJECTS_II_TAB_CLUSTERS = "datacenter.l2tab.relatedObjects.clusters";

   public final static String NID_DATACENTER_RELATED_OBJECTS_II_TAB_DS_CLUSTERS = "datacenter.l2tab.relatedObjects.datastoreClusters";
   public final static String NID_DATACENTER_RELATED_OBJECTS_II_TAB_HOSTS = "datacenter.l2tab.relatedObjects.hosts";
   public final static String NID_DATACENTER_RELATED_OBJECTS_II_TAB_DVSES = "datacenter.l2tab.relatedObjects.dvsForDatacenter";
   public final static String NID_DATACENTER_RELATED_OBJECTS_II_TAB_DVPGS = "datacenter.l2tab.relatedObjects.dvPortgroups";
   public final static String NID_DATACENTER_RELATED_OBJECTS_II_TAB_VAPPS = "datacenter.l2tab.relatedObjects.vapps";

   public final static String NID_DATACENTER_NETWORKS_II_TAB_DVPGS = "datacenter.l2tab.networks.dvPortgroups";
   public final static String NID_DATACENTER_NETWORKS_II_TAB_DVSES = "datacenter.l2tab.networks.dvsForDatacenter";

   public static final String NID_LIBRARY_ITEM_MANAGE = "library.item.l1tab.manage";
   public static final String NID_LIBRARY_ITEM_MANAGE_POLICIES = "library.item.l2tab.manage.policies";

   // UIDs
   private static final String LINK_UID_VM_STORAGE_POLICIES = "TreeNodeItem_vsphere.core.navigator.pbmStorageProfiles";

   private static final String LINK_UID_CONSUMERS = "TreeNodeItem_vsphere.core.navigator.consumers";

   private static final String LINK_UID_RESOURCES = "TreeNodeItem_vsphere.core.navigator.resources";

   private static final String LINK_UID_VCENTERS = "TreeNodeItem_vsphere.core.navigator.virtualInfrastructure";

   private static final String LINK_UID_LIBRARIES = "TreeNodeItem_vsphere.core.navigator.contentLibraries";

   private static final String LINK_UID_RULES_AND_PROFILES = "TreeNodeItem_vsphere.core.navigator.rulesAndProfiles";

   private static final String LINK_UID_ADMINISTRATION = "TreeNodeItem_vsphere.core.navigator.administration";

   private static final String LINK_UID_TAGS = "TreeNodeItem_vsphere.core.navigator.tagsAndCustomAttributesManager";

   private static final String LINK_UID_NEW_SEARCH = "TreeNodeItem_vsphere.core.navigator.search";

   private static final String LINK_UID_HOSTS_AND_CLUSTERS_TAB = "TreeNodeItem_vsphere.core.navigator.hostsClustersTree";

   private static final String LINK_UID_VMS_AND_TEMPLATES_TAB = "TreeNodeItem_vsphere.core.navigator.vmTemplatesTree";

   private static final String LINK_UID_STORAGE_TAB = "TreeNodeItem_vsphere.core.navigator.storageTree";

   private static final String LINK_UID_NETWORKING_TAB = "TreeNodeItem_vsphere.core.navigator.networkingTree";

   private static final String LINK_UID_ADMINISTRATION_LICENSES = "TreeNodeItem_vsphere.license.navigator.management";

   private static final String LINK_UID_ADMINISTRATION_ROLES = "TreeNodeItem_vsphere.core.navigator.accessRoleManager";

   private static final String LINK_UID_ORGANIZATIONS = "TreeNodeItem_vsphere.core.navigator.vilist.vsphere.core.vCloudOrganizations_AdminOrgType_";

   private static final String LINK_UID_VCENTER_VCS = "TreeNodeItem_vsphere.core.viVCenterServers.list.Folder_isRootFolder";

   private static final String LINK_UID_VCENTER_DCS = "TreeNodeItem_vsphere.core.viDatacenters.list.Datacenter_";

   private static final String LINK_UID_VCENTER_HOSTS = "TreeNodeItem_vsphere.core.viHosts.list.HostSystem_";

   private static final String LINK_UID_VCENTER_DATASTORES = "TreeNodeItem_vsphere.core.viDatastores.list.Datastore_!isSystemDatastore";

   private static final String LINK_UID_VCENTER_DS_CLUSTERS = "TreeNodeItem_vsphere.core.viDsCluster.list.StoragePod_";

   private static final String LINK_UID_VCENTER_DVSES = "TreeNodeItem_vsphere.core.viDvs.list.DistributedVirtualSwitch_";

   private static final String LINK_UID_VCENTER_HOSTSANDCLUSTERS = "TreeNodeItem_vsphere.core.navigator.hostsClustersTree";

   private static final String LINK_UID_VCENTER_VMSANDTEMPLATES = "TreeNodeItem_vsphere.core.navigator.vmTemplatesTree";

   private static final String LINK_UID_VCENTER_STORAGE = "TreeNodeItem_vsphere.core.navigator.storageTree";

   private static final String LINK_UID_VCENTER_CLUSTERS = "TreeNodeItem_vsphere.core.viClusters.list.ClusterComputeResource_";

   private static final String LINK_UID_VCENTER_RPS = "TreeNodeItem_vsphere.core.viResourcePools.list.ResourcePool_isNonRootRP";

   private static final String LINK_UID_VCENTER_VDC = "TreeNodeItem_vsphere.core.vdcs.list.com.vmware.vcenter.vdcs.Vdc_";

   private static final String LINK_UID_VCENTER_VM = "TreeNodeItem_vsphere.core.viVms.list.VirtualMachine_isNormalVM";

   private static final String LINK_VCENTER_TEMPLATES_IN_FOLDERS = "TreeNodeItem_vsphere.core.viVmTemplates.list.VirtualMachine_config.template";

   private static final String LINK_UID_VCENTER_VAPP = "TreeNodeItem_vsphere.core.viVapps.list.VirtualApp_";

   private static final String LINK_UID_CLOUD_VAPPS = "TreeNodeItem_vsphere.core.navigator.vilist.vcloud.core.vCloudVapps_VAppType_";

   private static final String LINK_UID_CLOUD_RESOURCE_POOLS = "TreeNodeItem_vsphere.core.navigator.vilist.vsphere.core.vCloudCrps_ProviderVdcType_";

   private static final String LINK_UID_HOST_PROFILES = "TreeNodeItem_vsphere.core.navigator.hostProfiles";

   private static final String LINK_UID_CEIP_ADMINISTRATION_LINK = "TreeNodeItem_com.vmware.vsphere.client.phonehomecollector.navigation.application";

   private static final String LINK_UID_ADMINISTRATION_CLIENT_PLUGINS = "TreeNodeItem_vsphere.core.navigator.solutionPluginManager";

   private static final String LINK_UID_ADMINISTRATION_VC_SERVER_EXTENSION = "TreeNodeItem_vsphere.core.navigator.solutionVcExtensionsManager";

   // Content library Manage second navigation tabs
   public final static String NID_LIBRARY_MANAGE_II_TAB_SETTINGS = "library.l2tab.manage.settings";
   public final static String NID_LIBRARY_MANAGE_II_TAB_TAGS = "library.l2tab.manage.tags";
   public final static String NID_LIBRARY_MANAGE_II_TAB_PERMISSIONS = "library.l2tab.manage.permission";
   public final static String NID_LIBRARY_MANAGE_II_TAB_STORAGE = "library.l2tab.manage.storage";

   // Storage policy Monitor tabs
   public final static String NID_STORAGE_POLICY_MONITOR_STORAGE_COMPATIBILITY = "storagePolicy.monitor.storageCompatibility";
   public final static String NID_STORAGE_POLICY_MONITOR_VMS_AND_DISKS = "storagePolicy.monitor.vmsAndDisks";

   // Storage Policy Component
   public static final String NID_STORAGE_POLICY_STORAGE_POLICY_COMPONENTS = "storagePolicy.storagePolicyComponents";

   // The following indices are used for tab navigation (they are not path
   // steps!)
   private static final int IDX_PRIMARY_TAB_GETTING_STARTED = 0;

   private static final int IDX_PRIMARY_TAB_OBJECTS = 1;


   private static final int IDX_PRIMARY_TAB_SUMMARY = 1;

   private static final int IDX_PRIMARY_TAB_MANAGE = 3;

   private static final int IDX_PRIMARY_TAB_RELATED_OBJECTS = 4;

   // PBM second navigation indices
   private static final int IDX_PBM_TAB_MONITOR_VMS = 0;
   private static final int IDX_PBM_TAB_MONITOR_TAGS = 1;

   private static final int IDX_PBM_TAB_RELATED_OBJECTS_CLUSTERS = 1;
   private static final int IDX_PBM_TAB_RELATED_OBJECTS_HOSTS = 0;
   private static final int IDX_PBM_TAB_RELATED_OBJECTS_POLICY_TAGS = 3;
   private static final int IDX_PBM_TAB_RELATED_OBJECTS_VMS = 2;

   // Administration first level navigation indices
   private static final int IDX_ADMINISTRATION_PRIMARY_TAB_LICENSES = 1;
   private static final int IDX_ADMINISTRATION_PRIMARY_TAB_ASSETS = 3;

   // Administration second level navigation indices
   private static final int IDX_ADMINISTRATION_SECONDARY_TAB_ASSETS_VCS = 0;
   private static final int IDX_ADMINISTRATION_SECONDARY_TAB_ASSETS_HOSTS = 1;

   // VM first level navigation indices
   private static final int IDX_VM_SUMMARY = 1;

   // VM second navigation indices
   private static final int IDX_VM_MONITOR_POLICIES = 2;
   private static final int IDX_VM_MANAGE_POLICIES = 4;

   // VM TOC navigation indices
   private static final int IDX_VM_MONITOR_POLICIES_TOC_STORAGE = 0;

   private static final int IDX_VM_MANAGE_POLICIES_TOC_STORAGE = 0;

   // Organization navigation indices
   private static final int IDX_ORG_TAB_MANAGE_SETTINGS = 0;

   private static final int IDX_ORG_TOC_MANAGE_SETTINGS_GENERAL = 0;

   private static final int IDX_ORG_TOC_MANAGE_SETTINGS_EMAIL = 1;

   private static final int IDX_ORG_TOC_MANAGE_SETTINGS_POLICIES = 2;

   // CRP navigation indices
   private static final int IDX_CRP_TAB_MANAGE_SETTINGS = 0;

   private static final int IDX_CRP_TAB_RELATED_OBJECTS_VDCS = 0;

   private static final int IDX_CRP_TAB_RELATED_OBJECTS_STORAGE_PROFILES = 1;

   private static final int IDX_CRP_TAB_RELATED_OBJECTS_CLUSTERS = 2;

   private static final int IDX_CRP_TAB_RELATED_OBJECTS_VMS = 3;

   private static final int IDX_CRP_TAB_RELATED_OBJECTS_HOSTS = 4;

   private static final int IDX_CRP_TAB_RELATED_OBJECTS_DATASTORES = 5;

   // Cloud vApp navigation indices
   private static final int IDX_CLOUDVAPP_TAB_MANAGE_SETTINGS = 0;

   private static final int IDX_CLOUDVAPP_TAB_MANAGE_START_STOP_VMS = 1;

   private static final int IDX_CLOUDVAPP_TOC_MANAGE_SETTINGS_LEASES = 0;

   // Datastore Cluster second navigation indices
   private static final int IDX_DS_CLUSTER_TAB_MONITOR_SDRS = 4;
   private static final int IDX_DS_CLUSTER_TAB_RELATED_OBJECTS_DATASTORES = 0;

   // Cluster second navigation indices
   private static final int IDX_CLUSTER_TAB_MONITOR_ISSUES = 0;
   private static final int IDX_CLUSTER_TAB_MONITOR_PROFILE_COMPLIANCE = 2;
   private static final int IDX_CLUSTER_TAB_MONITOR_TASKS = 3;
   private static final int IDX_CLUSTER_TAB_MONITOR_EVENTS = 4;
   private static final int IDX_CLUSTER_TAB_MONITOR_RESOURCE_ALLOCATION = 5;
   private static final int IDX_CLUSTER_TAB_MONITOR_UTILIZATION = 6;
   private static final int IDX_CLUSTER_TAB_MONITOR_STORAGE_REPORTS = 7;

   private static final int IDX_CLUSTER_TAB_MANAGE_SETTINGS = 0;
   private static final int IDX_CLUSTER_TAB_MANAGE_ALARM_DEFINITIONS = 1;
   private static final int IDX_CLUSTER_TAB_MANAGE_PERMISSIONS = 3;
   private static final int IDX_CLUSTER_TAB_MANAGE_SCHEDULED_TASKS = 4;

   private static final int IDX_CLUSTER_TAB_REL_OBJ_TOP_OBJECTS = 0;
   private static final int IDX_CLUSTER_TAB_REL_OBJ_HOSTS = 1;
   private static final int IDX_CLUSTER_TAB_REL_OBJ_VMS = 2;
   private static final int IDX_DATASTORE_TAB_REL_OBJ_VMS = 0;
   private static final int IDX_VCENTER_TAB_REL_OBJ_VMS = 4;
   private static final int IDX_VCENTER_TAB_REL_OBJ_VAPPS = 6;
   private static final int IDX_DATACENTER_TAB_REL_OBJ_VMS = 3;
   private static final int IDX_CLUSTER_TAB_REL_OBJ_VAPPS = 3;
   private static final int IDX_CLUSTER_TAB_REL_OBJ_DATASTORES = 4;
   private static final int IDX_CLUSTER_TAB_REL_OBJ_DATASTORE_CLUSTERS = 5;
   private static final int IDX_CLUSTER_TAB_REL_OBJ_NETWORKS = 6;
   private static final int IDX_CLUSTER_TAB_REL_OBJ_DVSES = 7;

   private static final int IDX_CLUSTER_TOC_MANAGE_SETTINGS_VSPHERE_DRS = 1;
   private static final int IDX_CLUSTER_TOC_MANAGE_SETTINGS_VSAN_GENERAL = 4;
   private static final int IDX_CLUSTER_TOC_MANAGE_SETTINGS_VSAN_DISK_MANAGEMENT = 5;
   private static final int IDX_CLUSTER_TOC_MANAGE_SETTINGS_GENERAL = 8;
   private static final int IDX_CLUSTER_TOC_MANAGE_SETTINGS_CONFIGURATION_LICENSING = 10;
   private static final int IDX_TOC_MONITOR_PERFORMANCE_ADVANCED = 1;
   private static final int IDX_TOC_MONITOR_PERFORMANCE_OVERVIEW = 0;

   // Host second navigation indices
   private static final int IDX_HOST_TAB_MANAGE_SETTINGS = 0;

   private static final int IDX_HOST_TAB_REL_OBJ_VMS = 0;
   private static final int IDX_HOST_TAB_REL_OBJ_NETWORKS = 2;
   private static final int IDX_HOST_TAB_REL_OBJ_VAPPS = 3;
   private static final int IDX_HOST_TAB_REL_OBJ_DATASTORES = 4;

   private static final int IDX_HOST_TOC_MANAGE_SETTINGS_LICENSING = 6;
   private static final int IDX_HOST_TOC_MANAGE_SETTINGS_ADVANCED_SYSTEM_SETTINGS = 12;
   private static final int IDX_HOST_TOC_MANAGE_SETTINGS_GRAPHICS = 19;

   // Dvs second navigation indices
   private static final int IDX_DVS_TAB_MANAGE_SETTINGS = 0;
   private static final int IDX_DVS_TAB_MANAGE_PORTS = 5;
   private static final int IDX_DVS_TAB_MANAGE_RESOURCE_ALLOCATION = 6;
   private static final int IDX_DVS_TAB_REL_OBJ_DVPGS = 3;
   private static final int IDX_DVS_TAB_MANAGE_SETTINGS_TOC_TOPOLOGY = 1;

   // vApp second navigation indices
   private static final int IDX_VAPP_TAB_REL_OBJ_VMS = 1;
   private static final int IDX_VAPP_TAB_REL_OBJ_VAPPS = 2;
   private static final int IDX_VAPP_TAB_REL_OBJ_NWS = 4;
   // Resource pool second navigation indices
   private static final int IDX_RESPOOL_TAB_REL_OBJ_VMS = 1;
   private static final int IDX_RESPOOL_TAB_REL_OBJ_VAPPS = 2;
   // Libraries primary navigation indices
   private static final int IDX_CL_TAB_SUMMARY = 1;
   private static final int IDX_CL_TAB_MANAGE = 2;
   private static final int IDX_CL_TAB_RELATED_OBJECTS = 3;

   // Library second navigation indices
   private static final int IDX_CL_RELATED_OBJECTS_TEMPLATES = 0;
   private static final int IDX_CL_RELATED_OBJECTS_OTHER_TYPES = 1;
   private static final int IDX_CL_RELATED_OBJECTS_DATASTORES = 2;
   // -------------------------------------------------------------------------
   // Datacenter second navigation indices
   private static final int IDX_DATACENTER_TAB_MONITOR_TASKS = 2;
   private static final int IDX_DATACENTER_TAB_MONITOR_EVENTS = 3;
   private static final int IDX_DATACENTER_TAB_REL_OBJ_TOP_LEVEL_OBJECTS = 0;
   private static final int IDX_DATACENTER_TAB_REL_OBJ_CLUSTERS = 1;
   private static final int IDX_DATACENTER_TAB_REL_OBJ_HOSTS = 2;
   private static final int IDX_DATACENTER_TAB_REL_OBJ_DVSES = 9;
   private static final int IDX_DATACENTER_TAB_REL_OBJ_DS_CLUSTERS = 7;
   private static final int IDX_DATACENTER_TAB_REL_OBJ_DVPGS = 10;
   private static final int IDX_DATACENTER_TAB_REL_OBJ_VAPPS = 5;

   // -------------------------------------------------------------------------------------
   // Library Item Tab and Toc Indexes
   private static final int IDX_LIBRARY_ITEM_MANAGE = 1;
   private static final int IDX_LIBRARY_ITEM_MANAGE_POLICIES = 0;

   // -------------------------------------------------------------------------------------
   // Library manage tab Indexes
   private static final int IDX_LIBRARY_MANAGE_SETTINGS = 0;
   private static final int IDX_LIBRARY_MANAGE_TAGS = 1;
   private static final int IDX_LIBRARY_MANAGE_PERMISSIONS = 2;
   private static final int IDX_LIBRARY_MANAGE_STORAGE = 3;

   // -------------------------------------------------------------------------------------
   // Storage policy monitor tab indices
   private static final int IDX_STORAGE_POLICY_MONITOR_STORAGE_COMPATIBILITY = 1;

   /**
    * Return an instance of the <code>VcdNavigator</code>. Currently a single
    * instance is always provided.
    *
    * @return <code>VcdNavigator</code> instance.
    */
   public static NGCNavigator getInstance() {
      if (s_vcdNavigator == null) {
         s_vcdNavigator = new NGCNavigator();
      }

      return s_vcdNavigator;
   }

   /**
    * Register the VCD specific navigation steps.
    */
   public NGCNavigator() {
      super();

      // Register home view steps
      registerNavStep(new HomeViewNavStep(NID_HOME_CONSUMERS,
            LINK_UID_CONSUMERS));

      registerNavStep(new HomeViewNavStep(NID_HOME_RESOURCES,
            LINK_UID_RESOURCES));

      registerNavStep(new HomeViewNavStep(NID_HOME_VCENTER, LINK_UID_VCENTERS));

      registerNavStep(new HomeViewNavStep(NID_HOME_LIBRARIES, LINK_UID_LIBRARIES));

      registerNavStep(new HomeViewNavStep(NID_HOME_RULES_AND_PROFILES,
            LINK_UID_RULES_AND_PROFILES));

      registerNavStep(new HomeViewNavStep(NID_VM_STORAGE_POLICIES,
            LINK_UID_VM_STORAGE_POLICIES));

      registerNavStep(new HomeViewNavStep(NID_HOME_ADMINISTRATION,
            LINK_UID_ADMINISTRATION));

      registerNavStep(new HomeViewNavStep(NID_HOME_TAGS, LINK_UID_TAGS));

      registerNavStep(new HomeViewNavStep(NID_HOME_CEIP,
            LINK_UID_CEIP_ADMINISTRATION_LINK));

      registerNavStep(new HomeViewNavStep(NID_HOME_NEW_SEARCH,
            LINK_UID_NEW_SEARCH));

      registerNavStep(new TreeTabNavStep(NID_HOME_HOSTS_AND_CLUSTERS_TAB,
            LINK_UID_HOSTS_AND_CLUSTERS_TAB));

      registerNavStep(new TreeTabNavStep(NID_HOME_VMS_AND_TEMPLATES_TAB,
            LINK_UID_VMS_AND_TEMPLATES_TAB));

      registerNavStep(new TreeTabNavStep(NID_HOME_STORAGE_TAB,
            LINK_UID_STORAGE_TAB));

      registerNavStep(new TreeTabNavStep(NID_HOME_NETWORKING_TAB,
            LINK_UID_NETWORKING_TAB));

      // Register entity view steps
      registerNavStep(new EntityViewNavStep(NID_CONSUMERS_ORGS,
            LINK_UID_ORGANIZATIONS));

      registerNavStep(new EntityViewNavStep(NID_RESOURCES_CRPS,
            LINK_UID_CLOUD_RESOURCE_POOLS));

      registerNavStep(new EntityViewNavStep(NID_CONSUMERS_CLOUD_VAPPS,
            LINK_UID_CLOUD_VAPPS));

      registerNavStep(new EntityViewNavStep(NID_HOST_PROFILES,
            LINK_UID_HOST_PROFILES));

      // -------------------------------------------------------------------------

      registerNavStep(new EntityViewNavStep(NID_ADMINISTRATION_LICENSES,
            LINK_UID_ADMINISTRATION_LICENSES));

      registerNavStep(new EntityViewNavStep(NID_ADMINISTRATION_ROLES,
            LINK_UID_ADMINISTRATION_ROLES));

      registerNavStep(new EntityViewNavStep(NID_ADMINISTRATION_CLIENT_PLUGINS,
            LINK_UID_ADMINISTRATION_CLIENT_PLUGINS));

      registerNavStep(new EntityViewNavStep(NID_ADMINISTRATION_VC_SERVER_EXTENSIONS,
            LINK_UID_ADMINISTRATION_VC_SERVER_EXTENSION));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_VCS,
            LINK_UID_VCENTER_VCS));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_DCS,
            LINK_UID_VCENTER_DCS));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_HOSTS,
            LINK_UID_VCENTER_HOSTS));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_DATASTORES,
            LINK_UID_VCENTER_DATASTORES));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_DVSES,
            LINK_UID_VCENTER_DVSES));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_DATASTORE_CLUSTERS,
            LINK_UID_VCENTER_DS_CLUSTERS));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_HOSTSANDCLUSTERS,
            LINK_UID_VCENTER_HOSTSANDCLUSTERS));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_VMSANDTEMPLATES,
            LINK_UID_VCENTER_VMSANDTEMPLATES));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_STORAGE,
            LINK_UID_VCENTER_STORAGE));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_CLUSTERS,
            LINK_UID_VCENTER_CLUSTERS));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_RPS,
            LINK_UID_VCENTER_RPS));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_VDC,
            LINK_UID_VCENTER_VDC));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_VM, LINK_UID_VCENTER_VM));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_TEMPLATES_IN_FOLDERS,
            LINK_VCENTER_TEMPLATES_IN_FOLDERS));

      registerNavStep(new EntityViewNavStep(NID_VCENTER_VAPP,
            LINK_UID_VCENTER_VAPP));

      // Primary tab steps
      registerNavStep(new PrimaryTabNavStep(NID_TAGS_GETTING_STARTED,
            CommonUtil.getLocalizedString("entity.tab.primary.gettingStarted")));

      registerNavStep(new PrimaryTabNavStep(NID_TAGS_ITEMS,
            CommonUtil.getLocalizedString("tags.items.tab")));

      registerNavStep(new SecondaryTabNavStep(NID_BACKING_TAGS_ITEMS,
            CommonUtil.getLocalizedString("tags.items.tags.tab")));

      registerNavStep(new SecondaryTabNavStep(NID_BACKING_CATEGORIES_ITEMS,
            CommonUtil.getLocalizedString("tags.items.categories.tab")));

      registerNavStep(new PrimaryTabNavStep(NID_CUSTOM_ATTRIBUTES,
            CommonUtil.getLocalizedString("ca.tab.name")));

      registerNavStep(new PrimaryTabNavStep(
            NID_ENTITY_PRIMARY_TAB_GETTING_STARTED,
            CommonUtil.getLocalizedString("entity.tab.primary.gettingStarted")));

      registerNavStep(new PrimaryTabNavStep(NID_ENTITY_PRIMARY_TAB_SUMMARY,
            CommonUtil.getLocalizedString("entity.tab.primary.summary")));

      registerNavStep(new PrimaryTabNavStep(NID_ENTITY_PRIMARY_TAB_MONITOR,
            CommonUtil.getLocalizedString("common.tabs.name.monitor")));

      registerNavStep(new PrimaryTabNavStep(NID_ENTITY_PRIMARY_TAB_MANAGE,
            CommonUtil.getLocalizedString("entity.tab.primary.manage")));

      registerNavStep(new PrimaryTabNavStep(NID_ENTITY_PRIMARY_TAB_HOSTS,
            CommonUtil.getLocalizedString("entity.tab.primary.hosts")));

      registerNavStep(new PrimaryTabNavStep(NID_ENTITY_PRIMARY_TAB_NETWORKS,
            CommonUtil.getLocalizedString("entity.tab.primary.networks")));

      registerNavStep(new PrimaryTabNavStep(NID_ENTITY_PRIMARY_TAB_VMS,
            CommonUtil.getLocalizedString("entity.tab.primary.vms")));

      registerNavStep(new PrimaryTabNavStep(NID_ENTITY_PRIMARY_TAB_RPS,
            CommonUtil.getLocalizedString("entity.tab.primary.rps")));

      registerNavStep(new PrimaryTabNavStep(
            NID_ENTITY_PRIMARY_TAB_RELATED_OBJECTS,
            IDX_PRIMARY_TAB_RELATED_OBJECTS));

      registerNavStep(new PrimaryTabNavStep(
            NID_ENTITY_PRIMARY_TAB_RELATED_OBJECTS_CL,
            IDX_CL_TAB_RELATED_OBJECTS));

      registerNavStep(new PrimaryTabNavStep(NID_ENTITY_PRIMARY_TAB_OBJECTS,
            IDX_PRIMARY_TAB_OBJECTS));

      registerNavStep(new PrimaryTabNavStep(NID_ENTITY_PRIMARY_TAB_DATASTORES,
            CommonUtil.getLocalizedString("entity.tab.primary.datastores")));

      // Register PBM steps
      registerNavStep(new SecondaryTabNavStep(NID_PBM_TAB_MONITOR_VMS,
            IDX_PBM_TAB_MONITOR_VMS));

      registerNavStep(new SecondaryTabNavStep(NID_PBM_TAB_MONITOR_TAGS,
            IDX_PBM_TAB_MONITOR_TAGS));

      registerNavStep(new SecondaryTabNavStep(
            NID_PBM_TAB_RELATED_OBJECTS_CLUSTERS,
            IDX_PBM_TAB_RELATED_OBJECTS_CLUSTERS));

      registerNavStep(new SecondaryTabNavStep(
            NID_PBM_TAB_RELATED_OBJECTS_HOSTS,
            IDX_PBM_TAB_RELATED_OBJECTS_HOSTS));

      registerNavStep(new SecondaryTabNavStep(
            NID_PBM_TAB_RELATED_OBJECTS_POLICY_TAGS,
            IDX_PBM_TAB_RELATED_OBJECTS_POLICY_TAGS));

      registerNavStep(new SecondaryTabNavStep(NID_PBM_TAB_RELATED_OBJECTS_VMS,
            IDX_PBM_TAB_RELATED_OBJECTS_VMS));

      // Administration
      registerNavStep(new PrimaryTabNavStep(
            NID_ADMINISTRATION_PRIMARY_TAB_LICENSES,
            IDX_ADMINISTRATION_PRIMARY_TAB_LICENSES));

      registerNavStep(new PrimaryTabNavStep(
            NID_ADMINISTRATION_PRIMARY_TAB_ASSETS,
            IDX_ADMINISTRATION_PRIMARY_TAB_ASSETS));

      registerNavStep(new SecondaryTabNavStep(
            NID_ADMINISTRATION_SECONDARY_TAB_ASSETS_VCS,
            IDX_ADMINISTRATION_SECONDARY_TAB_ASSETS_VCS));

      registerNavStep(new SecondaryTabNavStep(
            NID_ADMINISTRATION_SECONDARY_TAB_ASSETS_HOSTS,
            IDX_ADMINISTRATION_SECONDARY_TAB_ASSETS_HOSTS));

      // Register VM steps
      registerNavStep(new PrimaryTabNavStep(NID_VM_SUMMARY, IDX_VM_SUMMARY));

      registerNavStep(new SecondaryTabNavStep(NID_VM_MONITOR_POLICIES,
            IDX_VM_MONITOR_POLICIES));

      registerNavStep(new TocNavStep(NID_VM_MONITOR_POLICIES_TOC_STORAGE,
            IDX_VM_MONITOR_POLICIES_TOC_STORAGE));

      registerNavStep(new SecondaryTabNavStep(NID_VM_MANAGE_POLICIES,
            IDX_VM_MANAGE_POLICIES));

      registerNavStep(new TocNavStep(NID_VM_MANAGE_POLICIES_TOC_STORAGE,
            IDX_VM_MANAGE_POLICIES_TOC_STORAGE));

      // Register ORG steps
      registerNavStep(new SecondaryTabNavStep(NID_ORG_TAB_MANAGE_SETTINGS,
            IDX_ORG_TAB_MANAGE_SETTINGS));

      registerNavStep(new TocNavStep(NID_ORG_TOC_MANAGE_SETTINGS_GENERAL,
            IDX_ORG_TOC_MANAGE_SETTINGS_GENERAL));

      registerNavStep(new TocNavStep(NID_ORG_TOC_MANAGE_SETTINGS_EMAIL,
            IDX_ORG_TOC_MANAGE_SETTINGS_EMAIL));

      registerNavStep(new TocNavStep(NID_ORG_TOC_MANAGE_SETTINGS_POLICIES,
            IDX_ORG_TOC_MANAGE_SETTINGS_POLICIES));

      // Register CRP steps
      registerNavStep(new SecondaryTabNavStep(NID_CRP_TAB_MANAGE_SETTINGS,
            IDX_CRP_TAB_MANAGE_SETTINGS));

      registerNavStep(new SecondaryTabNavStep(NID_CRP_TAB_RELATED_OBJECTS_VDCS,
            IDX_CRP_TAB_RELATED_OBJECTS_VDCS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CRP_TAB_RELATED_OBJECTS_STORAGE_PROFILES,
            IDX_CRP_TAB_RELATED_OBJECTS_STORAGE_PROFILES));

      registerNavStep(new SecondaryTabNavStep(
            NID_CRP_TAB_RELATED_OBJECTS_CLUSTERS,
            IDX_CRP_TAB_RELATED_OBJECTS_CLUSTERS));

      registerNavStep(new SecondaryTabNavStep(NID_CRP_TAB_RELATED_OBJECTS_VMS,
            IDX_CRP_TAB_RELATED_OBJECTS_VMS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CRP_TAB_RELATED_OBJECTS_HOSTS,
            IDX_CRP_TAB_RELATED_OBJECTS_HOSTS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CRP_TAB_RELATED_OBJECTS_DATASTORES,
            IDX_CRP_TAB_RELATED_OBJECTS_DATASTORES));

      // Register cloud vApp steps
      registerNavStep(new SecondaryTabNavStep(
            NID_CLOUDVAPP_TAB_MANAGE_SETTINGS,
            IDX_CLOUDVAPP_TAB_MANAGE_SETTINGS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLOUDVAPP_TAB_MANAGE_START_STOP_VMS,
            IDX_CLOUDVAPP_TAB_MANAGE_START_STOP_VMS));

      registerNavStep(new TocNavStep(NID_CLOUDVAPP_TOC_MANAGE_SETTINGS_LEASES,
            IDX_CLOUDVAPP_TOC_MANAGE_SETTINGS_LEASES));

      // Register cluster steps
      registerNavStep(new PrimaryTabNavStep(
            NID_CLUSTER_PRIMARY_TAB_GETTING_STARTED,
            IDX_PRIMARY_TAB_GETTING_STARTED));

      registerNavStep(new PrimaryTabNavStep(NID_CLUSTER_PRIMARY_TAB_SUMMARY,
            IDX_PRIMARY_TAB_SUMMARY));

      registerNavStep(new PrimaryTabNavStep(NID_CLUSTER_PRIMARY_TAB_MONITOR,
            CommonUtil.getLocalizedString("common.tabs.name.monitor")));

      registerNavStep(new PrimaryTabNavStep(NID_CLUSTER_PRIMARY_TAB_MANAGE,
            IDX_PRIMARY_TAB_MANAGE));

      registerNavStep(new PrimaryTabNavStep(
            NID_CLUSTER_PRIMARY_TAB_RELATED_OBJECTS,
            IDX_PRIMARY_TAB_RELATED_OBJECTS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MONITOR_II_TAB_ISSUES, IDX_CLUSTER_TAB_MONITOR_ISSUES));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MONITOR_II_TAB_PERFORMANCE,
            CommonUtil
                  .getLocalizedString("common.tabs.name.monitor.performance")));

      registerNavStep(new TocNavStep(NID_TOC_MONITOR_PERFORMANCE_ADVANCED,
            IDX_TOC_MONITOR_PERFORMANCE_ADVANCED));

      registerNavStep(new TocNavStep(NID_TOC_MONITOR_PERFORMANCE_OVERVIEW,
            IDX_TOC_MONITOR_PERFORMANCE_OVERVIEW));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MONITOR_II_TAB_PROFILE_COMPLIANCE,
            IDX_CLUSTER_TAB_MONITOR_PROFILE_COMPLIANCE));

      registerNavStep(new SecondaryTabNavStep(NID_CLUSTER_MONITOR_II_TAB_TASKS,
            IDX_CLUSTER_TAB_MONITOR_TASKS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MONITOR_II_TAB_EVENTS, IDX_CLUSTER_TAB_MONITOR_EVENTS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MONITOR_II_TAB_RESOURCE_ALLOCATION,
            IDX_CLUSTER_TAB_MONITOR_RESOURCE_ALLOCATION));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MONITOR_II_TAB_UTILIZATION,
            IDX_CLUSTER_TAB_MONITOR_UTILIZATION));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MONITOR_II_TAB_VSPHERE_DRS,
            CommonUtil
                  .getLocalizedString("cluster.monitor.tabs.name.vSphereDRS")));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MONITOR_II_TAB_VSPHERE_HA,
            CommonUtil
                  .getLocalizedString("cluster.monitor.tabs.name.vSphereHA")));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MONITOR_II_TAB_STORAGE_REPORTS,
            IDX_CLUSTER_TAB_MONITOR_STORAGE_REPORTS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MANAGE_II_TAB_SETTINGS, IDX_CLUSTER_TAB_MANAGE_SETTINGS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MANAGE_II_TAB_ALARM_DEFINITIONS,
            IDX_CLUSTER_TAB_MANAGE_ALARM_DEFINITIONS));

      registerNavStep(new SecondaryTabNavStep(NID_CLUSTER_MANAGE_II_TAB_TAGS,
            CommonUtil.getLocalizedString("tags.items.tags.tab")));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MANAGE_II_TAB_PERMISSIONS,
            IDX_CLUSTER_TAB_MANAGE_PERMISSIONS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_MANAGE_II_TAB_SCHEDULED_TASKS,
            IDX_CLUSTER_TAB_MANAGE_SCHEDULED_TASKS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_RELATED_OBJECTS_II_TAB_TOP_OBJECTS,
            IDX_CLUSTER_TAB_REL_OBJ_TOP_OBJECTS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_RELATED_OBJECTS_II_TAB_HOSTS,
            IDX_CLUSTER_TAB_REL_OBJ_HOSTS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_RELATED_OBJECTS_II_TAB_VMS, IDX_CLUSTER_TAB_REL_OBJ_VMS));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATASTORE_RELATED_OBJECTS_II_TAB_VMS,
            IDX_DATASTORE_TAB_REL_OBJ_VMS));

      registerNavStep(new SecondaryTabNavStep(
            NID_HOSTSANDCLUSTERS_RELATED_OBJECTS_II_TAB_VMS,
            IDX_VCENTER_TAB_REL_OBJ_VMS));

      registerNavStep(new SecondaryTabNavStep(
            NID_VCENTER_RELATED_OBJECTS_II_TAB_VAPPS,
            IDX_VCENTER_TAB_REL_OBJ_VAPPS));

      registerNavStep(new SecondaryTabNavStep(
            NID_VCENTER_RELATED_OBJECTS_II_TAB_VMS, IDX_VCENTER_TAB_REL_OBJ_VMS));

      registerNavStep(new SecondaryTabNavStep(NID_VCENTER_VMS_II_TAB_VMS,
            CommonUtil.getLocalizedString("vcenter.vms.tabs.name.vms")));

      registerNavStep(new SecondaryTabNavStep(NID_VCENTER_VMS_II_TAB_VAPPS,
            CommonUtil.getLocalizedString("vcenter.vms.tabs.name.vapps")));

      registerNavStep(new SecondaryTabNavStep(
            NID_VCENTER_MANAGE_II_TAB_SETTINGS,
            CommonUtil.getLocalizedString("vcenter.manage.tabs.settings")));

      registerNavStep(new SecondaryTabNavStep(NID_VCENTER_MONITOR_II_TAB_VCHA,
            CommonUtil.getLocalizedString("vcenter.monitor.tabs.vcha")));

      registerNavStep(new SecondaryTabNavStep(NID_VCENTER_MONITOR_II_TAB_ISSUES,
            CommonUtil.getLocalizedString("vcenter.monitor.tabs.issues")));

      registerNavStep(new SecondaryTabNavStep(NID_VCENTER_MONITOR_II_TAB_EVENTS,
            CommonUtil.getLocalizedString("vcenter.monitor.tabs.events")));

      registerNavStep(new TocNavStep(
            NID_VCENTER_MONITOR_TOC_TRIGGERED_ALARMS,
            CommonUtil
                  .getLocalizedString("vcenter.monitor.vcha.toc.triggeredAlarms")));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATACENTER_RELATED_OBJECTS_II_TAB_VMS,
            IDX_DATACENTER_TAB_REL_OBJ_VMS));

      registerNavStep(new SecondaryTabNavStep(NID_DATACENTER_VMS_II_TAB_VMS,
            CommonUtil.getLocalizedString("datacenter.vms.tabs.name.vms")));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATACENTER_VMS_II_TAB_VM_FOLDERS,
            CommonUtil.getLocalizedString("datacenter.vms.tabs.name.vmfolders")));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_RELATED_OBJECTS_II_TAB_VAPPS,
            IDX_CLUSTER_TAB_REL_OBJ_VAPPS));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATACENTER_HOSTS_II_TAB_HOSTS,
            CommonUtil.getLocalizedString("datacenter.l2tab.hosts.hosts")));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_RELATED_OBJECTS_II_TAB_DATASTORES,
            IDX_CLUSTER_TAB_REL_OBJ_DATASTORES));

      registerNavStep(new TocNavStep(
            NID_DATASTORE_TOC_CONNECTIVITY_WITH_HOSTS,
            CommonUtil
                  .getLocalizedString(NID_DATASTORE_TOC_CONNECTIVITY_WITH_HOSTS)));

      registerNavStep(new SecondaryTabNavStep(NID_DATASTORE_VMS_II_TAB_VMS,
            CommonUtil.getLocalizedString("datastore.vms.tabs.name.vms")));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_RELATED_OBJECTS_II_TAB_DATASTORE_CLUSTERS,
            IDX_CLUSTER_TAB_REL_OBJ_DATASTORE_CLUSTERS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_RELATED_OBJECTS_II_TAB_NETWORKS,
            IDX_CLUSTER_TAB_REL_OBJ_NETWORKS));

      registerNavStep(new SecondaryTabNavStep(
            NID_CLUSTER_RELATED_OBJECTS_II_TAB_DVSES,
            IDX_CLUSTER_TAB_REL_OBJ_DVSES));

      registerNavStep(new SecondaryTabNavStep(NID_CLUSTER_HOSTS_II_HOSTS,
            CommonUtil.getLocalizedString("cluster.hosts.tabs.name.hosts")));

      registerNavStep(new SecondaryTabNavStep(NID_CLUSTER_VMS_II_VAPPS,
            CommonUtil.getLocalizedString("cluster.vms.tabs.name.vapps")));

      registerNavStep(new SecondaryTabNavStep(NID_CLUSTER_VMS_II_VMS,
            CommonUtil.getLocalizedString("cluster.vms.tabs.name.vms")));

      registerNavStep(new TocNavStep(
            NID_CLUSTER_TOC_MANAGE_SETTINGS_VSAN_GENERAL,
            IDX_CLUSTER_TOC_MANAGE_SETTINGS_VSAN_GENERAL));

      registerNavStep(new TocNavStep(
            NID_CLUSTER_TOC_MANAGE_SETTINGS_VSAN_DISK_MANAGEMENT,
            IDX_CLUSTER_TOC_MANAGE_SETTINGS_VSAN_DISK_MANAGEMENT));

      registerNavStep(new TocNavStep(NID_CLUSTER_TOC_MANAGE_SETTINGS_GENERAL,
            IDX_CLUSTER_TOC_MANAGE_SETTINGS_GENERAL));

      registerNavStep(new TocNavStep(
            NID_CLUSTER_TOC_MANAGE_SETTINGS_CONFIGURATION_LICENSING,
            IDX_CLUSTER_TOC_MANAGE_SETTINGS_CONFIGURATION_LICENSING));

      registerNavStep(new TocNavStep(
            NID_CLUSTER_TOC_MANAGE_SETTINGS_VSPHERE_DRS,
            IDX_CLUSTER_TOC_MANAGE_SETTINGS_VSPHERE_DRS));

      // Register host steps
      registerNavStep(new PrimaryTabNavStep(NID_HOST_PRIMARY_TAB_MANAGE,
            IDX_PRIMARY_TAB_MANAGE));

      registerNavStep(new PrimaryTabNavStep(
            NID_HOST_PRIMARY_TAB_RELATED_OBJECTS,
            IDX_PRIMARY_TAB_RELATED_OBJECTS));

      registerNavStep(new SecondaryTabNavStep(NID_HOST_MONITOR_TAB_PERFORMANCE,
            CommonUtil
                  .getLocalizedString("common.tabs.name.monitor.performance")));

      registerNavStep(new SecondaryTabNavStep(NID_HOST_MANAGE_TAB_NETWORKING,
            CommonUtil.getLocalizedString("host.tabs.networking")));

      registerNavStep(new SecondaryTabNavStep(NID_HOST_MANAGE_TAB_SETTINGS,
            IDX_HOST_TAB_MANAGE_SETTINGS));

      registerNavStep(new SecondaryTabNavStep(
            NID_ENTITY_MANAGE_TAB_II_TAB_CUSTOM_ATTRIBUTES,
            CommonUtil.getLocalizedString("ca.tab.name")));

      registerNavStep(new SecondaryTabNavStep(
            NID_HOST_RELATED_OBJECTS_II_TAB_VMS, IDX_HOST_TAB_REL_OBJ_VMS));

      registerNavStep(new SecondaryTabNavStep(
            NID_HOST_RELATED_OBJECTS_II_TAB_VAPPS, IDX_HOST_TAB_REL_OBJ_VAPPS));

      registerNavStep(new SecondaryTabNavStep(
            NID_HOST_RELATED_OBJECTS_II_TAB_DATASTORES,
            IDX_HOST_TAB_REL_OBJ_DATASTORES));

      registerNavStep(new SecondaryTabNavStep(
            NID_HOST_RELATED_OBJECTS_II_TAB_NETWORKS,
            IDX_HOST_TAB_REL_OBJ_NETWORKS));

      registerNavStep(new SecondaryTabNavStep(
            NID_HOST_NETWORKS_II_TAB_NETWORKS,
            CommonUtil.getLocalizedString("host.networks.tabs.name.networks")));

      registerNavStep(new SecondaryTabNavStep(NID_HOST_VMS_II_TAB_VMS,
            CommonUtil.getLocalizedString("host.vms.tabs.name.vms")));

      registerNavStep(new SecondaryTabNavStep(NID_HOST_VMS_II_TAB_VAPPS,
            CommonUtil.getLocalizedString("host.vms.tabs.name.vapps")));

      registerNavStep(new TocNavStep(
            NID_HOST_TOC_MANAGE_NETWORKING_VIRTUAL_SWITCHES,
            CommonUtil.getLocalizedString("configure.toc.networking.name.virturalSwitches")));

      registerNavStep(new TocNavStep(
            NID_HOST_TOC_MANAGE_NETWORKING_VMKERNEL_ADAPTERS,
            CommonUtil.getLocalizedString("configure.toc.networking.name.vmkernelAdapters")));

      registerNavStep(new TocNavStep(NID_TOC_CONFIGURE_OTHER_TAGS,
            CommonUtil.getLocalizedString("configure.toc.other.name.tags")));

      registerNavStep(new TocNavStep(NID_HOST_TOC_MANAGE_SETTINGS_LICENSING,
            IDX_HOST_TOC_MANAGE_SETTINGS_LICENSING));

      // Host Advanced system settings step
      registerNavStep(new TocNavStep(
            NID_HOST_TOC_MANAGE_SETTINGS_ADVANCED_SYSTEM_SETTINGS,
            IDX_HOST_TOC_MANAGE_SETTINGS_ADVANCED_SYSTEM_SETTINGS));

      // Host Graphics Settings step
      registerNavStep(new TocNavStep(NID_HOST_TOC_MANAGE_SETTINGS_GRAPHICS,
            IDX_HOST_TOC_MANAGE_SETTINGS_GRAPHICS));

      // Register dvs steps
      registerNavStep(new PrimaryTabNavStep(NID_DVS_PRIMARY_TAB_MANAGE,
            IDX_PRIMARY_TAB_MANAGE));

      registerNavStep(new PrimaryTabNavStep(
            NID_DVS_PRIMARY_TAB_RELATED_OBJECTS,
            IDX_PRIMARY_TAB_RELATED_OBJECTS));

      registerNavStep(new SecondaryTabNavStep(NID_DVS_MANAGE_TAB_SETTINGS,
            IDX_DVS_TAB_MANAGE_SETTINGS));
      registerNavStep(new SecondaryTabNavStep(NID_DVS_MANAGE_TAB_PORTS,
            IDX_DVS_TAB_MANAGE_PORTS));
      registerNavStep(new SecondaryTabNavStep(
            NID_DVS_MANAGE_TAB_RESOURCE_ALLOCATION,
            IDX_DVS_TAB_MANAGE_RESOURCE_ALLOCATION));
      registerNavStep(new SecondaryTabNavStep(
            NID_DVS_RELATED_OBJECTS_II_TAB_DVPGS, IDX_DVS_TAB_REL_OBJ_DVPGS));

      registerNavStep(new SecondaryTabNavStep(NID_DVS_NETWORKS_II_TAB_DVPGS,
            CommonUtil
                  .getLocalizedString("dvs.networks.tabs.name.dvPortgroups")));

      registerNavStep(new TocNavStep(NID_DVS_MANAGE_TAB_SETTINGS_TOC_TOPOLOGY,
            IDX_DVS_TAB_MANAGE_SETTINGS_TOC_TOPOLOGY));

      registerNavStep(new SecondaryTabNavStep(
            NID_HOST_RELATED_OBJECTS_II_TAB_NETWORKS,
            IDX_HOST_TAB_REL_OBJ_NETWORKS));

      // Register vApp steps
      registerNavStep(new PrimaryTabNavStep(
            NID_VAPP_PRIMARY_TAB_RELATED_OBJECTS,
            IDX_PRIMARY_TAB_RELATED_OBJECTS));

      registerNavStep(new PrimaryTabNavStep(NID_VAPP_PRIMARY_TAB_SUMMARY,
            IDX_PRIMARY_TAB_SUMMARY));

      registerNavStep(new PrimaryTabNavStep(NID_VAPP_PRIMARY_TAB_NETWORKS,
            CommonUtil.getLocalizedString("entity.tab.primary.networks")));

      registerNavStep(new SecondaryTabNavStep(
            NID_VAPP_RELATED_OBJECTS_II_TAB_VMS, IDX_VAPP_TAB_REL_OBJ_VMS));

      registerNavStep(new SecondaryTabNavStep(
            NID_VAPP_RELATED_OBJECTS_II_TAB_VAPPS, IDX_VAPP_TAB_REL_OBJ_VAPPS));

      registerNavStep(new SecondaryTabNavStep(
            NID_VAPP_RELATED_OBJECTS_II_TAB_NWS, IDX_VAPP_TAB_REL_OBJ_NWS));

      registerNavStep(new SecondaryTabNavStep(NID_VAPP_VMS_II_TAB_VMS,
            CommonUtil.getLocalizedString("vapp.vms.tabs.name.vms")));
      registerNavStep(new SecondaryTabNavStep(NID_VAPP_VMS_II_TAB_VAPPS,
            CommonUtil.getLocalizedString("vapp.vms.tabs.name.vapps")));

      // Register resource pool steps
      registerNavStep(new PrimaryTabNavStep(
            NID_RESPOOL_PRIMARY_TAB_RELATED_OBJECTS,
            IDX_PRIMARY_TAB_RELATED_OBJECTS));

      registerNavStep(new PrimaryTabNavStep(NID_RESPOOL_PRIMARY_TAB_SUMMARY,
            IDX_PRIMARY_TAB_SUMMARY));

      registerNavStep(new SecondaryTabNavStep(
            NID_RESPOOL_RELATED_OBJECTS_II_TAB_VMS, IDX_RESPOOL_TAB_REL_OBJ_VMS));

      registerNavStep(new SecondaryTabNavStep(
            NID_RESPOOL_RELATED_OBJECTS_II_TAB_VAPPS,
            IDX_RESPOOL_TAB_REL_OBJ_VAPPS));

      registerNavStep(new SecondaryTabNavStep(NID_RESPOOL_VMS_II_TAB_VMS,
            CommonUtil.getLocalizedString("respool.vms.tabs.name.vms")));
      registerNavStep(new SecondaryTabNavStep(NID_RESPOOL_VMS_II_TAB_VAPPS,
            CommonUtil.getLocalizedString("respool.vms.tabs.name.vapps")));

      // Libraries
      registerNavStep(new PrimaryTabNavStep(NID_CL_TAB_SUMMARY,
            IDX_CL_TAB_SUMMARY));

      registerNavStep(new PrimaryTabNavStep(NID_CL_TAB_MANAGE,
            IDX_CL_TAB_MANAGE));

      registerNavStep(new PrimaryTabNavStep(NID_CL_TAB_RELATED_OBJECTS,
            IDX_CL_TAB_RELATED_OBJECTS));

      registerNavStep(new PrimaryTabNavStep(NID_CL_TAB_TEMPLATES,
            CommonUtil.getLocalizedString("cl.tab.primary.templates")));

      registerNavStep(new PrimaryTabNavStep(NID_CL_TAB_OTHER_TYPES,
            CommonUtil.getLocalizedString("cl.tab.primary.otherTypes")));

      registerNavStep(new SecondaryTabNavStep(NID_CL_RELATED_OBJECTS_TEMPLATES,
            IDX_CL_RELATED_OBJECTS_TEMPLATES));

      registerNavStep(new SecondaryTabNavStep(
            NID_CL_RELATED_OBJECTS_OTHER_TYPES,
            IDX_CL_RELATED_OBJECTS_OTHER_TYPES));

      registerNavStep(new SecondaryTabNavStep(
            NID_CL_RELATED_OBJECTS_DATASTORES,
            IDX_CL_RELATED_OBJECTS_DATASTORES));

      // -------------------------------------------------------------------------
      // Datacenter navigation steps
      registerNavStep(new PrimaryTabNavStep(
            NID_DATACENTER_PRIMARY_TAB_GETTING_STARTED,
            IDX_PRIMARY_TAB_GETTING_STARTED));

      registerNavStep(new PrimaryTabNavStep(NID_DATACENTER_PRIMARY_TAB_SUMMARY,
            IDX_PRIMARY_TAB_SUMMARY));

      registerNavStep(new PrimaryTabNavStep(NID_DATACENTER_PRIMARY_TAB_MONITOR,
            CommonUtil.getLocalizedString("common.tabs.name.monitor")));

      registerNavStep(new PrimaryTabNavStep(NID_DATACENTER_PRIMARY_TAB_MANAGE,
            IDX_PRIMARY_TAB_MANAGE));

      registerNavStep(new PrimaryTabNavStep(
            NID_DATACENTER_PRIMARY_TAB_RELATED_OBJECTS,
            IDX_PRIMARY_TAB_RELATED_OBJECTS));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATACENTER_RELATED_OBJECTS_II_TAB_CLUSTERS,
            IDX_DATACENTER_TAB_REL_OBJ_CLUSTERS));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATACENTER_RELATED_OBJECTS_II_TAB_HOSTS,
            IDX_DATACENTER_TAB_REL_OBJ_HOSTS));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATACENTER_RELATED_OBJECTS_II_TAB_DVSES,
            IDX_DATACENTER_TAB_REL_OBJ_DVSES));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATACENTER_RELATED_OBJECTS_II_TAB_DVPGS,
            IDX_DATACENTER_TAB_REL_OBJ_DVPGS));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATACENTER_RELATED_OBJECTS_II_TAB_VAPPS,
            IDX_DATACENTER_TAB_REL_OBJ_VAPPS));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATACENTER_NETWORKS_II_TAB_DVPGS,
            CommonUtil
                  .getLocalizedString("datacenter.networks.tabs.name.dvPortgroups")));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATACENTER_NETWORKS_II_TAB_DVSES,
            CommonUtil.getLocalizedString("datacenter.networks.tabs.name.dvs")));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATACENTER_RELATED_OBJECTS_II_TAB_TOP_LEVEL_OBJECTS,
            IDX_DATACENTER_TAB_REL_OBJ_TOP_LEVEL_OBJECTS));

      registerNavStep(new SecondaryTabNavStep(NID_DATACENTER_MONITOR_TASKS,
            IDX_DATACENTER_TAB_MONITOR_TASKS));

      registerNavStep(new SecondaryTabNavStep(NID_DATACENTER_MONITOR_EVENTS,
            IDX_DATACENTER_TAB_MONITOR_EVENTS));

      registerNavStep(new SecondaryTabNavStep(
            NID_DATACENTER_RELATED_OBJECTS_II_TAB_DS_CLUSTERS,
            IDX_DATACENTER_TAB_REL_OBJ_DS_CLUSTERS));

      registerNavStep(new PrimaryTabNavStep(NID_LIBRARY_ITEM_MANAGE,
            IDX_LIBRARY_ITEM_MANAGE));

      registerNavStep(new SecondaryTabNavStep(NID_LIBRARY_ITEM_MANAGE_POLICIES,
            IDX_LIBRARY_ITEM_MANAGE_POLICIES));

      // -------------------------------------------------------------------------
      // Datastore cluster navigation steps
      registerNavStep(new PrimaryTabNavStep(
            NID_DS_CLUSTER_PRIMARY_TAB_GETTING_STARTED,
            IDX_PRIMARY_TAB_GETTING_STARTED));

      registerNavStep(new PrimaryTabNavStep(NID_DS_CLUSTER_PRIMARY_TAB_SUMMARY,
            IDX_PRIMARY_TAB_SUMMARY));

      registerNavStep(new PrimaryTabNavStep(NID_DS_CLUSTER_PRIMARY_TAB_MONITOR,
            CommonUtil.getLocalizedString("common.tabs.name.monitor")));

      registerNavStep(new PrimaryTabNavStep(NID_DS_CLUSTER_PRIMARY_TAB_MANAGE,
            IDX_PRIMARY_TAB_MANAGE));

      registerNavStep(new PrimaryTabNavStep(
            NID_DS_CLUSTER_PRIMARY_TAB_RELATED_OBJECTS,
            IDX_PRIMARY_TAB_RELATED_OBJECTS));

      registerNavStep(new SecondaryTabNavStep(NID_DS_CLUSTER_II_TAB_SDRS,
            IDX_DS_CLUSTER_TAB_MONITOR_SDRS));

      registerNavStep(new SecondaryTabNavStep(
            NID_DS_CLUSTER_MANAGE_II_TAB_SETTINGS,
            CommonUtil
                  .getLocalizedString("dscluster.manage.tabs.name.settings")));

      registerNavStep(new SecondaryTabNavStep(
            NID_DS_CLUSTER_RELATED_OBJECTS_II_TAB_DATASTORES,
            IDX_DS_CLUSTER_TAB_RELATED_OBJECTS_DATASTORES));

      // -------------------------------------------------------------------------
      // Datastore navigation steps
      registerNavStep(new PrimaryTabNavStep(NID_DATASTORE_PRIMARY_TAB_SUMMARY,
            IDX_PRIMARY_TAB_SUMMARY));

      registerNavStep(new PrimaryTabNavStep(NID_DATASTORE_PRIMARY_TAB_MANAGE,
            IDX_PRIMARY_TAB_MANAGE));

      registerNavStep(new PrimaryTabNavStep(NID_DATASTORE_PRIMARY_TAB_VMS,
            CommonUtil.getLocalizedString("datastore.tabs.name.vms")));

      registerNavStep(new SecondaryTabNavStep(NID_DATASTORE_MANAGE_II_TAB_TAGS,
            CommonUtil.getLocalizedString("datastore.manage.tabs.name.Tags")));

      // -------------------------------------------------------------------------
      // Library Manage navigation steps
      registerNavStep(new SecondaryTabNavStep(
            NID_LIBRARY_MANAGE_II_TAB_SETTINGS, IDX_LIBRARY_MANAGE_SETTINGS));

      registerNavStep(new SecondaryTabNavStep(NID_LIBRARY_MANAGE_II_TAB_TAGS,
            IDX_LIBRARY_MANAGE_TAGS));

      registerNavStep(new SecondaryTabNavStep(
            NID_LIBRARY_MANAGE_II_TAB_PERMISSIONS,
            IDX_LIBRARY_MANAGE_PERMISSIONS));

      registerNavStep(new SecondaryTabNavStep(
            NID_LIBRARY_MANAGE_II_TAB_STORAGE, IDX_LIBRARY_MANAGE_STORAGE));

      // -------------------------------------------------------------------------
      // Storage policy Monitor navigation steps
      registerNavStep(new SecondaryTabNavStep(
            NID_STORAGE_POLICY_MONITOR_STORAGE_COMPATIBILITY,
            IDX_STORAGE_POLICY_MONITOR_STORAGE_COMPATIBILITY));

      registerNavStep(new SecondaryTabNavStep(
            NID_STORAGE_POLICY_MONITOR_VMS_AND_DISKS,
            CommonUtil
                  .getLocalizedString(NID_STORAGE_POLICY_MONITOR_VMS_AND_DISKS)));

      // Host profiles
      registerNavStep(new SecondaryTabNavStep(
            NID_HOST_PROFILES_II_TAB_COMPLIANCE,
            CommonUtil
                  .getLocalizedString("hostProfiles.monitor.tabs.name.Compliance")));
      registerNavStep(new SecondaryTabNavStep(
            NID_HOST_PROFILES_II_TAB_SETTINGS,
            CommonUtil
                  .getLocalizedString("hostProfiles.manage.tabs.name.Settings")));

      // -------------------------------------------------------------------------
      // VC Manage Settings navigation steps
      registerNavStep(new TocNavStep(NID_VC_MANAGE_SETTINGS_TOC_VCHA,
            CommonUtil.getLocalizedString("vcenter.manage.settings.toc.vcha")));

      // VC Monitor VCHA navigation steps
      registerNavStep(new TocNavStep(NID_VC_MONITOR_VCHA_TOC_HEALTH,
            CommonUtil.getLocalizedString("vcenter.monitor.vcha.toc.health")));

      // Storage Policy Components navigation step
      // Storage Policy components navigation step
      registerNavStep(new PrimaryTabNavStep(
            NID_STORAGE_POLICY_STORAGE_POLICY_COMPONENTS,
            CommonUtil
                  .getLocalizedString(NID_STORAGE_POLICY_STORAGE_POLICY_COMPONENTS)));
   }
}
