/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * A spec containing the parameters used in "advanced search"
 */
public class AdvancedSearchSpec extends BaseSpec {
   public static final String ENTITY_TYPE_VM = "Virtual Machine";
   public static final String ENTITY_TYPE_CLUSTER = "Cluster";
   public static final String ENTITY_TYPE_VDC = "Virtual Datacenter";
   public static final String ENTITY_TYPE_POLICY = "Placement Policy";
   public static final String ENTITY_TYPE_POLICY_TAG = "Policy Tag";
   public static final String ENTITY_TYPE_HOST = "Host";

   public static final String PROPERTY_NAME_VM_POLICY =
         "Virtual Machine Placement Policy";
   public static final String PROPERTY_NAME_VM_POLICY_TAG = "Virtual Machine Policy Tag";
   public static final String PROPERTY_NAME_VM_VDC_NAME =
         "Virtual Machine Virtual Datacenter Name";
   public static final String PROPERTY_NAME_VM_OVERALL_COMPLIANCE =
         "Virtual Datacenter Virtual Machine Overall Compliance";

   public static final String PROPERTY_NAME_CLUSTER_VDC_NAME =
         "Cluster Virtual Datacenter Name";

   public static final String PROPERTY_NAME_HOST_VM_NAME =
         "Host Virtual Machine Name";

   public static final String PROPERTY_NAME_VDC_VM_OVERALL_COMPLIANCE =
         "Virtual Datacenter Virtual Machine Overall Compliance";
   public static final String PROPERTY_NAME_VDC_CLUSTER_NAME =
         "Virtual Datacenter Cluster Name";
   public static final String PROPERTY_NAME_VDC_NAME = "Virtual Datacenter Name";
   public static final String PROPERTY_NAME_VDC_PLACEMENT_POLICY_NAME =
         "Virtual Datacenter VM Placement Policy Name";
   public static final String PROPERTY_NAME_VDC_VM_NAME =
         "Virtual Datacenter Virtual Machine Name";
   public static final String PROPERTY_NAME_VDC_VM_POWER_STATE =
         "Virtual Datacenter Virtual Machine Power State";
   public static final String PROPERTY_NAME_VDC_VM_POLICY =
         "Virtual Datacenter Virtual Machine Placement Policy";
   public static final String PROPERTY_NAME_VDC_VM_POLICY_TAG =
         "Virtual Datacenter Virtual Machine Policy Tag";
   public static final String PROPERTY_NAME_VDC_VM_POLICY_COMPLIANCE =
         "Virtual Datacenter Virtual Machine VM Storage Policies Compliance";
   public static final String PROPERTY_NAME_VDC_VM_TEMPLATE =
         "Virtual Datacenter Virtual Machine Template";

   public static final String PROPERTY_NAME_POLICY_NAME = "Placement Policy Name";
   public static final String PROPERTY_NAME_POLICY_VDC_NAME =
         "Placement Policy Virtual Datacenter Name";

   public static final String PROPERTY_NAME_POLICY_TAG_NAME = "Policy Tag Name";
   public static final String PROPERTY_NAME_POLICY_TAG_VDC_NAME =
         "Policy Tag Virtual Datacenter Name";

   public static final String OPERATOR_CONTAINS = "contains";
   public static final String OPERATOR_IS = "is";
   public static final String OPERATOR_IS_NOT = "is not";

   public static final String COMPLIANCE_NON_COMPLIANT = "Noncompliant";
   public static final String COMPLIANCE_UNKNOWN = "Unknown";
   public static final String COMPLIANCE_COMPLIANT = "Compliant";
   public static final String COMPLIANCE_UPDATING = "Updating";

   public static final String POWER_STATE_ON = "Powered On";
   public static final String POWER_STATE_OFF = "Powered Off";

   /**
    * The type of the entity to search for
    */
   public DataProperty<String> entityType;

   /**
    * The property name to compare
    */
   public DataProperty<String> propertyName;

   /**
    * The operator used for the comparison
    */
   public DataProperty<String> operator;

   /**
    * The value of the property to search for
    */
   public DataProperty<String> propertyValue;

   /**
    * The value of the compliance to search for
    */
   public DataProperty<String> compliance;

   /**
    * The search results that are expected when executing the search
    */
   public DataProperty<ManagedEntitySpec> searchResults;

   /**
    * The search results that should not be found when executing the search
    */
   public DataProperty<ManagedEntitySpec> negativeResults;

}
